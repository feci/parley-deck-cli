---
idea: version-awareness-project-sync
author: codex
created: 2026-05-15
participants: [codex, gemini, hermes]
roles:
  codex: facilitator and implementation-shape reviewer
  gemini: installer and packaging compatibility reviewer
  hermes: protocol and startup-workflow reviewer
status: final
---

## Problem / idea

The user wants Parley Deck to become version-aware and self-checking:

- If a user asks a model which Parley Deck version is installed, the model should be able to answer reliably.
- The system installer should be able to report its own version and the versions installed into local agent runtimes.
- When the system installation is updated, project-local Parley Deck configuration, protocol snapshots, or skill copies should also be checked and updated when required.
- When a model starts a Parley Deck workflow, it should check for structural or process drift between the system skill/installer and the local project deck before proceeding.

Observed current state:

- Source package `parley-deck-skill/package.json` is `1.0.9`.
- Runtime skill markers under Codex, Claude, Gemini, and Hermes report installed skill version `1.0.9`.
- The `parley-deck-skill` command currently resolved from PATH is Homebrew `1.0.8`, because the local Homebrew Cellar was not writable during the prior upgrade.
- `parley version` reports `1.0.0`.
- `parley-deck-skill doctor --target all --json` already exposes per-runtime marker versions, but the plain `--version` only reports the system installer package version.
- Project-local `parley-deck/COOPERATION.md` is the canonical live protocol; bundled skill snapshots are fallbacks. The current skill already does a hash drift check, but it does not define a first-class project version/status metadata model.

## Requested output

Propose a concrete implementation plan for:

1. A clear version model: system installer version, installed runtime skill version, Parley CLI version, project deck/protocol version, and compatibility status.
2. Commands or command changes so both humans and models can answer "what Parley Deck version do I have installed?"
3. Startup checks in the skill/protocol so a facilitator checks project-local deck state against the installed system skill/installer before starting Parley Deck work.
4. A safe project-local update flow when the system skill/installer changes.
5. Tests and release implications.

## Constraints

- Keep the design vendor-neutral and usable by Codex, Claude, Gemini, Hermes, and other local CLI agents.
- Preserve the live project `parley-deck/COOPERATION.md` as canonical when present.
- Do not make network access required for normal startup.
- Avoid auto-overwriting project-local protocol files without a clear plan or explicit user approval.
- Keep the implementation small and incremental.
- English-only for all files under `parley-deck/`.

## Non-goals

- Do not implement code in this design idea.
- Do not require a cloud service or central registry lookup for every startup.
- Do not replace the existing installer marker model unless there is a strong reason.
- Do not solve npm publishing permissions or local Homebrew Cellar permissions in this design.
