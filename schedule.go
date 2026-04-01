package kit

import (
	"container/heap"
	"context"
	"fmt"
	"sync"
	"time"
)

type (
	event struct {
		id           uint64
		endTime      time.Time
		duration     time.Duration
		recurring    bool
		nextTime     func(time.Time) time.Time
		delayHandler DelayHandler
		index        int
	}
	eventHeap []*event
	// DelayHandler 调度程序 需要传递参数通过 struct 传递
	DelayHandler interface {
		OnTimer()
	}
	// Schedule 时间表
	// 命名也只是不想与标准库Timer重名而已
	// 定时时间精确到秒，不精确管理goroutine的退出
	// 主要逻辑只不过是对标准库Timer、Ticker封装管理而已
	Schedule struct {
		IDGen        uint64
		mutex        sync.Mutex
		events       map[uint64]*event
		heap         eventHeap
		intervalTime time.Duration
		wakeup       chan struct{}
	}
)

func NewSchedule() *Schedule {
	s := &Schedule{
		events:       make(map[uint64]*event),
		intervalTime: time.Second,
		wakeup:       make(chan struct{}, 1),
	}
	heap.Init(&s.heap)
	return s
}

func (s *Schedule) Run(ctx context.Context) {
	go s.Start(ctx)
}

// SetIntervalTime 仅为兼容保留，最小堆调度模型下不再依赖轮询间隔。
func (s *Schedule) SetIntervalTime(v time.Duration) {
	s.intervalTime = v
}

// Start 开始
func (s *Schedule) Start(ctx context.Context) {
	timer := time.NewTimer(time.Hour)
	defer timer.Stop()

	for {
		wait := s.nextWait()
		if wait < 0 {
			select {
			case <-s.wakeup:
			case <-ctx.Done():
				return
			}
			continue
		}

		resetTimer(timer, wait)
		select {
		case <-timer.C:
			s.fireDue()
		case <-s.wakeup:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		}
	}
}

// Add 添加
func (s *Schedule) Add(delayHandler DelayHandler, duration time.Duration, persistence bool) (TID uint64) {
	if duration < 0 {
		duration = 0
	}

	s.mutex.Lock()
	s.IDGen++
	TID = s.IDGen
	ev := &event{
		id:           TID,
		endTime:      time.Now().Add(duration),
		duration:     duration,
		recurring:    persistence,
		delayHandler: delayHandler,
		index:        -1,
	}
	s.pushEventLocked(ev)
	s.mutex.Unlock()
	s.notify()
	return
}

// AddDate 按指定日期时间添加一次性任务
func (s *Schedule) AddDate(delayHandler DelayHandler, date time.Time) (TID uint64) {
	duration := time.Until(date)
	if duration < 0 {
		duration = 0
	}
	return s.Add(delayHandler, duration, false)
}

// AddDaily 每天在指定时刻执行
func (s *Schedule) AddDaily(delayHandler DelayHandler, hour, minute, second int, persistence bool) (TID uint64) {
	now := time.Now()
	first := nextDailyTime(now, hour, minute, second)
	return s.addCalendarDate(delayHandler, first, persistence, func(last time.Time) time.Time {
		return nextDailyTime(last, hour, minute, second)
	})
}

// AddWeekly 每周在指定星期和时刻执行
func (s *Schedule) AddWeekly(delayHandler DelayHandler, weekday time.Weekday, hour, minute, second int, persistence bool) (TID uint64) {
	now := time.Now()
	first := nextWeeklyTime(now, weekday, hour, minute, second)
	return s.addCalendarDate(delayHandler, first, persistence, func(last time.Time) time.Time {
		return nextWeeklyTime(last, weekday, hour, minute, second)
	})
}

// AddMonthly 每月在指定日期和时刻执行
// 当目标月份不存在该日期时，自动使用该月最后一天。
func (s *Schedule) AddMonthly(delayHandler DelayHandler, day, hour, minute, second int, persistence bool) (TID uint64) {
	now := time.Now()
	first := nextMonthlyTime(now, day, hour, minute, second)
	return s.addCalendarDate(delayHandler, first, persistence, func(last time.Time) time.Time {
		return nextMonthlyTime(last, day, hour, minute, second)
	})
}

func (s *Schedule) addCalendarDate(delayHandler DelayHandler, first time.Time, persistence bool, next func(time.Time) time.Time) (TID uint64) {
	if first.Before(time.Now()) {
		first = time.Now()
	}

	s.mutex.Lock()
	s.IDGen++
	TID = s.IDGen
	ev := &event{
		id:           TID,
		endTime:      first,
		recurring:    persistence,
		nextTime:     next,
		delayHandler: delayHandler,
		index:        -1,
	}
	s.pushEventLocked(ev)
	s.mutex.Unlock()
	s.notify()
	return
}

func nextDailyTime(from time.Time, hour, minute, second int) time.Time {
	next := time.Date(from.Year(), from.Month(), from.Day(), hour, minute, second, 0, from.Location())
	if !next.After(from) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

func nextWeeklyTime(from time.Time, weekday time.Weekday, hour, minute, second int) time.Time {
	next := time.Date(from.Year(), from.Month(), from.Day(), hour, minute, second, 0, from.Location())
	diff := (int(weekday) - int(from.Weekday()) + 7) % 7
	next = next.AddDate(0, 0, diff)
	if !next.After(from) {
		next = next.AddDate(0, 0, 7)
	}
	return next
}

func nextMonthlyTime(from time.Time, day, hour, minute, second int) time.Time {
	year, month := from.Year(), from.Month()
	next := monthlyDate(year, month, day, hour, minute, second, from.Location())
	if !next.After(from) {
		year, month = nextMonth(year, month)
		next = monthlyDate(year, month, day, hour, minute, second, from.Location())
	}
	return next
}

func monthlyDate(year int, month time.Month, day, hour, minute, second int, loc *time.Location) time.Time {
	lastDay := daysInMonth(year, month)
	if day < 1 {
		day = 1
	}
	if day > lastDay {
		day = lastDay
	}
	return time.Date(year, month, day, hour, minute, second, 0, loc)
}

func nextMonth(year int, month time.Month) (int, time.Month) {
	if month == time.December {
		return year + 1, time.January
	}
	return year, month + 1
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// Remove 移除
func (s *Schedule) Remove(id uint64) {
	s.mutex.Lock()
	if ev, ok := s.events[id]; ok {
		if ev.index >= 0 {
			heap.Remove(&s.heap, ev.index)
		}
		delete(s.events, id)
	}
	s.mutex.Unlock()
	s.notify()
}

// Surplus 剩余
func (s *Schedule) Surplus(id uint64) (duration time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if ev, ok := s.events[id]; ok {
		duration = ev.endTime.Sub(time.Now())
	}
	if duration < 0 {
		duration = 0
	}
	return
}

func (s *Schedule) nextWait() time.Duration {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.heap) == 0 {
		return -1
	}
	wait := time.Until(s.heap[0].endTime)
	if wait < 0 {
		return 0
	}
	return wait
}

func (s *Schedule) fireDue() {
	var handlers []DelayHandler
	s.mutex.Lock()
	now := time.Now()
	for len(s.heap) > 0 {
		ev := s.heap[0]
		if ev.endTime.After(now) {
			break
		}

		heap.Pop(&s.heap)
		handlers = append(handlers, ev.delayHandler)

		if !ev.recurring {
			delete(s.events, ev.id)
			continue
		}

		s.scheduleNextLocked(ev, now)
		heap.Push(&s.heap, ev)
	}
	s.mutex.Unlock()

	for _, handler := range handlers {
		go safeCall(handler)
	}
}

func (s *Schedule) scheduleNextLocked(ev *event, now time.Time) {
	switch {
	case ev.nextTime != nil:
		next := ev.nextTime(ev.endTime)
		for !next.After(now) {
			next = ev.nextTime(next)
		}
		ev.endTime = next
	case ev.duration > 0:
		next := ev.endTime.Add(ev.duration)
		for !next.After(now) {
			next = next.Add(ev.duration)
		}
		ev.endTime = next
	default:
		ev.endTime = now
	}
}

func (s *Schedule) pushEventLocked(ev *event) {
	s.events[ev.id] = ev
	heap.Push(&s.heap, ev)
}

func (s *Schedule) notify() {
	select {
	case s.wakeup <- struct{}{}:
	default:
	}
}

func safeCall(d DelayHandler) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	d.OnTimer()
}

func resetTimer(timer *time.Timer, wait time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(wait)
}

func (h eventHeap) Len() int {
	return len(h)
}

func (h eventHeap) Less(i, j int) bool {
	return h[i].endTime.Before(h[j].endTime)
}

func (h eventHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *eventHeap) Push(x interface{}) {
	ev := x.(*event)
	ev.index = len(*h)
	*h = append(*h, ev)
}

func (h *eventHeap) Pop() interface{} {
	old := *h
	n := len(old)
	ev := old[n-1]
	ev.index = -1
	*h = old[:n-1]
	return ev
}
