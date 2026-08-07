# Add Automated D011 Failure-Path Test — Plan Brief

> Full plan: `context/changes/add-automated-d011-failure-path-test/plan.md`

## What & Why

The D011 error contract (stderr `error: <message>`, exit code `1`, stdout untouched) currently lives only in `main.go` and is verified by manually running the compiled binary. This plan extracts that logic into a testable function and adds an automated test, closing the gap flagged in the bootstrap-go-cli impl review (F6).

## Starting Point

`main.go:11-19` calls `root.Execute()` and, on error, writes `error: <err>\n` to stderr and calls `os.Exit(1)` — logic that can't run in-process in a test because of `os.Exit`. `internal/cli/root.go` builds the command tree with cobra's own error/usage output silenced, so all D011 formatting responsibility sits in `main.go`.

## Desired End State

`internal/cli.Run(deps *app.Deps, args []string) int` holds the D011 logic; `main.go` becomes `os.Exit(cli.Run(deps, os.Args[1:]))`. A new test calls `Run` with an unknown subcommand and asserts exit code `1`, stderr matching `^error: `, and empty stdout. `go test ./internal/cli/...`, `task test`, and `task doctor` all pass.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) |
| --- | --- | --- |
| Extraction location | `internal/cli.Run(deps, args) int` | Matches the task's stated test location and keeps `main.go` a thin wrapper |
| Failure trigger | Real "unknown command" cobra error via `SetArgs([]string{"bogus"})` | Exercises the actual D011 path with no synthetic test-only command needed |
| Return type | `int` exit code | Mirrors `main`'s existing `os.Exit` contract 1:1 |

## Scope

**In scope:** New `internal/cli/run.go` with `Run`; `main.go` rewired to call it; new failure-path test in `internal/cli/root_test.go`.

**Out of scope:** Changing the D011 message/exit-code contract itself; testing other failure modes (flag errors, subcommand `RunE` errors); modifying `NewRootCmd`.

## Architecture / Approach

`main.go` currently inlines the error-handling contract. This plan lifts it into `cli.Run`, a pure function of `(deps, args) -> exit code` that both `main` and tests can call — `main` supplies real `os.Args`/`os.Exit`, tests supply fake args and inspect the returned code plus buffer-backed stderr/stdout from `app.NewFakeDeps`.

## Phases at a Glance

| Phase | What it delivers | Key risk |
| --- | --- | --- |
| 1. Extract `cli.Run` and add the D011 test | Testable error-handling function + automated failure-path test | None — pure refactor plus one test, existing tests unaffected |

**Prerequisites:** None.
**Estimated effort:** Single small change, one session.

## Open Risks & Assumptions

- Assumes cobra's "unknown command" error text is stable enough to reliably trigger a non-nil `Execute()` error (not asserting on its exact wording, only on stderr's `error: ` prefix).

## Success Criteria (Summary)

- `internal/cli` test suite includes an automated D011 failure-path case; no manual binary invocation required.
- `task test` and `task doctor` pass with the refactored `main.go`/new `run.go`.
