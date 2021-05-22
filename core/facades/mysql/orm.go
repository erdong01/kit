package mysql

import (
	"github.com/erDong01/micro-kit/core"
	"gorm.io/gorm"
)

func Db() *gorm.DB {
	return core.New().Db
}
