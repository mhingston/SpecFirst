# CLI Reference

## Commands

- `specfirst init` initializes `.specfirst/` with defaults.
- `specfirst init --starter <name>` initializes with a specific starter kit.
- `specfirst starter list` lists available starter kits.
- `specfirst starter apply <name>` applies a starter kit to the current workspace.
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
- `specfirst attest <stage-id> --role <role> --status <status>` records attestations with rationale and conditions.
- `specfirst track create|list|switch|diff|merge` manages parallel futures (tracks).

### Cognitive Scaffold Commands

These commands generate **prompts only** — no state, no enforcement, no AI calls. They shape thinking, not execution.

- `specfirst diff <old-spec> <new-spec>` generates a change-analysis prompt comparing two specification files.
- `specfirst assumptions <spec-file>` generates a prompt to surface hidden assumptions.
- `specfirst review <spec-file> --persona <p>` generates a role-based review prompt. Personas: `security`, `performance`, `maintainer`, `accessibility`, `user`.
- `specfirst failure-modes <spec-file>` generates a failure-first interrogation prompt.
- `specfirst test-intent <spec-file>` generates a test **intent** prompt (not test code).
- `specfirst trace <spec-file>` generates a spec-to-code mapping prompt.
- `specfirst distill <spec-file> --audience <a>` generates an audience-specific summary prompt. Audiences: `exec`, `implementer`, `ai`, `qa`.
- `specfirst calibrate <artifact>` generates a comprehensive epistemic map for judgment calibration.

## Completion Options

- `--prompt-file <path>` hash an explicit prompt file when completing a stage.
- `--force` overwrite an existing stage completion (non-destructive; only removes old artifacts after new ones are successfully stored).

## Stage Execution Options

- `--protocol <path|name>` override active protocol (path to file or name in `.specfirst/protocols`).
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

## Track Options
 
 - `track create <name> --notes <text>` create a new track.
 - `track switch <name> --force` restore a track to the current workspace (overwrites existing data).
 - `track merge <source>` generate a merge plan prompt.
 
 ## Attestation Options
 
 - `--role <role>` (required) the role for the attestation.
 - `--status <status>` (required) status: `approved`, `approved_with_conditions`, `needs_changes`, `rejected`.
 - `--rationale <text>` rationale for the decision.
 - `--condition <text>` condition for conditional approval (repeatable).
 - `--by <name>` who attested (defaults to `$USER`).
