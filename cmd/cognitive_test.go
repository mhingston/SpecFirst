package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"specfirst/internal/prompts"
)

func TestDiffCommand(t *testing.T) {
	t.Run("generates diff prompt", func(t *testing.T) {
		oldFile := filepath.Join(t.TempDir(), "old.md")
		newFile := filepath.Join(t.TempDir(), "new.md")

		if err := os.WriteFile(oldFile, []byte("# Old\n\n- Feature A"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(newFile, []byte("# New\n\n- Feature B"), 0644); err != nil {
			t.Fatal(err)
		}

		prompt, err := prompts.Render("change-impact.md", prompts.DiffData{
			SpecBefore: "# Old\n\n- Feature A",
			SpecAfter:  "# New\n\n- Feature B",
		})
		if err != nil {
			t.Fatal(err)
		}

		// Check for prompt contract inclusion
		if !strings.Contains(prompt, "PROMPT CONTRACT") {
			t.Error("expected prompt to contain 'PROMPT CONTRACT'")
		}
		// Check for change impact structure
		if !strings.Contains(prompt, "Change Impact") {
			t.Error("expected prompt to contain 'Change Impact'")
		}
		if !strings.Contains(prompt, "OUTPUT SCHEMA") {
			t.Error("expected prompt to contain 'OUTPUT SCHEMA'")
		}
	})

	t.Run("command requires two args", func(t *testing.T) {
		var buf bytes.Buffer
		diffCmd.SetOut(&buf)
		diffCmd.SetErr(&buf)
		diffCmd.SetArgs([]string{})

		err := diffCmd.Execute()
		if err == nil {
			t.Error("expected error for missing args")
		}
	})
}

func TestAssumptionsCommand(t *testing.T) {
	t.Run("generates assumptions prompt", func(t *testing.T) {
		prompt, err := prompts.Render("assumptions-extraction.md", prompts.SpecData{
			Spec: "# Spec\n\n- Feature A",
		})
		if err != nil {
			t.Fatal(err)
		}

		// Check for prompt contract inclusion
		if !strings.Contains(prompt, "PROMPT CONTRACT") {
			t.Error("expected prompt to contain 'PROMPT CONTRACT'")
		}
		// Check for assumptions extraction structure
		if !strings.Contains(prompt, "Assumptions Extraction") {
			t.Error("expected prompt to contain 'Assumptions Extraction'")
		}
		if !strings.Contains(prompt, "OUTPUT SCHEMA") {
			t.Error("expected prompt to contain 'OUTPUT SCHEMA'")
		}
		if !strings.Contains(prompt, "impact_if_false") {
			t.Error("expected prompt to contain 'impact_if_false'")
		}
	})
}

func TestReviewCommand(t *testing.T) {
	t.Run("generates security persona prompt", func(t *testing.T) {
		prompt, err := generateReviewPrompt("# Spec", "test.md", "security")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(prompt, "security reviewer") {
			t.Error("expected prompt to contain 'security reviewer'")
		}
		if !strings.Contains(prompt, "Attack Surfaces") {
			t.Error("expected prompt to contain 'Attack Surfaces'")
		}
	})

	t.Run("generates performance persona prompt", func(t *testing.T) {
		prompt, err := generateReviewPrompt("# Spec", "test.md", "performance")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(prompt, "performance reviewer") {
			t.Error("expected prompt to contain 'performance reviewer'")
		}
		if !strings.Contains(prompt, "Bottlenecks") {
			t.Error("expected prompt to contain 'Bottlenecks'")
		}
	})

	t.Run("generates maintainer persona prompt", func(t *testing.T) {
		prompt, err := generateReviewPrompt("# Spec", "test.md", "maintainer")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(prompt, "maintainability reviewer") {
			t.Error("expected prompt to contain 'maintainability reviewer'")
		}
	})

	t.Run("rejects unknown persona", func(t *testing.T) {
		_, err := generateReviewPrompt("# Spec", "test.md", "unknown")
		if err == nil {
			t.Error("expected error for unknown persona")
		}
		if !strings.Contains(err.Error(), "unknown persona") {
			t.Errorf("expected 'unknown persona' in error, got: %v", err)
		}
	})
}

func TestFailureCommand(t *testing.T) {
	t.Run("generates failure modes prompt", func(t *testing.T) {
		prompt, err := prompts.Render("failure-modes.md", prompts.SpecData{
			Spec: "# Spec",
		})
		if err != nil {
			t.Fatal(err)
		}

		// Check for prompt contract inclusion
		if !strings.Contains(prompt, "PROMPT CONTRACT") {
			t.Error("expected prompt to contain 'PROMPT CONTRACT'")
		}
		// Check for failure modes structure
		if !strings.Contains(prompt, "Failure Modes") {
			t.Error("expected prompt to contain 'Failure Modes'")
		}
		if !strings.Contains(prompt, "OUTPUT SCHEMA") {
			t.Error("expected prompt to contain 'OUTPUT SCHEMA'")
		}
		if !strings.Contains(prompt, "misuse_scenarios") {
			t.Error("expected prompt to contain 'misuse_scenarios'")
		}
	})
}

func TestTestIntentCommand(t *testing.T) {
	t.Run("generates test intent prompt", func(t *testing.T) {
		prompt := generateTestIntentPrompt("# Spec", "test.md")

		if !strings.Contains(prompt, "test intent") {
			t.Error("expected prompt to contain 'test intent'")
		}
		if !strings.Contains(prompt, "Do NOT generate test code") {
			t.Error("expected prompt to contain 'Do NOT generate test code'")
		}
		if !strings.Contains(prompt, "Required Invariants") {
			t.Error("expected prompt to contain 'Required Invariants'")
		}
		if !strings.Contains(prompt, "Boundary Conditions") {
			t.Error("expected prompt to contain 'Boundary Conditions'")
		}
		if !strings.Contains(prompt, "Negative Cases") {
			t.Error("expected prompt to contain 'Negative Cases'")
		}
	})
}

func TestTraceCommand(t *testing.T) {
	t.Run("generates trace prompt", func(t *testing.T) {
		prompt := generateTracePrompt("# Spec", "test.md")

		if !strings.Contains(prompt, "traceability mapping") {
			t.Error("expected prompt to contain 'traceability mapping'")
		}
		if !strings.Contains(prompt, "Identify Code Modules Affected") {
			t.Error("expected prompt to contain 'Identify Code Modules Affected'")
		}
		if !strings.Contains(prompt, "Dead or Obsolete Code Risks") {
			t.Error("expected prompt to contain 'Dead or Obsolete Code Risks'")
		}
	})
}

func TestDistillCommand(t *testing.T) {
	t.Run("generates exec audience prompt", func(t *testing.T) {
		prompt, err := generateDistillPrompt("# Spec", "test.md", "exec")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(prompt, "executive audience") {
			t.Error("expected prompt to contain 'executive audience'")
		}
		if !strings.Contains(prompt, "Business Impact") {
			t.Error("expected prompt to contain 'Business Impact'")
		}
	})

	t.Run("generates implementer audience prompt", func(t *testing.T) {
		prompt, err := generateDistillPrompt("# Spec", "test.md", "implementer")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(prompt, "implementation-focused developer") {
			t.Error("expected prompt to contain 'implementation-focused developer'")
		}
	})

	t.Run("generates ai audience prompt", func(t *testing.T) {
		prompt, err := generateDistillPrompt("# Spec", "test.md", "ai")
		if err != nil {
			t.Fatal(err)
		}

		// AI audience now uses gold-standard template
		if !strings.Contains(prompt, "PROMPT CONTRACT") {
			t.Error("expected prompt to contain 'PROMPT CONTRACT'")
		}
		if !strings.Contains(prompt, "AI-Facing Distillation") {
			t.Error("expected prompt to contain 'AI-Facing Distillation'")
		}
		if !strings.Contains(prompt, "OUTPUT SCHEMA") {
			t.Error("expected prompt to contain 'OUTPUT SCHEMA'")
		}
	})

	t.Run("generates qa audience prompt", func(t *testing.T) {
		prompt, err := generateDistillPrompt("# Spec", "test.md", "qa")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(prompt, "QA engineer") {
			t.Error("expected prompt to contain 'QA engineer'")
		}
	})

	t.Run("rejects unknown audience", func(t *testing.T) {
		_, err := generateDistillPrompt("# Spec", "test.md", "unknown")
		if err == nil {
			t.Error("expected error for unknown audience")
		}
		if !strings.Contains(err.Error(), "unknown audience") {
			t.Errorf("expected 'unknown audience' in error, got: %v", err)
		}
	})
}

func TestCalibrateCommand(t *testing.T) {
	t.Run("generates calibration prompt", func(t *testing.T) {
		artifactFile := filepath.Join(t.TempDir(), "artifact.md")
		if err := os.WriteFile(artifactFile, []byte("# Artifact\n\nContent"), 0644); err != nil {
			t.Fatal(err)
		}

		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)
		rootCmd.SetArgs([]string{"calibrate", artifactFile})

		if err := rootCmd.Execute(); err != nil {
			// It might return error if command fails, but we want to check output too
			// t.Logf("Execute error: %v", err)
		}

		output := buf.String()

		// Check for key sections from the gold-standard template
		expected := []string{
			"Epistemic Calibration",
			"PROMPT CONTRACT",
			"epistemic_map",
			"known:",
			"assumed:",
			"uncertain:",
			"unknown:",
		}

		for _, exp := range expected {
			if !strings.Contains(output, exp) {
				t.Errorf("expected output to contain %q", exp)
			}
		}
	})
}
