// Copyright (c) 2020-present Douglass Kirkley All Rights Reserved.
// See LICENSE.txt for license information.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"

	"github.com/dougkirkley/trex/config"
	"github.com/dougkirkley/trex/controller/deploy"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys trex environment",
	RunE:  deployCmdF,
}

func init() {
	RootCmd.AddCommand(DeployCmd)
}

func deployCmdF(command *cobra.Command, args []string) error {
	configFile := viper.GetString("config")
	cfg := config.New(configFile)
	ctrl := deploy.New(cfg)
	if err := ctrl.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}
