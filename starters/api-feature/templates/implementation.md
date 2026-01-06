# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Task
Implement the code for this specific task from the task breakdown.

## Guidelines
- **Output Formatting**: For this implementation task, **provide the raw content/diff only**. Skip all preamble, conversational filler, and Markdown code blocks (unless the output *is* a Markdown file). Only provide reasoning if specifically asked or if there is a critical ambiguity.
- Follow the API design specs exactly
- Implement only what's in this task's scope
- Add tests for your changes
- Handle errors gracefully
- Add logging/monitoring as specified

## Output Format
Provide code as:
- Complete files (for new files)
- Unified diffs (for modifications)
- Test files
- Any configuration changes

## Verification
After implementation:
1. Run unit tests
2. Test API endpoint manually (if applicable)
3. Verify error cases
4. Check logs/metrics

## Assumptions
- Design has been approved
- Database schema exists (if needed)
- Dependencies are available
