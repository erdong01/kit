package model

import "time"

type StudentInfo struct {
	BaseModel
	StudentInfoId   int64      `gorm:"primary_key" json:"student_info_id,omitempty"`
	StudentUserNo   int64      `json:"student_user_no,omitempty"`
	StudentIdCard   string     `json:"student_id_card,omitempty"`
	StudentRealName string     `json:"student_real_name,omitempty"`
	StudentBirthday *time.Time `json:"student_birthday,omitempty"`
	StudentSex      int8       `json:"student_sex,omitempty"`
	CompanyNo       int64      `json:"company_no,omitempty"`
	CampusNo        int64      `json:"campus_no,omitempty"`
}

func (StudentInfo) TableName() string {
	return "rxt_student_user"
}
