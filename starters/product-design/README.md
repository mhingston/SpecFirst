# Product Design Protocol

A workflow for product design, from vision through to implementation handoff.

## Stages

| Stage | Output Artifacts |
|-------|-----------------|
| `product-vision` | `product/product-overview.md` |
| `product-roadmap` | `product/product-roadmap.md` |
| `data-model` | `product/data-model/data-model.md` |
| `design-tokens` | `product/design-system/colors.json`, `typography.json` |
| `design-shell` | `product/shell/spec.md`, `src/shell/components/*.tsx` |
| `shape-section` | `product/sections/<id>/spec.md` (repeatable) |
| `sample-data` | `product/sections/<id>/data.json`, `types.ts` (repeatable) |
| `design-screen` | `src/sections/<id>/*.tsx` (repeatable) |
| `screenshot-design` | `product/sections/<id>/*.png` (repeatable) |
| `export-product` | `product-plan/*` |

## Setup

To use this protocol in a new project:

1. Create a new directory and initialize it with Git:
   ```bash
   mkdir my-product && cd my-product
   git init
   ```

2. Initialize SpecFirst with the `product-design` starter:
   ```bash
   specfirst init --starter product-design
   ```

3. Run the workflow stages:
   ```bash
   opencode run "$(specfirst product-vision)" > product/product-overview.md
   ```

## Section Looping

This protocol supports iteration via `repeatable: true` stages. To design multiple sections:

1. Set `section_id` in your `.specfirst/config.yaml`:
   ```yaml
   custom_vars:
     section_id: invoices
     section_title: Invoices
   ```

2. Run the section stages:
   ```bash
   opencode run "$(specfirst shape-section)"
   # Complete with: specfirst complete shape-section product/sections/invoices/spec.md
   
   opencode run "$(specfirst sample-data)"
   # Complete with: specfirst complete sample-data product/sections/invoices/data.json product/sections/invoices/types.ts
   
   opencode run "$(specfirst design-screen)"
   # Complete with: specfirst complete design-screen src/sections/invoices/InvoiceList.tsx ...
   ```

3. Repeat for each section by changing `section_id`.

## Artifact Structure

```
product/
├── product-overview.md
├── product-roadmap.md
├── data-model/
│   └── data-model.md
├── design-system/
│   ├── colors.json
│   └── typography.json
├── shell/
│   └── spec.md
└── sections/
    ├── invoices/
    │   ├── spec.md
    │   ├── data.json
    │   └── types.ts
    └── dashboard/
        └── ...

src/
├── shell/
│   ├── components/*.tsx
│   └── ShellPreview.tsx
└── sections/
    └── invoices/
        └── *.tsx

product-plan/
└── ... (export output)
```

## Skills

This example uses the `readFile` template helper to include reusable skills:
- `skills/design-principles.md` - Design quality guidelines

Templates reference skills via `{{ readFile "design-principles.md" }}`.

