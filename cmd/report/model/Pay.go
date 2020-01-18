package model

import "github.com/jinzhu/gorm"

type Pay struct {
	gorm.Model
	PayId int `gorm:"primary_key"`
	UserNo int
	UserType int
	PayNO int64
	PayMethod int
	Pay_status int
}

func (Pay) TableName() string {
	return "rxt_pay"
}