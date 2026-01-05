# Example: API Feature with Approvals

A complete SpecFirst workflow demonstrating requirements, design with approvals, task decomposition, and implementation for a new API endpoint or feature.

## What This Demonstrates

- Full 4-stage workflow (requirements → design → decompose → implementation)
- Approval gates (architect and product must approve design)
- Task decomposition for parallel development
- Task-scoped implementation prompts

## Quick Start (Run in this repo)
 
You can run this example immediately using the `--protocol` override:
 
1. **Requirements**:
   ```bash
   specfirst --protocol examples/api-feature/protocol.yaml requirements
   ```
 
2. **Design**:
   ```bash
   specfirst --protocol examples/api-feature/protocol.yaml design
   ```
 
## Setup (For a new project)
 
To use this protocol in your own project:
 
1. Create a new directory and initialize:
   ```bash
   mkdir my-api-feature && cd my-api-feature
   specfirst init
   ```
 
2. Copy the protocol and templates:
   ```bash
   cp /path/to/specfirst/examples/api-feature/protocol.yaml .specfirst/protocols/
   cp -r /path/to/specfirst/examples/api-feature/templates/* .specfirst/templates/
   ```
 
3. Update config (optional) or use the flag:
   ```bash
   # Option A: Edit .specfirst/config.yaml to set protocol: api-feature
   # Option B: Use flag
   specfirst --protocol api-feature requirements
   ```
   ```yaml
   protocol: api-feature
   project_name: my-api-feature
   language: Go  # or your language
   framework: gin  # or your framework
   ```

## Workflow

### 1. Gather Requirements

Generate the requirements prompt:
```bash
specfirst requirements | claude -p > requirements.md
```

Complete the stage:
```bash
specfirst complete requirements ./requirements.md
```

### 2. Create Design

Generate the design prompt (automatically includes requirements):
```bash
specfirst design | claude -p > design.md
```

Complete the stage:
```bash
specfirst complete design ./design.md
```

### 3. Get Approvals

The protocol requires architect and product approval before proceeding:

```bash
# Architect approval
specfirst approve design --role architect --by "Jane Smith"

# Product approval  
specfirst approve design --role product --by "Bob Johnson"
```

Check status:
```bash
specfirst status
```

### 4. (Optional) Review Design Quality

Before decomposing, you can use cognitive commands:

```bash
# Security review
specfirst review ./design.md --persona security | claude -p

# Performance review
specfirst review ./design.md --persona performance | claude -p

# Surface assumptions
specfirst assumptions ./design.md | claude -p
```

### 5. Break Down into Tasks

Generate the decomposition prompt:
```bash
specfirst decompose | claude -p > tasks.yaml
```

Complete:
```bash
specfirst complete decompose ./tasks.yaml
```

### 6. List and Implement Tasks

See all tasks:
```bash
specfirst task
```

Generate prompt for a specific task:
```bash
specfirst task T1 | claude -p
```

After implementing, complete it:
```bash
specfirst complete implementation ./api/handler.go ./api/handler_test.go
```

Repeat for each task (can be done in parallel by team members).

### 7. Final Validation

Check everything is complete:
```bash
specfirst check
```

Archive the completed spec:
```bash
specfirst complete-spec --archive --version 1.0
```

## Timeline

**Solo developer**: 2-3 hours  
**Small team** (parallel tasks): 1-2 hours  
**With approvals/reviews**: Add 30-60 min

## When to Use This

- ✅ New API endpoints with cross-team coordination
- ✅ Features requiring architect/product sign-off
- ✅ Work that can be split among multiple developers
- ✅ When you need approval audit trail
- ❌ Tiny internal APIs (use simpler workflow)
