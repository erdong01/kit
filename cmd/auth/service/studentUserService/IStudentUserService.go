package studentUserService

import (
	"net"
	auth "rxt/cmd/auth/proto/student"
	"rxt/internal/core"
)

type container struct {
	I
}

func New(param ...*core.Core) *container {
	exam := &V1{}
	exam.Init(param...)
	return &container{exam}
}

type I interface {
	Login(param *Param) (*auth.LogicResponse, error)
}

type Param struct {
	StudentUserLoginName string
	Password             string
	Ip                    net.IP
}
