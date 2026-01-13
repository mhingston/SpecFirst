package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"specfirst/internal/domain"
	"specfirst/internal/repository"
)

func TestCompilePromptMissingTemplateListsAvailable(t *testing.T) {
	tmp := t.TempDir()
	repository.SetRootDir(tmp)
	t.Cleanup(repository.ResetRootDir)

	templatesDir := filepath.Join(tmp, ".specfirst", "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("mkdir templates: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templatesDir, "exists.md"), []byte("# Exists"), 0644); err != nil {
		t.Fatalf("write template: %v", err)
	}

	app := NewApplication(domain.Config{ProjectName: "Test"}, domain.Protocol{}, domain.State{})
	stage := domain.Stage{ID: "plan", Name: "Plan", Template: "missing.md"}

	_, err := app.CompilePrompt(stage, nil, CompileOptions{})
	if err == nil {
		t.Fatalf("expected error for missing template")
	}

	msg := err.Error()
	if !strings.Contains(msg, "available: exists.md") {
		t.Fatalf("expected available templates in error, got: %s", msg)
	}
	if !strings.Contains(msg, "specfirst init") {
		t.Fatalf("expected hint in error, got: %s", msg)
	}
}

func TestCompilePromptMissingTemplatesDirHint(t *testing.T) {
	tmp := t.TempDir()
	repository.SetRootDir(tmp)
	t.Cleanup(repository.ResetRootDir)

	app := NewApplication(domain.Config{ProjectName: "Test"}, domain.Protocol{}, domain.State{})
	stage := domain.Stage{ID: "plan", Name: "Plan", Template: "missing.md"}

	_, err := app.CompilePrompt(stage, nil, CompileOptions{})
	if err == nil {
		t.Fatalf("expected error for missing templates directory")
	}

	msg := err.Error()
	if !strings.Contains(msg, "templates directory not found") {
		t.Fatalf("expected templates directory message, got: %s", msg)
	}
	if !strings.Contains(msg, "specfirst init") {
		t.Fatalf("expected hint in error, got: %s", msg)
	}
}
