# SpecFirst

SpecFirst is a Go CLI for specification-driven workflows that focuses on **prompt quality, clarity, and human judgment**.

**Why SpecFirst?**
Ambiguity is the enemy of automated coding. SpecFirst exists to bridge the gap between vague human intent and precise AI execution. It provides a structured, file-centric protocol for communicating with LLMs, ensuring that requirements are clear, context is substantial, and humans remain in control of the creative process.

## What SpecFirst Is

- **Protocol-Driven**: Workflows are defined in `protocol.yaml` files, making them versionable and reproducible.
- **Staged Execution**: Complex tasks are broken down into logical stages (e.g., Scope -> Spec -> Implementation -> Review).
- **Template-Based**: Prompts are rendered from declarative Markdown templates, ensuring consistency.
- **Artifact-Centric**: Every stage produces tangible artifacts (files) that serve as the input for the next stage.
- **Human-in-the-Loop**: Designed to augment human reasoning, not replace it. You review and approve artifacts at each step.

## What SpecFirst Is Not

- **Not an Agent**: It does not execute code, run terminals, or make decisions for you. It prepares the context for you (or an agent) to act.
- **Not a Sandbox**: It runs in your shell. It does not provide an isolated execution environment.
- **Not a Test Framework**: While it can help write tests, it replaces neither unit tests nor integration tests.
- **Not Magic**: It requires you to articulate your intent clearly in the templates and inputs.

## Quick Start
Get up and running with a single copy/paste flow:

```bash
# 1. Install & Initialize
specfirst init

# 2. Run the "requirements" example
# Generates a prompt, pipes it to your clipboard (macos pbcopy)
specfirst requirements | pbcopy

# 3. Paste into your LLM to get requirements.md content
# ... (User Action: Paste result into requirements.md) ...

# 4. Complete the stage
specfirst complete requirements ./requirements.md

# 5. See your result
ls -l .specfirst/artifacts/
```

## Repository Layout
A standard SpecFirst project looks like this:

- `.specfirst/`
    - `artifacts/`: Stored outputs from completed stages (hashed and versioned).
    - `protocols/`: Your workflow definitions (e.g., `feature.yaml`, `bugfix.yaml`).
- `templates/`: Markdown templates for your prompts (e.g., `requirements.md`, `plan.md`).
- `protocol.yaml`: The default workflow definition for the project.
- `inputs/`: (Optional) Context files or static inputs for your workflows.

## Security

See [SECURITY.md](SECURITY.md) for guidance on data handling, credential safety, and using third-party models.

## Documentation

- [User Guide](docs/GUIDE.md): Detailed "how-to" and workflow examples.
- [CLI Reference](docs/REFERENCE.md): Commands, flags, and options.
- [Protocol Reference](docs/PROTOCOLS.md): YAML schema and stage definitions.
- [Template Reference](docs/TEMPLATES.md): Guide to authoring stage templates.
- [Template Context](docs/template-context.md): Variables and runner guarantees.
- [AI Integration](docs/AI_INTEGRATION.md): Patterns for using SpecFirst with AI tools.
- [Philosophy](docs/PHILOSOPHY.md): The "why" behind SpecFirst.

## Canonical Examples

Complete, runnable examples with protocols and templates:

- [Todo CLI](starters/todo-cli/README.md): Full workflow walkthrough for a simple CLI application
- [Bug Fix](starters/bug-fix/README.md): Minimal 2-stage workflow for systematic bug fixes
- [API Feature](starters/api-feature/README.md): Complete feature development with approvals and task decomposition
- [Spec Review](starters/spec-review/README.md): Using cognitive scaffold commands to improve specification quality
- [Refactoring](starters/refactoring/README.md): Structured code improvement with risk mitigation and verification
- [Code Review](starters/code-review/README.md): Staged, gated review workflow designed to converge
- [Database Migration](starters/database-migration/README.md): Safe schema changes with approval gates and rollback planning
- [Product Design](starters/product-design/README.md): Design flow from product vision to implementation handoff

## Install

### From Source

```bash
make build
# Or
go build .
```

### From GitHub Releases

Download the pre-built binary for your platform from [Releases](https://github.com/mhingston/SpecFirst/releases).

## Shell Completion

Generate completion scripts for your shell:

```bash
specfirst completion zsh
specfirst completion bash
specfirst completion fish
specfirst completion powershell
```

Common installs:

```bash
# zsh (one-time install)
specfirst completion zsh > "${fpath[1]}/_specfirst"

# bash (Homebrew path)
specfirst completion bash > /usr/local/etc/bash_completion.d/specfirst

# fish
specfirst completion fish > ~/.config/fish/completions/specfirst.fish
```
