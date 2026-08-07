# Improve error message for non-mapping conventional config value — Plan Brief

> Full plan: `context/changes/improve-error-message-for-non-mapping-conventional/plan.md`

## What & Why

When `.changes/config.yaml`'s `conventional:` key is given a YAML sequence or scalar instead of a mapping, `internal/config/config.go` currently leaks the raw yaml-library decode error (e.g. "cannot unmarshal !!seq into map[string][]string") instead of one of the file's own descriptive validation messages. This plan adds a pre-decode node-kind check so the error is actionable and consistent with the file's other validation errors.

## Starting Point

`fileConventionalFromYAML` (config.go:56-75) already distinguishes an absent `conventional:` key (`Kind == 0`) from a present one before decoding, and separately validates for an empty mapping after decoding. It has no check for "present but wrong shape" — that case currently falls straight through to `Decode`.

## Desired End State

Any non-mapping `conventional:` value returns `"conventional must be a mapping of keys to lists"`, matching the tone of the existing `"conventional must contain at least one key"` and `"unknown conventional key %q"` messages. Verified entirely by automated tests — no user-facing behavior changes for valid configs.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) | Source |
|---|---|---|---|
| Error message wording | `"conventional must be a mapping of keys to lists"` | Matches the task's suggested wording and the file's existing terse style | Plan |
| Shapes covered | Any non-mapping kind via `Kind != yaml.MappingNode` | One condition covers sequence, scalar, and alias uniformly | Plan |
| Test coverage | Two subtests — sequence and scalar | Cheap to add both, gives explicit coverage of the two example shapes named in the task | Plan |
| Check placement | Right after the existing absent-check, before `Decode` | Keeps all node-kind validation grouped at the top of the function, avoids parsing yaml error strings | Plan |

## Scope

**In scope:**
- Node-kind guard in `fileConventionalFromYAML` (config.go)
- Two new test subtests in config_test.go (sequence, scalar)

**Out of scope:**
- Changes to the empty-mapping or unknown-key validation
- Echoing the actual YAML kind in the error message
- Special-casing YAML aliases

## Architecture / Approach

Single guard clause inserted between the existing absent-check and the `Decode` call in `fileConventionalFromYAML`. No new functions, no API changes.

## Phases at a Glance

| Phase | What it delivers | Key risk |
|---|---|---|
| 1. Add node-kind guard and test coverage | Descriptive error + 2 new tests | None — isolated, single-file change with full existing test coverage as a regression guard |

**Prerequisites:** None — no dependencies on other in-flight work.
**Estimated effort:** ~15-30 minutes, single phase.

## Open Risks & Assumptions

- None identified — this is additive validation with no schema or behavior change for valid inputs.

## Success Criteria (Summary)

- `go test ./internal/config/...` passes, including two new tests for sequence and scalar `conventional:` values
- Non-mapping `conventional:` configs now surface a clear, actionable error instead of a raw yaml library error
