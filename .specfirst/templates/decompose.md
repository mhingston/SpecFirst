# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- if .Inputs }}
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}
{{- end }}

## Task
You are decomposing an approved design into implementation tasks.

### Rules
- Tasks must be independently reviewable.
- If assumptions are required, state them explicitly.
- If any acceptance criteria cannot be verified, add a task to define verification.

## Output Requirements
- Format: YAML
- Each task must include: id, title, goal, acceptance_criteria, dependencies, files_touched, risk_level, estimated_scope, test_plan
