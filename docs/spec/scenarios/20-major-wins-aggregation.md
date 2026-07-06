# Scenario 20 — Major increment wins in aggregation

When pending changes include both `minor` and `major`, computed release version uses a major bump.

**Scope:** [D004](../DECISIONS.md#d004-version-computation), boundary condition for increment precedence.

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

`.changes/0001-add-export.yaml`:

```yaml
id: "0001-add-export"
package: "."
event: change
increment: minor
date: "2026-06-17"
summary: "Add export endpoint"
details: |
  feat(export): add CSV download
breaking: false
source:
  - hash: "111aaaa"
    message: "feat(export): add CSV download"
```

`.changes/0002-remove-legacy-api.yaml`:

```yaml
id: "0002-remove-legacy-api"
package: "."
event: change
increment: major
date: "2026-06-18"
summary: "Remove legacy API"
details: |
  feat(api)!: remove legacy endpoint
breaking: true
source:
  - hash: "222bbbb"
    message: "feat(api)!: remove legacy endpoint"
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

## [Unreleased]

### Added

- Add export endpoint

  feat(export): add CSV download

### Removed

- Remove legacy API

  feat(api)!: remove legacy endpoint

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- Computed next release version is `2.0.0` (`major` wins over `minor`).

---

## Final filesystem state

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   ├── 0001-add-export.yaml
│   └── 0002-remove-legacy-api.yaml
└── CHANGELOG.md
```
