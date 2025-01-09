package testSchedule

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/erdong01/kit/schedule"
)

var s *schedule.Schedule

func TestXxx(t *testing.T) {
	s = schedule.New()
	var context = context.Background()
	s.Run(context)
	var job3 Job
	job3.TableId = 10000
	s.Add(&job3, time.Second*time.Duration(10), false)
	for i := 0; i < 10; i++ {
		var job Job
		job.TableId = i
		s.Add(&job, time.Second*time.Duration(i+1), false)

		var job2 Job
		job2.TableId = i + 10
		s.Add(&job2, time.Second*time.Duration(i), false)
	}
	time.Sleep(time.Second * 20)
}

type Job struct {
	TableId int
}

func (j *Job) OnTimer() {
	fmt.Println(j.TableId)
	panic(j.TableId)
}
