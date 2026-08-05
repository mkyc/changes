---
title: "Reject extra positional args on stub subcommands"
id: "0012"
status: in-progress
priority: medium
type: bug
tags: ["cli", "bootstrap-go-cli"]
created_at: "2026-08-05"
depends_on: ["0001"]
---

# Reject extra positional args on stub subcommands

## Objective

Harden stub subcommands so mistyped or extra positional arguments surface as
errors instead of silently succeeding.

## Steps to Reproduce

1. Build the `changes` binary.
2. Run `changes init extra-arg` (or any subcommand with trailing args).

## Expected Behavior

Command exits with a non-zero status and reports that unknown/extra arguments
are not accepted.

## Actual Behavior

Command exits 0 with no stdout or stderr output. Extra args are silently
ignored.

## Tasks

- [ ] Set `Args: cobra.NoArgs` on `init`, `propose`, `apply`, and `check`
      subcommands in `internal/cli/`.
- [ ] Add a test asserting unknown positional args return an error via fake
      deps.

## Acceptance Criteria

- `changes init foo`, `changes propose bar`, etc. exit non-zero with an error
  message on stderr.
- `go test ./internal/cli/...` passes including the new test.
- Existing silent-success smoke tests remain green when no extra args are
  passed.

## References

- Found during impl review F2: `context/changes/bootstrap-go-cli/reviews/impl-review.md`
- Location: `internal/cli/init.go`, `propose.go`, `apply.go`, `check.go`
