# Changelog

## v1.30.4 - 2026-06-20

- **Central per-user agent defaults (`~/.parley/agents.toml`).** A new
  user-global config lists each agent's command, model, and reasoning/effort
  level, inherited by every project. Wired into `config.LoadAgentSpecs` as the
  lowest config-override layer: built-in defaults → `~/.parley/agents.toml`
  (central) → `parley-deck/agents.toml` / `agents.local.toml` (per-project
  override) → `$PARLEY_HEADLESS_AGENT_CONFIG`. A deck overrides the central
  default; fields the deck leaves unset fall through to the central value.
  `parley init` seeds a starter `~/.parley/agents.toml` (never clobbering an
  existing one) and prints where to override per-project. `PARLEY_HOME`
  overrides the central dir (used for hermetic tests).
- **Reasoning/effort is now part of the deck-bootstrap confirmation.** §0 and
  the skill's deck-bootstrap step confirm the roster, each agent's model **and
  each agent's reasoning/effort level**; the default reasoning/effort is the
  **strongest (highest) level the agent supports**, falling back to
  `cli-default` only when it cannot be discovered. Protocol stays model- and
  reasoning-agnostic.

## v1.30.3 - 2026-06-19

- **Fix: roster & model confirmation is a deck-BOOTSTRAP gate, not per-idea/per-session.**
  Corrects 1.30.2, which placed the mandatory roster + per-agent model confirmation in
  the per-idea §9.0 readiness check. It now lives in **§0 (deck bootstrap)**: the
  confirmation fires **once, when `parley-deck/` is first created (`parley init`)** —
  not per idea, not per later session; an already-bootstrapped deck reuses the saved
  selection. §9.0 keeps only the per-idea agent **liveness** ping (it no longer
  re-selects models). Both COOPERATION.md copies; drift-guard lockstep; protocol stays
  model-agnostic. (Skill side: `parley-deck-skill` 1.3.3 moves the interactive flow to
  "Transport Selection / deck bootstrap".)

## v1.30.2 - 2026-06-19

- **Mandatory session-start roster & model confirmation.** §9.0 now states that at a
  session's first readiness check the facilitator MUST confirm the active roster and
  each agent's selected model with the user before the first idea; the user's
  persistent per-agent model choice is recorded in `meta/headless-agents.local.json`
  and reused until changed (later sessions show the saved picks for explicit
  confirmation). The protocol stays model-agnostic — it mandates the confirmation, not
  any specific model. (Both COOPERATION.md copies; drift-guard lockstep. The detailed
  interactive list-roster → confirm → list-models → pick flow lives in the
  `parley-deck-skill` SKILL.md Startup Flow / Selection Checkpoint.)

## v1.30.1 - 2026-06-19

- **Pin the `claude` participant to Opus 4.8 (1M context).** The built-in `claude`
  agent spec launched with `--model opus` — an **alias** the `claude` CLI resolves to
  "the latest opus", which on some installs/accounts landed on an older Opus (e.g.
  4.6). The spec now pins the exact model ID **`claude-opus-4-8[1m]`** (verified the
  CLI accepts it) so `parley run` always launches Opus 4.8 with the 1M-context window,
  not whatever the alias happens to resolve to. (Tradeoff: a future Opus bump must be
  re-pinned; the alias would auto-track but mis-resolved here.) The local
  `headless-agents.local.json` roster was pinned to match.

## v1.30.0 - 2026-06-19

Pre-idea readiness check (idea `meta-protocol-change-preflight-readiness`; 4-agent
deliberation, signoffs claude-1/codex-1/hermes-1, agy waived; Phase-6 review caught a
CRITICAL §1-bypass + 5 MAJORs, all fixed in fix-up cycle 1).

- **Protocol §9.0 "Pre-idea readiness check"** (both COOPERATION.md copies, drift-guard
  lockstep): at idea start the facilitator (a) checks protocol freshness —
  `source`=advisory/no-write, `consumer` additive bump=auto-sync (zone-preserving),
  breaking/unknown-role=confirm — and (b) hosted-PONG-pings the roster, gating per-idea
  exclude / re-include behind explicit user confirmation. Plus a §5 quorum-locks-at-
  Phase-0 sentence and a §7 carve-out (an upstream version sync is not a protocol change).
- **New `parley preflight` command** `[--dir][--json][--yes][--ping-timeout][--no-ping]`:
  freshness classifier + zone-preserving merge + bounded concurrent hosted-PONG probe
  (process-group-killed on timeout) + report/JSON + exit codes 0/1/2/3. Shared with the
  `parley run` pre-check, which runs **before idea creation**, defaults to hosted PONG
  (`--no-ping`/`--no-preflight` opt out), never auto-answers the new gates, and
  hard-stops unattended without reading stdin. The §1 non-solo hard-stop is evaluated on
  the exact `--participants` set; confirmed exclusions are recorded in `00-prompt.md`.
- `meta/version.json` gains `protocolRole` (`source`/`consumer`, fail-closed); `parley
  init` now writes `protocolRole: consumer`.
- Also bundles a 4-participant roster update (§2 tables → `claude-1`/`codex-1`/`hermes-1`/
  `antigravity-1`; backend map).
- Known follow-ups (deferred): roster-ID↔runtime-ID `-1` reconciliation in reports;
  preflight freshness-probe perf for source/`--no-ping`.

## v1.29.0 - 2026-06-19

Protocol: Fusion + ExecPlans inspiration (idea
`meta-protocol-change-fusion-execplans`; 4-agent deliberation, signoffs
claude/codex/hermes, agy waived on a tooling hang). Additive, **conditional-rigor**
guidance applied byte-identically to both `COOPERATION.md` copies (drift-guard
lockstep); no Go logic changed; embedded `parley init` default stays genericized.

- **`FINAL.md` gains static, self-contained design-time sections** (Phase 4): Purpose
  / user-visible outcome, Context & orientation, **Observable acceptance criteria**,
  **Idempotence & recovery**, Known risks / de-risking. `FINAL.md` stays immutable.
- **`IMPLEMENTATION.md` becomes a living execution doc** (Phase 5): Progress
  (timestamped), Decision Log, Surprises & Discoveries, Validation evidence, Outcomes
  & Retrospective — so a fresh headless agent or the auto-drive driver can resume
  **from the artifact alone**, and §13 `parley retro` gets richer evidence.
- **Advisory "Comparison & blind spots" lens** in `consensus.md` and
  `review/consensus.md` (Phase 3/7) — surfaces what *no* participant addressed. Not a
  gate; append-only signoffs remain the only gate.
- **Phase 6** reviewers may check observable acceptance criteria; severities
  (CRITICAL/MAJOR/MINOR/NIT) unchanged.
- **§13** gains a **confident-error** retro evidence signal (diagnostic only — never a
  new severity, blame label, or merge gate).
- Full living/static sections are required only for complex / `auto_implement` /
  driver-managed / pipeline ideas; trivial or design-only ideas may use `N/A`.
- Explicitly **rejected** (inspiration we did *not* adopt): confidence-by-breadth
  gates, a single-model judge with authority, hiding raw rounds behind a summary, the
  Fusion panel/recursion/cost/web-search machinery, collapsing the deck into one file,
  proceed-without-prompting autonomy across gates, and the anti-list prose maximalism.

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
