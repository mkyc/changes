# Scenario 02 — Init without prior release tag

First-time adoption in a repository that has never been tagged.

**Scope:** single repo, no release tags, [D006](../DECISIONS.md#d006-init-without-prior-tag).

## Applies decisions

| ID | How it shows up here |
|----|----------------------|
| [D004](../DECISIONS.md#d004-version-computation) | Baseline version becomes starting point for future computation |
| [D006](../DECISIONS.md#d006-init-without-prior-tag) | `summary: "0.0.0"` when no tag exists |
| [D007](../DECISIONS.md#d007-apply-side-effects) | Step 2: only `CHANGELOG.md` is written |

---

## Given

### Repository

- A git repository with default branch `main`.
- No release tags exist.
- Working tree is clean; current branch is `main`.
- No `.changes/` directory exists yet.

### Environment

- `changes init` runs at `2026-06-18T10:00:00Z`.
- `changes apply` runs at `2026-06-18T10:01:00Z`.

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
summary: "0.0.0"
details: |
  Start tracking changes with `changes`.
  Last release before adopting the tool was 0.0.0.
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.

---

## Step 2 — apply

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

## [0.0.0] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 0.0.0.
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- Computed next release version is `0.0.0` (no pending `change` entries).
- Only `CHANGELOG.md` is created or modified.

---

## Final filesystem state

```text
.
├── .changes/
│   └── 0000-init.yaml
└── CHANGELOG.md
```
