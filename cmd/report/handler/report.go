package handler

import (
	"context"
	"log"
	"rxt/cmd/report/proto/report"
	"rxt/cmd/report/service"
)

type Server struct{}

func (s *Server) Show(ctx context.Context, in *report.ReportRequest) (*report.ReportResponse, error) {
	log.Println("request: ")
	exam, err := service.Show(in.ExamId)
	p := &report.ReportResponse{
		ExamId: exam.Exam.ExamId,
	}
	return p, err
}
