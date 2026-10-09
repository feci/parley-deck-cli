---
from: codex-1
to: user
idea: meta-protocol-change-participant-dropout-review-gate-timing
phase: release
blocking: no
date: 2026-10-09
---

# Participant dropout review gates and timing — release delivery

CLI **1.53.0** and skill **2.17.0** are merged and released on GitHub. Both Homebrew formulae are published, upgraded and tested. All 15 managed targets plus the four requested generic locations passed doctor checks, with **114/114 SKILL.md hashes** matching the release payload. The skill-only WinGet PR is submitted. Core **2.17.0** is staged from staged 2.16.0 plus exactly the 14 reviewed substitutions. npm and global-core publication remain owner-only, with the still-missing commands below.

**AC14 delivery PASS**, independently verified by a separate zcode-1 process, PID24149, actual exit 0 after **658.682 s**. Its own [channel report](../ideas/meta-protocol-change-participant-dropout-review-gate-timing/source-context/channel-verification-zcode-1.md) has SHA256 `add849aefc6352e095b9488698c8f72458714cf020147ae1fd6fafa78d730269`. It verified live asset digests, downloaded files, build/source identities, formulae, all 114 installed skill hashes, the live WinGet PR and exact core reconstruction. Actual process window: `2026-10-09T09:33:28.032346+00:00` to `2026-10-09T09:44:26.723505+00:00`; the process receipt and GitHub API timestamps govern over minor date shorthand in the report. Owner npm/core publication and upstream WinGet merge are excluded from this delivery verdict by the brief.

## Behavior delivered

A validated recorded automatic dropout can reduce the independent-reviewer minimum from two to one when it is the cause of that reduction. The remaining non-implementer reviewer must have a known configured snapshot model distinct from the implementer's, including when the diversity flag is false. The shared causal predicate applies to prospective membership changes and the review, consensus, goal and close gates. Missing, manual-only, malformed or unproven history does not grant the exception. A fresh process of the remaining reviewer can perform the goal check. Strict findings, signatures, reservations, retained dissent, protected roles, the floor and current-tree evidence remain binding.

Eligible streaming participant steps default to a 120-second first-output watchdog, 300 seconds of stalled activity and a 60-second heartbeat. Headless signoffs now use the supervisor and the existing original-attempt-plus-one-retry path, with a five-second retry delay at the same hard ceiling. Valid BLOCK output remains dissent. Zcode and default Claude text adapters are declared buffered, so lack of streaming output alone cannot falsely trigger their soft watchdog; hard ceilings still apply. **The product goal-check ceiling remains 120 seconds.**

Compaction keeps the existing size limits and applicability map unchanged: phase-1 packet **68,846/70,000 bytes**, core SKILL.md **18,990/20,000 bytes**. All 14 protocol substitutions were reviewed and verified across the three protocol copies and staged-core transformation.

## Close and evidence

This run used the controlling override: codex-1 organized and implemented; zcode-1 independently reviewed in separate native processes. Kimi was excluded under explicit owner authority after two genuine HTTP403 readiness failures. No Claude task process participated, and no credential, model, gateway or roster setting was changed. The controlling brief is copied byte-identically at [source-context/ORGANIZER-BRIEF.md](../ideas/meta-protocol-change-participant-dropout-review-gate-timing/source-context/ORGANIZER-BRIEF.md), SHA256 `50cb660916b882965ff289dc2fc7d5ff9c8f234cac61d2e95e1e5c8849c8722a`.

This idea has no qualifying automatic history/snapshot; it closed under the brief's **attended authority**, not the new product exception. Owner intent: “participanti nie su nevyhnutne potrebny obaja, staci jeden a to by nemalo zaseknut parley-deck” — one participant reviewer is enough and should not stall Parley.

- Design PR [#77](https://github.com/feci/parley-deck-cli/pull/77), two design rounds, both design ACCEPTs, and a frozen FINAL preceded product implementation.
- One signed fix-up cycle was used, out of the maximum five. Final independent full-scope review has no CRITICAL, MAJOR or MINOR and zero agreed fixes; the two nonblocking NIT dispositions are recorded in final consensus/goal evidence.
- Final Zcode review process PID40945 exited 0 in 1772.599 s. Fresh independent goal process PID8391 exited 0 in 518.314 s with AC1–AC13 PASS; both current participants authored final review-consensus ACCEPT blocks. These are actual attended process results, not claims that the driver's 120-second auto-goal parser executed.
- Close commit `cf1c639` froze FINAL, IMPLEMENTATION, reviews, goal evidence and signoffs. Their nine recorded hashes remain unchanged. AC14 is the separate post-close delivery recorded here.
- Reviewed product commits: CLI `b89e2abaa056813a4b238c2fa4278195b0f7096b`; skill `74cc831b18ce33e48e705fc8992c89be85cadf6f`. Later CLI changes through release are this idea's evidence and merge metadata.
- Full current-source Go suite `go test ./... -json -count=1 -timeout 45m`: exit 0 in **1388.005 s**, **34 passing packages, 3330 passing test/subtest events, four built-in skips, zero failures**. Log `.parley-runtime/review-gate-timing/full-host-tests-cycle1.jsonl`, SHA256 `04d2f0529c7e38046ac64c7be04603fbdd05dd0cc461b796d5d8c65f27724db5`. Build, vet and focused lifecycle/causal/protocol checks passed.
- Full skill suite at reviewed 74cc831: **399 Node tests, 54 Python tests, six manifests**, exit 0. Log `skill-tests-v3.log`, SHA256 `3509b5ce4c7db449fb6f72ab9055215ce6ad7d0c29bbe367ff83a0aa0cc35a4b`. Failed earlier runs remain retained and are not counted as passes. Final cycle receipts were atomically written, synced and read back identically in native and shared locations.

The [Parley Deck skill](/Users/tomasfecko/.codex/skills/parley-deck/SKILL.md) governed the run and its GitHub review mirrors. Native reviews use COMMENT with attribution because participant identities share the PR author's GitHub account. Canonical participant artifacts and their own signatures remain authoritative. Driver run/continue/status/wait/consensus were used; recorded D6 pre-dispatch accounting and phase/reopen gaps required the configured native fallbacks. No accounting record or dropout event was fabricated.

## Release channels

Evidence root: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-review-gate-timing`. A compact evidence index is committed as [source-context/release-delivery.json](../ideas/meta-protocol-change-participant-dropout-review-gate-timing/source-context/release-delivery.json).

| Channel | Delivered state | Evidence |
| --- | --- | --- |
| CLI GitHub | [v1.53.0](https://github.com/feci/parley-deck-cli/releases/tag/v1.53.0); [PR #78](https://github.com/feci/parley-deck-cli/pull/78) merged at `128e30b479065a98437844886d295ae83e9b0e52`; six binaries plus sha256.json | `cli-build-verification.json`, `hosted-asset-verification.json` |
| Skill GitHub | [v2.17.0](https://github.com/feci/parley-deck-skill/releases/tag/v2.17.0); [PR #10](https://github.com/feci/parley-deck-skill/pull/10) merged at `8ce4dec0fde0f563f5110e99d8b7e4550d939770`; npm tarball and macOS/Windows portable assets | `skill/payload-verification.json`, `hosted-asset-verification.json` |
| Homebrew | Both formulae at tap commit `6964b1188eb0ed81d07f251a4ca8b3852f43dad6`; style, strict online audit, upgrade and formula tests passed | `homebrew-source-archives.json`, `homebrew-installed-verification.json`, `logs/homebrew-*.log` |
| Runtime skills | 15 managed + generic Hermes ldx/librade/testprofile + generic OpenCode; 19 healthy targets, all 114 SKILL.md hashes match | `runtime-verification.json`, doctor logs |
| Skill WinGet | [PR #449419](https://github.com/microsoft/winget-pkgs/pull/449419), head `2afe423e66299c133ae5c6b9ae9097c8dbb696a1`; exactly three added manifests with final live Windows installer hashes | `winget-submission.json`, `winget-validation.json`, `channel-status-latest.json` |
| Core | Staged `/Users/tomasfecko/.parley/staging/COOPERATION-2.17.0.md`; global publication reserved for owner | `core-staging.json` |
| npm | Exact 2.17.0 tarball prepared and GitHub-hosted; registry publication reserved for owner | `owner-actions.json` |

All **11 live GitHub assets** were downloaded and SHA256-verified against hosted digests; the six CLI binaries also match the published checksum manifest. Each CLI binary carries the clean merged VCS revision and `vcs.modified=false`. All 210 packed skill files match both reviewed and merged source. Skill portable workflow [37910828042](https://github.com/feci/parley-deck-skill/actions/runs/37910828042) passed macOS and Windows tests and built/uploaded both Windows assets.

Active installed commands report `parley 1.53.0` and `parley-deck-skill 2.17.0`. The first unqualified Homebrew upgrade failed because two taps provide the same formula name; selecting `feci/parley/parley-deck-cli` and `feci/parley/parley-deck-skill` resolved it. No tap trust setting was changed. The private pre-install backup remains mode0600 at `runtime-skills-before-2.17.0.tar.gz`.

Skill tarball SHA256: `40b6eb424cfc003864f06fa17065ea1ac997637e47d76f0f86aef0fbcfbae3ca`. Core SKILL.md: `0b9769e8a58fd83cdd9d1d92be5f0d755f3fe7cbba30be20126d5eea9cd6af0c`. Staged core 2.17.0: **124,976 bytes**, SHA256 `6073c311b592653528bc675056b04b8454b491ab2702b37c298402ba2cb87957`; reproduced from staged 2.16.0 SHA256 `8c6b95b611901d3d1460ccb8554e26adac67f58d82ea028713791116a3e65a63` plus exactly 14 reviewed substitutions. Hunk-file SHA256 `54d0dd07d747cf54ce9d0e234b9d65c34a8282b525bd8667090aaf19cbffcc39`; generic zones preserved.

At `2026-10-09T09:46:43.981503+00:00`, merged-main CLI [run 37910668038](https://github.com/feci/parley-deck-cli/actions/runs/37910668038) records go build & test (ubuntu-latest): success; go build & test (macos-latest): success; go build & test (windows-latest): in_progress. This is **not an all-CI-green claim**. Baseline Windows failures were independently observed in prior main runs; the still-pending result here is not assigned a cause or counted as passed. Four superseded self-triggered deck-only runs were cancelled only after checking zero product diff from b89e2ab; cancellations remain nonpasses (`superseded-ci-cancellations.json`).

WinGet PR #449419 is **OPEN** at that snapshot: 09. Installer Metadata Validation: QUEUED; 10. Validation Completed: QUEUED; 08. Installation Validation: IN_PROGRESS; 01. Pull Request Validation: SUCCESS; 02. Manifest Validation: SUCCESS; 03. URLs Validation: SUCCESS; 04. URL Domain Validation: SUCCESS; 05. Manifest Policy Validation: SUCCESS; 06. Catalog Content Verification: SUCCESS; 07. Installers Scan: SUCCESS; license/cla: SUCCESS. No upstream merge or completed native installation is claimed. Automated skipped auxiliary checks are not presented as passing checks.

## Owner-only publication commands

Rechecked at `2026-10-09T09:48:02.329050+00:00`: npm latest is **2.13.0**; installed global cores are **2.10.0, 2.13.0**. Only still-unpublished commands are listed, in required version order. The `!` prefix is the owner's execution convention. Codex executed none of these commands and bypassed no OTP/TTY gate. Artifact paths and hashes are recorded in `owner-actions.json`.

```text
! npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-08-quota-auto-exclude/skill/parley-deck-skill-2.15.0.tgz' --access public
! parley protocol publish --version 2.15.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.15.0.md

! npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-participant-dropout/skill/parley-deck-skill-2.16.0.tgz' --access public
! parley protocol publish --version 2.16.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.16.0.md

! npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-review-gate-timing/skill/parley-deck-skill-2.17.0.tgz' --access public
! parley protocol publish --version 2.17.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.17.0.md
```

## Limits and deferred work

- **Real goal duration versus product bound:** the healthy attended goal check took 518.314 s while the product bound remains 120 s. Timeout tuning is a separate follow-up; this release does not silently raise it.
- The optional duplicate-snapshot-identity branch fixture remains deferred to the next membership-test touch, date TBD, as accepted by independent review. The underlying malformed/unknown snapshot behavior remains fail-closed. The final receipt-record timing-overlap NIT is resolved by committed evidence.
- Investigate the shared-volume ENOENT receipt anomaly if it recurs. Final cycle dual receipts are readable; earlier recovered log evidence and failure history are retained. No infrastructure failure is rewritten as a passing run.
- Windows CLI remains experimental and CLI WinGet remains held. Native WinGet validation/installation was unavailable on this Mac; only actual upstream checks are reported. Local supervisor fixtures do not establish real-provider latency or cross-host/PID correctness.
- D6 legacy accounting, native-positive legacy quota recognition, alias/manual-edit guidance, persistent pre-idea proposal/resume/abandon lifecycle and inherited portability hygiene remain separate follow-ups. Prior dissent and shared integrity gates were not relaxed to get a passing close.
- Both owner worktrees remain allocated on `review-gate-timing`; neither was pruned or declared. Private invocation logs, failed test history, backup and delivery evidence are retained.

## Usage

Capture `2026-10-09T09:46:44.064689+00:00`; exact source records are retained in `usage-final.json` and the committed evidence index. The idea-filtered measured ledger records **2 successful Zcode starts**, 1045.875 s combined, and **2 failed Kimi starts**, 3.075 s combined (both HTTP403). Configured native fallbacks add **7 successful Zcode processes**, 5630.519 s, plus **one deliberately cancelled process**, PID20815, exit143, 40.358 s. Thus there were **10 actual Zcode task processes: 9 successful and one cancelled**. The cancellation corrected an erroneous prompt premise; it was neither a timeout nor a provider/dropout event, and no output from it was accepted. D6 control-plane refusals started no task child and are not billed as child processes here. Readiness probes without this idea's identity are excluded.

Configured participant models were Zcode `zai/glm-5.3` and Kimi `kimi-code/k3`. No own/Zcode provider timeout, 429 or 503 required the brief's delayed relaunch policy. The separate shared-memory indexing wait timed out but direct readback succeeded; that is not a participant provider failure. Native token/cost fields are unavailable, not zero. Codex client telemetry is reported separately:

| Session | Snapshot UTC | Input | Cached input | Output | Reasoning output | Total |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| `01a11f65-b2d6-73d0-b316-d9f190327573` | 2026-10-09T09:46:13.338Z | 51,848,445 | 49,713,664 | 203,638 | 60,870 | 52,052,083 |

These are cumulative client-session counters, not a marginal task bill. Cached input is a subset of input and reasoning output is a subset of output; snapshots are not summed. No reliable total monetary cost is claimed. Final documentation and publication operations follow this telemetry snapshot.

Shared memory `viking://resources/projects/parley-deck/review-gate-timing-release-20261009.md` was written and read back byte-identically. The server’s wait-for-index operation timed out at 60 seconds and the first scoped semantic search did not yet return this URI, so indexing remains unconfirmed. The exact note and verification are retained as `shared-memory-note.md` and `shared-memory-verification.json`.
