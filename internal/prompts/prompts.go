// Package prompts provides embedded prompt templates for cognitive scaffold commands.
package prompts

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed *.md
var promptFS embed.FS

// PromptTemplates holds all parsed prompt templates.
var PromptTemplates *template.Template

func init() {
	var err error
	// Parse prompt-contract.md first as it's referenced by other templates
	PromptTemplates, err = template.New("prompts").ParseFS(promptFS, "prompt-contract.md")
	if err != nil {
		panic(fmt.Sprintf("failed to parse prompt-contract.md: %v", err))
	}
	// Parse all other templates
	PromptTemplates, err = PromptTemplates.ParseFS(promptFS, "*.md")
	if err != nil {
		panic(fmt.Sprintf("failed to parse prompt templates: %v", err))
	}
}

// SpecData is the data structure for single-spec prompts.
type SpecData struct {
	Spec   string
	Source string
}

// DiffData is the data structure for diff/comparison prompts.
type DiffData struct {
	SpecBefore string
	SpecAfter  string
}

// Render executes a named prompt template with the given data.
func Render(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := PromptTemplates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("rendering prompt %s: %w", name, err)
	}
	return buf.String(), nil
}
