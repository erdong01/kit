package app

import (
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"rxt/cmd/artisan/http"
)

var env, version, name, configFile string
var (
	g errgroup.Group
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "app",
		Short: "A Http Api Gateway for RXT module",
		Run: func(cmd *cobra.Command, args []string) {
			http.New()
		},
	}

	return c
}
