---
title: "Normalize nested conventional decode errors to one line"
id: "0021"
status: completed
priority: medium
type: bug
tags: ["config", "cli"]
created_at: "2026-08-08"
depends_on: ["0019"]
---

# Normalize nested conventional decode errors to one line

## Objective

When a `conventional` mapping has a value of the wrong YAML type (e.g.
`major: feat` instead of a list), return a single-line, actionable error
instead of leaking the multi-line yaml-library unmarshal message — so
`cli.Run` can keep the D011 stderr contract (`error: <message>` on one line).

## Steps to Reproduce

1. Write:

   ```yaml
   conventional:
     major: feat
   ```

2. Load config and print `err.Error()`, or run a stub subcommand that loads
   that config through `cli.Run`.

## Expected Behavior

A single-line descriptive error (e.g. `conventional.major must be a list of
strings`, or a similar one-line wrap). Via `cli.Run`, stderr is exactly one
line matching `error: …\n`.

## Actual Behavior

Error is multi-line, e.g.:

```
yaml: unmarshal errors:
  line 2: cannot unmarshal !!str `feat` into []string
```

Through `cli.Run` this becomes multi-line stderr and violates D011.

## Tasks

- [#] In `fileConventionalFromYAML`, wrap or replace `Decode` failures with a
      single-line message (prefer naming the offending key when practical).
- [#] Add a unit test for a nested wrong-type value asserting a one-line
      error with no embedded newlines.
- [#] Optionally add a `cli.Run` case with memfs invalid conventional config
      asserting exit `1`, empty stdout, and single-line `error: …` stderr.

## Acceptance Criteria

- [#] Nested conventional type errors never embed newlines in `err.Error()`.
- [#] Top-level non-mapping / empty-block messages from 0017/0019 remain
      unchanged.
- [#] `go test ./internal/config/...` (and `./internal/cli/...` if a Run case
      is added) passes.

## References

- Branch review finding F2 (tasks 0011–0019 on `fixes`)
- Related: task 0019 improved top-level shape only
- Location: `internal/config/config.go` (~line 71); surface via
  `internal/cli/run.go`
