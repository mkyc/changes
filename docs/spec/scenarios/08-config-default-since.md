# Scenario 08 — Config provides default since ref

When `.changes/config.yaml` sets `since`, propose uses it instead of the built-in `main` default.

**Scope:** configuration resolution, [D005](../DECISIONS.md#d005-default-since-ref).

---

## Given

### Repository

- Default branch `main`.
- Tag `v1.0.0` on an older commit.
- Tag `v1.2.3` on `main` (current release).
- Branch `feat/report` created from `v1.2.3` with one commit not reachable from `v1.0.0`:

  | Hash (short) | Author date | Message |
  |--------------|-------------|---------|
  | `ddd4444` | `2026-06-17` | `feat(report): add CSV export` |

- Working tree clean; current branch `feat/report`.

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

`.changes/config.yaml`:

```yaml
changelog: CHANGELOG.md
packages:
  - path: "."
    tag_prefix: "v"
since: v1.0.0
conventional:
  major: [ "feat!", "BREAKING CHANGE" ]
  minor: [ "feat" ]
  patch: [ "fix", "perf", "docs", "style", "refactor", "test", "build", "chore" ]
  none: [ "ci" ]
```

---

## Step 1 — propose

**When**

```text
changes propose
```

(no `--since` flag)

**Then** file `.changes/0001-add-csv-export.yaml` is created with exactly:

```yaml
id: "0001-add-csv-export"
package: "."
event: change
increment: minor
date: "2026-06-17"
summary: "Add CSV export"
details: |
  feat(report): add CSV export
breaking: false
source:
  - hash: "ddd4444"
    message: "feat(report): add CSV export"
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.

---

## Final filesystem state

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   ├── 0001-add-csv-export.yaml
│   └── config.yaml
└── (no CHANGELOG.md)
```
