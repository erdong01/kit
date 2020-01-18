package main

import (
	"fmt"
	"os"
	"rxt/cmd/rxsc/app/cmd"
)

func main() {
	cmd := cmd.NewAPICommand()

	if err := cmd.GetCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
