# Scenario 13 — Check rejects duplicate sequence

Two change files sharing the same numeric prefix fail validation.

**Scope:** concurrent/conflicting operations, [D002](../DECISIONS.md#d002-concurrent-pr-sequence-collisions), error handling.

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

`.changes/0005-add-widget-api.yaml`:

```yaml
id: "0005-add-widget-api"
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

`.changes/0005-fix-widget-cache.yaml`:

```yaml
id: "0005-fix-widget-cache"
package: "."
event: change
increment: patch
date: "2026-06-18"
summary: "Fix widget cache invalidation"
details: |
  fix(widget): invalidate cache on update
breaking: false
source:
  - hash: "def5678"
    message: "fix(widget): invalidate cache on update"
```

`CHANGELOG.md` (in sync with `.changes/` if duplicate did not exist — present but irrelevant to this failure):

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
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
error: duplicate sequence 0005
```

**And**

- No files are created or modified.

---

## Final filesystem state

Unchanged from Given.
