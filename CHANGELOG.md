# Changelog

## v1.28.1 - 2026-06-16

- **`parley retro` precision fix.** The deterministic scanner matched signal
  patterns in free text, so it false-positived on prose that merely *discussed*
  them — e.g. it flagged `rho-retro-tooling` as "blocked-or-abandoned" because its
  own review consensus quoted ``Verdict: BLOCK``. Blocker detection is now
  anchored to structure: a real `Status: ❌` signoff line, or a `## Verdict`
  heading whose leading token is `BLOCK`/`BLOCKER` (or contains ❌) — not the
  substring "block" in prose, and not a `REQUEST-CHANGES`/"no blocking issues"
  explanation. NOT-FIXED is counted only in review round files, dismissed-findings
  only in consensus files. Regression test included
  (`TestBlockerDetectionIgnoresProse`). Surfaced by dogfooding `parley retro` on
  this repo right after the 1.28.0 ship.

## v1.28.0 - 2026-06-16

Retrospective optimization (RHO adoption — two reviewed ideas,
`meta-protocol-change-rho-retrospective-optimization` + `rho-retro-tooling`):

- **Protocol §13 "Retrospective optimization"** added to both COOPERATION.md
  copies (drift-guard lockstep). A retrospective pass mines the deck's own history
  to *propose* improvements but **applies nothing**: proposals enter as a normal
  idea (protocol-text changes via a meta-protocol-change idea + human approval),
  acceptance is the normal multi-agent gate (consensus + all-participant signoff +
  no-regression), and RHO-style self-preference is a diagnostic note only. Defines
  the layered harness (protocol / runtime "Repository Instruction Files" / local
  "Agent Local Memory" / evidence corpus) and the guardrails (audit,
  adversarial-trajectory hygiene, reversibility, multi-agent diagnosis).
- **New `parley retro` command** — read-only mining of the deck's structured
  artifacts: `scan` (failure-density signals per idea), `select` (type-diverse
  "hard cases" coreset), `diagnose` (grouped report), and `propose --slug` (which
  scaffolds **only** a single new `ideas/<slug>/00-prompt.md`, fail-closed). No
  raw session transcripts and no DPP/embeddings/re-rollout in v1; deterministic.
  Inspired by RHO (arXiv:2606.05922) but replacing its single-model self-preference
  with the deck's multi-agent quorum.

## v1.27.0 - 2026-06-15

- **Auto-drive now works on every transport.** The driver's auto-advance was
  hard-gated to `local-dir`, so a `github-pr` / `gitlab-mr` run stalled at
  round-01 even with auto-drive on. The gate is now transport-independent: the
  canonical artifacts (rounds, consensus, FINAL, …) are the source of truth under
  every transport, so auto-drive advances them everywhere. Only `--auto` /
  `--no-auto` gates it now. The driver still does NOT create PR/MR branches — that
  mirroring stays a manual, ergonomic step.

## v1.26.0 - 2026-06-13

- **New TUI `/run` command.** Advance the protocol on demand from inside the live
  TUI — it kicks the auto-driver (cross-review → consensus → finalize → opted-in
  implementation) for the current run. Most useful with `--no-auto` runs; under
  the default auto-drive it is a no-op once driving has started (idempotent). The
  command appears in `/help` and slash autocomplete.

## v1.25.0 - 2026-06-13

Auto-drive is now the default.

- **`parley run` auto-drives by default.** After round-01 the protocol now
  advances automatically — cross-review rounds, consensus draft, signoff
  requests, and finalize — without you running the next step. Pass **`--no-auto`**
  to opt out (stop after round-01 and advance manually). The flipped flag also
  governs the launch prompt: a default run launches and drives unattended, while
  `--no-auto` (without `--yes`) restores the pre-launch confirmation.
- **Auto-drive now runs inside the TUI.** Previously the driver only ran on the
  `--no-tui` path, so a TUI run stalled at round-01. The driver now runs in the
  background while the live TUI shows it advancing (its output is discarded so it
  never corrupts the render; quitting the TUI stops it).
- **Code-mutation stays gated.** The implementation/fix-up phases (Phase 5–8) are
  still only auto-driven when the idea opts in via `auto_implement`; flipping the
  auto default does not auto-write code. `--no-implement` still stops the driver
  at `FINAL.md`.
- `parley continue` is unchanged: it still prints the next action by default and
  executes it only with `--auto`.

## v1.24.1 - 2026-06-13

Maintenance (idea `embedded-default-protocol-resync`, PR #47):

- **Embedded default protocol resynced** with the live deck. The `parley init`
  bootstrap template (`internal/protocol/defaults/COOPERATION.md`) gained the
  missing `## 12. Pipeline blocks & action stages` section (byte-identical to the
  live deck) and was **genericized**: header `Workspace`/`Created` are now
  placeholders and both §2 tables ship empty bodies, so a freshly `parley init`-ed
  project no longer inherits this repo's roster/workspace.
- **Anti-drift guard**: a fail-closed Go test (`TestEmbeddedDefaultMatchesLiveDeck`)
  asserts the embedded default stays in sync with `parley-deck/COOPERATION.md`
  (modulo five documented, anchored project-specific zones) and that the embedded
  bootstrap shape holds — so a protocol edit landing in only one copy now breaks
  the build. Plus `TestDefaultCooperationForInit` for the init output.
- Synced the project deck to `parley-deck-skill` 1.3.1 (§12 was already present).

## v1.24.0 - 2026-06-12

Adopted from the MIT-licensed "kindly" skill (ideas `runner-hardening-kindly` +
`meta-protocol-change-review-gate-honesty`):

- **Agent supervision**: first-output watchdog (120s, one retry), stall guard
  (30m, output-growth based), persisted `agent.heartbeat` events (60s; excluded
  from transcripts/triggers); counting writers — zero healthy-path I/O; typed
  `agent.no_first_output`/`agent.stalled` events appended BEFORE the kill.
  Config: `first_event_timeout_ms`, `stall_timeout_ms`, `heartbeat_ms`.
- **Failure classification**: `agent.failed` now carries `failure_class` +
  `recovery_hint` (rate-limit/auth/billing/overloaded/…); surfaced in the TUI
  narrator and agent headers.
- **Artifact beats exit code**: a validated artifact with an ordinary nonzero
  exit finishes with `agent_exit` instead of failing (removes the agy
  wrote-then-exit-1 flake); ACP validation now respects the run phase; fix-ups
  validate IMPLEMENTATION.md instead of trusting exit 0; `Result.Success()`.
- **Review snapshots**: Phase 6 reviewers read a disposable shared-clone
  checkout on local tmp (dirty trees become temp-index snapshot commits);
  artifacts move back via copy+fsync+rename; loud fallback events.
- **parley consult** + `parley consults list`: advisory cross-agent questions
  with durable artifacts under parley-deck/consults/ (never quorum evidence).
- Hardening: claude participants shed nested host markers; read-only git probes
  set GIT_OPTIONAL_LOCKS=0; `fsutil.AppendLine`; docs/agent-cli-mechanics.md.
- **Protocol**: Phase 6 "Review briefs and dispositions" (no-suppression),
  Phase 8 opt-in `strict_gate` + "Stopping judgment", §8 "Consults" standing;
  mirrored to the embedded default protocol.

## v1.23.0 - 2026-06-12

- Protocol visibility in the live TUI (idea `tui-protocol-visibility`):
  collapsible protocol ribbon on every tab (Ctrl+P), tab activity glyphs
  (spinner/silent/delivered/failed/STALE), woven narrator lines, a Protocol
  tab with pipeline/delivery/signoff/next panes, and a Home phase column.
- New `run.phase` event emitted by the driver after every phase-changing
  cursor commit; cursor save errors are no longer discarded.
- `driver.RebuildDetail` exports the phase evidence (review round, review
  consensus, implementation status) in one disk pass.
- Declared `buffers_stdout` agent flag (TOML + run.created runtime payload);
  silent buffered agents get a structured placeholder instead of a blank tab.
- Status line shows `ph=N:<phase> wait=<agents>` instead of `round=<status>`.

## v1.5.4 - 2026-05-27

- Treat ACP as a selectable launch mode on an existing agent instead of
  exposing duplicate `*-acp` agent IDs for Codex, Claude, and Hermes.
- Add the TUI `a` key for session-only ACP launch overrides and show ACP
  command details in the selected-agent panel.
- Add `acp_args` runtime configuration so local installs can enable ACP for
  CLIs when their concrete ACP launch args are known.
- Apply TUI launch-mode overrides to newly started runs and record effective
  launch metadata in run runtime events.

## v1.5.3 - 2026-05-27

- Add Antigravity CLI `agy` as a first-class headless agent and default
  replacement for Gemini.
- Mark Gemini as legacy compatibility while keeping existing overrides working.
- Prefer verified stronger defaults for Claude (`opus`/`max`) and Hermes
  (`xai/grok-4.3`) while keeping Antigravity model/thinking fields at
  `cli-default` until the CLI exposes flags.
- Update project and embedded protocol rosters, docs, examples, and runtime
  configuration defaults for the Antigravity migration.

## v1.5.2 - 2026-05-26

- Add TUI planner action execution and focus-aware action controls.
- Refresh dashboard and live TUI layouts with height-aware compact modes,
  two-column normal views, semantic badges, and short-terminal tests.
- Embed the full default `COOPERATION.md` protocol for workspace initialization
  while preserving local-dir bootstrap transport.
- Record concrete project roster metadata and ignore machine-local headless
  agent launch config.
- Improve `parley version --all` project status probing timeouts and fallback
  behavior.

## v1.5.1 - 2026-05-25

- Complete Parley review cycles for the continuous-run TUI planner slice and
  version-awareness project sync.
- Unify planner and manifest next-action serialization through a shared
  `internal/runaction` type.
- Make continuation planning round-aware and remove hardcoded `codex` ownership
  from generated continuation commands.
- Add `parley version --dir DIR --all` project targeting, indented JSON output,
  and cleaner missing-installer fallback errors.
