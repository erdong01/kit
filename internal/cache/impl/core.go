package impl

import (
	"github.com/erDong01/micro-kit/db/redis"
	"github.com/erDong01/micro-kit/internal/cache/I"
	"github.com/erDong01/micro-kit/internal/cache/config"
)

func New() *I.ICache {
	drive := config.GetDrive()
	if drive == "redis" {
		return &I.ICache{
			redis.NewCache(),
		}
	} else {
		panic("驱动不存在")
	}
}
