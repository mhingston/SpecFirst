package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	texttmpl "text/template"
	"time"

	"specfirst/internal/assets"
	"specfirst/internal/config"
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
	return protocol.Load(path)
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
	if s.Approvals == nil {
		s.Approvals = map[string][]state.Approval{}
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
	if output == "" {
		return "", fmt.Errorf("output path is empty")
	}
	clean := filepath.Clean(output)
	baseDir := repoRoot()
	if baseDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		baseDir = wd
	}
	baseEval, err := filepath.EvalSymlinks(baseDir)
	if err == nil {
		baseDir = baseEval
	}
	abs := clean
	if !filepath.IsAbs(abs) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		abs = filepath.Join(wd, clean)
	}
	if absEval, err := filepath.EvalSymlinks(abs); err == nil {
		abs = absEval
	} else if dirEval, err := filepath.EvalSymlinks(filepath.Dir(abs)); err == nil {
		abs = filepath.Join(dirEval, filepath.Base(abs))
	}
	rel, err := filepath.Rel(baseDir, abs)
	if err != nil {
		return "", err
	}
	clean = rel
	if clean == "." {
		return "", fmt.Errorf("output path resolves to current directory")
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("output path escapes workspace: %s", output)
	}
	return clean, nil
}

func resolveOutputPath(output string) (string, error) {
	rel, err := outputRelPath(output)
	if err != nil {
		return "", err
	}
	baseDir := repoRoot()
	if baseDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		baseDir = wd
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
	}

	templatePath := store.TemplatesPath(stage.Template)
	return tmplpkg.Render(templatePath, data)
}

func promptHash(prompt string) string {
	hash := sha256.Sum256([]byte(prompt))
	return "sha256:" + hex.EncodeToString(hash[:])
}

func fileHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(hash[:]), nil
}

func writeOutput(path string, content string) error {
	if err := ensureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func formatPrompt(format string, stageID string, prompt string) (string, error) {
	if format == "text" {
		return prompt, nil
	}
	if format == "json" {
		payload := map[string]string{
			"stage":  stageID,
			"prompt": prompt,
		}
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data) + "\n", nil
	}
	if format == "yaml" {
		// Simple YAML output without external dependencies
		escapedPrompt := strings.ReplaceAll(prompt, "\\", "\\\\")
		escapedPrompt = strings.ReplaceAll(escapedPrompt, "\"", "\\\"")
		return fmt.Sprintf("stage: %s\nprompt: |\n  %s\n", stageID, strings.ReplaceAll(escapedPrompt, "\n", "\n  ")), nil
	}
	if format == "shell" {
		// Escape single quotes for shell safety
		escapedStageID := strings.ReplaceAll(stageID, "'", "'\"'\"'")
		delimiter := heredocDelimiter(prompt)
		return fmt.Sprintf("SPECFIRST_STAGE='%s'\nSPECFIRST_PROMPT=$(cat <<'%s'\n%s\n%s\n)\n", escapedStageID, delimiter, prompt, delimiter), nil
	}
	return "", errors.New("unsupported format: " + format)
}

func applyMaxChars(prompt string, maxChars int) string {
	if maxChars <= 0 {
		return prompt
	}
	// Use runes to avoid truncating mid-UTF8 character
	runes := []rune(prompt)
	if len(runes) <= maxChars {
		return prompt
	}
	return string(runes[:maxChars])
}

func renderInlineTemplate(tmpl string, data any) (string, error) {
	parsed, err := texttmpl.New("inline").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := parsed.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func collectFileHashes(root string) (map[string]string, error) {
	files := map[string]string{}
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", root)
	}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		hash, err := fileHash(path)
		if err != nil {
			return err
		}
		files[rel] = hash
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
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

func hasApproval(records []state.Approval, role string) bool {
	for _, record := range records {
		if record.Role == role {
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
		if !hasApproval(s.Approvals[approval.Stage], approval.Role) {
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
	if s.Approvals == nil {
		s.Approvals = map[string][]state.Approval{}
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
	lines, err := gitCmd("rev-parse", "--show-toplevel")
	if err != nil {
		return ""
	}
	if len(lines) == 0 {
		return ""
	}
	root := strings.TrimSpace(lines[0])
	if root == "" {
		return ""
	}
	return root
}
