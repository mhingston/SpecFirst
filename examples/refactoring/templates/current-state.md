# {{ .StageName }} - {{ .ProjectName }}

## Task
Document the current state of the code you want to refactor.

This creates a baseline understanding before making changes.

## Output Requirements

Create `current-state.md` with:

### Code Location
- File paths
- Module/package names
- Lines of code involved

### Purpose
What does this code currently do?

### Structure
- Key functions/classes/components
- Dependencies between them
- External dependencies (libraries, APIs, databases)

### Problems
Why does this code need refactoring?
- Performance issues
- Maintainability concerns
- Code duplication
- Hardcoded values
- Missing abstractions
- Tight coupling
- Testing difficulties
- Security concerns

### Current Behavior
Document the expected behavior that MUST be preserved after refactoring.

### Test Coverage
- What tests exist?
- What's not tested?
- Test execution time

### Metrics (if available)
- Cyclomatic complexity
- Lines of code
- Number of dependencies
- Test coverage percentage

## Guidelines
- Be objective, not judgmental
- Focus on facts, not opinions
- Document what works well too

## Assumptions
- You have access to the codebase
- (List other assumptions)
