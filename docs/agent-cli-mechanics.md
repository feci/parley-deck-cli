# Verified agent-CLI invocation mechanics

Empirical, hard-won behaviors of the roster CLIs as parley invokes them
(runner-hardening-kindly D8; several entries adapted with attribution from the
MIT-licensed "kindly" skill's reviewers reference — see
`parley-deck/ideas/runner-hardening-kindly/reference/`). Re-verify against the
installed CLI when a flag misbehaves, then update THIS file and the runner spec
together — never work around the runner with hand-rolled calls.

## codex (`codex exec`)

- **`</dev/null` is mandatory in non-TTY shells.** With a prompt argument AND an
  open piped stdin, `codex exec` also reads stdin (appends it as a `<stdin>`
  block) and hangs forever waiting for EOF. Symptom: log stops at
  "Reading additional input from stdin...".
- `-o <file>` writes exactly the final assistant message; the file is absent on
  failure, so "exit 0 + non-empty file" is the success oracle for message-style
  outputs. parley does NOT use `-o`: our contract is an artifact FILE the agent
  writes itself (plus the stdout-fallback for print-only responses), which `-o`
  would not capture.
- `--json` streams NDJSON events on stdout; the session id arrives first as
  `{"type":"thread.started","thread_id":"…"}`. The human stderr banner
  disappears under `--json` — parse events, not the banner.
- Sandbox: always pass it explicitly (`--sandbox workspace-write` for
  participants); a user config can default to `danger-full-access`. Outbound
  network in workspace-write stays blocked unless
  `-c 'sandbox_workspace_write.network_access=true'`.
- `codex exec resume <session-id>` has NO `--sandbox` flag — only
  `-c sandbox_mode="…"` pins it; it also has no `-C`, so resume from the same
  cwd that created the session.
- Known sandbox artifact: `sysctl kern.boottime` is restricted under seatbelt,
  which fails parley's `TestDurableKillEndToEndRealProcess` inside codex runs.

## claude (`claude -p`)

- **Bind the prompt to `-p` directly — never as a trailing positional.**
  `--tools`, `--allowedTools` and `--add-dir` are variadic and silently swallow
  a trailing prompt; claude then parks forever in its event loop with zero
  output and near-zero CPU.
- **`--tools` REMOVES built-in tools from the model; `--allowedTools` merely
  pre-approves.** A visible-but-denied tool under `dontAsk` produces a denial
  loop that looks exactly like a hang. Restrict the toolset, pre-approve what
  remains, and tell the model what it has.
- **MCP tools bypass `--tools`** — configured servers expose their tools
  regardless; pre-approve per server (`mcp__<name>__*`) or expect silent
  denials under `dontAsk`.
- **Session resume is cwd-scoped:** `--resume <id>` finds the session only from
  the directory that created it ("No conversation found" otherwise).
- Never pass a spend cap (`--max-budget-usd`) to a deep run: it aborts mid-work
  with "Exceeded USD budget" before any output lands.
- parley sheds nested host markers (CLAUDECODE, CLAUDE_CODE_*, AI_AGENT*) when
  spawning claude as a participant — a participant must not inherit the
  facilitator's session identity (runner `cleanParticipantEnv`).

## agy (Antigravity CLI)

- `--print`/`--prompt` is VALUE-taking: the prompt must be the flag's value and
  the last token (`agy … --print "<prompt>"`); never feed the prompt on stdin
  with a bare `--print`.
- `buffers_stdout = true` (runner spec default for agy): ALL stdout is buffered
  until process exit — a silent transcript is expected, not a hang; stderr is
  live. The runner's stdout-fallback recovers artifacts that agy printed
  instead of writing.
- Known flake (now harmless): agy sometimes writes a valid artifact and then
  exits 1 — the runner's artifact-wins rule records success with `agent_exit`.
- Appending to existing shared files is flaky; prefer prompts that write new
  files, and expect retries on append tasks.

## hermes

- The prompt is the VALUE of `-z`/`--oneshot` (`hermes --yolo -z "<prompt>"`),
  with `</dev/null`. Putting `--yolo` after `-z` makes argparse treat it as the
  prompt ("expected one argument").
- **Silent-death mode:** hermes can go permanently mute mid-session — every
  invocation exits rc=0 with zero stdout/stderr while `--version` still works;
  its agent.log stops right after plugin discovery (init takes ~40s — do not
  kill early and misdiagnose). Cause observed: backend auth/credit outage.
  The runner's first-output watchdog catches this in ~2 minutes and classifies
  it `no_first_output`; probe with a cheap "PONG" prompt before counting on
  hermes.

## General

- Run every participant in the FOREGROUND of its own background task; a `&`
  inside a wrapper orphans the agent when the wrapper exits.
- Sequential signoff appends: invoke signers one at a time — concurrent
  read-modify-write appends to one consensus file clobber each other.
- All read-only git probes set `GIT_OPTIONAL_LOCKS=0` so they never write
  `.git` on the weakly-coherent virtio-fs mount.

See also: `docs/agent-runtime-configuration.md` (spec/TOML knobs, including the
supervision windows `first_event_timeout_ms`, `stall_timeout_ms`,
`heartbeat_ms` and the `buffers_stdout` flag).