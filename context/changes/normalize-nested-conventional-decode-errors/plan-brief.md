# Normalize nested conventional decode errors to one line — Plan Brief

> Full plan: `context/changes/normalize-nested-conventional-decode-errors/plan.md`

## What & Why

When a `conventional:` mapping value has the wrong YAML type (e.g. `major: feat` instead of a list), `fileConventionalFromYAML`'s whole-map `Decode` call always returns the yaml library's multi-line error format, even for a single bad key. Through `cli.Run`, this breaks the D011 contract requiring exactly one line of stderr matching `error: <message>`. This plan replaces the whole-map decode with per-key iteration to guarantee a single-line, key-named error.

## Starting Point

`fileConventionalFromYAML` (config.go:56-82) already validates the top-level shape (0017/0019/0020: mapping-kind, alias resolution, empty-block checks) before decoding. The decode step itself hands the whole mapping to the yaml library in one call; a spike confirmed even a single bad value still produces the library's 2-line `"yaml: unmarshal errors:\n  line N: ..."` format — there's no single-line variant to simply surface.

## Desired End State

A wrong-type nested value (e.g. `major: feat`) produces `conventional.major must be a list of strings` — one line, naming the key. Through `cli.Run`, stderr is exactly `error: conventional.major must be a list of strings\n`. All prior top-level validation messages (0017/0019/0020) are unchanged.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) | Source |
|---|---|---|---|
| Message construction | Per-key iteration, build our own message | The yaml library never returns a single-line type error, even for one bad value — wrapping its text isn't viable | Plan |
| Multiple bad keys | Report only the first bad key found | Sufficient per the task's "when practical" phrasing; keeps the fix simple | Plan |
| cli.Run test | Include it | Only way to prove the D011 contract holds end-to-end, not just that `err.Error()` lacks newlines | Plan |

## Scope

**In scope:**
- Per-key decode loop in `fileConventionalFromYAML`
- Config-level regression test (single-line, key-named error)
- `cli.Run` end-to-end D011-contract test

**Out of scope:**
- Changes to top-level non-mapping/alias/empty-block validation
- Reporting all bad keys at once (only the first)
- Key-name validation against `allowedConventionalKeys` (already handled elsewhere)

## Architecture / Approach

Replace the single `doc.Conventional.Decode(&fileConv)` call with a loop over the already-alias-resolved mapping node's `.Content` pairs (key node, value node). Read the key from the key node's raw `.Value`; decode each value individually into `[]string`; on the first failure, return `fmt.Errorf("conventional.%s must be a list of strings", key)`.

## Phases at a Glance

| Phase | What it delivers | Key risk |
|---|---|---|
| 1. Per-key decode with single-line errors, plus regression tests | Fixed decode loop + config test + cli.Run test | None — isolated fix, full existing test suite as regression guard |

**Prerequisites:** None.
**Estimated effort:** ~30-45 minutes, single phase.

## Open Risks & Assumptions

- None identified — per-key decode produces identical results to the whole-map decode for valid configs; only the error path changes.

## Success Criteria (Summary)

- `go test ./internal/config/...` and `go test ./internal/cli/...` pass, including the two new tests
- Nested wrong-type `conventional:` values produce a single-line, key-named error
- `cli.Run` stderr for this error case is exactly one line matching `error: ...`
