---
agent: zcode-1
idea: meta-protocol-change-participant-dropout
round: 1
date: 2026-10-08
---

# Round 1 — zcode-1

Protocol context attestation: `context_mode=full`,
`source_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`,
`packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`,
`fallback_reason` absent. I verified the packet file
`.parley-runtime/protocol-packets/full-phase1-deliberation-73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e.md`
hashes to the attested value and read the full body (identical to the live protocol supplied at
launch). Active transport: `github-pr`. `meta/protocol-changelog.md` checked for updates newer than
the supplied protocol: none found (PRIMARY, read 2026-10-08).

Scope note: Phase 1 design only. I wrote no product code, made no commits, launched nobody, and read
no other participant's round-01 file. All code locators below are my own executed reads at HEAD
`bf5a2d3` (worktree `participant-dropout`, base `d16ee9c`, CLI 1.51.0) unless provenance says
otherwise. Claims I originate (estimates, predictions) are marked UNVERIFIED and carry no verdict.

## Summary

The shipped 1.51.0 quota auto-exclusion machinery — batch evaluation, immutable history, durable
transition records, idea-lifetime lease, protected roles, notices, retained dissent — is
rule-agnostic in everything except its trigger predicate, its floor arithmetic, and two
quota-specific validation constants. I propose extending exactly that machinery with a second,
broader trigger ("participant dropout"): a non-protected participant whose dispatched step fails for
ANY terminal reason — after a fixed budget of the original attempt plus exactly one relaunch — is
dropped from this idea only, without asking the owner, provided the floor holds. The quota predicate
stays untouched as the stricter sub-case under its own knob (a positively recognized quota failure
is dropout-eligible immediately, without burning the retry, and keeps its reset hint). A new
presence-aware knob `[defaults].participant_dropout` (deck override, per-idea opt-out, on by default
for ideas created after delivery) governs it; saved quota-v1 policies decode with dropout absent and
are never widened on resume. I find no open owner decision that must block FINAL, and recommend no
proposal stop.

## Existing alternatives

What the toolchain already ships for each mechanism this proposal builds on (all PRIMARY, read at
`bf5a2d3`):

- **Automatic mid-idea membership reduction:** `internal/quota/quota.go:118` `Evaluate` (whole-batch,
  all-or-nothing, duplicate-safe), `internal/membership/membership.go:279` `Settle` (called only
  after all supervised writers stopped), `.go:430` `RequireStopped`, immutable history at
  `internal/quota/history.go:112` `ReadHistory`, durable write `internal/quota/history.go:188`.
- **Kickoff-side exclusion:** `internal/app/preflight.go:412–436` evaluates the batch over readiness
  probe results and writes the decision into idea creation (`internal/app/app.go:1946–1971`,
  `runcontrol.Create` with `QuotaPolicy`/`QuotaDecision`); standalone preflight never applies.
- **In-runner retry:** `internal/runner/runner.go:516–548` — an existing retry-ONCE attempt loop,
  currently gated on `no_first_output` watchdog kills only (both exec and ACP modes).
- **Terminal failure classification:** `internal/runner/failclass.go:30–56` (bounded regex table:
  rate-limit, auth, billing, overloaded, …) plus supervisor-assigned watchdog kinds
  (`no_first_output`, `stalled`, `timeout`); `internal/runner/runner.go:762–766` already treats
  "exit 0 but no valid artifact" as a failure (`"artifact missing or invalid"`).
- **Floor:** `internal/quota/quota.go:13` `Floor = 2` with usable-survivor counting that already
  excludes the declared facilitator (`.go:146–149`).
- **Protected roles:** `internal/quota/quota.go:96–102` (`Facilitator`, `Designee`,
  `PinnedImplementer`, `StartedDrafters` — global defaults deliberately unprotected) +
  `internal/membership/membership.go:349–382` `Roles`; a protected candidate blocks the whole batch
  (`.go:166–175`) and `internal/membership/gates.go:17–112` `CheckGates` re-evaluates reviewer
  counts, LE-7/LE-11, goal-checker and model diversity — escalating, never waiving.
- **Evidence, markers, notices, surfaces:** `quota.Evidence` (scrubbed, structured),
  `internal/quota/record.go:129–150` markers/notice, `internal/quota/notice.go:31` checked
  publication helper, `internal/app/wait.go`/`organizer.go` surfaces.
- **Policy layering:** `internal/config/runtime.go:284,321,543` presence-aware
  `quota_auto_exclude` with machine→deck merge; per-idea opt-out read at kickoff.

Scoped null: I searched `internal/runmanifest/manifest.go:34–64` and
`internal/protocol/workspace.go` for any recorded **organizer** identity — there is none. The only
CLI-identifiable non-counting seat is the declared `facilitator:` (`workspace.go:38`,
`membership.Roles`). This grounds DP3/DP4 below.

Relayed, not verified by me (from `source-context/ORGANIZER-BRIEF.md` and the prior idea's
artifacts): skill 2.15.0 state, staged-unpublished core 2.15.0, the AC2 owner waiver, the R8
residuals, the 70,000-byte packet guardrail history. Local code confirms the shipped recognizer
support table (`internal/telemetry/quota.go:18–25`: zcode `supported-owner-deviation`, codex/kimi/
claude `diagnostic-only`) and the zcode-only classification path (`.go:69–75`).

## Proposed approach

Name: **participant dropout** (`dropout.any-failure.v1`), a sibling trigger inside the shipped
automatic-exclusion mechanism. The owner's controlling words (verbatim in `00-prompt.md`): a
participant that drops out "for any reason" does not continue in that idea and can join from the
next one. The seven design points:

### DP1 — Reuse 1.51.0; the quota predicate stays the stricter sub-case

One mechanism, two triggers, one transition per batch. `quota.Evaluate` gains a rule profile: today's
quota path stays byte-identical in behavior when `participant_dropout` is off; the dropout path
reuses members, candidates, protected roles, all-or-nothing batch semantics, `Settle`, history,
lease, notices and surfaces unchanged in shape. When both knobs are on, one failure produces ONE
candidate with ONE rule id: if the terminal evidence is positively classified by the quota
recognizer (`Evidence.Eligible` with a `quota.*` rule id), that id and its reset hint are kept; every
other terminal failure carries `dropout.any-failure.v1` with reset unknown (and therefore no
relaunch suggestion — matching the shipped "unknown reset ⇒ no suggestion" rule,
`internal/quota/quota.go:82–87`). A practical side effect worth recording in FINAL: because the
dropout trigger needs no recognizer, it covers the zcode quota-exhaustion case even where the AC2
waiver left the native recognizer possibly inert; the `quota-zcode-native-exhaustion-capture`
follow-up stays as-is for the quota knob's own semantics and is not superseded.

### DP2 — "Any reason" vs a transient blip: the finite retry default

**Default: the original attempt plus exactly one relaunch (2 attempts total), fixed, not
configurable.** Mechanically, widen the existing retry-once loop (`internal/runner/runner.go:516–548`)
so the single relaunch fires for ANY terminal failure that left no valid artifact — provider error of
any kind (429/503/401), crash, hang/watchdog kill, timeout, or structurally invalid output — not just
`no_first_output`. Same inputs, same configured ceiling; the skill's longer-timeout relaunch guidance
(2400 s/3600 s) remains organizer-side tooling for protected agents and is not productized. A
participant becomes dropout-eligible only when the terminal attempt also failed with no valid
artifact and no later success in the batch (shipped precedence, `internal/quota/quota.go:135–145`,
unchanged). One carve-out preserving 1.51.0 quota semantics exactly: a terminal failure that is
positively quota-classified needs no retry — a ≥60-minute reset makes an immediate relaunch futile —
so it is candidate-eligible on the first failure. This matches the owner's own model for this run
("fails twice on the same step … is dropped", `00-prompt.md` temporary authorization) and makes the
owner's intent the product default.

### DP3 — The floor

The owner's minimum: "the implementer plus at least one other participant, and the organizer never
counts." Shipped arithmetic — at least 2 usable survivors, the declared facilitator never counting
(`internal/quota/quota.go:146–149,176–178`) — already yields exactly the owner's outcomes in the
stated fleet shape, because the implementer/designee is protected (never a candidate) and therefore
always one of the 2. One refinement is needed for correctness in batches that did not dispatch the
implementer (e.g. a signoff or review step over participants only): **count the protected
designee/pin seat structurally, not only when observed in the batch.** Concretely, for this idea
(`participants: [codex-1, kimi-1, zcode-1]`, `implementer: codex-1`, no facilitator): drop kimi-1 →
survivors codex-1 (structural) + zcode-1 → floor holds → continue with one non-blocking notice; drop
both kimi-1 and zcode-1 → only the protected codex-1 seat remains → apply nothing, ONE blocking
escalation. Below the floor, the escalation lists every candidate, the arithmetic, and the owner
options (DP-Escape below) — FINAL states them; the owner chooses. The floor stays fixed at 2 and
non-configurable (R3 pattern). The structural-count refinement is scoped to dropout decisions; the
quota rule keeps its observed-only count so 1.51.0 quota behavior is untouched.

**DP-Escape (escalation content, options presented, none chosen silently):** when the floor fails,
the blocking notice offers: (a) authorize a substitute non-implementer reviewer/implementer-side
agent for this idea (e.g. a separate claude-1 or codex-1 process) via the existing owner-confirmed
revision path (`internal/app/quota_revision.go`, `membership.Revise` with validated owner authority);
(b) an owner-confirmed reduced-quorum/solo continuation under §1; (c) pause or abandon the idea. My
recommendation is (a) when independent review still matters, because it is the only option that
preserves reviewer independence — but the notice presents all three and the owner picks.

### DP4 — Protected roles stay protected

Unchanged shipped list, extended naming: the declared facilitator/organizer, the per-idea designee,
the `IMPLEMENTATION.md` pin, and any started consensus/FINAL/review-consensus drafter
(`internal/quota/quota.go:96–102,166–175`; `membership.Roles`). A batch holding a protected
candidate applies nothing and opens the existing three-exit gate with evidence prefilled; global
`default_implementer` keeps today's fall-through (never a gate). After any drop, `CheckGates`
re-evaluates reviewer counts per track, the LE-7/LE-11 two-reviewer close, goal-checker eligibility,
`require_model_diversity` and `strict_gate` — shortfalls escalate, none is waived
(`internal/membership/gates.go:17–112`, reused as-is). Known residual, disclosed: an organizer who
is a participant, is NOT the designee, and has no declared `facilitator:` is not CLI-identifiable
(scoped null above); such a seat could in principle be dropped. FINAL should state this and
recommend that a participating pure organizer declare `facilitator:` (existing field, existing
preflight consistency gate), with `CheckGates` plus the floor escalation as backstops. I do not
propose a new `organizer:` field (see ALT-6).

### DP5 — Integrity: evidence, stopped writers, retry idempotence, retained dissent

- **Trigger boundary (anti-dissenter guard).** Candidacy arises ONLY from supervisor-observed
  terminal facts: exit code, structured failure, watchdog kind, and structural artifact validity
  (`validateArtifactForPhase`). Artifact **content and positions are never examined**: a
  structurally valid artifact that dissents, blocks or disputes is a SUCCESS, not a failure — the
  strongest possible statement that dropout can never remove a dissenter. "Unparseable or invalid
  output" means structural invalidity only. Slowness alone never triggers: a valid artifact inside
  the ceiling succeeds regardless of duration; only a watchdog/timeout kill with no artifact is a
  failure. This deliberately widens the shipped quota integrity list (which excluded hangs and
  timeouts from automation) — that is the owner's explicit "for any reason", and the form-not-position
  boundary is what keeps it safe.
- **Verbatim audit evidence per drop.** Reuse `quota.Evidence`: invocation id, adapter, rule id,
  provenance `"supervisor-terminal"`, scrubbed excerpt (failure class, exit code, watchdog kind,
  bounded scrubbed log tail — same scrubbing discipline as today), observation time; reset only when
  quota-classified. Evidence is recorded before application, inside the same durable batch record
  (shipped). Dropout provenance is trivially native — the supervisor itself is the witness of its own
  child's terminal state — so no per-adapter recognizer project is needed (contrast AC2).
- **Stopped writers.** `Settle` runs only after every affected writer stopped
  (`membership.RequireStopped`, crashed-writer settlement with host/boot/PID proof) — reused
  unchanged; a dropout never races a live writer.
- **Retry idempotence.** Each attempt already carries an `attempt_id` and immutable invocation
  identity (`internal/runner/acp.go:29–33`, `runner.go` attempt loop). One step = one terminal
  evaluation; the earlier failed invocation and any `round.incomplete`/partial artifact are preserved
  and never counted as completion (shipped AC13 semantics, extended to dropout); a valid artifact on
  any attempt wins; replay of a committed transition stays idempotent (shipped).
- **Retained dissent.** The whole shipped apparatus is rule-agnostic: excluded ids stay known
  signers, filed ❌ yields `TriageBlocked` not `TriageMalformed`, `DISPUTED` claims and open
  CRITICAL/MAJOR/strict-gate findings need explicit independent disposition, no new signoff append
  until owner-confirmed re-inclusion, re-probe at the next idea, no same-idea rejoin, no timers.

### DP6 — Default and opt-out; explicit migration from saved quota-v1 policy

- **Knob.** NEW presence-aware `[defaults].participant_dropout` (machine → deck merge, per-idea
  `participant_dropout: false` opt-out in `00-prompt.md`), default ON for ideas created after
  delivery and release — mirroring R2. `quota_auto_exclude` is untouched in name and semantics.
- **Recorded policy.** The kickoff record (`quota-kickoff.json`) and every batch carry the resolved
  policy; `quota.Policy` (`internal/quota/quota.go:18–21`) gains a third, presence-aware key
  (`dropout *bool`, omitempty). The custom `UnmarshalJSON` (`.go:24–42`) currently rejects any record
  whose fields ≠ {enabled, scope}; it is extended to accept both the 2-key legacy shape (dropout
  absent ⇒ nil ⇒ OFF) and the 3-key shape. `strictDecode`'s `DisallowUnknownFields` keeps unknown
  future keys blocking.
- **Migration, explicit.** (1) In-flight ideas with a saved 2-field quota-v1 policy decode with
  dropout absent and therefore run with dropout INACTIVE for their whole lifetime — a binary upgrade
  or resume never widens a recorded scope (shipped no-widening invariant, honored by construction).
  (2) 1.51.0 decks need no action: their next NEW idea records dropout-on from the config default;
  setting `[defaults].participant_dropout = false` (machine or deck) restores confirmation-only
  behavior for new ideas. (3) The only way to switch dropout on for an existing in-flight idea is the
  existing owner-confirmed revision path (`quota_revision.go` with validated authority) — same gate
  as widening quota scope today. (4) Malformed or ambiguous policy records fail closed (shipped
  pattern). The `reflect.DeepEqual(m.QuotaKickoff, h.Kickoff)` manifest checks
  (`internal/membership/membership.go:178–190`) are unaffected because both sides decode through the
  same extended struct.

### DP7 — Size: exact inventory

**Protocol hunks — identical in all three copies** (`parley-deck/COOPERATION.md` in this worktree;
`internal/protocol/defaults/COOPERATION.md` under `TestEmbeddedDefaultMatchesLiveDeck`; the skill's
`skills/parley-deck/references/COOPERATION.md` in worktree `participant-dropout-skill`, confirmed
present):

1. `COOPERATION.md:59` (§0 `[defaults]` sentence): add `participant_dropout` beside
   `quota_auto_exclude`.
2. `COOPERATION.md:295` (Phase 0 template): add `participant_dropout: false` opt-out line beside the
   quota one.
3. §9.0, in the automatic-exclusion bullet block at `COOPERATION.md:896+`: one sibling block
   "Participant dropout (any-failure)" covering knob/default/migration (DP6), trigger + fixed 2-attempt
   budget (DP2), floor + escalation options (DP3, DP-Escape), protected roles (DP4), evidence and
   form-not-position boundary (DP5), marker/notice reuse, return semantics (per-idea, next-idea
   re-probe, owner-confirmed re-inclusion), and the statement that the quota predicate remains the
   stricter sub-case. No new headings, so `meta/packet-applicability.yaml` needs no new entry (§5 and
   §9.0 are `include: always`); verify phase 0/5/8 packet rendering anyway.
4. §5 (`COOPERATION.md:748` region): widen "except the recorded §9.0 quota auto-exclusion" to "…§9.0
   automatic exclusion (quota or participant dropout)".
5. §4 Phase 5 designee/pin cross-reference: name both rules.
6. `meta/protocol-changelog.md`: one entry naming this idea.

**CLI files (worktree `participant-dropout`):**

- `internal/config/runtime.go` (+`runtime_test.go`): knob field + machine→deck merge (beside
  `:284`, `:543`).
- `internal/quota/quota.go` (+`quota_test.go`): Policy third key; `NewPolicy` over both knobs;
  dropout rule profile in `Evaluate` (structural backbone survivor count; same candidates/batch
  shape); quota path unchanged.
- `internal/quota/record.go`, `history.go` (+`record_test.go`, `history_test.go`): extended Policy on
  Kickoff/Batch/Transition; `Batch.Validate` branches the reset constraint by rule-id family
  (`MinimumReset` only for `quota.*`; dropout requires excerpt + failure facts, `ResetAt` nil);
  dropout marker/notice text (reset `unknown`, no relaunch hint, label `automatic dropout`).
- `internal/telemetry/dropout.go` (NEW, +`dropout_test.go`): classifier from supervisor facts
  (exit/structured/watchdog/artifact-structural-validity/failure class) → dropout `Evidence`;
  `quota.go`/`quota_zcode.go` untouched.
- `internal/runner/runner.go` (+`hardening_test.go`): widen retry-once trigger at `:516–548`;
  quota-eligible terminal evidence skips the retry; wire dropout evidence into `Result`; keep the
  `:777–786` artifact/cancellation precedence overrides.
- `internal/runner/quota.go` (+tests): settle under a dropout-active policy; dropout floor path.
- `internal/membership/membership.go` (+tests): dropout floor + escalation text with DP-Escape
  options; `Roles`/`RequireStopped`/`Settle` shapes unchanged.
- `internal/app/preflight.go`, `preflight_liveness.go` (+tests): kickoff ping attempted twice per
  agent under dropout-on; dropout evaluation in the `:412–436` block; non-dropout gates and `--yes`
  semantics unchanged.
- `internal/app/app.go` (`:1946`, `:2483`): resolve both knobs; pass through to `runcontrol.Create`.
- `internal/runcontrol/runcontrol.go`, `internal/runmanifest/manifest.go` (+tests): record extended
  policy at kickoff (manifest carries `quota.Kickoff`, so it follows automatically).
- `internal/app/wait.go`, `organizer.go`, status (+tests): show dropout exclusions with reset
  `unknown` and pending states.
- `internal/app/driver_impl.go:827`, `driver_consensus.go:148–155`, `consensus_request_signoffs.go:203`
  (+tests): the single-agent settle paths accept dropout evidence.
- `internal/app/quota_revision.go` (+tests): verify the owner revision path can set `dropout=true`
  for an in-flight idea with validated authority (plumb the field; no new authority class).

**Skill worktree (`participant-dropout-skill`):** bundled `references/COOPERATION.md` (byte-identical
to the deck copy); `SKILL.md`/`ROSTER_AND_PROTOCOL.md` one short note each: the organizer reads the
dropout notice, never re-includes silently, and never infers a dropout itself.

**Tests (beyond the per-file units above):** policy layering + legacy 2-field record migration
decode (dropout-off, no widening on replay); dropout classifier matrix (exit≠0; exit 0 + invalid
artifact; watchdog/timeout kill; hang; auth 401; 503; **structurally valid dissenting artifact =
success** — the dissent-preservation adversarial fixture; assistant-quoted text alone never
supervisory); retry widening (any failure retried once; quota-eligible skips; valid attempt-1
artifact wins; killed-by-user never retried); floor permutations (3→2 applies; 3→1 and 2→1 block
with exactly one escalation; designee-not-dispatched structural count; facilitator never counts);
both-knobs single transition with quota rule id + reset hint; kickoff ping double-failure drop;
protected roles incl. started drafter; retained ❌/`TriageBlocked`/`DISPUTED`/findings (existing
suite re-run); status/wait/brief agreement; drift test; phase 0/5/8 packets; skill byte-compare.

**Estimate:** ~350–500 product LOC plus ~700–1000 test LOC and ~40–60 changed protocol lines per
copy. UNVERIFIED (my estimate; §15 — I cannot verdict my own sizing). It is small because the
transition machinery is reused, not duplicated; the genuinely new logic is the classifier, the retry
trigger, the floor variant, and the policy plumbing.

### Observable acceptance criteria (proposal for FINAL)

1. **AC-D1 (trigger).** With dropout on, a dispatched non-protected participant whose step fails
   twice on the same step (any terminal class: 429, 503, 401, crash, watchdog, timeout, invalid
   artifact) is removed from this idea's `participants:` via one durable batch transition; a
   positively quota-classified failure drops it on the first failure with the quota rule id and
   reset hint.
2. **AC-D2 (retry).** Exactly one relaunch is attempted per failed step (attempt ids recorded); a
   valid artifact on either attempt is success; a user kill is never retried; replay/inspection shows
   the preserved first failure.
3. **AC-D3 (valid dissent is safe).** A structurally valid round artifact whose content dissents,
   blocks or disputes is NOT a failure and its author is never a dropout candidate (adversarial
   fixture).
4. **AC-D4 (floor).** 3 participants with 1 candidate → applied, one non-blocking notice;
   3→1 (or 2→1) → nothing applied and exactly one blocking escalation listing candidates, arithmetic
   and the three DP-Escape options; permutations are order-independent; the declared facilitator
   never counts; the designee's seat counts structurally even when not dispatched.
5. **AC-D5 (protected).** A batch containing designee, pin or started drafter applies nothing and
   opens the three-exit gate with evidence prefilled; post-drop `CheckGates` shortfalls (reviewers,
   LE-7/LE-11, goal-checker, model diversity, strict gate) escalate.
6. **AC-D6 (retained dissent).** After a dropout, the dropped id's ❌ yields `TriageBlocked`; its
   `DISPUTED` claims and findings persist with explicit-disposition duty; it cannot append signoffs;
   it re-probes at the next idea and rejoins this idea only via owner-confirmed revision.
7. **AC-D7 (evidence).** Every drop carries invocation id, adapter, rule id, supervisor provenance,
   scrubbed excerpt with failure class/exit/watchdog kind, UTC times — recorded before application,
   inside the durable batch record.
8. **AC-D8 (knob/migration).** Per-idea `participant_dropout: false` and deck override work;
   malformed config fails closed; a saved 2-field quota-v1 policy decodes dropout-off and no resume
   or upgrade widens it; owner revision can enable it for an in-flight idea.
9. **AC-D9 (quota untouched).** With dropout off, every 1.51.0 quota behavior is unchanged (existing
   quota suite passes unmodified); with both on, one failure produces one transition.
10. **AC-D10 (kickoff).** Under dropout-on, a readiness ping failing twice (any reason, floor
    holding) drops the id before the first `participants:` line/manifest/dispatch; standalone
    `parley preflight` still never applies.
11. **AC-D11 (surfaces).** `status`, `wait` and the organizer brief show dropout exclusions with
    reset `unknown` (or the quota hint), the same current set, and pending states; `wait` keeps its
    exit contract.
12. **AC-D12 (protocol + checks).** AC-D-checkable protocol hunks present in all three copies; drift
    test, phase 0/5/8 packet rendering, skill byte-compare and the affected-package Go tests pass on
    the current tree.

### Owner decision

None required, in my assessment. Every point the brief enumerates is settleable from the owner's two
verbatim messages plus shipped mechanism; DP-Escape presents options at runtime rather than choosing
at design time, and the undeclared-organizer residual is disclosed with a recommendation rather than
a new field. I recommend FINAL proceed without a proposal stop, per the brief's "no unnecessary
owner gate"; ratification of the protocol text (§7) and the attended publish remain the owner's
standing steps, not new decisions.

## Concerns / open questions

1. **Preflight ping retry shape (open, implementation-verifiable).** I did not read
   `preflight_liveness.go` in enough depth to confirm the ping layer can cleanly attempt a second
   probe per agent; DP-D6/AC-D10 assume a 2-attempt ping budget. If the ping is single-shot by
   construction, the kickoff drop may need the budget implemented at the `preflight.go` loop level —
   same semantics, slightly different hunk. Not design-blocking.
2. **Packet guardrail headroom.** The prior idea hit the 70,150/70,000-byte facilitator-packet limit
   and had to shorten new §9.0 prose (prior IMPLEMENTATION.md, "Packet-size correction"; relayed,
   not re-measured by me). The new §9.0 block should be written tight and the guardrail re-measured
   in implementation.
3. **Timeout-ceiling interaction.** A participant whose configured ceiling is systematically too
   short will fail twice and drop even though it is healthy. The retry absorbs single blips only.
   Acceptable under "any reason" (and the notice names the class), but FINAL's known-risks should say
   it plainly; per-agent timeout configuration is the owner's existing remedy.
4. **`Batch.Validate` rule-id branching.** The reset constraint must key off the rule-id family, not
   a boolean, or a hand-edited record mixing families could smuggle a sub-60-minute quota reset past
   validation. Implementation detail; called out so the reviewer checks it.
5. **Both-drop discovery order.** If kimi-1 and zcode-1 fail in DIFFERENT steps (not one batch), the
   first drop applies (floor holds), the second blocks at its own settle — two notices total (one
   non-blocking, one blocking). That is per-batch semantics working as shipped; confirming this
   reading with codex-1/kimi-1 in round 2 would be good.

## Risks

- **Dissent-removal perception.** The highest-severity hazard of any auto-drop rule. Mitigated by the
  form-not-position boundary (AC-D3), supervisor-only provenance, retained-dissent machinery (AC-D6)
  and the adversarial fixture; residual risk is a reviewer dispute over "structural" artifact
  validity, which the existing validator definitions already pin down.
- **Transient fleet-wide events drop too much.** A gateway blip can fail both participants; the floor
  catches it (apply nothing, escalate once) — degrades to today's behavior, no loss.
- **Migration grammar change.** Touching the strict policy decoder risks breaking saved records;
  mitigated by decode-both-shapes tests, DisallowUnknownFields kept, and fail-closed on ambiguity.
  This is the one hunk that edits shipped validation code rather than adding beside it.
- **Correlated agreement (§15.6b).** codex-1 (gpt-6-astra), kimi-1 (kimi k3) and I (glm-5.3) are
  three families, but we all inherit ONE shipped mechanism, so convergence here is family-adjacent,
  not independent evidence of safety. What would prove the design wrong: a content-triggered drop, a
  lost veto, a drop below floor without escalation, or a resume that widens a saved policy. The
  fixtures named in DP7 target exactly those.
- **Size drift.** If implementation discovers the floor or ping-retry work is larger than estimated
  (my LOC estimate is UNVERIFIED), the honest move is a staged delivery like 1.51.0's stage 1/2 —
  kickoff-only first — rather than silent scope growth; the brief's "as small as possible" wants the
  staged fallback stated in FINAL if it triggers.
- **Delayed benefit realism.** Like 1.51.0 (whose recorded incidents would not have automated), the
  benefit is prospective; unlike 1.51.0, the any-failure trigger needs no recognizer, so the rule is
  live for every adapter from day one — the practical gap that AC2 left is closed for the dropout
  knob.
