---
agent: codex-1
idea: meta-protocol-change-verification-integrity
round: 1
date: 2026-08-04
---

## Summary

The source run is absent from this deck, so I treat every reported incident, count, and
verdict as unverified testimony. My positions below do not depend on believing that testimony.
The structural problem is nevertheless real: the current protocol can show who authored a
statement, but it has no common contract for what a verification verdict means or what evidence
supports it.

I support a smaller, checkable verification contract rather than nine broad new rituals. It
should apply to factual claims deliberately placed into a verification pass, record claim owner,
verifier, method, locator, and evidence, and preserve unresolved material conflicts. Provenance is
metadata, not an automatic truth ranking. These rules compose with the ratified Phase 6
no-suppression rule: a prior verification disposition is context a reviewer may challenge, never
a reason not to report a finding. They also leave the optional strict gate and trajectory-based
stopping judgment unchanged.

Positions: CRITICAL-1, CRITICAL-2, CRITICAL-3, MAJOR-4, MAJOR-5, MAJOR-6, MAJOR-7, and MINOR-9
should be adopted only with the replacement text below. MINOR-8 is already covered and needs no
new rule.

## Proposed approach

### CRITICAL-1 — Adopt amended

The proposal correctly rejects self-certification, but “may only weaken” is too rigid. An owner
must be allowed to replace a malformed claim, supply primary evidence, or correct it in either
direction. What must remain forbidden is counting that owner action as independent verification.
“Claim ownership” also needs to cover joint authorship.

Replacement text:

> On every track, whenever a canonical artifact assigns a verification verdict to a factual
> claim, it MUST identify the claim by a stable identifier or exact quotation and record the
> claim owner and verifier. A claim owner MUST NOT supply the independent verdict for that claim.
> The owner MAY append an `OWNER-REVISION` that retracts, narrows, replaces, or supplies evidence
> for the claim, but owner action alone cannot raise its status above `UNVERIFIED`.
> `CONFIRMED` requires at least one non-owner verdict satisfying the provenance rule. If a claim
> has multiple substantive authors, each is an owner for this purpose.

This is enforceable by inspecting the named owner and verifier; it does not pretend that an owner
is incapable of finding useful evidence.

### CRITICAL-2 — Adopt amended

The proposed table improves on an unqualified `CONFIRMED`, but `SECONDARY` is circular when it
means “another agent agreed,” and venue/DOI language is too narrow for software, legal,
operational, or locally reproducible claims. A primary source can also be misread, so provenance
must not mechanically settle a conflict.

Replacement text:

> Every verification verdict MUST record `verifier`, `method`, `locator`, and a short
> `evidence/result`. Allowed methods are:
>
> - `DIRECT-CHECK`: a reproducible test, inspection, calculation, or derivation; include the
>   command or steps, inputs, and relevant result.
> - `PRIMARY-SOURCE`: an authoritative or original source was inspected; include a stable
>   locator and the section, identifier, or passage that supports the verdict.
> - `SECONDARY-SOURCE`: a named source independent of the claim owner was inspected; include its
>   locator and what it supports. Another agent's unsupported verdict is not a source.
> - `RECALL`: no source or direct check was consulted.
>
> A `RECALL` verdict, a missing locator, or missing evidence is capped at `UNVERIFIED`.
> `CONFIRMED` and `WRONG` require independently reviewable evidence. When no external primary
> source exists, a reproducible direct check or explicit derivation is admissible. `FINAL.md`
> MUST preserve any material `UNVERIFIED` or `DISPUTED` status rather than silently promoting it.

This records what was actually done without claiming that a tag proves correctness.

### CRITICAL-3 — Adopt amended

A conflict must be visible, but a new append-only `verdicts.md` is unnecessary ceremony and
creates a second place to reconcile with `consensus.md` and `FINAL.md`. “Higher provenance wins”
is unsound: a directly inspected source may be inapplicable or misinterpreted, while a lower-level
artifact may contain the decisive counterexample. Agent counting must not decide truth.

Replacement text:

> Any contradictory verdicts on the same identified claim MUST be recorded under
> `## Verification conflicts` in the next cross-review artifact or, if consensus is opening, in
> `consensus.md`. The entry MUST contain the claim identifier or exact text, every verdict with
> its verifier and evidence, the resolution, and the rationale. A material conflict MUST NOT be
> omitted merely because one side has more agents.
>
> Conflicts are resolved by reviewable evidence and argument, never by vote. Provenance methods
> are inputs to that judgment, not an automatic ordering. If no source, direct check, or
> derivation decides the issue, the claim remains `DISPUTED`. Consensus may close with a
> `DISPUTED` claim only when the agreed decision and acceptance criteria do not depend on that
> claim being true; `FINAL.md` MUST preserve the dispute and its impact. Otherwise the dispute
> blocks consensus or is escalated under the existing user-escalation rule.

The existing canonical files remain authoritative and artifact ownership is unchanged.

### MAJOR-4 — Adopt amended

The proposed witness requirement is useful, but a single counterexample is not generally enough
to prove that a proposal avoids an obstacle. The evidence must be logically sufficient for the
scope of the actual claim.

Replacement text:

> A canonical recommendation that claims to avoid a named known obstacle MUST identify the
> obstacle and provide at least one of: (a) an explicit mapping of the obstacle's relevant
> preconditions to the proposal, with evidence that a necessary precondition does not hold;
> (b) a reproducible witness or test that is logically sufficient for the scoped avoidance
> claim; or (c) a located authoritative result establishing the exemption. An adjective,
> analogy, or isolated example that does not entail the claimed scope is not a witness.
> Without sufficient evidence the artifact MUST label the assertion
> `OBSTACLE-CLAIM UNVERIFIED`, and the assertion MUST NOT be used as a reason to accept or
> implement the recommendation.

The facilitator can check that the mapping, witness, or source exists; participants still judge
whether it actually entails the claim.

### MAJOR-5 — Adopt amended

The anecdote-specific rule is disproportionate if applied to every software milestone or local
acceptance criterion. Most such items have no meaningful “proved/refuted/open” status. The valid
general rule concerns externally settled research, novelty, and feasibility claims on which a
recommendation depends.

Replacement text:

> Before consensus, a proposed research sub-goal or recommendation whose value depends on a
> claim that the result is open, novel, already proved, or already refuted MUST include a
> settledness check. Record the scoped proposition, check date, provenance under the verification
> rule, and one status: `KNOWN-PROVED`, `KNOWN-REFUTED`, `OPEN`, or `UNVERIFIED`.
> `RECALL`-only support remains `UNVERIFIED`. Such a proposition may enter `FINAL.md` only as a
> hypothesis with verification as a blocking first step; it MUST NOT be described as known-open
> or novel and MUST NOT justify the recommendation on that basis. Ordinary project-local
> milestones and acceptance criteria are exempt unless they depend on an external settledness or
> novelty claim.

This keeps the useful check without importing a literature-review gate into routine engineering.

### MAJOR-6 — Adopt amended

The current protocol already denies the facilitator unilateral adjudicative power: Phase 3,
Phase 7, and §5 require the track-applicable participant signoffs, and the ratified Phase 6 rule
allows disputes to close only by reviewer withdrawal, normal review consensus, or an operator
ruling. Requiring only one participant to “ratify” a facilitator ruling would actually be weaker.
There is still a disclosure and independent-review gap when one person holds several procedural
roles.

Replacement text:

> A facilitator has no dispute-adjudication authority beyond its position as a participant; any
> proposed resolution remains subject to the existing track-applicable consensus and signoff
> rules. On `standard` and `deliberation`, when the facilitator is also a participant and drafts
> `consensus.md` or `FINAL.md`, that artifact MUST include a `## Role concentration` section
> naming the roles held and listing each material facilitator-authored factual disposition,
> concession, or correction by claim identifier. At least one non-facilitator signoff MUST state
> that those listed items and their evidence were reviewed. If there are no such items, the
> section says so. This review creates no new veto, ownership transfer, or signoff weight.

I reject the uncheckable “procedural roles should be separable” wording and the preference for a
different drafter as protocol rules; disclosure plus an explicit independent review is the
enforceable safeguard.

### MAJOR-7 — Adopt amended

Unanimity is not evidence, but neither model relatedness nor independence of training data is
observable to a facilitator. The enforceable response is an adversarial artifact and comparison
of decision families, not a blanket assertion about model genealogy. This additional round is
proportionate only for non-mechanical judgment on `standard` and `deliberation`.

Replacement text:

> On `standard` and `deliberation`, when round 1 unanimously recommends the same decision family
> and the idea's output is primarily a judgment rather than a mechanically decidable artifact,
> consensus MUST NOT open until a later canonical round contains an
> `## Adversarial alternative` section. A participant MUST present the strongest materially
> different feasible alternative, its best supporting evidence, and an observation that would
> change the recommendation. If no credible alternative is found, record the search scope and
> why candidates failed.
>
> `consensus.md` `## Comparison & blind spots` MUST group nominally different proposals that
> share the same core premises or programme family and MUST state that agent count or unanimity
> is not independent evidence. This task does not change quorum, role weight, or artifact
> ownership.

The trigger and required evidence are visible on disk; “related models agreed” is not itself a
finding.

### MINOR-8 — Already covered by existing sections

No amendment is needed. Phase 1 states the behavioral rule; §11.A Phase 1 explicitly says
“Independence rule is a social one” and “There is no enforcement beyond agent discipline.”
§11.B “Independence in Round 1” repeats that the default is social and offers the sub-branch
protocol for stronger isolation. Requiring `parley-worktrees` globally would contradict the
existing optional, transport-aware strengthening and would add disproportionate ceremony.

### MINOR-9 — Adopt amended

§6.4 already says to copy a snippet when referring to material outside the deck, but it does not
require the kickoff author to disclose every source that materially shaped the brief. Conversely,
copying every source gathered while scoping is overbroad, may violate access or copyright
constraints, and rewards dumping irrelevant context.

Replacement text:

> Before invoking round 1, the kickoff author MUST include a `## Source packet` in `00-prompt.md`
> or a referenced sibling file. It MUST list every external source materially relied upon to
> frame the problem, constraints, alternatives, or factual premises. For each source, provide the
> relevant excerpt or a stable locator available to every participant, and identify whether the
> supplied text is verbatim or a summary. If material cannot be shared because of size, access,
> confidentiality, or rights constraints, disclose that asymmetry and do not present a
> source-dependent proposition as established evidence. Incidental sources that did not affect
> the brief need not be listed; existing no-secrets and access rules still apply.

This is a source-parity rule for material inputs, not a demand to copy an entire research trail.

## Concerns / open questions

### Installed-tool verification

I tested installed `parley 1.37.0` and `parley-deck-skill 2.1.0`. I did not run live
`parley roster init`, did not perform an actual project sync, and did not run the hosted-PONG
preflight because those actions would mutate the live deck or ping live agents. Mutating roster
tests used disposable decks under `/private/tmp` initialized with `git init` and `parley init`.

| Defect | Check run | Result and position |
|---|---|---|
| T1 | Fresh disposable deck: `parley init`, inspect §2, then `parley roster show` | **Confirmed.** `parley init` returned success but left both §2 tables empty; `roster show` failed with `could not read the §2 roster`. Endorse seeding §2 from the resolved central config, or fail-closing bootstrap until §2 is filled. |
| T2 | Live non-mutating `parley-deck-skill sync-project --project . --dry-run --json` | **Confirmed for field loss.** Input metadata reported `protocolRole: source`; the proposed replacement metadata omitted `protocolRole`. I did not execute the sync or the downstream unknown-role gate. Endorse preservation of known and unknown existing fields unless a migration explicitly owns them. |
| T3 | Copied the live protocol into a disposable initialized deck; ran `parley roster init -scope session -dry-run -json`, then the actual init there; ran `roster show` | **Not reproduced.** All four roster mappings were retained, including `kimi-1` with `AUTO=no`; no retired adapter appeared. Live and disposable `roster show` emitted no `roster init` hint. Do not endorse this defect or fix as a current fact without a reproducer/version history. |
| T4 | Live `parley preflight -no-ping -json` only | **Partially confirmed.** Output used adapter-family IDs (`codex`, `claude`, `agy`, `hermes`, `kimi`) rather than deck IDs and included non-rostered `agy`. I deliberately did not test hosted PONG, so `unavailable:no-pong` is unverified. Endorse roster-ID filtering; reserve judgment on no-PONG classification. |
| T5 | Live and disposable `parley roster show` | **Confirmed.** `claude-1` displayed `claude_opus-4.8-1m_max` while the resolved `MODEL` column was `claude-opus-5[1m]`. Endorse deriving the display name from the resolved model/effort rather than stale identity data. |
| T6 | Inspected the installed skill's Timeout Policy and timeout examples | **Partially confirmed.** The skill specifies 30-minute process defaults and says to poll long-running processes, but contains no concrete background-launch/session-polling pattern. I did not run a >10-minute foreground experiment, so the claimed host cap is unverified. Endorse documenting a harness-agnostic detached/session-polling pattern and distinguishing process timeout from one foreground tool-call limit; do not canonize “10 minutes” as universal without host evidence. |

These results are version-specific. In particular, T3 may describe an older build, but the brief
provides no versioned reproducer, so it cannot support a current defect finding by itself.

### Rule-scope questions

- The verification contract needs an explicit trigger: it governs claims the idea or a reviewer
  puts into a verification pass, not every descriptive sentence in every artifact. Otherwise
  claim-ID bookkeeping will overwhelm `fast` work.
- “Material” should mean that changing the claim's truth value could change a decision,
  acceptance criterion, severity, or close result. Any participant may challenge a drafter's
  classification in the next round; the facilitator should not decide it unilaterally.
- A locator proves that something was consulted, not that it was interpreted correctly. That is
  why CRITICAL-3 must reject automatic provenance ranking.
- The source-packet rule needs to preserve legal quotation limits and secret handling. A shared
  accessible locator is preferable to copying a large protected work.

## Risks

- A verdict ledger can become box-checking. Stable claim IDs should be required only when a
  verdict is actually assigned or a material conflict exists.
- `PRIMARY-SOURCE` can create false authority if the cited source does not entail the claim.
  Reviewable evidence and open challenge under the no-suppression rule remain necessary.
- `DISPUTED` must not become a way to close a decision that depends on the disputed fact; the
  dependency test in CRITICAL-3 is the actual gate.
- The MAJOR-6 role-concentration section can become ceremonial unless the non-facilitator names
  the listed claim IDs it reviewed.
- The MAJOR-7 adversarial round adds cost and may manufacture weak alternatives. Restricting it
  to unanimous, non-mechanical `standard`/`deliberation` judgments keeps that cost proportionate.
- Tool results can regress or be fixed independently of protocol text. CLI fixes should carry
  versioned tests; this protocol idea should not turn unverified anecdotes into permanent tooling
  requirements.
