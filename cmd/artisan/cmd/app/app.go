package app

import (
	"github.com/spf13/cobra"
	"github.com/erDong01/micro-kit/cmd/artisan/cmd/I"
	"github.com/erDong01/micro-kit/cmd/artisan/http"
)

func New() I.ICommand {
	c := &cobra.Command{
		Use:   "app",
		Short: "A Http Api Gateway for RXT module",
		Run: func(cmd *cobra.Command, args []string) {
			http.New()
		},
	}

	return &I.Command{
		Cmd: c,
	}

}
