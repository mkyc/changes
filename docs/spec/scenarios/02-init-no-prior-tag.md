# Scenario 02 — Init without prior release tag

First-time adoption in a repository that has never been tagged, then propose and apply the first change.

**Scope:** single repo, no release tags, [D006](../DECISIONS.md#d006-init-without-prior-tag). Pre-existing commit history on `main` is out of scope — see [open question #1](../README.md#open-questions).

## Applies decisions

| ID | How it shows up here |
|----|----------------------|
| [D004](../DECISIONS.md#d004-version-computation) | Baseline `0.0.0`; first change bumps to `0.1.0` |
| [D005](../DECISIONS.md#d005-default-since-ref) | Step 2: `changes propose` with no `--since` |
| [D006](../DECISIONS.md#d006-init-without-prior-tag) | `summary: "0.0.0"` when no tag exists |
| [D007](../DECISIONS.md#d007-apply-behavior) | Step 3: writes `CHANGELOG.md` with `## [0.1.0]` and moves changeset |
| [D008](../DECISIONS.md#d008-propose-granularity) | Step 2: one file for one commit since `main` |

Each step documents **one** `changes` command. Git commits and branch checkouts are ordinary repo setup; they are not part of what the tool validates.

---

## Given

### Repository

- A git repository with default branch `main`.
- No release tags exist.
- `main` contains prior history (two commits below); this scenario does **not** define how `propose` would treat thousands of such commits — see [open question #1](../README.md#open-questions).
- Branch `feat/logging` was created from `main` and contains one commit not on `main`:

  | Hash (short) | Author date | Message |
  |--------------|-------------|---------|
  | `aaa1111` | `2026-06-17` | `feat(log): add structured logging` |

- Commits already on `main` (not proposed in this scenario):

  | Hash (short) | Author date | Message |
  |--------------|-------------|---------|
  | `0000001` | `2026-01-10` | `chore: initial commit` |
  | `0000002` | `2026-03-05` | `docs: add README` |

- No `.changes/` directory exists yet.

### Environment

- `changes init` runs at `2026-06-18T10:00:00Z`.
- `changes propose` runs at `2026-06-19T10:00:00Z`.
- `changes apply` runs at `2026-06-19T10:01:00Z`.

---

## Step 1 — init

### Precondition

- Current branch: `main`.
- Working tree clean.
- No `.changes/` directory.

### When

```text
changes init
```

### Then

File `.changes/0000-init.yaml` is created with exactly:

```yaml
id: "0000-init"
package: "."
event: init
increment: none
date: "2026-06-18"
summary: "0.0.0"
details: |
  Start tracking changes with `changes`.
  Last release before adopting the tool was 0.0.0.
```

### And

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- No other files are created or modified.

---

## Step 2 — propose

### Precondition

- Current branch: `feat/logging`.
- `.changes/0000-init.yaml` exists with the content from step 1.
- Working tree clean except for what step 1 created (or init output is already on this branch).
- No `.changes/0001-*.yaml` yet.

### When

```text
changes propose
```

### Then

File `.changes/0001-add-structured-logging.yaml` is created with exactly:

```yaml
id: "0001-add-structured-logging"
package: "."
event: change
increment: minor
date: "2026-06-17"
summary: "Add structured logging"
details: |
  feat(log): add structured logging
breaking: false
source:
  - hash: "aaa1111"
    message: "feat(log): add structured logging"
```

### And

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- `.changes/0000-init.yaml` is unchanged.
- Commits `0000001` and `0000002` on `main` are not in `source`.
- `CHANGELOG.md` does not exist.

---

## Step 3 — apply

### Precondition

- `.changes/0000-init.yaml` and `.changes/0001-add-structured-logging.yaml` exist with the content from steps 1 and 2.
- No `CHANGELOG.md` yet.

### When

```text
changes apply
```

### Then

File `CHANGELOG.md` is created with exactly:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-06-19

### Added

- Add structured logging

  feat(log): add structured logging

## [0.0.0] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 0.0.0.
```

### And

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- Release version written to changelog is `0.1.0` ([D004](../DECISIONS.md#d004-version-computation)).
- `.changes/0001-add-structured-logging.yaml` is moved to `.changes/released/0001-add-structured-logging.yaml` with identical content ([D001](../DECISIONS.md#d001-consumed-changesets), [D007](../DECISIONS.md#d007-apply-behavior)).

---

## Step 4 — check

### Precondition

- Filesystem matches the end of step 3.

### When

```text
changes check
```

### Then

- Exit code is `0`.
- Stdout is exactly:

```text
ok
```

- Stderr is empty.

---

## Final filesystem state

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   └── released/
│       └── 0001-add-structured-logging.yaml
└── CHANGELOG.md
```

With `CHANGELOG.md` content from step 3.

---

## Out of scope for this scenario

- Backfilling or proposing change files for large pre-existing history on `main` at adoption time — [open question #1](../README.md#open-questions).
- Monorepo, invalid input.
