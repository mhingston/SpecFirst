package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"specfirst/internal/assets"
	"specfirst/internal/store"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a SpecFirst workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureDir(store.SpecPath()); err != nil {
			return err
		}
		if err := ensureDir(store.ArtifactsPath()); err != nil {
			return err
		}
		if err := ensureDir(store.GeneratedPath()); err != nil {
			return err
		}
		if err := ensureDir(store.ProtocolsPath()); err != nil {
			return err
		}
		if err := ensureDir(store.TemplatesPath()); err != nil {
			return err
		}
		if err := ensureDir(store.ArchivesPath()); err != nil {
			return err
		}

		if err := writeIfMissing(store.ProtocolsPath(assets.DefaultProtocolName+".yaml"), assets.DefaultProtocolYAML); err != nil {
			return err
		}
		if err := writeIfMissing(store.TemplatesPath("requirements.md"), assets.RequirementsTemplate); err != nil {
			return err
		}
		if err := writeIfMissing(store.TemplatesPath("design.md"), assets.DesignTemplate); err != nil {
			return err
		}
		if err := writeIfMissing(store.TemplatesPath("implementation.md"), assets.ImplementationTemplate); err != nil {
			return err
		}
		if err := writeIfMissing(store.TemplatesPath("decompose.md"), assets.DecomposeTemplate); err != nil {
			return err
		}

		configPath := store.ConfigPath()
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			projectName := filepath.Base(mustGetwd())
			protoName := assets.DefaultProtocolName
			if protocolFlag != "" {
				protoName = protocolFlag
			}
			cfg := fmt.Sprintf(assets.DefaultConfigTemplate, projectName, protoName)
			if err := os.WriteFile(configPath, []byte(cfg), 0644); err != nil {
				return err
			}
		}

		statePath := store.StatePath()
		if _, err := os.Stat(statePath); os.IsNotExist(err) {
			if err := os.WriteFile(statePath, []byte("{}\n"), 0644); err != nil {
				return err
			}
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Initialized .specfirst workspace")
		return nil
	},
}

func writeIfMissing(path string, content string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "project"
	}
	return wd
}
