package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"specfirst/internal/protocol"
	"specfirst/internal/store"
	tmplpkg "specfirst/internal/template"
	"specfirst/internal/workspace"
)

type CompileOptions struct {
	Granularity    string
	MaxTasks       int
	PreferParallel bool
	RiskBias       string
}

func (e *Engine) RequireStageDependencies(stage protocol.Stage) error {
	for _, dep := range stage.DependsOn {
		if !e.State.IsStageCompleted(dep) {
			return fmt.Errorf("missing dependency: %s", dep)
		}
	}
	return nil
}

func (e *Engine) CompilePrompt(stage protocol.Stage, stageIDs []string, opts CompileOptions) (string, error) {
	var inputs []tmplpkg.Input
	if stage.Intent == "review" {
		artifacts, err := e.ListAllArtifacts()
		if err != nil {
			return "", err
		}
		inputs = artifacts
	} else {
		inputs = make([]tmplpkg.Input, 0, len(stage.Inputs))
		for _, input := range stage.Inputs {
			path, err := workspace.ArtifactPathForInput(input, stage.DependsOn, stageIDs)
			if err != nil {
				return "", err
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return "", err
			}
			inputs = append(inputs, tmplpkg.Input{Name: input, Content: string(content)})
		}
	}

	// Apply options overrides
	if opts.Granularity != "" || opts.MaxTasks > 0 || opts.PreferParallel || opts.RiskBias != "" {
		if stage.Prompt == nil {
			stage.Prompt = &protocol.PromptConfig{}
		}
		// Deep copy prompt config if we are mutating it?
		// Better to just override in the data object passed to template, but template uses stage.Prompt
		// For now, mutation of the local stage copy (passed by value? No, stage is struct, passed by value, but Prompt is pointer).
		// We should clone Prompt config to avoid mutating the cached protocol structure if that matters.
		// Since we load protocol fresh in engine usually it might be fine, but safe to clone.
		p := *stage.Prompt // shallow copy of struct
		if opts.Granularity != "" {
			p.Granularity = opts.Granularity
		}
		if opts.MaxTasks > 0 {
			p.MaxTasks = opts.MaxTasks
		}
		if opts.PreferParallel {
			p.PreferParallel = opts.PreferParallel
		}
		if opts.RiskBias != "" {
			p.RiskBias = opts.RiskBias
		}
		stage.Prompt = &p
	}

	data := tmplpkg.Data{
		StageName:   stage.Name,
		ProjectName: e.Config.ProjectName,
		Inputs:      inputs,
		Outputs:     stage.Outputs,
		Intent:      stage.Intent,
		Language:    e.Config.Language,
		Framework:   e.Config.Framework,
		CustomVars:  e.Config.CustomVars,
		Constraints: e.Config.Constraints,

		StageType:      stage.Type,
		Prompt:         stage.Prompt,
		OutputContract: stage.Output,
	}

	templatePath := store.TemplatesPath(stage.Template)
	return tmplpkg.Render(templatePath, data)
}

func (e *Engine) ListAllArtifacts() ([]tmplpkg.Input, error) {
	artifactsRoot := store.ArtifactsPath()
	info, err := os.Stat(artifactsRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return []tmplpkg.Input{}, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("artifacts path is not a directory: %s", artifactsRoot)
	}
	relPaths := []string{}
	err = filepath.WalkDir(artifactsRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(artifactsRoot, path)
		if err != nil {
			return err
		}
		relPaths = append(relPaths, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(relPaths)

	inputs := make([]tmplpkg.Input, 0, len(relPaths))
	for _, rel := range relPaths {
		data, err := os.ReadFile(filepath.Join(artifactsRoot, rel))
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, tmplpkg.Input{Name: rel, Content: string(data)})
	}
	return inputs, nil
}
