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
   specfirst --protocol starters/api-feature/protocol.yaml requirements | gemini -i
   ```
 
2. **Design**:
   ```bash
   specfirst --protocol starters/api-feature/protocol.yaml design | gemini -i
   ```
 
## Setup (For a new project)
 
To use this protocol in your own project:
 
1. Create a new directory and initialize it with Git:
   ```bash
   mkdir my-api-feature && cd my-api-feature
   git init
   ```

2. Initialize SpecFirst with the `api-feature` starter:
   ```bash
   specfirst init --starter api-feature
   ```

3. Update project metadata in `.specfirst/config.yaml` (optional):
   ```yaml
   project_name: my-api-feature
   language: Go  # or your language
   framework: gin  # or your framework
   ```

## Workflow

### 1. Gather Requirements

Generate the requirements prompt:
```bash
specfirst requirements | gemini -i > requirements.md
```

Complete the stage:
```bash
specfirst complete requirements ./requirements.md
```

### 2. Create Design

Generate the design prompt (automatically includes requirements):
```bash
specfirst design | gemini -i > design.md
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
specfirst review ./design.md --persona security | gemini -i

# Performance review
specfirst review ./design.md --persona performance | gemini -i

# Surface assumptions
specfirst assumptions ./design.md | gemini -i
```

### 5. Break Down into Tasks

Generate the decomposition prompt:
```bash
specfirst decompose | gemini -i > tasks.yaml
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
specfirst task T1 | gemini -i
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
