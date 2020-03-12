package sc

import (
	"context"
	exam "rxt/cmd/exam/proto/sc"
	"rxt/cmd/exam/service/examService"
	"rxt/internal/core"
)

type Server struct {
}

func (s Server) Submit(ctx context.Context, request *exam.ExamRequest) (*exam.ExamResponse, error) {
	var res exam.ExamResponse
	core := core.GlobalTransaction()
	examService := examService.New(core)
	param, err := examService.Submit(request)
	if err != nil {
		core.Transaction.Rollback()
		return &res, err
	}
	reportErr := examService.ReportCreate(param)
	if reportErr != nil {
		core.Transaction.Rollback()
		return &res, err
	}
	CreateScExamStudentErr := examService.CreateScExamStudent(param.ExamNo, request.ScStudentUserNo, request.BookNo)
	if CreateScExamStudentErr != nil {
		core.Transaction.Rollback()
		return &res, err
	}
	core.Transaction.Commit()
	res.ExamNo = param.ExamNo
	return &res, nil
}
