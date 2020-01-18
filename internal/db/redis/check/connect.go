package check

import (
	"github.com/go-redis/redis/v7"
	"rxt/internal/core"
)

func Connect() *redis.Client {
	return core.New().Redis
}
