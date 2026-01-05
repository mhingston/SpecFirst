package prompt

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	schema := DefaultSchema()

	t.Run("valid prompt", func(t *testing.T) {
		prompt := "## Context\nSome context\n## Task\nSome task\n## Assumptions\nNone"
		result := Validate(prompt, schema)
		if len(result.Warnings) > 0 {
			t.Errorf("expected no warnings, got %v", result.Warnings)
		}
	})

	t.Run("missing section", func(t *testing.T) {
		prompt := "## Context\nSome context"
		result := Validate(prompt, schema)
		found := false
		for _, w := range result.Warnings {
			if strings.Contains(w, "missing required section: Task") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected warning about missing Task section, got %v", result.Warnings)
		}
	})

	t.Run("forbidden phrase", func(t *testing.T) {
		prompt := "## Context\n## Task\nPlease make it better."
		result := Validate(prompt, schema)
		found := false
		for _, w := range result.Warnings {
			if strings.Contains(w, "contains ambiguous phrase") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected warning about forbidden phrase, got %v", result.Warnings)
		}
	})
}

func TestContainsAmbiguity(t *testing.T) {
	t.Run("vague language", func(t *testing.T) {
		prompt := "Maybe add some stuff as needed."
		issues := ContainsAmbiguity(prompt)
		if len(issues) == 0 {
			t.Errorf("expected ambiguity warnings, got none")
		}
	})
}

func TestExtractHeader(t *testing.T) {
	prompt := "---\nintent: test\n---\nBody content"
	header, body := ExtractHeader(prompt)
	if header != "intent: test" {
		t.Errorf("expected header 'intent: test', got %q", header)
	}
	if body != "Body content" {
		t.Errorf("expected body 'Body content', got %q", body)
	}
}
