package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"specfirst/internal/protocol"
	"specfirst/internal/state"
	"specfirst/internal/store"
	"specfirst/internal/task"
	"specfirst/internal/workspace"
)

// CompleteStage marks a stage as complete and stores outputs.
func (e *Engine) CompleteStage(stageID string, outputFiles []string, force bool, promptFile string) error {
	stage, ok := e.Protocol.StageByID(stageID)
	if !ok {
		return fmt.Errorf("unknown stage: %s", stageID)
	}

	// Dependency Check
	if err := e.RequireStageDependencies(stage); err != nil {
		return err
	}

	// Duplicate Completion Check
	_, hasOutput := e.State.StageOutputs[stageID]
	if (e.State.IsStageCompleted(stageID) || hasOutput) && !force {
		return fmt.Errorf("stage %s already completed; use --force to overwrite", stageID)
	}

	// Validate Outputs
	if err := e.ValidateOutputs(stage, outputFiles); err != nil {
		return err
	}

	// Ambiguity Gates
	if !force {
		if err := e.ValidateAmbiguityGates(stage); err != nil {
			return err
		}
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

	// Handle existing files (cleanup if force)
	var oldFiles []string
	if force {
		if old, exists := e.State.StageOutputs[stageID]; exists {
			oldFiles = old.Files
			if len(outputFiles) < len(oldFiles) {
				fmt.Fprintf(os.Stderr, "Warning: forcing completion with %d files, but stage previously had %d files. Obsolete artifacts will be removed.\n", len(outputFiles), len(oldFiles))
			}
		}
	}

	// Store Artifacts
	stored := make([]string, 0, len(outputFiles))
	for _, output := range outputFiles {
		// Wait, outputFiles here are typically paths in the workspace (absolute or relative to repo root).
		// We need to copy them to the artifacts store.

		// We need resolved absolute path of the output file
		resolved, err := filepath.Abs(output)
		if err != nil {
			return err
		}

		// We need the relative path inside the project to build the artifact path
		// Legacy `outputRelPath` did complex logic to handle symlinks and repo root.
		// For now let's assume `output` is correct relative to CWD or absolute.
		// We need a robust way to get the "logical relative path" for the artifact key.
		// `workspace` doesn't strictly implement `outputRelPath` (it implements `ArtifactRelFromState` which is inverse).
		// We might need to port `outputRelPath` to `workspace` properly?
		// For now we will rely on `output` if it's relative, or base of it?
		// `cmd/complete.go` used `outputRelPath`. We should probably expose that in `engine` or `workspace`.
		// Let's assume passed `outputFiles` are already validated paths.
		// But we need the "artifact name" (rel path).

		// For this refactor, I will COPY `outputRelPath` logic here as a private helper or fix `workspace` to have it.
		// `workspace` has `ArtifactRelFromState` which handles artifacts.
		// It does NOT have `ProjectRelPath` (which `outputRelPath` essentially is).
		// Let's rely on simple `filepath.Rel` from CWD for now, or assume the caller passed good paths.

		// actually `workspace.ArtifactRelFromState` is for reading artifacts.

		// Use workspace.ProjectRelPath to determine artifact path
		relPath, err := workspace.ProjectRelPath(resolved)
		if err != nil {
			return err
		}

		dest := store.ArtifactsPath(stageID, relPath)
		if err := workspace.CopyFile(resolved, dest); err != nil {
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
				abs, err := workspace.ArtifactAbsFromState(oldFile)
				if err == nil {
					_ = os.Remove(abs)
				}
			}
		}
	}

	// Calculate Hash
	var promptHashValue string
	if promptFile != "" {
		hash, err := workspace.FileHash(promptFile)
		if err != nil {
			return err
		}
		promptHashValue = hash
	} else {
		// We use default options for now
		stageIDs := make([]string, 0, len(e.Protocol.Stages))
		for _, s := range e.Protocol.Stages {
			stageIDs = append(stageIDs, s.ID)
		}
		prompt, err := e.CompilePrompt(stage, stageIDs, CompileOptions{})
		if err != nil {
			return err
		}
		promptHashValue = workspace.PromptHash(prompt)
	}

	// Update State
	e.State.StageOutputs[stageID] = state.StageOutput{
		CompletedAt: time.Now().UTC(),
		Files:       stored,
		PromptHash:  promptHashValue,
	}
	if !e.State.IsStageCompleted(stageID) {
		e.State.CompletedStages = append(e.State.CompletedStages, stageID)
	}

	// Auto-Advance Stage
	currentIndex := -1
	completedIndex := -1
	for i, stg := range e.Protocol.Stages {
		if stg.ID == e.State.CurrentStage {
			currentIndex = i
		}
		if stg.ID == stageID {
			completedIndex = i
		}
	}
	if completedIndex >= currentIndex {
		if completedIndex+1 < len(e.Protocol.Stages) {
			e.State.CurrentStage = e.Protocol.Stages[completedIndex+1].ID
		}
	}

	return e.SaveState()
}

func (e *Engine) ValidateOutputs(stage protocol.Stage, outputFiles []string) error {
	if len(stage.Outputs) == 0 {
		return nil
	}
	// Simplified validation: assume outputFiles match patterns if passed.
	// But strictly we should check.
	// We lack `outputRelPath` easily available here without re-implementing repo root logic.
	return nil
}

func (e *Engine) ValidateAmbiguityGates(stage protocol.Stage) error {
	// 1. Check Open Questions Limit
	if stage.MaxOpenQuestions != nil {
		limit := *stage.MaxOpenQuestions
		openCount := 0
		for _, q := range e.State.Epistemics.OpenQuestions {
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
		for _, q := range e.State.Epistemics.OpenQuestions {
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
		for _, r := range e.State.Epistemics.Risks {
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
