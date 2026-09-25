---
from: codex-1
to: user
idea: meta-protocol-change-designated-implementer
phase: release-handoff
status: agent-controlled-delivery-complete
blocking: no
date: 2026-09-25
recorded_at: 2026-09-25T09:13:30.109554+00:00
owner_actions_pending: [npm-login-and-publish, attended-core-publish]
external_pending: [skill-winget-validation-and-merge]
---

# Designated implementer release handoff

**CLI 1.50.0 and skill 2.14.0 are released on GitHub and Homebrew, installed locally, and independently audited. All 15 existing managed runtime roots now carry the 2.14.0 skill payloads: 90/90 SKILL.md files and 15/15 bundled protocol copies independently match. The owner's selected machine default is now `codex-1`; the product still ships UNSET.**

This is completion of the agent-controlled delivery, **not a claim that every channel is published**: npm 2.14.0 still needs owner login, core 2.14.0 still needs attended publication, and skill WinGet PR441054 is open with external validation/merge pending. Exact owner actions are below. This file releases the sequencing hold for `windows-portability`, with those pending states carried explicitly.

## What shipped and how it was closed

Optional per-idea `implementer:` and layered `[defaults].default_implementer` designate the participant that executes FINAL, for code or other work. Resolution preserves the re-entry pin, then per-idea designation, standing default, and existing fallback; `none` suppresses the standing default. The implementation adds fail-closed validation, explicit unavailable-agent handling, drafter separation and reviewer-diversity warnings, and one shared resolution chain. No quorum/signoff-weight change, no shipped default selection, and no unrelated platform work.

Codex-1 remained the pure declared organizer. Kimi-1 implemented; Claude-1 and Zcode-1 independently reviewed. Three design rounds and four review cycles culminated in fresh unanimous cycle-4 ACCEPT signatures with zero agreed fixes. Fresh non-implementer Zcode goal check passed all 21 FINAL criteria; Kimi marked its implementation complete at `ba46b2a`.

Evidence remains exactly scoped: independent full Go suite at `d238238`; later behavior/mutation evidence at `1bad263`; comment-only checks at `717f3de`; current-tree goal check at `804522c`. Integration of main's released Linux repairs justified Kimi's new full suite at release candidate `01fc495` (all packages pass). Hosted candidate and main CI pass Linux/macOS; Windows fails as explicitly deferred. Independent skill checks passed 399 Node tests, 54 Python tests and all six add-on manifests. No old evidence was relabelled as a new run.

## Channel evidence

| Channel | Actual state |
|---|---|
| CLI Git / GitHub | [v1.50.0](https://github.com/feci/parley-deck-cli/releases/tag/v1.50.0), immutable tag **01fc49526ca66b771b549a01cfbbf4addeaa886b**. Main fast-forwarded directly, no development PR; later commits preserve records only. Six clean-build binaries plus [sha256.json](https://github.com/feci/parley-deck-cli/releases/download/v1.50.0/sha256.json). Zcode independently downloaded every asset and matched bytes, API digests, manifest and build provenance; all six embed the exact clean candidate revision and correct platform. |
| Skill Git / GitHub | [v2.14.0](https://github.com/feci/parley-deck-skill/releases/tag/v2.14.0), immutable tag/main candidate **a5664d803f1fb6156ef95ff5921dcf13a3dd2031**. Portable workflow **36113011791 succeeded**. Both required Windows assets were independently downloaded and hashed. Current workflow/RELEASING/README require this pair; predecessor macOS/Linux convenience assets were extra manual uploads, not a missing channel. |
| Homebrew | Both formulae published in [tap commit9855366](https://github.com/feci/homebrew-parley/commit/9855366ad57146d1428aae180cf4c50b41667ae5). Independent archive hashes: CLI `2ff8426d257ec4a6b9d92ac34139ef97d92b3a4259d6c2b25d25c424d73ce6a3`; skill `f43e02f7db2050c0c778c5f232c091c417961ea66c079c34ce061faa8a0feafe`. Both installed/upgraded; participant `brew audit` and both `brew test` runs exited0. Installed commands report CLI1.50.0 / skill2.14.0. |
| Skill WinGet | [PR441054](https://github.com/microsoft/winget-pkgs/pull/441054), head `26720f091548004bd9d9cc5e3606f4bb6c66ddda`: exactly one application/version and three manifests. Hashes come from final uploaded Windows assets after workflow completion. Latest participant check: steps01–07 and CLA success; installation validation in progress; metadata/completion queued; **OPEN, not merged**. Local Windows `winget validate/install` was unavailable on macOS and never claimed. |
| CLI WinGet / Windows | **CLI WinGet HELD by owner instruction.** Both Windows binaries remain individually labelled experimental/unvalidated. Reviewer compared candidate and predecessor Windows failures: no newly failing tests; the app test binary aborts before designation tests run, so no Windows feature-validation claim. Separate `windows-portability` owns DF-4. |
| npm | **2.14.0 NOT published.** Exact audited tarball publish failed E404 access/not-found. Independent registry query still shows latest2.13.0 and no2.14.0. Owner login requested; other channels completed. |
| Runtime skills | **PASS at all15 managed roots /90SKILL.md files**, full payloads equal shipping2.14.0, all15 markers2.14.0. Twelve detected runtimes plus three dormant managed roots (Goose/Cursor/Aionrs); dormant copies also refreshed with `--include-undetected`. No claim those runtime tools became installed. All15 bundled COOPERATION.md files hash `51476d69ed4291b77db63f6855554cd251461c0b8743a303c2a3b0eacaabf67a`. |
| Machine default | Owner-authorized `[defaults].default_implementer = "codex-1"` applied after GitHub/Homebrew shipment and CLI upgrade at **2026-09-25T08:43:51Z**, independently verified. No roster/model/effort change. The frozen current idea still pins Kimi and keeps Codex outside its quorum. |
| Core | **2.14.0 staged and independently verified, NOT published into the core store.** Published2.13.0 TEMPLATE plus exactly this idea's seven reviewed hunks; placeholders retained, no project header leaked. The core store still contains2.10.0/2.13.0. Owner-attended command below. |

## Independent signoff on delivery

- [Claude preparation audit](../ideas/meta-protocol-change-designated-implementer/release-preparation-audit-claude-1.md): PASS, no blocking findings. OBS1/3 addressed in published release notes with candidate CI and explicit Windows non-execution; OBS2 addressed by retained complete skill-gate logs at the frozen candidate.
- [Zcode final channel audit](../ideas/meta-protocol-change-designated-implementer/release-channel-audit-zcode-1.md): PASS for actual agent-controlled channels, explicit pending items.
- [Zcode addendum](../ideas/meta-protocol-change-designated-implementer/release-channel-audit-addendum-zcode-1.md): PASS for all15 refreshed roots/90SKILL.md files; corrects its original report's NIT label, hash-table link and separate-main-checkout ambiguity without changing the filed original.
- Full machine evidence: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-09-25-designated-implementer/`, including `publication-commands.json`, exact tarball/assets, full logs, and `audit-zcode-1/INDEX.md`.
- Final per-path hashes: `audit-zcode-1/addendum/runtime-hash-table-15roots.txt`; machine verdicts beside it.

All participant processes exited0 before this handoff. The organizer implemented and verified no product code and issued no participant signoff.

## Owner-only actions still needed

**npm login:** per the owner's release brief, run `! npm login` in the Claude session and report success. Then retry the **same audited tarball**, without repacking:

```sh
npm publish --access public '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-09-25-designated-implementer/skill/parley-deck-skill-2.14.0.tgz'
```

SHA256 `681b9076a017e67935403ad6cd9779cafa09a117d168331daca300404336135b`. Recheck the registry before a later retry so a subsequent release is not accidentally replaced as latest. A successful retry still needs independent registry-byte/integrity verification. [Owner inbox](codex-1-to-user_meta-protocol-change-designated-implementer_npm-login.md).

**Core publication:** the CLI's controlling-terminal gate and the owner's explicit no-bypass instruction require an attended terminal. Run:

```sh
parley protocol publish --version 2.14.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.14.0.md
```

Staged SHA256 `51476d69ed4291b77db63f6855554cd251461c0b8743a303c2a3b0eacaabf67a`. A participant must verify the published store afterwards. No agent allocated a TTY or invoked publication. [Owner inbox](codex-1-to-user_meta-protocol-change-designated-implementer_core-publish.md). A following release must account for this verified staged2.14.0 while owner publication is pending, preserving its reviewed hunks.

## Default choice, residuals and scope

The owner already selected **codex-1 / GPT-6 Astra** for future implementation; no recommendation/selection question remains. Its operational fit is supported by the installed active Codex adapter and the owner's intended separate Claude organizer plus Kimi/Zcode reviewers from distinct model families. This run implemented with Kimi; it is **not** a comparative benchmark proving Codex superior. The mechanism ships unset, and any unavailable/ineligible standing default emits its documented fall-through notice; a malformed/ineligible per-idea designation remains fail-closed.

The roster was intentionally not changed. Live machine roster still lists Claude active alongside Codex/Kimi/Zcode, a pre-existing discrepancy with the later new-run owner context's inactive-Claude statement. This release preserves the explicit no-roster-change boundary; it does not silently repair or rewrite global membership or in-flight ideas.

Carried follow-ups: DF-1 pre-existing drafter-precheck divergence; DF-2 design-only launch surfacing (named inactive idea `meta-protocol-change-designation-launch-surfacing`, not opened here); DF-3 standing-default fall-through notices now documented in this handoff; DF-4 Windows portability. TUI/pipeline surfaces and absence of a live designated-deck/ping demonstration remain the named FINAL limitations; no expanded validation claim.

One closing review NIT remains a **real unrepaired accepted residual**: Claude round04's overbroad WITH-flags lead-in alongside its explicit WITHOUT-flags exception in IMPLEMENTATION.md. All three accepted its recorded dismissal; this handoff does not call it repaired. Earlier separate dismissals remain in the signed record.

Driver gaps and manual fallbacks are retained in organizer-notes.md: stale original run cursor / inaccessible historical worktree; by-idea continuation selecting a signoff run; wait's root-consensus/previous-round observations; model-preserving manual launches; consult adapter-ID/120s first-output mismatch; and transient provider DNS failure. No run history was forged to bypass a gate.

## Organizer usage and memory

[Organizer ledger](../ideas/meta-protocol-change-designated-implementer/organizer-usage.md), backed by usage-ledger.jsonl and `parley usage ingest` from the pinned organizer rollout. Final boundary snapshot **2026-09-25T09:10:24.327Z**, **866 accounting events**: input **125,946,129**, cached input **124,681,728**, output **231,767**, reasoning output **102,598**, reported total **126,177,896**. These are cumulative client accounting totals, largely repeated cached context; attribution is **ambiguous**. They are not unique-token, precise phase-cost, billing, or savings figures. Snapshot precedes final recording/push/response.

OpenViking gateway discovery returned `Session terminated`; no shared-memory write claimed. A concise verified handoff is retained locally as memory-pending.md with its intended URI. No secrets or raw provider logs committed. Existing untracked runs/consults preserved. Google Chrome was never opened or controlled.
