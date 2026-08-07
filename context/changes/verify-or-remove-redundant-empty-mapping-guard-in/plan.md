# Verify Or Remove Redundant Empty-Mapping Guard Implementation Plan

## Overview

Resolve the either/or left open by impl review F1: remove the empty-mapping guard in `fileConventionalFromYAML` (`internal/config/config.go`), proven redundant by direct verification, and add a regression test for the `conventional: {}` form it was meant to catch.

## Current State Analysis

`fileConventionalFromYAML` (`internal/config/config.go:52-74`) has a guard at lines 62-64:

```go
if doc.Conventional.Kind == yaml.MappingNode && len(doc.Conventional.Content) == 0 {
    return nil, true, errors.New("conventional must contain at least one key")
}
```

This guard fires only for the literal `conventional: {}` form (an explicit empty YAML mapping). Direct verification against the `go.yaml.in/yaml/v3` decoder used by this file confirms:

- `conventional:` (bare, no value) parses to a null scalar node (`Kind=8`) — the guard's `Kind == yaml.MappingNode` condition is false, so it never fires for this form. This form is caught later by the post-decode `len(fileConv) == 0` check (lines 70-72).
- `conventional: {}` parses to a mapping node (`Kind=4`) with empty `Content` — the guard *does* fire for this form.
- Decoding that same empty-mapping node directly via `.Decode(&fileConv)` (i.e., what would happen if the guard were removed) produces `fileConv = map[string][]string{}` — non-nil, `len == 0` — which the existing post-decode `len(fileConv) == 0` check at lines 70-72 already catches, producing the identical `"conventional must contain at least one key"` error.

No existing test in `internal/config/config_test.go` exercises the `conventional: {}` form — `TestLoad_EmptyConventionalBlockReturnsError` (lines 119-130) only covers the bare `conventional:\n` form, which never touches the guard.

## Desired End State

The guard at `config.go:62-64` is removed; `fileConventionalFromYAML` relies solely on the post-decode `len(fileConv) == 0` check for both the bare and `{}` empty forms. A new test confirms `conventional: {}\n` still produces the `"conventional must contain at least one key"` error after the guard's removal. `go build ./internal/config/...` and `go test ./internal/config/...` pass with no behavior change.

### Key Discoveries:

- `internal/config/config.go:62-64` — the guard to remove.
- `internal/config/config.go:70-72` — the post-decode check that already covers both empty forms; verified experimentally to produce the same error for `{}` as the guard does.
- `internal/config/config_test.go:119-130` — `TestLoad_EmptyConventionalBlockReturnsError` is the existing pattern for this exact error message assertion (`strings.Contains(err.Error(), "conventional must contain at least one key")`), to be mirrored for the new `{}` case.

## What We're NOT Doing

- Not changing the post-decode `len(fileConv) == 0` check or any other validation logic in `config.go`.
- Not changing the error message text.
- Not adding coverage for other YAML node kinds (e.g., a sequence under `conventional:`) — scope is limited to the bare vs `{}` empty forms named in the task.

## Implementation Approach

Delete the guard block (`config.go:62-64`) since it is proven redundant. Add one new test to `internal/config/config_test.go`, mirroring `TestLoad_EmptyConventionalBlockReturnsError`, that writes `conventional: {}\n` and asserts `Load` returns an error containing `"conventional must contain at least one key"` — proving the post-decode path alone still handles this form correctly.

## Phase 1: Remove the guard and add the regression test

### Overview

Delete the redundant guard in `fileConventionalFromYAML` and add a test proving the `conventional: {}` form still errors correctly via the remaining post-decode check.

### Changes Required:

#### 1. Remove redundant guard

**File**: `internal/config/config.go`

**Intent**: Delete the empty-mapping guard (lines 62-64) since it's proven redundant — the post-decode `len(fileConv) == 0` check already produces the identical error for this input.

**Contract**: Remove the `if doc.Conventional.Kind == yaml.MappingNode && len(doc.Conventional.Content) == 0 { ... }` block from `fileConventionalFromYAML`. No other logic in the function changes; `doc.Conventional.Decode(&fileConv)` and the subsequent `len(fileConv) == 0` check remain as-is and now handle both the bare and `{}` forms.

#### 2. Regression test for the `{}` form

**File**: `internal/config/config_test.go`

**Intent**: Prove the `conventional: {}` form still produces the same error after the guard's removal, closing the "no existing test exercises this form" gap noted in impl review F1.

**Contract**: New test function (e.g. `TestLoad_EmptyMappingConventionalReturnsError`), placed alongside `TestLoad_EmptyConventionalBlockReturnsError`. Uses `writeConfigFile(t, fsys, "conventional: {}\n")`, calls `config.Load`, and asserts `err != nil` plus `strings.Contains(err.Error(), "conventional must contain at least one key")`, mirroring the existing test's structure exactly.

### Success Criteria:

#### Automated Verification:

- `go build ./internal/config/...` passes
- `task test` passes (runs `go test ./... -v`, includes `internal/config`)
- `task doctor` passes (gofmt, go vet, go mod tidy check, golangci-lint)

#### Manual Verification:

- None — this is a code-cleanup plus test-coverage change with no user-facing or runtime behavior change.

**Implementation Note**: After completing this phase and automated verification passes, this plan is complete — there is no further phase.

---

## Testing Strategy

### Unit Tests:

- New `TestLoad_EmptyMappingConventionalReturnsError` covers the `conventional: {}` form directly.
- Existing `TestLoad_EmptyConventionalBlockReturnsError` (bare form) and all other `TestLoad_*` tests continue to pass unchanged, confirming no regression from the guard's removal.

### Integration Tests:

- None needed; scope is limited to this one function and its existing test file.

### Manual Testing Steps:

- None.

## Performance Considerations

None — this is a dead-code removal plus one test.

## Migration Notes

None — no behavior change; both empty forms (`conventional:` and `conventional: {}`) continue to produce the same error as before.

## References

- Task: `tasks/0017-verify-or-remove-redundant-empty-mapping-guard-in.md`
- Impl review finding: `context/changes/fix-partial-conventional-config-not-merging-with-defaults/reviews/impl-review.md` (F1)
- Guard to remove: `internal/config/config.go:62-64`
- Post-decode check that already covers both forms: `internal/config/config.go:70-72`
- Existing test pattern to mirror: `internal/config/config_test.go:119-130` (`TestLoad_EmptyConventionalBlockReturnsError`)

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Remove the guard and add the regression test

#### Automated

- [x] 1.1 `go build ./internal/config/...` passes
- [x] 1.2 `task test` passes
- [x] 1.3 `task doctor` passes
