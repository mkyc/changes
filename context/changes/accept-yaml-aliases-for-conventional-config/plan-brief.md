# Accept YAML aliases for conventional config mapping — Plan Brief

> Full plan: `context/changes/accept-yaml-aliases-for-conventional-config/plan.md`

## What & Why

Task 0019 added a guard that rejects non-mapping `conventional:` YAML values (sequences, scalars) with a friendly error. It didn't account for YAML aliases (`conventional: *conv`), which have `Kind == AliasNode` and get rejected by the same guard even when they resolve to a valid mapping — a regression from `921cbd3`. This plan resolves the alias chain before running the guard so valid alias-to-mapping configs load again, while alias-to-non-mapping configs still get the friendly error.

## Starting Point

`fileConventionalFromYAML` (config.go:56-78) checks the raw node's `Kind` against `yaml.MappingNode` before calling `Decode`. `Decode` itself already resolves aliases transparently and always has — only the guard added in 921cbd3 is unaware of aliases.

## Desired End State

`conventional: *conv` where `conv` anchors a mapping loads and merges with defaults exactly as a literal mapping would. An alias resolving to a sequence or scalar still returns `"conventional must be a mapping of keys to lists"`. Bare/null `conventional:` is unaffected.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) | Source |
|---|---|---|---|
| Alias resolution depth | Loop through `.Alias` until non-alias | Handles chained aliases with no extra complexity; YAML parser already prevents cycles | Plan |
| Where the fix lives | Only the guard's kind check, not `Decode` | `Decode` already resolves aliases correctly; changing it would be unnecessary | Plan |
| Test coverage | Alias-to-mapping (fix) + alias-to-sequence (still errors) | Directly covers both halves of the task's acceptance criteria | Plan |

## Scope

**In scope:**
- Alias resolution loop in `fileConventionalFromYAML`'s guard check
- Two new regression tests (alias-to-mapping, alias-to-sequence)

**Out of scope:**
- Alias support for other config keys
- Changes to `Decode` itself
- Artificial alias-chain-depth limits

## Architecture / Approach

Single local loop added to `fileConventionalFromYAML`: walk `doc.Conventional.Alias` while `Kind == yaml.AliasNode` to find the effective node, then run the existing mapping-kind guard against that resolved node instead of the raw one. `Decode` continues to run on the original node unchanged.

## Phases at a Glance

| Phase | What it delivers | Key risk |
|---|---|---|
| 1. Resolve aliases before the kind guard, with regression tests | Fixed guard + 2 new tests | None — isolated fix with full existing test suite as a regression guard |

**Prerequisites:** None.
**Estimated effort:** ~15-30 minutes, single phase.

## Open Risks & Assumptions

- None identified — `Decode`'s existing alias-resolution behavior is unchanged; this is a narrowly-scoped fix to the guard added in 921cbd3.

## Success Criteria (Summary)

- `go test ./internal/config/...` passes, including the two new alias tests
- Alias-to-mapping `conventional:` configs load successfully again
- Alias-to-non-mapping configs still produce the 0019 error message
