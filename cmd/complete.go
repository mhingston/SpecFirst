package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"specfirst/internal/app"
	"specfirst/internal/repository"
)

var completeForce bool

var completeCmd = &cobra.Command{
	Use:               "complete <stage-id> <output-files...>",
	Short:             "Mark a stage as complete and store outputs",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: stageIDThenFilesCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		stageID := args[0]
		var outputFiles []string

		// Load App
		application, err := app.Load(protocolFlag)
		if err != nil {
			return err
		}

		var stdinContent []byte
		if len(args) > 1 {
			outputFiles = make([]string, len(args[1:]))
			copy(outputFiles, args[1:])
			// First pass: identify stdin and map to target filenames
			stage, ok := application.Protocol.StageByID(stageID)
			if !ok {
				return fmt.Errorf("unknown stage: %s", stageID)
			}

			for i, arg := range outputFiles {
				if arg == "-" || strings.HasSuffix(arg, "=-") {
					if stdinContent == nil {
						var err error
						stdinContent, err = io.ReadAll(cmd.InOrStdin())
						if err != nil {
							return fmt.Errorf("failed to read from stdin: %w", err)
						}
					}

					if arg == "-" {
						if len(stage.Outputs) == 1 && !strings.Contains(stage.Outputs[0], "*") {
							outputFiles[i] = stage.Outputs[0]
						} else {
							return fmt.Errorf("ambiguous stdin mapping: stage has multiple or wildcard outputs %v. Use filename=- syntax.", stage.Outputs)
						}
					} else {
						outputFiles[i] = strings.TrimSuffix(arg, "=-")
					}
				}
			}
		} else {
			// Auto-discovery requires stage info
			stage, ok := application.Protocol.StageByID(stageID)
			if !ok {
				return fmt.Errorf("unknown stage: %s", stageID)
			}

			// Auto-discover changed files
			discovered, err := repository.DiscoverChangedFiles()
			if err != nil {
				return err
			}
			if len(discovered) == 0 {
				return fmt.Errorf("no changed files detected (untracked or modified)") // Using legacy error msg
			}

			// Filter discovered files against stage outputs
			if len(stage.Outputs) > 0 {
				filtered := make([]string, 0, len(discovered))
				for _, file := range discovered {
					rel, err := repository.ProjectRelPath(file)
					if err != nil {
						continue // Skip invalid paths
					}
					if isOutputMatch(stage.Outputs, rel) {
						filtered = append(filtered, file)
					}
				}
				if len(filtered) == 0 {
					return fmt.Errorf("found %d changed files, but none match stage outputs %v: %v", len(discovered), stage.Outputs, discovered)
				}
				outputFiles = filtered
			} else {
				return fmt.Errorf("no outputs defined for stage %q; cannot auto-discover changed files. Please list files explicitly.", stageID)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Auto-detected %d changed files matching stage outputs: %v\n", len(outputFiles), outputFiles)
		}
		promptFile, _ := cmd.Flags().GetString("prompt-file")

		// Write stdin content to workspace before passing to engine
		if stdinContent != nil {
			for i, arg := range args[1:] {
				if arg == "-" || strings.HasSuffix(arg, "=-") {
					target := outputFiles[i]
					resolved, err := repository.ResolveOutputPath(target)
					if err != nil {
						return err
					}
					if err := os.WriteFile(resolved, stdinContent, 0644); err != nil {
						return fmt.Errorf("failed to write stdin to %s: %w", target, err)
					}
				}
			}
		}

		// Delegate to App (pass context for graceful shutdown)
		if err := application.CompleteStage(cmd.Context(), stageID, outputFiles, completeForce, promptFile); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Completed stage %s\n", stageID)
		return nil
	},
}

func init() {
	completeCmd.Flags().String("prompt-file", "", "path to the prompt used for this stage")
	completeCmd.Flags().BoolVar(&completeForce, "force", false, "force overwrite of existing stage completion")
}

func isOutputMatch(patterns []string, file string) bool {
	for _, pattern := range patterns {
		if repository.MatchOutputPattern(pattern, file) {
			return true
		}
	}
	return false
}
