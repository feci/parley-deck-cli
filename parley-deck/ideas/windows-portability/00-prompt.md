---
idea: windows-portability
author: user
facilitator: codex-1
created: 2026-09-24
track: deliberation
participants: [claude-1, kimi-1, zcode-1]
auto_implement: true
require_model_diversity: true
status: round-04
---

## Problem / idea

Make the CLI genuinely correct on native Windows. POSIX-only product behavior needs a real Windows implementation or an explicit, reviewed, user-visible refusal. A test may be excluded on Windows only when truly inapplicable and its reason is reviewed. Acceptance requires hosted windows-latest green with no new suppression, with macOS and Ubuntu still green. Nobody on this run has a native Windows machine: never claim validation beyond hosted runners.

## Owner direction and authorization

The owner chose both: release CLI 1.49.1 immediately for macOS and Linux via GitHub and Homebrew, mark Windows experimental, hold CLI winget; simultaneously open this separate reviewed Windows idea to fix the defects and remove the label. Translation from Slovak supplied by the owner. Standing authorization: if there are no owner questions, release and deploy through all channels. This authorizes the lifecycle through implementation, review, integration and release without asking again, subject to the explicit gates below.

The owner's current instruction overrides the deck's GitHub PR mechanics for this run: no development PRs, retain files as canonical, integrate directly and keep main linear. Use this existing windows-portability worktree and branch, base 868825f20bd0988abde1a6c38b9efe405c5123d1. Do not mutate other active worktrees or main during development.

## Track classification

Product/portability change, NOT a protocol change. Deliberation under section 4.0 because snapshot privacy/ACLs are security/privacy surfaces, scope plausibly exceeds 15 files, and auto_implement applies. If design finds a protocol rule that must change, flag that before implementation; do not silently modify protocol or skills. Skill changes are not expected; any necessary skill work requires a separate worktree from origin/main of parley-deck-skill.

## Evidence shared with all participants

Read source-context/README.md and its copied release-ci-revalidation-claude-1.md and release-repair-plan-zcode-1.md. The latter is a participant proposal, not a decision. Verify source locators at current HEAD. Hosted run IDs: 35987916696 and 36009912946. Reported Windows failure classes: internal/evidence test build (syscall.Mkfifo undefined), 14 red packages versus 16 ok in the later run, ACP drain guard, snapshot privacy/ACLs, invalid gate filenames, cross-process locks, and fixture portability. Historical run 35987916696 also failed Ubuntu; keep per-run claims separate. Existing reviewed Linux repairs are in this base and must be preserved.

## Constraints and acceptance

- Design real Windows security and synchronization semantics; do not disable privacy checks or turn product failure into skipped tests.
- Preserve supported Unix behavior and compatibility; document any refusal and persisted-name compatibility decision.
- Hosted matrix: unfiltered build/test Windows, macOS and Ubuntu green on the final released commit, with participant review of absence of new suppression. Review truly inapplicable test exclusions individually.
- Native hosted Windows x64 evidence is not Windows ARM64 execution evidence; state actual coverage and runner assumptions.
- codex-1 only organizes: never implements, verifies code, authors participant verdicts, or signs off. Participants own drafting, implementation, review, goal-done and independent channel verification.
- Use driver first: parley run/continue, wait, organizer brief, status --json. Record concrete driver gaps. Never proxy-write participant artifacts. Participants do not launch another organizer/driver.
- Participant phase invocations write their own requested artifact; do not advance phases, publish, merge or edit peers' files independently. One implementation owner at a time; claim and worktree mapping go in IMPLEMENTATION.md before code edits.
- English artifacts/commits; organizer commits start [codex-1] windows-portability:. Never copy secrets. Never open Google Chrome.

## Release gates and sequence

Release only after BOTH files exist:
1. /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/parley-deck/inbox/codex-1-to-user_release-1.49.1_done.md
2. /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer/parley-deck/inbox/codex-1-to-user_meta-protocol-change-designated-implementer_done.md

After reviewed closure and both prerequisite handoffs, integrate onto latest origin/main, select the next version above what is released, and release GitHub including Windows assets, Homebrew parley-deck-cli.rb, and the held CLI winget PR (one application per PR). Remove Windows experimental labeling IF AND ONLY IF the hosted Windows leg is green on the released commit. Otherwise retain the label, hold winget, and report. Skill channels only if the skill changed. A participant independently verifies every released channel before completion is reported. Do not move released tags.

## Startup capability and context record

Installed CLI: parley 1.49.0. Machine and effective deck roster: exactly claude-1 (Claude CLI, claude-opus-5[1m], max/deep), kimi-1 (Kimi CLI, kimi-code/k3, max/deep), zcode-1 (Zcode CLI, zai/glm-5.3 from config, max/deep); codex-1 inactive as participant and explicit organizer. All three hosted liveness probes ready. Runtime timeouts: Claude/Zcode 30 minutes, Kimi deck override 40 minutes. No exclusions.

Preflight reports source-role protocol drift as advisory. Installer/runtime skills 2.13.0; deck metadata 2.12.0 stale; dry-run sync recorded, no protocol adoption or metadata sync needed for this product task. A preflight gate points to unrelated completed historical meta-protocol-change-devx-speed facilitator metadata; record a scoped driver fallback rather than rewrite historical artifacts.

Protocol attestation: context_mode=packet; source_sha256=8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7; packet_sha256=e88a7edb516ba1ff87b939294698626a5610393637f757816258f118b0ac8b2e; audience=facilitator; phase=0; track=deliberation; fallback_reason absent. The initial exploratory request used an invalid track implementation and was replaced before launch by this valid packet.

Shared memory was unavailable at startup. MCPAnywhere/OpenViking became callable during round 4; a scoped project search returned only the historical 2026-09-23 release note, not current release or authorization evidence. Current local sources remain authoritative for this run. Other historical open ideas remain in their existing worktrees; owner explicitly assigns this workspace and priority. Both release prerequisites absent at startup.

## Completion handoff

Write parley-deck/inbox/codex-1-to-user_windows-portability_done.md with what shipped, hosted evidence, deferrals, organizer-usage.md, and owner actions left. Never represent an incomplete release as shipped.

## Owner default update — 2026-09-25

The owner replaced the global defaults with Claude organizing new runs and Codex implementing, but explicitly exempted this in-flight windows-portability run. This idea retains facilitator codex-1, quorum claude-1/kimi-1/zcode-1, and the original launch-model choices. Machine roster now reports Claude inactive/Opus 5.5 and Codex active; that new-run default does not change this idea's membership, signatures, assignments, or existing sessions. The idea's participant-owned rounds select zcode-1 as drafter/implementer, with Claude and Kimi reviewing. No global roster setting is rewritten by this organizer.

## Round-cap escalation — 2026-09-25

Round 4 is complete; status remains round-04 pending owner direction in inbox/codex-1-to-user_windows-portability_round-cap.md. A new material durability disagreement remains after the three permitted cross-review rounds. No consensus or implementation is claimed. Proposed bounded follow-up: durability-followup-brief.md.
