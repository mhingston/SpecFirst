package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"specfirst/internal/assets"
	"specfirst/internal/repository"
	"specfirst/internal/utils"
)

// validProtocolNamePattern matches safe protocol names (alphanumeric, hyphens, underscores)
var validProtocolNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

var protocolCmd = &cobra.Command{
	Use:   "protocol",
	Short: "Manage protocols",
}

var protocolListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available protocols",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := os.ReadDir(repository.ProtocolsPath())
		if err != nil {
			return err
		}
		protocols := []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if strings.HasSuffix(name, ".yaml") {
				protocols = append(protocols, strings.TrimSuffix(name, ".yaml"))
			}
		}
		sort.Strings(protocols)
		for _, name := range protocols {
			fmt.Fprintln(cmd.OutOrStdout(), name)
		}
		return nil
	},
}

var protocolShowCmd = &cobra.Command{
	Use:               "show <name>",
	Short:             "Show a protocol definition",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: protocolNameCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		// Validate protocol name to prevent path traversal
		if !validProtocolNamePattern.MatchString(name) {
			return fmt.Errorf("invalid protocol name: %q (must start with letter, contain only letters, numbers, hyphens, underscores)", name)
		}
		path := repository.ProtocolsPath(name + ".yaml")
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("protocol not found: %s", name)
			}
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	},
}

var protocolCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a protocol from the default template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if name == "" {
			return fmt.Errorf("protocol name is required")
		}
		// Validate protocol name with comprehensive checks
		if len(name) > 64 {
			return fmt.Errorf("protocol name too long (max 64 characters)")
		}
		if !validProtocolNamePattern.MatchString(name) {
			return fmt.Errorf("invalid protocol name: %q (must start with letter, contain only letters, numbers, hyphens, underscores)", name)
		}
		path := repository.ProtocolsPath(name + ".yaml")
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("protocol already exists: %s", name)
		}
		if err := utils.EnsureDir(filepath.Dir(path)); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(assets.DefaultProtocolYAML), 0644); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Created protocol %s\n", name)
		return nil
	},
}

func init() {
	protocolCmd.AddCommand(protocolListCmd)
	protocolCmd.AddCommand(protocolShowCmd)
	protocolCmd.AddCommand(protocolCreateCmd)
}
