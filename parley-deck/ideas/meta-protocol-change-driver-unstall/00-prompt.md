---
idea: meta-protocol-change-driver-unstall
author: user
created: 2026-10-09
participants: [codex-1, kimi-1, zcode-1]
implementer: codex-1
track: deliberation
auto_implement: true
require_model_diversity: true
checks: go test ./... -count=1 -timeout 45m
status: round-01
---

## Problem / idea

Fix the two recurring driver stalls: D6 legacy cross-review accounting and the 120-second goal-check ceiling. Preserve fail-closed accounting and LE-7/LE-11 close integrity. Assess the other small, demonstrated driver stalls in the three predecessor runs; include only narrow fixes and explicitly defer the rest.

The controlling [ORGANIZER-BRIEF.md](source-context/ORGANIZER-BRIEF.md) is copied byte-for-byte, SHA256 `61bbd09b54a5afa753e4dcac0e03a5232e24cc6248c4bc3ef18ae655b19d5d1f`. It assigns codex-1 as organizer and implementer, overriding the machine organizer default for this idea. Claude takes no part. Because the existing protocol explicitly specifies the goal ceiling, this uses the permitted `meta-protocol-change-` prefix and the deliberation track.

## Owner direction (verbatim; original language Slovak)

2026-10-04, on the recurring legacy D6 stop: **"Samostatná trvalá oprava (Recommended)"**. Translation: "A separate lasting fix (Recommended)."

2026-10-08: **"participanti nie su nevyhnutne potrebny obaja, staci jeden a to by nemalo zaseknut parley-deck"**. Translation: "The participants are not necessarily both needed; one is enough and that should not stall parley-deck."

2026-10-09: **"ano spusti ideu ako poriesit ze participant vypadne v strede idei, ved tam mame nejake timeouts…"**. Translation: "Yes, start the idea on how to handle a participant dropping out mid-idea; we have some timeouts…"

Verbatim source quotes are retained by explicit owner instruction; all authored discussion is English.

## Constraints

- CLI worktree `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/driver-unstall`, branch `driver-unstall`, base `128e30b`; skill worktree `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/driver-unstall-skill`, branch `driver-unstall`, base `8ce4dec`.
- D6 must become a lasting, auditable fix. No silent ignoring of unknown history, reset of charged cycles, fabricated historical identity, deletion of evidence, worktree pruning or worktree declarations. Evaluate durable one-time declarations, idea-scoped gates, and provable legacy migration.
- Goal checks remain fresh independent processes, fail closed on missing/failed/inconclusive/reserved results, and supplement current-tree independent acceptance evidence. Fit real execution duration; reuse shipped retry/watchdog machinery where applicable.
- Read predecessor organizer notes: `../meta-protocol-change-quota-auto-exclude/organizer-notes.md`, `../meta-protocol-change-participant-dropout/organizer-notes.md`, and `../meta-protocol-change-participant-dropout-review-gate-timing/organizer-notes.md`. Verify claims against current source under section 15. Historical assertions alone are testimony.
- No product changes before signed consensus, FINAL, and an IMPLEMENTATION plan. Each participant owns its artifacts. zcode-1 owns independent code reviews and fresh goal checks in separate processes. codex-1 does not review its own code.
- Try the driver first and record concrete limitations before configured CLI fallbacks. Dogfood the 1.53.0 dropout reviewer rule. If it does not mechanically qualify, the brief pre-authorizes an attended close only after no open CRITICAL/MAJOR in zcode-1's final review, both current ACCEPT blocks, fresh zcode-1 goal PASS, and current-tree AC evidence.
- Run through design, implementation, independent review, and release without a proposal stop unless an owner-only decision remains. Maximum five fix-up cycles. Provider 429/503 or timeout: wait 15 minutes, at most eight relaunches per step; timeout retries 2400 seconds then 3600 seconds. codex-1 auth/credit failure stops with one blocking note. Never alter credentials or provider settings.
- Release authority and all channels are specified in the copied brief. Protocol change targets CLI 1.54.0 and skill/core 2.18.0; stage core from staged 2.17.0. npm/core publication remains an owner command where the shipped attended controls require it. Preserve owner worktrees.

## Startup / selection checkpoint

Transport: the live `github-pr` authority, with the owner's named worktree/branch. All participant GitHub identities map to `feci`; self-approval limitations will be recorded and canonical files remain authoritative.

| Participant | CLI / selected model | Effort / speed | Initial task ceiling | Autonomous write |
| --- | --- | --- | --- | --- |
| codex-1 | Codex 0.161.0 / gpt-6-astra | max / deep | 1800 s | workspace-write / approval never |
| kimi-1 | Kimi 0.42.0 / kimi-code/k3 | max configured / deep | 2400 s | print mode `-p` |
| zcode-1 | Zcode 3.7.7-13 / zai/glm-5.3 | max configured / deep | 1800 s | `--mode yolo --cwd <workspace>` |

Zcode's native model configuration was read selectively: `model.main=zai/glm-5.3`; the CLI has no per-launch model flag. Snapshot/configuration identity is configured authority, not a native observed-model claim.

CLI 1.53.0, installer and all eight managed runtimes 2.17.0. Source metadata valid and matching packaged protocol. `parley preflight --no-ping --json` exits 3 solely for the historical unrelated devx-speed facilitator declaration. Standalone preflight has no participant filter, so presence-only avoids invoking the excluded Claude runtime; actual requested-peer readiness is measured separately and retained here. No global roster edit.

Protocol context: `full`; source_sha256 = packet_sha256 = `acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137`; fallback_reason absent. Full live authority read. Shared OpenViking project context consulted, with local files taking precedence. No graphify graph exists; targeted source navigation is used. Unrelated historical ideas and their paused requests are outside this explicit assignment; no open GitHub PR conflicted at startup.

## Non-goals

No broad budget rewrite, provider substitution, credential/gateway repair, automatic reviewer replacement, Windows CLI release, unrelated old-idea changes, or graph-build side project. Deferred driver gaps are listed explicitly in FINAL and the release handoff.
