package drive

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"rxt/internal/config"
)

func err_handler(err error) {
	fmt.Printf("err_handler, error:%s\n", err.Error())
	panic(err.Error())
}
func New() *redis.Client {
	var Redis *redis.Client
	redisConfig := config.GetRedisConfig()
	config.New().Get(&redisConfig, "redis")
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	pong, err := Redis.Ping().Result()
	if err != nil {
		fmt.Printf("ping error[%s]\n", err.Error())
		err_handler(err)
	}
	fmt.Printf("ping result: %s\n", pong)
	return Redis
}
