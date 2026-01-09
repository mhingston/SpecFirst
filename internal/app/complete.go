package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"specfirst/internal/domain"
	"specfirst/internal/repository"
	"specfirst/internal/utils"
)

// CompleteStage marks a stage as complete and stores outputs.
func (app *Application) CompleteStage(ctx context.Context, stageID string, outputFiles []string, force bool, promptFile string) error {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return err
	}

	stage, ok := app.Protocol.StageByID(stageID)
	if !ok {
		return fmt.Errorf("unknown stage: %s", stageID)
	}

	// Dependency Check
	if err := app.RequireStageDependencies(stage); err != nil {
		return err
	}

	// Duplicate Completion Check
	_, hasOutput := app.State.StageOutputs[stageID]
	if (app.State.IsStageCompleted(stageID) || hasOutput) && !force {
		return fmt.Errorf("stage %s already completed; use --force to overwrite", stageID)
	}

	// Validate Outputs
	if err := app.ValidateOutputs(stage, outputFiles); err != nil {
		return err
	}

	// Ambiguity Gates
	if !force {
		if err := app.ValidateAmbiguityGates(stage); err != nil {
			return err
		}
	}

	// Task List Validation (File I/O - check context)
	if stage.Type == "decompose" {
		for _, file := range outputFiles {
			if err := ctx.Err(); err != nil {
				return err
			}

			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read task file %s: %w", file, err)
			}
			taskList, err := domain.ParseTaskList(string(content))
			if err != nil {
				return fmt.Errorf("failed to parse task list in %s: %w", file, err)
			}
			if warnings := taskList.Validate(); len(warnings) > 0 {
				return fmt.Errorf("invalid task list in %s:\n%s", file, strings.Join(warnings, "\n"))
			}
		}
	}

	// Handle existing files (cleanup if force)
	var oldFiles []string
	if force {
		if old, exists := app.State.StageOutputs[stageID]; exists {
			oldFiles = old.Files
			if len(outputFiles) < len(oldFiles) {
				fmt.Fprintf(os.Stderr, "Warning: forcing completion with %d files, but stage previously had %d files. Obsolete artifacts will be removed.\n", len(outputFiles), len(oldFiles))
			}
		}
	}

	// Store Artifacts (File I/O - check context)
	stored := make([]string, 0, len(outputFiles))
	for _, output := range outputFiles {
		if err := ctx.Err(); err != nil {
			return err
		}

		resolved, err := filepath.Abs(output)
		if err != nil {
			return err
		}

		relPath, err := repository.ProjectRelPath(resolved)
		if err != nil {
			return err
		}

		dest := repository.ArtifactsPath(stageID, relPath)
		if err := utils.CopyFile(resolved, dest); err != nil {
			return err
		}
		stored = append(stored, filepath.ToSlash(filepath.Join(stageID, relPath)))
	}

	// Cleanup obsolete artifacts
	if force && len(oldFiles) > 0 {
		newFilesMap := make(map[string]bool)
		for _, f := range stored {
			newFilesMap[f] = true
		}
		for _, oldFile := range oldFiles {
			if !newFilesMap[oldFile] {
				abs, err := repository.ArtifactAbsFromState(oldFile)
				if err == nil {
					_ = os.Remove(abs)
				}
			}
		}
	}

	// Calculate Hash (Potentially expensive - check context)
	if err := ctx.Err(); err != nil {
		return err
	}

	var promptHashValue string
	if promptFile != "" {
		hash, err := utils.FileHash(promptFile)
		if err != nil {
			return err
		}
		promptHashValue = hash
	} else {
		// Use default compile
		stageIDs := make([]string, 0, len(app.Protocol.Stages))
		for _, s := range app.Protocol.Stages {
			stageIDs = append(stageIDs, s.ID)
		}
		prompt, err := app.CompilePrompt(stage, stageIDs, CompileOptions{})
		if err != nil {
			return err
		}
		promptHashValue = utils.PromptHash(prompt)
	}

	// Update State
	app.State.StageOutputs[stageID] = domain.StageOutput{
		CompletedAt: time.Now().UTC(),
		Files:       stored,
		PromptHash:  promptHashValue,
	}
	if !app.State.IsStageCompleted(stageID) {
		app.State.CompletedStages = append(app.State.CompletedStages, stageID)
	}

	// Auto-Advance Stage
	if next := app.Protocol.NextStage(stageID); next != nil {
		// Only advance if currently at the completed stage (or previous)
		// Logic: if current stage is <= completed stage, move to next.
		// Detailed logic:
		// existing logic checked if completedIndex >= currentIndex.
		// We can infer: if we just completed stageID, and stageID is the current stage, we advance.
		// if stageID was a past stage, we don't advance current (unless we are replaying).
		// Simplification: If State.CurrentStage == stageID, advance.
		if app.State.CurrentStage == stageID {
			app.State.CurrentStage = next.ID
		}
	}

	return app.SaveState()
}

func (app *Application) ValidateOutputs(stage domain.Stage, outputFiles []string) error {
	if len(stage.Outputs) == 0 {
		return nil
	}
	// Sev3 Fix: Ensure outputs are provided if required by the stage.
	if len(outputFiles) == 0 {
		return fmt.Errorf("stage %q requires output artifacts but none were provided", stage.ID)
	}
	return nil
}

func (app *Application) ValidateAmbiguityGates(stage domain.Stage) error {
	// 1. Check Open Questions Limit
	if stage.MaxOpenQuestions != nil {
		limit := *stage.MaxOpenQuestions
		openCount := 0
		for _, q := range app.State.Epistemics.OpenQuestions {
			if q.Status == "open" {
				openCount++
			}
		}
		if openCount > limit {
			return fmt.Errorf("ambiguity gate failure: %d open questions (max %d)", openCount, limit)
		}
	}

	// 2. Check Must Resolve Tags
	if len(stage.MustResolveTags) > 0 {
		for _, q := range app.State.Epistemics.OpenQuestions {
			if q.Status != "open" {
				continue
			}
			for _, tag := range q.Tags {
				for _, requiredTag := range stage.MustResolveTags {
					if tag == requiredTag {
						return fmt.Errorf("ambiguity gate failure: open question %q has must-resolve tag %q", q.ID, tag)
					}
				}
			}
		}
	}

	// 3. Check High Risks Limit
	if stage.MaxHighRisksUnmitigated != nil {
		limit := *stage.MaxHighRisksUnmitigated
		highRiskCount := 0
		for _, r := range app.State.Epistemics.Risks {
			if r.Status == "mitigated" || r.Status == "accepted" {
				continue
			}
			if r.Severity == "high" {
				highRiskCount++
			}
		}
		if highRiskCount > limit {
			return fmt.Errorf("ambiguity gate failure: %d unmitigated high risks (max %d)", highRiskCount, limit)
		}
	}

	return nil
}
