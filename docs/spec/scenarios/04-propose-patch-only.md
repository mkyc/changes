# Scenario 04 — Propose patch-only change

A single `fix` commit maps to `increment: patch`.

**Scope:** single repo, initialized, one fix commit since `main`.

---

## Given

### Repository

- Default branch `main`; latest tag `v1.2.3`.
- Branch `fix/token-expiry` with one commit not on `main`:

  | Hash (short) | Author date | Message |
  |--------------|-------------|---------|
  | `bbb2222` | `2026-06-17` | `fix(auth): correct token expiry` |

- Working tree clean; current branch `fix/token-expiry`.

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

**Then** file `.changes/0001-correct-token-expiry.yaml` is created with exactly:

```yaml
id: "0001-correct-token-expiry"
package: "."
event: change
increment: patch
date: "2026-06-17"
summary: "Correct token expiry"
details: |
  fix(auth): correct token expiry
breaking: false
source:
  - hash: "bbb2222"
    message: "fix(auth): correct token expiry"
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
│   └── 0001-correct-token-expiry.yaml
└── (no CHANGELOG.md)
```
