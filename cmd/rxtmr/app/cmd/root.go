package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var env, version, name, configFile string

// New 构造方法
func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "tmr",
		Short: "A Http Api Gateway for TMR module",
	}

	if env == "" {
		c.PersistentFlags().StringVar(&name, "name", "", "app's name variable.")
		c.PersistentFlags().StringVarP(&version, "version", "v", "", "app's version variable.")
		c.PersistentFlags().StringVarP(&env, "env", "e", "local", "environment variable.")
		c.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config path.")

		err := viper.BindPFlags(c.PersistentFlags())
		if err != nil {
			log.Fatal(err)
		}
	}

	c.AddCommand(
		NewCmdServer(),
	)

	return c
}
