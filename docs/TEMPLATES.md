# Template Reference

SpecFirst uses Go's `text/template` package to render prompts. Templates are the primary way to define the instructions sent to an LLM.

## Template Context

The following variables are available in the top-level template context (`.`):

| Variable | Type | Description |
| --- | --- | --- |
| `StageName` | string | Name of the current stage. |
| `ProjectName` | string | Name of the project (from config). |
| `Inputs` | []Input | List of artifacts from previous stages. |
| `Outputs` | []string | List of declared output filenames. |
| `Intent` | string | The semantic intent of the stage. |
| `Language` | string | Project language (from config). |
| `Framework` | string | Project framework (from config). |
| `CustomVars` | map[string]string | User-defined variables. |
| `Constraints` | map[string]string | Project constraints. |

### Input Object
Each item in `Inputs` has:
- `Name`: Filename of the artifact.
- `Content`: Full text content of the artifact.

## Common Patterns

### Embedding Artifacts
Use a range loop to include previous work as context:

```markdown
## Context
{{- range .Inputs }}
### {{ .Name }}
{{ .Content }}
{{- end }}
```

### Conditional Sections
Only show a section if constraints are defined:

```markdown
{{- if .Constraints }}
## Constraints
{{- range $key, $value := .Constraints }}
- {{ $key }}: {{ $value }}
{{- end }}
{{- end }}
```

### Whitespace Control
Use `{{-` and `-}}` to remove leading/trailing whitespace and prevent extra blank lines in your rendered prompts.

## Example Template

```markdown
# Implementation Prompt for {{ .ProjectName }}

## Requirements
{{- range .Inputs }}
{{- if eq .Name "requirements.md" }}
{{ .Content }}
{{- end }}
{{- end }}

## Task
Implement the following files:
{{- range .Outputs }}
- {{ . }}
{{- end }}
```
