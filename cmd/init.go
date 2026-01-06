package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"specfirst/internal/assets"
	"specfirst/internal/starter"
	"specfirst/internal/state"
	"specfirst/internal/store"
)

var (
	initStarter string
	initChoose  bool
	initForce   bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a SpecFirst workspace",
	Long: `Initialize a SpecFirst workspace in the current directory.

Creates .specfirst/ with default protocol, templates, and config.

Options:
  --starter <name>  Initialize with a specific starter kit workflow
  --choose          Interactively select a starter kit
  --force           Overwrite existing templates/protocols (only with --starter or --choose)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create workspace directories
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

		// Handle starter selection
		selectedStarter := initStarter
		if initChoose {
			chosen, err := interactiveSelectStarter(cmd)
			if err != nil {
				return err
			}
			selectedStarter = chosen
		}

		// Determine protocol name
		protocolName := assets.DefaultProtocolName
		if selectedStarter != "" {
			protocolName = selectedStarter
		} else if protocolFlag != "" {
			protocolName = protocolFlag
		}

		// Write config first
		configPath := store.ConfigPath()
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			projectName := filepath.Base(mustGetwd())
			cfg := fmt.Sprintf(assets.DefaultConfigTemplate, projectName, protocolName)
			if err := os.WriteFile(configPath, []byte(cfg), 0644); err != nil {
				return err
			}
		}

		// Apply starter if selected
		if selectedStarter != "" {
			// Apply the starter (copies protocol, templates, skills)
			// We pass updateConfig=false because we already wrote the correct config above
			// BUT starter.Apply might have extra logic for defaults.
			// Actually, starter.Apply expects to be able to update config.
			// Let's let it update if it needs to, but we started with a valid one.
			if err := starter.Apply(selectedStarter, initForce, true); err != nil {
				return fmt.Errorf("applying starter %q: %w", selectedStarter, err)
			}
		} else {
			// Default behavior: write default protocol and templates
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
		}

		// Write state file using NewState
		statePath := store.StatePath()
		if _, err := os.Stat(statePath); os.IsNotExist(err) {
			// Initialize with the determined protocol
			// Check if we can load it to get the real name from config if starter changed it?
			// But for now, using the name we derived is safer than empty string.
			s := state.NewState(protocolName)
			if err := state.Save(statePath, s); err != nil {
				return err
			}
		}

		if selectedStarter != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Initialized .specfirst workspace with starter %q\n", selectedStarter)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "Initialized .specfirst workspace")
		}
		return nil
	},
}

// interactiveSelectStarter prompts the user to select a starter.
func interactiveSelectStarter(cmd *cobra.Command) (string, error) {
	starters, err := starter.List()
	if err != nil {
		return "", err
	}

	if len(starters) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No starters found. Using default protocol.")
		return "", nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Available starters:")
	fmt.Fprintln(cmd.OutOrStdout(), "  0) [default] Use default multi-stage protocol")
	for i, s := range starters {
		fmt.Fprintf(cmd.OutOrStdout(), "  %d) %s\n", i+1, s.Name)
	}

	fmt.Fprint(cmd.OutOrStdout(), "\nSelect a starter [0]: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" || input == "0" {
		return "", nil // Use default
	}

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(starters) {
		return "", fmt.Errorf("invalid selection: %s", input)
	}

	return starters[choice-1].Name, nil
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

func init() {
	initCmd.Flags().StringVar(&initStarter, "starter", "", "initialize with a specific starter kit")
	initCmd.Flags().BoolVar(&initChoose, "choose", false, "interactively select a starter kit")
	initCmd.Flags().BoolVar(&initForce, "force", false, "overwrite existing templates/protocols")
}
