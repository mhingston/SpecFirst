# Refactoring Goals - outputRelPath

## Primary Goals
- **Architecture Alignment:** Relocate path normalization logic from the CLI layer (`cmd/helpers.go`) to the core domain layer (`internal/workspace/workspace.go`), ensuring business logic does not depend on CLI helpers.
- **Eliminate Code Duplication:** Provide a shared `ProjectRelPath` utility that resolves the `TODO` in `internal/engine/complete.go`, removing the need for local path manipulation hacks.
- **Dependency Decoupling:** Replace the external shell dependency on `git rev-parse` with a native Go implementation (`FindProjectRoot`) that traverses parent directories for `.specfirst` or `.git` markers.
- **Centralized Testing:** Consolidate path-related tests into `internal/workspace/workspace_test.go`, achieving 100% branch coverage for the new implementation.

## Non-Goals
- Changing the core path normalization behavior (inputs and outputs must remain consistent).
- Adding support for complex monorepo structures or multi-root workspaces.
- Refactoring other unrelated helpers in the `cmd` package.

## Success Criteria
- `cmd/helpers.go` no longer contains the `outputRelPath` function.
- `internal/engine/complete.go` and `cmd` functions successfully import and use `workspace.ProjectRelPath`.
- All migrated tests pass, and new tests verify `FindProjectRoot` correctly identifies the root in various directory depths.
- The `git` command is no longer invoked for standard path normalization tasks.

## Constraints
- **Compatibility:** Must maintain current behavior for all existing CLI commands.
- **Zero Dependencies:** Must not introduce new third-party libraries; use standard library `os` and `path/filepath`.
- **Error Handling:** Must return clear, actionable errors when a path is outside the project root or the root cannot be found.

## Benefits
- **Improved Maintainability:** Logical separation of concerns makes the codebase easier to navigate and modify.
- **Portability:** Removing the hard `git` dependency allows the core engine to function in environments where `git` might not be initialized or available.
- **Testability:** Core workspace logic can now be unit-tested in isolation from the CLI and the filesystem shell.

## Assumptions
- The presence of a `.specfirst` directory or a `.git` directory is a definitive indicator of the project root.
- All existing consumers of `outputRelPath` are within the same repository/module.
