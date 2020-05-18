package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/erDong01/micro-kit/internal/core"
)

func Db() *gorm.DB {
	return core.New().Db
}
