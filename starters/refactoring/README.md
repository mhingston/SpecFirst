# Example: Code Refactoring

A comprehensive 5-stage workflow for planning and executing code refactoring safely and systematically.

## What This Demonstrates

- Structured approach to improving existing code
- Using `trace` and `diff` commands to understand code

## Quick Start (Run in this repo)
 
You can run this example immediately using the `--protocol` override:
 
1. **Current State**:
   ```bash
   opencode run "$(specfirst --protocol starters/refactoring/protocol.yaml current-state)"
   ```
 
2. **Goals**:
   ```bash
   opencode run "$(specfirst --protocol starters/refactoring/protocol.yaml goals)"
   ```
 
## Setup (For a new project)
 
To use this protocol in your own project:
 
1. Create a new directory and initialize it with Git:
   ```bash
   mkdir my-refactoring && cd my-refactoring
   git init
   ```

2. Initialize SpecFirst with the `refactoring` starter:
   ```bash
   specfirst init --starter refactoring
   ```

## Workflow

### 1. Analyze Current State

Map existing code to understand what you're refactoring:
```bash
# Map code to specifications
opencode run "$(specfirst trace ./path/to/current-code.go)"

# Generate current state analysis
opencode run "$(specfirst current-state)" > current-state.md
specfirst complete current-state ./current-state.md
```

### 2. Define Goals

Set clear, measurable refactoring objectives:
```bash
opencode run "$(specfirst goals)" > goals.md
specfirst complete goals ./goals.md
```

### 3. (Optional) Identify Risks

Before planning, surface potential problems:
```bash
opencode run "$(specfirst failure-modes ./goals.md)"
```

### 4. Create Refactoring Plan

Generate detailed step-by-step plan:
```bash
opencode run "$(specfirst plan)" > plan.md
specfirst complete plan ./plan.md
```

### 5. Execute Refactoring

Follow the plan step by step:
```bash
opencode --prompt "$(specfirst execute)"
# Implement changes following the plan
specfirst complete execute ./path/to/refactored-code.go ./tests/
```

### 6. Verify Results

Confirm goals met and behavior preserved:
```bash
opencode run "$(specfirst verify)" > verification-report.md
specfirst complete verify ./verification-report.md
```

### 7. (Optional) Compare Before/After

Analyze the changes made:
```bash
opencode run "$(specfirst diff ./current-state.md ./verification-report.md)"
```

## Timeline

**Small refactoring** (single function): 1-2 hours  
**Medium refactoring** (module/class): 4-8 hours  
**Large refactoring** (subsystem): 1-3 days

## When to Use This

- ✅ Improving code quality without changing behavior
- ✅ Reducing technical debt
- ✅ Making code more maintainable/testable
- ✅ When you need to justify refactoring effort
- ❌ Quick, obvious improvements (just do them)
- ❌ Refactoring as part of new feature (use feature workflow)

## Key Benefits

- **Risk reduction**: Incremental steps with rollback points
- **Measurable progress**: Clear goals and metrics
- **Team alignment**: Documented rationale and plan
- **Audit trail**: Complete record of what changed and why
