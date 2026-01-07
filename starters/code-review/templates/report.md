# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>
{{- end }}

## Task
Synthesize the verified findings into a final, actionable review report.

## Output Requirements

Create `{{ index .Outputs 0 }}` with:

### 1. Executive Summary
- **Overall Quality Grade**: (A-F)
- **Major Risks**: Top 2-3 critical items.
- **Pass/Fail**: Based on the Definition of Done in `scope.md`.

### 2. Actionable Findings (Verified)
**CRITICAL**: Only include items from the "Filtered Action Items" section of the skeptic pass. 
- Look for the YAML-ish block or structured list in that section.
- **Ignore** any findings marked REMOVED or DOWNGRADED in the verification matrix.
Grouped by severity (Sev0-Sev3). 
- Include file, identifier, description, and proposed fix for each.

### 3. Review Coverage
- Which modules were thoroughly reviewed?
- Which areas were excluded or had "No Findings"?

### 4. Next Steps
- What should the developer do first?
- Are there any unanswered questions that need manual verification?

---

## Guidelines
- **Convergence**: This report marks the end of the review cycle. 
- **Clarity**: Ensure fixes are copy-pasteable where appropriate.

## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Start directly with the markdown content.
