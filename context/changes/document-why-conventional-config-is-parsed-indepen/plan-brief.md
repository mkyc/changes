# Document Why Conventional Config Is Parsed Independently — Plan Brief

> Full plan: `context/changes/document-why-conventional-config-is-parsed-indepen/plan.md`

## What & Why

`internal/config/config.go` parses the config file's YAML twice: once via Viper for scalar keys, once directly via `go.yaml.in/yaml/v3` in `fileConventionalFromYAML` for the `conventional` block. This is deliberate — Viper can't distinguish "key absent" from "key present but empty" — but undocumented, so impl review F2 flagged it as something a future reader could mistake for an oversight.

## Starting Point

`fileConventionalFromYAML` (`internal/config/config.go:52`) has no comment explaining its independent parsing. The rationale (Viper's typed getters collapse the null-vs-empty distinction that the function's `len(fileConv) == 0` empty-block check depends on) was already confirmed experimentally in the sibling change `verify-or-remove-redundant-empty-mapping-guard-in`.

## Desired End State

`fileConventionalFromYAML` has a doc comment stating why it bypasses Viper for this one key. No behavior change; `go build` and `go test` on `internal/config` pass exactly as before.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) |
| --- | --- | --- |
| Comment placement | Doc comment on `fileConventionalFromYAML` itself | Co-located with the exact code doing the independent parsing, matching the task's primary suggestion |

## Scope

**In scope:** One doc comment added to `internal/config/config.go`.

**Out of scope:** Any parsing logic change; comments at other call sites (e.g. `v.ReadConfig` in `Load`).

## Architecture / Approach

Standard Go doc comment above the function signature, stating the Viper-limitation rationale directly.

## Phases at a Glance

| Phase | What it delivers | Key risk |
| --- | --- | --- |
| 1. Add the doc comment | Documented rationale for the dual-parse pattern | None — comment-only change |

**Prerequisites:** None.
**Estimated effort:** Single trivial edit, one session.

## Open Risks & Assumptions

- None — this is a documentation-only change with no behavioral surface.

## Success Criteria (Summary)

- Comment added; no behavior change.
- `go build ./internal/config/...` and `go test ./internal/config/...` pass.
