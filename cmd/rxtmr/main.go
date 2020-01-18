package main

import (
	"fmt"
	"os"
	"rxt/cmd/rxtmr/app/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
