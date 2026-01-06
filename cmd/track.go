package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"specfirst/internal/snapshot"
	"specfirst/internal/store"
)

var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Manage parallel futures (tracks)",
}

var trackCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new track from current state",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		notes, _ := cmd.Flags().GetString("notes")

		mgr := snapshot.NewManager(store.TracksPath())
		if err := mgr.Create(name, []string{"track"}, notes); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Created track %s\n", name)
		return nil
	},
}

var trackListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tracks",
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := snapshot.NewManager(store.TracksPath())
		tracks, err := mgr.List()
		if err != nil {
			return err
		}
		for _, t := range tracks {
			fmt.Fprintln(cmd.OutOrStdout(), t)
		}
		return nil
	},
}

var trackSwitchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch workspace to a specific track (restores it)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		force, _ := cmd.Flags().GetBool("force")

		// Safety check similar to archive restore
		if !force {
			if _, err := os.Stat(store.ConfigPath()); err == nil {
				// We should probably explicitly warn that switching OVERWRITES current workspace.
				// Ideally we'd auto-snapshot "backup" track?
				// For now, consistent with archive logic: require --force if data exists.
				return fmt.Errorf("workspace has data; use --force to overwrite with track contents")
			}
		}

		mgr := snapshot.NewManager(store.TracksPath())
		if err := mgr.Restore(name, force); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Switched to track %s\n", name)
		return nil
	},
}

var trackDiffCmd = &cobra.Command{
	Use:   "diff <track-a> <track-b>",
	Short: "Compare artifacts between two tracks",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := snapshot.NewManager(store.TracksPath())
		added, removed, changed, err := mgr.Compare(args[0], args[1])
		if err != nil {
			return err
		}

		if len(added)+len(removed)+len(changed) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No differences detected.")
			return nil
		}
		if len(added) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Added:")
			for _, item := range added {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", item)
			}
		}
		if len(removed) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Removed:")
			for _, item := range removed {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", item)
			}
		}
		if len(changed) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Changed:")
			for _, item := range changed {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", item)
			}
		}
		return nil
	},
}

var trackMergeCmd = &cobra.Command{
	Use:   "merge <source-track>",
	Short: "Generate a merge plan to merge source track into current workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceTrack := args[0]

		// 1. Validate source track exists
		mgr := snapshot.NewManager(store.TracksPath())
		// We can't easily "compare" current workspace vs track using snapshot.Compare directly
		// because snapshot.Compare expects two *snapshot names*.
		// But "current workspace" isn't a snapshot.

		// Strategy: Create a temporary snapshot of current workspace?
		// Yes, "MERGE_HEAD" equivalent.
		currentSnapshot := "merge-target-temp"
		_ = os.RemoveAll(store.TracksPath(currentSnapshot)) // Clean up previous if any

		if err := mgr.Create(currentSnapshot, []string{"temp"}, "Temporary snapshot for merge"); err != nil {
			return fmt.Errorf("failed to snapshot current workspace for comparison: %w", err)
		}
		defer func() {
			_ = os.RemoveAll(store.TracksPath(currentSnapshot))
		}()

		// 2. Diff source track vs current (temp)
		// We want to see what is in Source that is different from Target.
		// Compare(left, right) -> changes from left to right?
		// compareHashes: right is "new".
		// We want to merge Source INTO Target.
		// So we want to see diff(Target, Source).
		added, removed, changed, err := mgr.Compare(currentSnapshot, sourceTrack)
		if err != nil {
			return err
		}

		if len(added)+len(removed)+len(changed) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Tracks are identical. Nothing to merge.")
			return nil
		}

		// 3. Generate Merge Prompt
		// We need to render internal/prompts/merge.md
		// We'll use a hack or simple replacement since we handle template rendering elsewhere differently?
		// internal/prompts templates are usually loaded via embedded FS or file?
		// Actually, `internal/prompts` is where the source code lives.
		// `store.TemplatesPath` user templates.
		// We should probably just hardcode the prompt usage here or use `prompt` package if it supports dynamic?
		// Let's manually construct the prompt context data and use `text/template` or the existing `template` pkg?
		// `specfirst/internal/template` pkg `Render` function takes a filename.

		// We need to point to the template file.
		// Ideally we ship standard prompts in binary.
		// For now, I'll attempt to assume it's available or write a temporary one.
		// Or better: Just generate the text directly here to avoid dependency on "installed" templates.

		mergePromptPath := "MERGE_PLAN_PROMPT.md"
		promptContent := fmt.Sprintf(`# Merge Plan for %s into Current Workspace

## Context
- **Source Track**: %s

## Differences

### Added
%s

### Removed
%s

### Changed
%s

## Instructions
1. Review the differences above.
2. Generate a 'MERGE_PLAN.md' checklist to safely merge these changes.
3. Identify any conflicts or high-risk files.
`, sourceTrack, sourceTrack, formatList(added), formatList(removed), formatList(changed))

		if err := os.WriteFile(mergePromptPath, []byte(promptContent), 0644); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Merge prompt generated at %s\n", mergePromptPath)
		fmt.Fprintln(cmd.OutOrStdout(), "Run this prompt with your LLM to generate a granular merge strategy.")

		return nil
	},
}

func formatList(items []string) string {
	if len(items) == 0 {
		return "(none)"
	}
	res := ""
	for _, item := range items {
		res += fmt.Sprintf("- %s\n", item)
	}
	return res
}

func init() {
	rootCmd.AddCommand(trackCmd)
	trackCmd.AddCommand(trackCreateCmd)
	trackCmd.AddCommand(trackListCmd)
	trackCmd.AddCommand(trackSwitchCmd)
	trackCmd.AddCommand(trackDiffCmd)
	trackCmd.AddCommand(trackMergeCmd)

	trackCreateCmd.Flags().String("notes", "", "notes for the track")
	trackSwitchCmd.Flags().Bool("force", false, "force overwrite of existing workspace data")
}
