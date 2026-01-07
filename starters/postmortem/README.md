# Incident Postmortem Starter

A workflow to generate a blameless postmortem draft from incident logs and timelines.

## Workflow

1.  **Analyze**: SRE persona analyzes logs/timeline and produces a structured `postmortem.md`.

## Usage

```bash
specfirst init --starter postmortem

# 1. Analyze
# Generates the postmortem draft based on inputs
opencode run "$(specfirst analyze)" > postmortem.md

# 2. Complete
# Register the artifact to close the loop
specfirst complete analyze ./postmortem.md
```
