package cmd

import (
	"fmt"
	"os"
	"specfirst/internal/app"

	"github.com/spf13/cobra"
)

var attestCmd = &cobra.Command{
	Use:               "attest <stage-id>",
	Short:             "Record an attestation for a stage",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: stageIDCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		stageID := args[0]
		role, _ := cmd.Flags().GetString("role")
		attestedBy, _ := cmd.Flags().GetString("by")
		status, _ := cmd.Flags().GetString("status")
		notes, _ := cmd.Flags().GetString("rationale")
		conditions, _ := cmd.Flags().GetStringSlice("condition")

		if role == "" {
			return fmt.Errorf("role is required")
		}
		if status == "" {
			return fmt.Errorf("status is required")
		}
		if attestedBy == "" {
			attestedBy = os.Getenv("USER")
		}
		if attestedBy == "" {
			attestedBy = os.Getenv("USERNAME")
		}

		application, err := app.Load(protocolFlag)
		if err != nil {
			return err
		}

		warnings, err := application.AttestStage(stageID, role, attestedBy, status, notes, conditions)
		if err != nil {
			return err
		}
		for _, w := range warnings {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: %s\n", w)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Recorded attestation for %s (role: %s, status: %s)\n", stageID, role, status)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(attestCmd)
	attestCmd.Flags().String("role", "", "role providing the attestation")
	attestCmd.Flags().String("status", "approved", "status: approved|approved_with_conditions|needs_changes|rejected")
	attestCmd.Flags().String("by", "", "who is attesting (defaults to $USER)")
	attestCmd.Flags().String("rationale", "", "rationale or notes")
	attestCmd.Flags().StringSlice("condition", []string{}, "conditions for approval (repeatable)")

	_ = attestCmd.RegisterFlagCompletionFunc("role", attestRoleCompletions)
	_ = attestCmd.RegisterFlagCompletionFunc("status", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return filterPrefix([]string{"approved", "approved_with_conditions", "needs_changes", "rejected"}, toComplete), cobra.ShellCompDirectiveNoFileComp
	})
}
