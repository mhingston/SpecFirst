package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"specfirst/internal/prompts"
)

var diffCmd = &cobra.Command{
	Use:   "diff <old-spec> <new-spec>",
	Short: "Generate a prompt analyzing spec changes",
	Long: `Generate a prompt that describes the delta between two specification files.

This command helps evaluate behavioral differences, backward compatibility risks,
required code changes, and tests that must be updated when specifications change.

The output is a structured prompt suitable for AI assistants or human reviewers.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldPath := args[0]
		newPath := args[1]

		oldContent, err := os.ReadFile(oldPath)
		if err != nil {
			return fmt.Errorf("reading old spec %s: %w", oldPath, err)
		}

		newContent, err := os.ReadFile(newPath)
		if err != nil {
			return fmt.Errorf("reading new spec %s: %w", newPath, err)
		}

		prompt, err := prompts.Render("change-impact.md", prompts.DiffData{
			SpecBefore: string(oldContent),
			SpecAfter:  string(newContent),
		})
		if err != nil {
			return fmt.Errorf("rendering diff prompt: %w", err)
		}

		prompt = applyMaxChars(prompt, stageMaxChars)
		formatted, err := formatPrompt(stageFormat, "diff", prompt)
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
