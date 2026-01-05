# Epistemic Calibration Prompt

{{- /* Include Prompt Contract at top of every prompt */ -}}
{{- template "prompt-contract.md" . -}}

---

## PURPOSE

Classify statements in the specification by epistemic status: what is known, what is assumed, what is uncertain, and what is unknown.

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
epistemic_map:
  known:
    - id: K-001
      statement: "<statement explicitly defined in spec>"
      source_reference: "<exact quote or section>"
      confidence: 1.0

  assumed:
    - id: AS-001
      statement: "<statement implied but not explicit>"
      source_reference: "<section that implies this>"
      confidence: "<0.6 - 0.9>"
      invalidation_condition: "<what would prove this false>"

  uncertain:
    - id: U-001
      statement: "<statement with ambiguous or conflicting signals>"
      source_reference: "<relevant section>"
      confidence: "<0.3 - 0.6>"
      ambiguity_source: "<why this is uncertain>"
      invalidation_condition: "<what would prove this false>"

  unknown:
    - id: UK-001
      statement: "<question the spec does not address>"
      required_for: "<what decision or component needs this>"
      blocking: "<true | false>"

confidence_distribution:
  known_count: <N>
  assumed_count: <N>
  uncertain_count: <N>
  unknown_count: <N>
```

---

## CONSTRAINTS

1. Each item MUST include a source reference or state `source_reference: "not present"`.
2. Confidence values MUST be numeric (0.0 - 1.0) for known/assumed/uncertain.
3. Do not collapse categories. If something is uncertain, do not list it as assumed.
4. Do not provide recommendations for resolving unknowns unless explicitly requested.
5. Limit each category to 15 items.
6. If more than 15 exist in a category, append: `<category>_overflow_count: <N>`.

---

## CONFIDENCE TIERS

| Tier | Range | Meaning |
|------|-------|---------|
| Known | 1.0 | Explicitly stated in spec with no ambiguity |
| High Confidence | 0.8 - 0.9 | Strongly implied, single interpretation |
| Medium Confidence | 0.5 - 0.7 | Implied but multiple interpretations possible |
| Low Confidence | 0.3 - 0.4 | Weak signal, significant ambiguity |
| Unknown | N/A | Not addressed in spec |

---

## STOP CONDITIONS

- Stop after processing all sections of the spec.
- If spec is under 100 words, state: `insufficient_content: true`.

---
