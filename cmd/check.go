package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"specfirst/internal/prompt"
	"specfirst/internal/store"
	"specfirst/internal/task"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run all non-blocking validations (lint, tasks, approvals, outputs)",
	RunE: func(cmd *cobra.Command, args []string) error {
		failOnWarnings, _ := cmd.Flags().GetBool("fail-on-warnings")

		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		proto, err := loadProtocol(activeProtocolName(cfg))
		if err != nil {
			return err
		}
		s, err := loadState()
		if err != nil {
			return err
		}

		warnings := make(map[string][]string)
		addWarning := func(category, msg string) {
			warnings[category] = append(warnings[category], msg)
		}

		// 1. Protocol Drift / Missing Approvals & Outputs
		if s.Protocol != "" && s.Protocol != proto.Name {
			addWarning("Protocol", fmt.Sprintf("Protocol drift: state=%s protocol=%s", s.Protocol, proto.Name))
		}

		for _, stage := range proto.Stages {
			if stage.Intent == "review" && len(stage.Outputs) > 0 {
				addWarning("Protocol", fmt.Sprintf("Review stage %s declares outputs", stage.ID))
			}
			if !s.IsStageCompleted(stage.ID) {
				continue
			}

			// Collect stored artifact paths for wildcard matching
			storedRel := []string{}
			if output, ok := s.StageOutputs[stage.ID]; ok {
				for _, file := range output.Files {
					rel, err := artifactRelFromState(file)
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
						if matchOutputPattern(output, rel) {
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

		for _, approval := range proto.Approvals {
			if s.IsStageCompleted(approval.Stage) {
				records := s.Approvals[approval.Stage]
				if !hasApproval(records, approval.Role) {
					addWarning("Approvals", fmt.Sprintf("Missing approval for stage %s (role: %s)", approval.Stage, approval.Role))
				}
			}
		}

		// 2. Task List Validation

		for _, stage := range proto.Stages {
			if stage.Type == "decompose" && s.IsStageCompleted(stage.ID) {
				output, ok := s.StageOutputs[stage.ID]
				if ok {
					for _, file := range output.Files {
						artifactPath, err := artifactAbsFromState(file)
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
		for _, stage := range proto.Stages {
			// Skip prompt checks if dependencies aren't met to avoid "missing input" errors
			// for stages the user hasn't reached yet.
			if err := requireStageDependencies(s, stage); err != nil {
				continue
			}

			compiledPrompt, err := compilePrompt(stage, cfg, stageIDList(proto))
			if err != nil {
				addWarning("Prompts", fmt.Sprintf("Prompt compile (%s): %v", stage.ID, err))
				continue
			}
			schema := prompt.DefaultSchema()
			schema.Merge(proto.Lint)
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
			fmt.Fprintln(cmd.OutOrStdout(), "No issues found.")
			return nil
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Warnings (advisory):")
		var categories []string
		for cat := range warnings {
			categories = append(categories, cat)
		}
		sort.Strings(categories)

		totalWarnings := 0
		for _, cat := range categories {
			list := warnings[cat]
			if len(list) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "\n* %s (%d)\n", cat, len(list))
				for _, w := range list {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", w)
				}
				totalWarnings += len(list)
			}
		}

		if failOnWarnings {
			return fmt.Errorf("check failed with %d warnings", totalWarnings)
		}
		return nil
	},
}

func init() {
	checkCmd.Flags().Bool("fail-on-warnings", false, "exit with code 1 if warnings are found")
}
