---
agent: zcode-1
idea: meta-protocol-change-ste100-output-format
review-round: 2
date: 2026-10-13
reviewed-commit: 99d04b3ecd0740fb834e70ec109c9384c7b7cd4d
reviewed-skill-commit: 92e3dc78da2f1c55d35c905cabe57fc69ede0670
responding-to: [zcode-1/review/round-01]
---

## Summary

I reviewed the complete diff from d51e1c7 to 99d04b3, with fix-up cycle 1. All three round-1 findings are resolved. AC1–AC6 and AC8 survive refutation again. AC7 waits only for the pending local full-suite receipt; both ubuntu CI jobs passed the whole suite at this commit. My own launch used the new private-file path on macOS, and the file held exact bytes under the recorded hash.

## Position changes since prior review round

- Round 1 requested changes: one CRITICAL, one MAJOR, one MINOR. Round 2 finds all three resolved. My position moves from REQUEST CHANGES to ACCEPT, with one new NIT.
- The round-1 CRITICAL proposed file or stdin delivery under a bound. The consensus chose the private file. I accepted that shape in my cycle-1 signoff. The source now matches it.
- Round-1 open question 1 is closed: the consensus concurred with the `probePrompt` exclusion.
- Round-1 open question 2 is closed: `source-sizes.json` now records the full request argv for both packet variants. I reproduced both recorded hashes with those argv.
- Round-1 open question 3 is partly closed: `validation-status.json` records exit codes and log hashes. The fresh fix-up build and vet logs are still chained behind the active suite. The pending receipt must carry them.

## Responses to other reviewers

There are no other reviewers. Codex is organizer, drafter, and implementer, and casts no review vote. One fresh consult advised the repair; it holds no vote. Its disposition is accurate. The implementer adopted the shared threshold and the all-adapter tests. Unique file creation replaced the fixed-name file. The optional sidecar was dropped because the argument carries the digest. The consult confused one hash; the record corrected it. I re-hashed the defaults copy myself: `ea8f26f11acfe8834061fe286a0812bcf7555ebb14337b501d90006ab893bd05`.

## Refutation attempts

I challenged every criterion again at the new head. Tags follow §15.2.

### AC1 — exact two-hunk protocol copies and drift guard

Attempted: I re-hashed all three copies at 99d04b3. I checked that the fix-up touched no protocol byte. I ran the protocol suite.

- `shasum -a 256` gives `aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec` for the deck copy and the skill `references/COOPERATION.md`. The defaults copy gives `ea8f26f11acfe8834061fe286a0812bcf7555ebb14337b501d90006ab893bd05`.
- The fix-up diff lists no COOPERATION.md path. Only runner files, the packet test, and canonical records changed.
- `go test ./internal/protocol/ -count=1` passed.

Result: not broken. No H1 or H2 byte changed. Verdict: CONFIRMED (PRIMARY).

### AC2 — rule-body parity and semantics

Attempted: I re-checked the skill worktree state. I re-ran the identity and guard tests myself. I read the fresh skill evidence.

- The worktree stands at `92e3dc78da2f1c55d35c905cabe57fc69ede0670`, clean. `references/STE100.md` hashes to `f9eab1935e9681bf13360140cd3673f59598b66e64d94cb13b0ac6e35449ab63`.
- `node --test test/ste100.test.js`: 3/3 pass. `node --test test/lean-organizer.test.js`: 8/8 pass.
- `skill-test-full-r2.log` records 402 node tests, 54 python tests, and 6 payload manifests, all green.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### AC3 — thirteen builders, preserved inputs and syntax

Attempted: I recounted the instruction sites. I re-ran both policy tests.

- `grep -c "protocol.OutputTextInstruction"` finds 2 in `runner.go`, 4 in `phase58.go`, 2 in `consult.go`, 1 in `steer.go`, 2 in `driver_consensus.go`, 1 in `consensus_request_signoffs.go`, 1 in `trajectory_verify.go`. Total: 13.
- `TestOutputTextPolicyAcrossRunnerPrompts` and `TestOutputTextPolicyAcrossDriverPrompts` passed in my runs. The fix-up changed no builder.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### AC4 — map classification and size guards

Attempted: I reproduced both packet hashes from the recorded argv. I re-ran the map tests at this head.

- Unflagged variant: 68,935 bytes, sha `a380478f62aedecd2b8e7b74868cdb3acca578c36267b282862e1bf8920b1e45`. Meta-idea variant: 71,934 bytes, sha `b224f24610106f17ebcb20c93558d40620f01d3ee4c56d16e84e8a754fd1388f`. Both match `source-sizes.json`.
- `TestLiveMapIsCompleteHonorsNeverCutAndRendersEveryPhase` passed with 70 blocks. `TestLiveSTEAppendixIsOnDemandAndItsPointerIsRetained` passed for both audiences. It asserts the on-demand class, the trigger, and the retained section 6 pointer.
- The guard fixture stays 68,874 bytes under the 70,000-byte limit. The skill core stays 19,482 bytes under 20,000.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### AC5 — package content, template, versions, hashes

Attempted: I checked which files the fix-up touched. I read the skill package evidence.

- The fix-up changed no release metadata. The round-1 evidence carries: `VERSION` and `version.go` say 1.55.0; `package.json`, its lock, and `compatibility.json` say 2.19.0; the pack dry-run lists `references/STE100.md`.
- `skill-build-portable.log` records a native portable build, exit 0. `skill-fixture-doctor.json` reports `"ok": true`.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### AC6 — staged core 2.19.0

Attempted: I re-hashed both staged files.

- `~/.parley/staging/COOPERATION-2.18.0.md` gives `9ead4475df6d6fa500317e7d347a0b6422bf3c202aca49f1d4b57a5fe14dc512`. `COOPERATION-2.19.0.md` gives `16323e04ffb076db3b2019169010033c6668199ebab0ca6896e79924e74bfe3b`. Both match FINAL and `protocol-transform.json`.
- The round-1 reconstruction proved the result equals base plus only H1 and H2. No later commit touched the staged files.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### AC7 — required tests and build pass

Attempted: I ran build, vet, and every focused suite myself. I fetched the PR checks through `gh`. I read both fix-up status files.

- `go build ./...` and `go vet ./...` exited 0 in my runs.
- All focused suites passed: transport, policy, map, protocol, interactive bound, and the quota signoff test.
- Both ubuntu jobs at 99d04b3 completed SUCCESS. The CI job runs `go build ./...` and `go test ./... -count=1 -timeout 45m`. The retained pre-fix log `ci-failed-115311350336.log` shows the real E2BIG failure inside that job. The new jobs pass the same full suite.
- The local full suite still runs (PID 72351, started 15:32). Its log holds three packages, all ok. The chained build and vet logs do not exist yet. The organizer marks the state "active; no pass claim yet". I treat the partial log as no pass, per the brief.

Result: no failure evidence exists at this head. Ubuntu CI proves the suite on Linux. The local full-suite pass stays PENDING until the receipt lands. Verdict: CONFIRMED on ubuntu CI and my own runs (PRIMARY); local full suite UNVERIFIED so far.

### AC8 — out-of-scope exclusions

Attempted: I listed the fix-up files. I searched for new parsers, linters, and scripts.

- The fix-up adds `prompt_argument.go`, `prompt_argument_test.go`, wiring in `runner.go` and `launch.go`, the `packet_test.go` update, and canonical records. Nothing else changed.
- `ste_output` appears in Go code only inside the instruction constant. No parser reads it. No gate uses it.
- No new script, linter, dictionary, or metadata parser entered the diff. No protocol byte changed.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### Transport repair — the CRITICAL fix

Attempted: I broke each required property in turn.

- Both headless boundaries: `execAgentProcess` (`runner.go:1180`) and `commandForLaunch` (`launch.go:52`) now pass the invocation directory. `telemetry.Begin` creates that directory with mode 0700 before either call. The real-child test covers both, plus `CommandFor` through the stdin case.
- Every built-in adapter: `TestLargePromptDeliveryCoversEveryBuiltinArgumentTemplate` covers all 18 builtin and ACP specs. Each arg adapter gets one complete file plus digest. Each stdin-only adapter stays unchanged, with no file.
- Custom placeholders: the real-child test covers `{prompt}` and `--prompt={prompt}`. The bounds test proves prompt data never acts as a placeholder. Root substitution never re-expands prompt text.
- Short inputs: the bounds test asserts byte-identical args and no file.
- File privacy: `os.CreateTemp` writes mode 0600. The test asserts it. My own launch file shows `-rw-------` inside a `drwx------` directory.
- Exact full bytes and hash: the real child parses the instruction, verifies the digest, and prints the file. Its output equals the prepared prompt on both boundaries and both templates. The adapter test asserts the file body equals the prompt and the digest appears in argv.
- Refusal paths: a missing directory, a large literal without `{prompt}`, and a still-oversized argument after substitution each return an error and no args.
- Unchanged stdin: `TestLargeStdinRemainsExactWithInheritedPromptArgument` asserts bounded argv and exact stdin bytes.
- Host-independent bound: `headlessPromptArgs` has no GOOS branch. The same 122,880-byte constant applies on every host.
- Interactive guard: it shares the constant and keeps the refusal `interactive-prompt-argument-too-large-use-file`. `TestInteractiveArgumentDeliveryAndBound` passed.
- No model-read claim: the helper comment, the test comment, and the record all state that a file write proves transport only.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### Native use of the fixed path

My own review launch used the candidate CLI `parley-candidate`, version 1.55.0. The candidate was built after commit 99d04b3. The launch wrote `headless-prompt-399840839.md`, 134,710 bytes, into my private invocation directory. That size exceeds both the 122,880-byte bound and Linux's 131,072-byte single-string limit. The file carries mode 0600. Its SHA256 is `a3896c65b551bcd394fff7df2b3a7a4f5258a8f15534749ce3bb03b944076e56`. I verified that hash before any task action. The instruction text matches the `prompt_argument.go` format string exactly. `requested.json` and `started.json` record `context_mode: full` with source and packet sha `aae6e830…`. This is real native use of the fixed path on macOS, not Linux proof. The ubuntu CI supplies the Linux proof. Verdict: CONFIRMED (PRIMARY).

## Disposition of round-1 findings

1. `zcode-1/review/round-01` [CRITICAL] Linux headless launch regression — RESOLVED. The repair bounds every headless argument below 120 KiB on all hosts. Both boundaries use it. Real-child tests prove exact bytes, hash, privacy, refusals, and stdin. Ubuntu CI moved from two red runs at 4aab9f3 to green at 99d04b3 on the same suite. My own launch exercised the file path for real. The fix is not fixture-only: the quota fixture never reads the prompt, and the fix-up did not touch it. The production argument itself is now bounded.
2. `zcode-1/review/round-01` [MAJOR] Stale block-count constant — RESOLVED. The test expects 70 and adds the on-demand and pointer assertions. It passed in my run and in ubuntu CI.
3. `zcode-1/review/round-01` [MINOR] Canonical record omissions — RESOLVED. `IMPLEMENTATION.md` now records both failures, both CI run links, and explicit exit states. `source-sizes.json` records the exact request argv for both packet variants. I reproduced both recorded hashes.

## Updated findings

### [NIT] IMPLEMENTATION.md frontmatter head-commit is stale

The frontmatter says `status: fix-up-cycle-1` but still reads `head-commit: 4aab9f339e81646b5f2bc2325562cfc1f89e3fcf`. Phase 8 requires the head-commit bump with the status bump. The cycle block defers to "the following source commit", which is 99d04b3. Fix it in the receipt commit: set `head-commit` to the cycle head, and name the commit.

Open severity counts: CRITICAL 0, MAJOR 0, MINOR 0, NIT 1.

## Archive of the accepted consensus

I agree: archive the accepted three-fix consensus byte-for-byte as `review/consensus-cycle-1.md` before a fresh zero-fix consensus. The archive must be a pure copy or rename. No existing signoff may change. The fresh consensus must cite both review rounds and list zero agreed fixes. Close still needs the four brief conditions. This review carries no open CRITICAL or MAJOR. Both participants must file review-consensus ACCEPTs. A fresh goal check must PASS. Current-tree AC evidence must be complete; the pending full-suite receipt belongs to it.

## Open questions

1. The local full suite, build, and vet still run behind PID 72351. Record all three final exit states in the receipt before any close claim.
2. The macos and windows legs of CI runs 38406035407 and 38406028243 were in progress at review time. Record their final states with the receipt. The Windows baseline stays a disclosed limit; keep CLI WinGet held.

## Protocol context attestation

Verbatim from my launch instruction (argv):

```text
The complete Parley launch prompt is in the private UTF-8 file below. It includes the full live protocol, its attestation, and the task. Before any task action, read the entire file and verify its SHA256. Follow that exact prompt. If the file is missing, unreadable, or has a different hash, report a blocker and stop.
Prompt file (quoted path): "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100/.parley-runtime/invocations/a2332e7d-807e-47bc-b664-abfd970c901c/headless-prompt-399840839.md"
Prompt SHA256: a3896c65b551bcd394fff7df2b3a7a4f5258a8f15534749ce3bb03b944076e56
```

Verbatim from the private file:

```text
Protocol context attestation: {"context_mode":"full","source_sha256":"aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec","packet_sha256":"aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec"}
Shadow packet audit: {"packet_sha256":"9d8b2219d8a8a949761c9f46deb359f6471f38aecd165c9cbb7ce45527203f1f","packet_bytes":81966,"source_bytes":130042,"included_blocks":37,"omitted_blocks":33}
```

Mode `full`; renderer reachable; no fallback. `Transport: github-pr`. I read `meta/protocol-changelog.md`; it carries the 2.19.0 entry.

## Commands run

Exact commands I executed. All were read-only, except tests and packet renders, which wrote only temporary build and cache paths.

```text
shasum -a 256 <deck COOPERATION.md, defaults COOPERATION.md, skill references/COOPERATION.md, staged core 2.18.0 and 2.19.0, STE100.md, SKILL.md, own headless-prompt-399840839.md>
git -C <cli worktree> diff --stat d51e1c7..99d04b3
git -C <cli worktree> diff --stat 4aab9f3..99d04b3
git -C <cli worktree> diff 4aab9f3..99d04b3 -- internal/runner internal/protocolpacket
git -C <skill worktree> rev-parse HEAD; git status --short; git log --oneline -3
grep -c "protocol.OutputTextInstruction" <seven builder files>
grep -rn "ste_output" internal/ --include="*.go"
grep -rn "{prompt}" internal/ --include="*.go"   (non-test audit)
go build ./...
go vet ./...
go test ./internal/protocol/ -count=1
go test ./internal/protocolpacket/ -run "TestLiveMapIsCompleteHonorsNeverCutAndRendersEveryPhase|TestLiveSTEAppendixIsOnDemandAndItsPointerIsRetained" -count=1 -v
go test ./internal/runner/ -run "TestLargeHeadlessPromptReachesRealChildWithExactAttestedBytes|TestHeadlessArgumentBoundsPreserveShortInputsAndRefuseUnsafeDelivery|TestLargeStdinRemainsExactWithInheritedPromptArgument|TestLargePromptDeliveryCoversEveryBuiltinArgumentTemplate" -count=1 -v
go test ./internal/runner/ -run "TestOutputTextPolicyAcrossRunnerPrompts|TestInteractiveArgumentDeliveryAndBound" -count=1
go test ./internal/app/ -run "TestOutputTextPolicyAcrossDriverPrompts|TestQuotaSignoffProcessSettlesBatchAndKeepsCurrentRequiredSet|TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail" -count=1
node --test test/ste100.test.js; node --test test/lean-organizer.test.js   (skill worktree)
parley-candidate protocol packet --dir . --phase 1 --track deliberation --audience facilitator --json
parley-candidate protocol packet --dir . --phase 1 --track deliberation --audience facilitator --idea meta-protocol-change-ste100-output-format --flag auto_implement --flag protocol_change --json
gh pr checks 84 --repo feci/parley-deck-cli
gh pr view 84 --repo feci/parley-deck-cli --json headRefOid,state,title
gh run view 38406035407 --repo feci/parley-deck-cli --json jobs
gh run list --repo feci/parley-deck-cli --branch ste100 --workflow Tests --limit 6 --json …
ps aux | grep "go test ./..."
```

## Verification limits

- I ran on macOS only. No local Linux or Windows. Linux rests on CI evidence; the kernel constant stays RECALL.
- I did not duplicate the active full suite. I ran focused suites only.
- Full logs for CI runs 38406035407 and 38406028243 stay gated until each run completes. I used the per-job check conclusions.
- The macOS launch proves the file path on darwin. It does not prove Linux. The ubuntu job proves Linux at this commit.
- DC1–DC6 stay out of scope. A separate fresh process verifies them after release and install.
- No evidence proves a model read the file. File write and process start prove transport only. My own process read the file and verified its hash before action, as the instruction required.

## Disposition

ACCEPT — fix-up cycle 1 resolves all three agreed fixes. No open CRITICAL, MAJOR, or MINOR. One NIT remains; fix it in the receipt commit. Do not close before the full-suite receipt and the fresh goal check.

Open severity counts: CRITICAL 0, MAJOR 0, MINOR 0, NIT 1.
