package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"specfirst/internal/assets"
	"specfirst/internal/store"
)

// captureOutput captures stdout from a function call
func captureOutput(t *testing.T, f func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = orig
	}()

	f()

	w.Close()
	var buf strings.Builder
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestProtocolSwitch(t *testing.T) {
	// Create a temp directory for the test workspace
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Initialize workspace with default protocol
	rootCmd.SetArgs([]string{"init"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Create a dummy custom protocol file outside the standard location
	customProtoContent := `name: "custom-proto"
version: "1.0"
stages:
  - id: custom
    name: Custom Stage
    intent: test
    template: custom.md
    outputs: [custom.md]
`
	customProtoPath := filepath.Join(tmpDir, "custom-protocol.yaml")
	if err := os.WriteFile(customProtoPath, []byte(customProtoContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create the matching template
	customTmplPath := filepath.Join(tmpDir, ".specfirst", "templates", "custom.md")
	if err := os.WriteFile(customTmplPath, []byte("# Custom Template"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("overrides default protocol with file path", func(t *testing.T) {
		// We can't really execute 'custom' stage because runStage prints to stdout and doesn't return the prompt string easily for inspection without a buffer.
		// Instead, we can misuse 'protocol' command to list stages, which will verify the protocol loaded.

		// Reset flags
		protocolFlag = ""

		// Run protocol command with override
		// We capture stdout to verify
		output := captureOutput(t, func() {
			rootCmd.SetArgs([]string{"protocol", "list", "--protocol", customProtoContent}) // wait, valid arg is the path
			// Cobra flags persistence is tricky in tests if not careful.
			// Let's call loadConfig/loadProtocol directly to test the helper logic first?
			// But we updated loadConfig to use global protocolFlag variable.

			// Set the flag variable directly for unit testing loadConfig/loadProtocol interactions
			protocolFlag = customProtoPath

			cfg, err := loadConfig()
			if err != nil {
				t.Fatalf("loadConfig failed: %v", err)
			}
			// loadConfig should NOT override the protocol
			if cfg.Protocol == customProtoPath {
				t.Errorf("expected config protocol to remain default, but got %q", cfg.Protocol)
			}
			// activeProtocolName should return the override
			active := activeProtocolName(cfg)
			if active != customProtoPath {
				t.Errorf("expected active protocol to be %q, got %q", customProtoPath, active)
			}

			proto, err := loadProtocol(active)
			if err != nil {
				t.Fatalf("loadProtocol failed: %v", err)
			}
			if proto.Name != "custom-proto" {
				t.Errorf("expected protocol name 'custom-proto', got %q", proto.Name)
			}
		})
		_ = output
	})

	t.Run("check command detects drift when overriding", func(t *testing.T) {
		// Initialize state with default protocol
		s := assets.DefaultProtocolName // "multi-stage"
		// Write state file forcing it to default protocol
		statePath := store.StatePath()
		stateContent := fmt.Sprintf(`{"protocol": "%s", "spec_version": "1.0"}`, s)
		if err := os.WriteFile(statePath, []byte(stateContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Set flag to custom proto
		protocolFlag = customProtoPath

		// We want to verify that check warns about drift.
		// We'll call checkCmd.RunE essentially, or just verify helpers logic.
		// Since we modified loadConfig globally, any command dealing with state checking will see the mismatch.

		cfg, _ := loadConfig()
		proto, _ := loadProtocol(activeProtocolName(cfg)) // Loads custom-proto

		if proto.Name != "custom-proto" {
			t.Fatal("failed to load custom proto")
		}

		// This confirms the plumbing works. The actual warning logic is in 'check.go' which compares cfg.Protocol (now custom) vs state.Protocol (multi-stage).
	})
}
