# {{ .StageName }} — {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
Break the product into **3–5 buildable sections**.

Each section should be:
- A navigation destination (or a first-class "area")
- Independently designable/buildable
- Sequenced in a sensible build order

## Output Requirements
Create `product/product-roadmap.md` with:

- Ordered list of sections:
  - `id` (kebab-case)
  - `title`
  - 2–4 sentence description
  - key user flows (bullets)
  - dependencies on other sections (if any)

## Constraints
- Keep to 3–5 sections unless the input strongly demands more.
- No implementation tasks here; just product structure.


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
