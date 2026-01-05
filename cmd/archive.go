package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"specfirst/internal/config"
	"specfirst/internal/protocol"
	"specfirst/internal/state"
	"specfirst/internal/store"
)

// validVersionPattern matches safe archive version names (alphanumeric, dots, hyphens, underscores)
var validVersionPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)

// validateArchiveVersion ensures the version string is safe for use as a directory name
func validateArchiveVersion(version string) error {
	if version == "" {
		return fmt.Errorf("archive version is required")
	}
	if len(version) > 128 {
		return fmt.Errorf("archive version too long (max 128 characters)")
	}
	if strings.Contains(version, "..") {
		return fmt.Errorf("invalid archive version: %q (contains path traversal)", version)
	}
	if !validVersionPattern.MatchString(version) {
		return fmt.Errorf("invalid archive version: %q (must be alphanumeric, may contain dots, hyphens, underscores)", version)
	}
	return nil
}

type archiveMetadata struct {
	Version         string    `json:"version"`
	Protocol        string    `json:"protocol"`
	ArchivedAt      time.Time `json:"archived_at"`
	StagesCompleted []string  `json:"stages_completed"`
	Tags            []string  `json:"tags,omitempty"`
	Notes           string    `json:"notes,omitempty"`
}

var archiveCmd = &cobra.Command{
	Use:   "archive <version>",
	Short: "Archive spec versions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		tags, _ := cmd.Flags().GetStringSlice("tag")
		notes, _ := cmd.Flags().GetString("notes")
		if err := createArchive(version, tags, notes); err != nil {
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
		entries, err := os.ReadDir(store.ArchivesPath())
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		versions := []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				versions = append(versions, entry.Name())
			}
		}
		sort.Strings(versions)
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
		if err := validateArchiveVersion(args[0]); err != nil {
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
		if err := validateArchiveVersion(args[0]); err != nil {
			return err
		}
		forceRestore, _ := cmd.Flags().GetBool("force")
		archiveRoot := store.ArchivesPath(args[0])
		if info, err := os.Stat(archiveRoot); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("archive not found: %s", args[0])
			}
			return err
		} else if !info.IsDir() {
			return fmt.Errorf("archive is not a directory: %s", args[0])
		}

		existingPaths := []string{
			store.ArtifactsPath(),
			store.GeneratedPath(),
			store.ProtocolsPath(),
			store.TemplatesPath(),
			store.ConfigPath(),
			store.StatePath(),
		}

		// Validate required archive directories before any destructive actions
		requiredDirs := []string{
			filepath.Join(archiveRoot, "protocols"),
			filepath.Join(archiveRoot, "templates"),
		}
		for _, dir := range requiredDirs {
			info, err := os.Stat(dir)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("archive is incomplete or corrupt: missing %s", filepath.Base(dir))
				}
				return fmt.Errorf("cannot access archive directory %s: %w", filepath.Base(dir), err)
			}
			if !info.IsDir() {
				return fmt.Errorf("archive is incomplete or corrupt: %s is not a directory", filepath.Base(dir))
			}
		}

		// Check if workspace has existing data and warn/require --force
		hasExisting := false
		for _, path := range existingPaths {
			if _, err := os.Stat(path); err == nil {
				hasExisting = true
				break
			} else if !os.IsNotExist(err) {
				return err
			}
		}
		if hasExisting {
			if !forceRestore {
				return fmt.Errorf("workspace already has data; use --force to overwrite existing workspace data")
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: overwriting existing workspace data\n")
		}
		// Validate archive integrity before restore
		requiredFiles := []string{
			filepath.Join(archiveRoot, "state.json"),
			filepath.Join(archiveRoot, "config.yaml"),
			filepath.Join(archiveRoot, "metadata.json"),
		}
		for _, required := range requiredFiles {
			if _, err := os.Stat(required); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("archive is incomplete or corrupt: missing %s", filepath.Base(required))
				}
				return fmt.Errorf("cannot access archive file %s: %w", filepath.Base(required), err)
			}
		}

		// Validate config/protocol consistency in archive
		archivedConfigPath := filepath.Join(archiveRoot, "config.yaml")
		metadataPath := filepath.Join(archiveRoot, "metadata.json")
		metadataData, err := os.ReadFile(metadataPath)
		if err != nil {
			return fmt.Errorf("cannot read archive metadata: %w", err)
		}
		var metadata archiveMetadata
		if err := json.Unmarshal(metadataData, &metadata); err != nil {
			return fmt.Errorf("cannot parse archive metadata: %w", err)
		}

		archivedCfg, err := config.Load(archivedConfigPath)
		if err != nil {
			return fmt.Errorf("cannot load archived config: %w", err)
		}
		if strings.TrimSpace(archivedCfg.Protocol) == "" {
			return fmt.Errorf("archive is incomplete or corrupt: config missing protocol")
		}

		archivedProtoPath := filepath.Join(archiveRoot, "protocols", archivedCfg.Protocol+".yaml")
		if _, err := os.Stat(archivedProtoPath); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("archive is incomplete or corrupt: missing protocol file %s", filepath.Base(archivedProtoPath))
			}
			return fmt.Errorf("cannot access archived protocol file %s: %w", filepath.Base(archivedProtoPath), err)
		}
		archivedProto, err := protocol.Load(archivedProtoPath)
		if err != nil {
			return fmt.Errorf("cannot load archived protocol: %w", err)
		}
		if metadata.Protocol != "" && archivedProto.Name != metadata.Protocol {
			return fmt.Errorf("archive metadata protocol mismatch: metadata=%s protocol=%s", metadata.Protocol, archivedProto.Name)
		}

		// Verify all referenced artifacts exist in the archive
		archivedState, err := state.Load(filepath.Join(archiveRoot, "state.json"))
		if err != nil {
			return fmt.Errorf("cannot load archived state: %w", err)
		}
		for stageID, output := range archivedState.StageOutputs {
			for _, artifactPath := range output.Files {
				rel, err := artifactRelFromState(artifactPath)
				if err != nil {
					return fmt.Errorf("archive is corrupt: invalid artifact path for stage %s: %s (%v)", stageID, artifactPath, err)
				}

				archivedArtifactPath := filepath.Join(archiveRoot, "artifacts", filepath.FromSlash(rel))
				if _, err := os.Stat(archivedArtifactPath); err != nil {
					if os.IsNotExist(err) {
						return fmt.Errorf("archive is corrupt: stage %s references missing artifact %s", stageID, rel)
					}
					return fmt.Errorf("cannot access archived artifact %s: %w", rel, err)
				}
			}
		}

		// Perform robust restore via staging area
		restoreStaging := store.SpecPath() + "_restore.tmp"
		_ = os.RemoveAll(restoreStaging)
		if err := ensureDir(restoreStaging); err != nil {
			return fmt.Errorf("failed to create restore staging directory: %w", err)
		}
		defer func() {
			_ = os.RemoveAll(restoreStaging)
		}()

		// Stage all components from archive
		if err := copyDir(filepath.Join(archiveRoot, "artifacts"), filepath.Join(restoreStaging, "artifacts")); err != nil {
			return fmt.Errorf("failed to stage artifacts: %w", err)
		}
		if err := copyDir(filepath.Join(archiveRoot, "generated"), filepath.Join(restoreStaging, "generated")); err != nil {
			return fmt.Errorf("failed to stage generated: %w", err)
		}
		if err := copyDirWithOpts(filepath.Join(archiveRoot, "protocols"), filepath.Join(restoreStaging, "protocols"), true); err != nil {
			return fmt.Errorf("failed to stage protocols: %w", err)
		}
		if err := copyDirWithOpts(filepath.Join(archiveRoot, "templates"), filepath.Join(restoreStaging, "templates"), true); err != nil {
			return fmt.Errorf("failed to stage templates: %w", err)
		}
		if err := copyFile(filepath.Join(archiveRoot, "config.yaml"), filepath.Join(restoreStaging, "config.yaml")); err != nil {
			return fmt.Errorf("failed to stage config: %w", err)
		}
		if err := copyFile(filepath.Join(archiveRoot, "state.json"), filepath.Join(restoreStaging, "state.json")); err != nil {
			return fmt.Errorf("failed to stage state: %w", err)
		}

		// Perform component-wise swap with rollback capability
		type backupEntry struct {
			original string
			backup   string
		}
		var backups []backupEntry
		success := false

		// Deferred rollback function
		defer func() {
			if !success {
				fmt.Fprintln(cmd.ErrOrStderr(), "Restore failed, rolling back changes...")
				// Restore in reverse order
				for i := len(backups) - 1; i >= 0; i-- {
					entry := backups[i]
					_ = os.RemoveAll(entry.original) // Remove partial restore
					if err := os.Rename(entry.backup, entry.original); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "Critical: failed to restore backup %s -> %s: %v\n", filepath.Base(entry.backup), filepath.Base(entry.original), err)
					}
				}
			} else {
				// Clean up backups on success
				for _, entry := range backups {
					_ = os.RemoveAll(entry.backup)
				}
			}
		}()

		for _, path := range existingPaths {
			stagedPath := filepath.Join(restoreStaging, filepath.Base(path))

			// If staged path doesn't exist, it means the archive doesn't have this component.
			// For a clean restore, we should effectively "remove" the component from the workspace.
			// However, we can't just delete it yet because we need rollback support.
			// So, if missing, we effectively stage an "empty" version or handle removal in the swap.
			// To simplify, if missing in staging, we treat it as "restore to empty/absent".

			if _, err := os.Stat(stagedPath); err != nil {
				if os.IsNotExist(err) {
					// Component missing in archive.
					// If the workspace has it, we should remove it.
					// We'll mark this by creating an empty placeholder flag or just handling it in the swap logic?
					// Simpler: If it's a directory that should exist (like artifacts), create an empty dir in staging.
					// If it's a file, we might not want to create an empty file.

					base := filepath.Base(path)
					if base == "artifacts" || base == "generated" || base == "protocols" || base == "templates" {
						// Ensure empty dir exists in staging so we swap to an empty dir
						if err := ensureDir(stagedPath); err != nil {
							return fmt.Errorf("failed to create empty staging dir for %s: %w", base, err)
						}
					} else {
						// For config.yaml and state.json, they are REQUIRED inside the archive (validated earlier).
						// So we shouldn't reach here for them.
						// If we do, it's safe to skip or error.
						continue
					}
				} else {
					return err
				}
			}

			// Rename existing out of the way
			if _, err := os.Stat(path); err == nil {
				oldPath := path + ".old"
				_ = os.RemoveAll(oldPath)
				if err := os.Rename(path, oldPath); err != nil {
					return fmt.Errorf("failed to backup existing %s: %w", filepath.Base(path), err)
				}
				backups = append(backups, backupEntry{original: path, backup: oldPath})
			}

			// Rename staged into place
			if err := os.Rename(stagedPath, path); err != nil {
				return fmt.Errorf("failed to restore %s: %w", filepath.Base(path), err)
			}
		}

		// One final check: did we restore the protocol file?
		// (Already validated in integrity check, but good to be sure after swap)
		newCfg, err := config.Load(store.ConfigPath())
		if err == nil {
			if _, err := os.Stat(store.ProtocolsPath(newCfg.Protocol + ".yaml")); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Warning: restored workspace may be inconsistent: protocol %q file missing\n", newCfg.Protocol)
			}
		}

		success = true
		fmt.Fprintf(cmd.OutOrStdout(), "Restored archive %s\n", args[0])
		return nil
	},
}

var archiveCompareCmd = &cobra.Command{
	Use:   "compare <version-a> <version-b>",
	Short: "Compare archived artifacts between versions",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateArchiveVersion(args[0]); err != nil {
			return err
		}
		if err := validateArchiveVersion(args[1]); err != nil {
			return err
		}
		leftArchive := store.ArchivesPath(args[0])
		rightArchive := store.ArchivesPath(args[1])
		for _, archive := range []string{leftArchive, rightArchive} {
			info, err := os.Stat(archive)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("archive not found: %s", filepath.Base(archive))
				}
				return err
			}
			if !info.IsDir() {
				return fmt.Errorf("archive is not a directory: %s", filepath.Base(archive))
			}
		}
		leftRoot := filepath.Join(leftArchive, "artifacts")
		rightRoot := filepath.Join(rightArchive, "artifacts")

		left, err := collectFileHashes(leftRoot)
		if err != nil {
			return err
		}
		right, err := collectFileHashes(rightRoot)
		if err != nil {
			return err
		}

		added, removed, changed := compareHashes(left, right)
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

func compareHashes(left map[string]string, right map[string]string) ([]string, []string, []string) {
	added := []string{}
	removed := []string{}
	changed := []string{}

	for path, hash := range right {
		if leftHash, ok := left[path]; ok {
			if leftHash != hash {
				changed = append(changed, path)
			}
		} else {
			added = append(added, path)
		}
	}
	for path := range left {
		if _, ok := right[path]; !ok {
			removed = append(removed, path)
		}
	}

	sort.Strings(added)
	sort.Strings(removed)
	sort.Strings(changed)
	return added, removed, changed
}

func createArchive(version string, tags []string, notes string) error {
	if err := validateArchiveVersion(version); err != nil {
		return err
	}
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

	// Validate that all referenced artifacts actually exist before writing the archive
	for stageID, output := range s.StageOutputs {
		for _, artifactPath := range output.Files {
			abs, err := artifactAbsFromState(artifactPath)
			if err != nil {
				return fmt.Errorf("invalid artifact path for stage %s: %s (%v)", stageID, artifactPath, err)
			}
			if _, err := os.Stat(abs); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("missing artifact for stage %s: %s (archive aborted)", stageID, artifactPath)
				}
				return fmt.Errorf("cannot access artifact for stage %s: %s (%v)", stageID, artifactPath, err)
			}
		}
	}

	archiveRoot := store.ArchivesPath(version)
	tmpArchiveRoot := archiveRoot + ".tmp"

	// Ensure parent directory exists first
	if err := ensureDir(filepath.Dir(archiveRoot)); err != nil {
		return err
	}

	// Clean up any stale temp directory
	_ = os.RemoveAll(tmpArchiveRoot)

	// Atomic check-and-create for the final destination to prevent overwriting
	if _, err := os.Stat(archiveRoot); err == nil {
		return fmt.Errorf("archive already exists: %s", version)
	}

	if err := os.Mkdir(tmpArchiveRoot, 0755); err != nil {
		return fmt.Errorf("failed to create temporary archive directory: %w", err)
	}

	cleanupArchive := true
	defer func() {
		if cleanupArchive {
			_ = os.RemoveAll(tmpArchiveRoot)
		}
	}()

	if err := copyDir(store.ArtifactsPath(), filepath.Join(tmpArchiveRoot, "artifacts")); err != nil {
		return err
	}
	if err := copyDir(store.GeneratedPath(), filepath.Join(tmpArchiveRoot, "generated")); err != nil {
		return err
	}
	if err := copyDirWithOpts(store.ProtocolsPath(), filepath.Join(tmpArchiveRoot, "protocols"), true); err != nil {
		return err
	}
	if err := copyDirWithOpts(store.TemplatesPath(), filepath.Join(tmpArchiveRoot, "templates"), true); err != nil {
		return err
	}
	if err := copyFile(store.ConfigPath(), filepath.Join(tmpArchiveRoot, "config.yaml")); err != nil {
		return err
	}
	if err := copyFile(store.StatePath(), filepath.Join(tmpArchiveRoot, "state.json")); err != nil {
		return err
	}

	metadata := archiveMetadata{
		Version:         version,
		Protocol:        proto.Name,
		ArchivedAt:      time.Now().UTC(),
		StagesCompleted: s.CompletedStages,
		Tags:            tags,
		Notes:           notes,
	}
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.WriteFile(filepath.Join(tmpArchiveRoot, "metadata.json"), data, 0644); err != nil {
		return err
	}

	// Final atomic rename
	if err := os.Rename(tmpArchiveRoot, archiveRoot); err != nil {
		return fmt.Errorf("failed to finalize archive: %w", err)
	}

	cleanupArchive = false
	return nil
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
