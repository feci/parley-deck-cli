---
from: codex-1
to: user
idea: meta-protocol-change-participant-dropout
phase: release
blocking: no
date: 2026-10-09
---

# Participant dropout — release delivery

CLI **1.52.0** and skill **2.16.0** are merged and released. Both Homebrew formulae are updated and installed. All 15 managed targets plus the four requested generic locations passed doctor checks; all **114 SKILL.md hashes** match the packed release payload. Core **2.16.0** is staged from staged 2.15.0 plus the exact reviewed hunks. npm and global core publication remain owner-only actions below; the WinGet deliverable is the submitted skill PR.

**AC14 delivery PASS**, independently verified by a separate zcode-1 process, invocation `fbf2d597-50e2-4306-9a94-52d4648f8b88`, exit 0 after 571.909 s. Its own [channel report](../ideas/meta-protocol-change-participant-dropout/source-context/channel-verification-zcode-1.md) has SHA256 `e5cca084619f3510208beef4a316e9a440819ec1fd2805e1a5e00bc1f756c051`. It checked all 11 live asset digests, representative downloads, both Homebrew source archives and an independent strict online audit, the live WinGet PR, every installed SKILL.md against the tarball, and exact core reconstruction. No material delivery inconsistency or suppressed finding was reported. Owner npm/core publication and upstream WinGet merge are excluded from this verdict by the brief.

## Shipped behavior

New ideas use the participant-failure trigger through the existing automatic-exclusion mechanism and historical `quota_auto_exclude` setting. An eligible dispatched participant step gets its original attempt plus one retry after five seconds at the same ceiling. Two genuine failures can permanently remove that participant from this idea; the next idea probes it again. Durable idea/agent/logical-step identity prevents a third attempt after restart or changed run/input. Invalid own output can qualify; valid disagreement, cancellation, control-plane refusal and shared-file tampering cannot.

The organizer, implementer/designee and protected draft authors remain protected. A batch requires two positively usable non-organizer participants, including the implementer. Prior vetoes, disputes and findings remain binding. Existing reviewer, diversity, strict and goal-check gates still apply: **auto_implement 3→2 still blocks with only one independent reviewer**. This release supplies no automatic waiver. Saved legacy policies remain quota-only unless explicitly revised; per-idea false remains an opt-out.

## Close and independent evidence

The owner-appointed roles were codex-1 organizer/implementer and zcode-1 independent reviewer. Kimi was excluded only from this idea after two genuine round-01 process failures under the controlling brief. No Claude deliberation, implementation or review process participated. The original brief is preserved byte-identically at `../ideas/meta-protocol-change-participant-dropout/source-context/ORGANIZER-BRIEF.md` (SHA256 `54448074f5a4dac2587ee3ed9c4e08866d84e50d42c02df7dc9b1e6af70651b5`).

The GitHub-PR workflow completed two design rounds, two implementation-review rounds, one signed fix-up cycle of the maximum five, final ACCEPT blocks from both current participants, and a fresh independent goal check. Zcode owns every Zcode artifact and its own signoff blocks. Native review mirrors use COMMENT because the shared GitHub login cannot approve its own PR. The Parley Deck skill governed this workflow: `/Users/tomasfecko/.codex/skills/parley-deck/SKILL.md`.

- Reviewed product: CLI `cb78e9f43e70cb59284a21ba6e0e8a331cb447a1`; skill `efe296c7acf13a147ab820ce6cbf8e6705b68691`. Subsequent changes through release are deck evidence or merge metadata.
- Independent final review: `../ideas/meta-protocol-change-participant-dropout/review/round-02/zcode-1.md`, measured invocation `b569b516-c1f3-4bae-a81f-f1e20b50cfb9`, exit 0. It found zero agreed fixes and no open CRITICAL/MAJOR.
- Fresh goal check: `../ideas/meta-protocol-change-participant-dropout/goal-check-zcode-1.md`, invocation `eddf724d-2365-4cd3-a09f-e7494525026f`, exit 0, AC1–AC13 PASS. This is an attended assessment, not a claim that the driver's auto-close parser executed.
- Both final review-consensus ACCEPT blocks exist; close commit `ef1061ce0b73f237dab6d01834061e53cafaad4a` sets IMPLEMENTATION complete. AC14 was expressly post-close delivery, completed by this release record. Closed FINAL, IMPLEMENTATION, reviews and signoffs remain byte-identical.
- Qualifying full native suite: `go test ./... -json -count=1 -timeout 45m`, exit 0 in 731.058 s; 34 passing packages, three without tests, 3,262 passing test/subtest events and four built-in skips. The log is `.parley-runtime/participant-dropout-launches/full-host-tests-cycle1-local.jsonl`, SHA256 `131f8ef7c53971e1cedcfa0402db7c118cd7244fbd764425e22c26deeb37e503`, with its matching `-local-result.json`. Earlier shared-mount and disk-full runs remain failed evidence and are not counted as passes.
- Focused dropout/goal/readiness regressions, go vet/build and the skill suite passed. Skill: 399 Node tests, 54 Python tests and all six manifests; Zcode independently reran the skill checks and adversarial product probes.

## Release channels

Evidence root: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-participant-dropout`.

| Channel | Delivered state | Evidence |
| --- | --- | --- |
| CLI GitHub | [v1.52.0](https://github.com/feci/parley-deck-cli/releases/tag/v1.52.0), [PR #76](https://github.com/feci/parley-deck-cli/pull/76) merged at `2f4c9afc0c52f78ac1b23f733a494aa4385b1695`; six binaries plus sha256.json | `cli-build-verification.json`, `cli-assets/sha256.json` |
| Skill GitHub | [v2.16.0](https://github.com/feci/parley-deck-skill/releases/tag/v2.16.0), [PR #9](https://github.com/feci/parley-deck-skill/pull/9) merged at `dbdb91972ab8e3a0883aac3e8898fa8c118de108`; npm tarball and macOS/Windows portable assets | `skill/payload-verification.json`, `github-release-verification.json` |
| Homebrew | Both formulae published at tap commit `69b962d9701c633fb700a24fa6051f8ef517161b`; style, strict online audit, upgrade and brew test all passed | `homebrew-verification.json`, `homebrew-checks.json` |
| Runtime skills | 15 managed + generic Hermes ldx/librade/testprofile + generic OpenCode; 19 targets validated by doctor and 114 matching SKILL.md hashes | `runtime-verification.json` |
| Skill WinGet | [PR #449194](https://github.com/microsoft/winget-pkgs/pull/449194), head `265b8878acb3ab6921dee092ec91cab9943a1c19`; three schema-validated manifests with exact hosted asset digests | `winget-validation.json`, `winget-checks.json` |
| Core | Staged `/Users/tomasfecko/.parley/staging/COOPERATION-2.16.0.md`; not globally published | `core-staging.json`, `core-candidate-verification.json` |
| npm | Exact 2.16.0 tarball prepared and GitHub-hosted; registry publication reserved for owner | `owner-actions.json` |

All 11 CLI/skill GitHub assets were downloaded and hash-verified. The CLI binaries were built from the clean merged source, carrying its exact VCS revision with `vcs.modified=false`. Hosted skill portable workflow [37861735927](https://github.com/feci/parley-deck-skill/actions/runs/37861735927) passed. Active commands report `parley 1.52.0` and skill `2.16.0`.

The 2.16.0 npm tarball SHA256 is `80f3172a29f293f26952475e218876bd67301f9021114f5a7f0d000643607db6`; core SKILL.md is `37e36b5bbff6b7e4a440663b103eb191a48ba7ad99654429f0af4fdb9042cca8`. The private pre-install backup remains at `runtime-skills-before-2.16.0.tar.gz` with mode 0600.

Core 2.16.0 SHA256 is `8c6b95b611901d3d1460ccb8554e26adac67f58d82ea028713791116a3e65a63`. It reproduces from staged 2.15.0 SHA256 `0d81fd807114e4ff67f5fa98bba73622b58c86d66a2bf6b57d1def68cd099c09` plus 35 exact sequential reviewed substitutions (P1–P6 and required compaction), preserving generic zones. Hunk file SHA256: `055e7fd6e3e5f04205aebf16d72d04967716de947be697b2ab9bbda8f7b297e9`.

CI snapshot at `2026-10-09T00:14:48.165128+00:00`: candidate push run 37858412333 and PR run 37858416744 both finished with Linux/macOS success and Windows failure. The merged-main run 37861593668 has Linux/macOS success; Windows is still in progress at this snapshot. This is not an all-CI-green claim. Candidate Windows logs include directory-sync Access is denied failures in quota/dropout paths and other inherited Windows limitations; `logs/cli-candidate-windows-failed.log` retains the evidence. Six superseded self-triggered deck-only metadata runs were cancelled after verifying zero product diff; they are cancellations, not passes (`superseded-ci-cancellations.json`).

WinGet PR #449194 remains OPEN. At `2026-10-09T00:14:48.165128+00:00`, upstream pull-request, manifest, URL, domain, policy, catalog, installer-scan and CLA checks passed; installation validation was in progress and metadata/completion checks queued. No upstream merge or completed native installation is claimed.

## Owner-only publication commands

At verification, npm latest is 2.13.0 and neither 2.15.0 nor 2.16.0 is published. Locally published global cores remain 2.10.0 and 2.13.0. Run the staged 2.15.0 commands first, then 2.16.0. The `!` prefix below is the owner's execution convention. These commands were not executed by codex-1 and no OTP/TTY gate was bypassed.

```text
! npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-08-quota-auto-exclude/skill/parley-deck-skill-2.15.0.tgz' --access public
! parley protocol publish --version 2.15.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.15.0.md

! npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-participant-dropout/skill/parley-deck-skill-2.16.0.tgz' --access public
! parley protocol publish --version 2.16.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.16.0.md
```

The prior 2.15.0 tarball SHA256 is `64c274a86d0b4f2e136decbb1cc2810a18de89191192c6b9e680c5ec9af89b2a`. All four commands and artifact digests are in `owner-actions.json`.

## Limits and deferred work

- The precommit reviewer-count limitation above remains explicit. A separate `participant-dropout-review-gate-timing` idea could consider changing that contract; no change is silently included here.
- Z1 (former MAJOR) and Z3 (MINOR) are fixed and independently verified. Z2 was explicitly withdrawn by its reviewer against FINAL D2: pre-idea attempts are bounded within one proposed readiness batch; stable ProbeID replay never replenishes them, while a distinct new proposal probes afresh. After creation, the per-idea cap survives run/input changes. A persistent pre-idea proposal/resume/abandon lifecycle remains TBD.
- Z4 is an accepted maintenance limit: phase-1 packet 69,963/70,000 bytes and core SKILL.md 19,995/20,000 bytes. The next additive change requires compaction; no size limits were raised.
- Windows CLI remains experimental and CLI WinGet remains held. Native WinGet validate/install was unavailable on this Mac; upstream PR checks are the separate Windows evidence.
- Real provider timing and cross-host/PID behavior are not established by local supervisor fixtures. D6 legacy accounting, alias/plain-edit guidance, legacy native-positive quota recognition and inherited fixture portability hygiene remain separate follow-ups.
- The driver could not carry every attended phase: recorded D6 pre-dispatch refusal required one configured native Zcode round-02 fallback; the phase-pointer and ready-consensus reopen limitations are retained in organizer notes. No control-plane refusal was counted as participant dropout, no accounting was fabricated, and no roster/model was changed.
- Both owner worktrees remain allocated on `participant-dropout`; neither was pruned or declared. Private invocation logs, failed test history and delivery evidence are retained.

## Usage

Capture: `2026-10-09T00:16:15.410709+00:00`; exact records and source paths are retained in delivery `usage-final.json`. The idea-filtered measured ledger has **8 successful Zcode child starts**, 4714.030 s combined child duration, and **2 failed Kimi starts**, 3.245 s. A separate Zcode round-02 configured native fallback succeeded in 350.476 s, giving **nine actual Zcode task processes in total**. One additional round-02 control-plane budget refusal started no child and is excluded from the process count. Readiness probes without this idea identity are not included.

Configured participant models were native Zcode `zai/glm-5.3` and Kimi `kimi-code/k3`; measured launches selected the `deep` profile and the recorded 1,800 s default ceiling. No participant timeout needed a retry in the completed run. Native Zcode/Kimi token and monetary cost fields are unavailable; unavailable does not mean zero. Codex's client session telemetry is reported separately:

| Session | Last snapshot UTC | Input | Cached input | Output | Reasoning output | Total |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| `01a11d58-cc42-7d01-8edf-66796ef69bb1` | 2026-10-08T21:08:44.318Z | 15,198 | 9,856 | 6 | 0 | 15,204 |
| `01a11d58-24d0-7823-92d6-e4a0fea6c1b6` | 2026-10-09T00:15:24.658Z | 59,425,460 | 55,283,328 | 317,503 | 101,857 | 59,742,963 |

These are cumulative client-session counters at the stated snapshot, not a marginal task bill. Cached input is a subset of input; reasoning output is a subset of output. Snapshots are not summed. No reliable total monetary cost is claimed. The final documentation/commit operations occur after this snapshot.

Shared memory: `viking://resources/projects/parley-deck/participant-dropout-release-20261009.md` was written and read back byte-identically. The server's wait-for-index operation timed out at 60 s and the first scoped search did not yet return it, so semantic indexing remains unconfirmed. The exact saved note and verification result are retained as `shared-memory-note.md` and `shared-memory-verification.json`.
