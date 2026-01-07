# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
Perform a detailed review of individual files within the scoped area. **Architecture is frozen; do not suggest structural rewrites.**

## Output Requirements

Create `{{ index .Outputs 0 }}` with:

### 1. Findings List
- **Constraint**: Max 12 issues total across all files.
- **Format for each issue**:
  - **Severity**: Sev0 (Bug), Sev1 (Correctness/Security), Sev2 (Maintainability), Sev3 (Polish)
  - **Location**: `path/to/file.ext:Identifier` (Function/Class)
  - **Issue**: Precise description of the problem.
  - **Evidence**: Quote the code or explain the reasoning chain.
  - **Fix**: Suggested patch or concrete change.

### 2. Questions / Suspicions
If you suspect an issue but can't point to definitive code evidence, list it here as a **QUESTION**, not an issue.

---

## Guidelines
- **Cite Concrete Code**: Only raise an issue if you can cite a file and identifier.
- **Adhere to Priorities**: Look for issues matching the priorities defined in `scope.md`.
- **Budget Enforcement**: Stop exactly when the budget is spent. Do not invent extra issues.

## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Start directly with the markdown content.
