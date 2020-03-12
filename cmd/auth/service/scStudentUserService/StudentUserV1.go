package scStudentUserService

import (
	"errors"
	"net/http"
	"rxt/cmd/auth/model"
	"rxt/cmd/auth/service/baseService"
	"rxt/internal/jwt"
	"rxt/internal/wrong"
)

type StudentuserV1 struct {
	baseService.Service
}

// 登录
func (sc *StudentuserV1) Login(l *LoginRequest) (*LoginResult, error) {
	var expiresIn int32
	sc.Cache.Set("test", 36000, 100)
	ScStudentUser := model.ScStudentUser{}
	err := sc.Db.Where("student_user_mobile = ? ", l.Mobile).First(&ScStudentUser).Error
	if err != nil {
		return nil, wrong.New(http.StatusExpectationFailed, err, "用户不存在！")
	}
	User := jwt.User{
		Sub: ScStudentUser.StudentUserId,
		Prv: "App\\Modes\\Sc\\StudentUser",
	}
	token, _ := jwt.GenerateToken(User)
	sc.Cache.Get("test", &expiresIn)
	var tokenResult = TokenResult{
		Token:     token,
		TokenType: ScStudentUser.StudentUserMobile,
		ExpiresIn: expiresIn,
	}
	return &LoginResult{
		Token: &tokenResult,
	}, nil
}

// ValidateToken 验证Token
func (sc *StudentuserV1) ValidateToken(id int) (*model.ScStudentUser, error) {
	ScStudentUser := model.ScStudentUser{}
	err := sc.Db.Where("student_user_id = ? ", id).First(&ScStudentUser).Error
	if err != nil {
		return nil, errors.New("用户不存在！")
	}

	return &ScStudentUser, nil
}
