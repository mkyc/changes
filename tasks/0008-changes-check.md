---
title: "Implement changes check validation and drift detection"
id: "0008"
status: pending
priority: high
effort: large
type: feature
tags: ["check", "validation", "changelog"]
dependencies: ["0002", "0007"]
created_at: "2026-08-04"
---

# Implement changes check validation and drift detection

## Objective

Implement the read-only `changes check` command for CI. Validate every pending,
init, and released record, enforce sequence integrity, and prove that the
current changelog matches the canonical result without mutating the repository.

## Tasks

- [ ] Validate change-file YAML structure, required fields, enums, and package scope.
- [ ] Validate unique sequence prefixes across `.changes/` and
      `.changes/released/` while treating the root init record correctly.
- [ ] Detect unexpected pending changes after an applied release.
- [ ] Render the expected changelog in memory using the same release engine as
      `apply` and compare it byte-for-byte with the file on disk.
- [ ] Report actionable non-zero failures for schema, sequence, and drift errors.
- [ ] Emit only `ok` followed by a newline on successful validation.
- [ ] Add tests proving the command never changes checked files.

## Acceptance Criteria

- On Scenario 01's final filesystem, `changes check` exits 0, writes exactly
  `ok\n` to stdout, and writes nothing to stderr.
- `0000` in the root and `0001` only in `released/` pass sequence validation.
- Schema errors, duplicate sequence prefixes, unexpected pending records, and
  changelog drift each cause a non-zero exit.
- `CHANGELOG.md`, `.changes/0000-init.yaml`, and the released widget change have
  identical bytes before and after the command.
- The check path creates, removes, or renames no filesystem entries.
