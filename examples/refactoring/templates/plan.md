# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Task
Create a detailed, step-by-step refactoring plan.

This plan should minimize risk and allow incremental progress.

## Output Requirements

Create `plan.md` with:

### Refactoring Steps
Break down the work into small, safe increments:

**Step 1: Example**
- What: Extract method `foo()` from `bar()`
- Why: Reduce function length, improve testability
- Risk: Low - pure function
- Verification: Unit test for `foo()`
- Rollback: Inline the function

List all steps in dependency order.

### Testing Strategy
- What tests to write first
- What tests to update
- How to verify behavior is preserved

### Migration Path
If this affects APIs or interfaces:
- How to maintain compatibility during transition
- Deprecation strategy
- Feature flags (if needed)

### Risk Mitigation
For each high-risk step:
- What could go wrong
- How to detect issues
- Rollback plan

### Incremental Checkpoints
Define merge points where code is stable:
- After step 3: Mergeable, all tests pass
- After step 7: Mergeable, metrics improved
- After step 10: Complete

### Timeline Estimate
- Step 1-3: 2 hours
- Step 4-7: 4 hours
- Step 8-10: 2 hours
- Total: ~8 hours

### Dependencies
- Need code freeze? When?
- Need deployment window?
- Coordination with other teams?

## Guidelines
- Keep steps small (1-2 hours each)
- Each step should be independently testable
- Prefer additive changes over deletions initially
- Plan for pause/resume points

## Assumptions
- Goals have been approved
- Have time allocated for refactoring
- (List other assumptions)
