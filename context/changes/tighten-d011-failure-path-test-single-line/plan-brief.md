# Tighten D011 failure-path test for single-line stderr — Plan Brief

> Full plan: `context/changes/tighten-d011-failure-path-test-single-line/plan.md`

## What & Why

`TestRun_UnknownSubcommandFollowsD011Contract` only asserts stderr matches `^error: `, which passes even for an empty message or embedded newlines — so a regression like the multi-line config decode errors fixed in task 0021 wouldn't be caught by this test. Per explicit instruction, this plan does not modify that existing test; it adds a new one asserting the full tightened contract.

## Starting Point

`internal/cli/root_test.go:104-118` has the existing loose test: `cli.Run(deps, []string{"bogus"})`, asserting exit code 1, empty stdout, and a `^error: ` prefix match only.

## Desired End State

A new `TestRun_UnknownSubcommandStderrIsSingleLine` test exists alongside the original, asserting: non-empty message after `error: `, exactly one trailing newline, and no embedded newlines — using the same unknown-subcommand scenario. The original test is untouched.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) | Source |
|---|---|---|---|
| Modify vs add | Add a new test, leave the original alone | Explicit user instruction overriding the task file's literal wording | Plan |
| Scenario | Reuse the unknown-subcommand path | Task 0023 specifically calls out tightening that path's coverage | Plan |
| Assertion style | Regex `^error: .+\n$` + `strings.Count(...) == 1` | Matches the task's own suggested wording; the count check makes "no embedded newlines" unambiguous | Plan |

## Scope

**In scope:**
- New test `TestRun_UnknownSubcommandStderrIsSingleLine` in `internal/cli/root_test.go`

**Out of scope:**
- Modifying the existing `TestRun_UnknownSubcommandFollowsD011Contract`
- Any production code changes
- Testing other error scenarios (already covered elsewhere, e.g. `run_test.go`)

## Architecture / Approach

Single new test function, same setup pattern as the existing test (`app.NewFakeDeps` + `cli.Run(deps, []string{"bogus"})`), with two stderr assertions replacing the original's single loose one.

## Phases at a Glance

| Phase | What it delivers | Key risk |
|---|---|---|
| 1. Add tightened D011 test | New passing test, existing test untouched | None — test-only, current behavior already satisfies the tightened contract |

**Prerequisites:** None.
**Estimated effort:** ~10 minutes, single phase.

## Open Risks & Assumptions

- None — the unknown-subcommand path already produces well-formed single-line stderr, so the new test is expected to pass immediately without any production change.

## Success Criteria (Summary)

- `go test ./internal/cli/...` passes, including the new test
- The tightened test would fail if stderr had an empty message or embedded newlines
- The original test remains unmodified
