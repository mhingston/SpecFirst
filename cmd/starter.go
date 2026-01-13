package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"specfirst/internal/starter"
)

var (
	starterForce    bool
	starterNoConfig bool
	starterSource   string
)

var starterCmd = &cobra.Command{
	Use:   "starter",
	Short: "Manage starter kit workflows",
	Long: `Manage starter kits - pre-built workflow bundles with protocols and templates.

Use 'starter list' to see available starters.
Use 'starter apply <name>' to install a starter into the current workspace.`,
}

var starterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available starter kits",
	RunE: func(cmd *cobra.Command, args []string) error {
		starters, err := starter.List()
		if err != nil {
			return err
		}

		// Filter by source
		var filtered []starter.Starter
		for _, s := range starters {
			if starterSource == "all" ||
				(starterSource == "builtin" && s.IsBuiltin) ||
				(starterSource == "local" && !s.IsBuiltin) {
				filtered = append(filtered, s)
			}
		}

		if len(filtered) == 0 {
			if starterSource != "all" {
				fmt.Fprintf(cmd.OutOrStdout(), "No starters found for source %q.\n", starterSource)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "No starters found.")
			}
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tSOURCE\tSKILLS\tDEFAULTS\tDESCRIPTION")
		for _, s := range filtered {
			source := "local"
			if s.IsBuiltin {
				source = "builtin"
			}
			hasSkills := "-"
			if s.SkillsDir != "" {
				hasSkills = "✓"
			}
			hasDefaults := "-"
			if s.DefaultsPath != "" {
				hasDefaults = "✓"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", s.Name, source, hasSkills, hasDefaults, s.Description)
		}
		w.Flush()

		return nil
	},
}

var starterApplyCmd = &cobra.Command{
	Use:   "apply <name>",
	Short: "Apply a starter kit to the current workspace",
	Long: `Apply a starter kit to the current workspace.

This copies the starter's protocol to .specfirst/protocols/<name>.yaml
and templates to .specfirst/templates/.

By default, existing files are not overwritten. Use --force to overwrite.
By default, config.yaml is updated to use the new protocol. Use --no-config to skip.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: starterNameCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Check if workspace exists
		if _, err := os.Stat(".specfirst"); os.IsNotExist(err) {
			return fmt.Errorf("no .specfirst workspace found; run 'specfirst init' first")
		}

		updateConfig := !starterNoConfig
		if err := starter.Apply(name, starterForce, updateConfig); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Applied starter %q\n", name)
		if updateConfig {
			fmt.Fprintf(cmd.OutOrStdout(), "Updated config.yaml to use protocol: %s\n", name)
		}
		if !starterForce {
			fmt.Fprintln(cmd.OutOrStdout(), "(Existing files were preserved. Use --force to overwrite.)")
		}

		return nil
	},
}

func init() {
	starterCmd.AddCommand(starterListCmd)
	starterCmd.AddCommand(starterApplyCmd)

	starterListCmd.Flags().StringVar(&starterSource, "source", "all", "filter by source: all, builtin, local")

	starterApplyCmd.Flags().BoolVar(&starterForce, "force", false, "overwrite existing templates and protocols")
	starterApplyCmd.Flags().BoolVar(&starterNoConfig, "no-config", false, "do not update config.yaml with new protocol")

	_ = starterListCmd.RegisterFlagCompletionFunc("source", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return filterPrefix([]string{"all", "builtin", "local"}, toComplete), cobra.ShellCompDirectiveNoFileComp
	})
}
