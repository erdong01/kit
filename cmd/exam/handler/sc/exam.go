package sc

import (
	"context"
	exam "rxt/cmd/exam/proto/sc"
	"rxt/cmd/exam/service/examService"
	"rxt/internal/core"
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
