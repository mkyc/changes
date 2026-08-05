---
id: "0010"
title: "create taskfile"
status: completed
priority: medium
dependencies: []
tags: []
created_at: 2026-08-05
---

# create taskfile

## Objective

I need taskfile that wraps nicely at least:
 - build
 - test
 - doctor
 - cleanup

## Acceptance Criteria

- run `task build` to produce binary
- run `task test` to execute tests
- run `task doctor` to execute all kinds of checks (formatting, static analysis, tidy, etc.)
- run `task cleanup` to remove artifacts
