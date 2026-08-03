# Scenario 10 — Version aggregation across pending changes

Multiple pending increments aggregate with precedence: major > minor > patch > none.

**Scope:** [D004](../DECISIONS.md#d004-version-computation), data consistency.

---

## Given

### Repository

- Initialized single repo; no git operations required beyond filesystem state.

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

`.changes/0001-fix-typo.yaml`:

```yaml
id: "0001-fix-typo"
package: "."
event: change
increment: patch
date: "2026-06-17"
summary: "Fix typo in error message"
details: |
  fix(ui): correct button label
breaking: false
source:
  - hash: "111aaaa"
    message: "fix(ui): correct button label"
```

`.changes/0002-update-ci.yaml`:

```yaml
id: "0002-update-ci"
package: "."
event: change
increment: none
date: "2026-06-17"
summary: "Update CI cache key"
details: |
  ci: bump cache key
breaking: false
source:
  - hash: "222bbbb"
    message: "ci: bump cache key"
```

`.changes/0003-add-export.yaml`:

```yaml
id: "0003-add-export"
package: "."
event: change
increment: minor
date: "2026-06-18"
summary: "Add export endpoint"
details: |
  feat(export): add CSV download
breaking: false
source:
  - hash: "333cccc"
    message: "feat(export): add CSV download"
```

---

## Step 1 — apply

**When**

```text
changes apply
```

**Then** file `CHANGELOG.md` is created with exactly:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2026-06-19

### Added

- Add export endpoint

  feat(export): add CSV download

### Fixed

- Fix typo in error message

  fix(ui): correct button label

### Changed

- Update CI cache key

  ci: bump cache key

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- Release version written to changelog is `1.3.0` (`minor` wins over `patch` and `none`).
- `.changes/0001-fix-typo.yaml`, `.changes/0002-update-ci.yaml`, and `.changes/0003-add-export.yaml` are moved to `.changes/released/` with identical content.

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

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   └── released/
│       ├── 0001-fix-typo.yaml
│       ├── 0002-update-ci.yaml
│       └── 0003-add-export.yaml
└── CHANGELOG.md
```
