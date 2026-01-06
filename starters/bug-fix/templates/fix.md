# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Task
Implement the bug fix based on the analysis above.

Generate the minimal code changes required to:
1. Fix the root cause
2. Prevent regression
3. Handle edge cases identified in the analysis

## Implementation Guidelines
- **Output Formatting**: For this fix, **provide the raw content/diff only**. Skip preamble, conversation, and Markdown code fences.
- **Minimize scope**: Only change what's necessary
- **Add tests**: Cover the bug scenario and edge cases
- **Preserve behavior**: Don't introduce new features
- **Document changes**: Add comments explaining the fix

## Output Format
Provide code changes as unified diffs with:
- File paths
- Line numbers
- Clear diff markers (+/-)
- Explanatory comments

##Expected Output
Produce the actual code files or diffs needed:
{{- if .Outputs }}
{{- range .Outputs }}
- {{ . }}
{{- end }}
{{- else }}
- Code files with fixes applied
- Test files (if writing new tests)
{{- end }}

## Verification Plan
Describe how to verify the fix works:
1. Run existing tests
2. Test reproduction steps from analysis
3. Test edge cases
4. Smoke test related functionality

## Assumptions
- Analysis has identified the correct root cause
- Test environment is available
- (List any other assumptions explicitly)


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
