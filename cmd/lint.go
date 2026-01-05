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
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Run non-blocking checks on the workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
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
			storedRel := []string{}
			if output, ok := s.StageOutputs[stage.ID]; ok {
				for _, file := range output.Files {
					rel, err := artifactRelFromState(file)
					if err != nil {
						addWarning("Artifacts", fmt.Sprintf("Invalid stored artifact path for stage %s: %s (%v)", stage.ID, file, err))
						continue
					}
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
				if output == "" || strings.Contains(output, "*") {
					if output != "" && strings.Contains(output, "*") {
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
					}
					continue
				}
				expected := store.ArtifactsPath(stage.ID, output)
				if _, err := os.Stat(expected); err != nil {
					if os.IsNotExist(err) {
						addWarning("Outputs", fmt.Sprintf("Missing output for stage %s: %s", stage.ID, expected))
					} else {
						addWarning("Outputs", fmt.Sprintf("Unable to read output %s: %v", expected, err))
					}
				}
			}
		}

		for stageID, output := range s.StageOutputs {
			for _, file := range output.Files {
				abs, err := artifactAbsFromState(file)
				if err != nil {
					addWarning("Artifacts", fmt.Sprintf("Invalid stored artifact path for stage %s: %s (%v)", stageID, file, err))
					continue
				}
				info, err := os.Stat(abs)
				if err != nil {
					addWarning("Artifacts", fmt.Sprintf("Missing stored artifact for stage %s: %s", stageID, abs))
					continue
				}
				if info.Size() == 0 {
					addWarning("Artifacts", fmt.Sprintf("Empty artifact for stage %s: %s", stageID, file))
				}
			}

			if output.PromptHash == "" {
				addWarning("Artifacts", fmt.Sprintf("Missing prompt hash for stage %s", stageID))
			}
		}

		for _, approval := range proto.Approvals {
			if !s.IsStageCompleted(approval.Stage) {
				continue
			}
			records := s.Approvals[approval.Stage]
			if !hasApproval(records, approval.Role) {
				addWarning("Approvals", fmt.Sprintf("Missing approval for stage %s (role: %s)", approval.Stage, approval.Role))
			}
		}

		for key, value := range cfg.Constraints {
			if strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
				addWarning("Config", "Empty constraint key/value found")
				break
			}
		}

		// Prompt quality checks for each stage
		for _, stage := range proto.Stages {
			// Compile prompt and validate
			compiledPrompt, err := compilePrompt(stage, cfg, stageIDList(proto))
			if err != nil {
				// Skip if prompt can't compile (may be missing dependencies)
				continue
			}

			// Schema validation
			schema := prompt.DefaultSchema()
			schema.Merge(proto.Lint)
			if stage.Prompt != nil {
				schema.Merge(stage.Prompt.Lint)
			}
			result := prompt.Validate(compiledPrompt, schema)
			for _, w := range result.Warnings {
				addWarning("Prompts", fmt.Sprintf("Quality (%s): %s", stage.ID, w))
			}

			// Structure validation based on stage type
			stageType := stage.Type
			if stageType == "" {
				stageType = "spec"
			}
			structResult := prompt.ValidateStructure(compiledPrompt, stageType)
			for _, w := range structResult.Warnings {
				addWarning("Prompts", fmt.Sprintf("Structure (%s): %s", stage.ID, w))
			}

			// Ambiguity detection
			ambiguities := prompt.ContainsAmbiguity(compiledPrompt)
			for _, a := range ambiguities {
				addWarning("Prompts", fmt.Sprintf("Ambiguity (%s): %s", stage.ID, a))
			}
		}

		if len(warnings) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No lint warnings.")
			return nil
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Lint warnings:")
		var categories []string
		for cat := range warnings {
			categories = append(categories, cat)
		}
		sort.Strings(categories)

		for _, cat := range categories {
			list := warnings[cat]
			fmt.Fprintf(cmd.OutOrStdout(), "\n* %s (%d)\n", cat, len(list))
			for _, w := range list {
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", w)
			}
		}
		return nil
	},
}
