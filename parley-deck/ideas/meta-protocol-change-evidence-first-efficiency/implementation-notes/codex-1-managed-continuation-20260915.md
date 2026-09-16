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


## User-authorized continuation and Zcode candidate — 2026-09-15T21:12:42.482844+00:00

The user approved both pending steps: prospective Hermes-to-Zcode replacement
in the remaining audit quorum, and preparation of the full-six to full-four pilot
amendment. Codex remains organizer with Claude, Kimi and Zcode. Original kickoff,
FINAL, rounds and signatures remain unchanged. The resolved escalation was moved
byte-for-byte to inbox/archived. The new design is draft PR #74, branch
idea/meta-protocol-change-evidence-first-efficiency-v2. It has a Codex-owned
independent round; other owned analyses and signatures remain pending. This is
not ratified experiment execution or whole-implementation acceptance.

User ceilings: USD 15 TOTAL prospective live-experiment spending (pilot, separate
packet trial, retries and grading); 15 minutes per task per arm, shared across
that arm's workflow. Implementation/review spending is separate. No treatment
calls started. Enforced spend feasibility and grader disjointness remain design
work, with the original twelve tasks and packet trial intact.

Zcode invocation 07efbf63-7203-4f25-93f3-59f1d4a5a8a1 terminated exit 0 after
952.927 seconds wrapper time. Its owned note hash is
05394c040a861434424bc02dbcc040cdb2f47159ee1e0edfdbd6e5bc754a8c95.
It modified only its claimed usage.go and added its claimed test and note.
The note explicitly reports that Bash had no permission client; Zcode ran no
Go tests, git diff or gofmt. Its static fidelity statements are participant
testimony. Coordinator pre-correction hashes and exact-byte source snapshots
are retained in runtime zcode-candidate-verification/candidate-source/.

Coordinator execution against those unchanged candidate bytes:

- `go test -count=1 ./internal/telemetry -v`: exit 0, 3.206 seconds wrapper;
  package output 2.498 seconds. New envelope cases actually executed.
- First old-source and additional-test overlays used /var paths while Go used
  /private/var. They were ineffective: the first ran the candidate again and the
  second reported no tests to run. Both original logs remain retained and are
  rejected as control evidence.
- Corrected canonical-path old-source overlay from git 39107b1 returned exit 1
  (0.829 seconds) at `actual envelope not recognized`, as intended.
- Corrected canonical-path cross-adapter probe returned exit 1 (0.529 seconds):
  the Zcode envelope became claude/codex/opencode/kimi/unknown reported usage,
  both without and with a warning preamble. All ten cases actually executed.
- `gofmt -l` listed the new Zcode test file; formatting remains outstanding.

This candidate is NOT integrated or accepted. The observed cross-adapter failure
is sent to Zcode for both dispatch hooks to be restricted to the intended adapter
and for negative tests. Its new owned correction note and independent v2 round
are being generated in a measured follow-up, without shell or extra model calls.
Original candidate/note/test and ineffective controls are retained unchanged.

Terminal Zcode stdout still has unavailable normalized usage under the original
launcher binary. Raw provider-reported fields for this invocation are input
4,631,168; output 56,642; total 4,687,810; cache-read 4,521,536; cache-write 0.
No reported price or resolved model is inferred. These are reported aggregate
counters, not independently invoiced charges. The request-schema inspection of
its retained model IO found ten tool definitions (previous rejected helper run
had 1,873). This count supports the reduced tool-context claim only; it does not
prove absence of every descendant process. No raw request headers were copied.

Terminal inventory now contains 44 measured CLI attempts, 23 unknown monetary
costs, and USD 69.3177235 known CLI estimates. Three current calls (Claude R1/R2,
Kimi refused-verifier recovery, Zcode correction/amendment) are pending and
excluded. Total spending remains unknown; no sum of internal helper and parent
usage is asserted. Historical failed/empty calls remain in the inventory.


## Candidate verification and provider limits — 2026-09-15T21:38:57.490015+00:00

Claude R1/R2 invocation20d32972-9b5b-4239-b426-065495c03f39 terminated
exit0 after1342.614s, with USD7.287771 CLI-estimated cost and its own note hash
368bb33ff699e206a6d396d73f261a09a222a4376753f89d892a39f9fc9e1823.
Exact candidate source and tests were snapshotted before any follow-up. The
coordinator independently executed the three new R1 tests: pass11.727s wrapper,
10.211s package. The old guarded source with the disclosed check seam adapter
fails at intended live-Finish/concurrent-apply/guard-held assertions20.288s.
The existing registered-stop case passed on the same candidate in coordinator
default environment6.1s wrapper/4.707s package; restricted-environment failures
on both source versions remain retained and unexplained. Both-fail alone does
not establish causal exclusion. No original manual-launch failure is relabelled
fixed by R2 diagnostics.

Zcode correction invocationabb0ea55-1af2-4ea6-8b7a-c322f51df32c terminated
exit0 after399.633s, owned note hash
d8a6ca8c562308aa0a01305e4b1741e0318c94e1c30a1e50667e3f54b3f7ddd4.
Both parser hooks now restrict the Zcode envelope to that adapter. Independent
other-adapter controls pass all ten cases. Its new canonical test nevertheless
failed because it expected an empty Source rather than the existing normalized
value unavailable; those source bytes and failure are preserved. After a
recorded sequential claim, Codex corrected only that expectation. All telemetry
tests then passed3.404s wrapper/2.585s package; gofmt was silent. Final native
candidate source hash94cd6d4f520376fe7e56114bd5b120ed1536f7dd2de8cc8e20aeb9b84b211f9e;
test hashf9ca76ac5c00bacf844bf4d8e4e822c014011c4d31c45857717c49882feab2e8.
Zcode's notes remain byte-preserved and disclose its own lack of execution.

Kimi implementation invocationf071a935-befd-42ab-b20f-a5d34f5a85bd produced
no code or artifact. Its CLI printed HTTP429/APIProviderRateLimitError, K3
quota100%, then waited for a10757000ms retry; the measured process finally
hit30minutes and exited143. Preserve the provider observation AND normalized
timeout. Approximate reset01:51 CEST September16; no repeated retry or quorum
exclusion. Verifier-refusal recovery remains unimplemented by this call.

Claude follow-up60c845b1-c0b9-446d-8742-3869d939d6f5 ended rate-limit after
248.657s, with USD1.595971 CLI estimate, no review and no amendment round. Its
terminal JSON says subtype success but is_error=true; normalized failure is
correct. The provider reports reset02:30 Europe/Berlin September16. No review
or owner self-correction is inferred from this failed call.

Zcode supporting R1/R2 review179f3982-87d5-42e4-b732-a7512de962ea ended
exit0 after392.726s, source unchanged, note hash
49fe3322a1cee5f0c9effd33a2446013a8a345ef7ea84201c670b8c6c1a30128.
It reports no CRITICAL/MAJOR, while retaining three MINOR items (unguarded
archive-content invariant with no reachable writer identified, uncharged-intent
test coverage, unexplained registered-stop environment correlation) and a NIT
(error-path preview status). These are open supporting-review inputs, not
withdrawals or whole-scope acceptance. Its copied new-test hash is malformed;
participant-owned correction is still required. Its production-source hash
matches the retained candidate, but no corrected test-hash attribution is
invented. The reviewer ran no tests.

The same invocation appended its own v2 SELF-CORRECTION. Exact previous bytes
are retained as an unchanged prefix. Updated round SHA256
465f67a93ed402ef23612865efeb3f35ee89e80a5354975d48e7a76089c86e9e.
It concedes p95 reservations are not a provider spending bound, withdraws using
a largely-cached source-audit token count as pilot cost/impossibility evidence,
and corrects repeated duo seats by repetition half. Residual feasibility/role
balance assertions remain proposals awaiting cross-review. No treatment calls.

Terminal inventory:49 measured CLI attempts,26 unknown monetary costs,
USD78.2014655 known CLI estimates; total spend unknown. No internal-helper and
parent counters are summed. Quota failures remain in that inventory. Experiment
budget remains unspent by treatment: USD15 total, not these implementation and
review costs. Kimi and Claude remain in quorum despite provider limits.

A separate native validation branch now assembles the five disjoint changed/new
Go files into the430-file source/module manifest. Full tests and three-package
race run serially, then vet; compiler cache/build scratch are on shared storage
because native disk was low, source and default test temp roots remain native.
No candidate code has yet been copied into the integration branch. A verified
identical copy of the completed Claude task's208922547-byte compiler cache was
retained on shared storage before its disposable native copy was removed;
source, logs and test evidence were not deleted.


## Combined validation failed — 2026-09-15T22:07:50.478353+00:00

Build passed in 14.582s. Full `go test -count=1 -p 2 -json ./...` exited 1 after 1548.977s. There are 340 failed test events across seven failed packages. Race and vet were not run. Every source/module hash remained unchanged. Saved diagnostics show historical-worktree-unavailable errors, executable format errors for agents/runcontrol binaries on shared storage, and 10-minute package timeouts in app/runner/trajectory. Causes are under investigation; no failure is waived. Full logs and result hashes remain in `.parley-runtime/managed-continuation-20260915/combined-validation/`. Prior tool sessions 46232 and 99037 are no longer registered; their observed reporting edits are present, and running status has been corrected. Zcode owner correction invocation c428c10c-2ea2-47f4-a835-357b025ef5b5 failed exit 143 after 60.782s without output or changes; cause remains unknown. No treatment run occurred. OpenViking scoped recall again failed HTTP400, no valid session ID; no memory persistence claimed.


## Reproduced test isolation defect — September 16, 2026

Correction to the earlier validation-environment description: Go 1.27.1 `testing.TempDir` explicitly creates fixtures under GOTMPDIR (`testing/testing.go:1573,1613`), even when TMPDIR is native. The shared compiler-build directory was inside the integration Git worktree. Tests that expected a non-Git temporary root therefore inherited the real repository and its historical missing-worktree inventory. A control using the SAME compiled budget test binary (SHA256 2f4c8163a22cb0776fa6eeb6ce315df3def2b19d8e62f2b1fabf1216edc1136b) reproduced all three `TestActionIdentityActiveSessionPinsSemanticFields` failures with shared GOTMPDIR (10.582s), then passed with only GOTMPDIR moved to a native directory outside all repositories (1.368s). `git rev-parse` confirms shared-root discovery versus native-root non-repository refusal. This is a demonstrated test-configuration defect for those failures; it does not by itself explain every timeout or executable-format failure. No production source, assertion, historical worktree or ledger was changed.

A justified full/race/vet validation is now running from the same 430-file source manifest. A retained Go `-exec` wrapper copies each compiled test executable byte-for-byte to native storage, verifies its SHA256, sets native GOTMPDIR only for test execution, and records binary identity and exit status. Compilation/cache stay on shared storage because native disk has roughly 550 MiB free. Raw failed validation remains failed and retained. The wrapper changes fixture/executable placement, not source or test selection. Results remain pending in `.parley-runtime/managed-continuation-20260915/combined-validation-isolated/`.


## Zcode owner correction completed — 2026-09-15T22:16:20.785695+00:00

One targeted retry completed exit 0 in 337.297s, invocation `b288b45b-fcc8-4911-8108-90aa73004272`, with no source changes. Zcode appended its own correct full test digest and narrowed its causal claim to occurrence on both versions, without excluding candidate-specific contribution. All original review bytes remain an exact prefix; the corrected review and handoff were imported byte-for-byte. A minor explanatory typo remains in its appended character-count/deletion description (62 characters and missing e7, not 61/e78); the full corrected digest itself matches. No facilitator proxy edit or signature. The preceding failed invocation remains retained. Inventory is now 51 terminal CLI attempts, 28 unknown monetary costs, USD 78.2014655 known CLI estimates; total spend unknown, no treatment calls.


## Amendment round published; later validation failure retained

Draft PR #74 now publishes Zcode's independent round and appended owner corrections at `047ea492ccac284675f7e26eef26d9f6ac5616a7`. Round SHA256 remains `465f67a93ed402ef23612865efeb3f35ee89e80a5354975d48e7a76089c86e9e`. Claude/Kimi round-01 and cross-review/signatures are pending; no freeze or treatment call.

The isolated full suite passed agents, budget and driver, but `TestEvidenceVerifierProductionClosure/real-independent-execution` subsequently failed because the independent helper reported criterion unit as fail. The saved test output does not expose the underlying criterion diagnostics. No environmental or source cause is assigned yet. A scratch-only diagnostic overlay extends that helper error string; it changes no candidate files and will not count as candidate acceptance. Full validation continues to collect other package results; race/vet remain pending, and will not run on a failed full result.


## Current-source native validation — September 16, 2026

The unchanged 430-file Go/module source passed full validation (554.398s), six-package race (659.592s) and vet (1.247s). All source hashes remained unchanged. Compilation, cache, executables and test temporary roots were native. These three checks used no diagnostic overlay, test wrapper, assertion change or enlarged deadline.

The earlier native-fixture/shared-toolchain full run failed after 1106.286s: 28 package passes, three failed packages (app, runaction, telemetry), and the CLI no-test-files skip. Retained failures include three closure assertions, the app 10-minute timeout, runaction signal 9 (wrapper exit 247), and telemetry importcfg/link errors. They are not silently reclassified. The first diagnostic scaffold used the wrong package working directory and failed before material execution in 0.444s; it remains invalid control setup. The corrected shared-cache focused helper diagnostic passed in 32.049s and the matched native-cache control passed in 8.646s. Both used a scratch-only helper error-string overlay; neither is candidate acceptance or proof of the original failure cause. The subsequent full/race/vet checks used unmodified source.

The first native race command exited 0, but its shared-drive stdout later had 1460 NUL bytes and a hash different from the initially recorded digest. The observed damaged bytes and both digest observations remain retained; bytes matching the initial digest are unavailable. That trace is not complete durable evidence. Only race was repeated, with native stdout/stderr capture, fsync, complete package-terminal validation and byte-verified copying after termination. The full and vet logs retained their recorded hashes and had no malformed JSON.

Exact structured results, hashes, commands, the 430-file manifest and limitations are published in parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/implementation-notes/codex-1-validation-20260916.json. All 88 protected historical/protocol files remain unchanged. No candidate source change was introduced during diagnosis. Fresh complete-scope participant reviews/signatures, verifier-refusal and wider recovery work, live trials, final HTML and delivery-based follow-ups remain outstanding. No treatment call has started; implementation CLI activity is tracked separately in the current runtime handoff.


## September 16 report and verifier-refusal checkpoint

Integration HEAD b55eda9a6a568c43db72df78c7b491dda12af20a retains the same 430 Go/module bytes as the published passing candidate. Zcode invocation d67f256b-7201-4f81-9b18-fe0114826418 exited 0 after 1173.278 seconds, producing four owned Go files and its original note (SHA256 d45e32c8321101ba0ce3bebbf8f9c2cbccf47e51290249331bcf1db71cea574a). It did not execute tests. Coordinator native testing then failed in 1.033 seconds: trajectory test lines 325–327 use an undefined err and line 447 declares unused path. Runner compilation separately reported no space left on device. No test cases executed. The code defects and storage failure are distinct; initial source and all failed logs remain preserved. The candidate is unintegrated.

A bounded owner correction is running on the same four Go paths with a new owned note, also addressing the need to revalidate actual refused requested/terminal bytes on every recovered authority read. Parent reconciliation and a successful actual runner relaunch with bound replacement authority remain open. Terminal audit inventory: 52 attempts, 29 unknown costs, USD 78.2014655 known CLI estimates; total unknown. One additional Zcode correction is pending and excluded. No treatment has started.

The separate evaluation workspace now contains delivery/2026-09-16/report.html, built from unchanged assessments/2026-09-05. The generator accepts an optional --delivery-date and retains its old default. Historical design calls remain a separate subset from the continuing-audit inventory. Coordinator checks are explicitly distinct from pending independent acceptance. Ten existing report tests and fifteen new-output static checks passed; all 156 preserved historical assessment/report files match their prior hashes. The first generator command used Python without markdown_it and failed before output; the documented existing Hermes virtualenv built successfully without any environment change.

Ego-browser space 21 checked all panels, nonblank charts, four sizes (1440x900, 1280x540, 390x844, 320x720), search/filtering, keyboard tabs, print-media visibility and actual embedded checkpoint state: 29 checks passed. Desktop/mobile/low-height screenshots were inspected. Emulation does not dispatch beforeprint; physical print/PDF pagination and canvas scaling remain unverified. Browser space 21 was finished. Report SHA256 df5f37603ff82dbe3de0b2f64e0708c1adf93e8ec25603e0f9d2eaa078b35adf, 563594 bytes. This is a partial checkpoint, not final delivery; it starts no 14/30-day clock.

Kimi's independent v2 round is prepared with peer rounds read-denied, but may launch only after the reported 23:51:13 UTC reset. Claude's own independent round and current-source review remain required after the 00:30 UTC reset. Neither is excluded or substituted. OpenViking remains unavailable; no successful shared-memory write is claimed. Live protocol drift from the packaged skill is expected for this source-role deck; installed CLI/skill remain unchanged.


### Corrected source check and next allocation — September 16, 2026

Zcode correction invocation 4e8fd511-a882-43f7-b81f-5ba5918fd045 exited 0 after 735.479 seconds and left its original note unchanged. Its new note SHA256 is dc343e101014307aaaae902aea74375a3e6c565cb4b03f0d219bf0ef1df74671. The corrected trajectory selection passed with 93 test/subtest PASS events (110.536 seconds package; overall command 112.716 seconds). Runner compilation separately failed because the new ReadCapturedVerification call received two return values instead of three. No blanket failure attribution to storage applies to this run.

After preserving that source, Codex corrected only the assignment and ran the two runner tests. The first passed; the mutation test failed before the intended assertion because it omitted the actual charged-patch prerequisite. Its pre-fix source/log remains retained. Codex added the same real local RunMeasured changed-source attempt used by the first fixture, preserving the assertions. The mutation test then passed in 3.067 seconds. An exact initial shape-only-reader overlay failed at the intended stale-stop lineage assertion in 3.001 seconds; the failure is not a compilation, zero-test or setup failure. All these checks used native temp outside Git, native log capture and verified copies. No model was used for the mechanical corrections, no production source was changed by Codex, and no source was integrated. The separate gofmt change is whitespace only in the trajectory test file; the 93-event run preceded that formatting.

The corrected source snapshot is retained in managed-continuation-20260915/zcode-coordinator-corrected-source. Parent resolution still always follows the original refused invocation. Zcode now owns a bounded sequential parent-lineage implementation in the same isolated candidate; its new allocation includes reconcile.go and a new parent recovery test. Actual successful runner launch with the intended replacement identity and attended app wiring remain separate obligations. No independent acceptance or completion is implied.

Updated terminal inventory: 53 unique attempts, 30 unknown monetary costs, USD 78.2014655 known CLI estimates; total unknown. The additional parent-lineage call is active and excluded. The earlier HTML is an explicitly dated checkpoint and its recorded figures are not silently rewritten. Kimi and Claude independent amendment rounds remain prepared until their provider-reset timestamps. All 88 protected protocol/historical files remain unchanged.


### Print defect and post-reset Kimi round — September 16

Kimi's prepared independent v2 round launched at 2026-09-15T23:51:37Z, after the recorded provider-reset timestamp; observed launcher PID 35819, tool session 9330. This confirms the local attempt, not provider health or an authored artifact. Zcode's parent-lineage attempt remains under its existing 30-minute ceiling. No new model or quorum substitution is involved.

Source inspection exposed a real HTML print defect behind the earlier scaling limitation: draw() checked overview.hidden even when print CSS displayed that panel. The old report reproduces stale 354x230 bitmaps stretched to 630x210 when printing from Implementation after mobile viewing. The corrected template checks computed CSS visibility and listens for print-media transitions. Ego-browser space 22 reproduced the old defect, verified corrected 630x210 bitmaps with nonblank pixels under the same trigger, and verified correct dimensions on return to screen. The corrected screenshot was visually inspected and space 22 finished. Physical paper/PDF pagination remains unverified.

The prior report, template, inputs and hash-bound checks are preserved in delivery/2026-09-16/revisions/before-print-fix/. New report SHA256 d5abed7fd03edcd7d69dacd1893ff15d4d030e88a5064ba412c59befedbafdab. Its payload equals the prior checkpoint except build timestamp; ten existing report checks executed against the new output pass, and all 156 historical September 5 report/assessment files remain unchanged. Earlier 29 browser checks retain their original report hash, not silently reassigned to the new file. Three new old/new/return-to-screen controls cover this fix. The report is still a partial checkpoint, not final delivery.

Read-only CLI help discovery found Claude's advertised --max-budget-usd and Zcode's --max-turns; Kimi help had no matched dollar/token/turn-cap option. These are help observations, not an enforcement witness or a claim no other control exists. Complete prospective experiment dollar enforcement remains unverified. No treatment was launched by this discovery.


### Kimi independent round received; scoped owner correction — September 16

Kimi invocation `3b90ef32-593b-44af-a240-3ea0f008da0c` completed exit 0 in 206.57 seconds after quota reset. It authored its independent v2 round directly, SHA256 `589830c3baa5602e3866cf00254876d4c0163adcd5be2ddbda6db8e935d09140`, with no tracked-source drift. Requested model/effort were kimi-code/k3/max; reported model, usage and cost remain unavailable, not inferred. Terminal inventory is 54 attempts, 31 unknown costs, USD 78.2014655 known CLI estimates; total spend remains unknown. The HTML remains its explicitly dated earlier checkpoint.

Coordinator review found A2 incorrectly charges implementation/readiness history to the prospective experiment ledger; A1 proposes waiving the user's packet task-arm time cap; S9/S10 and ALT-B overgeneralize mixed source provenance and infer absence from an omitted LaunchBudget file; P6/R6 propose twice-observed-cost reservations and stopping after overrun, neither a provider-spending upper bound. The base 346-call arithmetic also does not include all proposed retries. These are disputed proposal details, not adopted requirements. A bounded Kimi-owned append-only correction is running with the original bytes preserved and peer rounds read-denied. Published runner/telemetry source and the actual LaunchBudget definition were supplied with exact provenance. No facilitator edits another participant's interpretation or signature.

The user already authorized USD 15 total prospective experiments (pilot, packet, grading and retries), separate from implementation/review spending, and 15 minutes shared per task per arm. These limits are preserved without asking again. No treatment has started; dollar enforcement remains an open gate. Zcode's original parent-lineage invocation remains active and is not restarted. Claude's independent round remains prepared for its recorded 00:30 UTC reset.


### Owner correction and parent validation outcome — September 16

Kimi appended its own correction in invocation c8ff6a7e-631e-4445-b81f-dc713cc7e10b (exit 0, 173.105s wrapper), preserving every original round byte as an exact prefix. Corrected round SHA256 acf1f9df26189a87984101ed59c00cc3d8f173d3f8bd9d39d9add28a00ea3a0d; new owned handoff SHA256 41fc4737fdf26846f2908739a0eecafe93fa1f380cb2d9d4817e9f173e7abc0e. It withdraws A1/A2's authorization errors, corrects mixed source provenance, and withdraws observed-cost reservations as proof of a hard provider-spend bound. Base/retry arithmetic remains a proposal requiring cross-review; not a frozen reduced task plan. The first correction wrapper exited 2 in 0.991s before any provider invocation because agents exec requires a new --artifact; the second used a new owned handoff while appending the old round. Both observations are preserved; only the actual provider invocation is counted.

Zcode parent-lineage invocation a513ddc3-6b1d-40f8-bd96-3f2b186eb3bd completed exit 0 in 1184.719s; owned note SHA256 d0903a2eded965ca37853e3ddd462a85b53f1fe831eb07964a333b393acb3bc4. It authored code and tests but ran none. Coordinator native tests then failed after 195.018s with source unchanged: 94 test/subtest passes, but the new clean positive lineage and dependent adversarial setups reject lawful replacement requested timestamps after recovery. The recovered predicate incorrectly reuses the ordinary requested-before-reservation window. This is a source defect, not an environmental attribution. Runner separately could not link because native storage was exhausted; no runner tests executed in this command. Retained JSON/stdout/config/manifest: managed-continuation-20260915/zcode-parent-focused/.

A sequential Zcode-owned slice is now implementing the chronology correction plus a narrow real RunMeasured bridge to the replacement invocation already bound by immutable recovery. The allocated files are telemetry/record.go, runner/telemetry.go, runner/trajectory_verification.go, trajectory/reconcile.go, their new bounded tests, the parent test and a new owned note. It may not modify app CLI wiring, other budgets, prior notes or historical protocol files. Attended app wiring remains a subsequent obligation. Source stays isolated and unintegrated until actual checks pass. All 88 protected files still match.

Terminal inventory: 56 measured CLI attempts, 33 unknown monetary costs, USD 78.2014655 known CLI estimates; total unknown. The active Zcode call is excluded. No treatment run or experiment spending has occurred. Claude remains a required participant with a prepared independent round after its recorded 00:30 UTC reset. Whole audit, amendment consensus, freeze, independent acceptance, experiments and final delivery remain incomplete.


### Kimi cap-source audit and actual local controls — September 16

Kimi invocation ab5dfc8a-1b3b-42d8-b077-4d7aa185a487 completed exit 0 in 94.377s, authoring its own source audit (SHA256 7b4fa0e661656d0e3c78ecac40844b1265b4d175cc750e74f4509451125542b6). It used a custom main-agent Read/Grep/Glob/Write/Edit allowlist with no subagents, twenty requested steps and one requested attempt per step. Cost and reported model remain unavailable. The tracked-drift entry names only IMPLEMENTATION.md, deliberately updated by Codex concurrently; it is not attributed to Kimi. No Go source changed. This brings the measured real CLI inventory to 57 attempts, 34 unknown monetary costs, USD78.2014655 known estimates; total remains unknown. Active/queued Zcode/Claude calls and local fake-provider controls are excluded.

The exact public documentation URL advertised by installed Kimi0.42.0 exposed controls not present in help: max_steps_per_turn, max_attempts_per_step, custom main-agent tool lists, background-print limits, and a kimi-provider max_completion_tokens clamp distinct from anthropic-only max_output_size. HTML, extracted text, URLs and hashes are retained in experiment-cap-discovery/kimi-docs-provenance.json. This discovery is documentation evidence, not a provider-spend guarantee.

Codex then executed five isolated local CLI controls with a dummy credential, separate KIMI_CODE_HOME and a single loopback listener. Sandbox canaries allowed that endpoint and denied external connections. The initial profile used an unsupported literal-IP network predicate and failed before Kimi started; this invalid setup remains recorded, and the corrected profile uses the supported localhost predicate. The success control made one request, sent max_completion_tokens17 and advertised only Read. Forced real Read sequences stopped after exactly2 and3 requests with the respective max-steps error. Controlled HTTP500 cases made exactly1 and2 attempts under their corresponding caps. Expected terminal exit1 is retained, not mislabeled as a model success. All five client controls and their assertions passed; raw requests/logs and scripts remain hashed.

Structured scoped evidence is codex-1-kimi-cap-controls-20260916.json. These controls prove only the observed local Kimi transport/configuration behavior. They do not prove real k3 output enforcement, input/compaction bounds, all background/error behavior, provider billing or total USD15 feasibility; they are not packet/pilot/grading treatment. Kimi's proposed print_max_turns0 is not used: the documented environment field accepts positive integers, so the control used print_background_mode=exit instead. No participant note was proxy-edited.


### Storage recovery, independent round and attended candidate — September 16

Zcode bound-runner invocation `15b1ca5f-07d1-45bc-9a3e-a3233e2f146b` has requested/started records but no normalized terminal. Its wrapper failed with ENOSPC while persisting collected output; wrapper and child were observed gone. Empty native output files do not establish the child exit, duration, usage or cost. The source and owned note were salvaged byte-for-byte in runtime `zcode-bound-runner-salvaged/`; no terminal was synthesized. Claude's queued prelaunch also failed on ENOSPC before any provider call, preserved separately. An unused synthetic sizing archive outside the fixture Git repo was copied, fsynced, hashed and checked unused before removal of only its redundant native placement. The retained archive and relocation manifest preserve its content; no historical Git evidence was pruned.

The immutable pre-app 436-file candidate check failed: telemetry passed (0.437s), trajectory passed (110.703s), runner failed (22.871s). Its sole test failure was parent-result comparison in the actual bound-relaunch positive test. A resolved-path, error-text-only diagnostic showed the sole difference was request_path: fixture /var versus canonical ticket origin /private/var. The first diagnostic overlay did not engage because its key used the alias; both failed runs remain retained. This establishes a fixture-origin mismatch, not justification to weaken production equality.

After storage recovery, Claude independently authored v2 round-01 in invocation `33cefc6c-f74b-4b4f-be3e-b3b1b21d4183`, exit 0, 677.39s wrapper, artifact SHA256 `95689ca42b652d684877e4f435f227fe54f65d5703a5018ddcecd55690392a47`, no tracked drift. CLI cost estimate USD 2.923397; singular reported model remains null, model list contains Haiku 4.5 and Opus 5. This is one participant, not additional quorum membership. Its genuine source-access gap was remedied with common exact public task/helper copies and provenance; hidden grading keys/tests were not copied. Published budget source and local client-control evidence were also copied. Draft PR #74 commit `5ecac7b` contains these inputs, Claude's exact round and Codex's own round-02 counterproposal; all three peers now have bounded own round-02 allocations. No amendment consensus or experiment freeze is claimed. Three trailing-whitespace lines in exact copied documentation are retained to preserve provenance.

Kimi completed attended app recovery candidate invocation `779d36ec-597c-4b36-afb2-1c11573fd3de`, exit 0, 1216.285s wrapper, own note SHA256 `f8307642a1dcd57f4ce39ce1514929c4b7d222556f14c7949535fe59d117170d`; requested k3/max, reported model/usage/cost unavailable. It implements a guarded read-only recovery reader, attended recovery preview/apply, separate bound relaunch through the real runner and immutable recovered-parent evidence retaining the original failed parent result. The note and sources are candidates, not executed-test proof or acceptance. Codex prepared a 439-file native snapshot, canonicalized only the runner fixture origin in that snapshot, and started focused telemetry/trajectory/runner/app tests with native cache and temporary files. Original candidate, participant notes and failed evidence remain unchanged.

Stopped real CLI inventory at this checkpoint: 59 normalized terminals plus one stopped invocation without terminal; 35 terminal costs unknown plus the incomplete unknown, USD 81.1248625 known CLI estimates. Three newly active cross-review calls are excluded. Fake-provider controls and provider-free prelaunch failures are not real-model inventory. Total monetary spend remains unknown. Prospective treatment remains not-run, and its separate USD 15 total / 900-second task-arm ceilings remain binding. All 88 protected historical files were checked unchanged. Shared OpenViking remains unavailable; no successful memory write claimed.


### Cross-review admission refusal and corrected test scope — September 16

The three newly prepared amendment round-02 attempts completed with `budget_refused` BEFORE any CLI process started; all started_at/PID values are null and no participant round was written. These are three additional normalized pre-provider requests, separate from the 59 stopped real-CLI terminals plus one missing-terminal attempt above. No usage or provider cost is inferred. Read-only `budget migrate inspect --kind cross-review` identifies the first unavailable registered historical worktree (`.../scratchpad/f2repo`); Git also retains unavailable `/private/tmp/revert-test2`. Registrations/history are preserved. No phase relabeling, uncharged direct launch, clone-based accounting reset or synthesized historical checkout is used to get past this refusal. Amendment peer cross-review remains pending. Claude now has a separate original-audit SOURCE review of supported recovery and N1 accounting; Zcode separately reviews Kimi's new recovery layer. Neither task may author the blocked amendment round or advance its phase.

The first 439-file snapshot check completed FAIL in 213.964s with 64 passing test/subtest events: runner/app could not compile because Codex's snapshot omitted embedded `internal/protocol/defaults/COOPERATION.md`; trajectory's rewritten-terminal adversary correctly failed an earlier exact requested/terminal guard than its expected message. The initial snapshot and failed logs remain. A new snapshot restores the exact embedded file and corrects that one expected error string while retaining mandatory rejection, unresolved-state and no-resolution assertions. The original candidate is unchanged. The corrected focused run remains active; telemetry and trajectory package completion so far are passing (0.459s / 201.143s), not whole-run acceptance.


### Current source reviews and resumable publication — September 16

Codex added a minimal candidate CLI route `trajectory recover-verifier-parent` after an actual app publication-obstruction control showed that completed replacement evidence lacked a standalone CLI publication path. The original failed control is retained (15.658s); the corrective command and exact terminal-guard case passed in 27.234s, three test/subtest events, source unchanged. It uses existing revalidated library APIs, publishes no new model work and preserves charges and the original failed result. Full source still awaits integration/acceptance.

Kimi independently reviewed those coordinator changes in invocation `2656c86f-a331-4c77-bc56-fdf5ea6f4898` (exit0,569.09s); own review SHA256 `7fbe553371886990856f0d357f6ae22d2cda838ca78bd38b5ed6068a340bcc48`. It concurs that the missing CLI route is a recovery gap and changes its earlier ergonomics position. Two NITs (usage output and exact refusal assertions) have been staged as F1/F2 source corrections, with the active full-test source untouched. Zcode invocation `b5f95c7f-602d-43d5-8c44-e2858545ed94` completed exit0 in664.697s, own review SHA256 `f56ffeb91dd577b065a0011049ad0763fa37c7185400b6a664c34b1c68b7983b`. Its alleged concurrent-publish race explicitly omitted the enclosing withState guard; a narrow owner follow-up now examines that actual callback serialization and its two non-defect/unknown NIT classifications. No facilitator edits the reviewer artifact or counts an unresolved finding as withdrawn. Both CLI costs remain unavailable; Zcode's structured token usage is retained.

Two earlier Zcode source-review attempts failed before authoring: unsupported `--allowed-tools`, then native Unix IPC socket creation refused by the sandbox. Corrected invocation used the supported deny-list and a short private task TMPDIR, with a checked owned-file/Unix-socket canary. A too-long socket canary also failed before any CLI launch and is not a provider attempt. These failures and the source-adapter gap remain visible.

Claude's owned history-admission review (`4d45120f-b183-4566-9a66-9625f4c4b2b2`,exit0,475.543s,USD3.4236355 CLI estimate) identifies that no existing attended command can enumerate the missing historical registrations for first binding. A read-only observation found27 registered worktrees and exactly2 unavailable paths, preserving registration bytes and unknown history. The user has been asked about original backups; no answer or attestation inferred. Claude is implementing only the read-only `budget worktree inspect` first stage. It does not change launchScope refusal, create a policy, attest missing history, prune, reconstruct roots or grant a count. The eventual amendment peer round should use the existing grouped runner to share one logical cross-review charge.

Stopped real CLI inventory now:64 normalized terminal records plus1 stopped invocation missing terminal;39 terminal costs unknown plus that incomplete unknown;USD84.548498 known CLI estimates. Three additional normalized requests were refused before provider/CLI spawn and are reported separately. Active Claude inventory and Zcode lock-review calls are excluded. Total monetary spend remains unknown. No prospective treatment call, experiment freeze or final audit delivery has occurred. All published protocol/historical authority remains unchanged.

The current full Go test attempt includes the new publication code but omitted VERSION and the live protocol/applicability-map fixtures from its coordinator snapshot. It is therefore non-passing and will not run race/vet. Missing inputs are prepared byte-exact for the next complete snapshot; the active source and failed evidence are preserved. Prior published430-file full/race/vet evidence is not extended to these candidates.
