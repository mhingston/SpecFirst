package cmd

import (
	"fmt"
	"os"
	"specfirst/internal/prompts"

	"github.com/spf13/cobra"
)

var traceCmd = &cobra.Command{
	Use:   "trace <spec-file>",
	Short: "Generate a prompt for spec-to-code mapping",
	Long: `Generate a prompt that asks for mapping between specification sections and code areas.

This command helps identify:
- Which code modules implement which spec sections
- Missing implementation coverage
- Dead or obsolete code risks
- Refactoring impact areas

The output is a structured prompt suitable for AI assistants or human reviewers.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		specPath := args[0]

		content, err := os.ReadFile(specPath)
		if err != nil {
			return fmt.Errorf("reading spec %s: %w", specPath, err)
		}

		prompt, err := prompts.Render("trace.md", prompts.SpecData{
			Spec:   string(content),
			Source: specPath,
		})
		if err != nil {
			return fmt.Errorf("rendering prompt: %w", err)
		}

		prompt = applyMaxChars(prompt, stageMaxChars)
		formatted, err := formatPrompt(stageFormat, "trace", prompt)
		if err != nil {
			return err
		}

		if stageOut != "" {
			if err := writeOutput(stageOut, formatted); err != nil {
				return err
			}
		}
		_, err = cmd.OutOrStdout().Write([]byte(formatted))
		return err
	},
}
