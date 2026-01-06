# {{ .StageName }} — {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
Write a section spec for ONE section.

### Section selection
- Prefer `CustomVars.section_id` if present: {{ index .CustomVars "section_id" }}
- Otherwise, pick the **next** section in roadmap that does not yet have a spec.

## Output Requirements
Create `product/sections/<section-id>/spec.md` with:

# <Section Title>

## Overview
2–3 sentences.

## In Shell?
Yes/No. (Default Yes.)

## Primary User Flows
List 3–8 flows with step-by-step bullets.

## UI Requirements
Be specific about:
- Layout patterns (list/detail, dashboard, form, wizard)
- Data density expectations
- Filters/sort/search, pagination, bulk actions
- Empty/loading/error states
- Modals/drawers/toasts if relevant

## Objects & Data Needed
List what the UI needs (entity names/fields in plain language). Reference the global data model terms when possible.

## Permissions / Roles (if any)
What different users can see/do.

## Out of Scope
Explicit exclusions.

## Open Questions
Only if truly unknown.

## Notes
- Do NOT design visuals here; define behavior and interface requirements.
- Avoid backend implementation details.


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
