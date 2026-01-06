# {{ .StageName }} — {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
For ONE section spec, generate:
- realistic sample data (`data.json`)
- TypeScript interfaces (`types.ts`) including:
  - entity interfaces
  - props interface for the main component
  - callback signatures for actions in the spec
 
## Guidelines
- **Output Formatting**: For this task, **provide the raw content only**. Skip all preamble, conversational filler, and Markdown code fences.

### Section selection
- Prefer `CustomVars.section_id` if present: {{ index .CustomVars "section_id" }}
- Otherwise infer from the included section spec.

## Output Requirements

### `product/sections/<section-id>/data.json`
- Include 5–10 realistic records per primary entity (where applicable).
- Use plausible names, dates, statuses, amounts, etc. (no lorem ipsum).
- Include edge cases (empty arrays, very long text, "archived" item, etc.).
- You MAY include a top-level `_meta` explaining fields, but keep it small.

### `product/sections/<section-id>/types.ts`
- Export types for each entity present in `data.json`.
- Export a `<SectionName>Props` interface with:
  - data props (arrays/objects)
  - UI state props (selectedId, query, etc.) only if needed by flows
  - callback props for actions (`onView`, `onEdit`, `onDelete`, `onCreate`, etc.)
- Keep types framework-agnostic and reusable.


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
