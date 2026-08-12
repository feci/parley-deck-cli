---
idea: meta-protocol-change-phase-packet-and-fixup-budget
status: final
drafted-by: claude-1
date: 2026-08-11
participants: [claude-1, codex-1, hermes-1, kimi-1]
consensus: accepted — codex-1 ACCEPT, hermes-1 ACCEPT, kimi-1 ACCEPT
track: deliberation
rounds: 4
follow-up: meta-protocol-change-track-gate-enforcement-audit
---

# FINAL — phase-scoped protocol packet, and finite deliberation budgets

A §7 meta-protocol change. Four cross-review rounds, two consensus passes. Every participant
reversed at least one of its own positions; every reversal is recorded where it was made.

## Why

Measured in `protocol-read-cost-regression` and `speedup-tooling-evaluation`:

```
NOT the CLI: every command under a second
per call : reading COOPERATION.md in full costs 3.3x median wall clock (n=3/arm)
per idea : review rounds 1.6 -> 5.1 (max 24); review bytes 20,237 -> 146,290 (7.2x)
protocol : 720 -> 1,359 lines in ten weeks; MUST 22 -> 37
```

Cost of a round × number of rounds. The packet cuts the first; the budgets bound the second.

**Structural constraint (@hermes-1, PRIMARY, verified in three ideas):** the Go runner never reads
`COOPERATION.md`. The cost lives in **instructions** — the skill's standing line and hand-written
facilitator prompts. A change that does not touch the instruction layer touches nothing.

## 1. The packet

- **Generated on demand; never committed.** A committed packet is a stale copy of protocol text.
- **Rendered from the live resolved protocol only** — the source `parley protocol check` resolves —
  bound by `sourceSha256`. **No embedded, bundled or frozen snapshot is an admissible source.**
- **Complete omission index** on every packet: each omitted block's stable locator, classification,
  and the trigger that would require it. Inclusion may be curated; the index may not.
- **Fail open.** Parser failure, unknown phase/track/flag, source drift, unresolved dependency or
  hash mismatch → the complete protocol, recorded as `context-mode=full-fallback` with the reason.
  An unclassified normative block is included in every packet and fails `packet check`.
- **One shared renderer**, exposed as `parley protocol packet`, called by the prompt builders. The
  builders never read `COOPERATION.md` themselves.
- **Three instruction paths change together**: the skill's standing line, §9's session-start
  checklist, the prompt templates. Official launches require packet attestation; an unconditional
  "read all of COOPERATION.md" must be marked `full-fallback` with a reason.

**§15 is load-bearing in Phases 5 and 8.** The verdict kernel — §15.1–§15.4 and §15.7 — is present
before an implementer authors any validation, resolution or completion claim. @codex-1 reversed its
own round-1 classification:

> The reason is temporal: an on-demand rule cannot prevent an implementer from already having
> written “met,” “proved,” “resolved,” “verified,” or “complete” as a self-verdict.

**The honest limit, stated by all four:** a rule missing from **both** the packet and the omission
index produces no in-loop signal, and the Phase 2 silence bullet (`COOPERATION.md:350`) reads that
as agreement. Detection exists only at generation time and ex post. No generator can prove a
human-authored applicability annotation is semantically correct. Therefore: applicability changes
are themselves §7 changes, the conservative default is always-include, and the first release runs in
shadow/audit mode against a full-protocol comparison.

## 2. The budgets

| §4.0 cell, `deliberation` | Text said | Code did | Becomes |
| --- | --- | --- | --- |
| Fix-up (Phase 8) | unbounded | driver default **3** | **5 inclusive published cycles** |
| Cross-review (Phase 2) | unbounded | driver default **1** | **3 rounds after round 1** |

At either boundary: **blocking user escalation with a trajectory payload. Never auto-close.** No
severity floor — fresh MAJORs at rounds 19–24 make "late findings are trivial" false in this deck.
An extension is a recorded finite grant that **never resets the count**; silence never extends a
budget.

**Why 5, and how it was chosen.** @codex-1 withdrew its own 6 (it had compared *review rounds* to
*fix-up cycles* — different units). @hermes-1 withdrew 8, then withdrew the 6 it had adopted in
parallel, and re-ran the distribution itself:

```text
n=69 ideas, count of '^## Fix-up cycle' per IMPLEMENTATION.md, sorted:
0×17, 1×34, 2×7, 3×2, 4×3, 5×2, then 9, 14, 15, 25
```

Every value above 5 is in {9, 14, 15, 25}. **Nothing has ever closed in the 6–8 band**, so 5, 6, 7
and 8 escalate an identical set across the deck's whole history. No evidence separates them, so the
choice fell to error asymmetry: a too-low cap costs one recorded escalation a human can grant; a
too-high cap costs another cycle of the pathology the cap exists to interrupt. §15.3 was honoured —
the number was not chosen by how many held it.

**Binding implementation conditions:**

1. **Text and code in the same patch.** @kimi-1 measured a `standard` idea running **15 cycles
   against a printed cap of 2 with no recorded escalation**. Printed caps bind only where
   enforcement lives.
2. **Inclusive boundary, with a test.** Under the `>=` guard, naively setting `MaxFixupCycles: 5`
   delivers four published cycles. Boundary tests at 5/6 and 3/4.
3. `CapCrossReviewRounds` wired for `deliberation`; today only `standard` carries that path.

## 3. The pre-registered experiment

**This section is pre-registered. It is written here before the experiment exists, and the result is
reported against it unchanged.**

| Item | Value |
| --- | --- |
| Unit | `R = median(packet_wall_clock / full_wall_clock)` over the pairs, per phase |
| Phases | 1 and 6 — lightest and heaviest packets bracket the range |
| Runs | 6 paired runs per phase, counterbalanced AB/BA; agent, model/effort, task, output cap and workspace snapshot held constant; **packet generation time counted inside the packet arm** |
| Canary | a task whose correct execution needs a rule the packet omits (an `auto_implement` idea requiring §14); **3 packet-arm replicates, all must pass**, plus a full-arm control |
| Obligations | §6, §14, §15 obligations seeded and checked in **every** run of both phases |
| **Ship** | `R ≤ 0.50` in **both** phases AND canary passed AND zero obligation misses |
| **Refute** | **any correctness miss, at any speed** — plus the speed cut in the open conflict below |
| Middle band | **returns to the user with the measured number**, which replaces the planning estimate verbatim here and is **never rounded up to it** |
| Runner | the Phase 5 implementer, on the implementation branch, before release; a **non-implementer recomputes both ratios from the raw log** |
| Source | live resolved protocol, hash-bound, both arms |

**The correctness veto is standalone.** A packet arm that misses a required rule fails the change at
any speed. A packet that is 3× faster and misses one §14 obligation is not a win.

**Planning estimates, superseded by the measurement when it exists:** @codex-1 50–70% reduction
(2.0–3.3×); @hermes-1 ~70% (2.3×), with the caveat that the original A/B never recorded its excerpt
size; @kimi-1 ~0.5 ratio, explicitly not defended ("n=3 per arm, arm-A range 27.3–105.3, two-point
calibration, linearity assumed").

**Middle band — @hermes-1 withdrew its own alternative and said why:**

> My replan-and-re-run was post-hoc adjustment masquerading as rigor — "change the intervention and
> try again" is the optimization pre-registration exists to prevent.

### Open conflict — the refute speed cut

**Unresolved, recorded, and not resolved by count.** Neither side treats it as a blocker; all three
signoffs accept `FINAL.md` carrying both.

- **@codex-1 and @kimi-1 — refute if `R > 0.80` in either phase.** A measured 0.70 is a 30% cut on
  the protocol-read term: below the estimate but a large absolute saving on the heaviest phase.
  Refuting at 0.67 "refutes a measured 1.4–1.5× speedup without the owner ever seeing the number."
  Above 0.80 the saving is under 20%, too thin against the omission-risk surface.
- **@hermes-1 — refute if `R > 0.67` in either phase.** Below 1.5× the optimization does not justify
  a new generator, a new failure mode and a packet system; 0.80 (1.25×) is too loose, and "a middle
  band that wide spends the experiment's credibility on a range where the answer is nearly always
  'ship anyway.'"

The disputed region is exactly `(0.67, 0.80]`, and the two treatments differ in one respect: whether
a measured 1.25–1.5× saving reaches the owner or is auto-killed before anyone sees it. **The
implementer must not resolve this by picking one.** If the measurement lands in that band, it goes
to the user with both positions and the number.

## 4. Scope — two ideas

This idea ships the packet and the two budget cells. The general audit becomes
**`meta-protocol-change-track-gate-enforcement-audit`** (unanimous after @codex-1 returned to the
long slug on the merits: "enforcement" names the deliverable — an inventory of which cells have an
enforcing code path).

It is separated because it would **confound the experiment**, which measures the packet, and because
dispositioning `Timeout per agent`, `Reviewers (Phase 6)` and `Review consensus (Phase 7)` is the
same anchoring work the fix-up cap just cost two rounds — times N cells, with no known-correct
per-track values.

**Seed inventory for the follow-up** — divergences already found, PRIMARY:

```text
$ rg -n 'MaxFixupCycles|MaxRounds|CrossReviewRounds' internal/app/*.go
internal/app/driver_impl.go:315:	// ... rather than spinning fresh rounds to MaxFixupCycles.   (comment only)
internal/app/app.go:1209:		CrossReviewRounds: driver.ReadCrossReviewRounds(ideaDir),
internal/app/app.go:1941:				CrossReviewRounds: driver.ReadCrossReviewRounds(created.Idea.Path),
internal/app/app.go:1995:					CrossReviewRounds: driver.ReadCrossReviewRounds(created.Idea.Path),
```

The app layer passes **only** `CrossReviewRounds`. `MaxFixupCycles` and `MaxRounds` are never passed,
so the driver's defaults of 3 and 4 stand on every run regardless of track — while the §4.0 table
declares itself "the single authoritative per-track gate."

**A table that calls itself the single authoritative gate has cells the tool does not read.** That
is the same shape as the finding that closed rank 2 of `protocol-read-cost-regression`: the
normative text and the implementation disagreed, and the text lost silently.

Binding on the follow-up:

- **No further §4.0 cell edits land until the audit runs**, so the divergence list is published once
  and complete.
- **No claim that the §4.0 table is code-enforced may be made before the audit.**
- The audit delivers the enumeration, every cell's enforcing path or an explicit
  non-machine-enforceable disposition, the published divergence list, and a **structural test that
  fails when a per-track cell has no enforcing path** — the technique that caught the ledger leak in
  the sibling idea, where the behaviour was unobservable through output.

## 5. Scope guard — required, and it limits what may be claimed

This change cuts **only the protocol-read term.** The other term — re-reading prior rounds via
`gatherPriorRounds` / `gatherReviewContext` — was rank 2 of `protocol-read-cost-regression`, and it
was implemented and then **deleted** in v1.43.1 (`aecbc1c` added the frontier machinery, `41e6cd6`
removed it). It stands today at its full pre-idea cost.

**The whole-idea saving is therefore smaller than the per-call saving.** The per-call ratio is
reported as per-call and must never be presented as an idea-level number.

## Never cut

The authoritative protocol; applicable modals, negations, conditions and exceptions; §4.0 overrides
and invariants; the exact current phase block; the non-solo and files-canonical close guards; §6
rule 3, status re-read, English-only and no-secrets; the escalation mechanism; §14; the active
transport's current-phase mechanics; every close/cap/strict-gate condition applicable to the phase;
§15 for Phases 1, 2, 3, 6, 7 and the verdict kernel for 5 and 8; §7 in every phase of a
protocol-change idea; and on-demand access to raw historical artifacts.

## Recorded deviation

`claude-1` filed no `round-01`. Its positions are in `00-prompt.md` — the drafter's framing, not an
independent Phase 1 analysis — and it wrote its round-02 participant contribution after reading all
six prior files. **The round-1 record holds three independent analyses, not four.**

## How this was decided

Four rounds. @codex-1 withdrew 6 and its §15 classification; @hermes-1 withdrew 8, then 6, then its
bare speed threshold, then its replan band and its one-idea scope; @kimi-1 withdrew its n=5 design
and tightened its own ship threshold from 0.60 to 0.50 — against itself. @codex-1 blocked the first
consensus over three record errors by the drafter, all three of which were real.

The number that decided the cap came from a participant re-running the distribution across all 69
ideas rather than defending an anchor. That is what the process is for.
