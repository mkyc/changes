# Fix Partial Conventional Config Merge Implementation Plan

## Overview

Fix `config.Load` so a partial `conventional` section in `.changes/config.yaml`
deep-merges per key onto built-in defaults (D010), instead of Viper replacing
the entire default map. Add validation for unknown keys and an empty
`conventional` block, per planning decisions.

## Current State Analysis

`internal/config/config.go` sets `conventional` defaults via
`v.SetDefault("conventional", defaultConventional())`, then reads the config
file. Scalar keys (`changelog`, `since`, etc.) partial-fallback correctly via
Viper defaults. Nested `conventional` does not: when the file specifies any
`conventional` key, Viper replaces the whole map and only file keys survive.

`TestLoad_PartialConfigFallsBackToDefaults` covers scalar partial fallback but
does not assert `Conventional`. There is no test for partial conventional
merge.

### Key Discoveries:

- Bug location: `internal/config/config.go:69-72` — builds `Conventional`
  solely from `v.GetStringMapStringSlice("conventional")` with no merge step.
- Existing pattern to mirror: scalar partial fallback in
  `TestLoad_PartialConfigFallsBackToDefaults` (`config_test.go:61-88`).
- Allowed conventional keys are fixed in `docs/DESIGN.md:121-125`:
  `major`, `minor`, `patch`, `none`.
- No validation-error pattern exists yet in `internal/` — this change
  introduces the first config validation errors.

## Desired End State

`config.Load` returns a four-key `Conventional` map (plus any explicitly
set known keys) when the config file specifies a subset of conventional keys.
File values replace defaults per key (including empty slices). Unknown
conventional keys and an empty `conventional:` block produce clear errors.
All existing config tests remain green.

### Verification

- `go test ./internal/config/...` passes, including new edge-case tests.
- Manual: create `.changes/config.yaml` with only `conventional.major:
  [custom]` and confirm resolved config has `major: [custom]` plus default
  `minor`, `patch`, `none`.

## What We're NOT Doing

- Fixing unrelated config issues (empty `--since` flag — task 0013, extra
  positional args — task 0012).
- Adding CLI flag overrides for conventional keys (no flags exist today).
- Building a generic nested-map merge helper for future config keys.
- Changing documented default conventional mappings in `docs/DESIGN.md`.

## Implementation Approach

After `ReadConfig`, resolve `Conventional` in three steps:

1. Start from a copy of `defaultConventional()`.
2. If a config file was read, read the file's conventional map from Viper.
   - If the map is empty (explicit `conventional:` with no nested keys),
     return an error.
   - For each file key: reject if not one of `major`, `minor`, `patch`,
     `none`; otherwise replace the corresponding entry in the merged map
     (honor empty slices as explicit disable).
3. Assign the merged map to `Config.Conventional`.

When no config file is present, or the file omits `conventional` entirely,
Viper retains defaults and the overlay is a no-op.

## Phase 1: Deep-Merge Conventional Defaults

### Overview

Replace the naive map copy in `Load` with default-first merge logic and
validation for unknown keys and empty conventional blocks.

### Changes Required:

#### 1. Conventional merge and validation

**File**: `internal/config/config.go`

**Intent**: After layering the config file, merge file `conventional` keys
onto built-in defaults per key (slice replacement, not whole-map
replacement). Reject unknown keys and an empty `conventional:` block with
descriptive errors.

**Contract**: Extract allowed keys (`major`, `minor`, `patch`, `none`) as a
package-level set or helper. Replace lines 69–72 with logic that:
- clones `defaultConventional()` as the merge base;
- when the config file was read, reads `v.GetStringMapStringSlice("conventional")`;
- returns an error if the file was read and that map is empty;
- for each file entry, validates the key and assigns the slice to the merge
  base (empty slice is valid);
- assigns the result to `Config.Conventional`.

Error messages should name the offending key or state that `conventional`
must contain at least one key.

### Success Criteria:

#### Automated Verification:

- Package builds: `go build ./internal/config/...`
- Unit tests pass: `go test ./internal/config/...`

#### Manual Verification:

- None beyond automated

**Implementation Note**: After completing this phase and all automated
verification passes, pause here for manual confirmation from the human that
the manual testing was successful before proceeding to the next phase.

---

## Phase 2: Tests

### Overview

Add test coverage for partial conventional merge and the validation edge
cases decided during planning.

### Changes Required:

#### 1. Partial conventional fallback test

**File**: `internal/config/config_test.go`

**Intent**: Lock in the primary bug fix — partial `conventional.major` in
config file merges with defaults for the other three keys.

**Contract**: Add `TestLoad_PartialConventionalFallsBackToDefaults` mirroring
`TestLoad_PartialConfigFallsBackToDefaults`: write a config file with only
`conventional.major: [custom]`, call `Load`, assert `major == [custom]` and
`minor`, `patch`, `none` match `defaultConventional()`.

#### 2. Validation edge-case tests

**File**: `internal/config/config_test.go`

**Intent**: Cover the validation rules chosen during planning so regressions
are caught early.

**Contract**: Add tests that assert:
- `conventional:` with no nested keys returns an error from `Load`;
- an unknown key (e.g. `conventional.unknown: [x]`) returns an error;
- `conventional.major: []` resolves to an empty slice for `major` while
  other keys retain defaults.

Also extend `TestLoad_PartialConfigFallsBackToDefaults` to assert
`Conventional` still matches full defaults when the file omits `conventional`
entirely (guards against regressions on the scalar-only path).

### Success Criteria:

#### Automated Verification:

- Unit tests pass: `go test ./internal/config/...`
- New test `TestLoad_PartialConventionalFallsBackToDefaults` exists and passes
- Existing tests `TestLoad_Defaults`,
  `TestLoad_PartialConfigFallsBackToDefaults`, and
  `TestLoad_CLIFlagWinsOverConfigFile` remain green

#### Manual Verification:

- None beyond automated

---

## Testing Strategy

### Unit Tests:

- Partial conventional merge (primary bug)
- Empty `conventional:` block → error
- Unknown conventional key → error
- Empty slice for a known key → honored
- Scalar-only partial config still yields full conventional defaults

### Integration Tests:

- None — config package is self-contained

### Manual Testing Steps:

1. Create `.changes/config.yaml` with `conventional.major: [custom]` only.
2. Run a small Go snippet or test calling `config.Load` and inspect output.
3. Confirm four keys present with expected values.

## Performance Considerations

Negligible — one extra map copy and a handful of key lookups per `Load` call.

## Migration Notes

Configs with a partial `conventional` section that previously lost default
keys will now pick up defaults for unset keys — this is the intended fix.
Configs with unknown conventional keys or an empty `conventional:` block will
now error instead of silently producing unexpected mappings.

## References

- Task: `tasks/0011-fix-partial-conventional-config-not-merging-with-d.md`
- Impl review F1: `context/changes/bootstrap-go-cli/reviews/impl-review.md`
- D010: `docs/spec/DECISIONS.md:123-130`
- Default mappings: `docs/DESIGN.md:121-125`
- Bootstrap Phase 3 criterion: `context/changes/bootstrap-go-cli/plan.md:240-241`

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles.

### Phase 1: Deep-Merge Conventional Defaults

#### Automated

- [x] 1.1 Package builds: `go build ./internal/config/...`
- [x] 1.2 Unit tests pass: `go test ./internal/config/...`

#### Manual

- [ ] 1.3 None beyond automated

### Phase 2: Tests

#### Automated

- [ ] 2.1 Unit tests pass: `go test ./internal/config/...`
- [ ] 2.2 New test `TestLoad_PartialConventionalFallsBackToDefaults` exists and passes
- [ ] 2.3 Existing tests `TestLoad_Defaults`, `TestLoad_PartialConfigFallsBackToDefaults`, and `TestLoad_CLIFlagWinsOverConfigFile` remain green

#### Manual

- [ ] 2.4 None beyond automated
