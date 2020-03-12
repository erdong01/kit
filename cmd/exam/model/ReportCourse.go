package model

type ReportCourse struct {
	BaseModel
	ReportCourseId   int64 `gorm:"primary_key" json:"report_course_id,omitempty"`
	ExamNo           int64 `json:"exam_no,omitempty"`
	ReportCourseSort int   `json:"report_course_sort,omitempty"`
}

func (ReportCourse) TableName() string {
	return "rxt_report_course"
}
