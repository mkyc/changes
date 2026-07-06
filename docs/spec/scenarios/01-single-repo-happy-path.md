# Scenario 01 — Single-repo happy path

Adopt `changes` in an existing single-package repository, propose a changeset from conventional commits, edit and commit it, then apply to produce an up-to-date changelog.

**Scope:** single repo (`package: "."`), no monorepo, no concurrent PRs, all dependencies available.

**Changelog format (pinned for this scenario):** [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) 1.1.0.

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

- No `.changes/config.yaml` is present; tool defaults apply:
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
- No other files are created or modified.
- `CHANGELOG.md` does not exist.

**When** the user commits `.changes/0000-init.yaml`.

**Then** filesystem contains:

```text
.
├── .changes/
│   └── 0000-init.yaml
└── (no CHANGELOG.md)
```

---

## Step 2 — propose

**When**

```text
changes propose --since main
```

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
- `.changes/0000-init.yaml` is unchanged.
- `.changes/0001-add-widget-api.yaml` is unchanged.
- Implied next release version is `1.3.0` (baseline `1.2.3` + one `minor` increment); it is not written into `CHANGELOG.md` until release.

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

- Exit code is `0`.
- Stdout is exactly:

```text
ok
```

- Stderr is empty.

**And**

- `CHANGELOG.md` is unchanged.
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

With file contents exactly as specified in steps 1–4.

---

## Out of scope for this scenario

- Monorepo / multiple packages
- Invalid YAML or schema violations
- Missing git, dirty apply drift, concurrent ID collision
- `major`, `patch`-only, or `none`-only increments
- Breaking-change commits (`feat!`, `BREAKING CHANGE` footer)
- Post-release archival of consumed changesets

These will appear in later scenario files.
