# The Spec-First Philosophy

Spec-First is not just a tool; it's a methodology that prioritizes **intentionality** and **cognitive scaffolding** over raw output.

## Thinking as Default

In modern software development, the "code-first" approach often leads to high-entropy changes where the fundamental logic is obscured by implementation details. Spec-First flips this.

*   **Code is the side effect.** The specification is the primary artifact.
*   **Structure the thought, then the work.** By defining strict stages (protocols) and templates, we force the engineer (and the AI) to slow down and validate assumptions before a single line of code is committed.

## Cognitive Scaffolding

We believe that even the best engineers benefit from "scaffolds"—structured ways of thinking that prevent common failure modes.

*   **Epistemic Calibration:** Moving from "I think this works" to "I know why this works and where it might fail."
*   **Assumptions Extraction:** Explicitly surfacing the "unspoken" requirements that lead to most bugs.
*   **Failure Analysis:** Designing for error states from the beginning, not as an afterthought.

## Intent-Centrism

Tradition version control tracks *what* changed (the diff). Spec-First tracks *why* it changed (the intent).

By capturing the protocol state—approvals, assumptions, and calibrations—we create a **Canonical History of Intent**. This makes long-lived projects easier to maintain because future maintainers don't just see the code; they see the reasoning that led to it.

## The Role of AI

In Spec-First, AI is an **adversarial collaborator**. It shouldn't just write code; it should challenge your assumptions, find gaps in your specifications, and help you distill complex logic into verifiable steps.
