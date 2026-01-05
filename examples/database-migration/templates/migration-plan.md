# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}

## Task
Create a detailed migration plan with SQL scripts and execution steps.

This will be reviewed by DBA before execution.

## Output Requirements

Create `migration-plan.md` with:

### Migration Scripts

**1. Schema Changes:**
```sql
-- Add email column
ALTER TABLE users ADD COLUMN email VARCHAR(255);

-- Add index
CREATE INDEX idx_users_email ON users(email);

-- Add constraint (done separately to avoid locking)
-- ALTER TABLE users ALTER COLUMN email SET NOT NULL;
```

**2. Data Migration:**
```sql
-- Backfill email column from separate table
UPDATE users u
SET email = ue.email
FROM user_emails ue
WHERE u.id = ue.user_id;
```

**3. Cleanup:**
```sql
-- Add NOT NULL constraint after backfill
ALTER TABLE users ALTER COLUMN email SET NOT NULL;

-- Drop old table if needed
-- DROP TABLE user_emails;
```

### Execution Steps
Step-by-step execution plan:

**Step 1: Pre-migration validation**
- Time: 5 min
- Actions:
  - Verify backup exists
  - Check current row count
  - Verify disk space
- Verification: `SELECT COUNT(*) FROM users;`

**Step 2: Add column (non-blocking)**
- Time: 30 sec
- Actions: Run schema change script
- Verification: `DESCRIBE users;`
- Impact: None (nullable column)

**Step 3: Backfill data**
- Time: 2-4 hours (batched)
- Actions: Run data migration in 10k row batches
- Verification: Check progress every 15 min
- Impact: Elevated DB load

**Step 4: Add constraint**
- Time: 1-2 min
- Actions: Add NOT NULL
- Verification: Try INSERT without email (should fail)
- Impact: Brief table lock

### Rollback Points
- After Step 2: Can rollback safely
- After Step 3: Can rollback with data loss
- After Step 4: Should use rollback plan

### Monitoring
- Query performance during migration
- Replication lag
- Disk I/O
- Lock wait time
- Error logs

### Risk Assessment
- **High risk**: Adding NOT NULL without backfill
- **Medium risk**: Index creation load
- **Low risk**: Adding nullable column

### Timeline
- Total estimated time: 3-5 hours
- Maintenance window needed: No (online migration)
- Best execution time: Off-peak hours (2-6 AM)

### Pre-requisites
- [ ] Application code deployed with email handling
- [ ] Backup completed
- [ ] DBA approval received
- [ ] Ops team notified

## Guidelines
- Batch large data operations
- Add constraints last
- Monitor throughout
- Have DBA available during execution

## Assumptions
- Requirements approved
- Have production access
- (List other assumptions)
