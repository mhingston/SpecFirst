# SpecFirst Blind Spots Remediation Plan

## Purpose
Address identified blind spots in SpecFirst by adding:
1) Lightweight early-stage workflows
2) Epistemic state (assumptions/unknowns/risks/disputes)
3) Richer approvals (attestations with scope, conditions, dissent)
4) First-class forks/tracks for parallel futures
5) Protocol-level ambiguity gates (carry-forward rules)

This plan prioritizes changes that preserve SpecFirst’s strengths:
- reproducibility
- lifecycle governance
- prompt-based collaboration (not forced codegen)
- safety-by-default

---

## Guiding Principles

### 1. Epistemic Honesty
Specs must explicitly represent uncertainty and disagreement, not erase them.

### 2. Progressive Formalization
Users should start messy and only later “harden” specs into approval-gated stages.

### 3. Accountable Governance
Approvals should encode scope, rationale, and conditions—not just a boolean.

### 4. Parallel Futures Are Normal
Exploration needs forks; execution needs convergence.

### 5. Human + AI Co-Workflows
Commands should generate structured prompts and checklists rather than forcing automation.

---

## Workstreams

## A) Sketch Mode (lightweight / pre-spec)

### Goal
Enable fast ideation without upfront protocol overhead, while keeping a clean upgrade path.

### Deliverables
- New stage type: `sketch`
- `specfirst init --mode sketch`
- `specfirst promote sketch --to <protocol>`

### Behavior
- Sketch artifacts are marked non-authoritative.
- Promotion maps sketch into the first formal stage artifact (e.g., requirements/analysis).

### Acceptance
- A user can go from zero → sketch → promoted protocol without deleting work.

---

## B) Epistemic Ledger in State

### Goal
Make assumptions, unknowns, decisions, risks, and disputes first-class and persistent.

### State Model (proposal)
Extend `state.json` with:

```json
{
  "epistemics": {
    "assumptions": [
      { "id": "A1", "text": "...", "status": "open|validated|invalidated", "owner": "", "created_at": "" }
    ],
    "open_questions": [
      { "id": "Q1", "text": "...", "tags": ["security"], "status": "open|resolved|deferred" }
    ],
    "decisions": [
      { "id": "D1", "text": "...", "rationale": "...", "alternatives": [], "status": "proposed|accepted|reversed" }
    ],
    "risks": [
      { "id": "R1", "text": "...", "severity": "low|medium|high", "mitigation": "", "status": "open|mitigated|accepted" }
    ],
    "disputes": [
      {
        "id": "X1",
        "topic": "...",
        "positions": [{ "owner": "", "claim": "" }],
        "status": "open|resolved"
      }
    ],
    "confidence": {
      "overall": "low|medium|high",
      "by_stage": { "requirements": "medium", "design": "low" }
    }
  }
}
```

### CLI (proposal)

* `specfirst assume add|list|close`
* `specfirst question add|list|resolve|defer`
* `specfirst decision add|accept|reverse`
* `specfirst risk add|mitigate|accept`
* `specfirst dispute add|update|resolve`

### Template Integration

Templates should automatically include:

* open questions (with tags)
* high risks
* pending disputes
* confidence indicators

### Acceptance

* “What do we not know?” is visible without reading every artifact.

---

## C) Approvals → Attestations

### Goal

Replace binary approvals with scoped attestations that can be conditional and reasoned.

### Model (proposal)

```json
{
  "attestations": {
    "design": [
      {
        "role": "architect",
        "scope": ["architecture", "security-boundaries"],
        "status": "approved|approved_with_conditions|rejected|needs_changes",
        "conditions": ["Add rate limiting plan"],
        "rationale": "Meets scalability goals; auth details pending.",
        "dissenting_opinions": [],
        "created_at": ""
      }
    ]
  }
}
```

### CLI (proposal)

* `specfirst attest --stage design --role architect --status approved_with_conditions --condition "..."`
* `specfirst approvals status` prints blocking items + conditions to resolve

### Acceptance

* Blocking approvals tell you exactly what to do next and what the reviewer cared about.

---

## D) Forks / Tracks (parallel futures)

### Goal

Allow exploration branches without losing SpecFirst statefulness and provenance.

### Option A (recommended v1)

* `specfirst track create <name>`: create a named snapshot workspace
* `specfirst track list`
* `specfirst track diff <a> <b>`: show artifact + state differences
* `specfirst track merge <source> --strategy manual`: produce merge prompt + checklist

### Acceptance

* Teams can hold two designs in parallel and compare them with a first-class command.

---

## E) Ambiguity Gates (carry-forward rules)

### Goal

Support intentionally unresolved questions while preventing dangerous ambiguity.

### Protocol extensions (proposal)

Per stage:

* `max_open_questions`
* `must_resolve_tags`
* `max_high_risks_unmitigated`

Example:

```yaml
stages:
  - id: design
    max_open_questions: 5
    must_resolve_tags: [security, compliance]
    max_high_risks_unmitigated: 0
```

### Validation

A stage can complete if:

* outputs exist AND
* epistemic gates pass

### Acceptance

* Exploration stays fluid, but safety-critical unknowns block progress.

---

## Roadmap

### Phase 1: Ledger foundation

* State schema changes + helpers
* CLI for epistemic items
* Template rendering of ledger context

### Phase 2: Attestations

* Replace/extend approvals model
* Update gating + summaries

### Phase 3: Sketch mode + promotion

* Add stage type + init mode
* Promote workflow

### Phase 4: Tracks

* Snapshot + diff + merge prompt tooling

### Phase 5: Ambiguity gates

* Protocol extensions + validate integration

---

## Success Metrics

* Onboarding time from init → useful artifact < 10 minutes (sketch path)
* Reduced “unknown unknowns” at implementation stage (tracked via ledger closure rate)
* Approvals become actionable (conditions resolved time)
* Parallel designs supported without chaos (tracks adoption)

---

## Open Questions

* Should ledger items be stored only in `state.json` or also mirrored into human-readable artifacts?
* How tightly should track management integrate with git branches?
* Do we want per-artifact confidence or per-stage only?
