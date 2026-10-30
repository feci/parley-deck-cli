---
agent: zcode-1
idea: meta-protocol-change-credit-exhaustion-failover
review-round: 3
date: 2026-10-30
reviewed-commit: b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666
responding-to: [kimi-1/review/round-02]
---

## Summary

I performed a fresh full-scope review of the complete CLI diff `c2d2431..b7dce40` and the skill diff `f557c5e..406b6976`, covering both fix-up cycles. F7, F8, and F9 are applied as signed. I executed my own refutations at the fixed source: a 100-run concurrent-budget repetition, the new publication fault tests, the focused quota, membership, protocol, telemetry, runner, readiness, and packet suites, a `gofmt` scan over every changed Go file, and hash checks on all protocol copies. The round-02 CRITICAL concurrent-admission race no longer reproduces in 100 consecutive runs. I found no new CRITICAL, MAJOR, or MINOR product defect. One NIT follows on evidence receipts. Two gates stay open and I do not close them: the organizer's required full Go command is still running, and AC-15 release evidence does not exist yet by design.

Provenance: every verdict below carries its tag. My executed commands and quoted outputs are `PRIMARY`. Organizer evidence I only read is named as such. I made no provider call and induced no exhaustion. Synthetic fixtures only.

Protocol attestation: context_mode full; source_sha256 and packet_sha256 `b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709` (verified by `shasum -a 256 parley-deck/COOPERATION.md`); fallback_reason none. I read the complete live protocol, including §15 and Appendix B, before I wrote this file.

Scope facts, all `PRIMARY`: `git rev-parse HEAD` prints `511e2eb`; `git diff --stat b7dce40..HEAD -- ':!parley-deck'` is empty and `git status --porcelain -- ':!parley-deck'` is empty, so the reviewed product tree equals `b7dce40`. The fix-up-2 product diff `31c039c..b7dce40` changes six files only: `docs/credit-exhaustion.md`, `internal/app/app.go`, `internal/driver/impl.go`, `internal/quota/history.go`, `internal/telemetry/record.go`, and the new `internal/telemetry/record_publication_test.go`.

## Position changes since prior review round

1. My round-02 supplemental CRITICAL finding (common exhaustion admission races with concurrent record creation) is resolved at `b7dce40`. I withdraw it. Evidence: my own 100-run repetition passed, and the writer now publishes complete bytes exclusively. See Refutation attempts, AC-05.
2. My round-02 NIT (capture status deletes expired evidence) is resolved by F8 through the documentation alternative I proposed. I verified the text in the CLI help, the CLI guide, and the packaged skill reference.
3. My round-02 AC-14 position changes from "evidence negative" to "concurrency portion restored; full-suite portion pending". The organizer's `go test ./... -count=1 -timeout 45m` process is still alive and its log `fixup02-full-go.log` holds three package lines. Pending is never passed. I record AC-14 as still open on that full-suite result, not as failed.
4. No other position changes. My round-02 durability addendum's required corrections (directory barriers, store barrier) are what F7 implemented; my assessment of them stands, now with the executed fault tests that were missing then.

## Responses to other reviewers

### @kimi-1 — review/round-02

Your CRITICAL finding was correct at `31c039c`, and your durability correction went beyond my overlay assessment. You were right that the proposal I tested covered visibility only. F7 implements both barriers you named, and I verified them in source and through the new fault tests. Your `quota.SyncPath(i.Dir)` proposal is the exact call now at `internal/telemetry/record.go:450`. I concur with your finding and its fix; it is resolved at the current source.

Your MINOR on capture-status documentation is resolved. `docs/credit-exhaustion.md:102-103` now reads: `It also deletes expired files, removes oldest files for the total cap, and updates `cleanup.json`.` The same text is in `internal/app/app.go:259-260` and in the packaged reference at `skills/parley-deck/references/EXHAUSTION_AND_HANDOVER.md:100-103`. All three `PRIMARY` by direct read.

Your NIT on `gofmt` is resolved. `gofmt -l` over all 74 Go files changed in `c2d2431..b7dce40` printed nothing. The import now sits with the project imports in `internal/driver/impl.go`.

Your round-02 open questions, answered at this source: the implementer did add exclusive atomic publication and both directory barriers (verified in source and tests); the concurrent repetition and race checks pass (my own run below), the full Go command is still running; the fresh goal check and channel verification remain open post-consensus processes under D12 and AC-15.

### Signed fix-plan dispositions (F7–F9, review/consensus.md)

I concur with all three fixes as applied. Details under Refutation attempts. The F1–F6 dispositions from the archived `consensus-cycle-01.md` remain satisfied across the full current diff; I re-verified the F1 protection sites and reran the readiness tests this round.

### Windows disposition

I concur with the documented platform scope, as in round-02. `docs/credit-exhaustion.md:91-92` states: `Windows cannot supply these POSIX modes through the current storage check. Capture refuses there and writes no provider data.` and `CLI Windows support remains experimental.` No check was weakened for any platform: the fix-up-2 diff returns every persistence error on every platform, keeps the `Perm()`-based store checks, and adds no `GOOS`-gated softening. Platform failures establish no Windows correctness. CLI WinGet stays held under D12.

### Driver refusal and native fallback

`fixup02-driver.log` records the phase-6 refusal: `continue failed: idea "meta-protocol-change-credit-exhaustion-failover" has no runs yet`. The owner answer authorizes the native fallback. I state its boundary as before: it supplies no driver accounting and no automatic membership authority. This review is a native-process artifact under that authorization.

## Refutation attempts

I assumed each criterion was broken and tried to break it. Commands ran at HEAD `511e2eb`, whose product tree equals `b7dce40`.

### F7 — durable complete invocation records (both round-02 CRITICALs)

Attempt: reproduce the race; break publication atomicity, refusal semantics, durability barriers, or error handling through the new fault hooks; make crash residue poison a reader; make a store-sync failure admit a reservation.

Writer structure, `PRIMARY` by source read of `internal/telemetry/record.go:409-451`: `os.CreateTemp(i.Dir, ".record-*")` at `:415` stages with mode 0600; write at `:421`; `fsutil.SyncFile` at `:427`; `Close` at `:429`; fault hooks at `:420`, `:424`, `:431`, `:439`, `:446` for stages write, file-sync, close, link, directory-sync; exclusive `os.Link(file.Name(), path)` at `:443`; `quota.SyncPath(i.Dir)` at `:450`. A deferred `os.Remove` at `:419` cleans the staging name. All three record kinds share this path: `requested.json` at `:212`, `started.json` at `:332`, `terminal.json` at `:392`.

Store barrier, `PRIMARY` by source read of `internal/telemetry/record.go:186-200`: `quota.DurableMkdir(directory)` at `:189`; real-directory and symlink refusal at `:192-195`; `quota.SyncPath(directory)` at `:199` after a store-sync fault hook at `:196`. Both constructors call it before reservation: `Begin` at `:223` before `Mkdir` at `:227`; `BeginBound` at `:259` before `Mkdir` at `:263` with a fresh-directory `Lstat` check at `:266`.

The barrier semantics, `PRIMARY` by source read of `internal/quota/history.go:284-296`: `SyncPath` opens the path, syncs it, then syncs its parent through `syncDir(filepath.Dir(path))`. Called with `i.Dir`, it syncs the invocation directory, which persists the published record's directory entry, then syncs the store directory, which persists the invocation directory's own entry. These are the two missing barriers from my round-02 addendum. `DurableMkdir` at `history.go:299-313` is the old private `durableMkdir` exported with a comment change only: stat short-circuit, recursive parent creation, `Mkdir 0700` tolerating `IsExist`, then `syncDir(parent)`. `DurableWrite` now calls it at `history.go:217`. No behavior change to the quota helper.

Race refutation, `PRIMARY` executed:

```sh
go test ./internal/runner -run '^TestStepBudgetConcurrentIndependentGroupsShareCap$' -count=100 -timeout 8m
```

Output: `ok parley-deck-cli/internal/runner 27.152s`, exit 0. Log: `.parley-runtime/credit-failover/zcode1-round03-concurrent.log`. The same command failed 7 of 100 runs at `31c039c` in my round-02 repetition. Zero failures in 100 runs now.

Fault and collision refutation, `PRIMARY` executed:

```sh
go test ./internal/telemetry -run 'TestInvocationPublication|TestInvocationStoreSync|TestInvocationStagingResidue|TestWriteFailuresAreReturned|TestBeginBound' -count=1 -v
```

Output: `--- PASS` for `TestInvocationPublicationFaults` (write, file-sync, close, link, directory-sync), `TestInvocationPublicationRefusesExistingNames` (file, symlink, directory), `TestInvocationTerminalPublicationFaults` (all five stages; repeated `Finish` keeps the error), `TestInvocationStoreSyncFailurePreventsReservation`, `TestInvocationStagingResidueDoesNotEnterAdmission`, `TestBeginBoundRefusesEveryReuseAndCollision`, `TestWriteFailuresAreReturned`, and the other matched tests. The full telemetry package also passed: `ok parley-deck-cli/internal/telemetry 1.934s` (`-count=1 -timeout 10m`).

Why the race is structurally gone, `PRIMARY` by source read: the published name appears only through `os.Link` after complete bytes are written, synced, and closed. A reader sees no file or a complete file. `ExhaustionAdmission` at `internal/telemetry/exhaustion_admission.go:82-95` therefore can no longer parse a partial `requested.json`; a parse failure at `:94` now indicates genuine corruption and still refuses. The same benefit covers `participantRecords` at `internal/runner/dropout.go:404-418` and the budget history readers. `os.Link` fails with `EEXIST` for any existing name form, so the old `O_EXCL` refusal is preserved for files, symlinks, and directories; the collision test verifies all three and that the existing target keeps its bytes.

Launch-permission control flow, `PRIMARY` by source read of `internal/runner/telemetry.go`: `exhaustionAdmission` at `:119` runs before `recordLaunchRequest` at `:179`, and its error returns without a launch. A record-write failure at `:179` also returns the error, so a post-publication sync failure publishes a readable record, launches no child, and grants no permission. The directory-sync fault test asserts the readable published state directly.

D3 verdict: `CONFIRMED (PRIMARY)` for the fault-injection and structural scope. The invocation ledger now meets the same standard as `quota.DurableWrite`: staged complete bytes, checked write, file sync, close, exclusive publication, and both directory barriers. Stated limit: no test on this machine injects real power loss; the fault stages plus the code structure are the evidence, as in the quota package's own tests. A crash that leaves a published `requested.json` without `terminal.json` still trips the designed pending gate at `exhaustion_admission.go:102-105` and blocks for repair, which D3 requires; the `BeginBound` recovery bridge exists for exactly that reconciliation.

### AC-01 — v2 freeze and saved-policy compatibility

Attempt: break legacy byte stability or widen a saved policy. Result: `CONFIRMED (PRIMARY)` for the executed scope. `go test ./internal/quota ./internal/membership ./internal/protocol -count=1 -timeout 12m` printed `ok parley-deck-cli/internal/quota 2.060s`, `ok parley-deck-cli/internal/membership 21.422s`, `ok parley-deck-cli/internal/protocol 0.543s`. The F1 protection sites are intact across the current diff: `Roles.Protected` at `internal/quota/quota.go:144-146` excludes `Assignment.Implementer`, and both succession call sites pass `exhaustedCandidateIDs` (`internal/quota/quota.go:275`, `internal/quota/roles.go:210`).

### AC-02 — native evidence and reset separation

Attempt: feed near-miss framing and uncertain resets. Result: `CONFIRMED (PRIMARY)` for synthetic fixtures. The full telemetry suite above includes the recognizer and capture fixtures. No native positive is claimed; the Zcode native-positive limit and the exact-version Codex limit stay binding.

### AC-03 — one child and durable suppression

Attempt: replay suppression under changed context. Result: `CONFIRMED (PRIMARY)`. `go test ./internal/runner -run 'TestExhaustion|TestSelectedHandover|TestActiveExecution|TestDraft|TestParticipant' -count=1 -timeout 10m` printed `ok parley-deck-cli/internal/runner 0.921s`. Suppression derives from the durable ledger, and the ledger's writer is now durable under F7.

### AC-04 — terminal/validation recovery and stale output

Attempt: corrupt receipts; present stale artifacts. Result: `CONFIRMED (PRIMARY)` for the tested cases. The runner focused suite above passed. Validation receipts still use `quota.DurableWrite`; the invocation side now matches it.

### AC-05 — transient parity and concurrency

Attempt: break transient retries; reproduce the concurrent refusal; check attempt accounting after a spurious refusal. Result: `CONFIRMED (PRIMARY)` for the executed scope. The readiness suite `go test ./internal/app -run 'TestExhaustionReadiness' -count=1 -timeout 10m` printed `ok parley-deck-cli/internal/app 10.987s`; it counts children and five-second spacing through real synthetic processes. The grouped-accounting check `go test ./internal/runner -run 'TestLaunchPreflightChargedGroupedSessionReusableAtCap' -count=1` printed `ok parley-deck-cli/internal/runner 0.306s`. My 100-run concurrent repetition passed. The round-02 risk I could not verify then — a spurious admission refusal inside a participant step interacting with attempt accounting — is gone at the source: the spurious-refusal cause was the partial-record read, which atomic publication removed, and an admission refusal writes no record at all. The organizer's own repeated check also passed: `fixup02-concurrent.log` holds `ok parley-deck-cli/internal/runner 132.081s`.

### AC-06 and AC-07 — owner-order succession and frozen facilitator

Attempt: synthesize a transition from an ordinary loss; count the facilitator toward the floor. Result: `CONFIRMED (PRIMARY)`. The quota and membership suites above passed, including the whole-batch, floor, and external-facilitator tests verified in round-02 and untouched since.

### AC-08 — active execution authority

Attempt: find a fixed `IMPLEMENTATION.md` reader that bypasses the resolver. Result: `CONFIRMED (PRIMARY)`. The protocol suite passed; the resolver and its consumers are unchanged since my round-02 audit, and the fix-up-2 diff does not touch them.

### AC-09 — shared-draft integrity

Result: `CONFIRMED (PRIMARY)`. Covered by the protocol and membership suites above; unchanged since round-02.

### AC-10 — retained obligations and independent review

Result: `CONFIRMED (PRIMARY)` for synthetic succession coverage. The membership suite passed. No fresh goal PASS exists yet; that AC-15 gate stays open.

### AC-11 — exact Codex support

Result: `CONFIRMED (PRIMARY)` within the source-backed limits: exact version `0.165.0`, explicit JSON argv, the four exact messages. The telemetry suite passed. The smoke receipt remains ordinary-behavior evidence only; no exhaustion was captured.

### AC-12 — private capture controls

Result: `CONFIRMED (PRIMARY)`. The telemetry suite passed, including permissions, symlink, scrub, retention, and capacity cases. Capture behavior and machine-only enablement are unchanged by F8; the fix-up-2 diff touches no capture code.

### AC-13 — exact hunks, staged core, size limits

Attempt: find a copy or limit that drifted. Result: `CONFIRMED (PRIMARY)`. Executed `shasum -a 256`: live `parley-deck/COOPERATION.md` = `b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709`; bundled `internal/protocol/defaults/COOPERATION.md` = `43f2be0ef5fb03aa1a669697362446616925d647563f0bafa985d89e41c96994`; skill reference = `b9284d7d…` (same as live); staged core `~/.parley/staging/COOPERATION-2.20.0.md` = `6c437ed5fa0805bc8929b1aa25371664f9c4df1fc34c731a9df02bf02f2bd357`. The signed hunk set hashes to `f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c`. Both reviewers reconstructed all four copies from baselines plus exactly these hunks in round-02; the hashes are unchanged since. Packet test executed: `facilitator body=68448 B` and `phase-1 protocol_change facilitator body=69880 B`, both under 70000, `PASS`. `wc -c` on the skill `SKILL.md` prints `19989`, under 20000.

### AC-14 — versions, guidance, full suites

Result: guidance and versions `CONFIRMED (PRIMARY)`; full-suite portion open, pending. `VERSION` prints `1.56.0`. F8 text verified in all three places. Skill suite: the organizer log `fixup02-skill-tests.log` ends with `pass 402`, `fail 0`, `python 3.15: 54 tests OK across 7 files`, and payload manifests ok; its `parley-deck` manifest line `sha256:821fb014…fcd341` matches the `aggregate` in the tree's `parley-addon.json`, which I compared myself. The portable build log holds only the header, but the artifact exists: `dist/parley-deck-skill-v2.20.0-macos-arm64`, 65521856 bytes, timestamped 00:35, matching the log time. The required `go test ./... -count=1 -timeout 45m` is still running in the organizer process (alive; `fixup02-full-go.log` holds three package lines). Pending is never passed; I record AC-14 as open on that result only.

### AC-15 — independent close and release

Result: `UNVERIFIED`, correctly open. No release, channel publication, or installation exists. D12 orders delivery after this review round, a review consensus with zero agreed fixes, and a fresh independent goal check; channel evidence follows the authorized publication in a separate post-release process. This ordering is FINAL's design, not a defect. CLI WinGet stays held; npm and global-core publication stay owner-only.

### Decision spot checks D1–D12

D1: `CONFIRMED (PRIMARY)` through the quota and membership suites. D2: `CONFIRMED (PRIMARY)` within AC-02's synthetic limits. D3: `CONFIRMED (PRIMARY)` for the F7 fault-injection scope; see the stated power-loss limit. D4, D5: `CONFIRMED (PRIMARY)` through AC-06/AC-07 suites. D6, D7: `CONFIRMED (PRIMARY)` through the protocol and membership suites. D8: `CONFIRMED (PRIMARY)` for synthetic coverage; no goal PASS yet. D9: `CONFIRMED (PRIMARY)` within the exact-version, source-backed limits. D10: `CONFIRMED (PRIMARY)`; F8 changed documentation only. D11: `CONFIRMED (PRIMARY)` through the hash and packet checks above. D12: versions and notes staged; delivery `UNVERIFIED` and correctly gated.

## Findings

### [NIT] Fix-up-cycle-2 validation evidence lacks command-bearing receipts

Defect: the fixup02 evidence under `.parley-runtime/credit-failover/` records outcomes without commands. `fixup02-concurrent.log` holds one line, `ok parley-deck-cli/internal/runner 132.081s`, with no count or test name. `fixup02-skill-portable.log` holds only the build header. No `fixup02-*.receipt.json` exists, unlike cycle 1's `fixup01-full-go-v2.receipt.json`, which recorded the command, exit code, commit, and log hash.

Impact: an auditor cannot derive from retained evidence which repetition the organizer ran or at which source; the mapping relies on the organizer's word. My own 100-run pass makes the concurrent result moot for this reviewer's verdict, but the close record should stand on its own receipts.

Concrete fix: before close, write receipt files for the repeated concurrent check, the full Go command, and the portable build, each with command, exit code, reviewed commit, and log sha256, matching the cycle-1 pattern. Locations: `.parley-runtime/credit-failover/fixup02-*` and `IMPLEMENTATION.md` fix-up cycle 2, `### Current verification and remaining gates`.

No other finding. No CRITICAL, MAJOR, or MINOR defect remains open from any prior round at this source.

## Updated findings

- Zcode round-02 CRITICAL (concurrent admission race): withdrawn as resolved at `b7dce40`. My own 100-run repetition passed; the writer publishes exclusively; readers cannot observe partial bytes.
- Zcode round-02 NIT (capture status deletes expired evidence): resolved by F8 as proposed. Verified in `internal/app/app.go:259-260`, `docs/credit-exhaustion.md:102-103`, and the packaged reference.
- Kimi round-02 CRITICAL (records visible before complete publication): resolved by F7, same evidence.
- Kimi round-02 MINOR (capture-status documentation): resolved by F8, verified.
- Kimi round-02 NIT (`gofmt`): resolved by F9. My `gofmt -l` over all 74 changed Go files printed nothing.
- Round-01 findings F1–F6: all remain resolved across the current diff. F1 sites re-verified this round; readiness tests re-run and passed; `VERSION` re-read as `1.56.0`; F3, F5, F6 sites are untouched by the fix-up-2 diff, which changes only the six files listed in the Summary.
- AC-05 and AC-14, the two round-02 WRONG verdicts: AC-05 is `CONFIRMED (PRIMARY)` for my executed scope now; AC-14 keeps only its full-suite portion open, pending the organizer's running command.

## Open questions

1. Will the organizer's running `go test ./... -count=1 -timeout 45m` exit 0 at the fixed source, and will its receipt be written? Until then AC-14 stays open and no close may cite it.
2. Which fresh independent process performs the pre-release goal check after a zero-agreed-fixes consensus, and which separate process verifies release channels, hosted Windows skill hashes, and installed runtime hashes under D12?
3. Does the owner want the NIT above fixed in this cycle through receipts, or accepted as the cycle-2 evidence record with my independent repetitions standing in?
4. The driver still refuses with `has no runs yet`. The release inbox note must keep listing the pending owner commands, versions 2.15.0 through 2.20.0, with native-fallback receipts kept separate from driver accounting.

Independence: `review/round-03/` was empty when I checked; I have not read any peer round-03 output. I wrote this file in a separate native process under the owner-authorized fallback. The product tree stayed unchanged; my checks wrote only under `.parley-runtime/credit-failover/` (`zcode1-round03-concurrent.log`).

## Supplemental evidence assessment — 2026-10-30 (round-03, unchanged source)

Scope: this supplement assesses the new fix-up-cycle-2 evidence receipts only. It replaces no part of my full-scope review above. The product source did not change. `git rev-parse HEAD` prints `511e2eb`; `git diff --stat b7dce40..HEAD -- ':!parley-deck'` and `git status --porcelain -- ':!parley-deck'` stay empty. The skill worktree HEAD is `406b6976cc75344d76f908297fe46ec2443e62cb`, with a clean tree. Each check below is my own execution, unless I name the organizer as executor.

### NIT disposition — resolved and withdrawn

My NIT asked for three command-bearing receipts before close. All three exist now under `.parley-runtime/credit-failover/`. I verified each against the artifacts myself.

1. `fixup02-concurrent.receipt.json` records the exact command, both test names, `count=100`, exit 0, product commit `b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666`, cwd, and log hash. The command is `go test ./internal/runner -run 'Test(StepBudgetConcurrentIndependentGroupsShareCap|StepBudgetGroupedRunnerCountsOnceAndRetainsRefusals)' -count=100 -timeout 10m`. I hashed `fixup02-concurrent.log` myself: `b6af0c1487bed310cf30fb6afe7bf25e3e113f6b83dd2d2899112f7cd85e189d`. It matches. Both tests exist at this source: `internal/runner/step_budget_test.go:23` and `:119`.
2. `fixup02-skill-portable.receipt.json` records `npm run build:portable:current`, exit 0, source commit `406b6976cc75344d76f908297fe46ec2443e62cb`, cwd, log hash, binary hash, size, and the executed version result. I hashed the binary myself: `dd821c1646c0f2843db05b146a28764a1f4e21e47fcd1cec19e70515eb1ca8f9`, 65521856 bytes. Both match. I also executed `dist/parley-deck-skill-v2.20.0-macos-arm64 --version` myself. It printed `2.20.0`, exit 0.
3. `fixup02-full-go.receipt.json` came from the command wrapper when the process ended at 23:46:12Z. It records `go test ./... -count=1 -timeout 45m`, product `b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666`, exit 0, 34 packages, start and end times, and the log hash. I hashed `fixup02-full-go.log` myself: `68cd0dedd2609cf803e8e24d8ad741a09e48d2193d84f911f90e58864e6172d5`. It matches. I counted 34 `ok` lines and zero `FAIL` lines myself. Its shape matches the cycle-1 receipt `fixup01-full-go-v2.receipt.json`.

The canonical `source-context/fixup02-validation.json` embeds both command receipts and the full Go package output. I recomputed the seven check hashes, the full-Go log hash, and the binary hash; each matched. Receipt times agree with the file times I read with `stat`. `IMPLEMENTATION.md` records completion under `### Fixed-source verification complete — 2026-10-29 23:46Z`. The concrete fix I demanded is delivered in full. I withdraw the NIT.

Provenance and limits. The concurrent and portable receipts state that the organizer recorded them after its own tool run; they are not driver accounting. I did not execute those two commands and I claim no such run. Their exit codes stay organizer statements. The log hashes, the log bytes, the source mapping, and the version result are my own `PRIMARY` checks. The JSON top-level `at` field keeps its first collection time, 23:38:09Z; each embedded record carries its own later time. I rate that no defect. My own executions in the full-scope review still stand: the 100-run concurrent repetition, the fault tests, the focused suites, and the gofmt scan.

### AC-14 update — full-suite portion

Position change: the AC-14 full-suite portion moves from `pending` to `passed on organizer evidence`. Evidence: the wrapper-written `fixup02-full-go.receipt.json` and the completed `fixup02-full-go.log`, both hash-verified by me, with 34 `ok` packages and no failure line. Limit: I did not execute the full suite in this round; my executed coverage stays the focused suites named in my review. This correction needs no duplicate full suite, and I ran none. All other AC and D verdicts stay unchanged at the same source.

## Supplemental findings

### [NIT] Fix-up-2 plan checklist keeps a completed check item open

`IMPLEMENTATION.md:192` still reads `- [ ] Run persistence, grouped accounting, repeated concurrency, race, full Go, skill, package, and portable checks.` The same file records that item as complete at `IMPLEMENTATION.md:234-241`, `### Fixed-source verification complete — 2026-10-29 23:46Z`. Impact: the living record states one status in its plan and another in its evidence section. Concrete fix: tick the item at line 192. The item at line 193 stays open with cause: reviews, consensus, and the goal check remain outstanding.

### Boundaries

This supplement changes no signoff and no gate. AC-15 delivery stays open by D12's order; no release approval follows from it. The fresh goal check follows a zero-fix consensus. I read no peer round-03 file; the directory holds only my own. I edited no code and no other artifact, and I committed nothing.

## Checklist-correction assessment — 2026-10-30 (round-03, second supplement)

Scope: this second supplement assesses the checklist correction and the delivery statements of the proposed cycle-3 consensus. It replaces no earlier section. The product source did not change. Every check below is my own execution, `PRIMARY`.

### Checklist NIT disposition — resolved and withdrawn

Executed `git diff -- parley-deck/ideas/meta-protocol-change-credit-exhaustion-failover/IMPLEMENTATION.md` at HEAD `511e2eb`. The diff shows three hunks. Two are the correction:

1. `IMPLEMENTATION.md:27` changed `- [ ] Checks: focused tests and fault cases, relevant race tests, full Go suite, skill suite, package checks, portable builds.` to the same text with `[x]`.
2. `IMPLEMENTATION.md:192` changed `- [ ] Run persistence, grouped accounting, repeated concurrency, race, full Go, skill, package, and portable checks.` to the same text with `[x]`.

The third hunk adds `### Fixed-source verification complete — 2026-10-29 23:46Z`. That section predates this assessment. My first supplement cited it at `IMPLEMENTATION.md:234-241` and verified its receipts.

The ticks state no new result. Every named check is complete in that section and in the receipts I hash-verified myself: full Go exit 0 with 34 packages, the 100-run concurrent and grouped checks, the portable build hash, and version output `2.20.0`. The plan and the evidence section now agree.

The still-open items stay open, correctly. `IMPLEMENTATION.md:28-30` keep AC-15, release, and delivery verification open. `IMPLEMENTATION.md:193` keeps reviews, consensus, and the goal check open. The goal check has not run, so line 193 stays open as a whole.

No product edit occurred. Executed now: `git rev-parse HEAD` prints `511e2eb`; `git diff --stat b7dce40..HEAD -- ':!parley-deck'` and `git status --porcelain -- ':!parley-deck'` are both empty. The skill worktree is clean at `406b6976cc75344d76f908297fe46ec2443e62cb`. The correction edits the living execution record only. It matches the proposed disposition: a factual execution-record correction with no product edit.

I withdraw the checklist NIT as resolved.

### Consensus delivery statements

The proposed consensus states the delivery scope. I verified its counts and order, all `PRIMARY`:

- Skill `lib/installer.js:18-123` lists exactly 15 managed targets: codex, claude, agy, gemini, hermes, qwen, codebuddy, goose, kimi, droid, vibe, cursor, opencode, aionrs, zcode.
- Executed `find skills -name SKILL.md` in the skill worktree. It prints exactly six files: parley-bidding, parley-deck, parley-design-check, parley-design, parley-tracker, parley-worktrees.
- `source-context/ORGANIZER-BRIEF.md:147-151` names the same four generic destinations and requires hash verification of every installed `SKILL.md`. `:152-159` keeps npm and `parley protocol publish` owner-only and orders the pending commands from version 2.15.0 through 2.20.0. `:154` assigns channel verification to a separate Kimi or Zcode process.
- The consensus gate order matches FINAL D12 and the brief: signatures, then the fresh independent goal check, then completion and merges, then channels, installs, channel verification, and the final owner-command note.

### Boundaries

This supplement changes no AC or D verdict and no gate. My side holds zero open findings at CLI `b7dce40` and skill `406b6976`. AC-15, the fresh independent goal check, and all D12 delivery gates stay open. For this consensus action I read `review/round-03/kimi-1.md` and the proposed `review/consensus.md`; my earlier independence statement applied at its authoring time. I edited no code and no other artifact, and I committed nothing.
