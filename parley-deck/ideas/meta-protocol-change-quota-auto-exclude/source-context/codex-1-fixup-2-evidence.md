# codex-1 producer evidence — Phase 8 fix-up cycle 2

Date: 2026-10-06. Role: implementer ONLY. This is producer evidence, not independent acceptance, a criterion grade, a signoff, or authorization to close, merge, publish or release. The separate claude-1 re-review and both Phase-8 stages remain required. The organizer owns normative protocol/skill edits and commits after this process exits.

CLI workspace: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude`.
Skill workspace: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude-skill`, read-only in this cycle.
Every relative locator below is relative to the CLI workspace unless explicitly stated otherwise.
Raw log prefix **CHECKS** = `.parley-runtime/quota-implementation/fixup-2/checks/`.

## Authority and complete live protocol attestation

I read the controlling `runs-handoff/quota-auto-exclude-2026-10-03/IMPL-ORGANIZER-BRIEF.md` from the main repository; immutable FINAL (902 lines); living IMPLEMENTATION (464 lines at read); all four matching `parley-deck/inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_*.md` owner answers and the user-to-all ratification; complete `review/round-04/claude-1.md` including appended erratum/conditions; and the signed `review/consensus.md`, including the complete original claude-1 ACCEPT-WITH-RESERVATIONS V1–V9. I used V4's actual pre-change path, not the superseded generic interpretation. The source hashes and exact owner-answer names are in `CHECKS/authority-sha256.json`.

Obtained the live packet with:

```sh
parley protocol packet --dir . --phase 8 --track deliberation --idea meta-protocol-change-quota-auto-exclude --flag protocol_change --json
```

Attestation: `CHECKS/protocol-attestation.json`; `context_mode=full`; source and packet SHA-256 both `73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`; no fallback reason. The full body is `.parley-runtime/protocol-packets/full-phase8-deliberation-73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e.md`. I read the entire 1–1502-line body in bounded consecutive chunks (1–380, 381–760, 761–1140, 1141–1502), including the long normative lines, not just the index. This includes Phase 8, quorum/roles, stopping and churn, human-dependent threads, protocol-change publication and reservation obligations. The packet's transport is github-pr; the explicit run-specific owner brief controls this child's local-file/no-commit operation. No normative copy was edited. The final live protocol hash still equals the packet hash.

Skills used/read: parley-deck, graphify and openviking-memory. Graph context was consulted without rewriting it. Shared memory's scoped call failed with the exact tool error `MCP tool call requires approval, but approval policy is never`; local authorities supplied the evidence. No successful memory read/write is claimed.

## G10 / V1–V3 — semantic retries and receipt-bound threshold

Implementation: `internal/telemetry/quota_framing.go`, `quota_zcode.go`, `quota.go`, and `internal/runner/telemetry.go`. Complete retry records may differ in countdown while agreeing in exhaustion class and absolute reset. Exact `lastError` fingerprint agreement remains required; its duplicated dump is parsed/checked but is not a new observation. Every record still validates its machine duration, header, display/countdown and reset internally. Unknown framing is not discarded.

The receipt of the terminal line is the threshold anchor. The collector's stream offsets identify it, independently of later process teardown. Absolute reset minus duration infers an attempt observation only for consistency. Observations must not go backwards beyond one second (a high-water mark prevents accumulating tolerated backwards steps), be later than their receipt, or predate a recorded invocation start. An aggregate's receipt is not assigned to each nested attempt. The runner now supplies the recorded start. The minimum remains exactly 60 minutes. Zcode reports actual framing/consistency rejection text rather than the diagnostic-only adapter fallback.

Regressions: `internal/telemetry/quota_cycle2_test.go` tests 61-minute-first/59-minute-terminal rejection, exactly 60 minutes, backwards/future/pre-start observations, body/header/class contradictions, two complete top-level records in one write and real distinct SDK retry errors. Existing fingerprint/adversarial coverage remains. `quota_framing_test.go` now expects the terminal receipt when earlier records arrived separately.

Fixtures are SOURCE-DERIVED: `zcode-retained-retries-source-derived.stderr` and `zcode-retained-top-level-source-derived.stderr` use both retained 176930/176890 countdowns; `zcode-backoff-retries-source-derived.stderr` uses 176936/176934/176930 (6/4/0-second deltas). The unchanged identical-retry fixture is only a parser control. The fixture README and harness make these distinctions explicit.

Reviewer counterexamples: actual `.parley-runtime/claude1-r4/zretry/main.go`, copied into `CHECKS/reviewer-zretry/main.go` with ONLY its output-file locator changed by `CHECKS/prepare-differential.py`. `CHECKS/reviewer-zretry.log` shows B/B2/C-final/D/E now eligible, while B-first/C-first reject contradictory future observations. D uses separately timed collector writes; E uses a single aggregate write with 6/4/0 deltas. `CHECKS/reviewer-zadv.log` retains all 54 negative probe results (all ineligible); the partial native tail remains rejected in `CHECKS/reviewer-tail.log`. Their historical comments do not turn reconstructed inputs into native captures.

## G11 / V4–V5 — preserve the actual policy-off CLI path

Implementation: `internal/protocol/quota_manual.go`, `internal/quota/revision_snapshot.go`, `revision.go`, `history.go`, `internal/protocol/quota_retained.go`, and the narrowly scoped notice compatibility helper in `internal/membership/membership.go`.

Ordinary policy-off exclusions accept bracketed and unbracketed §9.0 confirmations, with reasons containing em dashes. The trailing ` — confirmed YYYY-MM-DD` is the delimiter. A trailing note, hyphen separators and inbox-reference-only forms reject with the accepted grammar and corrective action. A known return requires only a `participants:` edit. A catch-up join requires the existing valid late round-1 artifact plus the participants edit; there is no new mandatory included marker, return/join record, owner schema/directive, CLI command or read-priors/join-from metadata on that path. Reading priors and joining from round 2 remain protocol duties.

Each imported change is an immutable manual revision. It is never owner-confirmed authority. Retained-veto withdrawal now explicitly requires committed owner authority, so a plain return or optional included marker cannot unlock it. Existing owner-ruling/explicit authority paths remain. Historical batches, snapshots and obligations are retained. Exact pre-fix notices with the former owner-confirmed display label are preserved and receive a separate immutable manual-authority clarification; arbitrary altered notices still gate. This clarification is not authority and does not rewrite history.

Differential oracle: exported the ACTUAL FINAL commit `27e42b8` via `git archive` into `/tmp/quota-cycle2-prechange-27e42b8`, without creating a worktree. The exact full commit, archive command and original source hashes are in `CHECKS/baseline-identity.txt`. Read original protocol §5 (lines 744–754) and §9.0 (878–887): recorded excluded grammar and explicit re-inclusion duty, but no new return record. A common probe invokes actual `app.Run` → `consensus signoff` on both versions. Only setup differs: baseline `CreateIdeaFull`, current policy-off kickoff plus its genuine run manifest. Identical edits and late-round bytes exercise exclusion with em-dash reason, plain known return, and late round-1/plain join. All three CLI operations return 0 on both. The current side additionally asserts three immutable manual revisions and absence of owner authority. No return/join record was supplied.

Reproduce the differential preparation with `python3 CHECKS/prepare-differential.py`, then run `go run ./.cycle2-differential /tmp/NEW-BASELINE-ROOT` in the archive and `go run ./CHECKS/differential-current /tmp/NEW-CURRENT-ROOT` in this module (expand CHECKS; choose fresh scratch roots). Original exact invocations/outputs: `CHECKS/differential-baseline.log`, `differential-current.log`, `shared-v4-cli.log`, `local-v4-cli.log` and the preparation script. Source identities are hashed; this is not an emulation of the old CLI.

Regressions: `TestQuotaFixupKnobOffRecordedManualChanges`, `TestQuotaFixupOffModeCatchupKeepsPreChangePath`, consensus retained-obligation tests, and `TestQuotaCycle2LegacyManualNoticePreservedWithoutAuthority`. Actual review R2 grammar probes ran on both filesystems, using its `knob-only` argument to avoid the unrelated authority probe's scratch commits. Accepted and malformed outcomes are in `shared-grammar.log` / `local-grammar.log`.

## G12 / V6–V7 — later edit versus interrupted projection

`internal/protocol/quota.go:InspectQuota` checks durable applied receipts. A fully applied revision followed by a changed participants set, including an older recorded set, is an owner edit requiring escalation; content matching history alone is not evidence of interruption. The error explicitly names `parley quota revise`, preserves the edit, and directs preservation/restoration before explicit revision. A latest unapplied/corrupt receipt permits checked projection recovery only for that transition's before/after membership and policy. Unrelated edits remain intact and blocked.

`TestQuotaCycle2AppliedEditAndUnknownPendingNeverRewritten` covers applied old-set edits and missing/corrupt receipts with unrelated membership, verifies bytes survive two attempts and one escalation. Existing interrupted projection/recovery tests cover idempotent notice/evaluation replay. CLI docs contain the command guidance. No protocol or skill bytes changed in this child; the exact organizer wording is below.

## G13 / V8 — crash and missing-manifest recovery

`internal/membership/crash.go`, `membership.go`, `revision.go`, `internal/telemetry/record.go`, and `internal/pidlease/{lease,live_unix,live_windows}.go` implement the repairs.

New started records retain host, boot, supervisor PID and supervised process group. A missing terminal may be settled only with the same known host/boot and proven absent supervisor/writer PIDs plus absent writer-owned process group. PID absence on this Mac never settles an unknown/foreign writer. The shared dead-local predicate also governs stale leases. Windows fails closed. A distinct immutable `crash-settlement.json` records invocation, start-record SHA-256, identity, proof and time. Concurrent publication validates the winner; replay checks canonical bytes and binding and syncs them. It fabricates no terminal outcome, usage, quota evidence, completed artifact, owner permission or budget settlement.

All affected original/current manifests are checked before automatic, manual or explicit revision commit. Missing kickoff manifest fails before adding history. A pre-existing committed incomplete revision remains pending until the authentic original manifest is restored and checked recovery succeeds. No manifest, identity, accounting history or authorization is synthesized.

Regressions: `TestQuotaCycle2MissingManifestBeforeManualCommitAndCheckedRecovery` reproduces the missing-manifest import, proves no new commit, deliberately constructs the historical incomplete state, verifies signoff/dispatch gates, restores the exact original bytes and repeats recovery. `TestQuotaCycle2CrashSettlementProofDurabilityAndReplay` covers foreign/unknown/legacy/live identities, failed durable publication, concurrent idempotence and altered-start digest rejection. `TestQuotaCycle2NativeCrashWriter` uses this test executable and a local `/bin/sh` sleep stub, exits before terminal persistence, and tests actual OS proof.

The native crash test was EXECUTED on shared and local scratch storage and under race; it fails closed because this child sandbox denies the boot query. Existing `TestLeaseProcess` fails for the same missing boot identity. Exact raw evidence is preserved, not skipped or weakened: `CHECKS/sandbox-boot-probe.log` records `sysctl: sysctl fmt -1 1024 1: Operation not permitted`; `CHECKS/final-race.log` records `boot=""`. Host rerun remains necessary. Injected-proof tests and all other membership quota tests ran without additional reported failures; that is not a substitute for native host verification.

Operator limits: run checked recovery on the original host/boot only after the supervisor, writer and its process group have stopped. Unknown/foreign/identity-less starts stay blocked; investigate there and restore authentic terminal evidence if it exists. A lost original manifest needs trustworthy original bytes. No force-unlock or fabricated proof route is provided. D6 legacy/accounting work remains untouched.

## G14 / V9 — exact source and sink evidence, AC2 still open

Source: `/opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs`, installed package 3.7.7-13. SHA-256 `3e3433d90fa502e5d02498dfde6c2090df898331359bcfe5f3dbc9a1d00b685f`. Read the organizer source leads including the credential-free default-console reproduction; leads are not independent verdicts.

Reproducible extraction command: `python3 CHECKS/source-audit.py` (expand CHECKS). This reads source bytes only, saves exact bounded slices and hashes in `g14-extraction.json`, and runs only isolated SDK classes. Coordinates are zero-based Unicode CHARACTER offsets, end exclusive, not byte offsets.

| Source locus | Located behavior and bounded evidence |
| --- | --- |
| SUr, 4900280; console.error at 4900711 | `g14-stream-default-error-sink.txt`: SDK streamText defaults onError to `console.error(q)`; no inspect options at this call. |
| w9r 7119632; R9r 7123374 | `g14-stream-options.txt`, `g14-undefined-option-filter.txt`: streaming adapter options contain `maxRetries:0`, no onError override; filter only removes undefined values. |
| k7r 7167750, streamText call near 7169744; T7r 7185858 | `g14-adapter-call.txt`, `g14-adapter-runtime.txt`: own loop calls runtime.streamText with w9r options; default runtime is SUr. |
| ozr near 4442739; b7r 7179790 | `g14-own-retry-defaults.txt`, `g14-stream-retry.txt`: own retry defaults include base delay 2000, factor 2, jitter, maximum 11 attempts / 60000ms delay; error handling schedules retries or returns terminal error. Configuration names were read as source strings only. |
| SDK retry near 4854200 | `g14-sdk-retry.txt`: maxRetries zero throws raw error; RetryError grammar remains relevant to the SDK but does not establish an aggregate on this adapter path. |
| runPrompt hZt catch near 13085739 | `g14-run-error-sinks.txt`: headless catch writes `Error: ${T}` plus optional traceId and newline through e.stderr; verbose additionally writes Cause/stack. This is the final CLI terminal sink, distinct from SDK console.error. |
| `new Console` 12437593; console setting 7542650 | `g14-console-constructor.txt` belongs to a workflow child script; `g14-console-assignment.txt` configures a structured logger stream. Neither located site changes the SDK sink's inspect depth. |
| depth:null at 10660301 and 11968625 | Separate YAML debug output and metric exporter, retained in `g14-depth-unrelated-{1,2}.txt`. Literal searches found no inspect.defaultOptions/inspectOptions/console.error assignment. This bounded source finding is not proof against every possible external runtime setting. |

`source-audit.py` invokes absolute `/opt/homebrew/bin/node` with ONLY `LC_ALL=C`, `LANG=C`, `TZ=UTC`; no credential variables or NODE_OPTIONS are inherited. It never executes zcode's entry point, loads its configuration, imports a network API or calls a provider. Node 26.10.0 reports default inspect depth 2. All four sterile reproduction outputs are byte-identical to the saved fixtures (`g14-sterile-results.json`). Realistic nested messages yield `[Object]` at the default console sink. Deep inspect positives use depth 10 explicitly and are labelled source-derived; they are not claimed as the default native output.

| Input shape | Post-G10 behavior | Evidence class |
| --- | --- | --- |
| Complete APICallError record(s), exact allowed structure, consistent 429/exhaustion/reset and terminal receipt >=60m | Eligible, including decreasing countdowns and multiple top-level records sharing one receipt | SDK/source-derived and reviewer reconstructed probes |
| Complete RetryError aggregate with exact final lastError fingerprint, consistent distinct retries and receipt >=60m | Eligible | SDK/source-derived; not established as this adapter's native aggregate path |
| Default-console realistic nested request containing `[Object]` | Rejected for framing | Isolated source-derived default-sink reproduction |
| Retained four decisive incident lines or retained native tail beginning mid-object | Rejected for incomplete/unknown SDK framing | Partial native excerpts only |
| Exact-lastError mismatch, mixed error classes, contradictory header/countdown/reset, future/backwards/pre-start observations, below-60m terminal receipt, unknown/verbose extra framing | Rejected | Meaningful existing/new negative tests and original reviewer scripts |
| Complete original 24,833-byte native capture | Unavailable; no claim | AC2 remains open |

Concrete later owner decision, for organizer use after safe repairs/re-review: retain this fail-closed recognizer and decide whether to (a) defer native AC2 disposition pending an explicitly authorized complete capture with invocation/start/terminal-receipt facts, or (b) explicitly accept the bounded evidence deviation knowing the actual default-console `[Object]` shape and partial tails gate. Neither option is selected here. No parser relaxation, capture fabrication or AC2 completion is authorized by these tests. Future native capture would need separate owner authority and the organizer's launch policy; this child launched no provider.

## Settled operator wording for the organizer's skill guidance

The following text is ready for the organizer to apply to command-level skill guidance after this implementer exits. It matches `docs/quota-membership.md`. No COOPERATION change is needed for these implementation details. If the organizer nevertheless changes normative copies, V7 requires normative neutrality, identical copies, and disclosure in IMPLEMENTATION and the attended-close request; this producer has not made that decision or edit.

> With quota_auto_exclude false, edit the actual participants list. Record an exclusion as `excluded: [agent-id — reason — confirmed YYYY-MM-DD]` or the same text without brackets. Anchor parsing on the trailing ` — confirmed YYYY-MM-DD`; reasons may contain em dashes. A suffix after the date, hyphen separators, or an inbox reference alone is not this recorded form; rewrite it using the displayed grammar.
>
> A known participant returns through an ordinary participants edit. A new participant follows the existing catch-up procedure: read priors, write a valid late round-1 artifact, edit participants, and join from round 2. Do not require a new included marker, return/join record, committed-answer schema, exact directive or quota command for that compatibility path. The CLI imports these edits as immutable manual revisions, never owner-confirmed authority. Explicit owner confirmation remains a protocol duty for re-inclusion. Such a manual revision cannot release a retained veto; use the existing quoted owner ruling or committed owner authority route for that authority-sensitive gate.
>
> For policy-on ideas, preserve a proposed prompt edit for the owner, restore the recorded participants and policy, then use `parley quota revise --dir WORKSPACE --idea IDEA --run EXISTING_RUN --request revision.json` with the committed owner decision. An edit back to an older set after an applied revision is still an edit. `parley quota recover --dir WORKSPACE --idea IDEA --run EXISTING_RUN` only replays genuinely pending, checked projections; it does not silently apply owner edits or infer permission.
>
> Stop writers before checked recovery. A crashed invocation can be settled only with recorded matching host/boot, proven dead supervisor/writer and absent writer-owned process group. The separate immutable crash settlement binds invocation, proof and time. Old records without identity and foreign/unknown owners remain blocked; investigate on the original host and restore authentic terminal evidence if available. A missing original kickoff manifest needs trustworthy original bytes before recovery; never invent a run or owner authorization. Windows native crash recovery remains unverified/unavailable.
>
> Source-derived zcode fixtures are not complete native evidence. Complete consistent allowlisted records may qualify, but partial tails and default-console `[Object]` truncation remain blocked. Do not report AC2 met from these fixtures.

## Commands, outcomes and storage

All check commands use `GOCACHE=/tmp/quota-auto-exclude-go-cache`. `CHECKS/run-checks.py` and `CHECKS/final-checks.py` retain exact argv, working directory, elapsed time, status and environment overrides in every log and their JSON result ledgers. A task-local PATH guard denies installed provider execution; discovery stubs and test-owned `/bin/sh` fixtures are the only dispatch used. No actual zcode/provider/participant was launched. No models, effort, credential, roster or provider settings were changed in the task workspaces.

Shared acceptance base: `CHECKS/shared-acceptance/`, on `/Volumes/My Shared Files/.../quota-auto-exclude/`; local base: `/tmp/quota-cycle2-local-acceptance/`. Their `tmp/` directories are the explicit TMPDIR values for final membership/runner/app quota tests. The actual runner regression `TestQuotaRunnerProcessFailureTransitionsAfterWritersStop` uses `/bin/sh` and synthetic stderr; it passes in both final runs. These are producer filesystem/dispatch checks, not independent acceptance.

| Final command (with environment above) | Outcome | Raw locator under CHECKS |
| --- | --- | --- |
| `go test ./internal/telemetry ./internal/quota ./internal/protocol ./internal/consensus ./internal/loop -count=1 -timeout 3m` | exit 0; all five packages pass | `final-core.log` |
| `go test ./internal/membership ./internal/runner ./internal/app -run TestQuota -count=1 -timeout 3m`, shared TMPDIR | exit 1: only reported membership failure is native crash boot proof; runner and app pass | `final-shared-membership-dispatch.log` |
| Same command, local TMPDIR | same boundary; runner and app pass | `final-local-membership-dispatch.log` |
| `go test -race ./internal/membership ./internal/pidlease ./internal/telemetry -run 'TestQuota|TestLease' -count=1 -timeout 3m` | exit 1: native crash and existing native stale-lease test fail on unavailable boot identity; telemetry passes; no data-race report | `final-race.log` |
| `go build ./...` | exit 0 | `final-build.log` |
| `go vet ./internal/telemetry ./internal/quota ./internal/protocol ./internal/membership ./internal/consensus ./internal/pidlease ./internal/runner ./internal/app ./internal/loop` | exit 0 | `final-vet.log` |
| `gofmt -w` on changed/new Go files only | exit 0; exact file list in log | `gofmt.log` |
| `git diff --check` | exit 0 | `final-diff-check.log` |
| `sysctl -n kern.boottime` | exit 1, sandbox Operation not permitted | `sandbox-boot-probe.log` |
| Actual 27e42b8 CLI common differential / current CLI | all three signoff operations exit 0 on both, current immutable/manual assertions hold | `differential-baseline.log`, `differential-current.log`, `baseline-identity.txt` |
| Reviewer zretry, tail and zadv Go probes | tool processes exit 0; classifications detailed above | `reviewer-zretry.log`, `reviewer-tail.log`, `reviewer-zadv.log` |
| Reviewer grammar plus actual current CLI differential on both filesystems | process exit 0, accepted/rejected outputs preserved | `shared-grammar.log`, `local-grammar.log`, `shared-v4-cli.log`, `local-v4-cli.log` |
| Sterile isolated SDK-class modes (all four) | exit 0, exact fixture equality | `g14-sterile-results.json`, `g14-sterile-*.log`, `g14-sterile-*.stderr` |

Earlier attempts remain intact: `initial-focused.log`, `cycle2-first-regressions.log`, `cycle2-regressions.log`, `affected-tests.log`, `race-synchronization.log`, `build.log`, `vet.log`, `check-results.json`. The first regression run exposed stale pre-fix assertions and a test-template reuse bug; both were repaired before final checks. The first broad affected command also mistakenly named nonexistent `internal/preflight` (preflight lives in app); the corrected final vet/focused command covers app. Its full runner/app tests encountered many sandbox-denied `/Users/tomasfecko/Library/Caches/parley/budget-locks/.lock-identity-*` writes and process-identity restrictions. Exact errors and dependent assertion failures are retained, not labelled a passing suite and not worked around by changing accounting or weakening tests. The organizer's requested full HOST `go test ./... -count=1 -timeout 45m` plus skill suite remains after this child exits.

Native membership and lease regressions have not been skipped, disabled or relaxed to turn these runs green. Host runtime proof and Windows proof remain outstanding. No Windows runtime check was performed here. No full-native AC2 proof is available. The later owner decision described above is prepared; this process does not solicit or select it.

Source integrity: `CHECKS/product-source-sha256.json` hashes changed implementation/tests/docs/fixtures; `CHECKS/authority-sha256.json` hashes read authorities; `CHECKS/g14-extraction.json` hashes source slices; `CHECKS/raw-checks-sha256.json` hashes top-level scripts/logs/results. `CHECKS/reviewer-source-comparison.json` confirms all four committed reviewer scripts under `source-context/review-round-04-relaunch-20261005/independent-probes/` are byte-identical to the runtime scripts inspected. `prepare-differential.py` shows the only zretry copy transformation, its output path. Reviewer sources, outputs, reviews and signoffs were not rewritten.

Code was formatted and frozen before final checks; subsequent work only recorded offline source/check evidence and this producer document. No task commits, push, merge, publication, worktree creation/declaration/pruning, protocol/skill edit, or D6 repair was performed. The initial unrelated untracked historical run directories were left alone. Organizer/reviewer assessment remains open; these results establish what was executed and what still cannot be proven in the child sandbox.
