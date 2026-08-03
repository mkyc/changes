# Scenario 12 — Apply is idempotent

Running `apply` twice with no new pending changes produces identical output.

**Scope:** repeated execution / idempotency, [D007](../DECISIONS.md#d007-apply-behavior).

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

`.changes/released/0001-add-widget-api.yaml`:

```yaml
id: "0001-add-widget-api"
package: "."
event: change
increment: minor
date: "2026-06-17"
summary: "Add widget API endpoint"
details: |
  Adds REST endpoint for widget creation.
breaking: false
source:
  - hash: "abc1234"
    message: "feat(widget): add API endpoint"
```

`CHANGELOG.md`:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2026-06-19

### Added

- Add widget API endpoint

  Adds REST endpoint for widget creation.

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

(no pending change files under `.changes/` root)

---

## Step 1 — apply (second run)

**When**

```text
changes apply
```

**Then**

- `CHANGELOG.md` is byte-for-byte identical to the Given content.
- `.changes/0000-init.yaml` is unchanged.
- `.changes/released/0001-add-widget-api.yaml` is unchanged.

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.

---

## Step 2 — check

**When**

```text
changes check
```

**Then**

- Exit code is `0`.
- Stdout is exactly:

```text
ok
```

- Stderr is empty.

---

## Final filesystem state

Unchanged from Given.
