# Organizer execution notes

Organizer: claude-1 (declared facilitator), launched headlessly by the owner's relay session with
`claude -p --model claude/claude-opus-5-5[1m] --effort max`. This process writes orchestration files only.
Every claude-1 participant artifact comes from a separately launched claude-1 process (§15.5, recorded in
`00-prompt.md`).

## Phase 0 (2026-10-03, 12:17 to 12:40 CEST)

- Read the controlling brief and the evidence file in full. Both are copied byte-identical to
  `source-context/` (sha256 `e33b222c…af8b6` and `9d46b27d…f7bf`). A secret-pattern scan of both was clean.
- Applied the parley-deck skill. Protocol context: `parley protocol packet --phase 0 --track deliberation
  --audience facilitator --flag protocol_change` gave context_mode=packet, source_sha256 `b273af1e…f388`,
  packet_sha256 `0744ea22…b5ca`, with no fallback reason. §3 and §15.5 were read from the full source on
  demand. `parley organizer brief` reports `full (full fallback: facilitator-participates)` for this idea,
  because the facilitator participates. The organizer reads omitted sections by their triggers. It
  adjudicates no verdict.
- OpenViking shared memory is not connected in this session (no MCPAnywhere tools loaded), so local sources
  are used and nothing is claimed as saved.
- Preflight: `parley preflight --dir . --json` ran at 12:23 CEST in 94 s and exited 3. claude-1 and codex-1
  are ready. kimi-1 is `deadline-after-output`. zcode-1 is `provider-failure:rate-limit`. The facilitator-
  declaration gate on `meta-protocol-change-devx-speed` is unrelated, as the designated-implementer run also
  recorded. That idea is not touched. The raw JSON is in `source-context/`.
- Exclusions follow the owner's relayed decision ("codex-1 + claude-1 (Recommended)"). They are written in
  the CLI's own `excluded: <id> — <reason> — confirmed <date>` line format. Preflight was not run with
  `--yes`.
- Version state: CLI 1.50.0. The skill installer and the codex and claude runtime markers are 2.14.0.
  Freshness is source-advisory with no write.
- Session-start state: no inbox note is addressed to claude-1. The windows-portability run is in flight in
  its own worktree and is not touched.

### Driver gaps (Phase 0)

1. **No existing-kickoff entry point.** `parley run TASK` always mints a timestamped slug
   (`internal/protocol/workspace.go` `CreateIdeaFull` calls `timestampedSlug`), and §7 requires
   `ideas/meta-protocol-change-<topic>/`. `parley continue <slug>` on a pre-created idea fails with
   `idea "<slug>" has no runs yet` (`internal/runstate/runstate.go` `ResolveRun`). `parley resume` only opens
   the TUI, and `parley steer` only queues messages. The designated-implementer run hit the same gap.
2. **The driver never launches the current round, only the next one.** `advanceRound` returns `ActionAwait`
   while round N is incomplete and calls `RunRound(N+1)` only after it completes
   (`internal/driver/driver.go`). Round 1 is launched only inside `parley run` (`runner.RunRoundOne`).
   - Fallback used: seed orchestration-only run metadata (`runs/<id>/run.json` plus an empty `events.jsonl`),
     then launch round 1 manually with each participant's configured headless invocation. The prompt is the
     runner's own round-1 template (`BuildRoundOnePrompt`), wrapped in the runner's own protocol-context
     envelope, and it is identical for both participants apart from agent ID, output path and launch-config
     lines. The claude child environment is sanitized as `cleanParticipantEnv` does. No events or
     participant content are forged. From round 2 on, `parley continue --auto --no-implement` drives.
3. **No safe roster freeze for a seeded run.** The run manifest's `roster_snapshot` replaces model, effort,
   speed and the whole headless argv at every continuation (`applyRosterSnapshot`). No CLI verb prints the
   resolved argv, so a hand-written snapshot could silently change launches. The seeded manifest therefore
   carries no snapshot, and continuations use the live configuration. That configuration at launch, from
   `parley agents list`:
   - codex-1: `codex exec --skip-git-repo-check --cd {root} --sandbox workspace-write -c approval_policy="never" -c model_reasoning_effort=max -m gpt-6-astra -`, with a 1800000 ms timeout.
   - claude-1: `claude -p --model claude/claude-opus-5-5[1m] --effort max --output-format text --permission-mode bypassPermissions --add-dir {root}`, with a 1800000 ms timeout.
   - The organizer makes no roster or config change during this run.
4. **The driver defaults to one cross-review round** (`ReadCrossReviewRounds`, default 1). It then drafts
   consensus. A BLOCK reopens another round, up to the §4.0 deliberation cap of 3 after round 1, and then
   escalates (`internal/driver/consensus.go`). The default is kept. The signoffs are the convergence gate,
   and the owner pre-authorized the Phase 0 to Phase 4 transitions in the brief.
