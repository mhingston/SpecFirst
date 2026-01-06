package cmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"specfirst/internal/assets"
	"specfirst/internal/snapshot"
	"specfirst/internal/store"

	"github.com/spf13/cobra"
)

func TestArchiveRestoreCleansWorkspace(t *testing.T) {
	// Setup workspace with "dirty" state (extra artifacts)
	wd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(wd) })
	tmp := t.TempDir()
	os.Chdir(tmp)

	// Create a "dirty" artifact that should be removed by restore
	dirtyArtifact := store.ArtifactsPath("requirements", "dirty.md")
	os.MkdirAll(filepath.Dir(dirtyArtifact), 0755)
	os.WriteFile(dirtyArtifact, []byte("should be gone"), 0644)

	// Create valid structure for Archive
	os.MkdirAll(store.ProtocolsPath(), 0755)
	os.WriteFile(store.ProtocolsPath(assets.DefaultProtocolName+".yaml"), []byte(assets.DefaultProtocolYAML), 0644)
	os.MkdirAll(store.TemplatesPath(), 0755)
	os.WriteFile(store.TemplatesPath("requirements.md"), []byte("# Req"), 0644)
	os.WriteFile(store.ConfigPath(), []byte("protocol: "+assets.DefaultProtocolName+"\n"), 0644)

	// Create Clean State
	cleanStateJSON := `{"completed_stages": []}`
	os.WriteFile(store.StatePath(), []byte(cleanStateJSON), 0644)

	// Create snapshot manager
	mgr := snapshot.NewManager(store.ArchivesPath())

	// Create Archive "clean-v1" (Snapshot of this state, but WITHOUT the dirty artifact in state, so createArchive won't include it?)
	// update: createArchive only includes artifacts referenced in state.json.
	// The `dirty.md` is NOT in state.json, so it won't be in the archive.
	// But it IS on disk.
	// So `createArchive` will make an archive that *doesn't* have `dirty.md`.

	err := mgr.Create("clean-v1", nil, "")
	if err != nil {
		t.Fatalf("createArchive failed: %v", err)
	}

	// Verify archive exists
	if _, err := os.Stat(store.ArchivesPath("clean-v1")); err != nil {
		t.Fatalf("archive not created")
	}

	// Now, create another file "new_dirty.md" just to be sure we are modifying workspace
	os.WriteFile(store.ArtifactsPath("requirements", "new_dirty.md"), []byte("garbage"), 0644)

	// Restore "clean-v1"
	// With the FIX, this should replace the entire `artifacts` directory with the one from the archive.
	// Since the archive has NO artifacts (or just empty dirs?), the restored `artifacts` dir should not contain `new_dirty.md`.
	// Wait, `createArchive` copies `store.ArtifactsPath()` (entire dir).
	// So `dirty.md` WAS included in the archive because `createArchive` copies the whole directory:
	// `copyDir(store.ArtifactsPath(), filepath.Join(tmpArchiveRoot, "artifacts"))`

	// Ah! `createArchive` DOES validate referenced artifacts, but the actual copy is `copyDir` of the root.
	// So `dirty.md` (untracked) IS in the archive.
	// This means my test setup is flawed. I need to make an archive that definitely DOESN'T have the file.

	// Correct Approach:
	// 1. Clean workspace (no dirty files).
	// 2. Create Archive "clean".
	// 3. Add dirty file to workspace.
	// 4. Restore "clean" -> Dirty file should be gone.

	// Reset
	os.RemoveAll(store.ArtifactsPath())
	// Create minimal valid artifact structure (maybe empty?)
	os.MkdirAll(store.ArtifactsPath(), 0755)

	// Create clean archive
	err = mgr.Create("clean-v2", nil, "")
	if err != nil {
		t.Fatalf("createArchive v2 failed: %v", err)
	}

	// Create dirty file
	os.MkdirAll(store.ArtifactsPath("requirements"), 0755)
	dirtyPath := store.ArtifactsPath("requirements", "dirty.md")
	os.WriteFile(dirtyPath, []byte("trash"), 0644)

	// Restore "clean-v2" --force
	cmd := &cobra.Command{}
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.Flags().Bool("force", false, "")
	cmd.Flags().Set("force", "true")

	err = archiveRestoreCmd.RunE(cmd, []string{"clean-v2"})
	if err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	// Assertions
	if _, err := os.Stat(dirtyPath); !os.IsNotExist(err) {
		t.Errorf("dirty file `dirty.md` still exists after restore! The workspace was not cleanly reset.")
	}
}
