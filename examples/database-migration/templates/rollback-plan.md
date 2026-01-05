# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Task
Create a comprehensive rollback strategy in case the migration fails or causes issues.

This will be reviewed by Ops before migration execution.

## Output Requirements

Create `rollback-plan.md` with:

### Rollback Triggers
When should we rollback?
- Migration script fails with error
- Performance degrades beyond X threshold
- Data corruption detected
- Application errors spike above Y%
- Manual decision after monitoring

### Rollback Scripts

**Level 1: Quick Rollback (before Step 4)**
```sql
-- Remove new column
ALTER TABLE users DROP COLUMN email;

-- Drop index
DROP INDEX idx_users_email;
```
**Time**: 30 seconds  
**Data loss**: Email data added during migration

**Level 2: Full Rollback (after Step 4)**
```sql
-- Remove NOT NULL constraint
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- Then follow Level 1 steps
```
**Time**: 1-2 minutes  
**Data loss**: Email data

**Level 3: Restore from Backup**
If data corruption occurs:
1. Stop application writes
2. Restore from backup taken at [timestamp]
3. Replay transaction logs to current time
4. Restart application

**Time**: 30-60 minutes  
**Data loss**: Potential data after backup

### Rollback Decision Matrix

| Situation | Action | Time | Data Loss |
|-----------|--------|------|-----------|
| Script error during Step 1-2 | Level 1 | 30s | None |
| Performance issue during Step 3 | Pause, optimize, continue | Varies | None |
| Constraint failure in Step 4 | Level 2 | 2m | Email data |
| Data corruption detected | Level 3 | 60m | Minimal |

### Health Checks
Before declaring rollback complete:
- [ ] Application stable
- [ ] Database responsive
- [ ] No error spikes
- [ ] Queries performing normally
- [ ] Data integrity verified

### Communication Plan
Who to notify and when:
- **Before rollback**: Ops lead, team lead
- **During rollback**: Engineering team, product
- **After rollback**: Post-mortem scheduled

### Post-Rollback Actions
1. Document what went wrong
2. Fix the root cause
3. Update migration plan
4. Schedule retry (if appropriate)

### Prevention Measures
How to avoid needing rollback:
- Test on staging first with production-like data volume
- Run migration during low-traffic period
- Monitor actively throughout
- Use batched operations
- Have DBA available

## Guidelines
- Make rollback fast and safe
- Prefer partial rollback to full restore
- Test rollback scripts on staging
- Keep ops team in the loop

## Assumptions
- Migration plan approved
- Backup is fresh and verified
- (List other assumptions)
