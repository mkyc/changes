# Scenario 15 — Check rejects invalid change file

Schema validation failure on a malformed change file.

**Scope:** invalid input, error handling.

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

`.changes/0001-missing-summary.yaml`:

```yaml
id: "0001-missing-summary"
package: "."
event: change
increment: minor
date: "2026-06-17"
details: |
  feat(widget): add API endpoint
breaking: false
source:
  - hash: "abc1234"
    message: "feat(widget): add API endpoint"
```

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
error: .changes/0001-missing-summary.yaml: missing required field "summary"
```

**And**

- No files are created or modified.

---

## Step 2 — apply

**When**

```text
changes apply
```

**Then**

- Exit code is `1`.
- Stdout is empty.
- Stderr is exactly:

```text
error: .changes/0001-missing-summary.yaml: missing required field "summary"
```

**And**

- `CHANGELOG.md` is not created or modified.

---

## Final filesystem state

Unchanged from Given.
