# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Philosophy
{{ readFile "production-mindset.md" }}

## Task
Review the Draft Specification specifically for **Assertiveness** and **Observability**.

## Output Requirements
Create `hardening-report.md` with the following sections:

### 1. System Invariants (The "Crash" List)
Identify states that should be theoretically impossible.
*Example: "A transaction cannot exist without a timestamp."*
- **Requirement**: Explicitly state that the code MUST panic/crash if these occur.

### 2. Observability Contracts
Define the **Context** that must be present in every log line for this feature.
*Example: "Every log must include `request_id`, `tenant_id`, and `step_name`."*

### 3. Error Handling Strategy
- Identify **Recoverable Errors** (e.g., API timeout) -> Retry/Backoff.
- Identify **Corrupt State Errors** (e.g., Data mismatch) -> Crash/Alert.

### 4. Complexity Red Flags
Identify any logic in the draft that looks "clever" or overly complex. Propose a "boring" alternative.

## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
