package mysql

import (
	"github.com/jinzhu/gorm"
	"rxt/internal/db/mysql/check"
	"rxt/internal/db/mysql/drive"
)

// New 初始化数据库ORM
func New() *gorm.DB {
	var db *gorm.DB
	db = check.Connect()
	if db != nil {
		return db
	}
	db = drive.New()
	return db
}



