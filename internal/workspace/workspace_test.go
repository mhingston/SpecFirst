package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"specfirst/internal/store"
)

func TestFindProjectRoot(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "specfirst-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Structure:
	// tmpDir/
	//   .specfirst/
	//   subdir/
	//     deep/
	//   other/
	//     .git/
	//     child/

	specRootDir := filepath.Join(tmpDir, "spec_root")
	if err := os.MkdirAll(filepath.Join(specRootDir, ".specfirst"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(specRootDir, "subdir", "deep"), 0755); err != nil {
		t.Fatal(err)
	}

	gitRootDir := filepath.Join(tmpDir, "git_root")
	if err := os.MkdirAll(filepath.Join(gitRootDir, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(gitRootDir, "child"), 0755); err != nil {
		t.Fatal(err)
	}

	noRootDir := filepath.Join(tmpDir, "no_root")
	if err := os.MkdirAll(noRootDir, 0755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		startDir string
		want     string
		wantErr  bool
	}{
		{
			name:     "Root with .specfirst",
			startDir: specRootDir,
			want:     specRootDir,
			wantErr:  false,
		},
		{
			name:     "Subdirectory in .specfirst project",
			startDir: filepath.Join(specRootDir, "subdir"),
			want:     specRootDir,
			wantErr:  false,
		},
		{
			name:     "Deep subdirectory in .specfirst project",
			startDir: filepath.Join(specRootDir, "subdir", "deep"),
			want:     specRootDir,
			wantErr:  false,
		},
		{
			name:     "Root with .git",
			startDir: gitRootDir,
			want:     gitRootDir,
			wantErr:  false,
		},
		{
			name:     "Subdirectory in .git project",
			startDir: filepath.Join(gitRootDir, "child"),
			want:     gitRootDir,
			wantErr:  false,
		},
		{
			name:     "No root found",
			startDir: noRootDir,
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.FindProjectRoot(tt.startDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindProjectRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Resolve symlinks for comparison as /var/folders/ usually involves symlinks on Mac
			if got != "" {
				gotEval, _ := filepath.EvalSymlinks(got)
				wantEval, _ := filepath.EvalSymlinks(tt.want)
				if gotEval != wantEval {
					t.Errorf("FindProjectRoot() = %v, want %v", gotEval, wantEval)
				}
			}
		})
	}
}

func TestProjectRelPath(t *testing.T) {
	// We need to mock the current working directory or ensure our test runs inside a "project"
	// easiest way is to change the working directory for the test
	tmpDir, err := os.MkdirTemp("", "specfirst-rel-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup valid project root
	if err := os.Mkdir(filepath.Join(tmpDir, ".specfirst"), 0755); err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd) // Restore WD

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "Simple file in root",
			path:    "README.md",
			want:    "README.md",
			wantErr: false,
		},
		{
			name:    "File in subdirectory",
			path:    "src/main.go",
			want:    "src/main.go",
			wantErr: false,
		},
		{
			name:    "Absolute path inside root",
			path:    filepath.Join(tmpDir, "docs", "design.md"),
			want:    "docs/design.md",
			wantErr: false,
		},
		{
			name:    "Path escape",
			path:    "../outside.txt",
			wantErr: true,
		},
		{
			name:    "Root itself",
			path:    ".",
			wantErr: true, // "path resolves to project root" error
		},
		{
			name:    "Absolute path outside root",
			path:    "/tmp/somewhere/else",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProjectRelPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectRelPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ProjectRelPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
