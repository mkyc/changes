# changes — design notes (draft)

Small, language-independent Go CLI for changelog generation.

## Goals

| Goal | Notes |
|------|-------|
| Generate changesets from commits | Since a base ref (branch, tag, merge-base) |
| Derive `CHANGELOG.md` | From changeset files, not raw git history |
| Compute semantic versions | Aggregate bump type across changesets |
| Support no-release changes | Docs, chores — tracked but no version bump |
| Single repo + monorepo | Package-scoped changesets where needed |

## Core concepts

### Change file

One file per logical change, stored under `.changes/`:

```yaml
# .changes/0001-add-widget-api.yaml
id: "0001-add-widget-api"
package: "."                    # monorepo: path or package name; single repo: "."
event: change                   # change | init
increment: minor                # major | minor | patch | none
date: "2026-06-17"
summary: "Add widget API endpoint"
details: |
  Optional longer description for changelog.
breaking: false
issues: ["#42"]
authors: ["@alice"]
source:
  - hash: "abc1234"
    message: "feat(widget): add API"
  - hash: "def5678"
    message: "fix(widget): handle edge case"
```

Use `event: init` when adopting the tool in an existing repo:

```yaml
# .changes/0000-init.yaml
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

- **`increment: none`** — included in changelog but does not contribute to version bump.
- **`event: init`** — records the version at the point the tool was introduced.
- **Sequential IDs** (`0001`, `0002`, …) keep ordering stable and diff-friendly.
- Files are the **source of truth** for release notes; commits are hints for `propose`, not the changelog itself.

### Version bump aggregation

```
major > minor > patch > none
```

If any changeset is `major`, release is major. Else if any is `minor`, minor. Else if any is `patch`, patch. If all are `none`, no release.

Monorepo: aggregate **per package** independently.

## CLI (sketch)

| Command | Purpose |
|---------|---------|
| `changes propose [--since <ref>]` | Inspect commits since ref; emit draft `.changes/NNNN-*.yaml` |
| `changes apply [--package <path>]` | Regenerate `CHANGELOG.md` from `.changes/` |
| `changes check` | Validate schema, IDs, ordering; verify `apply` is clean (no drift) |

`propose` maps conventional commit prefixes to suggested `increment` (configurable):

| Prefix | Suggested increment |
|--------|---------------------|
| `feat!`, `BREAKING CHANGE` | major |
| `feat` | minor |
| `fix` | patch |
| `docs`, `chore`, `refactor`, `test`, `ci` | none |

User always edits the draft before committing.

## Workflow

```
┌─────────────┐   propose   ┌───────────────────┐
│ new commits │ ──────────► │ .changes/000N.yaml │  (draft, edit manually)
└─────────────┘             └─────────┬─────────┘
                                      │ apply
                                      ▼
                             ┌─────────────────┐
                             │  CHANGELOG.md   │
                             └─────────────────┘
```

1. `changes propose` → review/edit `.changes/000N-*.yaml`
2. `changes apply` → commit changesets + changelog together
3. CI: `changes check` validates schema and drift

Monorepo: same flow with `--package`, one `CHANGELOG.md` per package.

## Configuration (`.changes/config.yaml`, optional)

```yaml
changelog: CHANGELOG.md
packages:
  - path: "."
    tag_prefix: "v"
since: main
conventional:
  minor: ["feat"]
  patch: ["fix"]
  none: ["docs", "chore", "refactor", "test", "ci"]
```

## Open questions

1. **Consumed changesets** — archive/delete after release, or move to `.changes/released/`?
2. **Concurrent PRs** — how to handle two PRs both allocating `0005-*.yaml`?
3. **Changelog format** — Keep a Changelog style, or configurable template?

---

*Draft for internal discussion. Commands, paths, and schema are proposals, not implementation.*

## Inspirations

- [changesets](https://github.com/changesets/changesets)
- [semantic-release](https://github.com/semantic-release/semantic-release)
- [keepachangelog](https://github.com/olivierlacan/keep-a-changelog)
