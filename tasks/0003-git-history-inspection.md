---
title: "Implement Git history and release tag inspection"
id: "0003"
status: pending
priority: high
effort: medium
type: feature
tags: ["git", "repository"]
dependencies: ["0001"]
created_at: "2026-08-04"
---

# Implement Git history and release tag inspection

## Objective

Provide a small Git adapter that can inspect an existing repository for the
latest semantic release tag and for commits introduced since a base ref. Keep
ordering, short hashes, messages, and author dates deterministic for downstream
`init` and `propose` commands.

## Tasks

- [ ] Verify that the working directory is inside a readable Git repository.
- [ ] Resolve the latest reachable tag matching the configured prefix and parse
      its semantic version without the prefix.
- [ ] Read commits on the current branch that are not reachable from a supplied
      base ref, in stable oldest-to-newest order.
- [ ] Capture each commit's short hash, full conventional-commit message, and
      author date.
- [ ] Return typed, actionable errors for missing Git, invalid refs, malformed
      release tags, and command failures.
- [ ] Add integration tests around temporary Git repositories and tags.

## Acceptance Criteria

- A repository tagged `v1.2.3` resolves baseline version `1.2.3`.
- Against base ref `main`, the Scenario 01 fixture returns `abc1234` followed by
  `def5678`, including their exact messages and the date `2026-06-17`.
- Commits already reachable from the base ref are excluded.
- Git command failures are propagated without writing partial change files.
- Tests do not depend on the developer's global Git configuration.
