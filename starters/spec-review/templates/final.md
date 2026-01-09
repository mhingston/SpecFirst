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
- **Production hardening report**

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

### Technical Reliability
(Incorporate findings from the hardening report)

**Invariants (Crash Conditions)**:
- List the system invariants identified in the hardening report
- These are "impossible states" - if violated, the system is corrupt
- Specify that violation results in immediate termination (panic/crash)

**Observability Contracts**:
- Define the standard log context fields for this feature
- Every log line must include enough data to reproduce issues locally

**Error Handling Strategy**:
- Recoverable errors (retry/backoff)
- Corrupt state errors (crash/alert)

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
- **Ensure all hardening requirements are incorporated**

## Assumptions
- Reviews have been completed
- Stakeholders have provided input
- Production hardening has been performed


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).

