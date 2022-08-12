package drive

import (
	"context"
	"fmt"
	"github.com/erDong01/micro-kit/config"
	"github.com/go-redis/redis/v8"
	"runtime"
	"time"
)

func New() *redis.Client {
	var Redis *redis.Client
	redisConfig := config.GetRedisConfig()
	config.New().Get(&redisConfig, "redis")
	Redis = redis.NewClient(&redis.Options{
		Addr:         redisConfig.Addr,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		PoolSize:     runtime.NumCPU() * 512,
		MinIdleConns: runtime.NumCPU() * 10,
		IdleTimeout:  time.Minute * 5,
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
	})
	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("ping error[%s]\n", err.Error())
	}
	return Redis
}
