# Scenario 06 — Propose with no new commits

No commits since the since-ref means propose creates nothing.

**Scope:** single repo, initialized, branch has no commits ahead of `main`.

---

## Given

### Repository

- Default branch `main`; latest tag `v1.2.3`.
- Current branch `main`.
- No commits on `main` after tag `v1.2.3`.
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

## Step 1 — propose

**When**

```text
changes propose
```

**Then**

- No new files are created under `.changes/`.
- `.changes/0000-init.yaml` is unchanged.

**And**

- Exit code is `0`.
- Stdout is empty.
- Stderr is empty.

---

## Final filesystem state

```text
.
├── .changes/
│   └── 0000-init.yaml
└── (no CHANGELOG.md)
```
