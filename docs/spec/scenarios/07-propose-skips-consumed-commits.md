# Scenario 07 — Propose skips already-consumed commits

Commits listed in an existing change file's `source` are not proposed again.

**Scope:** single repo, idempotency / partial propose, [D008](../DECISIONS.md#d008-propose-granularity).

---

## Given

### Repository

- Default branch `main`; latest tag `v1.2.3`.
- Branch `feat/widget-api` with two commits not on `main`:

  | Hash (short) | Author date | Message |
  |--------------|-------------|---------|
  | `abc1234` | `2026-06-17` | `feat(widget): add API endpoint` |
  | `def5678` | `2026-06-18` | `fix(widget): handle nil input` |

- Working tree clean; current branch `feat/widget-api`.

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

`.changes/0001-add-widget-api.yaml`:

```yaml
id: "0001-add-widget-api"
package: "."
event: change
increment: minor
date: "2026-06-17"
summary: "Add widget API endpoint"
details: |
  feat(widget): add API endpoint
breaking: false
source:
  - hash: "abc1234"
    message: "feat(widget): add API endpoint"
```

---

## Step 1 — propose

**When**

```text
changes propose
```

**Then** file `.changes/0002-handle-nil-input.yaml` is created with exactly:

```yaml
id: "0002-handle-nil-input"
package: "."
event: change
increment: patch
date: "2026-06-18"
summary: "Handle nil input"
details: |
  fix(widget): handle nil input
breaking: false
source:
  - hash: "def5678"
    message: "fix(widget): handle nil input"
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- `.changes/0000-init.yaml` is unchanged.
- `.changes/0001-add-widget-api.yaml` is unchanged.
- Commit `abc1234` does not appear in any new file's `source`.

---

## Final filesystem state

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   ├── 0001-add-widget-api.yaml
│   └── 0002-handle-nil-input.yaml
└── (no CHANGELOG.md)
```
