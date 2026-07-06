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

## Core concepts

| Concept | Role |
|---------|------|
| **Change file** | One YAML file under `.changes/` describing a single logical change; source of truth for release notes |
| **Package** | Scope of a change; `"."` for a single-repo project |
| **Event** | `change` (normal entry) or `init` (adoption baseline) |
| **Increment** | Suggested semver bump: `major`, `minor`, `patch`, or `none` |
| **Init** | Detect latest release tag; emit `.changes/0000-init.yaml` |
| **Propose** | Inspect git commits since a ref; emit draft change file(s) |
| **Apply** | Regenerate `CHANGELOG.md` from committed change files |
| **Check** | Validate schema, unique sequence IDs, ordering, and that apply would not drift committed output |

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
- **Committed `.changes/`** — authoritative change history for apply/check
- **`CHANGELOG.md`** — derived artifact; must match apply output when check passes
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
9. **Single-repo default package** — when no monorepo config exists, all changes use `package: "."`.

## Scenario index

| # | File | Category | Status |
|---|------|----------|--------|
| 01 | [single-repo-happy-path](scenarios/01-single-repo-happy-path.md) | Happy path | draft |

## Decisions

Product and behavior decisions live in [DECISIONS.md](DECISIONS.md). Scenario-specific pins (exact file contents, stdout) stay in individual scenario files.

## Open questions

None currently. Add new items here when scenarios surface unresolved behavior; move to [DECISIONS.md](DECISIONS.md) once decided.
