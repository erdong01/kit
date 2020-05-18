package check

import (
	"github.com/go-redis/redis/v7"
	"github.com/erDong01/gin-kit/internal/core"
)

func Connect() *redis.Client {
	return core.New().Redis
}
