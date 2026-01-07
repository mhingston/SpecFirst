# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
Conduct a high-level architectural and structural review. **Do not perform file-by-file nitpicks at this stage.**

## Output Requirements

Create `{{ index .Outputs 0 }}` with:

### 1. Structural Analysis
- **Core Pattern**: (e.g., Layered, Hexagonal, Spaghetti)
- **Dependency Flow**: Are there circular dependencies? High coupling?

### 2. High-Level Findings
- **Constraint**: Max 10 findings total.
- **Focus**: Cross-cutting concerns, module boundaries, data flow, entrypoints.
- **Rule**: Every finding must reference a concrete structural element.

### 3. Proposed Target Architecture
- 8-12 bullets describing the ideal structure for this component.
- Risks/Tradeoffs of moving to this architecture.

---

## Guidelines
- **No Style Comments**: Ignore naming, spacing, or local logic issues.
- **Evidence-Driven**: If you can't point to a relationship between modules/files, it's not a valid architectural finding.
- **Gatekeeper**: This stage "freezes" the architecture for subsequent review stages.

## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Start directly with the markdown content.
