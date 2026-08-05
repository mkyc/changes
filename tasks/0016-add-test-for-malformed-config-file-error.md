---
title: "Add test for malformed config file error"
id: "0016"
status: pending
priority: low
type: chore
tags: ["test", "config", "bootstrap-go-cli"]
created_at: "2026-08-05"
depends_on: ["0001"]
---

# Add test for malformed config file error

## Objective

Cover the error boundary when `.changes/config.yaml` contains invalid YAML.

## Tasks

- [ ] Add `TestLoad_InvalidConfigReturnsError` in `internal/config/config_test.go`.
- [ ] Write invalid YAML to the in-memory filesystem and assert `Load` returns
      a non-nil error.

## Acceptance Criteria

- Test confirms malformed config does not silently succeed.
- `go test ./internal/config/...` passes.

## References

- Found during impl review F7: `context/changes/bootstrap-go-cli/reviews/impl-review.md`
- Location: `internal/config/config_test.go`
