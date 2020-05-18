package core

import (
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/erDong01/gin-kit/internal/cache/I"
	"github.com/erDong01/gin-kit/internal/config"
	"sync"
)

type info struct {
	port    int
	env     string
	Name    string
	version string
}
type Core struct {
	Db          *gorm.DB
	Transaction *gorm.DB
	Redis       *redis.Client
	Config      *config.Config
	Cache       *I.ICache
	Info        info
	opts        []Option
	once        sync.Once
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

func Close() {
	if c.Redis != nil {
		defer c.Redis.Close()
	}
	if c.Db != nil {
		defer c.Db.Close()
	}
}

