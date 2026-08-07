# Assert Exactly Four Subcommands in Discovery Test Implementation Plan

## Overview

Tighten `TestNewRootCmd_CommandDiscovery` in `internal/cli/root_test.go` so it fails not only when an expected subcommand (`init`/`propose`/`apply`/`check`) is missing, but also when an extra, unexpected subcommand is registered on the root command.

## Current State Analysis

`TestNewRootCmd_CommandDiscovery` (`internal/cli/root_test.go:11-26`) builds a `want` map of the four expected subcommand names, iterates `root.Commands()` marking each as found, then asserts every entry in `want` was found. It never checks the total count, so a fifth subcommand added to `NewRootCmd` (`internal/cli/root.go:20-23`) would pass this test silently — which is exactly the gap flagged in the bootstrap-go-cli impl review (F5).

`NewRootCmd` (`internal/cli/root.go:12-26`) explicitly registers exactly 4 commands (`init`, `propose`, `apply`, `check`) via `root.AddCommand(...)` and returns the built tree without calling `Execute()`, so `root.Commands()` reflects only those 4 — no cobra-injected `help` or `completion` command is present at this point.

## Desired End State

`TestNewRootCmd_CommandDiscovery` fails if `root.Commands()` does not have exactly 4 entries, in addition to the existing presence checks. `go test ./internal/cli/...` (and the broader `task test` / `task doctor`) pass.

### Key Discoveries:

- `internal/cli/root_test.go:16-20` — the presence loop to extend, insert the count check immediately after it (line 20).
- `internal/cli/root.go:20-23` — confirms exactly 4 `AddCommand` calls; no other commands are attached before `root.Commands()` would be read in this test.

## What We're NOT Doing

- Not restructuring the presence-loop into a sorted-slice/name-set comparison — the existing loop shape stays, only the length check is appended (per acceptance criteria and to keep this change minimal).
- Not touching any other test in `root_test.go`.
- Not changing `internal/cli/root.go`.

## Implementation Approach

Add a single `if len(root.Commands()) != 4 { t.Errorf(...) }` check right after the existing presence-verification loop, before the function returns. The error message includes both the actual count and the actual command names (via `root.Commands()`) so a future failure is self-diagnosing without needing to add debug prints.

## Phase 1: Add the count assertion

### Overview

Add the exact-count check to `TestNewRootCmd_CommandDiscovery`.

### Changes Required:

#### 1. Discovery test count assertion

**File**: `internal/cli/root_test.go`

**Intent**: After the existing loop that verifies all four expected subcommands are present (ending at line 20, `want[cmd.Name()] = true`), add a check that `root.Commands()` has exactly 4 entries, so any additional/unexpected subcommand registration causes this test to fail.

**Contract**: Insert between the presence loop (ends at the closing `}` of the `for _, cmd := range root.Commands()` loop, current line 20) and the existing `for name, found := range want` loop (current line 21). The failure message must report both the actual count and the actual command names, e.g.:

```go
if got := len(root.Commands()); got != 4 {
    t.Errorf("expected exactly 4 subcommands, got %d: %v", got, root.Commands())
}
```

(`%v` on `[]*cobra.Command` will print pointer-ish default formatting; use a name-extracting helper inline instead — build a `[]string` of `cmd.Name()` from `root.Commands()` and format that, so the message is actually readable.)

### Success Criteria:

#### Automated Verification:

- `task test` passes (runs `go test ./... -v`, includes `internal/cli`)
- `task doctor` passes (gofmt, go vet, go mod tidy check, golangci-lint)

#### Manual Verification:

- None — this is a pure test-assertion change with no user-facing or runtime behavior.

**Implementation Note**: After completing this phase and automated verification passes, this plan is complete — there is no further phase.

---

## Testing Strategy

### Unit Tests:

- `TestNewRootCmd_CommandDiscovery` itself is the test being tightened; verify it still passes against the current `NewRootCmd` (4 commands) and would fail (mentally/by temporary local check) if a 5th command were added — no need to add a throwaway 5th command to the codebase to prove this, the assertion logic is straightforward enough to review by inspection.

### Integration Tests:

- None needed; this is scoped to one existing unit test.

### Manual Testing Steps:

- None.

## Performance Considerations

None — test-only change.

## Migration Notes

None.

## References

- Task: `tasks/0014-assert-exactly-four-subcommands-in-discovery-test.md`
- Impl review finding: `context/changes/bootstrap-go-cli/reviews/impl-review.md` (F5)
- Test location: `internal/cli/root_test.go:11-26`
- Root command construction: `internal/cli/root.go:12-26`

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Add the count assertion

#### Automated

- [x] 1.1 `task test` passes
- [x] 1.2 `task doctor` passes
