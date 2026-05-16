---
agent: hermes
idea: version-awareness-project-sync
round: 1
date: 2026-05-15
---

## Summary
Define a clear version model covering system installer, runtime skills, Parley CLI, project deck/protocol, and compatibility. Add CLI commands for reliable version reporting. Introduce startup checks for drift and required updates. Provide a safe sync flow for project-local state when the installer changes. Outline test and release needs.

## Proposed approach
Introduce a single source-of-truth version file in the packaged skill. Expose it via parley-deck-skill --version and parley version. Add parley-deck-skill doctor --version to report all installed markers plus compatibility status. At startup the facilitator runs a lightweight check against COOPERATION.md and local deck structure, prompting sync only on detected drift. Updates flow from Homebrew to runtime markers then to project-local copies when the protocol changes.

## Concerns / open questions
How to handle partial Homebrew upgrade failures without leaving mixed versions. Whether project-local sync should be automatic or require explicit confirmation. Best way to surface compatibility warnings to both humans and models.

## Risks
Stale markers after failed installs. Over-eager sync overwriting custom project state. Version drift between packaged reference and live COOPERATION.md.