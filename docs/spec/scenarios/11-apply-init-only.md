# Scenario 11 — Apply with init only

Changelog from init baseline alone; no `[Unreleased]` section when no pending changes exist.

**Scope:** boundary condition, [D003](../DECISIONS.md#d003-changelog-format).

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

## [1.2.3] - 2026-06-18

- Start tracking changes with `changes`.
  Last release before adopting the tool was 1.2.3.
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- No `## [Unreleased]` heading appears in `CHANGELOG.md`.
- Computed next release version is `1.2.3`.

---

## Final filesystem state

```text
.
├── .changes/
│   └── 0000-init.yaml
└── CHANGELOG.md
```
