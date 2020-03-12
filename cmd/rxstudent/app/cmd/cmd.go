package cmd

import (
	"github.com/spf13/cobra"
)

// ICommand command接口
type ICommand interface {
	GetCommand() *cobra.Command
	addCommand(cmds ...*cobra.Command)
}

type command struct {
	cmd *cobra.Command
}

func (c *command) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *command) addCommand(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		c.cmd.AddCommand(cmd)
	}
}

// NewAPICommand 命令行入口
func NewAPICommand() ICommand {
	cmd := newCmd()
	cmd.addCommand(
		newCmdVersion().GetCommand(),
		NewCmdServer().GetCommand(),
	)

	return cmd
}
