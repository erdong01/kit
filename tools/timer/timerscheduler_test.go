package timer

/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/9 10:14
* @Mail: danbing.at@gmail.com
*
*  时间轮定时器调度器单元测试
 */

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// 触发函数
func foo(args ...interface{}) {
	fmt.Printf("I am No. %d function, delay %d ms\n", args[0].(int), args[1].(int))
}

// 手动创建调度器运转时间轮
func TestNewTimerScheduler(t *testing.T) {
	timerScheduler := NewTimerScheduler()
	timerScheduler.Start()

	//在scheduler中添加timer
	for i := 1; i < 2000; i++ {
		f := NewDelayFunc(foo, []interface{}{i, i * 3})
		tID, err := timerScheduler.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			log.Println("create timer error", tID, err)
			break
		}
	}

	//执行调度器触发函数
	go func() {
		delayFuncChan := timerScheduler.GetTriggerChan()
		for df := range delayFuncChan {
			df.Call()
		}
	}()

	//阻塞等待
	select {}
}

// 采用自动调度器运转时间轮
func TestNewAutoExecTimerScheduler(t *testing.T) {
	autoTS := NewAutoExecTimerScheduler()

	//给调度器添加Timer
	for i := 0; i < 2000; i++ {
		f := NewDelayFunc(foo, []interface{}{i, i * 3})
		tID, err := autoTS.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			log.Println("create timer error", tID, err)
			break
		}
	}

	//阻塞等待
	select {}
}

// 测试取消一个定时器
func TestCancelTimerScheduler(t *testing.T) {
	Scheduler := NewAutoExecTimerScheduler()
	f1 := NewDelayFunc(foo, []interface{}{3, 3})
	f2 := NewDelayFunc(foo, []interface{}{5, 5})
	timerID1, err := Scheduler.CreateTimerAfter(f1, time.Duration(3)*time.Second)
	if nil != err {
		t.Log("Scheduler.CreateTimerAfter(f1, time.Duration(3)*time.Second)", "err：", err)
	}
	timerID2, err := Scheduler.CreateTimerAfter(f2, time.Duration(5)*time.Second)
	if nil != err {
		t.Log("Scheduler.CreateTimerAfter(f1, time.Duration(3)*time.Second)", "err：", err)
	}
	log.Printf("timerID1=%d ,timerID2=%d\n", timerID1, timerID2)
	Scheduler.CancelTimer(timerID1) //删除timerID1

	//阻塞等待
	select {}
}

var timerScheduler *TimerScheduler
var TID uint32
var TID2 uint32

func TAddTimer(v ...interface{}) {
	i := v[0].(int)
	fmt.Println(i)
	i++
	f1 := NewDelayFunc(TAddTimer, []interface{}{i, 1})
	timerScheduler.AddTimer(TID, f1, time.Duration(10)*time.Millisecond)
}

func TAddTimer2(v ...interface{}) {
	i := v[0].(int)
	fmt.Println(i)
	i += 10
	f1 := NewDelayFunc(TAddTimer2, []interface{}{i, 1})
	timerScheduler.AddTimer(TID2, f1, time.Duration(10)*time.Microsecond)

}

// 采用自动调度器运转时间轮
func TestAddTimer(t *testing.T) {
	timerScheduler = NewAutoExecTimerScheduler()
	f1 := NewDelayFunc(TAddTimer, []interface{}{1})
	TID, _ = timerScheduler.CreateTimerAfter(f1, time.Duration(2)*time.Second)
	f2 := NewDelayFunc(TAddTimer2, []interface{}{10})
	TID2, _ = timerScheduler.CreateTimerAfter(f2, time.Duration(2)*time.Second)
	time.Sleep(time.Second * 100)
	//阻塞等待
}
