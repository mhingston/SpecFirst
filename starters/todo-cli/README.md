# Example: Todo CLI

This example demonstrates a full SpecFirst workflow, including requirements, design, and task-scoped implementation for a simple Todo CLI application.

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

1. **Requirements**:
   ```bash
   opencode run "$(specfirst --protocol starters/todo-cli/protocol.yaml reqs)"
   ```
   
2. **Design**:
   ```bash
   opencode run "$(specfirst --protocol starters/todo-cli/protocol.yaml design)"
   ```
