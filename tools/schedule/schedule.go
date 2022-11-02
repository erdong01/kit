package schedule

import (
	"context"
	"sync"
	"time"

	"github.com/erDong01/micro-kit/base"
)

type (
	event struct {
		timer     *time.Timer
		ticker    *time.Ticker
		endTime   time.Time
		duration  time.Duration
		scheduler Scheduler
	}
	// Scheduler 调度程序
	Scheduler interface {
		OnTimer()
	}
	// Schedule 时间表
	// 命名也只是不想与标准库Timer重名而已
	// 定时时间精确到秒，不精确管理goroutine的退出
	// 主要逻辑只不过是对标准库Timer、Ticker封装管理而已
	Schedule struct {
		IDGen  int
		mutex  sync.Mutex
		events map[int]*event
	}
)

func New() *Schedule {
	return &Schedule{
		events: make(map[int]*event),
	}
}
func (s *Schedule) Run(ctx context.Context) {
	go s.Start(ctx)
}

// Start 开始
func (s *Schedule) Start(ctx context.Context) {
	var timer = time.NewTimer(1 * time.Second)
	for {
		timer.Reset(1 * time.Second)
		select {
		case <-timer.C:
			s.mutex.Lock()
			for k, v := range s.events {
				if v.expire() {
					go func(s Scheduler) {
						defer func() {
							if err := recover(); err != nil {
								base.TraceCode(err)
							}
						}()
						s.OnTimer()
					}(s.events[k].scheduler)
					if v.timer != nil {
						delete(s.events, k)
					}
					if v.ticker != nil {
						v.endTime = v.endTime.Add(v.duration)
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
func (s *Schedule) Add(scheduler Scheduler, duration time.Duration, persistence bool) (TID int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.IDGen++
	TID = s.IDGen
	ev := event{}
	if !persistence {
		ev.timer = time.NewTimer(duration)
	} else {
		ev.ticker = time.NewTicker(duration)
	}
	ev.duration = duration
	ev.scheduler = scheduler
	ev.endTime = time.Now().Add(duration)
	s.events[s.IDGen] = &ev
	return
}

// Remove 移除
func (s *Schedule) Remove(id int) {
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
func (s *Schedule) Surplus(id int) (duration time.Duration) {
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
