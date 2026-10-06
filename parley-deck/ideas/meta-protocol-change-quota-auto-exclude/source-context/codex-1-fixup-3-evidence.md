# codex-1 fix-up cycle 3 — producer evidence

This is implementation/producer evidence, not independent review, acceptance, or signoff. Work was activated after the 2026-10-06 15:15Z independent plan signoff. This subprocess uses the existing CLI worktree only. The parent owns IMPLEMENTATION, orchestration, G17 and documentation. No participant/provider launch, zcode invocation, capture, parser relaxation, roster/configuration/credential change, worktree operation, commit, merge or release is authorized here.

## Scope constraint discovered: W2 thread stopped

The actual `27e42b8` CLI differential now includes two decline attempts in private copies of the independent D1–D3 program. On `/tmp`, literal `--status '❌ NON-PARTICIPANT'` exits 1 (`unknown signoff status`); the existing supported `--status block --notes '❌ NON-PARTICIPANT' --counter 'Continue without me'` exits 0 and records an ordinary `### Signoff: e` / `❌ BLOCK`, yielding `Consensus: blocked`. The next existing participant's append exits 0. The joiner has no late round-1 artifact.

Restoring that exact canonical signoff while keeping G15's binding condition that an incomplete joiner is not a known signer requires a distinct decline-note interpretation in consensus validation. `internal/consensus/consensus.go` has no decline-note type: `AppendSignoff` accepts only a current participant and `validateDocumentAwaiting` rejects every unknown signoff. That package is outside this subprocess's edit ownership. Treating the joiner as known would contradict G15; inventing a new required join marker or writing a fake owner authorization is forbidden. I have therefore stopped W2's implementation thread, retained the rejection, and disclosed the differential deviation rather than claiming G15 complete. The parent must resolve the precise scope/representation issue. No owner waiver or independent resolution is inferred.

Logs: `.parley-runtime/quota-implementation/fixup-3/producer/logs/baseline-local.log` and `current-local.log`; identical extended probe main files, with only setup differing. The reviewer originals are untouched. Remaining checks and final coverage are appended below when frozen.

## Protocol context and authority read

- `context_mode=full`
- `source_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`
- `packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`
- `fallback_reason=null`
- Body: `.parley-runtime/protocol-packets/full-phase8-deliberation-73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e.md`.
- Actual reading: every line in bounded chunks 1–200, 201–420, 421–620, 621–850, 851–1050, 1051–1260, 1261–1501. Long lines 59, 452 and 1002 were included in full in those chunks. Python hashlib reproduced the packet hash and byte comparison against live `parley-deck/COOPERATION.md` was true. No optimized-context claim or shadow packet.
- Read the current signed `review/consensus.md` in full, including G15–G17 and W1–W5; `review/round-05/claude-1.md` in full (its leading coverage section separately re-read after batched output truncation); the complete owner round05-answer and scope-reset-answer. Read FINAL §§1–11 and AC1–AC21 plus its idempotence/recovery contract. Q2's native-positive AC2 waiver is release-scoped; R5-MAJOR-2 is accepted/deferred, never fixed or passed by this work.
- Consulted the Parley Deck skill, with the user's explicit implementer role and no-provider/no-other-participant instructions overriding organizer startup actions. Consulted the graphify and OpenViking skill instructions; no graph rebuild/provider call was made. The scoped OpenViking query through the configured MCPAnywhere connector returned `MCP tool call requires approval, but approval policy is never`. No memory was read or written; current local authorities supplied the context.

## Implementation plan within remaining scope

1. Permit only an unknown policy-off joiner's own canonical late round-1 launch; retain identity/path/history/policy-on/closed/excluded-member gates and allow an incomplete own stub to be retried.
2. Represent an incomplete manual edit as uncommitted pending catch-up, exposing the intended quorum but retaining the prior known-signer set. Existing CLI signers may append; ordinary CLI completion imports the checked revision. No owner authority or kickoff membership is synthesized.
3. Admit the existing plain participant edit for a kickoff-excluded return during round 1 using the existing prompt exclusion record; later return retains catch-up requirements.
4. Use checked applied receipts as prior notice-delivery proof. Validate any extant live/archive notice; replay a missing or corrupt receipt only after checked reconciliation. Preserve historical manual display text with a separately receipted correction.
5. Add local-stub CLI and notice regressions, execute on both filesystems, report host-proof failures without weakening assertions.


## Frozen producer changes and file coverage (2026-10-06 15:43Z)

Only the following thirteen Go files were edited/created by this subprocess. All are inside the permitted packages. Their exact final bytes are recorded in `producer-file-sha256.json`.

| Files | Purpose / executed coverage |
| --- | --- |
| `internal/protocol/quota_manual.go`, `quota.go` | Checked policy-off launch permission, pending manual catch-up preview, no pending known-signer expansion, closed/identity/path gates; live/archive notice validation on read paths. CLI regressions exercise actual callers. |
| `internal/quota/revision.go`, `revision_snapshot.go` | Round-1 return from the existing exclusion marker and plain participants edit; immutable revision revalidates the same snapshot. No kickoff mutation, new join/return marker, or owner-authority inference. |
| `internal/quota/receipt.go`, `notice.go` (new) | Exact receipt bytes and regular-file checks; read-only validation of live/archive notices and exact historical manual notice wording. |
| `internal/runner/telemetry.go` | The sole dispatch exception is an unknown policy-off id's own canonical late round-1 target. Known exclusions, policy-on, committed pending history, and mismatched launch identity retain their gates. |
| `internal/app/agents_exec.go`, `app.go` | Permit retry of the same incomplete own stub; validate the completed catch-up and import its pending edit after the real local process completes. Existing CLI signers append during pending catch-up; unknown pending signers are rejected before writing. |
| `internal/membership/membership.go`, `notice.go` (new) | Ordinary driving waits for catch-up; durable notice delivery survives normal archival/deletion; checked replay retains one terminal evaluation. Legacy manual wording is preserved and its separate honest clarification has its own delivery receipt. |
| `internal/app/quota_cycle3_test.go` (new) | Real CLI-dispatched local shell stubs for artifact-first, edit-first, failed-stub retry, multiple pending joiners, round-1 kickoff-excluded return and later catch-up; policy-on, foreign/later path, known/kickoff exclusion, closed idea, wrong idea, symlink, foreign stub, completed-file overwrite negatives. Existing signoffs remain partial until the new member is imported. No fabricated owner-authority object; kickoff bytes are compared unchanged. |
| `internal/membership/cycle3_test.go` (new) | Applied archive/delete, before-receipt archive/delete, before-notice interruption, corrupt receipt recovery with archive/delete, corrupt/mismatched/symlink notice and receipt rejection, archived historical manual notice plus archived/deleted clarification, repeated Before and exactly one terminal evaluation. |

The pending catch-up list is internal transient state (`json:"-"`), not a new persisted membership schema. It previews the requested participants for quorum calculation while immutable `ReadHistory` and `Known` remain unchanged. The catch-up CLI command is actionable, and the late artifact remains the existing round-01 contract. No capture/grammar/fixture-label change or D6 work was made.

Changes simultaneously appearing in CHANGELOG, docs, IMPLEMENTATION, organizer usage, usage ledger, release-note draft, the native-capture follow-up prompt and `internal/telemetry/quota.go` were outside my file writes and belong to the parent/other concurrent work. Existing untracked run directories were left alone. I did not touch the skill worktree or any review/signoff file.

## Executed differential and notice evidence

Private probe copies live below `.parley-runtime/quota-implementation/fixup-3/producer/`. The baseline private program is `/tmp/claude1-r5-prechange-27e42b8/.parley-runtime/codex1-fixup3/`; it is an archive, never a worktree. `git archive 27e42b8 internal cmd go.mod go.sum` was compared against 489 archived Go/build inputs: zero mismatches (`baseline-identity.json`). Current/baseline extended `main.go` files are byte-identical, SHA-256 `7a87344bd7f4420372a48c711fafdc1b099df1923dcdc97e55f9a6de968b1c12`. Only their setup differs, as in the independent probe. D4 is the explicitly attributed producer W2 extension. The notice program is byte-identical to the reviewer original. Source hashes are in `probe-source-hashes.json`; reviewer originals were not edited.

| Probe | `/tmp` and shared-volume observations |
| --- | --- |
| D1 actual agents exec before edit | Baseline and frozen current exit 0, with artifact written by the CLI-dispatched local stub. No tested late artifact was prewritten. |
| D2 edit first, existing signer | Baseline and frozen current append exit 0. Current status shows pending catch-up with the own-path command, current requested ids, and unchanged known ids. The probe's unchanged current setup creates a manifest without an events ledger, so status also reports that missing ledger. It is retained evidence, not repaired D6 accounting; permanent app fixtures create their ledger and show the catch-up diagnostic without that setup error. |
| D3 kickoff-excluded round-1 return | Both signoffs exit 0 on baseline and frozen current. Current durable history is revision 1; the original kickoff has no returned id. The revision is manual with nil owner authority. |
| D4 decline | Baseline supported BLOCK/NON-PARTICIPANT note exits 0 and blocks; current exits 1 pending catch-up. Both reject the literal unsupported status. **Unresolved W2 deviation**, not an accepted knob-off change. |
| Original notice probe | Both filesystems: 1 notice after settle, 0 after archive, 0 after each of two Before calls; `round.completed` remains 1. |

Permanent tests additionally execute ordinary edit-first completion and incomplete-stub retry on both filesystems. New transition receipts are written only after checked projections, round evaluation and notice publication/validation. A missing/corrupt receipt is **not** evidence of delivery: checked replay may replace it after validation, preserving the existing recovery tests unchanged. A valid applied receipt plus deletion suppresses republication. Before receipt, a validated archived notice suppresses it; deletion with no receipt permits one benign republication, then the receipt prevents another. Corrupt extant notices gate even with an applied receipt. Historical manual corrections stay explicitly non-authoritative and are independently deduplicated after archive/delete.

## Commands, exits and failures

All Go/probe commands used a private writable `GOCACHE` at `producer/gocache` and `PATH` prefixed with `producer/no-provider-bin`. That directory shadows claude, codex, kimi, zcode, opencode, hermes and agy; version queries return a local stub and any other invocation is denied and logged. `guard-denied.log` does not exist. Test processes used local stubs only. Actual commands, working directories, exact TMPDIR values, exits and elapsed times are recorded without omission in `commands.jsonl`; stdout/stderr are in `logs/<label>.log`.

- Frozen affected-package command, on **both** filesystems: `go test ./internal/protocol ./internal/quota ./internal/membership ./internal/runner ./internal/app -run 'Quota|Cycle2' -count=1` exits **1**. The only reported final failing test is the unchanged `TestQuotaCycle2NativeCrashWriter`: local supervisor/writer/group death cannot be proven. Other selected membership checks run; quota, runner and app packages report `ok`. Protocol had no matching tests under this name filter.
- Full protocol/quota command `go test ./internal/protocol ./internal/quota -count=1` exits 0, with both packages reporting `ok`.
- New-regression command `go test ./internal/app ./internal/membership -run TestQuotaCycle3 -count=1` exits 0 on `/tmp` and the shared volume. The final frozen focused commands execute these regressions again. These are producer test results, not an independent verdict.
- Concrete native-proof environment limit: `sysctl kern.boottime` exits **1**, `Operation not permitted`. I did not replace the boot/death proof, skip the native test or weaken its assertion. A HOST rerun is required; no HOST/native pass is claimed.
- The initial focused run also failed `TestQuotaFixupCorruptReceiptRecoveryChecksAllProjections` and `TestQuotaFixupCLIRevisionRecovery` because my first receipt implementation hard-gated corrupt receipts. I changed the implementation to checked replay, without editing those older tests. They pass in the subsequent focused runs; the native crash failure remains.
- `gofmt -l` over all thirteen producer Go files exits 0 with empty output; `git diff --check -- <producer files>` exits 0 with empty output. Exact arguments/results are in `gofmt-check.json` and `diff-check.json`.

The following table is the complete checked-command ledger (exact argv; environment and cwd are in the JSONL). Setup/copy/read-only source inspection commands are described above and were not test verdicts.

| Label | Exit | Command |
| --- | ---: | --- |
| baseline-local | 0 | `go run ./.parley-runtime/codex1-fixup3 /tmp/codex1-fixup3-baseline-local` |
| format-initial | 0 | `gofmt -w internal/protocol/quota.go internal/protocol/quota_manual.go internal/quota/revision.go internal/quota/revision_snapshot.go internal/quota/receipt.go internal/membership/membership.go internal/membership/notice.go internal/runner/telemetry.go internal/app/agents_exec.go internal/app/app.go` |
| focused-initial | 1 | `go test ./internal/protocol ./internal/quota ./internal/membership ./internal/runner ./internal/app -run 'Quota\|Cycle2' -count=1` |
| current-local | 0 | `go run ./.parley-runtime/quota-implementation/fixup-3/producer/diff-current /tmp/codex1-fixup3-current-local` |
| format-tests | 0 | `gofmt -w internal/app/quota_cycle3_test.go internal/app/agents_exec.go internal/quota/receipt.go internal/membership/membership.go internal/membership/notice.go` |
| cycle3-app-local | 0 | `go test ./internal/app -run TestQuotaCycle3 -count=1 -v` |
| format-notice-tests | 0 | `gofmt -w internal/protocol/quota.go internal/quota/notice.go internal/membership/notice.go internal/membership/cycle3_test.go` |
| cycle3-local | 0 | `go test ./internal/app ./internal/membership ./internal/protocol ./internal/quota ./internal/runner -run 'TestQuotaCycle3\|TestQuotaFixupCorruptReceipt\|TestQuotaCycle2LegacyManualNotice\|TestQuotaFixupCLIRevisionRecovery' -count=1 -v` |
| format-final | 0 | `gofmt -w internal/app/app.go` |
| baseline-shared | 0 | `go run ./.parley-runtime/codex1-fixup3 '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude/.parley-runtime/quota-implementation/fixup-3/producer/baseline-shared'` |
| notice-local | 0 | `go run ./.parley-runtime/quota-implementation/fixup-3/producer/notice /tmp/codex1-fixup3-notice-local` |
| notice-shared | 0 | `go run ./.parley-runtime/quota-implementation/fixup-3/producer/notice '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude/.parley-runtime/quota-implementation/fixup-3/producer/notice-shared'` |
| current-shared | 0 | `go run ./.parley-runtime/quota-implementation/fixup-3/producer/diff-current '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude/.parley-runtime/quota-implementation/fixup-3/producer/current-shared'` |
| focused-final-local | 1 | `go test ./internal/protocol ./internal/quota ./internal/membership ./internal/runner ./internal/app -run 'Quota\|Cycle2' -count=1` |
| focused-final-shared | 1 | `go test ./internal/protocol ./internal/quota ./internal/membership ./internal/runner ./internal/app -run 'Quota\|Cycle2' -count=1` |
| format-freeze | 0 | `gofmt -w internal/protocol/quota_manual.go internal/app/quota_cycle3_test.go` |
| final-cycle3-local | 0 | `go test ./internal/app ./internal/membership -run TestQuotaCycle3 -count=1` |
| final-cycle3-shared | 0 | `go test ./internal/app ./internal/membership -run TestQuotaCycle3 -count=1` |
| gofmt-last | 0 | `gofmt -w internal/app/agents_exec.go internal/protocol/quota.go` |
| boot-proof | 1 | `sysctl kern.boottime` |
| frozen-current-local | 0 | `go run ./.parley-runtime/quota-implementation/fixup-3/producer/diff-current /tmp/codex1-fixup3-frozen-current-local` |
| frozen-current-shared | 0 | `go run ./.parley-runtime/quota-implementation/fixup-3/producer/diff-current '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude/.parley-runtime/quota-implementation/fixup-3/producer/frozen-current-shared'` |
| frozen-focused-local | 1 | `go test ./internal/protocol ./internal/quota ./internal/membership ./internal/runner ./internal/app -run 'Quota\|Cycle2' -count=1` |
| frozen-focused-shared | 1 | `go test ./internal/protocol ./internal/quota ./internal/membership ./internal/runner ./internal/app -run 'Quota\|Cycle2' -count=1` |
| protocol-quota-full-local | 0 | `go test ./internal/protocol ./internal/quota -count=1` |

## Remaining work and proof limits

W2 remains blocked on the scope/representation issue stated first. It must not be relabeled resolved or owner-accepted. The baseline comparison is now executed, not inference. The current CLI protects unknown pending signers; `internal/consensus` has not been adapted into a distinct decline-note channel, and this report makes no such lower-level API guarantee.

The parent still owns full HOST Go/build/vet/gofmt/race/shared/local checks, skill suites, G17/W4/W5 documentation, independent re-review, signoffs and the new attended close. Full `go test ./...`, full build, vet, race, skill tests, Windows runtime and an end-to-end kill of a real parley process were **not executed by this subprocess**. Native-positive AC2 remains NOT MET and owner-waived for this release; no native-positive evidence or capture was obtained.

No self-acceptance, independent green verdict, completed review cycle, commit, release or close is claimed. The producer code is frozen for parent inspection with W2 outstanding; this report is the only canonical artifact authored by this subprocess.
