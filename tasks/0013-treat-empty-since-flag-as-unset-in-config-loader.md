---
title: "Treat empty --since flag as unset in config loader"
id: "0013"
status: completed
priority: medium
type: bug
tags: ["config", "bootstrap-go-cli"]
created_at: "2026-08-05"
depends_on: ["0001"]
---

# Treat empty --since flag as unset in config loader

## Objective

Prevent an empty `--since` CLI flag from overriding the `"main"` default ref
with an empty string.

## Steps to Reproduce

1. Run `changes propose --since=""` (or invoke `config.Load` with a flag set
   where `since` is changed but empty).

## Expected Behavior

When `--since` is not given a meaningful value, config resolution falls back
to the file value or the built-in default `"main"`.

## Actual Behavior

`flags.Changed("since")` is true even for `--since=""`, so the loader sets
`Since` to `""`, bypassing defaults.

## Tasks

- [x] In `config.Load`, skip the CLI override when `since == ""` even if
      `flags.Changed("since")` is true (recommended Fix A from impl review).
- [x] Add a unit test asserting empty `--since` preserves the default `"main"`.

## Acceptance Criteria

- `config.Load` with `since` changed to `""` resolves `Since` to `"main"` (or
  the config-file value if set).
- `go test ./internal/config/...` passes including the new test.
- Existing CLI-override test (`TestLoad_CLIFlagWinsOverConfigFile`) remains
  green for non-empty values.

## References

- Found during impl review F3: `context/changes/bootstrap-go-cli/reviews/impl-review.md`
- Location: `internal/config/config.go:64-66`
- Recommended fix: Fix A (treat empty string as unset)
