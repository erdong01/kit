package timer

import (
	"github.com/erDong01/micro-kit/wrong"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TimeNearShift  = 8
	TimeNear       = (1 << TimeNearShift)
	TimeLevelShift = 6
	TimeLevel      = (1 << TimeLevelShift)
	TimeNearMask   = (TimeNear - 1)
	TimeLevelMask  = (TimeLevel - 1)
	TickInterval   = 10 * time.Millisecond
)

//先搞清楚下面的单位
//1秒=1000毫秒 milliseconds
//1毫秒=1000微秒 microseconds
//1微秒=1000纳秒 nanoseconds
//整个timer中毫秒的精度都是10ms，
//也就是说毫秒的一个三个位，但是最小的位被丢弃
type (
	Handle func()
	Node   struct {
		next   *Node
		expire uint32
		handle Handle
		id     *int64
		time   uint32
		bOnce  bool
	}

	//这个队列可以换成无锁队列
	LinkList struct {
		head Node
		tail *Node
	}
	Timer struct {
		near         [TimeNear]LinkList     //临近的定时器数组
		t            [4][TimeLevel]LinkList //四个级别的定时器数组
		lock         sync.Mutex             //锁
		time         uint32                 //计数器
		startTime    uint32                 //程序启动的时间点，timestamp，秒数
		current      uint64                 //从程序启动到现在的耗时，精度10毫秒级
		currentPoint uint64                 //当前时间，精度10毫秒级
		pTimer       *time.Ticker           //定时器
		loopNode     []*Node
	}
	Op struct {
		bOnce bool
	}
	OpOption func(*Op)
)

var (
	TIMER *Timer
	gId   int64
)

func (this *Node) LoadId() int64 {
	return atomic.LoadInt64(this.id)
}

func (op *Op) applyOpts(opts []OpOption) {
	for _, opt := range opts {
		opt(op)
	}
}

func withOnce() OpOption {
	return func(op *Op) {
		op.bOnce = true
	}
}

func init() {
	TIMER = &Timer{}
	TIMER.Init()
}

func uuid() int64 {
	return atomic.AddInt64(&gId, 1)
}

//清空链表，返回链表第一个结点
func linkClear(list *LinkList) *Node {
	ret := list.head.next
	list.head.next = nil
	list.tail = &list.head
	return ret
}

func link(list *LinkList, node *Node) {
	list.tail.next = node
	list.tail = node
	node.next = nil
}

func (this *Timer) Init() {
	for i := 0; i < TimeNear; i++ {
		linkClear(&this.near[i])
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < TimeLevel; j++ {
			linkClear(&this.t[i][j])
		}
	}

	this.current = 0
	this.pTimer = time.NewTicker(TickInterval)
	this.currentPoint = uint64(time.Now().UnixNano()) / uint64(TickInterval)
	go this.run()
}

func (this *Timer) addNode(node *Node) {
	time := node.expire      //去看一下它是在哪赋值的
	currentTime := this.time //当前计数
	//没有超时，或者说时间点特别近了
	if (time | TimeNearMask) == (currentTime | TimeNearMask) {
		link(&this.near[time&TimeNearMask], node)
	} else { //这里有一种特殊情况，就是当time溢出，回绕的时候
		i := 0
		mask := uint32(TimeNear << TimeLevelShift)
		for i = 0; i < 3; i++ { //看到i<3没，很重要很重要
			if (time | (mask - 1)) == (currentTime | (mask - 1)) {
				break
			}
			mask <<= TimeLevelShift // mask越来越大
		}
		link(&this.t[i][(time>>uint(TimeNearShift+i*TimeLevelShift))&TimeLevelMask], node)
	}
}

//Add 添加一个定时器
func (this *Timer) Add(id *int64, time uint32, handle Handle, opts ...OpOption) *Node {
	op := Op{}
	op.applyOpts(opts)
	node := &Node{expire: time + this.time, handle: handle, time: time, bOnce: op.bOnce, id: id} //超时时间+当前计数
	this.lock.Lock()
	defer func() { this.lock.Unlock() }()
	this.addNode(node)
	return node
}

//Delete 删除一个定时器
func (this *Timer) Delete(id *int64) {
	atomic.StoreInt64(id, 0)
}

//moveList 移动某个级别的链表内容
func (this *Timer) moveList(level int, idx int) {
	current := linkClear(&this.t[level][idx])
	for current != nil {
		temp := current.next
		this.addNode(current)
		current = temp
	}
}

//这是一个非常重要的函数
//定时器的移动都在这里
func (this *Timer) shift() {
	mask := uint32(TimeNear)
	this.time += 1
	ct := this.time
	if ct == 0 { //time溢出了
		this.moveList(3, 0) //这里就是那个很重要的3
	} else { //time正常
		time := ct >> TimeNearShift
		i := 0
		for (ct & (mask - 1)) == 0 {
			idx := time & TimeLevelMask
			if idx != 0 {
				this.moveList(i, int(idx))
				break
			}
			mask <<= TimeLevelShift //mask越来越大
			time >>= TimeLevelShift //time越来越小
			i += 1
		}
	}

}

//派发消息到目标服务消息队列
func (this *Timer) dispatch(current *Node) {
	for current != nil {
		id := current.LoadId()
		if id != 0 {
			current.handle()
			if !current.bOnce {
				this.loopNode = append(this.loopNode, current)
			}
		}
		current = current.next
	}
}

//派发消息
func (this *Timer) execute() {
	idx := this.time & TimeNearMask
	for this.near[idx].head.next != nil {
		current := linkClear(&this.near[idx])
		this.lock.Unlock()
		this.dispatch(current)
		this.lock.Lock()
		for _, v := range this.loopNode {
			v.expire = v.time + this.time
			this.addNode(v)
		}
	}
	this.loopNode = []*Node{}
}

//时间更新好了以后，这里检查调用各个定时器
func (this *Timer) advace() {
	this.lock.Lock()
	this.execute()
	this.shift()
	this.execute()
	this.lock.Unlock()

}

//在线程中不断被调用
//调用时间 间隔为微秒
func (this *Timer) update() {
	cp := uint64(time.Now().UnixNano()) / uint64(TickInterval)
	if cp < this.currentPoint {
		this.currentPoint = cp
	} else if cp != this.currentPoint {
		diff := cp - this.currentPoint
		this.currentPoint = cp
		this.current += diff
		for i := uint64(0); i < diff; i++ {
			this.advace() //注意这里
		}
	}
}

func (this *Timer) loop() bool {
	defer func() {
		if err := recover(); err != nil {
			wrong.TraceCode(err)
		}
	}()
	select {
	case <-this.pTimer.C:
		this.update()
	}
	return true
}

func (this *Timer) run() {
	for {
		if !this.loop() {
			break
		}
	}
	this.pTimer.Stop()
}

func RegisterTimer(id *int64, duration time.Duration, handle Handle, opts ...OpOption) {
	TIMER.Add(id, uint32(duration/TickInterval), handle, opts...)
}

func StopTimer(id *int64) {
	if id != nil {
		TIMER.Delete(id)
	}
}
