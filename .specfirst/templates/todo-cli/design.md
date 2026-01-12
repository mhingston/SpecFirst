# Design: {{ .ProjectName }}

## Requirements
{{- range .Inputs }}
{{- if eq .Name "requirements.md" }}
{{ .Content }}
{{- end }}
{{- end }}

## Requirements Contract
You must follow the clarified requirements exactly.

- Treat In Scope and Out of Scope / Non-Goals as hard boundaries.
- Use Acceptance Criteria as the definition of done.
- Respect all stated Constraints.
- Treat listed Assumptions as true unless explicitly challenged.

If you believe the requirements are incorrect, incomplete, or unsafe:
- Stop.
- Explain the issue.
- Propose a revision to requirements.md.
- Do NOT silently override requirements.

## Task
Design the system. Focus on the internal data structures and the command-line interface implementation.


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
