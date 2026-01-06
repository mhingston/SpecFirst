package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"specfirst/internal/protocol"
	"specfirst/internal/state"
	"specfirst/internal/store"
)

func TestOutputRelPath(t *testing.T) {
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

	// Initialize git repo so repoRoot() works as expected
	if _, err := gitCmd("init"); err != nil {
		t.Logf("Skipping git-dependent test parts: init failed: %v", err)
	}

	abs := filepath.Join(tmp, "src", "main.go")
	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	rel, err := outputRelPath("src/main.go")
	if err != nil {
		t.Fatalf("outputRelPath relative: %v", err)
	}
	if rel != filepath.Join("src", "main.go") {
		t.Fatalf("expected relative path, got %q", rel)
	}

	rel, err = outputRelPath(abs)
	if err != nil {
		t.Fatalf("outputRelPath absolute: %v", err)
	}
	if rel != filepath.Join("src", "main.go") {
		t.Fatalf("expected relative path from abs, got %q", rel)
	}

	if _, err := outputRelPath(filepath.Join("..", "outside.txt")); err == nil {
		t.Fatalf("expected error for escaping path")
	}

	// Test CWD-relative path resolution from a subdirectory
	subDir := filepath.Join(tmp, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}
	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("chdir subdir: %v", err)
	}
	rel, err = outputRelPath("foo.md")
	if err != nil {
		t.Fatalf("outputRelPath in subdir: %v", err)
	}
	// It should resolve foo.md in subdir to "subdir/foo.md" relative to repo root
	expected := filepath.Join("subdir", "foo.md")
	if rel != expected {
		t.Fatalf("expected %q, got %q", expected, rel)
	}
}

func TestValidateOutputs(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	// Initialize workspace root
	if err := os.Mkdir(filepath.Join(tmp, ".specfirst"), 0755); err != nil {
		t.Fatal(err)
	}

	stage := protocol.Stage{
		Outputs: []string{"src/*", "requirements.md"},
	}
	outputs := []string{"src/main.go", "requirements.md"}
	if err := validateOutputs(stage, outputs); err != nil {
		t.Fatalf("expected outputs to validate, got %v", err)
	}

	withBackslash := []string{"src\\main.go", "requirements.md"}
	if err := validateOutputs(stage, withBackslash); err != nil {
		t.Fatalf("expected outputs to validate with backslash, got %v", err)
	}

	missing := []string{"src/main.go"}
	if err := validateOutputs(stage, missing); err == nil {
		t.Fatalf("expected missing outputs error")
	}
}

func TestMissingApprovals(t *testing.T) {
	p := protocol.Protocol{
		Approvals: []protocol.Approval{
			{Stage: "design", Role: "reviewer"},
		},
	}
	s := state.State{
		CompletedStages: []string{"design"},
		Attestations:    map[string][]state.Attestation{},
	}

	missing := missingApprovals(p, s)
	if len(missing) != 1 || missing[0] != "design:reviewer" {
		t.Fatalf("expected missing approval, got %v", missing)
	}

	s.Attestations["design"] = []state.Attestation{{Role: "reviewer", Status: "approved"}}
	missing = missingApprovals(p, s)
	if len(missing) != 0 {
		t.Fatalf("expected no missing approvals, got %v", missing)
	}
}

func TestArtifactPathForInput(t *testing.T) {
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

	artifactDir := store.ArtifactsPath("design")
	if err := os.MkdirAll(artifactDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	artifactPath := store.ArtifactsPath("design", "notes.md")
	if err := os.WriteFile(artifactPath, []byte("ok"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	stageIDs := []string{"design"}
	got, err := artifactPathForInput(filepath.Join("design", "notes.md"), nil, stageIDs)
	if err != nil {
		t.Fatalf("artifactPathForInput: %v", err)
	}
	if got != artifactPath {
		t.Fatalf("expected %q, got %q", artifactPath, got)
	}

	if _, err := artifactPathForInput(filepath.Join("design", "..", "secrets.txt"), nil, stageIDs); err == nil {
		t.Fatalf("expected error for path traversal")
	}

	abs, err := filepath.Abs("abs.txt")
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	if _, err := artifactPathForInput(abs, nil, stageIDs); err == nil {
		t.Fatalf("expected error for absolute path")
	}

	nestedArtifact := store.ArtifactsPath("design", "nested", "notes.md")
	if err := os.MkdirAll(filepath.Dir(nestedArtifact), 0755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	if err := os.WriteFile(nestedArtifact, []byte("ok"), 0644); err != nil {
		t.Fatalf("write nested: %v", err)
	}
	got, err = artifactPathForInput(filepath.Join("nested", "notes.md"), []string{"design"}, stageIDs)
	if err != nil {
		t.Fatalf("artifactPathForInput nested: %v", err)
	}
	if got != nestedArtifact {
		t.Fatalf("expected nested %q, got %q", nestedArtifact, got)
	}
}

func TestCopyDirWithOpts_Symlinks(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	dst := filepath.Join(tmp, "dst")
	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a regular file
	regPath := filepath.Join(src, "regular.txt")
	if err := os.WriteFile(regPath, []byte("regular content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Case 1: Valid Relative Symlink (inside root)
	// src/valid_link.txt -> regular.txt
	validLink := filepath.Join(src, "valid_link.txt")
	err := os.Symlink("regular.txt", validLink)
	if err != nil {
		t.Logf("Skipping symlink test: %v", err)
		return
	}

	// Case 2: Insecure Absolute Symlink (to outside) - Should act as a blocker or be skipped?
	// Our copyDir logic errors on insecure symlinks.
	// We'll create a separate source dir for the failure case to avoid messing up the first copy.

	// Run CopyDir on valid setup
	if err := copyDirWithOpts(src, dst, true); err != nil {
		t.Fatal(err)
	}

	// Verify valid link copied
	dstLink := filepath.Join(dst, "valid_link.txt")
	info, err := os.Lstat(dstLink)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("destination valid_link is not a symlink")
	}

	// Case 3: Verify Insecure Link Failure
	srcBad := filepath.Join(tmp, "src_bad")
	os.Mkdir(srcBad, 0755)
	badLink := filepath.Join(srcBad, "bad_link.txt")
	// Points to /tmp/outside.txt (absolute)
	targetPath := filepath.Join(tmp, "outside.txt")
	os.WriteFile(targetPath, []byte("secrets"), 0644)
	os.Symlink(targetPath, badLink)

	dstBad := filepath.Join(tmp, "dst_bad")
	if err := copyDirWithOpts(srcBad, dstBad, true); err == nil {
		t.Error("expected error for absolute symlink, got nil")
	} else if !strings.Contains(err.Error(), "absolute links not allowed") {
		t.Errorf("expected absolute link error, got: %v", err)
	}
}
