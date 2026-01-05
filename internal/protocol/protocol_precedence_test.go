package protocol

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProtocolImportPrecedence(t *testing.T) {
	// Setup temporary directory structure
	tmpDir := t.TempDir()

	// 1. Create Base Protocol (defines stage 'check')
	baseProto := `
name: base
stages:
  - id: check
    intent: base check
    template: base/tpl
    depends_on: []
`
	err := os.WriteFile(filepath.Join(tmpDir, "base.yaml"), []byte(baseProto), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create Override Protocol (redefines stage 'check')
	overrideProto := `
name: override
stages:
  - id: check
    intent: override check
    template: override/tpl
    depends_on: []
`
	err = os.WriteFile(filepath.Join(tmpDir, "override.yaml"), []byte(overrideProto), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Create Main Protocol (uses [base, override])
	// Expected: valid override means 'check' stage comes from 'override' protocol
	mainProto := `
name: main
uses:
  - base
  - override
stages: []
approvals: []
`
	mainPath := filepath.Join(tmpDir, "main.yaml")
	err = os.WriteFile(mainPath, []byte(mainProto), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Load Protocol
	p, err := Load(mainPath)
	if err != nil {
		t.Fatalf("Failed to load protocol: %v", err)
	}

	// Assertions
	stage, found := p.StageByID("check")
	if !found {
		t.Fatal("Stage 'check' not found")
	}

	if stage.Intent != "override check" {
		t.Errorf("Expected intent 'override check', got '%s'.\nThis indicates Base overrode Override (Precedence Bug).", stage.Intent)
	}
}
