package schedule

import (
	"context"
	"sync"
	"time"

	"github.com/erDong01/micro-kit/base"
)

type (
	event struct {
		timer        *time.Timer
		ticker       *time.Ticker
		endTime      time.Time
		duration     time.Duration
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
		IDGen  uint64
		mutex  sync.Mutex
		events map[uint64]*event
	}
)

func New() *Schedule {
	return &Schedule{
		events: make(map[uint64]*event),
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
					go func(d DelayHandler) {
						defer func() {
							if err := recover(); err != nil {
								base.TraceCode(err)
							}
						}()
						d.OnTimer()
					}(s.events[k].delayHandler)

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
