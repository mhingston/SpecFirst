# {{ .StageName }} - {{ .ProjectName }}

## Task
Gather and document requirements for a new API endpoint or feature.

Define what the API should do from the user's perspective, without diving into implementation details.

## Output Requirements

Create `requirements.md` with these sections:

### Feature Overview
Brief description of what this API feature does and why it's needed.

### User Stories
List 3-7 user stories in the format:
- As a [user type], I want [action] so that [benefit]

### Functional Requirements
What must the API do? Be specific:
- Input parameters
- Output format
- Behavior/logic
- Edge cases

### Non-Functional Requirements
- Performance expectations (latency, throughput)
- Security requirements
- Scalability needs
- Monitoring/observability

### Constraints
- Technical limitations
- Business rules
- Compliance requirements (GDPR, etc.)

### Out of Scope
Explicitly list what this feature will NOT do.

### Success Criteria
How will we know this feature is successful?

## Guidelines
- Focus on WHAT, not HOW
- Be concrete with examples
- Call out any assumptions
- Identify unknowns or questions

## Assumptions
- (List assumptions here - e.g., "We have auth infrastructure in place")


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
