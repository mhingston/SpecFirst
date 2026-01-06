# {{ .StageName }} — {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
Create a lightweight visual record of the designed screen(s) for the section.

## Output Requirements
Create `product/sections/<section-id>/<name>.png` screenshots for any completed views.

## Notes
- If the environment cannot render UI to image, output a placeholder note explaining how to capture screenshots manually.


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
