---
from: codex-1
to: user
idea: meta-protocol-change-driver-unstall
phase: complete
blocking: no
date: 2026-10-09
---

# Driver unstall — released channels and owner handoff

CLI 1.54.0 and skill 2.18.0 are published on GitHub and Homebrew, and all requested skill destinations are installed and verified. The separate fresh Zcode channel process exited 0 with **CHANNEL-VERIFICATION: PASS** after 702.956 seconds (invocation 8b01e946-b799-4d0b-b576-5874d73f731c). Its own artifact is [channel-verification-zcode-1.md](../ideas/meta-protocol-change-driver-unstall/source-context/channel-verification-zcode-1.md); the matching measured receipt is retained beside it.

The one-time real-repository D6 declaration and npm/global-core publication remain **pending owner acts**, as the controlling brief and signed design require. No terminal attendance was emulated and no real legacy declaration was applied.

## Delivered channels

| Channel | Verified delivery |
| --- | --- |
| CLI GitHub | [v1.54.0](https://github.com/feci/parley-deck-cli/releases/tag/v1.54.0), six binaries plus checksum map. All hosted hashes match the clean build; VCS revision is bed5ffd1049dc0f207eb39beaa3dd3ba879181c3 and vcs.modified=false. |
| Skill GitHub | [v2.18.0](https://github.com/feci/parley-deck-skill/releases/tag/v2.18.0), npm-format tarball, macOS arm64 portable and Windows x64/arm64 portable assets. All 11 CLI/skill hosted asset digests were verified. GitHub publication does not imply npm publication. |
| Homebrew | [bd2a364](https://github.com/feci/homebrew-parley/commit/bd2a3640a9d217d4c3f48d7121c6669690ed61f0) bumps only the two formula URLs/hashes. Both installed versions, formula tests and style checks passed: parley 1.54.0 and skill 2.18.0. Old installations were preserved. |
| Runtime skills | Actual installer used for 15 managed targets, including undetected targets, plus the four required generic locations. All **114 SKILL.md hashes across 19 targets** match the six release payload hashes; doctor passed and the private pre-install backup is retained locally. |
| Skill WinGet | [microsoft/winget-pkgs#449541](https://github.com/microsoft/winget-pkgs/pull/449541), current state **OPEN**. Exactly three Feci.ParleyDeckSkill 2.18.0 manifests, schemas/CRLF/package/version/live installer hashes verified. This is a skill-only submission; CLI WinGet stays held. Current upstream check states are preserved in winget-live-status.json. |
| Global core | 2.18.0 staged from staged 2.17.0, SHA256 9ead4475df6d6fa500317e7d347a0b6422bf3c202aca49f1d4b57a5fe14dc512. Publication remains an attended owner action. |

Generic destinations are ~/.hermes/profiles/ldx/skills/parley-deck, ~/.hermes/profiles/librade/skills/parley-deck, ~/.hermes/profiles/testprofile/skills/parley-deck and ~/.config/opencode/skills/parley-deck. The five companion skills are installed alongside parley-deck under each target's skills directory.

The [skill portable workflow](https://github.com/feci/parley-deck-skill/actions/runs/37930155652) passed macOS tests, then Ubuntu tests and cross-built/uploaded both Windows installers. Its job named “windows” runs on Ubuntu. No local native Windows execution, winget validate or winget install result is claimed from this Mac. An open WinGet PR is not yet a distributed package.

## Source completion and verification

[Design PR80](https://github.com/feci/parley-deck-cli/pull/80), [CLI implementation PR81](https://github.com/feci/parley-deck-cli/pull/81), and [skill PR11](https://github.com/feci/parley-deck-skill/pull/11) are merged. The implementation merges are bed5ffd1049dc0f207eb39beaa3dd3ba879181c3 and 054954458dbfe2d764bf6234a6d6bbeee12e2494, each with two parents. Product/docs/protocol bytes equal reviewed CLI 134ac40; the skill merge tree equals reviewed e46e551. Later CLI changes are this idea's audit records only.

Scope is exactly D1 and D2. D1 adds an attended, durable declaration for a strictly eligible legacy manifest, with unknown-history disclosure, immutable adoption, copy validation and preserved charges/caps. D2 derives the goal ceiling from the minimum of the track's 5/15/30 minutes and a positive configured checker timeout, using the existing original-plus-one retry path with a frozen bound across restart. Semantic FAIL is final, failed-process PASS cannot close, and independent completion/membership duties remain intact.

One fix-up cycle out of five was used: two documentation clarifications and a stronger timeout-test oracle. Three independent review rounds found no remaining CRITICAL/MAJOR; both current participants authored fresh ACCEPTs on the zero-fix review consensus. The fresh independent source goal process passed AC1–AC8 and AC9 checks/drift, including real May bytes in a two-worktree fixture and a real 121.16-second goal process. The full Go suite passed (34 packages), build/vet and focused races passed, macOS/Linux CI passed at 134ac40, and the skill suite passed (399 Node tests, 54 Python tests, six manifests). The merged main macOS/Linux rerun also passed; its exact snapshot is retained.

Native CLI Windows CI remains **unresolved/experimental**. Both final product checks failed with 356/355 direct failing names versus the 356-name baseline, with zero added names and the same 45-minute timeout. The later merged-main Windows run also failed, with 355 direct names and zero added names versus baseline, plus the same 45-minute timeout (main-windows-comparison.json). Matching names prove neither equal causes nor per-test execution. CLI WinGet stays held. Four superseded audit-only CI reruns were explicitly cancelled after product-byte identity verification to free macOS capacity; they are recorded as CANCELLED, never as passes. The original product checks and main-branch run were retained.

The mechanical single-reviewer exception returned allowed:false. Source close used only the controlling brief's expressly pre-authorized path: final Zcode review without open CRITICAL/MAJOR, both own ACCEPTs, separate fresh goal PASS and independent current-tree AC evidence. Kimi had two real readiness HTTP403 connection-allowlist failures and was excluded for this idea only; credentials and global roster were unchanged.

The source-close IMPLEMENTATION explicitly left live AC9 evidence pending. This separate post-publication channel PASS completes that delivery condition. FINAL, IMPLEMENTATION and both consensus records remain byte-identical to the source merge; no closed artifact was rewritten. The old three-fix consensus was archived byte-identically after explicit reviewer concurrence because the CLI refused reopening triage=ready; the product reopen gap remains deferred.

## Owner action: one-time D6 activation

The release is installed, but the real legacy declaration has not been applied. The owner must stop all Parley writers for this repository, inspect the current preview, then execute the exact declaration in an attended terminal. History remains unknown, never zero; no budget charge or cap is reset.

The read-only preview covers 24 matching visible copies. Manifest SHA256: `7283348e38b86004c304e55befa95949d276d247cfd0166c4b98ed759333c6b6`. Preview SHA256: `d7d6d89007e9ecfe3d03f235c37d1aa861b40b3c90c107a1fdf3181adb4c5b1b`. Original May bytes: `ee4f52b717159e54aa8f144161d14467160cfc6846aed97abbdd6768683c6906`. No `legacy-*` authority exists in the git-common budget store.

```sh
parley budget legacy inspect --dir '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/driver-unstall' --run parley-deck/runs/20260510T194003Z
parley budget legacy apply --dir '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/driver-unstall' --run parley-deck/runs/20260510T194003Z --expected-preview-sha256 d7d6d89007e9ecfe3d03f235c37d1aa861b40b3c90c107a1fdf3181adb4c5b1b --decision-id driver-unstall-20260510-smoke-v1 --reason 'Record the 2026-05-10 smoke run as unknown history under the reviewed driver-unstall D6 repair; preserve original bytes, charges and cycle caps.' --writers-stopped --acknowledge-unknown-history --yes
```

If the inspection digest differs, review the new preview and use its exact digest; do not reuse this command blindly. This command does not declare, prune or remove a worktree.

## Owner actions: npm and global core, in version order

The fresh check at 2026-10-09T13:37:18.917665+00:00 found **8 still-unpublished commands**. All source tarball and staged-core hashes were re-verified. Only still-missing publications are listed below.

Run these in your attended terminal. Only versions still absent at the recorded check time are listed. Neither npm nor the global core was published as part of this run.

### 2.15.0

npm tarball SHA256: `64c274a86d0b4f2e136decbb1cc2810a18de89191192c6b9e680c5ec9af89b2a`

Staged core SHA256: `0d81fd807114e4ff67f5fa98bba73622b58c86d66a2bf6b57d1def68cd099c09`

```sh
npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-08-quota-auto-exclude/skill/parley-deck-skill-2.15.0.tgz' --access public
parley protocol publish --version 2.15.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.15.0.md
```

### 2.16.0

npm tarball SHA256: `80f3172a29f293f26952475e218876bd67301f9021114f5a7f0d000643607db6`

Staged core SHA256: `8c6b95b611901d3d1460ccb8554e26adac67f58d82ea028713791116a3e65a63`

```sh
npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-participant-dropout/skill/parley-deck-skill-2.16.0.tgz' --access public
parley protocol publish --version 2.16.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.16.0.md
```

### 2.17.0

npm tarball SHA256: `40b6eb424cfc003864f06fa17065ea1ac997637e47d76f0f86aef0fbcfbae3ca`

Staged core SHA256: `6073c311b592653528bc675056b04b8454b491ab2702b37c298402ba2cb87957`

```sh
npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-review-gate-timing/skill/parley-deck-skill-2.17.0.tgz' --access public
parley protocol publish --version 2.17.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.17.0.md
```

### 2.18.0

npm tarball SHA256: `2d1e79016764d2e9e2ea5a4b5665c25da36ad7138ad15fecf0a93864f478eaf0`

Staged core SHA256: `9ead4475df6d6fa500317e7d347a0b6422bf3c202aca49f1d4b57a5fe14dc512`

```sh
npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-driver-unstall/skill/parley-deck-skill-2.18.0.tgz' --access public
parley protocol publish --version 2.18.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.18.0.md
```

Checked at 2026-10-09T13:37:18.917665+00:00. npm latest: `2.13.0`; installed cores: 2.10.0, 2.13.0.

Informational observation outside the requested 2.15–2.18 publication list: the verifier also found npm 2.14.0 absent. No 2.14.0 publication is instructed or changed by this handoff.

## Deferred work

- Auxiliary phase-pointer/run resolution and missing planner actions (FINAL ALT8); choosing an older run alone can reuse stale membership or budgets.
- Checked ready/partial consensus reopen tied to a completed fix cycle and new review (ALT9). This idea's transparent byte-preserving archive is a recorded fallback, not new product semantics.
- Existing-kickoff/fixed-slug startup, missing-round launch, seeded snapshots, placeholder-before-exit validation, generic briefs, wider reporting and unrelated quota reporting (ALT10).
- Broad native Windows baseline repair, including internal/evidence build failure, the pre-existing dropout-test hang and filesystem/process/fixture portability. Name-set comparisons do not diagnose causes.
- Owner recovery after both durable goal attempts are exhausted, and a future explicitly authorized cycle-scoped goal identity. No reset or third attempt was added.

Follow-up slugs remain TBD; none was silently implemented or promised in this release.

## Workflow and usage

Timing clarification: the peer artifact's approximate “15:35Z start” heading does not match the measured receipt. The actual measured process ran from 13:22:41.150979 UTC to 13:34:24.080463 UTC; the receipt is authoritative. The peer artifact remains unchanged.

The first manual channel-check invocation timed out after 1800.101 seconds while observing a fresh main-branch Windows CI run and produced no canonical verdict. Its receipt is retained as a failure. The organizer honored the controlling brief's 15-minute wait before the authorized 2400-second retry; no result from the failed process was counted as a pass. A known orphaned read-only polling shell was terminated after timeout, and the previously identified task PIDs were verified absent. This manual channel-process incident does not establish product RunParticipantStep behavior or alter the shipped goal retry policy.

Transport was github-pr. codex-1 organized and implemented this idea under the explicit brief override; that does not change the global Claude-organizer default. Codex CLI 0.161.0 / configured gpt-6-astra and Zcode CLI 3.7.7-13 / configured zai/glm-5.3 used the recorded max/deep defaults. Native observed model/token/cost telemetry was unavailable. Zcode authored its own rounds, reviews, signoffs, fresh goal and separate channel record. Claude did not participate.

Two design rounds, three implementation review rounds, one fix-up cycle, a fresh independent goal, and two post-publication channel attempts were recorded. The first channel attempt timed out; the authorized fresh retry passed. Additional signoff and host-mirror processes have their actual receipts retained.

| Measured agent | Child starts | Successful children | Failed children | Control-plane refusals | Child seconds |
| --- | ---: | ---: | ---: | ---: | ---: |
| kimi-1 | 2 | 0 | 2 | 0 | 3.577 |
| zcode-1 | 12 | 11 | 1 | 1 | 7276.724 |

One additional configured native Zcode round-02 process completed separately: exit 0, **280.608 seconds**, from 10:11:36.211684 UTC to 10:16:16.823739 UTC. It is not included in the measured-process table above; its full receipt is retained in usage-final.json.

Codex client cumulative snapshot at 2026-10-09T13:34:34.033Z: input 52,550,149 (including 50,431,616 cached input), output 294,173 (including 137,587 reasoning output), total 52,844,322. These are session counters, **not a task bill**; cached input and reasoning output are subsets, not additional totals.

Native participant token and cost fields are **unavailable, not zero**. The usage-final.json client snapshots are cumulative Codex session telemetry, not a marginal task bill; snapshots must not be added together. No reliable total dollar cost is claimed.

## Audit and shared memory

Compact channel/owner/usage receipts are under [source-context/release/](../ideas/meta-protocol-change-driver-unstall/source-context/release/), including an evidence hash index and retained-log hash index. Full logs, release assets, temporary verification environment and the private runtime backup remain at /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-driver-unstall. No backup payload or credentials were uploaded.

Shared OpenViking was unavailable through the configured MCPAnywhere connection (HTTP400: No valid session ID provided). No successful write is claimed; the concise intended memory note is retained locally with URI viking://resources/projects/parley-deck/driver-unstall-2026-10-09.md. No gateway, credential or direct-connection workaround was attempted.

The worktrees and branches were preserved, with no pruning or worktree declarations. The four local orchestration runs remain retained; no execution history was invented.
