---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-15
source-commit: 39107b138e77de064f21fec4aafa6073887c2783
status: in-progress
---

# Managed continuation: independent findings and native implementation slices

The user requested that Codex conserve tokens, primarily coordinate, and invite
Claude, Kimi and Zcode. The existing audit scope and draft PR #73 remain open.
No historical quorum, pilot arm, numeric ceiling, signed artifact or release
was changed. Three disjoint native implementation candidates now run; none is
integrated or accepted by this report.

## Imported evidence and its limits

Each artifact below was authored by the named CLI, copied byte-for-byte only
after its command terminated, and checked against its terminal artifact hash.
All 428 original Go/module files and 88 protected historical/protocol files
matched the prior checkpoint during import. Full source manifest:
`6149bf1db913498f6b4495eaac24853facd5df0e7c1bb91c37f1e66020ff5995`.

- `claude-1-current-hardening-review-20260915.md`, SHA256
  `e9a12e6ae78948ec3542f9ef7ae3d259b734a664b2301788626aa27eedc57dcc`:
  supporting source review, not full Phase-6 acceptance or signoff. It narrows
  N1/N2 and dispositions N3 for mismatches present at check time. Remaining R1
  is a full-history guard hold in reservation recovery; R2 is lost content-free
  diagnostic stdout on manual-launch test failure; R3 is the scope-check-to-
  ticket-publication window. Cross-review/step precharge residuals remain open.
- `kimi-1-recovery-next-slice-20260915.md`, SHA256
  `650e10e09b5ac8a5377c4cf7ecd28303b7b58ea294eb604203be10c8cff91280`:
  source proposal for verifier-launch refusal recovery; its initial assertion
  that missing request bytes imply no retained ticket handle is not accepted.
- `kimi-1-recovery-reachability-20260915.md`, SHA256
  `613426c5ab7186237253a239cc188f2b5016c5b7da181b2c397000a09897e771`:
  two native tests were reported as executed against unchanged production source.
  The retained test SHA is
  `b614a41d2c97ad2db2fd2abfbc3f036cc6ef831c293296933165ce1ce290321f`.
  The test source constructs an actual changed-source attempt and captured ticket,
  then a real pre-start launch-budget refusal; it asserts retained accounting,
  no child spawn and refusal of the sampled later public paths. A second test
  retains a real handle before removing/corrupting its fixture-local request,
  refuting the stronger no-handle inference. Its concurrent Prepare case is an
  unscheduled two-goroutine race, not a deterministic pause inside publication.

Kimi's owner-issued `PROVEN` strengthening is not independent acceptance under
protocol Section 15.1. The test code and terminal logs are retained as evidence
pending non-owner disposition. The report's claim that the first build.log was
preserved under that name is inaccurate: run_tests.py overwrote build.log. The
original compiler diagnostic does survive in the immutable invocation stdout,
including `invalid operation: tickets[i] != zero (struct containing
trajectory.CapturedRequest cannot be compared)`. No log has been relabelled.
Observed old-handle refusal before replacement does not establish safe
reactivation/quarantine after new authority is published. Quarantine is not
being implemented from that inference.

## Zcode attempts rejected; a concrete telemetry gap

The first invocation `30292ee0-2a1d-4413-991f-d8051cc025f1` failed because the
installed runtime rejected --allowed-tools although its help advertises it.
The corrected invocation `b25d8e16-271c-4a1e-8a57-02e176a19934` exited0 but left
an EMPTY artifact (SHA256 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855).
It was not imported and contributes no completed review. Its own final response
reports background transcription through a helper, contrary to the no-helper
brief. Session metadata records two internal Agent sessions: one completed after
creating an empty file; another still said running after the parent exited.
An actual process inventory found no task-matching Zcode process; unrelated Zcode
processes belonged to another workspace and were left alone. This is not a
same-UID writer fence or a general descendant-inactivity assertion.

The actual corrected CLI stdout contains a warning prefix and JSON `usage`:
source=provider, modelRequestCount=42, inputTokens=4119566, outputTokens=18299,
totalTokens=4137865, cacheReadTokens=3996160, cacheWriteTokens=0. Normalized Parley
terminal usage is nevertheless unavailable. This retained payload is the input
for a new parser candidate. No monetary value or single resolved model is
inferred. The completed helper separately reports301284 tokens; overlap with
the parent aggregate is unknown, so these are NOT added together.

The old call exposed1873 tool schemas in its retained main requests. The new
Zcode worktree configuration disables MCP, subagents, skills, memory and plugins,
plus explicitly disallows the actual Agent/SendMessage tools. This is an isolated
launch correction; global configuration is unchanged. Reduced actual tool/context
counts have not yet been observed for the running new call.

## Allocated implementation candidates

Exact paths, branches and native worktrees were claimed in IMPLEMENTATION before
edits. File-set intersections were checked empty.

- Claude: reservation recovery validation outside the control guard, deterministic
  contention and authority/replay checks, plus diagnostic stdout in the manual
  launch test. R3 and N1 spending residual remain separate.
- Kimi: explicit recovery of a fully evidenced pre-start budget-refused verifier
  launch. Retain original ticket/launch/accounting; new invocation separately
  checked and charged; stale stop cannot control replacement. No quarantine,
  automatic retry, refund, custody or completion permission.
- Zcode: recognize the actual warning-prefixed CLI usage envelope with positive
  and discriminating negative parser tests; preserve unknown cost/model and
  existing format behavior.

Focused candidate checks come first, followed by independent review and one
appropriate integrated validation. This report performs no repeat of accepted
Go suites. All original unexplained test failures remain open as recorded.

## Invocation inventory and remaining full scope

Five new terminal invocations are recorded: Claude source review1094.809s,
Kimi source proposal296.014s, failed Zcode2.534s, empty Zcode561.527s, and Kimi
native probe397.570s. Only Claude reports a monetary estimate, USD6.013931.
Its reported_models lists Opus5[1m] and Haiku4.5; no single resolved model is
claimed. Added to the prior38 terminal attempts, the checkpoint is43 terminal
attempts,22 unknown costs, and USD69.3177235 known CLI estimates; total cost
remains unknown. Three subsequent implementation calls are pending and excluded
from terminal totals. Internal Zcode sessions are disclosed separately, not
fabricated as measured CLI invocations. The corrected Zcode attempt is linked
here to the failed one; normalized retry_of remains null in the retained record.

The full goal still requires current complete-scope reviews/signatures,
recovery beyond these slices, real closure/concurrency evidence, the frozen
packet comparison, the12-task pilot, final offline HTML with ego-browser QA,
and delivery-based14/30-day follow-ups. Missing numeric experiment ceilings and
explicit historical quorum/pilot decisions were asked asynchronously again on
September15; no answer was present when this note was written. Independent
implementation continues. OpenViking scoped recall returned no matching context.

Runtime evidence and live handles:
`.parley-runtime/managed-continuation-20260915/`.
