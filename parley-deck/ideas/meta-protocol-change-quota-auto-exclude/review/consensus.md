---
idea: meta-protocol-change-quota-auto-exclude
review-cycle: 4
outstanding_agreed_fixes: 3
blocked: false
drafted-by: codex-1
date: 2026-10-07
reviewed-commit: 25ea1da73ab0eab5e6099f759e2720aff54ada88
product-commit: fac40aaa1bba5350351bab0310ea8ede30ed3e6e
skill-commit: e976f7c9515f3250761380c1e09295e0f6c59985
---

## Scope and review basis

The newest owner finish-now direction authorizes cycle 4 for round-06's R6-MAJOR-1
and R6-MINOR-1/2, and if needed one narrow cycle 5 without asking. The G18–G20 plan
below is unchanged. A scope/FINAL-changing CRITICAL or findings remaining after cycle 5
requires escalation. Close is pre-confirmed after a final review with no open
CRITICAL/MAJOR, both signoffs and current AC1–AC21 evidence (AC2 owner-waived).

Role concentration (§15.5): codex-1 organizes, implements and drafts this plan;
claude-1 independently evaluates it and later reviews the code in a separate process.
This plan is not a code verdict. Prior reviewer files and archived signoffs are unchanged.

## Agreed fixes

- **G18 — owner-owned notices never gate membership** (R6-MAJOR-1). Remove notice
  inspection from `protocol.InspectQuota`. A valid applied receipt ends ordinary notice
  publication: no later notice content, archive, deletion, formatting or pathname condition
  can invalidate membership, signoff or driving. Before a receipt, preserve any regular
  live/archived notice as owner-owned; its bytes are not authority. If absent, publish with
  the existing checked exclusive write. Never overwrite annotations or follow a symlink or
  non-regular destination. Publication-only failures must be reported as non-blocking
  diagnostics and must not prevent the checked membership/evaluation receipt; distinguish
  a completed publication attempt from proven message delivery in comments/docs. Preserve
  genuine immutable-history, receipt-path, writer and projection integrity gates. Keep the
  historical manual-authority clarification accurate and non-authoritative, with the same
  owner-ownership rule. Prefer deleting validation coupling over new queues/schemas/retries.
  Tests: live/archived appended answers and arbitrary edits, before/after receipts;
  archive/delete/repeated recovery; safe refusal of unsafe publication paths without
  blocking read/signoff/dispatch; true receipt/history corruption still gates; one terminal
  evaluation. Run on shared storage and /tmp, including actual status/signoff probes.
- **G19 — round-1 return from immutable kickoff evidence** (R6-MINOR-1). Preserve the
  initially confirmed excluded identities at new kickoff in a small optional immutable
  field, disjoint from the kickoff quorum. Existing automatic kickoff transition evidence
  is already immutable. The policy-off round-1 participant edit can then restore that id
  even when its stale display marker is removed. Apply the same predicate during manual
  import and immutable-history validation; the id is known only from the new revision.
  Do not manufacture historical membership or owner authority. Old kickoff records did
  not retain manual exclusions: leave those bytes untouched; if no authentic immutable
  evidence exists and the marker is removed, require the existing catch-up path and
  disclose that legacy limit. Preserve already recorded manual revisions and the existing
  marker-kept compatibility path. Test kept/removed markers, genuine new identities,
  missing legacy evidence, later rounds, policy-on and known-excluded boundaries and
  retained vetoes. No migration or general membership redesign.
- **G20 — truthful policy-off release wording** (R6-MINOR-2). Replace preservation claims
  with the actual differences against 27e42b8 in CHANGELOG, the release draft, CLI docs
  and skill guidance: C1 and bare 503; new joiner's late-round-1 import before signing or
  satisfying quorum; exact NON-PARTICIPANT decline remains in Missing; post-round-1
  kickoff return needs catch-up; pending catch-up stops ordinary driving until manual
  exec; final/closed membership edits are frozen; any G19 legacy evidence limitation.
  Include G18's non-blocking publication-failure behavior. Preserve the explicit native
  AC2 NOT MET / owner-waived wording and inactive capture follow-up. No recognizer change,
  provider invocation, capture, roster/model/settings change or normative core rewrite.

## Deferred follow-ups

- R5-MAJOR-2 / native AC2: owner-accepted and deferred, never fixed or PASS, to
  `../quota-zcode-native-exhaustion-capture/00-prompt.md` (inactive candidate). No native
  capture or grammar relaxation now. Manual private stderr preservation remains as stated.
- D6 legacy driver accounting, unsupported adapter provenance and Windows runtime remain
  separate. Core publication remains owner-only. Release follows the newest pre-confirmed close conditions.

## Dismissed findings

None. G18–G20 are proposed repairs/disclosures; only claude-1's re-review may verify them.

## Coverage & blind spots

The raw round-06 review is the current independent verdict: 0 CRITICAL, 1 MAJOR, 2 MINOR.
It includes full product coverage and broad independent macOS checks. Tests alone did
not catch the owner-annotation defect. Re-review the full product diff since FINAL and
all fixes, weighing every issue freely. Repeat meaningful new regression probes and
required current-tree full Go/build/vet/gofmt/skill/packet checks. Windows runtime and
container/PID namespaces remain unverified. This single independent reviewer still
requires both final signoffs and the newest pre-confirmed close conditions.

## Drafter position changes

Relative to cycle-3 G16 (`review/round-06/consensus.md`), the condition "ambiguous state
stays gated" no longer includes any notice contents or publication-only failure. Owner
round06-answer expressly requires the notice never to block; immutable membership
records retain their gates. G19 records that the old manual kickoff schema actually
contains no excluded identities (`internal/quota/record.go:Kickoff`), so the prior mutable
marker assumption cannot establish an exclusion after deletion. A small prospective
field plus an explicit legacy limit replaces that assumption. No other material design
position changes from `round-03/codex-1.md` are proposed beyond the recorded owner rulings.

## User direction

> ## Owner answer to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_round06-trajectory.md`
>
> Relayed by the owner's Claude Code session. The answer arrived 2026-10-07 at about 09:32 CEST. The question
> and the answer are quoted verbatim in Slovak, with a translation of each below:
>
> > Question: "Po treťom cykle opráv review našlo 0 CRITICAL, 1 MAJOR a 2 MINOR. Počet nálezov klesá 15 → 6 →
> > 3 → 3. MAJOR: keď do notifikácie o vyradení dopíšeš odpoveď, systém ju vyhodnotí ako poškodenú a zablokuje
> > stav, podpisy aj driver. Zmena už zasahuje 117 súborov CLI, a to je ďaleko od „malej zmeny“. Ako ďalej?"
> > Selected: **"Cyklus 4, posledný (Recommended)"**. The option read: "Úzko ohraničená oprava: notifikácia o
> > vyradení už nebude nikdy blokovať, návrat agenta v 1. kole bude fungovať a release notes budú pravdivé.
> > Ak sa objaví ďalší MAJOR, codex už nebude opravovať a ozve sa ti s možnosťou vydať iba etapu 1 alebo
> > zmenu zastaviť."
>
> Translation of the question: "After the third fix-up cycle, the review found 0 CRITICAL, 1 MAJOR and 2 MINOR.
> The finding count is falling, 15 → 6 → 3 → 3. The MAJOR: when you append an answer to an exclusion notice,
> the system judges it corrupt and blocks status, signoffs and the driver. The change already touches 117 CLI
> files, far from a 'small change'. How to continue?"
>
> Translation of the answer: "Cycle 4, the last one. A narrowly bounded fix: the exclusion notice never blocks,
> returning an agent in round 1 works, and the release notes are truthful. If another MAJOR appears, codex
> stops fixing and reports back with the option to release only stage 1 or to stop the change."
>
> ## What this authorizes
>
> - **One narrow fix-up cycle 4**, exactly within the bounds of your note: notice ownership, round-1 return
>   and honest release wording. The usual sequence applies: Phase-7 plan, claude-1 signoff, implementation,
>   then claude-1's full re-review. Keep the diff minimal. Prefer removing gating over adding machinery.
> - **Cycle 4 is the last fix-up cycle.** If the re-review finds a new CRITICAL or MAJOR, do not open cycle 5.
>   Stop and write a blocking note with two options: release stage 1 only, with a split assessment, or stop and
>   park the change. If the re-review is clean of CRITICAL and MAJOR, continue to the review signoffs and the
>   NEW attended-close request.
> - Both standing retry permissions still apply to failed invocations.

## Newest user direction — finish now

> ## Owner direction: finish now (supersedes the wait in `…_long-quota-answer.md` and `…_provider-stop-answer.md`)
>
> Relayed by the owner's Claude Code session on 2026-10-07 at about 23:20 CEST. Verbatim (Slovak):
>
> > "sakra tak to fixni a dokonci a deployni cez vsetky kanaly, taha sa to dlho"
>
> Translation: "Damn, then fix it, finish it and deploy it through all channels, this is dragging on."
>
> ## Relay facts (PRIMARY, 23:18 CEST)
>
> - `ANTHROPIC_BASE_URL` points every Claude CLI call at the OmniRoute gateway, so there is no direct route.
> - Two probes ran, each with a 103,660-byte prompt (FINAL + consensus + review consensus): one with
>   `--model 'claude-opus-5-5[1m]'` (plain id) and one with `--model 'claude/claude-opus-5-5[1m]'`
>   (prefixed id). **Both returned `PONG`.** Large requests pass right now. The 22:01 failure came after
>   857 s of an agentic session, so the gateway pool is intermittent, not hard down.
> - The relay's own long sessions use the plain id and kept working through the pool errors that hit the
>   prefixed id.
>
> ## What the owner authorizes now
>
> 1. **Start immediately.** Do not wait for 2026-10-09. The relay killed the auto-resume daemon.
> 2. **claude-1 stays the reviewer and the model stays Opus 5.5.** For the rest of this idea, launch claude-1
>    with the plain model id `claude-opus-5-5[1m]` instead of `claude/claude-opus-5-5[1m]`. This is the
>    same model at the same max effort; only the gateway route id changes. Record it in
>    `organizer-notes.md`. Do not edit `agents.toml`.
> 3. **Provider errors no longer stop you.**
>    - If a claude-1 or codex-1 step fails with a 429 or 503 quota/unavailable error, whatever reset it
>      states, wait 15 minutes and relaunch the same step. Do this at most 8 times per step.
>    - Silent timeouts keep the earlier rule: relaunch with 2400 s, then 3600 s.
>    - Stop only on an auth or credit error, or when the attempts for a step run out.
> 4. **Fix-up.** Do cycle 4 now: the plan signoff, the implementation and the full re-review. If the cycle-4
>    re-review finds new findings, fix them narrowly in cycle 5 without asking. That is the last cycle the
>    protocol's deliberation cap allows. Escalate only if a CRITICAL needs a change of scope or FINAL, or
>    if findings remain after cycle 5.
> 5. **Close is pre-confirmed.** The owner's "dokonci a deployni" is the attended-close confirmation,
>    provided all of the following hold:
>    - the final re-review has no open CRITICAL or MAJOR;
>    - both review-consensus signoffs exist;
>    - current-tree evidence for AC1 to AC21 is recorded, with AC2 owner-waived as already decided.
>
>    Record this verbatim as the close authority. Do not write a separate close-request note.
> 6. **Release immediately after the close,** on all channels, as `IMPL-ORGANIZER-BRIEF.md` says:
>    - merge to `main`;
>    - CLI 1.51.0, skill 2.15.0 and core 2.15.0 staged;
>    - GitHub releases;
>    - both Homebrew formulae;
>    - a winget PR for the skill only;
>    - install the skill into all local runtimes and verify each by hash;
>    - claude-1 verifies every channel independently, using a short brief.
>
>    The two steps only the owner can run (`! npm publish …` and `! parley protocol publish …`) go in ONE
>    final note, `codex-1-to-user_meta-protocol-change-quota-auto-exclude_released.md`, with the exact
>    commands.
> 7. **Unchanged:** the owner's earlier answers (the zcode stderr rule, the AC2 waiver and the follow-up),
>    no other reviewer, no quorum change, and English artifacts.

## Protocol context

Resumption uses the unchanged live phase-7 full packet, source and packet SHA256
`73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`; fallback_reason absent.
Lean reorientation uses IMPLEMENTATION, raw round-06, prior signed plan, organizer tail,
computed brief/status and relevant Phase 7/8/§15.5 provisions. The separate signer reads
its rendered full protocol body before attesting. The newest finish-now retry rules bind.

## Signoffs

<!-- Each participant appends only its own block. Agreement is to the plan, not code. -->

### Signoff: codex-1 — 2026-10-07
Status: ✅ ACCEPT
Notes: I accept the unchanged G18–G20 repairs under the newest finish-now direction.
This is implementer plan acceptance, not independent verification. Cycle 5 is authorized
if needed; final independent review, both signoffs and the recorded close conditions
precede release. The newest note supersedes historical schedule and last-cycle rules.
