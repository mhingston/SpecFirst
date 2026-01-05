package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"specfirst/internal/config"
	"specfirst/internal/prompt"
	"specfirst/internal/protocol"
	"specfirst/internal/task"
	tmplpkg "specfirst/internal/template"
)

var taskCmd = &cobra.Command{
	Use:   "task [task-id]",
	Short: "Generate implementation prompt for a specific task",
	Long: `Generate an implementation prompt for a specific task from a completed decomposition stage.
If no task ID is provided, it lists all available tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		proto, err := loadProtocol(cfg.Protocol)
		if err != nil {
			return err
		}

		s, err := loadState()
		if err != nil {
			return err
		}

		// Find the decompose stage
		var decomposeStageID string
		for _, stage := range proto.Stages {
			if stage.Type == "decompose" {
				decomposeStageID = stage.ID
				break
			}
		}

		if decomposeStageID == "" {
			return fmt.Errorf("no stage of type 'decompose' found in protocol %q", proto.Name)
		}

		if !s.IsStageCompleted(decomposeStageID) {
			return fmt.Errorf("decompose stage %q has not been completed", decomposeStageID)
		}

		// Find the artifact for the decompose stage
		output, ok := s.StageOutputs[decomposeStageID]
		if !ok || len(output.Files) == 0 {
			return fmt.Errorf("no artifacts found for decompose stage %q", decomposeStageID)
		}

		// Search through all artifacts to find a valid task list
		var taskList task.TaskList
		var foundTaskList bool
		for _, file := range output.Files {
			artifactPath, err := artifactAbsFromState(file)
			if err != nil {
				continue
			}

			content, err := os.ReadFile(artifactPath)
			if err != nil {
				continue
			}

			parsed, err := task.Parse(string(content))
			if err == nil && len(parsed.Tasks) > 0 {
				taskList = parsed
				foundTaskList = true
				break
			}
		}

		if !foundTaskList {
			return fmt.Errorf("no valid task list found in artifacts for decompose stage %q", decomposeStageID)
		}

		// Find the task_prompt stage that refers to this decompose stage to get common inputs
		var taskPromptStage protocol.Stage
		var foundTaskPrompt bool
		for _, stg := range proto.Stages {
			if stg.Type == "task_prompt" && stg.Source == decomposeStageID {
				taskPromptStage = stg
				foundTaskPrompt = true
				break
			}
		}

		var artifactInputs []tmplpkg.Input
		if foundTaskPrompt {
			// Gather artifacts for the task_prompt stage (requirements, design, etc.)
			stageIDs := stageIDList(proto)
			artifactInputs = make([]tmplpkg.Input, 0, len(taskPromptStage.Inputs))
			for _, input := range taskPromptStage.Inputs {
				path, err := artifactPathForInput(input, taskPromptStage.DependsOn, stageIDs)
				if err != nil {
					// Skip missing artifacts for tasks to avoid hard failure
					continue
				}
				content, err := os.ReadFile(path)
				if err != nil {
					continue
				}
				artifactInputs = append(artifactInputs, tmplpkg.Input{Name: input, Content: string(content)})
			}
		}

		if len(args) == 0 {
			// List tasks
			fmt.Println("Available tasks:")
			sort.Slice(taskList.Tasks, func(i, j int) bool {
				return taskList.Tasks[i].ID < taskList.Tasks[j].ID
			})
			for _, t := range taskList.Tasks {
				fmt.Printf("- %-10s: %s\n", t.ID, t.Title)
			}

			// Surface validation warnings
			if warnings := taskList.Validate(); len(warnings) > 0 {
				fmt.Fprintln(os.Stderr, "\nWarnings:")
				for _, w := range warnings {
					fmt.Fprintf(os.Stderr, "- %s\n", w)
				}
			}
			return nil
		}

		// Surface validation warnings before generating specific task prompt
		if warnings := taskList.Validate(); len(warnings) > 0 {
			fmt.Fprintln(os.Stderr, "Warnings:")
			for _, w := range warnings {
				fmt.Fprintf(os.Stderr, "- %s\n", w)
			}
			fmt.Fprintln(os.Stderr)
		}

		taskID := args[0]
		var targetTask *task.Task
		for _, t := range taskList.Tasks {
			if t.ID == taskID {
				targetTask = &t
				break
			}
		}

		if targetTask == nil {
			return fmt.Errorf("task %q not found in decomposition output", taskID)
		}

		// Generate the prompt
		promptText := generateTaskPrompt(*targetTask, artifactInputs, cfg, proto)

		// Surface ambiguity warnings if found in the generated prompt
		if issues := prompt.ContainsAmbiguity(promptText); len(issues) > 0 {
			fmt.Fprintln(os.Stderr, "Ambiguity Warnings in generated prompt:")
			for _, issue := range issues {
				fmt.Fprintf(os.Stderr, "- %s\n", issue)
			}
			fmt.Fprintln(os.Stderr)
		}

		// Respect output format flags if added later, for now just print raw
		fmt.Println(promptText)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
}

func generateTaskPrompt(t task.Task, inputs []tmplpkg.Input, cfg config.Config, proto protocol.Protocol) string {
	var sb strings.Builder

	// Header
	sb.WriteString("---\n")
	sb.WriteString("intent: implementation\n")
	sb.WriteString("expected_output: code_diff\n")
	sb.WriteString("determinism: medium\n")
	sb.WriteString("allowed_creativity: low\n")
	sb.WriteString("---\n\n")

	sb.WriteString(fmt.Sprintf("# Implement %s: %s\n\n", t.ID, t.Title))

	if len(inputs) > 0 {
		sb.WriteString("## Context\n")
		for _, input := range inputs {
			sb.WriteString(fmt.Sprintf("<artifact name=\"%s\">\n", input.Name))
			sb.WriteString(input.Content)
			if !strings.HasSuffix(input.Content, "\n") {
				sb.WriteString("\n")
			}
			sb.WriteString("</artifact>\n\n")
		}
	}

	sb.WriteString("## Goal\n")
	sb.WriteString(t.Goal + "\n\n")

	if len(t.AcceptanceCriteria) > 0 {
		sb.WriteString("## Acceptance Criteria\n")
		for _, ac := range t.AcceptanceCriteria {
			sb.WriteString(fmt.Sprintf("- %s\n", ac))
		}
		sb.WriteString("\n")
	}

	if len(t.FilesTouched) > 0 {
		sb.WriteString("## Known Files\n")
		for _, f := range t.FilesTouched {
			sb.WriteString(fmt.Sprintf("- %s\n", f))
		}
		sb.WriteString("\n")
	}

	if len(t.Dependencies) > 0 {
		sb.WriteString("## Dependencies\n")
		sb.WriteString("This task depends on the completion of:\n")
		for _, dep := range t.Dependencies {
			sb.WriteString(fmt.Sprintf("- %s\n", dep))
		}
		sb.WriteString("\n")
	}

	if len(t.TestPlan) > 0 {
		sb.WriteString("## Test Plan\n")
		for _, tp := range t.TestPlan {
			sb.WriteString(fmt.Sprintf("- %s\n", tp))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Instructions\n")
	sb.WriteString("- Produce ONLY the minimal changes required.\n")
	sb.WriteString("- Maintain project standards and existing architecture.\n\n")

	sb.WriteString("## Expected Output\n")
	sb.WriteString("- Format: unified diff\n")
	sb.WriteString("- Scope: only listed files unless explicitly justified\n")
	sb.WriteString("- Tests: added or updated if behavior changes\n\n")

	sb.WriteString("## Assumptions\n")
	sb.WriteString("- (List explicitly before implementation if any)\n")

	return sb.String()
}
