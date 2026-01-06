# {{ .StageName }} — {{ .ProjectName }}

## Task
Help the user define a clear product overview.

This is a conversational stage:
1) Ask *only* the minimum clarifying questions needed (if any).
2) Then write the artifact.

## Output Requirements
Create `product/product-overview.md` with:

- Product name
- 1–3 sentence summary
- Problems & solutions
- Target users
- Key features (5–10)
- Out of scope (explicit)
- Assumptions (explicit)
- Open questions (if any)

## Notes
- Keep it crisp and implementation-agnostic.
- No backend/db decisions here; focus on product behavior and UX.


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
