// Copyright (c) 2020-present Douglass Kirkley All Rights Reserved.
// See LICENSE.txt for license information.

package cmd

import (
	"os"
	"log"

	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/dougkirkley/trex/terraform"
)

var DestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Deletes trex environment",
	RunE: destroyCmdF,
}

func init() {
	RootCmd.AddCommand(DestroyCmd)
}

func destroyCmdF(command *cobra.Command, args []string) error {
	cmd := terraform.DestroyCommand(terraform.ApplyArguments{
		NoColor: aws.Bool(true),
	})

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}
