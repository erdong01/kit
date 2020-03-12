package main

import (
	"fmt"
	"os"
	"rxt/cmd/artisan/cmd"
)

func main() {
	cmd := cmd.NewApiCommand()
	if err := cmd.GetCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
