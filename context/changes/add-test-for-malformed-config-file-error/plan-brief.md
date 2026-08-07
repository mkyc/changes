# Add Test For Malformed Config File Error — Plan Brief

> Full plan: `context/changes/add-test-for-malformed-config-file-error/plan.md`

## What & Why

`config.Load` silently short-circuits with an error when `.changes/config.yaml` contains invalid YAML, but nothing in the test suite confirms this — the bootstrap-go-cli impl review (finding F7) flagged this error boundary as untested.

## Starting Point

`internal/config/config.go:117-121` calls `v.ReadConfig` on the config file's raw bytes and returns `nil, err` immediately if parsing fails — before the `conventional`-specific validation logic that already has its own error tests. `internal/config/config_test.go:46-54` already provides a `writeConfigFile` helper for writing arbitrary YAML into an in-memory fs, reused as-is.

## Desired End State

A new `TestLoad_InvalidConfigReturnsError` in `internal/config/config_test.go` writes syntactically broken YAML, calls `Load`, and asserts both a non-nil error and a nil `*Config`. `go test ./internal/config/...`, `task test`, and `task doctor` all pass.

## Key Decisions Made

| Decision | Choice | Why (1 sentence) |
| --- | --- | --- |
| Malformed YAML content | Unterminated/broken mapping (`changelog: [unterminated`) | Unambiguously invalid at the syntax level regardless of YAML library quirks |
| Assertion scope | Non-nil error + nil `*Config` | Matches acceptance criteria and guards against a half-populated Config leaking out alongside an error |

## Scope

**In scope:** One new test function in `internal/config/config_test.go`.

**Out of scope:** Changing `config.go`; asserting on viper's exact error message text; testing other malformed-input variants (valid YAML, wrong types).

## Architecture / Approach

Follows the existing `TestLoad_EmptyConventionalBlockReturnsError` pattern exactly: write bad content via the existing `writeConfigFile` fixture, call `config.Load`, assert on the returned error and config.

## Phases at a Glance

| Phase | What it delivers | Key risk |
| --- | --- | --- |
| 1. Add the malformed-config test | Automated coverage for the YAML-parse-failure path | None — pure test-only addition, zero production code touched |

**Prerequisites:** None.
**Estimated effort:** Single small addition, one session.

## Open Risks & Assumptions

- None identified — the failure path is a direct, well-isolated `return nil, err` in existing code.

## Success Criteria (Summary)

- `go test ./internal/config/...` includes a case proving malformed config does not silently succeed.
- `task test` and `task doctor` pass.
