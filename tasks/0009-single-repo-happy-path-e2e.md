---
title: "Add the single-repo happy-path end-to-end test"
id: "0009"
status: pending
priority: high
effort: large
type: feature
tags: ["e2e", "scenario", "git"]
dependencies: ["0004", "0005", "0007", "0008"]
created_at: "2026-08-04"
---

# Add the single-repo happy-path end-to-end test

## Objective

Turn `docs/spec/scenarios/01-single-repo-happy-path.md` into an executable
black-box acceptance test covering adoption, proposal generation, manual edit,
release application, and final validation in a real temporary Git repository.

## Tasks

- [ ] Build the CLI once for the test suite and create an isolated repository
      with `main`, tag `v1.2.3`, and branch `feat/widget-api`.
- [ ] Create the two dated conventional commits from the scenario with stable
      identities and capture their actual short hashes.
- [ ] Run `changes init` under a fixed `2026-06-18T10:00:00Z` clock and assert
      its complete filesystem snapshot and streams.
- [ ] Run `changes propose` under a fixed `2026-06-19T10:00:00Z` clock and
      compare its generated YAML after substituting the fixture hashes.
- [ ] Replace the proposal with the exact user-edited YAML from Step 3.
- [ ] Run `changes apply` and compare changelog content, released-file bytes,
      and the complete filesystem layout.
- [ ] Run `changes check`, assert exact process output, and prove all tracked
      artifacts remain unchanged.
- [ ] Make the scenario test hermetic across supported operating systems.

## Acceptance Criteria

- One automated test executes all five Scenario 01 steps in order against a
  real temporary Git repository.
- Every specified exit code, stdout value, stderr value, file path, and file
  body is asserted at the step where it becomes observable.
- The test proves init stays in `.changes/`, the edited change moves unchanged
  to `.changes/released/`, and final `CHANGELOG.md` is byte-exact.
- The final `changes check` emits only `ok\n` and performs no writes.
- The test passes without network access, global Git identity, wall-clock
  dependence, or leakage into the source checkout.
