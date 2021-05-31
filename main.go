package main

import "github.com/erDong01/micro-kit/network"

func main() {
	var s network.ServerSocket
	s.Init()
	s.StartTcpServer()
}
