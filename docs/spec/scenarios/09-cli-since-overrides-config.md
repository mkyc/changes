# Scenario 09 — CLI since overrides config

Command-line `--since` takes precedence over `.changes/config.yaml`.

**Scope:** configuration precedence.

---

## Given

### Repository

- Default branch `main`; latest tag `v1.2.3`.
- Branch `feat/search` with two commits:

  | Hash (short) | On `main`? | Author date | Message |
  |--------------|------------|-------------|---------|
  | `eee5555` | yes | `2026-06-16` | `docs: update README` |
  | `fff6666` | no | `2026-06-17` | `feat(search): add full-text index` |

- Working tree clean; current branch `feat/search`.

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

`.changes/config.yaml`:

```yaml
since: v1.0.0
```

---

## Step 1 — propose with CLI override

**When**

```text
changes propose --since main
```

**Then** file `.changes/0001-add-full-text-index.yaml` is created with exactly:

```yaml
id: "0001-add-full-text-index"
package: "."
event: change
increment: minor
date: "2026-06-17"
summary: "Add full-text index"
details: |
  feat(search): add full-text index
breaking: false
source:
  - hash: "fff6666"
    message: "feat(search): add full-text index"
```

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.
- Commit `eee5555` is not in `source` (config `since: v1.0.0` would have included it; CLI `--since main` excludes it).

---

## Final filesystem state

```text
.
├── .changes/
│   ├── 0000-init.yaml
│   ├── 0001-add-full-text-index.yaml
│   └── config.yaml
└── (no CHANGELOG.md)
```
