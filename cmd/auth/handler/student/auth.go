package student

import (
	"context"
	"net"
	"net/http"
	auth "rxt/cmd/auth/proto/student"
	"rxt/cmd/auth/service/studentUserService"
	"rxt/internal/wrong"
)

type Server struct{}

func (s *Server) Login(ctx context.Context, request *auth.LogicRequest, ) (*auth.LogicResponse, error) {
	response, err := studentUserService.New().Login(&studentUserService.Param{
		StudentUserLoginName: request.StudentUserLoginName,
		Password:             request.Password,
		Ip:                   net.ParseIP(request.Ip),
	})
	if err != nil {
		return nil, wrong.New(http.StatusExpectationFailed, err)
	}
	return response, nil
}
