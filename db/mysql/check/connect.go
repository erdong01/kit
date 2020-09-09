package check

import (
	"github.com/erDong01/micro-kit/core"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	return core.New().Db
}
