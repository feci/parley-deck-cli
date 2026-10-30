---
from: codex-1
to: user
idea: meta-protocol-change-credit-exhaustion-failover
phase: implementation
blocking: no
date: 2026-10-30
---

CLI **1.56.0** and skill **2.20.0** are released through the authorized channels below.
Core **2.20.0** is staged for owner publication.
All 114 runtime `SKILL.md` files match the released payload.
The independent goal check and the separate release-channel check both passed.

Recognized provider exhaustion stops that identity without retry, probe, or reconnection in the current idea.
Role succession follows the owner's order. Transient failures retain the existing two-attempt path.
Durable records protect restart admission, predecessor artifacts, filed dissent, and independent review.
The usable-participant floor and the adapter evidence limits still apply.

Both implementation PRs merged with two-parent merge commits. The tags point to those exact commits.
The clean CLI build contains only the reviewed product and later idea audit records.
All 212 skill package files match the merged source.

| Channel | Result | Evidence |
| --- | --- | --- |
| CLI GitHub release | Published, `v1.56.0` | [Release](https://github.com/feci/parley-deck-cli/releases/tag/v1.56.0); [PR 87](https://github.com/feci/parley-deck-cli/pull/87); merge `59e78377ed00dadd940104f271bb13b8fce5d1a1` |
| Skill GitHub release | Published, `v2.20.0` | [Release](https://github.com/feci/parley-deck-skill/releases/tag/v2.20.0); [PR 13](https://github.com/feci/parley-deck-skill/pull/13); merge `7cdfb241b65a3eea855a9ed61e860d75acdd2f4a` |
| Both Homebrew formulae | Published, installed, tested, and audited | Tap commit `25230804f6161e3a01ecc3e9e5fa75393d98fa5c`; [receipt](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/homebrew-verification.json) |
| Managed runtime installs | 15 targets, including undetected targets | Six skills per target; 90 hashes match |
| Generic runtime installs | Four required destinations | Six skills per destination; 24 hashes match |
| Skill WinGet | PR submitted; upstream work pending | [PR 465809](https://github.com/microsoft/winget-pkgs/pull/465809); head `daf9930c4d6653c6d503c4acbdfc6d410ee59f15` |
| CLI WinGet | Held while CLI Windows is experimental | No CLI PR submitted |
| npm and global core | Owner-only; pending for six versions | Exact commands below; [verified input hashes](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/pending-owner-commands.json) |

The default `/opt/homebrew/bin/parley` reports `parley 1.56.0`.
The default `/opt/homebrew/bin/parley-deck-skill` reports `2.20.0`.
Both Homebrew tests and `brew audit --strict` passed.
The installed skill subtree matches the released package bytes.
Homebrew built its CLI from the tagged source archive. The GitHub binaries have separate clean-build evidence.
The report makes no claim that these binary bytes are identical.

The Homebrew CLI source archive SHA256 is `24b4c2ee613c6aaed23516ffbbbf0c35074a2485061e9e42d72cca3ac05a3514`.
The Homebrew skill source archive SHA256 is `004f134b5726aca710fe7e8c1aeed3a5869538673fb779eeb943b80077920aae`.
[The formula receipt](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/homebrew-publication.json) records both source URLs and pins.

Zcode downloaded all 11 hosted assets in a separate process.
Each actual byte hash matched its live GitHub digest and the organizer's record.
[The asset receipt](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/hosted-assets-verification.json) also records URLs and byte counts.

| Hosted asset | SHA256 |
| --- | --- |
| `parley-v1.56.0-darwin-arm64` | `03eaed39456e88e17a91d50ed9aafabeca9245739be2c6d08aa974d011e1a584` |
| `parley-v1.56.0-darwin-x64` | `511bf0d2c34262dc535cf62758cb962f59166e64df170be0c2456d0182cbac32` |
| `parley-v1.56.0-linux-arm64` | `f411d19eb7d208b25f603ed7e4f52e7c095e7ec8804c71abc50fc2e3d0e7ffa4` |
| `parley-v1.56.0-linux-x64` | `3fe7d87217609365cbe9a131ca8999f25a44bc92382b506a99022d2ee462f755` |
| `parley-v1.56.0-windows-arm64.exe` | `8f636f5c1374825664e7064887a9cbe1882adcf9f36153dcee02c7f09be15f25` |
| `parley-v1.56.0-windows-x64.exe` | `a73c9ed944dd54c2e1bd0183acb1a3fa7bb9afd1483537c25b67495e317bb945` |
| `sha256.json` | `cb597998b33d15f6d492488cbc04260f837ca6e9e0e723166829e777c845c791` |
| `parley-deck-skill-2.20.0.tgz` | `4f75dd45fc1c5c9960f027093f422d825e5beb1d9e2fb5c4c8708afe169116f0` |
| `parley-deck-skill-v2.20.0-macos-arm64` | `0478b8a44a294b79cbb7e0b306b60cf1d2523ad081fc27ca58eae9b3aa72e2fe` |
| `parley-deck-skill-v2.20.0-windows-arm64.exe` | `b93dd523f1f6d4b75ae5fb40378fe2a6dd7401f5c54e0d963e14a21888d70a69` |
| `parley-deck-skill-v2.20.0-windows-x64.exe` | `239a34fd11dfc8b27f8f5a7572c9c8f236664ad30e47de673147663006bc21d3` |

The staged core `~/.parley/staging/COOPERATION-2.20.0.md` has SHA256 `6c437ed5fa0805bc8929b1aa25371664f9c4df1fc34c731a9df02bf02f2bd357`.
It derives from staged core 2.19.0 plus the exact reviewed 29-hunk set.
The hunk-set SHA256 is `f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c`.
The core skill `SKILL.md` is 19,989 bytes; its SHA256 is `33154bead53b8808a057fe5194c669acfcbac1a587ad38ac618d600f2c9c0ad6`.
The phase packets measure 68,448 and 69,880 bytes. The size limits remain unchanged.

The installer-status enumeration plus the four generic destinations is the **runtime inventory of record**.
The status command includes all managed targets: `parley-deck-skill status --target all --include-undetected --json`.
[The direct runtime receipt](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/runtime-verification.json) records 114 unique paths and 114 matching hashes.
All doctor checks passed. The generic installations used the required `install --target generic --dest <dir> --force` commands.

- `~/.hermes/profiles/ldx/skills/parley-deck`
- `~/.hermes/profiles/librade/skills/parley-deck`
- `~/.hermes/profiles/testprofile/skills/parley-deck`
- `~/.config/opencode/skills/parley-deck`

Each destination contains all six sibling skills.
The 12 Gemini files under `plugins` and `extensions` are part of the direct 114 checks.
The requested `*/skills/parley-*/SKILL.md` pattern cannot match those files.

The exact whole-home `find` timed out after 180 seconds.
The narrower scan excluded two historical cloud-trash paths and completed in 37.42 seconds.
It found 132 files: 102 runtime files and 30 source/archive files.
All 30 historical files retain their original hashes.
[The copy-inventory receipt](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/runtime-inventory-after-hashes.json) preserves the scan results and limits.

`OneDrive-LabourDigitalGmbH/.Trash` exists and remains outside the completed inventory.
`OneDrive-NordlichtITSolutions,s.r.o/.Trash` was historically excluded as unavailable; it was absent at final verification.
The old receipts remain unchanged. This report does not claim complete whole-home coverage.
These statements apply the reviewer's two requested record corrections.
The reviewer-owned NITs remain in the review; this report claims no withdrawal.

A private backup preserves runtime files from before installation. Its mode is 0600.
The backup SHA256 is `1b4b82813427a68d550e11b8ef74d266bd4a09bfb8e8625ddf6fd4fba85269d7`.
Sixty destination files were absent before installation.
The first install-helper attempt stopped on an absent status field before backup or runtime mutation.
The corrected helper derived sibling destinations from installer source. Zcode independently verified that mapping and all resulting hashes.

The skill [portable workflow](https://github.com/feci/parley-deck-skill/actions/runs/40980088483) completed successfully.
Its `macos-tests` job ran on macOS. Its `windows` job ran on Ubuntu and cross-built Windows binaries.
This evidence does not establish native Windows execution.
The final hosted Windows bytes passed hash and PE-architecture checks after that workflow completed.

The final remote snapshot is dated `2026-10-30T00:59:25.104Z`.
Skill WinGet PR 465809 remains OPEN and unmerged in that snapshot.
Its three manifests passed schema 1.12.0, CRLF, package/version, hosted-hash, and architecture checks.
The CLA check passed. Upstream supplied no completed validation verdict at the final check.
The catalog still has 2.19.0; this report does not claim 2.20.0 catalog availability.
Native WinGet validation was unavailable on this Mac.

The [reviewed CLI CI run](https://github.com/feci/parley-deck-cli/actions/runs/40974357405) passed Ubuntu and macOS; Windows failed.
CLI Windows remains experimental. CLI WinGet remains held.
The [merged-main CI run](https://github.com/feci/parley-deck-cli/actions/runs/40979932958) was still in progress at the final snapshot.
This report does not count it as passed.
Six superseded runs covered only later audit changes. They were canceled after merge to free release-workflow capacity.
[The cancellation record](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/superseded-ci-cancellations.json) preserves their identities and source-diff checks.
No canceled run counts as passed.

The required `go test ./... -count=1 -timeout 45m` passed all 34 packages at reviewed product `b7dce40`.
Its log SHA256 is `68cd0dedd2609cf803e8e24d8ad741a09e48d2193d84f911f90e58864e6172d5`.
Concurrent-budget and grouped-accounting checks each passed 100 repetitions.
Relevant race checks, `go vet`, and `gofmt` passed.
The skill passed 402 Node tests, 54 Python tests, package checks, and portable builds.
[Current-tree validation evidence](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/fixup02-validation.json) preserves the receipts and their attribution.

Parley used GitHub PR transport and separate native CLI processes for Kimi and Zcode.
The recorded participant defaults were Codex `gpt-6-astra/xhigh`, Kimi `omniroute/k3.1-pro/max`, and Zcode `zai/glm-5.3/max`.
These labels describe configuration; they do not prove the provider's observed model.
[The selection checkpoint](../ideas/meta-protocol-change-credit-exhaustion-failover/00-prompt.md) records the original launch details.
Three design rounds and three full-scope review rounds completed. Two fix-up cycles resolved F1 through F9.
All three participants signed plain ACCEPT on the zero-fix [review consensus](../ideas/meta-protocol-change-credit-exhaustion-failover/review/consensus.md).
Its SHA256 is `f64b7be182a3c663575318b7d9587ff1eaf88ece2a0ec59a85c6189f509e17cf`.

Zcode's fresh goal process completed at `2026-10-30T00:27:03.272436+00:00` with exit 0 and `GOAL-CHECK: PASS`.
Its separate channel process completed at `2026-10-30T00:52:08.703011+00:00` with exit 0 and `RELEASE-CHANNEL-CHECK: PASS`.
Neither process timed out.
The [canonical review](../ideas/meta-protocol-change-credit-exhaustion-failover/review/round-03/zcode-1.md) owns both verdicts and their exact scope.
Its final SHA256 is `cc75643357f01326da0c0e1336e5d04f4ad2d425fd35b759f0a369ef3f09bd34`.
The [goal receipt](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/goalcheck-zcode-1-native.receipt.json) and [channel receipt](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/channel-review-zcode-1-native.receipt.json) preserve process evidence.
Exact-prefix checks prove that each append preserved the earlier review bytes.
Codex mirrored the final Zcode review as COMMENT reviews on both merged implementation PRs.
The shared `feci` account supplies no distinct GitHub approval identities.
[The mirror receipt](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/channel-review-pr-mirrors.json) records both review URLs and the exact body hash.

The host reboot interrupted two earlier peer processes. Those interruptions count as neither participant failures nor provider failures.
The [reboot record](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/resume3-host-reboot.json) preserves partial logs and invalid private scratch-file observations.
All three participants remained active through close.
The owner authorized native CLI fallback for the recorded driver refusal, `has no runs yet`.
Fallback receipts supply no driver accounting or automatic membership authority.
No worktree was pruned or declared. Credentials and gateway settings were unchanged.

The frozen FINAL and closed IMPLEMENTATION retain their release-merge bytes.
[Their hash receipt](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/frozen-artifact-hashes.json) proves both comparisons.
The old delivery index remains an unchanged pre-check snapshot; its pending label records that earlier time.
[The final delivery summary](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/final-summary.json) records the completed channel check and the current limits.
Raw process logs, session content, and the runtime backup remain private.

Owner-only publication remains pending for every version from 2.15.0 through 2.20.0.
The npm registry still ends at 2.13.0. Installed global cores remain 2.10.0 and 2.13.0.
All six tarballs and all six staged cores exist and match their recorded hashes.
The organizer executed no command below and did not modify the relay's owner-publication script.
The controlling brief reserves these actions for the owner's OTP and TTY gates.

Run these commands in the stated order: npm, then core, for each version.
The relay retains ownership of the attended publication script.

```sh
npm publish "/Users/tomasfecko/AI_WORKSPACE/parley-deck/release-delivery/2026-10-08-quota-auto-exclude/skill/parley-deck-skill-2.15.0.tgz" --access public
parley protocol publish --version 2.15.0 --from "/Users/tomasfecko/.parley/staging/COOPERATION-2.15.0.md"

npm publish "/Users/tomasfecko/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-participant-dropout/skill/parley-deck-skill-2.16.0.tgz" --access public
parley protocol publish --version 2.16.0 --from "/Users/tomasfecko/.parley/staging/COOPERATION-2.16.0.md"

npm publish "/Users/tomasfecko/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-review-gate-timing/skill/parley-deck-skill-2.17.0.tgz" --access public
parley protocol publish --version 2.17.0 --from "/Users/tomasfecko/.parley/staging/COOPERATION-2.17.0.md"

npm publish "/Users/tomasfecko/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-driver-unstall/skill/parley-deck-skill-2.18.0.tgz" --access public
parley protocol publish --version 2.18.0 --from "/Users/tomasfecko/.parley/staging/COOPERATION-2.18.0.md"

npm publish "/Users/tomasfecko/AI_WORKSPACE/parley-deck/release-delivery/2026-10-13-ste100/skill/parley-deck-skill-2.19.0.tgz" --access public
parley protocol publish --version 2.19.0 --from "/Users/tomasfecko/.parley/staging/COOPERATION-2.19.0.md"

npm publish "/Users/tomasfecko/AI_WORKSPACE/parley-deck/release-delivery/2026-10-29-credit-failover/skill/parley-deck-skill-2.20.0.tgz" --access public
parley protocol publish --version 2.20.0 --from "/Users/tomasfecko/.parley/staging/COOPERATION-2.20.0.md"
```

Deferred scope remains explicit:

- Remaining work in `quota-kickoff-reporting-and-alias-guidance`; this release includes the required alias guidance.
- Native-positive provenance and completeness in `quota-zcode-native-exhaustion-capture`; this release includes diagnostic capture.
- Shared-account inference and a machine exhaustion cache.
- Additional native adapters and broader Codex versions. Codex recognition requires explicit JSON and exact version `0.165.0`.
- Zcode's bounded provenance limits and SDK-internal retries before observable terminal exhaustion.
- CLI WinGet, pending owner-only publication, and external skill WinGet validation and merge.

Usage was collected at `2026-10-30T00:56:51.468747+00:00`.
[The usage summary](../ideas/meta-protocol-change-credit-exhaustion-failover/source-context/release/usage-summary.json) deduplicates receipt copies and preserves the measurement limits.
Native completed processes and measured child starts are separate populations.

| Agent | Native completed | Native success | Native timeout | Native ms | Measured starts | Measured success / nonzero exit | Measured ms | Control refusals |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| codex-1 | 0 | 0 | 0 | 0 | 2 | 2 / 0 | 27,365 | 0 |
| kimi-1 | 10 | 9 | 1 | 10,733,082 | 4 | 2 / 2 | 1,343,540 | 2 |
| zcode-1 | 15 | 15 | 0 | 9,435,870 | 2 | 2 / 0 | 989,498 | 1 |

The two measured Kimi nonzero exits preceded the reboot. They occurred at about 18:17 and 19:04 UTC.
Two additional reboot-interrupted processes have unknown terminal duration and usage. They are excluded from the failure totals.
The ordinary Codex JSON smoke was separate: 15,644 input tokens, including 10,368 cached tokens, and 10 output tokens.
Native token and cost totals are unavailable; they are not zero.

The organizer's cumulative session counters at `2026-10-30T00:56:33.765Z` were:

| Counter | Tokens |
| --- | ---: |
| Input | 143,917,155 |
| Cached input, included in input | 140,958,720 |
| Output | 580,460 |
| Reasoning output, included in output | 255,940 |
| Total input plus output | 144,497,615 |

These are cumulative session counters, not a marginal bill for this resume.
Snapshots and nested subsets must not be added. No reliable dollar total exists.
