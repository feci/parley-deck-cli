---
idea: meta-protocol-change-quota-auto-exclude
review-cycle: 2
outstanding_agreed_fixes: 5
blocked: false
drafted-by: codex-1
date: 2026-10-05
reviewed-commit: d8b729af187a18ebe3891a36dfc379c4b04c8b95
product-commit: 906857b9af33306158d3b0e91b94e6ef4b6f7458
skill-commit: dcb7d593130247709cfad19715aa1f53e5cc8be0
---

## Scope and review basis

The separate claude-1 round-04 review is the binding full-scope re-review after fix-up cycle 1.
It completed the single owner-authorized relaunch without quota/auth failure and found 0 CRITICAL,
3 MAJOR, 2 MINOR and 1 NIT. Both earlier CRITICALs and most G1–G9 are resolved within the reviewer's
stated scopes. codex-1 accepts the findings as implementer; it supplies no independent code verdict.
This draft authorizes cycle 2 only after both participants sign. It authorizes no close or release.

Trajectory: 15 round-03 findings became six round-04 findings after one fix-up. Both CRITICALs are
resolved. R4-MAJOR-1 concerns new retry-framing fix code; R4-MAJOR-2 is the remaining G2/R2 compatibility
issue; R4-MAJOR-3 is the previously disclosed native-evidence gap. Continue the bounded unrelated repairs
under Phase 8 stopping judgment, rather than treating either green tests or a pass count as acceptance.
Cycle 2 of the five-cycle cap is proposed. No D6/legacy-accounting change, roster change or excluded-agent
invocation is included.

## Agreed fixes

- **G10 — real retry semantics and actionable recognizer reasons** (round-04/claude-1 R4-MAJOR-1,
  R4-NIT-1): accept complete allowlisted retry reports whose exhaustion class and absolute machine reset
  agree within the owner's one-second bound even when their countdown text changes. Check each record's
  internal duration/header/countdown consistency, distinguish its observation from one aggregated dump's
  receipt time, and apply the fixed 60-minute threshold against the actual observation. Keep exact
  lastError fingerprint equality and all mixed-status/class/contradiction, artifact and later-success
  gates. Add the two retained countdowns and realistic backoff deltas as explicitly source-derived
  fixtures. Return the specific framing/consistency rejection instead of calling a supported recognizer
  unsupported. No change to the owner exception, no relabeling of fixtures as native.
- **G11 — preserve the ordinary knob-off path** (R4-MAJOR-2 option **(b)**, R4-MINOR-1): restore the
  pre-existing protocol's recorded confirmation and catch-up forms for policy-off ideas. Do not require
  a new included marker, a new exact directive, a mandatory committed answer schema or a new command
  for this existing path. Preserve explicit owner-confirmed re-inclusion, late round-1/read-priors/join
  from round-2 catch-up, historical-known/current-required membership, retained obligations and immutable
  revisions. Do not treat an unconfirmed prompt edit as owner authorization. Use the pre-change path
  and existing protocol forms as the compatibility oracle. Make exclusion parsing tolerate reasons
  containing an em dash and report the accepted grammar/action on malformed input. Document the real
  accepted forms in CLI and skill guidance. Choosing (b) avoids proposing a third knob-off behavior
  change; if preserving this contract proves infeasible, report the precise blocker before changing it.
- **G12 — visible policy-on revision/recovery boundary** (R4-MINOR-2): name parley quota revise in the
  contradictory-projection diagnostic. A prompt edit after a fully applied authoritative revision must
  not be silently overwritten merely because it matches an older revision. Distinguish a genuinely
  interrupted pending projection (normal checked recovery) from a later manual edit; the latter gets
  an owner-visible escalation with the required correction/revise path. Do not invent authority from
  the edited prompt or break idempotent pending-transition replay. Clarify the policy-on path in the
  protocol/skill without broadening owner policy.
- **G13 — close the recovery leads with evidence** (round-04 open questions 2 and 3): reproduce a writer
  whose parley process crashed before terminal.json, and manual import with a missing kickoff manifest.
  Establish and document a checked recovery route within FINAL's durability/recovery requirements, or
  fix a confirmed defect with regression coverage. Never infer that an unknown/foreign writer stopped
  from local PID absence, and never dispatch with an unresolved writer. Validate prerequisites before
  immutable commit where possible; already committed incomplete state remains fail-closed until checked
  recovery. Report any action requiring owner judgment rather than adding an unsafe escape hatch.
- **G14 — prepare native-framing evidence; retain the AC2 owner gate** (R4-MAJOR-3 and open question 1):
  locate the actual installed zcode terminal-error sink and inspect settings from source; use local
  source inspection and offline SDK/Node reproductions only. Label source-derived framing, partial native
  tails and complete native capture distinctly. Do not invoke zcode, a provider or another participant,
  loosen the allowlist, or declare AC2 met. The unrecovered complete native capture remains an OPEN
  owner decision after the authorized repair work is concrete. Either owner-authorized native capture
  or an explicit owner ruling is required; no assumption or signoff here substitutes for it.

## Owner-dependent thread

R4-MAJOR-3 / AC2 remains unresolved and is not deferred out of this idea or waived. Only its native-capture
thread is held; unrelated G10–G13 work may proceed after this fix-list consensus under Phase 8 stopping
judgment. Before requesting a decision, finish the authorized repairs and source/evidence preparation
so the owner sees a concrete result. The later blocking note must state exact capture scope or exact
proposed evidence ruling and the remaining risk. No zero-fix/close claim can be based on this draft.

## Deferred follow-ups

- D6 / driver gap 11: separate lasting legacy-accounting fix, follow-up idea TBD after this idea.
- Windows runtime remains unexecuted; compile-only evidence is not a runtime pass. CLI winget stays held.
- Unsupported-adapter native provenance and the Claude JSON-output change remain FINAL's follow-ups.

## Dismissed findings

None of the new round-04 findings is dismissed. The reviewer's prior dispositions stand: G1, G3,
G5–G9 and R1/R3 are resolved in the scopes stated in round-04; G2/G4 and R2 retain the specific new
residues above. Completed-drafter protection and literal AC1 skill/deck equality retain their signed
round-03 dispositions. No deferred item is used to erase a current requirement.

## Coverage & blind spots

claude-1 read the entire current protocol body, all 24 CLI product chunks and both skill product chunks,
and filed coverage over every changed product file. Its 13-package focused checks and actual shared/local
probes are PRIMARY within their recorded scopes. Full host suite/race/vet/skill checks are producer
SECONDARY evidence; Windows runtime and complete native zcode capture remain unexecuted/unavailable.
There is one independent reviewer. codex-1 organizes and implements under the recorded §15.5 role
concentration. A new full-scope re-review, final review consensus with both signoffs and a NEW attended
owner close are still required. No acceptance of code is inferred from agreement to this fix list.

## Signoffs

<!-- Each participant appends only its own block. Signoff accepts this fix list, not the current code. -->

### Signoff: codex-1 — 2026-10-05
Status: ✅ ACCEPT
Notes: I accept G10–G14 as the cycle-2 implementer plan. G11 uses option (b) to preserve existing knob-off behavior without a new owner-policy deviation. AC2/native capture stays open for owner judgment after concrete repairs and source evidence. This is not a code verdict or permission to close, merge or release.

### Signoff: claude-1 — 2026-10-06
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: I accept G10–G14 as the cycle-2 fix list, subject to reservations V1–V9 below; the cycle-2 re-review checks each one. G11's option (b) is my own R4-MAJOR-2 option, and I agree with it. AC2/native evidence stays an OPEN owner decision. This signoff waives no native evidence, approves no code, and authorizes no close, merge or release.

This is the independent Phase-7 signoff of the only non-implementer reviewer. It is timeout relaunch 1 of 2 of the
same step, with a 2400 s ceiling; the product snapshot and the fix draft are unchanged. No quota, credit or
authentication error occurred in this process. It made no retries, probed no provider, invoked nobody, and edited
no product, draft or other agent's text.

**Protocol context.** `context_mode=full`,
`source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`, and
`fallback_reason` is absent. I used the attached attestation. I did not re-render, and I used no shadow or
optimized packet.

- Body: `.parley-runtime/protocol-packets/full-phase7-deliberation-73613f95….md`.
- PRIMARY: its SHA-256 matches, and `cmp` against the deck `COOPERATION.md` exits 0 (1,501 lines, 125,862 B).
- **Actual coverage: the entire body**, in four bounded chunks (lines 1–380, 381–760, 761–1140 and 1141–1501).
  The three over-long lines (59, 452 and 1002) were displayed in full, and I checked their endings separately.

**Snapshot (PRIMARY).**

- HEAD is `914d072`. `git diff --name-only 906857b HEAD -- internal cmd docs go.mod go.sum
  parley-deck/COOPERATION.md parley-deck/meta` is empty, and `5c7da6c..HEAD` touches only owner-answer and
  orchestration files.
- The CLI product is `906857b`. The skill worktree is at `dcb7d59` and clean.
- Before this append, this file's SHA-256 was `82fddfef…235c`.
- This signoff is not a code re-review.

**What I read.**

- In full: this file; my `review/round-04/claude-1.md`; the scope-reset, review-quota, quota-standing and
  timeout-standing owner answers; and the ratification note.
- In part: FINAL lines 1–514 and 670–873 (§1–§14, AC1–AC21, Idempotence & recovery, Known risks), and the Phase
  7/8 text with its stopping judgment in the packet. Not read this run: FINAL lines 515–669 (Purpose and Context).
- Organizer leads only, not verdicts: `source-context/fix-consensus-2-*` and `organizer-notes.md:900-910`.

**Erratum.** I appended a dated erratum to my round-04 review. My quote reads `IMPLEMENTATION.md`; the owner's
answer reads `IMPL-ORGANIZER-BRIEF.md`. The original line is kept. The file's SHA-256 went from `499afc78…3453`
to `8b4093df…afcc`. No verdict changed.

#### User direction

Below are verbatim quotes of both owner answers, without their frontmatter. The quoted question and option are
in Slovak; the translations are as relayed. I apply them as written. The organizer, not this process, manages
relaunches: up to 3 quota relaunches per step after the reset plus 2 minutes, for a stated reset of 60 minutes or
less, and up to 2 timeout relaunches per step. A terminal quota, credit or auth error stops this process with the
verbatim error and no retry.

From `parley-deck/inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_quota-standing-permission.md`
(SHA-256 `593512a3…7dc04`):

> ## Owner answer to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_fix-consensus-quota-stop-20261005.md`
>
> Relayed by the owner's Claude Code session. Received 2026-10-06 at about 12:18 CEST. The question and the
> answer, verbatim (Slovak), with a translation:
>
> > Question: "Zopakované review prešlo: claude-1 našiel 0 CRITICAL, 3 MAJOR, 2 MINOR a 1 NIT a codex na ne
> > pripravil 5 opráv. Podpis plánu opráv od claude-1 však o 21:57 znova spadol na „429 All claude accounts
> > have exhausted their quota (reset after 5m)“. Claude pool na gatewayi sa opakovane vyčerpá na pár minút a
> > každá takáto zastávka teraz čaká hodiny na teba. Ako ďalej?"
> > Selected: **"Trvalé povolenie: počkať a opakovať (Recommended)"**. The option read: "Pri chybe kvóty claude
> > s resetom do 60 minút codex počká do resetu plus 2 minúty a krok zopakuje, najviac 3-krát na jeden krok.
> > Až potom sa zastaví a ozve sa. Platí len pre túto ideu a nikoho to nevyradí."
>
> Translation: "Standing permission: wait and retry. On a claude quota error with a reset within 60 minutes,
> codex waits until the reset plus 2 minutes and repeats the step, at most 3 times per step. Only then does it
> stop and report. This applies to this idea only and excludes nobody."
>
> ## Standing rule for the rest of this idea (it replaces the brief's "do not spin retries" for this case only)
>
> - **When it applies.** The claude-1 participant (review, signoff or channel verification) fails with a
>   quota or rate-limit error that carries a stated reset of **60 minutes or less**. Example: `429 … All claude
>   accounts have exhausted their quota … (reset after 5m)`.
> - **What to do.** Wait with blocking tool calls until the stated reset plus 2 minutes, then relaunch **the
>   same step** with the same inputs.
> - **Limit.** At most **3 relaunches per step**. After the third failure, or on any error without a reset of
>   60 minutes or less, or on an auth or credit error, stop and write a blocking note as before.
> - **Records.** Record every attempt in `organizer-notes.md` with the verbatim error.
> - **Never** change the model or provider, substitute a reviewer or exclude an agent. The same rule applies if
>   codex-1's own model hits a short quota window.
>
> ## Now
>
> The claude-1 Phase-7 signoff on the cycle-2 fix plan in `review/consensus.md` is authorized under this rule.
> Then continue with fix-up cycle 2, re-review, the attended-close request and the release, all as
> `IMPL-ORGANIZER-BRIEF.md` says.

From `parley-deck/inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_timeout-standing-permission.md`
(SHA-256 `fc5b5f95…a2613`):

> ## Owner answer to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_fix-consensus-timeout-20261006.md`
>
> Relayed by the owner's Claude Code session, 2026-10-06, about 13:26 CEST. Verbatim (Slovak), with a
> translation:
>
> > Question: "Podpis claude-1 tentoraz nespadol na kvótu, ale do 20-minútového limitu bez jedinej chyby
> > nevrátil nič. Codex sa preto znova zastavil a pýta sa. Rozšíriť trvalé povolenie aj na takéto timeouty,
> > aby sa na teba pri každom nemuselo čakať?"
> > Selected: **"Áno, aj timeouty (Recommended)"**. The option read: "Pri timeoute codex krok zopakuje s dlhším
> > limitom (40, potom 60 minút), najviac 2-krát na jeden krok. Rovnaký postup pri timeoute už predpisuje
> > skill. Kvóty platia ako doteraz. Zastaví sa iba pri chybe autentifikácie alebo kreditov, alebo keď vyčerpá
> > pokusy."
>
> Translation: "Yes, timeouts too. On a timeout, codex repeats the step with a longer limit (40, then 60
> minutes), at most twice per step. The skill already prescribes this for timeouts. Quota handling stays as
> before. It stops only on an auth or credit error, or when it runs out of attempts."
>
> ## Standing rule, extending `…_quota-standing-permission.md` for the rest of this idea
>
> - **Silent timeouts.** If any participant step (claude-1 or codex-1) times out with no provider error,
>   relaunch the same step with the same inputs and a longer ceiling: first **2400 s**, then **3600 s**. That
>   is at most **2 timeout relaunches per step**. Record each attempt in `organizer-notes.md`.
> - **Short quota windows.** The rule in `…_quota-standing-permission.md` still applies unchanged: a reset of
>   60 minutes or less, wait for the reset plus 2 minutes, at most 3 relaunches per step.
> - **When to stop and write a blocking note.**
>   - an auth or credit error;
>   - a quota error with no reset, or a reset longer than 60 minutes;
>   - the attempts for that step are used up;
>   - any decision that changes scope, FINAL, the quorum or the owner's prior rulings.
> - **Never** change the model or provider, substitute a reviewer or exclude an agent.
>
> ## Now
>
> Relaunch claude-1's Phase-7 signoff on the unchanged cycle-2 plan with a 2400 s ceiling. Then continue with
> fix-up cycle 2, re-review, the attended-close request and the release, all as `IMPL-ORGANIZER-BRIEF.md`
> says. Routine procedural steps that these rules or the skill already cover need no owner question.

#### Evaluation of G10–G14 (my own judgment)

codex-1's ACCEPT records its agreement, as implementer, to carry out this list. It is not an independent verdict.
Each item below is my own evaluation against my round-04 findings, FINAL and the owner's rules. No prior
disposition limited what I checked.

- **G10. Accept, with V1–V3.** It is my R4-MAJOR-1 fix and stays inside the owner's text. Records are compared by
  exhaustion class and absolute machine reset within one second, never by human countdown text. The
  `lastError` fingerprint, every mixed/contradiction gate and the artifact/later-success gates stay. Fixtures
  stay labeled source-derived. It also covers R4-NIT-1.
- **G11. Accept, with V4–V5.** Option (b) is my own R4-MAJOR-2(b). It removes the unlisted third knob-off behavior
  change without a new owner decision. FINAL §11 allows exactly two such changes: C1 and the bare 503.
- **G12. Accept, with V6–V7.** It is my R4-MINOR-2 fix: the diagnostic names `parley quota revise`, and a later
  owner edit gets an escalation instead of a silent revert.
- **G13. Accept, with V8.** It answers my open questions 2 and 3. It stays fail-closed and adds no unsafe escape
  hatch.
- **G14. Accept, with V9.** It answers my open question 1, keeps AC2 OPEN and does not loosen the allowlist.
- **Owner-dependent thread, deferrals, dismissals. Concur.**
  - Holding only the AC2 native-capture thread matches Phase 8: "pause that finding's thread until the operator
    answers; unrelated fixes may continue".
  - The deferred items match FINAL's follow-ups.
  - "Dismissed findings: none" matches my round-04 review.
- **Stopping judgment. Concur with cycle 2 of the five-cycle cap.**
  - Findings fell from 15 to 6, and both CRITICALs are resolved. One MAJOR (R4-MAJOR-1) landed on cycle-1 fix
    code.
  - If the cycle-2 re-review finds a new CRITICAL or MAJOR on cycle-2 fix code, above all in G10's framing or
    G11's membership path, treat it as churn. Stop and escalate with the trajectory; do not open cycle 3 by
    default.

#### Reservations V1–V9 (conditions the cycle-2 re-review checks)

- **V1 (G10), threshold anchor.**
  - Measure the 60-minute threshold from the latest observation of the failing invocation, which is the
    terminal record's receipt. Never measure it from an earlier record's inferred observation.
  - A record's inferred observation is its `reset_at` minus its `retry_after` or countdown. Use it only to check
    that record's internal consistency.
  - Why: anchoring on an earlier attempt can qualify a reset that is already under 60 minutes away at receipt,
    which breaks the fixed constant (R3).
  - Required negative fixture: the first record is 61 minutes before the reset and the terminal receipt is 59
    minutes before it. The result must be ineligible.
- **V2 (G10), sanity of inferred observations.** Inferred observations must not decrease in record order and
  must not be later than the receipt, within the one-second tolerance. Where the invocation's start time is
  recorded, none may be earlier than it. A violation gates as contradictory.
- **V3 (G10), the native shape.**
  - Cover several top-level records parsed from one captured stderr with a single receipt clock. That is my
    round-04 P2 case C; the fix must not cover only the nested RetryError dump.
  - PRIMARY, a bounded `grep`: the installed `zcode-app-cli/vendor/zcode.cjs` (SHA-256 `3e3433d9…685f`)
    contains `allowSystemInMessages:!0,maxRetries:0` twice. The organizer's G14 lead attributes that option set
    to `w9r`. On that path, the SDK's own RetryError aggregate may therefore not be the native shape.
  - UNVERIFIED and left to G14: which path produced the incident's errors, and how zcode's own retries emit
    records.
  - Let G14's sink and retry findings decide which shapes G10 accepts, and label each fixture's shape provenance.
- **V4 (G11), the reading under which I accept.**
  - The compatibility oracle follows the existing protocol. Where it documents a recorded form, that form is the
    oracle: the §9.0 `excluded:` marker, in G11's tolerant grammar. Where it documents none, the pre-change CLI
    path is the oracle: a known agent's return by a `participants:` edit, and a §5 catch-up join (late round-1
    plus the `participants:` edit).
  - A differential test against the pre-change code path checks these forms.
  - A change accepted this way enters immutable history as a manual revision. It is never labeled
    owner-confirmed.
  - On its own, such a change never satisfies a gate that the ratified §5/§9.0 text ties to owner-confirmed
    re-inclusion, such as a re-included author withdrawing a retained ❌ (G6). That keeps the existing route: a
    quoted owner ruling or committed owner authority.
  - If the implementer would instead require any further record for return or join on a policy-off idea, that is
    a knob-off behavior change. It falls under G11's own "report the precise blocker" clause and is not left to
    the implementer.
- **V5 (G11), grammar.**
  - Anchor the em-dash tolerance on the trailing `— confirmed YYYY-MM-DD` (my R4-MINOR-1).
  - The error prints the accepted grammar and the action to take.
  - The docs list the variants that are still rejected (a note after the date, hyphen separators), so a
    rejection is never a surprise.
- **V6 (G12), the discriminator.**
  - Tell an interrupted projection from a later edit by durable transition state: the latest committed
    revision's applied/receipt state. Never decide from prompt content alone.
  - If that state is unreadable and the prompt matches neither the pending revision's before-set nor its
    after-set, escalate instead of rewriting.
  - Replay keeps G7's single notice and single terminal evaluation.
- **V7 (G11/G12), governance of the protocol text.**
  - The owner ratified the §9.0 text on 2026-10-04, and it is not yet published.
  - Put command-level guidance in the CLI docs and the skill. That covers `parley quota revise` and the accepted
    marker grammar.
  - If any `COOPERATION.md` wording changes anyway, it must stay normative-neutral and byte-identical across the
    deck, embedded and skill copies (AC1).
  - It must also be listed as a post-ratification text change in `IMPLEMENTATION.md` and in the attended-close
    request, so the owner sees the exact text before `parley protocol publish`.
- **V8 (G13), crashed writers.**
  - Reuse the R3 predicate: only a proven-dead local owner (same host and boot, dead PID) may be settled.
  - A record without host or boot identity stays fail-closed, with a route the owner can see.
  - Each settlement leaves a durable record (the invocation, the proof and the time), so a crash settlement can
    be told apart from a normal terminal. It is idempotent.
- **V9 (G14), the evidence boundary.**
  - Run offline reproductions against the AI SDK, or against functions extracted into an isolated script, with
    no network and no provider credentials.
  - Never run the zcode CLI entry point or load its configuration. That would be a zcode invocation, which needs
    owner authorization.
  - The later AC2 question must state which native-shaped inputs the post-G10 recognizer accepts and rejects, so
    the owner rules on the behavior that ships.

**Non-blocking notes on the draft's accuracy.**

- "13-package focused checks": 13 packages were invoked, but `internal/runstate` had no matching tests in my T2,
  so 12 actually ran tests.
- `00-prompt.md` declares `facilitator: codex-1` with `facilitator_participates: true`, and codex-1 drafted this
  file. The draft records the role concentration in one line. It has no `## Drafter position changes` section,
  and neither did the signed round-03 review consensus. If §15.5 is read to cover `review/consensus.md`, add that
  section to the next review consensus; it is likely `None`.

**Boundary.** codex-1 organizes and implements; I am the only independent reviewer. Still mandatory:

- the full-scope cycle-2 re-review;
- review consensus with both signoffs;
- the owner's AC2 decision (native capture or an explicit ruling);
- a NEW attended owner close.

There is no unattended close. No acceptance of code is implied.
