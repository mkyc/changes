# Scenario 17 — Release moves changesets to released

After release, pending change files move to `.changes/released/` and `apply` promotes them into a versioned changelog section.

**Scope:** state transition, [D001](../DECISIONS.md#d001-consumed-changesets), [D004](../DECISIONS.md#d004-version-computation).

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
issues: ["#42"]
authors: ["@alice"]
source:
  - hash: "abc1234"
    message: "feat(widget): add API endpoint"
```

`CHANGELOG.md` (pre-release, with `[Unreleased]`):

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add widget API endpoint (#42) (@alice)

  Adds REST endpoint for widget creation.

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

### Environment

- Release move and `changes apply` run at `2026-06-20T10:00:00Z`.

---

## Step 1 — move consumed changeset

**When** the user runs:

```text
mkdir -p .changes/released
mv .changes/0001-add-widget-api.yaml .changes/released/0001-add-widget-api.yaml
```

**Then** filesystem contains:

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   └── released/
│       └── 0001-add-widget-api.yaml
└── CHANGELOG.md
```

**And**

- `.changes/released/0001-add-widget-api.yaml` content is byte-for-byte identical to the Given `0001` file.

---

## Step 2 — apply after release

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

## [1.3.0] - 2026-06-20

### Added

- Add widget API endpoint (#42) (@alice)

  Adds REST endpoint for widget creation.

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- No `## [Unreleased]` heading appears in `CHANGELOG.md`.
- Computed next release version is `1.3.0` (no pending changes remain).
- Only `CHANGELOG.md` is modified; released change file is not moved or edited.

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
│       └── 0001-add-widget-api.yaml
└── CHANGELOG.md
```
