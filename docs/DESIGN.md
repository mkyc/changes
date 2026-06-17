# changes — design notes (draft)

Small, language-independent Go CLI for changelog generation.

## Goals

| Goal | Notes                                                                       |
|------|-----------------------------------------------------------------------------|
| Generate changesets from commits | Since a base ref (branch, tag, merge-base)                                  |
| Derive `CHANGELOG.md` | From changeset files, not from raw git history                              |
| Compute semantic versions | Aggregate bump type across changesets                                       |
| Support no-release changes | Docs, chores, internal refactors — tracked but no version bump              |
| Single repo + monorepo | Package-scoped changesets where needed                                      |
| Language independent | YAML/JSON artifacts; no runtime coupling                                    |
| Change is event | Change entry might be ie.: code change entry, manual release tag event, etc |

## Core concepts

### Change file

One file per logical change, stored under `.changes/` (name/format TBD):

```yaml
# .changes/0001-add-widget-api.yaml  (example shape)
id: "0001-add-widget-api"
package: "."                    # monorepo: path or package name; single repo: "."
event: change                     # change | init
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

Another example for repository initialization (`event: init`):

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

- **`increment: none`** — included in changelog (optional) but does **not** contribute to version bump.
- **`event: change`** — change entry is a code change.
- **`event: init`** — change entry is an initial version number ie.: when we introduce changelog in existing repo with previous versions.
- **Sequential IDs** (`0001`, `0002`, …) keep ordering stable and diff-friendly.
- Files are the **source of truth** for release notes; commits are hints for `propose`, not the changelog itself.

### Version bump aggregation

When multiple changesets land before a release:

```
major > minor > patch > none
```

If any changeset is `major`, release is major. Else if any is `minor`, minor. Else if any is `patch`, patch. If all are `none`, no release.

Monorepo: aggregate **per package** independently.

## CLI (sketch)

| Command | Purpose |
|---------|---------|
| `changes propose [--since <ref>] [--base <ref>]` | Inspect commits since ref; emit draft `.changes/NNNN-*.yaml` |
| `changes apply [--package <path>]` | Regenerate `CHANGELOG.md` (and per-package changelogs in monorepo) from `.changes/` |
| `changes check` | Validate changesets schema, IDs, ordering; verify `apply` is clean (no drift) |

`propose` maps conventional commit prefixes to suggested `increment` (configurable):

| Prefix | Suggested increment |
|--------|----------------|
| `feat!`, `BREAKING CHANGE` | major |
| `feat` | minor |
| `fix` | patch |
| `docs`, `chore`, `refactor`, `test`, `ci` | none (default) |

User always edits the draft before committing.

## Single-repo workflow

```
┌─────────────┐     propose      ┌──────────────────┐
│ new commits │ ───────────────► │ .changes/000N.yaml│ (draft, edit manually)
└─────────────┘                  └────────┬─────────┘
                                          │ apply
                                          ▼
                                 ┌──────────────────┐
                                 │  CHANGELOG.md    │
                                 └────────┬─────────┘
                                          │ commit + push PR
                                          ▼
                                 ┌──────────────────┐
                                 │  CI: check       │
                                 └────────┬─────────┘
                                          │ merge to main
                                          ▼
                                 ┌──────────────────┐
                                 │  CI: release     │ → tag vX.Y.Z, notify automations
                                 └──────────────────┘
```

1. **Before PR:** `changes propose` → review/edit `.changes/000N-*.yaml` → `changes apply` → commit changesets + changelog.
2. **PR CI:**
   - At least one **new** file in `.changes/` (unless explicitly exempt — e.g. docs-only PR policy TBD).
   - `changes check` passes (schema, no duplicate IDs).
   - `changes apply` is idempotent: re-running produces **no diff** on `CHANGELOG.md`.

## Monorepo workflow

Same flow, scoped by `package` field:

```
.changes/
  0001-core-fix.yaml          # package: "packages/core"
  0002-cli-feat.yaml          # package: "packages/cli"
packages/core/CHANGELOG.md
packages/cli/CHANGELOG.md
.changes/config.yaml          # package roots, tag prefix per package (optional)
```

- Independent version lines per package (`@scope/pkg` or path).
- `propose` / `apply` accept `--package` or operate on all changed packages.

## CI integration

### PR checks

```bash
changes check --require-new          # fail if no new .changes file in PR
changes apply && git diff --exit-code CHANGELOG.md
```

Downstream systems (build, publish, deploy, Slack, etc.) consume changelog. 

## Configuration (`.changes/config.yaml` or `.changesrc.json`, etc.)

Minimal, optional:

```yaml
changelog: CHANGELOG.md
packages:
  - path: "."
    tag_prefix: "v"
since: main                    # default base for propose
conventional:
  minor: ["feat"]
  patch: ["fix"]
  none: ["docs", "chore", "refactor", "test", "ci"]
require_changeset_on_pr: true
```

## Open questions

1. **One changeset per PR vs many?** — Sequential files allow stacked PRs; need merge/conflict story for ID allocation.
2. **Consumed changesets** — After release, archive/delete merged changesets or move to `.changes/released/`?
3. **ID generation** — Auto-increment from max existing, or content-hash based?
4. **Changelog format** — Keep a simple Keep-a-Changelog style, or configurable template?
5. **Merge queue** — How to handle concurrent PRs both adding `0005-*.yaml`?

---

*Draft for internal discussion. Commands, paths, and schema are proposals, not implementation.*


## inspirations

 - [changesets](https://github.com/changesets/changesets)
 - [semantic-release](https://github.com/semantic-release/semantic-release)
 - [keepachangelog](https://github.com/olivierlacan/keep-a-changelog)
