package examClient

import (
	"fmt"
	exam "github.com/erDong01/gin-kit/cmd/exam/proto/sc"
	"testing"
)

func TestExamRpc_ScExam(t *testing.T) {
	var r *exam.ExamRequest
	test, _ := New().Submit(r)
	fmt.Println(test.ExamNo)
}
