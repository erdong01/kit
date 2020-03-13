package util

import (
	"time"
)

// WeekStartDay 获取本周第一天（周一）
func WeekStartDay() time.Time {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}
