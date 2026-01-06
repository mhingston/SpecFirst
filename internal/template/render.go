package template

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Input struct {
	Name    string
	Content string
}

type Data struct {
	StageName   string
	ProjectName string
	Inputs      []Input
	Outputs     []string
	Intent      string
	Language    string
	Framework   string
	CustomVars  map[string]string
	Constraints map[string]string

	// Protocol v2 fields
	StageType      string
	Prompt         any
	OutputContract any
	Epistemics     any
}

// Render renders a template file with the given data.
func Render(templatePath string, data Data) (string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("reading template %s: %w", templatePath, err)
	}

	// Use a custom delimiter to avoid conflicts with potential code in prompts?
	// Standard Go templates use {{ }} which is fine for most things.
	// SpecFirst might want strict control.
	// For now assuming standard.

	tmpl, err := template.New(templatePath).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// RenderInline renders an inline template string.
func RenderInline(tmpl string, data any) (string, error) {
	parsed, err := template.New("inline").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := parsed.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
