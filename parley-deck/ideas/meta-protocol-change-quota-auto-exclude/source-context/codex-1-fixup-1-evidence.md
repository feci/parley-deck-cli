---
producer: codex-1
idea: meta-protocol-change-quota-auto-exclude
phase: 8
fixup_cycle: 1
role: implementer-only
context_mode: full
source_sha256: 73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e
packet_sha256: 73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e
---

# Fix-up cycle 1 producer evidence

This file records implementation and executed checks. It is not an independent review, a code grade,
a signoff, or authority to close, merge or release. The organizer remains the other codex-1 process;
claude-1 independently reviews. Both stages and their remaining review/owner gates remain required.

All paths below are relative to the CLI worktree unless absolute. `C/` means
`.parley-runtime/quota-implementation/fixup-1/checks/`. The source snapshot checked by the final
commands is enumerated with SHA-256 hashes in `C/checked-source-manifest.json` (53 product, test,
fixture and CLI-documentation paths, including two deleted flock implementations). Product code
was held stable after final checks started at approximately 2026-10-04 04:18 UTC.

## Complete protocol and authority attestation

I read the full emitted Phase-8 protocol packet, all 1,501 lines / 125,862 bytes, in sequential
chunks, not only selected sections. The command was:

```sh
parley protocol packet --dir . --phase 8 --track deliberation --idea meta-protocol-change-quota-auto-exclude --flag protocol_change --json
```

The returned packet metadata is preserved in `C/protocol.json`; the complete emitted body is
`.parley-runtime/protocol-packets/full-phase8-deliberation-73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e.md`.
Attestation:

- `context_mode: full`; `optimize: false`; no fallback/fallback reason was returned.
- `source_sha256` and `packet_sha256` both equal
  `73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`.
- Source: this worktree's `parley-deck/COOPERATION.md`; role `source`; authority
  `live source file (protocolRole: source); global core drift is expected and never substituted`.
- Request: phase 8, deliberation track, transport `github-pr`, this idea slug, flag
  `protocol_change`, emitted `facilitator_participates: true`. These are the emitted metadata;
  they do not override this assignment's implementer-only role or appoint me organizer.
- Final read-only hashing confirmed identical source and packet bytes; see
  `C/environment-protocol-final.log`. `cmp` against the read-only skill worktree's
  `skills/parley-deck/references/COOPERATION.md` returned 0.

I also read in full the controlling
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/runs-handoff/quota-auto-exclude-2026-10-03/IMPL-ORGANIZER-BRIEF.md`,
immutable `FINAL.md` (902 lines), living `IMPLEMENTATION.md`, the owner's
`parley-deck/inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_scope-reset-answer.md`,
`review/round-03/claude-1.md`, and the signed `review/consensus.md`, including the ENTIRE raw
claude-1 ACCEPT-WITH-RESERVATIONS signoff and its R1–R4. I treated every CRITICAL/MAJOR and G1–G9
as binding, and R1–R4 as implementation conditions, not optional advice. Final authority hashes
and line counts are in `C/authority-hashes-final.log`.

The parley-deck, graphify and openviking-memory skills were read. The graph query is recorded in
`C/graph-query.log`; implementation relationships were then inspected in source because this graph
principally indexes documents. Shared memory discovery succeeded, but the read was denied with
`MCP tool call requires approval, but approval policy is never`. I continued from local authorities;
no successful OpenViking read or write is claimed.

## G1 — checked shared-volume durability

Sources: `internal/fsutil/sync_darwin.go:12`, `internal/quota/record.go`,
`internal/protocol/quota.go`, `internal/app/quota.go`, `internal/runcontrol/runcontrol.go`, and the
new lease publication paths in `internal/pidlease/lease.go`. Newly introduced raw file/directory
sync calls now use `fsutil.SyncFile`. Darwin ENOTTY invokes ordinary fsync; other errors and
fallback errors remain fatal. The only `.Sync()` in the membership implementation is the existing
higher-level `store.Store.Sync()` durability operation, not a raw file call.

Executed tests: `internal/fsutil/sync_darwin_test.go:10` exercises full success, wrapped ENOTTY,
fallback EIO and fatal initial EIO. `internal/app/quota_fixup_test.go:51` executes actual `Run`
creation and round-1 stub dispatch with both `--quota-auto-exclude=true` and `false`. Each command is
printed verbatim with stdout/stderr in `C/shared-volume-final.log:5` and
`C/local-control-final.log:5`. The commands have this exact argument sequence (each concrete root
is in the log):

```text
parley run --dir SCRATCH --no-tui --no-auto --no-ping --no-preflight --yes --quota-auto-exclude=true --participants codex "Shared filesystem acceptance"
parley run --dir SCRATCH --no-tui --no-auto --no-ping --no-preflight --yes --quota-auto-exclude=false --participants codex "Shared filesystem acceptance"
```

Both creation commands returned 0 on AppleVirtIOFS and /tmp; the created kickoff policy and actual
round artifact were validated. Shared-volume paths are under this shared workspace's `C/shared-temp/`.
No APFS-only t.TempDir inference is used for that acceptance.

An additional native diagnostic, outside product code, invoked F_FULLFSYNC directly and then the
checked helper on the same file and directory:

```sh
GOCACHE=/tmp/quota-auto-exclude-go-cache go run .parley-runtime/quota-implementation/fixup-1/checks/sync-native/main.go "$PWD/.parley-runtime/quota-implementation/fixup-1/checks/shared-temp" /tmp
```

Exit 0, `C/native-sync-final.log`:

```text
shared file: F_FULLFSYNC=inappropriate ioctl for device; fsutil.SyncFile=<nil>
shared directory: F_FULLFSYNC=inappropriate ioctl for device; fsutil.SyncFile=<nil>
/tmp file: F_FULLFSYNC=<nil>; fsutil.SyncFile=<nil>
/tmp directory: F_FULLFSYNC=<nil>; fsutil.SyncFile=<nil>
```

The log preserves the complete concrete paths. `C/environment-protocol-final.log` records
`/dev/disk0 on /Volumes/My Shared Files (AppleVirtIOFS, ...)` and the APFS local control.

## G2 / R1 / R2 — immutable authorized revisions and ordinary off-mode changes

Sources: `internal/quota/authority.go:90`, `internal/quota/revision.go`,
`internal/quota/revision_snapshot.go`, `internal/quota/history.go`,
`internal/protocol/quota_manual.go`, `internal/membership/revision.go:19`,
`internal/membership/revision.go:78`, `internal/app/quota_revision.go:18`, and the current policy
surface in `internal/app/quota.go`. Product instructions and the exact JSON contract are in
`docs/quota-membership.md`, linked from `docs/cli-reference.md`.

`parley quota revise --dir ROOT --idea IDEA --run EXISTING_RUN --request FILE` appends an immutable
owner revision. It requires a committed user-answer path, actual commit object, blob ID, full-content
SHA-256 and verbatim quote with a standalone exact current-set/policy directive. Duplicate JSON keys,
trailing records, unsafe IDs, unavailable/changed/wrong-idea/self-authored evidence and absent original
run identity are rejected. The kickoff and older batches are never rewritten. The current policy is
replayed from authorized revisions; merely resuming a kickoff-only idea does not widen it.

R1: permitted inbox move/delete works because committed bytes and the quote remain authoritative;
a contradictory existing original/archived copy is rejected. R2: normal off-mode `participants:` and
recorded `excluded: ... — confirmed DATE` / known-return `included:` confirmations are imported by
ordinary mutation paths. No new CLI command is mandatory. New joins require the existing catch-up
conditions, valid referenced priors, a late round-1 snapshot and owner answer. A new off-mode join
can record that answer's committed binding in its late-round frontmatter without `quota revise`.
Replaying immutable history revalidates the captured confirmation/catch-up content. Closed ideas
refuse membership revisions until the existing reopen procedure is used.

Executed coverage (all passed in `C/focused-final.log`, also in the relevant volume logs):

- `internal/membership/revision_test.go:17`: owner revision, known return, current/known sets,
  immutable prior bytes, current manifests, missing/wrong/changed/fabricated authority, allowed
  archive/delete and new-identity catch-up refusal/acceptance.
- `:101`: explicit scope widening and competing-run refusal; no binary/resume inference.
- `:127`, `:167`: ordinary off-mode exclusion, return and new catch-up join with owner evidence;
  snapshot replay survives permitted inbox cleanup. No off-mode driving singleton.
- `:197`: closed membership frozen.
- `internal/quota/authority_test.go:17`, `:76`: committed self-author/wrong-idea negatives,
  non-commit Git object, fabricated path, changed archive, unquoted content, malformed manual
  snapshots and fabricated catch-up proof.
- `internal/app/quota_fixup_test.go:89`, `:312`: actual revision/recovery CLI and normal off-mode
  signoff CLI after already recorded confirmation. Exact CLI output is in the logs.
- `internal/consensus/quota_test.go:224`: off-mode exclusion → known return → author withdrawal →
  fresh signoff, without the revision CLI.

Historical known/current consumers use replayed membership; reconciliation updates all extant
idea-bound manifests, not only the batch-originating run. Status surfaces explicitly show current
and known sets and label a past exclusion as historical after return. No roster or provider is edited.

## G3 / R3 — deck-filesystem PID/token leases

Sources: `internal/pidlease/lease.go:48`, `:53`, `:98`, `live_unix.go`, `live_windows.go`,
`internal/membership/lock.go:28`, `:97`. The former membership flock files were removed.
The lease retains the repository's exclusive-publication/PID/token approach, with complete owner
publication by exclusive hard link, random generation token, idea/run identity and host/boot identity.
It deliberately does not copy the old read-dead/unlink loop. A permanent O_EXCL claim for each stale
generation elects one reaper; a delayed second reader cannot delete the new winner. Partial owners,
unknown/different host or boot, liveness uncertainty and interrupted reclamation fail closed.

Locks remain on the deck filesystem under ignored `.parley-runtime/membership/<idea-hash>/`.
`git check-ignore -v` returned `.gitignore:23:/.parley-runtime/` for both lifetime and projection
leases (`C/environment-protocol-final.log`). Projection waiting is bounded to two seconds, separate
from provider retry/polling. Off, legacy and kickoff-only ordinary driving has no singleton.

Executed on BOTH volumes: cross-process lifetime refusal; same-PID competing run refusal; nested
context checks; released-context rejection; cross-process projection refusal/wait/release;
old-release identity protection; partial/foreign/unknown owners; concurrent and delayed stale reapers.
Locators: `internal/membership/membership_test.go:221`, `:245`,
`internal/membership/fixup_test.go:18`, `:123`, `:135`, and
`internal/pidlease/lease_test.go:71`, `:126`.

`TestLeaseProcessSyntheticIdentity` (`lease_test.go:67`) uses actual separate processes and actual
kill/wait/liveness, but explicitly injected **synthetic host/boot identity**. Stale takeover passes on
both AppleVirtIOFS and /tmp and under `-race`. This is algorithm/volume evidence, not proof of native
host identity acquisition.

**Native stale takeover remains an environmental verification gap.** The unweakened
`TestLeaseProcess` (`lease_test.go:15`) fails on both volumes and under race because this sandbox
returns an empty boot identity. Exact shared failure (`C/shared-volume-final.log:304`):

```text
stale takeover lease held: idea/child-run pid=96398 host="Tomass-Virtual-Machine.local" boot="" at /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude/.parley-runtime/quota-implementation/fixup-1/checks/shared-temp/TestLeaseProcess11622563/001/driver.lease; only a proven dead local owner may be reclaimed
```

The environmental commands were executed and preserved verbatim:

```text
$ /usr/sbin/sysctl -n kern.boottime
sysctl: sysctl fmt -1 1024 1: Operation not permitted
exit=1
$ /usr/sbin/sysctl -n kern.uuid
sysctl: sysctl fmt -1 1024 1: Operation not permitted
exit=1
```

No unknown-boot relaxation, skip or altered native expectation was used. Organizer host execution is
still needed for native stale takeover. Windows compilation passed; Windows runtime/liveness and
shared-volume behavior were not executed. Windows automatic stale reclamation remains conservative.

## G4 / R4 — complete allowlisted SDK framing, honest native-capture limit

Sources: `internal/telemetry/quota_framing.go:32`, `internal/telemetry/quota_zcode.go:18`, `:33`.
The parser consumes complete SDK error objects and a terminal failure. It validates stack shape,
class/symbol identity, fields, containers, headers, matching response/header message, data/cause,
all retries, attempt count/reason and exact last-error correspondence. It rejects arbitrary free
prose, quoted wrappers, mixed JS errors, unknown/duplicate fields, extra records, malformed/nested
inconsistencies, inspect truncation and incomplete captures. Quoted `'undefined'` is not undefined.
The exhaustion message itself is allowlisted so a warning about a limit is not exhaustion.
Per-record observed receipt timing, failure/artifact/success precedence, reset consistency and the
60-minute threshold remain enforced.

I investigated the installed SDK source, including actual APICallError/RetryError definitions,
responseHeaders, request fields, nested data/errors, symbols and terminal trace. Source:
`/opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs`, SHA-256
`3e3433d90fa502e5d02498dfde6c2090df898331359bcfe5f3dbc9a1d00b685f`.
The offline `internal/telemetry/testdata/quota/sdk-framing-harness.cjs` extracts those error classes
into an isolated VM without bootstrapping the CLI or accessing configuration/network. Commands:

```sh
node internal/telemetry/testdata/quota/sdk-framing-harness.cjs api
node internal/telemetry/testdata/quota/sdk-framing-harness.cjs retry
```

`C/sdk-source-fixed.log` records the successful harness execution; the initial failed harness
attempt is retained separately. The two resulting `.stderr` files are explicitly **OFFLINE
SOURCE-DERIVED SDK fixtures**, with synthetic request/stack fields and the retained owner's real
positive response body/reset. Their hashes are in `C/environment-protocol-final.log`; provenance
and limitations are in `internal/telemetry/testdata/quota/README.md`.

The original `/tmp/probe-20261002/zcode.err` (24,833 bytes) is gone. It was NOT recovered. The
four-line historical excerpt is incomplete and is rejected alone. The SECOND retained excerpt,
`source-context/provenance-review/codex-1-retained-native-tail.md`, starts inside responseHeaders;
it confirms later retry-after 176890, reset_at `2026-10-04T22:14:57.003Z`, actual symbols/data and
terminal trace. It is the same incident, not a second complete capture. No excluded provider was
invoked to replace it. **Full-native AC2 capture evidence remains missing for independent review
and owner disposition.** This is not silently waived or reclassified as native coverage.

Executed: `quota_framing_test.go:10` includes every reviewer wrong-positive family in complete
framing before/between/after records; `:68` nested retry consistency; `:95` differing observed receipt
clocks (including the second tail's decreasing reset); `:120` complete-framing semantics, 3599/3600
seconds, generic/negated/warning messages, missing/invalid resets and lastError mutation. The retained
positive body/reset still classifies at `2026-10-04T22:14:57.003Z` in complete source-derived framing.
All telemetry tests passed in the full affected-package run.

The reviewer's ORIGINAL unmodified `.parley-runtime/claude1-r3/zadv/main.go` was executed with
`GOCACHE=/tmp/quota-auto-exclude-go-cache go run .parley-runtime/claude1-r3/zadv/main.go` (exit 0,
`C/reviewer-zadv-final.log`). Every row is ineligible, including its historical incomplete positive
and boundary excerpts. The complete-framing positive/boundary tests above distinguish this from a
recognizer that merely rejects everything. Accepted residual subagent-source ambiguity is unchanged;
other adapters remain diagnostic-only.

## G5 — active context through signoff handoff

Sources: `internal/app/consensus_request_signoffs.go:584`, `:613`, interactive handoff path, and
`internal/runner/handoff.go`. Manual and interactive handoffs receive the active context through
`HandoffOptions.Context` and carry the canonical artifact target.

`internal/app/quota_fixup_test.go:134` creates recorded mid-idea history revision 1 before exercising
both modes. Both handoffs succeed with the owner context; a competing run and an unleased handoff
are refused. Executed on both volumes and under race, with exact results in their final logs.

## G6 / R1 — retained-obligation authority and return/withdrawal

Sources: `internal/protocol/quota_retained.go:129`, `:243`,
`internal/consensus/consensus.go:1181`, and committed authority validation above. A ruling requires
attributable canonical evidence from a known identity, an actual committed owner answer and its
verbatim quote under User direction. Path/idea/author mismatches and duplicate disposition fields
fail closed. Human identity/truth remains outside structural/hash verification.

Author withdrawal requires the original author, an author-owned later canonical round, exact
immutable re-inclusion batch/revision, current membership and the named obligation. Ordinary
owner-confirmed off-mode returns qualify through their validated immutable confirmation snapshot.
The old filed bytes remain untouched. Withdrawal is not a fresh signoff and cannot erase a later
new veto, including an identical veto filed again on the same day or signoffs containing blank lines.

Executed: `internal/consensus/quota_test.go:168`, `:224`, `:264`, plus
`TestQuotaRetainedVetoOnlyOwnerRulingCanDispose` and committed-authority negatives in
`internal/quota/authority_test.go:17`. Re-inclusion → bound withdrawal → fresh signoff reaches the
expected consensus state in tests; absent/fabricated binding remains blocked. Permitted inbox
deletion preserves the committed ruling/re-inclusion authority. Full consensus package passed.

## G7 — checked immutable-history receipt recovery

Sources: `internal/protocol/quota.go` (bad receipts mean pending),
`internal/membership/membership.go:79` (no receipt-existence shortcut), and
`internal/app/quota_revision.go:18` (`quota recover`). Documentation:
`docs/quota-membership.md`, Checked receipt recovery and locking.

`internal/membership/fixup_test.go:66` executes empty and contradictory receipts, confirms actions
are blocked, supplies a contradictory manifest and verifies recovery refuses without marking the
receipt applied, restores the authoritative projection, recovers twice and verifies one notice and
one terminal evaluation. Immutable history is unchanged. `internal/app/quota_fixup_test.go:89`
executes the actual CLI recovery twice and prints exact outputs. The tests passed locally, on the
shared volume, and under race; existing persistence-fault tests also passed in the full membership
package run. No legacy budget migration or original-run synthesis was added.

## G8 — canonical launch and complete driving lifetime

Sources: `internal/runner/quota_target.go:14`, `internal/runner/telemetry.go`,
`internal/app/agents_exec.go`, `internal/app/app.go` (driver join before lease release),
`internal/app/quota_pipeline.go:16`, `internal/app/pipeline_cmd.go`, and
`consensus.FinalizeContext`.

The TUI lead is corrected by cancelling and joining the driver goroutine before releasing the
outer idea lease. `internal/app/quota_fixup_test.go:172` executes the actual run path with local
round dispatch and deterministic UI/driver operation seams: UI returns, driver observes cancellation
but blocks in cleanup, competing acquisition is refused, then cleanup completes and acquisition
succeeds. This reproduced the lifetime boundary rather than relying on comments.

The pipeline lead is corrected by one idea-bound run/lease across the block, including between-stage
intervals. `quota_fixup_test.go:370` executes the actual implementation block driver with local stage
callbacks: implementation → review-1 → consensus-1 → fixup → review-2 → consensus-2. Every callback
checks the same live run, refuses a competitor, and the post-return acquisition succeeds. `:249`
checks the intervening stage ownership, including signoff; finalization accepts the same scoped
context. These are stubbed stage-dispatch lifecycle tests, not provider executions. They passed on
both volumes and under race.

A canonical artifact target under this deck infers the idea when `--idea` is omitted. Conflicting
idea/preflight targets and symlink ambiguity are rejected. Arbitrary unbound work stays outside this
idea. `quota_fixup_test.go:332` runs the actual `agents exec` CLI: excluded canonical launch without
`--idea` refuses before the marker process runs; generic unbound exec succeeds with the local stub.
`internal/runner/quota_test.go:188` covers the lower-level target/identity gates. Exact outputs are
preserved in `C/focused-final.log` and both volume logs.

## G9 — named gates and matching provider diagnostics

Sources: `internal/quota/record.go` notice, `internal/telemetry/provider.go`,
`internal/app/preflight_liveness.go`, `internal/runner/failclass.go`.
Notices name reviewer count, LE-7/LE-11, independent goal-checker eligibility,
require_model_diversity, strict_gate and retained veto/DISPUTED/finding gates with applicability.
Preflight and runner use the same bare-503 grammar: a whole line or explicit HTTP/status/error
framing, not any occurrence of the number in sizes/source positions.

Executed: `internal/quota/authority_test.go:66`,
`internal/app/quota_fixup_test.go:284`, `internal/runner/quota_test.go:212`. Exact bare 503,
whitespace, 503 Service Unavailable, HTTP/2 503, statusCode and API error forms pass their provider
gate. `size=503 bytes`, `source.js:503:1`, `processed 503 records`, `HTTP/2 200 (503 bytes)`,
`line 503` and `build 1503` remain noise. Focused checks passed.

## Final executed commands and outcomes

Exact expanded argument arrays, environment overrides, elapsed times and exit codes are in
`C/final-checks-results.json`; every `.log` begins with its command and environment.
Reproducible orchestration is `C/final-checks.py`. All final commands used
`GOCACHE=/tmp/quota-auto-exclude-go-cache`. Test discovery is guarded by `C/no-provider-bin/`;
test-specific fixtures provide their own stubs. No `C/guard-denied.log` was produced.

| Log in C/ | Command (in CLI worktree) | Observed result |
| --- | --- | --- |
| focused-final.log | `go test ./internal/app ./internal/runner ./internal/driver ./internal/consensus ./internal/membership ./internal/quota ./internal/telemetry ./internal/fsutil ./internal/runcontrol ./internal/runstate ./internal/runmanifest ./internal/protocol -run 'TestQuota\|TestSyncFile\|TestEmbeddedDefaultMatchesLiveDeck\|TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail' -count=1 -v` | Exit 0; all selected tests passed. |
| shared-volume-final.log | `TMPDIR="$PWD/.parley-runtime/quota-implementation/fixup-1/checks/shared-temp" go test ./internal/app ./internal/membership ./internal/pidlease ./internal/fsutil -run 'TestQuotaFixup\|TestQuotaLifetimeLockDifferentRunAndOffScope\|TestQuotaCrossProcessLease\|TestLease\|TestSyncFile' -count=1 -v` | Exit 1 only for native `TestLeaseProcess` empty boot identity. All creation, owner/recovery, handoff, lifetime/projection, synthetic stale, race-of-reapers and fsutil cases in this run passed. |
| local-control-final.log | Same command with `TMPDIR=/tmp/quota-auto-exclude-local-control` | Same native boot-identity failure; other selected checks passed. |
| race-final.log | `go test -race ./internal/pidlease ./internal/membership ./internal/app -run 'TestLease\|TestQuotaFixup\|TestQuotaCrossProcessLease\|TestQuotaLifetimeLock' -count=1 -v` | Exit 1 for the same native boot test; membership and app passed, synthetic stale/reaper cases passed, no race-detector report. Not a blanket race-suite pass. |
| affected-full-final.log | `go test ./internal/quota ./internal/telemetry ./internal/fsutil ./internal/membership ./internal/consensus ./internal/runcontrol ./internal/runstate ./internal/runmanifest ./internal/store ./internal/protocol ./internal/app ./internal/runner ./internal/driver -count=1` | Exit 1. First ten packages passed in full; app/runner/driver failed with sandbox budget-lock/boot denials and downstream failures. Full pass is NOT claimed. |
| build-final.log | `go build ./...` | Exit 0. |
| vet-final.log | `go vet ./internal/...` | Exit 0. |
| windows-build-final.log | `GOOS=windows GOARCH=amd64 go build ./...` | Exit 0; compilation only, no Windows execution. |
| format-diff-final.log | `gofmt -l` on `C/changed-go-files.txt`; `git diff --check` | Both exit 0; no listed formatting defects. Only changed Go files were formatted. |
| reviewer-zadv-final.log | `go run .parley-runtime/claude1-r3/zadv/main.go` | Exit 0; all incomplete/adversarial rows ineligible; native excerpt acceptance gap remains explicit. |

The actual regex arguments contain normal `|` alternation (escaped above only for Markdown table
rendering); the JSON and log headers are the exact shell/argv authority.

The full affected run preserved 167 top-level failing test groups. Many explicitly report the
host-wide budget-lock denial; others report downstream missing policy/evidence, zero starts or
boot verification failures. I did not silently classify every secondary assertion as independently
proven environmental, alter those tests, or repair the excluded legacy-budget/driver-gap scope.
Representative verbatim failures in `C/affected-full-final.log`:

```text
persistent driver budget: open /Users/tomasfecko/Library/Caches/parley/budget-locks/.lock-identity-723126391: operation not permitted
process verification failed (no recorded boot id); not killed
```

Organizer host action remains: execute the same full affected command and the unweakened native
lease/volume/race commands outside this child sandbox. Do not infer these results from the focused
pass. Native stale behavior on both filesystems and the full affected host run remain unverified here.

## Scope, deviations and limits left visible

- No project commit, merge, release, owner close, worktree creation/declaration/pruning, roster or
  credential mutation, skill/protocol-copy edit, reviewer artifact or signoff was made by this process.
  Synthetic Git commits exist only inside isolated temporary authority-test repositories because
  committed-evidence validation must be exercised. The real project's commits are left to the organizer.
- The only canonical idea artifact I wrote is this producer evidence file. Organizer notes/usage were
  concurrently changed by the organizer; I did not edit them. The two pre-existing untracked run
  directories `20261003T103229.878442000Z` and `20261003T185354.660080000Z` were left alone.
- No task/inference/provider probe was dispatched to a real participant. Earlier development test
  discovery did invoke installed CLIs for read-only version inventory (visible in `C/app-e2e-dev.log`).
  Final acceptance then isolated all discovery to missing/local stub commands and a version-only guard;
  no non-fixture guard invocation occurred. No browser, alternate participant, model or provider was
  launched for implementation/review. SDK fixture generation is offline source extraction.
- R1 committed quote/path/blob/digest survives permitted inbox cleanup. R2 existing off-mode recorded
  confirmations work without a mandatory new CLI. R3 locks stay ignored and deck-local; unknown/foreign
  identity never becomes stale from local PID absence. R4 is an allowlist and the missing full-native
  AC2 capture remains an acceptance limitation. None is treated as an optional suggestion.
- Native capture acceptance, native stale takeover under an available boot identity, the unsuccessful
  full affected host checks, and Windows runtime verification remain explicit limits. Hash/attribution
  checks do not prove the human's identity, truth, prior reading or freedom from cooperative forgery.
- No new owner policy, legacy budget migration or driver-gap repair was inferred. This delivery is a
  checked producer patch for independent Phase-8 re-review; it does not finish the two-stage process.

Initial work began approximately 2026-10-04 03:07 UTC. Final checks completed around 04:20 UTC;
this evidence was assembled after the 04:24 UTC clock check, within the recorded 7,200-second ceiling.
