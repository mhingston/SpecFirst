# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Task
Execute the database migration according to the approved plan.

Document each step as you go, noting any deviations or issues.

## Output Format

Provide migration execution log with:

### Pre-Migration Checklist
- [ ] Backup verified (timestamp, size)
- [ ] Application code deployed
- [ ] Maintenance window scheduled (if needed)
- [ ] Ops team notified
- [ ] Monitoring dashboards open
- [ ] Rollback scripts ready

### Execution Log

**Step 1: Pre-validation**
- Started: [timestamp]
- Actions taken: [what you did]
- Results: [output/metrics]
- Status: ✅ Success / ❌ Failed / ⚠️ Issues
- Duration: X minutes
- Notes: [any observations]

**Step 2: Schema changes**
- Started: [timestamp]
- SQL executed:
  ```sql
  ALTER TABLE users ADD COLUMN email VARCHAR(255);
  ```
- Results: [rows affected, execution time]
- Status: ✅ Success
- Duration: 0.5 minutes
- Notes: No locks detected

(Continue for each step...)

### Monitoring Observations
- Peak replication lag: X seconds
- Query performance: Normal/Degraded
- Error rate: Y%
- Lock wait time: Z ms

### Issues Encountered
For each issue:
- **Issue**: Description
- **When**: Step and timestamp
- **Severity**: Critical/High/Medium/Low
- **Resolution**: What was done
- **Impact**: User/system impact

### Verification
- [ ] All migration steps completed
- [ ] Row counts match expected
- [ ] Indexes created successfully
- [ ] Constraints applied
- [ ] Application tests passing
- [ ] No unexpected errors in logs
- [ ] Performance within acceptable range

### Final State
- Rows migrated: X
- Total time: Y hours
- Rollback needed: Yes/No
- Production impact: None/Minimal/Significant

### Post-Migration Actions
- [ ] Application monitoring - 24 hours
- [ ] Performance review
- [ ] Cleanup old data (scheduled)
- [ ] Documentation updated
- [ ] Post-mortem (if issues)

## Guidelines
- Log everything in real-time
- Don't skip verification steps
- If unsure, pause and consult DBA
- Use rollback plan if needed

## Assumptions
- Plans have DBA and Ops approval
- Have required access and permissions
- (List other assumptions)
