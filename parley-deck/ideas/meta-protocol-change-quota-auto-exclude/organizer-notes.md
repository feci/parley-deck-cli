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

## Phase 1 launch (2026-10-03, 12:32 CEST)

- Phase 0 was committed as `04b22e2`. The seeded run is `runs/20261003T103229.878442000Z/`, with `run.json`
  (no roster snapshot, see gap 3), an empty `events.jsonl` and `questions/`. It is untracked runtime state,
  as in the designated-implementer run. `parley continue <slug>` now resolves the idea to this run.
- Participant protocol context for round 1: `parley protocol packet --phase 1 --track deliberation
  --flag protocol_change` gave context_mode=full, with source_sha256 = packet_sha256 =
  `b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388` and no fallback. The body hash was
  checked before it was embedded in the runner's envelope. The prompts are private run state at
  `runs/<id>/agents/<agent>/prompt.txt`: codex-1 is 131320 bytes (sha256 prefix `9ad7f636044b65a9`) and
  claude-1 is 131347 bytes (`105be7085660974d`). They differ only in agent ID, output path and
  launch-config lines.
- claude-1 participant: launched 12:32:30 CEST as a separate headless process with a 1800 s alarm and a
  sanitized env. The organizer passed it nothing beyond the shared prompt.
- codex-1: the first attempt at 12:32 never started the agent. It was an organizer launcher bug: under
  `set -u`, macOS bash 3.2 treats an empty array as unbound, and no log was created. This was **not a
  provider error**. Fixed and relaunched at 12:33. The codex banner confirms model gpt-6-astra, provider
  omniroute, approval never, sandbox workspace-write and reasoning effort max.
- The launch order differs by about one minute. Round-1 independence is the protocol's social rule (each
  prompt says not to read the other's round-01 file). No isolation beyond that is claimed.
- **Observer gap 5.** `parley wait --for round` exited 0 ("boundary reached (round complete)") at about
  12:34. The cause is the runner template's own rule, "Immediately use your file-writing tool to create
  that exact file with its required frontmatter". Both agents first wrote placeholder stubs (codex-1 421 B,
  claude-1 370 B, bodies like "Independent analysis in progress."), and the validator accepts those as
  filed and valid. Both processes were still running. So the organizer treats **process exit plus final
  content** as the round boundary, and does not start `parley continue --auto` until both round-1
  processes have exited. Otherwise the driver would open round 2 over the stubs.

## Phase 1 complete (2026-10-03, 12:57 CEST)

- codex-1 exited 0 at 12:49 after 960 s, with a 28178 B artifact. Its stderr reports "tokens used 162,290".
  The claude-1 participant exited 0 at 12:56 after 1449 s, with a 24569 B artifact. Neither process hit a
  quota or credit error. Neither filed a question under `runs/<id>/questions/`.
- `parley wait --for round` then showed 2/2 valid, with `unparsed=false` for both. The mechanical stance
  flags were codex-1 block 18 / accept 7 / escalate 2 and claude-1 block 17 / counter 1 / accept 4 /
  escalate 9. The raw files were opened: the counts come from the topic's own vocabulary (blocking gates,
  escalation paths). No round-1 file is a BLOCK or a blocking owner escalation. Both record the attestation
  (full, `b273af1e…f388`, no fallback).
- `git status` showed no participant edit outside its own round file.
- Procedural reading only (no verdict). The round-1 positions overlap on per-idea scope only, a
  positively classified signal, a floor of 2, dissent surviving and fail-closed. They differ on the
  threshold (codex-1 15 min, claude-1 60 min), on what proves exhaustion (codex-1 requires allowance
  semantics, claude-1 accepts a reset window of at least 60 min or an allowlisted phrase), on the mid-idea
  trigger point, and on how a participating facilitator counts toward the floor. claude-1 lists five
  candidate owner decisions. These go to cross-review, and nothing is escalated from round 1.
- Committed as `f3345c6` (codex-1) and `f878287` (claude-1), each under its author's prefix and marked as
  swept by the organizer.
