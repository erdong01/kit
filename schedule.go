package kit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type (
	event struct {
		timer        *time.Timer
		ticker       *time.Ticker
		endTime      time.Time
		duration     time.Duration
		recurring    bool
		nextTime     func(time.Time) time.Time
		delayHandler DelayHandler
	}
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
		intervalTime time.Duration
	}
)

func NewSchedule() *Schedule {
	return &Schedule{
		events:       make(map[uint64]*event),
		intervalTime: time.Second,
	}
}

func (s *Schedule) Run(ctx context.Context) {
	go s.Start(ctx)
}

func (s *Schedule) SetIntervalTime(v time.Duration) {
	s.intervalTime = v
}

// Start 开始
func (s *Schedule) Start(ctx context.Context) {
	secTicker := time.NewTicker(s.intervalTime)
	defer secTicker.Stop()
	for {
		select {
		case <-secTicker.C:
			s.mutex.Lock()
			for k, v := range s.events {
				if v.expire() {
					handler := v.delayHandler
					go func(d DelayHandler) {
						defer func() {
							if err := recover(); err != nil {
								fmt.Println(err)
							}
						}()
						d.OnTimer()
					}(handler)

					if v.timer != nil && !v.recurring {
						v.timer.Stop()
						delete(s.events, k)
					}
					if v.ticker != nil {
						v.endTime = v.endTime.Add(v.duration)
					}
					if v.timer != nil && v.recurring && v.nextTime != nil {
						next := v.nextTime(v.endTime)
						delay := time.Until(next)
						if delay < 0 {
							delay = 0
						}
						v.endTime = next
						v.timer.Reset(delay)
					}
				}
			}
			s.mutex.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (s *event) expire() bool {
	if s.timer != nil {
		select {
		case <-s.timer.C:
			return true
		default:
		}
	}

	if s.ticker != nil {
		select {
		case <-s.ticker.C:
			return true
		default:
		}
	}
	return false
}

// Add 添加
func (s *Schedule) Add(delayHandler DelayHandler, duration time.Duration, persistence bool) (TID uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.IDGen++
	TID = s.IDGen
	var ev event
	if !persistence {
		ev.timer = time.NewTimer(duration)
	} else {
		ev.ticker = time.NewTicker(duration)
	}
	ev.duration = duration
	ev.delayHandler = delayHandler
	ev.endTime = time.Now().Add(duration)
	s.events[s.IDGen] = &ev
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	duration := time.Until(first)
	if duration < 0 {
		duration = 0
	}

	s.IDGen++
	TID = s.IDGen
	ev := &event{
		timer:        time.NewTimer(duration),
		endTime:      first,
		delayHandler: delayHandler,
	}
	if persistence {
		ev.recurring = true
		ev.nextTime = next
	}
	s.events[TID] = ev
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
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
}

// Remove 移除
func (s *Schedule) Remove(id uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if ev, ok := s.events[id]; ok {
		if ev.timer != nil && !ev.timer.Stop() {
			<-ev.timer.C
		}
		if ev.ticker != nil {
			ev.ticker.Stop()
		}
		delete(s.events, id)
	}
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
