# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Task
Verify that the refactoring achieved its goals without breaking behavior.

## Output Requirements

Create `verification-report.md` with:

### Goals Achievement
For each goal from the plan:
- ✅ Goal: Reduce complexity from 25 to 10
  - Result: Reduced to 8
  - Evidence: Cyclomatic complexity metrics
  
- ✅ Goal: Increase test coverage to 80%
  - Result: 85% coverage
  - Evidence: Coverage report

### Behavior Preservation
- [ ] All existing tests pass
- [ ] Regression testing complete
- [ ] No new bugs reported
- [ ] Performance unchanged (or improved)

### Code Quality Metrics

**Before Refactoring:**
- Cyclomatic complexity: 25
- Lines of code: 500
- Test coverage: 40%
- Duplicated code blocks: 4

**After Refactoring:**
- Cyclomatic complexity: 8
- Lines of code: 350
- Test coverage: 85%
- Duplicated code blocks: 0

### Test Results
- Unit tests: X/X passed
- Integration tests: Y/Y passed
- Regression suite: Z/Z passed

### Code Review Feedback
Summary of peer review comments and resolutions.

### Issues Found
Any problems discovered during verification:
- Issue: Description
- Severity: Critical/High/Medium/Low
- Status: Fixed/Deferred/Won't Fix
- Resolution: What was done

### Lessons Learned
- What went well
- What was harder than expected
- What would you do differently next time

## Guidelines
- Be honest about results
- Document any compromises made
- Note any new technical debt created

## Assumptions
- Refactoring has been completed
- Tests have been run
- (List other assumptions)
