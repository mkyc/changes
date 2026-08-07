# Normalize nested conventional decode errors to one line Implementation Plan

## Overview

When a `conventional:` mapping value has the wrong YAML type (e.g. `major: feat` instead of a list), `fileConventionalFromYAML` currently decodes the whole mapping in one `Decode` call, and the yaml library's resulting error is always multi-line (`"yaml: unmarshal errors:\n  line N: ..."`), even for a single bad key. Through `cli.Run`, this breaks the D011 contract that stderr must be exactly one line matching `error: <message>\n`. This plan replaces the whole-map decode with per-key iteration so the fix can both guarantee a single-line message and name the offending key.

## Current State Analysis

`fileConventionalFromYAML` (config.go:56-82) already has an alias-resolution step and a mapping-kind guard (from tasks 0019/0020) before decoding. The decode step itself (`doc.Conventional.Decode(&fileConv)`, config.go:75) hands the whole node to the yaml library at once; any single bad value produces the library's standard `*yaml.TypeError`-style multi-line message, which is returned unwrapped.

Confirmed via spike: even decoding a *single* value node directly (`valNode.Decode(&s)`) still returns the same `"yaml: unmarshal errors:\n  line N: cannot unmarshal ..."` two-line format — the yaml library never produces a single-line variant for a type mismatch. Wrapping or wrapping around the library's own message is not viable; the fix must construct its own message instead of surfacing the library's.

`cli.Run` (internal/cli/run.go:14-19) writes `fmt.Fprintf(deps.Stderr, "error: %s\n", err)` — a multi-line `err.Error()` therefore produces multi-line stderr, violating D011.

## Desired End State

A `conventional:` mapping with a wrong-type value (e.g. `major: feat`) produces the single-line error `conventional.major must be a list of strings`, naming the specific offending key. Through `cli.Run`, this surfaces as exactly one line of stderr matching `^error: `. Top-level non-mapping / empty-block messages from 0017/0019/0020 are unchanged. Verified by `go test ./internal/config/...` and `go test ./internal/cli/...` passing, including new tests for the nested-error case (config-level, no embedded newlines) and the `cli.Run` D011-contract case.

### Key Discoveries:

- `internal/config/config.go:74-77` — exact location of the whole-map `Decode` call to replace.
- `go.yaml.in/yaml/v3`'s `yaml.Node` mapping nodes expose pairs via `Content []*Node`, alternating key node then value node (`Content[i]` = key, `Content[i+1]` = value). Scalar key nodes carry their raw text in `.Value`, so the key name can be read directly without a `Decode` call.
- `internal/cli/root_test.go:104-118` (`TestRun_UnknownSubcommandFollowsD011Contract`) establishes the exact pattern for asserting the D011 contract via `cli.Run`: `app.NewFakeDeps(time.Now())` for isolated buffers/memfs, assert exit code, empty stdout, and `regexp.MustCompile(`^error: `)` against stderr.
- `internal/app/fake.go` — `NewFakeDeps` returns `deps.FS` as an `afero.NewMemMapFs()`, so a config file can be written to it the same way `internal/config/config_test.go`'s `writeConfigFile` helper does, using `config.ConfigFilePath`.
- The `resolved` node computed by the existing alias-resolution loop (config.go:66-69) is exactly the mapping node whose `.Content` should be iterated — for an alias, `doc.Conventional.Content` is empty but `resolved.Content` holds the real pairs.

## What We're NOT Doing

- Not changing the top-level non-mapping guard, empty-block validation, or alias resolution (0017/0019/0020) — those stay as-is and are covered by existing tests.
- Not attempting to name *every* bad key when multiple values are wrong — returning on the first bad key found (in mapping order) is sufficient per the task's "when practical" phrasing and keeps the fix simple.
- Not validating key names against `allowedConventionalKeys` in this function — that check already happens later in `mergeConventional` and is unaffected.

## Implementation Approach

Replace the single `doc.Conventional.Decode(&fileConv)` call with a loop over `resolved.Content` in pairs: read the key from `Content[i].Value`, attempt `Content[i+1].Decode(&val)` into `[]string`, and on the first decode error return `fmt.Errorf("conventional.%s must be a list of strings", key)`. On success, accumulate into the `fileConv` map as before. Add a config-level regression test asserting the error has no embedded newlines, and a `cli.Run` test asserting the full D011 contract (exit 1, empty stdout, single-line stderr).

## Phase 1: Per-key decode with single-line errors, plus regression tests

### Overview

Replaces the whole-map decode with per-key iteration that produces a single-line, key-named error, and adds tests proving both the config-level fix and the end-to-end D011 contract through `cli.Run`.

### Changes Required:

#### 1. Per-key decode in fileConventionalFromYAML

**File**: `internal/config/config.go`

**Intent**: Decode each key/value pair of the `conventional:` mapping individually instead of the whole map at once, so a type mismatch on any single value can be reported with a single-line message naming that key, instead of surfacing the yaml library's multi-line error.

**Contract**: Replace the block at config.go:74-77 (`var fileConv map[string][]string` / `doc.Conventional.Decode(&fileConv)` / error return) with a loop over `resolved.Content` in steps of 2: `key := resolved.Content[i].Value`; decode `resolved.Content[i+1]` into a `[]string`; on error, `return nil, true, fmt.Errorf("conventional.%s must be a list of strings", key)`; on success, set `fileConv[key] = val` (initializing `fileConv` as `make(map[string][]string, len(resolved.Content)/2)` before the loop). The subsequent `len(fileConv) == 0` empty-mapping check is unchanged.

#### 2. Config-level regression test

**File**: `internal/config/config_test.go`

**Intent**: Prove a nested wrong-type value produces a single-line, key-named error with no embedded newlines.

**Contract**: Add `TestLoad_NestedWrongTypeConventionalReturnsSingleLineError`, writing a config with `"conventional:\n  major: feat\n"`, asserting `err != nil`, `strings.Contains(err.Error(), "conventional.major must be a list of strings")`, and `!strings.Contains(err.Error(), "\n")`. Place alongside the existing non-mapping/alias tests in `config_test.go`.

#### 3. cli.Run D011-contract test

**File**: `internal/cli/run_test.go` (new file, following the `root_test.go` pattern)

**Intent**: Prove the fix holds end-to-end through `cli.Run`'s D011 error contract — exit code 1, empty stdout, and stderr matching exactly one line of `error: ...`.

**Contract**: Add `TestRun_InvalidNestedConventionalFollowsD011Contract` in a new `internal/cli/run_test.go` (package `cli_test`), following `TestRun_UnknownSubcommandFollowsD011Contract`'s structure: build deps via `app.NewFakeDeps(time.Now())`, write `"conventional:\n  major: feat\n"` to `deps.FS` at `config.ConfigFilePath` (creating the `.changes/` directory first, mirroring `config_test.go`'s `writeConfigFile` helper), call `cli.Run(deps, []string{"check"})` (or another stub subcommand that loads config), and assert exit code `1`, empty stdout, and stderr matches `^error: conventional\.major must be a list of strings\n$`.

### Success Criteria:

#### Automated Verification:

- Unit tests pass: `go test ./internal/config/...`
- Unit tests pass: `go test ./internal/cli/...`
- Full test suite passes: `go test ./...`
- Vet passes: `go vet ./...`

#### Manual Verification:

- None — this is a pure error-message fix with full automated coverage, including an end-to-end D011-contract test.

---

## Testing Strategy

### Unit Tests:

- Nested wrong-type conventional value (`major: feat`) returns a single-line error naming the key, with no embedded newlines.
- `cli.Run` with the same invalid config returns exit code 1, empty stdout, and exactly one line of stderr matching `error: conventional.major must be a list of strings`.
- All existing tests (0017/0019/0020's non-mapping, empty-block, alias, unknown-key, absent, partial-mapping, bare-null cases) continue to pass unchanged, confirming per-key iteration preserves prior behavior for well-typed mappings.

### Manual Testing Steps:

None required — behavior is fully exercised by automated tests, including the end-to-end CLI contract test.

## Performance Considerations

None — iterating a handful of mapping pairs instead of one batch `Decode` call is negligible for config-sized YAML.

## Migration Notes

None — this is a strictly additive error-message improvement; valid configs continue to decode identically (same resulting `map[string][]string`).

## References

- Task: `tasks/0021-normalize-nested-conventional-decode-errors.md`
- D011 contract: `internal/cli/run.go:10-19`
- Existing D011 test pattern: `internal/cli/root_test.go:104-118`
- Related changes: `context/changes/improve-error-message-for-non-mapping-conventional/`, `context/changes/accept-yaml-aliases-for-conventional-config/`

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles. See `references/progress-format.md`.

### Phase 1: Per-key decode with single-line errors, plus regression tests

#### Automated

- [ ] 1.1 Unit tests pass: `go test ./internal/config/...`
- [ ] 1.2 Unit tests pass: `go test ./internal/cli/...`
- [ ] 1.3 Full test suite passes: `go test ./...`
- [ ] 1.4 Vet passes: `go vet ./...`
