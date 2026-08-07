---
title: "Improve error message for non-mapping conventional config value"
id: "0019"
status: completed
priority: low
type: chore
tags: ["config"]
created_at: "2026-08-05"
---

# Improve error message for non-mapping conventional config value

## Description

In `internal/config/config.go`, if `conventional:` is given a YAML sequence
or scalar instead of a mapping (e.g. `conventional: [a, b]`), the code
doesn't pre-validate the node kind before `doc.Conventional.Decode(&fileConv)`
(line 64); it returns the underlying yaml library error (e.g. "cannot
unmarshal !!seq into map[string][]string") instead of one of this file's own
descriptive messages (compare to the custom "conventional must contain at
least one key" and "unknown conventional key %q" messages). Functionally
correct — no panic, error is returned — but the message is less actionable
for users than the rest of the file's validation errors.

This was out of scope for the original plan (only "unknown key" and "empty
block" validation were specified), so this is optional polish, not a bug fix.

Source: impl review F3 —
`context/changes/fix-partial-conventional-config-not-merging-with-defaults/reviews/impl-review.md`

## Tasks

- [x] Add a check on `doc.Conventional.Kind` before `Decode` and return a
      descriptive error (e.g. "conventional must be a mapping of keys to
      lists") when it isn't a mapping node.
- [x] Add a test asserting the improved error message for a non-mapping
      `conventional:` value.

## Acceptance Criteria

- [x] Non-mapping `conventional:` values produce a clear, actionable error
      message consistent with the file's other validation errors.
- [x] `go test ./internal/config/...` passes, including the new test.
