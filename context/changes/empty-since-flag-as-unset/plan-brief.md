# Treat Empty Since Values as Unset — Plan Brief

> Full plan: `context/changes/empty-since-flag-as-unset/plan.md`

## What & Why

An explicitly empty `--since` currently enters Viper's highest-precedence override layer and replaces both a configured ref and the built-in `main` default with an unusable empty string. This plan treats exact-empty `since` values from CLI and YAML as unset while preserving meaningful values and established config precedence.

## Starting Point

`config.Load` already layers built-in defaults, `.changes/config.yaml`, and CLI flags in the correct order. The defect is limited to empty strings being considered meaningful values; existing tests cover defaults and non-empty overrides but not empty CLI or YAML inputs.

## Desired End State

Empty CLI values fall back to a non-empty config-file value and then to `main`. Quoted empty YAML also falls back to `main`, while non-empty CLI values continue to win and whitespace-only values remain explicit.

## Key Decisions Made

| Decision | Choice | Why |
|---|---|---|
| Empty CLI semantics | Treat exact `""` as unset | Matches the ticket without silently trimming malformed refs |
| Empty YAML semantics | Treat quoted empty as unset | Provides consistent empty-value behavior across both configuration sources |
| Fallback coverage | Test file and built-in fallback layers | Prevents an incorrect fix that always forces `main` |
| Verification boundary | Config-loader tests only | `propose` currently discards the resolved config, so CLI/binary tests cannot observe the value |
| Delivery shape | One implementation phase | The two-file fix and its precedence tests form one cohesive change |

## Scope

**In scope:**

- Exact-empty CLI `since` values are treated as absent.
- Quoted empty YAML `since` values are treated as absent.
- Regression tests cover fallback to a config-file value and to `main`.
- Existing non-empty CLI precedence remains protected.

**Out of scope:**

- Whitespace trimming or Git-ref validation.
- Other config keys or general empty-value normalization.
- CLI, binary, or process-level integration tests.
- Changes to how `propose` consumes configuration.

## Architecture / Approach

Retain the existing defaults → file → CLI pipeline. Skip exact-empty CLI overrides so the file layer remains visible, and normalize any exact-empty final file-derived value to the built-in default without interfering with a later non-empty CLI override.

## Phases at a Glance

| Phase | What it delivers | Key risk |
|---|---|---|
| 1. Normalize empty since values and add regression coverage | Empty-aware resolution in `config.Load` plus CLI/file fallback tests | Incorrect ordering could force `main` and bypass a configured non-empty ref |

**Prerequisites:** Existing config loader and unit-test helpers; no external dependencies.
**Estimated effort:** One session, one phase.

## Open Risks & Assumptions

- Whitespace-only refs remain explicit and may be rejected later by the planned Git adapter.
- Bare YAML `since:` already falls through to the default; the new YAML regression specifically targets quoted `since: ""`.
- The expanded YAML behavior is intentional even though task 0013 originally named the CLI flag path.

## Success Criteria (Summary)

- Empty CLI values fall back to `develop` from the file when present, otherwise to `main`.
- Quoted empty YAML resolves to `main`, and non-empty CLI override behavior remains green.
- Config tests, the full suite, and project doctor checks all pass.
