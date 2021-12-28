package core

import (
	"sync"

	"github.com/erDong01/micro-kit/config"
	"github.com/go-redis/redis/v8"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type info struct {
	Port    int
	env     string
	Name    string
	version string
}

type Core struct {
	Db          *gorm.DB
	Transaction *gorm.DB
	Redis       *redis.Client
	Config      *config.Config
	Info        info
	opts        []Option
	once        sync.Once
	Mongo       *mongo.Client
	Etcd        *clientv3.Client
}

var (
	c    *Core
	once sync.Once
)

func New() *Core {
	once.Do(func() {
		c = &Core{}
	})
	return c

}

func Copy() *Core {
	var core *Core = new(Core)
	*core = *New()
	return core
}

func Set(newCore *Core) {
	once.Do(func() {
		c = newCore
	})
}

func Close() {
	if c.Redis != nil {
		defer c.Redis.Close()
	}

}
