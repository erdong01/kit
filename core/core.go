package core

import (
	"sync"

	"github.com/erDong01/micro-kit/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Info struct {
	Port    int
	Env     string
	Name    string
	version string
}

type Core struct {
	Db     *gorm.DB
	Config *config.Config
	Info   *Info
	opts   []Option
	Mongo  *mongo.Client
	Etcd   *clientv3.Client
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
	var core = new(Core)
	*core = *New()
	return core
}

func Set(newCore *Core) {
	c = newCore
}
