# Current State Analysis - outputRelPath Refactoring

## Code Location
- **Source:** `cmd/helpers.go` (Function: `outputRelPath`)
- **Target:** `internal/workspace/workspace.go` (Proposed: `ProjectRelPath`)
- **Consumer:** `internal/engine/complete.go` (Currently contains TODO to copy logic)
- **Tests:** `cmd/helpers_test.go` (Test: `TestOutputRelPath`)

## Purpose
The `outputRelPath` function calculates a file's path relative to the project root. This is critical for normalizing artifact paths stored in the `.specfirst/artifacts` directory, ensuring they are portable and consistent regardless of where the CLI is executed.

## Structure
- `cmd/helpers.go`: Contains `outputRelPath` which depends on `gitRoot()` (wraps `git rev-parse`).
- `internal/engine/complete.go`: Contains logic for completing a stage. It currently needs to store artifacts but lacks a robust way to determine the "project relative" path, relying on a local `filepath.Rel` or having a TODO to fix it.

## Problems
1.  **Architecture Violation:** Logic for workspace path normalization resides in the CLI layer (`cmd`), but is needed by the core business logic (`internal/engine`).
2.  **Code Duplication/Debt:** `internal/engine/complete.go` has a TODO explicitly stating: `// For this refactor, I will COPY outputRelPath logic here as a private helper or fix workspace to have it.`
3.  **Hard Dependency on Git:** The current implementation relies on executing `git` commands (`git rev-parse`), which makes it slower and less portable than a pure Go directory traversal (scanning for `.specfirst` or `.git`).
4.  **Testing Split:** Tests for core path logic are in `cmd` tests.

## Current Behavior
- The function accepts an absolute or relative path.
- It determines the project root (currently via `git`).
- It returns the path relative to the root.
- It returns an error if the path is outside the project root.
- This behavior MUST be preserved.

## Refactoring Plan Strategy
1.  Implement `FindProjectRoot` in `internal/workspace` to locate the project root (looking for `.specfirst` directory or `.git`).
2.  Move `outputRelPath` logic to `internal/workspace`, renaming it to `ProjectRelPath`.
3.  Refactor `ProjectRelPath` to use `FindProjectRoot` instead of shelling out to `git`.
4.  Move relevant tests from `cmd/helpers_test.go` to `internal/workspace/workspace_test.go`.
5.  Update `cmd` and `internal/engine` to use the new `workspace.ProjectRelPath`.

## Test Coverage
- `TestOutputRelPath` in `cmd/helpers_test.go` covers:
    - Relative paths inside root.
    - Absolute paths inside root.
    - Paths outside root (error).
    - Paths in subdirectories.
