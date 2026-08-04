---
title: "Implement changes apply and archive released changes"
id: "0007"
status: pending
priority: high
effort: large
type: feature
tags: ["apply", "changelog", "changesets"]
dependencies: ["0002", "0006"]
created_at: "2026-08-04"
---

# Implement changes apply and archive released changes

## Objective

Implement `changes apply` as the release mutation: validate pending change
files, compute the next version, write the canonical changelog, and consume
pending changes by moving them into `.changes/released/` while preserving the
init baseline.

## Tasks

- [ ] Load and validate the init record, released history, and pending changes
      for package `.`.
- [ ] Compute the release version and render the complete changelog with the
      command date.
- [ ] Write `CHANGELOG.md` safely so a render or write failure does not leave a
      truncated file.
- [ ] Create `.changes/released/` and move only successfully applied change
      records into it.
- [ ] Preserve `.changes/0000-init.yaml` and each moved file's exact bytes.
- [ ] Return successful CLI status without incidental output.
- [ ] Add command tests for the Scenario 01 apply transition.

## Acceptance Criteria

- Applying the edited Scenario 01 change at `2026-06-19T10:00:00Z` writes the
  exact Step 4 `CHANGELOG.md` with release version `1.3.0`.
- `.changes/0001-add-widget-api.yaml` moves to
  `.changes/released/0001-add-widget-api.yaml` byte-for-byte.
- `.changes/0000-init.yaml` remains unchanged in the root and no pending change
  file remains there.
- The command exits 0 with empty stdout and stderr.
- A validation or changelog-write failure does not archive a pending change.
