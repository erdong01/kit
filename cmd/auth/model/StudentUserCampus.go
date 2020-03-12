package model

import "time"

type StudentUserCampus struct {
	BaseModel
	StudentUserCampusId     int64      `gorm:"primary_key" json:"student_user_campus_id,omitempty"`
	StudentUserNo           int64      `json:"student_user_no,omitempty"`
	CompanyNo               int64      `json:"company_no,omitempty"`
	CampusNo                int64      `json:"campus_no,omitempty"`
	StudentUserCampusStatus int64      `json:"student_user_campus_status,omitempty"`
	GradeId                 int64      `json:"grade_id,omitempty"`
	StudentSchoolId         int64      `json:"student_school_id,omitempty"`
	StudentSchool           string     `json:"student_school,omitempty"`
	FirstSignDate           *time.Time `json:"first_sign_date,omitempty"`
	StudentManagerId        int64      `json:"student_manager_id,omitempty"`
	CounselorId             int64      `json:"counselor_id,omitempty"`
	ProvinceId              int64      `json:"province_id,omitempty"`
	CityId                  int64      `json:"city_id,omitempty"`
	DistrictId              int64      `json:"district_id,omitempty"`
	Address                 string     `json:"address,omitempty"`
	StudentIdCard           string     `json:"student_id_card,omitempty"`
	StudentRealName         string     `json:"student_real_name,omitempty"`
	StudentBirthday         *time.Time `json:"student_birthday,omitempty"`
	StudentSex              int8       `json:"student_sex,omitempty"`
	StudentEconomicStatus   int8       `json:"student_economic_status,omitempty"`
	FrontUserCampusId       int8       `json:"front_user_campus_id,omitempty"`
	StudentSigningStatus    int8       `json:"student_signing_status,omitempty"`
	TotalClassTime          float64    `json:"total_class_time,omitempty"`
	UseClassTime            float64    `json:"use_class_time,omitempty"`
}

func (StudentUserCampus) TableName() string {
	return "rxt_student_user_campus"
}
