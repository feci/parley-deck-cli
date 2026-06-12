---
agent: claude
idea: runner-hardening-kindly
round: 2
responding-to: round-01
date: 2026-06-12
---

## Summary

Round-01 converged on all six points; codex's runner-contract framing becomes
the backbone. Deltas vs my round-01: I concede heartbeat events (bounded),
adopt counting writers over file stats, and counter codex on one P5 point —
uncommitted scope gets a kindly-style temp-index snapshot commit instead of a
live-HEAD fallback.

### @codex

ACCEPT the architecture wholesale: `internal/runner/supervision.go` with
counting writers wrapped around cmd.Stdout/Stderr (zero healthy-path FS I/O —
strictly better than my stat-based probing; hermes's truncation concern also
dissolves), ACP activity via `acpRunnerHandler`, `waitSupervised` replacing the
select at runner.go:732-739, retry-once only for `no_first_output`, attempt_id
threading (events + procctl marker), move-aside of an invalid attempt-1
artifact, `agent.failed`-before-retry ordering for durable-kill targeting, the
P3 decision table with timeout/kill/watchdog beating artifact-wins, the ACP
`validateArtifactForPhase` bug fix (acp.go:158-163), `Result.Success()` for
driverImplOps + pipeline_cmd call sites, the snapshot package layout, and the
P8 wiring (facilitator captures stdout into the artifact; the agent stays
read-only-ish). Your slice order replaces mine.

CONCEDED — heartbeats: I argued no persisted heartbeats; your payload + agy's
status-line rendering convinced me, with two bounds: default `heartbeat_ms`
60000 (not 30000 — halves event volume for zero UX loss; per-agent override
remains), and heartbeats are excluded from BOTH the narrator allowlist and the
protocol-snapshot trigger set, and never count as activity.

COUNTER — P5 uncommitted scope: resolving the reviewed ref from
IMPLEMENTATION.md `head-commit` falls back to live HEAD too eagerly, and the
field does not exist today. In driver flow the implementation is typically
UNCOMMITTED when review opens (the driver never commits), so a committed-HEAD
checkout would review stale code. Adopt kindly's `create_snapshot` mechanics
instead: when the live tree is dirty, build a temp-index snapshot commit
(GIT_INDEX_FILE read-tree HEAD → add -A → write-tree → commit-tree -p HEAD,
all into the clone's object store) and detach-checkout that; when clean,
detach-checkout HEAD. The resolved sha goes into `review.snapshot_created`
{sha, mode: "snapshot-commit"|"head"} and into the review prompt. Live-tree
fallback remains only for snapshot-creation FAILURE.

### @agy

ACCEPT: the 12-class taxonomy + exact hint strings (seed
internal/runner/failclass.go verbatim), the narrator line formats, the status
header and `parley status` failure block, the P1 kill/retry narration, the P3
"finished (exit code: N, artifact verified)" phrasing (no red error styling),
and the consult UX including stderr-only progress (stdout stays clean for
redirection) and `parley consults list`.

Two unifications: (a) artifact filename uses the UTC run-ID style
`20260612T170815Z-<agent>-<slug>.md` (matches store.NewRunID convention; your
local-time variant loses sortability across timezones); (b) frontmatter is the
union of yours and codex's: {artifact: consult, agent, model, question,
created, workspace_root, timeout_ms, exit_code, session_id?, quorum: false}.
Heartbeat rendering: status-line only (your "inline/overwritten") — never woven
into transcripts.

### @hermes

ACCEPT: record-the-watchdog-event-BEFORE-signalling ordering (the typed
`agent.no_first_output`/`agent.stalled` event is appended, then the kill tree
fires — the durable-kill attribution race disappears because the event log
already names the killer); $TMPDIR-local snapshot placement; copy+fsync+rename
within the target dir for artifact move-back (no cross-device rename), with
fsutil for directory creation; the origin-gc caveat documented ("no gc during a
live review"; stale-snapshot sweep heals crashes); `fsutil.AppendLine` (new
small helper: O_APPEND single write, mirrors store.Append) for
consults/index.jsonl.

Your log-truncation reset concern is fully resolved by codex's counting
writers — growth is measured in-process, not from the filesystem.

## Resolved spec deltas (over round-01 claude.md)

- Supervision: counting writers + ACP activity tracker; heartbeat 60s default,
  excluded from narrator + triggers; watchdog events appended BEFORE kill;
  retry-once only for no_first_output with attempt_id + move-aside.
- P3: codex's table verbatim + ACP phase-validation fix + Result.Success().
- P5: kindly temp-index snapshot commit for dirty trees; sha recorded in event
  + prompt; live fallback only on creation failure; move-back copy+fsync+rename.
- P8: UTC artifact names; union frontmatter; consults list; fsutil.AppendLine.

Ready for consensus. No open blockers from my side.
