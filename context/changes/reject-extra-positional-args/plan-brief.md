# Reject Extra Positional Args on Stub Subcommands — Plan Brief

> Full plan: `context/changes/reject-extra-positional-args/plan.md`

## What & Why

`changes init foo`, `changes propose bar`, etc. currently exit 0 with no output, silently ignoring mistyped or extra positional arguments. Fix: set `Args: cobra.NoArgs` on `init`, `propose`, `apply`, and `check` so extra args produce a non-zero exit and an error message instead of failing silently.

## Starting Point

None of the four stub subcommands set an `Args` validator, so cobra defaults to `cobra.ArbitraryArgs` (accepts and discards anything). `root.go` and `main.go` already propagate and print `RunE`/`Args` errors correctly — no changes needed there. `root_test.go` already has a table-driven test (`TestSubcommands_SilentSuccessWithNoConfigFile`) covering all four commands with no args.

## Desired End State

Extra positional args on any of the four commands cause a non-zero exit with an error on stderr. Zero-arg invocations (including `propose --since X`) behave exactly as before.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) | Source |
|---|---|---|---|
| Error message | Use cobra's built-in `cobra.NoArgs` message (`unknown command %q for %q`) | Zero custom code, matches cobra idioms, ticket doesn't ask for custom wording. | Plan |
| Test structure | Table-driven, one subtest per command in `root_test.go` | Mirrors the existing `TestSubcommands_SilentSuccessWithNoConfigFile` pattern exactly. | Plan |

## Scope

**In scope:**
- `Args: cobra.NoArgs` on `init.go`, `propose.go`, `apply.go`, `check.go`
- New table-driven regression test in `root_test.go`

**Out of scope:**
- Custom error messages/validators
- Changes to `root.go` / `main.go` error handling (already correct)
- `Args` validation on the root `changes` command itself
- Any change to `propose`'s `--since` flag behavior

## Architecture / Approach

Single flat change: add one field (`Args: cobra.NoArgs`) to each of the four `cobra.Command` literals, then extend the existing test file with a mirrored table-driven test asserting extra positional args now error.

## Phases at a Glance

| Phase | What it delivers | Key risk |
|---|---|---|
| 1. Harden stub subcommands + regression test | All 4 commands reject extra args; new test proves it; existing no-args test stays green | Low — cobra's `Args` validation is well-understood; main risk is a typo in one of the 4 files |

**Prerequisites:** None — self-contained within `internal/cli/`.
**Estimated effort:** Single session, one phase.

## Open Risks & Assumptions

- Assumes `cobra.NoArgs`'s default error text is acceptable to the user (not a custom message) — confirmed during planning questions.

## Success Criteria (Summary)

- `changes init foo` / `propose bar` / `apply foo` / `check foo` all exit non-zero with an error.
- `go test ./internal/cli/...` passes including the new test.
- Existing silent-success behavior with no extra args is unchanged.
