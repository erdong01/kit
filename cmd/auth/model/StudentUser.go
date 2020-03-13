package model

import "time"

type StudentUser struct {
	BaseModel
	StudentUserId     int64      `gorm:"primary_key" json:"student_user_id,omitempty"`
	StudentUserNo     int64      `gorm:"unique;not null" json:"student_user_no,omitempty"`
	StudentUserName   *string     `gorm:"unique;not null" json:"student_user_name,omitempty"`
	Password          string     `json:"password,omitempty"`
	StudentUserStatus int8       `json:"student_user_status,omitempty"`
	StudentUserMobile int64      `json:"student_user_mobile,omitempty"`
	StudentUserHead   string      `json:"student_user_head,omitempty"`
	LastLoginIp       int64      `json:"last_login_ip,omitempty"`
	LastLoginTime     *time.Time `json:"last_login_time,omitempty"`
	StudentUserCampusOne  StudentUserCampus  `json:"student_user_campus,omitempty" gorm:"foreignkey:StudentUserNo;association_foreignkey:StudentUserNo"`
}

func (StudentUser) TableName() string {
	return "rxt_student_user"
}
