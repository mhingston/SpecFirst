# {{ .StageName }} — {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Design Quality Skill
{{- readFile "design-principles.md" -}}

## Task
Implement the **section screen** as production-grade React components.

### Hard requirements (do not violate)
- Exportable components MUST accept all data via props (do NOT import `data.json` inside exportable components).
- Exportable components MUST accept callback props for all user actions from the section spec.
- Implement ALL user flows + UI requirements from the section spec.
- Use Tailwind utility classes.
- Support light + dark mode via `dark:` variants.
- Apply design tokens from `colors.json` and `typography.json`.
- Do not re-implement global navigation inside section screens (shell owns navigation).

### What to create
- Exportable component(s) under: `src/sections/<section-id>/components/`
- A preview wrapper screen under: `src/sections/<section-id>/<ViewName>.tsx`
  - Preview wrapper MAY import `product/sections/<id>/data.json` for demo purposes
  - Preview wrapper MUST pass data into exportable components as props

### Multiple views
If the section spec implies multiple views (list + detail + form):
- Build the primary view first unless `CustomVars.view_name` specifies otherwise.
- Prefer list/dashboard as primary if ambiguous.
 
## Guidelines
- **Output Formatting**: For this task, **provide the raw code only**. Skip all preamble, conversational filler, and Markdown code fences.

## Output Requirements
Create files matching:
{{- range .Outputs }}
- {{ . }}
{{- end }}


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
