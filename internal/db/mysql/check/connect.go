package check

import (
	"github.com/jinzhu/gorm"
	"github.com/erDong01/micro-kit/internal/core"
)

func Connect() *gorm.DB {
	return core.New().Db
}
