package check

import (
	"github.com/jinzhu/gorm"
	"github.com/erDong01/gin-kit/internal/core"
)

func Connect() *gorm.DB {
	return core.New().Db
}
