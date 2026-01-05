package protocol

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeProtocolFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "protocol.yaml")
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatalf("write protocol: %v", err)
	}
	return path
}

func TestLoadRejectsInvalidStageID(t *testing.T) {
	path := writeProtocolFile(t, `name: "test"
version: "1.0"
stages:
  - id: "bad/name"
    name: Bad
    intent: test
    template: requirements.md
    outputs: []
`)
	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "invalid stage id") {
		t.Fatalf("expected invalid stage id error, got %v", err)
	}
}

func TestLoadRejectsInvalidTemplatePath(t *testing.T) {
	path := writeProtocolFile(t, `name: "test"
version: "1.0"
stages:
  - id: "build"
    name: Build
    intent: test
    template: "../templates/requirements.md"
    outputs: []
`)
	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "invalid template path") {
		t.Fatalf("expected invalid template path error, got %v", err)
	}
}

func TestLoadRejectsSelfReference(t *testing.T) {
	path := writeProtocolFile(t, `name: "test"
version: "1.0"
stages:
  - id: "build"
    name: Build
    intent: test
    template: requirements.md
    depends_on: [build]
    outputs: []
`)
	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "cannot depend on itself") {
		t.Fatalf("expected self-reference error, got %v", err)
	}
}

func TestLoadRejectsEmptyApprovalRole(t *testing.T) {
	path := writeProtocolFile(t, `name: "test"
version: "1.0"
stages:
  - id: "build"
    name: Build
    intent: test
    template: requirements.md
    outputs: []
approvals:
  - stage: "build"
    role: ""
`)
	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "approval role is required") {
		t.Fatalf("expected approval role error, got %v", err)
	}
}

func TestLoadWithMixins(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specfirst-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a base protocol
	baseYAML := `name: base
version: 1.0
stages:
  - id: s1
    name: Stage 1
    template: s1.md
`
	if err := os.WriteFile(filepath.Join(tmpDir, "base.yaml"), []byte(baseYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a protocol that uses the base
	childYAML := `name: child
version: 1.0
uses:
  - base
stages:
  - id: s2
    name: Stage 2
    template: s2.md
    depends_on: [s1]
`
	childPath := filepath.Join(tmpDir, "child.yaml")
	if err := os.WriteFile(childPath, []byte(childYAML), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := Load(childPath)
	if err != nil {
		t.Fatalf("failed to load child protocol: %v", err)
	}

	if len(p.Stages) != 2 {
		t.Errorf("expected 2 stages, got %d", len(p.Stages))
	}

	// Verify order (imported stages should be first)
	if p.Stages[0].ID != "s1" {
		t.Errorf("expected first stage to be s1, got %s", p.Stages[0].ID)
	}
	if p.Stages[1].ID != "s2" {
		t.Errorf("expected second stage to be s2, got %s", p.Stages[1].ID)
	}
}

func TestValidateStageTypes(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid types",
			yaml: `name: test
stages:
  - {id: s1, type: spec, template: t.md}
  - {id: s2, type: decompose, template: t.md}
`,
			wantErr: false,
		},
		{
			name: "invalid type",
			yaml: `name: test
stages:
  - {id: s1, type: unknown, template: t.md}
`,
			wantErr: true,
		},
		{
			name: "task_prompt with valid source",
			yaml: `name: test
stages:
  - {id: s1, type: decompose, template: t.md}
  - {id: s2, type: task_prompt, source: s1, template: t.md}
`,
			wantErr: false,
		},
		{
			name: "task_prompt with invalid source type",
			yaml: `name: test
stages:
  - {id: s1, type: spec, template: t.md}
  - {id: s2, type: task_prompt, source: s1, template: t.md}
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := filepath.Join(t.TempDir(), "test.yaml")
			os.WriteFile(tmpFile, []byte(tt.yaml), 0644)

			_, err := Load(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadRejectsCircularImports(t *testing.T) {
	tmpDir := t.TempDir()

	// p1 uses p2, p2 uses p1
	p1YAML := `name: p1
uses: [p2]
stages: [{id: s1, template: t.md}]
`
	p2YAML := `name: p2
uses: [p1]
stages: [{id: s2, template: t.md}]
`
	if err := os.WriteFile(filepath.Join(tmpDir, "p1.yaml"), []byte(p1YAML), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "p2.yaml"), []byte(p2YAML), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(filepath.Join(tmpDir, "p1.yaml"))
	if err == nil || !strings.Contains(err.Error(), "circular protocol import detected") {
		t.Fatalf("expected circular import error, got %v", err)
	}
}
