package core

import (
	"github.com/gin-gonic/gin"
	rds "github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/erDong01/micro-kit/internal/config"
	drive2 "github.com/erDong01/micro-kit/internal/db/mysql/drive"
	"github.com/erDong01/micro-kit/internal/db/redis/drive"
	"github.com/erDong01/micro-kit/internal/http"
)

// ICore 服务核心接口
type ICore interface {
	Init()
	GetEnv() string
	GetName() string
	GetVersion() string
	GetDb() *gorm.DB
	GetRedis() *rds.Client
	GetPort() int
}

// Init 启动其他服务
func (c *Core) Init() {
	c.once.Do(func() {
		for _, o := range c.opts {
			o(c)
		}
	})
}

// GetEnv 获取当前环境
func (c *Core) GetEnv() string {
	return c.Info.env
}

// GetName 获取当前服务名称
func (c *Core) GetName() string {
	return c.Info.Name
}

// GetVersion 获取当前版本号
func (c *Core) GetVersion() string {
	return c.Info.version
}

// GetDb 获取当前数据库实例
func (c *Core) GetDb() *gorm.DB {
	return c.Db
}

// GetDb 获取当前数据库实例
func (c *Core) GetPort() int {
	return c.Info.port
}

// GetRedis 获取当前Redis实例
func (c *Core) GetRedis() *rds.Client {
	return c.Redis
}

// Option 返回匿名函数 供初始化执行时
type Option func(*Core)

// Make 构建新核心
func Make(opts ...Option) ICore {
	once.Do(func() {
		c = &Core{
			opts: opts,
		}
	})
	return c
}

// Env 设置环境变量
func Env(env string) Option {
	return func(c *Core) {
		c.Info.env = env
	}
}

// Name 设置项目名称
func Name(name string) Option {
	return func(c *Core) {
		c.Info.Name = name
	}
}

// Port 设置端口号
func Port(port int) Option {
	return func(c *Core) {
		c.Info.port = port
	}
}

// Version 设置项目版本
func Version(version string) Option {
	return func(c *Core) {
		c.Info.version = version
	}
}

// Engine 设置项目版本
func Engine(route func(g *gin.Engine) *gin.Engine) Option {
	return func(c *Core) {
		http.Init(route(GetEngine()), c.Info.port)
	}
}

// DbRegister 设置数据库
func DbRegister() Option {
	return func(c *Core) {
		c.Db = drive2.New()
	}
}

// RedisRegister 设置Redis
func RedisRegister() Option {
	return func(c *Core) {
		c.Redis = drive.New()
	}
}

func ConfigRegister() Option {
	return func(c *Core) {
		config.Init("config")
	}
}

// GetEngine 获取当前gin引擎
func GetEngine() *gin.Engine {
	return gin.Default()
}

// GetEngine 获取当前gin引擎
func SetPort(port int) {
	c.Info.port = port
}
