# RFC Review Starter

A staged workflow for reviewing technical designs (RFCs/Design Docs) using SpecFirst.

## Workflow

1.  **Scope**: Reviewer sets the focus areas and exclusions.
2.  **Review**: Principal Engineer persona performs the deep dive.
3.  **Skeptic**: Red Team persona challenges the review findings.
4.  **Report**: Synthesized final report with verdicts.

## Usage

```bash
specfirst init --starter rfc-review

# 1. Scope
# Define reviewer focus and exclusions
opencode run "$(specfirst scope)" > scope.md
specfirst complete scope ./scope.md

# 2. Review
# Principal Engineer persona performs deep dive
opencode run "$(specfirst review)" > review.md
specfirst complete review ./review.md

# 3. Skeptic
# Red Team persona challenges findings
opencode run "$(specfirst skeptic)" > skeptic_findings.md
specfirst complete skeptic ./skeptic_findings.md

# 4. Report
# Synthesize final verdicts
opencode run "$(specfirst report)" > final_review.md
specfirst complete report ./final_review.md
```
