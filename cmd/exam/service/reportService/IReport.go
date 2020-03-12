package reportService

type Report struct {
	IReport
}

func New() Report {
	report := &V1{}
	report.Init()
	return Report{report}
}

type IReport interface {
}
