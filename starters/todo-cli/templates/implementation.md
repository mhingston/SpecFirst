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
