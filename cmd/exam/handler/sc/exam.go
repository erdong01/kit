package sc

import (
	"context"
	exam "github.com/erDong01/micro-kit/cmd/exam/proto/sc"
	"github.com/erDong01/micro-kit/cmd/exam/service/examService"
	"github.com/erDong01/micro-kit/internal/core"
)

type Server struct {
}

// 提交
func (s Server) Submit(ctx context.Context, request *exam.ExamRequest) (*exam.ExamResponse, error) {
	var res exam.ExamResponse
	core := core.GlobalTransaction()
	examService := examService.New(core)
	result, err := examService.Submit(request)
	if err != nil {
		core.Transaction.Rollback()
		return &res, err
	}
	res.ExamNo = result
	return &res, nil
}
