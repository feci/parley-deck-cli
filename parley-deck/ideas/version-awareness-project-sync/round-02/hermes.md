---
agent: hermes
idea: version-awareness-project-sync
round: 2
date: 2026-05-15
responding-to: [codex/round-02, gemini/round-01]
---

## Position changes since round 1
I now support a compatibility manifest approach instead of enforcing a single lockstep VERSION source across CLI and skill.

## Responses to others

### codex
I agree with the compatibility manifest model over lockstep versioning. The manifest approach correctly allows independent release cadences while still providing clear compatibility checks. I also fully support the proposed commands: `parley-deck-skill status`, `parley-deck-skill sync-project`, and the optional `parley version --all`. These give humans and models reliable, actionable status without false unification.

### gemini
The four-component model is sound, but I disagree with using a single VERSION file as source of truth for both CLI and installer. A compatibility manifest is the safer and more flexible choice, as Codex proposed. This avoids reporting inaccuracies when components update separately. The `parley version --all` command should surface the manifest results rather than force alignment.

## New concerns / questions
None at this stage.

## Current proposal
Adopt a compatibility manifest plus the commands `parley-deck-skill status`, `parley-deck-skill sync-project`, and optional `parley version --all`. Startup checks read the manifest, report actual component versions, and warn on drift without requiring lockstep numbers.