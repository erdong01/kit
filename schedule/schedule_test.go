package schedule

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var schedule *Schedule

func TestXxx(t *testing.T) {
	schedule = New()
	var context = context.Background()
	schedule.Run(context)
	var job3 Job
	job3.TableId = 10000
	schedule.Add(&job3, time.Second*time.Duration(10), false)
	for i := 0; i < 10; i++ {
		var job Job
		job.TableId = i
		schedule.Add(&job, time.Second*time.Duration(i+1), false)

		var job2 Job
		job2.TableId = i + 10
		schedule.Add(&job2, time.Second*time.Duration(i), false)
	}
	time.Sleep(time.Second * 20)
}

type Job struct {
	TableId int
}

func (j *Job) OnTimer() {
	fmt.Println(j.TableId)
}
