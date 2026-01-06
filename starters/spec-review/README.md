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
   opencode run "$(specfirst --protocol starters/spec-review/protocol.yaml draft)"
   ```
 
2. **Finalize**:
   ```bash
   opencode run "$(specfirst --protocol starters/spec-review/protocol.yaml finalize)"
   ```
 
## Setup (For a new project)
 
To use this protocol in your own project:
 
1. Create a new directory and initialize it with Git:
   ```bash
   mkdir my-spec && cd my-spec
   git init
   ```

2. Initialize SpecFirst with the `spec-review` starter:
   ```bash
   specfirst init --starter spec-review
   ```
 
## Workflow

### 1. Create Initial Draft

Generate the draft prompt:
```bash
opencode run "$(specfirst draft)" > spec-draft.md
```

Complete the stage:
```bash
specfirst complete draft ./spec-draft.md
```

### 2. Surface Hidden Assumptions

Identify what you're assuming but haven't stated:
```bash
opencode run "$(specfirst assumptions ./spec-draft.md)" > assumptions-found.md
```

Review the output and update your draft with explicit assumptions.

### 3. Run Role-Based Reviews

Get different perspectives on your spec:

**Security review:**
```bash
opencode run "$(specfirst review ./spec-draft.md --persona security)" > security-review.md
```

**Performance review:**
```bash
opencode run "$(specfirst review ./spec-draft.md --persona performance)" > performance-review.md
```

**Maintainability review:**
```bash
opencode run "$(specfirst review ./spec-draft.md --persona maintainer)" > maintainer-review.md
```

Address concerns in your draft.

### 4. Failure Mode Analysis

Identify what could go wrong:
```bash
opencode run "$(specfirst failure-modes ./spec-draft.md)" > failure-analysis.md
```

Add risk mitigation to your spec based on findings.

### 5. Calibrate Confidence

Identify areas where you're uncertain:
```bash
opencode run "$(specfirst calibrate ./spec-draft.md --mode confidence)" > confidence-report.md
```

Strengthen low-confidence areas or mark them as open questions.

### 6. Check for Ambiguity

Surface vague language:
```bash
opencode run "$(specfirst calibrate ./spec-draft.md --mode uncertainty)" > ambiguity-check.md
```

Clarify any ambiguous statements.

### 7. Finalize

Incorporate all feedback and create final spec:
```bash
opencode run "$(specfirst finalize)" > spec-final.md
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
