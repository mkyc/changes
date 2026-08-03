# Scenario 14 — Check rejects changelog drift

Committed `CHANGELOG.md` that differs from `apply` output fails validation.

**Scope:** data consistency, error handling, recovery.

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

`.changes/0001-add-widget-api.yaml`:

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

`CHANGELOG.md` (stale — missing the pending change):

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

### Environment

- `changes apply` runs at `2026-06-19T10:00:00Z`.

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
error: CHANGELOG.md is out of date; run changes apply
```

**And**

- No files are created or modified.

---

## Step 2 — recovery via apply

**When**

```text
changes apply
```

**Then** file `CHANGELOG.md` is replaced with exactly:

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

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- `.changes/0001-add-widget-api.yaml` is moved to `.changes/released/0001-add-widget-api.yaml`.

---

## Step 3 — check after recovery

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

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   └── released/
│       └── 0001-add-widget-api.yaml
└── CHANGELOG.md
```

With `CHANGELOG.md` content from step 2.
