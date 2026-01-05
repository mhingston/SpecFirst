package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"specfirst/internal/prompts"
)

var failureCmd = &cobra.Command{
	Use:   "failure-modes <spec-file>",
	Short: "Generate a prompt focused on how the spec could fail",
	Long: `Generate a prompt that enumerates all plausible failure modes of a specification.

This command directly addresses happy-path bias by forcing consideration of:
- Partial failures
- Misuse cases
- Ambiguous interpretations
- Race conditions
- Operational edge cases

The output is a structured prompt suitable for AI assistants or human reviewers.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		specPath := args[0]

		content, err := os.ReadFile(specPath)
		if err != nil {
			return fmt.Errorf("reading spec %s: %w", specPath, err)
		}

		prompt, err := prompts.Render("failure-modes.md", prompts.SpecData{
			Spec: string(content),
		})
		if err != nil {
			return fmt.Errorf("rendering failure-modes prompt: %w", err)
		}

		prompt = applyMaxChars(prompt, stageMaxChars)
		formatted, err := formatPrompt(stageFormat, "failure-modes", prompt)
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
