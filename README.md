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

- [User Guide](docs/GUIDE.md): Detailed "how-to" and workflow examples.
- [Common Workflows](docs/WORKFLOWS.md): Practical patterns for real-world scenarios.
- [Protocol Reference](docs/PROTOCOLS.md): YAML schema and stage definitions.
- [Template Reference](docs/TEMPLATES.md): Guide to authoring stage templates.

## Canonical Examples

- [Todo CLI](examples/todo-cli/README.md): A full walkthrough of a simple CLI application.

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

## AI CLI Integration

SpecFirst outputs prompts to stdout, making it composable with AI CLIs in the Unix tradition.

**Tools that support direct piping:**

```bash
# Pipe to Claude Code (requires -p flag for headless mode)
specfirst requirements | claude -p

# Pipe to Codex (use - to read from stdin)
specfirst design | codex -

# Pipe to Gemini CLI
specfirst implementation | gemini
```

**Tools that require the prompt as an argument** (use command substitution or `--out`):

```bash
# GitHub Copilot requires the prompt as a string argument
copilot -p "$(specfirst requirements)" --allow-all-tools

# Or write to a file first
specfirst requirements --out prompt.txt
copilot -p "$(cat prompt.txt)" --allow-all-tools
```

**For any tool that reads from files**, use `--out`:

```bash
specfirst requirements --out prompt.txt
some-ai-tool --input prompt.txt
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
- `specfirst approve <stage-id> --role <role>` records approvals.

### Cognitive Scaffold Commands

These commands generate **prompts only** — no state, no enforcement, no AI calls. They shape thinking, not execution.

- `specfirst diff <old-spec> <new-spec>` generates a change-analysis prompt comparing two specification files.
- `specfirst assumptions <spec-file>` generates a prompt to surface hidden assumptions.
- `specfirst review <spec-file> --persona <p>` generates a role-based review prompt. Personas: `security`, `performance`, `maintainer`, `accessibility`, `user`.
- `specfirst failure-modes <spec-file>` generates a failure-first interrogation prompt.
- `specfirst test-intent <spec-file>` generates a test **intent** prompt (not test code).
- `specfirst trace <spec-file>` generates a spec-to-code mapping prompt.
- `specfirst distill <spec-file> --audience <a>` generates an audience-specific summary prompt. Audiences: `exec`, `implementer`, `ai`, `qa`.
- `specfirst calibrate <artifact> [--mode <m>]` generates an epistemic annotation prompt for judgment calibration. Modes: `default`, `confidence`, `uncertainty`, `unknowns`.

## Completion Options

- `--prompt-file <path>` hash an explicit prompt file when completing a stage.
- `--force` overwrite an existing stage completion (non-destructive; only removes old artifacts after new ones are successfully stored).

## Stage Execution Options

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

## Approval Options

- `--role <role>` (required) the role for the approval.
- `--by <name>` who approved (defaults to `$USER`).
- `--notes <text>` optional notes for the approval.

## Workspace Layout

```
.specfirst/
  artifacts/     # Stage outputs stored by stage ID
  generated/     # Generated files (e.g., compiled prompts)
  protocols/     # Protocol YAML definitions
  templates/     # Prompt templates
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
