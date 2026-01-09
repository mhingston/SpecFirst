package cmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"specfirst/internal/assets"
	"specfirst/internal/domain"
	"specfirst/internal/repository"

	"github.com/spf13/cobra"
)

// testRestoreCreateParams creates a CreateParams for testing.
func testRestoreCreateParams(t *testing.T) repository.CreateParams {
	t.Helper()
	cfg, err := repository.LoadConfig(repository.ConfigPath())
	if err != nil {
		cfg = domain.Config{Protocol: assets.DefaultProtocolName}
	}
	if cfg.Protocol == "" {
		cfg.Protocol = assets.DefaultProtocolName
	}
	proto, err := repository.LoadProtocol(repository.ProtocolsPath(cfg.Protocol + ".yaml"))
	if err != nil {
		t.Fatalf("load protocol: %v", err)
	}
	s, err := repository.LoadState(repository.StatePath())
	if err != nil {
		s = domain.NewState(proto.Name)
	}
	return repository.CreateParams{
		Config:   cfg,
		Protocol: proto,
		State:    s,
	}
}

func TestArchiveRestoreCleansWorkspace(t *testing.T) {
	// Setup isolated workspace using dependency injection (no os.Chdir)
	tmp := t.TempDir()

	// Inject the temp dir as the project root for this test execution
	repository.SetRootDir(tmp)
	t.Cleanup(func() { repository.ResetRootDir() })

	// Create a "dirty" artifact that should be removed by restore
	dirtyArtifact := repository.ArtifactsPath("requirements", "dirty.md")
	os.MkdirAll(filepath.Dir(dirtyArtifact), 0755)
	os.WriteFile(dirtyArtifact, []byte("should be gone"), 0644)

	// Create valid structure for Archive
	os.MkdirAll(repository.ProtocolsPath(), 0755)
	os.WriteFile(repository.ProtocolsPath(assets.DefaultProtocolName+".yaml"), []byte(assets.DefaultProtocolYAML), 0644)
	os.MkdirAll(repository.TemplatesPath(), 0755)
	os.WriteFile(repository.TemplatesPath("requirements.md"), []byte("# Req"), 0644)
	os.WriteFile(repository.ConfigPath(), []byte("protocol: "+assets.DefaultProtocolName+"\n"), 0644)

	// Create Clean State
	cleanStateJSON := `{"completed_stages": []}`
	os.WriteFile(repository.StatePath(), []byte(cleanStateJSON), 0644)

	// Create snapshot manager
	mgr := repository.NewSnapshotRepository(repository.ArchivesPath())
	params := testRestoreCreateParams(t)

	// Create Archive "clean-v1"
	err := mgr.Create("clean-v1", nil, "", params)
	if err != nil {
		t.Fatalf("createArchive failed: %v", err)
	}

	// Verify archive exists
	if _, err := os.Stat(repository.ArchivesPath("clean-v1")); err != nil {
		t.Fatalf("archive not created")
	}

	// Now, create another file "new_dirty.md" just to be sure we are modifying workspace
	os.WriteFile(repository.ArtifactsPath("requirements", "new_dirty.md"), []byte("garbage"), 0644)

	// Reset for clean-v2
	os.RemoveAll(repository.ArtifactsPath())
	os.MkdirAll(repository.ArtifactsPath(), 0755)

	// Create clean archive
	params = testRestoreCreateParams(t)
	err = mgr.Create("clean-v2", nil, "", params)
	if err != nil {
		t.Fatalf("createArchive v2 failed: %v", err)
	}

	// Create dirty file
	os.MkdirAll(repository.ArtifactsPath("requirements"), 0755)
	dirtyPath := repository.ArtifactsPath("requirements", "dirty.md")
	os.WriteFile(dirtyPath, []byte("trash"), 0644)

	// Restore "clean-v2" --force
	cmd := &cobra.Command{}
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.Flags().Bool("force", false, "")
	cmd.Flags().Set("force", "true")

	// We can't call archiveRestoreCmd.RunE directly because it depends on args/flags parsing that cobra does.
	// But `archiveRestoreCmd` logic invokes `repository.Restore`.
	// Ideally we unit test `repository.Restore` directly in repository package tests.
	// BUT `archive_restore_test` is an integration test for the cmd.
	// The problem is `archiveRestoreCmd` in `archive.go` does:
	// 	mgr := repository.NewSnapshotRepository(repository.ArchivesPath())
	// 	if err := mgr.Restore(version, force); err != nil { ... }
	// So we can re-implement the call locally to test the logic if we want, OR invoke the command.

	// The original test invoked `archiveRestoreCmd.RunE`.
	// We need to make sure `archive.go` has `archiveRestoreCmd` exported or accessible. It is local var?
	// `var archiveRestoreCmd = ...` in `archive.go`.
	// Yes, `archiveRestoreCmd` is package level variable in `cmd`.
	// So we can use it.

	err = archiveRestoreCmd.RunE(cmd, []string{"clean-v2"})
	if err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	// Assertions
	if _, err := os.Stat(dirtyPath); !os.IsNotExist(err) {
		t.Errorf("dirty file `dirty.md` still exists after restore! The workspace was not cleanly reset.")
	}
}
