package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"specfirst/internal/protocol"
	"specfirst/internal/state"
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

		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		proto, err := loadProtocol(activeProtocolName(cfg))
		if err != nil {
			return err
		}
		if !approvalDeclared(proto.Approvals, stageID, role) {
			return fmt.Errorf("approval not declared in protocol: stage=%s role=%s", stageID, role)
		}

		s, err := loadState()
		if err != nil {
			return err
		}
		s = ensureStateInitialized(s, proto)

		// Warn if stage is not yet completed
		if !s.IsStageCompleted(stageID) {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: stage %s is not yet completed; approval recorded preemptively\n", stageID)
		}

		// Check for existing approval with same role and update instead of duplicating
		newApproval := state.Approval{
			Role:       role,
			ApprovedBy: approvedBy,
			ApprovedAt: time.Now().UTC(),
			Notes:      notes,
		}
		updated := false
		for i, existing := range s.Approvals[stageID] {
			if existing.Role == role {
				s.Approvals[stageID][i] = newApproval
				updated = true
				fmt.Fprintf(cmd.ErrOrStderr(), "Note: updating existing approval for role %s\n", role)
				break
			}
		}
		if !updated {
			s.Approvals[stageID] = append(s.Approvals[stageID], newApproval)
		}

		if err := saveState(s); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Recorded approval for %s (role: %s)\n", stageID, role)
		return nil
	},
}

func approvalDeclared(approvals []protocol.Approval, stageID string, role string) bool {
	for _, approval := range approvals {
		if approval.Stage == stageID && approval.Role == role {
			return true
		}
	}
	return false
}

func init() {
	approveCmd.Flags().String("role", "", "role required for approval")
	approveCmd.Flags().String("by", "", "who approved")
	approveCmd.Flags().String("notes", "", "approval notes")
}
