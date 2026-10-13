---
agent: zcode-1
idea: meta-protocol-change-ste100-output-format
review-round: 1
date: 2026-10-13
reviewed-commit: 4aab9f339e81646b5f2bc2325562cfc1f89e3fcf
reviewed-skill-commit: 92e3dc78da2f1c55d35c905cabe57fc69ede0670
---

## Summary

I reviewed both frozen worktrees in full. The design deltas are exact: three
protocol copies, the skill reference, the map, the staged core, and the metadata
all match their bases plus only H1 and H2. The thirteen prompt paths carry the
shared instruction, and the focused tests protect exact inputs and syntax.

AC1–AC6 and AC8 survive refutation. AC7 does not. The full CLI suite fails at
the reviewed HEAD, in two ways. First, a real Linux regression: the larger
protocol now crosses the kernel limit for one exec argument, so headless
launches of argv-mode agents fail with `argument list too long`. Second, a
stale test constant expects 69 protocol blocks where the live source now parses
into 70. Ubuntu CI is red on this branch and green on main for the same tests.

I filed one CRITICAL, one MAJOR, and one MINOR finding. Fix-up cycle 1 is
required before any close. I created this file only. I changed no product file,
and I made no commit, merge, or release.

## Refutation attempts

I tried to break each criterion against the real trees. Tags follow §15.2.

### AC1 — exact two-hunk protocol copies and drift guard

Attempted: I audited `verify-source.py` line by line and confirmed it is
read-only. I ran it. I also re-measured every hash with `shasum -a 256` on my
own, counted diff hunks per copy, and ran the drift guard suite myself.

- `python3 verify-source.py`: all three copies equal base + H1 + H2, byte-exact.
- Hashes: deck and skill copies both `aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec`; defaults copy `ea8f26f11acfe8834061fe286a0812bcf7555ebb14337b501d90006ab893bd05`. Both match `protocol-transform.json`.
- Each copy's `git diff` shows exactly two hunks: the rule-6 replace and the appendix append. No other byte changed.
- `go test ./internal/protocol/ -count=1` passes. This includes `TestEmbeddedDefaultMatchesLiveDeck`.
- `parley protocol packet check` returns `"ok": true` with 70 blocks.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### AC2 — rule-body parity and semantics

Attempted: I compared the Appendix B body with `references/STE100.md` after its
title. I ran the identity test, which also asserts sixteen semantic markers:
limits, exemptions, authority, binding summary, adoption, exception, advisory
review, and copyright bounds.

- `node --test test/ste100.test.js`: 3/3 pass. The bodies are identical.
- The old authority text still appears in no copy: `Appendix B` occurs once per protocol copy.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### AC3 — thirteen builders, preserved inputs and syntax

Attempted: I listed every prompt builder in `internal/` and counted the real
call sites. I ran both new tests. I classified each excluded builder by hand.

- `grep -rn "protocol.OutputTextInstruction" internal/runner internal/app` finds exactly 13 non-test sites. They match the FINAL table one for one: 2 in `runner.go`, 4 in `phase58.go`, 2 in `consult.go`, 1 in `steer.go`, 2 in `driver_consensus.go`, 1 in `consensus_request_signoffs.go:871`, 1 in `trajectory_verify.go:294`.
- `go test ./internal/runner/ -run TestOutputTextPolicyAcrossRunnerPrompts` and `go test ./internal/app/ -run TestOutputTextPolicyAcrossDriverPrompts` pass. The fixtures keep CRLF bytes, a Slovak quote, code, paths with spaces, frontmatter, and every required report and signoff token.
- Excluded builders are outside the design boundary: `SeedBlockPrompt`, `retroPromptBody`, and `writeQuotaPrompt` are driver-authored kickoff scaffolds, which FINAL exempts. `probePrompt` is a liveness sentinel, not a prose artifact. See Open questions.

Result: not broken. Verdict: CONFIRMED (PRIMARY), with one boundary question.

### AC4 — map classification and size guards

Attempted: I rendered both packet variants with the built CLI and checked the
index entries myself. I ran the guard tests and read their thresholds in source.

- Unflagged variant: 68,935 bytes, packet sha `a380478f62aedecd2b8e7b74868cdb3acca578c36267b282862e1bf8920b1e45`. Exact match with `source-sizes.json`.
- Meta-idea variant: 71,934 bytes, packet sha `b224f24610106f17ebcb20c93558d40620f01d3ee4c56d16e84e8a754fd1388f`. Exact match. The request was phase 1, `deliberation`, `github-pr`, `--audience facilitator`, `--flag auto_implement`, `--optimize`.
- Both indexes classify `## Appendix B — ASD-STE100 output text` as `on-demand`, `included: false`, with the trigger text. `## 6. Conflict-avoidance mechanics` stays `always`.
- `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail` passes. The measured phase-1 body is 68,874 bytes under the 70,000-byte guard. The guard value in the test source is unchanged.
- `SKILL.md` is 19,482 bytes, under 20,000. `test/lean-organizer.test.js` passes: 8/8.

Result: not broken. The fixture guard was not widened. Verdict: CONFIRMED (PRIMARY).

### AC5 — package content, template, versions, hashes

Attempted: I searched the pack dry-run for the reference, re-hashed every
metadata target, and read the template change.

- `skill-pack-dry-run.json` lists `references/STE100.md`.
- `parley-deck/meta/version.json` hashes match my `shasum` results: protocol and packaged protocol `aae6e830…`, skill `64d5236b…`, compatibility manifest `34779a4b…`.
- `VERSION` and `internal/app/version.go` say `1.55.0`. `package.json`, `package-lock.json`, and `compatibility.json` say `2.19.0`; `recommendedCli` is `>=1.55.0`.
- `ARTIFACT_TEMPLATES.md` now says `Preserve verbatim source quotes, including non-English quotes.` The old translate instruction is gone. The template test asserts both.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### AC6 — staged core 2.19.0

Attempted: I hashed the staged base and result myself and checked the
reconstruction equality through the audited helper.

- `~/.parley/staging/COOPERATION-2.18.0.md` hashes to `9ead4475df6d6fa500317e7d347a0b6422bf3c202aca49f1d4b57a5fe14dc512`. This equals `core_base_sha256` in the payload and FINAL.
- `COOPERATION-2.19.0.md` hashes to `16323e04ffb076db3b2019169010033c6668199ebab0ca6896e79924e74bfe3b`. It equals base plus only H1 and H2. `protocol-transform.json` records base, payload, and result hashes.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

### AC7 — required tests and build pass

Attempted: I ran the focused suites myself, read every organizer log, and
compared CI results between main and this branch.

- Full CLI suite: `go-test-full.log` ends in `FAIL`, `internal/protocolpacket` fails. I reproduced the failure myself at HEAD. See MAJOR-1.
- Ubuntu CI: runs `38401898310` and `38401910409` on `4aab9f3` failed. Main run `37938864253` on `24fe38f` passed ubuntu. The new failure is `TestQuotaSignoffProcessSettlesBatchAndKeepsCurrentRequiredSet` with `fork/exec …/zcode: argument list too long`. See CRITICAL-1.
- Skill side: `skill-test-full-r2.log` records 402 node tests and 54 python tests, all pass, zero `✖`. The per-skill manifest checks pass. Package dry-run, portable build, installer dry-run, and doctor runs are recorded with exit 0. `skill-doctor-preinstall.json` shows the pre-install state of 2.18.0 runtimes; the status file discloses it, and the delivery phase owns the fix.
- My own runs: `go build ./...` passes; `go vet ./...` exits 0; both new Go test files pass.
- The organizer's `go-build.log` and `go-vet.log` are empty at freeze. Empty output is not a pass. My runs above cover both checks independently.

Result: BROKEN. The AC claim "Required CLI tests and build pass" is false at HEAD. Verdict: WRONG (PRIMARY).

### AC8 — out-of-scope exclusions

Attempted: I grepped both diffs for linter, dictionary, script, and checker
additions. I searched all Go code for `ste_output`. I compared the full
changed-file list with FINAL plus the two documented additions.

- `ste_output` appears in Go code only inside the instruction constant text at `internal/protocol/outputtext.go:6`. No parser reads it. No gate uses it.
- No new `.py`, `.sh`, lint, or dictionary file enters either diff.
- Changed files: all 23 CLI paths and all 12 skill paths are on the FINAL list, are idea artifacts, or are the documented additions (`parley-addon.json`, `test/installer.test.js`).
- No historical rewrite: each protocol copy changed by exactly the two hunks.

Result: not broken. Verdict: CONFIRMED (PRIMARY).

## Findings

### [CRITICAL] Linux headless launch regression: the full protocol now exceeds one exec argument

The protocol grew from 125,551 to 130,042 bytes. `prepareProtocolPrompt`
embeds the whole protocol into every launch prompt
(`internal/runner/protocol_context.go:111`; full context is the only enabled
treatment). For argv-mode agents, `buildAgentInvocation` then substitutes the
entire wrapped prompt into ONE argv element
(`internal/runner/runner.go:1313`). The zcode adapter uses
`HeadlessArgs: []string{"--prompt={prompt}", "--mode", "yolo", "--cwd", "{root}"}`
with `PromptMode: PromptArg` (`internal/agents/discover.go:413,417`).

Linux limits one argument string to 131,072 bytes (`MAX_ARG_STRLEN`). The
wrapped prompt is now about 130.6 KB plus the launch task, so almost every
real task crosses the limit. Before this change the same construction stayed
under it.

Evidence:

- Main `24fe38f`, run `37938864253`, ubuntu: success, same test passes.
- Branch `4aab9f3`, run `38401898310`, ubuntu: `--- FAIL: TestQuotaSignoffProcessSettlesBatchAndKeepsCurrentRequiredSet (0.13s)` with `fork/exec /tmp/…/zcode: argument list too long`, plus the block-count failure below.
- Branch `4aab9f3`, run `38401910409`, ubuntu: fail.
- I fetched all three results through `gh` myself.

Why it blocks: every headless launch of an argv-mode adapter fails on Linux
once the wrapped prompt exceeds the limit. `agy`, `hermes`, `kimi`, `opencode`,
and `zcode` all pass the prompt in argv (`internal/agents/discover.go:264,316,352,384,417`).
Rounds, reviews, signoffs, and goal checks all break on a shipped platform.
The required ubuntu CI gate is red.

Suggested fix: deliver large prompts to argv-mode agents through a prompt file
plus a short reference argument, or through stdin where the adapter accepts it.
Bound the fallback so no argv element can approach the kernel limit. The fix
needs a fix-up decision because FINAL froze the hunk bytes, not this runner
behavior.

Attribution: the CI pass-to-fail flip is PRIMARY. The kernel constant is
RECALL; I could not run Linux locally. The flip on its own proves the change
crossed a real limit.

### [MAJOR] Stale block-count constant fails the required CLI suite

`internal/protocolpacket/packet_test.go:107` reads `if rep.Blocks != 69 {`.
The live protocol now parses into 70 blocks. The test fails:

```text
live protocol parsed into 70 blocks; the map was written for 69 — reconcile both
```

I reproduced it at HEAD with
`go test ./internal/protocolpacket/ -run TestLiveMapIsCompleteHonorsNeverCutAndRendersEveryPhase -count=1`.
The full suite log `go-test-full.log` shows the same failure and ends in `FAIL`.

The map itself is correct: `parley protocol packet check` returns `"ok": true`
with 70 blocks, and the map's new Appendix B entry classifies the seventieth
block. Only the hardcoded count is stale.

Why it blocks: AC7 requires the CLI tests to pass. `go test ./... -count=1`
exits 1 at the reviewed HEAD. The organizer's `validation-status.json`
discloses this sole local failure and defers the fix to review disposition.

Suggested fix: change the constant to 70, or derive the expected count from the
map so the test cannot drift again.

### [MINOR] The canonical record omits the failing suite and the red ubuntu CI

`IMPLEMENTATION.md` `## Validation evidence` says the full checks are "in
progress". It does not record the suite failure or the red ubuntu CI. Those
facts live only in `release-delivery/2026-10-13-ste100/validation-status.json`,
outside the repository, and that file names only the local failure. It omits
the CI quota failure entirely.

AC7 requires disclosure of failures in the canonical record. Update
`## Validation evidence` with both failures and the CI links before the next
review round.

## Open questions

1. `probePrompt` (`internal/app/app.go:2273`) is model-facing. It asks for a
   sentinel file, not for protocol prose. I assess its exclusion from the
   thirteen as correct under FINAL's scope. The review consensus should concur
   or extend the list.
2. `source-sizes.json` records packet bytes and hashes without the request
   parameters. I had to recover them to reproduce the meta-idea hash. Record
   the full request next time.
3. The organizer's `go-build.log` and `go-vet.log` are empty placeholders at
   freeze. My independent runs pass both. The organizer should record the
   completed logs before close.

## Protocol context attestation

Verbatim from my launch:

```text
Protocol context attestation: {"context_mode":"full","source_sha256":"aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec","packet_sha256":"aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec"}
Shadow packet audit: {"would_fallback_reason":"unknown-phase:-1","packet_bytes":0,"source_bytes":130042,"included_blocks":27,"omitted_blocks":43}
```

Mode `full`; renderer reachable; no fallback. `Transport: github-pr`. I read
`meta/protocol-changelog.md`; it carries the 2.19.0 entry.

## Commands run

Exact commands I executed, all read-only except the test runs:

```text
python3 verify-source.py
shasum -a 256 <protocol copies, STE100.md, SKILL.md, protocol-hunks.json, staged core files>
git -C <worktree> diff --stat <base>..<HEAD>   (both worktrees)
git -C <worktree> diff <base>..<HEAD -- <paths> (both worktrees, full diff)
grep -rn "protocol.OutputTextInstruction" internal/runner internal/app
grep -rn "ste_output" internal/
go build ./...
go build -o /tmp/parley-ste100-review ./cmd/parley
go vet ./...
go test ./internal/protocol/ -count=1
go test ./internal/runner/ -run "TestOutputTextPolicyAcrossRunnerPrompts" -count=1
go test ./internal/app/ -run "TestOutputTextPolicyAcrossDriverPrompts" -count=1
go test ./internal/app/ -run "TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail|TestLiveDeckFacilitatorAcrossPhases|TestProtocolPacketCheckPassesOnTheFixtureAndFailsOnUnknownApplicability|TestProtocolPacketCheckPassesOnThisRepository" -count=1 -v
go test ./internal/protocolpacket/ -run TestLiveMapIsCompleteHonorsNeverCutAndRendersEveryPhase -count=1 -v
/tmp/parley-ste100-review protocol packet check --dir . --json
/tmp/parley-ste100-review protocol packet … --optimize --json   (variant sweep, phases 0–8, both audiences, flags)
node --test test/ste100.test.js
node --test test/lean-organizer.test.js
gh pr checks 84
gh run view 37938864253 / 38401910409 --repo feci/parley-deck-cli [--log-failed]
gh run list --repo feci/parley-deck-cli --workflow Tests --branch main --limit 5 --json …
```

## Verification limits

- I ran on macOS only. No local Linux or Windows. The E2BIG mechanism rests on
  CI logs, code reading, and byte arithmetic. The kernel constant is RECALL.
- I did not repeat the two full suites end to end. I read the organizer logs
  and ran focused suites. Both failures reproduce under my focused runs.
- Single full-suite runs on each side; I compared main against branch but ran
  no repeat-variance study.
- Delivery criteria DC1–DC6 are out of scope here. A separate fresh process
  verifies them after publication and installation.
- The staged core lives outside the repositories. I verified its bytes by hash
  and reconstruction, not by a release-channel check.

## Disposition

REQUEST CHANGES — fix-up cycle 1 required.

Open severity counts: CRITICAL 1, MAJOR 1, MINOR 1, NIT 0.
