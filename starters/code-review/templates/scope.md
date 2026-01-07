# {{ .StageName }} - {{ .ProjectName }}

## Task
Define the scope, priorities, and constraints for this code review to ensure it converges and remains actionable.

## Output Requirements

Create `{{ index .Outputs 0 }}` with the following sections:

### 1. Target Code
- **Paths**: (e.g., `./src/auth`, `cmd/cli/`)
- **Exclusions**: (e.g., `*_test.go`, `vendor/`, `generated/`)
- **CRITICAL**: Never include secrets/credentials in included files; add them to exclusions (.env, .pem, id_rsa, credentials.json, etc.).
- **Git Context**: (e.g., "Review changes in branch `feature/new-auth` compared to `main`")

### 2. Review Priorities
Rank the following in order of importance (1-5):
- [ ] Correctness & Logic
- [ ] Security & Privacy
- [ ] Performance & Resource Usage
- [ ] Maintainability & Readability
- [ ] Test Coverage & Quality

### 3. Review Budget
- **Max Findings per Category**: (e.g., 2 items per category max)
- **Top Goal**: (e.g., "Find the top 3 items that would block a production release")

### 4. Definition of Done
- When are we finished? (e.g., "When all Sev0 and Sev1 items from the File Review are addressed")

---

## Guidelines
- **Be Selective**: Don't review everything. Focus on core modules and high-risk areas.
- **Set Limits**: Explicitly state that the model should stop when the budget is reached.
- **Freeze Architecture**: Remind the reviewer that after Stage 2, the architecture is frozen.

## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text.
- Do NOT include markdown code block fences.
- Start directly with the markdown content.
