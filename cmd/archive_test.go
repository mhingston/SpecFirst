package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"specfirst/internal/assets"
	"specfirst/internal/snapshot"
	"specfirst/internal/store"
)

func TestCreateArchiveRequiresTemplates(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	if err := os.MkdirAll(store.ProtocolsPath(), 0755); err != nil {
		t.Fatalf("mkdir protocols: %v", err)
	}
	if err := os.WriteFile(store.ProtocolsPath(assets.DefaultProtocolName+".yaml"), []byte(assets.DefaultProtocolYAML), 0644); err != nil {
		t.Fatalf("write protocol: %v", err)
	}

	mgr := snapshot.NewManager(store.ArchivesPath())
	err = mgr.Create("1.0", nil, "")
	if err == nil {
		t.Fatalf("expected error when templates directory is missing")
	}
	if !strings.Contains(err.Error(), "required directory missing") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestArchiveRestoreRejectsProtocolMismatch(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	archiveRoot := store.ArchivesPath("mismatch-test")
	if err := os.MkdirAll(filepath.Join(archiveRoot, "protocols"), 0755); err != nil {
		t.Fatalf("mkdir protocols: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(archiveRoot, "templates"), 0755); err != nil {
		t.Fatalf("mkdir templates: %v", err)
	}

	protocol := `name: "proto-a"
version: "1.0"
stages:
  - id: requirements
    name: Requirements
    intent: exploration
    template: requirements.md
    outputs: [requirements.md]
`
	if err := os.WriteFile(filepath.Join(archiveRoot, "protocols", "proto-a.yaml"), []byte(protocol), 0644); err != nil {
		t.Fatalf("write protocol: %v", err)
	}
	if err := os.WriteFile(filepath.Join(archiveRoot, "templates", "requirements.md"), []byte("# Requirements\n"), 0644); err != nil {
		t.Fatalf("write template: %v", err)
	}

	config := `project_name: mismatch
protocol: proto-a
language: ""
framework: ""

custom_vars: {}
constraints: {}
`
	if err := os.WriteFile(filepath.Join(archiveRoot, "config.yaml"), []byte(config), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(archiveRoot, "state.json"), []byte("{}\n"), 0644); err != nil {
		t.Fatalf("write state: %v", err)
	}

	metadata := `{
  "version": "mismatch-test",
  "protocol": "proto-b",
  "archived_at": "2025-01-01T00:00:00Z",
  "stages_completed": []
}
`
	if err := os.WriteFile(filepath.Join(archiveRoot, "metadata.json"), []byte(metadata), 0644); err != nil {
		t.Fatalf("write metadata: %v", err)
	}

	cmd := &cobra.Command{}
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.Flags().Bool("force", false, "")
	if err := cmd.Flags().Set("force", "true"); err != nil {
		t.Fatalf("set force flag: %v", err)
	}

	err = archiveRestoreCmd.RunE(cmd, []string{"mismatch-test"})
	if err == nil || !strings.Contains(err.Error(), "archive metadata protocol mismatch") {
		t.Fatalf("expected metadata mismatch error, got %v", err)
	}
}

func TestCreateArchive_MissingArtifact(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Setup necessary directories and files
	if err := os.MkdirAll(store.ProtocolsPath(), 0755); err != nil {
		t.Fatalf("mkdir protocols: %v", err)
	}
	if err := os.WriteFile(store.ProtocolsPath(assets.DefaultProtocolName+".yaml"), []byte(assets.DefaultProtocolYAML), 0644); err != nil {
		t.Fatalf("write protocol: %v", err)
	}
	if err := os.MkdirAll(store.TemplatesPath(), 0755); err != nil {
		t.Fatalf("mkdir templates: %v", err)
	}
	if err := os.WriteFile(store.TemplatesPath("requirements.md"), []byte("# Requirements"), 0644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	if err := os.WriteFile(store.ConfigPath(), []byte("protocol: "+assets.DefaultProtocolName+"\n"), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Setup state with a missing artifact
	stateJSON := fmt.Sprintf(`{
  "protocol": %q,
  "spec_version": "1.0",
  "completed_stages": ["requirements"],
  "stage_outputs": {
    "requirements": {
      "files": ["requirements/foo.md"]
    }
  }
}
`, assets.DefaultProtocolName)
	if err := os.MkdirAll(filepath.Dir(store.StatePath()), 0755); err != nil {
		t.Fatalf("mkdir state dir: %v", err)
	}
	if err := os.WriteFile(store.StatePath(), []byte(stateJSON), 0644); err != nil {
		t.Fatalf("write state: %v", err)
	}

	mgr := snapshot.NewManager(store.ArchivesPath())

	// Try to create archive - should fail because requirements/foo.md is missing
	err = mgr.Create("1.0", nil, "")
	if err == nil {
		t.Fatalf("expected error for missing artifact")
	}
	if !strings.Contains(err.Error(), "missing artifact for stage requirements") {
		t.Fatalf("unexpected error: %v", err)
	}

	// Now create the artifact and it should succeed
	artifactPath := store.ArtifactsPath("requirements", "foo.md")
	if err := os.MkdirAll(filepath.Dir(artifactPath), 0755); err != nil {
		t.Fatalf("mkdir artifact dir: %v", err)
	}
	if err := os.WriteFile(artifactPath, []byte("content"), 0644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}

	err = mgr.Create("1.0", nil, "")
	if err != nil {
		t.Fatalf("unexpected error after creating artifact: %v", err)
	}
}
