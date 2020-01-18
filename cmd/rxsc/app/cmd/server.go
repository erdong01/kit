package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"rxt/cmd/rxsc/bootstrap"
	"rxt/cmd/rxsc/route"
	"rxt/internal/core"
	"rxt/internal/http"
)

func NewCmdServer() ICommand {
	c := &cobra.Command{
		Use:   "server",
		Short: "A http gateway server powered by gin",
		Run: func(cmd *cobra.Command, args []string) {
			bootstrap.App(name, env, version)
			defer core.Close()
			http.Init(route.Init(gin.Default()), core.New().GetPort())
		},
	}
	return &command{
		cmd: c,
	}
}
