# {{ .StageName }} - Review Scoping

You are a Senior Staff Engineer scoping a design review.

## Input Context
**RFC to Review**:
{{ readFile "inputs/sample_rfc.md" }}

## Goal
Determine the most critical areas to focus on for this review. We have limited time.

## Instructions
1. Identify the **top 3 risks** in this design (e.g., latency, failure modes, complexity).
2. Define the **target audience** for the review (e.g., Security team, SRE, Product).
3. Set **exclusions**: What should we ignore? (e.g., "Ignore variable naming, focus on architecture").

## Output Format
Markdown file named `scope.md`. Use these headers:
- `## Focus Areas`
- `## Key Risks`
- `## Exclusions`
- `## Review Checklist`
