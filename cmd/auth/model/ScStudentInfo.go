package model

import "time"

type ScStudentInfo struct {
	BaseModel
	StudentInfoId   int64      `gorm:"primary_key" json:"student_info_id,omitempty"`
	StudentUserNo   int64      `json:"student_user_no,omitempty"`
	StudentIdCard   string     `json:"student_id_card,omitempty"`
	StudentRealName string     `json:"student_real_name,omitempty"`
	StudentSex      int8       `json:"student_sex,omitempty"`
	SchoolId        int64      `json:"school_id,omitempty"`
	GradeId         int64      `json:"grade_id,omitempty"`
	CompanyNo       int64      `json:"company_no,omitempty"`
	CampusNo        int64      `json:"campus_no,omitempty"`
	DistrictId      int64      `json:"district_id,omitempty"`
	CityId          int64      `json:"city_id,omitempty"`
	ProvinceId      int64      `json:"province_id,omitempty"`
	StudentBirthday *time.Time `json:"student_birthday,omitempty"`
}

func (ScStudentInfo) TableName() string {
	return "rxt_sc_student_info"
}
