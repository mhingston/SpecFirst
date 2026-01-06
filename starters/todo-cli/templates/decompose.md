# Decompose: {{ .ProjectName }}

## Context
{{- range .Inputs }}
### {{ .Name }}
{{ .Content }}
{{- end }}

## Task
Break the design down into at most 5 implementation tasks.
Output must be a valid YAML file `tasks.yaml`.
Each task needs: id, title, goal, dependencies, files_touched, acceptance_criteria, and test_plan.


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
