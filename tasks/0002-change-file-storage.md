---
title: "Implement change file model and storage"
id: "0002"
status: pending
priority: high
effort: large
type: feature
tags: ["changesets", "yaml", "storage"]
dependencies: ["0001"]
created_at: "2026-08-04"
---

# Implement change file model and storage

## Objective

Implement the YAML change-file domain model and filesystem repository that
serve as the source of truth for releases. Support init and change events,
stable sequential IDs, pending files, and released files without losing user
edits or changing serialized content during moves.

## Tasks

- [ ] Define typed fields for IDs, package, event, increment, date, summary,
      details, breaking state, issues, authors, and source commits.
- [ ] Parse and validate `event: init|change` and
      `increment: major|minor|patch|none`, including required fields.
- [ ] Discover pending `.changes/*.yaml` and consumed
      `.changes/released/*.yaml` files while keeping init in the root.
- [ ] Allocate zero-padded sequential IDs across both pending and released
      locations and derive safe summary slugs for filenames.
- [ ] Detect duplicate sequence prefixes across both locations.
- [ ] Write new files deterministically and move released change files without
      altering their bytes.
- [ ] Add unit tests using the init and widget API examples from Scenario 01.

## Acceptance Criteria

- The exact `0000-init.yaml` and `0001-add-widget-api.yaml` documents in
  Scenario 01 parse into the expected typed values and serialize deterministically.
- The next ID after `0000-init.yaml` is `0001`, even when prior change files are
  under `.changes/released/`.
- A duplicate numeric prefix in pending and released storage is reported as a
  validation error.
- Moving `.changes/0001-add-widget-api.yaml` to
  `.changes/released/0001-add-widget-api.yaml` preserves its content byte-for-byte.
- Optional `issues`, `authors`, and `source` fields round-trip without data loss.
