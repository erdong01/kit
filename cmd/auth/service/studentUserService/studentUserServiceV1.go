package studentUserService

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"rxt/cmd/auth/model"
	auth "rxt/cmd/auth/proto/student"
	"rxt/cmd/auth/service/baseService"
	"rxt/internal/wrong"
	"time"
)

type V1 struct {
	baseService.Service
}

func (c *V1) Login(param *Param) (*auth.LogicResponse, error) {
	var result *auth.LogicResponse
	var studentLogic model.StudentUserLogin
	c.Db.Where("student_user_login_name = ? ", param.StudentUserLoginName).First(&studentLogic)
	password := []byte(param.Password)
	dataPassword := []byte(studentLogic.Password)
	//密码验证
	err := bcrypt.CompareHashAndPassword(dataPassword, password)
	if err != nil {
		return nil, wrong.New(http.StatusExpectationFailed, errors.New("密码错误或账号已经被禁用！"))
	}
	var studentUser model.StudentUser
	c.Db.Where("student_user_no = ? ", studentLogic.StudentUserNo).First(&studentUser)

	studentUser.LastLoginIp = int64(ip2long(param.Ip))
	timeNow := time.Now()
	studentUser.LastLoginTime = &timeNow
	c.Db.Save(&studentUser)
	return result, nil
}

func ip2long(ip net.IP) uint32 {
	a := uint32(ip[12])
	b := uint32(ip[13])
	c := uint32(ip[14])
	d := uint32(ip[15])
	return uint32(a<<24 | b<<16 | c<<8 | d)
}

func long2ip(ip uint32) net.IP {
	a := byte((ip >> 24) & 0xFF)
	b := byte((ip >> 16) & 0xFF)
	c := byte((ip >> 8) & 0xFF)
	d := byte(ip & 0xFF)
	return net.IPv4(a, b, c, d)
}
