package redis

import (
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Redis *redis.Client `name:"db0"`
}
