<!-- IMPL-REVIEW-REPORT -->
# Implementation Review: Fix Partial Conventional Config Merge Implementation Plan

- **Plan**: context/changes/fix-partial-conventional-config-not-merging-with-defaults/plan.md
- **Scope**: Phase 1 and Phase 2 (all phases, both complete)
- **Date**: 2026-08-05
- **Verdict**: APPROVED
- **Findings**: 0 critical, 0 warnings, 3 observations

## Verdicts

| Dimension | Verdict |
|-----------|---------|
| Plan Adherence | PASS |
| Scope Discipline | PASS |
| Safety & Quality | PASS |
| Architecture | PASS |
| Pattern Consistency | PASS |
| Success Criteria | PASS |

## Findings

### F1 — Redundant empty-mapping guard is only partially reachable

- **Severity**: OBSERVATION
- **Impact**: 🏃 LOW — quick decision; fix is obvious and narrowly scoped
- **Dimension**: Architecture
- **Location**: internal/config/config.go:59-60
- **Detail**: The guard `doc.Conventional.Kind == yaml.MappingNode && len(doc.Conventional.Content) == 0` was verified experimentally: `conventional:\n` (bare key, no value) parses as a null scalar node (`Kind=8`), not a mapping node, so this branch never fires for that input — the actual empty-block rejection for that case happens later via `len(fileConv) == 0` after `Decode` (line 67-68). The guard only matters for the literal `conventional: {}` form, which parses as `Kind=4` (MappingNode) with 0 content and is not exercised by any test.
- **Fix**: Either add a test case for `conventional: {}` to justify keeping the guard, or remove it and rely solely on the post-decode `len(fileConv) == 0` check, which already covers both forms correctly.
- **Decision**: TRACKED as tasks/0017-verify-or-remove-redundant-empty-mapping-guard-in.md

### F2 — Config file parsed independently by Viper and raw YAML

- **Severity**: OBSERVATION
- **Impact**: 🏃 LOW — quick decision; fix is obvious and narrowly scoped
- **Dimension**: Architecture
- **Location**: internal/config/config.go:53 (yaml.Unmarshal) vs config.go:115 (v.ReadConfig)
- **Detail**: The same config file bytes are parsed twice by two different YAML decoders — once by Viper for scalar keys, once directly via `go.yaml.in/yaml/v3` for the `conventional` block (deliberate substitution for the plan's suggested `v.GetStringMapStringSlice`, since Viper can't distinguish "key absent" from "key present but empty"). This is a sound, intentional design choice, not a bug, but the duplication isn't documented in a comment, so a future reader might assume it's accidental.
- **Fix**: Add a one-line comment at `fileConventionalFromYAML` noting why it parses raw YAML instead of reading from `v` (to preserve the null-vs-empty-vs-populated distinction Viper collapses).
- **Decision**: TRACKED as tasks/0018-document-why-conventional-config-is-parsed-indepen.md

### F3 — Non-mapping `conventional:` value yields a generic YAML error

- **Severity**: OBSERVATION
- **Impact**: 🏃 LOW — quick decision; fix is obvious and narrowly scoped
- **Dimension**: Safety & Quality
- **Location**: internal/config/config.go:64 (doc.Conventional.Decode)
- **Detail**: If `conventional:` is given a YAML sequence or scalar instead of a mapping (e.g. `conventional: [a, b]`), the code doesn't pre-validate the node kind before `Decode`; it returns the underlying library error ("cannot unmarshal !!seq into map[string][]string") rather than one of this file's own descriptive messages. Functionally correct (error is returned, no panic), just a less-actionable message than the two custom ones already in the file.
- **Fix**: Not required — out of scope of the planned validation rules (only "unknown key" and "empty block" were specified). Leave as-is unless a user reports confusion.
- **Decision**: TRACKED as tasks/0019-improve-error-message-for-non-mapping-conventional.md

