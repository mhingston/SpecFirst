# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Task
Finalize the specification after review and iteration.

Incorporate feedback from:
- Assumption surfacing
- Role-based reviews
- Failure mode analysis
- Confidence calibration

## Output Requirements

Create `spec-final.md` with refined:

### Executive Summary
Clear 2-3 sentence summary of what this spec defines.

### Problem Statement
What problem are we solving and why?

### Solution Overview
High-level approach.

### Detailed Requirements
Organized by category:
- Functional requirements
- Non-functional requirements (performance, security, etc.)
- Constraints

### Architecture/Design
Technical approach with diagrams if helpful.

### Risk Mitigation
How we address identified risks and failure modes.

### Success Criteria
How we'll measure success.

### Assumptions (Explicit)
List all assumptions clearly marked.

### Out of Scope
What we're explicitly NOT doing.

## Guidelines
- Address concerns raised in reviews
- Strengthen low-confidence areas
- Be explicit about assumptions
- Make it reviewable by others

## Assumptions
- Reviews have been completed
- Stakeholders have provided input


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
