package assets

const DefaultProtocolName = "multi-stage"

const DefaultProtocolYAML = `name: "multi-stage"
version: "2.1"

stages:
  - id: clarify
    name: Requirements Clarification
    type: spec
    intent: exploration
    template: clarify.md
    depends_on: []
    outputs: [requirements.md]
    output:
      format: markdown
      sections:
        - Problem Statement
        - Users & Primary Use Cases
        - In Scope
        - Out of Scope / Non-Goals
        - Acceptance Criteria
        - Constraints
        - Open Questions & Assumptions

  - id: design
    name: System Design
    type: spec
    intent: decision
    template: design.md
    depends_on: [clarify]
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
    depends_on: [clarify, design]
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
    depends_on: [clarify, design, decompose]
    inputs: [requirements.md, design.md, tasks.yaml]
    source: decompose
    outputs: []
`

const ClarifyTemplate = `# {{ .StageName }} - {{ .ProjectName }}

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
You are performing a requirements clarification step for a software-building task.

Your goal is to transform the user's request into a clear, bounded set of requirements
that can be safely designed, decomposed, and implemented without guesswork.

## Rules
- Do NOT design the solution.
- Do NOT propose architecture or implementation details.
- Focus only on scope, intent, constraints, and definition of done.
- Prefer explicit statements over assumptions.
- If something is unclear, surface it as an open question.
- If you must proceed without an answer, state the assumption explicitly.
- Keep the output concise and structured.

## Fast-Path Check (Internal)
Before writing the full document, determine whether the user input already contains:
- Clearly defined scope
- Explicit acceptance criteria or definition of done
- Known constraints or stated assumptions

If all are present:
- Use the FAST-PATH.
- Produce a compressed version of the document using the standard headings.
- Begin the document with the line: "FAST-PATH USED".

If any are missing:
- Perform full clarification.
- Begin the document with the line: "FULL CLARIFICATION REQUIRED".

## Output Format
Produce one document named "requirements.md" with the following sections,
in this exact order and with these exact headings:

### 1. Problem Statement
- 1-3 sentences describing the problem being solved and why.

### 2. Users & Primary Use Cases
- Bullet list of user types and what they need to do.

### 3. In Scope
- Explicit list of what will be built or changed.

### 4. Out of Scope / Non-Goals
- Explicit exclusions.
- This section is mandatory.

### 5. Acceptance Criteria
- Testable conditions.
- Use checklists or Given/When/Then where possible.

### 6. Constraints
- Technical constraints (stack, APIs, data, performance).
- Non-technical constraints (security, compliance, timelines).

### 7. Open Questions & Assumptions
- Blocking questions that affect scope or behavior.
- If unanswered, list the assumption being made.

## Stopping Conditions
If there are unresolved blocking questions and no safe assumptions can be made:
- Stop.
- Ask the user the minimum number of questions required to proceed.
- Do not continue to downstream stages.

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

You are facilitating a SpecFirst workflow. Ask any clarifying questions first,
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
