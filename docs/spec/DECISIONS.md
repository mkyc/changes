# Decision log

Recorded decisions for `changes`. Supersedes the open questions in [DESIGN.md](../DESIGN.md).

New entries append at the bottom. When a decision changes, add a new entry rather than editing history.

| ID | Date | Summary |
|----|------|---------|
| [D001](#d001-consumed-changesets) | 2026-07-06 | Move released changesets to `.changes/released/` |
| [D002](#d002-concurrent-pr-sequence-collisions) | 2026-07-06 | Duplicate sequence IDs are a merge-time concern; `check` rejects them |
| [D003](#d003-changelog-format) | 2026-07-06 | Keep a Changelog 1.1.0 only |
| [D004](#d004-version-computation) | 2026-07-06 | Version is computed from `init` + increments |
| [D005](#d005-default-since-ref) | 2026-07-06 | Default `--since` is `main` |
| [D006](#d006-init-without-prior-tag) | 2026-07-06 | Baseline version `0.0.0` when no release tag exists |
| [D007](#d007-apply-side-effects) | 2026-07-06 | `apply` writes only `CHANGELOG.md` |
| [D008](#d008-propose-granularity) | 2026-07-06 | One change file per propose run |
| [D009](#d009-check-success-output) | 2026-07-06 | `check` prints `ok` on success |

---

## D001: Consumed changesets

**Date:** 2026-07-06  
**Status:** accepted

After a release, consumed change files are moved to `.changes/released/` (not deleted, not left in `.changes/` root).

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

Release version is always computed from the `init` baseline plus aggregated `increment` values across pending change files. No config or CLI override for now.

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

## D007: Apply side effects

**Date:** 2026-07-06  
**Status:** accepted

`changes apply` writes only `CHANGELOG.md`. It does not bump `package.json`, tags, or any other files.

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
