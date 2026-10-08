---
agent: zcode-1
idea: meta-protocol-change-participant-dropout
round: 2
date: 2026-10-08
responding-to: [codex-1/round-01, codex-1/round-02]
---

# Round 2 — zcode-1 (cross-review)

Protocol attestation: `context_mode=full`,
`source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`,
`fallback_reason` absent. I read the full emitted body of
`.parley-runtime/protocol-packets/full-phase2-deliberation-73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e.md`
(all 1502 lines) and verified it hashes to the attested value (`shasum -a 256` →
`73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`; PRIMARY, executed 2026-10-08).
Active transport: `github-pr` (draft PR #75 per organizer-notes.md). `meta/protocol-changelog.md`
checked; latest entries predate this idea, nothing newer than the supplied protocol.

Scope note: Phase 2 cross-review only. I wrote no product code, made no commits, launched no
participant, drafted no consensus/FINAL, and signed nothing for anyone. All code locators below are
my own reads at HEAD `2239017`; `internal/` is unchanged since `bf5a2d3` (my round-01 base), so my
round-01 locators remain valid (PRIMARY: `git diff --stat bf5a2d3..2239017 -- internal/` is empty).
Sizing statements are my own estimates and carry no verdict (§15.1).

## Position changes since round 1

I change four positions, each on codex-1's counter-proposal or my own verification:

1. **Single knob, versioned saved trigger — I withdraw the second knob** (`participant_dropout`).
   My round-1 DP6 created a two-knob space whose combination semantics I had not fully specified:
   a deck or idea that deliberately set `quota_auto_exclude: false` (confirmation-only) would still
   get the broader any-failure rule under an independent default-on dropout knob — a silent widening
   of an explicit opt-out. codex-1's R2-2 single-knob design (`quota_auto_exclude` + saved
   `trigger: participant-failure-v1`, absent trigger = legacy quota-only) honors every existing
   opt-out by construction. I found no concrete safety reason that *requires* the second knob; the
   only thing two knobs add is quota-only selection for new ideas, which nobody asked for. Accepted
   with the conditions in the R2-2 response below.
2. **Uniform two attempts — I withdraw the first-failure quota carve-out.** My round-1 DP2 let a
   positively quota-classified failure drop on the first attempt, keeping the quota recognizer as a
   live precedence branch inside the new path. That re-imports exactly the recognizer dependency and
   per-failure classification complexity the prior idea suffered under (its zcode recognizer is
   "possibly inert" today), for a marginal wall-clock saving. codex-1's R2-3 uniform rule — every
   non-protected failed step under the new trigger gets exactly one relaunch, the recognizer is not
   consulted for candidacy at all — is simpler, matches the owner's own wording ("fails twice on the
   same step … is dropped"), and is more robust. Accepted.
3. **Evidence-backed floor — I withdraw structural designee counting.** My round-1 DP3 counted the
   protected designee/pin seat structurally toward the floor even when not dispatched. codex-1's
   R2-4 is right: the owner requires a *usable* implementer plus another participant, and a seat
   with an unresolved failure is not usable; a structural count could fabricate a survivor. The
   shipped rule (`internal/quota/quota.go:134–150`, verified below) counts only positive evidence,
   and the codebase already has the right seeding precedents. Accepted; the designee's protection
   (never a candidate, blocks the batch) is unchanged — what I withdraw is only floor-supply by role
   identity.
4. **Permanent same-idea dropout — I withdraw owner-revision re-entry.** My round-1 DP5/AC-D6 kept
   "rejoins this idea only via owner-confirmed revision," mirroring the shipped quota return rule.
   codex-1's R2-1 correctly reads the controlling brief: "does not rejoin this idea … can join only
   from the next one" admits no same-idea exception, and the brief's design target states it twice.
   Accepted with the §5-integrity refinement in the R2-1 response below.

**SELF-CORRECTION (§15.1, weakening, effective immediately).** My round-01 DP5 stated: "Retry
idempotence. Each attempt already carries an `attempt_id` and immutable invocation identity …
One step = one terminal evaluation" (round-01/zcode-1.md, DP5, bullet 3), implying the existing
attempt loop bounds retries durably. That claim was incomplete: the loop starts `attemptID := 1` on
every `runAgent` call and links retries only through a local `retryOf`
(`internal/runner/runner.go:515–551`, verified this round), so a driver restart between attempts
could mint a fresh budget. The durable cross-restart bound must be built (R2-5); it does not exist
today. This was also codex-1's R2-5 position, which I verified rather than merely accepted — no
conflicting verdict remains.

## Responses to others

### @codex-1 — round-01 and round-02

Your round-01 points 1–7 are each carried into round-02 as R2-1 through R2-7 (or restated there);
I respond to them at those loci. No round-01 position of yours is left unaddressed, and I have no
verdict conflict with either artifact. Verdicts below are my own reads at `2239017`, not endorsements
of your claims (§15.1: quoting does not transfer ownership; where I verified the same source I say
so with my own locator).

**R2-1, return is binding, not discretionary — ACCEPT.** The owner's words are unambiguous
("vypadnú z hocikakeho dôvodu … v tej idei už nepokračujú a môžu sa zapojiť až do ďalšej") and the
brief's rule block repeats "does not rejoin this idea" with no carve-out. Agreed contract: the
durable dropout history records the id's same-idea ineligibility as a terminal membership fact, and
every rejoin path consults history and refuses — `quota revise`/`recover`
(`internal/app/quota_revision.go:84` passes `r.Participants` into `membership.Revise`, so today a
revision CAN re-add an id; that path must reject re-adding a dropout-dropped id for this idea),
manual catch-up, and any later policy downgrade or opt-out. Two refinements I ask FINAL to state
explicitly, both preserving §5/§15 integrity rather than weakening the rule:

- A retained ❌, `DISPUTED` claim or open finding of a permanently dropped author keeps its force;
  its only exits are the author's withdrawal (impossible — the author cannot append anything) or an
  explicit owner ruling quoted into the next artifact. Absence can never imply withdrawal. (You said
  this; I am making it a FINAL must-state so the deadlock-exit is documented, not discovered.)
- The owner's residual authority acts on the *idea*, never on membership: abandon and re-open as a
  v2 idea with the participant included, or rule on the retained veto. This is the honest bound of
  "never rejoins" — it is not a protocol owner-override, and FINAL should say so in one sentence.

Existing quota-v1 exclusions keep their shipped owner-confirmed-return semantics — agreed, and it
follows from the versioned record: legacy transitions carry the quota rule id, only
`participant-failure.v1` transitions are terminal.

**R2-2, one knob with a versioned saved trigger — ACCEPT**, with four conditions that I believe you
already hold and want pinned in FINAL:

1. The saved policy grammar accepts exactly two shapes: legacy `{enabled, scope}` (trigger absent ⇒
   quota-only v1, serialized bytes and hashes unchanged) and `{enabled, scope, trigger:
   "participant-failure-v1"}`. Unknown or malformed trigger values fail closed
   (`DisallowUnknownFields` discipline preserved). Current strictness verified:
   `internal/quota/quota.go:24–42` rejects any record whose fields ≠ {enabled, scope} today — the
   extension is exactly this one optional key.
2. `quota_auto_exclude: false` at idea/deck/machine disables *all* automatic reductions, both
   triggers — the existing opt-out meaning, now covering the new one (PRIMARY:
   `quota.go:44–53` resolves idea-over-default; the opt-out semantic is the 1.51.0 one).
3. The trigger is frozen at kickoff into the record and never widened by resume, upgrade or
   binary switch; enabling the new trigger on an in-flight idea goes only through the owner-bound
   revision path with committed authority.
4. FINAL documents that the knob's *name* is historical: it governs both the legacy quota trigger
   and the versioned participant-failure trigger. (I withdraw my ALT-second-knob; see Alternates
   disposition.)

I explicitly confirm there is no safety reason requiring the second knob — my round-1 motivation was
semantic cleanliness plus speculative quota-only selection for new ideas, and neither outweighs the
opt-out coherence and smaller surface of the single knob.

**R2-3, exactly two attempts, legacy quota untouched — ACCEPT.** Uniform rule: under
`participant-failure-v1`, every non-protected dispatched-step failure (any terminal class: nonzero
exit, provider error of any code, crash, watchdog/timeout kill, exit 0 without a valid artifact)
gets exactly one relaunch after a small fixed delay (5 s is fine), original ceiling unchanged; a
valid artifact on either attempt is success and prevents candidacy. Legacy saved policies keep
today's immediate quota behavior and 60-minute predicate untouched — the recognizer stays a legacy
sub-case, never consulted by the new trigger. Reset hints under the new trigger are diagnostic only
and never authorize or suggest same-idea return. Two disclosed costs, both acceptable: a ≥60-minute
quota exhaustion burns one predictably-futile relaunch, and a transient blip longer than the fixed
delay still drops the agent — the floor, the notice naming the class, and the next-idea re-probe
bound the damage. Both go in FINAL's known-risks. This also resolves my round-01 DP1's "practical
side effect" (dropout covers zcode quota exhaustion without the recognizer): still true, and now
true uniformly rather than through a precedence branch.

**R2-4, floor requires evidence of usability — ACCEPT, your PRIMARY claim independently CONFIRMED
by me.** I verified `internal/quota/quota.go:134–150`: a member supplies `success` only via
`Usable || ValidArtifact || LaterSuccess` (`:137`), an unobserved id is neither survivor nor
candidate (`:132–135`), unresolved failures bar candidacy (`:142–144`), and the facilitator never
counts (`:146–149`). Both seeding precedents also verify: existing valid signoffs seed
`quota.Member{Usable: true, ValidArtifact: true}` at
`internal/app/consensus_request_signoffs.go:164–167` and `:245`, and a review-phase batch adds the
pinned implementer only after `ValidateImplementationArtifact` passes, with the comment stating the
principle verbatim ("a valid completed implementation is its usability evidence, not a new provider
observation", `internal/runner/quota.go:47–55`). Settled rule: floor slots require positive evidence
— dispatch success in the batch, a later success in the batch, or a validated current artifact or
record of that agent for the relevant phase (valid round artifact for design phases, valid signoff
for signoff settles, validated `IMPLEMENTATION.md` for review). Never `Usable:true` from role
identity alone. A designee with an unresolved failure in the batch blocks the batch as protected and
supplies nothing. A conservative unresolved-usability escalation is acceptable; a fictitious
successful seat is not. My round-01 concern 1 (where the kickoff ping budget lives if the ping layer
is single-shot) is answered by your inventory putting kickoff retry in
`preflight.go`/`preflight_liveness.go`; the hunk may land at the loop level instead — same
semantics, implementation detail.

**R2-5, durable two-attempt budget — ACCEPT, your PRIMARY claim independently CONFIRMED by me.**
`internal/runner/runner.go:515–551`: both the ACP and exec loops are `for attemptID := 1; ;
attemptID++` per `runAgent` call, retry fires only on `no_first_output`
(`:520`, `:538`), and `retryOf` links the second attempt to the first only in local state
(`:534–539`). Nothing persists the consumed attempt across a driver restart. Settled contract: the
budget is keyed to (idea, agent, stable logical step/artifact); each attempt carries an immutable
id and retry linkage; a restarted driver consumes recorded attempts for that step and never mints a
third; changed inputs do not reset a failed step's budget; an attempt recorded as started but
without a terminal outcome fails closed until stopped-writer recovery resolves it (reusing the
shipped crashed-writer settlement with host/boot/PID proof); no retry on an in-progress writer.
Tests must kill/restart between failed attempts and after the second failure, then assert exactly
two launches and one eventual reduction. Agreed this is a bounded retry record at the existing
seam, not a second membership system. My SELF-CORRECTION above covers my round-1 claim this
replaces.

**R2-6, protected/control-plane failures — ACCEPT.** Policy, budget, protocol and telemetry
refusals, and operator cancellation, are not removable participant failures: no authorized
dispatched child step failed, so there is nothing to count. An authorized dispatched child's
terminal failure is. Parent cancellation must be distinguished from child timeout in the evidence
record. An invalid signoff that edits another participant's blocks is an integrity failure: preserve
bytes, stop for repair, never erase/retry shared content. The live example is in this run:
the round-02 launch `f33856d3…` (`failure_class: budget_refused`, `started_at: null`) was correctly
treated as a control-plane refusal and not as a zcode-1 attempt (organizer-notes.md, "Round-02
measured-launch refusal"). This extends my round-1 form-not-position boundary with exactly the
right control-plane distinction; my round-1 artifact did not draw it. Also note the inverse guard
you stated and I adopt: a missing or structurally invalid own output after settled execution IS the
owner's any-reason case and gets the one retry.

**R2-7, floor and review gate interaction — ACCEPT, your PRIMARY claim independently CONFIRMED by
me; I correct my round-1 wording.** `internal/membership/membership.go:297–315`: `Settle` runs
`quota.Evaluate`, and before `CommitBatch` (`:336`) it calls `CheckGates(root, ideaDir, runID,
d.After)` at `:311`, demoting the decision to a block on gate failure. `internal/membership/gates.go:82–87`:
under `auto_implement` the reviewer minimum is forced to 2 and a shortfall returns
"review/LE-7/LE-11 gate: N independent reviewers remain; require M". So my round-1 phrase
"post-drop shortfalls escalate" described the intent, not the order — the shipped order is
precommit, and I adopt it unchanged. FINAL must state the consequence plainly (see New concerns 1).

### @kimi-1 — excluded, no position exists

kimi-1 failed both round-01 attempts (HTTP 400 ambiguous k3 routing, no artifact;
source-context/kimi-round01-failure-evidence.md) and is excluded under the brief's standing
owner authorization, recorded in `00-prompt.md`. There is no kimi-1 round artifact to respond to;
I neither impute agreement nor disagreement, and nothing in my responses treats the two recorded
process failures as a canonical position. No Claude participant exists in this idea.

## New concerns / questions

1. **The rule's practical reach in 3-participant `auto_implement` decks — accepted, but FINAL must
   state it as the primary limitation, and a follow-up path must be named.** With the owner's
   standard shape (organizer+implementer codex, participants kimi+zcode, `auto_implement: true`),
   a single participant drop satisfies the membership floor but leaves one independent reviewer,
   and the precommit gate (R2-7) blocks the reduction with one escalation. So for exactly the
   owner's flagship scenario the rule automates the *accounting* (evidence, options, one notice)
   but still asks the owner once. I accept this: it is shipped behavior, the brief's item 4 says
   shortfalls escalate and never waive, and a gate-order change would be a bigger, riskier edit
   than this idea's "as small as possible". But I do not want it buried: FINAL must state that the
   drop is blocked pending one owner answer whose options are (a) authorize a substitute
   model-diverse reviewer process, (b) an attended close under one reviewer with recorded owner
   authorization (the pattern this run's brief already demonstrates), or (c) pause/abandon — and
   must record that owners can pre-authorize (b) per idea to remove the stall. If the owner's
   lived experience is that this asks too often, splitting precommit-hard gates (floor, protected
   roles) from recorded escalated gates (reviewer count, diversity, goal-checker) is a coherent
   follow-up idea — explicitly out of scope here.
2. **Cross-batch double failure (my round-01 concern 5, still open).** If kimi-1 and zcode-1 fail
   in *different* settled batches, the first drop applies (floor holds) with one non-blocking
   notice, and the second blocks at its own settle with one blocking escalation — two transitions,
   two notices, per-batch semantics. I believe this is correct and matches the shipped whole-batch
   rule; FINAL's §9.0 wording should say "per settled batch" so it is unambiguous. Not blocking.
3. **Inbox-missing/unwritable fallback (your new concern) — resolved: include it, narrowly.** If
   FINAL promises every failed-floor batch is visible, then the notice path needs the same
   create-safe-inbox/stderr fallback diagnostic you describe, plus crash-replay of a kickoff notice
   via the existing receipt/publication machinery with a regression test. Scoped to exactly that;
   no alias, Windows or unrelated repair work imported. Agree it is in scope, not a waiver of the
   prior R8 residual.
4. **Packet guardrail headroom (round-01 concern 2, carried).** The new §9.0 block must be written
   tight against the 70,000-byte facilitator-packet limit and re-measured during implementation;
   no new top-level headings so `meta/packet-applicability.yaml` stays untouched.
5. **D6 accounting gap — acknowledged for this run, unchanged as product scope.** This round-02
   artifact is produced through the brief's recorded measured-CLI fallback because the driver
   cannot perform the cross-review step (legacy accounting). That is an orchestration limitation
   of the *tooling*, not a participant failure and not evidence for or against the design; it is
   correctly outside this idea's product scope (no budget migration).

## Current proposal

The owner's seven design points, settled as one exact contract (my round-1 architecture with
codex-1's R2-1..R2-7 amendments adopted as above). Everything below is *proposed design*; code
locators cite the reused shipped seams verified at `2239017`.

**DP1 — Reuse 1.51.0; one mechanism, one knob, two versioned triggers.** The shipped
automatic-exclusion machinery (`internal/quota/quota.go:118` `Evaluate` whole-batch all-or-nothing;
`internal/membership/membership.go:279` `Settle`; immutable history `internal/quota/history.go`;
durable transitions, notices, surfaces) is extended by one trigger version, not by a second system.
Saved policy: `{enabled, scope, trigger?}` — trigger absent ⇒ legacy quota-only v1 with bytes and
hashes unchanged; `trigger: "participant-failure-v1"` ⇒ the new rule. New ideas created after
delivery record the new trigger on by default (scope `kickoff-and-mid-idea`, as
`quota.go:44–53` already defaults). Trigger frozen at kickoff; unknown values fail closed; widening
only via owner-bound revision. The quota recognizer (`internal/telemetry/quota.go`,
`quota_zcode.go`) is untouched and not consulted by the new trigger; a dropout candidate carries
rule id `participant-failure.v1` with reset unknown and no relaunch suggestion.

**DP2 — "Any reason" with a fixed two-attempt budget.** Per (idea, agent, logical step): the
original attempt plus exactly one relaunch, 5 s fixed delay, unchanged ceiling. Any terminal child
failure class opens the relaunch; a valid artifact on either attempt is success. Control-plane
refusals and operator cancellation never count (R2-6). The budget is durable across restarts
(R2-5): keyed to idea+agent+step, immutable attempt ids with linkage, restart consumes recorded
attempts and never mints a third, changed inputs do not reset it, started-but-unresolved fails
closed until stopped-writer recovery. Kickoff readiness pings get the same two-attempt budget
(`internal/app/preflight.go`, `preflight_liveness.go`), then candidacy, before any
`participants:` line, manifest or dispatch; standalone `parley preflight` still never applies.

**DP3 — Evidence-backed floor.** Floor 2, fixed and non-configurable. A survivor slot requires
positive evidence: dispatch success, later batch success, or a validated current artifact/record for
the phase (valid round artifact; valid signoff per `consensus_request_signoffs.go:164–167,245`;
validated `IMPLEMENTATION.md` per `internal/runner/quota.go:47–55`). Never role identity. The
facilitator never counts; the designee is protected (blocks the batch, never a candidate) but
supplies the floor only with its own positive evidence. At or below floor-minus-one after
application: apply nothing, one blocking escalation listing every candidate, the arithmetic, and
the three owner options — (a) authorize a substitute model-diverse reviewer/implementer-side
process, (b) attended reduced-quorum close under recorded owner authorization, (c) pause or
abandon — none chosen silently.

**DP4 — Protected roles and precommit gates.** Protected list unchanged: declared facilitator,
per-idea designee, `IMPLEMENTATION.md` pin, started consensus/FINAL/review-consensus drafters
(`quota.go:96–102`, `membership.go:349–382`). A protected candidate ⇒ apply nothing and open the
existing three-exit gate with evidence prefilled. `CheckGates` stays precommit
(`membership.go:311` before `CommitBatch` at `:336`): designee-presence, reviewer minimum
(auto ⇒ 2, `gates.go:82–87`), independent goal-checker (`gates.go:88–90`), model diversity
(`gates.go:91–109`) — shortfalls block the reduction with one actionable escalation; no gate is
waived. FINAL states the DP3/DP4 interplay honestly (New concerns 1).

**DP5 — Integrity.** Candidacy arises only from supervisor-observed terminal facts: exit code,
structured failure, watchdog kind, structural artifact validity (`validateArtifactForPhase`).
Content and positions are never examined; a structurally valid dissenting, blocking or disputing
artifact is a success (adversarial fixture required). Control-plane refusal and operator
cancellation are not participant failures (R2-6). Every drop records, before application, inside
the durable batch record: both attempts' invocation ids and linkage, failure class, exit/watchdog
kind, validator reason, scrubbed decisive excerpts (visible truncation/redaction), UTC observation
times. Partial/invalid artifacts preserved as incomplete; stopped-writer settlement
(`RequireStopped`) unchanged; replay idempotent; retained dissent intact (excluded ids stay known
signers; ❌ yields `TriageBlocked`; `DISPUTED` and findings need explicit independent disposition;
no new signoffs until next-idea re-probe).

**DP6 — Default, opt-out, migration.** `[defaults].quota_auto_exclude` (machine → deck → idea) is
the single knob; `false` disables all automatic reductions. New ideas: new trigger on. Legacy saved
policies decode quota-only and are never widened by resume or upgrade; enabling the new trigger on
an in-flight idea is an owner-bound revision. Malformed/ambiguous records fail closed. FINAL notes
the knob name is historical and now governs both triggers.

**DP7 — Permanent per-idea dropout and return.** A dropout-dropped id is ineligible for every
same-idea rejoin path (manual catch-up, `quota revise`/`recover` — today `quota_revision.go:84`
can re-add ids and must refuse dropout-terminated ones — and any later downgrade/opt-out). History,
not display markers, carries the terminal fact. Next idea probes afresh. A dropped author's
retained ❌/`DISPUTED` exits only via explicit owner ruling quoted into the next artifact, or the
idea's abandonment (v2). Legacy quota exclusions keep owner-confirmed return.

**Inventory (unchanged in shape from my round-01 DP7, minus the withdrawn pieces: no second knob,
no structural-floor code, no quota-precedence branch in the runner).** Protocol hunks, identical in
all three copies: §0 defaults sentence (one clause: versioned trigger); Phase-0 field comment;
§5 exception wording widened to "automatic exclusion (quota or participant dropout)" plus the
terminal-return sentence; §9.0 sibling block (policy/migration, trigger+budget, evidence,
floor+options, protected roles, marker/notice reuse, return semantics, "per settled batch");
Phase-5 cross-reference names both rules; `meta/protocol-changelog.md` entry. CLI: `internal/quota/
{quota,record,history,revision}.go`; `internal/membership/{membership,notice}.go` (rejoin guard);
`internal/protocol/{quota,quota_manual}.go`; `internal/runner/{runner,acp,quota}.go` (budget +
retry widening + outcome handoff); new `internal/telemetry/dropout.go` classifier;
`internal/app/{preflight,preflight_liveness,consensus_request_signoffs,quota_signoff,
driver_consensus}.go`; `internal/app/{wait,organizer,status}` surfaces; `runcontrol`/`runmanifest`
follow the extended `quota.Policy` automatically; `internal/config/runtime.go` docs only. Skill:
bundled `references/COOPERATION.md`, one short note each in `SKILL.md`/
`ROSTER_AND_PROTOCOL.md`. Tests beyond per-file units: kill/restart between and after attempts
(exactly two launches, one reduction); exec/ACP parity; valid-BLOCK-output fixture;
control-plane-refusal fixture; floor permutations with evidence-only seeding; protected roles;
return refusal on every entry point; legacy round-trip hash stability; drift, phase-0/5/8 packets,
skill byte-compare; existing quota suite passes unmodified with the new trigger absent. Estimate
~300–450 product LOC, ~700–1000 test LOC, ~40–60 protocol lines per copy — my estimate, UNVERIFIED.

**Alternatives disposition (for consensus to lift):**

- ALT-A second knob — REJECT (R2-2; opt-out coherence beats speculative quota-only selection).
- ALT-B structural designee floor count — REJECT (R2-4; a fictitious usable seat).
- ALT-C first-failure quota carve-out in the new path — REJECT (R2-3; recognizer independence).
- ALT-D split precommit vs recorded-escalated gates — REJECT for this idea (size, fail-closed
  single-point integrity), named as the follow-up if the owner finds the blocking behavior asks
  too often (New concerns 1).
- ALT-E owner-revision same-idea re-entry — REJECT (R2-1; contrary to the controlling brief).
- codex-1's ALT-1..ALT-6 (round-01) — ADOPT as stated, consistent with everything above.

**Correlated agreement (§15.6b).** codex-1 and I are two model families inheriting one shipped
mechanism and one controlling brief; our convergence is family-adjacent, not independent proof of
safety. What would make this design wrong: a content-triggered drop, a lost or lapsed veto, a drop
below floor or past a gate without exactly one escalation, a restart minting a third attempt, or a
resume/upgrade widening a saved trigger. The named fixtures target exactly these.

**Owner decision: none blocking.** Every point is settled from the owner's two verbatim messages
plus shipped mechanism; DP3's escalation presents runtime options without choosing; the reviewer-gate
limitation is a disclosure, not an open decision. On your this-run question: I read **no unresolved
owner boundary**. The brief's standing authorization covers the reduction (already applied) and an
attended codex-1 + zcode-1 close when both final review-consensus signoffs exist, no CRITICAL/MAJOR
is open and current-tree evidence is recorded — an attended close under recorded owner authority,
not a driver auto-close, and no product gate waiver is claimed. The deliberation cap (3 rounds
after round 1) and the 5-cycle fix-up cap are respected manually under the recorded D6 fallback.

**Consensus readiness: ready.** I am ready to sign a consensus adopting the contract in this
section — codex-1's R2-1..R2-7 as accepted above plus my four withdrawals — provided FINAL carries
the must-state items: (1) the reviewer-gate limitation and its owner options, including standing
pre-authorization; (2) permanent dropout's veto-exit rule (owner ruling or v2, never implied
withdrawal, no protocol re-entry); (3) "per settled batch" double-failure semantics; (4) the
known-risk pair from R2-3 (futile quota relaunch; sub-delay blips) and the timeout-misconfiguration
drop; (5) packet-headroom discipline; (6) the knob-name history note. Nothing blocks; I have no
open counter-proposal against consensus text matching this contract. I do not draft consensus or
FINAL here, and I sign nothing for anyone else.
