---
title: "Bootstrap the Go CLI and default configuration"
id: "0001"
status: completed
priority: high
effort: medium
type: feature
tags: ["cli", "foundation", "config"]
created_at: "2026-08-04"
---

# Bootstrap the Go CLI and default configuration

## Objective

Create the Go application skeleton and command wiring needed to expose the
`changes init`, `changes propose`, `changes apply`, and `changes check`
workflows. Centralize default configuration so a single-repository project
works without `.changes/config.yaml`.

## Tasks

- [ ] Initialize the Go module and add an executable `changes` entry point.
- [ ] Add a root command with `init`, `propose`, `apply`, and `check` subcommands.
- [ ] Define configuration values for package path, changelog path, base ref,
      tag prefix, and conventional-commit increment mappings.
- [ ] Implement defaults of package `.`, changelog `CHANGELOG.md`, base ref
      `main`, and tag prefix `v` when no config file exists.
- [ ] Make command dependencies such as the clock, filesystem, Git runner, and
      standard streams injectable for deterministic tests.
- [ ] Add smoke tests for command discovery, default loading, and exit behavior.

## Acceptance Criteria

- `go build` produces a runnable `changes` binary.
- `changes --help` exposes the four scenario commands and `propose` accepts an
  optional `--since <ref>` flag.
- With no `.changes/config.yaml`, resolved configuration is package `.`,
  changelog `CHANGELOG.md`, base ref `main`, and tag prefix `v`.
- Successful commands can return exit code 0 without writing incidental output
  to stdout or stderr.
- Automated tests can supply fixed dates and isolated repositories without
  changing process-global state.
