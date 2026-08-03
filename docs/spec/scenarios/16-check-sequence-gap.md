# Scenario 16 — Check rejects sequence gap

Missing a sequence number between `0000` and the highest pending ID fails validation.

**Scope:** boundary condition, ordering invariant, error handling.

---

## Given

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

`.changes/0002-skip-sequence.yaml`:

```yaml
id: "0002-skip-sequence"
package: "."
event: change
increment: patch
date: "2026-06-17"
summary: "Fix off-by-one error"
details: |
  fix(calc): correct boundary check
breaking: false
source:
  - hash: "abc1234"
    message: "fix(calc): correct boundary check"
```

(no `0001-*.yaml` exists)

---

## Step 1 — check

**When**

```text
changes check
```

**Then**

- Exit code is `1`.
- Stdout is empty.
- Stderr is exactly:

```text
error: missing sequence 0001
```

**And**

- No files are created or modified.

---

## Final filesystem state

Unchanged from Given.
