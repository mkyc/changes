# Treat Empty Since Values as Unset Implementation Plan

## Overview

Fix `config.Load` so exact-empty `since` values are treated as unset whether they come from an explicitly changed CLI flag or from `.changes/config.yaml`. Preserve the existing precedence of non-empty CLI values over non-empty file values over the built-in `main` default, and add focused regression coverage for every affected fallback layer.

## Current State Analysis

`internal/config/config.go` initializes Viper with `since: main`, reads the optional config file, then applies the CLI value whenever the `since` flag is marked changed. Pflag marks `--since=` and `--since ""` as changed even though their value is empty, so the current override replaces both a configured ref and the built-in default with `""`.

Viper also treats quoted YAML `since: ""` as an explicit config value, so it suppresses the default. Bare YAML `since:` already decodes as null and falls through to `main`. The existing tests cover defaults, partial config fallback, and a non-empty CLI value overriding a configured value, but they do not cover either empty-value source.

## Desired End State

`Config.Since` follows this resolution contract:

- An exact-empty CLI value is absent for precedence purposes.
- An exact-empty YAML value is absent for precedence purposes.
- A non-empty CLI value wins over every file state.
- A non-empty file value is used when the CLI value is absent or empty.
- The built-in `main` value is used when neither source provides a non-empty value.
- Whitespace-only values remain explicit values and are not normalized.

The behavior is verified entirely at the config-loader boundary because `propose` currently discards the resolved config and cannot expose its `Since` value through an in-process CLI or binary assertion.

### Key Discoveries:

- `internal/config/config.go:103-123` implements defaults → file → CLI layering; the current CLI condition checks `Changed("since")` but not whether the value is empty.
- `internal/config/config.go:131-136` materializes `Config.Since` from Viper after all layers have been applied, which is the observable contract under test.
- `internal/config/config_test.go:159-176` already protects the non-empty CLI-over-file rule and should remain unchanged and green.
- D005 defines `main` as the default ref, while D010 defines CLI > config file > built-in defaults (`docs/spec/DECISIONS.md:70-75,123-130`).
- Quoted YAML `since: ""` suppresses Viper's default; bare `since:` already falls through to the default.

## What We're NOT Doing

- Not trimming or otherwise normalizing whitespace-only refs.
- Not validating Git ref syntax or changing downstream Git error behavior.
- Not changing the behavior of a bare `--since` flag with no value; pflag rejects it before `config.Load` runs.
- Not changing config keys other than `since`.
- Not adding CLI, binary, or process-level integration tests.
- Not changing `propose` to consume or expose the resolved config as part of this bug fix.

## Implementation Approach

Keep the existing defaults → file → CLI loading structure. Prevent an exact-empty changed CLI flag from entering Viper's highest-precedence override layer, then ensure an exact-empty value remaining after file resolution falls back to the built-in `main` value. Preserve non-empty values verbatim at both sources.

## Critical Implementation Details

Ordering is load-bearing: an empty CLI value must be skipped rather than immediately replaced with `main`, because a non-empty file value must remain visible. Empty-file normalization must not run in a way that can override a later non-empty CLI value.

## Phase 1: Normalize Empty Since Values and Add Regression Coverage

### Overview

Update the loader's `since` resolution contract and add unit tests for the CLI and YAML empty-value paths while preserving established precedence for non-empty values.

### Changes Required:

#### 1. Empty-aware since resolution

**File**: `internal/config/config.go`

**Intent**: Treat exact-empty values as absent at both the CLI and config-file layers without disturbing the existing precedence of meaningful values.

**Contract**: `Load` must resolve `Config.Since` according to the following matrix: empty or unchanged CLI plus no meaningful file value resolves to `main`; empty or unchanged CLI plus a non-empty file value resolves to the file value; quoted empty YAML resolves to `main`; any non-empty CLI value overrides the file. Comparisons are exact and must not trim whitespace.

#### 2. Regression tests

**File**: `internal/config/config_test.go`

**Intent**: Lock both fallback layers and the expanded YAML behavior so future changes cannot replace precedence with an unconditional default assignment.

**Contract**: Add config-loader unit coverage proving:

- a changed empty CLI flag with no config file resolves `Since` to `main`;
- a changed empty CLI flag with `since: develop` resolves `Since` to `develop`;
- quoted YAML `since: ""` with no meaningful CLI override resolves `Since` to `main`;
- the existing non-empty CLI-over-config test remains green.

### Success Criteria:

#### Automated Verification:

- Config package tests pass: `go test ./internal/config/...`
- Full test suite passes: `task test`
- Formatting, vet, module tidiness, and lint checks pass: `task doctor`

#### Manual Verification:

- None beyond automated verification; the resolved value is fully observable at the config-loader unit-test boundary.

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the verification was successful before proceeding to close out the change.

---

## Testing Strategy

### Unit Tests:

- Exact-empty changed CLI value falls back to the built-in `main` default when the config file is absent.
- Exact-empty changed CLI value falls back to a non-empty config-file value.
- Quoted empty YAML value falls back to the built-in `main` default.
- Existing non-empty CLI value continues to override a non-empty config-file value.

### Integration Tests:

None. `Config.Since` is not currently consumed or exposed by `propose`, so CLI and binary tests cannot assert the resolved fallback value and would add no meaningful coverage.

### Manual Testing Steps:

None beyond reviewing the automated test results.

## Performance Considerations

Negligible. Resolution adds only constant-time empty-string checks during configuration loading.

## Migration Notes

Quoted YAML `since: ""` will now resolve to `main` instead of remaining empty. This is the intentionally expanded behavior chosen during planning. Non-empty CLI and file values are unchanged.

## References

- Task: `tasks/0013-treat-empty-since-flag-as-unset-in-config-loader.md`
- Originating finding: `context/changes/bootstrap-go-cli/reviews/impl-review.md` (F3)
- Follow-up register: `context/changes/bootstrap-go-cli/follow-ups/review-fixes.md`
- Config precedence decisions: `docs/spec/DECISIONS.md` (D005, D010)
- Loader: `internal/config/config.go:99-137`
- Existing config tests: `internal/config/config_test.go:56-176`

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Normalize Empty Since Values and Add Regression Coverage

#### Automated

- [ ] 1.1 Config package tests pass: `go test ./internal/config/...`
- [ ] 1.2 Full test suite passes: `task test`
- [ ] 1.3 Formatting, vet, module tidiness, and lint checks pass: `task doctor`

#### Manual

- [ ] 1.4 None beyond automated verification; review the automated test results
