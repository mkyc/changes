# Bootstrap the Go CLI and default configuration — Implementation Plan

## Overview

Create the Go application skeleton for `changes`: a cobra-based root command
exposing `init`, `propose`, `apply`, and `check` as stub subcommands, an
injectable `Deps` struct (clock, filesystem, standard streams — Git runner
deferred), viper-backed configuration with documented defaults and
CLI-over-config-over-defaults precedence, shared D011-compliant error
handling, and smoke tests. No business logic for any subcommand is
implemented here — that's tasks 0002-0009.

## Current State Analysis

The repository has no Go code yet: no `go.mod`, no `main.go`, no `internal/`
or `cmd/` directories. `docs/DESIGN.md` and `docs/spec/DECISIONS.md` define
the target CLI shape and behavior; `tasks/0001-bootstrap-go-cli.md` scopes
this task to skeleton + config only, with `tasks/0002`-`0009` depending on it
for the real init/propose/apply/check logic and the Git adapter.

## Desired End State

`go build ./...` produces a runnable `changes` binary. `changes --help` lists
`init`, `propose`, `apply`, `check`; `propose --help` shows an optional
`--since <ref>` flag. Running any of the four subcommands with no
`.changes/config.yaml` present exits 0, writes nothing to stdout/stderr, and
resolves configuration to package `.`, changelog `CHANGELOG.md`, base ref
`main`, tag prefix `v`, and the default conventional-commit increment
mappings from `docs/DESIGN.md`. All commands accept dependencies (clock, FS,
stdout, stderr) through an injected struct so tests can run without touching
the real filesystem, real time, or process-global streams.

**Verification**: `go build ./...`, `go vet ./...`, `go test ./...` all pass;
`./changes --help` and `./changes propose --help` show the expected commands
and flags; manually running each subcommand in an empty temp dir exits 0 with
no output.

### Key Discoveries:

- Config defaults and conventional-commit mappings are fully specified in
  `docs/DESIGN.md:113-126` (`Configuration` section) and confirmed by
  `tasks/0001-bootstrap-go-cli.md:24-28` (Acceptance Criteria).
- D010 (`docs/spec/DECISIONS.md:123-130`) sets precedence as CLI flags >
  `.changes/config.yaml` > built-in defaults.
- D011 (`docs/spec/DECISIONS.md:134-139`) fixes the error contract:
  stderr `error: <message>`, exit code 1, stdout untouched on failure.
- `tasks/0003-git-history-inspection.md` owns the Git adapter; this task
  must not invent that interface. The `Deps` struct is intentionally left
  open (a plain, extensible struct, not a fixed/sealed type) so task 0003
  can add a `GitRunner` field without changing the DI pattern established
  here.
- Scenario 01 (`docs/spec/scenarios/01-single-repo-happy-path.md`) confirms
  successful commands produce empty stdout/stderr — reinforced by the
  chosen "silent success stub" behavior for all four subcommands in this
  task.

## What We're NOT Doing

- No real `init`, `propose`, `apply`, or `check` business logic (tasks
  0004-0008).
- No change-file YAML model or `.changes/` filesystem repository (task
  0002).
- No Git adapter/interface of any kind — not even a stub type (task 0003
  owns its shape entirely).
- No `.changes/config.yaml` schema validation beyond what's needed to load
  the documented default keys — schema/format validation belongs to `check`
  (task 0008).
- No release/versioning logic, no changelog rendering.

## Implementation Approach

Standard Go CLI layout: `main.go` at the repo root wires a real `Deps` and
delegates to `internal/cli.NewRootCmd(deps)`. `internal/cli` holds the cobra
command tree; `internal/config` holds the viper-backed loader and defaults;
`internal/app` holds the `Deps` struct and its real/fake constructors. Cobra
+ viper + afero are added as the only new dependencies, matching the
spf13-ecosystem choice already directed by the task (`spf13 cobra`).

## Phase 1: Module & Dependency Setup

### Overview

Initialize the Go module and pull in cobra, viper, and afero so later
phases have a compiling foundation.

### Changes Required:

#### 1. Go module

**File**: `go.mod` (new)

**Intent**: Initialize the module so the project builds and dependencies
resolve.

**Contract**: Module path `github.com/mkyc/changes`, Go directive matching
the toolchain in use (1.26).

#### 2. Dependencies

**File**: `go.mod`, `go.sum` (updated by tooling)

**Intent**: Add `github.com/spf13/cobra`, `github.com/spf13/viper`, and
`github.com/spf13/afero` as direct dependencies.

**Contract**: Use `go get` for each; no manual version pinning beyond what
`go get` resolves.

#### 3. Entry point

**File**: `main.go` (new, repo root)

**Intent**: Minimal executable entry point — construct a real `Deps`,
build the root command via `internal/cli.NewRootCmd`, execute it, and hand
any returned error to the shared error-formatting logic before exiting.

**Contract**: `func main()` only; no business logic. Exit code comes from
the shared error handler in Phase 4, not from ad hoc `os.Exit` calls
scattered through commands.

### Success Criteria:

#### Automated Verification:

- Module builds: `go build ./...`
- Module is tidy: `go mod tidy` produces no diff

#### Manual Verification:

- `go run . --help` runs without error (command tree not yet wired — a
  bare cobra root or placeholder output is acceptable at this checkpoint;
  full behavior lands in Phase 4)

---

## Phase 2: Deps/App Scaffolding

### Overview

Define the injectable dependency struct used by every command, with real
and fake constructors.

### Changes Required:

#### 1. Deps struct

**File**: `internal/app/deps.go` (new)

**Intent**: Hold the dependencies commands need for deterministic
behavior: a clock, a filesystem, and stdout/stderr writers. Left as a
plain, extensible struct so task 0003 can add a `GitRunner` field later
without touching this DI pattern.

**Contract**: A `Deps` struct with fields `Clock func() time.Time`
(or a small `Clock` interface — implementer's choice, either satisfies
"injectable clock"), `FS afero.Fs`, `Stdout io.Writer`, `Stderr io.Writer`.

#### 2. Real constructor

**File**: `internal/app/deps.go`

**Intent**: Build a `Deps` wired to real time, the real OS filesystem, and
`os.Stdout`/`os.Stderr`, for use by `main.go`.

**Contract**: `func NewRealDeps() *Deps` using `time.Now`,
`afero.NewOsFs()`, `os.Stdout`, `os.Stderr`.

#### 3. Fake constructor

**File**: `internal/app/deps.go` or `internal/app/deps_test.go`

**Intent**: Build a `Deps` suitable for tests — fixed clock, in-memory
filesystem, buffer-backed streams — so Phase 5 smoke tests can assert
exact output and avoid touching the real filesystem or process-global
state.

**Contract**: A test helper returning `*Deps` plus handles to the output
buffers, e.g. `func NewFakeDeps(fixedTime time.Time) (*Deps, *bytes.Buffer, *bytes.Buffer)`.

### Success Criteria:

#### Automated Verification:

- Package builds: `go build ./internal/app/...`
- Unit test confirms `NewFakeDeps` produces isolated, independent buffers
  across two calls: `go test ./internal/app/...`

#### Manual Verification:

- None beyond automated (pure scaffolding, no user-facing behavior)

---

## Phase 3: Config Resolution

### Overview

Load and merge configuration from built-in defaults, `.changes/config.yaml`,
and CLI flags, per D010.

### Changes Required:

#### 1. Config struct and defaults

**File**: `internal/config/config.go` (new)

**Intent**: Define the resolved configuration shape and its built-in
defaults, matching `docs/DESIGN.md:113-126` exactly.

**Contract**: A `Config` struct with fields for package path (default
`.`), changelog path (default `CHANGELOG.md`), base/since ref (default
`main`), tag prefix (default `v`), and conventional-commit increment
mappings (default: `major: [feat!, BREAKING CHANGE]`,
`minor: [feat]`, `patch: [fix, perf, docs, style, refactor, test, build, chore]`,
`none: [ci]`).

#### 2. Loader

**File**: `internal/config/config.go`

**Intent**: Resolve a `Config` by layering, in order, built-in defaults,
then `.changes/config.yaml` if present, then CLI flag overrides — CLI
wins per D010.

**Contract**: `func Load(fs afero.Fs, flags *pflag.FlagSet) (*Config, error)`
using viper with `viper.SetFs`-equivalent handling for the injected `afero.Fs`
(viper reads via its own file access; since afero is required for testability,
read the config file through the injected `fs` and feed viper the bytes via
`viper.ReadConfig`, rather than relying on viper's own OS-level file reads).
Missing `.changes/config.yaml` is not an error — defaults apply silently.

### Success Criteria:

#### Automated Verification:

- Package builds: `go build ./internal/config/...`
- Unit tests pass: `go test ./internal/config/...`
  - No config file present → resolved `Config` matches documented defaults
    exactly (package `.`, changelog `CHANGELOG.md`, base ref `main`, tag
    prefix `v`, default conventional mappings)
  - Config file present with a subset of keys → unset keys still fall back
    to defaults
  - CLI flag set alongside a conflicting config file value → CLI value wins

#### Manual Verification:

- None beyond automated

---

## Phase 4: Cobra Command Wiring

### Overview

Build the root command and four stub subcommands, wired to `Deps` and
`Config`, with shared D011 error formatting.

### Changes Required:

#### 1. Root command constructor

**File**: `internal/cli/root.go` (new)

**Intent**: Construct the cobra root command (`changes`) and attach the
four subcommands, closing over the injected `Deps`.

**Contract**: `func NewRootCmd(deps *app.Deps) *cobra.Command`. Set
`SilenceErrors: true` and `SilenceUsage: true` on the root command so
cobra's default error/usage printing is suppressed in favor of the shared
handler in `main.go`.

#### 2. Subcommand stubs

**File**: `internal/cli/init.go`, `internal/cli/propose.go`,
`internal/cli/apply.go`, `internal/cli/check.go` (new)

**Intent**: One constructor per subcommand (`newInitCmd`, `newProposeCmd`,
`newApplyCmd`, `newCheckCmd`), each returning a `*cobra.Command` whose
`RunE` loads config via `internal/config.Load` and then returns `nil` —
silent success, no incidental output, per the acceptance criteria and
Scenario 01. `propose` additionally declares a `--since` string flag
(no default value forced here beyond what config resolution supplies).

**Contract**: Each `RunE` signature is
`func(cmd *cobra.Command, args []string) error`. `propose`'s flag:
`cmd.Flags().String("since", "", "...")`.

#### 3. Shared error formatting

**File**: `main.go`

**Intent**: Implement D011 at the single point where `Execute()`'s error
return is handled — not duplicated per command.

**Contract**: On a non-nil error from `rootCmd.Execute()`, write
`fmt.Fprintf(deps.Stderr, "error: %s\n", err)` and exit with code 1; on
nil, exit 0.

### Success Criteria:

#### Automated Verification:

- Module builds: `go build ./...`
- `go vet ./...` passes
- Unit tests (command tree shape) pass: `go test ./internal/cli/...`

#### Manual Verification:

- `./changes --help` lists `init`, `propose`, `apply`, `check`
- `./changes propose --help` shows `--since <ref>`
- Running each subcommand in an empty temp directory exits 0 with no
  stdout/stderr output
- `go run . nonexistent-command` exits 1 with stderr matching
  `error: ...` and empty stdout

**Implementation Note**: After completing this phase and all automated
verification passes, pause here for manual confirmation from the human
that the manual testing was successful before proceeding to the next
phase.

---

## Phase 5: Smoke Tests

### Overview

Cover the three items named in the task's acceptance criteria: command
discovery, default config loading, and exit behavior — using the fakes
built in Phase 2.

### Changes Required:

#### 1. Command discovery tests

**File**: `internal/cli/root_test.go` (new)

**Intent**: Assert the root command exposes exactly `init`, `propose`,
`apply`, `check`, and that `propose` has a `since` flag.

**Contract**: Table-driven test iterating `rootCmd.Commands()` names;
separate assertion for `rootCmd.Find([]string{"propose"})` flag set
containing `since`.

#### 2. Default config loading tests

**File**: `internal/config/config_test.go` (new, if not already covered
in Phase 3)

**Intent**: Already specified under Phase 3's Automated Verification —
this phase just confirms these tests exist and pass as part of the full
suite; no new test content required if Phase 3 delivered them.

**Contract**: N/A — verification only.

#### 3. Exit behavior tests

**File**: `internal/cli/root_test.go`

**Intent**: For each of the four subcommands, execute the root command
(via `Deps` built with `app.NewFakeDeps`) in an isolated in-memory
filesystem with no `.changes/config.yaml`, and assert exit is
success (`Execute()` returns nil) with both fake stdout and stderr
buffers empty.

**Contract**: Use `rootCmd.SetArgs([]string{"init"})` (etc.) then
`rootCmd.Execute()`; assert `err == nil` and buffer contents are empty.

### Success Criteria:

#### Automated Verification:

- Full suite passes: `go test ./...`
- `go vet ./...` passes

#### Manual Verification:

- None beyond automated — this phase is test-only

---

## Testing Strategy

### Unit Tests:

- Config default resolution (no file, partial file, CLI override wins)
- Command tree shape (`init`/`propose`/`apply`/`check` present, `propose`
  has `--since`)
- Stub exit behavior (exit 0, empty stdout/stderr) for all four
  subcommands using fake `Deps`

### Integration Tests:

- None at this stage — no real business logic exists yet to integrate
  against. In-process cobra execution (via `Execute()` against fake
  `Deps`) is sufficient per the "Test scope" decision made during
  planning.

### Manual Testing Steps:

1. `go build -o changes .` then `./changes --help` — confirm all four
   subcommands and their one-line descriptions appear.
2. `./changes propose --help` — confirm `--since` is listed.
3. In an empty temp directory (no `.changes/`), run each of
   `./changes init`, `./changes propose`, `./changes apply`,
   `./changes check` — confirm each exits 0 with no output.
4. Run `./changes badcommand` — confirm exit 1 and stderr
   `error: ...` with empty stdout.

## Performance Considerations

None — this is a CLI skeleton with no I/O-heavy or hot-path logic.

## Migration Notes

Not applicable — greenfield module, no existing users or data to migrate.

## References

- Task: `tasks/0001-bootstrap-go-cli.md`
- Design: `docs/DESIGN.md` (Configuration section, lines 113-126)
- Decisions: `docs/spec/DECISIONS.md` (D010, D011)
- Scenario: `docs/spec/scenarios/01-single-repo-happy-path.md`
- Downstream dependents: `tasks/0002-change-file-storage.md`,
  `tasks/0003-git-history-inspection.md`,
  `tasks/0004-changes-init.md`

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Module & Dependency Setup

#### Automated

- [x] 1.1 Module builds: `go build ./...`
- [x] 1.2 Module is tidy: `go mod tidy` produces no diff

#### Manual

- [x] 1.3 `go run . --help` runs without error

### Phase 2: Deps/App Scaffolding

#### Automated

- [ ] 2.1 Package builds: `go build ./internal/app/...`
- [ ] 2.2 Unit test confirms `NewFakeDeps` produces isolated buffers: `go test ./internal/app/...`

### Phase 3: Config Resolution

#### Automated

- [ ] 3.1 Package builds: `go build ./internal/config/...`
- [ ] 3.2 Unit tests pass (defaults, partial config, CLI override): `go test ./internal/config/...`

### Phase 4: Cobra Command Wiring

#### Automated

- [ ] 4.1 Module builds: `go build ./...`
- [ ] 4.2 `go vet ./...` passes
- [ ] 4.3 Command tree unit tests pass: `go test ./internal/cli/...`

#### Manual

- [ ] 4.4 `./changes --help` lists `init`, `propose`, `apply`, `check`
- [ ] 4.5 `./changes propose --help` shows `--since <ref>`
- [ ] 4.6 Each subcommand exits 0 with no output in an empty temp directory
- [ ] 4.7 `go run . nonexistent-command` exits 1 with `error: ...` on stderr and empty stdout

### Phase 5: Smoke Tests

#### Automated

- [ ] 5.1 Full suite passes: `go test ./...`
- [ ] 5.2 `go vet ./...` passes
