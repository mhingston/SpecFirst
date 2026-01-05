package task

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Task represents a decomposed work item.
type Task struct {
	ID                 string   `yaml:"id" json:"id"`
	Title              string   `yaml:"title" json:"title"`
	Goal               string   `yaml:"goal" json:"goal"`
	AcceptanceCriteria []string `yaml:"acceptance_criteria" json:"acceptance_criteria"`
	Dependencies       []string `yaml:"dependencies" json:"dependencies"`
	FilesTouched       []string `yaml:"files_touched" json:"files_touched"`
	RiskLevel          string   `yaml:"risk_level" json:"risk_level"`           // low, medium, high
	EstimatedScope     string   `yaml:"estimated_scope" json:"estimated_scope"` // S, M, L
	NonGoals           []string `yaml:"non_goals,omitempty" json:"non_goals,omitempty"`
	TestPlan           []string `yaml:"test_plan,omitempty" json:"test_plan,omitempty"`
}

// TaskList represents a collection of tasks, typically the output of a decompose stage.
type TaskList struct {
	Tasks []Task `yaml:"tasks" json:"tasks"`
}

// Parse attempts to read a TaskList from a string (YAML or JSON).
func Parse(content string) (TaskList, error) {
	var tl TaskList

	// Try YAML first (which also handles JSON)
	err := yaml.Unmarshal([]byte(content), &tl)
	if err == nil && len(tl.Tasks) > 0 {
		return tl, nil
	}

	// If it fails or returns no tasks, try finding a code block if it's Markdown
	if strings.Contains(content, "```") {
		extracted := extractCodeBlock(content)
		if extracted != "" {
			err = yaml.Unmarshal([]byte(extracted), &tl)
			if err == nil && len(tl.Tasks) > 0 {
				return tl, nil
			}
		}
	}

	// Try direct JSON unmarshal for safety
	err = json.Unmarshal([]byte(content), &tl)
	if err == nil && len(tl.Tasks) > 0 {
		return tl, nil
	}

	return TaskList{}, fmt.Errorf("failed to parse task list from content")
}

// extractCodeBlock tries to extract the content inside a markdown code block.
func extractCodeBlock(content string) string {
	start := strings.Index(content, "```")
	if start == -1 {
		return ""
	}

	// Find the end of the opening backticks line
	firstLineEnd := strings.Index(content[start:], "\n")
	if firstLineEnd == -1 {
		return ""
	}
	start += firstLineEnd + 1

	end := strings.Index(content[start:], "```")
	if end == -1 {
		// No closing backticks, take the rest of the string
		return strings.TrimSpace(content[start:])
	}

	return strings.TrimSpace(content[start : start+end])
}

// Validate checks for duplicate IDs and common errors in the task list.
func (tl TaskList) Validate() []string {
	var warnings []string
	seen := make(map[string]bool)

	for i, t := range tl.Tasks {
		taskRef := t.ID
		if taskRef == "" {
			taskRef = fmt.Sprintf("at index %d", i)
			warnings = append(warnings, "task with empty ID found "+taskRef)
		} else {
			if seen[t.ID] {
				warnings = append(warnings, fmt.Sprintf("duplicate task ID: %s", t.ID))
			}
			seen[t.ID] = true
		}

		if t.Title == "" {
			warnings = append(warnings, fmt.Sprintf("task %s has empty title", taskRef))
		}
		if t.Goal == "" {
			warnings = append(warnings, fmt.Sprintf("task %s has empty goal", taskRef))
		}
	}

	visitedGlobal := make(map[string]bool)
	for _, t := range tl.Tasks {
		for _, dep := range t.Dependencies {
			if !seen[dep] {
				warnings = append(warnings, fmt.Sprintf("task %s depends on unknown task %s", t.ID, dep))
			}
		}

		if t.ID != "" {
			if err := checkTaskCycles(tl, t.ID, []string{}, visitedGlobal); err != nil {
				warnings = append(warnings, fmt.Sprintf("circular dependency: %v", err))
			}
		}
	}

	return warnings
}

func checkTaskCycles(tl TaskList, current string, path []string, visitedGlobal map[string]bool) error {
	if visitedGlobal[current] {
		return nil // Already verified safe
	}

	for _, visited := range path {
		if visited == current {
			return fmt.Errorf("%s -> %s", strings.Join(path, " -> "), current)
		}
	}
	path = append(path, current)

	var targetTask *Task
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == current {
			targetTask = &tl.Tasks[i]
			break
		}
	}

	if targetTask != nil {
		for _, dep := range targetTask.Dependencies {
			if err := checkTaskCycles(tl, dep, path, visitedGlobal); err != nil {
				return err
			}
		}
	}
	visitedGlobal[current] = true
	return nil
}
