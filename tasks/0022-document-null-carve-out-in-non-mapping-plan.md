---
title: "Document !!null carve-out in non-mapping conventional plan"
id: "0022"
status: pending
priority: low
type: chore
tags: ["config", "docs"]
created_at: "2026-08-08"
depends_on: ["0019"]
---

# Document !!null carve-out in non-mapping conventional plan

## Description

Task 0019's implementation correctly excludes `Tag == "!!null"` from the
non-mapping kind guard so bare `conventional:` still fails with
`conventional must contain at least one key` instead of
`conventional must be a mapping of keys to lists`. The commit message
explains this; the plan's Contract / "What We're NOT Doing" do not.

Update the plan so future reviews don't treat the carve-out as accidental
drift.

## Tasks

- [ ] Amend
      `context/changes/improve-error-message-for-non-mapping-conventional/plan.md`
      Contract (and NOT Doing if needed) to document that `!!null` /
      bare `conventional:` is intentionally excluded from the non-mapping
      error and still uses the empty-key validation path.
- [ ] No code changes required.

## Acceptance Criteria

- [ ] Plan text matches the implemented `Kind != MappingNode && Tag != "!!null"`
      behavior and states why.
- [ ] No behavioral change; `go test ./internal/config/...` still green if run.

## References

- Branch review finding F3 (tasks 0011–0019 on `fixes`)
- Implementation: `921cbd3`
- Location: plan under
  `context/changes/improve-error-message-for-non-mapping-conventional/`
