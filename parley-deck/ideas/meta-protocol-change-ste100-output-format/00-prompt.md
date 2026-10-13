---
idea: meta-protocol-change-ste100-output-format
author: user
created: 2026-10-13
participants: [codex-1, zcode-1]
excluded: [kimi-1 — owner-authorized dropout after two measured readiness failures; no same-idea return — confirmed 2026-10-13]
implementer: codex-1
track: deliberation
auto_implement: true
require_model_diversity: true
checks: go test ./... -count=1 -timeout 45m
status: final
---

## Problem / idea

Make ASD-STE100 the default writing standard for new English text that agents produce under the Parley Deck protocol.
Complete design, implementation, independent review, and release through the authorized channels.

Read the controlling [ORGANIZER-BRIEF.md](source-context/ORGANIZER-BRIEF.md) in full.
Its SHA256 is `8ab5c82ef9ba7295ebb939ae94ab9182ad7f2ede757fa620938d4718ab848dfa`.
The copy is byte-identical to the owner's brief.
It assigns codex-1 as organizer and implementer for this idea.
Claude takes no part.

## Owner direction

Verbatim source, in Slovak:

> "pridaj do parley-deck protokolo language format ASDSTE-100, vsetko co komunikujes, ako formulujes text a vytvaras obsah, bude to default output text format for the protocol, sprav ideu na to ako to zapracovat do protokolu a potom rovno vydaj novu verziu cez vsetky kanaly"

Translation: "Add the ASD-STE100 language format to the parley-deck protocol. Everything you communicate, how you formulate text and create content: it will be the default output text format for the protocol. Make an idea on how to work it into the protocol, and then release the new version directly through all channels."

The owner's instruction requires this verbatim quote. All new discussion uses English.

## Constraints

- Use the owner-provided CLI and skill worktrees, both on branch `ste100`.
- Target CLI 1.55.0, skill 2.19.0, and staged core 2.19.0.
- Use staged `~/.parley/staging/COOPERATION-2.18.0.md` plus exactly the reviewed protocol hunks for core 2.19.0.
- Preserve the 70,000-byte packet limit and the 20,000-byte core skill limit.
- Put detailed rules in an on-demand reference. Keep the binding pointers short.
- Verify the current ASD issue through an authoritative source. Record section 15 provenance.
- Summarize rules in original words. Do not copy the ASD dictionary or long specification passages.
- Protect quotes, code, commands, paths, identifiers, logs, errors, schema keys, frontmatter, generated text, and validator output.
- State the English-only relationship, default, any opt-out, and adoption for open ideas and older decks.
- Preserve existing artifacts and existing protocol prose except the exact reviewed hunks.
- Decide whether an advisory mechanical check earns its size and cost. A full dictionary checker is out of scope.
- Update runner and driver briefs. Use the drafted rules in this idea from round 2.
- zcode-1 must author independent reviews and a fresh goal check in separate processes. codex-1 does not review its own code.
- Apply the brief's Kimi exclusion, close authority, five-cycle limit, retry rules, release channels, and final handoff duties exactly.
- Do not apply the D6 legacy declaration. Do not declare or prune worktrees. Do not change credentials or gateway settings.

## Startup / selection checkpoint

Transport is `github-pr`. The owner's named branch overrides the standard branch-name convention.
All participant host identities map to `feci`; canonical signoffs remain authoritative where GitHub refuses self-approval.

| Agent | CLI | Configured model | Effort / speed | Task ceiling | Write mode |
| --- | --- | --- | --- | --- | --- |
| codex-1 | Codex 0.162.1 | gpt-6-astra | max / deep | 1800 s | workspace-write, approval never |
| kimi-1 | Kimi 0.42.0 | kimi-code/k3 | max / deep | 2400 s | print mode |
| zcode-1 | Zcode 3.14.5-34 | zai/glm-5.3 | max configured / deep | 1800 s | yolo, workspace cwd |

The Zcode model is configured authority. It is not native observed-model evidence.
CLI 1.54.0 and the installer plus eight detected runtimes use skill 2.18.0.
Source metadata is valid and matches the packaged protocol.

`parley preflight --no-ping --json` reports the unrelated historical devx-speed facilitator declaration.
The command has no participant filter. Separate measured readiness probes avoid a Claude task launch.
Readiness receipts will record actual requested-peer results before Phase 0 closes.

Full live protocol context was read.
`context_mode: full`; source and packet SHA256: `0357d504982f92b713f2f86604129f46e1ef3276664e4da91d39ecf420253f0c`.
There is no fallback reason. The complete attestation is in `source-context/phase0-packet.json`.
Shared project memory and the current driver-unstall release handoff were consulted.
No graphify graph exists. Targeted source navigation avoids an unrelated graph-build project.
There are no open CLI PRs. Historical unrelated ideas do not change this explicit assignment.

## Non-goals

No full dictionary, certification claim, historical rewrite, broad driver repair, new roster member, or CLI WinGet publication.
No npm or core publication through an emulated attended terminal.
The final note must list the still-pending owner commands in the required version order.

## Phase 0 readiness result

Kimi had two failed measured readiness processes. Both exited 1.
The brief pre-authorizes its exclusion for this idea. No global roster changed.
Zcode returned PONG and exited 0. Codex is this active process.
The original proposed three-member prompt remains in commit `116e64c`.
This manual owner-authorized exclusion supplies no automatic reviewer-exception evidence.
The brief's conditional close authority remains available.
