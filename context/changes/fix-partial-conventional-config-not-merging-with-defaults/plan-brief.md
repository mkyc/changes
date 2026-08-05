# Fix Partial Conventional Config Merge — Plan Brief

> Full plan: `context/changes/fix-partial-conventional-config-not-merging-with-defaults/plan.md`
> Research: (none — task 0011 + impl review F1 provided sufficient context)

## What & Why

When `.changes/config.yaml` specifies only a subset of `conventional` keys
(e.g. just `major`), Viper replaces the entire default map and silently drops
`minor`, `patch`, and `none`. This violates D010 partial-fallback semantics
already working for scalar config keys. We fix the loader to deep-merge
conventional keys onto built-in defaults.

## Starting Point

`internal/config/config.go` sets conventional defaults via Viper then reads
the config file. Lines 69–72 copy Viper's conventional map directly into
`Config.Conventional` with no merge step. Scalar partial fallback is tested;
conventional partial fallback is not.

## Desired End State

`config.Load` returns all four default conventional keys when the config file
overrides only some of them. File slices replace defaults per key (empty slice
= explicit disable). Unknown keys and an empty `conventional:` block produce
clear errors. Tests lock in the behavior.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) | Source |
| -------- | ------ | ---------------- | ------ |
| Merge strategy | Per-key slice replacement | File values fully replace defaults for that key; matches scalar D010 semantics | Plan |
| Empty slice | Honor as explicit disable | `major: []` is intentional — user wants no major triggers | Plan |
| Unknown keys | Reject with error | Only `major`/`minor`/`patch`/`none` are documented; fail fast on typos | Plan |
| Empty `conventional:` block | Return validation error | Empty block is almost certainly a mistake, not a valid config | Plan |
| Scope | Conventional merge only | Bug is isolated; no other nested Viper defaults exist today | Plan |

## Scope

**In scope:**
- Deep-merge logic in `config.Load` for `conventional`
- Validation for unknown keys and empty conventional block
- Unit tests for merge + edge cases
- Extend scalar partial-fallback test to assert conventional defaults

**Out of scope:**
- Empty `--since` flag fix (task 0013)
- Extra positional args validation (task 0012)
- CLI overrides for conventional keys
- Generic nested-map merge infrastructure

## Architecture / Approach

After Viper reads the config file, clone `defaultConventional()` as the merge
base, overlay each file key (validating against the allowed set), and assign
the result. Empty file map after a config read triggers an error. No Viper
behavior changes — just post-processing of its output.

## Phases at a Glance

| Phase | What it delivers | Key risk |
| ----- | ---------------- | -------- |
| 1. Deep-merge conventional defaults | Merge + validation in `config.go` | Distinguishing "key absent from file" vs "empty block in file" |
| 2. Tests | Partial merge + validation edge-case tests | Over-testing Viper internals vs behavior |

**Prerequisites:** Bootstrap Go CLI (task 0001) complete; `internal/config` package exists.
**Estimated effort:** ~1 session, 2 phases

## Open Risks & Assumptions

- Detecting empty `conventional:` relies on "config file read + zero keys in
  Viper map" — assumes Viper does not strip the key when the block is empty
  (verified behavior in impl review repro).
- First validation errors in `internal/` — error message wording is new
  convention for this codebase.

## Success Criteria (Summary)

- Partial `conventional.major` in YAML yields four-key map with defaults for
  unset keys
- Unknown keys and empty `conventional:` block error clearly
- `go test ./internal/config/...` passes including new tests
