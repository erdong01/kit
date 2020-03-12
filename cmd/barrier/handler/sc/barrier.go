package sc

import (
	"context"
	pb "rxt/cmd/barrier/proto/sc"
	"rxt/cmd/barrier/service/skip"
)

type Server struct{}

// Skip 跳关
func (s *Server) Skip(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	result, err := skip.New().Skip(in)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Result: result}, nil
}
