package register

import (
	"github.com/erDong01/micro-kit/config"
	"github.com/erDong01/micro-kit/core"
	gongDbDrive "github.com/erDong01/micro-kit/db/mongoDB/drive"
)

// GlobalInit 全局初始化
func GlobalInit() *Register {
	return new(Register).ConfigRegister().
		DbRegister().
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

// MongoRegister 注册Mongo
func (r *Register) MongoRegister() *Register {
	core.New().Mongo = gongDbDrive.Init()
	return r
}

// SetName 设置名称
func (r *Register) SetName(name string) *Register {
	core.New().Info.Name = name
	return r
}

// SetPort 设置端口号
func (r *Register) SetPort(port int) *Register {
	core.New().Info.Port = port
	return r
}
