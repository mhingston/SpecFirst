# Example: Todo CLI

This example demonstrates a full SpecFirst workflow, including requirements, design, and task-scoped implementation for a simple Todo CLI application.

## Setup

1. Create a new directory for your project and enter it:
   ```bash
   mkdir my-todo-app && cd my-todo-app
   ```

2. Copy the example files:
   ```bash
   cp -r /path/to/specfirst/examples/todo-cli/templates .specfirst/
   cp /path/to/specfirst/examples/todo-cli/protocol.yaml .specfirst/protocols/
   ```

3. Initialize SpecFirst with this protocol:
   ```bash
   specfirst init --protocol todo-cli-protocol
   ```

## Workflow Walkthrough

### 1. Requirements
Generate the requirements prompt:
```bash
specfirst reqs | claude -p
```
Save the output as `requirements.md`, then complete the stage:
```bash
specfirst complete reqs ./requirements.md
```

### 2. Design
Generate the design prompt (it will include `requirements.md` automatically):
```bash
specfirst design | claude -p
```
Save the output as `design.md`, then complete the stage:
```bash
specfirst complete design ./design.md
```

### 3. Task Decomposition
Break the design into specific tasks:
```bash
specfirst breakdown | claude -p
```
Save the YAML output as `tasks.yaml`, then complete:
```bash
specfirst complete breakdown ./tasks.yaml
```

### 4. Implementation
List the generated tasks:
```bash
specfirst task
```
Generate a prompt for the first task:
```bash
specfirst task T1 | claude -p
```

### 5. Finalize
Verify your progress at any time with `specfirst status`.
When finished, validate the whole spec:
```bash
specfirst check
```
