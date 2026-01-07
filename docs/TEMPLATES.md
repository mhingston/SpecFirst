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
| `StageType` | string | The type of the current stage (`spec`, `decompose`, etc.). |
| `Prompt` | PromptConfig | The detailed prompt configuration for the stage. |
| `OutputContract` | OutputContract | The expected structure of the stage output. |

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

## Template Functions

Templates have access to these helper functions:

| Function | Usage | Description |
| --- | --- | --- |
| `join` | `{{ join .Outputs ", " }}` | Joins a slice of strings with a delimiter |
| `readFile` | `{{ readFile "skill.md" }}` | Includes a file from `.specfirst/skills/` |

## Skills (Reusable Prompt Chunks)

Skills are reusable markdown files stored in `.specfirst/skills/`. They help share common guidance across templates without copy-pasting.

### Example: Including a Skill

```markdown
# {{ .StageName }}

## Design Guidelines
{{ readFile "design-principles.md" }}

## Task
Design the user interface...
```

This includes the content of `.specfirst/skills/design-principles.md` directly in the rendered prompt.

### Creating Skills

Create any `.md` file in `.specfirst/skills/` and reference it via `{{ readFile "your-skill.md" }}`.

See the [product-design example](../starters/product-design/) for a working demonstration of skills.

**Security Note**: The `readFile` function only reads from `.specfirst/skills/` and rejects paths containing `..` or absolute paths.

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
