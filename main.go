package main

import (
	"fmt"
	"github.com/erDong01/micro-kit/examples/account"
	"github.com/erDong01/micro-kit/examples/netgate"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	netgate.SERVER.Init()
	account.SERVER.Init()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	t := <-c
	fmt.Printf("server【%s】 exit ------- signal:[%v]", t)
}
