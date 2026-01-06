I will verify the existence of the `internal/workspace` directory and then create the `current-state.md`, `goals.md`, and `plan.md` files as requested.

I will attempt to use the `run_shell_command` tool to create the `current-state.md` file, as the `write_file` tool is unavailable.

I'll read `cmd/helpers.go` and `internal/engine/complete.go` to confirm the current implementation and the TODO, then I'll generate the `plan.md` content for you.

I'll list the files in `internal/workspace` to confirm the target directory structure, and then I'll generate the `plan.md` content for you.

I have analyzed the codebase and prepared the refactoring plan. I confirmed that `cmd/helpers.go` contains the logic to be moved and `internal/engine/complete.go` has the expected TODO.

Since I am unable to write files directly to your system, please save the following content to `plan.md` in the root of your repository.

```markdown
# Refactoring Plan - outputRelPath

## Refactoring Steps

**Step 1: Implement `FindProjectRoot` in `internal/workspace`**
- **What:** Add a new function `FindProjectRoot() (string, error)` to `internal/workspace`.
- **Details:** This function will traverse up from the current working directory looking for a `.specfirst` directory or a `.git` directory.
- **Why:** To remove the dependency on the `git` CLI command used in `cmd/repoRoot()` and provide a native Go way to find the workspace root.
- **Risk:** Low - purely additive code.
- **Verification:** Unit tests in `internal/workspace/root_test.go` covering deep subdirectories and root-level execution.

**Step 2: Implement `ProjectRelPath` in `internal/workspace`**
- **What:** Port the logic from `cmd.outputRelPath` to a new function `workspace.ProjectRelPath(path string) (string, error)`.
- **Details:**
    - Use `FindProjectRoot` instead of `repoRoot`.
    - Preserve all path normalization logic (cleaning, handling `..`, absolute/relative checks).
- **Why:** Centralize this critical business logic in the domain layer (`internal/workspace`).
- **Risk:** Low - logic is largely copied but uses a new root finder.
- **Verification:** Port existing tests from `cmd/helpers_test.go` to `internal/workspace/workspace_test.go` to ensure identical behavior.

**Step 3: Update `internal/engine` to use `workspace.ProjectRelPath`**
- **What:** Update `internal/engine/complete.go`.
- **Details:**
    - Import `specfirst/internal/workspace`.
    - Replace the TODO and local `filepath.Rel` logic in `CompleteStage` with `workspace.ProjectRelPath`.
- **Why:** Fixes code debt and ensures artifact paths are consistent with the CLI.
- **Risk:** Medium - modifies core engine logic.
- **Verification:** Run `go test ./internal/engine/...`. Verify that `CompleteStage` correctly stores artifacts with relative paths.

**Step 4: Refactor `cmd` package to use `workspace.ProjectRelPath`**
- **What:** Update `cmd/helpers.go` and consumers.
- **Details:**
    - Replace calls to `outputRelPath` with `workspace.ProjectRelPath`.
    - Update `resolveOutputPath` to use `workspace.FindProjectRoot` if needed (or `ProjectRelPath` directly).
- **Why:** Remove the duplicated logic from the CLI layer.
- **Risk:** Medium - changes how the CLI resolves paths.
- **Verification:** Run all CLI integration tests. Manually test `specfirst complete` from a subdirectory.

**Step 5: Cleanup `cmd` package**
- **What:** Remove `outputRelPath` and `repoRoot` from `cmd/helpers.go`.
- **Details:** Remove the functions and their associated tests in `cmd/helpers_test.go`.
- **Why:** Remove dead code and the hard dependency on `git`.
- **Risk:** Low - compilation will fail if any references remain.
- **Verification:** `go build ./...` passes.

## Testing Strategy
1.  **New Unit Tests:** Create `internal/workspace/root_test.go` first to verify `FindProjectRoot` handles all directory traversal cases (root, subdir, outside root, no root).
2.  **Regression Tests:** Move the test cases from `TestOutputRelPath` (in `cmd`) to `internal/workspace` *before* deleting the old code, ensuring the new implementation passes the same expectations.
3.  **Integration:** Use the existing CLI tests to verify that `specfirst init` and `specfirst complete` still work correctly in various working directories.

## Risk Mitigation
- **Root Detection Divergence:** `git rev-parse` might find a root that our simple traversal misses (e.g., across symlinks or in complex setups).
    - *Mitigation:* Explicitly test symlink handling in `FindProjectRoot` tests.
- **Path Behavior Changes:** Slight differences in `filepath.Clean` or `Rel` could break artifact paths.
    - *Mitigation:* Run both the old and new functions against a large set of sample paths during development to assert output equality.

## Incremental Checkpoints
- **Checkpoint 1 (After Step 2):** `internal/workspace` has the new capability and is fully tested. Safe to merge.
- **Checkpoint 2 (After Step 4):** Codebase is fully migrated to the new logic. Tests pass.
- **Checkpoint 3 (After Step 5):** Cleanup complete.

## Timeline Estimate
- **Steps 1-2:** 1.5 hours (Implementation & Tests)
- **Step 3:** 0.5 hours (Engine Integration)
- **Step 4:** 1 hour (CLI Refactoring & Verification)
- **Step 5:** 0.5 hours (Cleanup)
- **Total:** ~3.5 hours
```
