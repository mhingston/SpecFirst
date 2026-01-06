# {{ .StageName }} — {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
Generate a complete **handoff package** under `product-plan/` so another coding agent/LLM can implement the real product in a separate codebase.

## Step 1: Check prerequisites
Minimum required:
- `product/product-overview.md`
- `product/product-roadmap.md`
- at least one section with screen designs in `src/sections/<section-id>/`

Recommended (warn if missing, but continue):
- `product/data-model/data-model.md`
- `product/design-system/colors.json`
- `product/design-system/typography.json`
- `product/shell/spec.md` and shell components

If minimum required are missing, output ONLY:
- a short message listing what's missing and which stage(s) to run.

## Step 2: Produce export directory structure
Create:

```
product-plan/
├── README.md
├── product-overview.md
├── prompts/
│   ├── one-shot-prompt.md
│   └── section-prompt.md
├── instructions/
│   ├── one-shot-instructions.md
│   └── incremental/
│       ├── 01-foundation.md
│       ├── 02-shell.md
│       └── NN-<section-id>.md
├── design-system/
│   ├── colors.json
│   └── typography.json
├── data-model/
│   ├── data-model.md
│   └── types.ts
├── shell/
│   ├── README.md
│   └── components/...
└── sections/
    └── <section-id>/
        ├── README.md
        ├── tests.md
        ├── types.ts
        ├── sample-data.json
        └── components/...
```

## Step 3: Fill in contents

### `product-plan/README.md`
- What this package is
- Prereqs (React 18, Tailwind)
- How to implement (one-shot vs incremental)
- Where tokens/types/components are

### `product-plan/product-overview.md`
Summarize:
- Product description
- Sections in order
- Data model entity list (or "not defined")
- Design tokens summary

### `product-plan/prompts/*.md`
Provide prompts that:
- Instruct an implementation LLM to build the product in a new codebase
- Require props-based components, tokens usage, and consistent types

### `product-plan/instructions/*`
Write step-by-step implementation instructions:
- Foundation: project setup, tailwind, tokens wiring, routing
- Shell: implement shell
- Per-section: implement from section spec + types + component designs

### Copy-through assets
Copy the latest versions of:
- tokens JSON
- data model md
- shell spec + components
- each section's spec, sample data, types, components

## Output Requirements
Write all files under `product-plan/` (at least one level deep so it matches `product-plan/*`).


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
