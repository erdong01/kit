package cmd

import (
	"github.com/spf13/cobra"
	"rxt/cmd/rxtmr/bootstrap"
	"rxt/internal/core"
)

func NewCmdServer() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "A http gateway server powered by gin",
		Run: func(cmd *cobra.Command, args []string) {
			bootstrap.App(env, name, version)
			defer core.Close()
		},
	}
}
