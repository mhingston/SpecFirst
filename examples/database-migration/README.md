# Example: Database Schema Migration

A 4-stage workflow for planning and executing database migrations safely with proper approvals, rollback strategy, and risk mitigation.

## What This Demonstrates

- High-risk workflow with approval gates (DBA and Ops)
- Explicit rollback planning before execution
- Using `failure-modes` to identify migration risks
- Incremental execution with monitoring
- Balance between safety and minimal downtime

## Setup

1. Create a new directory and initialize:
   ```bash
   mkdir my-migration && cd my-migration
   specfirst init
   ```

2. Copy the database-migration protocol and templates:
   ```bash
   cp /path/to/specfirst/examples/database-migration/protocol.yaml .specfirst/protocols/
   cp -r /path/to/specfirst/examples/database-migration/templates/* .specfirst/templates/
   ```

3. Set the protocol in `.specfirst/config.yaml`:
   ```yaml
   protocol: database-migration
   project_name: my-migration
   ```

## Workflow

### 1. Document Requirements

Define what needs to change:
```bash
specfirst requirements | claude -p > requirements.md
specfirst complete requirements ./requirements.md
```

### 2. Create Migration Plan

Generate detailed migration SQL and steps:
```bash
specfirst migration-plan | claude -p > migration-plan.md
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
specfirst failure-modes ./migration-plan.md | claude -p
```

Review output and incorporate into rollback strategy.

### 5. Create Rollback Plan

Plan for what to do if migration fails:
```bash
specfirst rollback-plan | claude -p > rollback-plan.md
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
specfirst execute | claude -p
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
