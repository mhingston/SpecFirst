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
Design the API structure, endpoints, data models, and system architecture.

This design will be reviewed by architect and product teams before implementation.

## Output Requirements

Create `design.md` with these sections:

### API Specification

#### Endpoints
For each endpoint, define:
- Method: GET/POST/PUT/DELETE
- Path: `/api/v1/resource`
- Request parameters (query, body, headers)
- Response format (200, 400, 500, etc.)
- Example requests/responses

#### Data Models
Define request and response schemas (use JSON schema or examples):
```json
{
  "field": "type and description"
}
```

### Architecture

#### Components
- Which services/modules are involved?
- How do they interact?
- What databases or external systems?

#### Sequence Diagram
Describe the request flow from client → API → backend → response.

### Error Handling
- What error cases exist?
- How are they communicated to clients?
- Retry/fallback strategies?

### Security Considerations
- Authentication/authorization approach
- Rate limiting
- Input validation
- Sensitive data handling

### Performance Considerations
- Expected load/throughput
- Caching strategy
- Database query optimization

### Trade-offs
Document any design decisions and alternatives considered.

## Guidelines
- Be specific enough for implementation
- Leave room for engineer judgment on details
- Highlight risk areas

## Assumptions
- Requirements have been finalized
- (List other assumptions)


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
