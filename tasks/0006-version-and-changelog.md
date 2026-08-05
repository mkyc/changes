---
title: "Implement semantic version aggregation and changelog rendering"
id: "0006"
status: pending
priority: high
effort: large
type: feature
tags: ["semver", "changelog", "rendering"]
dependencies: ["0002"]
created_at: "2026-08-04"
---

# Implement semantic version aggregation and changelog rendering

## Objective

Build reusable, side-effect-free release calculation and Keep a Changelog
rendering. Compute a package's next semantic version from its init baseline and
pending changes, then render deterministic Markdown from change-file data.

## Tasks

- [ ] Parse and compare semantic versions from init records.
- [ ] Aggregate increments using `major > minor > patch > none` and calculate
      the next version from the baseline.
- [ ] Map change records into Keep a Changelog sections such as `Added`.
- [ ] Format summaries with optional issue and author suffixes.
- [ ] Render multiline details, release dates, init baseline entries, headings,
      spacing, and preamble deterministically.
- [ ] Expose an in-memory render path that `apply` and `check` can share.
- [ ] Add golden tests for the exact Scenario 01 changelog.

## Acceptance Criteria

- Baseline `1.2.3` plus a minor change produces `1.3.0`.
- The edited widget API record renders as `Add widget API endpoint (#42)
  (@alice)` under `### Added`, followed by its two detail lines.
- The init record renders as the `1.2.3` release dated `2026-06-18`.
- Rendering at `2026-06-19T10:00:00Z` matches the complete Step 4
  `CHANGELOG.md` block byte-for-byte.
- Repeated renders of the same inputs are identical and do not touch the filesystem.
