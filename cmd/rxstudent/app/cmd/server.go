package cmd

import (
	"rxt/cmd/rxstudent/bootstrap"
	"rxt/cmd/rxstudent/route"
	"rxt/internal/core"
	"rxt/internal/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func NewCmdServer() ICommand {
	c := &cobra.Command{
		Use:   "server",
		Short: "A http gateway server powered by gin",
		Run: func(cmd *cobra.Command, args []string) {
			bootstrap.App()
			defer core.Close()
			http.Init(route.Init(gin.Default()), core.New().GetPort())
		},
	}
	return &command{
		cmd: c,
	}
}
