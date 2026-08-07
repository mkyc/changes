# Improve error message for non-mapping conventional config value Implementation Plan

## Overview

`internal/config/config.go`'s `fileConventionalFromYAML` currently checks only whether `conventional:` is absent (`Kind == 0`) before calling `doc.Conventional.Decode(&fileConv)`. If the value is present but is a YAML sequence or scalar instead of a mapping (e.g. `conventional: [a, b]`), `Decode` fails with the underlying yaml library's error (e.g. "cannot unmarshal !!seq into map[string][]string"), which is less actionable than this file's own descriptive validation errors. This plan adds a node-kind guard that returns a message consistent with the file's existing style.

## Current State Analysis

`fileConventionalFromYAML` (config.go:56-75) already distinguishes "absent" (`Kind == 0`) from "present" before decoding, and after decoding it separately validates for an empty mapping. It has no check for "present but not a mapping" — that case falls through to `Decode`, whose error is returned as-is (config.go:68-70).

## Desired End State

A non-mapping `conventional:` value (sequence, scalar, or any other non-mapping YAML kind) produces the error `"conventional must be a mapping of keys to lists"`, returned before `Decode` is ever called. Verified by `go test ./internal/config/...` passing, including two new subtests for a sequence value and a scalar value.

### Key Discoveries:

- `internal/config/config.go:63-70` — the exact insertion point: right after the existing `if doc.Conventional.Kind == 0` absent-check, before the `Decode` call.
- `internal/config/config_test.go:119-156` — three existing subtests (`TestLoad_EmptyConventionalBlockReturnsError`, `TestLoad_EmptyMappingConventionalReturnsError`, `TestLoad_UnknownConventionalKeyReturnsError`) follow an identical pattern: write a config file via `writeConfigFile`, call `config.Load`, assert `err != nil`, assert `strings.Contains(err.Error(), "...")`. New tests follow this exact pattern.
- `go.yaml.in/yaml/v3` exposes `yaml.MappingNode` as the constant for mapping-kind nodes, importable in `config.go` (already imports `yaml "go.yaml.in/yaml/v3"`).

## What We're NOT Doing

- Not changing the empty-mapping validation (`len(fileConv) == 0` check, config.go:71-73) — that already has its own descriptive message and is unaffected by this change.
- Not echoing the actual YAML kind (e.g. "got sequence") in the error message — kept consistent with the file's existing terse style.
- Not adding special-case handling for `yaml.AliasNode` — a single `Kind != yaml.MappingNode` check covers all non-mapping shapes uniformly.

## Implementation Approach

Add one guard clause in `fileConventionalFromYAML`, immediately after the existing absent-check, that returns the new descriptive error when the node kind isn't a mapping. Add two new test subtests mirroring the existing three, one for a sequence value and one for a scalar value.

## Phase 1: Add node-kind guard and test coverage

### Overview

Adds the validation check to `config.go` and covers it with two new tests in `config_test.go`.

### Changes Required:

#### 1. Node-kind guard

**File**: `internal/config/config.go`

**Intent**: After confirming `conventional:` is present, reject any value that isn't a YAML mapping before attempting to decode it, returning a descriptive error consistent with the file's other validation messages.

**Contract**: Insert a check between the existing `if doc.Conventional.Kind == 0 { return nil, false, nil }` block and the `Decode` call: if `doc.Conventional.Kind != yaml.MappingNode`, return `nil, true, errors.New("conventional must be a mapping of keys to lists")`.

#### 2. Test coverage for non-mapping values

**File**: `internal/config/config_test.go`

**Intent**: Assert the new error message is returned for both a sequence and a scalar `conventional:` value, following the existing subtest pattern used for the empty-block and unknown-key cases.

**Contract**: Add `TestLoad_SequenceConventionalReturnsError` (config: `"conventional:\n  - a\n  - b\n"`) and `TestLoad_ScalarConventionalReturnsError` (config: `"conventional: scalar\n"`), each asserting `err != nil` and `strings.Contains(err.Error(), "conventional must be a mapping of keys to lists")`, placed alongside the existing conventional-validation subtests (after `TestLoad_UnknownConventionalKeyReturnsError`).

### Success Criteria:

#### Automated Verification:

- Unit tests pass: `go test ./internal/config/...`
- Full test suite passes: `go test ./...`
- Vet passes: `go vet ./...`

#### Manual Verification:

- None — this is a pure error-message change with full automated coverage.

---

## Testing Strategy

### Unit Tests:

- Sequence `conventional:` value returns the new descriptive error.
- Scalar `conventional:` value returns the new descriptive error.
- Existing tests (empty block, empty mapping, unknown key, partial mapping, absent block) continue to pass unchanged, confirming the new check doesn't regress mapping-kind handling.

### Manual Testing Steps:

None required — behavior is fully exercised by automated tests.

## Performance Considerations

None — a single additional comparison on an already-parsed YAML node.

## Migration Notes

None — this is a strictly additive validation improving an error message; no config schema or behavior change for valid inputs.

## References

- Task: `tasks/0019-improve-error-message-for-non-mapping-conventional.md`
- Source discussion: `context/changes/fix-partial-conventional-config-not-merging-with-defaults/reviews/impl-review.md` (impl review F3)
- Existing pattern: `internal/config/config_test.go:119-156`

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Add node-kind guard and test coverage

#### Automated

- [x] 1.1 Unit tests pass: `go test ./internal/config/...`
- [x] 1.2 Full test suite passes: `go test ./...`
- [x] 1.3 Vet passes: `go vet ./...`
