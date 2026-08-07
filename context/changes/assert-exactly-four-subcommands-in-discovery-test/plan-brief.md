# Assert Exactly Four Subcommands in Discovery Test — Plan Brief

> Full plan: `context/changes/assert-exactly-four-subcommands-in-discovery-test/plan.md`

## What & Why

`TestNewRootCmd_CommandDiscovery` currently only checks that the four expected subcommands (`init`, `propose`, `apply`, `check`) are present — it never fails if a fifth, unexpected subcommand gets registered. This was flagged in the bootstrap-go-cli impl review (finding F5) as a gap versus the plan's "exactly four subcommands" requirement.

## Starting Point

`internal/cli/root_test.go:11-26` loops over `root.Commands()` marking a `want` map of 4 names as found, then asserts all 4 were found. `internal/cli/root.go:20-23` registers exactly 4 commands today, so the test currently passes — but nothing would catch a regression if a 5th command were added later.

## Desired End State

The test also asserts `len(root.Commands()) == 4`, with a readable failure message listing the actual command names if it ever doesn't. `task test` and `task doctor` both pass.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) |
| --- | --- | --- |
| Failure message content | Count + actual command name list | Self-diagnosing on failure, matches debugging style used elsewhere in this test file |
| Test restructuring | None — inline addition only | Matches acceptance criteria exactly, avoids scope creep on a low-risk passing test |
| Verification commands | `task test` + `task doctor` | Project's standard test/lint/vet gate, per user preference |

## Scope

**In scope:** One assertion added to `TestNewRootCmd_CommandDiscovery` in `internal/cli/root_test.go`.

**Out of scope:** Any change to `internal/cli/root.go`; any restructuring of the existing presence-loop logic; any other test in the file.

## Architecture / Approach

Single `if len(root.Commands()) != 4 { t.Errorf(...) }` check inserted right after the existing presence loop, before the existing `for name, found := range want` block.

## Phases at a Glance

| Phase | What it delivers | Key risk |
| --- | --- | --- |
| 1. Add the count assertion | Exact-count check in the discovery test | None — pure test-only change, near-zero risk |

**Prerequisites:** None.
**Estimated effort:** Single small edit, one session.

## Open Risks & Assumptions

- None identified — `NewRootCmd` is called without `Execute()`, so no cobra-injected `help`/`completion` command inflates the count.

## Success Criteria (Summary)

- Test fails if a 5th subcommand is ever registered without updating the test.
- `task test` and `task doctor` pass.
