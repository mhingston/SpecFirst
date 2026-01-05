package store

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	SpecDir      = ".specfirst"
	ArtifactsDir = "artifacts"
	GeneratedDir = "generated"
	ProtocolsDir = "protocols"
	TemplatesDir = "templates"
	ArchivesDir  = "archives"
	StateFile    = "state.json"
	ConfigFile   = "config.yaml"
)

func SpecPath(elem ...string) string {
	parts := append([]string{baseDir(), SpecDir}, elem...)
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

func StatePath() string {
	return SpecPath(StateFile)
}

func ConfigPath() string {
	return SpecPath(ConfigFile)
}

func baseDir() string {
	if root, err := gitRoot(); err == nil && root != "" {
		return root
	}
	if wd, err := os.Getwd(); err == nil && wd != "" {
		return wd
	}
	return "."
}

func gitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("empty git root")
	}
	return root, nil
}
