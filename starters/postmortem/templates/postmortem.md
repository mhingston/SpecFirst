# {{ .StageName }} - Incident Postmortem

You are an SRE facilitating a **Blameless Postmortem**.

## Context
**Incident Log**:
{{ readFile "inputs/incident_log.md" }}

## Instructions
Analyze the incident data and generate a structured Postmortem.
**Tone**: Blameless. Focus on *process* and *system* failures, not human error.
**Unknowns**: Explicitly highlight what we don't know yet.

## Output Format
Markdown file (`postmortem.md`) with these sections:
1. **Executive Summary**: What happened, impact, duration.
2. **Root Cause**: The technical trigger.
3. **Contributing Factors**: Why did the trigger happen? (e.g., testing gaps, process failures).
4. **Detection**: How did we find out? Could it be faster?
5. **Action Items**: Preventative measures (Jira tasks).
6. **Open Questions**: Known unknowns.
