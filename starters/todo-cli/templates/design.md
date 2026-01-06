# Design: {{ .ProjectName }}

## Requirements
{{- range .Inputs }}
{{- if eq .Name "requirements.md" }}
{{ .Content }}
{{- end }}
{{- end }}

## Task
Design the system. Focus on the internal data structures and the command-line interface implementation.


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
