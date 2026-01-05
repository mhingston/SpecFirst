package assets

const DefaultProtocolName = "multi-stage"

const DefaultProtocolYAML = `name: "multi-stage"
version: "2.0"

stages:
  - id: requirements
    name: Requirements Gathering
    type: spec
    intent: exploration
    template: requirements.md
    depends_on: []
    outputs: [requirements.md]
    output:
      format: markdown
      sections:
        - Goals
        - Non-Goals
        - Assumptions
        - Constraints
        - Open Questions

  - id: design
    name: System Design
    type: spec
    intent: decision
    template: design.md
    depends_on: [requirements]
    inputs: [requirements.md]
    outputs: [design.md]
    prompt:
      intent: design_outline
      rules:
        - Do not propose implementation details
        - Prefer interfaces over concrete choices
    output:
      format: markdown
      sections:
        - Architecture
        - Components
        - Interfaces
        - Trade-offs

  - id: decompose
    name: Task Decomposition
    type: decompose
    intent: planning
    template: decompose.md
    depends_on: [requirements, design]
    inputs: [requirements.md, design.md]
    outputs: [tasks.yaml]
    prompt:
      intent: task_decomposition
      granularity: ticket
      max_tasks: 12
      prefer_parallel: true
      rules:
        - Prefer parallelizable tasks
        - Surface unknowns explicitly
        - Tasks must be independently reviewable
      required_fields:
        - id
        - title
        - goal
        - acceptance_criteria
        - dependencies
        - files_touched
        - risk_level
        - estimated_scope
        - test_plan
    output:
      format: yaml

  - id: implement
    name: Implementation
    type: task_prompt
    intent: execution
    template: implementation.md
    depends_on: [requirements, design, decompose]
    inputs: [requirements.md, design.md, tasks.yaml]
    source: decompose
    outputs: []
`

const RequirementsTemplate = `# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- if .Inputs }}
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}
{{- else }}
(No prior artifacts)
{{- end }}

## Task
Gather requirements, ask clarifying questions, and enumerate constraints.

{{- if .Outputs }}
## Output Requirements
{{- range .Outputs }}
- {{ . }}
{{- end }}
{{- end }}
`

const DesignTemplate = `# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- if .Inputs }}
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}
{{- end }}

## Task
Produce a system design that satisfies the requirements. Make binding decisions and note trade-offs.

{{- if .Outputs }}
## Output Requirements
{{- range .Outputs }}
- {{ . }}
{{- end }}
{{- end }}
`

const ImplementationTemplate = `# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- if .Inputs }}
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}
{{- end }}

## Task
Implement the specified system. Produce concrete artifacts only.

{{- if .Outputs }}
## Output Requirements
{{- range .Outputs }}
- {{ . }}
{{- end }}
{{- end }}
`

const DecomposeTemplate = `# {{ .StageName }} - {{ .ProjectName }}

## Context
{{- if .Inputs }}
{{- range .Inputs }}
<artifact name="{{ .Name }}">
{{ .Content }}
</artifact>

{{- end }}
{{- end }}

## Task
You are decomposing an approved design into implementation tasks.

### Rules
- Tasks must be independently reviewable.
- If assumptions are required, state them explicitly.
- If any acceptance criteria cannot be verified, add a task to define verification.

## Output Requirements
- Format: YAML
- Each task must include: id, title, goal, acceptance_criteria, dependencies, files_touched, risk_level, estimated_scope, test_plan
`

const DefaultConfigTemplate = `project_name: %s
protocol: %s
language: ""
framework: ""


custom_vars: {}
constraints: {}
`

const InteractiveTemplate = `# Interactive Spec Session - {{ .ProjectName }}

You are facilitating a spec-first workflow. Ask any clarifying questions first,
then produce the requested artifacts for each stage in order.

## Workflow
Protocol: {{ .ProtocolName }}

{{- range .Stages }}
- {{ .ID }}: {{ .Name }} (intent: {{ .Intent }})
  outputs:
  {{- if .Outputs }}
  {{- range .Outputs }}
    - {{ . }}
  {{- end }}
  {{- else }}
    - (none)
  {{- end }}
{{- end }}

## Instructions
- Start with any questions needed to clarify requirements or constraints.
- Then generate artifacts for each stage, in order, using the filenames listed.
- Do not invent new stages or outputs.
- Keep outputs structured and concise.

## Project Context
- Language: {{ .Language }}
- Framework: {{ .Framework }}
{{- if .Constraints }}
- Constraints:
{{- range $key, $value := .Constraints }}
  - {{ $key }}: {{ $value }}
{{- end }}
{{- end }}
{{- if .CustomVars }}
- Custom Vars:
{{- range $key, $value := .CustomVars }}
  - {{ $key }}: {{ $value }}
{{- end }}
{{- end }}
`

// Calibrate templates for epistemic annotation
// Note: Calibrate templates have been moved to internal/prompts/epistemic-calibration.md
