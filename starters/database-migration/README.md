# Example: Database Schema Migration

A 4-stage workflow for planning and executing database migrations safely with proper approvals, rollback strategy, and risk mitigation.

## What This Demonstrates

- High-risk workflow with approval gates (DBA and Ops)
- Explicit rollback planning before execution
- Using `failure-modes` to identify migration risks
- Incremental execution with monitoring
- Balance between safety and minimal downtime

## Quick Start (Run in this repo)
 
You can run this example immediately using the `--protocol` override:
 
1. **Requirements**:
   ```bash
   opencode run "$(specfirst --protocol starters/database-migration/protocol.yaml requirements)"
   ```
 
2. **Migration Plan**:
   ```bash
   opencode run "$(specfirst --protocol starters/database-migration/protocol.yaml migration-plan)"
   ```
 
## Setup (For a new project)
 
To use this protocol in your own project:
 
1. Create a new directory and initialize it with Git:
   ```bash
   mkdir my-migration && cd my-migration
   git init
   ```

2. Initialize SpecFirst with the `database-migration` starter:
   ```bash
   specfirst init --starter database-migration
   ```

## Workflow

### 1. Document Requirements

Define what needs to change:
```bash
opencode run "$(specfirst requirements)" > requirements.md
specfirst complete requirements ./requirements.md
```

### 2. Create Migration Plan

Generate detailed migration SQL and steps:
```bash
opencode run "$(specfirst migration-plan)" > migration-plan.md
specfirst complete migration-plan ./migration-plan.md
```

### 3. Get DBA Approval

Protocol requires DBA review before proceeding:
```bash
specfirst approve migration-plan --role dba --by "Jane DBA"
```

### 4. (Optional) Identify Failure Modes

Surface risks before creating rollback plan:
```bash
opencode run "$(specfirst failure-modes ./migration-plan.md)"
```

Review output and incorporate into rollback strategy.

### 5. Create Rollback Plan

Plan for what to do if migration fails:
```bash
opencode run "$(specfirst rollback-plan)" > rollback-plan.md
specfirst complete rollback-plan ./rollback-plan.md
```

### 6. Get Ops Approval

Protocol requires Ops review of rollback strategy:
```bash
specfirst approve rollback-plan --role ops --by "Bob Ops"
```

### 7. Execute Migration

Run the migration with monitoring:
```bash
opencode --prompt "$(specfirst execute)"
# Follow the plan step by step, documenting as you go
specfirst complete execute ./migration-log.md ./scripts/*.sql
```

### 8. Verify Completion

Validate all approvals and outputs:
```bash
specfirst check
```

## Timeline

**Small migration** (add column): 1-2 hours  
**Medium migration** (schema change): 3-5 hours  
**Large migration** (data transformation): 1-2 days

Add buffer time for approvals and coordination.

## When to Use This

- ✅ Production database schema changes
- ✅ High-volume table modifications  
- ✅ Migrations requiring coordination
- ✅ When downtime must be minimized
- ❌ Dev/staging schema changes (just do them)
- ❌ Trivial index additions (use simpler workflow)

## Key Safety Features

- **Approval gates**: DBA and Ops must sign off
- **Explicit rollback**: Plan the "undo" before executing
- **Incremental steps**: Break down into small, safe operations
- **Monitoring**: Track performance throughout
- **Risk assessment**: Use failure-modes to identify issues

## Best Practices

1. Always test on staging with production-like data volume first
2. Run during off-peak hours
3. Have DBA available during execution
4. Monitor replication lag, query performance, and lock wait time
5. Use batched operations for large data changes
6. Verify backups before starting
