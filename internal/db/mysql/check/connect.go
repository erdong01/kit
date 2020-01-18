package check

import (
	"github.com/jinzhu/gorm"
	"rxt/internal/core"
)

func Connect() *gorm.DB {
	return core.New().Db
}
