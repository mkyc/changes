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
| **Check** | Validate schema, ID ordering, and that apply would not drift committed output |

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
8. **Single-repo default package** — when no monorepo config exists, all changes use `package: "."`.

## Scenario index

| # | File | Category | Status |
|---|------|----------|--------|
| 01 | [single-repo-happy-path](scenarios/01-single-repo-happy-path.md) | Happy path | draft |

## Open questions

Decisions that must be clarified before implementation. New scenarios may add items here.

1. **Consumed changesets** — archive/delete after release, or move to `.changes/released/`?
2. **Concurrent PRs** — how to handle two PRs both allocating `0005-*.yaml`?
3. **Changelog format configurability** — scenario 01 pins Keep a Changelog 1.1.0; should other templates be supported later?
4. **Version string source** — is the numeric version always computed from `init` + increments, or can it be overridden in config or CLI?
5. **Default `--since`** — when config omits `since`, is the default `main`, current branch upstream, or merge-base with default branch?
6. **Init without prior tag** — what should `0000-init.yaml` contain when the repo has no releases yet (`0.0.0`, `0.1.0`, or require explicit user input)?
7. **Apply side effects** — does `apply` only write `CHANGELOG.md`, or also emit/version-bump other files (e.g. `package.json`)?

**Resolved in scenario 01**

- **Propose granularity** — one change file per propose run, grouping all commits since `--since` that are not already referenced in an existing change file's `source`.
- **Changelog format** — Keep a Changelog 1.1.0 (exact template pinned in scenario file).
- **Check output** — exit code `0`, stdout `ok\n`, stderr empty.
