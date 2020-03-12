package examClient

import (
	"fmt"
	exam "rxt/cmd/exam/proto/sc"
	"testing"
)

func TestExamRpc_ScExam(t *testing.T) {
	var r *exam.ExamRequest
	test, _ := New().Submit(r)
	fmt.Println(test.ExamNo)
}
