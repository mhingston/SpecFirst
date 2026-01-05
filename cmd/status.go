package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"specfirst/internal/store"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show workflow status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		proto, err := loadProtocol(activeProtocolName(cfg))
		if err != nil {
			return err
		}
		s, err := loadState()
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Project: %s\n", cfg.ProjectName)
		fmt.Fprintf(cmd.OutOrStdout(), "Protocol: %s\n", proto.Name)
		fmt.Fprintf(cmd.OutOrStdout(), "State: %s\n", store.StatePath())

		if len(s.CompletedStages) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Completed stages: (none)")
		} else {
			// Preserve completion order (chronological) rather than sorting alphabetically
			fmt.Fprintf(cmd.OutOrStdout(), "Completed stages: %v\n", s.CompletedStages)
		}

		if s.CurrentStage != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Current stage: %s\n", s.CurrentStage)
		}
		return nil
	},
}
