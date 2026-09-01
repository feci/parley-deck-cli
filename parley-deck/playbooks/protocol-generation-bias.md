---
playbook: protocol-generation-bias
distilled-from: ideas/protocol-generation-bias
distilled: (set at commit)
status: advisory
---

## When to use

Work resembling **protocol-generation-bias** — see ideas/protocol-generation-bias for the concrete precedent. (Generalize this line:
what class of task does this playbook cover? Strip the idea-specific specifics.)

## Proven shape

- Track: deliberation
- Participants: [claude-1, codex-1, hermes-1, kimi-1, zcode-1]
- Fix-up cycles taken: 3
- Lifecycle actually run: round-01 → cross-review (if divergent) → consensus + signoffs
  → FINAL → implement → refutation review → fix-up → complete.

## Step checklist

- [ ] Frame the idea narrowly in 00-prompt.md; set the track deliberately.
- [ ] Independent round-01 from every participant before anyone reads the others.
- [ ] Open a cross-review round only where positions genuinely diverge.
- [ ] Draft consensus; collect every participant's signoff; then FINAL.
- [ ] Implement strictly to FINAL; record deviations in IMPLEMENTATION.md.
- [ ] Refutation-default review; fix-up until zero agreed findings; then complete.

## Gotchas & fixes

(Fill from this idea's review consensus + IMPLEMENTATION.md deviations — the mistakes
that were caught and how. Keep only the transferable ones.)

## Verification pattern

(How "done" was proven — the checks/evidence used. If a completion contract was used,
note the shape here.)
