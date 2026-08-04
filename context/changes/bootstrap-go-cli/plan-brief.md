# Bootstrap the Go CLI and default configuration — Plan Brief

> Full plan: `context/changes/bootstrap-go-cli/plan.md`

## What & Why

Stand up the `changes` Go CLI skeleton: module init, a cobra command tree
for `init`/`propose`/`apply`/`check`, default configuration, and injectable
dependencies (clock, filesystem, streams). This unblocks tasks 0002-0009,
which depend on it for the real domain logic and Git adapter.

## Starting Point

No Go code exists in the repo yet — no `go.mod`, no `main.go`. The target
shape is fully specified in `docs/DESIGN.md` and `docs/spec/DECISIONS.md`;
this plan only wires the skeleton, not the business logic.

## Desired End State

`go build` produces a `changes` binary. `changes --help` lists all four
subcommands; `propose --help` shows `--since`. With no
`.changes/config.yaml`, every command resolves config to the documented
defaults and exits 0 silently. All commands run through injected
dependencies so tests never touch the real filesystem or clock.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) | Source |
|---|---|---|---|
| Config library | viper + cobra flags | Free CLI>config>default precedence (D010), no custom merge logic | Plan |
| Dependency injection | `Deps` struct passed into `RunE` closures | Explicit, type-safe, trivially fakeable — avoids `context.Value` anti-pattern | Plan |
| Git adapter | Deferred entirely — no type defined yet | Task 0003 owns the interface shape; `Deps` stays an extensible plain struct so it can add a field later | Plan |
| Stub command behavior | Silent success (exit 0, no output) | Matches acceptance criteria and Scenario 01's empty-stdout/stderr expectation | Plan |
| Filesystem abstraction | spf13/afero | Pairs with existing spf13 ecosystem choice; in-memory fs for isolated tests | Plan |
| Error output wiring | D011 (`error: <msg>`, exit 1) implemented now as shared `main.go` plumbing | Every future task inherits correct formatting instead of reimplementing it | Plan |
| Smoke test scope | In-process cobra tests only (command tree, config defaults, stub exit codes) | Matches acceptance criteria exactly; no business logic exists yet to integration-test | Plan |

## Scope

**In scope:**
- Go module init, cobra/viper/afero dependencies
- Root command + 4 stub subcommands, `propose --since` flag
- Config defaults + `.changes/config.yaml` + CLI precedence (D010)
- Injectable `Deps` (clock, FS, stdout, stderr)
- Shared D011 error formatting
- Smoke tests for command discovery, config defaults, exit behavior

**Out of scope:**
- Any real init/propose/apply/check logic (tasks 0004-0008)
- Change-file YAML model/storage (task 0002)
- Git adapter of any kind (task 0003)
- Config schema validation beyond default keys (task 0008 / `check`)
- Versioning/changelog rendering

## Architecture / Approach

`main.go` builds a real `Deps` and calls `internal/cli.NewRootCmd(deps)`,
then handles the returned error per D011. `internal/cli` holds the cobra
tree; `internal/config` holds the viper-backed loader; `internal/app` holds
`Deps` and its real/fake constructors.

## Phases at a Glance

| Phase | What it delivers | Key risk |
|---|---|---|
| 1. Module & Dependency Setup | `go.mod`, cobra/viper/afero deps, minimal `main.go` | Low — mechanical |
| 2. Deps/App Scaffolding | `Deps` struct + real/fake constructors | Low — struct shape may need a field later (expected, by design) |
| 3. Config Resolution | Default + file + CLI precedence loader | Medium — viper + afero integration needs the config file read through the injected FS, not viper's own OS reads |
| 4. Cobra Command Wiring | Root + 4 stub subcommands, D011 error handling | Low — mostly wiring |
| 5. Smoke Tests | Discovery/config/exit-code tests | Low |

**Prerequisites:** None — greenfield.
**Estimated effort:** ~1 session, 5 phases.

## Open Risks & Assumptions

- The `GitRunner` field is deliberately absent from `Deps`; task 0003 must
  extend the struct rather than replace the DI pattern. If task 0003's
  interface needs differ significantly from a simple struct-field
  extension, this plan's DI choice should be revisited.
- Viper's native file-reading bypasses the injected `afero.Fs`; the plan
  requires reading `.changes/config.yaml` through `fs` and feeding bytes to
  `viper.ReadConfig` to preserve test isolation — this is a specific
  integration detail called out in Phase 3, not left to chance.

## Success Criteria (Summary)

- `go build`, `go vet`, `go test ./...` all pass
- `changes --help` / `changes propose --help` show the right commands/flags
- Every subcommand exits 0 silently with no config file present
- Config defaults match `docs/DESIGN.md` exactly with no `.changes/config.yaml`
