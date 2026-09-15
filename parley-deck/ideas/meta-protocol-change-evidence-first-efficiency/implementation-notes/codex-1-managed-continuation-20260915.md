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
