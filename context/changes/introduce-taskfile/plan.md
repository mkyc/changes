# Introduce Taskfile Implementation Plan

## Overview

Add a `Taskfile.yml` (for [go-task/task](https://taskfile.dev)) that wraps the project's build, test, and quality-check workflows into four commands: `task build`, `task test`, `task doctor`, and `task cleanup`.

## Current State Analysis

The project is a small single-module Go CLI (`github.com/mkyc/changes`) with no existing build automation — no `Makefile`, no `Taskfile.yml`, no CI workflow. Developers currently run raw `go build` / `go test` / `go vet` commands by hand. `task` (v3.48.0) and `golangci-lint` (v2.12.2) are both already installed in the local dev environment, but there is no `.golangci.yml` config, so `golangci-lint run` will use its built-in default linter set. `.gitignore` currently only excludes `.idea/`, `.vscode/`, and `/changes` (a checked-in-looking but gitignored binary at repo root from ad-hoc `go build` runs).

## Desired End State

A `Taskfile.yml` exists at the repo root with four tasks:

- `task build` — compiles the CLI to `bin/changes`
- `task test` — runs the full test suite verbosely
- `task doctor` — runs formatting, vet, module-tidiness, and lint checks
- `task cleanup` — removes build output and clears the Go build/test cache

Verification: running each of the four commands from a clean checkout succeeds (or, for `doctor`, correctly surfaces any real issues) and `bin/` is not tracked by git.

### Key Discoveries:

- No prior Makefile/Taskfile/CI convention to follow — this is a greenfield addition (confirmed via repo-wide file listing).
- `.gitignore` (repo root) excludes `/changes` but not `bin/` — needs updating since build output is moving to `bin/changes`.
- `golangci-lint` v2.x is installed; v2's default command is `golangci-lint run` (no config file required, uses default enabled linters).
- Module path is `github.com/mkyc/changes`; entrypoint is `main.go` at repo root (`internal/app`, `internal/cli`, `internal/config` house the logic).

## What We're NOT Doing

- Not adding a `.golangci.yml` config — using golangci-lint's defaults is sufficient for this task; a custom lint config is a separate concern.
- Not setting up CI (GitHub Actions etc.) — this task is scoped to the local Taskfile only.
- Not adding coverage reporting/thresholds to `task test`.
- Not removing or replacing the existing ad-hoc `./changes` binary at repo root (it stays gitignored as-is; new builds go to `bin/`).
- Not adding a `task install` or release/packaging task — out of scope per the ticket's four acceptance criteria.

## Implementation Approach

Single flat `Taskfile.yml` using the standard Taskfile schema (`version: '3'`), with one task per acceptance criterion. Each task shells out to standard `go` toolchain commands (`go build`, `go test`, `gofmt`, `go vet`, `go mod tidy`, `golangci-lint run`) rather than introducing new dependencies. `.gitignore` gets a `bin/` entry so build output never gets committed.

## Phase 1: Add Taskfile and supporting gitignore entry

### Overview

Create `Taskfile.yml` with the four required tasks and update `.gitignore` so the new build output directory is excluded from version control.

### Changes Required:

#### 1. Taskfile

**File**: `Taskfile.yml` (new, repo root)

**Intent**: Define the four tasks the ticket requires, each wrapping standard Go tooling.

**Contract**:
- `version: '3'` schema.
- `build`: runs `go build -o bin/changes .`
- `test`: runs `go test ./... -v`
- `doctor`: runs, in sequence, a gofmt check that fails if any file is unformatted (`gofmt -l .` with non-empty output causing failure — e.g. via a shell one-liner or `test -z "$(gofmt -l .)"`), `go vet ./...`, a `go mod tidy` drift check (run `go mod tidy` then `git diff --exit-code go.mod go.sum`, or equivalent diff-based check), and `golangci-lint run`.
- `cleanup`: removes `bin/` (e.g. `rm -rf bin`) and runs `go clean -testcache -cache`.

#### 2. Gitignore

**File**: `.gitignore`

**Intent**: Exclude the new `bin/` build output directory from version control, consistent with the existing `/changes` exclusion.

**Contract**: Add a `bin/` entry (new line) to the existing three-line ignore list.

### Success Criteria:

#### Automated Verification:

- `task build` succeeds and produces an executable at `bin/changes`: `task build && test -x bin/changes`
- `task test` succeeds: `task test`
- `task doctor` succeeds on a clean tree (no formatting/vet/tidy/lint issues introduced by this change): `task doctor`
- `task cleanup` removes `bin/`: `task cleanup && test ! -d bin`
- `git status --short` shows `bin/` is not tracked after a build: `task build && git status --short -- bin/` produces no output

#### Manual Verification:

- Run `bin/changes --help` after `task build` and confirm it behaves like the existing binary (same subcommands as `internal/cli`).
- Introduce a deliberately unformatted or lint-failing line temporarily and confirm `task doctor` reports it, then revert.

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful before proceeding to the next phase.

---

## Testing Strategy

### Unit Tests:

N/A — this change adds tooling, not application code. No new unit tests required.

### Integration Tests:

- Run each of the four `task` commands end-to-end from a clean working tree as described in Success Criteria above.

### Manual Testing Steps:

1. From a clean checkout, run `task build`; confirm `bin/changes` is produced and executable.
2. Run `bin/changes --help` (or equivalent) and confirm output matches running `go run . --help`.
3. Run `task test`; confirm it matches `go test ./... -v` output.
4. Run `task doctor`; confirm it passes cleanly on the current tree.
5. Run `task cleanup`; confirm `bin/` is removed.

## Performance Considerations

None — these are developer-facing convenience commands with no runtime/production impact.

## Migration Notes

None — no existing build tooling to migrate from or deprecate.

## References

- Ticket: `tasks/0010-create-taskfile.md`

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Add Taskfile and supporting gitignore entry

#### Automated

- [x] 1.1 `task build` succeeds and produces an executable at `bin/changes` — b68a092
- [x] 1.2 `task test` succeeds — b68a092
- [x] 1.3 `task doctor` succeeds on a clean tree — b68a092
- [x] 1.4 `task cleanup` removes `bin/` — b68a092
- [x] 1.5 `bin/` is not tracked by git after a build — b68a092

#### Manual

- [x] 1.6 `bin/changes --help` behaves like the existing binary — b68a092
- [x] 1.7 `task doctor` correctly flags a deliberately introduced issue, then passes again after revert — b68a092
