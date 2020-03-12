package model

type ReportClassQuestionAnswerOption struct {
	BaseModel
	ReportClassQuestionAnswerOptionId int64   `gorm:"primary_key" json:"report_class_question_answer_option_id,omitempty"`
	ReportClassQuestionId             int64   `json:"report_class_question_id,omitempty"`
	QuestionSmallId                   int64   `json:"question_small_id,omitempty"`
	StudentUserNo                     int64   `json:"student_user_no,omitempty"`
	StudentUserCampusId               float64 `json:"student_user_campus_id,omitempty"`
	QuestionSmallIsRight              int8    `json:"question_small_is_right,omitempty"`
	OptionName                        string  `json:"option_name,omitempty"`
}

func (ReportClassQuestionAnswerOption) TableName() string {
	return "rxt_report_class_question_answer_option"
}
