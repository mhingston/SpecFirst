package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completeSpecCmd = &cobra.Command{
	Use:   "complete-spec",
	Short: "Validate spec completion and optionally archive",
	Long: `Validate that all stages in the protocol are completed and all required approvals are present.
This is a validation tool to ensure rigor, but it is not a strict workflow requirement.
Use --warn-only to report missing stages or approvals without failing the command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		archiveFlag, _ := cmd.Flags().GetBool("archive")
		warnOnly, _ := cmd.Flags().GetBool("warn-only")
		version, _ := cmd.Flags().GetString("version")
		tags, _ := cmd.Flags().GetStringSlice("tag")
		notes, _ := cmd.Flags().GetString("notes")

		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		proto, err := loadProtocol(cfg.Protocol)
		if err != nil {
			return err
		}
		s, err := loadState()
		if err != nil {
			return err
		}
		s = ensureStateInitialized(s, proto)

		missing := []string{}
		for _, stage := range proto.Stages {
			if !s.IsStageCompleted(stage.ID) {
				missing = append(missing, stage.ID)
			}
		}
		if len(missing) > 0 {
			err := fmt.Errorf("spec is not complete, missing stages: %v", missing)
			if warnOnly {
				fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
			} else {
				return err
			}
		}

		missingApprovalRecords := missingApprovals(proto, s)
		if len(missingApprovalRecords) > 0 {
			err := fmt.Errorf("spec is not approved, missing approvals: %v", missingApprovalRecords)
			if warnOnly {
				fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
			} else {
				return err
			}
		}

		if len(missing) == 0 && len(missingApprovalRecords) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "All stages completed.")
		}

		if archiveFlag {
			if version == "" {
				version = s.SpecVersion
			}
			if version == "" {
				return fmt.Errorf("archive version is required (set --version or state.spec_version)")
			}
			if err := createArchive(version, tags, notes); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Archived version %s\n", version)
		}
		return nil
	},
}

func init() {
	completeSpecCmd.Flags().Bool("archive", false, "create an archive snapshot after completion")
	completeSpecCmd.Flags().Bool("warn-only", false, "report missing stages/approvals as warnings without failing")
	completeSpecCmd.Flags().String("version", "", "archive version (defaults to state.spec_version)")
	completeSpecCmd.Flags().StringSlice("tag", nil, "tag to apply to the archive (repeatable)")
	completeSpecCmd.Flags().String("notes", "", "notes for the archive")
}
