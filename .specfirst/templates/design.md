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
Produce a system design that satisfies the requirements. Make binding decisions and note trade-offs.

{{- if .Outputs }}
## Output Requirements
{{- range .Outputs }}
- {{ . }}
{{- end }}
{{- end }}
