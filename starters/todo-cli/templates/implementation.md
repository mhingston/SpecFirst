# Implementation Task: {{ .ProjectName }}

## Task Details
(This is a task-scoped prompt. The specific task details will be appended by specfirst.)

## Context
{{- range .Inputs }}
### {{ .Name }}
{{ .Content }}
{{- end }}

## Instructions
- **Output Formatting**: For this task, **provide the raw code only**. Skip all preamble, conversation, and Markdown code fences.
- Implement ONLY what is required for the specific task.
- Follow the project language and framework standards.


## Output Format Constraints
CRITICAL: You must output ONLY the raw markdown content for the file.
- Do NOT include any conversational text (e.g. "Here is the file...", "I will now...").
- Do NOT include markdown code block fences (```markdown ... ```) around the content.
- Start directly with the markdown content (e.g. # Title).
