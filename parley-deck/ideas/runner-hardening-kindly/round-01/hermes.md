---
agent: hermes
idea: runner-hardening-kindly
round: 1
date: 2026-06-12
---

## P1 — Watchdog timing edge cases

- Clock granularity: 1s sleep in wait_with_progress is acceptable; macOS time(2) + date(1) give second resolution. No sub-second requirement.
- Log rotation/truncation: kindly tracks bytes via wc -c on the event file. A truncate-to-zero would reset growth and trigger stall kill after KINDLY_STALL_SECS. Runner already writes per-agent stdout/stderr logs; if the TUI growth cache (1.23.0) is reused, it must also tolerate reset.
- ACP vs exec first-output: ACP streams JSON events on stdout immediately; exec writes to -o file. First-event signal must be "any stdout byte or event file byte growth". ACP can emit the initial {"type":"thread.started"} within <5s; exec waits for the first real turn.
- Interaction with TimeoutMS and procctl kill gates: watchdog kill (TERM then KILL via signal_tree/kill_tree) must precede the durable-kill attribution check. A race where the watchdog fires and procctl's kill attribution writes the wrong killer ID is possible if both act on the same process group within the same scheduler tick. Recommendation: watchdog records its own kill reason in the event log before signalling; the durable gate should treat an explicit "watchdog" marker as authoritative.
- Heartbeat events: reuse existing agent.heartbeat or emit agent.progress with {elapsed, bytes, last_event}. Config surface: KINDLY_* env vars for test harness only; production uses fixed defaults (120s grace, 30m stall).

## P5 — Shared-clone snapshot on LOCAL tmp

- Origin .git on virtio-fs: alternates file points at the origin objects; read path is correct as long as the origin remains mounted. If the mount hiccups mid-review the clone becomes a dangling reference; recovery is a fresh clone on next run (pid marker healing already present in kindly-agent.sh).
- Pid-marker healing races: stable path per (repo, reviewer) + stepped-aside unique path on collision. Healing on next launch is racy only if two instances race the mkdir; kindly already serialises via the stable path.
- Artifact move-back atomicity: snapshot lives under $TMPDIR (local SSD). Final artifact is written inside the snapshot then copied to the canonical location. rename(2) across devices fails; must be copy+fsync+rename within the target directory. fsutil retry wrapper is not currently used for this path and should be added for the final move.
- Origin GC risk: a gc that removes objects referenced only via the alternates link will corrupt the snapshot. Mitigation already present: clone is --shared so the origin must stay alive for the clone lifetime; no GC during a live review.

## P4 — GIT_OPTIONAL_LOCKS=0 and env -u

- Scope: already exported at the top of kindly-agent.sh for every git probe inside the snapshot (status, diff, rev-parse, commit-tree). No behavior change for read-only probes; the flag only suppresses optional lock writes on the weakly-coherent mount.
- Marker shedding completeness: env -u covers CLAUDECODE, CLAUDE_CODE_SESSION_ID, CLAUDE_CODE_ENTRYPOINT, AI_AGENT and the Codex equivalents. The list matches the detection block; no omissions.

## P8 — Consult artifact durability

- index.jsonl append semantics: kindly writes one provenance line per run under the report directory. Concurrent consults on different slugs are independent. Same-slug concurrent consults step aside via the stable path mechanism already used for audits.
- No fsutil wrapping today; the append is a plain >> so partial writes are possible on virtio-fs. Add a small fsutil.AppendLine helper if we adopt the consult command.

## Position

Adopt P1, P4, P5, P8 with the timing, race, and atomicity notes above. P5 snapshot location must be $TMPDIR (local) not the shared mount. Add fsutil retry for the final artifact move. No cuts required; all six points remain in scope.