package register

import (
	"github.com/erDong01/micro-kit/cache/impl"
	"github.com/erDong01/micro-kit/config"
	"github.com/erDong01/micro-kit/core"
	gongDbDrive "github.com/erDong01/micro-kit/db/mongoDb/drive"
	"github.com/erDong01/micro-kit/db/mysql"
	"github.com/erDong01/micro-kit/db/redis/drive"
)

// GlobalInit 全局初始化
func GlobalInit() *Register {
	return new(Register).ConfigRegister().
		RedisRegister().
		DbRegister().
		FacadeCacheRegister().
		SetPort(5001)
}

type Register struct {
}

// ConfigRegister 注册 配置
func (r *Register) ConfigRegister() *Register {
	config.Init("config")
	core.New().Config = config.New()
	return r
}

// DbRegister 注册 orm
func (r *Register) DbRegister() *Register {
	core.New().Db = mysql.New()
	return r
}

// MongoRegister 注册Mongo
func (r *Register) MongoRegister() *Register {
	core.New().Mongo = gongDbDrive.Init()
	return r
}

// FacadeCacheRegister 注册 缓存中心
func (r *Register) FacadeCacheRegister() *Register {
	core.New().Cache = impl.New()
	return r
}

// RedisRegister 注册 缓存中心
func (r *Register) RedisRegister() *Register {
	core.New().Redis = drive.New()
	return r
}

// SetName 设置名称
func (r *Register) SetName(name string) *Register {
	core.New().Info.Name = name
	return r
}

// SetPort 设置端口号
func (r *Register) SetPort(port int) *Register {
	core.SetPort(port)
	return r
}
