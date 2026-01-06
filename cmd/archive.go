package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"specfirst/internal/snapshot"
	"specfirst/internal/store"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <version>",
	Short: "Archive spec versions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		tags, _ := cmd.Flags().GetStringSlice("tag")
		notes, _ := cmd.Flags().GetString("notes")

		mgr := snapshot.NewManager(store.ArchivesPath())
		if err := mgr.Create(version, tags, notes); err != nil {
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
		mgr := snapshot.NewManager(store.ArchivesPath())
		versions, err := mgr.List()
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
	Use:   "show <version>",
	Short: "Show archive metadata",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := snapshot.ValidateName(args[0]); err != nil {
			return err
		}
		path := filepath.Join(store.ArchivesPath(args[0]), "metadata.json")
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	},
}

var archiveRestoreCmd = &cobra.Command{
	Use:   "restore <version>",
	Short: "Restore an archive snapshot into .specfirst",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			if _, err := os.Stat(store.ConfigPath()); err == nil {
				return fmt.Errorf("workspace has data; use --force to overwrite")
			}
		}

		mgr := snapshot.NewManager(store.ArchivesPath())
		if err := mgr.Restore(version, force); err != nil {
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
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := snapshot.NewManager(store.ArchivesPath())
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

func init() {
	archiveCmd.Flags().StringSlice("tag", nil, "tag to apply to the archive (repeatable)")
	archiveCmd.Flags().String("notes", "", "notes for the archive")

	archiveRestoreCmd.Flags().Bool("force", false, "force overwrite of existing workspace data")

	archiveCmd.AddCommand(archiveListCmd)
	archiveCmd.AddCommand(archiveShowCmd)
	archiveCmd.AddCommand(archiveRestoreCmd)
	archiveCmd.AddCommand(archiveCompareCmd)
}
