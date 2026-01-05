# Example: Bug Fix Workflow

A minimal 2-stage SpecFirst workflow for analyzing and fixing bugs systematically.

## What This Demonstrates

- Simple linear workflow (analysis → fix)
- Structured bug investigation before coding
- Using cognitive scaffold commands for risk assessment
- Minimal overhead for quick fixes

## Quick Start (Run in this repo)
 
You can run this example immediately using the `--protocol` override:
 
1. **Analysis**:
   ```bash
   specfirst --protocol examples/bug-fix/protocol.yaml analysis
   ```
 
2. **Fix**:
   ```bash
   specfirst --protocol examples/bug-fix/protocol.yaml fix
   ```
 
## Setup (For a new project)
 
To use this protocol in your own project:
 
1. Create a new directory and initialize:
   ```bash
   mkdir my-bugfix && cd my-bugfix
   specfirst init
   ```
 
2. Copy the protocol and templates:
   ```bash
   cp /path/to/specfirst/examples/bug-fix/protocol.yaml .specfirst/protocols/
   cp -r /path/to/specfirst/examples/bug-fix/templates/* .specfirst/templates/
   ```
 
3. Update config (optional) or use the flag:
   ```bash
   # Option A: Edit .specfirst/config.yaml to set protocol: bug-fix
   # Option B: Use flag
   specfirst --protocol bug-fix analysis
   ```

## Workflow

### 1. Analyze the Bug

Generate the analysis prompt:
```bash
specfirst analysis | claude -p > analysis.md
```

This will prompt you to document:
- Root cause
- Impact assessment
- Reproduction steps
- Proposed fix approach
- Risks

Complete the stage:
```bash
specfirst complete analysis ./analysis.md
```

### 2. (Optional) Check for Failure Modes

Before implementing, surface potential issues:
```bash
specfirst failure-modes ./analysis.md | claude -p
```

Review the output and update your analysis if needed.

### 3. Implement the Fix

Generate the implementation prompt (includes analysis automatically):
```bash
specfirst fix | claude -p
```

Save the generated code changes, then complete:
```bash
specfirst complete fix ./src/bugfix.go ./tests/bugfix_test.go
# Or whatever files you modified
```

### 4. Validate

Check that everything is complete:
```bash
specfirst check
```

## Timeline

**Quick fix**: 10-15 minutes  
**Complex bug**: 30-60 minutes

## When to Use This

- ✅ Bug reports from users or QA
- ✅ Production incidents requiring root cause analysis
- ✅ When you want structured documentation of the fix
- ❌ Trivial typos or one-line fixes (use git commit message)
