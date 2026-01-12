# Example: Todo CLI

This example demonstrates a full SpecFirst workflow, including requirements clarification, design, and task-scoped implementation for a simple Todo CLI application.

## Setup

1. Create a new directory for your project and initialize it with Git:
   ```bash
   mkdir my-todo-app && cd my-todo-app
   git init
   ```

2. Initialize SpecFirst with the `todo-cli` starter:
   ```bash
   specfirst init --starter todo-cli
   ```

## Quick Start (Run in this repo)

You can run this example immediately using the `--protocol` override:

1. **Clarify Requirements**:
   ```bash
   opencode run "$(specfirst --protocol starters/todo-cli/protocol.yaml clarify)"
   ```
   
2. **Design**:
   ```bash
   opencode run "$(specfirst --protocol starters/todo-cli/protocol.yaml design)"
   ```

## Workflow

After initializing the project, follow these steps:

### 1. Clarify Requirements

Generate application requirements:
```bash
opencode run "$(specfirst clarify)" > requirements.md
specfirst complete clarify ./requirements.md
```

### 2. Design

Create technical design based on requirements:
```bash
opencode run "$(specfirst design)" > design.md
specfirst complete design ./design.md
```

### 3. Break Down Tasks

Decompose the work into manageable tickets:
```bash
opencode run "$(specfirst breakdown)" > tasks.yaml
specfirst complete breakdown ./tasks.yaml
```

### 4. Implementation

Implement tasks one by one:
```bash
# View tasks
specfirst task

# Generate prompt for a task (e.g., T1)
opencode run "$(specfirst codes T1)"

# Mark as complete
specfirst complete codes ./main.go
```
