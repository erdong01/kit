package main

import (
	"fmt"
	"github.com/erDong01/micro-kit/cmd/rxstudent/app/cmd"
	"os"
)

func main() {
	cmd := cmd.NewAPICommand()
	if err := cmd.GetCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
