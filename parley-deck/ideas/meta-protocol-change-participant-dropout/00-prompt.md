---
idea: meta-protocol-change-participant-dropout
author: user
created: 2026-10-08
participants: [codex-1, zcode-1]
excluded: [kimi-1 — owner-authorized two failures on round-01 (HTTP 400 ambiguous k3); no same-idea return — confirmed 2026-10-08]
implementer: codex-1
track: deliberation
auto_implement: true
require_model_diversity: true
checks: go test ./... -count=1 -timeout 45m
status: round-02
---

## Problem / idea

Broaden CLI 1.51.0 quota auto-exclusion to drop a non-protected participant whose dispatched step fails for any reason, for this idea only, using the shipped transition mechanism. The minimum is the implementer and at least one other usable participant; an organizer never supplies an extra floor slot. Preserve every filed objection, finding and DISPUTED claim. The controlling source is [ORGANIZER-BRIEF.md](source-context/ORGANIZER-BRIEF.md), copied byte-for-byte (SHA256 54448074f5a4dac2587ee3ed9c4e08866d84e50d42c02df7dc9b1e6af70651b5).

## Owner direction (verbatim, original language Slovak)

> "ok, cize je v parley deck automaticke pravidlo, ze ak niektory z participant modelov nie je dostupny tak sa preskoci? ak mame organizatora, implementatora a potom participantov, tak participanti nie su nevyhnutne potrebny obaja, staci jeden a to by nemalo zaseknut parley-deck"

Translation: "OK, so is there an automatic rule in parley deck that if one of the participant models is unavailable, it is skipped? If we have an organizer, an implementer and then participants, the participants are not necessarily both needed. One is enough, and that should not stall parley-deck."

> "organizator a implementator su vacsinou codex a claude, takze ucastnici su kimi a zcode, ak vypadnu z hocikakeho dovodu, tak v tej idei uz nepokracuju a mozu sa zapojit az do dalsej"

Translation: "The organizer and the implementer are usually codex and claude, so the participants are kimi and zcode. If they drop out for any reason, they do not continue in that idea, and they can join only from the next one."

The current owner specifically assigns codex-1 as organizer and implementer for this idea. No facilitator field and no claude-1 participant. This overrides the global organizer default. Verbatim quotations are retained as expressly requested; all authored discussion is English.

## Constraints

- Settle all seven numbered design points in the controlling brief. Reuse existing quota machinery, keep changes small, list exact code/test/protocol hunks and migration. Protect organizer, implementer/designee, and started consensus/FINAL drafter. No same-idea rejoin for a dropped participant. Evaluate the whole batch; below floor apply nothing and escalate once with owner options.
- Prior art: ../meta-protocol-change-quota-auto-exclude/{FINAL.md,IMPLEMENTATION.md,review/}; live COOPERATION.md sections 5 and 9.0. The prior release is delivered; do not redo its native quota recognizer project. Accepted limitations remain evidence, not new passes.
- CLI worktree: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/participant-dropout (base d16ee9c). Skill worktree: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/participant-dropout-skill (base 352a475).
- No code before signed consensus, FINAL and a written IMPLEMENTATION plan. Non-implementers own independent reviews; at most five fix-up cycles. No self-review, no proxy-written participant artifacts, no secrets, no Chrome, no pruning or worktree declarations.
- Continue without proposal stop if FINAL matches the owner rule and leaves no owner decision. Close and releases pre-authorized under the brief's explicit conditions. Target CLI 1.52.0, skill/core 2.16.0. Core based on staged 2.15.0 plus exactly reviewed hunks; npm/core publication commands remain owner-only.

## This run's temporary failure authorization

The new rule is not in force. Shipped rules apply. The controlling brief additionally authorizes dropping kimi-1 or zcode-1 after two failures of the same step (original plus one relaunch), provided codex-1 plus one other participant remain. Record an owner-authorized excluded entry and non-blocking notice. Both failing blocks once. This specific two-attempt participant authorization governs over the general retry ceiling. Provider 429/503 or timeout relaunch waits 15 minutes; timeout ceilings 2400 s then 3600 s. General protected-agent retry ceiling is eight per step; codex-1 auth/credit stops. A valid completed artifact takes precedence over a failed process.

## Startup / capabilities

`parley preflight --dir . --json` ran before kickoff. CLI 1.51.0; installer and managed runtime skills 2.15.0; source protocol metadata valid and matching. Preflight reported codex-1 and zcode-1 ready, kimi-1 process-failure exit 1. No exclusion is inferred from that single observation. The global preflight also probed claude-1 (ready), but it takes no part in this idea. An unrelated historical devx-speed facilitator declaration gate is untouched. Existing unrelated ideas in separate worktrees do not supersede this explicit task.

| Participant | CLI / model | Effort / speed | Timeout | Autonomous write |
| --- | --- | --- | --- | --- |
| codex-1 | Codex / gpt-6-astra | max / deep | 1800 s default | workspace-write, approval never |
| kimi-1 | Kimi 0.42.0 / kimi-code/k3 | max configured / deep | 2400 s | print mode -p |
| zcode-1 | Zcode 3.7.7-13 / zai/glm-5.3 | max configured / deep | 1800 s | --mode yolo --cwd workspace |

Transport is the existing github-pr deck transport; this owner-designated participant-dropout branch is the run branch. Canonical files govern; PR mirrors cannot approve as the same feci account that creates a PR. Driver-first orchestration is used, with recorded fallbacks for named-idea startup and the known D6 accounting gap. No global roster edit.

Protocol attestation for startup: context_mode=full; source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e; fallback_reason absent. Full live source read. Relevant OpenViking release memory was read and checked against local shipped artifacts; local files control. No graph exists in this worktree; source navigation uses direct searches.

## Non-goals

No provider/model/credential change, automatic replacement reviewer, organizer failover, native zcode capture, Windows portability fix, D6 budget migration or broad roster redesign. Existing global membership stays unchanged.
