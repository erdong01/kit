package mysql

import (
	"github.com/jinzhu/gorm"
	"rxt/internal/core"
)

func Db() *gorm.DB {
	return core.New().Db
}
