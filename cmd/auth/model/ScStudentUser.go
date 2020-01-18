package model

import "time"

type ScStudentUser struct {
	BaseModel
	StudentUserId     int64      `gorm:"primary_key" json:"student_user_id,omitempty"`
	StudentUserNo     int64      `json:"student_user_no,omitempty"`
	StudentUserName   string     `json:"student_user_name,omitempty"`
	Password          string     `json:"password,omitempty"`
	StudentUserStatus int8       `json:"student_user_status,omitempty"`
	StudentUserMobile string     `json:"student_user_mobile,omitempty"`
	StudentUserHead   string     `json:"student_user_head,omitempty"`
	LastLoginIp       int64      `json:"last_login_ip,omitempty"`
	LastLoginTime     *time.Time `json:"last_login_time,omitempty"`
}

func (ScStudentUser) TableName() string {
	return "rxt_sc_student_user"
}
