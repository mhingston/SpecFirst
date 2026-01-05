package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// discoverChangedFiles returns a list of files that are either untracked or modified in git.
// It filters out files within .specfirst directory.
func discoverChangedFiles() ([]string, error) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("git not found: cannot auto-discover changed files")
	}
	root, err := gitRoot()
	if err != nil {
		return nil, fmt.Errorf("auto-discovery failed: %w", err)
	}

	// Get untracked files
	untracked, err := gitCmd("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, fmt.Errorf("failed to list untracked files: %w", err)
	}

	// Get modified files
	modified, err := gitCmd("diff", "--name-only")
	if err != nil {
		return nil, fmt.Errorf("failed to list modified files: %w", err)
	}

	// Get staged files
	staged, err := gitCmd("diff", "--name-only", "--cached")
	if err != nil {
		return nil, fmt.Errorf("failed to list staged files: %w", err)
	}

	all := append(untracked, modified...)
	all = append(all, staged...)
	unique := make(map[string]bool)
	filtered := make([]string, 0, len(all))

	for _, f := range all {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		// Ignore specfirst internal files
		if strings.HasPrefix(f, ".specfirst/") {
			continue
		}
		abs := filepath.Join(root, f)
		info, err := os.Stat(abs)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if info.IsDir() {
			continue
		}
		if !unique[abs] {
			unique[abs] = true
			filtered = append(filtered, abs)
		}
	}

	return filtered, nil
}

func gitCmd(args ...string) ([]string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	lines := strings.Split(out.String(), "\n")
	return lines, nil
}

func gitRoot() (string, error) {
	lines, err := gitCmd("rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("current directory is not a git repository (use 'git init')")
	}
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return "", fmt.Errorf("failed to determine git root")
	}
	return strings.TrimSpace(lines[0]), nil
}
