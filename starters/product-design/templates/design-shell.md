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
Design and implement an **application shell** that wraps the product sections.

The shell includes:
- Global layout (sidebar/topbar/etc.)
- Primary navigation for the roadmap sections
- User menu placement
- Responsive behavior notes (mobile)
- A minimal set of shell components

## Output Requirements

### Spec
Create `product/shell/spec.md` describing:
- Shell layout choice and rationale
- Navigation structure (section ids/titles)
- Default landing route (which section)
- User menu items
- Responsive behavior and constraints
- Accessibility notes (keyboard focus, etc.)

### Code
Create:
- `src/shell/components/AppShell.tsx`
- `src/shell/components/MainNav.tsx`
- `src/shell/components/UserMenu.tsx`
- `src/shell/components/index.ts`
- `src/shell/ShellPreview.tsx`

## Implementation rules (shell code)
- React + Tailwind.
- Use tokens from `colors.json` and `typography.json` as the styling source of truth.
- Shell components should accept props where reasonable (nav items, current section id, user).
- ShellPreview can be a simple wrapper for local previewing (not export-critical).


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
