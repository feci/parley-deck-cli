---
agent: antigravity-1
idea: meta-protocol-change-preflight-readiness
round: 1
lens: docs / consistency / UX
date: 2026-06-19
---

# Pre-idea Readiness Check — Round-01 Analysis (Docs / Consistency / UX)

## Summary
This document analyses the mechanics of the proposed pre-idea readiness check from a documentation integrity, project-wide consistency, and operator user experience (UX) perspective. While the goals and core decisions of **auto-freshness** and **roster liveness pings** are locked, the operational details require careful framing. We identify and resolve several critical edge cases: preventing the regression of source repositories, ensuring the clear recording of exclusions in both machine-readable and human-legible formats, preventing silent or mid-idea quorum expansion, and structuring the `parley preflight` CLI interface to support clean interactive prompts and unattended automated runs.

---

## Refinement of Mechanics

### 1. Freshness Sync & its Records (Topic A)

#### Source vs Consumer Repository Identification
To prevent the "source-vs-consumer" inversion where the local protocol is ahead of the published package version (as is the case in this repository), the tool must identify the project's role.
- We refine the `protocolRole` setting in [version.json](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/meta/version.json) to support `source` and `consumer`.
- If the field is missing, it defaults to `consumer`.
- When set to `source`:
  - The freshness check runs in **advisory-only** mode.
  - If drift is detected between `parley-deck/COOPERATION.md` and the embedded [defaults/COOPERATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/defaults/COOPERATION.md), the tool prints a warning and displays a dry-run diff but never performs a sync or blocks execution.

#### Preservation of Project-Specific Zones
When performing a sync in a `consumer` repository, the tool must isolate custom project zones. We standardise the preserved zones to align with the normalizer in [TestEmbeddedDefaultMatchesLiveDeck](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/drift_test.go#L46):
1. **Header Metadata Block** (lines 1 to 10 containing `Workspace:`, `Parley deck:`, `Transport:`, etc.)
2. **Section 0** (`## 0. Choose the transport`)
3. **Section 2** (`## 2. Active agents (roster)`)

#### In-Place Synced Header Format
The sync tool updates the header in [COOPERATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md) in-place:
```markdown
**Protocol synced:** YYYY-MM-DD — parley-deck-skill <version> (<facilitator-agent-id>)
```
In a `source` repo, the embedded reference template [defaults/COOPERATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/defaults/COOPERATION.md) must omit this line. This is handled gracefully by [TestEmbeddedDefaultMatchesLiveDeck](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/drift_test.go#L46) which explicitly tests that the default template contains zero sync provenance lines.

#### Synced Record Meta Document Template
When a sync finishes, a record is saved to `parley-deck/meta/protocol-sync_<YYYY-MM-DD>_v<version>.md`. If multiple syncs happen on the same day, a suffix is appended: `_N.md`. The template is:

```markdown
# Protocol sync — parley-deck-skill <version>

**Date:** YYYY-MM-DD
**By:** <facilitator-agent-id>
**Type:** [Additive | Breaking (operator confirmed)]
**Previous Version:** <old-version>
**Previous Protocol SHA:** <old-sha256>
**New Protocol SHA:** <new-sha256>

## Summary of Changes
<A short description of text additions, modifications, or rules changed.>

## Applied Diffs
```diff
<diff output of the changes made to COOPERATION.md, excluding preserved zones>
```

## Verification
- Checked that header metadata, transport selection, and roster tables were preserved.
- Status verification: `parley-deck-skill status` returned `metadataStatus: valid`.
```

---

### 2. Operator UX of Confirmation Gates (Topic B)

The liveness ping runs at Phase 0. We define the console interaction for the two gates.

#### Gate 1: Exclude Gate Prompt (Unavailable Agent)
If a rostered agent fails to respond to the pong probe within a 5.0-second limit:
```
⚠️ Warning: Agent 'hermes-1' failed to respond to the liveness ping.
Reason: Connection timed out (5s limit).

Would you like to exclude 'hermes-1' from this idea ('meta-protocol-change-preflight')?
Excluding them will temporarily adjust the quorum to 3/4 agents.

[y] Yes, exclude 'hermes-1' for this idea (temporary)
[n] No, abort preflight (keep trying to connect)
```

#### Gate 2: Re-include Gate Prompt (Previously-Excluded Agent Online)
If an agent was excluded in a previous idea but is now responding online at the start of a new idea:
```
ℹ️ Info: Agent 'hermes-1' (previously excluded in 'some-prior-idea') is now online.

Would you like to re-include 'hermes-1' in the quorum for this new idea?

[y] Yes, re-include them (quorum expands back to 4/4 agents)
[n] No, keep them excluded for this idea
```

#### Recorded Shape in [00-prompt.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/meta-protocol-change-preflight-readiness/00-prompt.md)
To ensure the audit trail is both machine-parsable and human-legible, exclusions are written to two locations in `00-prompt.md`:

1. **YAML Frontmatter** (for programmatic parsing during rounds):
```yaml
excluded:
  - agent: hermes-1
    reason: connection timeout
    confirmed: 2026-06-19
```

2. **Body Markdown** (under a new standard `## Readiness` section at the bottom of the file):
```markdown
## Readiness
- **Protocol freshness:** Consumer deck up-to-date with parley-deck-skill 1.4.0.
- **Roster ping:**
  - claude-1 ✅ (facilitator)
  - codex-1 ✅ PONG
  - hermes-1 ❌ UNAVAILABLE (reason: connection timeout — operator confirmed exclusion 2026-06-19)
  - antigravity-1 ✅ PONG
```

#### Quorum Lock-In Rule (Mid-Idea Invariant)
To maintain consistency across rounds:
- **Rule:** The quorum for a specific idea is fixed once Phase 0 is complete and `00-prompt.md` is committed. No agent can be excluded or re-included mid-idea.
- If an excluded agent comes back online during an idea, they cannot write files or vote. If the operator wants them to participate, they must restart Phase 0.
- If an active agent goes offline mid-idea, the existing async timeout escalation rules (§5) apply.

---

### 3. 'parley preflight' CLI Output (Topic C)

#### Human-Readable Output (Stdout)
```
================================================================================
                           PARLEY PREFLIGHT REPORT
================================================================================

[+] Protocol Freshness Check
    Role: Consumer
    Local Version: 1.3.1
    Installed Version: 1.4.0
    Status: Drifted (Additive changes only)
    Action: Automatically synced COOPERATION.md (Roster and Transport preserved).

[+] Roster Liveness Ping Check
    -------------------------------------------------------------------------
    Agent ID        Role                     Ping    Latency    Disposition
    -------------------------------------------------------------------------
    claude-1        facilitator+participant  OK      12ms       Included
    codex-1         participant              OK      45ms       Included
    hermes-1        participant              FAIL    Timeout    Excluded (Operator)
    antigravity-1   participant              OK      8ms        Included
    -------------------------------------------------------------------------

[+] Summary
    Status: READY
    Quorum Adjusted: 3/4 participants active.
================================================================================
```

#### JSON Output Shape (`--json`)
```json
{
  "timestamp": "2026-06-19T10:17:02Z",
  "ready": true,
  "freshness": {
    "role": "consumer",
    "status": "synced",
    "compatibility": "additive",
    "localVersion": "1.4.0",
    "installedVersion": "1.4.0",
    "localSha256": "59905d3df97b22b960f16f28e3f4edec2c444c1d50ae75185e07ff1169f8f89e",
    "installedSha256": "59905d3df97b22b960f16f28e3f4edec2c444c1d50ae75185e07ff1169f8f89e"
  },
  "roster": [
    {
      "agentId": "claude-1",
      "role": "facilitator+participant",
      "status": "online",
      "latencyMs": 12,
      "quorumAction": "include"
    },
    {
      "agentId": "hermes-1",
      "role": "participant",
      "status": "offline",
      "error": "connection timeout",
      "quorumAction": "exclude"
    }
  ],
  "quorumCount": 3,
  "totalRosterCount": 4
}
```

#### Exit Codes
- `0`: Success / Ready (all checks passed; or additive sync succeeded).
- `1`: Invocation or system error (e.g. invalid arguments, missing workspace directory).
- `2`: Protocol freshness gate blocked (breaking changes detected, operator rejected sync, or sync failed).
- `3`: Roster availability gate blocked (roster ping failed/re-include found, operator rejected exclusion/re-inclusion).
- `4`: Invariant block (the adjusted quorum violates §1 non-solo execution).

#### Unattended / Automation Mode
To prevent deadlocking in CI pipelines or during unattended runs:
- Running with `--yes` or `--non-interactive` disables all interactive prompts.
- If a gate requires user input (breaking protocol bump, unavailable agent, or re-inclusion candidate), the CLI must immediately exit non-zero (Exit Code 2 or 3) rather than waiting or applying silent defaults.

---

## Invariant Alignment & Rationale

- **§9 Session-start Checklist Alignment:**
  Step 8 is modified to state:
  > 8. **Run `parley preflight`** at the beginning of every session to verify protocol freshness and ping the roster. The facilitator runs this command automatically at idea kickoff.
- **§1 Non-solo Requirement Alignment:**
  If the exclude gate would result in all non-facilitator agents being excluded, the preflight tool must prevent the exclusion. It must fail with Exit Code 4 and instruct the operator to either fix the agent connection or write a manual solo exception to `parley-deck/inbox/`.
- **§5 Quorum Alignment:**
  We add a sentence to §5:
  > Quorum is determined at Phase 0 by adjusting the roster based on the liveness check. Excluded agents do not count towards the quorum of the current idea. The quorum is locked once Phase 0 is complete.
- **§7 Carve-out:**
  We add a paragraph to §7:
  > Carve-out: Automatically applying an upstream-ratified protocol sync does not require a protocol-change idea, provided the sync is compatible/additive and preserves project-specific zones.

---

## Risks

| Risk | Impact | Mitigation |
| :--- | :--- | :--- |
| **False-positive Agent Failures:** Network spikes or slow model startup cause liveness pings to time out, triggering false exclusions. | High | Set a configurable ping timeout (default 5.0s) and allow the operator to retry before choosing to exclude. |
| **Branch Inconsistencies (B/C Transports):** One agent syncs the protocol on a local branch, but another agent on a different branch is still on the old version. | Medium | Git-based transports pull remote updates before checking freshness, and the preflight tool validates that the workspace has no uncommitted protocol modifications. |
| **Drift Normalization Divergence:** A sync overwrite changes files in a way that breaks local code references or CLI behavior. | High | Test protocol changes against both [TestEmbeddedDefaultMatchesLiveDeck](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/drift_test.go#L46) and CLI behavior before publishing upstream. |
