# {{ .StageName }} - {{ .ProjectName }}

## Task
Document the requirements for this database migration.

Be specific about what needs to change and why.

## Output Requirements

Create `requirements.md` with:

### Migration Goal
What database change are we making?
- Adding columns/tables
- Modifying schema
- Data transformation
- Index changes
- Constraint changes

### Business Context
Why is this migration needed?
- New feature requirement
- Performance improvement
- Compliance requirement
- Bug fix

### Scope
**In Scope:**
- Tables affected
- Estimated rows impacted
- Expected data volume

**Out of Scope:**
- What we're NOT changing

### Current Schema
Document the current state:
```sql
CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100)
);
```

### Desired Schema
Document the target state:
```sql
CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100),
  email VARCHAR(255) NOT NULL  -- NEW
);
```

### Data Considerations
- Existing data: X million rows
- Growth rate: Y rows/day
- Read/write load: Z QPS
- Nullable vs NOT NULL requirements
- Default values needed

### Constraints
- Downtime window available? How long?
- Can we do this online?
- Performance requirements during migration
- Compliance requirements (GDPR, data retention, etc.)

### Dependencies
- Application code changes needed?
- API version changes?
- Coordination with other teams?

## Guidelines
- Be explicit about data volume
- Identify breaking changes
- Note performance concerns

## Assumptions
- Have database access
- (List other assumptions)
