package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"specfirst/internal/prompts"
)

var calibrateCmd = &cobra.Command{
	Use:   "calibrate <artifact>",
	Short: "Generate an epistemic annotation prompt for judgment calibration",
	Long: `Generate a prompt that helps calibrate confidence in a specification or artifact.

This command helps identify:
- What is known with high confidence
- What is assumed but not proven
- What is explicitly uncertain
- What would invalidate the spec

The output is a structured prompt suitable for AI assistants or human reviewers.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		artifactPath := args[0]

		content, err := os.ReadFile(artifactPath)
		if err != nil {
			return fmt.Errorf("reading artifact %s: %w", artifactPath, err)
		}

		prompt, err := prompts.Render("epistemic-calibration.md", prompts.SpecData{
			Spec: string(content),
		})
		if err != nil {
			return fmt.Errorf("rendering calibrate prompt: %w", err)
		}

		prompt = applyMaxChars(prompt, stageMaxChars)
		formatted, err := formatPrompt(stageFormat, "calibrate", prompt)
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

func init() {
	// No flags needed anymore as we consolidated on a single gold-standard prompt
}
