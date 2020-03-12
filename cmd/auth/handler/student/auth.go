package student

import (
	"context"
	auth "rxt/cmd/auth/proto/student"
)

type Server struct{}

func (s *Server) Login(ctx context.Context, request *auth.LogicRequest, ) (*auth.LogicResponse, error) {

	return nil, nil
}
