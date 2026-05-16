---
agent: gemini
idea: version-awareness-project-sync
round: 1
date: 2026-05-15
---

## Summary
The current versioning model is fragmented, with the CLI reporting 1.0.0 while the installer and skills are at 1.0.9. This proposal introduces a unified, multi-tier versioning model to ensure that models, users, and the system can reliably determine compatibility across the system installer, CLI, runtime skills, and project-local protocol.

## Proposed approach

### 1. Unified Version Model
We will distinguish between four distinct versioned components:
- **System Installer (`parley-deck-skill`)**: Manages installations and environment health.
- **CLI (`parley`)**: The core orchestration binary.
- **Runtime Skills**: The specialized agent capabilities (Codex, Claude, etc.) installed via the system installer.
- **Project Protocol (`COOPERATION.md`)**: The local project's governance document.

### 2. Enhanced CLI Commands
- **`parley version --all`**: Provide a comprehensive report including CLI version, system installer version, and the project-local protocol "signature" (version/hash).
- **`parley-deck-skill version`**: Standardize to report both the installer version and a list of all installed runtime skill versions.
- **`parley agents verify --json`**: Include a `compatibility_status` field for each agent, comparing its reported version (via `--version`) against the expected runtime skill version.

### 3. Startup & Continuity Checks
- **Protocol Drift Detection**: On `parley` startup, the facilitator should compare the local `parley-deck/COOPERATION.md` with the "reference" version bundled with the current CLI.
- **Graceful Warning**: If a version mismatch is detected (e.g., CLI is 1.0.0 but protocol expects 1.1.0 features), the CLI should issue a high-signal warning to the model/user before starting any round.
- **Agent Skill Verification**: Before a round starts, the facilitator should verify that the selected participants' versions match the project's requirements.

### 4. Project-Local Sync Flow
- **`parley sync` Command**: An automated command to update `parley-deck/` directory structure, labels, and `COOPERATION.md`.
- **Protocol-Aware Updates**: For breaking protocol changes, `parley sync` should be able to propose a "meta-protocol-change" idea (Phase 0) to migrate the project legally under the current rules.

### 5. Deployment and Verification
- **Versioning Lockstep**: Ensure that the `VERSION` file in the CLI repository is used as the source of truth for both `parley` and the `parley-deck-skill` installer.
- **Installer Robustness**: Update `parley-deck-skill` to handle non-writable environments (like locked Homebrew Cellars) by reporting a "skew" status instead of silently failing or reporting outdated versions.

## Concerns / open questions
- How should we handle "dirty" protocol changes where a user has manually edited `COOPERATION.md`? Should we use content hashing or a semantic version in the frontmatter?
- Should the `parley` binary itself be able to perform a "self-upgrade" via the system installer, or should it only ever advise the user to run `parley-deck-skill upgrade`?

## Risks
- **Automatic Sync Hazards**: Automated updates to `COOPERATION.md` might overwrite intentional project-local customizations.
- **Model Confusion**: If a model detects a version mismatch and stops work, it might cause a "deadlock" where the model refuses to proceed even for trivial mismatches.
- **Skewed Reporting**: If the installer reports 1.0.9 but the binary is actually 1.0.8, we might create a false sense of security. The system must verify the actual executable version.
