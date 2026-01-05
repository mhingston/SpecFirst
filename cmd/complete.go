package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"specfirst/internal/protocol"
	"specfirst/internal/state"
	"specfirst/internal/store"
	"specfirst/internal/task"
	"strings"
)

var completeForce bool

var completeCmd = &cobra.Command{
	Use:   "complete <stage-id> <output-files...>",
	Short: "Mark a stage as complete and store outputs",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stageID := args[0]
		var outputFiles []string

		// Load config and protocol once
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		proto, err := loadProtocol(activeProtocolName(cfg))
		if err != nil {
			return err
		}
		stage, ok := proto.StageByID(stageID)
		if !ok {
			return fmt.Errorf("unknown stage: %s", stageID)
		}

		if len(args) > 1 {
			outputFiles = args[1:]
		} else {
			// Auto-discover changed files
			discovered, err := discoverChangedFiles()
			if err != nil {
				return err
			}
			if len(discovered) == 0 {
				return fmt.Errorf("no changed files detected (untracked or modified)")
			}

			// Filter discovered files against stage outputs
			if len(stage.Outputs) > 0 {
				filtered := make([]string, 0, len(discovered))
				for _, file := range discovered {
					rel, err := outputRelPath(file)
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
				// If no outputs defined, we don't know what to include automatically.
				return fmt.Errorf("no outputs defined for stage %q; cannot auto-discover changed files. Please list files explicitly.", stageID)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Auto-detected %d changed files matching stage outputs: %v\n", len(outputFiles), outputFiles)
		}
		promptFile, _ := cmd.Flags().GetString("prompt-file")

		s, err := loadState()
		if err != nil {
			return err
		}
		s = ensureStateInitialized(s, proto)

		if !stageNoStrict {
			if err := requireStageDependencies(s, stage); err != nil {
				return err
			}
		}

		// Check for duplicate completion. A stage is considered completed if it exists
		// in the completed_stages list or has an entry in stage_outputs.
		_, hasOutput := s.StageOutputs[stageID]
		if (s.IsStageCompleted(stageID) || hasOutput) && !completeForce {
			return fmt.Errorf("stage %s already completed; use --force to overwrite", stageID)
		}

		if err := validateOutputs(stage, outputFiles); err != nil {
			return err
		}

		// Task List Validation
		if stage.Type == "decompose" {
			for _, file := range outputFiles {
				content, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("failed to read task file %s: %w", file, err)
				}
				taskList, err := task.Parse(string(content))
				if err != nil {
					return fmt.Errorf("failed to parse task list in %s: %w", file, err)
				}
				if warnings := taskList.Validate(); len(warnings) > 0 {
					return fmt.Errorf("invalid task list in %s:\n%s", file, strings.Join(warnings, "\n"))
				}
			}
		}

		var oldFiles []string
		if completeForce {
			if old, exists := s.StageOutputs[stageID]; exists {
				oldFiles = old.Files
				if len(outputFiles) < len(oldFiles) {
					fmt.Fprintf(cmd.ErrOrStderr(), "Warning: forcing completion with %d files, but stage previously had %d files. Obsolete artifacts will be removed.\n", len(outputFiles), len(oldFiles))
				}
			}
		}

		stored := make([]string, 0, len(outputFiles))
		for _, output := range outputFiles {
			rel, err := outputRelPath(output)
			if err != nil {
				return err
			}
			resolved, err := resolveOutputPath(output)
			if err != nil {
				return err
			}
			dest := store.ArtifactsPath(stageID, rel)
			if err := copyFile(resolved, dest); err != nil {
				return err
			}
			stored = append(stored, filepath.ToSlash(filepath.Join(stageID, rel)))
		}

		// Cleanup obsolete artifacts after successful copy
		if completeForce && len(oldFiles) > 0 {
			newFilesMap := make(map[string]bool)
			for _, f := range stored {
				newFilesMap[f] = true
			}
			for _, oldFile := range oldFiles {
				if !newFilesMap[oldFile] {
					abs, err := artifactAbsFromState(oldFile)
					if err == nil {
						_ = os.Remove(abs)
					}
				}
			}
		}

		var promptHashValue string
		if promptFile != "" {
			hash, err := fileHash(promptFile)
			if err != nil {
				return err
			}
			promptHashValue = hash
		} else {
			prompt, err := compilePrompt(stage, cfg, stageIDList(proto))
			if err != nil {
				return err
			}
			promptHashValue = promptHash(prompt)
		}

		s.StageOutputs[stageID] = state.StageOutput{
			CompletedAt: time.Now().UTC(),
			Files:       stored,
			PromptHash:  promptHashValue,
		}
		if !s.IsStageCompleted(stageID) {
			s.CompletedStages = append(s.CompletedStages, stageID)
		}

		// Update CurrentStage to the NEXT stage only if we are progressing
		currentIndex := -1
		completedIndex := -1
		for i, stg := range proto.Stages {
			if stg.ID == s.CurrentStage {
				currentIndex = i
			}
			if stg.ID == stageID {
				completedIndex = i
			}
		}

		// If we completed the current stage (or a later one), bump CurrentStage to the next one
		if completedIndex >= currentIndex {
			if completedIndex+1 < len(proto.Stages) {
				s.CurrentStage = proto.Stages[completedIndex+1].ID
			} else {
				// We reached the end. Optional: verify if we want to signal "done" differently.
				// For now, we leave it at the last stage, but since it's in CompletedStages, users know it's done.
			}
		}

		if err := saveState(s); err != nil {
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

func validateOutputs(stage protocol.Stage, outputFiles []string) error {
	if len(stage.Outputs) == 0 {
		return nil
	}
	relOutputs := make([]string, 0, len(outputFiles))
	for _, output := range outputFiles {
		rel, err := outputRelPath(output)
		if err != nil {
			return err
		}
		relOutputs = append(relOutputs, rel)
	}
	missing := []string{}
	for _, expected := range stage.Outputs {
		if expected == "" {
			continue
		}
		found := false
		for _, output := range relOutputs {
			if matchOutputPattern(expected, output) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, expected)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing expected outputs: %v", missing)
	}
	return nil
}

func isOutputMatch(patterns []string, file string) bool {
	for _, pattern := range patterns {
		if matchOutputPattern(pattern, file) {
			return true
		}
	}
	return false
}
