package report

import (
	"rxt/cmd/report/model"
)

type Register struct {
	IReport
}

func New() Register {
	reportService := &ServiceV1{}
	reportService.Init()
	return Register{reportService}
}

type IReport interface {
	Show(examNo int64) (ShowReport, error)
}

type ShowReport struct {
	Exam model.Exam `json:"exam"`
}
