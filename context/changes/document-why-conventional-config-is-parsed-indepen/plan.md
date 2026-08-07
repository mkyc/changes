# Document Why Conventional Config Is Parsed Independently Implementation Plan

## Overview

Add a doc comment to `fileConventionalFromYAML` explaining why the `conventional` block is parsed directly via `go.yaml.in/yaml/v3` instead of read from Viper, so a future reader doesn't mistake the duplicated YAML parsing for an oversight.

## Current State Analysis

`internal/config/config.go` parses the config file bytes twice: once by Viper for scalar keys (`v.ReadConfig`, line 115), and once directly via `go.yaml.in/yaml/v3` in `fileConventionalFromYAML` (line 52) for the `conventional` block specifically. This is a deliberate substitution for the plan's originally suggested `v.GetStringMapStringSlice` approach — Viper collapses "key absent" and "key present but empty" into the same observable state, which breaks the empty-block validation at lines 67-68 (`len(fileConv) == 0` → error). `fileConventionalFromYAML` needs the raw `yaml.Node.Kind` to tell these cases apart (verified during the sibling change `verify-or-remove-redundant-empty-mapping-guard-in`: bare `conventional:` parses to `Kind=8` null-scalar, `conventional: {}` parses to `Kind=4` mapping — a distinction Viper's typed getters don't expose). This reasoning is not currently documented anywhere in the file, per impl review F2.

## Desired End State

`fileConventionalFromYAML` (`internal/config/config.go:52`) has a doc comment explaining the rationale for its independent YAML parsing. No behavior changes; `go build ./internal/config/...` and `go test ./internal/config/...` pass exactly as before.

### Key Discoveries:

- `internal/config/config.go:52` — `fileConventionalFromYAML`, the function to document.
- `internal/config/config.go:115` — `v.ReadConfig`, the Viper parse path this function's approach deliberately bypasses for the `conventional` key.
- `internal/config/config.go:67-68` — the `len(fileConv) == 0` empty-block check that depends on distinguishing null from empty-mapping, which is the actual reason raw YAML parsing is needed.

## What We're NOT Doing

- Not changing any parsing logic or behavior — this is a comment-only change.
- Not adding comments elsewhere in the file (e.g., at the `v.ReadConfig` call site) — the task specifies one comment at `fileConventionalFromYAML`, confirmed as the placement of choice.

## Implementation Approach

Add a Go doc comment directly above `func fileConventionalFromYAML(...)` stating that it parses the config's raw YAML independently of Viper because Viper's typed getters (e.g. `GetStringMapStringSlice`) can't distinguish a key that's absent from one that's present but empty — a distinction the empty-block validation below requires.

## Phase 1: Add the doc comment

### Overview

Document the rationale for `fileConventionalFromYAML`'s independent YAML parsing.

### Changes Required:

#### 1. Doc comment on `fileConventionalFromYAML`

**File**: `internal/config/config.go`

**Intent**: Explain why this function re-parses the config's raw bytes with `go.yaml.in/yaml/v3` instead of reading `conventional` from the already-populated Viper instance, so the duplication reads as deliberate rather than an oversight.

**Contract**: Add a standard Go doc comment immediately above `func fileConventionalFromYAML(data []byte) (map[string][]string, bool, error) {` (currently line 52), e.g.: `// fileConventionalFromYAML parses the conventional block directly instead of` / `// reading it from v, because Viper's typed getters can't distinguish an` / `// absent key from one that's present but empty — exactly the distinction` / `// the empty-block validation below needs.`

### Success Criteria:

#### Automated Verification:

- `go build ./internal/config/...` passes
- `go test ./internal/config/...` passes (via `task test`)
- `task doctor` passes (gofmt, go vet, go mod tidy check, golangci-lint)

#### Manual Verification:

- None — this is a comment-only change with no behavior or output to manually verify.

**Implementation Note**: After completing this phase and automated verification passes, this plan is complete — there is no further phase.

---

## Testing Strategy

### Unit Tests:

- None added — this is a documentation-only change. Existing `internal/config` tests must continue to pass unchanged, confirming no behavior was altered.

### Integration Tests:

- None needed.

### Manual Testing Steps:

- None.

## Performance Considerations

None — comment-only change.

## Migration Notes

None.

## References

- Task: `tasks/0018-document-why-conventional-config-is-parsed-indepen.md`
- Impl review finding: `context/changes/fix-partial-conventional-config-not-merging-with-defaults/reviews/impl-review.md` (F2)
- Function to document: `internal/config/config.go:52`
- Related prior verification: `context/changes/verify-or-remove-redundant-empty-mapping-guard-in/plan.md` (confirmed the null-vs-empty-mapping YAML node distinction experimentally)

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Add the doc comment

#### Automated

- [ ] 1.1 `go build ./internal/config/...` passes
- [ ] 1.2 `task test` passes
- [ ] 1.3 `task doctor` passes
