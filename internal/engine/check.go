package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"specfirst/internal/prompt"
	"specfirst/internal/store"
	"specfirst/internal/task"
	"specfirst/internal/workspace"
)

// Check runs all non-blocking validations (lint, tasks, approvals, outputs).
// Returns a list of warning messages grouped by category.
func (e *Engine) Check(failOnWarnings bool) error {
	warnings := make(map[string][]string)
	addWarning := func(category, msg string) {
		warnings[category] = append(warnings[category], msg)
	}

	// 1. Protocol Drift / Missing Approvals & Outputs
	if e.State.Protocol != "" && e.State.Protocol != e.Protocol.Name {
		addWarning("Protocol", fmt.Sprintf("Protocol drift: state=%s protocol=%s", e.State.Protocol, e.Protocol.Name))
	}

	for _, stage := range e.Protocol.Stages {
		if stage.Intent == "review" && len(stage.Outputs) > 0 {
			addWarning("Protocol", fmt.Sprintf("Review stage %s declares outputs", stage.ID))
		}
		if !e.State.IsStageCompleted(stage.ID) {
			continue
		}

		// Collect stored artifact paths for wildcard matching
		storedRel := []string{}
		if output, ok := e.State.StageOutputs[stage.ID]; ok {
			for _, file := range output.Files {
				rel, err := workspace.ArtifactRelFromState(file)
				if err != nil {
					addWarning("Artifacts", fmt.Sprintf("Invalid stored artifact path for stage %s: %s (%v)", stage.ID, file, err))
					continue
				}
				// Clean up the path relative to the stage artifact root if necessary
				relPath := filepath.FromSlash(rel)
				cleanRel := relPath
				stagePrefix := stage.ID + string(os.PathSeparator)
				if strings.HasPrefix(relPath, stagePrefix) {
					cleanRel = strings.TrimPrefix(relPath, stagePrefix)
				}
				storedRel = append(storedRel, cleanRel)
			}
		}

		for _, output := range stage.Outputs {
			if output == "" {
				continue
			}
			if strings.Contains(output, "*") {
				found := false
				for _, rel := range storedRel {
					if workspace.MatchOutputPattern(output, rel) {
						found = true
						break
					}
				}
				if !found {
					addWarning("Outputs", fmt.Sprintf("Missing output for stage %s: %s (no stored artifacts match)", stage.ID, output))
				}
				continue
			}
			expected := store.ArtifactsPath(stage.ID, output)
			if _, err := os.Stat(expected); os.IsNotExist(err) {
				addWarning("Outputs", fmt.Sprintf("Missing output for stage %s: %s", stage.ID, expected))
			} else if stage.Output != nil && len(stage.Output.Sections) > 0 {
				// Check for required sections
				content, err := os.ReadFile(expected)
				if err == nil {
					sContent := string(content)
					for _, sectionHeader := range stage.Output.Sections {
						// Check for markdown header
						// We check for "# Header" or "## Header"
						if !strings.Contains(sContent, "# "+sectionHeader) && !strings.Contains(sContent, "## "+sectionHeader) {
							addWarning("Structure", fmt.Sprintf("Missing section %q in %s", sectionHeader, expected))
						}
					}
				}
			}
		}
	}

	for _, approval := range e.Protocol.Approvals {
		if e.State.IsStageCompleted(approval.Stage) {
			if !e.State.HasApproval(approval.Stage, approval.Role) {
				addWarning("Approvals", fmt.Sprintf("Missing approval for stage %s (role: %s)", approval.Stage, approval.Role))
			}
		}
	}

	// 2. Task List Validation
	for _, stage := range e.Protocol.Stages {
		if stage.Type == "decompose" && e.State.IsStageCompleted(stage.ID) {
			output, ok := e.State.StageOutputs[stage.ID]
			if ok {
				for _, file := range output.Files {
					artifactPath, err := workspace.ArtifactAbsFromState(file)
					if err == nil {
						content, err := os.ReadFile(artifactPath)
						if err == nil {
							taskList, err := task.Parse(string(content))
							if err == nil {
								taskWarnings := taskList.Validate()
								for _, tw := range taskWarnings {
									addWarning("Tasks", fmt.Sprintf("[%s]: %s", file, tw))
								}
							}
						}
					}
				}
			}
		}
	}

	// 3. Prompt Quality Checks
	for _, stage := range e.Protocol.Stages {
		// Skip prompt checks if dependencies aren't met to avoid "missing input" errors
		// for stages the user hasn't reached yet.
		if err := e.RequireStageDependencies(stage); err != nil {
			continue
		}

		stageIDs := make([]string, 0, len(e.Protocol.Stages))
		for _, s := range e.Protocol.Stages {
			stageIDs = append(stageIDs, s.ID)
		}

		compiledPrompt, err := e.CompilePrompt(stage, stageIDs, CompileOptions{})
		if err != nil {
			addWarning("Prompts", fmt.Sprintf("Prompt compile (%s): %v", stage.ID, err))
			continue
		}
		schema := prompt.DefaultSchema()
		schema.Merge(e.Protocol.Lint)
		if stage.Prompt != nil {
			schema.Merge(stage.Prompt.Lint)
		}
		result := prompt.Validate(compiledPrompt, schema)
		for _, w := range result.Warnings {
			addWarning("Prompts", fmt.Sprintf("Quality (%s): %s", stage.ID, w))
		}
		ambiguities := prompt.ContainsAmbiguity(compiledPrompt)
		for _, a := range ambiguities {
			addWarning("Prompts", fmt.Sprintf("Ambiguity (%s): %s", stage.ID, a))
		}
	}

	if len(warnings) == 0 {
		fmt.Println("No issues found.")
		return nil
	}

	fmt.Println("Warnings (advisory):")
	var categories []string
	for cat := range warnings {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	totalWarnings := 0
	for _, cat := range categories {
		list := warnings[cat]
		if len(list) > 0 {
			fmt.Printf("\n* %s (%d)\n", cat, len(list))
			for _, w := range list {
				fmt.Printf("  - %s\n", w)
			}
			totalWarnings += len(list)
		}
	}

	if failOnWarnings {
		return fmt.Errorf("check failed with %d warnings", totalWarnings)
	}
	return nil
}
