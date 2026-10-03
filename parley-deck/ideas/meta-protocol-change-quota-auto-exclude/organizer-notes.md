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

## Phase 2: driver halt and round-2 fallback (2026-10-03, 13:00 CEST)

- `nohup parley continue --dir <worktree> --auto --no-implement meta-protocol-change-quota-auto-exclude`
  (pid 59040) reconstructed round 1's `round.completed` event and then halted in about 1 s. Verbatim:
  `continue --auto: run round-02: cross-review accounting: historical worktree is unavailable:
  /private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/f2repo`.
  Before that it printed a notice: `designated run leaves a single non-implementer reviewer (claude-1);
  codex-1 implements and never reviews itself`. That is a Phase 5/6 consideration recorded for the
  proposal.
- **Driver gap 6 (same as the designated-implementer run).** `internal/budget/binding.go` enumerates
  `git worktree list` and refuses when a registered worktree path is missing, unless it is
  operator-declared unavailable. The shared repository has **13 stale worktree registrations**: the
  `/private/tmp/.../scratchpad/f2repo` and `/private/tmp/revert-test2` checkouts plus 11 `parley-*`
  checkouts under `/private/var/folders/.../T/`. The gate is reached only through the driver's
  cross-review round launch (`internal/driver/cycle_budget.go`, `admitCrossReviewCycle` /
  `cycleRoundRunner.RunRound`). The consensus, signoff and FINAL launches use `runHeadlessSignoffAgent`
  instead.
- No worktree registration is pruned and no unavailable worktree is declared. Either act rewrites shared
  git or budget history that the budget system treats as evidence, and that is the owner's call. It goes
  in the proposal as an owner action item. The driver's blocking note
  `inbox/claude-to-user_meta-protocol-change-quota-auto-exclude_driver-error.md` (author `claude`, the
  driver default) is left in place unedited, because its cause is real and unresolved. The driver does
  not gate on inbox notes, and `parley wait` reports a pre-existing note only as a digest note.
- Fallback, as in round 1: the runner's own cross-review template (`BuildRoundPrompt` plus
  `gatherPriorRounds`, both round-01 files inline, 52816 B) in the runner's envelope. The phase-2 packet
  is full, with source = packet = `b273af1e…f388` and no fallback. The organizer set
  `status: round-02` (the driver's promotion step) and launched both participants at 13:02 CEST with
  1800 s alarms. The prompts are 170233 B (codex-1) and 170234 B (claude-1). Logs are in
  `runs/<id>/round-02/<agent>/`.

## Phase 2: interruption and round-2 relaunch (2026-10-03, from 13:54 CEST)

- **Interruption, not a provider error.** The previous organizer (pid 1554) and both round-2 participant
  processes were killed at about 13:02 to 13:08 CEST, when the owner's Claude Code session that had launched
  them exited. The relay checked this at 13:52 CEST (`runs-handoff/.../RESUME-BRIEF-1.md`): codex-1's
  round-2 stderr ends mid-work with no 429, 503 or limit text, and claude-1's round-2 logs are empty. The only
  round-2 file is codex-1's own scaffold (456 B, "Cross-review in progress; this scaffold is not a completed
  round or signoff"). It is not a completed round. claude-1 has no round-2 file. The owner asked for the
  resume ("pokracuj", "continue").
- Resumed at 13:54 CEST as a new organizer process (pid 57867) in the detached tmux session `qae-organizer`,
  so it survives the owner's session. Re-oriented from `parley organizer brief`, `parley status` and these
  notes. No round-2 participant process was alive.
- **Owner decision on driver gap 6** (non-blocking inbox note `user-to-claude-1_…_worktree-prune.md`, 13:56):
  "Yes, prune with a backup (Recommended)". The owner's session backed up the 13 admin directories to
  `runs-handoff/worktree-prune-backup-20261003/` and ran `git worktree prune -v`. Checked here:
  `git worktree list --porcelain | grep -c '^prunable'` returns 0, all 21 registered paths exist and none
  is locked. Gap 6's cause is resolved, so the driver's note `claude-to-user_…_driver-error.md` is archived
  under `inbox/archived/`, as the owner permitted. The organizer pruned and declared nothing.
- **Driver gap 7 (new): the driver cannot relaunch an interrupted round over a participant's scaffold.**
  `parley continue` (print-only) now recommends "Open round-02 (cross-review) before drafting consensus", so
  `--auto` would call `RunRound(2)`. But `runAgent` skips any participant whose artifact already exists unless
  `Overwrite` is set (`internal/runner/runner.go:406-418`, event `agent.skipped`, "artifact already
  exists"), so only claude-1 would launch. And the round-2 gate `roundComplete` (`internal/driver/driver.go`)
  checks only validity, `responding-to` and a `### @<other>` heading. codex-1's scaffold has all three, so the
  driver would count the placeholder as codex-1's completed cross-review. The organizer may not edit or
  remove a participant file, so round 2 is relaunched with the recorded runner-template fallback, as
  RESUME-BRIEF-1 step 2 prescribes. The driver takes over again for the next step.
- **Relaunch inputs.** The re-rendered phase-2 context (`parley protocol packet --phase 2 --track
  deliberation --idea <slug> --flag protocol_change`) is full, with source = packet = `b273af1e…f388` and no
  fallback. Its 115166 B body is byte-identical to the body embedded in the original round-2 prompts. Both
  prompts also end with the two committed round-01 files verbatim (52816 B, regenerated as
  `gatherPriorRounds` builds them). Each relaunch prompt is the agent's original round-2 prompt plus one
  notice of orchestration facts, inserted after the "You are <agent>…" line. The notice says the earlier
  processes were killed by the launching session's exit, not by a provider error. For codex-1 it says that the
  existing file is its own interrupted scaffold and that it replaces it with its completed round, which
  overrides the template's "Do not overwrite the file if it already exists" for that one file only. For
  claude-1 it says that it has no round-2 file yet. No position or opinion is passed. The prompts are
  `runs/<id>/round-02/<agent>/relaunch-1/prompt.txt`: codex-1 is 171300 B (sha256 prefix
  `ded885798af6b7cc`, original plus 1067 B) and claude-1 is 170770 B (`f2f2b18342b2cf14`, original plus
  536 B). The original attempt's prompts and logs stay where they were, as evidence.
- Launcher: `runs/<id>/round-02/relaunch-1-launcher.sh`. It uses the same argv as `parley agents list`,
  the prompt on stdin, a fresh 1800 s alarm (`perl alarm`), the runner's `PARLEY_RUN_ID` /
  `PARLEY_AGENT_ID` / `PARLEY_PROC_MARKER` env, and for claude only the `cleanParticipantEnv` strip. Each
  participant starts in its own session (`setsid`), so it survives this organizer process.
- **Relaunched at 14:02:25 CEST.** codex-1 has wrapper pid 74002 and agent pid 74008. claude-1 has wrapper
  pid 74009 and agent pid 74019. Each wrapper is a session leader reparented to launchd. The codex banner
  confirms model gpt-6-astra, provider omniroute, approval never, sandbox workspace-write and reasoning effort
  max. The claude child env carries the three `PARLEY_*` marker keys and no `CLAUDECODE`, `CLAUDE_CODE_*` or
  `AI_AGENT*` key. Logs and `exit.json` go to `runs/<id>/round-02/<agent>/relaunch-1/`. The round boundary is
  process exit plus final content (observer gap 5).
- **claude-1 participant quota failure (14:04:32 CEST). Launching stopped.** The claude-1 round-2 process
  exited 1 after 127 s, with empty stderr and no artifact. Its entire stdout, verbatim: `API Error: 503
  [claude/claude-opus-5-5] Unavailable (reset after 55m 29s). This is a server-side issue, usually
  temporary — try again in a moment. If it persists, check your inference gateway (omniroute.marao.sk).`
  It has the same shape as evidence incident 4 (an omniroute 503 "reset after", a quota window). The reset
  is due at about 15:00 CEST. As the brief requires (the current protocol applies, and dropping either agent
  would go below the owner's minimum of 2), the organizer wrote the blocking note
  `inbox/claude-1-to-user_…_quota.md` with the verbatim error and launches nothing more. There is no retry,
  and the driver is not started. The organizer process uses the same model and gateway and still gets
  responses. This is recorded as an observation only. codex-1's round-2 process keeps running normally.
  The one 503 line in its stderr is its own read of the evidence file (incident 4 quoted), not an error.

## Phase 2: owner answer and claude-1 relaunch 2 (2026-10-03, from 18:04 CEST)

- **Organizer gap 8: a headless organizer cannot wait in prose.** The previous resumed organizer (pid 57867)
  committed `76dc2b7`, ended its turn with "I'm now waiting for codex-1's round-2 exit and the owner's
  answer", and exited 0 at about 14:07 CEST. A headless `claude -p` session ends when its turn ends, so its
  stated 16:30 watch never ran. From now on the organizer waits only with blocking tool calls (bounded
  `until` loops on process exit plus file, under 10 minutes per call) and ends its turn only at the end of
  the session's work or on a blocking escalation.
- Resumed at about 18:04 CEST as organizer pid 1534 in the tmux session `qae-organizer`, per
  `runs-handoff/.../RESUME-BRIEF-2.md`. Re-oriented from `parley organizer brief`, `parley status` and
  these notes. No participant process was alive.
- **codex-1 round 2 is complete.** Its relaunch exited 0 at 14:14:29 CEST after 724 s, with stderr reporting
  "tokens used 139,878" and no provider error. `round-02/codex-1.md` is 18720 B of final content (it ends
  "This file is a completed round-02 position") and validates. It replaced codex-1's own scaffold. Swept
  unmodified as `7595ee1` under codex-1's prefix. `git status` showed no edit outside its own file.
- **Owner answer** (non-blocking note `inbox/user-to-claude-1_…_quota-answer.md`, written 18:02 CEST):
  "Relaunch once after 15:05 (Recommended)". It arrived after pid 57867 had exited, so the relaunch had not
  happened. The relay reports a one-word liveness probe on `claude/claude-opus-5-5[1m]` through the same
  gateway at about 18:04 CEST that returned `PONG` in 4.6 s. That is the relay's fact, not a participant
  launch.
- **Relaunch 2 inputs.** `runs/<id>/round-02/claude-1/relaunch-2/prompt.txt` is the relaunch-1 prompt with
  only its notice block replaced: 172333 B, sha256 prefix `fd4e6a70e500d863`, and the diff against
  relaunch 1 is 1 changed and 9 added notice lines, so the protocol context and the template are
  unchanged. The notice adds three facts: the 14:04 failure with its verbatim stdout; the owner's answer
  verbatim, which the participant is asked to quote under `## User direction` (§4), as the escalation was
  claude-1's and the organizer writes no round file; and that codex-1's round 2 is complete and was
  written from round 1 only, so claude-1 also writes from round 1 and does not read `round-02/codex-1.md`
  in this round, with responses to it belonging to round 3 if one opens (Phase 2: "React in your own file
  in the next round"). No position or opinion is passed.
- **Relaunched at 18:07:21 CEST, once.** The launcher is `runs/<id>/round-02/relaunch-2-launcher.sh`,
  claude-1 only, with the same argv, alarm (1800 s), marker env and `cleanParticipantEnv` strip as
  relaunch 1. It also writes `launcher.pid`. macOS has no `setsid` binary, so the session is created with
  `perl` fork plus `POSIX::setsid`. Wrapper pid 10966 is a session leader reparented to launchd, and agent
  pid 10975 carries the three `PARLEY_*` marker keys and no `CLAUDECODE`, `CLAUDE_CODE_*` or `AI_AGENT*`
  key. Per the owner's answer, a second quota or credit failure stops launching, with a new blocking note
  that carries the verbatim error.

## Phase 2: resume 3, round-2 sweep, owner instruction and round 3 (2026-10-03, from 19:00 CEST)

- **Organizer pid 1534 died with the whole tmux server** after its relaunch-2 entry above (relay check at
  18:57 CEST, `runs-handoff/.../RESUME-BRIEF-3.md`). There is no `tmux-exit.txt`, so its command never
  exited. The participants launched with `perl` plus `POSIX::setsid` survived. Resumed at about 19:00 CEST
  as organizer pid 53360, launched the same way (a session leader reparented to launchd, not under tmux).
  Re-oriented from `parley organizer brief`, `parley status` and these notes. No participant process was
  alive.
- **claude-1 round 2 is complete.** Relaunch 2 exited 0 (`exit.json`: 16:07:21Z to 16:22:42Z, 921 s) with
  an empty stderr. `round-02/claude-1.md` is 19522 B of final content. It records the attestation (full,
  `b273af1e…f388`, no fallback), quotes the owner's quota answer verbatim under `## User direction`
  (checked against the inbox note), and states that it did not read `round-02/codex-1.md`. The only other
  files changed after the 18:07 launch are this file (the relaunch-2 entry, 18:08:21) and the git-ignored
  `graphify-out/` cache. That cache was rebuilt at 18:07:51 by the graphify post-commit hook of the
  codex-1 sweep `7595ee1`, committed at 18:07:48. No participant edited outside its own file. Swept
  unmodified as `bc7a8cb` under the participant's prefix.
- The answered quota escalation `claude-1-to-user_…_quota.md` and the owner's answer
  `user-to-claude-1_…_quota-answer.md` are archived under `inbox/archived/`. The answer is quoted verbatim
  in `round-02/claude-1.md`.
- **New owner instruction** (relayed in RESUME-BRIEF-3, about 18:55 CEST, Slovak, verbatim): "sakra tak
  pouzi len claude a codex" ("Damn, then just use claude and codex."). The relay's reading: the quorum is
  codex-1 and claude-1 only, through FINAL. The owner proposal plans implementation and review with codex-1
  and claude-1 only, without kimi-1 or zcode-1. With codex-1 implementing, claude-1 is the single
  non-implementer reviewer, and the §15.5 role concentration applies. The instruction is quoted under
  `## User direction` in `00-prompt.md` as the round-03 kickoff notice, which asks each participant to
  quote it in its round-03 file. It supersedes the kickoff sentence that let a later implementation or
  review phase add kimi-1 or zcode-1.
- **Round 3 opens.** This is a procedural call, provisional under §15.5. The reasons: neither round-02
  file responds to the other (codex-1 wrote from round 1 only, and the claude-1 participant was told not
  to read `round-02/codex-1.md`), and both end with a list of remaining disagreements. Phase 2 says
  "Address every other active agent explicitly" and "Continue until nobody has new substantive
  objections". Round 3 is the second of the three cross-review rounds the §4.0 deliberation cap allows.
  No position is adjudicated here.
- `00-prompt.md` also gains a correction of the organizer's own record. The zcode-1 `excluded:` line's
  "reset 2026-10-05 06:14 local" is the provider's UTC+8 wall clock. `reset_at` 2026-10-04T22:14:57Z is
  2026-10-05 00:14:57 CEST (raised by codex-1 in round 1, confirmed by the claude-1 participant in round 2
  as C5). The frontmatter line stays as recorded.
- **Driver.** `git worktree list --porcelain | grep -c '^prunable'` returns 0. `parley continue`
  (print-only) still recommends "Open round-02 (cross-review) before drafting consensus", because round 2
  ran through the fallback and the driver cursor is at round 1. In `internal/driver/driver.go`
  `advanceRound`, the re-entry check `roundComplete(2)` sees both valid round-02 files and promotes without
  re-dispatch. It then opens round 3 only if `1 + cross_review_rounds > 2`. The default is 1
  (`internal/driver/transport.go` `ReadCrossReviewRounds`), and the deliberation track caps it at 3
  (`internal/track/track.go`). So the organizer set `cross_review_rounds: 2` in `00-prompt.md`. After round
  3 the driver drafts consensus, requests signoffs and authors FINAL, and `--no-implement` stops it there
  (`internal/driver/loop.go`). Its round launches use the runner's template inside the runner's
  protocol-context envelope, with all prior rounds inline (`gatherPriorRounds`).
- **Driver gap 9: the driver has no channel for an organizer notice.** `parley steer` records intent only.
  `internal/steer/steer.go`: "records a QUEUED new attempt; it does NOT execute delivery". So the owner's
  instruction cannot reach driver-launched prompts through steer. The runner's task text tells participants
  to read `00-prompt.md`, so the notice is placed there.
- **Driver gap 10: the drafting prompts are thinner than the round prompts.** `runHeadlessSignoffAgent`
  (`internal/app/consensus_request_signoffs.go`) passes the drafting prompt unwrapped, with no
  protocol-context envelope or attestation. The consensus prompt (`internal/app/driver_consensus.go`)
  carries the §15.7 duties but not the §15.5 one-line role-concentration record. The FINAL prompt points
  the drafter at `consensus.md` and the round files, not at the brief's FINAL list in `00-prompt.md`. The
  organizer accepts this and records it. It checks `consensus.md` and `FINAL.md` against both when they
  land.
- Fact: `parley agents list` now shows codex-cli 0.160.0 (0.159.3 at preflight). Model and effort are
  unchanged: gpt-6-astra at max, and claude/claude-opus-5-5[1m] at max. The organizer made no roster or
  config change.
