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
