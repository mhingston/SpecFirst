package cmd

import (
	"specfirst/internal/engine"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run all non-blocking validations (lint, tasks, approvals, outputs)",
	RunE: func(cmd *cobra.Command, args []string) error {
		failOnWarnings, _ := cmd.Flags().GetBool("fail-on-warnings")

		eng, err := engine.Load(protocolFlag)
		if err != nil {
			return err
		}

		return eng.Check(failOnWarnings)
	},
}

func init() {
	checkCmd.Flags().Bool("fail-on-warnings", false, "exit with code 1 if warnings are found")
}
