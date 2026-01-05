# Failure Modes Analysis Prompt

{{- /* Include Prompt Contract at top of every prompt */ -}}
{{- template "prompt-contract.md" . -}}

---

## PURPOSE

Enumerate how the system described in the specification can fail, including partial failures, silent failures, and misuse scenarios.

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
failure_modes:
  - id: FM-001
    description: "<what fails and how>"
    category: "<one of: FUNCTIONAL | DATA | INTEGRATION | RESOURCE | SECURITY | OPERATOR | TIMING>"
    trigger: "<condition or sequence that causes failure>"
    detection:
      detectable: <true | false>
      detection_method: "<how failure is observed, or 'silent' if undetectable>"
      detection_latency: "<immediate | delayed | post-mortem | never>"
    impact:
      scope: "<one of: LOCAL | PROPAGATING | SYSTEM-WIDE>"
      severity: "<one of: DEGRADED | PARTIAL_OUTAGE | FULL_OUTAGE | DATA_LOSS | SECURITY_BREACH>"
      affected_components: ["<component>"]
    recovery:
      automatic: <true | false>
      manual_steps: "<required intervention, or 'none'>"
    source_reference: "<section of spec relevant to this failure>"

misuse_scenarios:
  - id: MS-001
    actor: "<who misuses: USER | OPERATOR | EXTERNAL_SYSTEM | ATTACKER>"
    action: "<the misuse action>"
    consequence: "<resulting failure or harm>"
    spec_gap: "<what the spec fails to address>"

silent_failures:
  - id: SF-001
    description: "<failure that produces no observable error>"
    consequence: "<downstream effect>"
    detection_gap: "<why this is undetectable>"

gap_notices:
  - "<areas where spec provides insufficient detail for failure analysis>"
```

---

## CONSTRAINTS

1. Include at least one failure mode per category present in the spec.
2. Silent failures MUST be enumerated separately.
3. Misuse scenarios MUST include both accidental and intentional misuse.
4. Do not invent failures for components not mentioned in the spec.
5. Limit to 15 failure modes, 5 misuse scenarios, 5 silent failures.
6. If limits exceeded, append: `<section>_overflow_count: <N>`.

---

## CATEGORY DEFINITIONS

| Category | Definition |
|----------|------------|
| FUNCTIONAL | Core feature does not work as specified |
| DATA | Data corruption, loss, or inconsistency |
| INTEGRATION | Failure at system boundary or external dependency |
| RESOURCE | Exhaustion of compute, memory, storage, connections |
| SECURITY | Unauthorized access, data exposure, integrity violation |
| OPERATOR | Failure caused by incorrect configuration or operation |
| TIMING | Race conditions, timeouts, ordering violations |

---

## STOP CONDITIONS

- Stop after limits are reached.
- If spec describes fewer than 3 components, state: `limited_scope: true`.

---
