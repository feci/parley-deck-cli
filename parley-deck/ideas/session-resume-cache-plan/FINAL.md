---
idea: session-resume-cache-plan
status: final
author: codex
consensus-date: 2026-05-25
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

Parley Deck session resume support must prioritize workflow-level recovery over
native model conversation continuity. The reliable promise is: Parley can reopen
an older run, inspect canonical artifacts and local run state, determine the
next safe step, and retry or resume only the missing participant action. Native
agent session continuation is opportunistic and reported per agent.

### Core invariants

- Canonical outcomes stay under the workspace `parley-deck/` directory.
- Sensitive local state, including full rendered prompts and native session
  handles, stays in local cache such as `~/.parley-deck` or another explicit
  local Parley home.
- Committed workspace metadata may include prompt hashes, redacted launch
  configuration, bounded diagnostics, expected artifact paths, validation
  results, and schema versions.
- Do not depend on private external CLI cache formats as the only source of
  truth.
- `headless-agents.local.json` is uncommitted machine-local configuration. It
  may describe models, launch flags, resume capabilities, and capability tiers.
- Every durable schema introduced by this roadmap must include a version field
  from the first release.

### Roadmap

#### Slice 1 - complete

The existing implementation is accepted as the foundation:

- repository-local `parley-deck/runs/<run-id>/run.json`;
- read-only `parley sessions list`;
- read-only `parley sessions inspect`;
- graceful handling of legacy runs where `run.json` is missing.

#### Slice 2 - per-agent attempts and prompt hashes

Add attempt capture for each agent invocation. The implementation should record:

- agent ID;
- selected model/profile/effort;
- redacted CLI launch configuration;
- prompt SHA-256 hash and prompt byte length;
- expected artifact path;
- process duration and exit status;
- bounded stdout/stderr diagnostics with redaction;
- validation result for the expected artifact;
- minimal run status updates needed for coherent `sessions inspect` output.

Full prompt bodies should be cached locally, not committed by default.

#### Slice 3 - native resume handles and capability config

Add explicit resume capability configuration for Codex, Claude, Gemini, and
Hermes. Capability tiers should distinguish at least:

- `none`;
- `interactive-only`;
- `headless-session`;
- `fork`.

Capture native session handles only when the launched CLI exposes or accepts a
stable handle. Store sensitive handle values in local cache. Workspace metadata
may state that a handle exists and which capability tier applied.

#### Slice 4 - cache rebuild and workspace rebind

Add commands to rebuild local Parley state from workspace run manifests and to
repair stale workspace paths after a repository move. Cache loss must degrade to
prompt replay or fresh agent invocation, not corrupt canonical artifacts.

#### Slice 5 - TUI resume, retry, and manual handoff

Expose recovery actions in the TUI and CLI after the lower-level attempt,
capability, and cache-rebuild data is reliable. The UI should show whether an
agent can native-resume, fresh-retry, or needs manual/user handoff.

### Deferred details

- Exact attempt record path and retry naming.
- Bounded-log size and redaction policy.
- Locking and stale-heartbeat behavior for concurrent TUI/CLI access.
- Exact rebuild/rebind command names.
- Whether a workspace-local ignored cache is useful in addition to
  `~/.parley-deck`.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/, ./round-02/, ./round-03/
- Implementation slice 1: ./IMPLEMENTATION.md
- Historical review exception: ./review/consensus.md
