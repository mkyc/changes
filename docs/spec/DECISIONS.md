# Decision log

Recorded decisions for `changes`. Supersedes the open questions in [DESIGN.md](../DESIGN.md).

New entries append at the bottom. When a decision changes, add a new entry rather than editing history.

| ID | Date | Summary |
|----|------|---------|
| [D001](#d001-consumed-changesets) | 2026-07-06 | `apply` moves consumed changesets to `.changes/released/` |
| [D002](#d002-concurrent-pr-sequence-collisions) | 2026-07-06 | Duplicate sequence IDs are a merge-time concern; `check` rejects them |
| [D003](#d003-changelog-format) | 2026-07-06 | Keep a Changelog 1.1.0 only |
| [D004](#d004-version-computation) | 2026-07-06 | Version is computed from published baseline + pending increments |
| [D005](#d005-default-since-ref) | 2026-07-06 | Default `--since` is `main` |
| [D006](#d006-init-without-prior-tag) | 2026-07-06 | Baseline version `0.0.0` when no release tag exists |
| [D007](#d007-apply-behavior) | 2026-07-06 | `apply` writes `CHANGELOG.md` and moves pending files to `released/` |
| [D008](#d008-propose-granularity) | 2026-07-06 | One change file per propose run |
| [D009](#d009-check-success-output) | 2026-07-06 | `check` prints `ok` on success |
| [D010](#d010-cli-overrides-config) | 2026-07-06 | CLI flags override `.changes/config.yaml` |
| [D011](#d011-error-output-format) | 2026-07-06 | Errors on stderr as `error: …`, exit code `1` |
| [D012](#d012-sequence-gap-validation) | 2026-07-06 | `check` fails on missing sequence numbers |

---

## D001: Consumed changesets

**Date:** 2026-07-06  
**Status:** accepted

When `changes apply` runs, each pending `event: change` file under `.changes/` is moved to `.changes/released/` (same filename, same content). Files are not deleted. `0000-init.yaml` and `config.yaml` stay in `.changes/` root.

---

## D002: Concurrent PR sequence collisions

**Date:** 2026-07-06  
**Status:** accepted

Two PRs both adding `0005-*.yaml` is a merge conflict the team resolves at merge time — the tool does not allocate or reserve IDs across branches.

`changes check` must fail when more than one change file shares the same numeric sequence prefix (e.g. both `0005-foo.yaml` and `0005-bar.yaml` exist under `.changes/` or `.changes/released/`).

---

## D003: Changelog format

**Date:** 2026-07-06  
**Status:** accepted

Use [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) 1.1.0 only. No configurable or alternate templates for now.

---

## D004: Version computation

**Date:** 2026-07-06  
**Status:** accepted

When `changes apply` runs with pending change files, the release version is computed from the **current published version** plus aggregated `increment` values across those pending files. No config or CLI override for now.

**Current published version** is:
- the `summary` from `0000-init.yaml` when `.changes/released/` is empty, or
- the semver of the latest version section already in `CHANGELOG.md` from a prior apply.

`apply` writes pending entries under `## [X.Y.Z] - <date>` using the computed version. There is no `[Unreleased]` section.

If all pending increments are `none`, the version number does not bump; new entries are appended under the current version section in `CHANGELOG.md`.

---

## D005: Default `--since` ref

**Date:** 2026-07-06  
**Status:** accepted

When `--since` is not passed and `.changes/config.yaml` omits `since`, default to `main`.

---

## D006: Init without prior tag

**Date:** 2026-07-06  
**Status:** accepted

When `changes init` runs in a repo with no release tags, `0000-init.yaml` uses `summary: "0.0.0"`.

---

## D007: Apply behavior

**Date:** 2026-07-06  
**Status:** accepted

`changes apply`:

1. Regenerates `CHANGELOG.md` from `0000-init.yaml`, `.changes/released/`, and any pending files under `.changes/`.
2. Writes computed version heading(s) per [D004](#d004-version-computation).
3. Moves every pending `event: change` file from `.changes/` to `.changes/released/`.

Does not create git tags or modify files outside `.changes/` and `CHANGELOG.md`.

When there are no pending change files, `apply` only regenerates `CHANGELOG.md` from existing state (idempotent if nothing changed).

---

## D008: Propose granularity

**Date:** 2026-07-06  
**Status:** accepted

`changes propose` emits one change file per run, grouping all commits since `--since` that are not already referenced in an existing change file's `source`.

---

## D009: Check success output

**Date:** 2026-07-06  
**Status:** accepted

On success, `changes check` exits `0`, writes exactly `ok` to stdout (plus trailing newline), and leaves stderr empty.

---

## D010: CLI overrides config

**Date:** 2026-07-06  
**Status:** accepted

When the same setting is provided on the CLI and in `.changes/config.yaml`, the CLI value wins. Config values win over built-in defaults.

Precedence: **CLI flags > `.changes/config.yaml` > built-in defaults**.

---

## D011: Error output format

**Date:** 2026-07-06  
**Status:** accepted

On failure, commands exit `1`, write a single line to stderr in the form `error: <message>`, and leave stdout empty unless a scenario specifies otherwise.

---

## D012: Sequence gap validation

**Date:** 2026-07-06  
**Status:** accepted

Sequence prefixes must be unique across `.changes/` and `.changes/released/` combined. Pending files under `.changes/` (excluding init) must form a contiguous tail continuing from the highest sequence in `.changes/released/` (or start at `0001` when `released/` is empty). `changes check` fails with `error: missing sequence NNNN` when a gap exists.
