package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"specfirst/internal/app"
	"specfirst/internal/domain"
	"specfirst/internal/repository"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <version>",
	Short: "Archive spec versions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		tags, _ := cmd.Flags().GetStringSlice("tag")
		notes, _ := cmd.Flags().GetString("notes")

		application, err := app.Load(protocolFlag)
		if err != nil {
			return err
		}

		if err := application.CreateSnapshot(version, tags, notes); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Archived version %s\n", version)
		return nil
	},
}

var archiveListCmd = &cobra.Command{
	Use:   "list",
	Short: "List archives",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ListSnapshots is repository logic, but accessed via App mostly for convenience.
		// Since List doesn't require loaded state, we can use app.ListSnapshots if implemented statically?
		// No, app methods are instance methods.
		// But loading app entire state just to list archives is heavy?
		// Let's use repository directly for LIST if we want speed, or instantiate dummy app?
		// Or update app to have static utility?
		// For now, let's load app. It's safe.

		// Actually, loading app requires .specfirst to exist.
		// If listing archives, maybe .specfirst exists.
		// But what if we are listing invalid workspace?
		// Using repository direct is safer for LIST.
		repo := repository.NewSnapshotRepository(repository.ArchivesPath())
		versions, err := repo.List()
		if err != nil {
			return err
		}
		for _, version := range versions {
			fmt.Fprintln(cmd.OutOrStdout(), version)
		}
		return nil
	},
}

var archiveShowCmd = &cobra.Command{
	Use:               "show <version>",
	Short:             "Show archive metadata",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: archiveVersionCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !domain.IsValidSnapshotName(args[0]) {
			return fmt.Errorf("invalid snapshot name: %s", args[0])
		}
		// Direct file access via repo path helper
		path := filepath.Join(repository.ArchivesPath(args[0]), "metadata.json")
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	},
}

var archiveRestoreCmd = &cobra.Command{
	Use:               "restore <version>",
	Short:             "Restore an archive snapshot into .specfirst",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: archiveVersionCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			// Sev2 Fix: Unsafe check. Previously only checked ConfigPath.
			// Now we check if .specfirst exists and has any content.
			specDir := repository.SpecPath()
			if info, err := os.Stat(specDir); err == nil && info.IsDir() {
				entries, err := os.ReadDir(specDir)
				if err == nil && len(entries) > 0 {
					return fmt.Errorf("workspace is not empty; use --force to overwrite")
				}
			}
		}

		// Restore doesn't need "App" loaded since it overwrites it.
		// Uses repository directly or app helper?
		// We can't load app if config is corrupt/missing, which Restore might fix.
		// So using repo directly is better or static helper.
		repo := repository.NewSnapshotRepository(repository.ArchivesPath())
		if err := repo.Restore(version, force); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Restored archive %s\n", version)
		return nil
	},
}

var archiveCompareCmd = &cobra.Command{
	Use:   "compare <version-a> <version-b>",
	Short: "Compare archived artifacts between versions",
	Args:  cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 1 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return filterPrefix(loadArchiveVersions(), toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewSnapshotRepository(repository.ArchivesPath())
		added, removed, changed, err := repo.Compare(args[0], args[1])
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

func init() {
	archiveCmd.Flags().StringSlice("tag", nil, "tag to apply to the archive (repeatable)")
	archiveCmd.Flags().String("notes", "", "notes for the archive")

	archiveRestoreCmd.Flags().Bool("force", false, "force overwrite of existing workspace data")

	archiveCmd.AddCommand(archiveListCmd)
	archiveCmd.AddCommand(archiveShowCmd)
	archiveCmd.AddCommand(archiveRestoreCmd)
	archiveCmd.AddCommand(archiveCompareCmd)
}
