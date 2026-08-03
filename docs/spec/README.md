# Behavioral specification

Stress-tests the ideas in [DESIGN.md](../DESIGN.md) before implementation. Each file under `scenarios/` describes observable behavior in **Given / When / Then / And** form so a future implementation can be validated against it.

This is not an architecture or API design document.

## How to read scenarios

| Section | Meaning |
|---------|---------|
| **Given** | Initial state, inputs, assumptions, environment |
| **When** | A user action, system event, or operation |
| **Then** | Observable outcomes |
| **And** | Guarantees, side effects, constraints, or absence of side effects |

Scenario files are numbered and named by theme (`01-single-repo-happy-path.md`, …). Add new files as coverage grows; do not cram every category into one file.

Each step documents **one** `changes` command. Use a **Precondition** subsection for repo state required before that command. Git commits and branch checkouts are ordinary setup outside tool validation unless a scenario explicitly tests them.

## Core concepts

| Concept | Role |
|---------|------|
| **Change file** | One YAML file under `.changes/` describing a single logical change; source of truth for release notes |
| **Package** | Scope of a change; `"."` for a single-repo project |
| **Event** | `change` (normal entry) or `init` (adoption baseline) |
| **Increment** | Suggested semver bump: `major`, `minor`, `patch`, or `none` |
| **Init** | Detect latest release tag; emit `.changes/0000-init.yaml` |
| **Propose** | Inspect git commits since a ref; emit draft change file(s) |
| **Apply** | Regenerate `CHANGELOG.md`, compute version, move pending change files to `.changes/released/` |
| **Check** | Validate schema, unique sequence IDs, ordering, and that apply would not drift committed output |

## Workflow

Minimal day-to-day flow ([DESIGN.md](../DESIGN.md)):

```
propose  →  edit & commit .changes/000N.yaml  →  apply  →  commit CHANGELOG.md + moved changesets
```

`apply` does both changelog generation and archival: pending files land in `.changes/released/` and the changelog gets a `## [X.Y.Z] - <date>` section for the computed version ([D004](DECISIONS.md#d004-version-computation), [D007](DECISIONS.md#d007-apply-behavior)).

| Location | Contents |
|----------|----------|
| `.changes/0000-init.yaml` | Adoption baseline; always in root |
| `.changes/000N-*.yaml` | **Pending** changes awaiting the next `apply` |
| `.changes/released/000N-*.yaml` | **Applied** changes; already reflected in `CHANGELOG.md` |

## Inputs

- Git history (commits, refs, merge-base)
- Existing `.changes/` directory (change files, optional `config.yaml`)
- Optional CLI flags (`--since`, `--package`)
- User edits to draft change files before commit

## Outputs

- Draft or committed `.changes/NNNN-*.yaml` files
- Updated `CHANGELOG.md`
- Exit code and human-readable messages (success, validation error, dependency failure)
- Implied release version (derived from `init` + aggregated increments)

## State

- **Working tree** — draft files from propose, uncommitted edits
- **Committed `.changes/`** — pending change files (root) and released archive (`released/`)
- **`CHANGELOG.md`** — derived artifact; version sections written by `apply`
- **Git refs** — used only by propose to discover new commits, not to build the changelog directly

## External dependencies

- `git` executable and a valid git repository
- Filesystem read/write under repo root
- Optional `.changes/config.yaml` for defaults (`since`, paths, conventional-commit mapping)

## Fundamental invariants

These must hold regardless of user actions or execution order:

1. **Change files are authoritative** — `apply` and `check` never derive changelog content from raw git history; only from `.changes/`.
2. **Stable ordering** — change file IDs are sequential (`0000`, `0001`, …) and determine changelog order.
3. **Increment precedence** — when aggregating bumps for a release: `major` > `minor` > `patch` > `none`.
4. **Init is baseline** — `event: init` records the version at tool adoption; subsequent releases build on aggregated increments from `change` entries after init.
5. **None does not bump** — `increment: none` entries appear in the changelog but do not affect semver.
6. **Propose is non-destructive** — propose creates or updates draft change files only; it does not modify `CHANGELOG.md` or rewrite existing committed change files without explicit user action.
7. **Check detects drift** — if committed `CHANGELOG.md` differs from what `apply` would produce, `check` fails.
8. **Unique sequence IDs** — at most one change file per numeric sequence prefix (`0000`, `0001`, …) under `.changes/` and `.changes/released/` combined; duplicates fail `check`.
9. **Contiguous sequences** — sequence numbers are unique across `.changes/` and `.changes/released/` combined ([D002](#d002-concurrent-pr-sequence-collisions)). Pending files continue the sequence after the highest released number with no gaps ([D012](#d012-sequence-gap-validation)).
10. **Single-repo default package** — when no monorepo config exists, all changes use `package: "."`.

## Scenario index

| # | File | Category |
|---|------|----------|
| 01 | [single-repo-happy-path](scenarios/01-single-repo-happy-path.md) | Happy path |
| 02 | [init-no-prior-tag](scenarios/02-init-no-prior-tag.md) | Init + propose + apply (no tag) |
| 03 | [propose-breaking-major](scenarios/03-propose-breaking-major.md) | Propose / conventional commits |
| 04 | [propose-patch-only](scenarios/04-propose-patch-only.md) | Propose / conventional commits |
| 05 | [propose-none-increment](scenarios/05-propose-none-increment.md) | Propose / no-release |
| 06 | [propose-no-new-commits](scenarios/06-propose-no-new-commits.md) | Propose / boundary |
| 07 | [propose-skips-consumed-commits](scenarios/07-propose-skips-consumed-commits.md) | Propose / idempotency |
| 08 | [config-default-since](scenarios/08-config-default-since.md) | Configuration |
| 09 | [cli-since-overrides-config](scenarios/09-cli-since-overrides-config.md) | Configuration precedence |
| 10 | [version-aggregation-mixed-increments](scenarios/10-version-aggregation-mixed-increments.md) | Version aggregation |
| 11 | [apply-init-only](scenarios/11-apply-init-only.md) | Apply / boundary |
| 12 | [apply-idempotent](scenarios/12-apply-idempotent.md) | Apply / idempotency |
| 13 | [check-duplicate-sequence](scenarios/13-check-duplicate-sequence.md) | Check / conflict |
| 14 | [check-changelog-drift](scenarios/14-check-changelog-drift.md) | Check / drift + recovery |
| 15 | [check-invalid-schema](scenarios/15-check-invalid-schema.md) | Check / invalid input |
| 16 | [check-sequence-gap](scenarios/16-check-sequence-gap.md) | Check / ordering |
| 17 | [apply-second-change](scenarios/17-apply-second-change.md) | Apply / second cycle |
| 18 | [not-a-git-repository](scenarios/18-not-a-git-repository.md) | Missing dependency |
| 19 | [init-already-exists](scenarios/19-init-already-exists.md) | Init / error |
| 20 | [major-wins-aggregation](scenarios/20-major-wins-aggregation.md) | Version aggregation |

### Coverage map

| Original category | Scenarios |
|-------------------|-----------|
| Happy path | 01 |
| Invalid input | 15, 19 |
| Missing / unavailable dependencies | 18 |
| Boundary conditions | 02, 06, 11, 20 |
| Error handling | 13, 15, 16, 18, 19 |
| Recovery behavior | 14 |
| Repeated execution / idempotency | 07, 12 |
| Configuration resolution and precedence | 08, 09 |
| State transitions | 17 |
| Concurrent or conflicting operations | 13 |
| Data consistency | 10, 14, 20 |
| Security / permissions | — (deferred) |
| Upgrade / migration / compatibility | — (deferred) |
| Monorepo | — (deferred) |

## Decisions

Product and behavior decisions live in [DECISIONS.md](DECISIONS.md). Scenario-specific pins (exact file contents, stdout) stay in individual scenario files.

## Open questions

1. **Large history on adoption** — when `changes init` runs in a repo that already has extensive history (e.g. 10 000 commits on `main`, no tags), what should `changes propose` do?
   - Ignore all pre-init history and only track commits after adoption (current assumption in [scenario 02](scenarios/02-init-no-prior-tag.md))?
   - Backfill one or many change files from the full history?
   - Require an explicit `--since` (or other flag) before proposing against old history?
   - Refuse / warn when commit count exceeds a threshold?
2. **Cleanup of released change files** — may users delete archives under `.changes/released/` after some time? If so, how does `apply`/`check` reconstruct history?
3. Why 0000-init.yaml stays in root?
4. in scenario 05, I'm not sure if I like that we release 1.2.3 even if that was already in repo state. I think this also touches what should we do with "none" increments on apply. That is also with scenario 11. 
5. Do we need `id` field in change files?

Add new items here when scenarios surface unresolved behavior; move to [DECISIONS.md](DECISIONS.md) once decided.
