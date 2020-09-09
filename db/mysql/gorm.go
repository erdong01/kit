package mysql

import (
	"github.com/erDong01/micro-kit/db/mysql/check"
	"github.com/erDong01/micro-kit/db/mysql/drive"
	"gorm.io/gorm"
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
