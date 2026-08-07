---
title: "Add automated D011 failure-path test"
id: "0015"
status: completed
priority: low
type: chore
tags: ["test", "bootstrap-go-cli"]
created_at: "2026-08-05"
depends_on: ["0001"]
---

# Add automated D011 failure-path test

## Objective

Automate verification of the D011 error contract (stderr `error: <message>`,
exit failure, stdout untouched) that is currently manual-only (Progress 4.7).

## Tasks

- [x] Add an in-process test executing an unknown subcommand via fake deps.
- [x] Assert stderr matches `^error: ` and stdout is empty.
- [x] Optionally extract `main`'s error handler for direct unit testing if
      needed to avoid `os.Exit` in tests.

## Acceptance Criteria

- `go test ./internal/cli/...` includes a failure-path case for D011.
- Test passes without requiring manual binary invocation.

## References

- Found during impl review F6: `context/changes/bootstrap-go-cli/reviews/impl-review.md`
- D011: `docs/spec/DECISIONS.md`
- Location: `internal/cli/root_test.go`
