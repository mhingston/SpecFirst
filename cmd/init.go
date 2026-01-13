package cmd

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"specfirst/internal/app"
	"specfirst/internal/starter"
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
		// Handle starter selection
		selectedStarter := initStarter
		if initChoose {
			chosen, err := interactiveSelectStarter(cmd.InOrStdin(), cmd.OutOrStdout())
			if err != nil {
				return err
			}
			selectedStarter = chosen
		}

		opts := app.InitOptions{
			Starter: selectedStarter,
			Force:   initForce,
		}

		if err := app.InitializeWorkspace(opts); err != nil {
			return err
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
// Sev3 Fix: Use io.Reader/Writer for testability instead of direct os.Stdin/Stdout.
func interactiveSelectStarter(in io.Reader, out io.Writer) (string, error) {
	starters, err := starter.List()
	if err != nil {
		return "", err
	}

	if len(starters) == 0 {
		fmt.Fprintln(out, "No starters found. Using default protocol.")
		return "", nil
	}

	fmt.Fprintln(out, "Available starters:")
	fmt.Fprintln(out, "  0) [default] Use default multi-stage protocol")
	for i, s := range starters {
		fmt.Fprintf(out, "  %d) %s\n", i+1, s.Name)
	}

	fmt.Fprint(out, "\nSelect a starter [0]: ")

	reader := bufio.NewReader(in)
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

func init() {
	initCmd.Flags().StringVar(&initStarter, "starter", "", "initialize with a specific starter kit")
	initCmd.Flags().BoolVar(&initChoose, "choose", false, "interactively select a starter kit")
	initCmd.Flags().BoolVar(&initForce, "force", false, "overwrite existing templates/protocols")

	_ = initCmd.RegisterFlagCompletionFunc("starter", starterNameCompletions)
}
