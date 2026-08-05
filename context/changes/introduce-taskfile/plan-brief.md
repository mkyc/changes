# Introduce Taskfile — Plan Brief

> Full plan: `context/changes/introduce-taskfile/plan.md`

## What & Why

Add a `Taskfile.yml` that wraps `build`, `test`, `doctor`, and `cleanup` into single `task <name>` commands, so contributors don't need to remember the raw `go build` / `go test` / `go vet` / `golangci-lint` incantations by hand.

## Starting Point

No build automation exists today — no Makefile, Taskfile, or CI. `task` and `golangci-lint` are already installed locally, but there's no `.golangci.yml`, so linting will run with default rules. `.gitignore` excludes `/changes` (an ad-hoc built binary) but not a `bin/` directory.

## Desired End State

Running `task build`, `task test`, `task doctor`, or `task cleanup` from the repo root does the expected thing: produces `bin/changes`, runs the test suite verbosely, runs a battery of static checks, or wipes build artifacts — with `bin/` gitignored throughout.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) | Source |
|---|---|---|---|
| Build output location | `bin/changes` | Keeps repo root clean vs. the existing ad-hoc `./changes` binary; standard Go convention. | Plan |
| Doctor scope | fmt + vet + `go mod tidy` drift check + `golangci-lint run` | Layers explicit module-hygiene checks on top of golangci-lint's defaults rather than relying on lint alone. | Plan |
| Cleanup scope | Remove `bin/` + `go clean -testcache -cache` | Covers everything this Taskfile itself creates without touching unrelated caches. | Plan |
| Test command | `go test ./... -v` | Simple, verbose, matches typical local dev feedback loop; no coverage tracking added. | Plan |

## Scope

**In scope:**
- `Taskfile.yml` with `build`, `test`, `doctor`, `cleanup` tasks
- `.gitignore` update to exclude `bin/`

**Out of scope:**
- `.golangci.yml` custom lint configuration
- CI setup (GitHub Actions etc.)
- Coverage reporting/thresholds
- Removing the existing ad-hoc `./changes` binary
- Install/release/packaging tasks

## Architecture / Approach

A single flat `Taskfile.yml` (schema `version: '3'`) at repo root, one task per acceptance criterion, each shelling out to standard Go toolchain commands — no new dependencies introduced.

## Phases at a Glance

| Phase | What it delivers | Key risk |
|---|---|---|
| 1. Add Taskfile and gitignore entry | Working `Taskfile.yml` + updated `.gitignore` | `go mod tidy` drift check or golangci-lint defaults may surface pre-existing issues unrelated to this change |

**Prerequisites:** `task` and `golangci-lint` installed locally (already the case in this environment).
**Estimated effort:** Single session, one phase.

## Open Risks & Assumptions

- Assumes golangci-lint's default linter set (no `.golangci.yml`) passes cleanly on the current codebase; if not, `task doctor` will fail on pre-existing issues outside this change's scope, and those would need to be triaged separately.
- Assumes `go mod tidy` produces no diff on the current `go.mod`/`go.sum`.

## Success Criteria (Summary)

- `task build`, `task test`, `task doctor`, `task cleanup` all run successfully and do what their names imply.
- `bin/` never appears in `git status`.
- `bin/changes --help` behaves identically to `go run . --help`.
