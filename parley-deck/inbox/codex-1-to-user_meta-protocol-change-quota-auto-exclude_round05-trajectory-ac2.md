---
from: codex-1
to: user
idea: meta-protocol-change-quota-auto-exclude
phase: review-round-05
blocking: yes
date: 2026-10-06
---

## Question

Please decide two substantive items after claude-1's complete cycle-2 re-review:

1. **Recommended: authorize one narrow fix-up cycle 3** for the policy-off catch-up/return regression
   (R5-MAJOR-1) and the reappearing archived notice (R5-MINOR-1), followed by the normal signed fix plan
   and separate full review. The alternative is an explicit acceptance of the membership regression as
   a new policy-off behavior change; I do not recommend that.
2. **Recommended: authorize one evidence-only native zcode capture and narrowly bounded grammar fixes
   justified by that capture**, under the exact capture and parser boundaries below. This combines
   reviewer options (a) and (c), keeps the owner's existing stderr rule, and does not waive AC2. You may
   instead authorize capture only and retain a separate grammar decision, explicitly accept AC2 as unmet
   with the inert-support risk, or defer the feature/release.

This is not a close request. A later NEW attended close still requires both final review-consensus
signoffs and current-tree criterion evidence. No cycle 3, capture, parser relaxation, merge or release
has been started on the strength of this proposal.

## Context

The newly authorized 2400-second Phase-7 relaunch succeeded: claude-1 signed the cycle-2 plan with
V1–V9. Cycle 2 was then implemented and tested. Its independent round-05 process finished normally
on the first attempt (1055.4 seconds, exit 0, timeout false); neither standing retry rule was needed
for this review. There is no quota, credit or authentication failure at this boundary.

Current product commits are CLI `0ee18889977a29d6baf53d3016b501112760ae12` and skill
`e2f3649eb938e870367c76943b441382fe7acf65`. Reviewed CLI snapshot:
`e7bf96c2f8914fea151ed1b862ba658b6f3852b9` (same product bytes).
Canonical independent review: `parley-deck/ideas/meta-protocol-change-quota-auto-exclude/review/round-05/claude-1.md`,
SHA-256 `73305207c6c57f9cb8bd76af2b22d3b14e86ea43185f8994de1d1cbd5e6d4e4f`.
It is filed-and-valid; the organizer read all 313 lines after process exit and did not edit it.

| Review | CRITICAL | MAJOR | MINOR | NIT | Total |
| --- | ---: | ---: | ---: | ---: | ---: |
| Round 03, before cycle 1 | 2 | 4 | 5 | 4 | 15 |
| Round 04, after cycle 1 | 0 | 3 | 2 | 1 | 6 |
| Round 05, after cycle 2 | 0 | 2 | 1 | 0 | 3 |

The count is falling, and both CRITICALs remain resolved. However, R5-MAJOR-1 is a fresh MAJOR on G11's
cycle-2 membership path. The signed cycle-2 consensus expressly says:

> If the cycle-2 re-review finds a new CRITICAL or MAJOR on cycle-2 fix code, above all in G10's framing or
> G11's membership path, treat it as churn. Stop and escalate with the trajectory; do not open cycle 3 by
> default.

Source: `review/round-04/consensus.md`, claude-1's own signoff. This applies now. Phase 8's five-cycle cap
is an escalation ceiling, not permission to ignore that condition. The standing retry permissions
cover failed participant invocations, not this product/scope decision.

### The exact regression and proposed narrow cycle 3

claude-1 ran the same real `app.Run` probes against actual pre-change `27e42b8` and the current code,
on both the shared volume and `/tmp`. All these policy-off paths succeed before the change and fail now:

- Dispatching joiner `e` to write its own late `round-01/e.md` fails with
  `excluded participant e cannot dispatch`.
- Adding `e` to `participants:` first blocks existing signers with a missing `round-01/e.md` error.
- Returning a kickoff-excluded `d` through a plain `participants:` edit during round 1 blocks signers
  with a missing `round-01/d.md` error.

The producer's differential pre-created the late round-1 bytes directly, so it missed the CLI dispatch
deadlock. Its broader compatibility claim is corrected in IMPLEMENTATION.md. The proposed repair is:

- Restore the policy-off CLI route for an agent to author only its own late round-1 artifact.
- Represent an otherwise valid join awaiting that artifact as pending catch-up with an actionable
  diagnostic, without falsely granting known-signer status, owner authority or a completed quorum.
- Restore the pre-change plain edit for a kickoff-excluded agent's return; it becomes known from that
  manual revision, never retroactively at kickoff. Manual history cannot authorize retained-veto withdrawal.
- Keep the policy-on, identity, protected-role, retained-finding and closure gates intact. Add regressions
  through CLI dispatch on both volumes, instead of pre-writing the artifact under test.
- Fix R5-MINOR-1: after the owner archives an applied transition notice, the next `Before` currently
  recreates it in the inbox. Preserve one durable publication and avoid resurfacing an archived notice;
  replay must still emit just one terminal evaluation.

These are proposed fix-plan bounds, not signed Agreed fixes. If authorized, codex-1 drafts Phase 7,
claude-1 appends its own signoff, codex-1 implements, and claude-1 independently re-reviews the complete
diff. A further fresh CRITICAL/MAJOR on that fix code again escalates. No D6/legacy-driver work is added.

### AC2: what the current recognizer actually accepts and rejects

The original 24,833-byte native stderr is unavailable. The offline audit locates the default SDK
`console.error` sink, default Node inspect depth 2, and the adapter's own retries with SDK `maxRetries:0`.
The reviewer reproduced the source harness byte-for-byte without invoking zcode or loading its config.

| Input shape | Current behavior | Evidence |
| --- | --- | --- |
| Complete allowlisted top-level APICallError records, including decreasing countdowns, consistent exhaustion/reset, terminal receipt at least 60 minutes before reset | Eligible | Source-derived fixtures and independent reconstructed probes |
| Complete allowlisted RetryError aggregate with exact lastError identity and the same gates | Eligible | Source-derived; not established as this adapter's native retry form |
| Default console output with realistic nested request messages rendered as `[Object]` | Rejected | Independent offline reproduction; not a complete native capture |
| Retained native tails beginning mid-object or only the decisive four lines | Rejected | Partial native evidence |
| Bundle-path, processTicksAndRejections and `at async k7r (...)` stack frames | Accepted within otherwise valid fixtures | Independent stack substitution probes |
| `at async Promise.all (index 0)`, `at <anonymous>`, `at Generator.next (<anonymous>)` | Rejected | Independent stack substitution probes; occurrence in native zcode remains unverified |
| Mixed/contradictory errors, missing reset, receipt below 60 minutes, incomplete/unknown framing, artifact or later success | Rejected | Existing and independent negative checks |

R5-MAJOR-2 remains open: the sole supported adapter may be effectively inert on realistic native
output. This is a functional/evidence gap, not evidence of a false automatic exclusion. The owner's
earlier stderr exception does not authorize us to call these reconstructed fixtures native or waive AC2.

### Proposed capture and grammar authority for question 2

Authorize **one** configured zcode evidence invocation (`zai/glm-5.3`, configured max effort), with a
2400-second ceiling and the minimal prompt: `Reply exactly PONG. Do not call tools, read or write files,
or execute commands.` It is an evidence probe, not a Parley participant/reviewer, quorum change or
substitute for claude-1. No other model/provider or roster/config change is proposed.

Capture complete stdout/stderr plus invocation identity, start, terminal receipt, exit status and output
facts in private ignored runtime files; retain hashes and publish only a checked scrubbed fixture and
provenance. Do not read, print, copy or alter credentials. No retry loop, forced quota consumption or
repeated requests to manufacture an exhaustion. If it succeeds, times out, or reports a different failure,
record that result honestly and leave native positive AC2 open for a further owner decision.

If the capture contains the qualifying native exhaustion, authorize a Phase-7 plan for only the
observed benign framing differences, such as inspect placeholders confined to `requestBodyValues`
and an exact allowlist of observed real V8 stack frames. Never accept a placeholder in provider error,
status, responseBody or reset evidence; never ignore unknown/mixed error records or extra terminal text.
Preserve nonzero exit, no valid artifact/no later success, every-error agreement, fixed 60-minute
terminal threshold, exact lastError identity when applicable, floor two, role and notification gates.
Keep quoted/tool-emitted spoofing negatives. Any broader change requires another explicit owner ruling.
claude-1 must review the actual native fixture, framing changes and adversarial regressions independently.

### Verification already available and remaining

Producer HOST checks pass: full Go suite (607.553 seconds), build/vet/race, shared/local/native crash
and lease checks, 86 changed Go files formatted, and the skill suite (399 Node, 54 Python, six manifests).
No cycle-2 normative protocol bytes changed. claude-1 read every changed product file (109 CLI, four
skill), ran focused checks in 13 packages (12 with matching tests) on both volumes, and authored the
counterexamples above. Full suite/race/vet/skill results are explicitly SECONDARY to that reviewer.
Windows is compilation evidence only; runtime/native crash settlement remains unverified/unavailable.

Before a later close, the next review should independently execute the relevant broad checks and a real
process-crash recovery exercise, or record any remaining evidence decision in the attended-close note.
The reviewer also raises a non-finding, unverified container/PID-namespace limitation; retain it as a
question, not a proven defect. No current close claim relies on accepting either residual.

Selected raw evidence is copied without edits and hashed under
`ideas/meta-protocol-change-quota-auto-exclude/source-context/round-05-review-20261006/`.
`organizer-usage.md` and `usage-ledger.jsonl` include the successful review and organizer snapshots;
cumulative snapshots replace earlier values, and no monetary estimate is claimed.

## What I need from you

- **Q1:** authorize the narrow cycle 3 above, or explicitly choose a different disposition of R5-MAJOR-1.
- **Q2:** choose **capture plus bounded evidence-based grammar fixes** (recommended), capture only,
  an explicit AC2 waiver accepting the practical limitation, or defer.

An answer such as **“Q1 yes; Q2 capture plus bounded fixes”** is unambiguous. It authorizes only those
next steps. Both final review signoffs and a NEW attended-close answer remain necessary before release.
The organizer deliberately ends this turn after preserving this blocking note; no process is left waiting.
