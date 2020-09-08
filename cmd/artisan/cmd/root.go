package cmd

import (
	"fmt"
	"github.com/erDong01/micro-kit/cmd/artisan/cmd/I"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var env, version, name, configFile string

func newCmd() I.ICommand {
	c := &cobra.Command{
		Use:   "api",
		Short: "A Http Api Gateway",
	}

	if env == "" {
		c.PersistentFlags().StringVar(&name, "name", "", "handler's name variable.")
		c.PersistentFlags().StringVarP(&version, "version", "v", "", "handler's version variable.")
		c.PersistentFlags().StringVarP(&env, "env", "e", "local", "environment variable.")
		c.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config path.")

		err := viper.BindPFlags(c.PersistentFlags())
		if err != nil {
			fmt.Println(err)
		}
	}

	return &I.Command{
		Cmd: c,
	}
}
