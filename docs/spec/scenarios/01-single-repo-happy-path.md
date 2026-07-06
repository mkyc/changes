# Scenario 01 — Single-repo happy path

Adopt `changes` in an existing single-package repository, propose a changeset from conventional commits, edit and commit it, then apply to produce an up-to-date changelog.

**Scope:** single repo (`package: "."`), no monorepo, no concurrent PRs, all dependencies available.

## Applies decisions

| ID | How it shows up here |
|----|----------------------|
| [D001](../DECISIONS.md#d001-consumed-changesets) | Pending changes stay under `.changes/`; `.changes/released/` does not exist yet |
| [D002](../DECISIONS.md#d002-concurrent-pr-sequence-collisions) | Step 5: `check` passes with unique sequences `0000` and `0001` |
| [D003](../DECISIONS.md#d003-changelog-format) | Step 4: exact Keep a Changelog 1.1.0 output |
| [D004](../DECISIONS.md#d004-version-computation) | Step 1: `init` reads `1.2.3` from tag; step 4: next version computed as `1.3.0` |
| [D005](../DECISIONS.md#d005-default-since-ref) | Step 2: `changes propose` with no `--since` (defaults to `main`) |
| [D006](../DECISIONS.md#d006-init-without-prior-tag) | Not exercised — repo has tag `v1.2.3` |
| [D007](../DECISIONS.md#d007-apply-side-effects) | Step 4: only `CHANGELOG.md` is written |
| [D008](../DECISIONS.md#d008-propose-granularity) | Step 2: one file grouping both commits |
| [D009](../DECISIONS.md#d009-check-success-output) | Step 5: exit `0`, stdout `ok` |

---

## Given

### Repository

- A git repository with default branch `main`.
- Latest release tag on `main`: `v1.2.3`.
- Feature branch `feat/widget-api` was created from `main` and contains two commits not on `main`:

  | Hash (short) | Author date | Message |
  |--------------|-------------|---------|
  | `abc1234` | `2026-06-17` | `feat(widget): add API endpoint` |
  | `def5678` | `2026-06-17` | `fix(widget): handle nil input` |

- Working tree is clean; current branch is `feat/widget-api`.
- No `.changes/` directory exists yet.
- No `CHANGELOG.md` exists yet.

### Configuration

- No `.changes/config.yaml` is present; tool defaults apply ([D005](../DECISIONS.md#d005-default-since-ref)):
  - `package: "."`
  - `changelog: CHANGELOG.md`
  - `since: main`

### Environment

- `git` is installed and the repository is readable.
- User has write access to the working tree.
- Command timestamps for deterministic `date` fields:
  - `changes init` runs at `2026-06-18T10:00:00Z`
  - `changes propose`, `changes apply`, and `changes check` run at `2026-06-19T10:00:00Z`

### Initial filesystem

```text
.
├── .git/
└── (no .changes/, no CHANGELOG.md)
```

---

## Step 1 — init

**When**

```text
changes init
```

**Then** file `.changes/0000-init.yaml` is created with exactly:

```yaml
id: "0000-init"
package: "."
event: init
increment: none
date: "2026-06-18"
summary: "1.2.3"
details: |
  Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- `summary` is derived from the latest release tag `v1.2.3` ([D004](../DECISIONS.md#d004-version-computation)).
- No other files are created or modified.
- `CHANGELOG.md` does not exist.
- `.changes/released/` does not exist ([D001](../DECISIONS.md#d001-consumed-changesets)).

**When** the user commits `.changes/0000-init.yaml`.

**Then** filesystem contains:

```text
.
├── .changes/
│   └── 0000-init.yaml
└── (no CHANGELOG.md, no .changes/released/)
```

---

## Step 2 — propose

**When**

```text
changes propose
```

(no `--since`; defaults to `main` per [D005](../DECISIONS.md#d005-default-since-ref))

**Then** file `.changes/0001-add-widget-api.yaml` is created with exactly:

```yaml
id: "0001-add-widget-api"
package: "."
event: change
increment: minor
date: "2026-06-17"
summary: "Add widget API endpoint"
details: |
  feat(widget): add API endpoint

  fix(widget): handle nil input
breaking: false
source:
  - hash: "abc1234"
    message: "feat(widget): add API endpoint"
  - hash: "def5678"
    message: "fix(widget): handle nil input"
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- Exactly one new change file is created ([D008](../DECISIONS.md#d008-propose-granularity)).
- `.changes/0000-init.yaml` is unchanged.
- `CHANGELOG.md` does not exist.

---

## Step 3 — edit and commit

**When** the user replaces the contents of `.changes/0001-add-widget-api.yaml` with:

```yaml
id: "0001-add-widget-api"
package: "."
event: change
increment: minor
date: "2026-06-17"
summary: "Add widget API endpoint"
details: |
  Adds REST endpoint for widget creation.
  Includes nil-input guard on the handler.
breaking: false
issues: ["#42"]
authors: ["@alice"]
source:
  - hash: "abc1234"
    message: "feat(widget): add API endpoint"
  - hash: "def5678"
    message: "fix(widget): handle nil input"
```

**And** commits `.changes/0001-add-widget-api.yaml`.

**Then** filesystem contains:

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   └── 0001-add-widget-api.yaml
└── (no CHANGELOG.md)
```

**And**

- `.changes/0001-add-widget-api.yaml` on disk matches the YAML block above byte-for-byte.

---

## Step 4 — apply

**When**

```text
changes apply
```

**Then** file `CHANGELOG.md` is created with exactly:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add widget API endpoint (#42) (@alice)

  Adds REST endpoint for widget creation.
  Includes nil-input guard on the handler.

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- Only `CHANGELOG.md` is created or modified; no other files are written ([D007](../DECISIONS.md#d007-apply-side-effects)).
- `.changes/0000-init.yaml` is unchanged.
- `.changes/0001-add-widget-api.yaml` is unchanged.
- Computed next release version is `1.3.0` — baseline `1.2.3` from init ([D004](../DECISIONS.md#d004-version-computation)) plus one `minor` increment from `0001`; not written into `CHANGELOG.md` until release.

**When** the user commits `CHANGELOG.md`.

**Then** filesystem contains:

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   └── 0001-add-widget-api.yaml
└── CHANGELOG.md
```

---

## Step 5 — check

**When**

```text
changes check
```

**Then**

- Exit code is `0` ([D009](../DECISIONS.md#d009-check-success-output)).
- Stdout is exactly:

```text
ok
```

- Stderr is empty.

**And** all of the following hold (otherwise `check` would fail):

- Change file schema is valid.
- Sequence prefixes are unique: `0000` and `0001` each appear on exactly one file under `.changes/` ([D002](../DECISIONS.md#d002-concurrent-pr-sequence-collisions)).
- Sequence prefixes have no gaps: `0000`, `0001`.
- `CHANGELOG.md` matches exactly what `changes apply` would write (no drift).
- `CHANGELOG.md` is unchanged on disk after `check`.
- `.changes/0000-init.yaml` is unchanged.
- `.changes/0001-add-widget-api.yaml` is unchanged.

---

## Final filesystem state

After all steps and commits:

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   └── 0001-add-widget-api.yaml
└── CHANGELOG.md
```

With file contents exactly as specified in steps 1–4. No `.changes/released/` directory ([D001](../DECISIONS.md#d001-consumed-changesets)).

---

## Out of scope for this scenario

- Monorepo / multiple packages
- Invalid YAML or schema violations
- Missing git, changelog drift, duplicate sequence collision ([D002](../DECISIONS.md#d002-concurrent-pr-sequence-collisions) failure case)
- `init` with no prior release tag ([D006](../DECISIONS.md#d006-init-without-prior-tag))
- `major`, `patch`-only, or `none`-only increments
- Breaking-change commits (`feat!`, `BREAKING CHANGE` footer)
- Post-release move of consumed changesets to `.changes/released/` ([D001](../DECISIONS.md#d001-consumed-changesets))

These will appear in later scenario files.
