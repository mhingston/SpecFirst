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

### 5. Overriding Protocols
You can override the active protocol for any command using the `--protocol` flag. This is useful for testing different workflows or using protocols stored outside `.specfirst/protocols/`:
```bash
specfirst --protocol path/to/custom-protocol.yaml status
```

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

## Team Collaboration

**Use case**: Multiple developers working on different tasks from decomposition.

```bash
# Team Lead:
specfirst init
specfirst requirements | claude -p > requirements.md
specfirst complete requirements ./requirements.md
specfirst decompose | claude -p > tasks.yaml
specfirst complete decompose ./tasks.yaml

# Share task list with team
specfirst task  # Shows all available tasks

# Developer A:
specfirst task T1 | claude -p  # Work on task T1

# Developer B:
specfirst task T2 | claude -p  # Work on task T2 (parallel)

# Team Lead (after completion):
specfirst check  # Validate all outputs
specfirst complete-spec --archive --version 0.1.0
```

## CI/CD Integration

**Use case**: Automated quality checks in build pipeline.

```bash
# In your CI pipeline (e.g., .github/workflows/spec-check.yml):

- name: Check spec quality
  run: |
    specfirst check --fail-on-warnings
    
- name: Validate all stages complete
  run: |
    specfirst complete-spec --warn-only || exit 1
    
- name: Archive on release
  if: startsWith(github.ref, 'refs/tags/')
  run: |
    specfirst complete-spec --archive --version ${{ github.ref_name }}
```

## Best Practices

1. **Always run `specfirst check` before archiving** - Catches missing outputs or approvals
2. **Use `--prompt-file` when completing stages** - Enables prompt hash verification
3. **Archive at milestones** - Creates rollback points for major versions
4. **Leverage cognitive commands iteratively** - Run assumptions/review multiple times as specs evolve
5. **Prefer `--out` for complex prompts** - Easier to review before sending to AI
