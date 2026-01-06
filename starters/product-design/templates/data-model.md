# {{ .StageName }} — {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
Define a **global data model**: the core entities ("nouns") and how they relate.

## Output Requirements
Create `product/data-model/data-model.md` with:

- Entities list (5–15 typical)
  - name
  - description
  - key fields (with types in plain English)
  - ownership / lifecycle notes
- Relationships section (1:many, many:many, etc.)
- Shared enums / statuses (if any)
- Notes on invariants (e.g., uniqueness, required fields)
- Explicit non-goals / unknowns

## Guidelines
- This model is used to keep sections consistent.
- Avoid DB/ORM specifics; think "product and API shape".


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
