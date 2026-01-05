# Change Impact Analysis Prompt

{{- /* Include Prompt Contract at top of every prompt */ -}}
{{- template "prompt-contract.md" . -}}

---

## PURPOSE

Evaluate risk introduced by changes between two specification versions.

---

## INPUT

```
<specification_before>
{{.SpecBefore}}
</specification_before>

<specification_after>
{{.SpecAfter}}
</specification_after>
```

---

## OUTPUT SCHEMA

Produce a YAML document with the following structure:

```yaml
changes:
  - id: C-001
    type: "<ADD | REMOVE | MODIFY>"
    location: "<section or component affected>"
    description: "<what changed>"
    scope: "<LOCAL | CROSS_CUTTING | BREAKING>"
    backward_compatible: <true | false>
    compatibility_detail: "<why compatible or incompatible>"
    risk_level: "<LOW | MEDIUM | HIGH | CRITICAL>"
    affected_components: ["<component>"]
    validation_required: "<what must be tested>"

breaking_changes:
  - id: BC-001
    change_ref: "<C-XXX reference>"
    breaking_reason: "<why this breaks existing behavior>"
    migration_required: <true | false>
    migration_complexity: "<TRIVIAL | MODERATE | SIGNIFICANT | UNKNOWN>"

cross_cutting_impacts:
  - id: CCI-001
    change_ref: "<C-XXX reference>"
    propagation: ["<component affected>"]
    secondary_effects: "<indirect consequences>"

summary:
  total_changes: <N>
  breaking_count: <N>
  high_risk_count: <N>
  backward_compatible_count: <N>
```

---

## CONSTRAINTS

1. Every change MUST be classified by type, scope, and backward compatibility.
2. Breaking changes MUST be enumerated separately with explicit reasoning.
3. Do not invent changes not present in the diff.
4. If changes cannot be determined due to ambiguity, add to `gap_notices`.
5. Limit to 25 changes, 10 breaking changes, 10 cross-cutting impacts.
6. If limits exceeded, append: `<section>_overflow_count: <N>`.

---

## SCOPE DEFINITIONS

| Scope | Definition |
|-------|------------|
| LOCAL | Change affects single component, no external contracts change |
| CROSS_CUTTING | Change affects multiple components or shared abstractions |
| BREAKING | Change invalidates existing contracts, APIs, or assumptions |

---

## RISK LEVEL CRITERIA

| Level | Criteria |
|-------|----------|
| LOW | Additive change, fully backward compatible |
| MEDIUM | Modification with clear migration path |
| HIGH | Breaking change requiring coordinated updates |
| CRITICAL | Breaking change with security, data integrity, or availability implications |

---

## STOP CONDITIONS

- If no changes detected, output: `no_changes_detected: true`.
- If specs are identical, output: `specs_identical: true`.

---
