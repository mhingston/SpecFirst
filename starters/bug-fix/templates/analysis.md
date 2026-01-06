# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- if .Inputs }}
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}
{{- else }}
(No prior context)
{{- end }}

## Task
Analyze the reported bug and create a structured analysis that identifies:
1. **Root cause** - What is causing the bug?
2. **Impact** - Who/what is affected?
3. **Reproduction steps** - How to trigger the bug reliably
4. **Proposed fix** - High-level approach to fixing it
5. **Risks** - What could break if we fix this?
 
## Guidelines
- **Output Formatting**: For this analysis, **provide the raw content only**. Skip all preamble, conversational filler, and Markdown code fences (the output *is* a Markdown file).
- Be thorough but concise
- Use technical language appropriate for a root cause analysis
- Surface hidden risks even if they seem unlikely

## Output Requirements
Create `analysis.md` with the following structure:

### Bug Report
- **Title**: Brief description
- **Severity**: Critical / High / Medium / Low
- **Reported By**: Name/source
- **Date**: When reported

### Root Cause
Explain what's causing the bug technically.

### Impact Assessment
- Who is affected?
- How often does this occur?
- What is the business/user impact?

### Reproduction Steps
Clear steps to reproduce the bug consistently.

### Proposed Fix
High-level approach to fixing without implementation details yet.

### Risks & Side Effects
What could potentially break or need testing?

### Dependencies
Any blockers or prerequisites for the fix?

## Assumptions
- The bug has been reported by users or identified in testing
- We have access to reproduction steps or environment
- (List any other assumptions explicitly)


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
