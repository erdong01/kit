package mysql

import (
	"github.com/erDong01/micro-kit/core"
	"github.com/jinzhu/gorm"
)

func Db() *gorm.DB {
	return core.New().Db
}
