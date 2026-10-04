# codex-1 implementation evidence — quota auto-exclusion

Producer: codex-1, bounded implementer-only child invocation, 2026-10-04. Base HEAD `550fbf8`, branch `quota-auto-exclude`. No commit, merge, release, participant launch, or review/signoff performed. Both implementation stages are present in the working tree. This is producer evidence, not independent acceptance.

## Authority and attestation

Read the frozen FINAL (AC1–AC21), IMPLEMENTATION plan/resumption direction, the owner scope/reset answer, and review/round-02/claude-1.md in full. Rendered and read the full emitted body of:

```sh
parley protocol packet --dir . --phase 5 --track deliberation --idea meta-protocol-change-quota-auto-exclude --flag protocol_change --json
```

- `context_mode: full`
- `source_sha256: 8e9213bd45059069d484bd10e5ca1a1c509297039dfd8fc67d5e9ebda7590416`
- `packet_sha256: 8e9213bd45059069d484bd10e5ca1a1c509297039dfd8fc67d5e9ebda7590416`
- No fallback reason emitted.
- Body: `.parley-runtime/protocol-packets/full-phase5-deliberation-8e9213bd45059069d484bd10e5ca1a1c509297039dfd8fc67d5e9ebda7590416.md`.
- Captured command result: `/tmp/quota-implementation-packet.json`. The body was not edited. No driver phase was invoked.

The owner answer controls both prior MAJOR scope choices. Its explicit exception says: “For zcode, the `responseBody` JSON in stderr is enough when the process ended with an error, left no output, and every error record is a 429 'Limit Exhausted' with `reset_at`.” The detailed authorization includes same-record `reset_at` or `retry_after`, all error records agreeing, the terminal turn-failure line, no valid artifact/later success, and the fixed 60-minute threshold. Residual subagent-source ambiguity is accepted by the owner; no native root channel is claimed or invented. Display clocks use only offsets -12:00 through +14:00 in 15-minute steps, with all interpreted instants agreeing within one second.

## Delivered areas

- Stage 1: enabled the bounded zcode stderr recognizer, tightened malformed/mixed/duplicate/incomplete record gates, fixed reset parenthesis capture, retained the provider display and canonical machine UTC reset, and timestamped observed stderr records before teardown. Other adapters remain diagnostic-only. The preserved prototype's C1 filtering, bare-503 gate, standalone report-only preflight, presence-aware configuration, and fixed whole-batch floor remain wired into kickoff.
- Stage 2: immutable hash-bound batch history (`quota-history/NNNNNN.json`), checked file/directory durability, immutable replay, pending applied receipts, conservative contradiction/truncation handling, prompt/manifest recovery, terminal round re-evaluation, and one deduplicated notice or blocking escalation with arithmetic. Old events and failed/partial participant artifacts remain untouched by transitions.
- Lifetime serialization: an idea lease for enabled mid-idea scope, with same-context nested reuse and cross-process/cross-run exclusion. Separate short projection locks serialize signoff appends and commits. Unsettled prior writers block recovery until terminal telemetry establishes completion. Off/legacy/kickoff-only runs do not acquire the new lease.
- Consumers: runner batches and single-phase helpers, driver round/consensus/implementation adapters, draft/signoff/close paths, goal-check failure handling, replay, next-action planning, status, wait, and organizer brief use current membership. Common launch admission also rejects stale/excluded identities and a lease for another idea or run. Known signers come from immutable kickoff/history; required signers are current. Kickoff exclusions never become known.
- Retained obligations: immutable snapshots of excluded members' filed vetoes, disputes, and blocking findings survive later synthesis changes. Status remains blocked until an explicit disposition supplies authority, rationale and independent evidence; historical veto disposition specifically needs an owner ruling quoted in another artifact. Signoff bytes are not rewritten.
- Role and close guards: per-idea designee, implementation pin, started canonical draft attribution, track reviewer count, diversity, independent checker availability and strict findings are re-evaluated. No global roster mutation, timer, automatic rejoin, or repair of legacy gap 11 was added.

Only Go implementation/tests under `internal/`, quota fixtures/README, and this evidence file were edited by this invocation. Concurrent organizer-owned protocol/skill/orchestration changes are preserved. In particular this invocation did not edit any COOPERATION.md, IMPLEMENTATION.md, inbox or reviewer/signoff artifact.

## Implementation representations needing review

These are bounded representations of FINAL's required history, pending state and retained findings, not new dependencies or a manifest schema-version bump:

- Optional `quota_revision` in run manifests; history remains authority. Optional quota membership/pending fields in PhaseDigest are emitted for quota-backed ideas.
- Each batch can carry `retained` obligation snapshots with deterministic id, kind, author, canonical path, source hash and text. This prevents exclusion or a rewritten synthesis from silently deleting a filed obligation.
- A disposition is a contiguous paragraph in the current canonical consensus containing the obligation id and `Disposition: resolved|withdrawn|operator-ruling|no-dependency`, `Rationale:`, `Authority:` and `Evidence:` (relative canonical artifact path). For finding/dispute dispositions, independent evidence must name the obligation, identify an author other than the excluded author or pinned implementer, and contain an Evidence/PRIMARY marker. A veto specifically requires `operator-ruling`, `Authority: owner`, and `## User direction` in the linked next artifact. The gate checks attribution/structure; it does not independently establish the truth of an author's claim. Explicit owner rulings affect triage without deleting or rewriting the filed veto.
- `quota-driver.lock` and `quota-projection.lock` reuse OS locking primitives from existing dependencies. `quota-applied/<batch-id>` is the checked-durable reconciliation receipt. Incomplete writers and missing manifests gate; they are not silently reconstructed.

## AC-to-test map

| AC | Producer coverage |
| --- | --- |
| AC1 | Organizer owns protocol/skill wording; implementation support table and fixtures updated here. |
| AC2 | `TestQuotaZcodeRecordedIncident`, `TestQuotaIncidentOneReadinessWithRecordedProviderReplay`, `TestQuotaZcodeEveryRecordAndObservedClock`, `TestQuotaDisplayClockBoundsAndParentheses`. See incident-1 evidence limit below. |
| AC3–4, AC18 | `TestQuotaStrictSemantics`, `TestQuotaSuccessAndWatchdogsWin`, `TestQuotaUnsupportedAndAdversarialProvenance`, adversarial stderr fixtures, duplicate-key/wrapper/clock tests; valid artifact/later-success runner tests. Only supported invocation evidence reaches automatic candidacy. |
| AC5 | `TestQuotaBatchPermutations`, `TestQuotaFloorRolesAndSuccess`, `TestQuotaPreflightWholeBatchAndOneBlock`; process round and signoff batch integrations. |
| AC6–7 | `TestQuotaC1KickoffReplayAndFrozenScope`, `TestQuotaLegacyC1NoPolicyAndNoHistoricalExcludedSigner`, `TestQuotaAutomaticKickoffRecordsAndNotice`, `TestQuotaPreflightReportsOnlyAndBare503`, existing shared failure-classification tests. |
| AC8 | `TestQuotaStubSurvivorAndProtectedRolesBlockWholeBatch`, `TestQuotaFloorRolesAndSuccess` (including unpinned global default and facilitator). |
| AC9 | `TestQuotaExcludedVetoRemainsBlockedKnownAndAppendGateNarrow`, `TestQuotaKickoffExcludedIdentityNeverKnown`, `TestQuotaRetainedVetoOnlyOwnerRulingCanDispose`, strict retained-MINOR test. |
| AC10 | `TestQuotaReviewDiversityAndStrictGatesOnProspectiveMembers`, current consumer rebinding test, existing model-diversity/goal-check/consensus tests. Host-only driver budget/strict-close paths remain in the host suite. |
| AC11 | `TestQuotaBatchCommitFaultsAndTruncation` (create/write/sync/publish/directory-sync), `TestQuotaProjectionFaultsPendingReadOnlyAndRecovery` (prompt/manifest/evaluation/notice/applied, repeated faults and dedup), pending signoff/close tests, stale low-level dispatch refusal. |
| AC12 | `TestQuotaLifetimeLockDifferentRunAndOffScope`, `TestQuotaCrossProcessLease`, `TestQuotaUnsettledWriterFromPriorRunBlocksRecovery`. |
| AC13–14 | `TestQuotaTransitionReplayHistoryAndIncompletePreservation`, `TestQuotaRunnerProcessFailureTransitionsAfterWritersStop`, `TestQuotaLaterRunReadsHistoryAndContradictionsBlock`, `TestQuotaDriverReconcilesPendingBeforeRebindingEveryConsumer`. |
| AC15 | `TestQuotaStatusWaitBriefAgreePendingAndAppliedReadOnly`, transition notice/replay/fault tests; wait returns 3 while pending and 0 after reduction, never 4 merely for a successful reduction. |
| AC16 | `TestQuotaPresenceAwareDefaults`, `TestQuotaDeckOverridesMachine`, `TestQuotaNewScopeAndOldScopeImmutable`, `TestQuotaPolicyJSONRequiresPresenceAndHistoryRejectsDuplicateKeys`, frozen kickoff/resume tests. |
| AC17 | Existing `TestQuotaPolicyAndHints`; immutable scope/history tests. No roster-writing, rejoin or timer path introduced. |
| AC19–21 | Organizer-owned implementation status, full HOST checks and the separately configured claude-1 review/attended close remain outside this implementer invocation. |

## Evidence limits and remaining organizer work

Incident 2's four decisive stderr lines are retained verbatim in the local incident evidence and copied into the fixture. Incident 1 retains readiness classifications/gates, **not its raw stderr**. Its fixture honestly pairs those recorded readiness facts with the recorded incident-2 provider body; it is not mislabeled as a second native capture. The actual incident-1 survivor arithmetic still blocks (only codex-1 usable after the declared facilitator and unresolved kimi are excluded from the count). No provider was invoked to fill this historical evidence gap.

No new dependency, provider/model/effort/auth setting or roster change was made. The owner-accepted stderr source ambiguity remains. Independent review must assess the complete diff and these storage/disposition choices. The organizer owns serial commits, protocol/skill completion, full HOST suite, criterion acceptance, and attended close. This producer file does not issue an acceptance verdict.

Shared-memory recall was unavailable: `MCP tool call requires approval, but approval policy is never`. Local frozen sources were used; no memory write is claimed.

## Checks

The full `go test ./...` HOST suite was deliberately not run.

All successful commands below used `GOCACHE=/tmp/quota-auto-exclude-go-cache` because the default Go cache was denied by the sandbox.

1. `go test ./internal/driver ./internal/runner ./internal/app ./internal/config ./internal/consensus ./internal/quota ./internal/telemetry ./internal/membership ./internal/runcontrol -run '^TestQuota' -count=1` — PASS, all nine packages. Latest output `/tmp/quota-final-focused.log`.
2. `go test ./internal/consensus ./internal/protocol ./internal/runplan ./internal/runstate ./internal/runmanifest ./internal/runcontrol ./internal/telemetry ./internal/quota ./internal/membership -count=1` — PASS, all nine packages. Output `/tmp/quota-final-packages.log`.
3. `go test ./internal/app -run '^(TestQuota|TestConsensus|TestWait|TestPhaseDigest|TestOrganizerBrief|TestCheckModelDiversity|TestGoalCheck)' -count=1` — PASS (23.172 s). Output `/tmp/quota-app-regression2.log`. Two implementation regressions discovered by the earlier run (legacy missing-prompt fallback and an unsolicited legacy digest field) were fixed before this passing run.
4. `go test -race ./internal/membership ./internal/runner ./internal/app -run '^TestQuota' -count=1` — PASS, all three packages; includes cross-process serialization and local process fixtures. Output `/tmp/quota-race.log`. The later parser-only SDK-wrapper refinement also passed the final focused suite.
5. `go build ./...` — PASS, latest output `/tmp/quota-build.log` empty.
6. `go vet ./...` — PASS, latest output `/tmp/quota-vet.log` empty.
7. `git diff --check -- internal cmd` — PASS. Gofmt applied to changed/new Go files only.

Sandbox-limited checks (not weakened, not retried with privilege changes):

- Default Go cache initially failed: `open /Users/tomasfecko/Library/Caches/go-build/...: operation not permitted`. Subsequent commands used the permitted temporary Go cache above.
- A broader affected-package test command included full `internal/driver` and `internal/runner` package tests. Existing budget/trajectory launch helpers were denied with errors such as `open /Users/tomasfecko/Library/Caches/parley/budget-locks/.lock-identity-2244717876: operation not permitted`. This causes cascading failures in those host-dependent tests; their completion is deferred to the organizer's HOST run. No cache/lock behavior or tests were weakened to bypass this constraint.
- The app selection that also included `TestDriverPrecheckConsensusDraftRefusalAndRestore` hit the same restriction: `open /Users/tomasfecko/Library/Caches/parley/budget-locks/.lock-identity-4192918985: operation not permitted`. The rest of the selected app regression coverage passed after the two real regressions above were corrected.

No genuine provider quota/credit/auth failure occurred in this invocation. Provider-shaped failures in tests were offline shell/record fixtures only.

Final targeted follow-ups: `go test ./internal/consensus ./internal/membership ./internal/app -run '^TestQuota' -count=1` passed after making retained severity detection case/spacing tolerant and keeping the added review-status diversity check scoped to mid-idea policy. `TestQuotaRunnerLostRoundEventBlocksReduction` passed after gating a reduction on an infrastructure failure to persist the original round event. The final focused suite, build and vet were then rerun on the completed tree. No host-only test was skipped or weakened in source.
