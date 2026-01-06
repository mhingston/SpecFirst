package cmd

import (
	"fmt"
	"strings"

	"specfirst/internal/state"

	"github.com/spf13/cobra"
)

var (
	epistemicOwner   string
	epistemicTags    string
	epistemicContext string
)

// -- Assume --

var assumeCmd = &cobra.Command{
	Use:   "assume",
	Short: "Manage assumptions",
}

var assumeAddCmd = &cobra.Command{
	Use:   "add [text]",
	Short: "Add a new assumption",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		id := s.AddAssumption(args[0], epistemicOwner)
		if err := state.Save(".specfirst/state.json", s); err != nil {
			return err
		}
		fmt.Printf("Added assumption %s\n", id)
		return nil
	},
}

var assumeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List assumptions",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		for _, a := range s.Epistemics.Assumptions {
			fmt.Printf("[%s] %s (%s)\n", a.ID, a.Text, a.Status)
		}
		return nil
	},
}

// -- Question --

var questionCmd = &cobra.Command{
	Use:   "question",
	Short: "Manage open questions",
}

var questionAddCmd = &cobra.Command{
	Use:   "add [text]",
	Short: "Add a new open question",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		tags := []string{}
		if epistemicTags != "" {
			tags = strings.Split(epistemicTags, ",")
		}
		id := s.AddOpenQuestion(args[0], tags, epistemicContext)
		if err := state.Save(".specfirst/state.json", s); err != nil {
			return err
		}
		fmt.Printf("Added question %s\n", id)
		return nil
	},
}

var questionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List open questions",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		for _, q := range s.Epistemics.OpenQuestions {
			fmt.Printf("[%s] %s (Tags: %v, Status: %s)\n", q.ID, q.Text, q.Tags, q.Status)
		}
		return nil
	},
}

// -- Decision --

var decisionCmd = &cobra.Command{
	Use:   "decision",
	Short: "Manage decisions",
}

var decisionAddCmd = &cobra.Command{
	Use:   "add [text]",
	Short: "Record a decision",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		// In a real implementation we might want flags for rationale/alternatives
		id := s.AddDecision(args[0], "No rationale provided via CLI yet", nil)
		if err := state.Save(".specfirst/state.json", s); err != nil {
			return err
		}
		fmt.Printf("Recorded decision %s\n", id)
		return nil
	},
}

// -- Risk --

var riskCmd = &cobra.Command{
	Use:   "risk",
	Short: "Manage risks",
}

var riskAddCmd = &cobra.Command{
	Use:   "add [text] [severity]",
	Short: "Add a risk",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		severity := "medium"
		if len(args) > 1 {
			severity = args[1]
		}
		id := s.AddRisk(args[0], severity)
		if err := state.Save(".specfirst/state.json", s); err != nil {
			return err
		}
		fmt.Printf("Added risk %s\n", id)
		return nil
	},
}

// -- Dispute --
var disputeCmd = &cobra.Command{
	Use:   "dispute",
	Short: "Manage disputes",
}

var disputeAddCmd = &cobra.Command{
	Use:   "add [topic]",
	Short: "Log a dispute",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		id := s.AddDispute(args[0])
		if err := state.Save(".specfirst/state.json", s); err != nil {
			return err
		}
		fmt.Printf("Logged dispute %s\n", id)
		return nil
	},
}

func init() {
	// Assume
	assumeAddCmd.Flags().StringVar(&epistemicOwner, "owner", "", "Owner of the assumption")
	assumeCmd.AddCommand(assumeAddCmd)
	assumeCmd.AddCommand(assumeListCmd)

	// Question
	questionAddCmd.Flags().StringVar(&epistemicTags, "tags", "", "Comma-separated tags")
	questionAddCmd.Flags().StringVar(&epistemicContext, "context", "", "Context reference")
	questionCmd.AddCommand(questionAddCmd)
	questionCmd.AddCommand(questionListCmd)

	// Decision
	decisionCmd.AddCommand(decisionAddCmd)

	// Risk
	riskCmd.AddCommand(riskAddCmd)

	// Dispute
	disputeCmd.AddCommand(disputeAddCmd)

	// Lifecycle
	assumeCmd.AddCommand(assumeCloseCmd)
	questionCmd.AddCommand(questionResolveCmd)
	decisionCmd.AddCommand(decisionUpdateCmd)
	riskCmd.AddCommand(riskMitigateCmd)
	disputeCmd.AddCommand(disputeResolveCmd)
}

// Lifecycle Commands

var assumeCloseCmd = &cobra.Command{
	Use:   "close [id]",
	Short: "Close an assumption",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		status, _ := cmd.Flags().GetString("status")
		if status == "" {
			return fmt.Errorf("status is required")
		}
		if !s.CloseAssumption(args[0], status) {
			return fmt.Errorf("assumption %s not found", args[0])
		}
		return state.Save(".specfirst/state.json", s)
	},
}

var questionResolveCmd = &cobra.Command{
	Use:   "resolve [id]",
	Short: "Resolve an open question",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		answer, _ := cmd.Flags().GetString("answer")
		if answer == "" {
			return fmt.Errorf("answer is required")
		}
		if !s.ResolveOpenQuestion(args[0], answer) {
			return fmt.Errorf("question %s not found", args[0])
		}
		return state.Save(".specfirst/state.json", s)
	},
}

var decisionUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update decision status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		status, _ := cmd.Flags().GetString("status")
		if status == "" {
			return fmt.Errorf("status is required")
		}
		if !s.UpdateDecision(args[0], status) {
			return fmt.Errorf("decision %s not found", args[0])
		}
		return state.Save(".specfirst/state.json", s)
	},
}

var riskMitigateCmd = &cobra.Command{
	Use:   "mitigate [id]",
	Short: "Mitigate a risk",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		mitigation, _ := cmd.Flags().GetString("mitigation")
		status, _ := cmd.Flags().GetString("status")
		if mitigation == "" {
			return fmt.Errorf("mitigation is required")
		}
		if status == "" {
			status = "mitigated"
		}
		if !s.MitigateRisk(args[0], mitigation, status) {
			return fmt.Errorf("risk %s not found", args[0])
		}
		return state.Save(".specfirst/state.json", s)
	},
}

var disputeResolveCmd = &cobra.Command{
	Use:   "resolve [id]",
	Short: "Resolve a dispute",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load(".specfirst/state.json")
		if err != nil {
			return err
		}
		if !s.ResolveDispute(args[0]) {
			return fmt.Errorf("dispute %s not found", args[0])
		}
		return state.Save(".specfirst/state.json", s)
	},
}

func init() {
	// Flags for lifecycle
	assumeCloseCmd.Flags().String("status", "validated", "Status: validated|invalidated")
	questionResolveCmd.Flags().String("answer", "", "Answer to the question")
	decisionUpdateCmd.Flags().String("status", "accepted", "New status")
	riskMitigateCmd.Flags().String("mitigation", "", "Mitigation plan")
	riskMitigateCmd.Flags().String("status", "mitigated", "New status")
}
