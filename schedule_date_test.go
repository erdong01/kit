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
