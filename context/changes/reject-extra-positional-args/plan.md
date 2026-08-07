# Reject Extra Positional Args on Stub Subcommands Implementation Plan

## Overview

Harden the `init`, `propose`, `apply`, and `check` stub subcommands so extra or mistyped positional arguments produce a non-zero exit and an error on stderr instead of being silently ignored, and add regression coverage.

## Current State Analysis

Each of the four subcommand constructors (`internal/cli/init.go`, `propose.go`, `apply.go`, `check.go`) builds a `*cobra.Command` with a `RunE` that only calls `config.Load` — none set an `Args` validator, so cobra defaults to `cobra.ArbitraryArgs`, which accepts any positional args and silently discards them. `internal/cli/root_test.go` has `TestSubcommands_SilentSuccessWithNoConfigFile` (lines 41-59), a table-driven test looping over all four command names that asserts exit success with empty stdout/stderr when run with no args. `root.go` sets `SilenceErrors: true` and `SilenceUsage: true` on the root command, and `main.go` prints `RunE` errors via `fmt.Fprintf(deps.Stderr, "error: %s\n", err)` before `os.Exit(1)` — so any error returned from `Args` or `RunE` already surfaces correctly through the existing error-handling path with no changes needed there.

## Desired End State

Running `changes init foo`, `changes propose bar`, `changes apply foo`, or `changes check foo` exits non-zero and prints an error to stderr. Running any of the four commands with zero positional args (with or without valid flags, e.g. `changes propose --since HEAD~1`) continues to succeed silently exactly as today.

### Key Discoveries:

- `cobra.NoArgs` (`github.com/spf13/cobra@v1.10.2/args.go:42-47`) returns `fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())` when `len(args) > 0` — this is the exact error text that will surface on stderr.
- `cobra.NoArgs` only validates positional args (`args []string`, post flag-parsing) — `propose`'s `--since` flag is unaffected and continues to work with zero positional args.
- `main.go` and `root.go` already propagate and print `RunE`/`Args` errors correctly; no changes needed outside `internal/cli/`.
- `TestSubcommands_SilentSuccessWithNoConfigFile` in `root_test.go:41-59` is the established table-driven pattern for exercising all four commands together via `app.NewFakeDeps` and `root.SetArgs`.

## What We're NOT Doing

- Not writing a custom `Args` validator or custom error message — using cobra's built-in `cobra.NoArgs` as specified by the ticket.
- Not changing `root.go`, `main.go`, or the error-printing/exit-code path — it already works correctly for `RunE` errors and `Args` errors alike.
- Not adding `Args` validators to the root `changes` command itself — only the four named stub subcommands per the ticket.
- Not changing `propose`'s `--since` flag behavior.

## Implementation Approach

Add `Args: cobra.NoArgs` to each of the four `*cobra.Command` literals, then extend `root_test.go` with a new table-driven test mirroring the existing `TestSubcommands_SilentSuccessWithNoConfigFile` pattern, asserting an error and non-empty stderr when each command is invoked with one extra positional arg.

## Phase 1: Harden stub subcommands and add regression test

### Overview

Set `Args: cobra.NoArgs` on all four subcommands and add a test proving extra positional args are rejected while the existing no-args success path stays green.

### Changes Required:

#### 1. Stub subcommand hardening

**File**: `internal/cli/init.go`

**Intent**: Reject positional arguments on `init`.

**Contract**: Add `Args: cobra.NoArgs` field to the `*cobra.Command` literal returned by `newInitCmd`.

**File**: `internal/cli/propose.go`

**Intent**: Reject positional arguments on `propose` without affecting the `--since` flag.

**Contract**: Add `Args: cobra.NoArgs` field to the `*cobra.Command` literal returned by `newProposeCmd`, alongside the existing `--since` flag registration.

**File**: `internal/cli/apply.go`

**Intent**: Reject positional arguments on `apply`.

**Contract**: Add `Args: cobra.NoArgs` field to the `*cobra.Command` literal returned by `newApplyCmd`.

**File**: `internal/cli/check.go`

**Intent**: Reject positional arguments on `check`.

**Contract**: Add `Args: cobra.NoArgs` field to the `*cobra.Command` literal returned by `newCheckCmd`.

#### 2. Regression test

**File**: `internal/cli/root_test.go`

**Intent**: Prove each stub subcommand now rejects an extra positional argument with a non-zero-equivalent error and non-empty stderr, matching the acceptance criteria in `tasks/0012-...md`.

**Contract**: Add a new test function (e.g. `TestSubcommands_RejectExtraPositionalArgs`) using the same table-driven, `t.Run` per-command structure as `TestSubcommands_SilentSuccessWithNoConfigFile` (`root_test.go:41-59`): for each of `init`, `propose`, `apply`, `check`, call `root.SetArgs([]string{name, "extra-arg"})`, assert `root.Execute()` returns a non-nil error, and assert `stderr.Len() != 0` (via `main.go`'s error-printing convention being exercised at the `cobra.Command.Execute()` level — note `Execute()` returns the error directly rather than writing to stderr itself, since `SilenceErrors`/`SilenceUsage` are set on root; the test should assert on the returned error's content, not on `stderr` buffer state, mirroring how cobra actually surfaces `Args` validation failures here).

### Success Criteria:

#### Automated Verification:

- All CLI tests pass: `go test ./internal/cli/...`
- New test passes and existing test still passes: `go test ./internal/cli/... -run 'TestSubcommands' -v`
- Full test suite passes: `task test` (or `go test ./...`)
- Lint/vet/fmt clean: `task doctor`

#### Manual Verification:

- Build the binary and run `bin/changes init foo`; confirm non-zero exit and an error message printed.
- Run `bin/changes propose --since HEAD~1` (no positional args) and confirm it still succeeds silently as before.

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful before proceeding to the next phase.

---

## Testing Strategy

### Unit Tests:

- Table-driven test over `init`, `propose`, `apply`, `check`, each invoked with one extra positional arg, asserting `Execute()` returns an error.
- Existing `TestSubcommands_SilentSuccessWithNoConfigFile` continues to pass unchanged, proving no-args behavior is unaffected.

### Integration Tests:

- N/A — CLI unit tests using `app.NewFakeDeps` already exercise the full command tree end-to-end within the test process.

### Manual Testing Steps:

1. `task build`, then `bin/changes init foo` — expect non-zero exit and an error on stderr.
2. `bin/changes propose bar` — expect non-zero exit and an error on stderr.
3. `bin/changes apply foo` — expect non-zero exit and an error on stderr.
4. `bin/changes check foo` — expect non-zero exit and an error on stderr.
5. `bin/changes propose --since HEAD~1` — expect silent success (exit 0, no output), confirming the flag still works with zero positional args.

## Performance Considerations

None — this is a validation-only change with no runtime cost beyond cobra's existing arg-parsing.

## Migration Notes

None.

## References

- Ticket: `tasks/0012-reject-extra-positional-args-on-stub-subcommands.md`
- cobra.NoArgs implementation: `github.com/spf13/cobra@v1.10.2/args.go:42-47`
- Existing test pattern: `internal/cli/root_test.go:41-59`

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Harden stub subcommands and add regression test

#### Automated

- [x] 1.1 All CLI tests pass: `go test ./internal/cli/...` — 8f73635
- [x] 1.2 New test passes and existing test still passes: `go test ./internal/cli/... -run 'TestSubcommands' -v` — 8f73635
- [x] 1.3 Full test suite passes: `task test` — 8f73635
- [x] 1.4 Lint/vet/fmt clean: `task doctor` — 8f73635

#### Manual

- [x] 1.5 `bin/changes init foo` exits non-zero with an error message — 8f73635
- [x] 1.6 `bin/changes propose --since HEAD~1` (no positional args) still succeeds silently — 8f73635
