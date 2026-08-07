# Accept YAML aliases for conventional config mapping Implementation Plan

## Overview

The 0019 change added a mapping-kind guard in `fileConventionalFromYAML` that rejects any `conventional:` value whose `yaml.Node.Kind` isn't `MappingNode`. It didn't account for YAML aliases: `conventional: *conv` has `Kind == AliasNode`, so the guard rejects it before `Decode` ever gets a chance — even though `Decode` already resolves aliases transparently and would succeed. This plan resolves the alias chain to determine the *effective* kind before running the guard.

## Current State Analysis

`fileConventionalFromYAML` (config.go:56-78) does, in order: absent-check (`Kind == 0`), non-mapping guard (`Kind != MappingNode && Tag != "!!null"`, added by 921cbd3 / task 0019), then `Decode` into `map[string][]string`, then an empty-mapping check. The non-mapping guard operates on the raw node without resolving `AliasNode` first, so any alias — even one pointing at a valid mapping — gets rejected by the guard's error message instead of reaching `Decode`.

Confirmed via spike: for `conventional: *conv` where `conv` anchors a mapping, `doc.Conventional.Kind == 16` (`AliasNode`) and `doc.Conventional.Alias` points to the resolved node (`Kind == 4`, `MappingNode`). `doc.Conventional.Decode(&fileConv)` called directly on the alias node already resolves correctly today (`map[major:[feat]]`) — the fix only needs to change what the *guard* inspects, not the `Decode` call itself.

## Desired End State

A `conventional:` value that's a YAML alias resolving to a mapping loads successfully, with the mapping's values used as before (merged with defaults per existing `mergeConventional` behavior). A `conventional:` value that's an alias resolving to a non-mapping (sequence, scalar) still returns `"conventional must be a mapping of keys to lists"`. Bare `conventional:` / `!!null` still hits the empty-key path unchanged. Verified by `go test ./internal/config/...` passing, including two new tests: alias-to-mapping (succeeds) and alias-to-sequence (still errors).

### Key Discoveries:

- `internal/config/config.go:63-68` — exact location of the guard; the fix inserts alias resolution between the absent-check and the guard's kind comparison.
- `go.yaml.in/yaml/v3`'s `yaml.Node` has an `Alias *Node` field, populated when `Kind == yaml.AliasNode`, pointing to the resolved target node. Chained aliases (alias-to-alias) resolve the same way — following `.Alias` repeatedly until a non-alias node is reached is safe, since the YAML parser itself rejects alias cycles.
- `internal/config/config_test.go` — the existing `TestLoad_SequenceConventionalReturnsError` / `TestLoad_ScalarConventionalReturnsError` (added in 921cbd3) establish the subtest pattern: `writeConfigFile`, call `config.Load`, assert on `err.Error()` substring. New tests follow this pattern, using a YAML anchor/alias in the config string.

## What We're NOT Doing

- Not changing the `Decode` call itself — it already resolves aliases correctly and is untouched by this fix.
- Not adding alias support anywhere else in the config file (e.g., other top-level keys) — scope is limited to `conventional:` per the task.
- Not imposing an artificial alias-chain-depth limit — the resolution loop runs until a non-alias node is found, relying on the YAML parser's own cycle prevention.

## Implementation Approach

Introduce a small local loop in `fileConventionalFromYAML` that walks `doc.Conventional.Alias` while `Kind == yaml.AliasNode`, producing an "effective" node reference used only for the kind/tag guard. The absent-check and `Decode` call continue operating on `doc.Conventional` as before (aliases were never a problem there). Add two regression tests mirroring the existing non-mapping test pattern.

## Phase 1: Resolve aliases before the kind guard, with regression tests

### Overview

Fixes the guard to look at the alias-resolved node's kind, and adds tests proving both the fix (alias-to-mapping succeeds) and that the 0019 protection still holds (alias-to-sequence still errors).

### Changes Required:

#### 1. Alias resolution before the guard

**File**: `internal/config/config.go`

**Intent**: Determine the effective YAML node kind by following `.Alias` through any chain of alias nodes, then run the existing non-mapping guard against that effective node instead of the raw (possibly alias) node.

**Contract**: Between the existing absent-check (`if doc.Conventional.Kind == 0`) and the guard (`if doc.Conventional.Kind != yaml.MappingNode && ...`), add a loop that walks a local node pointer through `.Alias` while its `Kind == yaml.AliasNode`, and change the guard to check that resolved pointer's `Kind`/`Tag` instead of `doc.Conventional`'s directly. The subsequent `Decode` call keeps operating on `doc.Conventional` unchanged (it already resolves aliases).

#### 2. Regression tests

**File**: `internal/config/config_test.go`

**Intent**: Prove an alias resolving to a mapping loads successfully (the task's reported bug), and that an alias resolving to a non-mapping still produces the 0019 error message.

**Contract**: Add `TestLoad_AliasToMappingConventionalSucceeds`, writing a config with an anchor/alias pair (e.g. `"defs:\n  conv: &conv\n    major:\n      - feat\nconventional: *conv\n"`), asserting `err == nil` and `cfg.Conventional["major"]` equals `["feat"]` merged with defaults (via `assertConventional`, following the `TestLoad_PartialConventionalFallsBackToDefaults` pattern). Add `TestLoad_AliasToSequenceConventionalReturnsError`, writing a config with an anchor/alias pair pointing to a sequence (e.g. `"defs:\n  conv: &conv\n    - a\n    - b\nconventional: *conv\n"`), asserting `err != nil` and `strings.Contains(err.Error(), "conventional must be a mapping of keys to lists")`. Place both alongside the existing `TestLoad_SequenceConventionalReturnsError` / `TestLoad_ScalarConventionalReturnsError` tests.

### Success Criteria:

#### Automated Verification:

- Unit tests pass: `go test ./internal/config/...`
- Full test suite passes: `go test ./...`
- Vet passes: `go vet ./...`

#### Manual Verification:

- None — this is a pure bug fix with full automated coverage.

---

## Testing Strategy

### Unit Tests:

- Alias resolving to a mapping loads successfully and merges correctly with defaults.
- Alias resolving to a sequence still returns the non-mapping error.
- All existing tests (0019's sequence/scalar cases, empty-block, empty-mapping, unknown-key, absent, partial mapping, bare-null) continue to pass unchanged, confirming the alias resolution doesn't affect non-alias node handling.

### Manual Testing Steps:

None required — behavior is fully exercised by automated tests.

## Performance Considerations

None — resolving an alias chain is a handful of pointer dereferences on an already-parsed, small YAML node.

## Migration Notes

None — this restores previously-working behavior (aliases resolving to mappings worked before 921cbd3); no config schema change.

## References

- Task: `tasks/0020-accept-yaml-aliases-for-conventional-config.md`
- Regression introduced by: `921cbd3` (task 0019)
- Related change: `context/changes/improve-error-message-for-non-mapping-conventional/`
- Existing pattern: `internal/config/config_test.go` (TestLoad_SequenceConventionalReturnsError, TestLoad_ScalarConventionalReturnsError)

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Resolve aliases before the kind guard, with regression tests

#### Automated

- [ ] 1.1 Unit tests pass: `go test ./internal/config/...`
- [ ] 1.2 Full test suite passes: `go test ./...`
- [ ] 1.3 Vet passes: `go vet ./...`
