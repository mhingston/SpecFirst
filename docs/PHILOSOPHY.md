# The SpecFirst Philosophy

SpecFirst is not just a tool; it's a methodology that prioritizes **intentionality** and **cognitive scaffolding** over raw output.

## Thinking as Default

In modern software development, the "code-first" approach often leads to high-entropy changes where the fundamental logic is obscured by implementation details. SpecFirst flips this.

*   **Code is the side effect.** The specification is the primary artifact.
*   **Structure the thought, then the work.** By defining strict stages (protocols) and templates, we force the engineer (and the AI) to slow down and validate assumptions before a single line of code is committed.

## Cognitive Scaffolding

We believe that even the best engineers benefit from "scaffolds"—structured ways of thinking that prevent common failure modes.

*   **Epistemic Calibration:** Moving from "I think this works" to "I know why this works and where it might fail."
*   **Assumptions Extraction:** Explicitly surfacing the "unspoken" requirements that lead to most bugs.
*   **Failure Analysis:** Designing for error states from the beginning, not as an afterthought.

## Intent-Centrism

Tradition version control tracks *what* changed (the diff). SpecFirst tracks *why* it changed (the intent).

By capturing the protocol state—approvals, assumptions, and calibrations—we create a **Canonical History of Intent**. This makes long-lived projects easier to maintain because future maintainers don't just see the code; they see the reasoning that led to it.

## The Role of AI

In SpecFirst, AI is an **adversarial collaborator**. It shouldn't just write code; it should challenge your assumptions, find gaps in your specifications, and help you distill complex logic into verifiable steps.

## Design Principles

The principles below are not incidental — they are design constraints that guide every feature.

> **Litmus Test**: If a proposed feature could change project outcomes without a human making an explicit decision, it does not belong in SpecFirst.

### 1. No Execution

SpecFirst never executes the code it helps specify. It operates entirely in the space of intent, structure, and verification, leaving execution to the developer or external tools (editors, CI, AI CLIs).

### 2. No Automated Planning

SpecFirst does not decide what to do next.

It can generate prompts that help decompose work into tasks, but:

* task lists are human-authored artifacts
* ordering is human-governed
* dependencies are descriptive, not prescriptive

SpecFirst describes work; it does not plan it.

### 3. No Task State Machines

SpecFirst records facts (e.g. “this stage was marked complete by a human”), but it does not implement a state machine that automatically advances a workflow.

There is no implicit progression, no automatic transitions, and no hidden lifecycle logic. SpecFirst is a record-keeper, not a workflow engine.

> State in SpecFirst represents recorded human attestations, not automated workflow progression.

### 4. Human Judgment Is the Source of Truth

Whenever judgment is required — “is this task finished?”, “is this design acceptable?”, “does this output meet the intent?” — SpecFirst defers to the human.

Approvals are attestations of human judgment, not the result of automated checks.

### 5. Warnings, Not Enforcement

Validation, linting, and completion checks are advisory by default.

They exist to surface:

* ambiguity
* missing information
* weak specifications
* structural inconsistencies

They are meant to **encourage rigor**, not enforce compliance.

### 6. Prompt Infrastructure, Not Automation

SpecFirst provides infrastructure for generating and validating prompts:

* stage prompts
* decomposition prompts
* task-scoped implementation prompts

Everything SpecFirst produces is text.
SpecFirst never acts on that text.

This makes it composable with any editor, any AI tool, and any delivery process — and keeps humans firmly in control.

## Non-Goals

SpecFirst will never:

- Execute prompts or call LLM APIs
- Decide task order or auto-advance workflows
- Score correctness or claim completeness
- Make decisions without explicit human attestation
