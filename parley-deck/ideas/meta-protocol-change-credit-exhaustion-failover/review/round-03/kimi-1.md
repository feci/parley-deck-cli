---
agent: kimi-1
idea: meta-protocol-change-credit-exhaustion-failover
review-round: 3
date: 2026-10-30
reviewed-commit: b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666
responding-to: [zcode-1/review/round-02]
---

## Summary

The reviewed implementation has no objective finding: `CONFIRMED (PRIMARY)` for AC-01 through AC-14 within the limits below. F7-F9 satisfy the signed second fix-up plan: `CONFIRMED (PRIMARY)`. F1-F6 remain resolved across the complete diff: `CONFIRMED (PRIMARY)`.

AC-15 remains `UNVERIFIED (PRIMARY — gate inspection)`. Review consensus, a fresh independent goal PASS, release, installation, and channel review remain open. No release or installation exists.

CLI target: `b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666`, compared against `c2d2431`.
Skill target: `406b6976cc75344d76f908297fe46ec2443e62cb`, compared against `f557c5e`.

Protocol attestation:

context_mode: full
source_sha256: b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709
packet_sha256: b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709
fallback_reason: none

I reviewed the full product and skill diffs, the staged core, both signed fix plans, and both complete round-02 files. Current metadata changes do not alter the reviewed product. I made no provider call and induced no exhaustion. I did not read any round-03 peer output.

## Position changes since prior review round

- I withdraw my round-02 CRITICAL finding: `CONFIRMED (PRIMARY)` as resolved. Atomic publication, directory barriers, and 100-run concurrent and grouped checks pass.
- I withdraw my capture-status MINOR finding: `CONFIRMED (PRIMARY)` as resolved. CLI help and both references now disclose deletion and `cleanup.json` updates.
- I withdraw my `gofmt` NIT: `CONFIRMED (PRIMARY)` as resolved. My changed-file `gofmt -l` check printed nothing.
- My F1-F6 dispositions stand: `CONFIRMED (PRIMARY)`. The readiness, version, unreadable-record, umask, wait, and changelog checks pass at current source.
- AC-05 and AC-14 change from `WRONG` to `CONFIRMED (PRIMARY)` within the tested scope. The former race did not recur, and the required full suite passed at the fixed commit.
- AC-15 stays `UNVERIFIED (PRIMARY — gate inspection)`. D12 still places delivery after independent review, consensus, and the fresh goal check. Channel verification follows release.

## Responses to other reviewers

### @zcode-1 — review/round-02

I concur with your corrected durability assessment and the signed F7 resolution: `CONFIRMED (PRIMARY)`. Parity with the old writer could not satisfy D3. The new writer uses complete staged records, exclusive links, and `quota.SyncPath(i.Dir)`.

Your store-creation addendum is covered: `CONFIRMED (PRIMARY)`. `internal/telemetry/record.go:188-200` now calls `quota.DurableMkdir(directory)`, injects `store-sync`, and returns `quota.SyncPath(directory)`. This barrier precedes invocation-directory reservation at `record.go:218-230` and `:255-270`.

I also concur with your refusal to tolerate malformed published JSON: `CONFIRMED (PRIMARY)`. `TestInvocationStagingResidueDoesNotEnterAdmission` keeps crash residue inert and rejects a corrupt published request. The persistence, collision, repeated-failure, and race tests passed in my execution.

Your capture-status documentation alternative is implemented without a behavior change: `CONFIRMED (PRIMARY)`. Your Windows disposition remains sound within the experimental scope: `CONFIRMED (PRIMARY)`. Neither old platform failures nor safe-refusal fixtures establish Windows correctness. CLI WinGet remains held under D12.

I concur with your F1-F6 assessments and the eager read-only candidate deviation: `CONFIRMED (PRIMARY)`. It grants no transient protection or stored succession authority. Immutable history binds a successor only for an exhaustion-backed transition.

## Refutation attempts

### AC-01 — v2 freeze and saved-policy compatibility

Attempt: trace new-idea constructors, reserialize legacy policies, reject malformed triggers, and try an unbound widening revision.

Result: `CONFIRMED (PRIMARY)` for source and synthetic coverage. `internal/quota/quota.go:76-85` separates `NewExhaustionPolicy` from the retained v1 constructor. Legacy golden and owner-widening tests passed.

Executed:

```sh
go test ./internal/quota ./internal/membership ./internal/protocol -count=1 -timeout 10m
```

Output:

```text
ok  parley-deck-cli/internal/quota       1.613s
ok  parley-deck-cli/internal/membership  19.111s
ok  parley-deck-cli/internal/protocol    0.476s
```

### AC-02 — native evidence and reset separation

Attempt: exercise exact Codex messages, source near-misses, malformed framing, quotes, mixed terminals, and uncertain Zcode resets.

Result: `CONFIRMED (PRIMARY)` for synthetic fixtures. The recognizers rejected the tested near-misses. Uncertain resets did not negate valid Zcode exhaustion. Valid output and in-flight success logic remains bounded by validation and overlap checks.

`internal/quota/exhaustion.go:59-67` admits only the four exact Codex strings. `internal/telemetry/exhaustion_admission.go:130-140` requires success to start before and finish after the exhaustion completion.

Executed:

```sh
go test ./internal/telemetry -run '^(TestInvocationPublicationFaults|TestInvocationPublicationRefusesExistingNames|TestInvocationTerminalPublicationFaults|TestInvocationStoreSyncFailurePreventsReservation|TestInvocationStagingResidueDoesNotEnterAdmission|TestExhaustion.*|TestZcodeCapture.*)$' -count=1 -timeout 5m
```

Output: `ok parley-deck-cli/internal/telemetry 0.379s`.

No captured native Codex exhaustion or Zcode native-positive claim follows.

### AC-03 — one child and durable suppression

Attempt: exhaust a synthetic child, replay with changed policy and launch context, and inspect the common pre-launch gate.

Result: `CONFIRMED (PRIMARY)` for the synthetic lifecycle scope. `internal/runner/telemetry.go:118-121` calls exhaustion admission before request recording at `:175-180`. The lifecycle tests observed one child and durable suppression.

Executed:

```sh
go test ./internal/runner -run '^(TestExhaustion.*|TestSelectedHandoverChildAndReplay)$' -count=1 -timeout 5m
```

Output: `ok parley-deck-cli/internal/runner 0.763s`.

### AC-04 — terminal/validation recovery and stale output

Attempt: remove and corrupt validation receipts, replay recovery, and present stale implementation output.

Result: `CONFIRMED (PRIMARY)` for the tested recovery cases. The runner tests above passed. Recovery validates existing output without a child. `internal/runner/dropout.go:163-179` limits missing-receipt recovery to the last invocation.

The tests do not inject a fatal crash after every persistence instruction. I make no claim beyond the exercised cases and source inspection.

### AC-05 — transient parity and protected-step compatibility

Attempt: verify F1 through synthetic readiness children, repeat the former concurrent failure, and exercise grouped accounting refusal paths.

Result: `CONFIRMED (PRIMARY)`. Transient readiness kept two attempts and five-second spacing. Protected and code-mutating paths retained single-attempt gates. Both 100-run checks passed.

Executed:

```sh
go test ./internal/runner -run '^TestStepBudgetConcurrentIndependentGroupsShareCap$' -count=100 -timeout 5m
go test ./internal/runner -run '^TestStepBudgetGroupedRunnerCountsOnceAndRetainsRefusals$' -count=100 -timeout 10m
```

Output:

```text
ok  parley-deck-cli/internal/runner  25.932s
ok  parley-deck-cli/internal/runner  51.088s
```

The readiness check passed `legacy/transient`, `legacy/exhaustion`, `global-default/transient`, and `global-default/exhaustion` in 10.79 seconds.

### AC-06 — owner-order succession

Attempt: exercise ordinary, organizer, implementer, both-role, simultaneous-loss, original-order, and unusable-counterpart cases.

Result: `CONFIRMED (PRIMARY)` for reducer and immutable-history coverage. The quota, membership, and protocol packages passed in AC-01.

`internal/quota/roles.go:144-148` rejects a lost, unusable, or out-of-scope selected successor. It does not search for a substitute.

### AC-07 — frozen facilitator exception and floor

Attempt: transfer implementation to the frozen facilitator, test role-only organizer loss, and try to count that facilitator toward the floor.

Result: `CONFIRMED (PRIMARY)` for the tested reducer cases. `TestExhaustionRoleBatchFloorAndExternalFacilitator` passed within the quota package. The tested transition adds only the facilitator implementer and retains two non-facilitator usable survivors.

No live external-facilitator takeover was attempted.

### AC-08 — active execution authority

Attempt: test chained handovers, filename adversaries, changed bases, symlinks, early handovers, and unreadable active records.

Result: `CONFIRMED (PRIMARY)` for the inspected and tested scope. `internal/protocol/execution.go:25-73` selects validated history and blocks corrupt authority. Protocol and runner handover tests passed.

Executed:

```sh
go test ./internal/driver -run '^TestRebuildDetailSurfacesReadErrors/unreadable_IMPLEMENTATION.md$' -count=1 -v
```

Output: `PASS` and `ok parley-deck-cli/internal/driver 0.232s`.

### AC-09 — shared-draft integrity

Attempt: recover incomplete unsigned drafts, then add signoffs, malformed identity, duplicate metadata, and truncated constructs.

Result: `CONFIRMED (PRIMARY)` for structural fixtures. `internal/protocol/draft_recovery.go:71-78` preserves complete drafts, rejects signed or frozen incomplete drafts, and binds exact archive hashes. Protocol and membership tests passed.

### AC-10 — retained obligations and independent review

Attempt: move a reviewer into implementation, retain its finding, vary Before/After sets, and require distinct snapshot models.

Result: `CONFIRMED (PRIMARY)` for synthetic succession and review-gate coverage. The membership suite passed. `internal/membership/review_gate.go:161-164` rejects empty, equal, or unknown implementer and reviewer snapshot models.

No fresh independent goal PASS exists. That separate AC-15 gate remains open.

### AC-11 — exact Codex support

Attempt: test text mode, wrong versions, malformed effective argv, exact messages, source near-misses, and quoted content.

Result: `CONFIRMED (PRIMARY)` for exact-version, explicit-JSON, source-backed support only. The telemetry tests passed.

I inspected `source-context/codex-json-smoke-receipt.json:49`. It states: `normal explicit-JSON argv/envelope only; no exhaustion induced; no native failure capture or binary equivalence proved`.

### AC-12 — private capture controls and security

Attempt: test enablement precedence, scrub failures, permissions, symlinks, retention, capacity, status, and purge under both umasks.

Result: `CONFIRMED (PRIMARY)` for product behavior. The focused telemetry tests passed. The unsafe-storage test also passed under both umasks. Capture remains machine-only, private, bounded, and non-authoritative.

Executed:

```sh
sh -c 'umask 0077; go test ./internal/telemetry -run "^TestZcodeCaptureRejectsUnsafeStorage$" -count=1'
sh -c 'umask 0022; go test ./internal/telemetry -run "^TestZcodeCaptureRejectsUnsafeStorage$" -count=1'
```

Output: `ok parley-deck-cli/internal/telemetry 0.304s` and `ok parley-deck-cli/internal/telemetry 0.110s`.

### AC-13 — exact protocol hunks and packet limits

Attempt: apply every signed hunk independently to each base and compare all resulting bytes with the four targets.

Result: `CONFIRMED (PRIMARY)`. My read-only Python reconstruction printed 29 hunks and exact byte equality for every target:

```text
hunks 29 sha256 f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c
live occurrences [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1] exact_bytes True after_sha256 b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709
bundled occurrences [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1] exact_bytes True after_sha256 43f2be0ef5fb03aa1a669697362446616925d647563f0bafa985d89e41c96994
skill occurrences [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1] exact_bytes True after_sha256 b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709
staged occurrences [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1] exact_bytes True after_sha256 6c437ed5fa0805bc8929b1aa25371664f9c4df1fc34c731a9df02bf02f2bd357
```

Executed:

```sh
go test ./internal/app -run '^(TestExhaustionReadinessUnpinnedDefault|TestWaitImplementationAbsentOrEarlyHandover|TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail|TestVersionFileMatchesBinaryVersion)$' -count=1 -timeout 10m -v
```

Output included `facilitator body=68448 B`, `phase-1 protocol_change facilitator body=69880 B`, and `PASS`. `SKILL.md` has 19,989 bytes.

### AC-14 — versions, guidance, and full suites

Attempt: run fixed-source focused checks, inspect the full-suite receipt, run the skill suite, and review platform and operator guidance.

Result: `CONFIRMED (PRIMARY)`. Version, guidance, race, persistence, package, portable-build, and full-suite evidence covers this criterion at the fixed commits.

The required full command completed in the organizer's separate process. I did not duplicate it. `.parley-runtime/credit-failover/fixup02-full-go.receipt.json:2-9` records:

```text
"command": "go test ./... -count=1 -timeout 45m",
"product_commit": "b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666",
"exit_code": 0,
"passing_packages": 34
```

The complete log has 34 `ok` package lines and no failure. My SHA256 check printed the receipt value `68cd0dedd2609cf803e8e24d8ad741a09e48d2193d84f911f90e58864e6172d5`.

My independent `npm test` passed 402 Node tests and 54 Python tests. All six payload manifests passed. My focused race command passed telemetry, quota, and runner. Organizer package dry-run and native portable logs also match their exit-0 receipt.

All changed Go files passed my `gofmt -l` check with empty output. `VERSION` and `internal/app/version.go` both say `1.56.0`. The skill package is `2.20.0`.

Windows scope: `CONFIRMED (PRIMARY)` for documentation and fail-closed behavior only. `internal/telemetry/capture_test.go:54-62` expects private-mode refusal and no capture file. `docs/credit-exhaustion.md:90-93` discloses the POSIX-mode limit and experimental status. Invocation link and sync failures remain errors. No Windows correctness claim follows, and CLI WinGet remains held.

### AC-15 — independent close and release

Attempt: inspect review, goal-check, release, installation, and channel evidence.

Result: `UNVERIFIED (PRIMARY — gate inspection)`. `IMPLEMENTATION.md:230-241` keeps fresh review open and states: `These execution results grant no independent review or close verdict.` The fresh independent goal check and all release/install evidence remain open.

D12 places release after independent review, fix-up, and a fresh goal check. AC-15 channel evidence belongs to the separate post-release process. Native fallback receipts supply no driver accounting or automatic membership authority. The phase-6 refusal changes none of these gates.

### Decision spot checks D1-D12

- D1: `CONFIRMED (PRIMARY)` through AC-01 and AC-05. Legacy goldens, owner widening, transient retries, and grouped accounting passed.
- D2: `CONFIRMED (PRIMARY)` within AC-02's synthetic limits. Exact native framing, reset separation, and output precedence passed.
- D3: `CONFIRMED (PRIMARY)` through AC-03, AC-04, and F7. Durable admission, checked publication, and no-child recovery passed.
- D4: `CONFIRMED (PRIMARY)` through AC-06 and the F1 readiness cases. Whole-batch selection follows the owner order.
- D5: `CONFIRMED (PRIMARY)` through AC-07. The sole facilitator addition and fixed floor passed synthetic checks.
- D6: `CONFIRMED (PRIMARY)` through AC-08. Resolver, handover, base-hash, and error-path checks passed.
- D7: `CONFIRMED (PRIMARY)` through AC-09. Structural and byte-binding fixtures passed.
- D8: `CONFIRMED (PRIMARY)` through AC-10's synthetic causal proof. No goal PASS follows from that result.
- D9: `CONFIRMED (PRIMARY)` only within AC-11's exact-version and source-backed limits.
- D10: `CONFIRMED (PRIMARY)` through AC-12 and F8. Capture behavior and cleanup documentation match.
- D11: `CONFIRMED (PRIMARY)` through AC-13's exact reconstruction and packet checks.
- D12: `UNVERIFIED (PRIMARY — gate inspection)` for delivery. Versions are staged, but the fresh goal check, release, installation, and channel evidence do not exist.

## Findings

No CRITICAL, MAJOR, MINOR, or NIT finding remains at the reviewed commits: `CONFIRMED (PRIMARY)` through the refutations above.

## Updated findings

- F1: `CONFIRMED (PRIMARY)` as resolved. Fallback/default candidates no longer gain protected-implementer status. Exhaustion-only succession and transient readiness cases pass.
- F2: `CONFIRMED (PRIMARY)` as resolved. Root `VERSION` and binary metadata both say `1.56.0`; the version test passed.
- F3: `CONFIRMED (PRIMARY)` as resolved. The unreadable active-record regression test passed. Positive stat evidence survives the integrity error.
- F4: `CONFIRMED (PRIMARY)` as resolved. The unsafe-mode fixture passed under umasks `0077` and `0022`. Private product checks remain intact.
- F5: `CONFIRMED (PRIMARY)` as resolved. The real wait test covers absent records and validated early handovers.
- F6: `CONFIRMED (PRIMARY)` as resolved. The CLI changelog starts with `## 1.56.0 — 2026-10-29` and has no empty heading above it.
- F7: `CONFIRMED (PRIMARY)` as resolved within the tested durability scope. Request, started, and terminal records share checked atomic publication. Store barriers precede reservation. Faults, collisions, modes, cleanup, crash residue, and repeated terminal failure passed. I concur with D3 and the signed durability contract.
- F8: `CONFIRMED (PRIMARY)` as resolved. `internal/app/app.go:256-260`, `docs/credit-exhaustion.md:102-104`, and `EXHAUSTION_AND_HANDOVER.md:100-102` disclose status cleanup and `cleanup.json`.
- F9: `CONFIRMED (PRIMARY)` as resolved. The driver import is formatted, and my changed-file `gofmt -l` output was empty.
- Windows disposition: `CONFIRMED (PRIMARY)` within the documented experimental scope. No directory-sync or private-permission check was weakened. Hosted Windows correctness remains unclaimed.
- AC-15: `UNVERIFIED (PRIMARY — gate inspection)` as a process gate. It supplies no reason to hide or downgrade a product finding.

## Open questions

- Which fresh independent process will perform the pre-release goal check after review consensus?
- Which separate post-release process will verify channel bytes, hosted Windows skill hashes, and installed runtime hashes?
- Will the final delivery record retain the owner-only npm and global-core commands in version order 2.15.0 through 2.20.0?

## Cycle-3 consensus assessment — 2026-10-30

Scope: I assessed the proposed `../consensus.md` and the current `IMPLEMENTATION.md` diff. This supplement starts no new product review cycle. My full-scope review above remains the product evidence.

### Source identity and process evidence

`CONFIRMED (PRIMARY)`: the product sources still match CLI `b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666` and skill `406b6976cc75344d76f908297fe46ec2443e62cb`.

Executed in the CLI worktree:

```sh
git rev-parse HEAD
git diff --stat b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666..HEAD -- ':!parley-deck'
git status --porcelain -- ':!parley-deck'
```

Output:

```text
511e2ebe8de5ec35770c28cdee1ab9d41eab4e61
```

Both product-tree checks printed nothing and exited 0. The uncommitted CLI files are idea execution, review, and audit records. The live protocol hash remains `b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709`.

Executed in the skill worktree:

```sh
git rev-parse HEAD
git status --short
```

Output: `406b6976cc75344d76f908297fe46ec2443e62cb`. The status check printed nothing.

`CONFIRMED (PRIMARY)`: both full-scope native reviews exited 0 without timeout. `source-context/review03-kimi-1-native.receipt.json:18-23` and `review03-zcode-1-native.receipt.json:20-23` each record:

```json
"timed_out": false,
"exit_code": 0
```

`CONFIRMED (PRIMARY)`: `consensus-cycle-02.md` preserves the signed cycle-2 bytes. I executed:

```sh
cmp <(git show HEAD:parley-deck/ideas/meta-protocol-change-credit-exhaustion-failover/review/consensus.md) parley-deck/ideas/meta-protocol-change-credit-exhaustion-failover/review/consensus-cycle-02.md
shasum -a 256 parley-deck/ideas/meta-protocol-change-credit-exhaustion-failover/review/consensus-cycle-02.md
```

`cmp` printed nothing and exited 0. The archive SHA256 is `3ed7a36b3c17ab6651553876360d06f1670327b60d7fae6a7dfdbe102380131f`.

### Zero-fix and evidence dispositions

My F1-F9 resolutions stand: `CONFIRMED (PRIMARY)` through the refutations and `Updated findings` above.

Zcode's corresponding resolutions and both evidence-record withdrawals are `CONFIRMED (SECONDARY — zcode-1/review/round-03)`. Its first supplement states `I withdraw the NIT.` at line 190. Its second supplement states `I withdraw the checklist NIT as resolved.` at line 227. Those supplements cite its own PRIMARY receipt, hash, diff, and source checks.

The closure mechanism is author withdrawal after correction, not organizer dismissal or severity suppression. I concur with zero `Agreed fixes`.

`CONFIRMED (PRIMARY)`: the evidence attribution is accurate. `.parley-runtime/credit-failover/fixup02-full-go.receipt.json:2-9` binds the required command, reviewed CLI commit, exit 0, 34 packages, and log hash. The other two receipts use schema `organizer-command-receipt/v1`. Each states: `This retrospective receipt is not a driver ledger record.`

I executed:

```sh
shasum -a 256 .parley-runtime/credit-failover/fixup02-full-go.log .parley-runtime/credit-failover/fixup02-concurrent.log .parley-runtime/credit-failover/fixup02-skill-portable.log
```

Output:

```text
68cd0dedd2609cf803e8e24d8ad741a09e48d2193d84f911f90e58864e6172d5  .parley-runtime/credit-failover/fixup02-full-go.log
b6af0c1487bed310cf30fb6afe7bf25e3e113f6b83dd2d2899112f7cd85e189d  .parley-runtime/credit-failover/fixup02-concurrent.log
04a31d1123443efb29c75d8f2a2fac4a0eb9367fc8c94b3768d15548bbf91949  .parley-runtime/credit-failover/fixup02-skill-portable.log
```

All three hashes match their receipts. These remain organizer execution receipts, not independent full-suite claims. I ran no duplicate full suite.

### Checklist correction

`CONFIRMED (PRIMARY)`: the proposed factual-record disposition holds. Executed `git diff -- parley-deck/ideas/meta-protocol-change-credit-exhaustion-failover/IMPLEMENTATION.md` shows these two corrections:

```diff
-- [ ] Checks: focused tests and fault cases, relevant race tests, full Go suite, skill suite, package checks, portable builds.
+- [x] Checks: focused tests and fault cases, relevant race tests, full Go suite, skill suite, package checks, portable builds.
-- [ ] Run persistence, grouped accounting, repeated concurrency, race, full Go, skill, package, and portable checks.
+- [x] Run persistence, grouped accounting, repeated concurrency, race, full Go, skill, package, and portable checks.
```

`IMPLEMENTATION.md:234-241` and `source-context/fixup02-validation.json` substantiate both completed items. The diff also retains that completion-evidence section. It changes no product source, as the source-identity checks above establish.

The open items remain open at `IMPLEMENTATION.md:28-30` and `:193`. Zcode's second supplement satisfies the draft's pending author-assessment condition. I concur with this disposition.

### Limitations and remaining gates

I retain every limit in my full-scope review. Native evidence remains synthetic or source-backed within FINAL's exact contracts. Fault injection does not prove every power-loss sequence or filesystem. Windows remains experimental, and CLI WinGet remains held. No provider exhaustion was induced.

AC-15 remains `UNVERIFIED (PRIMARY — gate inspection)`. `IMPLEMENTATION.md:28-30` keeps AC-15, release, installation, and separate delivery verification open. FINAL D12 states: `Release CLI 1.56.0 and skill/core 2.20.0 after independent review, fix-up, and fresh goal check.`

`CONFIRMED (PRIMARY)`: the proposed gate order follows FINAL D12 and the controlling brief. All signatures precede the fresh independent goal process. Implementation completion and merges precede channels and installs. A separate channel review follows publication. The final note retains the owner-only commands.

The delivery counts also match: `lib/installer.js:18-128` lists 15 managed targets. `ORGANIZER-BRIEF.md:147-151` names four generic destinations. The skill glob `skills/**/SKILL.md` returned six files: parley-bidding, parley-deck, parley-design, parley-design-check, parley-tracker, and parley-worktrees.

`ORGANIZER-BRIEF.md:155-160` requires `the still-pending owner commands, in version order: npm and core for each still-unpublished version from 2.15.0 to 2.20.0`. The owner alone executes those commands.

`CONFIRMED (PRIMARY)`: the native fallback boundary is correct. The owner answer at `inbox/user-to-codex-1_meta-protocol-change-credit-exhaustion-failover_decision-answer.md:15-18` states: `Do not claim driver accounting or automatic membership authority for fallback steps.` The recorded refusal changes no gate.

Position: ACCEPT. The scheduled post-consensus gates need no reservation. This position supplies no new verification and no AC-15 pass.
