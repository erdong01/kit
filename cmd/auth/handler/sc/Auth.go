package sc

import (
	"context"
	auth "rxt/cmd/auth/proto/sc"
	"rxt/cmd/auth/service/scStudentUser"
)

type Server struct{}

func (s *Server) Login(ctx context.Context, in *auth.AuthRequest) (*auth.AuthResponse, error) {
	token, err := scStudentUser.New().Login(&scStudentUser.LoginRequest{
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
