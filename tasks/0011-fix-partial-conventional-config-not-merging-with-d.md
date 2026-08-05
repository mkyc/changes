---
title: "Fix partial conventional config not merging with defaults"
id: "0011"
status: in-progress
priority: high
type: bug
tags: ["config", "bootstrap-go-cli"]
created_at: "2026-08-05"
depends_on: ["0001"]
---

# Fix partial conventional config not merging with defaults

## Objective

Fix D010 config precedence for nested `conventional` keys so a partial
`.changes/config.yaml` merges with built-in defaults instead of replacing the
entire map.

## Steps to Reproduce

1. Create a `.changes/config.yaml` with only a subset of conventional keys:

   ```yaml
   conventional:
     major:
       - custom
   ```

2. Run `config.Load` (or any subcommand that loads config).

## Expected Behavior

Per Phase 3 of the bootstrap plan and D010, unset conventional keys fall back
to built-in defaults. The resolved `Config.Conventional` map should contain
all four keys (`major`, `minor`, `patch`, `none`) with `major` overridden and
the rest from defaults.

## Actual Behavior

Viper replaces the entire `conventional` default map. Only `major: [custom]`
remains; `minor`, `patch`, and `none` defaults are lost.

## Tasks

- [ ] After `ReadConfig`, deep-merge file `conventional` keys onto
      `defaultConventional()` (per-key slice replacement, not whole-map
      replacement).
- [ ] Add `TestLoad_PartialConventionalFallsBackToDefaults` mirroring the
      existing scalar partial-fallback test.

## Acceptance Criteria

- YAML with only `conventional.major` yields a four-key map with defaults for
  `minor`, `patch`, and `none`.
- `go test ./internal/config/...` passes including the new test.
- Existing scalar partial-fallback and CLI-override tests remain green.

## References

- Found during impl review F1: `context/changes/bootstrap-go-cli/reviews/impl-review.md`
- Location: `internal/config/config.go:69-72`
