# Common Workflows

This guide shows practical patterns for using SpecFirst in real-world scenarios.

---

## Quick Fix Workflow

**Use case**: Fixing a bug or small enhancement with clear scope.

```bash
# 1. Initialize and analyze the issue
specfirst init
specfirst requirements > fix.prompt.txt
# Use the prompt to define the issue, then save output as requirements.md

# 2. Complete requirements
specfirst complete requirements ./requirements.md --prompt-file fix.prompt.txt

# 3. Surface potential issues before implementing
specfirst failure-modes ./requirements.md | claude -p
# Review failure modes, adjust requirements if needed

# 4. Generate implementation prompt and execute
specfirst implementation | claude -p > implementation.txt
# Review and apply changes

# 5. Complete and validate
specfirst complete implementation --prompt-file implementation.txt
specfirst check
```

**Timeline**: 10-15 minutes  
**Output**: Documented fix with failure analysis

---

## Full Feature Workflow

**Use case**: Building a new feature with multiple tasks and team coordination.

```bash
# 1. Gather requirements
specfirst init
specfirst requirements | claude -p > requirements.md
specfirst complete requirements ./requirements.md

# 2. Create detailed design
specfirst design | claude -p > design.md
specfirst complete design ./design.md

# 3. Get stakeholder approval
specfirst approve design --role architect --by "Jane Doe"
specfirst approve design --role product --by "Bob Smith"

# 4. Break down into tasks
specfirst decompose | claude -p > tasks.yaml
specfirst complete decompose ./tasks.yaml

# 5. Review task structure
specfirst task
# Lists all tasks from decomposition

# 6. Implement task by task
specfirst task T1 | claude -p
# Repeat for each task, or distribute to team members

# 7. Final validation and archive
specfirst check
specfirst complete-spec --archive --version 1.0 --notes "Initial release"
```

**Timeline**: 1-3 hours  
**Output**: Fully documented, approved, and archived feature

---

## Review & Iterate Workflow

**Use case**: Improving specification quality through structured review.

```bash
# 1. Draft initial spec
specfirst requirements | claude -p > requirements-draft.md

# 2. Surface hidden assumptions
specfirst assumptions ./requirements-draft.md | claude -p > assumptions.md
# Review assumptions, update requirements

# 3. Role-based reviews
specfirst review ./requirements.md --persona security | claude -p > security-review.md
specfirst review ./requirements.md --persona performance | claude -p > perf-review.md
# Address concerns from reviews

# 4. Validate calibration
specfirst calibrate ./requirements.md --mode confidence | claude -p
# Identify low-confidence areas, strengthen them

# 5. Final check and commit
specfirst complete requirements ./requirements.md
specfirst lint
```

**Timeline**: 30-60 minutes  
**Output**: High-quality, reviewed specification

---

## Change Impact Analysis Workflow

**Use case**: Understanding effects of modifying an existing spec.

```bash
# 1. Compare versions
specfirst diff requirements-v1.md requirements-v2.md | claude -p > change-analysis.md
# Review what changed and why

# 2. Trace to existing code
specfirst trace ./requirements-v2.md | claude -p > trace-report.md
# Identify which code needs updating

# 3. Assess test impact
specfirst test-intent ./requirements-v2.md | claude -p > test-intent.md
# Update test strategy based on changes

# 4. Archive old version and proceed
specfirst archive restore v1.0  # If needed to reference
```

**Timeline**: 20-30 minutes  
**Output**: Impact analysis and migration plan

---

## Team Collaboration Workflow

**Use case**: Multiple developers working on different tasks from decomposition.

```bash
# Team Lead:
specfirst init
specfirst requirements | claude -p > requirements.md
specfirst complete requirements ./requirements.md
specfirst decompose | claude -p > tasks.yaml
specfirst complete decompose ./tasks.yaml

# Share task list with team
specfirst task  # Shows all available tasks

# Developer A:
specfirst task T1 | claude -p  # Work on task T1

# Developer B:
specfirst task T2 | claude -p  # Work on task T2 (parallel)

# Team Lead (after completion):
specfirst check  # Validate all outputs
specfirst complete-spec --archive --version 0.1.0
```

**Timeline**: Varies by team size  
**Output**: Coordinated, parallel development

---

## CI/CD Integration Workflow

**Use case**: Automated quality checks in build pipeline.

```bash
# In your CI pipeline (e.g., .github/workflows/spec-check.yml):

- name: Check spec quality
  run: |
    specfirst check --fail-on-warnings
    
- name: Validate all stages complete
  run: |
    specfirst complete-spec --warn-only || exit 1
    
- name: Archive on release
  if: startsWith(github.ref, 'refs/tags/')
  run: |
    specfirst complete-spec --archive --version ${{ github.ref_name }}
```

**Purpose**: Enforce specification discipline automatically  
**Output**: Build failures on spec drift or missing stages

---

## Distill for Stakeholders Workflow

**Use case**: Communicating technical specs to different audiences.

```bash
# After completing design
specfirst complete design ./design.md

# Generate executive summary
specfirst distill ./design.md --audience exec > exec-summary.md

# Generate implementation guide for developers
specfirst distill ./design.md --audience implementer > dev-guide.md

# Generate QA brief
specfirst distill ./design.md --audience qa > qa-brief.md
```

**Timeline**: 5 minutes  
**Output**: Audience-appropriate documentation from single source

---

## Best Practices

1. **Always run `specfirst check` before archiving** - Catches missing outputs or approvals
2. **Use `--prompt-file` when completing stages** - Enables prompt hash verification
3. **Archive at milestones** - Creates rollback points for major versions
4. **Leverage cognitive commands iteratively** - Run assumptions/review multiple times as specs evolve
5. **Prefer `--out` for complex prompts** - Easier to review before sending to AI
