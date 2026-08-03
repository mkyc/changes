# Scenario 18 — Not a git repository

Commands that require git fail when run outside a repository.

**Scope:** missing dependency, error handling.

---

## Given

### Filesystem

```text
/tmp/not-a-repo/
└── README.md
```

- Directory `/tmp/not-a-repo` is not a git repository.
- Current working directory is `/tmp/not-a-repo`.

---

## Step 1 — init

**When**

```text
changes init
```

**Then**

- Exit code is `1`.
- Stdout is empty.
- Stderr is exactly:

```text
error: not a git repository
```

**And**

- No `.changes/` directory is created.

---

## Step 2 — propose

**When**

```text
changes propose
```

**Then**

- Exit code is `1`.
- Stdout is empty.
- Stderr is exactly:

```text
error: not a git repository
```

**And**

- No files are created or modified.

---

## Final filesystem state

```text
/tmp/not-a-repo/
└── README.md
```
