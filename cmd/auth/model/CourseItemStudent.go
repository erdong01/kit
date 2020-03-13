package model

type CourseItemStudent struct {
	BaseModel
	CourseItemStudentId     int64 `gorm:"primary_key" json:"course_item_student_id,omitempty"`
	CourseItemId            int64 `json:"course_item_id,omitempty"`
	StudentUserCampusId     int64 `json:"student_user_campus_id,omitempty"`
	CourseItemStudentStatus int8  `json:"course_item_student_status,omitempty"`
}

func (CourseItemStudent) TableName() string {
	return "rxt_course_item_student"
}
