---
from: codex-1
to: user
idea: meta-protocol-change-lean-organizer
phase: release-handoff
status: blocked-on-owner-input
blocking: yes
date: 2026-09-24
---

# Delivery report — release is NOT complete

This is the requested `done.md` delivery record, explicitly marked blocked rather than a completion claim. Approved A–D implementation reached Phase-8 closure; post-deploy audit exposed pre-existing platform defects, and distribution work remains unfinished. No permission or waiver is inferred from elapsed time.

## A–D delivered and reviewed

Zcode implemented; Claude and Kimi independently reviewed. Three fix-up cycles and four review rounds ended in cycle-4 zero-fix consensus, all three unconditional ACCEPTs, and Kimi's separately commissioned LE-7 PASS. Frozen FINAL remains at its approved scope. Closure commit: 3d57c59. The organizer did not implement, test, verify product code, or sign off.

Features: pure-organizer defaults; deterministic Go wait/digest; facilitator protocol packet and read-only brief; slim 17,802-byte skill core with on-demand references; phase handoff and client-accounting ingestion; participant telemetry repairs. No controlled token-savings or cost-reduction claim is made.

## Published channels and evidence

| Channel | Actual status |
|---|---|
| CLI GitHub | [1.49.0](https://github.com/feci/parley-deck-cli/releases/tag/v1.49.0), immutable tag 06e563e, six assets including Windows. All hashes independently downloaded/verified. Mutable notes correct the unsupported Windows-coverage claim and disclose validation failures. |
| Skill GitHub | [2.13.0](https://github.com/feci/parley-deck-skill/releases/tag/v2.13.0), immutable tag 8161e5e, five assets; final Windows assets from successful CI run35983892815. All hashes independently verified. |
| Homebrew | BOTH formulae updated in feci/homebrew-parley commit697901e. Independent archive URL/SHA verification passed. Installed CLI1.49.0 and skill2.13.0. Installer resolves into `/opt/homebrew/Cellar/parley-deck-skill/2.13.0/`. |
| Runtime skills | `install --target all --force` completed. All12 runtime SKILL.md copies independently match released source SHA256 `db71eca5a4f778cb6cb0d6aee0660474b4dbfc0f8cf19c6585edaf8b0010c03f`. Doctor reports valid. |
| Skill winget | [PR440360](https://github.com/microsoft/winget-pkgs/pull/440360) MERGED at2026-09-24T11:16:23Z. Independent audit verified its manifests and final Windows hashes while open; organizer later checked merge state. |
| CLI winget | NOT submitted. Held on unresolved Windows distribution decision; no catalog-availability claim. |
| npm | NOT published. Latest still2.12.1 at final check. Exact audited2.13.0 tarball preserved; owner publish verification required. |
| Global protocol core | 2.13.0 is now installed in the published core store (observed outside organizer actions). `parley protocol status --json` lists2.13.0; published bytes exactly match the independently reviewed staged template. No core action remains. |

Independent channel audit: `ideas/meta-protocol-change-lean-organizer/release-audit-kimi-1.md`. Release evidence, downloads, logs, checksums and exact npm pack: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-09-24-lean-organizer`. The final delta audit is still owed after npm and the remaining CLI distribution work.

## Post-deploy findings and prepared repair

Hosted CI exposed older defects: Linux stderr loss and a Linux post-Start empty-cmdline identity-capture race; Windows privacy/ACL, filename, locking and fixture portability problems. Workflow setup fixes alone were insufficient. Raw independent incident reports remain canonical.

Linux U1/U2 repairs were implemented by Zcode, independently challenged by both reviewers, reproduced before/after in Linux containers, then confirmed on hosted x86_64 Ubuntu. Run[36009912946](https://github.com/feci/parley-deck-cli/actions/runs/36009912946): Ubuntu PASS, macOS PASS, Windows FAIL. Claude's acceptance report is `release-linux-hosted-acceptance-claude-1.md`. Identity checks remain fail-closed; a bounded100ms empty-cmdline retry does not relax comparison facets. Explanatory comment commit519951e was independently verified behavior-equivalent.

CLI1.49.1 candidate source/metadata: **481fb657f0e219535edc83ca286bd33891cf83dc**. Independent metadata verdict PASS: `release-candidate-1.49.1-review-kimi-1.md`. Six candidate assets built from that clean commit are staged in `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-09-24-lean-organizer/rc-1.49.1/`, with `sha256.json`. They are **NOT tagged, published or installed**. v1.49.0 remains unchanged. Candidate wording records the pending Windows decision; it must reflect the owner's actual decision before publication. No Windows architecture work began.

## Owner decisions/actions required

1. **Windows scope:** choose explicit deferral with Windows labelled experimental/unvalidated, or authorize a separate reviewed Windows-portability track. Omitting Windows assets is a separate deviation requiring explicit direction. See `codex-1-to-user_meta-protocol-change-lean-organizer_scope.md` and the concrete `release-repair-plan-zcode-1.md`. This decision is required because the brief requires audit findings fixed before completion, while native-Windows architecture changes exceed A–D. No added mandatory protocol obligation is proposed.
2. **npm:** complete the exact-tarball command in `codex-1-to-user_meta-protocol-change-lean-organizer_npm-login.md`, or be present for a fresh publish-verification URL. Expired links cannot be reused.
## Global core publication observed complete

A final store check found2.13.0 published outside this organizer's actions (file mtime2026-09-24T11:19:35Z). `parley protocol status --json` lists installed2.10.0 and2.13.0. Published file `/Users/tomasfecko/.parley/protocol/core/2.13.0/COOPERATION.md` and staged file are byte-identical:109,772 bytes, SHA256 `fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f`. This matches the participant-reviewed combined template, including prior2.11.0, released1.48.0 and lean-organizer hunks. The organizer did not invoke publish or bypass its TTY gate. Earlier participant reports correctly describe their earlier observation; this current handoff supersedes their pending-core status. No repeat publish is requested.

After the scope decision: finalize candidate wording/metadata as needed, rebuild from final pinned commit, integrate/tag/publish new immutable1.49.1 through authorized channels, update CLI Homebrew and its own winget PR, install, then request participant delta audit of every changed/pending channel. No development PR is required; no released tag moves.

## Deferred by the signed original consensus

Inactive, unstaffed candidates (no new implementation launched):
- DF-1 `meta-protocol-change-consensus-duty-gates`: unratified hard gates removed; future protocol decision needed.
- DF-2 `meta-protocol-change-facilitator-integrity-phase-coverage`: further §15 phase pinning.
- DF-3 `facilitator-packet-per-phase-bounds`: only phase1 has the ratified hard byte guardrail.
- DF-4 `release-binary-reproducibility`: build-path/VCS provenance beyond current recorded method.
- DF-5 `wait-boundary-vs-published-fixup`: stale prior-review boundary and quoted-status handling; current workaround is create the intended next-round directory before waiting.

Native Windows work is **awaiting** owner decision, not silently deferred. U3 identity fallback hardening remains a proposal and was not implemented. The original acceptance disclosures, including ambiguous historical phase attribution, remain recorded.

## Organizer accounting

Full per-boundary cumulative client accounting: `ideas/meta-protocol-change-lean-organizer/organizer-usage.md`, from the pinned own Codex rollout (not model self-report). Last observation:

| Boundary | Accounting timestamp | Input | Cached input | Output | Reasoning | Accounting request events |
|---|---|---:|---:|---:|---:|---:|
| Owner-blocked release handoff; all independent preparation finished (observed 2026-09-24T14:46:24.664Z) | 2026-09-24T14:46:24.664Z | 44910726 | 44147840 | 166206 | 81192 | 326 |

Totals include the substantial unplanned post-release CI investigation; they are not a controlled savings experiment. No raw rollout or credentials were committed.

## Resume state

Remote CLI main last pushed9134c7a contains reviewed Linux fixes; local lean-organizer additionally contains comment/candidate metadata and this handoff. Skill branch/tag remains8161e5e/2.13.0. Preserve untracked runs and usage-ledger smoke artifacts; do not blindly stage them. No participant process remains running. Continue from the candidate and evidence above after the owner's response, without repeating completed A–D deliberation or moving published tags.
