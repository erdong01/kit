package model

type ReportClassQuestion struct {
	BaseModel
	ReportClassQuestionId int64   `gorm:"primary_key" json:"report_class_question_id,omitempty"`
	ExamClassId           int64   `json:"exam_class_id,omitempty"`
	QuestionNo            int64   `json:"question_no,omitempty"`
	QuestionScore         float64 `json:"question_score,omitempty"`
	QuestionAvgScore      float64 `json:"question_avg_score,omitempty"`
	CorrectStudentCount   int32   `json:"correct_student_count,omitempty"`
	ErrorStudentCount     int32   `json:"error_student_count,omitempty"`
	DownloadCount         int32   `json:"download_count,omitempty"`
}

func (ReportClassQuestion) TableName() string {
	return "rxt_report_class_question"
}
