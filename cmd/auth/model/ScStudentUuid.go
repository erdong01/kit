package model

type ScStudentUuid struct {
	BaseModel
	StudentUuidId int64  `gorm:"primary_key" json:"student_uuid_id,omitempty"`
	StudentUuid   string `json:"student_uuid,omitempty"`
	StudentUserNo int64  `json:"student_user_no,omitempty"`
}

func (ScStudentUuid) TableName() string {
	return "rxt_sc_student_uuid"
}
