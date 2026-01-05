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

## Quick Start (Run in this repo)

You can run this example immediately using the `--protocol` override:

1. **Requirements**:
   ```bash
   specfirst --protocol examples/todo-cli/protocol.yaml reqs
   ```
   
2. **Design**:
   ```bash
   specfirst --protocol examples/todo-cli/protocol.yaml design
   ```

## Setup (For a new project)

To use this protocol in your own project:

1. Create a new directory and initialize:
   ```bash
   mkdir my-todo-app && cd my-todo-app
   specfirst init
   ```

2. Copy the protocol and templates:
   ```bash
   cp /path/to/specfirst/examples/todo-cli/protocol.yaml .specfirst/protocols/
   cp -r /path/to/specfirst/examples/todo-cli/templates/* .specfirst/templates/
   ```

3. Update config (optional) or use the flag:
   ```bash
   # Option A: Edit .specfirst/config.yaml to set protocol: todo-cli-protocol
   # Option B: Use flag
   specfirst --protocol todo-cli-protocol reqs
   ```
