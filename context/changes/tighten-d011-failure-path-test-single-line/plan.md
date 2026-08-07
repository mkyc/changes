# Tighten D011 failure-path test for single-line stderr Implementation Plan

## Overview

`TestRun_UnknownSubcommandFollowsD011Contract` (internal/cli/root_test.go) only asserts stderr matches `^error: `, which would still pass even if the message were empty after the prefix or contained embedded newlines — so a regression like the multi-line config decode errors fixed in task 0021 wouldn't be caught here. Per explicit instruction, this plan does not modify that existing test; it adds a new test alongside it that asserts the full tightened D011 contract on the same unknown-subcommand scenario.

## Current State Analysis

`TestRun_UnknownSubcommandFollowsD011Contract` (internal/cli/root_test.go:104-118) builds `deps` via `app.NewFakeDeps(time.Now())`, calls `cli.Run(deps, []string{"bogus"})`, and asserts: exit code `1`, empty stdout, and `stderr` matches `regexp.MustCompile(`^error: `)`. That regex only checks the prefix — it doesn't require any content after it, and doesn't forbid additional newlines.

### Key Discoveries:

- `internal/cli/root_test.go:104-118` — the existing test, left untouched per instruction; its structure (deps setup, `cli.Run(deps, []string{"bogus"})` call, three assertions) is the template for the new test.
- `internal/cli/run.go:14-19` — `cli.Run` writes `fmt.Fprintf(deps.Stderr, "error: %s\n", err)`, so a well-behaved single-line error naturally satisfies `^error: .+\n$` with exactly one `\n`.
- `internal/cli/run_test.go` (added in the `normalize-nested-conventional-decode-errors` change) already establishes the pattern of a second, separate `cli.Run`-based D011 test living alongside `root_test.go`'s original — this task's "add, don't modify" approach is consistent with that precedent.

## What We're NOT Doing

- Not modifying `TestRun_UnknownSubcommandFollowsD011Contract` — explicit instruction overriding the task file's literal wording ("Update ... in root_test.go"); a new test is added instead.
- Not changing `cli.Run` or any production code — this is a test-only addition; the unknown-subcommand path already produces single-line stderr today, so the new test is expected to pass immediately.
- Not adding a second scenario (e.g., the config-decode-error path already covered by `TestRun_InvalidNestedConventionalFollowsD011Contract` in run_test.go) — task 0023 specifically calls out tightening coverage of the unknown-subcommand path.

## Implementation Approach

Add `TestRun_UnknownSubcommandStderrIsSingleLine` to `internal/cli/root_test.go`, right after the existing `TestRun_UnknownSubcommandFollowsD011Contract`. It repeats the same `deps`/`cli.Run(deps, []string{"bogus"})` setup, then asserts: exit code `1`, empty stdout (parity with the existing test), stderr matches `^error: .+\n$`, and `strings.Count(stderr.String(), "\n") == 1`.

## Phase 1: Add tightened D011 test

### Overview

Adds one new test function asserting the full tightened D011 contract, without touching the existing test.

### Changes Required:

#### 1. New tightened-contract test

**File**: `internal/cli/root_test.go`

**Intent**: Prove the unknown-subcommand failure path produces stderr with a non-empty message and exactly one line, so a future regression to multi-line or empty-message output would be caught here, independent of the existing looser test.

**Contract**: Add `TestRun_UnknownSubcommandStderrIsSingleLine` immediately after `TestRun_UnknownSubcommandFollowsD011Contract` (after line 118). Reuse the same setup (`app.NewFakeDeps(time.Now())`, `cli.Run(deps, []string{"bogus"})`) and the same exit-code/empty-stdout assertions. Replace the loose stderr assertion with two checks: `regexp.MustCompile(`^error: .+\n$`).MatchString(stderr.String())` (non-empty message, exactly one trailing newline) and `strings.Count(stderr.String(), "\n") == 1` (no embedded newlines before the trailing one). Add `"strings"` to the existing import block (alongside `"regexp"`, `"testing"`, `"time"`).

### Success Criteria:

#### Automated Verification:

- Unit tests pass: `go test ./internal/cli/...`
- Full test suite passes: `go test ./...`
- Vet passes: `go vet ./...`

#### Manual Verification:

- None — this is a test-only addition with no production code change.

---

## Testing Strategy

### Unit Tests:

- New `TestRun_UnknownSubcommandStderrIsSingleLine` passes against current `cli.Run` behavior (unknown-subcommand path already produces well-formed single-line stderr).
- Existing `TestRun_UnknownSubcommandFollowsD011Contract` and all other `internal/cli` tests remain unchanged and passing.

### Manual Testing Steps:

None required.

## Performance Considerations

None — one additional table-free unit test.

## Migration Notes

None — test-only addition, no behavior or schema change.

## References

- Task: `tasks/0023-tighten-d011-failure-path-test-single-line.md`
- Existing test: `internal/cli/root_test.go:104-118` (`TestRun_UnknownSubcommandFollowsD011Contract`)
- D011 contract: `internal/cli/run.go:10-19`
- Related change: `context/changes/normalize-nested-conventional-decode-errors/` (added `internal/cli/run_test.go` with a second D011 test)

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Add tightened D011 test

#### Automated

- [x] 1.1 Unit tests pass: `go test ./internal/cli/...` — c7103c0
- [x] 1.2 Full test suite passes: `go test ./...` — c7103c0
- [x] 1.3 Vet passes: `go vet ./...` — c7103c0
