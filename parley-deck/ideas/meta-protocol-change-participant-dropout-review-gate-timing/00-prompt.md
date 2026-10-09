---
idea: meta-protocol-change-participant-dropout-review-gate-timing
author: user
created: 2026-10-09
participants: [codex-1, zcode-1]
excluded: [kimi-1 — two measured readiness failures HTTP 403 connection allowlist / quota scope; controlling brief authorizes this exclusion — confirmed 2026-10-09]
implementer: codex-1
track: deliberation
auto_implement: true
require_model_diversity: true
checks: go test ./... -count=1 -timeout 45m
status: final
---

## Problem / idea

Complete the follow-up to shipped participant dropout: after a recorded eligible dropout or exclusion leaves exactly one model-diverse non-implementer reviewer, allow that reviewer to satisfy review, review consensus, goal check and auto-close. Without a recorded dropout/exclusion the gates stay unchanged. Decide precise admissible evidence, fresh-process goal-check independence, strict_gate interpretation, bounded watchdog timing and settlement in the middle of a round/consensus/signoff batch. Preserve dissent, protected roles and the implementer-plus-one floor. Compact existing guidance without losing rules or raising size limits.

The controlling brief is [source-context/ORGANIZER-BRIEF.md](source-context/ORGANIZER-BRIEF.md), copied byte-identically (SHA256 50cb660916b882965ff289dc2fc7d5ff9c8f234cac61d2e95e1e5c8849c8722a). Read it in full. This run explicitly overrides the usual Claude organizer default: codex-1 organizes and implements as an ordinary participant; no facilitator field; claude-1 takes no part. Zcode is an independent separate-process reviewer; the implementer never reviews itself.

## Owner direction

Verbatim source quotations are retained at the owner's express request; the surrounding protocol artifacts are English.

> "organizator a implementator su vacsinou codex a claude, takze ucastnici su kimi a zcode, ak vypadnu z hocikakeho dovodu, tak v tej idei uz nepokracuju a mozu sa zapojit az do dalsej"

Translation: "The organizer and implementer are usually codex and claude, so the participants are kimi and zcode. If they drop out for any reason, they do not continue in that idea and can join only the next one."

> "participanti nie su nevyhnutne potrebny obaja, staci jeden a to by nemalo zaseknut parley-deck"

Translation: "The participants are not necessarily both needed; one is enough, and that should not stall parley-deck."

> "ano spusti ideu ako poriesit ze participant vypadne v strede idei, ved tam mame nejake timeouts, takze kludne rozved participant-dropout-review-gate-timing"

Translation: "Yes, start the idea on how to handle a participant dropping out in the middle of an idea. We do have some timeouts there, so feel free to expand participant-dropout-review-gate-timing."

## Constraints

- Prior art: ../meta-protocol-change-participant-dropout/ (especially FINAL.md, organizer-notes.md and its released note in ../../inbox/), and ../meta-protocol-change-quota-auto-exclude/. Reuse the shipped reducer, durable two-attempt ledger, history, retained obligations and recovery. Do not redo those features.
- Work only in the owner-allocated CLI review-gate-timing worktree and sibling review-gate-timing-skill; never prune or declare worktrees. Base CLI 85c5a8f / 1.52.0; skill dbdb919 / 2.16.0. Transport github-pr; owner-selected branch names are retained.
- Design before implementation; FINAL freezes the precise protocol hunks, CLI files, tests, compaction and ACs. Write IMPLEMENTATION.md plan before product edits. At most five fix-up cycles; unresolved CRITICAL/MAJOR after cycle five requires a blocking owner note.
- Run readiness and record Kimi's actual failure. No credential/key/gateway changes. If unavailable, the owner's explicit per-idea exclusion authority leaves codex-1 + zcode-1. No global roster mutation.
- Driver first (run/continue/wait/status/consensus), with measured CLI fallback only for a recorded driver gap (D6 is known). No artificial accounting or fabricated participant artifacts.
- Own/Zcode 429, 503 or timeout: wait 15 minutes, then retry at most eight times per step; timeout relaunch ceilings 2400s then 3600s. Codex auth/credit error: stop with a blocking note. These run-specific instructions do not define the product's retry rule.
- The owner-attended close is pre-authorized by this brief: no open CRITICAL/MAJOR in final Zcode review, both current participants' review-consensus ACCEPT blocks, fresh Zcode goal-check process PASS, current-tree AC evidence. Quoted owner intent above supplies the rationale. This authority is specific to this run before the new rule ships; it never fabricates driver auto-close success.
- Release after close: merge main, CLI 1.53.0 / skill 2.17.0 / staged core 2.17.0. Core is staged 2.16.0 plus exactly reviewed hunks. GitHub releases, both Homebrew formulae, skill-only WinGet PR; install all managed skills and four named generic targets and verify every SKILL.md hash. Separate Zcode channel verifier. npm/core remain owner-only ordered publication commands in the final released note.

## Readiness and provenance

Full live protocol read. Phase-0 full attestation source_sha256=packet_sha256=091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf; fallback_reason absent. Live source-role deck and installed skill agree; installer/runtimes 2.16.0, CLI 1.52.0; no metadata sync needed. Shared OpenViking prior release recall checked against local released notes.

Standalone parley preflight --no-ping reports all three requested CLIs installed and an unrelated existing facilitator-declaration gate on meta-protocol-change-devx-speed. It has no participant-filter flag; presence-only avoids launching the expressly excluded Claude. Kimi's real bounded readiness is measured separately through parley agents exec; final result will be recorded here. No unrelated historical gate is repaired.

Effective participant defaults: codex-1 gpt-6-astra/max/deep (current organizer session; CLI available); kimi-1 kimi-code/k3/max/deep (effort from config); zcode-1 zai/glm-5.3/max/deep (native model/effort from config). Zcode effective argv is --prompt={prompt} --mode yolo --cwd {root}; fresh one-shot processes, default 1800s ceiling. Kimi configured ceiling 2400s; readiness probe 90s. External backend context limited to task and repository, never secrets.

## Non-goals

No roster/model/credential changes, reviewer impersonation, silent removal of dissent, blanket single-reviewer waiver, new retry service, size-limit increase, D6/accounting repair, Windows CLI promotion or unrelated legacy quota/alias remediation. Published global core remains owner-controlled.

Readiness settled: Kimi failed twice with the relayed HTTP 403 allowlist/quota-scope error; see source-context/kimi-readiness-evidence.md. The original proposed participants were [codex-1, kimi-1, zcode-1]. Under the controlling brief Kimi is excluded permanently for THIS idea; codex-1 and zcode-1 remain. This is the owner-attended recorded exclusion, not a forged automatic transition or a global roster edit.
