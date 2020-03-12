package main

import (
	"log"
	"net"
	report2 "rxt/cmd/report/handler"
	"rxt/cmd/report/proto/report"
	"rxt/internal/core"

	"google.golang.org/grpc"
)

const (
	PORT = ":50001"
)

var name, env, version string

func main() {
	core.Make(
		core.DbRegister(),
	).Init()
	defer core.Close()

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	report.RegisterReportRpcServer(s, &report2.Server{})
	log.Println("rpc服务已经开启")
	s.Serve(lis)
}
