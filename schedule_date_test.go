package kit

import (
	"context"
	"testing"
	"time"
)

type notifyJob struct {
	ch chan time.Time
}

func (j *notifyJob) OnTimer() {
	j.ch <- time.Now()
}

func nextSecondTime() time.Time {
	now := time.Now()
	return now.Truncate(time.Second).Add(2 * time.Second)
}

func TestScheduleAddDate(t *testing.T) {
	s := NewSchedule()
	s.SetIntervalTime(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run(ctx)

	fireAt := time.Now().Add(80 * time.Millisecond)
	job := &notifyJob{ch: make(chan time.Time, 1)}
	id := s.AddDate(job, fireAt)

	if id == 0 {
		t.Fatal("expected non-zero task id")
	}

	select {
	case firedAt := <-job.ch:
		if firedAt.Before(fireAt.Add(-50 * time.Millisecond)) {
			t.Fatalf("timer fired too early: fireAt=%v firedAt=%v", fireAt, firedAt)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for AddDate task")
	}
}

func TestScheduleAddDatePastTime(t *testing.T) {
	s := NewSchedule()
	s.SetIntervalTime(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run(ctx)

	job := &notifyJob{ch: make(chan time.Time, 1)}
	s.AddDate(job, time.Now().Add(-time.Second))

	select {
	case <-job.ch:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected past AddDate task to run immediately")
	}
}

func TestScheduleAddDailyOnce(t *testing.T) {
	s := NewSchedule()
	s.SetIntervalTime(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run(ctx)

	fireAt := nextSecondTime()
	job := &notifyJob{ch: make(chan time.Time, 1)}
	id := s.AddDaily(job, fireAt.Hour(), fireAt.Minute(), fireAt.Second(), false)
	if id == 0 {
		t.Fatal("expected non-zero task id")
	}

	select {
	case firedAt := <-job.ch:
		if firedAt.Before(fireAt.Add(-time.Second)) {
			t.Fatalf("daily task fired too early: fireAt=%v firedAt=%v", fireAt, firedAt)
		}
	case <-time.After(3500 * time.Millisecond):
		t.Fatal("timed out waiting for daily task")
	}

	time.Sleep(50 * time.Millisecond)
	if remaining := s.Surplus(id); remaining != 0 {
		t.Fatalf("expected one-time daily task removed after execution, got %v", remaining)
	}
}

func TestScheduleAddDailyPersistence(t *testing.T) {
	s := NewSchedule()
	s.SetIntervalTime(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run(ctx)

	fireAt := nextSecondTime()
	job := &notifyJob{ch: make(chan time.Time, 1)}
	id := s.AddDaily(job, fireAt.Hour(), fireAt.Minute(), fireAt.Second(), true)
	if id == 0 {
		t.Fatal("expected non-zero task id")
	}

	select {
	case firedAt := <-job.ch:
		if firedAt.Before(fireAt.Add(-time.Second)) {
			t.Fatalf("daily task fired too early: fireAt=%v firedAt=%v", fireAt, firedAt)
		}
	case <-time.After(3500 * time.Millisecond):
		t.Fatal("timed out waiting for daily task")
	}

	remaining := s.Surplus(id)
	if remaining < 23*time.Hour {
		t.Fatalf("expected next daily task about one day later, got %v", remaining)
	}
}

func TestScheduleAddWeeklyOnce(t *testing.T) {
	s := NewSchedule()
	s.SetIntervalTime(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run(ctx)

	fireAt := nextSecondTime()
	job := &notifyJob{ch: make(chan time.Time, 1)}
	id := s.AddWeekly(job, fireAt.Weekday(), fireAt.Hour(), fireAt.Minute(), fireAt.Second(), false)
	if id == 0 {
		t.Fatal("expected non-zero task id")
	}

	select {
	case firedAt := <-job.ch:
		if firedAt.Before(fireAt.Add(-time.Second)) {
			t.Fatalf("weekly task fired too early: fireAt=%v firedAt=%v", fireAt, firedAt)
		}
	case <-time.After(3500 * time.Millisecond):
		t.Fatal("timed out waiting for weekly task")
	}

	time.Sleep(50 * time.Millisecond)
	if remaining := s.Surplus(id); remaining != 0 {
		t.Fatalf("expected one-time weekly task removed after execution, got %v", remaining)
	}
}

func TestScheduleAddWeeklyPersistence(t *testing.T) {
	s := NewSchedule()
	s.SetIntervalTime(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run(ctx)

	fireAt := nextSecondTime()
	job := &notifyJob{ch: make(chan time.Time, 1)}
	id := s.AddWeekly(job, fireAt.Weekday(), fireAt.Hour(), fireAt.Minute(), fireAt.Second(), true)
	if id == 0 {
		t.Fatal("expected non-zero task id")
	}

	select {
	case firedAt := <-job.ch:
		if firedAt.Before(fireAt.Add(-time.Second)) {
			t.Fatalf("weekly task fired too early: fireAt=%v firedAt=%v", fireAt, firedAt)
		}
	case <-time.After(3500 * time.Millisecond):
		t.Fatal("timed out waiting for weekly task")
	}

	remaining := s.Surplus(id)
	if remaining < 6*24*time.Hour {
		t.Fatalf("expected next weekly task about one week later, got %v", remaining)
	}
}

func TestScheduleAddMonthlyOnce(t *testing.T) {
	s := NewSchedule()
	s.SetIntervalTime(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run(ctx)

	fireAt := nextSecondTime()
	job := &notifyJob{ch: make(chan time.Time, 1)}
	id := s.AddMonthly(job, fireAt.Day(), fireAt.Hour(), fireAt.Minute(), fireAt.Second(), false)
	if id == 0 {
		t.Fatal("expected non-zero task id")
	}

	select {
	case firedAt := <-job.ch:
		if firedAt.Before(fireAt.Add(-time.Second)) {
			t.Fatalf("monthly task fired too early: fireAt=%v firedAt=%v", fireAt, firedAt)
		}
	case <-time.After(3500 * time.Millisecond):
		t.Fatal("timed out waiting for monthly task")
	}

	time.Sleep(50 * time.Millisecond)
	if remaining := s.Surplus(id); remaining != 0 {
		t.Fatalf("expected one-time monthly task removed after execution, got %v", remaining)
	}
}

func TestScheduleAddMonthlyPersistence(t *testing.T) {
	s := NewSchedule()
	s.SetIntervalTime(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run(ctx)

	fireAt := nextSecondTime()
	job := &notifyJob{ch: make(chan time.Time, 1)}
	id := s.AddMonthly(job, fireAt.Day(), fireAt.Hour(), fireAt.Minute(), fireAt.Second(), true)
	if id == 0 {
		t.Fatal("expected non-zero task id")
	}

	select {
	case firedAt := <-job.ch:
		if firedAt.Before(fireAt.Add(-time.Second)) {
			t.Fatalf("monthly task fired too early: fireAt=%v firedAt=%v", fireAt, firedAt)
		}
	case <-time.After(3500 * time.Millisecond):
		t.Fatal("timed out waiting for monthly task")
	}

	remaining := s.Surplus(id)
	if remaining < 27*24*time.Hour {
		t.Fatalf("expected next monthly task about one month later, got %v", remaining)
	}
}
