package protocol

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Protocol represents a workflow definition with stages and approvals.
type Protocol struct {
	Name      string      `yaml:"name"`
	Version   string      `yaml:"version"`
	Uses      []string    `yaml:"uses,omitempty"` // Protocol imports/mixins
	Stages    []Stage     `yaml:"stages"`
	Approvals []Approval  `yaml:"approvals"`
	Lint      *LintConfig `yaml:"lint,omitempty"` // Protocol-level schema additions
}

// Stage represents a workflow step with optional type, modifiers, and contracts.
type Stage struct {
	ID        string   `yaml:"id"`
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type,omitempty"` // "spec", "decompose", "task_prompt"
	Intent    string   `yaml:"intent"`
	Template  string   `yaml:"template"`
	DependsOn []string `yaml:"depends_on"`
	Inputs    []string `yaml:"inputs"`
	Outputs   []string `yaml:"outputs"`

	// Stage modifiers
	Optional   bool `yaml:"optional,omitempty"`
	Repeatable bool `yaml:"repeatable,omitempty"`
	Terminal   bool `yaml:"terminal,omitempty"`

	// Prompt configuration
	Prompt *PromptConfig `yaml:"prompt,omitempty"`

	// Output contract
	Output *OutputContract `yaml:"output,omitempty"`

	// For task_prompt type - reference to decompose stage
	Source string `yaml:"source,omitempty"`

	// Ambiguity Gates
	MaxOpenQuestions        *int     `yaml:"max_open_questions,omitempty"`
	MustResolveTags         []string `yaml:"must_resolve_tags,omitempty"`
	MaxHighRisksUnmitigated *int     `yaml:"max_high_risks_unmitigated,omitempty"`
}

// PromptConfig defines prompt generation parameters.
type PromptConfig struct {
	Intent            string      `yaml:"intent,omitempty"`
	ExpectedOutput    string      `yaml:"expected_output,omitempty"`
	Determinism       string      `yaml:"determinism,omitempty"`        // high, medium, low
	AllowedCreativity string      `yaml:"allowed_creativity,omitempty"` // high, medium, low
	Granularity       string      `yaml:"granularity,omitempty"`        // feature, story, ticket, commit
	MaxTasks          int         `yaml:"max_tasks,omitempty"`
	PreferParallel    bool        `yaml:"prefer_parallel,omitempty"`
	RiskBias          string      `yaml:"risk_bias,omitempty"` // conservative, balanced, fast
	Rules             []string    `yaml:"rules,omitempty"`
	RequiredFields    []string    `yaml:"required_fields,omitempty"`
	Lint              *LintConfig `yaml:"lint,omitempty"` // Stage-level schema additions
}

// LintConfig defines additional validation rules for prompts.
type LintConfig struct {
	RequiredSections []string `yaml:"required_sections,omitempty"`
	ForbiddenPhrases []string `yaml:"forbidden_phrases,omitempty"`
}

// OutputContract defines expected output structure.
type OutputContract struct {
	Format         string   `yaml:"format,omitempty"` // markdown, yaml, json
	Sections       []string `yaml:"sections,omitempty"`
	RequiredFields []string `yaml:"required_fields,omitempty"`
}

type Approval struct {
	Role  string `yaml:"role"`
	Stage string `yaml:"stage"`
}

var validStageIDPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// validStageTypes defines allowed stage types.
var validStageTypes = map[string]bool{
	"":            true, // Empty defaults to "spec"
	"spec":        true,
	"decompose":   true,
	"task_prompt": true,
}

// Load reads and validates a protocol from a YAML file.
// It resolves protocol imports (uses) and validates all stage references.
func Load(path string) (Protocol, error) {
	return LoadWithResolver(path, nil)
}

// LoadWithResolver loads a protocol with a custom resolver for imports.
// If resolver is nil, imports are resolved relative to the protocol's directory.
func LoadWithResolver(path string, resolver func(name string) (Protocol, error)) (Protocol, error) {
	// visitedStack tracks the current import chain to detect cycles
	// processedCache tracks already loaded files to avoid reloading and handle diamond deps
	return loadWithDepth(path, resolver, make([]string, 0), make(map[string]Protocol))
}

func loadWithDepth(path string, resolver func(name string) (Protocol, error), visitedStack []string, processedCache map[string]Protocol) (Protocol, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return Protocol{}, err
	}

	// 1. Cycle Detection (Check current stack)
	for _, visited := range visitedStack {
		if visited == abs {
			return Protocol{}, fmt.Errorf("circular protocol import detected: %s", abs)
		}
	}

	// 2. Memoization (Check if already processed)
	if cached, ok := processedCache[abs]; ok {
		return cached, nil
	}

	// 3. Load and Parse
	data, err := os.ReadFile(path)
	if err != nil {
		return Protocol{}, err
	}

	var p Protocol
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Protocol{}, err
	}

	// 4. Resolve Imports
	if len(p.Uses) > 0 {
		baseDir := filepath.Dir(path)
		// Push current path to stack for recursive calls
		newStack := append(visitedStack, abs)

		// Iterate backwards so that later imports in the list take precedence over earlier ones.
		// Since we prepend imported stages, processing the last import *first* puts its stages
		// "deepest" in the prepended list.
		// Wait.
		// Uses = [Base, Override]
		// Reverse Loop:
		// 1. Override. Stages = [Override, Local]
		// 2. Base. Stages = [Base, Override, Local]
		// Dedupe (Backwards): Local, Override, (Base ignored). -> Override wins.
		for i := len(p.Uses) - 1; i >= 0; i-- {
			importPath := p.Uses[i]
			var imported Protocol
			var importErr error

			if resolver != nil {
				// Resolver handles its own recursion/caching/paths if needed,
				// but typically we'd pass the stack/cache if the resolver allows.
				// For the default simple resolver pattern here, we assume it recurses back to us or similar.
				// To keep it simple for the default case:
				imported, importErr = resolver(importPath)
			} else {
				// Default: resolve relative to protocol directory
				importFile := filepath.Join(baseDir, importPath+".yaml")
				imported, importErr = loadWithDepth(importFile, nil, newStack, processedCache)
			}

			if importErr != nil {
				return Protocol{}, fmt.Errorf("importing %q: %w", importPath, importErr)
			}

			// Prepend imported stages (so they can be depended upon)
			p.Stages = append(imported.Stages, p.Stages...)
			// Merge approvals
			p.Approvals = append(imported.Approvals, p.Approvals...)
		}
	}

	// Deduplicate Stages
	// Strategy: Keep the LAST occurrence of a stage ID.
	// This allows local protocol to override imports, and later imports to override earlier ones.
	// It also safely handles diamond dependencies where the same stage is imported multiple times.
	uniqueStages := make(map[string]bool)
	var dedupedStages []Stage
	// Iterate backwards to find the "winner" (most local/latest) first
	for i := len(p.Stages) - 1; i >= 0; i-- {
		s := p.Stages[i]
		if !uniqueStages[s.ID] {
			uniqueStages[s.ID] = true
			// Prepend to maintain relative order of the "winners"
			dedupedStages = append([]Stage{s}, dedupedStages...)
		}
	}
	p.Stages = dedupedStages

	// Deduplicate Approvals (Key: Stage + Role)
	uniqueApprovals := make(map[string]bool)
	var dedupedApprovals []Approval
	for i := len(p.Approvals) - 1; i >= 0; i-- {
		a := p.Approvals[i]
		key := a.Stage + "::" + a.Role
		if !uniqueApprovals[key] {
			uniqueApprovals[key] = true
			dedupedApprovals = append([]Approval{a}, dedupedApprovals...)
		}
	}
	p.Approvals = dedupedApprovals

	// Validate stage ID uniqueness
	seen := make(map[string]bool)
	for _, stage := range p.Stages {
		if err := validateStageID(stage.ID); err != nil {
			return Protocol{}, err
		}
		if err := validateTemplatePath(stage.Template); err != nil {
			return Protocol{}, err
		}
		if err := validateStageType(stage.Type); err != nil {
			return Protocol{}, fmt.Errorf("stage %q: %w", stage.ID, err)
		}
		if seen[stage.ID] {
			return Protocol{}, fmt.Errorf("duplicate stage ID: %s", stage.ID)
		}
		seen[stage.ID] = true
	}

	// Validate depends_on references point to existing stages and no self-references
	stageMap := make(map[string]Stage)
	for _, stage := range p.Stages {
		stageMap[stage.ID] = stage
		for _, dep := range stage.DependsOn {
			if dep == stage.ID {
				return Protocol{}, fmt.Errorf("stage %q cannot depend on itself", stage.ID)
			}
			if !seen[dep] {
				return Protocol{}, fmt.Errorf("stage %q depends on unknown stage %q", stage.ID, dep)
			}
		}
	}

	// Validate inputs against outputs of dependencies
	for _, stage := range p.Stages {
		for _, input := range stage.Inputs {
			if input == "" {
				continue
			}
			found := false

			// First, try to match input against outputs of DependsOn stages
			// This works for both simple filenames and multi-level paths
			for _, depID := range stage.DependsOn {
				depStage := stageMap[depID]
				for _, out := range depStage.Outputs {
					if matchPattern(out, input) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}

			// If not found and input has a path separator, try stage-qualified format
			// Format: <stage-id>/<output-path> (e.g., "requirements/readme.md")
			if !found && (strings.Contains(input, "/") || strings.Contains(input, string(os.PathSeparator))) {
				parts := strings.SplitN(filepath.ToSlash(input), "/", 2)
				stageID := parts[0]
				filename := parts[1]

				// Only try stage-qualified if the first component looks like a stage ID
				if targetStage, ok := stageMap[stageID]; ok {
					for _, out := range targetStage.Outputs {
						if matchPattern(out, filename) {
							found = true
							break
						}
					}
				}
			}

			if !found {
				var depInfo []string
				for _, depID := range stage.DependsOn {
					depStage := stageMap[depID]
					depInfo = append(depInfo, fmt.Sprintf("%s outputs: %v", depID, depStage.Outputs))
				}
				return Protocol{}, fmt.Errorf("stage %q: input %q not found in outputs of any dependency. Dependencies: %v", stage.ID, input, depInfo)
			}
		}
	}

	// Validate task_prompt source references
	for _, stage := range p.Stages {
		if stage.Type == "task_prompt" && stage.Source != "" {
			if !seen[stage.Source] {
				return Protocol{}, fmt.Errorf("stage %q references unknown source stage %q", stage.ID, stage.Source)
			}
			// Verify source is a decompose stage
			sourceStage, _ := p.StageByID(stage.Source)
			if sourceStage.Type != "decompose" {
				return Protocol{}, fmt.Errorf("stage %q source must be a decompose stage, got %q", stage.ID, sourceStage.Type)
			}
		}
	}

	// Detect circular dependencies
	for _, stage := range p.Stages {
		if err := checkCycles(p, stage.ID, []string{}); err != nil {
			return Protocol{}, err
		}
	}

	// Validate approval stage references
	for _, approval := range p.Approvals {
		if strings.TrimSpace(approval.Stage) == "" {
			return Protocol{}, fmt.Errorf("approval references empty stage")
		}
		if strings.TrimSpace(approval.Role) == "" {
			return Protocol{}, fmt.Errorf("approval role is required for stage %q", approval.Stage)
		}
		if !seen[approval.Stage] {
			return Protocol{}, fmt.Errorf("approval references unknown stage %q", approval.Stage)
		}
	}

	// Cache the fully resolved protocol
	processedCache[abs] = p
	return p, nil
}

func (p Protocol) StageByID(id string) (Stage, bool) {
	for _, stage := range p.Stages {
		if stage.ID == id {
			return stage, true
		}
	}
	return Stage{}, false
}

func checkCycles(p Protocol, current string, path []string) error {
	for _, visited := range path {
		if visited == current {
			return fmt.Errorf("circular dependency detected: %v -> %s", path, current)
		}
	}
	path = append(path, current)

	stage, _ := p.StageByID(current)
	for _, dep := range stage.DependsOn {
		if err := checkCycles(p, dep, path); err != nil {
			return err
		}
	}
	return nil
}

func validateStageID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("stage id is required")
	}
	if id != strings.ToLower(id) {
		return fmt.Errorf("stage id %q must be lowercase", id)
	}
	if !validStageIDPattern.MatchString(id) {
		return fmt.Errorf("invalid stage id %q (must be alphanumeric and may contain hyphens or underscores)", id)
	}
	return nil
}

func matchPattern(pattern, file string) bool {
	if pattern == "" {
		return false
	}
	if strings.Contains(pattern, "*") {
		ok, _ := filepath.Match(pattern, file)
		return ok
	}
	return pattern == file
}

func validateTemplatePath(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("stage template is required")
	}
	clean := filepath.Clean(filepath.FromSlash(value))
	if filepath.IsAbs(clean) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("invalid template path: %s", value)
	}
	return nil
}

func validateStageType(stageType string) error {
	if !validStageTypes[stageType] {
		return fmt.Errorf("invalid stage type %q (must be one of: spec, decompose, task_prompt)", stageType)
	}
	return nil
}
