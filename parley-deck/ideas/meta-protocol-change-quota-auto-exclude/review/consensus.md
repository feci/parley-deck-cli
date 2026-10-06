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
