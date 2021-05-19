// Copyright (c) 2020-present Douglass Kirkley All Rights Reserved.
// See LICENSE.txt for license information.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/dougkirkley/trex/config"
	"log"
)

type Command = cobra.Command

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use:   "trex",
	Short: "CLI tool for creating serverless api using terraform and aws lambda",
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "trex.yaml", "Configuration file to use.")

	viper.SetEnvPrefix("trex")
	viper.SetConfigName("trex")
	viper.SetConfigType("yaml")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cfg := config.New("trex.yaml")
			err := cfg.GenerateDefaultConfig()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
