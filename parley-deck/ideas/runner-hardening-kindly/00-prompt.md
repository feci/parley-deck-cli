---
idea: runner-hardening-kindly
author: user
created: 2026-06-12
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + synthesis owner; implementer
  codex: Go correctness — runner/process supervision, event semantics, git plumbing, tests
  agy: UX + operator ergonomics — failure messages, recovery hints, consult output, TUI surfacing
  hermes: filesystem/process correctness — snapshot clones on virtio-fs, watchdog timing, kill semantics
transport: local-dir
cross_review_rounds: 1
status: final
---

## Task (owner's directive)

Adopt the six CLI-side improvements identified in the analysis of the MIT-licensed
"kindly" skill (Softmachine/Meta — copies with attribution in `reference/`; see
especially kindly-reviewers.md and kindly-agent.sh) into parley-deck-cli, and ship
them in release 1.24.0 together with the sibling protocol idea
`meta-protocol-change-review-gate-honesty`. All six points are IN SCOPE by owner
decision — deliberation refines the design, not the scope.

## The six points

**P1 — First-event watchdog + stall guard + heartbeats (runner supervision).**
Today the runner has only a hard per-agent timeout (default 30m). Adopt kindly's
three-layer supervision: (a) no first output/event within a grace window (default
~120s) → kill the agent process tree, retry once, then fail loudly with a typed
event; (b) after first output, a long silence with zero output-byte growth
(default ~30m, but bounded by the existing timeout) → classify as STALLED, kill
with diagnostics (distinct from a healthy deep run); (c) periodic heartbeat
events (elapsed, bytes) so the TUI narrator can show progress. Note the runner
already writes per-agent stdout/stderr log files — growth tracking can stat them
(the TUI's 1.23.0 growth cache is precedent). Decide: which knobs are config
(agents.toml fields? defaults?), which events are emitted (e.g.
`agent.stalled`, `agent.no_first_output`, heartbeat — or reuse/extend existing
types), and how this composes with ACP-mode agents (events vs file growth).

**P2 — Failure classification + recovery hints.**
On agent failure, scan the agent's stderr + captured output with bounded regex
classifiers (rate-limit / overloaded / auth / billing / invalid-request /
model-not-found / context-window / sandbox / budget — see kindly-agent.sh
stderr_failure_details + recovery-hint functions) and attach
`failure_class` + `recovery_hint` to the `agent.failed` event payload. TUI
narrator + agent status header surface the hint ("hermes appears auth-broken —
run hermes auth"). Keep classifiers data-driven and bounded (tail N lines).

**P3 — Usable artifact beats exit code.**
Today `eventType = agent.failed` when ExitError != "" OR !ArtifactOK. Adopt
kindly's nuance: when the artifact validates but the process exited nonzero,
record SUCCESS (`agent.finished`) with an `agent_exit` field carrying the code
(+ a narrator-visible note), instead of failing the step. This directly removes
the known agy flake (valid artifact then exit 1). Decide: does the same apply to
review/fixup/steer phases; any case where nonzero exit + valid artifact should
still fail?

**P4 — Small hardening batch.**
(a) When launching the `claude` CLI as a participant from a Claude Code host,
shed nested host markers (env -u CLAUDECODE, CLAUDE_CODE_SESSION_ID,
CLAUDE_CODE_ENTRYPOINT, AI_AGENT...) — a participant is independent by
definition; (b) set GIT_OPTIONAL_LOCKS=0 on all read-only git probes
(driver gitTreeClean, status/diff probes) so probes never write .git on the
weakly-coherent virtio-fs mount; (c) document the verified per-CLI invocation
mechanics (codex `</dev/null` + `-o` final-message file; claude `-p` prompt
binding gotcha, `--tools` removes vs `--allowedTools` pre-approves, MCP bypasses
--tools, cwd-scoped resume; agy --print value-taking; hermes -z) in
`docs/agent-cli-mechanics.md` + referenced from the skill; decide whether codex
output capture should switch to `-o`.

**P5 — Snapshot checkout isolation for Phase 6 reviews.**
Reviewers currently read the live working tree. Adopt kindly's snapshot model:
before opening a review round, create a disposable shared-clone checkout
(`git clone --shared --no-checkout` + detached checkout of the reviewed commit)
under the user temp dir; reviewers run with that checkout as their root; teardown
deletes it. The reviewed-commit pin becomes physical and the live tree stays
free. Mind virtio-fs (clone target should live on the FAST local tmp, not the
shared mount — decide), crash healing (pid markers), and which runner phases use
it (review rounds yes; fix-up no — implementer needs the live tree). Fallback to
live tree with a loud event when snapshot creation fails.

**P8 — `parley consult` command.**
A lightweight advisory cross-agent question with repo context — the gap between
an inbox ping and a full idea. Shape: `parley consult <agent> "<question>"`
(or stdin) → invokes the agent read-only-ish with the workspace as context, the
question as brief; writes the answer as a durable artifact (proposal:
`parley-deck/consults/<ts>-<agent>-<slug>.md` with provenance frontmatter +
an index/event line) and prints it. No protocol ceremony, no signoffs; the
protocol idea (sibling) defines its standing. Decide: artifact location +
frontmatter, timeout, whether agents may write (kindly: consult is read-only —
but our agents need write access to produce their artifact file; resolve, e.g.
facilitator captures stdout into the artifact for consult mode).

## Constraints

- One release (1.24.0) ships everything, together with the sibling protocol idea.
- Workspace lives on virtio-fs — bounded I/O; reuse fsutil where mkdir is involved.
- Honor existing safety invariants: durable-kill attribution gates, process-group
  spawning (procctl), events.jsonl as source of truth, gitTreeClean gates.
- Keep healthy paths zero-cost; supervision must not add per-tick I/O beyond the
  existing log tailing.
- English-only under parley-deck/.

## Questions for round-01

1. P1: exact watchdog timings + event names + config surface; how to detect
   "first output" for exec-mode vs ACP-mode agents; interaction with the
   existing TimeoutMS and durable-kill machinery.
2. P3: precise decision table (artifact_ok × exit × phase) — where does
   artifact-wins apply?
3. P5: snapshot location (local tmp vs shared mount), lifetime, healing, and
   how runner threads the alternate root into agent invocation (cwd, --add-dir,
   prompt paths).
4. P8: command UX + artifact schema + how consult answers are kept out of
   quorum/protocol state.
5. Sizing: propose implementation slices that ship as ONE release.
