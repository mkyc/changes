# Add Test For Malformed Config File Error Implementation Plan

## Overview

Add `TestLoad_InvalidConfigReturnsError` to `internal/config/config_test.go` to cover the error boundary when `.changes/config.yaml` contains syntactically invalid YAML, closing the gap flagged in the bootstrap-go-cli impl review (F7).

## Current State Analysis

`config.Load` (`internal/config/config.go:102-145`) reads `.changes/config.yaml` via `afero.ReadFile`, and if present, calls `v.ReadConfig(bytes.NewReader(data))` (line 118). If viper's YAML parser fails, `Load` returns `nil, err` immediately (line 119-120) — before ever reaching the `conventional`-specific validation logic that already has dedicated error tests (`TestLoad_EmptyConventionalBlockReturnsError`, `TestLoad_UnknownConventionalKeyReturnsError`). This top-level parse-failure path is currently untested.

`internal/config/config_test.go:46-54` already has a `writeConfigFile(t, fsys, yaml string)` helper that writes arbitrary YAML content into an in-memory `afero.Fs` at `config.ConfigFilePath` — this is the exact fixture the new test needs, no new test infrastructure required.

## Desired End State

`internal/config/config_test.go` includes `TestLoad_InvalidConfigReturnsError`, which writes syntactically broken YAML via the existing `writeConfigFile` helper, calls `config.Load`, and asserts both a non-nil error and a nil `*Config`. `go test ./internal/config/...` (and `task test` / `task doctor`) pass.

### Key Discoveries:

- `internal/config/config.go:117-121` — the exact failure path: `v.ReadConfig` error short-circuits `Load` with `return nil, err`.
- `internal/config/config_test.go:46-54` — `writeConfigFile` fixture, reused as-is.
- `internal/config/config_test.go:119-130` — `TestLoad_EmptyConventionalBlockReturnsError` is the closest existing pattern for an error-path test in this file (write invalid content, call `Load`, assert non-nil error).

## What We're NOT Doing

- Not asserting on the exact error message/wording from viper's YAML parser — only that `Load` returns an error, per the acceptance criteria and to avoid coupling the test to viper's internal error text.
- Not changing `config.go` — this is a test-only addition.
- Not testing other malformed-input variants (e.g., valid YAML with wrong types) — scope is limited to the syntactically-invalid-YAML case named in the task.

## Implementation Approach

Add one new test function to `internal/config/config_test.go`, following the existing error-path test pattern (`TestLoad_EmptyConventionalBlockReturnsError`): build an in-memory fs, write syntactically broken YAML via `writeConfigFile`, call `config.Load`, and assert `err != nil` and `cfg == nil`.

## Phase 1: Add the malformed-config test

### Overview

Add `TestLoad_InvalidConfigReturnsError` to `internal/config/config_test.go`.

### Changes Required:

#### 1. Malformed config test

**File**: `internal/config/config_test.go`

**Intent**: Cover the error boundary when `.changes/config.yaml` contains invalid YAML, verifying `Load` fails loudly instead of silently falling back to defaults.

**Contract**: New test function `TestLoad_InvalidConfigReturnsError`, placed alongside the other `TestLoad_*` error-path tests (e.g. after `TestLoad_BareNullConfigValueFallsBackToDefault`). Builds `fsys := afero.NewMemMapFs()`, writes broken YAML via `writeConfigFile(t, fsys, "changelog: [unterminated\n")`, calls `cfg, err := config.Load(fsys, newFlagSet())`, and asserts `err != nil` (via `t.Fatal` if nil) and `cfg != nil` is an error (i.e. asserts `cfg == nil`).

### Success Criteria:

#### Automated Verification:

- `task test` passes (runs `go test ./... -v`, includes `internal/config`)
- `task doctor` passes (gofmt, go vet, go mod tidy check, golangci-lint)

#### Manual Verification:

- None — this is a pure test-assertion change with no user-facing or runtime behavior.

**Implementation Note**: After completing this phase and automated verification passes, this plan is complete — there is no further phase.

---

## Testing Strategy

### Unit Tests:

- `TestLoad_InvalidConfigReturnsError` is the test being added; it directly exercises the `v.ReadConfig` failure path in `config.Load`.
- Existing `TestLoad_*` tests in `config_test.go` continue to pass unchanged — no production code is touched.

### Integration Tests:

- None needed; this is scoped to one new unit test.

### Manual Testing Steps:

- None.

## Performance Considerations

None — test-only change.

## Migration Notes

None.

## References

- Task: `tasks/0016-add-test-for-malformed-config-file-error.md`
- Impl review finding: `context/changes/bootstrap-go-cli/reviews/impl-review.md` (F7)
- Test location: `internal/config/config_test.go`
- `Load` parse-failure path: `internal/config/config.go:112-121`

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Add the malformed-config test

#### Automated

- [ ] 1.1 `task test` passes
- [ ] 1.2 `task doctor` passes
