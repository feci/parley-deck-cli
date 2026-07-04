---
agent: codex-1
idea: completion-contracts-evidence-ledger
review-round: 1
date: 2026-07-04
reviewed-commit: e0f2b45
---

## Summary

The parser and list runner are close: absent/scalar `checks:` stays on the legacy dispatch path, syntactically valid malformed lists fail closed, list output is bounded, and `replaceSection` preserves following sections. The requested checks pass: `go build ./...`, `go vet ./internal/...`, `go test ./internal/driver ./internal/app ./internal/protocol`, and the explicit drift guard `go test ./internal/protocol -run TestEmbeddedDefaultMatchesLiveDeck -count=1`.

I found blocking issues: the Phase-8 close path does not actually call `RunChecks` before `Complete`, and the scrubber can persist common secrets in `IMPLEMENTATION.md`. Also, syntactically broken list-form frontmatter can fall through to legacy behavior instead of failing closed. Raw `COOPERATION.md` hashes differ only in drift-allowlisted header/provenance and roster zones; the drift guard is green.

## Findings

### [CRITICAL] List checks do not veto the zero-fixes completion path

`internal/driver/impl.go:201-263` handles `OutstandingAgreedFixes == 0` by applying strict gate, LE-11, goal check, and then `Complete`. There is no `RunChecks` call on this path. A fake impl with `checksOK: false` would still complete here once review consensus is ready; `TestPhaseReviewDraftsThenCompletes` already exercises completion with `fakeImpl.checksOK` left at its zero value and only asserts `complete`.

This violates FINAL.md's Phase-8 rule: with list-form `checks:`, `status: complete` requires the latest run ALL-PASS at current HEAD. The existing pre-review/post-fixup checks do not prove current HEAD at close and cannot veto a stale or newly failing criterion when review consensus has zero agreed fixes.

Suggested fix: add a completion-contract gate in the zero-fixes branch before `GoalCheck`/`Complete`, scoped to list-form checks so scalar/absent legacy behavior is not changed. Add a driver test where review consensus is ready with zero fixes, `RunChecks` returns false, and `ActionComplete` is blocked.

### [MAJOR] Evidence scrubber can leak Authorization bearer tokens

`internal/app/driver_checks.go:27` redacts a single token-shaped word after `authorization`, and `scrubAndTruncate` writes the resulting tail into the canonical evidence section at `internal/app/driver_checks.go:91-101`. For common output like `Authorization: Bearer real-token`, the regex consumes and redacts only `Authorization: Bearer`, leaving `real-token` in the evidence. It also does not cover common unlabelled tokens such as `sk-...`, `ghp_...`, or AWS access-key-shaped values.

The output is bounded, so this is not a raw unbounded dump, but it is not safe enough for the no-secret-leak requirement because leaked tails are written into `IMPLEMENTATION.md` and can be committed.

Suggested fix: either stop persisting command output text and record metadata/digests only, or harden scrubbing with targeted patterns for full Authorization header values, standalone bearer tokens, and common provider token shapes. Add tests for `Authorization: Bearer ...`, standalone `Bearer ...`, `sk-...`, `ghp_...`, and AWS-style keys.

### [MAJOR] YAML syntax errors in list-form `checks:` fall back to legacy

`internal/driver/checks.go:41-42` treats any YAML unmarshal error as `(nil, false, nil)`, which sends the caller to the legacy scalar/go-test fallback. That preserves old behavior for malformed legacy frontmatter, but it also means a syntactically broken list-form `checks:` can fail open instead of blocking as an invalid contract.

The valid-YAML malformed cases are covered (duplicate names, empty names, empty commands, empty list), but a broken list is still a malformed list from the author's perspective and should not be silently replaced by `go test ./...` or "nothing to check."

Suggested fix: when the raw frontmatter has a `checks:` key that appears to start block/list form, treat YAML parse failure as `isList=true` with an error. Keep scalar legacy fallback only for absent checks or a clear inline scalar. Add a test with a broken list under `checks:` and assert `RunChecks` returns false.

## Open questions

- The live and embedded `COOPERATION.md` files are not byte-identical because the drift guard deliberately normalizes project-specific header/provenance and roster zones. Should this idea require raw byte identity anyway, or is the existing normalized drift guard the intended definition?
