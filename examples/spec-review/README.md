# Example: Specification Review Cycle

A workflow focused on using SpecFirst's cognitive scaffold commands to improve specification quality through structured review and iteration.

## What This Demonstrates

- Using cognitive commands to refine specs
- Multiple review perspectives (security, performance, maintainability)
- Surfacing hidden assumptions
- Calibrating confidence in specifications
- Iterative improvement process

## Quick Start (Run in this repo)
 
You can run this example immediately using the `--protocol` override:
 
1. **Draft**:
   ```bash
   specfirst --protocol examples/spec-review/protocol.yaml draft
   ```
 
2. **Finalize**:
   ```bash
   specfirst --protocol examples/spec-review/protocol.yaml finalize
   ```
 
## Setup (For a new project)
 
To use this protocol in your own project:
 
1. Create a new directory and initialize:
   ```bash
   mkdir my-spec && cd my-spec
   specfirst init
   ```
 
2. Copy the protocol and templates:
   ```bash
   cp /path/to/specfirst/examples/spec-review/protocol.yaml .specfirst/protocols/
   cp -r /path/to/specfirst/examples/spec-review/templates/* .specfirst/templates/
   ```
 
3. Update config (optional) or use the flag:
   ```bash
   # Option A: Edit .specfirst/config.yaml to set protocol: spec-review
   # Option B: Use flag
   specfirst --protocol spec-review draft
   ```

## Workflow

### 1. Create Initial Draft

Generate the draft prompt:
```bash
specfirst draft | claude -p > spec-draft.md
```

Complete the stage:
```bash
specfirst complete draft ./spec-draft.md
```

### 2. Surface Hidden Assumptions

Identify what you're assuming but haven't stated:
```bash
specfirst assumptions ./spec-draft.md | claude -p > assumptions-found.md
```

Review the output and update your draft with explicit assumptions.

### 3. Run Role-Based Reviews

Get different perspectives on your spec:

**Security review:**
```bash
specfirst review ./spec-draft.md --persona security | claude -p > security-review.md
```

**Performance review:**
```bash
specfirst review ./spec-draft.md --persona performance | claude -p > performance-review.md
```

**Maintainability review:**
```bash
specfirst review ./spec-draft.md --persona maintainer | claude -p > maintainer-review.md
```

Address concerns in your draft.

### 4. Failure Mode Analysis

Identify what could go wrong:
```bash
specfirst failure-modes ./spec-draft.md | claude -p > failure-analysis.md
```

Add risk mitigation to your spec based on findings.

### 5. Calibrate Confidence

Identify areas where you're uncertain:
```bash
specfirst calibrate ./spec-draft.md --mode confidence | claude -p > confidence-report.md
```

Strengthen low-confidence areas or mark them as open questions.

### 6. Check for Ambiguity

Surface vague language:
```bash
specfirst calibrate ./spec-draft.md --mode uncertainty | claude -p > ambiguity-check.md
```

Clarify any ambiguous statements.

### 7. Finalize

Incorporate all feedback and create final spec:
```bash
specfirst finalize | claude -p > spec-final.md
specfirst complete finalize ./spec-final.md
```

### 8. Validate

Run quality checks:
```bash
specfirst check
specfirst lint
```

## Timeline

**Quick review**: 30-45 minutes  
**Thorough review**: 1-2 hours

## Cognitive Commands Reference

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `assumptions` | Surface hidden assumptions | Before finalizing any spec |
| `review --persona X` | Get perspective-specific feedback | When you need expert eyes |
| `failure-modes` | Identify what could go wrong | For risky or complex features |
| `calibrate --mode confidence` | Gauge certainty levels | When uncertain about parts |
| `calibrate --mode uncertainty` | Find vague language | Before stakeholder review |
| `diff old new` | Understand changes | When updating existing spec |
| `trace` | Map spec to code | When implementing |
| `distill --audience X` | Create summaries | For different stakeholders |

## When to Use This

- ✅ Complex features requiring deep thought
- ✅ High-risk projects
- ✅ Specs that will guide long-term work
- ✅ When you want thorough peer review
- ❌ Simple, well-understood changes
