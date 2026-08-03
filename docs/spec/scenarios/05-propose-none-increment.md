# Scenario 05 — Propose no-release change

A `ci` commit maps to `increment: none` — tracked in changelog later but no version bump.

**Scope:** single repo, initialized, one `ci` commit since `main`.

---

## Given

### Repository

- Default branch `main`; latest tag `v1.2.3`.
- Branch `ci/workflow` with one commit not on `main`:

  | Hash (short) | Author date | Message |
  |--------------|-------------|---------|
  | `ccc3333` | `2026-06-17` | `ci: update GitHub Actions workflow` |

- Working tree clean; current branch `ci/workflow`.

### Filesystem

`.changes/0000-init.yaml`:

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

---

## Step 1 — propose

**When**

```text
changes propose
```

**Then** file `.changes/0001-update-github-actions-workflow.yaml` is created with exactly:

```yaml
id: "0001-update-github-actions-workflow"
package: "."
event: change
increment: none
date: "2026-06-17"
summary: "Update GitHub Actions workflow"
details: |
  ci: update GitHub Actions workflow
breaking: false
source:
  - hash: "ccc3333"
    message: "ci: update GitHub Actions workflow"
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.

---

## Step 2 — apply

**When** the user commits `.changes/0001-update-github-actions-workflow.yaml` unchanged, then runs:

```text
changes apply
```

**Then** file `CHANGELOG.md` is created with exactly:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.

### Changed

- Update GitHub Actions workflow

  ci: update GitHub Actions workflow
```

**And**

- Release version remains `1.2.3` ([D004](../DECISIONS.md#d004-version-computation) — `none` does not bump).
- `.changes/0001-update-github-actions-workflow.yaml` is moved to `.changes/released/0001-update-github-actions-workflow.yaml`.

---

## Final filesystem state

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   └── released/
│       └── 0001-update-github-actions-workflow.yaml
└── CHANGELOG.md
```
