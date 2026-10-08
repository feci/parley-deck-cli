---
idea: meta-protocol-change-quota-auto-exclude
review-cycle: 5
outstanding_agreed_fixes: 0
blocked: false
drafted-by: codex-1
date: 2026-10-08
reviewed-commit: 2705a1e850132f74aed2dbb87149df5491bfe94b
skill-commit: 99b3f3f9fee161e8e61e585ad6c8e5b1bbbd4d2b
closing_review_round: round-08
strict_gate_clean: false
---

## Scope and review basis

This is the final review-consensus draft after the full independent round-08 on
CLI 2705a1e / skill 99b3f3f, all 124 CLI product files and five skill files since
FINAL. The review reports 0 CRITICAL, 0 MAJOR, 1 MINOR and 2 NIT and verifies every
G21–G25 fix and signed reservation. Its raw file is unchanged, SHA256
6ea871262a87731b2c5273b28128b851c7a0ad797abaf26601e4c845d404ce06. Validator: 1/1 valid.
The signed cycle-5 plan is preserved byte-for-byte at review/round-08/consensus.md
and source-context/cycle5-plan-20261008/signed-plan.md. Those were plan signoffs;
the blocks below are the separate final consensus signoffs.

The newer binding round08-answer accepts/defer all three residual findings with
exact disclosures and a linked follow-up, including AC5/AC15 exceptions. It says
explicitly that the owner's Claude relay made the decision under standing direction
without asking the owner again. The current user instruction makes every matching
inbox note binding, newest wins. That is the authority applied here; no fresh direct
owner answer is invented. No sixth code cycle or clean-review claim is made.

Role concentration (§15.5): codex-1 organizes, implements and drafts, and supplies
no independent code verdict. claude-1 alone independently reviewed and owns its
final block. The idea has no strict_gate: true; strict_gate_clean is deliberately
false because accepted/deferred findings still exist. Owner pre-confirmed close
conditions remain separate from this draft and require both final blocks.

## Agreed fixes

None. Zero agreed fixes remain in this idea. The three round-08 findings below
are explicitly accepted/deferred under the binding relay direction, not fixed,
withdrawn or silently waived. The authorized disclosure text has been applied
exactly to both changelogs, CLI docs and the release draft. No Go, schema,
recognizer, protocol text, model, provider, roster or runtime behavior changed
after the frozen round-08 product commits.

## Deferred follow-ups

- **R8-MINOR-1, kickoff escalation:** accepted/deferred by round08-answer. AC5 is
  NOT MET in full: a missing/read-only inbox loses the escalation and decision
  details, although the gate fails closed and no idea is created. Follow-up:
  [quota-kickoff-reporting-and-alias-guidance](../../quota-kickoff-reporting-and-alias-guidance/00-prompt.md).
- **R8-NIT-1, kickoff notice crash window:** accepted/deferred by the same note.
  AC15 is NOT MET in full: a crash after manifest and before publication can lose
  the notice permanently; other exclusion surfaces remain. Source evidence only,
  no injected crash. The same linked follow-up owns the recovery/semantics work.
- **R8-NIT-2, aliased-deck guidance:** accepted/deferred, with the exact disclosure
  that plain participant/confirmed-exclusion edits require a physical scope even
  when policy is off and can block all signers/driving. The generic disable-policy
  remedy is inapplicable there. The same linked follow-up owns contextual guidance.
- **R5-MAJOR-2 / native AC2:** NOT MET, expressly owner-waived by round05-answer Q2,
  accepted/deferred to [quota-zcode-native-exhaustion-capture](../../quota-zcode-native-exhaustion-capture/00-prompt.md).
  Real zcode auto-exclusion may be inert. No capture or grammar relaxation ships.
- **Windows:** known directory-sync failures in creation (including policy-off),
  scoped driving/signing, manual imports, revisions and transitions. Windows assets
  remain experimental; CLI winget is held. Concrete site/CI handoff is
  ../../inbox/codex-1-to-all_windows-portability_quota-durable-sites.md; the
  windows-portability track is on its unmerged branch, not present in this tree,
  and promises no fix for these future durable sites.
- **Other existing limits:** D6 stale driver accounting and unsupported adapter
  provenance remain separate follow-ups (TBD if not opened); container/PID namespace
  behavior is unverified. No parent participant is excluded or replaced here.

## Dismissed findings

None. All R7 findings and the additional signed G25 reservation are independently
verified resolved in round-08. Earlier fix dispositions remain in the archived
consensuses and reviews. The AC2 waiver and the three R8 deferrals remain explicit.

## Coverage & blind spots

The reviewer independently read every full CLI/skill diff chunk and the complete
live protocol. Its full Go suite passes (616.8s, 34 packages), as do build/vet,
race (411 PASS), shared/local focused tests (244 PASS each), drift/packet checks,
all 100 changed Go files' formatting and skill tests (399 Node, 54 Python, six
manifests). Real CLI overlay probes cover kickoff inboxes/floor/aliases and prior
notice/plain-edit cases. Producer host checks separately pass, as do Linux/macOS
and skill CI; Windows hosted tests fail and are never relabeled as passing.

Retain the raw review's limits: earlier interrupted attempts had shared-volume
failures that later passed; filesystem coherence is a reasoned cause, not proof.
The baseline binary is reused from round-07, Windows runtime is not executed by
this reviewer, GitHub CI is producer evidence, and R8-NIT-1 is source-only.
IMPLEMENTATION.md's current Validation evidence table records every AC1–AC21,
including AC2/AC5/AC15 exceptions and AC21 pending the final blocks and close.

## Drafter position changes

The initial cycle-5-limit escalation was appropriate under finish-now point 4.
The newer binding relay note now disposes of the three findings as accepted/deferred,
so no additional code repair is proposed. Exact disclosures plus the inactive linked
follow-up make the release limitations reviewable. The answered escalation is
archived unchanged in source-context/round08-cycle5-owner-gate-20261008/answered-escalation.md.
This draft does not rewrite the independent reviewer, infer a clean review or
replace the required final signature with the earlier signed plan.

## User direction

### Finish-now, verbatim

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

### Newest round08-answer, verbatim

> ## Answer to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_round08-cycle5-limit.md`
>
> Written by the owner's Claude Code session (the relay) at about 02:30 CEST on 2026-10-08. The relay did
> **not** put this question to the owner. It decides under the owner's standing direction of 2026-10-07,
> verbatim: "sakra tak to fixni a dokonci a deployni cez vsetky kanaly, taha sa to dlho" ("then fix it,
> finish it and deploy it through all channels, this is dragging on"). The decision is reported to the owner
> in the same turn, and the owner may override it.
>
> **Decision: accept and defer** the three remaining findings, R8-MINOR-1, R8-NIT-1 and R8-NIT-2, as explicit
> release limitations. Use exactly the disclosures in
> `source-context/codex-1-round08-disclosure-proposal.md` and the linked follow-up. The basis:
>
> - the independent round-08 review has 0 CRITICAL and 0 MAJOR;
> - every cycle-5 fix is verified;
> - each residual fails closed or is a pre-existing narrow window.
>
> This includes the AC5 escalation-detail exception and the AC15 kickoff crash-window exception. Neither is
> reported as a PASS or as fixed. The alias diagnostic caveat stays visible, and Windows CI failing stays
> disclosed (Windows remains experimental, and CLI winget stays held). There is no sixth code cycle.
>
> Then:
>
> 1. Get both final review-consensus signoffs. claude-1 owns its own signoff.
> 2. The close is pre-confirmed per `…_finish-now.md` point 5.
> 3. Release on all channels per point 6.
> 4. Finish with the single `codex-1-to-user_meta-protocol-change-quota-auto-exclude_released.md` note,
>    which carries the two owner-only commands.
>
> The retry rules in `…_finish-now.md` point 3 still apply.

## Signoffs

### Signoff: codex-1 — 2026-10-08
Status: ✅ ACCEPT
Notes: I accept the final dispositions as organizer/implementer, not as an independent
code verdict. Zero agreed fixes remain under the binding round08-answer; AC2 is
NOT MET/owner-waived, AC5 and AC15 are NOT MET in full/accepted-deferred, and all
three R8 limitations are disclosed exactly with a linked inactive follow-up. No
sixth code cycle or strict-clean claim. The independent final block remains required.
Phase-7 attestation: context_mode=full,
source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e,
fallback_reason absent. Current full body bytes match the previously read live
protocol. Both product commits remain the round-08 baseline; only authorized
disclosure and canonical workflow records have changed since that review.

### Signoff: claude-1 — 2026-10-08
Status: ✅ ACCEPT
Notes: This is my independent final signoff. It rests on my full round-08 review and on a fresh check of
everything changed since then. It grants no close, release or channel verification.

- **Dispositions: I concur with each one, judged on its own.** The relay's three bases hold against my raw
  review. First, there is no CRITICAL or MAJOR finding. Second, G21–G25 and my signed reservations are
  verified. Third, every residual fails closed or is a pre-existing narrow window. More precisely:
  - R8-MINOR-1 is both. It is stage-1 code, and when it fires nothing is applied and no idea is created.
  - R8-NIT-1 is the pre-existing window. It is source-only and not crash-injected.
  - R8-NIT-2 is not pre-existing: in 1.50.0 the signoff path had no alias-dependent code. It does fail
    closed, refusing without writes.
  - Disclosure plus a follow-up is option (c) of my R8-MINOR-1 fix list, and it is the "disclose" option
    of each NIT.
- **Acceptance-criteria labels.** I accept the stricter labels: AC5 and AC15 are NOT MET in full and
  accepted-deferred, never PASS or fixed. My raw "PASS, except" meant the same thing. AC2 stays NOT MET
  and owner-waived (round05-answer Q2).
- **What stays explicit.** The Windows CI failure, the broken Windows durable operations, the experimental
  Windows assets and the CLI winget hold. There are zero agreed fixes, no cycle 6 and no clean-review
  claim. `strict_gate_clean: false` is correct, because 00-prompt sets no strict_gate.
- **Authority, recorded as it stands.** round08-answer is the relay's decision under the standing
  finish-now direction. The relay made it without asking the owner, and the note says so.
  - I add no owner testimony, and I cannot verify the relay's same-turn report to the owner.
  - The three findings close through this normal signoff process (the Phase 6 dispositions rule), not by
    relay fiat.
  - Accepting AC5 and AC15 as unmet at close is a risk-tolerance call, and that call rests on the relay
    decision. If the owner overrides it, this acceptance does not carry over, and the findings stand as
    filed in review/round-08/claude-1.md (SHA256 still 6ea87126…).
- **Draft accuracy, ratified under §15.5.**
  - Agreed fixes, Deferred follow-ups, Dismissed findings and Coverage match my raw review.
  - Both User direction quotes are verbatim. After stripping the quote markers, cmp exits 0 against
    finish-now lines 10–64 and round08-answer lines 10–38.
  - The role-concentration line is present.
  - Drafter position changes is accurate and complete. codex-1's escalation already recommended
    accept/defer (answered-escalation.md lines 12–16 and 71–72), so there is no material change.
- **Delta since round-08, from my own runs.**
  - CLI: `git diff 2705a1e d3a846a`, outside parley-deck/{ideas,inbox}, touches only CHANGELOG.md
    (+25 lines) and docs/quota-membership.md (+24 lines). Both are pure additions.
  - Skill: 99b3f3f..bd01ed0 touches only CHANGELOG.md (+25 lines).
  - Proposal paragraphs 1–3 (lines 8–25, sha256 8f2536fd…) are byte-identical (cmp exit 0) at four places:
    CLI CHANGELOG.md:7, docs/quota-membership.md:266, source-context/codex-1-release-notes-draft.md:75
    and skill CHANGELOG.md:7. The release draft also states that AC5 and AC15 remain unmet.
  - The wording matches the code. Floor and role blocks share one plain O_EXCL writer
    (internal/app/preflight.go:419–423 and internal/app/quota.go:27). The kickoff notice is a single
    attempt after run.created and the manifest (runcontrol.Create).
  - The candidate quota-kickoff-reporting-and-alias-guidance is backlog only: no run has started and
    auto_implement is false. It covers all three findings, including crash injection before any fix claim
    and the preflight-ordering point, and its relative links resolve.
- **Why the evidence transfers to CLI d3a846a / skill bd01ed0.**
  - No *.go, go.mod, go.sum or testdata file changed. The only go:embed is defaults/COOPERATION.md, which
    is unchanged. No tracked Go file or skill file refers to either changelog or to the changed doc, and
    the skill package's `files` list omits CHANGELOG.md.
  - These are unchanged: FINAL (f90577f1…), deck and skill COOPERATION (73613f95…), both agents.toml files
    (f8cc2ab5… and f52a0a77…), and the archived signed plan. review/round-08/consensus.md is still
    byte-identical to signed-plan.md (e5b56a7f…).
  - Fresh runs on the current tree:
    - `go build ./...` and `go vet ./...` exit 0.
    - TestEmbeddedDefaultMatchesLiveDeck passes.
    - `parley protocol packet check` is ok (69 blocks).
    - The phase 0, 5 and 8 packets are full, with source = packet = 73613f95….
    - Skill `npm test` at bd01ed0 exits 0: 399 Node tests pass and 0 fail, 54 Python tests pass, and all
      six add-on manifests are ok.
  - The guard stubs logged no provider call. Git status and the roster hashes are identical before and
    after.
  - I did not repeat the round-08 full suite (616.8 s, 34 packages), the race run or the probes. Their
    results transfer because nothing they read has changed.
- **Close conditions (finish-now point 5), with this block in place.** There is no open CRITICAL or MAJOR
  finding, and both final blocks exist. IMPLEMENTATION's Validation evidence records AC1–AC21: AC2 is
  waived, AC5 and AC15 are accepted-deferred, and AC21 is pending the close.

Phase-7 attestation: context_mode=full,
source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e,
fallback_reason absent.
- I rendered it at 2026-10-08T00:02:53Z with `parley protocol packet --dir . --phase 7 --track
  deliberation --idea meta-protocol-change-quota-auto-exclude --flag protocol_change --json`.
- I ran it with three binaries: a current-tree build, my round-08 build of 2705a1e and the installed
  1.50.0. Each exited 0 with empty stderr, and the three JSON outputs are byte-identical (26,812 B).
- The shadow packet 01ba034e… (41 included and 28 omitted blocks) was not used.
- The emitted body, full-phase7-deliberation-73613f95….md, has 1,501 lines and 125,862 B. Its hash
  reproduces, and cmp against the deck and skill COOPERATION.md files exits 0.
- I read the entire body, lines 1–1501, and confirmed the ends of the long lines 59, 452 and 1002.

Limits:
- I did not re-run the full Go suite, the race run or the probes for this docs-only delta.
- I did not execute anything on Windows.
- GitHub CI is producer evidence, and I did not access GitHub.
- The round-07 baseline binary and the source-only basis of R8-NIT-1 carry over as limits.
- OpenViking was not consulted.

Non-blocking note: the changelog follow-up links sit outside the exact paragraphs and point at
blob/quota-auto-exclude/. The local origin/quota-auto-exclude ref (1961675) does not contain that file
until d3a846a is pushed, and the links would 404 if the branch is deleted after merge. A blob/main/ URL
would survive the merge. This is optional, not a condition.
