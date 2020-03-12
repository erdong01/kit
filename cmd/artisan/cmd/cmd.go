package cmd

import (
	"rxt/cmd/artisan/cmd/I"
	"rxt/cmd/artisan/cmd/app"
)

func NewApiCommand() I.ICommand {
	cmd := newCmd()
	cmd.AddCommand(
		app.New().GetCommand(),
	)
	return cmd
}
