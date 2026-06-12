---
agent: claude
idea: runner-hardening-kindly
round: 1
date: 2026-06-12
---

## Summary

Adopt kindly's supervision/classification/snapshot machinery into the runner with
parley-native shapes: watchdogs ride the existing per-agent log files and procctl
kill-tree; classification enriches `agent.failed` payloads; artifact-beats-exit
removes the agy flake; snapshots give Phase 6 a physically pinned review root on
LOCAL tmp (off virtio-fs); `parley consult` reuses the single-agent invocation
path with a new durable artifact class. No new tick loops; no persisted
heartbeat events (the 1.23.0 TUI growth cache already covers liveness display).

## P1 — Watchdogs

- **Where:** the runner's existing wait path for exec-mode agents; ACP mode uses
  event arrival instead of file growth.
- **First-output watchdog:** no growth on stdout+stderr (exec) / no acp event
  (ACP) within `first_output_secs` (default 120; agents.toml per-agent override;
  0 disables) → kill the process tree via the existing durable-kill machinery,
  retry ONCE (same artifact path, Overwrite=false semantics), then emit
  `agent.failed` with `failure_class: "no_first_output"`. This catches the
  hermes silent-death mode in 2 minutes instead of the 30m timeout.
- **Stall guard:** after first output, zero byte-growth for `stall_secs`
  (default 1800, clamped to < TimeoutMS) → kill + `agent.failed` with
  `failure_class: "stalled"` and a diagnostics tail in the payload. Growth
  probing = stat of the two log files on a coarse cadence (≥5s), piggybacked on
  the existing wait loop — no new goroutine storms.
- **No persisted heartbeats.** kindly's heartbeats are stderr prints for a
  human supervisor; our supervisor is the TUI, which already renders liveness
  from the growth cache + glyphs. Persisting heartbeat events would bloat
  events.jsonl (~30 events/agent/run) for no consumer. If reviewers disagree,
  the fallback is a single `agent.progress` event at most every 5 minutes.
- **Config surface:** `first_output_secs`, `stall_secs` on agents.toml override
  + Spec defaults; both recorded in run.created runtime payload.

## P2 — Failure classification + recovery hints

- New `internal/runner/failclass.go`: an ordered table of bounded regex
  classifiers (rate-limit, overloaded/5xx, auth, billing, invalid-request,
  model-not-found, context-window, sandbox, budget — seeded from
  reference/kindly-agent.sh `stderr_failure_details`/`stderr_recovery_hint`),
  applied to the LAST 4 KiB of stderr + the ExitError text. Output:
  `failure_class` + `recovery_hint` strings.
- Attached to `agent.failed` (and the P1 classes above). `runstate.SummarizeEvent`
  includes the class; the TUI narrator/header shows the hint
  ("hermes: auth/billing — run hermes auth").
- Unmatched → class "unknown", no hint. Table is data, trivially extensible.

## P3 — Usable artifact beats exit code

Decision table (artifact-bearing phases: round, review, fixup, implementation):

| artifact validates | exit | result |
|---|---|---|
| yes | 0 | agent.finished (today) |
| yes | non-0 | **agent.finished + `agent_exit: N` + narrator note** (changed) |
| no | any | agent.failed (today; now with failure_class) |

Steer replies (no artifact contract) keep current semantics. The driver's
roundComplete/ReviewRoundComplete gates already key on artifact validation, so
the change is confined to the event-type decision in runner.go (and its fixup
twin in phase58.go). Removes the recurring agy "wrote artifact, exit 1" retry.

## P4 — Small hardening batch

- **Marker shedding:** when the spawned participant is the `claude` CLI, build
  cmd.Env without CLAUDECODE, CLAUDE_CODE_SESSION_ID, CLAUDE_CODE_ENTRYPOINT,
  CLAUDE_CODE_ENABLE_TASKS, AI_AGENT (both exec and ACP paths). A participant
  must not inherit the host-session identity.
- **GIT_OPTIONAL_LOCKS=0** on every read-only git probe we spawn (driver
  gitTreeClean's two commands; any status/diff probes) — probes must never
  write `.git` on the weakly-coherent mount.
- **Docs:** new `docs/agent-cli-mechanics.md` with the verified mechanics per
  roster CLI (codex `</dev/null` + `-o`; claude `-p` binding, `--tools` removes
  vs `--allowedTools` pre-approves, MCP bypasses --tools, cwd-scoped resume;
  agy value-taking `--print`, buffers_stdout; hermes `-z` + silent-death mode),
  referenced from the skill. **Do NOT switch codex capture to `-o`:** kindly's
  report-on-stdout model differs from our artifact-file contract; `-o` would
  capture the final message, not the artifact. Document it; don't adopt it.

## P5 — Snapshot checkout isolation (Phase 6)

- **Creation:** per review round, `git clone --shared --no-checkout <repo>
  <localtmp>/parley-snapshots/<repoHash>-<runID>` + detached checkout of the
  reviewed commit. The clone lives on the LOCAL temp filesystem (os.TempDir()),
  NOT the virtio-fs mount — worktree reads become local and fast; the alternates
  link reads objects from the original .git (read-only).
- **Artifact flow:** reviewers write their artifact INSIDE the snapshot at the
  usual relative path (sandboxes scope writes to cwd), and the runner moves it
  to the canonical deck path after validation (copy+rename across filesystems,
  fsutil for dirs). Validation runs on the canonical copy.
- **Lifecycle:** pid-marker file like kindly; teardown on round end; sweep of
  stale snapshots (dead pid) at the next review; `review.snapshot_created` /
  `review.snapshot_fallback` events. Any creation failure → loud fallback to
  the live tree (current behavior).
- **Scope:** review rounds only. Fix-up and implementation keep the live tree
  (they must mutate it). The reviewed-commit pin in review/consensus frontmatter
  becomes physically enforced.

## P8 — `parley consult`

- CLI: `parley consult <agent> "<question>"` (or `-` for stdin). Reuses the
  single-agent invocation machinery (spec resolution, timeouts, env isolation,
  P1 watchdogs, P2 classification) outside any run.
- Artifact: `parley-deck/consults/<UTC-ts>-<agent>-<slug>.md` with frontmatter
  {agent, model, date, question, elapsed} + the answer body; the body also
  prints to stdout. Append a provenance line to `parley-deck/consults/index.jsonl`.
- Consult prompt preamble: advisory, read-only intent, "lead with your
  recommendation", never a pass/fail verdict (adapted from kindly's
  consult_prompt). The agent writes its answer to the artifact path (our agents
  write files reliably; no stdout-capture contortion), EXCEPT agy-style
  buffers_stdout agents where stdout capture is the fallback if the file is
  missing but stdout is non-empty.
- Standing: consults are advisory, non-canonical, never quorum inputs — one
  clarifying sentence belongs in the protocol (flagging to the sibling idea's
  participants; if they decline, it lives in docs + the artifact header).

## Slices (one release, 1.24.0, together with the sibling protocol idea)

1. P4 (markers, locks, docs) + P2 (failclass) — small, independent, testable.
2. P1 watchdogs (uses P2 classes).
3. P3 decision-table change + tests.
4. P5 snapshots.
5. P8 consult.

## Risks

- P1 false stalls on legitimately silent deep runs → stall window stays large
  (30m), clamped under TimeoutMS, per-agent override, buffers_stdout agents get
  first-output measured on stderr too (agy's stderr is live).
- P5 disk usage on local tmp → shared clone keeps objects in the origin; only
  the worktree materializes; sweep on next run.
- P3 masking real failures → only when the artifact VALIDATES; exit code is
  preserved in the event and surfaced by the narrator.
