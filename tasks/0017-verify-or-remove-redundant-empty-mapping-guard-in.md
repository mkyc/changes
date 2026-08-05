---
title: "Verify or remove redundant empty-mapping guard in conventional config parsing"
id: "0017"
status: pending
priority: low
type: chore
tags: ["config", "cleanup"]
created_at: "2026-08-05"
---

# Verify or remove redundant empty-mapping guard in conventional config parsing

## Description

`internal/config/config.go:59-60` has a guard
`doc.Conventional.Kind == yaml.MappingNode && len(doc.Conventional.Content) == 0`
that only fires for the literal `conventional: {}` form. Verified experimentally:
bare `conventional:` (no value) parses as a null scalar node (`Kind=8`), not a
mapping node, so this branch never runs for that input — the actual
empty-block rejection for that case happens later via the `len(fileConv) == 0`
check after `Decode` (line 67-68). No existing test exercises the
`conventional: {}` form this guard is meant to catch.

Source: impl review F1 —
`context/changes/fix-partial-conventional-config-not-merging-with-defaults/reviews/impl-review.md`

## Tasks

- [ ] Add a test case for `conventional: {}\n` (empty mapping) to confirm the
      guard is reachable and correct, OR remove the guard and rely solely on
      the post-decode `len(fileConv) == 0` check, which already covers both
      the null and empty-mapping forms.
- [ ] Run `go test ./internal/config/...` to confirm no regressions.

## Acceptance Criteria

- [ ] Either a passing test exists for `conventional: {}` exercising the
      guard, or the guard is removed and behavior is unchanged (same tests
      pass).
- [ ] `go build ./internal/config/...` and `go test ./internal/config/...`
      pass.
