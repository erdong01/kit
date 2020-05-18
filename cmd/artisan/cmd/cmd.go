package cmd

import (
	"github.com/erDong01/micro-kit/cmd/artisan/cmd/I"
	"github.com/erDong01/micro-kit/cmd/artisan/cmd/app"
)

func NewApiCommand() I.ICommand {
	cmd := newCmd()
	cmd.AddCommand(
		app.New().GetCommand(),
	)
	return cmd
}
