package impl

import (
	"rxt/internal/cache/I"
	"rxt/internal/cache/config"
	"rxt/internal/db/redis"
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
