---
agent: kimi-1
idea: roster-membership-overlay
phase: signoff
date: 2026-08-19
---

### kimi-1 — 🟡 accept with reservations

The substance — §1's fix shapes, the 3–2 split, §7's escalation — I accept. The reservations are
two attribution corrections (one inside the "unanimous" §1) and three (a)-side omissions, each
stated with its evidence below. None changes my position: **(a) — fix the gestures, keep the
authority model, overlay deferred behind the instrumented trigger.**

**1. §1 (the unanimous block)**

Confirmed, with one contested sentence. §1.1's defect table matches my reproductions. §1.2's fix
shapes I sign: D-A's gate keyed to the resolver's before/after member sets (the attributed
diagnosis is mine; locators re-verified this session — `membershipChange`,
`internal/app/roster_set.go:287-290`, keys on block existence in the file; the authority branch,
`internal/config/runtime.go:182-186`, keys on deck-block presence — PRIMARY, code read); D-B's
atomic triple of renderer, guard anchor and embedded default (`roster_render.go:73` compact
4-column vs. `drift_test.go:28` padded 3-column — PRIMARY, re-verified); D-C's outcome-honesty.
§1.3 and §1.5 I sign as written.

**Contested: §1.4's "@hermes-1, @kimi-1 and @zcode-1 each proposed it independently in round 1"
is false as to me** (SECONDARY — I re-read my round-01 file this session). My round-1 approach
proposes fixing F1/F2, the §2 generated-view marker, and the deferred trigger; migration appears
only as a caveat ("any migration plan must measure first"). I proposed the fleet migration in
**round 2** (fix-batch item 4), extending @hermes-1's round-1 proposal with the instrumented
trigger. The D-C note repeats the error and points at my round-1 item #3 — which is the marker
proposal, not a migration proposal. Correct text: "@hermes-1 and @zcode-1 proposed it in round 1;
@kimi-1 adopted and instrumented it in round 2." The correction is load-bearing, not cosmetic:
§1.4's "three independent round-1 proposals" supports the draft's independence framing, and the
true round-1 count is two.

One clarification I need on the concatenated record for §1.2 D-A: "silent materialisation must
not be the default" I sign — and my position additionally requires that an **explicit, previewed
materialize-all-then-apply path exists** (round-2: a bare refusal exports operators to the
gitignored layer, the least auditable file in the system). The draft's text does not foreclose
this; I record it so the final artifact cannot be read as foreclosing it.

**2. D-C**

**Reproduced and confirmed — PRIMARY, this session.** Isolated copy under `/tmp`
(`kimi-signoff-dc.*`, since removed), shared tree untouched, `git status --porcelain` empty after.
Five `[roster.*]` blocks each carrying only a redundant `adapter`: BEFORE, five rows STATUS `ok`
(deck-declared), `zcode-1` absent; `roster sync --yes` prints "5 redundant override(s) and 0
deliberate pin(s) removed; **the deck now inherits**"; AFTER, the same five rows STATUS `ok` —
still deck-declared, not `inherited-roster` — and all five `[roster.*]` headers survive. The claim
is false exactly as @claude-1 states. One sharpening the note did not state: the mis-statement is
in the **preview** too — `roster_sync.go:114` prints "removing these makes the deck inherit"
before any write, so the lie is told twice per run (PRIMARY, code read). The success line is
`roster_sync.go:169-170`; the mechanism — fields removed, headers left byte-identical — is
`removeRosterField`, `roster_sync.go:174-205`, so @codex-1's round-2 source inspection holds
(PRIMARY, independently re-read).

Does it change my sequencing? **One amendment, no reordering.** My round-2 batch was D-B → D-A →
transition verb + §2 disambiguation → migration. Amended: the D-C message fix joins the batch as a
peer of D-A (both small, neither writes §2, so neither waits on the D-B header decision), and the
census for decks previously "migrated" via `sync` (draft §2 item 2, currently unmeasured) becomes
a named **precondition** of the migration rather than a nicety — any deck an operator believes
was migrated and was not is a false premise in every subsequent fleet count. What D-C changes most
is the status of @zcode-1's `roster inherit` verb: until it exists, the only honest transition
paths are the manual two-file edit (the rule-2 trap that opened this idea) or nothing. The verb an
operator would reach for today claims the transition and does not perform it. For @hermes-1's and
@zcode-1's migration: it is no longer "executable with care" — it is **blocked on tooling**, and
the draft's §1.4 ("currently has no working instrument") states that correctly.

Does it change my (a)/(c)? **No.** D-C is a third tooling-layer honesty defect. Per @claude-1's
own round-2 framing, citing defects against the overlay is a sequencing argument, not a design one
— my (a) is a sequencing argument and owns that. All three verbs §2 documents as the way to change
a roster now misreport their own effects; layering a second membership semantics on this base is
worse sequencing today than it was yesterday. The design question is untouched.

**3. §2 — is my side stated at full strength?**

Every sentence the draft attributes to me is accurate (SECONDARY — I re-read all ten round
artifacts and the D-C note this session): §1.1's finder/reproducer credits; §1.2 D-A's diagnosis
and fix shape; §2's (a)-side seams argument (near-verbatim from my round-2, PRIMARY locators
re-verified above); §4's "@kimi-1 via the D-A/D-B reproductions" inside @hermes-1's reading. The
split table signs me under (a) — correct.

My side is nonetheless stated below full strength. Three omissions:

(i) **The draft never carries my round-2 PRIMARY narrowing of (c)'s central claim.** claude-1's
"no way to change one local setting" is false as stated: the gitignored `agents.local.toml` layer
carries value-only overrides with membership untouched (I ran it in round 2;
`runtime.go:150-164` bars non-machine layers from membership — PRIMARY, re-read this session). The
live gap is the narrower *"no committed, tool-supported way."* That narrowing decides what §2.1 is
even testing: a values-only-no-membership mechanism **exists today at the gitignored layer**; the
open question is the committed layer. The draft quotes @zcode-1's narrowed sentence but not the
evidence that forced the narrowing, so a reader cannot see that (c)'s claim was already corrected
once by measurement.

(ii) **The draft quotes @claude-1's "a trigger that cannot fire" at full strength and omits the
standing (a)-side answer.** My round-2 instrumented trigger: the attended migration already asks
each deck "deliberate override or stale copy?", so recording the answers converts the demand
question from an argument into a number ("migration records ≥2 deliberate ±1 deviations"). A
trigger bound to a questionnaire the migration already runs *can* fire; claude-1's objection
targets the passive version (@hermes-1's round-2 trigger is likewise passive; mine is the
amendment that makes it active). Stating the objection without its answer weakens (a).

(iii) **"What would settle it" #3 misattributes the value-case trigger amendment to @hermes-1.**
It is @zcode-1's — round-2 position-changes #2: "The trigger must also cover the value case: a
deck that must track machine membership live while overriding at least one value" (SECONDARY —
read this session). @hermes-1's round-2 trigger went *stricter* (three limbs), not value-phrased.

And one counter-witness the draft's (c) form does not carry: its "that demand is not zero" rests
on the owner's *originating* sentence. The owner's **dated instruction of record** in this deck's
`agents.toml` (2026-08-19) says nothing local at all (PRIMARY — I read the file in round 1;
@zcode-1's round-01 #6 independently). The instruction post-dates the originating sentence and
points the opposite way. §2 need not argue my side — but the record should show the witness
exists.

**4. §2.1 — the decisive unrun experiment**

**(i) Technically possible?** Yes, in two shapes, and the choice between them is the whole
question:

- **Content-keying** — a block constitutes membership only if it carries some marker field (the
  natural candidate is `adapter`, which D-A's written block lacked): makes a values-only block not
  collapse the roster. What breaks: authority rule 1 stops being "block presence"
  (`runtime.go:182-186`); the gate's keying (`roster_set.go:287-290`); and — decisively — every
  one of the 37 full-declared decks whose blocks all carry values becomes a semantic coin-flip
  unless the chosen key provably never reclassifies an existing block. A content rule that flips
  any unmodified file is the §1.3 sin this idea has already sworn off. Proving it doesn't requires
  a fleet census of block **contents** — every census to date, mine included, counted blocks, not
  their fields. That census is part of the experiment. (If the census finds even one existing deck
  whose membership depends on adapter-less blocks, the cheap shape is dead on arrival.)
- **Explicit marking** — an opt-in stanza or per-block key separating value from membership:
  back-compatible by construction, and it **is** @codex-1's recorded `[membership]` schema minus
  add/remove — the overlay's core separation mechanism wearing a smaller hat. "Separation without
  the overlay" converges on the overlay's mechanism; what remains in dispute is only whether
  deltas ship with it.

Also breaking in every shape: rule-2 fall-through when a deck's blocks stop constituting
membership (does it drop to legacy §2 or machine?); `roster init`/`sync`/`render` assumptions that
blocks are membership; the skill's documented semantics. So: possible, not free, and the cheapest
shape is the one §1.3 forbids until measured safe.

**(ii) Would a positive result change my (a)/(c) answer? No.** A positive result *strengthens*
(a): the value-override gap gets serviced inside the current authority model, (c)'s live case is
absorbed into the D-A fix, and the overlay's residue shrinks to membership ±1 deltas — whose
measured demand is zero across five censuses. A negative result does not flip me either: it
confirms @codex-1's schema as the only back-compatible separation shape and raises the value of
keeping it on file; the demand gate (my round-2 falsifiers: one named deck with a live need,
evidence that exclusions are deliberate, an explicit owner instruction) is untouched either way.
**I therefore do not declare the split premature on my own account.** But both (c) signers wrote
falsifiers a positive result engages directly — @claude-1 ("what would have made me stay at NO
CHANGE: if a `[roster.*]` block that sets only values did not constitute membership") and
@codex-1 ("I would answer (a) instead if someone demonstrated a committed, value-only override
mechanism…", though conjunctive with post-migration demand) (SECONDARY — both read this session).
The experiment is cheap, is decisive for the two signers whose votes compose the split, and
should run before FINAL — as an isolated-copy evidence task including the block-contents census —
not because my answer depends on it, but because theirs do.

**5. What I do not sign**

- §1.4's sentence as written (the round-1 migration attribution to me) — correction in §1 above.
- "What would settle it" #3's attribution of the value-case trigger to @hermes-1 — it is
  @zcode-1's; correction in §3(iii).
- Any reading of §1.2's D-A shape that forecloses the explicit, previewed materialize path —
  clarification in §1.
- "That demand is not zero" as uncontested record, while the owner instruction of record says
  nothing-local — §3's counter-witness stands.
- §5's treatment of @opencode-1 (incomplete participation, no inference either way) is correct
  and I sign it as written.

**Verification basis.** D-C reproduction: PRIMARY (ran this session in a `/tmp` copy; shared tree
never written, `git status --porcelain` empty after). Code locators — `roster_sync.go:114`,
`:169-170`, `:174-205`; `roster_set.go:287-290`; `runtime.go:150-164`, `:182-186`;
`roster_render.go:73`; `drift_test.go:28`: PRIMARY (read this session, quoted). Attribution
audits: SECONDARY (all ten round artifacts, the prompt, and the D-C note re-read this session).
No RECALL claims relied on.
