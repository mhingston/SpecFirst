# SpecFirst User Guide

SpecFirst is a CLI for specification-driven development. It helps you maintain a disciplined workflow by using declarative protocols to guide the creation of requirements, designs, and code.

## Core Concepts

- **Protocol**: A YAML file defining the stages of your workflow (e.g., Requirements -> Design -> Implementation).
- **Stage**: A single unit of work defined in the protocol (Note: stage IDs must be **lowercase**).
- **Template**: A markdown or text file using Go template syntax to render a prompt for a stage.
- **Artifact**: The output of a completed stage (e.g., a `.md` file for design, or `.go` files for implementation).
- **State**: Tracked in `.specfirst/state.json`, recording completed stages, approvals, and prompt hashes.

## Getting Started

### 1. Initialize a Project
Run `specfirst init` in your project root. This creates the `.specfirst` directory with a default protocol, templates, and configuration.

### 2. Check Status
Use `specfirst status` to see your current progress in the workflow. It shows which stages are completed and what's next.

### 3. Generate a Prompt
To work on a stage (e.g., `requirements`), run:
```bash
specfirst requirements
```
This renders the template for that stage to `stdout`, embedding any needed context from previous stages. You can pipe this to an AI CLI:
```bash
specfirst requirements | claude -p
```

### 4. Complete a Stage
Once you have the output from the LLM, record it:
```bash
specfirst complete requirements ./requirements.md
```
This moves the file into the artifact store and updates the project state.

## Advanced Workflow

### Task Decomposition
Protocols can include a `decompose` stage that breaks down a design into a list of specific tasks.
1. Run `specfirst decompose` and save the LLM output to `tasks.yaml`.
2. Complete the stage: `specfirst complete decompose ./tasks.yaml`.
3. List tasks: `specfirst task`.
4. Generate a prompt for a specific task: `specfirst task T1`.

### Approvals
Stages can required approvals from specific roles (e.g., `architect`, `product`).
```bash
specfirst approve requirements --role architect --by "Jane Doe"
```

### Validation
Run `specfirst lint` or `specfirst check` to find issues like protocol drift, missing artifacts, or vague prompts.

### Archiving
When a spec version is finalized, archive it:
```bash
specfirst complete-spec --archive --version 1.0
```
This creates a snapshot of the entire workspace.
