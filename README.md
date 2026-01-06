# SpecFirst

SpecFirst is a Go CLI for specification-driven workflows that focuses on **prompt quality, clarity, and human judgment**.

It compiles structured prompts from declarative protocols and templates, stores artifacts, and records decisions — but it does **not** plan work, execute prompts, or decide what happens next. SpecFirst is LLM-agnostic: it emits text to stdout and stays out of the execution loop.

Think of SpecFirst as **prompt infrastructure**: a discipline amplifier that helps humans and AI reason clearly together before implementation.

## Features

- Protocol-defined workflow stages with dependency gating.
- **Protocol Composability**: Import and mixin common stages using the `uses` field.
- **Prompt Quality Infrastructure**: Schema validation, ambiguity detection, and structure checks integrated into `lint`.
- **Task Decomposition**: Break down designs into structured units of work via the `decompose` stage.
- **Task-Scoped Prompts**: Generate focused implementation prompts for specific tasks using `specfirst task <id>`.
- Template-based prompt rendering with artifact embedding.
- Durable artifact store with prompt hashing for reproducibility.
- Explicit state tracking with approvals and prompt hashes.

## Philosophy

SpecFirst takes a deliberately different approach to specification-driven workflows.

Many tools in this space focus on **automation**: planning work, advancing stages, executing prompts, or deciding what should happen next. SpecFirst intentionally avoids those responsibilities.

Instead, SpecFirst focuses on a narrower problem:

> **Turning structured human intent into clear, deterministic prompts that humans and AI can reason about together.**

SpecFirst is designed as a **discipline amplifier**, not a process enforcer. It helps you think clearly *before* you act, without automating away judgment, context, or responsibility.

The principles below are not incidental — they are design constraints that guide every feature.

> **Litmus Test**: If a proposed feature could change project outcomes without a human making an explicit decision, it does not belong in SpecFirst.

---

### 1. No Execution

SpecFirst never executes the code it helps specify. It operates entirely in the space of intent, structure, and verification, leaving execution to the developer or external tools (editors, CI, AI CLIs).

---

### 2. No Automated Planning

SpecFirst does not decide what to do next.

It can generate prompts that help decompose work into tasks, but:

* task lists are human-authored artifacts
* ordering is human-governed
* dependencies are descriptive, not prescriptive

SpecFirst describes work; it does not plan it.

---

### 3. No Task State Machines

SpecFirst records facts (e.g. “this stage was marked complete by a human”), but it does not implement a state machine that automatically advances a workflow.

There is no implicit progression, no automatic transitions, and no hidden lifecycle logic. SpecFirst is a record-keeper, not a workflow engine.

> State in SpecFirst represents recorded human attestations, not automated workflow progression.

---

### 4. Human Judgment Is the Source of Truth

Whenever judgment is required — “is this task finished?”, “is this design acceptable?”, “does this output meet the intent?” — SpecFirst defers to the human.

Approvals are attestations of human judgment, not the result of automated checks.

---

### 5. Warnings, Not Enforcement

Validation, linting, and completion checks are advisory by default.

They exist to surface:

* ambiguity
* missing information
* weak specifications
* structural inconsistencies

They are meant to **encourage rigor**, not enforce compliance.

---

### 6. Prompt Infrastructure, Not Automation

SpecFirst provides infrastructure for generating and validating prompts:

* stage prompts
* decomposition prompts
* task-scoped implementation prompts

Everything SpecFirst produces is text.
SpecFirst never acts on that text.

This makes it composable with any editor, any AI tool, and any delivery process — and keeps humans firmly in control.

## Non-Goals

SpecFirst will never:

- Execute prompts or call LLM APIs
- Decide task order or auto-advance workflows
- Score correctness or claim completeness
- Make decisions without explicit human attestation

## Documentation

- [Philosophy](docs/PHILOSOPHY.md): The "why" behind SpecFirst and the cognitive scaffold approach.
- [Architecture & Mechanics](docs/ARCHITECTURE.md): Conceptual overview, state semantics, and archive philosophy.
- [User Guide](docs/GUIDE.md): Detailed "how-to" and workflow examples.
- [Protocol Reference](docs/PROTOCOLS.md): YAML schema and stage definitions.
- [Template Reference](docs/TEMPLATES.md): Guide to authoring stage templates.

## Canonical Examples

Complete, runnable examples with protocols and templates:

- [Todo CLI](starters/todo-cli/README.md): Full workflow walkthrough for a simple CLI application
- [Bug Fix](starters/bug-fix/README.md): Minimal 2-stage workflow for systematic bug fixes
- [API Feature](starters/api-feature/README.md): Complete feature development with approvals and task decomposition
- [Spec Review](starters/spec-review/README.md): Using cognitive scaffold commands to improve specification quality
- [Refactoring](starters/refactoring/README.md): Structured code improvement with risk mitigation and verification
- [Database Migration](starters/database-migration/README.md): Safe schema changes with approval gates and rollback planning
- [Product Design](starters/product-design/README.md): Design flow from product vision to implementation handoff

## Install

### From Source

```bash
make build
```

Or:

```bash
go build .
```

### Add to PATH

```bash
# Option 1: Install to /usr/local/bin (may require sudo)
sudo mv specfirst /usr/local/bin/

# Option 2: Install to ~/bin (add to PATH if not already)
mkdir -p ~/bin
mv specfirst ~/bin/
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc  # or ~/.bashrc

# Option 3: Use make install (installs to $GOBIN or $GOPATH/bin)
make install
```

### From GitHub Releases

Download the pre-built binary for your platform from [Releases](https://github.com/mhingston/SpecFirst/releases), then move it to your PATH as shown above.

## Quick Start

```bash
specfirst init
specfirst status
specfirst requirements
```

## Starter Kits

SpecFirst includes pre-built workflow starter kits for common scenarios. Starters bundle a protocol, templates, and optional skill files into a single package.

### List Available Starters

```bash
specfirst starter list
```

### Create Project with Starter

```bash
mkdir my-project && cd my-project
specfirst init --starter api-feature
```

This creates a workspace with the `api-feature` protocol and templates pre-configured.

### Apply Starter to Existing Workspace

```bash
specfirst init
specfirst starter apply bug-fix               # Apply bug-fix workflow
specfirst starter apply api-feature --force   # Overwrite existing templates
```

### Interactive Selection

```bash
specfirst init --choose
```

This prompts you to select from available starters interactively.

### Available Starters

| Name | Description |
|------|-------------|
| `api-feature` | Full workflow with approvals for API features |
| `bug-fix` | Lightweight 2-stage workflow for bug fixes |
| `database-migration` | Migration workflow with rollback planning |
| `product-design` | Full product design from brief to handover |
| `refactoring` | Code refactoring workflow with risk assessment |
| `spec-review` | Specification review workflow |
| `todo-cli` | Simple example workflow |

## AI CLI Integration

SpecFirst outputs prompts to stdout, making it composable with AI CLIs. 

### Interactive vs. Non-Interactive Modes

Most AI CLIs support a **one-shot (non-interactive)** mode for automation and an **interactive** mode for refinement.

#### 1. Interactive Refinement (Recommended)
To maintain an interactive session where you can refine the output, use your system clipboard or tools that support stdin-to-interactive:

```bash
# Copy prompt to clipboard and paste into your AI tool
specfirst requirements | pbcopy # macOS
specfirst requirements | xclip -sel clip # Linux

# Or use command substitution in the tool's interactive prompt
# (Works if the tool allows starting a session with an initial prompt)
copilot -p "$(specfirst requirements)" --allow-all-tools
```

#### 2. One-Shot / Piped (Non-Interactive)
Use these for quick generations or scripting. Note that flags like `-p` or `--print` usually exit after one response.

```bash
# Claude Code (headless mode)
specfirst requirements | claude -p

# GitHub Copilot (non-interactive)
copilot -p "$(specfirst requirements)" --allow-all-tools

# Gemini CLI (non-interactive prompt)
opencode run "$(specfirst implementation)"

# Gemini CLI (via stdin)
specfirst implementation | opencode run
```

#### 3. Pipelining Back to SpecFirst
You can pipe AI output directly into `specfirst complete` using `-` to read from stdin:

```bash
# Example: Generate requirements with Claude and complete the stage in one go
specfirst requirements | claude -p | specfirst complete requirements -
```

**For any tool that reads from files or requires non-interactive output for redirection**, use the tool's one-shot mode:

```bash
# Gemini one-shot to file
opencode run "$(specfirst requirements)" > requirements.md
```

## End-to-End Example (Default Protocol)

The default protocol (`multi-stage`) ships with `requirements`, `design`, and
`implementation` stages and templates. A typical run looks like this:

```bash
specfirst init

specfirst requirements > requirements.prompt.txt
cat requirements.prompt.txt
# Use the prompt with your LLM and save its output as requirements.md

specfirst complete requirements ./requirements.md --prompt-file requirements.prompt.txt

specfirst design --out design.prompt.txt
# Use the prompt and save output as design.md

specfirst complete design ./design.md --prompt-file design.prompt.txt

specfirst implementation --out implementation.prompt.txt
# Use the prompt and save output as generated code files (e.g., src/main.go)

# Complete the stage (automatically detects changed files):
specfirst complete implementation --prompt-file implementation.prompt.txt

# Task Decomposition and Scoped Prompts
specfirst decompose --out tasks.prompt.txt
# Save output as tasks.yaml or tasks.md (structured tasks)
specfirst complete decompose ./tasks.yaml

# List tasks found in decomposition:
specfirst task

# Generate a prompt for a specific task (automatically searches all stage artifacts):
specfirst task T1 | claude -p

# Validate completion and optionally archive:
specfirst complete-spec --archive --version 2.0
```

## Usage Examples

Initialize a workspace and generate the first prompt:

```bash
specfirst init
specfirst requirements > requirements.prompt.txt
```

Complete a stage and store artifacts:

```bash
specfirst complete requirements ./requirements.md --prompt-file requirements.prompt.txt
```

Render JSON output for tooling:

```bash
specfirst design --format json --out design.prompt.json
```

Generate an interactive prompt to a file:

```bash
specfirst --interactive --out interactive.prompt.txt
```

## Commands

- `specfirst init` initializes `.specfirst/` with defaults.
- `specfirst init --starter <name>` initializes with a specific starter kit.
- `specfirst starter list` lists available starter kits.
- `specfirst starter apply <name>` applies a starter kit to the current workspace.
- `specfirst status` shows current workflow status.
- `specfirst <stage-id>` renders a stage prompt to stdout.
- `specfirst complete <stage-id> <output-files...>` records completion and stores artifacts.
- `specfirst task [task-id]` lists tasks or generates a prompt for a specific task (requires a completed `decompose` stage).
- `specfirst complete-spec [--archive|--warn-only]` validates completion and optionally archives. It is a validation tool, not a strict workflow requirement.
- `specfirst --interactive` generates a meta-prompt for an end-to-end session.
- `specfirst lint` runs non-blocking checks, including **prompt quality and ambiguity detection**.
- `specfirst check [--fail-on-warnings]` runs a **preflight / hygiene report** including all non-blocking validations (lint, tasks, approvals, outputs).
- `specfirst archive <version>` manages workspace archives.
- `specfirst protocol list|show|create` manages protocol definitions.
- `specfirst attest <stage-id> --role <role> --status <status>` records attestations with rationale and conditions.
- `specfirst track create|list|switch|diff|merge` manages parallel futures (tracks).

### Cognitive Scaffold Commands

These commands generate **prompts only** — no state, no enforcement, no AI calls. They shape thinking, not execution.

- `specfirst diff <old-spec> <new-spec>` generates a change-analysis prompt comparing two specification files.
- `specfirst assumptions <spec-file>` generates a prompt to surface hidden assumptions.
- `specfirst review <spec-file> --persona <p>` generates a role-based review prompt. Personas: `security`, `performance`, `maintainer`, `accessibility`, `user`.
- `specfirst failure-modes <spec-file>` generates a failure-first interrogation prompt.
- `specfirst test-intent <spec-file>` generates a test **intent** prompt (not test code).
- `specfirst trace <spec-file>` generates a spec-to-code mapping prompt.
- `specfirst distill <spec-file> --audience <a>` generates an audience-specific summary prompt. Audiences: `exec`, `implementer`, `ai`, `qa`.
- `specfirst calibrate <artifact>` generates a comprehensive epistemic map for judgment calibration.

## Completion Options

- `--prompt-file <path>` hash an explicit prompt file when completing a stage.
- `--force` overwrite an existing stage completion (non-destructive; only removes old artifacts after new ones are successfully stored).

## Stage Execution Options

- `--protocol <path|name>` override active protocol (path to file or name in `.specfirst/protocols`).
- `--format text|json|yaml|shell` output format (default: `text`).

- `--out <file>` write prompt to a file.
- `--max-chars <n>` truncate output.
- `--no-strict` bypass dependency gating.
- `--interactive` generate an interactive meta-prompt.

## Decomposition Options

- `--granularity feature|story|ticket|commit` set task size (default: `ticket`).
- `--max-tasks <n>` limit the number of tasks generated.
- `--prefer-parallel` favor tasks that can be implemented concurrently.
- `--risk-bias conservative|balanced|fast` tune implementation risk (default: `balanced`).

## Complete-Spec Options

- `--archive` create an archive snapshot after completion.
- `--warn-only` report missing stages/approvals as warnings to stderr without failing the command (exit 0).
- `--version <v>` explicit version for the archive.
- `--tag <tag>` tags for the archive (repeatable).
- `--notes <text>` notes for the archive.

## Archive Options

- `archive <version> --tag <tag>` apply tags to the archive (repeatable).
- `archive <version> --notes <text>` add notes to the archive.
- `archive restore <version> --force` overwrite existing workspace data when restoring (strict restore; removes existing workspace data before restore). Restore now fails if required archive directories (like `protocols/` or `templates/`) are missing.
- `archive <version>` requires `.specfirst/protocols/` and `.specfirst/templates/` to exist (run `specfirst init` if missing).

## Track Options
 
 - `track create <name> --notes <text>` create a new track.
 - `track switch <name> --force` restore a track to the current workspace (overwrites existing data).
 - `track merge <source>` generate a merge plan prompt.
 
 ## Attestation Options
 
 - `--role <role>` (required) the role for the attestation.
 - `--status <status>` (required) status: `approved`, `approved_with_conditions`, `needs_changes`, `rejected`.
 - `--rationale <text>` rationale for the decision.
 - `--condition <text>` condition for conditional approval (repeatable).
 - `--by <name>` who attested (defaults to `$USER`).

## Workspace Layout

```
.specfirst/
  artifacts/     # Stage outputs stored by stage ID
  generated/     # Generated files (e.g., compiled prompts)
  protocols/     # Protocol YAML definitions
  templates/     # Prompt templates
  skills/        # Reusable prompt chunks (for readFile helper)
  archives/      # Archived spec versions
  state.json     # Workflow state
  config.yaml    # Project configuration
```

## Defaults

`specfirst init` installs:

- Protocol: `.specfirst/protocols/multi-stage.yaml`
- Templates: `.specfirst/templates/requirements.md`, `design.md`, `implementation.md`
- Config: `.specfirst/config.yaml`

Users can edit or replace these defaults in their workspace.

## Protocol Format

Protocols are YAML DAGs of stages:

```yaml
name: "multi-stage"
version: "2.0"
uses:
  - shared-stages
stages:
  - id: requirements
    name: Requirements Gathering
    type: spec
    template: requirements.md
    outputs: [requirements.md]
    output:
      format: markdown
      sections: [Goals, Constraints]
  - id: design
    name: System Design
    type: spec
    template: design.md
    depends_on: [requirements]
    inputs: [requirements.md]
    outputs: [design.md]
  - id: decompose
    name: Task Decomposition
    type: decompose
    template: decompose.md
    depends_on: [design]
    inputs: [design.md]
    outputs: [tasks.yaml]
    prompt:
      granularity: ticket
      max_tasks: 10
```

### Output Pattern Matching

Output patterns in protocols support single-level wildcards only:

- ✅ `src/*` - matches files directly under `src/`
- ✅ `*.md` - matches markdown files
- ❌ `src/**/*.go` - recursive patterns are **not supported**

For complex directory structures, use flat output organization or enumerate specific files.
Lint will warn if a stage declares wildcard outputs but no stored artifacts match.

### Stage-Qualified Inputs

When the same filename exists in multiple stage artifacts, use stage-qualified paths:

```yaml
inputs:
  - requirements/requirements.md  # Explicit stage
  - design/notes.md
```

## Template Authoring

Templates are Go `text/template` files with full access to the variables listed below. Stage inputs are automatically embedded as artifacts.

## Template Variables

These variables are available to templates:

| Variable | Type | Description |
| --- | --- | --- |
| `StageName` | string | Human-readable stage name. |
| `ProjectName` | string | Project name from `config.yaml` or the working directory. |
| `Inputs` | []Input | Inputs attached to the stage (each has `Name` and `Content`). |
| `Outputs` | []string | Expected output filenames declared in the protocol. |
| `Intent` | string | Stage intent (e.g. `exploration`, `decision`, `execution`, `review`). |
| `Language` | string | Optional project language from config. |
| `Framework` | string | Optional project framework from config. |
| `CustomVars` | map[string]string | Arbitrary user-defined variables from config. |
| `Constraints` | map[string]string | Constraints map from config. |
| `StageType` | string | The type of the current stage (`spec`, `decompose`, etc.). |
| `Prompt` | PromptConfig | The detailed prompt configuration for the stage. |
| `OutputContract` | OutputContract | The expected structure of the stage output. |


Example template:

```markdown
# {{ .StageName }} - {{ .ProjectName }}

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

{{- if .Constraints }}
## Constraints
{{- range $key, $value := .Constraints }}
- {{ $key }}: {{ $value }}
{{- end }}
{{- end }}

## Output Requirements
{{- range .Outputs }}
- {{ . }}
{{- end }}
```

### Template Functions

Templates have access to these helper functions:

| Function | Usage | Description |
| --- | --- | --- |
| `join` | `{{ join .Outputs ", " }}` | Joins a slice with a delimiter |
| `readFile` | `{{ readFile "skill-name.md" }}` | Includes a skill file from `.specfirst/skills/` |

### Skills (Reusable Prompt Chunks)

Skills are reusable markdown files stored in `.specfirst/skills/`. Use them to share common guidance across templates:

```markdown
# In your template:
{{ readFile "spec-writing.md" }}
```

This includes the content of `.specfirst/skills/spec-writing.md` directly in the rendered prompt. Skills help maintain consistency across protocols without copy-pasting.

## Config and State

`config.yaml` sets project metadata and the active protocol. `state.json` tracks
completed stages, approvals, and prompt hashes.

## Build Notes

If your Go build cache is sandboxed, set a local cache directory:

```bash
GOCACHE=./.gocache go build ./...
```

## Makefile Targets

- `make build` builds the local binary.
- `make test` runs Go tests.
- `make lint` runs `go vet`.
- `make install` installs to your `$GOBIN`.
- `make dist` builds cross-platform binaries into `dist/`.
- `make clean` removes `dist/`.
