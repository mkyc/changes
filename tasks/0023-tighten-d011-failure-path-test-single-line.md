---
title: "Tighten D011 failure-path test for single-line stderr"
id: "0023"
status: in-progress
priority: low
type: chore
tags: ["cli", "testing"]
created_at: "2026-08-08"
depends_on: ["0015"]
---

# Tighten D011 failure-path test for single-line stderr

## Description

`TestRun_UnknownSubcommandFollowsD011Contract` only asserts stderr matches
`^error: `. That passes even when the message is empty after the prefix or
when stderr contains embedded newlines — so multi-line config decode errors
(finding F2 / task 0021) would not fail this test.

Tighten the assertion to the full D011 shape: non-empty message, exactly one
trailing newline, no embedded newlines.

## Tasks

- [ ] Update `TestRun_UnknownSubcommandFollowsD011Contract` in
      `internal/cli/root_test.go` to assert something equivalent to
      `^error: .+\n$` and `strings.Count(stderr, "\n") == 1` (or a clear
      equivalent).
- [ ] Keep the existing exit-code `1` and empty-stdout checks.

## Acceptance Criteria

- [ ] The D011 test fails if stderr has an empty message after `error: ` or
      contains more than one newline.
- [ ] Unknown-subcommand path still passes under the tightened contract.
- [ ] `go test ./internal/cli/...` passes.

## References

- Branch review finding F4 (tasks 0011–0019 on `fixes`)
- Related: task 0015 (`cli.Run` + D011 failure-path test); pairs with 0021
- Location: `internal/cli/root_test.go` (`TestRun_UnknownSubcommandFollowsD011Contract`)
