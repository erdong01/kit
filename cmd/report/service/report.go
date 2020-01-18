package service

import (
	"rxt/cmd/report/service/report"
)

func Show(examNo int64) (report.ShowReport, error) {
	return report.New().Show(examNo)
}
