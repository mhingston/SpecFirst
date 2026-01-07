# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
Act as a skeptic to verify findings and eliminate hallucinations or weak claims.

## Output Requirements

Create `{{ index .Outputs 0 }}` with:

### 1. Verification Matrix
For each finding in `file-findings.md`:
- **Finding**: (Summary)
- **Skeptic's Challenge**: Attempt to disprove this finding given the code (see `scope.md` or analyzed files). Is there a valid reason it was written this way?
- **Confidence Score**: 0-100%
- **Status**: VERIFIED / DOWNGRADED / REMOVED

### 2. Filtered Action Items
Only include items with a **Confidence Score ≥ 70%**.
**Format**:
```yaml
- severity: (Sev0-Sev3)
  file: (path)
  identifier: (function/class)
  description: (summary)
  fix: (proposed fix)
```
- List them by severity.
- Include the refined reasoning and fix.

---

## Guidelines
- **Be Ruthless**: It is better to miss a minor issue than to report a hallucinated or incorrect one.
- **Look for context**: Check if the "issue" is actually handled elsewhere in the code (e.g., error handling in a parent caller).

## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Start directly with the markdown content.
