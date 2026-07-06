# Scenario 17 — Apply second change

After a prior apply cycle, a new propose + apply bumps the version and archives the next changeset.

**Scope:** [D004](../DECISIONS.md#d004-version-computation), [D007](../DECISIONS.md#d007-apply-behavior), sequence continuity.

---

## Given

### Repository

- Default branch `main`; latest tag `v1.3.0`.
- Branch `fix/typo` with one commit not on `main`:

  | Hash (short) | Author date | Message |
  |--------------|-------------|---------|
  | `xyz9999` | `2026-06-21` | `fix(ui): correct button label` |

- Working tree clean; current branch `fix/typo`.

### Filesystem

State after [scenario 01](01-single-repo-happy-path.md):

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
  Includes nil-input guard on the handler.
breaking: false
issues: ["#42"]
authors: ["@alice"]
source:
  - hash: "abc1234"
    message: "feat(widget): add API endpoint"
  - hash: "def5678"
    message: "fix(widget): handle nil input"
```

`CHANGELOG.md`:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2026-06-19

### Added

- Add widget API endpoint (#42) (@alice)

  Adds REST endpoint for widget creation.
  Includes nil-input guard on the handler.

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

### Environment

- `changes propose` runs at `2026-06-21T10:00:00Z`.
- `changes apply` runs at `2026-06-21T10:01:00Z`.

---

## Step 1 — propose

**When**

```text
changes propose
```

**Then** file `.changes/0002-correct-button-label.yaml` is created with exactly:

```yaml
id: "0002-correct-button-label"
package: "."
event: change
increment: patch
date: "2026-06-21"
summary: "Correct button label"
details: |
  fix(ui): correct button label
breaking: false
source:
  - hash: "xyz9999"
    message: "fix(ui): correct button label"
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- Sequence is `0002` (continues after `0001` in `.changes/released/`).

**When** the user commits `.changes/0002-correct-button-label.yaml` unchanged.

---

## Step 2 — apply

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

## [1.3.1] - 2026-06-21

### Fixed

- Correct button label

  fix(ui): correct button label

## [1.3.0] - 2026-06-19

### Added

- Add widget API endpoint (#42) (@alice)

  Adds REST endpoint for widget creation.
  Includes nil-input guard on the handler.

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- Release version written to changelog is `1.3.1` (`1.3.0` + one `patch`).
- `.changes/0002-correct-button-label.yaml` is moved to `.changes/released/0002-correct-button-label.yaml`.

---

## Step 3 — check

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
│       ├── 0001-add-widget-api.yaml
│       └── 0002-correct-button-label.yaml
└── CHANGELOG.md
```
