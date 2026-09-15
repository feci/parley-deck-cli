---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
reviewed-commit: 89a4305
artifact-kind: independent slice verification
not-a-signoff: true
---

# Independent verification of Kimi's retained-execution correction

This is Codex's independent source inspection and execution of Kimi's slice,
not whole-idea acceptance or Phase-6 consensus. The owner handoff is preserved.
At this checkpoint the slice remains unintegrated because the probes below fail.

## Independently executed checks

In `worktrees/evidence-first-kimi` at clean source commit 89a4305:

`TMPDIR=/private/tmp GIT_CEILING_DIRECTORIES=/private/tmp go test -count=1 -p 1 -json ./internal/evidence ./internal/app -run 'Evidence|Criterion|Envelope|TreeDigest|Typed|Verifier|Barrier|AttestExecution|HelperProcess'`

PASS: 42 terminal test-pass events, zero skips/failures. Evidence package 1.180s,
app 4.895s. `TestTreeDigestUnsupportedEntryFails` actually passed (0.07s) on
local storage; this resolves the owner's shared-volume SKIP as a coverage gap.
The portable serial/barrier fixture passed (2.67s). Full structured output is
retained in the worktree's ignored
`.parley-runtime/codex-independent-evidence-20260911.jsonl`.

Additional tests were supplied through Go's file overlay, without editing any
participant-owned source. Command:

`TMPDIR=/private/tmp GIT_CEILING_DIRECTORIES=/private/tmp go test -count=1 -overlay .parley-runtime/codex-evidence-probes-overlay.json ./internal/evidence -run '^TestCodexProbe' -v`

The overlay source and output are retained in `.parley-runtime/` as
`codex_evidence_probes_test.go` and `codex-independent-probes-20260911.log`.

## Confirmed corrections

- Old unrelated-verifier-command probe now refuses: PASS.
- Old duplicate `executed_cases:0` then `executed_cases:1` envelope probe now
  refuses: PASS. The strict envelope parser addresses this earlier finding.

## Remaining findings

### [MAJOR] Structured verifier skipped counts can be invalid while closure passes

`reconcileRerun` only compares skipped counts when both are nonnegative. A
Go-test verifier rerun with `SkippedCases` equal to -1, -2 or -100 is accepted
by `AttestExecution`, and `Evaluate` then returns no reasons. All three values
reproduced. Apply count validity to retained executions as well as original
records; retain -1 only where the envelope format really lacks the field.
Ensure format compatibility cannot bypass reconciliation.

### [MAJOR] Duplicate Go-test Action fields turn a failed event into passing evidence

`ParseGoTestJSON` accepts
`{"Action":"fail","Action":"pass","Test":"TestX"}` as one executed,
zero-failure test. The envelope parser is strict but the other supported
structured format still silently overwrites contradictory fields. Reject
duplicate/case-aliased semantic fields and malformed recognized event streams
without falling back to opaque shell PASS. Preserve supported real test2json
fields and normal noise handling; do not invent universal shell semantics.

### [MINOR] The capture boundary rejects valid output based on chunk size

`cappedWriter.Write` computes remaining capacity after adding the new bytes,
then compares the original chunk length against that diminished capacity.
With a 16-byte limit, writes of 16 bytes, two writes of 8 bytes, and one write
of 9 bytes all incorrectly set overflow. Compute overflow from the pre-write
remaining capacity. Cover exact-bound, split-bound and one-byte-over cases.

## Integration limit

Actual independent execution/attestation at the production close call site is
still absent and remains Codex-owned integration work. Metadata consistency
does not authenticate a human or establish that a caller actually ran a test.
The owner's causal explanation for earlier `signal: killed` output was not
independently established; the serial package run above only establishes its
own observed result.
