# Round 2 — cross-review

Write `round-02/<your-agent-id>.md`. Address every other participant by name. Disagreement requires
a counter-proposal, not just an objection.

## Locked by round 1 — do NOT re-derive, do NOT re-litigate

All three round-1 files independently converged on these. Treat them as settled and build on them:

- The packet is **generated on demand from the single live `COOPERATION.md`, never committed.** A
  committed packet is a fourth stale copy of protocol text, and this deck has been bitten twice by
  exactly that.
- Every packet carries a **machine-generated, always-complete omission index**: every omitted block's
  stable locator, its classification, and the trigger that would require it. Inclusion may be
  curated; the index may not.
- **Any detector failure falls open to the complete resolved protocol**, labelled, with the reason.
- **Three instruction paths change together** — the skill's standing line, §9's session-start
  checklist, and the prompt templates. Text alone cannot govern a hand-written prompt; that limit is
  acknowledged, not solved.
- **At the fix-up cap: escalate, never auto-close. No severity floor** — fresh MAJORs at rounds 19-24
  make "late findings are trivial" empirically false in this deck.
- **The honest limit, stated by all three:** a rule missing from BOTH the packet and the omission
  index produces no in-loop signal, and under the Phase 2 silence bullet that reads as agreement.
  Detection exists only at generation time and ex post.

## Verified since round 1 — the facilitator's own PRIMARY

@hermes-1 and @kimi-1 each found the same protocol/code divergence independently. It reproduces:

```text
$ sed -n '229p' parley-deck/COOPERATION.md
| Fix-up (Phase 8) | cap 1 cycle; fix-only verification ok | cap 2 cycles; ... | unbounded; `strict_gate` available |

$ sed -n '150,153p' internal/track/track.go
	case Deliberation:
		// Deliberation == today's full lifecycle (backward-compat constraint), but
		// still subject to the non-solo floor checked above.
		return Policy{Track: Deliberation, ApplyOverrides: false, CrossReviewRounds: -1}, nil

$ sed -n '103,105p' internal/driver/driver.go
	if cfg.MaxFixupCycles <= 0 {
		cfg.MaxFixupCycles = 3
	}
```

So "unbounded" is true only for hand-driven runs; a driver-managed deliberation is silently capped
at **3** today. **This reframes rank 3.** It is not "add a cap where none exists" — it is "the text
and the tool already disagree, and the tool is stricter than the text." That is the same shape as the
finding that closed `protocol-read-cost-regression` rank 2. Say what follows from it.

## What round 2 must settle

**D1 — The deliberation fix-up cap. You gave three numbers: @kimi-1 5, @codex-1 6, @hermes-1 8.**
§15.3 forbids settling this by count. Each of you anchored differently — @kimi-1 on the measured 5.1
mean, @codex-1 on "just above the mean, below the tail", @hermes-1 on convergence-versus-cost. Either
adopt another participant's number and say what changed your mind, or defend yours against their
specific anchor. Whoever is right, the answer must survive the fact that the driver already enforces
3 and nobody noticed.

**D2 — @hermes-1's instruction-layer claim contradicts the constraint @hermes-1 established.**
Its summary puts the change in "the runner's prompt builders (`internal/runner/runner.go` and
`phase58.go`)", but the round-1 brief records @hermes-1's own PRIMARY finding that the Go runner
never reads `COOPERATION.md` — zero references in those exact files. @hermes-1: resolve this in your
own words. @codex-1 and @kimi-1: say whether the prompt-builder path can carry a packet at all, or
whether rank 1 lives entirely in the three text paths plus a new `parley protocol packet` command.

**D3 — §15 on-demand in Phases 5 and 8.** @codex-1 classifies §15 as an on-demand verdict trigger
there and then flags in its own risks that if agents routinely issue completion verdicts in those
phases, §15 must simply be load-bearing. This idea's own cycle-4 record is evidence: the
implementer's phase is exactly where self-verdicts get issued. Decide it.

**D4 — the missing number, and it is the one the owner actually asked for.** Nobody estimated how
much faster this makes anything. The established A/B is 3.3× median wall clock for *full protocol vs
none* — a packet sits somewhere between. Each of you: give an expected saving with your reasoning,
and give the measurement that would confirm or refute it **before** the change ships, not after. If
your honest answer is that the saving cannot be estimated in advance, say that and say what the
smallest experiment is that would produce the number.

**D5 — is Phase 8 even the right lever?** The measured 7.2× growth is in *review volume*, and §4.0
also lists cross-review rounds as `unbounded` for deliberation. A fix-up cap bounds one loop while
another stays open. Say whether rank 3 should cover both, or why bounding fix-up alone is sufficient.

## Recorded deviation

`claude-1` filed no `round-01` file for this idea. Its positions are in `00-prompt.md`, which is the
drafter's framing, not an independent Phase 1 analysis. This is recorded rather than hidden: the
round-1 record has three independent analyses, not four, and `claude-1` had read all three before
writing anything, so it cannot now supply an independent one. It participates from round 2.

## Constraints

- Read-only: no edits outside your own `round-02/<agent-id>.md`, no git write commands.
- English only. Redact obvious secrets.
- §15.2 provenance: tag every factual claim. `PRIMARY` needs a stable locator or quoted command
  output; **untagged reads as `RECALL`**. A compliance claim — "X satisfies Y" — is itself a factual
  claim and needs the same tag and the same check.
- Do not quote another participant non-literally. If you condense, say you condensed.
