---
idea: meta-protocol-change-designated-implementer
author: user
facilitator: codex-1
created: 2026-09-24
track: deliberation
participants: [claude-1, kimi-1, zcode-1]
auto_implement: true
require_model_diversity: true
status: final
---

## Problem / idea

Should the owner be able to designate the participant who executes FINAL (code or other work), while the others evaluate it? Design the mechanism, independently verify the current behavior and both resolver copies at HEAD, then implement and review the agreed solution. The full controlling owner brief is copied verbatim to `source-context/owner-brief.md`; every participant must read it.

## The owner's words (verbatim, Slovak) and translation

> "cize organizator ma specialne miesto v parley-deck protokole, rozmyslam ze este by nebolo zle setnut aj
> implementatora, ktory vlastne napise kod, ak treba, alebo teda spravi tu pracu co treba a ostatny to
> ohodnotia, kto to robi teraz? je to v protokole predpisane?"

"So the organizer has a special place in the parley-deck protocol. I am thinking it would be good to also set
the implementer, who actually writes the code if needed, or does whatever work is needed, while the others
evaluate it. Who does it now? Is it prescribed in the protocol?"

> "ano otvor na to ideu a navrhnite to cez /parley-deck cim skor"

"Yes, open an idea for it and design it through /parley-deck as soon as possible."

> "a ked na mna nebudete mat otazky tak to rovno releasni a deployni cez vsetky kanaly"

"And if you have no questions for me, release it and deploy it through all channels directly."


## Constraints

- codex-1 is the declared pure organizer, never a participant, signer, implementer or code verifier. Participants own every canonical round, signoff, FINAL, implementation and review artifact.
- Ship with the global default UNSET. No agent is selected as the owner's persistent default. Recommend one with evidence after release for the owner to decide.
- Cover role scope, per-idea/global precedence, unavailable-agent behavior, drafter separation and self-contained FINAL, reviewer diversity, claim fallback, driver/preflight fail-closed validation, design/signoff participation and resolver consistency.
- Lightest permitted track is deliberation: §4.0 explicitly forces protocol changes (§7), and auto_implement also triggers it. Speed is achieved by driver orchestration and bounded rounds, not reducing the mandatory track.
- The owner authorizes normal phase transitions and release when no owner questions remain. Use local canonical files and direct main integration, without development PRs, as an explicit per-run override of the existing github-pr transport mechanics. Do not change the global transport header.
- CLI worktree is this directory, branch designated-implementer. Skill worktree is `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer-skill`, same branch. Change all three protocol copies. No participant publishes, merges to main, changes global defaults or invokes core publication before the organizer's release step.
- All facts and hypotheses in the brief are unverified testimony until participants inspect HEAD and record §15 evidence. The organizer makes no code-verification verdict.
- FINAL must be implementable by a non-drafter. The current protocol, not the proposed mechanism, governs this run until ratified.
- Observe all §15 duties and refutation-default review; the runner may omit some duties in its prompts. Existing alternatives must use exact `## Existing alternatives` heading.
- All new artifacts/commits are English; owner's expressly requested original Slovak quotations are preserved as source evidence alongside their English translations. Every commit uses `[codex-1] meta-protocol-change-designated-implementer: ...` per owner instruction; canonical file authorship remains each participant's.

## Readiness and context

Machine roster and hosted preflight on 2026-09-24: claude-1 / Claude Opus 5[1m] / max; kimi-1 / kimi-code/k3 / max; zcode-1 / zai/glm-5.3 / max. All installed, autonomous, deep speed, hosted PONG ready. Codex inactive in quorum intentionally. Kimi's configured round timeout is 40 minutes, others 30 minutes; organizer wait calls remain below 25 minutes.

CLI 1.49.0, skill installer and runtime markers 2.13.0. Project is protocolRole: source; stale 2.12.0 metadata and differing bundled reference are advisory. No protocol replacement or metadata refresh was performed. Preflight reports an unrelated historical facilitator declaration conflict in meta-protocol-change-devx-speed. This run has valid separate facilitator/quorum and will not change that other idea.

Organizer packet: context_mode=full, source_sha256=8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7, packet_sha256=8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7, fallback_reason absent. OpenViking tools are unavailable in this session; local sources are used.

## Release boundary

The previous release's required done.md exists and process 21386 has exited. Its copied `source-context/prior-release-handoff.md` reports incomplete channels and an unresolved Windows scope decision; do not infer that release succeeded or inherit an authorization to fix unrelated platform architecture. Participants must assess inherited limitations against this idea's delivery criteria. Integrate latest origin/main only after reviewed implementation. Select the next minor versions above actual releases. Global core 2.13.0 was reported published; verify live at staging time. Core publication remains owner-attended; never allocate a TTY to bypass that gate.

## Non-goals

No quorum/roster change, model/effort downgrade, unrelated code, or organizer-rule changes beyond this role's needs. No persistent default implementer selection.

## Subsequent owner direction — release order and Windows

The controlling inbox `inbox/user-to-codex-1_meta-protocol-change-designated-implementer_release-order.md` supersedes the original release wait rule. Wait for the lean-organizer inbox `codex-1-to-user_release-1.49.1_done.md`; then this idea; then the separate windows-portability run. Keep labelled experimental/unvalidated Windows assets, hold CLI winget, and leave the skill channels unaffected. No Windows architecture work here. Full direction copied to source-context/release-order.md.

## Subsequent owner direction — 2026-09-25 global defaults

Read `source-context/owner-default-2026-09-25.md`, the owner's latest verbatim direction. It explicitly exempts this already-running idea: codex-1 remains organizer and this idea's three participants/roles stay unchanged. The mechanism must still ship UNSET as the product default. The owner has now personally chosen codex-1 as the future machine-global default implementer and directs setting it when the mechanism ships; this is post-release owner configuration, not a choice by the organizer or a shipped default. No owner-choice question remains. New-run organizer/quorum changes do not retroactively rewrite this run, existing signatures or roster snapshot.
