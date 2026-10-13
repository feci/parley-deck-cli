---
agent: zcode-1
idea: meta-protocol-change-ste100-output-format
role: independent source goal checker
date: 2026-10-13
reviewed-cli-commit: 99d04b3ecd0740fb834e70ec109c9384c7b7cd4d
current-cli-head: bcf1140ff64175dc21e89386a0c0bc8491c0ea72
reviewed-skill-commit: 92e3dc78da2f1c55d35c905cabe57fc69ede0670
process: fresh headless goal check, separate from all review and signoff processes
---

## Scope and freshness

I am zcode-1, the independent source goal checker. This process is fresh. It is separate from every review and signoff process. Codex-1 organized and implemented the idea. I did not implement any part of it.

I read these artifacts in full before any check: `source-context/ORGANIZER-BRIEF.md`, `FINAL.md`, `IMPLEMENTATION.md`, `review/round-01/zcode-1.md`, `review/round-02/zcode-1.md`, `review/consensus.md`, and `review/consensus-cycle-1.md`. I also read the source-context receipts named below.

I changed no product file, signoff, FINAL, IMPLEMENTATION, credential, gateway setting, or worktree declaration. I created this file only. I did not merge, publish, install a managed runtime, apply legacy history, or prune a worktree.

## Protocol context attestation

The launcher supplied one private prompt file. I hashed it before any task action.

```text
/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100/.parley-runtime/invocations/55dbd0fd-6648-48bc-b991-696f49b13655/headless-prompt-3314907935.md
SHA256 db530f811744c59c01c21d4e113580c9c6adbe8804d3075c5033297d334afe4a
Size 134988 bytes, mode 0600, inside a 0700 invocation directory.
```

The hash matched the launcher value. I then read the complete file. Its internal attestation, verbatim:

```text
Protocol context attestation: {"context_mode":"full","source_sha256":"aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec","packet_sha256":"aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec"}
Shadow packet audit: {"would_fallback_reason":"unknown-phase:-1","packet_bytes":0,"source_bytes":130042,"included_blocks":27,"omitted_blocks":43}
```

Mode `full`. Renderer reachable. No fallback. `Transport: github-pr`. The shadow audit line is an unapplied diagnostic. The live source hash equals the deck copy hash I measured below. I read `parley-deck/meta/protocol-changelog.md` through the prompt file. It carries the 2.19.0 entry.

## Tree identity

I verified both worktrees myself. I did not infer identity from the brief.

- CLI worktree HEAD: `bcf1140ff64175dc21e89386a0c0c0bc8491c0ea72` prefix `bcf1140`, branch `ste100`, clean. One commit follows the reviewed commit `99d04b3`. `git diff --stat 99d04b3..HEAD` lists only idea audit artifacts, the goal-launch inbox note, and two `runs/` event files. No product byte changed.
- Skill worktree HEAD: `92e3dc78da2f1c55d35c905cabe57fc69ede0670`, clean. `git diff --stat 0549544..92e3dc7` shows the reviewed skill payload only.

Refutation attempt: I searched the post-review commit for hidden product edits. The file list contains none. The protocol hashes below prove the product bytes directly.

## AC1 — three protocol copies equal base plus H1 and H2; drift guard passes

Attempted refutation: I rebuilt each copy from its git base and the signed payload, then compared bytes.

Commands and locators:

```text
git show d51e1c7:parley-deck/COOPERATION.md                       base 0357d504982f92b713f2f86604129f46e1ef3276664e4da91d39ecf420253f0c
git show d51e1c7:internal/protocol/defaults/COOPERATION.md        base 93bfa57936964e409faf76bed320d1f947563edad1568dbe83a49057a09432d7
git show 0549544:skills/parley-deck/references/COOPERATION.md     base 0357d504982f92b713f2f86604129f46e1ef3276664e4da91d39ecf420253f0c
```

My Python reconstruction loaded `source-context/protocol-hunks.json`. It required exactly one match of the old rule 6 and no prior `## Appendix B`. It replaced rule 6 and appended the appendix. Results:

```text
deck copy      130042 bytes  aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec  byte-equals live file
defaults copy  129796 bytes  ea8f26f11acfe8834061fe286a0812bcf7555ebb14337b501d90006ab893bd05  byte-equals live file
skill copy     130042 bytes  aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec  byte-equals live file
```

All three results match the recorded rows.

Drift guard: `go test ./internal/protocol/ -count=1` exited `ok`. `parley protocol packet check --dir . --json` returned `"ok": true` with 70 blocks. Verdict: CONFIRMED (PRIMARY).

## AC2 — identical rule bodies with the required semantics

Attempted refutation: I compared the appendix body with `references/STE100.md` byte by byte.

The reference starts with `# ASD-STE100 output text`. The protocol appendix starts with `## Appendix B — ASD-STE100 output text`. After each title, the bodies are identical. My direct string compare differed by one leading newline only. The test uses the exact marker split. `node --test test/ste100.test.js` passed 3/3 in my run. It asserts one exact body and sixteen semantic markers: limits, exemptions, authority, binding force, adoption, exception, advisory review, and copyright bounds.

The old rule-6 text (`**English only.**`) appears in no copy. `## Appendix B` occurs exactly once per copy. Verdict: CONFIRMED (PRIMARY).

The Issue 9 currency rests on the round-1 and FINAL records. I did not fetch the external page in this process. Provenance: SECONDARY (`zcode-1/review/round-01`, `FINAL.md`). No acceptance criterion depends on it.

## AC3 — thirteen builders, preserved inputs, required syntax

Attempted refutation: I recounted the call sites and re-ran both policy tests.

`grep -rn "protocol.OutputTextInstruction" internal/ --include="*.go" | grep -v _test.go` returns exactly 13 sites: two in `runner.go`, four in `phase58.go`, two in `consult.go`, one in `steer.go`, two in `driver_consensus.go`, one in `consensus_request_signoffs.go:871`, one in `trajectory_verify.go:294`. The list matches the FINAL table one for one.

`TestOutputTextPolicyAcrossRunnerPrompts` and `TestOutputTextPolicyAcrossDriverPrompts` passed in my run. Each case requires the instruction exactly once, as the prompt suffix. The fixtures inject CRLF frontmatter, a Slovak quote, code, a path with spaces, and an ERROR line. Each case asserts those bytes and the required tokens survive. Required tokens include `## Refutation attempts`, `reviewed-commit:`, `outstanding_agreed_fixes: 0`, `### Signoff:`, `GOAL-CHECK: PASS`, and `GOAL-CHECK: FAIL —`.

`probePrompt` stays excluded. The cycle-1 consensus concurred. It requests a liveness sentinel, not protocol prose. Verdict: CONFIRMED (PRIMARY).

## AC4 — map classification and unchanged size guards

Attempted refutation: I rendered both packet variants myself with a CLI built from the current tree. I compared bytes and hashes with `source-sizes.json`.

```text
parley protocol packet --dir . --phase 1 --track deliberation --audience facilitator --json
  68935 bytes, packet a380478f62aedecd2b8e7b74868cdb3acca578c36267b282862e1bf8920b1e45
parley protocol packet --dir . --phase 1 --track deliberation --audience facilitator \
  --idea meta-protocol-change-ste100-output-format --flag auto_implement --flag protocol_change --json
  71934 bytes, packet b224f24610106f17ebcb20c93558d40620f01d3ee4c56d16e84e8a754fd1388f
```

Both renders matched the recorded sizes and hashes. Both indexes classify the appendix as `on-demand`, `included: false`, with a non-empty trigger. Section 6 stays `always`, `included: true`.

`TestLiveSTEAppendixIsOnDemandAndItsPointerIsRetained` and `TestLiveMapIsCompleteHonorsNeverCutAndRendersEveryPhase` passed in my run. `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail` passed. The guard value stays `70_000` in `internal/app/facilitator_packet_live_test.go:21`. The measured fixture body is 68,874 bytes. `SKILL.md` is 19,482 bytes, under 20,000. The meta-idea variant exceeds 70,000 bytes before this change. FINAL discloses that baseline. Verdict: CONFIRMED (PRIMARY).

## AC5 — package content, template, versions, hashes

Attempted refutation: I hashed every metadata target and read the template and package list myself.

`npm pack --dry-run` lists `skills/parley-deck/references/STE100.md` among 211 files. The third `ste100.test.js` case performs a real generic install. It asserts the installed reference equals the source byte for byte, and `SKILL.md` names it.

Versions: `VERSION` and `internal/app/version.go` say `1.55.0`. `package.json`, `package-lock.json`, and `compatibility.json` say `2.19.0`; `recommendedCli` is `>=1.55.0`. `parley-deck/meta/version.json` hashes match my direct hashes: protocol and packaged protocol `aae6e830…`, skill `64d5236b…`, compatibility manifest `34779a4b…`.

`ARTIFACT_TEMPLATES.md` contains `Preserve verbatim source quotes, including non-English quotes.`. The old translate instruction is absent. Verdict: CONFIRMED (PRIMARY).

## AC6 — staged core 2.19.0

Attempted refutation: I rebuilt the staged core from its base and compared bytes.

`~/.parley/staging/COOPERATION-2.18.0.md` hashes to `9ead4475df6d6fa500317e7d347a0b6422bf3c202aca49f1d4b57a5fe14dc512`. My reconstruction from that base plus H1 and H2 equals `COOPERATION-2.19.0.md` byte for byte. Its hash is `16323e04ffb076db3b2019169010033c6668199ebab0ca6896e79924e74bfe3b`. Base, payload, and result hashes are recorded in `protocol-transform.json`. Verdict: CONFIRMED (PRIMARY).

The staged files live outside the repositories. I verified bytes and reconstruction. A release-channel check belongs to delivery, not this check.

## AC7 — required tests, build, and skill evidence

Attempted refutation: I read the retained logs, re-hashed them, and searched for hidden failures. I then ran focused suites, build, and vet myself at the current head.

Local CLI evidence at `99d04b3`:

```text
go-test-fixup1-full.log  sha 4fd83eb99af8b15a6d29805193a4568dc2eb9742fddb45049eaa928e7d2ea972
  every package "ok"; zero "FAIL" lines; internal/app 730s; internal/runner 149s
go-build-fixup1.log      sha e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 (empty)
go-vet-fixup1.log        sha e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 (empty)
```

The empty logs match the recorded hashes for empty files. My own `go build ./...` and `go vet ./...` exited 0 with no output at `bcf1140`. The empty logs are consistent with silent success, and my runs confirm the outcome on the same product bytes.

Skill evidence: `skill-test-full-r2.log` records 402 node tests and 54 python tests, all pass, zero `✖`. Six payload manifests pass. `skill-build-portable.log` records a real portable build. `skill-install-dry-run.json` reports `"ok": true` for the all-target dry run. `skill-fixture-doctor.json` reports `"ok": true` on a fresh fixture. `skill-pack-dry-run-parsed.json` records `parley-deck-skill@2.19.0` with full integrity fields. I opened each file and checked its content. No pass rests on a filename alone.

My own current-tree runs: `go test ./internal/protocol/`, the two packet-map tests, the four transport tests, both policy tests, `TestQuotaSignoffProcessSettlesBatchAndKeepsCurrentRequiredSet`, and the guard tests all passed. Build and vet exited 0.

Live CI, recorded honestly:

- Runs `38406035407` and `38406028243` at `99d04b3`: ubuntu SUCCESS twice, macos SUCCESS twice, windows IN_PROGRESS twice at check time.
- Main run `37938864253` at `24fe38f`: ubuntu SUCCESS, macos SUCCESS, windows FAILURE. The Windows failure predates this change.
- Runs `38409356907` and `38409349656` at `bcf1140`: ubuntu pass; macos and windows pending at check time.

This is not all-green CI. The Windows baseline is a disclosed, pre-existing failure. FINAL and IMPLEMENTATION disclose it. CLI Windows support stays experimental and CLI WinGet stays held. Verdict: CONFIRMED (PRIMARY), with the Windows limitation stated above.

## AC8 — out-of-scope exclusions

Attempted refutation: I listed every changed file and searched for a parser, linter, dictionary, or script.

`ste_output` appears in non-test Go code only inside the instruction constant at `internal/protocol/outputtext.go:6`. No parser reads it. No gate uses it. No new `.py` or `.sh` file entered the product diff. The changed-file lists match FINAL, the two documented deviations (`parley-addon.json`, `test/installer.test.js`), the authorized fix-up files, and idea artifacts. The reconstruction in AC1 proves no protocol byte changed beyond the two hunks. Verdict: CONFIRMED (PRIMARY).

## Private-file transport challenge

I challenged the transport at both launch boundaries.

- Wiring: `execAgentProcess` (`internal/runner/runner.go:1180`) and `commandForLaunch` (`internal/runner/launch.go:52`) both call `buildAgentInvocation`, which calls `headlessPromptArgs` with the invocation directory. `internal/telemetry/record.go` creates that directory with mode 0700.
- Bound: `inlinePromptArgumentLimit = 120 << 10` = 122,880 bytes. No GOOS branch exists. The value stays below the Linux single-argument limit. The Linux limit itself is RECALL knowledge here; the ubuntu CI flip from red to green is PRIMARY proof.
- Real children: `TestLargeHeadlessPromptReachesRealChildWithExactAttestedBytes` passed in my run. It covers both boundaries and both templates. The child verifies the SHA256 before it prints. Its output equals the prepared prompt exactly. The fixture exceeds 128 KiB.
- Refusals: the bounds test covers a missing directory, a large literal without `{prompt}`, and a still-oversized argument. Each refuses with an error and no args. Short inputs stay byte-identical and create no file.
- Stdin: `TestLargeStdinRemainsExactWithInheritedPromptArgument` passed. Stdin bytes stay exact while argv stays bounded.
- Coverage: `TestLargePromptDeliveryCoversEveryBuiltinArgumentTemplate` passed for every builtin spec. Each argument adapter receives one complete file plus digest. Stdin-only adapters stay unchanged, with no file.
- Privacy: the tests assert mode 0600. My own launch directory shows `drwx------` and `-rw-------`.
- Native use: my own launch wrote `headless-prompt-3314907935.md`, 134,988 bytes. That size exceeds the 122,880 bound. The argv carried the short instruction with the path and hash. I verified the hash before any task action. This is real darwin use. The ubuntu CI supplies the Linux proof.
- Limit: a written file proves transport only. No evidence proves a model read it. My own process did read and verify it.

## Completed-evidence verification

I verified the decisive receipts in `release-delivery/2026-10-13-ste100` and in `source-context`:

- `fixup1-validation-status.json` records exit 0 for the full suite, build, and vet. The three log hashes match my measurements. The suite log contains no FAIL.
- `fixup1-ci-status.json` matches the live CI state I fetched through `gh`.
- `source-verification-fixup1.json` claims: three exact copies, identical rule body, 13 sites, matching metadata. My own measurements confirm each claim.
- `review-consensus-archive.json` records the archive hash `775aee16…`. The archived `consensus-cycle-1.md` hashes to that value. It also equals `git show 99d04b3:…/review/consensus.md` byte for byte. No signoff changed.
- The round-2 review artifact hashes to `c4ade109…` and matches the recorded value.
- `review02-signoff-exec.json` records the zcode-1 signoff process, exit 0. The signoff log shows the appended ACCEPT.

## Close conditions

The controlling brief sets four conditions for source close. I state each outcome. The mechanical single-reviewer exception is false.

1. The final zcode-1 review has no open CRITICAL or MAJOR. HOLDS. Round 2 counts CRITICAL 0, MAJOR 0, MINOR 0, NIT 1. The NIT named a stale frontmatter head. The receipt commit corrected both fields to `99d04b3ecd0740fb834e70ec109c9384c7b7cd4d`. I verified both fields in the current `IMPLEMENTATION.md`.
2. Both current participants authored review-consensus ACCEPTs. HOLDS. `review/consensus.md` carries `codex-1 ✅ ACCEPT` and `zcode-1 ✅ ACCEPT`. The archive preserved the earlier cycle byte-exact.
3. A fresh zcode-1 goal check passes. HOLDS through this artifact. See the verdict line at the end.
4. Current-tree independent AC evidence covers AC1–AC8. HOLDS through this artifact. Every criterion above carries my own current-tree command or byte measurement.

All four conditions hold. I did not merge, close, or mark the idea complete. A PASS permits consideration of close by the owner. Delivery DC1–DC6 remains future work. GitHub releases, tap updates, the skill WinGet PR, runtime installs, the separate channel check, and the final handoff are not complete. npm and core publication and the D6 legacy apply stay attended owner actions.

## Coverage limits

- I ran locally on macOS only. Linux rests on CI logs. The kernel constant stays RECALL. No native Windows validation exists anywhere in this evidence.
- Windows CI is pending at both branch commits and fails on main before this change. This is not all-green CI.
- I did not repeat the full twelve-minute suite. I verified the retained same-tree log by hash and content. I ran every focused suite, build, and vet myself. The two failures from round 1 no longer reproduce: the map test expects 70 and the quota signoff test passes.
- The staged core was verified by hash and reconstruction. No release channel carried it yet.
- The attribution of the 20/25/6 numbers to the complete ASD standard remains UNVERIFIED, as FINAL records. They are local working rules. No criterion depends on the attribution.
- File delivery proves transport. It does not prove that any model read the file.
- This file uses ASD-STE100 working rules for its new prose. Quotations, commands, paths, hashes, and frontmatter stay exact.

## Commands run

```text
shasum -a 256 <prompt file, three protocol copies, staged core files, STE100.md, SKILL.md,
  compatibility.json, protocol-hunks.json, consensus-cycle-1.md, round-02/zcode-1.md, three fixup1 logs>
git rev-parse HEAD; git status --short; git log --oneline 99d04b3..HEAD; git diff --stat 99d04b3..HEAD   (CLI worktree)
git show d51e1c7:<two protocol copies> | shasum -a 256
git -C ste100-skill rev-parse HEAD; git status --short; git diff --stat 0549544..92e3dc7; git show 0549544:<skill copy> | shasum
python3 <reconstruction of four targets from bases plus H1/H2, hash comparison>
git show 99d04b3:…/review/consensus.md | shasum -a 256   (archive integrity)
grep -rn "protocol.OutputTextInstruction" internal/ --include="*.go" | grep -v _test.go
grep -rn "ste_output" internal/ --include="*.go"
go test ./internal/protocol/ -count=1
go test ./internal/protocolpacket/ -run "TestLiveMapIsCompleteHonorsNeverCutAndRendersEveryPhase|TestLiveSTEAppendixIsOnDemandAndItsPointerIsRetained" -count=1 -v
go test ./internal/runner/ -run "TestLargeHeadlessPromptReachesRealChildWithExactAttestedBytes|TestHeadlessArgumentBoundsPreserveShortInputsAndRefuseUnsafeDelivery|TestLargeStdinRemainsExactWithInheritedPromptArgument|TestLargePromptDeliveryCoversEveryBuiltinArgumentTemplate" -count=1 -v
go test ./internal/runner/ -run "TestOutputTextPolicyAcrossRunnerPrompts|TestInteractiveArgumentDeliveryAndBound" -count=1
go test ./internal/app/ -run "TestOutputTextPolicyAcrossDriverPrompts|TestQuotaSignoffProcessSettlesBatchAndKeepsCurrentRequiredSet|TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail" -count=1
go build ./... ; go vet ./... ; go build -o /tmp/parley-goalcheck ./cmd/parley
/tmp/parley-goalcheck --version ; /tmp/parley-goalcheck protocol packet check --dir . --json
/tmp/parley-goalcheck protocol packet <both recorded variants> --json   (byte and hash comparison)
node --test test/ste100.test.js ; node --test test/lean-organizer.test.js   (skill worktree)
npm pack --dry-run   (skill worktree)
gh pr checks 84 --repo feci/parley-deck-cli ; gh pr view 84 --repo feci/parley-deck-cli --json headRefOid,state,title
gh run view 38406035407 / 38406028243 / 37938864253 --repo feci/parley-deck-cli --json headSha,jobs,status,conclusion
gh run list --repo feci/parley-deck-cli --branch ste100 / --branch main --workflow Tests --limit 6 --json …
```

GOAL-CHECK: PASS
