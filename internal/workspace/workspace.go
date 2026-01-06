package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"specfirst/internal/store"
)

// EnsureDir ensures a directory exists.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("opening source file %s: %w", src, err)
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return fmt.Errorf("stat source file %s: %w", src, err)
	}
	mode := info.Mode()

	dstDir := filepath.Dir(dst)
	if err := EnsureDir(dstDir); err != nil {
		return fmt.Errorf("creating destination directory %s: %w", dstDir, err)
	}

	// Use atomic write pattern: temp file + rename
	tmp, err := os.CreateTemp(dstDir, ".copyfile.*.tmp")
	if err != nil {
		return fmt.Errorf("creating temp file in %s: %w", dstDir, err)
	}
	tmpPath := tmp.Name()

	// Clean up temp file on any error
	success := false
	defer func() {
		if !success {
			_ = os.Remove(tmpPath)
		}
	}()

	// Set the correct mode on the temp file
	if err := tmp.Chmod(mode); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("setting file mode: %w", err)
	}

	if _, err := io.Copy(tmp, in); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("copying content: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("syncing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	// Atomic rename
	if runtime.GOOS == "windows" {
		if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("removing existing destination file for atomic rename: %w", err)
		}
	}
	if err := os.Rename(tmpPath, dst); err != nil {
		return fmt.Errorf("renaming temp file to %s: %w", dst, err)
	}
	success = true
	return nil
}

// CopyDir copies a directory recursively.
func CopyDir(src, dst string) error {
	return CopyDirWithOpts(src, dst, false)
}

// CopyDirWithOpts copies a directory with options.
func CopyDirWithOpts(src, dst string, required bool) error {
	info, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			if required {
				return fmt.Errorf("required directory missing: %s", src)
			}
			fmt.Fprintf(os.Stderr, "Warning: source directory does not exist, skipping copy: %s\n", src)
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}
	if err := EnsureDir(dst); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return EnsureDir(target)
		}
		if entry.Type()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("reading symlink %s: %w", path, err)
			}
			// Security check: ensure symlink doesn't point/traverse outside the target root
			// We only allow relative links that stay within the tree.
			if filepath.IsAbs(linkTarget) {
				return fmt.Errorf("insecure symlink %s -> %s: absolute links not allowed in archives", path, linkTarget)
			}
			// Check for traversal
			if strings.HasPrefix(linkTarget, "..") || strings.Contains(linkTarget, "/../") || strings.Contains(linkTarget, "\\..\\") {
				return fmt.Errorf("insecure symlink %s -> %s: directory traversal not allowed", path, linkTarget)
			}

			// Replicate the symlink
			if err := os.Symlink(linkTarget, target); err != nil {
				return fmt.Errorf("creating symlink %s -> %s: %w", target, linkTarget, err)
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			// Skip other non-regular files to prevent data leaks and ensure portability
			fmt.Fprintf(os.Stderr, "Warning: skipping non-regular file (socket/pipe): %s\n", path)
			return nil
		}
		return CopyFile(path, target)
	})
}

// ArtifactPathForInput resolves an input artifact path.
func ArtifactPathForInput(filename string, priorityStages []string, stageIDs []string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("invalid input artifact path: %s", filename)
	}

	clean := filepath.Clean(filename)
	if filepath.IsAbs(clean) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid input artifact path: %s", filename)
	}
	for _, part := range strings.Split(filepath.ToSlash(clean), "/") {
		if part == ".." {
			return "", fmt.Errorf("invalid input artifact path: %s", filename)
		}
	}

	// Support stage-qualified paths like "stage-id/filename" when stage-id matches a known stage.
	stageQualified := false
	if strings.Contains(clean, "/") || strings.Contains(clean, string(os.PathSeparator)) {
		first := strings.SplitN(filepath.ToSlash(clean), "/", 2)[0]
		for _, id := range stageIDs {
			if id == first {
				stageQualified = true
				break
			}
		}
	}
	if stageQualified {
		path := store.ArtifactsPath(clean)
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				return "", fmt.Errorf("missing input artifact: %s", filename)
			}
			return "", err
		}
		return path, nil
	}

	// Priority Check Level 1: Check in the priority stages first (explicit dependencies)
	for _, stageID := range priorityStages {
		path := store.ArtifactsPath(stageID, clean)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Glob-based resolution for unqualified filenames (fallback)
	pattern := store.ArtifactsPath("*", clean)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("missing input artifact: %s", filename)
	}
	if len(matches) > 1 {
		// Extract stage names from paths for helpful suggestion
		stages := make([]string, 0, len(matches))
		for _, match := range matches {
			rel, _ := filepath.Rel(store.ArtifactsPath(), match)
			if parts := strings.SplitN(rel, string(os.PathSeparator), 2); len(parts) > 0 {
				stages = append(stages, parts[0])
			}
		}
		sort.Strings(stages) // Alphabetical for consistent, deterministic output
		return "", fmt.Errorf("ambiguous input artifact %q found in multiple stages: %v\nHint: use a stage-qualified path like %q", filename, stages, stages[0]+"/"+clean)
	}
	return matches[0], nil
}

func ArtifactRelFromState(value string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("invalid artifact path: %q", value)
	}
	normalized := strings.ReplaceAll(value, "\\", "/")
	clean := filepath.Clean(filepath.FromSlash(normalized))
	if filepath.IsAbs(clean) || isWindowsAbs(normalized) {
		rel, ok := relFromArtifactsPath(clean)
		if !ok {
			return "", fmt.Errorf("artifact path is outside artifacts dir: %s", value)
		}
		return filepath.ToSlash(rel), nil
	}
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid artifact path: %s", value)
	}
	return filepath.ToSlash(clean), nil
}

func ArtifactAbsFromState(value string) (string, error) {
	rel, err := ArtifactRelFromState(value)
	if err != nil {
		return "", err
	}
	return filepath.Join(store.ArtifactsPath(), filepath.FromSlash(rel)), nil
}

func isWindowsAbs(value string) bool {
	if len(value) >= 3 && value[1] == ':' && value[2] == '/' {
		return true
	}
	return strings.HasPrefix(value, "//")
}

func relFromArtifactsPath(abs string) (string, bool) {
	clean := filepath.Clean(abs)
	parts := strings.Split(clean, string(os.PathSeparator))
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == store.ArtifactsDir {
			if i+1 >= len(parts) {
				return "", false
			}
			rel := filepath.Join(parts[i+1:]...)
			if rel == "" || rel == "." {
				return "", false
			}
			return rel, true
		}
	}
	return "", false
}

// NormalizeMatchPath normalizes a path for pattern matching.
func NormalizeMatchPath(value string) string {
	normalized := filepath.ToSlash(value)
	normalized = strings.ReplaceAll(normalized, "\\", "/")
	return strings.TrimPrefix(normalized, "./")
}

// MatchOutputPattern checks if a file matches an output pattern.
func MatchOutputPattern(pattern string, file string) bool {
	if pattern == "" || file == "" {
		return false
	}
	cleanPattern := NormalizeMatchPath(pattern)
	cleanFile := NormalizeMatchPath(file)
	if strings.Contains(cleanPattern, "*") {
		// Try matching against relative path and basename
		candidates := []string{cleanFile, path.Base(cleanFile)}
		for _, candidate := range candidates {
			if ok, _ := path.Match(cleanPattern, candidate); ok {
				return true
			}
		}
		return false
	}
	return cleanFile == cleanPattern
}

// PromptHash returns the SHA256 hash of a prompt string.
func PromptHash(prompt string) string {
	hash := sha256.Sum256([]byte(prompt))
	return "sha256:" + hex.EncodeToString(hash[:])
}

// FileHash returns the SHA256 hash of a file's content.
func FileHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(hash[:]), nil
}

// WriteOutput writes content to a file, ensuring the directory exists.
func WriteOutput(path string, content string) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// CollectFileHashes returns a map of relative paths to file hashes for a directory.
func CollectFileHashes(root string) (map[string]string, error) {
	files := map[string]string{}
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", root)
	}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		hash, err := FileHash(path)
		if err != nil {
			return err
		}
		files[rel] = hash
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// ProjectRelPath returns the path relative to the project root.
func ProjectRelPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	clean := filepath.Clean(path)

	// Determine project root
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	root, err := store.FindProjectRoot(wd)
	if err != nil {
		return "", err
	}

	// Resolve absolute path
	abs := clean
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(wd, clean)
	}

	// Evaluate symlinks for canonical path comparison
	// We must resolve the root to its canonical form
	if rootEval, err := filepath.EvalSymlinks(root); err == nil {
		root = rootEval
	}

	// We must also resolve abs to its canonical form as much as possible
	// Walk up until we find a directory that exists
	base := abs
	suffix := ""
	for {
		if _, err := os.Stat(base); err == nil {
			// Found existing path
			break
		}
		parent := filepath.Dir(base)
		if parent == base {
			break
		}
		suffix = filepath.Join(filepath.Base(base), suffix)
		base = parent
	}

	if baseEval, err := filepath.EvalSymlinks(base); err == nil {
		if suffix != "" {
			abs = filepath.Join(baseEval, suffix)
		} else {
			abs = baseEval
		}
	}

	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return "", err
	}

	// Safety checks
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes project root: %s", path)
	}
	if rel == "." {
		return "", fmt.Errorf("path resolves to project root")
	}

	return rel, nil
}
