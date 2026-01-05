package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"specfirst/internal/prompts"
)

var assumptionsCmd = &cobra.Command{
	Use:   "assumptions <spec-file>",
	Short: "Generate a prompt to extract implicit assumptions",
	Long: `Generate a prompt that forces the surfacing of hidden assumptions in a specification.

This command helps identify implicit assumptions that could lead to:
- Misunderstandings between stakeholders
- Incorrect implementations
- Untested edge cases
- Deployment failures

The output is a structured prompt suitable for AI assistants or human reviewers.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		specPath := args[0]

		content, err := os.ReadFile(specPath)
		if err != nil {
			return fmt.Errorf("reading spec %s: %w", specPath, err)
		}

		prompt, err := prompts.Render("assumptions-extraction.md", prompts.SpecData{
			Spec: string(content),
		})
		if err != nil {
			return fmt.Errorf("rendering assumptions prompt: %w", err)
		}

		prompt = applyMaxChars(prompt, stageMaxChars)
		formatted, err := formatPrompt(stageFormat, "assumptions", prompt)
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
