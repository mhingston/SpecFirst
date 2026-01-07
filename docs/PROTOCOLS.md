# Protocol Reference

A Protocol defines the stages, dependencies, and outputs of a workflow. It is the "source of truth" for the discipline SpecFirst amplifies.

## Root Fields

| Field | Type | Description |
| --- | --- | --- |
| `name` | string | Unique name for the protocol. |
| `version` | string | Version of the protocol definition. |
| `uses` | []string | Optional list of protocols to import stages from. |
| `stages` | []Stage | List of stages in the workflow. |
| `approvals` | []Approval | Required approvals for specific stages. |

## Stage Fields

| Field | Type | Description |
| --- | --- | --- |
| `id` | string | Unique identifier for the stage (must be **lowercase**). |
| `name` | string | Human-readable name. |
| `type` | string | `spec`, `decompose`, or `task_prompt`. |
| `intent` | string | Semantic intent (e.g., `exploration`, `execution`). |
| `template` | string | Template filename in `.specfirst/templates/`. |
| `depends_on` | []string | IDs of stages that must be completed first. |
| `inputs` | []string | Filenames of artifacts from previous stages. Must match an entry in the `outputs` list of one of the stages in `depends_on`. |
| `outputs` | []string | Expected filenames to be produced. |

### Decomposition Fields
Used when `type: decompose` or `type: task_prompt`.

| Field | Type | Description |
| --- | --- | --- |
| `source` | string | For `task_prompt`, the ID of the `decompose` stage providing tasks. |
| `prompt` | PromptConfig | Configuration for task generation (granularity, etc.). |

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

## Example Protocol

```yaml
name: "fast-track"
version: "1.0"

stages:
  - id: outline
    name: Feature Outline
    type: spec
    template: outline.md
    outputs: [outline.md]

  - id: implement
    name: Implementation
    type: spec
    template: implement.md
    depends_on: [outline]
    inputs: [outline.md]

approvals:
  - stage: outline
    role: lead
```

## Protocol Evolution & Versioning

As projects grow, your protocols will need to evolve. Proper versioning ensures that active work isn't disrupted by changes to the underlying workflow.

### Guidance on Versioning

1.  **Semantic Versioning:** Treat protocols like APIs.
    *   **Patch (1.0.1):** Typofixes in templates or human-readable stage names.
    *   **Minor (1.1.0):** Adding a new optional stage or approval.
    *   **Major (2.0.0):** Removing stages, changing `depends_on` relationships, or modifying `inputs`/`outputs` in a way that breaks existing `state.json` files.

2.  **Immutability:** Once a protocol is used in production, avoid making breaking changes to that specific version. Instead, create a new version and migrate.

### Migration Strategies

*   **Long-Lived Projects:** For projects that take weeks or months, it is often safer to finish the project using the protocol version it started with.
*   **Protocol Overrides:** Use the `--protocol` flag to explicitly point to an older version if you need to maintain compatibility with an archived state.
*   **Manual State Correction:** In extreme cases, you may need to manually update the `state.json` or move artifacts to match the new protocol's expectations before running `spec init` or `spec complete`.

