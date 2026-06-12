---
idea: runner-hardening-kindly
drafted-by: claude
date: 2026-06-12
---

## Agreed decisions

Source: round-01 + round-02 (all four participants ACCEPT; codex round-01 is the
architectural backbone; deltas converged in round-02).

- **D1 — Supervision architecture (P1).** New `internal/runner/supervision.go`:
  `SupervisionConfig{FirstEventTimeout, StallTimeout, HeartbeatInterval}` + an
  in-process `activityTracker`. Exec mode: counting writers wrap cmd.Stdout/
  Stderr (zero healthy-path FS I/O; no stat probing — truncation concerns moot).
  ACP mode: `acpRunnerHandler` marks activity on session updates and protocol
  events; `agent.started` itself never satisfies first-output. `waitSupervised`
  replaces the select at runner.go:732-739. Kill ordering (hermes): append the
  typed watchdog event FIRST, then `procctl.KillGroup`, then drain Wait — the
  event log names the killer before any signal race with durable-kill
  attribution.
- **D2 — Knobs.** `first_event_timeout_ms` (default 120000),
  `stall_timeout_ms` (default min(1800000, agent timeout)), `heartbeat_ms`
  (default 60000; 0 disables heartbeats only; other guards disable only when
  explicitly configured 0). Fields on agents.Spec near TimeoutMS + TOML
  overrides + recorded in the run.created runtime payload. `timeout_ms` stays
  the hard outer budget.
- **D3 — Retry.** Retry ONCE only for `no_first_output` (never for stalled,
  timeout, validation, auth/billing, user kill). `attempt_id` threads through
  agent.started, watchdog events, terminal events, and the procctl marker
  (runID:agentID:attemptID). An invalid attempt-1 artifact is moved aside to
  `<artifact>.attempt-1.invalid`; never delete a valid artifact, never
  overwrite a pre-existing one. `agent.failed` for the killed attempt is
  appended BEFORE the retry's `agent.started` (durable-kill targets the latest
  started attempt).
- **D4 — Heartbeats.** `agent.heartbeat` persisted, payload {agent, segment_id,
  attempt_id, phase, launch, elapsed_ms, timeout_ms, stdout_bytes,
  stderr_bytes, last_activity_ms_ago}; never counts as activity. View
  contract: EXCLUDED from narrator transcripts, protocol-snapshot trigger sets,
  review/consensus prompt context; only the latest heartbeat shows in live
  status surfaces and `parley status --verbose`; events.jsonl remains the audit
  trail.
- **D5 — Failure classification (P2).** `internal/runner/failclass.go`: ordered,
  data-driven, bounded regex table seeded with agy's 12 classes AND exact hint
  strings (round-01/agy.md table verbatim; hyphenated provider classes;
  watchdog classes named no_first_output/stalled to match event types).
  Applied to bounded stderr/stdout tails + ACP error events + ExitError.
  `agent.failed` gains failure_class, recovery_hint, exit_code, signal,
  stderr_tail_bytes. Watchdog/timeout classes are set before regex
  classification. ONE payload builder shared by exec, ACP, failEarly, and
  fixup terminal paths.
- **D6 — Classification consumers.** runstate.AgentState += FailureClass,
  RecoveryHint, AgentExit (projected in applyAgentEvent; summarized in
  SummarizeEvent). Narrator: `── <ts> <agent> failed: <class> — <hint> ──`;
  agent status header and `parley status` failure block per agy round-01.
  Watchdog events join narrator + snapshot trigger sets; heartbeats do not.
- **D7 — Artifact-wins decision table (P3).** New `finalizeAgentResult`
  helper, codex round-01 table verbatim: valid artifact + clean exit →
  finished; valid + ordinary nonzero exit → finished with `agent_exit` (Result
  gains AgentExit; ExitError cleared; new `Result.Success()`); valid + ACP
  prompt error after session start → finished with agent_exit.kind=acp_error
  (init/session-setup errors still fail); timeout / final no_first_output /
  stalled / user kill → failed regardless of artifact; invalid/missing →
  failed. Bug fix folded in: ACP validation switches from hard-coded
  ValidateRoundOneArtifact to validateArtifactForPhase (acp.go:158-163).
  Fix-up: new `ValidateFixupArtifact` (IMPLEMENTATION.md validates, contains
  the new fix-up section, review-ready status) gates artifact-wins in
  phase58. driverImplOps (driver_impl.go:99-105,178-187,218-224) and
  pipeline_cmd (780-784,817-820,845-848) switch to Result.Success().
- **D8 — Small hardening (P4).** `cleanParticipantEnv`: when spawning the
  claude CLI as a participant, drop CLAUDECODE, CLAUDE_CODE_SESSION_ID,
  CLAUDE_CODE_ENTRYPOINT, CLAUDE_CODE_ENABLE_TASKS, AI_AGENT and AI_AGENT_*
  (keep PARLEY_*), in both exec and ACP paths. GIT_OPTIONAL_LOCKS=0 on
  gitTreeClean's git invocations (the only Go git probe today; tests assert
  the env). New `docs/agent-cli-mechanics.md` (codex stdin/-o, claude -p
  binding + --tools vs --allowedTools + MCP bypass + cwd-scoped resume, agy
  --print value-taking + buffers_stdout, hermes -z + silent-death mode),
  linked from docs/agent-runtime-configuration.md. Codex `-o` documented,
  NOT adopted (artifact-file contract differs).
- **D9 — Review snapshots (P5).** New snapshot module (internal/runner/
  reviewsnapshot.go). Phase 6 review rounds only (Phase 5/8 keep the live
  tree). Location: LOCAL tmp `<tmp>/parley-review-snapshots/<repoHash>/<idea>/
  <round>/<agent>` (fsutil for parents). Mechanics: `git clone --shared
  --no-checkout` + clean tree → detached checkout of HEAD; dirty tree →
  kindly temp-index snapshot commit (GIT_INDEX_FILE read-tree HEAD → add -A →
  write-tree → commit-tree -p HEAD → detached checkout; codex verified the
  sequence with a local-tmp clone and a virtio-fs origin); staged/worktree
  divergence → live-tree review (kindly's pre-check); any creation failure →
  `review.snapshot_fallback` + live tree. Resolved sha + mode recorded in
  `review.snapshot_created` and the review prompt. The reviewer writes its
  artifact inside the snapshot at the usual relative path; the runner
  validates, then moves it back to the canonical deck path via copy + fsync +
  rename WITHIN the target directory (no cross-device rename); move-back
  failure → `review.snapshot_artifact_move_failed` with the recovery path and
  the snapshot retained. Lifecycle: `.pid` marker {run id, agent id, pid,
  boot id}; stale-marker sweep on create; live marker → step-aside suffixed
  path; teardown on round end. Terminal events report the LIVE artifact path.
  Documented caveats: origin objects must stay alive for the clone lifetime
  (no gc during a live review); a mid-review mount hiccup dangles the clone —
  recovery is a fresh snapshot (hermes).
- **D10 — parley consult (P8).** `parley consult [--dir DIR] [--timeout D]
  <agent> "<question>"` (stdin when the question arg is absent). The
  facilitator captures stdout into the artifact — the agent stays
  read-only-ish; progress goes to stderr so stdout stays redirectable.
  Artifact: `parley-deck/consults/<YYYYMMDDTHHMMSSZ>-<agent>-<slug>.md` with
  codex's canonical frontmatter {artifact: consult, agent, model, created,
  question_slug, question, workspace_root, timeout_ms, exit_code, session_id,
  stdout_log, stderr_log, quorum: false}. Provenance line appended to
  `parley-deck/consults/index.jsonl` via new `fsutil.AppendLine` (O_APPEND
  single write mirroring store.Append; cross-process concurrency guarded by
  an adjacent claim/lock because virtio-fs O_APPEND atomicity is weak —
  hermes). `parley consults list` table with a FILE column showing the bare
  filename. Consult runs under D1 supervision + D5 classification and never
  touches protocol/quorum state.
- **D11 — Slices (one release, 1.24.0, with the sibling protocol idea).**
  (1) shared result decision + P3 artifact-wins + tests; (2) supervision core
  exec+ACP; (3) P2 classifiers + consumer surfacing; (4) P4 env/git/docs;
  (5) P5 snapshots; (6) P8 consult + consults list.
- **D12 — Tests.** codex round-01 "Required tests" list adopted verbatim
  (runner exec + ACP matrices incl. retry/stall/heartbeat/artifact-wins edge
  cases; agents/config knob parsing + runtime payload; runstate projection;
  TUI narrator/heartbeat exclusion; driver/app Result.Success +
  GIT_OPTIONAL_LOCKS assertion; snapshot construction/healing/fallback/
  move-back-failure; consult usage/stdin/schema/no-mutation).

## Agreed trade-offs

- Heartbeat events add ~30-60 rows/agent/30min to events.jsonl — accepted for
  operator liveness; the strict view contract keeps transcripts clean.
- Snapshot clones depend on origin objects via alternates — accepted with the
  no-gc-during-review caveat documented; a hiccup costs only a re-snapshot.
- Artifact-wins can mask an agent's nonzero exit — bounded: only when the
  artifact VALIDATES, the exit is preserved in agent_exit, and watchdog/user
  kills always win.

## Open items deferred to implementation

- Exact counting-writer flush/locking details (single-goroutine writes per
  stream; atomic counters).
- consults claim/lock shape (adjacent .lock dir vs O_EXCL file) — implementer
  picks, test covers two concurrent consults.
- Optional future: artifact-size watcher for buffered CLIs; codex {final_output}
  placeholder; driver ReadStrictGate (sibling idea defines flag semantics).

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-12
Status: ✅ ACCEPT
Notes: Drafted from fully converged round-02 positions; D9 folds codex's local verification of the kindly temp-index sequence.

### Signoff: codex — 2026-06-12
Status: ✅ ACCEPT
Notes: Verified against codex round-01/round-02; consensus captures the converged runner-hardening decisions.

### Signoff: hermes — 2026-06-12
Status: ✅ ACCEPT
Notes: Faithful record of round-02 convergence.

### Signoff: agy — 2026-06-12
Status: ✅ ACCEPT
Notes: Verified against agy round-01/round-02; the consensus faithfully captures the converged supervision UX and consult list changes.
