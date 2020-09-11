package drive

import (
	"fmt"
	config2 "github.com/erDong01/micro-kit/config"
	"github.com/go-redis/redis/v7"
)

func err_handler(err error) {
	fmt.Printf("err_handler, error:%s\n", err.Error())
}
func New() *redis.Client {
	var Redis *redis.Client
	redisConfig := config2.GetRedisConfig()
	config2.New().Get(&redisConfig, "redis")
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		fmt.Printf("ping error[%s]\n", err.Error())
		err_handler(err)
	}
	return Redis
}
