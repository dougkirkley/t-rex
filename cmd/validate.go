// Copyright (c) 2020-present Douglass Kirkley All Rights Reserved.
// See LICENSE.txt for license information.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dougkirkley/trex/config"
)

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates trex config",
	RunE: validateCmdF,
}

func init() {
	RootCmd.AddCommand(ValidateCmd)
}

func validateCmdF(command *cobra.Command, args []string) error {
	configFile := viper.GetString("config")
	cfg := config.New(configFile)
	_, err := cfg.LoadConfig()
	if err != nil {
		return err
	}
	return nil
}
