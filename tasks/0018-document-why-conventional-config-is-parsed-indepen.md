---
title: "Document why conventional config is parsed independently of Viper"
id: "0018"
status: completed
priority: low
type: chore
tags: ["config", "docs"]
created_at: "2026-08-05"
---

# Document why conventional config is parsed independently of Viper

## Description

The config file bytes are parsed twice by two different YAML decoders in
`internal/config/config.go`: once by Viper for scalar keys (`v.ReadConfig`,
line 115), and once directly via `go.yaml.in/yaml/v3` in
`fileConventionalFromYAML` (line 53) for the `conventional` block. This is a
deliberate, sound substitution for the plan's originally suggested
`v.GetStringMapStringSlice` approach — Viper can't distinguish "key absent"
from "key present but empty," which is exactly the distinction the
empty-block validation needs. The duplication isn't currently documented in
a comment, so a future reader could mistake it for an oversight.

Source: impl review F2 —
`context/changes/fix-partial-conventional-config-not-merging-with-defaults/reviews/impl-review.md`

## Tasks

- [x] Add a one-line comment at `fileConventionalFromYAML` (or nearby)
      explaining why it parses raw YAML instead of reading `conventional`
      from `v`: to preserve the null-vs-empty-vs-populated distinction that
      Viper's `GetStringMapStringSlice` collapses.

## Acceptance Criteria

- [x] Comment added; no behavior change.
- [x] `go build ./internal/config/...` and `go test ./internal/config/...`
      pass.
