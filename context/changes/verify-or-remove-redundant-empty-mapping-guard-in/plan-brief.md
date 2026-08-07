# Verify Or Remove Redundant Empty-Mapping Guard — Plan Brief

> Full plan: `context/changes/verify-or-remove-redundant-empty-mapping-guard-in/plan.md`

## What & Why

Impl review F1 flagged an empty-mapping guard in `fileConventionalFromYAML` (`internal/config/config.go:62-64`) as possibly dead code, and asked to either verify it's reachable and correct with a test, or remove it in favor of the post-decode check that already covers the same ground. Direct verification during planning confirmed the guard is reachable (for `conventional: {}`) but genuinely redundant — the post-decode check catches the same input with the same error message.

## Starting Point

`fileConventionalFromYAML` has two paths that can both reject an empty `conventional` block: a `Kind == yaml.MappingNode && len(Content) == 0` guard (lines 62-64) that only fires for `conventional: {}`, and a post-decode `len(fileConv) == 0` check (lines 70-72) that fires for both `conventional:` (bare) and, as verified, `conventional: {}` too. No existing test exercises the `{}` form.

## Desired End State

The guard is removed; the post-decode check alone handles both empty forms. A new test (`conventional: {}`) proves the error message is unchanged after removal. `go build ./internal/config/...`, `task test`, and `task doctor` all pass with no behavior change.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) |
| --- | --- | --- |
| Verify vs. remove | Remove the guard + add regression test | Direct decoder verification proved the post-decode check alone already produces the identical error for `{}` |
| Test assertion strength | Exact error-message match via `strings.Contains` | Proves the post-decode path yields the *same* user-facing message, not just "some error" — the real regression risk of removing the guard |

## Scope

**In scope:** Delete the guard in `internal/config/config.go`; add one regression test in `internal/config/config_test.go` for the `conventional: {}` form.

**Out of scope:** Changing the post-decode check or error message; covering other YAML node kinds under `conventional:`.

## Architecture / Approach

`fileConventionalFromYAML` currently has a specific guard for one empty-input shape (`{}`) and a general post-decode guard for "empty after decode" that already subsumes it. Removing the specific guard collapses both empty-input shapes onto the single general check, verified via direct decoder testing to produce identical output.

## Phases at a Glance

| Phase | What it delivers | Key risk |
| --- | --- | --- |
| 1. Remove the guard and add the regression test | Dead-code removal + closed test-coverage gap | None — behavior proven identical via direct decoder verification before implementation |

**Prerequisites:** None.
**Estimated effort:** Single small change, one session.

## Open Risks & Assumptions

- None identified — the guard's redundancy was confirmed via a standalone script exercising the actual `go.yaml.in/yaml/v3` decode path used by production code, not inferred from reading alone.

## Success Criteria (Summary)

- Either a passing test exists for `conventional: {}`, or the guard is removed and behavior is unchanged — both are true here: the guard is removed and a new test covers `{}`.
- `go build ./internal/config/...` and `go test ./internal/config/...` pass.
