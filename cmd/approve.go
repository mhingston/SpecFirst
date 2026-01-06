package cmd

import (
	"fmt"
	"os"
	"specfirst/internal/engine"

	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:   "approve <stage-id>",
	Short: "Record an approval for a stage",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stageID := args[0]
		role, _ := cmd.Flags().GetString("role")
		approvedBy, _ := cmd.Flags().GetString("by")
		notes, _ := cmd.Flags().GetString("notes")

		if role == "" {
			return fmt.Errorf("role is required")
		}
		if approvedBy == "" {
			approvedBy = os.Getenv("USER")
		}
		if approvedBy == "" {
			approvedBy = os.Getenv("USERNAME")
		}

		eng, err := engine.Load(protocolFlag)
		if err != nil {
			return err
		}

		warnings, err := eng.ApproveStage(stageID, role, approvedBy, notes)
		if err != nil {
			return err
		}
		for _, w := range warnings {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: %s\n", w)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Recorded approval for %s (role: %s)\n", stageID, role)
		return nil
	},
}

func init() {
	approveCmd.Flags().String("role", "", "role required for approval")
	approveCmd.Flags().String("by", "", "who approved")
	approveCmd.Flags().String("notes", "", "approval notes")
}
