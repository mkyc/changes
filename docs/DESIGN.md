# changes — design notes (draft)

Small, language-independent Go CLI for versioning workflows. Combines ideas from **semver**, **conventional commits**, and **changesets** into automation-friendly artifacts that work locally and in CI.

## Goals

| Goal | Notes |
|------|-------|
| Generate changesets from commits | Since a base ref (branch, tag, merge-base) |
| Derive `CHANGELOG.md` | From committed changeset files, not from raw git history |
| Compute semantic versions | Aggregate bump type across changesets |
| Signal release intent | Machine-readable output for downstream automations |
| Support no-release changes | Docs, chores, internal refactors — tracked but no version bump |
| Single repo + monorepo | Package-scoped changesets where needed |
| Language independent | YAML/JSON artifacts; no runtime coupling |

## Core concepts

### Changeset file

One file per logical change, stored under `.changes/` (name/format TBD):

```yaml
# .changes/0001-add-widget-api.yaml  (example shape)
id: "0001-add-widget-api"
package: "."                    # monorepo: path or package name; single repo: "."
type: minor                     # major | minor | patch | none
summary: "Add widget API endpoint"
details: |
  Optional longer description for changelog.
breaking: false
issues: ["#42"]
authors: ["@alice"]
source:
  commits: ["abc1234", "def5678"]   # optional provenance from `propose`
  conventional: "feat(widget): add API"  # optional
```

- **`type: none`** — included in changelog (optional) but does **not** contribute to version bump.
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
| `changes plan [--since-tag <tag>]` | Output intended bump(s) and next version(s) without tagging |
| `changes release [--dry-run]` | Compute version, tag, optionally emit release metadata for CI |

`propose` maps conventional commit prefixes to suggested `type` (configurable):

| Prefix | Suggested type |
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
3. **Post-merge CI (main):**
   - `changes plan` / `changes release` on changesets merged since last tag.
   - Create semver tag if bump ≠ none.
   - Emit release signal (see below).

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
- CI can release only packages with pending non-`none` changesets.

## CI integration

### PR checks

```bash
changes check --require-new          # fail if no new .changes file in PR
changes apply && git diff --exit-code CHANGELOG.md
```

### Release pipeline (main)

```bash
changes plan --format json           # → { "package": ".", "current": "1.2.3", "next": "1.3.0", "bump": "minor" }
changes release --dry-run            # validate before tag push
changes release --format json        # → release manifest for downstream jobs
```

### Signaling other automations

Machine-readable **release manifest** (stdout or artifact), e.g.:

```json
{
  "releases": [
    {
      "package": ".",
      "previous": "1.2.3",
      "next": "1.3.0",
      "bump": "minor",
      "tag": "v1.3.0",
      "changesets": ["0001-add-widget-api", "0002-fix-login"],
      "changelog_path": "CHANGELOG.md"
    }
  ],
  "skipped": [
    { "package": ".", "reason": "all changesets type none" }
  ]
}
```

Downstream systems (build, publish, deploy, Slack, etc.) consume this without parsing git or markdown.

## Configuration (`.changes/config.yaml`)

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

## Non-goals (for now)

- Replacing git or hosting release pages.
- Enforcing commit message format (conventional commits are a **hint**, not a gate).
- Package registry publish logic (emit manifest; let other tools publish).

## Open questions

1. **One changeset per PR vs many?** — Sequential files allow stacked PRs; need merge/conflict story for ID allocation.
2. **Consumed changesets** — After release, archive/delete merged changesets or move to `.changes/released/`?
3. **Pre-release versions** — Support `-alpha`, `-rc` in plan/release?
4. **ID generation** — Auto-increment from max existing, or content-hash based?
5. **Changelog format** — Keep a simple Keep-a-Changelog style, or configurable template?
6. **Merge queue** — How to handle concurrent PRs both adding `0005-*.yaml`?

## Comparison (intent)

| | changesets (JS) | semantic-release | **changes** (this tool) |
|--|-----------------|------------------|-------------------------|
| Artifact | `.md` in `.changeset/` | none (git only) | `.yaml` in `.changes/` |
| Changelog | generated at release | generated at release | generated at **PR time** (`apply`) |
| Version bump | human declares in file | inferred from commits | declared in file, validated in CI |
| Monorepo | yes | plugins | first-class `package` field |
| Language | JS ecosystem | JS ecosystem | **agnostic** |

---

*Draft for internal discussion. Commands, paths, and schema are proposals, not implementation.*
