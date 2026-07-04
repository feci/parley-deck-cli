---
idea: completion-contracts-evidence-ledger
drafted-by: claude-1
date: 2026-07-04
track: deliberation
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Agreed decisions

Round-02 converged unanimously on hermes-1's minimalism (claude-1, codex-1, hermes-1,
antigravity-1 all accept). The proposal is NOT a new `done_when:` key + a separate
`review/evidence.md` artifact — it is a small extension of primitives that already exist.

1. **Extend `checks:`, do not add `done_when:`.** `checks:` in `00-prompt.md` accepts a
   scalar command (today's exact behavior) OR an optional named list of
   `{name, command, expect: exit 0}` criteria. The list form activates the completion
   contract. Never both shapes at once. Absent or scalar `checks:` is byte-for-byte
   today's behavior.

2. **The ledger IS `## Validation evidence` in IMPLEMENTATION.md.** No new artifact.
   When `checks:` is a list, the driver writes the per-criterion table
   (name, exit code, duration, scrubbed+truncated output tail) into that existing
   section, **overwriting each cycle** — git history preserves earlier cycles (no new
   append-only rule). Reviewers already receive IMPLEMENTATION.md as review input.

3. **One fail-closed veto in Phase 8.** Mirroring the strict_gate pattern: when
   `checks:` is a list, `status: complete` requires the latest driver run to be
   ALL-PASS at current HEAD. The veto can only fail a close claim, never auto-pass one.
   ~30 lines of Go reusing the existing RunChecks capture.

4. **v1 = exit-0 only.** No `stdout_contains` / regex / `path_exists` matchers — those
   are expressible inside the command (`grep -q`, `test -f`). Deferred.

5. **Safety (antigravity-1, adopted as MUSTs):** the recorded output is secret-scrubbed
   and truncated to a fixed cap (~100 lines / ~4KB). No crypto digest — exit code +
   scrubbed tail is enough, and the command lives in a reviewed artifact.

6. **Flaky/persistent failures** escalate via the existing §14 human-brake / stopping-
   judgment path — the driver records the failing evidence and stops, never auto-retries
   into an infinite loop. The Phase-8 failure message must be descriptive enough to fix
   a broken command string.

7. **strict_gate is independent.** List-`checks:` is the necessary condition (checks
   passed); strict_gate adds the sufficient condition (clean full-scope review). A green
   evidence table does NOT satisfy strict_gate.

8. **Implementation shape (codex-1):** `ReadChecksContract` (yaml.v3) in
   `internal/driver/transport.go` (absent/scalar → legacy; list → validated criteria);
   list-form loop in `internal/app/driver_impl.go`; one completion-path branch in
   `internal/driver/impl.go`. Protocol text in BOTH COOPERATION.md copies (LE-4 ~2
   sentences + Phase-5 template ~1 sentence + Phase-8 veto), scoped to the list shape.

## Deferred follow-ups

- `parley check-contract` pre-flight validation of commands (antigravity-1).
- Richer expectation matchers (stdout/regex/path_exists).

## Dismissed / non-goals

- No new `done_when:` frontmatter key. No separate `review/evidence.md` artifact.
- No append-only ledger rule (git history is the audit trail).
- No auto-close on green evidence (reviewer judgment stays).

## Signoffs

<!-- each participant appends its own block -->

### Signoff: claude-1 — 2026-07-04
Status: ✅ ACCEPT
Drafter. Round-02 converged unanimously on the minimalist design (extend checks: to a
list, ledger = ## Validation evidence, one Phase-8 veto, v1 exit-0, strict_gate independent).

### Signoff: codex-1 — 2026-07-04
Status: ✅ ACCEPT
I accept the minimalist list-form `checks:` design: extend existing `checks:`, keep the ledger in `IMPLEMENTATION.md` `## Validation evidence`, use exit-0 criteria only, and keep `strict_gate` independent.
My round-02 implementation plan stands: add `ReadChecksContract` with `yaml.v3`, run list criteria through the driver with scrubbed/truncated evidence, add the Phase-8 all-pass close veto, update both protocol copies, and test legacy scalar behavior plus list pass/fail and `strict_gate` independence.

### Signoff: hermes-1 — 2026-07-04
Status: ✅ ACCEPT
The adopted design is my minimalist position: extend the existing `checks:` primitive to a list, keep the ledger inside `IMPLEMENTATION.md` `## Validation evidence`, and add only one fail-closed Phase-8 veto — no new frontmatter key, no new artifact, exit-0 only in v1.

### Signoff: antigravity-1 — 2026-07-04
Status: ✅ ACCEPT
I accept the minimalist list-form `checks:` design and validation ledger, which ensures robust safety through secret scrubbing and explicit escalation for flaky failures.
