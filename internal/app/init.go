package app

import (
	"fmt"
	"os"
	"path/filepath"

	"specfirst/internal/assets"
	"specfirst/internal/domain"
	"specfirst/internal/repository"
	"specfirst/internal/starter"
	"specfirst/internal/utils"
)

// InitOptions defines parameters for workspace initialization.
type InitOptions struct {
	Starter string
	Force   bool
}

// InitializeWorkspace sets up a new SpecFirst workspace.
func InitializeWorkspace(opts InitOptions) error {
	// Create workspace directories
	dirs := []string{
		repository.SpecPath(),
		repository.ArtifactsPath(),
		repository.GeneratedPath(),
		repository.ProtocolsPath(),
		repository.TemplatesPath(),
		repository.ArchivesPath(),
	}
	for _, dir := range dirs {
		if err := utils.EnsureDir(dir); err != nil {
			return err
		}
	}

	// Determine protocol name
	protocolName := assets.DefaultProtocolName
	if opts.Starter != "" {
		protocolName = opts.Starter
	}

	// Write config first
	configPath := repository.ConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		projectName := filepath.Base(repository.BaseDir())
		cfg := fmt.Sprintf(assets.DefaultConfigTemplate, projectName, protocolName)
		if err := os.WriteFile(configPath, []byte(cfg), 0644); err != nil {
			return err
		}
	}

	// Apply starter if selected
	if opts.Starter != "" {
		if err := starter.Apply(opts.Starter, opts.Force, true); err != nil {
			return fmt.Errorf("applying starter %q: %w", opts.Starter, err)
		}
	} else {
		// Default behavior: write default protocol and templates
		if err := writeIfMissing(repository.ProtocolsPath(assets.DefaultProtocolName+".yaml"), assets.DefaultProtocolYAML); err != nil {
			return err
		}
		if err := writeIfMissing(repository.TemplatesPath("clarify.md"), assets.ClarifyTemplate); err != nil {
			return err
		}
		if err := writeIfMissing(repository.TemplatesPath("design.md"), assets.DesignTemplate); err != nil {
			return err
		}
		if err := writeIfMissing(repository.TemplatesPath("implementation.md"), assets.ImplementationTemplate); err != nil {
			return err
		}
		if err := writeIfMissing(repository.TemplatesPath("decompose.md"), assets.DecomposeTemplate); err != nil {
			return err
		}
	}

	// Write state file using NewState
	statePath := repository.StatePath()
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		s := domain.NewState(protocolName)
		if err := repository.SaveState(statePath, s); err != nil {
			return err
		}
	}

	return nil
}

func writeIfMissing(path string, content string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}
