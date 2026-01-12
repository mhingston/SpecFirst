# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- if .Inputs }}
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}
{{- else }}
(No prior artifacts)
{{- end }}

## Task
Gather requirements, ask clarifying questions, and enumerate constraints.

{{- if .Outputs }}
## Output Requirements
{{- range .Outputs }}
- {{ . }}
{{- end }}
{{- end }}
