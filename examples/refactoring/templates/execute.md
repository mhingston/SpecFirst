# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Task
Execute the refactoring plan step by step.

Follow the plan, documenting any deviations.

## Output Format
Provide code changes as:
- Unified diffs for modifications
- Complete new files
- List of deleted files
- Updated tests

For each change, explain:
- Which plan step this addresses
- Why (if deviating from plan)
- What to verify

## Implementation Guidelines
- Follow the step order in the plan
- Commit after each checkpoint
- Run tests after each step
- Don't skip steps to "save time"
- If you discover issues, pause and update the plan

## Example Output Structure

### Step 1: Extract Method `foo()`
```diff
--- a/src/bar.go
+++ b/src/bar.go
@@ -10,15 +10,7 @@
 func bar() {
-    // inline code here
-    // ...
-    // ...
+    result := foo()
+    return result
 }
+
+func foo() int {
+    // extracted code
+    return 42
+}
```

**Tests Added:**
- `TestFoo()` - verifies extraction logic
- `TestBarCallsFoo()` - integration test

**Verification:** `make test` - all pass ✓

---

Repeat for each step...

## Verification Checklist
- [ ] All existing tests pass
- [ ] New tests added and passing
- [ ] No behavioral changes (regression testing)
- [ ] Code metrics improved as planned
- [ ] No new lint warnings
- [ ] Documentation updated

## Assumptions
- Plan has been reviewed
- Tests exist for behavior preservation
- (List other assumptions)
