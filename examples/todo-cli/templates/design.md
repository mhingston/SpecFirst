# Design: {{ .ProjectName }}

## Requirements
{{- range .Inputs }}
{{- if eq .Name "requirements.md" }}
{{ .Content }}
{{- end }}
{{- end }}

## Task
Design the system. Focus on the internal data structures and the command-line interface implementation.
