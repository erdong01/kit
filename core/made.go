package core

import (
	"github.com/erDong01/micro-kit/config"
	gongDbDrive "github.com/erDong01/micro-kit/db/mongoDB/drive"
	mysqlDrive "github.com/erDong01/micro-kit/db/mysql/drive"
	redisDrive "github.com/erDong01/micro-kit/db/redis/drive"
	"github.com/erDong01/micro-kit/http"
	"github.com/gin-gonic/gin"
	rds "github.com/go-redis/redis/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
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
	ConfigRegister() *config.Config
	SetConfigFile(file string) *config.Config
	MongoRegister(Uri ...string) *mongo.Client
}

// Init 启动其他服务
func (*Core) Init() {
	New().once.Do(func() {
		for _, o := range c.opts {
			o(c)
		}
	})
}

// GetEnv 获取当前环境
func (*Core) GetEnv() string {
	return New().Info.env
}

// GetName 获取当前服务名称
func (*Core) GetName() string {
	return New().Info.Name
}

// GetVersion 获取当前版本号
func (*Core) GetVersion() string {
	return New().Info.version
}

// GetDb 获取当前数据库实例
func (*Core) GetDb() *gorm.DB {
	return New().Db
}

// GetDb 获取当前数据库实例
func (*Core) GetPort() int {
	return New().Info.port
}

// GetRedis 获取当前Redis实例
func (*Core) GetRedis() *rds.Client {
	return New().Redis
}

// ConfigRegister 注册 配置
func (*Core) ConfigRegister() *config.Config {
	config.Init("config")
	New().Config = config.New()
	return New().Config
}

// ConfigRegister 注册 配置
func (*Core) SetConfigFile(file string) *config.Config {
	config.SetConfigFile(file)
	New().Config = config.New()
	return New().Config
}

// MongoRegister 注册Mongo
func (*Core) MongoRegister(Uri ...string) *mongo.Client {
	New().Mongo = gongDbDrive.Init(Uri...)
	return c.Mongo
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
		New().Info.env = env
	}
}

// Name 设置项目名称
func Name(name string) Option {
	return func(c *Core) {
		New().Info.Name = name
	}
}

// Port 设置端口号
func Port(port int) Option {
	return func(c *Core) {
		New().Info.port = port
	}
}

// Version 设置项目版本
func Version(version string) Option {
	return func(c *Core) {
		New().Info.version = version
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
		New().Db = mysqlDrive.New()
	}
}

// RedisRegister 设置Redis
func RedisRegister() Option {
	return func(c *Core) {
		New().Redis = redisDrive.New()
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
	New().Info.port = port
}
