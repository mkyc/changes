---
change_id: tighten-d011-failure-path-test-single-line
title: Tighten D011 failure-path test for single-line stderr
status: implementing
created: 2026-08-08
updated: 2026-08-08
archived_at: null
---

## Notes

Address tasks/0023-tighten-d011-failure-path-test-single-line.md, but do not
update the existing TestRun_UnknownSubcommandFollowsD011Contract test — add
a new test instead that asserts the tightened D011 contract (non-empty
message, exactly one trailing newline, no embedded newlines).
