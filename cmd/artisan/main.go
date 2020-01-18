package main

import (
	"fmt"
	"os"
	"rxt/cmd/artisan/cmd/app"
)

func main() {

	if err := app.New().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
