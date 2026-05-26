---
agent: claude
idea: session-resume-cache-plan
round: 3
date: 2026-05-25
responding-to: [claude/round-02, codex/round-02, gemini/round-02, hermes/round-02]
---

## Position changes since prior round

None. The 5-slice roadmap that emerged in round 2 is the right plan. All four participants now agree on scope, ordering, and the core framing (workflow recovery promised; native conversation continuity best-effort). Ready to move to consensus.

## Responses to others

### @codex

Agree on all points in your round-03 artifact.

**Status updates folded into slice 2:** Correct call. A dedicated status-only slice would be under-scoped — the status transitions only become meaningful once there are attempt records to contextualize them. Folding minimal status writes into slice 2 keeps the slice self-contained without inflating it.

**Prompt sensitivity:** Strong agree. Full rendered prompts belong in `~/.parley-deck/runs/<run-id>/input-packs/` (local cache, never committed). Workspace attempt records should carry only: prompt SHA-256 hash, prompt byte length, model/profile, CLI invocation (with secrets redacted), expected artifact path, exit code, bounded stderr tail, and validation result. This keeps repo diffs clean and avoids leaking user instructions or API keys that may appear in rendered prompts.

**`headless-agents.local.json` uncommitted:** Agree. It is machine-local configuration — different operators may have different CLI paths, versions, or auth contexts. A `.gitignore` entry and a `headless-agents.local.example.json` template in the repo is sufficient.

### @gemini

Your round-02 pivot from a standalone registry to enriching existing agent config was the right move. The `capabilities.resume` block in `headless-agents.local.json` is a clean way to express per-agent resume semantics without a parallel source of truth. One refinement: the block should also declare the agent's capability *tier* from codex's round-01 matrix (`none | interactive-only | headless-session | fork`) so the facilitator can decide strategy before attempting invocation, not after a failed probe.

Your detailed slice-2 attempt record schema (`promptHash`, `config`, `cli`, `expectedArtifact`, `result`) is sound. It aligns with what codex and I proposed independently. The file-naming convention (`<round>-<agent-id>-<timestamp>.json`) is reasonable; I suggest dropping the timestamp suffix in favor of `<round>-<agent-id>.json` since a given agent runs at most once per round. If retries exist, append `-attempt-N`.

### @hermes

Agree on readiness for consensus. Your round-02 suggestion to add a minimal status updater is addressed by codex's proposal to fold it into slice 2 — no separate slice needed, but the concern was valid and is now covered.

Your round-01 concerns about registry races and multi-machine rebuild remain relevant for slices 4–5. They do not block consensus on the roadmap shape, but they should be called out as open design questions in the consensus document so they are not lost.

## New concerns / questions

- **Attempt record naming:** Gemini proposed `<round>-<agent-id>-<timestamp>.json`; I prefer `<round>-<agent-id>.json` with an `-attempt-N` suffix only on retry. Minor, but worth settling before implementation to avoid a rename later.
- **Consensus scope:** This idea's consensus should cover the roadmap (5 slices, ordering, boundaries) and the core invariants (workflow recovery promised, conversation continuity best-effort, prompts local-only, capabilities configured not discovered). It should NOT lock down internal schemas for slices 3–5 — those should be designed when their slice begins.

## Current proposal

Move to consensus on the following roadmap:

1. **Slice 1 — complete.** Repo-local `run.json` (schema_version 1), `parley sessions list`, `parley sessions inspect`. Read-only, no schema changes to existing files, graceful legacy-run handling.
2. **Slice 2 — per-agent attempt records + prompt capture.** Write attempt metadata to workspace (`<round>-<agent>.json` under `attempts/`), full prompts to local cache only. Fold minimal `run.json` status transitions into this slice. Prompt integrity via SHA-256 hash.
3. **Slice 3 — native resume handles + capability config.** `headless-agents.local.json` gains a `capabilities.resume` block with tier classification. Facilitator captures native session IDs into attempt records when the agent supports it. Resume is opportunistic, never required.
4. **Slice 4 — cache rebuild/rebind.** `parley sessions rebuild --from-workspace` reconstructs local index from repo artifacts. `parley sessions rebind` fixes stale workspace paths. Design must address hermes's multi-machine and registry-divergence concerns.
5. **Slice 5 — TUI resume/retry/handoff.** Expose actions in the run console. Consumed slices 2–4 state. Design deferred until lower layers are stable.

**Core invariants for consensus:**
- Parley promises workflow-level recovery and per-step retry. Native conversation continuity is best-effort, reported per agent, never required for correctness.
- Full prompt content is local cache only; workspace records carry hashes and bounded metadata.
- Agent capabilities are explicitly configured, not discovered by parsing help text.
- `headless-agents.local.json` is local, uncommitted configuration.
- Schema versioning (`schema_version` field) is mandatory on all new durable artifacts from slice 1 onward.

No blockers. Ready for FINAL.md.
