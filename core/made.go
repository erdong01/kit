package core

import (
	"reflect"

	"github.com/erDong01/micro-kit/config"
	gongDbDrive "github.com/erDong01/micro-kit/db/mongoDB/drive"
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
	GetPort() int
	ConfigRegister() *config.Config
	SetConfigFile(file string) *config.Config
	MongoRegister(Uri ...string) *mongo.Client
}

// Init 启动其他服务
func (*Core) Init() {
	for _, o := range c.opts {
		o(c)
	}
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
	return New().Info.Port
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
		New().Info.Port = port
	}
}

// Version 设置项目版本
func Version(version string) Option {
	return func(c *Core) {
		New().Info.version = version
	}
}

func ConfigRegister() Option {
	return func(c *Core) {
		config.Init("config")
	}
}

func Bind(a any, b any) {
	csValue := reflect.ValueOf(a).Elem()
	csValue.Set(reflect.ValueOf(b))
}
