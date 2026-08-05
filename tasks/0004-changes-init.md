---
title: "Implement changes init for an existing repository"
id: "0004"
status: pending
priority: high
effort: medium
type: feature
tags: ["init", "git", "changesets"]
dependencies: ["0002", "0003"]
created_at: "2026-08-04"
---

# Implement changes init for an existing repository

## Objective

Implement `changes init` for adopting the tool in an existing single-package
repository with a prior semantic release tag. Record that tag as the immutable
`0000` baseline without creating a changelog or released directory.

## Tasks

- [ ] Resolve the latest semantic release tag using the configured tag prefix.
- [ ] Create `.changes/` only when the initialization record can be written.
- [ ] Build the `event: init`, `increment: none` record with the command date,
      baseline version, and adoption details.
- [ ] Refuse to overwrite or create a second init record.
- [ ] Keep stdout and stderr empty on success.
- [ ] Add command tests for the tagged single-repository path in Scenario 01.

## Acceptance Criteria

- At a fixed time of `2026-06-18T10:00:00Z` in a repository tagged `v1.2.3`,
  `changes init` writes the exact Scenario 01 `.changes/0000-init.yaml` bytes.
- The command exits 0 with empty stdout and stderr.
- No `CHANGELOG.md`, `.changes/released/`, or other file is created or modified.
- The summary is `1.2.3`, while the explanatory details retain the same version.
- A Git or write failure leaves no partial init file.
