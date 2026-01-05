# AI-Facing Distillation Prompt

{{- /* Include Prompt Contract at top of every prompt */ -}}
{{- template "prompt-contract.md" . -}}

---

## PURPOSE

Produce a constraint-dense, machine-usable summary of the specification. No narrative prose.

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
distillation:
  invariants:
    - id: INV-001
      statement: "<condition that must always hold>"
      scope: "<SYSTEM | COMPONENT | DATA | INTERFACE>"
      violation_consequence: "<what breaks if violated>"

  constraints:
    - id: CON-001
      type: "<MUST | MUST_NOT | SHOULD | MAY>"
      statement: "<the constraint>"
      applies_to: "<component or scope>"
      enforcement: "<COMPILE_TIME | RUNTIME | EXTERNAL | MANUAL>"

  decision_boundaries:
    - id: DB-001
      condition: "<triggering condition>"
      outcome_true: "<behavior when condition is true>"
      outcome_false: "<behavior when condition is false>"
      edge_cases: ["<edge case>"]

  error_contracts:
    - id: EC-001
      trigger: "<what causes the error>"
      error_type: "<classification>"
      handling: "<PROPAGATE | RECOVER | FAIL_FAST | LOG_ONLY>"
      caller_expectation: "<what caller must handle>"

  data_invariants:
    - id: DI-001
      entity: "<data entity or type>"
      constraint: "<invariant on this data>"
      validation_point: "<where validated>"

  interface_contracts:
    - id: IC-001
      interface: "<interface name>"
      preconditions: ["<condition>"]
      postconditions: ["<condition>"]
      error_modes: ["<error>"]

metadata:
  distillation_coverage: "<COMPLETE | PARTIAL | MINIMAL>"
  ambiguous_sections: ["<section with insufficient precision>"]
  underconstrained_components: ["<component lacking constraints>"]
```

---

## CONSTRAINTS

1. No narrative prose. Use only structured fields.
2. Each invariant/constraint MUST be independently verifiable.
3. Decision boundaries MUST be complete (both outcomes specified).
4. Error contracts MUST specify caller expectations.
5. Do not invent constraints not implied by the spec.
6. Limit each section to 15 items.
7. If limits exceeded, append: `<section>_overflow_count: <N>`.

---

## CONSTRAINT TYPE DEFINITIONS

| Type | Meaning |
|------|---------|
| MUST | Required for correctness. Violation is a defect. |
| MUST_NOT | Prohibited. Occurrence is a defect. |
| SHOULD | Recommended. Deviation requires justification. |
| MAY | Optional. No enforcement. |

---

## STOP CONDITIONS

- If spec contains no verifiable constraints, output: `no_constraints_found: true`.
- If spec is purely narrative with no technical content, output: `narrative_only: true`.

---
