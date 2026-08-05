---
title: "Assert exactly four subcommands in discovery test"
id: "0014"
status: pending
priority: low
type: chore
tags: ["test", "bootstrap-go-cli"]
created_at: "2026-08-05"
depends_on: ["0001"]
---

# Assert exactly four subcommands in discovery test

## Objective

Tighten the command discovery test to match the plan's "exactly four
subcommands" requirement.

## Tasks

- [ ] In `TestNewRootCmd_CommandDiscovery`, add
      `if len(root.Commands()) != 4 { t.Errorf(...) }` after the presence
      loop.

## Acceptance Criteria

- Test fails if an extra subcommand is registered beyond init/propose/apply/check.
- `go test ./internal/cli/...` passes.

## References

- Found during impl review F5: `context/changes/bootstrap-go-cli/reviews/impl-review.md`
- Location: `internal/cli/root_test.go:11-25`
