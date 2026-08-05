<!-- IMPL-REVIEW-REPORT -->
# Implementation Review: Bootstrap the Go CLI and default configuration

- **Plan**: context/changes/bootstrap-go-cli/plan.md
- **Scope**: All 5 phases (complete)
- **Date**: 2026-08-05
- **Verdict**: NEEDS ATTENTION
- **Findings**: 0 critical, 3 warnings, 4 observations

## Verdicts

| Dimension | Verdict |
|-----------|---------|
| Plan Adherence | WARNING |
| Scope Discipline | PASS |
| Safety & Quality | WARNING |
| Architecture | PASS |
| Pattern Consistency | PASS |
| Success Criteria | WARNING |

## Findings

### F1 — Partial `conventional` config replaces defaults instead of merging

- **Severity**: ⚠️ WARNING
- **Impact**: 🔎 MEDIUM — real tradeoff; pause to reason through it
- **Dimension**: Plan Adherence
- **Location**: internal/config/config.go:69-72
- **Detail**: Phase 3 requires "Config file present with a subset of keys → unset keys still fall back to defaults." Scalar keys behave correctly (`TestLoad_PartialConfigFallsBackToDefaults`), but a YAML file with only `conventional.major` replaces the entire default map — `minor`, `patch`, and `none` are lost. Verified: partial `conventional:\n  major:\n    - custom\n` yields a map with one key, not four.
- **Fix**: After `ReadConfig`, deep-merge file `conventional` keys onto `defaultConventional()` (per-key slice merge), then add `TestLoad_PartialConventionalFallsBackToDefaults`.
  - Strength: Matches D010 and the existing scalar partial-fallback behavior; prevents silent loss of increment mappings.
  - Tradeoff: Slightly more loader logic; must decide whether file slices replace or append per key.
  - Confidence: HIGH — verified with a repro test; plan criterion is explicit.
  - Blind spot: None significant.
- **Decision**: DEFERRED — task 0011 created (`tasks/cli/0011-fix-partial-conventional-config-not-merging-with-d.md`)

### F2 — Extra positional args silently accepted

- **Severity**: ⚠️ WARNING
- **Impact**: 🏃 LOW — quick decision; fix is obvious and narrowly scoped
- **Dimension**: Safety & Quality
- **Location**: internal/cli/init.go:11-18 (also propose.go, apply.go, check.go)
- **Detail**: All four stub subcommands omit `Args` validation. `changes init extra-arg` exits 0 with no output — typos and misuse go undetected.
- **Fix**: Set `Args: cobra.NoArgs` on all four subcommands; add a test asserting unknown args return an error.
- **Decision**: DEFERRED — task 0012 created (`tasks/cli/0012-reject-extra-positional-args-on-stub-subcommands.md`)

### F3 — Empty `--since` flag overrides default ref

- **Severity**: ⚠️ WARNING
- **Impact**: 🔎 MEDIUM — real tradeoff; pause to reason through it
- **Dimension**: Safety & Quality
- **Location**: internal/config/config.go:64-66
- **Detail**: When `flags.Changed("since")` is true and the value is `""`, the loader sets `Since` to empty string, bypassing the `"main"` default. Likely to cause confusing git/ref errors in downstream tasks.
- **Fix A ⭐ Recommended**: Treat empty string as unset — skip the override when `since == ""`.
  - Strength: Preserves default semantics; matches user intent when flag is declared but not given a value.
  - Tradeoff: Cannot explicitly set since to empty via CLI (unlikely use case).
  - Confidence: HIGH — common CLI convention.
  - Blind spot: None significant.
- **Fix B**: Validate non-empty `since` before returning `Config`; return a descriptive error.
  - Strength: Explicit failure rather than silent bad state.
  - Tradeoff: Changes error contract for a flag edge case.
  - Confidence: MED — depends on whether empty since is ever valid.
  - Blind spot: Downstream git adapter behavior not yet implemented.
- **Decision**: DEFERRED — task 0013 created (`tasks/cli/0013-treat-empty-since-flag-as-unset-in-config-loader.md`); recommended Fix A

### F4 — `NewFakeDeps` split into unplanned files

- **Severity**: 💡 OBSERVATION
- **Impact**: 🏃 LOW — quick decision; fix is obvious and narrowly scoped
- **Dimension**: Scope Discipline
- **Location**: internal/app/fake.go, internal/app/fake_test.go
- **Detail**: Plan expected `NewFakeDeps` in `deps.go` or `deps_test.go`. Implementation uses separate `fake.go` / `fake_test.go`. Behavior and signature match the plan exactly.
- **Fix**: No action required — acceptable file organization drift.
- **Decision**: ACCEPTED — benign file organization; no change needed

### F5 — Command discovery test doesn't assert exactly four subcommands

- **Severity**: 💡 OBSERVATION
- **Impact**: 🏃 LOW — quick decision; fix is obvious and narrowly scoped
- **Dimension**: Plan Adherence
- **Location**: internal/cli/root_test.go:11-25
- **Detail**: Plan says "exactly init, propose, apply, check." Test checks presence of all four but does not assert no extra subcommands are registered.
- **Fix**: Add `if len(root.Commands()) != 4 { t.Errorf(...) }` after the presence loop.
- **Decision**: DEFERRED — task 0014 created (`tasks/cli/0014-assert-exactly-four-subcommands-in-discovery-test.md`)

### F6 — No automated D011 failure-path test

- **Severity**: 💡 OBSERVATION
- **Impact**: 🏃 LOW — quick decision; fix is obvious and narrowly scoped
- **Dimension**: Success Criteria
- **Location**: internal/cli/root_test.go
- **Detail**: D011 error contract (`error: <message>` on stderr, exit 1, empty stdout) is implemented in `main.go` and verified manually (Progress 4.7), but no in-process test covers failure paths. Plan Phase 4.7 is manual-only.
- **Fix**: Add a test executing an unknown subcommand via fake deps and asserting stderr matches `^error: ` with empty stdout.
- **Decision**: DEFERRED — task 0015 created (`tasks/cli/0015-add-automated-d011-failure-path-test.md`)

### F7 — Malformed config file boundary untested

- **Severity**: 💡 OBSERVATION
- **Impact**: 🏃 LOW — quick decision; fix is obvious and narrowly scoped
- **Dimension**: Safety & Quality
- **Location**: internal/config/config_test.go
- **Detail**: `Load` returns an error from `v.ReadConfig` on invalid YAML, but no test asserts this boundary behavior.
- **Fix**: Add `TestLoad_InvalidConfigReturnsError` with invalid YAML content.
- **Decision**: DEFERRED — task 0016 created (`tasks/cli/0016-add-test-for-malformed-config-file-error.md`)
