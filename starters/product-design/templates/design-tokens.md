# {{ .StageName }} — {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
Choose a lightweight design token set to be used across all designed screens.

Use Tailwind color names (e.g., slate, zinc, stone, blue, emerald) and Google Fonts.

## Output Requirements
Create two JSON files:

### `product/design-system/colors.json`
A JSON object like:
```json
{
  "neutral": "slate",
  "primary": "blue",
  "secondary": "emerald",
  "danger": "red",
  "warning": "amber"
}
```

### `product/design-system/typography.json`
A JSON object like:
```json
{
  "heading": "Inter",
  "body": "Inter",
  "mono": "JetBrains Mono"
}
```

## Rules
- Keep the palette small (neutral + 2 accents + status colors).
- Pick fonts that fit the product context.
- Do not invent custom hex values; just choose token "families".


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
