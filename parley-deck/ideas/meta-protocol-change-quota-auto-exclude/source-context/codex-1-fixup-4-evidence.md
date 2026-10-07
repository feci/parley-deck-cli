---
agent: codex-1
idea: meta-protocol-change-quota-auto-exclude
phase: fix-up-cycle-4
role: implementer-producer
status: implementation-evidence
base-commit: 267086c402e4880600487f824e2dac7a529d181d
date: 2026-10-07
---

## Scope and authority

This is producer evidence from the activated cycle-4 implementation process. It is not independent review,
a code acceptance, an AC verdict, a signoff or close authority. The parent owns G20, IMPLEMENTATION,
orchestration and final review artifacts. This process edited only the G18/G18b/G19 Go files listed below
and this report; scratch, fixture copies, caches and logs are private under
`.parley-runtime/quota-implementation/fixup-4/producer/` (abbreviated **P** below), with the requested local
filesystem control under `/tmp/quota-fixup4-producer/`.

Read in full: signed `review/consensus.md` (303 lines), including claude-1's complete signoff and every
reservation; `review/round-06/claude-1.md` (558 lines); the owner `finish-now`, `round06-answer` and
`round05-answer` notes. Read FINAL's relevant sections 5–11 (lines 124–289), AC1–AC21 and
Idempotence & recovery (670–834). Read this idea's kickoff status/implementation handoff and the current
protocol-changelog entry. The signed plan plus activation accepts Reservation 1 and selects Reservation 2's
smaller G19 route. G19 therefore adds no Kickoff field and has no legacy-marker limitation. The deliberate
change to earlier G15 is that a genuine new identity can join during round 1 by a plain participants edit.

No additional agent was spawned, delegated to or invoked. No provider, zcode or native-exhaustion capture
ran. No real repository commit, merge, release, worktree/declaration/prune, roster/settings/credentials,
D6 accounting, normative COOPERATION, version, docs or skill edit was performed by this process. Synthetic
Git commits exist only inside isolated test fixtures for immutable-authority tests. Parent-owned docs and
orchestration edits appeared concurrently in the shared worktree and were left untouched. The two untracked
historical run directories present at entry were left untouched.

I read the Parley Deck and OpenViking-memory skill instructions. A scoped read-only OpenViking find returned
older round-04/05/06 checkpoints; the newer signed local plan and activation supersede their old stop rules.
No shared-memory write was made within this narrow producer scope. No graphify graph existed; no graph or
helper-agent pipeline was created.

## Full protocol context attestation

- `context_mode`: `full`
- `source_sha256`: `73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`
- `packet_sha256`: `73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`
- `fallback_reason`: `null`
- Body: `.parley-runtime/protocol-packets/full-phase8-deliberation-73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e.md`
- Active transport in the body: `github-pr`; no transport mutation in this subprocess.

Actual coverage: all 1,501 lines, read in bounded chunks 1–250, 251–500, 501–750, 751–1000, 1001–1250,
1251–1501. A combined tool-output cap clipped the end of the 251–500 call; lines 475–500 were re-read
separately. The adjacent review-file header was also re-read separately. The long lines 59, 452 and 1002
were displayed through their ends. No optimized or shadow packet was used. `shasum -a 256` of both the
unchanged body and live `parley-deck/COOPERATION.md` returned the hash above.

## Implementation

**G18.** `protocol.InspectQuota` has no notice inspection or notice-path access. The quota-layer exact-byte
notice validator was removed. Before an ordinary applied receipt, the publication step preserves regular
live and archived copies regardless of bytes and exclusively creates only if absent. It refuses symlink or
nonregular destinations and non-directory deck/inbox/archive parents. Notice-only failures emit stderr
`publication diagnostic (non-blocking; delivery unconfirmed)` and do not prevent membership reconciliation,
terminal evaluation or the checked applied receipt. After the receipt, ordinary notice publication does
not inspect notice files or paths at all. Immutable history, writer, projection, terminal-evaluation and
receipt checks remain in place.

A completed publication **attempt** is not proof that a message was delivered, read or acknowledged. A
refused unsafe destination may leave no notice, with a diagnostic and a valid reconciliation receipt;
there is no new retry queue. Regular owner files never need matching producer bytes or a fresh sync to
count as an existing copy. A crash before the receipt still allows checked replay; deletion of an
unreceipted notice can cause one benign re-publication. No historical records are migrated or rewritten.

The historical manual-authority clarification remains a separate, non-authoritative notice with its own
checked receipt. A retained legacy prefix, including appended owner text, is only a publication hint.
Existing clarification copies belong to the owner. Its valid receipt ends clarification inspection and
publication; an unsafe clarification receipt still returns an integrity error. Ordinary publication-only
failures remain diagnostics. The original inaccurate historical notice bytes are preserved.

**G18b.** Working-copy comparisons moved out of `authorityBytes`/`ValidateAuthority` into `BindAuthority`.
An extant superseded live or archived answer cannot be newly bound. After binding, history and replay
validate the Git commit, blob identity, SHA-256, attribution and verbatim quote, without comparing mutable
inbox bytes. Those immutable checks were not relaxed.

**G19.** `ManualRoundOneReturn` now checks only the frozen prompt snapshot's `status: round-01`; its unused
identity argument and marker condition were removed. Manual import and immutable-history replay use the
same predicate. The existing policy-off boundary, known-identity path and after-round-1 catch-up path
remain. No schema field, historical membership, owner authority or legacy limitation was added.

## File coverage and tests

Production files read and edited in their affected context:

| Files | Change |
|---|---|
| `internal/protocol/quota.go` | Delete all notice inspection from membership reads. |
| `internal/membership/notice.go` | Owner-owned copies; receipt short-circuit; safe publication-only checks and diagnostics; separate clarification receipt. |
| `internal/membership/membership.go` | Move the notice fault seam into publication so it cannot gate receipts. |
| `internal/quota/notice.go`, `internal/quota/receipt.go` | Remove exact-byte notice validation; clarify attempt versus delivery. |
| `internal/quota/authority.go` | Bind-time working-copy comparison only. |
| `internal/quota/revision_snapshot.go`, `internal/quota/revision.go`, `internal/protocol/quota_manual.go` | One immutable round-1 predicate at import and replay. |

New tests (read fully after writing):

- `internal/membership/cycle4_test.go`: before/after receipt × live/archive/both copies × appended answer or
  arbitrary replacement, preserved bytes and history, repeated recovery, exactly one terminal evaluation;
  symlink/nonregular notices, inbox/archive symlink parents and injected publication failure as diagnostics;
  receipt symlink/directory, changed receipt and corrupted history gates; manual clarification annotation,
  arbitrary replacement, archive, unsafe destination and independent receipt integrity.
- `internal/app/quota_cycle4_test.go`: real CLI status and signoff after notice edits and unsafe destinations,
  before/after receipt; actual `agents exec` command implementation dispatching a local stub with an authentic
  driver lifetime lease; real `quota revise`, bound live/archived owner-answer append, status/signoff/Before
  and missing-receipt replay; deleting the immutable blob still gates status/signoff/Before. Real `parley run`
  kickoff uses local shell agents and a controlled readiness-probe seam to exclude one identity; kept/removed
  markers and a genuinely new round-1 identity sign through the real CLI. Kickoff bytes remain unchanged,
  manual revisions have no owner authority/catch-up fiction, and later live round advancement does not
  invalidate the recorded round-1 snapshot. Policy-on and final/closed plain-edit boundaries remain gated.
- `internal/quota/authority_cycle4_test.go`: superseded live/archive copies reject new BindAuthority while
  bound authority remains valid; missing commit/blob, changed blob/digest and unquoted text still fail.
- `internal/quota/revision_cycle4_test.go`: immutable round-01 snapshot admits a new identity without a marker;
  round-02, consensus, final, closed and absent status do not bypass catch-up, with or without a marker.

Earlier tests were adapted only for changed signed semantics:

- `internal/app/quota_cycle3_test.go`, `internal/consensus/catchup_test.go` and the off-mode catch-up test in
  `internal/membership/revision_test.go` now set round-02 for later-round catch-up/decline assertions. Their
  gates, incomplete retries, atomic import and known-excluded assertions remain.
- `internal/membership/cycle2_test.go`, `cycle3_test.go`, `membership_test.go` no longer expect owner edits
  or notice-only faults to gate; the removed notice-corruption matrix is replaced by the broader cycle-4
  owner-copy/unsafe-path/real-integrity matrices. Archive/delete/replay assertions remain.
- `internal/quota/authority_test.go` no longer treats an archived owner edit as invalid immutable evidence;
  new bind-time coverage replaces it. The self-authored negative in `membership/revision_test.go` now
  supplies a committed self-authored blob with matching hashes, so attribution rather than a mutable-file
  comparison rejects it.

Additional full source reading: quota history/revision/snapshot/authority/receipt code, membership
reconciliation/revision code, protocol quota/manual code. Targeted reading: app command dispatch and
preflight, participant artifact validation, test fixtures, runner/driver lease call sites and retained-veto
tests. This is focused producer coverage, not the full product-diff coverage owed by independent review.

## Commands and observed outputs

All Go runs used P's private `cache`, inspected/copied provider guards, and `GOPROXY=off`; baseline probes
also use a private `PARLEY_HOME`. `P/no-provider-bin` is a copy of the requested
`fixup-3/host/no-provider-bin`, with only its deny-log destination changed to P. No `P/guard-denied.log`
was produced. The guard source was read before use. Test agents are local stubs; fixture authority is
explicitly synthetic.

Exact expanded commands, working directories, durations and exits are in `P/results.json`,
`P/final-results.json`, `P/driver-local-result.json`, `P/format-result.json` and each log's command header.
The reproducible runners are `P/checks.py` and `P/final-checks.py`. Summary:

| Command | Observed exit / evidence |
|---|---|
| `go test ./internal/quota ./internal/protocol ./internal/membership ./internal/consensus ./internal/app -run TestQuota -count=1 -timeout 5m -v` | **1 on both volumes**. 70 top-level test functions reported PASS, one reported FAIL: unchanged `TestQuotaCycle2NativeCrashWriter`; exact limit below. Logs `focused-shared.log`, `focused-local.log`. |
| `go test ./internal/quota ./internal/membership ./internal/consensus ./internal/app -run TestQuotaCycle4 -count=1 -timeout 5m -v` | **0 on both volumes** after the last test addition. 11 top-level tests, 81 PASS lines including subtests, zero FAIL lines per run. Logs `final-cycle4-shared.log`, `final-cycle4-local.log`. |
| `go test ./internal/runner ./internal/driver -run TestQuota -count=1 -timeout 5m -v` | **1 with shared TMPDIR**: runner exits successfully; driver binary reports `signal: killed` before any test starts. |
| `go test ./internal/driver -run TestQuota -count=1 -timeout 5m -v` | **0 with local TMPDIR**; `TestQuotaDriverReconcilesPendingBeforeRebindingEveryConsumer` executed. |
| `go build -o P/bin/current-<probe> ./P/probes-current/<probe>` for notice6, notice7, d14, vprobe, plain-edits | **0**, all five builds. See exact paths in results.json. |
| `go build -o P/bin/baseline-<probe> ./producer-probes/<probe>` in P/oracle-27e42b8 for d14, vprobe, plain-edits | **0**, all three builds from the authentic prechange source. |
| `gofmt -w` affected Go sources/tests; `gofmt -l` all 20 producer Go files; `git diff --check -- <those files>` | **0**, no unformatted paths or whitespace errors. Exact file list and hashes in producer-source-files.json. |

Shared TMPDIR: `P/shared-temp`. Local TMPDIR: `/tmp/quota-fixup4-producer`.
The narrower final cycle-4 runs supplement, and do not erase or reclassify, the broader failed runs.

The earlier development runs are retained too: `first-focused-shared.log` exited 1 for native crash plus
an old self-authored test that changed only the live answer (updated to committed negative evidence);
`second-focused-shared.log` exited 1 because the new dispatch probe lacked the existing required lifetime
lease. Its setup was corrected to carry the real lease and run id, without changing product lease logic.
`third-cli-shared.log` then exited 0. These are development failures, not suppressed test evidence.

## Copied reviewer probes and actual baseline

The originals `.parley-runtime/claude1-r6/{notice7,notice6,vprobe,d14}` remain read-only and unchanged.
Copies were made into `P/probes-original` before adaptation or execution. D14's runner points to the
unchanged prior producer `fixup-3/producer/diff-current` main/setup; those Go sources were also copied into
P. Original-source hashes are recorded in `P/reviewer-probe-original-hashes.json`. The existing
`fixup-4/precheck` was neither altered nor relabeled: it reproduces the old defect, not a passing check.

The actual archive `/tmp/claude1-r5-prechange-27e42b8` was the oracle. Before copying it into private scratch,
all 491 tracked internal/cmd/module files were compared byte-for-byte with `git cat-file --batch` for
`27e42b8`; every comparison matched. `P/oracle-manifest.json` records each SHA-256. No original archive
file was modified, and no baseline source was approximated. D14 and Vprobe use identical main files on
both trees, with only their established setup adapter changed for the two APIs. The new `plain-edits`
probe derives its helpers from the copied Vprobe and also uses one identical main on both trees.

All 20 probe executions (8 notice, 8 original D14/Vprobe, 4 plain-edits) exit 0. An original probe's exit
alone is not an assertion: the printed command outcomes were inspected. Notice6/7 on **both volumes**:

```text
Before #1 err=<nil>
Before #2 err=<nil>
status exit=0 integrity-gate=false pending=false corrupt=false
survivor a consensus signoff exit=0 stderr=""   # notice7
applied-receipt=true err=<nil>
round.completed events: 1
```

Archived notices stay at live=0/archive=1; deleted applied notices stay at 0/0; annotated live notices stay
at 1/0. This is the behavior the original precheck failed to produce. Full outputs are
`P/logs/probe-current-notice{6,7}-*.log`.

The executable plain-edit differential yields this table on **both volumes**:

| Requested join before late artifact | 27e42b8 joiner signoff | Current joiner signoff | Current after CLI catch-up |
|---|---|---|---|
| round-01, excluded marker kept | 0 | 0 | not required |
| round-01, excluded marker removed | 0 | 0 | not required |
| round-01, genuinely new identity | 0 | 0 | not required |
| round-02, excluded marker kept | 0 | 1, pending catch-up | 0 |
| round-02, excluded marker removed | 0 | 1, pending catch-up | 0 |
| round-02, genuinely new identity | 0 | 1, pending catch-up | 0 |

Current history checks also assert no kickoff membership for the new/returned id, then exactly one manual
revision and no owner authority. Full logs: `P/logs/probe-{baseline,current}-plain-edits-{shared,local}.log`.
The original D14/Vprobe logs remain available, including the deliberate round-1 change in their formerly
pending join/decline scenarios. Later-round decline assertions remain in the committed Go tests, now set
at round-02. Retained-veto, known-excluded and withdrawal tests executed in the broader focused run.

## Failures, host limits and checks still owed

Native crash evidence could not be established in this sandbox. Exact command/output:

```text
$ sysctl -n kern.boottime
sysctl: sysctl fmt -1 1024 1: Operation not permitted
exit=1
```

The unchanged native-crash test reported on the shared volume:

```text
cycle2_test.go:213: unsettled writer for idea native-crash: f8e87507-c78b-47ee-bc6d-39267c2504ed: only proven-dead local supervisor, writer and process group may be settled; stop/verify the original host's supervisor and writer, then run parley quota recover; legacy or unknown identity requires restoration of authentic terminal evidence, never a fabricated settlement
--- FAIL: TestQuotaCycle2NativeCrashWriter (0.30s)
```

The local run reports the same refusal with invocation `f3939bfe-58d1-4ae3-afd6-ff107afae45b` (0.23s).
Logs: `native-boot-limit.log`, `focused-shared.log`, `focused-local.log`. No native-lease, boot, writer-death
or crash-settlement assertion was weakened. The shared-TMPDIR driver launch additionally printed
`signal: killed` and `FAIL parley-deck-cli/internal/driver 0.003s`; no cause is established by that output.
The local-TMPDIR repetition exits 0, but does not prove the shared execution issue resolved.

Parent still owes full HOST Go suite/build/vet/race, shared/local host checks including native crash,
Windows cross-build, packet/drift/copy checks and the skill suite after the code freeze. None is claimed
from these focused runs. Windows runtime and container/PID namespaces remain unexecuted. AC2's native
positive evidence remains **NOT MET / owner-waived**, not fixed or independently accepted here.

## Deviations and handoff

No scope/authority change was needed. Reservation 1 and the selected smaller Reservation 2 route are the
signed implementation, not deviations. There is no new persisted schema, queue, rejoin timer or retry
mechanism. The new tests distinguish the existing lease requirement from notice handling; a standalone
policy-on `agents exec` without a driver lease retains its prior refusal.

The parent should use the G19 table and attempt-versus-delivery distinction for G20 wording. Producer code
hashes are in `P/producer-source-files.json`; they were unchanged through the final checks. This report
requests no self-acceptance and supplies no independent PASS. The parent may freeze these Go changes for
its host validation and the separate reviewer.
