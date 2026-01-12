# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Requirements Contract
You must follow the clarified requirements exactly.

- Treat In Scope and Out of Scope / Non-Goals as hard boundaries.
- Use Acceptance Criteria as the definition of done.
- Respect all stated Constraints.
- Treat listed Assumptions as true unless explicitly challenged.

If you believe the requirements are incorrect, incomplete, or unsafe:
- Stop.
- Explain the issue.
- Propose a revision to requirements.md.
- Do NOT silently override requirements.

## Task
Break down the API design into implementable tasks.

Each task should be:
- Independently completable
- Testable
- Reviewable
- Small enough to finish in 1-4 hours

## Output Format
YAML list of tasks with this structure:

```yaml
tasks:
  - id: T1
    title: "Short task name"
    goal: "What this task accomplishes"
    acceptance_criteria:
      - "Criterion 1"
      - "Criterion 2"
    dependencies: []  # IDs of tasks that must complete first
    files_touched:
      - "path/to/file.go"
    risk_level: "low|medium|high"
    estimated_scope: "S|M|L"
    test_plan:
      - "How to verify this task"
```

## Decomposition Guidelines

### Good Task Boundaries
- Implement single endpoint
- Add database migration
- Write integration test suite
- Add monitoring/alerts

### Avoid
- Tasks mixing infrastructure + business logic
- Tasks spanning multiple endpoints
- "Finish everything" tasks

### Ordering
Order tasks by:
1. Dependencies first (data models, migrations)
2. Core logic
3. Tests and monitoring
4. Documentation

## Parameters
- Granularity: {{ if .Prompt }}{{ .Prompt.Granularity }}{{ else }}ticket{{ end }}
- Max tasks: {{ if .Prompt }}{{ .Prompt.MaxTasks }}{{ else }}10{{ end }}

## Assumptions
- Design has been approved by architect and product
- (List other assumptions)


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
