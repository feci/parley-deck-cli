---
idea: meta-protocol-change-quota-auto-exclude
review-cycle: 3
outstanding_agreed_fixes: 3
blocked: false
drafted-by: codex-1
date: 2026-10-06
reviewed-commit: e7bf96c2f8914fea151ed1b862ba658b6f3852b9
product-commit: 0ee18889977a29d6baf53d3016b501112760ae12
skill-commit: e2f3649eb938e870367c76943b441382fe7acf65
---

## Scope and review basis

The binding round-05 review reports 0 CRITICAL, 2 MAJOR and 1 MINOR after cycle 2.
The owner has answered its trajectory stop: one narrow cycle 3 for R5-MAJOR-1 and
R5-MINOR-1, and AC2's native-positive evidence waived for this release with an explicit
known limitation and linked follow-up. This draft proposes only that work. It authorizes
implementation after both signoffs; no code verdict, attended close or release is inferred.

Role concentration (§15.5): codex-1 organizes, implements and drafts this consensus;
claude-1 is the sole independent reviewer and signs in a separate configured process.

## Agreed fixes

- **G15 — policy-off CLI catch-up and kickoff-excluded return** (round-05/claude-1
  R5-MAJOR-1, owner Q1): restore a narrowly scoped CLI route for a policy-off joiner to
  author only its own late round-01 artifact, even before the participant edit. If the
  edit comes first, show pending catch-up with an actionable command/path instead of a
  raw file error; do not grant known-signer status, owner authority or completed quorum
  while the required artifact is missing. Restore a kickoff-excluded agent's plain
  participant-edit return during round 1: known only from the manual revision, never
  retroactively at kickoff. Preserve policy-on membership/identity/protected-role gates,
  retained veto/finding force and immutable history. Manual revisions cannot authorize
  withdrawal of a retained veto. Run actual CLI-dispatched local stubs and differential
  probes against pre-change 27e42b8 on shared and local filesystems; do not pre-create the
  artifact whose dispatch is under test. Rejoining after round 1 retains §5 catch-up duties.
- **G16 — one durable notice after archival** (R5-MINOR-1, owner Q1): replay must not
  re-publish an already delivered transition notice after it is archived. Preserve
  checked recovery if publication was interrupted, durable/idempotent application and
  one terminal evaluation. Cover applied-receipt and archive-before-receipt recovery,
  including repeated Before calls on both filesystems; ambiguous state stays gated.
- **G17 — disclose the owner's AC2 waiver and link the evidence follow-up** (R5-MAJOR-2,
  owner Q2): record R5-MAJOR-2 as owner-accepted and deferred, never fixed or AC2 passed.
  State in the skill, support guidance and release notes that zcode auto-exclusion may
  not fire on real native output and unrecognized failures use the owner-confirmed path.
  Open `quota-zcode-native-exhaustion-capture`, whose grammar work requires a captured
  real native zcode exhaustion with invocation/exit/start/receipt facts and scrubbed
  complete streams. Do not invoke zcode, relax its parser or relabel fixtures now.
  Defer new private raw-stderr retention to that follow-up: existing diagnostics are not
  a verified complete, privately scrubbed capture facility, and adding that privacy and
  lifecycle contract exceeds this narrow repair. No protocol wording change is needed;
  the skill wording fulfills Q2 and leaves all three normative copies unchanged.

## Deferred follow-ups

- **R5-MAJOR-2 / AC2:** owner-accepted known limitation, linked to
  `../quota-zcode-native-exhaustion-capture/00-prompt.md`. No native capture is claimed.
  Source-derived positives and partial native excerpts keep their current labels.
- D6 / legacy driver accounting remains a separate owner-requested idea after this one.
- Windows runtime remains unexecuted; compile evidence is not a runtime pass. CLI winget
  stays held. Other unsupported-adapter provenance work remains as in FINAL.

## Dismissed findings

None. The AC2 owner waiver disposes of its release gate without withdrawing the finding.
Prior reviewer dispositions remain as recorded, including G11/V4's incomplete status until
G15 is independently verified.

## Coverage & blind spots

Round-05 covers every changed product file and focused checks on both volumes. Full host
suite/race/vet/skill checks were producer evidence, not independent reviewer execution.
The cycle-3 full-scope re-review should run these broad checks and a real-process crash
exercise where feasible; it must name any remaining limits. Container/PID-namespace
reasoning from round-05 is unverified, not a new proven defect. The two standing retry
permissions change only failed-invocation handling, never product scope or quorum.

A fresh CRITICAL or MAJOR on cycle-3 fix code requires another trajectory escalation.
Cycle 3 is one authorized cycle within the five-cycle ceiling, not automatic permission
for cycle 4. Both final signoffs and a NEW attended-close owner answer remain mandatory.

## Drafter position changes

No material change from my latest design round (`round-03/codex-1.md`) is independently
proposed here. Owner rulings supersede the frozen design as already recorded. Relative to
my cycle-2 signoff (`review/round-04/consensus.md`, “AC2/native capture stays open for owner
judgment”), the owner's new Q2 explicitly supplies that judgment: AC2 is waived for this
release, with the inert-support risk accepted and deferred. Relative to my round-05 owner
note's Q2 recommendation (native capture plus bounded grammar fixes), the owner selected
release with a limitation, so this plan includes neither invocation nor grammar change.

## User direction

The new owner answer is quoted verbatim below as requested; original Slovak and the relay's
English translation are retained. Earlier scope-reset and both standing permissions remain
binding; none is silently replaced by this plan.

> ## Owner answers to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_round05-trajectory-ac2.md`
>
> Relayed by the owner's Claude Code session on 2026-10-06 at about 15:50 CEST. The questions and answers
> below are verbatim (Slovak), each followed by a translation.
>
> **Q1.** "Po druhom kole opráv zostali 3 nálezy (predtým 15, potom 6). Prvý je nový MAJOR: pri vypnutom
> nastavení oprava pokazila neskoré pripojenie agenta k idei a návrat agenta, ktorý bol vyradený pri štarte.
> Povoliť tretí, úzko ohraničený cyklus opráv?"
>
> Selected: **"Áno, úzky cyklus 3 (Recommended)"**. The option read: "Vrátiť pôvodné správanie pri vypnutom
> nastavení a opraviť drobnosť s archivovanou notifikáciou. Plán podpíše claude-1 a potom urobí review."
>
> Translation of the answer: "Yes, narrow cycle 3. Restore the original behavior with the setting off and fix
> the small issue with the archived notice. claude-1 signs the plan and then reviews."
>
> **Q2.** "Druhý MAJOR: rozpoznávač chyby zcode môže byť v praxi neúčinný. Reálny stderr zcode vnorené objekty
> skracuje na „[Object]“ a rozpoznávač taký záznam odmietne. Codex navrhuje zachytiť skutočný výstup zcode.
> Overil som však, že zcode má teraz kvótu: o 14:57 vrátil PONG. Takto by sa zachytila len úspešná odpoveď,
> nie chyba kvóty. Čo s tým?"
>
> Selected: **"Vydať s obmedzením + follow-up (Recommended)"**. The option read: "Vydať so známym obmedzením:
> rozpoznávanie zcode môže byť neúčinné a namiesto automatiky sa agent opýta teba, čo je bezpečné. Výstup
> najbližšieho skutočného vyčerpania zcode sa zachytí a gramatika sa podľa neho upraví v nadväzujúcej idei."
>
> Translation of the answer: "Release with the limitation plus a follow-up. Release with a known limitation:
> zcode recognition may be inert, so instead of the automatic path you are asked, which is safe. The output of
> the next real zcode exhaustion is captured, and the grammar is adjusted to it in a follow-up idea."
>
> **Relay fact (PRIMARY).** At 14:57 CEST the relay ran one zcode invocation with the prompt "Reply exactly
> PONG. Do not call tools, read or write files, or execute commands." It exited 0 and printed `PONG`, with
> empty stderr. zcode is not exhausted now, so a capture today could not produce the qualifying native
> error.
>
> ## What this authorizes
>
> - **Q1.** One narrow fix-up cycle 3 within the bounds in your note (R5-MAJOR-1 and R5-MINOR-1), with the
>   normal Phase-7 plan signed by claude-1 and a separate full re-review. A fresh CRITICAL or MAJOR on the
>   cycle-3 fix code escalates again, as you proposed.
> - **Q2.** An explicit owner waiver of AC2's native positive evidence for this release, with no capture and no
>   grammar relaxation now. Record it in `IMPLEMENTATION.md` as an owner-accepted known limitation, and state
>   it in the protocol or skill wording and in the release notes. The wording should say that zcode
>   auto-exclusion may not fire on real native output, and that unrecognized failures fall back to the
>   owner-confirmed path. Open a linked follow-up idea whose precondition is a captured native zcode
>   exhaustion. If cheap and in scope, the CLI should retain the raw scrubbed stderr of a failed zcode
>   invocation privately, so the next real exhaustion becomes that capture. Otherwise, name this in the
>   follow-up. R5-MAJOR-2 is dispositioned as owner-accepted and deferred, not fixed.
> - Everything else is unchanged: the claude-1 signoffs, a NEW attended-close request with current-tree
>   evidence, then the release. Both standing retry permissions still apply to failed invocations.

## Signoffs

<!-- Each participant appends only its own block. Agreement is to this plan, not code acceptance. -->

### Signoff: codex-1 — 2026-10-06
Status: ✅ ACCEPT
Notes: I accept G15–G17 as the narrow cycle-3 plan implementing Q1/Q2. R5-MAJOR-2 is owner-accepted and deferred; AC2 is not met by native evidence. No grammar relaxation or provider capture is authorized. Independent full-scope re-review, both final signoffs and the NEW attended close remain required.

### Signoff: claude-1 — 2026-10-06
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: I accept G15–G17 as the narrow cycle-3 plan. It honestly implements the owner's round05-answer: Q1's two bounded repairs (R5-MAJOR-1, R5-MINOR-1) and Q2's release waiver of native-positive AC2, with R5-MAJOR-2 recorded as owner-accepted and deferred, never fixed or passed. Reservations W1–W5 below are open items for the cycle-3 implementation; the cycle-3 re-review checks each one. This signoff approves no code, does not say AC2 passes, and authorizes no close, merge or release.

**Protocol context.** `context_mode=full`, `source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`, `fallback_reason=null`. Body `.parley-runtime/protocol-packets/full-phase7-deliberation-73613f95….md`: `shasum -a 256` reproduces the hash, and `cmp` against the deck `COOPERATION.md` exits 0 (1,501 lines, 125,862 B). Actual coverage: the entire body, in four bounded chunks (1–380, 381–760, 761–1140, 1141–1501). The three over-long lines (59, 452, 1002) were shown in full and their tails printed separately. No optimized or shadow packet was used.

**Basis (PRIMARY unless marked).**
- Read in full: this draft; my `review/round-05/claude-1.md`; all five `inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_*.md` owner notes; `../../quota-zcode-native-exhaustion-capture/00-prompt.md`.
- Read in part: FINAL §5, §9, §11 and AC1–AC21; Phase 7/8 and the stopping judgment in the packet. Plan-only step: I did not re-read the product diff.
- Nothing changed since round 05: `git diff --name-only 0ee1888 HEAD` lists only `parley-deck/` paths; the skill worktree is clean at `e2f3649`; before my append, this file's working copy equalled its HEAD blob (sha256 `5525bb84…`).
- The drafter's position-change citations are exact: `review/round-04/consensus.md:108` and recommendation 2 at `inbox/archived/codex-1-to-user_…_round05-trajectory-ac2.md:18`.
- I invoked no provider, participant or zcode, edited nothing else, and made no retry.

**User direction.** The governing answer is `inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_round05-answer.md`. Q1 was answered "Áno, úzky cyklus 3" ("Yes, narrow cycle 3") and Q2 "Vydať s obmedzením + follow-up" ("Release with the limitation plus a follow-up"); originals in Slovak. The draft's `## User direction` quotes it verbatim, and I checked that quote against the file. That owner ruling disposes of R5-MAJOR-2 for this release. I do not withdraw the finding, and my round-06 review will cite the ruling. The scope-reset, review-quota, quota-standing and timeout-standing answers remain binding; I re-read each. This process is timeout relaunch 2 of 2.

**Evaluation.**
- G15 matches R5-MAJOR-1 and my suggested fix, keeps policy-on gates and veto protection, and requires real CLI-dispatched differentials against `27e42b8` on both volumes.
- G16 matches R5-MINOR-1.
- G17 stays within the owner's options. The owner allowed "protocol or skill wording", so skill-only wording is acceptable and keeps the three COOPERATION copies unchanged (AC1/V7). The plan includes no invocation, capture, parser relaxation or fixture relabeling.
- Deferring stderr retention is the owner's own "Otherwise, name this in the follow-up" branch. I concur: retention exceeds Q1's bound "(R5-MAJOR-1 and R5-MINOR-1)".
- Stopping judgment: concur. This is cycle 3 of 5. A fresh CRITICAL or MAJOR on cycle-3 code escalates, and there is no cycle 4 by default.

**Reservations.**
- **W1 (G15), what "restored" means.** Pass means the same outcomes as `27e42b8`, not only better messages:
  - D1: the joiner's dispatch before the edit exits 0 and writes its artifact.
  - D2: after an edit-first join, other participants' signoff appends still exit 0, and `status` shows pending catch-up with the command, not an integrity gate. The joiner's later dispatch then completes the import.
  - D3: both signoffs exit 0.
  List any remaining difference as an owner-visible knob-off change; FINAL §11 permits only C1 and the bare 503. Scope the route to the launching id's own `round-01/<id>.md` when that id is not a known member. A retry over its own incomplete stub before import must work (AC13). The route must never let an excluded known id dispatch (C1/AC6).
- **W2 (G15), the §5 decline route.** §5 offers a joiner two options: catch up, or "decline (❌ NON-PARTICIPANT note in consensus)". The plan covers only catch-up. No product code handles that note: non-test `grep -rn NON-PARTICIPANT internal cmd` matches only `internal/protocol/defaults/COOPERATION.md:754`. As worded ("no known-signer status while the artifact is missing"), G15 would bar a policy-off joiner from filing the decline. My inference from D2 is that the baseline allowed it (UNVERIFIED, not executed). Add this case to the differential and restore the baseline outcome, or list it as a deviation.
- **W3 (G16), deletion is normal.** §4 lets a recipient archive a notice *or delete it*. Also test that deleting a notice after the applied receipt never re-publishes it. Archival or deletion by the owner must never count as the "ambiguous state" that gates the idea. State the chosen behavior for an interrupted publication; one benign re-publication is better than a permanent gate.
- **W4 (G17), honest and consistent disclosure.**
  (a) The option the owner selected says the next real exhaustion's output "sa zachytí" ("is captured"). Retention is deferred, so the skill and release notes must say that capture is not automatic in this release. They must give a manual step: keep that run's private, unscrubbed per-agent `stderr.log`, never commit or copy it, and point to the follow-up. No code change is needed. By reading only, `internal/runner/runner.go:406,1131-1142` writes the stderr stream to that private file; its completeness and cleanup are UNVERIFIED.
  (b) Make every living record match the waiver: `docs/quota-membership.md:185` ("AC2 remains an owner decision") and IMPLEMENTATION.md's AC2 rows (e.g. `:166`, `:550`). They must say native AC2 is NOT MET, owner-waived for this release (Q2), a known limitation with the follow-up linked. Never PASS. No zcode support row may read as native-verified.
- **W5 (G17), the follow-up must stay inactive.** Its `00-prompt.md` has `status: open` (not a Phase 0 value) together with `participants: [codex-1, kimi-1, zcode-1]`. A participant's session-start check (§9 steps 3 and 5) could treat it as an open idea and start round-01, contrary to its "backlog only" text and this plan's "do not invoke zcode". Use the protocol's inactive shape, `status: candidate` with no participants quorum (LE-10, §12.11, §14.1), and keep the intended roster in prose.

**Still mandatory.**
- My full-scope cycle-3 re-review. Where feasible, I will run the full suite, race, vet, the skill suite and a real-process crash exercise myself, and name any limits.
- Both final review-consensus signoffs.
- A NEW attended owner close.
codex-1's ✅ is the implementer's; procedural calls stay provisional under §15.5.
