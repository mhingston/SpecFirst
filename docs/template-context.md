# Template Context & Runner Guarantees

This document describes the variables available in SpecFirst templates and the guarantees provided by the runtime environment.

## Template Context

SpecFirst templates are rendered using Go's `text/template` package. The following variables are available in every stage:

### Base Variables

| Variable | Type | Description |
| :--- | :--- | :--- |
| `.StageName` | `string` | The ID of the current stage (e.g., `requirements`, `plan`). |
| `.ProjectName` | `string` | The name of the project root directory. |
| `.Inputs` | `[]Input` | List of input files available to this stage. |
| `.Outputs` | `[]string` | List of expected output filenames for this stage. |

### Input Object

Each item in `.Inputs` has:

- `.Name`: The filename (e.g., `main.go`).
- `.Content`: The file content (string).
- `.Path`: The relative path to the file.

### Template Functions

Standard Go template functions are available, plus:

- `readFile <path>`: Reads a file from the workspace (relative to root).
- `upper`: Converts string to uppercase.
- `lower`: Converts string to lowercase.

---

## Runner Guarantees

The SpecFirst runner provides the following guarantees to ensure deterministic and reliable execution:

### 1. Artifact Determinism
- **Hashing**: Outputs are hashed. If the inputs (prompt hash) haven't changed, the stage is skipped (unless forced).
- **Paths**: Artifacts are always written to `.specfirst/artifacts/<StageID>/<Hash>/`.

### 2. Failure Modes
- **Missing Output**: If a stage completes but the expected output file is not found, the runner exits with an error status.
- **Template Error**: If a template fails to render (e.g., missing variable), the process aborts immediately.
- **Missing Input**: If a file specified in `inputs:` is missing, the runner halts before generating the prompt.

### 3. File Operations
- **Atomic Writes**: Artifact writes are effectively atomic (file is finalized after write).
- **Read-Only Context**: Templates cannot modify input files; they can only generate text to stdout (which is then captured).
