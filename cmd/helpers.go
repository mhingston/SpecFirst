package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"specfirst/internal/assets"
	"specfirst/internal/config"
	"specfirst/internal/prompt"
	"specfirst/internal/protocol"
	"specfirst/internal/state"
	"specfirst/internal/store"
	tmplpkg "specfirst/internal/template"
	"specfirst/internal/workspace"
)

func loadConfig() (config.Config, error) {
	cfg, err := config.Load(store.ConfigPath())
	if err != nil {
		return config.Config{}, err
	}
	if cfg.Protocol == "" {
		cfg.Protocol = assets.DefaultProtocolName
	}
	if cfg.ProjectName == "" {
		if wd, err := os.Getwd(); err == nil {
			cfg.ProjectName = filepath.Base(wd)
		} else {
			cfg.ProjectName = "project" // Fallback when working directory is unavailable
		}
	}
	if cfg.CustomVars == nil {
		cfg.CustomVars = map[string]string{}
	}
	if cfg.Constraints == nil {
		cfg.Constraints = map[string]string{}
	}
	return cfg, nil
}

func activeProtocolName(cfg config.Config) string {
	if protocolFlag != "" {
		return protocolFlag
	}
	return cfg.Protocol
}

func loadProtocol(name string) (protocol.Protocol, error) {
	// If name looks like a path or has .yaml extension, load it directly
	if filepath.IsAbs(name) || strings.Contains(name, string(os.PathSeparator)) || strings.HasSuffix(name, ".yaml") {
		return protocol.Load(filepath.Clean(name))
	}
	// Otherwise treat as a protocol name in the protocols directory
	path := store.ProtocolsPath(name + ".yaml")
	proto, err := protocol.Load(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Check if .specfirst directory exists at all
			specDir := store.SpecPath()
			if _, statErr := os.Stat(specDir); os.IsNotExist(statErr) {
				return protocol.Protocol{}, fmt.Errorf("project not initialized. Run 'specfirst init' in %s", store.BaseDir())
			}
			return protocol.Protocol{}, fmt.Errorf("protocol %q not found in %s", name, specDir)
		}
		return protocol.Protocol{}, err
	}
	return proto, nil
}

func loadState() (state.State, error) {
	s, err := state.Load(store.StatePath())
	if err != nil {
		return state.State{}, err
	}
	if s.CompletedStages == nil {
		s.CompletedStages = []string{}
	}
	if s.StageOutputs == nil {
		s.StageOutputs = map[string]state.StageOutput{}
	}
	if s.Attestations == nil {
		s.Attestations = make(map[string][]state.Attestation)
	}
	return s, nil
}

func saveState(s state.State) error {
	return state.Save(store.StatePath(), s)
}

func ensureDir(path string) error {
	return workspace.EnsureDir(path)
}

func copyFile(src, dst string) error {
	return workspace.CopyFile(src, dst)
}

func copyDir(src, dst string) error {
	return workspace.CopyDir(src, dst)
}

func copyDirWithOpts(src, dst string, required bool) error {
	return workspace.CopyDirWithOpts(src, dst, required)
}

func outputRelPath(output string) (string, error) {
	return workspace.ProjectRelPath(output)
}

func resolveOutputPath(output string) (string, error) {
	rel, err := outputRelPath(output)
	if err != nil {
		return "", err
	}
	baseDir := repoRoot()
	if baseDir == "" {
		// Fallback if repoRoot fails (shouldn't happen if outputRelPath worked)
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		// Try to find root again or default to WD
		if root, err := store.FindProjectRoot(wd); err == nil {
			baseDir = root
		} else {
			baseDir = wd
		}
	}
	return filepath.Join(baseDir, rel), nil
}

func artifactPathForInput(filename string, priorityStages []string, stageIDs []string) (string, error) {
	return workspace.ArtifactPathForInput(filename, priorityStages, stageIDs)
}

func artifactRelFromState(value string) (string, error) {
	return workspace.ArtifactRelFromState(value)
}

func artifactAbsFromState(value string) (string, error) {
	return workspace.ArtifactAbsFromState(value)
}

func isWindowsAbs(value string) bool {
	if len(value) >= 3 && value[1] == ':' && value[2] == '/' {
		return true
	}
	return strings.HasPrefix(value, "//")
}

func relFromArtifactsPath(abs string) (string, bool) {
	clean := filepath.Clean(abs)
	parts := strings.Split(clean, string(os.PathSeparator))
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == store.ArtifactsDir {
			if i+1 >= len(parts) {
				return "", false
			}
			rel := filepath.Join(parts[i+1:]...)
			if rel == "" || rel == "." {
				return "", false
			}
			return rel, true
		}
	}
	return "", false
}

func compilePrompt(stage protocol.Stage, cfg config.Config, stageIDs []string) (string, error) {
	var inputs []tmplpkg.Input
	if stage.Intent == "review" {
		artifacts, err := listAllArtifacts()
		if err != nil {
			return "", err
		}
		inputs = artifacts
	} else {
		inputs = make([]tmplpkg.Input, 0, len(stage.Inputs))
		for _, input := range stage.Inputs {
			path, err := artifactPathForInput(input, stage.DependsOn, stageIDs)
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

	// Override PromptConfig with CLI flags if set
	if stageGranularity != "" || stageMaxTasks > 0 || stagePreferParallel || stageRiskBias != "balanced" {
		if stage.Prompt == nil {
			stage.Prompt = &protocol.PromptConfig{}
		}
		if stageGranularity != "" {
			stage.Prompt.Granularity = stageGranularity
		}
		if stageMaxTasks > 0 {
			stage.Prompt.MaxTasks = stageMaxTasks
		}
		if stagePreferParallel {
			stage.Prompt.PreferParallel = stagePreferParallel
		}
		if stageRiskBias != "balanced" || stage.Prompt.RiskBias == "" {
			stage.Prompt.RiskBias = stageRiskBias
		}
	}

	data := tmplpkg.Data{
		StageName:   stage.Name,
		ProjectName: cfg.ProjectName,
		Inputs:      inputs,
		Outputs:     stage.Outputs,
		Intent:      stage.Intent,
		Language:    cfg.Language,
		Framework:   cfg.Framework,
		CustomVars:  cfg.CustomVars,
		Constraints: cfg.Constraints,

		StageType:      stage.Type,
		Prompt:         stage.Prompt,
		OutputContract: stage.Output,
		Epistemics:     nil, // Epistemics not available in this context (static check)
	}

	templatePath := store.TemplatesPath(stage.Template)
	return tmplpkg.Render(templatePath, data)
}

func promptHash(prompt string) string {
	return workspace.PromptHash(prompt)
}

func fileHash(path string) (string, error) {
	return workspace.FileHash(path)
}

func writeOutput(path string, content string) error {
	return workspace.WriteOutput(path, content)
}

func formatPrompt(format string, stageID string, promptString string) (string, error) {
	return prompt.Format(format, stageID, promptString)
}

func applyMaxChars(promptString string, maxChars int) string {
	return prompt.ApplyMaxChars(promptString, maxChars)
}

func renderInlineTemplate(tmpl string, data any) (string, error) {
	return tmplpkg.RenderInline(tmpl, data)
}

func collectFileHashes(root string) (map[string]string, error) {
	return workspace.CollectFileHashes(root)
}

func listAllArtifacts() ([]tmplpkg.Input, error) {
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

func requireStageDependencies(s state.State, stage protocol.Stage) error {
	for _, dep := range stage.DependsOn {
		if !s.IsStageCompleted(dep) {
			return fmt.Errorf("missing dependency: %s", dep)
		}
	}
	return nil
}

func hasApproval(records []state.Attestation, role string) bool {
	for _, record := range records {
		if record.Role == role && record.Status == "approved" {
			return true
		}
	}
	return false
}

func missingApprovals(p protocol.Protocol, s state.State) []string {
	missing := []string{}
	for _, approval := range p.Approvals {
		if !s.IsStageCompleted(approval.Stage) {
			continue
		}
		if !hasApproval(s.Attestations[approval.Stage], approval.Role) {
			missing = append(missing, fmt.Sprintf("%s:%s", approval.Stage, approval.Role))
		}
	}
	sort.Strings(missing)
	return missing
}

func ensureStateInitialized(s state.State, proto protocol.Protocol) state.State {
	if s.Protocol == "" {
		s.Protocol = proto.Name
	}
	if s.SpecVersion == "" {
		s.SpecVersion = proto.Version
	}
	if s.StartedAt.IsZero() {
		s.StartedAt = time.Now().UTC()
	}
	if s.CompletedStages == nil {
		s.CompletedStages = []string{}
	}
	if s.StageOutputs == nil {
		s.StageOutputs = map[string]state.StageOutput{}
	}
	if s.Attestations == nil {
		s.Attestations = map[string][]state.Attestation{}
	}
	return s
}

func stageIDList(proto protocol.Protocol) []string {
	ids := make([]string, 0, len(proto.Stages))
	for _, stage := range proto.Stages {
		ids = append(ids, stage.ID)
	}
	return ids
}

func normalizeMatchPath(value string) string {
	return workspace.NormalizeMatchPath(value)
}

func matchOutputPattern(pattern string, file string) bool {
	return workspace.MatchOutputPattern(pattern, file)
}

func heredocDelimiter(prompt string) string {
	base := "SPECFIRST_EOF"
	delimiter := base
	for i := 0; ; i++ {
		if i > 0 {
			delimiter = fmt.Sprintf("%s_%d", base, i)
		}
		line := delimiter + "\n"
		if !strings.Contains(prompt, "\n"+delimiter+"\n") &&
			!strings.HasPrefix(prompt, line) &&
			!strings.HasSuffix(prompt, "\n"+delimiter) &&
			prompt != delimiter {
			return delimiter
		}
	}
}

func repoRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	root, err := store.FindProjectRoot(wd)
	if err != nil {
		return ""
	}
	return root
}
