# Add Automated D011 Failure-Path Test Implementation Plan

## Overview

Automate verification of the D011 error contract (stderr `error: <message>`, exit code `1`, stdout untouched) for an unknown subcommand, which today is exercised only by manually running the compiled binary (Progress 4.7 in the bootstrap-go-cli plan). We do this by extracting `main.go`'s error-handling logic into a testable `cli.Run` function and adding an in-process test against it.

## Current State Analysis

`main.go:11-19` builds the root command via `cli.NewRootCmd`, calls `root.Execute()`, and on error writes `error: <err>\n` to `deps.Stderr` and calls `os.Exit(1)`. This is the only place the D011 stderr-format contract is implemented, and it can't be exercised in-process today because `os.Exit` would kill the test binary.

`internal/cli/root.go:12-26` builds the command tree with `SilenceErrors: true, SilenceUsage: true`, so cobra itself never writes error/usage output — all D011 formatting responsibility sits in `main.go`.

`internal/cli/root_test.go` already has the fixture pattern for this: `app.NewFakeDeps(time.Now())` returns buffer-backed `stdout`/`stderr`, and existing tests (e.g. `TestSubcommands_RejectExtraPositionalArgs`, lines 77-93) already call `root.SetArgs(...)` + `root.Execute()` and assert on the returned error — but none currently assert on the stderr-formatting/exit-code contract, since that logic isn't in the `cli` package yet.

## Desired End State

`internal/cli` exposes a `Run(deps *app.Deps, args []string) int` function that `main.go` delegates to. A new test in `internal/cli/root_test.go` invokes `Run` with an unrecognized subcommand and asserts: return value is `1`, `stderr` matches `^error: `, and `stdout` is empty. `go test ./internal/cli/...` passes, `task test` and `task doctor` pass, and manual binary invocation is no longer required to verify D011.

### Key Discoveries:

- `main.go:15-18` — the exact error-handling logic to extract (`if err := root.Execute(); err != nil { ...; os.Exit(1) }`).
- `internal/cli/root.go:12-26` — `NewRootCmd` already takes `*app.Deps` and returns `*cobra.Command`; `Run` wraps this.
- `internal/app/fake.go:14-24` — `NewFakeDeps` is the existing fixture for buffer-backed stdout/stderr, reused by the new test.
- Cobra's built-in "unknown command" error path fires naturally when `SetArgs` includes a subcommand name not registered on root — no synthetic test-only command is needed to trigger a D011 failure case.

## What We're NOT Doing

- Not changing the D011 contract itself (message format, exit code) — only making it testable.
- Not adding a test-only fake subcommand to synthesize failures — the real "unknown subcommand" cobra error path is used directly.
- Not modifying `internal/cli/root.go`'s `NewRootCmd` signature or behavior.
- Not testing every possible failure mode (e.g., flag-parsing errors, subcommand `RunE` errors) — scope is limited to the unknown-subcommand case named in the task.

## Implementation Approach

Add a new exported `Run` function to the `cli` package (new file `internal/cli/run.go`) that wraps `NewRootCmd`, sets args, executes, and on error writes the `error: <message>\n` line to `deps.Stderr` before returning `1`; returns `0` on success. Rewire `main.go` to call `os.Exit(cli.Run(deps, os.Args[1:]))`, removing the duplicated logic. Add a new test function to `internal/cli/root_test.go` that calls `cli.Run` with `[]string{"bogus"}` and asserts the D011 contract via a regex match on stderr, an empty stdout, and a returned exit code of `1`.

## Phase 1: Extract `cli.Run` and add the D011 failure-path test

### Overview

Move the error-handling contract out of `main.go` into a testable `cli.Run` function, then add the automated D011 failure-path test against it.

### Changes Required:

#### 1. Extracted, testable entry point

**File**: `internal/cli/run.go` (new file)

**Intent**: House the D011 error-handling contract (stderr `error: <message>\n` on failure, silent on success) as an exported function that both `main.go` and tests can call without triggering `os.Exit`.

**Contract**: `func Run(deps *app.Deps, args []string) int` — builds the root command via `NewRootCmd(deps)`, calls `root.SetArgs(args)` then `root.Execute()`; on non-nil error, writes `fmt.Fprintf(deps.Stderr, "error: %s\n", err)` and returns `1`; otherwise returns `0`. This is a direct lift of the current `main.go:15-18` body, parameterized by `args` instead of implicitly using `os.Args`.

#### 2. Thin main

**File**: `main.go`

**Intent**: Delegate all command execution and error handling to `cli.Run`, keeping `main` a thin wrapper around process args/exit.

**Contract**: `main` becomes `deps := app.NewRealDeps(); os.Exit(cli.Run(deps, os.Args[1:]))`. The `fmt` import is no longer needed in `main.go` once the `Fprintf` call moves to `run.go`.

#### 3. D011 failure-path test

**File**: `internal/cli/root_test.go`

**Intent**: Assert the full D011 contract — stderr prefix, empty stdout, exit code — for an unrecognized subcommand, closing the gap flagged in impl review F6.

**Contract**: New test function (e.g. `TestRun_UnknownSubcommandFollowsD011Contract`) that: builds `deps, stdout, stderr := app.NewFakeDeps(time.Now())`; calls `code := cli.Run(deps, []string{"bogus"})`; asserts `code == 1`; asserts `stdout.Len() == 0`; asserts `stderr.String()` matches `^error: ` via `regexp.MustCompile(`^error: `).MatchString(...)` (add the `regexp` import to the test file).

### Success Criteria:

#### Automated Verification:

- `task test` passes (runs `go test ./... -v`, includes the new `internal/cli` test)
- `task doctor` passes (gofmt, go vet, go mod tidy check, golangci-lint)

#### Manual Verification:

- None — this closes the gap that previously required manual binary invocation; no user-facing behavior changes (main's runtime behavior for real `os.Args`/`os.Exit` is unchanged, only refactored into a callable function).

**Implementation Note**: After completing this phase and automated verification passes, this plan is complete — there is no further phase.

---

## Testing Strategy

### Unit Tests:

- New `TestRun_UnknownSubcommandFollowsD011Contract` in `internal/cli/root_test.go`, covering: exit code `1`, stderr regex `^error: `, empty stdout.
- Existing tests in `root_test.go` continue to pass unchanged (they call `NewRootCmd` + `Execute` directly, not through `Run`, and are unaffected by the extraction).

### Integration Tests:

- None needed — this is an in-process unit test replacing the need for manual binary invocation.

### Manual Testing Steps:

- None.

## Performance Considerations

None — this is a pure refactor plus one test.

## Migration Notes

None — `main.go`'s real-world behavior (real `os.Args`, real `os.Exit`) is unchanged; only the internal structure changes.

## References

- Task: `tasks/0015-add-automated-d011-failure-path-test.md`
- Impl review finding: `context/changes/bootstrap-go-cli/reviews/impl-review.md` (F6)
- D011 decision: `docs/spec/DECISIONS.md` (D011: Error output format)
- Current error-handling logic: `main.go:11-19`
- Root command construction: `internal/cli/root.go:12-26`
- Fake deps fixture: `internal/app/fake.go:14-24`

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Extract `cli.Run` and add the D011 failure-path test

#### Automated

- [x] 1.1 `task test` passes
- [x] 1.2 `task doctor` passes
