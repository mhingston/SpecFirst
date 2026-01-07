# {{ .StageName }} - Final Review Consolidation

## Context
**RFC**:
{{ readFile "inputs/sample_rfc.md" }}

**Principal Review**:
{{ readFile "review.md" }}

**Skeptic Findings**:
{{ readFile "skeptic_findings.md" }}

## Instructions
Synthesize a single, actionable Final Review.
- Combine duplicate points.
- Elevate "Skeptic" findings if they are valid.
- Categorize feedback into: **Must Fix**, **Should Fix**, **Consider**, **Questions**.

## Output Format
Markdown file (`final_review.md`).
Start with a "Verdict": **Approve**, **Approve with Comments**, or **Request Changes**.
