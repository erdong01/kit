package model

type ScExamStudent struct {
	BaseModel
	ExamStudentId int64 `gorm:"primary_key" json:"exam_student_id,omitempty"`
	StudentUserNo int64 `json:"student_user_no,omitempty"`
	ExamNo        int64 `json:"exam_no,omitempty"`
	BookNo        int64 `json:"book_no,omitempty"`
}

func (ScExamStudent) TableName() string {
	return "rxt_sc_exam_student"
}
