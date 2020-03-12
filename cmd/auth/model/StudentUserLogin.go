package model

type StudentUserLogin struct {
	BaseModel
	StudentUserLoginId   int64  `gorm:"primary_key" json:"student_user_login_id,omitempty"`
	StudentUserLoginName int64  `json:"student_user_login_name,omitempty"`
	Password             string `json:"password,omitempty"`
	StudentUserNo        string `json:"student_user_no,omitempty"`
}

func (StudentUserLogin) TableName() string {
	return "rxt_student_user_login"
}
