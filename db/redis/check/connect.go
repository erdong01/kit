package check

import (
	"github.com/erDong01/micro-kit/core"
	"github.com/go-redis/redis/v7"
)

func Connect() *redis.Client {
	return core.New().Redis
}
