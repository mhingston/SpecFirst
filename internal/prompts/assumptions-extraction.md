# Assumptions Extraction Prompt

{{- /* Include Prompt Contract at top of every prompt */ -}}
{{- template "prompt-contract.md" . -}}

---

## PURPOSE

Extract implicit assumptions from the specification that, if false, could invalidate the design or implementation.

---

## INPUT

```
<specification>
{{.Spec}}
</specification>
```

---

## OUTPUT SCHEMA

Produce a YAML document with the following structure:

```yaml
assumptions:
  - id: A-001
    statement: "<atomic assumption in declarative form>"
    category: "<one of: ENVIRONMENTAL | BEHAVIORAL | TECHNICAL | TEMPORAL | RESOURCE | SECURITY | DATA>"
    source_reference: "<quote or section from spec that implies this assumption>"
    impact_if_false: "<specific consequence if assumption does not hold>"
    detection_signal: "<observable indicator that assumption may be false>"
    validation_action: "<concrete step to verify or falsify this assumption>"
    confidence: "<HIGH | MEDIUM | LOW>"

gap_notices:
  - "<description of missing information that prevents assumption extraction>"
```

---

## CONSTRAINTS

1. Each assumption MUST be atomic (one claim per entry).
2. Each assumption MUST be traceable to a source reference in the spec.
3. Do not invent assumptions about implementation details not implied by the spec.
4. If the spec is too vague to extract assumptions, populate `gap_notices` instead.
5. Limit output to the top 20 assumptions by impact severity.
6. If more than 20 exist, append a summary count: `additional_assumptions_count: <N>`.

---

## STOP CONDITIONS

- Stop after 20 assumptions.
- If fewer than 5 assumptions can be extracted, state: `insufficient_signal: true`.

---

## CATEGORY DEFINITIONS

| Category | Definition |
|----------|------------|
| ENVIRONMENTAL | Assumptions about runtime environment, infrastructure, or deployment context |
| BEHAVIORAL | Assumptions about user behavior, usage patterns, or operator actions |
| TECHNICAL | Assumptions about system capabilities, dependencies, or integrations |
| TEMPORAL | Assumptions about timing, sequencing, or duration |
| RESOURCE | Assumptions about availability of compute, memory, storage, or budget |
| SECURITY | Assumptions about trust boundaries, authentication, or threat model |
| DATA | Assumptions about data format, volume, quality, or availability |

---
