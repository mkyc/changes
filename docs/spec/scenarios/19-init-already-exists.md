# Scenario 19 — Init when already initialized

Running `init` twice is rejected.

**Scope:** invalid input / state conflict, error handling.

---

## Given

### Repository

- A git repository with default branch `main`.
- Latest release tag `v1.2.3`.
- Working tree clean.

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

## Step 1 — init

**When**

```text
changes init
```

**Then**

- Exit code is `1`.
- Stdout is empty.
- Stderr is exactly:

```text
error: .changes/0000-init.yaml already exists
```

**And**

- `.changes/0000-init.yaml` is unchanged.
- No additional files are created under `.changes/`.

---

## Final filesystem state

Unchanged from Given.
