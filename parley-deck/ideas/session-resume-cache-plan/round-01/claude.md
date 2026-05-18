---
agent: claude
idea: session-resume-cache-plan
round: 1
date: 2026-05-18
---

## Summary

Treat resume as two cooperating layers: a **parley workflow manifest** (what step is next, what artifact is missing) and **per-step pointers to each native agent CLI's own session** (codex/claude/gemini/hermes session IDs and jsonl paths). The manifest lives in `~/.parley-deck/runs/<run-id>/` (rebuildable cache) with a thin canonical mirror in `parley-deck/runs/<run-id>/` (survives cache wipe). Canonical idea artifacts under `parley-deck/ideas/...` remain the source of truth; reconciliation between artifacts on disk and manifest status drives "what's left to do." Resume = reconcile + re-invoke the next missing step, preferring the agent's native `--resume` path and falling back to a fresh session seeded with a cached prompt.

## Proposed approach

### Storage layout

`~/.parley-deck/` (local cache; respect `PARLEY_DECK_HOME` env override):

```
config.json
runs.json                       # v1 index; migrate to index.sqlite if/when it slows
runs/<run-id>/
  manifest.json                 # plan + step statuses, schema_version
  heartbeat.json                # pid, step, last_tick_at
  workspace.json                # absolute repo path, idea slug, transport
  agents/<round-NN>.<agent>.json  # { cli, session_id, session_path, started_at, ended_at, exit_code }
  input-packs/<round-NN>.<agent>.prompt.txt   # exact rendered prompt
  logs/<round-NN>.<agent>.stderr              # capped, diagnostic
cache/transcripts/              # optional opportunistic copy of native jsonls at end-of-step
```

`parley-deck/runs/<run-id>/` in the repo (canonical, lightweight):

```
events.jsonl                    # already exists, keep
run.json                        # NEW: { run_id, idea, started_at, transport, schema_version }
agents.json                     # NEW: same agent→session pointer map as cache (so cache is rebuildable)
```

Keeping the agent pointer map in the repo is the load-bearing choice: if `~/.parley-deck` is wiped or the user moves machines (with the repo), we can rebuild the index with `parley-deck sessions rebuild --from-workspace`. Without it, resume after cache loss collapses to "start the next step from a fresh agent session," which is still survivable but loses native-CLI conversation continuity.

### Workflow manifest

`manifest.json` enumerates the planned steps from `COOPERATION.md` for this idea: idea slug, transport, participants, rounds, and for each step `{ id, kind (round/review/consensus), round, agent, expected_artifact_path, status, started_at, ended_at, native_session_ref, prompt_hash }`. Status is the small set `pending | in_progress | succeeded | failed | timed_out`. Reconciliation always trusts the filesystem: if `expected_artifact_path` exists and parses, the step is `succeeded` regardless of stored status (handles "process died right after writing the file").

### Native CLI integration

One adapter per CLI, all conforming to a small interface `{ start(prompt, opts) -> session_ref, resume(session_ref, prompt?) -> session_ref, locate(session_id) -> path }`:

- **Codex**: capture `session_id` from `codex` stdout/`~/.codex/session_index.jsonl`; resume via `codex resume <id>` (or fork).
- **Claude**: pass `--session-id <uuid>` on start so we own the ID; resume via `claude --resume <id>`; locate file under `~/.claude/projects/<escaped>/`.
- **Gemini**: pass `--session-id` on start; resume via `--resume`; list via `--list-sessions` as a sanity probe.
- **Hermes**: capture id from `hermes sessions list`; resume via `--resume <id>`.

Where the CLI lets us set the session ID, do so (deterministic mapping parley-run-id × step → native id). Where it doesn't, capture it from output and persist immediately.

### Resume UX

- `parley-deck sessions list` — read `runs.json`; show `run-id | idea | workspace | last-activity | status | next-step`.
- `parley-deck sessions inspect <run-id>` — manifest + per-step native-session resolvability (`session file exists?`, `last touched`).
- `parley-deck sessions resume <run-id|--last>` — reconcile, print the next missing step, confirm, then re-invoke. `--yes` for non-interactive. `--from-step <id>` to retry a specific step (e.g., a single failed agent in round 2 without rerunning round 1).
- `parley-deck sessions rebuild --from-workspace <path>` — walk `parley-deck/runs/*/run.json` + `agents.json` to reconstruct `~/.parley-deck` index after cache loss.
- `parley-deck sessions rebind <run-id> --workspace <path>` — fix stale absolute paths when the user moves the repo.

### Liveness, locks, concurrency

- Each in-progress invocation writes `heartbeat.json` every ~5 s; on resume, stale heartbeat (>30 s with no live pid) flips status to `failed`/`timed_out` so the step is re-invocable.
- `flock` on `manifest.json` during mutation; refuse to resume a run already held by another parley process; warn on filesystems where advisory locks are unreliable (iCloud/Dropbox/NFS) by sniffing the mount.

### Re-invoke without restarting the idea

Re-invocation is per-step. The step's cached prompt is replayed against either the existing native session (best effort) or a fresh one. We hash the rendered prompt at creation time and again at re-invoke; on hash mismatch we surface a "prompt drift" warning so the user can choose `--use-cached` vs `--rerender`. This makes "agent died mid-round-02" recoverable without re-running round-01.

## Concerns / open questions

- **TUI binding**: the workspace sessions console added in `2be9df1` likely reads the existing `~/.parley-deck/sessions.json`. We need to decide whether to (a) extend that file in place with new fields, (b) introduce `runs.json` alongside and migrate, or (c) make the TUI read the manifest directly. I lean (b) with a one-shot migration on first launch.
- **Native session retention**: Codex/Claude/Gemini/Hermes each have their own pruning policies. We can't guarantee the native jsonl will exist a week later. Promise should be "step-level re-invoke with same prompt"; native conversation continuity is *best effort* and we say so in `--help`.
- **Transports (PR/MR)**: for GitHub/GitLab transports, the PR/MR is part of canonical state. Does resume re-fetch PR comments before deciding the next step, or trust local artifacts? v1: trust local; add `--refresh-transport` flag later.
- **Identity of a "run"**: is `run-id` per idea, per `(idea, branch)`, or per invocation? Per-invocation is simplest but produces many resumable runs for the same idea; users likely want "resume idea X" not "resume run 7af3". Suggest secondary index `idea -> latest_run_id`.
- **Multi-machine**: pointer map in the repo lets a user clone elsewhere, but the native session jsonls don't follow. Acceptable v1 limitation; document it.
- **Schema versioning**: manifest needs `schema_version` plus a migration step on load. Worth scoping a `migrations/` module from day one even if empty.

## Risks

- **Native CLI drift**: any of the four agent CLIs may change session paths/commands; each adapter needs a startup probe and a clear failure mode ("agent X doesn't expose --resume in this version → fresh-session fallback").
- **Stale absolute paths**: workspace moves break the index; mitigated by `rebind`, but easy to miss for users.
- **Cache wipe combined with native pruning**: if both happen, we lose conversation continuity entirely and can only fresh-invoke from prompt. Need to communicate this clearly so users don't expect magical recovery.
- **Concurrent processes**: TUI and CLI both resuming the same run is the most likely foot-gun; locking helps but lock-on-NFS-style filesystems may silently no-op.
- **Prompt drift across parley-deck versions**: cached prompts may not match current templates; surfacing drift as a warning rather than silently choosing is important.
- **SQLite-vs-JSON premature optimization**: starting with `runs.json` is simpler; risk is that under hundreds of runs or concurrent writers we'll regret it. Cheap to migrate later if we keep the access behind one `RunsIndex` interface.
- **Heartbeat overhead**: writing every 5 s × parallel agents is fine on SSD, painful on synced cloud directories. Make interval configurable and disable when `PARLEY_DECK_HOME` is under a known sync mount.
