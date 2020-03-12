package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCmdVersion() ICommand {
	cmd := &command{
		cmd: &cobra.Command{
			Use:   "version",
			Short: "Print the version of api",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(version)
			},
		},
	}

	return cmd
}
