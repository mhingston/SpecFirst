# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Task
Define clear, measurable goals for this refactoring.

Avoid vague goals like "make it better" - be specific about what success looks like.

## Output Requirements

Create `goals.md` with:

### Primary Goals
List 3-5 specific, measurable objectives:
- ✅ "Reduce cyclomatic complexity from 25 to 10"
- ✅ "Eliminate 4 duplicated code blocks"
- ✅ "Increase test coverage from 40% to 80%"
- ✅ "Reduce function length from 200 to 50 lines"
- ❌ "Make code cleaner" (too vague)

### Non-Goals
What are we explicitly NOT doing?
- Changing functionality
- Adding new features
- Performance optimization (if not a goal)

### Success Criteria
How will we know the refactoring succeeded?
- All existing tests still pass
- New tests added for X
- Metrics improved by Y%
- Code review approval

### Constraints
- Must maintain backward compatibility
- Must complete within X timeframe
- Cannot change public APIs
- Must not affect performance

### Benefits
Why is this refactoring worth doing?
- Easier to maintain
- Easier to test
- Easier to extend with new features
- Reduced bug surface area

## Guidelines
- Prioritize goals (most important first)
- Make goals testable/verifiable
- Consider effort vs benefit

## Assumptions
- Current behavior analysis is accurate
- (List other assumptions)
