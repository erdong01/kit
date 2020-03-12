package I

import "github.com/spf13/cobra"

type ICommand interface {
	GetCommand() *cobra.Command
	AddCommand(...*cobra.Command)
}
type Command struct {
	Cmd *cobra.Command
}

func (c *Command) GetCommand() *cobra.Command {
	return c.Cmd
}

func (c *Command) AddCommand(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		c.Cmd.AddCommand(cmd)
	}
}
