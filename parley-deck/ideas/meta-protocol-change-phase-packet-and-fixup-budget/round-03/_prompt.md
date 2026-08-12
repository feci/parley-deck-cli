# Round 3 — convergence

Write `round-03/<your-agent-id>.md`. This round exists to close three things. If you have nothing to
add on an item, say "settled, no change" and move on — length is not a virtue here.

## Settled in round 2 — do not reopen

- **D2** — one standalone generator, exposed as `parley protocol packet`, called by the prompt
  builders. The builders never read `COOPERATION.md` themselves. @hermes-1 withdrew the
  read-in-the-runner phrasing; all four agree.
- **D3** — §15 is **load-bearing** in Phases 5 and 8, not on demand. @codex-1 reversed its own
  round-1 classification; @kimi-1 and @claude-1 concur. The verdict kernel (§15.1–§15.4, §15.7) is
  always present before an implementer authors a validation, resolution or completion claim.
- All round-1 locks (generated never committed; complete omission index; fail open to full protocol;
  three instruction paths; escalate never auto-close; no severity floor).

## C1 — D1, the deliberation fix-up cap. One number, this round.

Positions after round 2: @codex-1 **5** (withdrew its own 6), @kimi-1 **5**, @hermes-1 **6**
(adopted @codex-1's 6 in parallel, before its withdrawal was visible), @claude-1 no independent
position.

@hermes-1: you rejected 8 because it "nearly tripled the existing driver cap" and "we have no data
showing that 3 is too tight." @claude-1's round 2 puts the direct question to you — that reason
argues for a small number but does not separate 6 from 5, and @kimi-1's dataset says nothing has
ever closed in the 6–8 band, so the two escalate an identical set. **Adopt 5, or state what 6 buys
that 5 does not.** Either answer closes this; a restatement does not.

@codex-1 and @kimi-1: if @hermes-1 produces a reason 6 beats 5, engage it on the merits. Do not
settle it by noting that two of you already said 5 — §15.3 forbids exactly that.

## C2 — D4, pre-register ONE threshold set, before any run

The estimates are compatible (50–70%, ~0.5×, ~70%). The **thresholds are not**, and they must be
fixed before data exists:

| | Ship if | Refute if |
| --- | --- | --- |
| @codex-1 | ≥50% median saving in both phases AND zero obligation misses | <20% in either, or any seeded rule missed |
| @kimi-1 | packet median ≤ ~60% of full AND canary passed | otherwise |
| @hermes-1 | ≥2× speedup | <1.5× |

Converge on one row. Specifically decide, each with a yes/no:

1. **Is @kimi-1's canary a veto on its own?** A task whose correct execution requires a rule the
   packet omits — if the packet arm misses it, does the change fail regardless of the speed number?
2. **n per arm, and which phases.** @kimi-1 says n=3 already proved too small; @codex-1 proposes six
   paired runs per phase; name the number and the phases.
3. **The exact speed threshold**, stated once, in one unit (median wall clock ratio, packet arm over
   full arm). Convert your own number into that unit rather than restating it in yours.
4. **Who runs it, and against which source** — the live resolved protocol or a snapshot.

Whatever you agree becomes pre-registered: it is written into `FINAL.md` before the experiment runs,
and the experiment's result is reported against it unchanged.

## C3 — scope: the §4.0 audit

@hermes-1 found that cross-review rounds are also "unbounded" in text but bounded in code.
@claude-1 verified it goes further — PRIMARY, `internal/app/*.go` passes only `CrossReviewRounds`, so
`MaxFixupCycles: 3` and `MaxRounds: 4` are driver defaults on every run regardless of track — while
the table declares itself "the single authoritative per-track gate."

@claude-1's counter-proposal: do not patch the two cells we noticed. Audit **every** cell of the
§4.0 table against its enforcing code, publish the divergence list, and add a structural test that
fails when a per-track cell has no enforcing path.

Decide: **in scope for this idea, a named follow-up idea, or rejected.** If in scope, say what it
does to this idea's size and whether the packet work should ship first. If a follow-up, name it and
say what blocks on it. This is a scoping decision and it is the one that determines whether this
idea ships in one piece or two.

## C4 — anything that would make you refuse to sign a consensus

Say it now. Round 4 is a signoff round, and an objection raised there that could have been raised
here costs everyone a cycle.

## Constraints

- Read all four `round-02/*.md` files, including `claude-1.md`, which was written last and responds
  to all three of yours.
- Read-only outside your own `round-03/<agent-id>.md`. No git write commands.
- English only. Redact obvious secrets.
- §15.2: tag provenance. `PRIMARY` needs a stable locator or quoted command output; untagged reads
  as `RECALL`. "X satisfies Y" is itself a factual claim and needs the same tag and the same check.
- Quote other participants literally, or say explicitly that you condensed.
