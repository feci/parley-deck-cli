---
from: codex-1
to: user
idea: meta-protocol-change-ste100-output-format
phase: released
blocking: no
date: 2026-10-13
---

# ASD-STE100 output format release

CLI 1.55.0 and skill 2.19.0 are released through the authorized channels. Core 2.19.0 is staged for attended owner publication.
The working summary is the default for new English protocol prose. Exact source text and required syntax stay intact.
The summary does not certify full ASD conformance.

## Channel evidence

| Channel | Result | Evidence |
| --- | --- | --- |
| CLI source and GitHub | Merged and released 1.55.0 | [PR 84](https://github.com/feci/parley-deck-cli/pull/84); [release](https://github.com/feci/parley-deck-cli/releases/tag/v1.55.0) |
| Skill source and GitHub | Merged and released 2.19.0 | [PR 12](https://github.com/feci/parley-deck-skill/pull/12); [release](https://github.com/feci/parley-deck-skill/releases/tag/v2.19.0) |
| Homebrew | Both formulae published, installed, tested, and audited | [Tap commit](https://github.com/feci/homebrew-parley/commit/8ef3b97bcf97b9d710fa5a26305f010bb9c33e2c) |
| Skill WinGet | PR open; upstream installation and metadata checks remain pending | [PR 451758](https://github.com/microsoft/winget-pkgs/pull/451758) |
| CLI WinGet | Held | Windows baseline remains unresolved |
| Runtime skills | 114 SKILL.md hashes match across 15 managed and 4 generic targets | runtime-verification.json; all doctor checks pass |
| npm and global core | Pending owner commands below | Live registry and installed-core inspection |

CLI merge: `88dc6c5da139d2deb039bda4b605acfbe42663b4`. Skill merge: `f557c5e679b0c1a03f57a885f41ae54f68348b2a`.
Both merges retain the reviewed product bytes. The CLI binaries record the clean merge revision. Hosted asset hashes match the build records.
The skill tarball includes STE100.md. Every packaged file matches the reviewed source.
Skill tarball SHA256: `f07102c862908b3f6da6fb65f6a01f948fd3b83ba6fafa3d014c51145dce21e4`.
Core 2.19.0 equals staged 2.18.0 plus only H1/H2. Its SHA256 is `16323e04ffb076db3b2019169010033c6668199ebab0ca6896e79924e74bfe3b`.

## Independent acceptance and checks

Codex organized and implemented. Zcode supplied two full reviews, its own signoffs, and a fresh separate source goal check.
Kimi had two measured HTTP403 readiness failures. The brief authorized its exclusion. Claude took no part.
Fix-up cycle 1 resolved the Linux argument-limit regression, stale block count, and validation-record omissions.
The final review has no open CRITICAL or MAJOR. Both current participants accepted the zero-fix consensus. The fresh source goal check passed.
Current-tree independent evidence covers AC1–AC8. The mechanical single-reviewer exception was false; the brief supplied close authority.
Full CLI tests, build, and vet passed. Both Linux and both macOS product-commit CI jobs passed.
Skill checks passed 402 Node tests, 54 Python tests, six manifests, package checks, portable build, and installer preflight.
Both Windows jobs failed. Their failure names all occur in the two pre-change baseline runs. Zcode checked the terminal evidence independently.
Windows outcomes are recorded in windows-final-comparison.json. Failure-name comparisons do not prove equal causes or complete test coverage.
CLI Windows remains experimental. The skill Windows assets come from an Ubuntu cross-build, with no native Windows installation claim.
Runtime doctor used a task-local Python 3.14 PATH. The bidding add-on requires Python 3.10 or newer; system Python 3.9.6 is insufficient.

## Final independent channel result

Zcode returned CHANNEL-CHECK: PASS in a fresh process after release and installation.
Invocation: 54e2db04-8c86-4cc8-ad05-ff2f901b8ea5. Exit: 0. Duration: 744.227 seconds.
Its report is source-context/channel-verification-zcode-1.md in the idea.
It independently checked hosted bytes, formula archives, all 202 installed skill payload files, and all 114 runtime SKILL.md files.
The Homebrew launcher has its standard interpreter-path rewrite. The skill payload remains exact.

Transport was github-pr. Configured participants were Codex gpt-6-astra and Zcode zai/glm-5.3, both max/deep.
Configured model identity is not native observed-model evidence. Kimi's two readiness failures supplied no code-review vote.
The design used two rounds. Source review used two rounds and one fix-up cycle, followed by separate goal and Windows-evidence processes.
All agents share GitHub account feci. COMMENT mirrors link canonical participant-owned artifacts; they do not represent separate native approvals.

The read-only D6 CLI preview was repeated after channel verification. It still returns decisions: [] and the same preview hash.
The channel verifier disclosed that it read the preview record; it did not independently recompute the manifest digest.
The final owner command uses the fresh CLI output. No declaration was applied.

## Pending owner commands

Run D6 first in an attended terminal. Stop all Parley writers for this repository before the declaration.
Inspect the current preview. History remains unknown; this declaration preserves original bytes, charges, and caps.
The recorded preview covers 25 roots. Manifest SHA256: `7283348e38b86004c304e55befa95949d276d247cfd0166c4b98ed759333c6b6`. Preview SHA256: `d78dfd90744fc3eb43e86726ba4df99e0bf5b5f08a0109bc0d84b9a29da868f1`.
No agent applied the declaration. If the inspection digest changes, review the new preview and use its exact digest.

```sh
parley budget legacy inspect --dir '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100' --run parley-deck/runs/20260510T194003Z
parley budget legacy apply --dir '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100' --run parley-deck/runs/20260510T194003Z --expected-preview-sha256 d78dfd90744fc3eb43e86726ba4df99e0bf5b5f08a0109bc0d84b9a29da868f1 --decision-id driver-unstall-20260510-smoke-v1 --reason 'Record the 2026-05-10 smoke run as unknown history under the reviewed driver-unstall D6 repair; preserve original bytes, charges and cycle caps.' --writers-stopped --acknowledge-unknown-history --yes
```

Next, publish only the versions still absent from npm and the installed global core. The checked source hashes follow each version.
Live check: 2026-10-13T14:58:04.779363+00:00. npm latest: `2.13.0`. Installed cores: 2.10.0, 2.13.0.

### 2.15.0

Tarball SHA256: `64c274a86d0b4f2e136decbb1cc2810a18de89191192c6b9e680c5ec9af89b2a`.
Staged core SHA256: `0d81fd807114e4ff67f5fa98bba73622b58c86d66a2bf6b57d1def68cd099c09`.

```sh
npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-08-quota-auto-exclude/skill/parley-deck-skill-2.15.0.tgz' --access public
parley protocol publish --version 2.15.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.15.0.md
```

### 2.16.0

Tarball SHA256: `80f3172a29f293f26952475e218876bd67301f9021114f5a7f0d000643607db6`.
Staged core SHA256: `8c6b95b611901d3d1460ccb8554e26adac67f58d82ea028713791116a3e65a63`.

```sh
npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-participant-dropout/skill/parley-deck-skill-2.16.0.tgz' --access public
parley protocol publish --version 2.16.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.16.0.md
```

### 2.17.0

Tarball SHA256: `40b6eb424cfc003864f06fa17065ea1ac997637e47d76f0f86aef0fbcfbae3ca`.
Staged core SHA256: `6073c311b592653528bc675056b04b8454b491ab2702b37c298402ba2cb87957`.

```sh
npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-review-gate-timing/skill/parley-deck-skill-2.17.0.tgz' --access public
parley protocol publish --version 2.17.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.17.0.md
```

### 2.18.0

Tarball SHA256: `2d1e79016764d2e9e2ea5a4b5665c25da36ad7138ad15fecf0a93864f478eaf0`.
Staged core SHA256: `9ead4475df6d6fa500317e7d347a0b6422bf3c202aca49f1d4b57a5fe14dc512`.

```sh
npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-driver-unstall/skill/parley-deck-skill-2.18.0.tgz' --access public
parley protocol publish --version 2.18.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.18.0.md
```

### 2.19.0

Tarball SHA256: `f07102c862908b3f6da6fb65f6a01f948fd3b83ba6fafa3d014c51145dce21e4`.
Staged core SHA256: `16323e04ffb076db3b2019169010033c6668199ebab0ca6896e79924e74bfe3b`.

```sh
npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-13-ste100/skill/parley-deck-skill-2.19.0.tgz' --access public
parley protocol publish --version 2.19.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.19.0.md
```

## Deferred items and limits

- Native Windows CLI baseline repair and CLI WinGet publication remain deferred.
- Upstream review controls the skill WinGet merge. Native WinGet validation remains an upstream check.
- Full ASD vocabulary verification, a prose linter, and historical re-authoring are outside this release.
- The registry also lacks 2.14.0. This brief limits the owner command list to 2.15.0–2.19.0.
- Fixed-slug startup, D6 activation, and ready-consensus reopening exposed recorded driver gaps. No broad driver change was added.
- The original packet fixture is 68,874 bytes under 70,000. The meta-idea variant is 71,934 bytes and already exceeded the limit before this change.
- The skill core is 19,482 bytes under its unchanged 20,000-byte limit.
- Private-file delivery preserves exact prompt bytes. It does not prove that a model read them.

## Usage

Usage snapshot: 2026-10-13T14:58:04.905414+00:00. Counts include only measured invocations that name this idea.

- kimi-1: 2 child starts, 0 successes, 2 failed children, 0 control-plane refusals; 4.406 child seconds.
- zcode-1: 11 child starts, 11 successes, 0 failed children, 1 control-plane refusals; 5420.128 child seconds.
- Native fallback round02-zcode-native: exit 0; 520.406 seconds. This is separate from measured CLI invocations.

Codex cumulative snapshot at 2026-10-13T14:58:03.856Z: {'input_tokens': 39952862, 'cached_input_tokens': 36774912, 'cache_write_input_tokens': 0, 'output_tokens': 196354, 'reasoning_output_tokens': 68561, 'total_tokens': 40149216}.

Native participant tokens and costs are unavailable, not zero. No reliable dollar total is claimed.
Codex client snapshots are cumulative session counters, not a marginal task bill. Cached input and reasoning output are subsets.
See usage-final.json for exact cumulative counters and their timestamps. Do not add snapshots.

## Evidence location

/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-13-ste100

Source artifacts live in parley-deck/ideas/meta-protocol-change-ste100-output-format/.
Closed source artifacts and participant-owned text remain unchanged. Delivery evidence and this handoff remain separate.

## Audit and shared memory

Release audit PR: https://github.com/feci/parley-deck-cli/pull/85.
The audit contains source receipts, exact raw logs, hosted hashes, installer records, manifests, and usage evidence.
Raw logs retain their whitespace. WinGet manifests retain CRLF. No private prompt, credential, or runtime backup was published.
OpenViking saved the note at viking://resources/projects/parley-deck/ste100-output-format-release-20261013.md.
The write wait timed out, but exact readback and scoped filename retrieval confirmed persistence.
Semantic search did not return the note at the recorded check. Index processing remains pending.
The receipt is source-context/release/openviking-verification.json. No duplicate write was made.
