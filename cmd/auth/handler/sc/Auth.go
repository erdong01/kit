package sc

import (
	"context"
	"errors"
	auth "rxt/cmd/auth/proto/sc"
	"rxt/cmd/auth/service/scStudentUserService"
	"rxt/internal/jwt"
)

type Server struct{}

func (s *Server) Login(ctx context.Context, in *auth.AuthRequest) (*auth.AuthResponse, error) {
	token, err := scStudentUserService.New().Login(&scStudentUserService.LoginRequest{
		Mobile:  in.Mobile,
		SmsCode: in.SmsCode,
	})
	if err != nil {
		return nil, err
	}
	return &auth.AuthResponse{
		Token: &auth.Token{
			Token:     token.Token.Token,
			TokenType: token.Token.TokenType,
			ExpiresIn: token.Token.ExpiresIn,
		},
	}, nil

}

// Validate 解析验证Token
func (s *Server) Validate(ctx context.Context, in *auth.TokenRequest) (*auth.UserResponse, error) {
	claims, err := jwt.ValidateToken(in.GetToken())
	if err != nil {
		return nil, err
	}

	uid, ok := claims["sub"].(float64)
	if !ok {
		return nil, errors.New("Token无效")
	}

	user, err := scStudentUserService.New().ValidateToken(int(uid))
	if err != nil {
		return nil, err
	}

	return &auth.UserResponse{
		StudentUserId:     user.StudentUserId,
		StudentUserNo:     user.StudentUserNo,
		StudentUserName:   user.StudentUserName,
		StudentUserStatus: int32(user.StudentUserStatus),
		StudentUserMobile: user.StudentUserMobile,
		StudentUserHead:   user.StudentUserHead,
		LastLoginIp:       user.LastLoginIp,
		LastLoginTime:     user.LastLoginTime.Format("2006-01-02 15:04:05"),
	}, nil
}
