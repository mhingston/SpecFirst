package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	stageFormat      string
	stageOut         string
	stageMaxChars    int
	stageNoStrict    bool
	stageInteractive bool

	stageGranularity    string
	stageMaxTasks       int
	stagePreferParallel bool
	stageRiskBias       string

	protocolFlag string
)

var rootCmd = &cobra.Command{
	Use:     "specfirst",
	Short:   "SpecFirst CLI for specification-driven workflows",
	Version: version,
	Args:    cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if stageInteractive {
			if len(args) > 0 {
				return errors.New("interactive mode does not accept a stage id")
			}
			return runInteractive(cmd.OutOrStdout())
		}
		if len(args) == 0 {
			return errors.New("missing command or stage id")
		}
		stageID := args[0]
		return runStage(cmd.OutOrStdout(), stageID)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&protocolFlag, "protocol", "", "override protocol (name or path)")
	rootCmd.PersistentFlags().StringVar(&stageFormat, "format", "text", "output format: text, json, yaml, or shell")
	rootCmd.PersistentFlags().StringVar(&stageOut, "out", "", "write compiled prompt to a file")
	rootCmd.PersistentFlags().IntVar(&stageMaxChars, "max-chars", 0, "truncate output to max chars")
	rootCmd.PersistentFlags().BoolVar(&stageNoStrict, "no-strict", false, "bypass dependency gating")
	rootCmd.Flags().BoolVar(&stageInteractive, "interactive", false, "generate interactive meta-prompt")

	rootCmd.PersistentFlags().StringVar(&stageGranularity, "granularity", "", "task granularity: feature, story, ticket, commit")
	rootCmd.PersistentFlags().IntVar(&stageMaxTasks, "max-tasks", 0, "maximum number of tasks to generate")
	rootCmd.PersistentFlags().BoolVar(&stagePreferParallel, "prefer-parallel", false, "prefer parallelizable tasks")
	rootCmd.PersistentFlags().StringVar(&stageRiskBias, "risk-bias", "balanced", "risk bias: conservative, balanced, fast")

	rootCmd.ValidArgsFunction = stageIDCompletions

	_ = rootCmd.RegisterFlagCompletionFunc("protocol", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return filterPrefix(loadProtocolNames(), toComplete), cobra.ShellCompDirectiveDefault
	})
	_ = rootCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return filterPrefix([]string{"text", "json", "yaml", "shell"}, toComplete), cobra.ShellCompDirectiveNoFileComp
	})
	_ = rootCmd.RegisterFlagCompletionFunc("granularity", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return filterPrefix([]string{"feature", "story", "ticket", "commit"}, toComplete), cobra.ShellCompDirectiveNoFileComp
	})
	_ = rootCmd.RegisterFlagCompletionFunc("risk-bias", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return filterPrefix([]string{"conservative", "balanced", "fast"}, toComplete), cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(starterCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(completeSpecCmd)
	rootCmd.AddCommand(lintCmd)
	rootCmd.AddCommand(protocolCmd)
	rootCmd.AddCommand(archiveCmd)
	rootCmd.AddCommand(attestCmd)
	rootCmd.AddCommand(checkCmd)

	// Cognitive scaffold commands
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(assumptionsCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(failureCmd)
	rootCmd.AddCommand(testIntentCmd)
	rootCmd.AddCommand(traceCmd)
	rootCmd.AddCommand(distillCmd)
	rootCmd.AddCommand(calibrateCmd)

	// Epistemic Ledger commands
	rootCmd.AddCommand(assumeCmd)
	rootCmd.AddCommand(questionCmd)
	rootCmd.AddCommand(decisionCmd)
	rootCmd.AddCommand(riskCmd)
	rootCmd.AddCommand(disputeCmd)
}
