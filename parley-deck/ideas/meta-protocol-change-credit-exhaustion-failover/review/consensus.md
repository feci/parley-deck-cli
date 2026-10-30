---
idea: meta-protocol-change-credit-exhaustion-failover
review-cycle: 3
drafted-by: codex-1
date: 2026-10-30
reviewed-commit: b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666
outstanding_agreed_fixes: 0
closing_review_round: 3
---

Codex-1 is the organizer, implementer, participant, and consensus drafter.
Kimi-1 and Zcode-1 wrote independent full-scope reviews in separate native processes.
Skill target: `406b6976cc75344d76f908297fe46ec2443e62cb`.
The complete CLI baseline is `c2d2431`; the skill baseline is `f557c5e`.
Later CLI commits contain audit records only. No product source changed during this review round.
The signed second-cycle consensus remains byte-exact at `consensus-cycle-02.md`.

## Agreed fixes

None remain proposed. Both reviewers withdrew the product findings after the two accepted fix-up cycles.
The two evidence-record NITs have concrete corrections below. Zcode must assess the checklist correction before its signoff.
This draft is provisional until all three participants append ACCEPT.
The fresh independent goal check remains required after these signatures.

## Resolved prior findings

| Finding | Correction and independent disposition |
| --- | --- |
| F1: pre-pin transient protection | Both round-03 reviews retain the resolved disposition. Synthetic readiness preserves two attempts and five-second spacing. |
| F2: root version | Both reviews retain the resolved disposition. Root and binary versions are 1.56.0. |
| F3: unreadable active record | Both reviews retain the resolved disposition. Positive stat evidence survives the read error. |
| F4: umask-dependent fixture | Both reviews retain the resolved disposition. Kimi reran both umasks. |
| F5: wait wording | Both reviews retain the resolved disposition. Tests cover absent and early handover records. |
| F6: empty changelog heading | Both reviews retain the resolved disposition. The release heading is first. |
| F7: partial invocation publication and durability | Both reviewers withdrew their CRITICAL findings. Each ran 100 concurrent-budget repetitions at fixed source. |
| F8: capture-status cleanup text | Both reviewers withdrew their documentation findings. Help and both references disclose deletion and cleanup receipts. |
| F9: driver import format | Both reviewers withdrew the formatting finding after their changed-file checks. |

The F7 checks cover complete exclusive publication, store barriers, request/terminal faults, collisions, permissions, cleanup, residue, and repeated failure.
The reviewers retain the limit: fault injection and source inspection do not prove every real power-loss sequence or filesystem.
The signed cycle-2 conflict resolutions remain in the archived consensus. No current decision depends on the rejected durability-parity premise.

Zcode's round-03 receipt NIT is withdrawn in its own supplemental assessment.
The full Go receipt came from the original process wrapper. The concurrent and portable receipts explicitly state their retrospective organizer provenance.
Those two receipts supply no driver accounting. Zcode checked their hashes and retained its independent focused refutations.

Zcode's supplemental checklist NIT requested a checked box for completed cycle-2 tests.
Codex checked that box and the initial completed-test box. No source or evidence result changed.
Proposed disposition: resolved by this record correction. Zcode must append its assessment to its own round-03 review before it signs.

## Current-tree criterion evidence

These entries summarize verifier-owned verdicts. They do not originate an implementer verification verdict.
The complete refutations, commands, output, provenance, and limits remain in both `round-03/` files.

| Criteria | Independent evidence at the reviewed products |
| --- | --- |
| AC-01 | Both reviewers ran quota, membership, and protocol checks for saved-policy and widening behavior. |
| AC-02 | Both ran synthetic native-framing and reset fixtures. Native-positive and exact-version limits remain. |
| AC-03–04 | Both ran runner suppression, recovery, stale-output, and handover checks. F7 adds checked publication evidence. |
| AC-05 | Both ran independent 100-count concurrent-budget tests. Kimi also ran 100 grouped-accounting repetitions and transient readiness. |
| AC-06–07 | Both ran whole-batch, original-order, facilitator, and fixed-floor fixtures. |
| AC-08–10 | Both checked active execution, shared drafts, retained obligations, and causal review proof. |
| AC-11 | Both retain source-backed Codex 0.165.0 with explicit JSON only. The ordinary smoke proves no exhaustion capture. |
| AC-12 | Both ran capture fixtures and inspected the status documentation. No capture-store authority exists. |
| AC-13 | Kimi reconstructed all four protocol/core targets from 29 signed hunks. Both checked hashes and packet limits. |
| AC-14 | Both inspected the completed full Go result. Kimi independently ran 402 Node and 54 Python tests. |
| AC-15 | Full-scope review exists. Consensus and the fresh goal check gate implementation close; delivery evidence follows D12. |

The organizer's required command passed at CLI product `b7dce40`: `go test ./... -count=1 -timeout 45m`.
Its wrapper receipt records exit 0 and 34 packages. Both reviewers checked the retained log against its hash.
Log SHA256: `68cd0dedd2609cf803e8e24d8ad741a09e48d2193d84f911f90e58864e6172d5`.
The organizer also ran relevant race checks, go vet, repeated concurrency, package checks, and the native portable build.
Canonical evidence: `../source-context/fixup02-validation.json`.
These organizer runs are attributed evidence, not a substitute for independent review.

## Deferred follow-ups

- `quota-zcode-native-exhaustion-capture` retains native-positive provenance and completeness work beyond this diagnostic facility.
- `quota-kickoff-reporting-and-alias-guidance` retains its remaining reporting work. Required account-alias guidance is delivered here.
- Peer account inference, machine exhaustion caches, other native adapters, and broader Codex versions remain deferred under FINAL.
- CLI Windows remains experimental. CLI WinGet stays held. No permission or durability check was weakened.
- npm and global-core publication remain owner-only. Their pending commands remain part of the final delivery note.

## Dismissed findings

None were dismissed by organizer judgment. Product findings were corrected and withdrawn by their authors.
The receipt NIT was corrected and withdrawn. The checklist disposition requires its author's own assessment and this signoff gate.
No participant dropout, silence rule, owner waiver, or severity downgrade closes a finding.

## Drafter position changes

No new design position changes. The changes since `../round-03/codex-1.md` remain recorded in `../consensus.md` under the same heading.
The cycle-2 change from visibility-only publication to checked directory durability remains recorded in `consensus-cycle-02.md`.
This draft now proposes zero remaining fixes, based on the complete current reviews and the stated evidence-record correction.
It replaces no frozen design or historical verdict.

## Alternatives disposition

- **ALT-F7a/F7c: retain the signed adoption.** Complete exclusive publication and checked directory barriers remain required.
- **ALT-F7b/F7d: retain the signed rejection.** Parse-error tolerance and the incomplete overlay cannot satisfy the accepted persistence contract.
- **ALT-F8a/F8b: retain the signed disposition.** Explicit cleanup documentation preserves D10's behavior.
- **ALT-CLOSE-1: adopt a fresh independent goal process after signatures.** Agreement alone supplies no current-tree goal proof.
- **ALT-CLOSE-2: reject a premature channel-pass claim.** D12 places publication after the pre-release gates and requires a separate channel process.

## Coverage & blind spots

Both reviewers independently inspected the full product diffs and tried to refute the acceptance criteria.
Each reproduced the former race before F7, then passed 100 repetitions after F7.
Kimi found the missing durability operation in cycle 2. Zcode found the receipt and checklist gaps in round 3.
Related models can share assumptions. Unanimity is a shared prior, not independent evidence by itself.
An untested persistence boundary, unsupported filesystem, incomplete native contract, or missed execution consumer could invalidate acceptance.
No provider exhaustion was induced. Native exhaustion fixtures remain synthetic or source-backed within FINAL's exact limits.
The full Go result is attributed to the organizer; neither review falsely claims an independent full-suite execution in this round.

## Remaining close and delivery gates

After all signatures, Zcode-1 will perform the fresh goal check in a new native process.
It must inspect the current products, all criteria, the review dispositions, and the completed validation evidence.
It must record its own unambiguous result in its review artifact. A reserved or incomplete result cannot close implementation.
If the pre-release gates pass, mark the execution record complete and merge both implementation PRs with merge commits.
Preserve the product trees and bind release assets to the clean merged commits.

D12 places publication and installation after the fresh goal check. AC-15 channel evidence therefore remains pending until delivery.
Publish both GitHub releases and Homebrew formulae. Submit the skill-only WinGet PR.
Install all 15 managed targets and four named generic destinations. Verify all six packaged SKILL.md files at each destination.
Preserve historical source/archive copies and runtime backups. State the two unavailable cloud-trash paths as an inventory limit if they persist.
A separate Kimi or Zcode process must inspect every authorized channel and runtime hash after publication.
The final inbox note must list npm then core commands for each pending version, 2.15.0 through 2.20.0.
The owner alone executes those publication commands. The relay owns the owner script.

The phase-7 driver refused: `continue failed: idea "meta-protocol-change-credit-exhaustion-failover" has no runs yet`.
The owner-authorized native fallback applies. It supplies no driver accounting or automatic membership authority.
The host reboot caused no participant exclusion. All three participants remain required signers.

## Protocol context

context_mode: full
source_sha256: b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709
packet_sha256: b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709
fallback_reason: none

## Signoffs

### Signoff: codex-1 — 2026-10-30
Status: ✅ ACCEPT
Notes: I accept the stated dispositions and zero-fix proposal, subject to both peer signatures and the fresh independent goal check.
This signoff supplies no independent verification or delivered-channel claim.

### Signoff: zcode-1 — 2026-10-30
Status: ✅ ACCEPT
Notes: I accept the zero-fix proposal and the stated dispositions. My round-03 review and both supplements hold zero open findings at CLI `b7dce40` and skill `406b6976`. I withdrew the receipt NIT in my first supplement. After I verified the two checkbox corrections, I withdrew the checklist NIT in my second supplement. Evidence attribution is accurate: organizer runs are attributed evidence; my AC-14 verdict rests on hash-verified organizer receipts, not an independent full-suite execution. The dispositions, retained limitations, delivery counts, and gate order match FINAL D12, AC-15, and the controlling brief.

This signoff is a position, not new verification. It passes no criterion alone and states no goal-check result. The fresh independent goal check in a new native process stays required before completion. Delivery, installation, and separate channel verification follow D12. npm and global-core publication stay owner-only; the final note must list versions 2.15.0 through 2.20.0.

### Signoff: kimi-1 — 2026-10-30
Status: ✅ ACCEPT
Notes: I accept the zero-fix proposal, evidence attribution, dispositions, limitations, and remaining gate order. My cycle-3 consensus assessment in `round-03/kimi-1.md` supplies the factual basis.
I concur with the checklist correction as a factual execution-record correction with no product edit. This signoff is a position, not new verification. It passes no criterion alone and states no goal-check result. The fresh independent goal check remains required before implementation completion. Delivery, installation, separate channel verification, and owner-only publication remain under D12.
