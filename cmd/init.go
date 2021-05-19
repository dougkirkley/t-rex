// Copyright (c) 2020-present Douglass Kirkley All Rights Reserved.
// See LICENSE.txt for license information.

package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dougkirkley/trex/config"
	initialize "github.com/dougkirkley/trex/controller/init"
	"github.com/dougkirkley/trex/terraform"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes new trex environment",
	RunE:  initCmdF,
}

func init() {
	RootCmd.AddCommand(InitCmd)
}

func initCmdF(command *cobra.Command, args []string) error {
	configFile := viper.GetString("config")
	CmdPrettyPrintln("Initializing trex environment")
	cfg := config.New(configFile)

	if !terraform.IsTerraformInstalled() {
		log.Println("terraform is not installed. Installing...")
		if err := terraform.InstallTerraform(); err != nil {
			log.Fatal(err)
		}
	}

	ctrl := initialize.New(cfg)
	if err := ctrl.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}