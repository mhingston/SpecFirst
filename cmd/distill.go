package cmd

import (
	"fmt"
	"os"
	"specfirst/internal/prompts"
	"strings"

	"github.com/spf13/cobra"
)

var distillAudience string

var distillCmd = &cobra.Command{
	Use:   "distill <spec-file>",
	Short: "Generate an audience-specific spec summary prompt",
	Long: `Generate a prompt that produces different cognitive views of the same specification.

Different audiences need different spec representations:
- exec: Business impact, risk, timeline
- implementer: Technical details, constraints, interfaces
- ai: Constraint-dense summary with invariants and decision boundaries
- qa: Test surface, acceptance criteria, edge cases

The output is a structured prompt suitable for AI assistants or human writers.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		specPath := args[0]

		content, err := os.ReadFile(specPath)
		if err != nil {
			return fmt.Errorf("reading spec %s: %w", specPath, err)
		}

		prompt, err := generateDistillPrompt(string(content), specPath, distillAudience)
		if err != nil {
			return err
		}

		prompt = applyMaxChars(prompt, stageMaxChars)
		formatted, err := formatPrompt(stageFormat, "distill", prompt)
		if err != nil {
			return err
		}

		if stageOut != "" {
			if err := writeOutput(stageOut, formatted); err != nil {
				return err
			}
		}
		_, err = cmd.OutOrStdout().Write([]byte(formatted))
		return err
	},
}

func init() {
	distillCmd.Flags().StringVar(&distillAudience, "audience", "ai",
		"Target audience: exec, implementer, ai, qa")
}

var distillAudiences = map[string]string{
	"exec": `Summarize this specification for an **executive audience**.

Focus on:

## Business Impact
- What problem does this solve?
- What is the expected value/ROI?
- What competitive advantage does it provide?

## Risk Assessment
- What are the biggest risks?
- What is the probability and impact of each?
- What mitigation strategies are proposed?

## Timeline & Resources
- What is the estimated effort?
- What dependencies could cause delays?
- What team/skills are required?

## Success Criteria
- How will we know this succeeded?
- What metrics will we track?
- What is the rollback plan?

Keep the summary to 1-2 pages. Avoid technical jargon.
Use bullet points for scanability.`,

	"implementer": `Summarize this specification for an **implementation-focused developer**.

Focus on:

## Technical Requirements
- What exact functionality must be built?
- What interfaces must be exposed?
- What data structures are involved?

## Constraints
- What performance requirements exist?
- What compatibility requirements exist?
- What security requirements must be met?

## Dependencies
- What external services are required?
- What internal APIs are consumed?
- What libraries or frameworks are mandated?

## Edge Cases
- What error handling is specified?
- What boundary conditions are defined?
- What undefined behavior exists?

## Implementation Hints
- What architecture patterns are suggested?
- What trade-offs have already been decided?
- What is explicitly out of scope?

Be precise and actionable. Include specific values, formats, and constraints.`,

	"qa": `Summarize this specification for a **QA engineer**.

Focus on:

## Test Surface
- What features need to be tested?
- What user journeys are critical?
- What integration points exist?

## Acceptance Criteria
For each feature:
- What does "done" look like?
- What is the expected behavior?
- What are the success metrics?

## Edge Cases & Boundaries
- What input boundaries exist?
- What error conditions should be tested?
- What race conditions are possible?

## Test Data Requirements
- What test data is needed?
- What test environments are required?
- What mocks or stubs are necessary?

## Regression Risks
- What existing functionality could break?
- What areas need extra coverage?
- What has historically been buggy?

Organize by testability. Flag areas that are hard to test.`,
}

func generateDistillPrompt(content, path, audience string) (string, error) {
	audienceLower := strings.ToLower(strings.TrimSpace(audience))

	// AI audience uses the gold-standard template
	if audienceLower == "ai" {
		return prompts.Render("ai-distillation.md", prompts.SpecData{Spec: content})
	}

	audienceInstructions, ok := distillAudiences[audienceLower]
	if !ok {
		validAudiences := make([]string, 0, len(distillAudiences)+1)
		validAudiences = append(validAudiences, "ai")
		for k := range distillAudiences {
			validAudiences = append(validAudiences, k)
		}
		return "", fmt.Errorf("unknown audience %q; valid options: %s", audience, strings.Join(validAudiences, ", "))
	}

	return fmt.Sprintf(`%s

---

## Specification
**Source**: %s

%s
`, audienceInstructions, path, content), nil
}
