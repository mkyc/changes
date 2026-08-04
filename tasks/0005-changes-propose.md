---
title: "Implement conventional commit parsing and changes propose"
id: "0005"
status: pending
priority: high
effort: large
type: feature
tags: ["propose", "git", "conventional-commits"]
dependencies: ["0002", "0003"]
created_at: "2026-08-04"
---

# Implement conventional commit parsing and changes propose

## Objective

Implement `changes propose` for a single repository by reading commits since a
base ref, classifying conventional commits, and emitting one editable pending
change file with the highest required semantic increment.

## Tasks

- [ ] Parse conventional commit types, breaking markers, scopes, subjects, and
      `BREAKING CHANGE` footers.
- [ ] Map commit prefixes to the configured major, minor, patch, and none
      increments and aggregate them with `major > minor > patch > none`.
- [ ] Resolve the base ref from `--since`, then configuration, then the default
      `main`.
- [ ] Group all qualifying commits into one proposal for package `.`.
- [ ] Derive the proposal date from commit author dates, a human-readable
      summary and slug from the lead change, and raw commit messages for details.
- [ ] Allocate the next change-file ID and include ordered source hashes and
      messages.
- [ ] Write no incidental stdout or stderr on success and add focused command tests.

## Acceptance Criteria

- With Scenario 01's two commits and no config or `--since`, the command reads
  from `main` and creates exactly one `.changes/0001-add-widget-api.yaml`.
- The generated file matches the Step 2 YAML exactly, including `increment:
  minor`, date `2026-06-17`, ordered details, and both source entries.
- The `feat` commit raises the aggregated increment to minor over the `fix`
  commit's patch suggestion.
- `.changes/0000-init.yaml` is unchanged and no changelog is created.
- Success exits 0 with empty stdout and stderr.
