# Scenario 03 — Propose breaking change as major

Conventional commit with breaking change maps to `increment: major`.

**Scope:** single repo, initialized, one breaking commit since `main`.

## Applies decisions

| ID | How it shows up here |
|----|----------------------|
| [D005](../DECISIONS.md#d005-default-since-ref) | `changes propose` with no `--since` |
| [D008](../DECISIONS.md#d008-propose-granularity) | One file for one commit |

---

## Given

### Repository

- Default branch `main`; latest tag `v2.0.0`.
- Branch `feat/break-api` with one commit not on `main`:

  | Hash (short) | Author date | Message |
  |--------------|-------------|---------|
  | `aaa1111` | `2026-06-17` | `feat(api)!: remove legacy endpoint` |

- Working tree clean; current branch `feat/break-api`.

### Filesystem

`.changes/0000-init.yaml`:

```yaml
id: "0000-init"
package: "."
event: init
increment: none
date: "2026-06-18"
summary: "2.0.0"
details: |
  Start tracking changes with `changes`.
  Last release before adopting the tool was 2.0.0.
```

### Environment

- `changes propose` runs at `2026-06-19T10:00:00Z`.

---

## Step 1 — propose

**When**

```text
changes propose
```

**Then** file `.changes/0001-remove-legacy-endpoint.yaml` is created with exactly:

```yaml
id: "0001-remove-legacy-endpoint"
package: "."
event: change
increment: major
date: "2026-06-17"
summary: "Remove legacy endpoint"
details: |
  feat(api)!: remove legacy endpoint
breaking: true
source:
  - hash: "aaa1111"
    message: "feat(api)!: remove legacy endpoint"
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- `.changes/0000-init.yaml` is unchanged.
- `CHANGELOG.md` is not created or modified.

---

## Final filesystem state

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   └── 0001-remove-legacy-endpoint.yaml
└── (no CHANGELOG.md)
```
