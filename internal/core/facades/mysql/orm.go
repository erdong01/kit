package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/erDong01/gin-kit/internal/core"
)

func Db() *gorm.DB {
	return core.New().Db
}
