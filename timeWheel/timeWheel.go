package timeWheel

import "time"

type TimeWheel struct {
	interval time.Duration
}

func NewTimeWheel(interval time.Duration, slotNum int) (*TimeWheel, error) {
	tw := &TimeWheel{
		interval: interval,
	}
	return tw, nil
}
