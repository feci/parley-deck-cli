---
agent: claude-1
idea: completion-contracts-evidence-ledger
round: 2
date: 2026-07-04
responding-to: [codex-1/round-01, hermes-1/round-01, antigravity-1/round-01]
---

## Position changes since prior round

I move toward **hermes-1's minimalism** and withdraw the separate `review/evidence.md`
artifact and the new `done_when:` key. hermes is right that the primitives already
exist (`checks:`, the `## Validation evidence` section, the strict_gate veto pattern),
and a second ledger artifact would drift from IMPLEMENTATION.md. I keep codex-1's
structured-validation rigor and antigravity-1's safety concerns, but fold them into the
smaller surface.

## Responses to others

### @hermes-1
Agreed on all four points, adopted as the spine of the design: (1) extend `checks:` from
a scalar to an optional named list rather than adding `done_when:`; (2) the ledger IS the
existing `## Validation evidence` section, driver-populated when `checks:` is a list;
(3) one fail-closed veto in §8 mirroring strict_gate (~20 lines); (4) protocol text
scoped to "when `checks:` is a list" so existing ideas are byte-for-byte unchanged. Your
"is it just a nudge?" question answers itself: make the existing section machine-filled
instead of prose-filled. That IS the fix.

One place I'd push back: overwrite-per-cycle vs append. You argue git history is the
audit trail. I agree for simplicity — driver **overwrites** `## Validation evidence` with
the latest run each cycle; git history preserves prior cycles. No new append-only rule.

### @codex-1
Your structured reader is the right shape, but scoped down: keep `ReadCompletionChecks`
as a NEW helper (don't retrofit the scalar `readFrontmatterField`), returning the list
form only. Drop the rich matcher set (`stdout_contains`, `stdout_regex`, `path_exists`,
`expect_exit`) from v1 — hermes and antigravity both flag matchers as speculative. v1 =
`{name, command, expect: exit 0}`; richer expectations are expressible inside the command
(`grep -q`, `test -f`). Your "`done_when` authoritative when present" becomes "list-form
`checks:` authoritative"; a scalar `checks:` keeps today's exact behavior.

### @antigravity-1
Adopt two of your safety points as MUSTs, drop one: (MUST) secret-scrub + fixed-size
truncation of the recorded output (first+last N lines, ~4KB cap) — cheap, safe by
construction; (MUST-note) the flaky-test paradox — a non-deterministic check can wedge
completion, so the driver records the failing evidence and escalates via the existing
§14 human-brake / stopping-judgment path rather than looping forever (no auto-retry).
Drop `parley check-contract` pre-validation as separate scope — nice, but a follow-up;
the driver already surfaces a failing check at Phase 5/8.

## New concerns / questions

- **strict_gate interaction** (raised by hermes): state them independent — list-`checks:`
  is the necessary condition (checks passed); strict_gate adds the sufficient condition
  (clean full-scope review). A green evidence table does NOT satisfy strict_gate.
- **Digest vs truncation**: agree with hermes — no crypto digest; exit code + truncated,
  secret-scrubbed tail is enough. The command is authored in a reviewed artifact.

## Current proposal

Extend `checks:` to accept a named list; driver runs each criterion, writes a per-
criterion table (name, exit, duration, scrubbed+truncated tail) into `##  Validation
evidence` in IMPLEMENTATION.md each cycle; §8 gains one fail-closed veto: with list-form
`checks:`, `status: complete` requires the latest evidence entry ALL-PASS at current
HEAD. Protocol text: ~6-10 sentences in §5/§8 scoped to the list shape. No new artifact,
no new frontmatter key, absent/scalar `checks:` byte-for-byte unchanged.
