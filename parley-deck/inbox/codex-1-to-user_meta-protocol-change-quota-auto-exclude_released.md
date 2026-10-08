---
from: codex-1
to: user
idea: meta-protocol-change-quota-auto-exclude
phase: release
blocking: no
date: 2026-10-08
---

# Released: quota auto-exclusion

CLI **1.51.0** and skill **2.15.0** are released to GitHub and Homebrew. The skill-only
WinGet PR is submitted. All 15 local managed skill installations are updated and
verified. npm and global core publication are staged for the two attended commands below.

The implementation was closed under finish-now point 5 and the binding round08-answer:
no open CRITICAL or MAJOR, both final review-consensus blocks ACCEPT, and current
AC1–AC21 evidence recorded. AC2 is NOT MET / owner-waived. AC5 and AC15 are NOT MET
in full / accepted-deferred. The three round-08 residuals remain release limitations.
No sixth code cycle was opened.

## Published and submitted channels

| Channel | Version and evidence | State |
| --- | --- | --- |
| CLI source/GitHub | [v1.51.0](https://github.com/feci/parley-deck-cli/releases/tag/v1.51.0), commit `343aab21a7d95083767260c01c66d5e69f557cd6`; six platform binaries and `sha256.json` | Published; main/tag pushed |
| Skill source/GitHub | [v2.15.0](https://github.com/feci/parley-deck-skill/releases/tag/v2.15.0), commit `352a475c4697d43e7582d00cc789bb6462867ee5`; npm tarball, macOS ARM64 portable, Windows x64/ARM64 portable | Published; main/tag pushed |
| Homebrew | Both formulae at [320a1d7](https://github.com/feci/homebrew-parley/commit/320a1d70ee59b9eb47a1c37d2f38360eba313809) | Published; active local CLI 1.51.0 and installer 2.15.0 |
| WinGet skill | [microsoft/winget-pkgs#448514](https://github.com/microsoft/winget-pkgs/pull/448514), exactly three 2.15.0 manifest files | PR submitted; latest upstream state recorded in evidence |
| WinGet CLI | Windows remains experimental with known durable-operation failures | Held |
| Local skills | codex, claude, agy, gemini, hermes, qwen, codebuddy, goose, kimi, droid, vibe, cursor, opencode, aionrs, zcode | All 15 marker versions 2.15.0; identical core SKILL.md hashes; doctor passes |
| npm | Exact 2.15.0 tarball staged and smoke-tested; registry latest observed as 2.13.0 | Owner publication pending |
| Global protocol core | `~/.parley/staging/COOPERATION-2.15.0.md` | Owner publication pending |

All 11 release assets match live GitHub SHA256 digests. CLI binaries report release
commit `343aab2` with `vcs.modified=false`. The final skill Windows EXE bytes came
from successful [portable workflow 37706962158](https://github.com/feci/parley-deck-skill/actions/runs/37706962158)
and supply the WinGet hashes. Local native skill and CLI smoke tests pass. Homebrew
style, strict online audit, version/help tests, and a real skill installation plus
doctor pass. Three WinGet 1.12.0 schema checks pass. Native Windows `winget validate`
and installation were not run on this macOS host; the upstream pipeline is authoritative.

Local core SKILL.md SHA256 for each of 15 runtimes:
`76765df83428a23288d5d1d09f0fc7b40d53197a832c0f3f51e11a3a3687a599`.
The exact packed payload was installed with `--include-undetected` to cover the
seven existing managed installations without currently detected executables.
The prior payload backup is retained privately in the delivery directory.

## Two owner-only commands

npm tarball SHA256:
`64c274a86d0b4f2e136decbb1cc2810a18de89191192c6b9e680c5ec9af89b2a`.

```sh
! npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-08-quota-auto-exclude/skill/parley-deck-skill-2.15.0.tgz' --access public
```

Staged core SHA256:
`0d81fd807114e4ff67f5fa98bba73622b58c86d66a2bf6b57d1def68cd099c09`.
Core 2.14.0 was never locally published. The staged 2.15.0 reproduces staged 2.14.0
plus exactly this idea's reviewed embedded-protocol patch, preserving the generic
template zones. Published local core versions remain 2.10.0 and 2.13.0.

```sh
! parley protocol publish --version 2.15.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.15.0.md
```

These commands retain the owner's npm OTP and core TTY gates.

## Accepted and deferred limitations

1. **Kickoff blocking escalation (R8-MINOR-1).** If the kickoff quota decision
   would leave fewer than two usable participants or excludes a protected role,
   an absent or unwritable `parley-deck/inbox/` prevents the blocking note from
   being saved. The command still stops without applying the exclusion or
   creating an idea, but prints only the file error rather than the candidates
   and arithmetic. Keep a writable inbox directory, including in a fresh clone,
   until this is fixed in the follow-up.
2. **Kickoff notice crash window (R8-NIT-1).** A process crash after writing the
   kickoff run manifest and before publishing its notice can leave the notice
   permanently absent. The kickoff record, marker, status and organizer brief
   still show the exclusion. This case was identified by source review; no crash
   was injected. Mid-idea receipt/replay behavior is unchanged.
3. **Aliased decks and plain edits (R8-NIT-2).** On a symlinked deck, a plain
   `participants:` edit or confirmed exclusion counts as a manual revision and
   requires a physical deck path even when `quota_auto_exclude` is false.
   Until that path is restored, the pending edit blocks all signers and driving.
   Disabling the policy again does not resolve this case; the generic diagnostic's
   suggestion to do so is inapplicable.

These three findings are accepted/deferred with AC5/AC15 exceptions under the binding
round08-answer, never fixed or passed. Follow-up: [quota-kickoff-reporting-and-alias-guidance](../ideas/quota-kickoff-reporting-and-alias-guidance/00-prompt.md), an inactive candidate, not a sixth cycle.

Native-positive AC2 remains owner-waived and unmet. Real zcode native exhaustion
may not be recognized; the human-confirmed path remains the fallback. Follow-up:
[quota-zcode-native-exhaustion-capture](../ideas/quota-zcode-native-exhaustion-capture/00-prompt.md).
No private stderr capture is copied or committed. Other adapters remain diagnostic-only;
provenance ambiguity and cross-host/container/PID-namespace recovery limits remain explicit.

Windows CLI cannot create new ideas because durable directory sync returns
`Access is denied`, including with the policy off. Scoped driving, signing, manual
imports, revisions and transitions are also affected. Windows CLI CI fails, assets
are experimental and CLI WinGet is held. The fix is routed to the separate, unmerged
`windows-portability` work. D6 legacy run-accounting
remains separately deferred.

## Independent delivery verification and usage

The separate claude-1 [channel report](../ideas/meta-protocol-change-quota-auto-exclude/release-channel-verification-claude-1.md)
concludes **agent-controlled delivery is complete; no channel blocker**. It records
0 CRITICAL, 0 MAJOR, 0 MINOR and 2 pre-existing NITs. Report SHA256:
`9eb887a90c9e48b9f2428a8cf1e2f80dff0b6d153ad7dc9d6eebb46a1bd9be98`.

The initial verifier timed out at 1800.4 seconds without a report. The authorized
2400-second retry completed in 982.7 seconds, exit 0, with its own fresh commands.
There were no provider retries. Plain `claude-opus-5-5[1m]`, max effort, was used.
Only codex-1 and claude-1 participated; codex-1 organized/implemented, and separately
launched Claude owned the reviews and channel report. The declared deck transport
is github-pr; the owner authorized local CLI work and direct main merges, with the
skill-only WinGet PR as the distribution contribution.

The verifier independently re-downloaded/re-hashed all assets, repeated both
Homebrew tests/audits, compared all 15 core trees and 75 add-on trees with released
payloads, ran npm publish in dry-run mode, and reproduced the staged core and its
base chain. Exact-release CLI CI [37706766781](https://github.com/feci/parley-deck-cli/actions/runs/37706766781)
is **failure overall: Linux and macOS pass, Windows fails**. No all-platform green
claim is made. The skill workflow cross-builds Windows EXEs on Linux; native Windows
execution remains pending in WinGet step 08. WinGet steps 01–07 and CLA pass;
08 is in progress, 09–10 queued, and the PR is OPEN at this finalization snapshot.

Pre-existing Homebrew nits remain visible: the local `feci/release-validation-20260918`
tap causes ambiguity for unqualified formula names, so use
`feci/parley/parley-deck-cli` and `feci/parley/parley-deck-skill`. Also, installs made
by the Homebrew skill installer omit core LICENSE/README from runtime directories
because Homebrew moves them to the keg root; 2.14.0 behaved the same way. These are
outside this release's diff and do not block delivery. The 15 installations delivered
here used the exact npm payload and include those files. No unrelated cleanup or
sixth implementation cycle was opened.

[Channel evidence](../ideas/meta-protocol-change-quota-auto-exclude/source-context/release-20261008/README.md)
contains release refs, asset and runtime hashes, formula checks, WinGet state,
reviewer launch/exit and report identity. Complete local delivery files are retained
at `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-08-quota-auto-exclude`.

Usage uses the latest cumulative snapshot per `(source, source_path)`, never the sum
of repeated snapshots. Failed and successful verifier attempts have distinct
sources. Cutoff: 2026-10-08T01:15:33.362077+00:00.

| Client | Sources | Events | Input | Cached input | Cache writes | Output | Client total |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| claude-jsonl | 29 | 3685 | 7,398 | 782,866,240 | 33,011,263 | 5,082,259 | 5,089,657 |
| codex-rollout | 16 | 2073 | 305,314,060 | 288,085,248 | 0 | 1,314,050 | 306,628,110 |

Codex input includes cached input, and reasoning output is a subset of output.
Claude cache reads/writes remain separate from its client total. These client
accounting conventions are not interchangeable. Attribution is ambiguous, so these
are transcript counters, not exclusive marginal costs for this release. No monetary
estimate is asserted. [Ledger](../ideas/meta-protocol-change-quota-auto-exclude/usage-ledger.jsonl)
and [deduplicated source summary](../ideas/meta-protocol-change-quota-auto-exclude/source-context/release-20261008/usage.json)
retain the audit trail.
