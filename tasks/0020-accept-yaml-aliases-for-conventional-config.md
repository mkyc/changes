---
title: "Accept YAML aliases for conventional config mapping"
id: "0020"
status: pending
priority: medium
type: bug
tags: ["config"]
created_at: "2026-08-08"
depends_on: ["0019"]
---

# Accept YAML aliases for conventional config mapping

## Objective

Restore support for YAML aliases that resolve to a mapping under
`conventional:`, which the 0019 node-kind guard now rejects before
`Decode`.

## Steps to Reproduce

1. Write a config that defines a mapping anchor and aliases it:

   ```yaml
   defs:
     conv: &conv
       major:
         - feat
   conventional: *conv
   ```

2. Load config (e.g. via `config.Load` or any stub subcommand that calls it).

## Expected Behavior

Config loads successfully; `conventional.major` is `["feat"]` (or merges
with defaults as usual).

## Actual Behavior

Load fails with `conventional must be a mapping of keys to lists` because
`doc.Conventional.Kind` is `AliasNode`, not `MappingNode`, and the 0019
guard rejects it before `Decode`.

## Tasks

- [ ] Resolve aliases (follow `Node.Alias`) before the mapping-kind check in
      `fileConventionalFromYAML`, or only emit the friendly non-mapping error
      when `Decode` fails for a non-mapping shape.
- [ ] Add a regression test for `conventional: *anchor` that expects a
      successful load (or a merge-equivalent result).

## Acceptance Criteria

- [ ] Valid alias-to-mapping `conventional:` configs load without error.
- [ ] Sequence/scalar non-mapping cases from 0019 still return
      `conventional must be a mapping of keys to lists`.
- [ ] Bare `conventional:` / `!!null` still hits the empty-key path.
- [ ] `go test ./internal/config/...` passes, including the new alias test.

## References

- Branch review finding F1 (tasks 0011–0019 on `fixes`)
- Introduced by: `921cbd3` (task 0019)
- Location: `internal/config/config.go` (~line 66)
