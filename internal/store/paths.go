package store

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	SpecDir      = ".specfirst"
	ArtifactsDir = "artifacts"
	GeneratedDir = "generated"
	ProtocolsDir = "protocols"
	TemplatesDir = "templates"
	ArchivesDir  = "archives"
	TracksDir    = "tracks"
	SkillsDir    = "skills"
	StateFile    = "state.json"
	ConfigFile   = "config.yaml"
)

func SpecPath(elem ...string) string {
	parts := append([]string{BaseDir(), SpecDir}, elem...)
	return filepath.Join(parts...)
}

func ArtifactsPath(elem ...string) string {
	parts := append([]string{ArtifactsDir}, elem...)
	return SpecPath(parts...)
}

func GeneratedPath(elem ...string) string {
	parts := append([]string{GeneratedDir}, elem...)
	return SpecPath(parts...)
}

func ProtocolsPath(elem ...string) string {
	parts := append([]string{ProtocolsDir}, elem...)
	return SpecPath(parts...)
}

func TemplatesPath(elem ...string) string {
	parts := append([]string{TemplatesDir}, elem...)
	return SpecPath(parts...)
}

func ArchivesPath(elem ...string) string {
	parts := append([]string{ArchivesDir}, elem...)
	return SpecPath(parts...)
}

func TracksPath(elem ...string) string {
	parts := append([]string{TracksDir}, elem...)
	return SpecPath(parts...)
}

func SkillsPath(elem ...string) string {
	parts := append([]string{SkillsDir}, elem...)
	return SpecPath(parts...)
}

func StatePath() string {
	return SpecPath(StateFile)
}

func ConfigPath() string {
	return SpecPath(ConfigFile)
}

func BaseDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	if root, err := FindProjectRoot(wd); err == nil {
		return root
	}
	return wd
}

// FindProjectRoot looks for .specfirst or .git directory walking up from startDir.
func FindProjectRoot(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		// Check for .specfirst
		if _, err := os.Stat(filepath.Join(dir, SpecDir)); err == nil {
			return dir, nil
		}
		// Check for .git
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found (no .specfirst or .git in parents of %s)", startDir)
		}
		dir = parent
	}
}
