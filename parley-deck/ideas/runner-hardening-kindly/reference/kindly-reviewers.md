# Reviewer Mechanics

`agent.sh` encapsulates everything below; read this when diagnosing a failed wrapper run or rechecking the CLI mechanics it relies on. The same read-only invocation backs both the audit gate and consult mode; an `--exec` audit is the one variation — it widens the sandbox and adds a shell. Pass `--reviewer claude` from Codex and `--reviewer codex` from Claude unless host-agent markers are known to reach the shell; without a detectable host, the script fails closed rather than guessing, and when the chosen reviewer matches the detected host it warns and records `same-reviewer: yes` in the header. These are empirical CLI behaviors — re-verify against the installed CLIs when a flag misbehaves, then fold any needed flag changes back into `agent.sh` instead of bypassing it.

## Launch and Supervision

A review or consult runs for a long time at full reasoning depth, so launch it with the host's reliable long-running task mechanism and check in at intervals. The script keeps report content in the final artifact, prints `report:` and `session:` when a usable report lands, and writes content-free progress heartbeats to stderr (elapsed time, event count, byte count, last event type) while the reviewer is still running. Treat quiet stdout as normal; use the stderr heartbeats, process state, and the eventual report path to supervise the run. Use the wait actively — the snapshot checkout keeps the live tree free, so the default wait activity is an independent adversarial self-review of the same scope.

- **Claude Code host** — start `agent.sh` as a background Bash task (`run_in_background`): the host returns a handle at once and notifies you when the run exits, then read the printed `report:` path. Treat this Bash-tool launch as tied to the host session unless you have reverified the current Claude background-session surface for this exact command shape.
- **Codex host** — start `agent.sh` as a normal long-running shell-tool call and poll the returned tool session until the `report:` line lands. The shell-tool wait window is a polling/output window, not the reviewer lifetime, so a quiet foreground run can keep going while you poll. Treat `nohup ... &` as a host-specific optimization to re-verify before use, not the default launch strategy; some shell-tool wrappers reap background children when the tool call closes. If the wrapper receives `INT` or `TERM`, `agent.sh` leaves a `*-interrupted.log` in the report directory rather than a silent empty log.

Two watchdogs back the supervision: no first event within the grace window (default 120s) means kill, retry once, then fail loudly; after the first event, a long silence with zero new event bytes (default 30 minutes) means the run is treated as stalled — killed with diagnostics saved — because a healthy deep review keeps emitting events. The windows take env overrides (`KINDLY_FIRST_EVENT_SECS`, `KINDLY_STALL_SECS`, `KINDLY_HEARTBEAT_SECS` for the heartbeat cadence, and `KINDLY_MCP_LIST_SECS` for the bounded MCP enumeration) so the test harness can exercise them quickly.

## Snapshot Checkouts

Read-only audits review a disposable snapshot checkout instead of the live tree, so the caller keeps working while the review runs and context reads cannot drift mid-pass.

- **Storage** — the snapshot is a shared clone of the repository: its object store reaches history through alternates, every object the snapshot creates lands inside the clone, and in steady state the audited repository's own `.git` is never written (healing a leftover from the earlier git-worktree implementation prunes its stale admin files once). Teardown is a plain delete with nothing left behind.
- **Scope capture** — uncommitted scope becomes a snapshot commit: a temporary index reads `HEAD`, stages the live working tree (staged, unstaged and untracked content) into the clone, and `commit-tree` parents the result on `HEAD` so the snapshot's own `HEAD` diff is exactly the review scope. Committed scopes get a clean detached checkout of `HEAD`, which keeps a dirty live tree from leaking into context reads.
- **Ignores** — ignored files and `.local/` stay out of the capture; the repo's `info/exclude` and any local excludes-file setting are carried into the clone so local-only ignore rules keep holding.
- **Stable path** — the snapshot lives at a stable path per repo and reviewer (under the user temp directory, whose shared base directory leaves with the last snapshot) because Claude session resume is cwd-scoped: recreating the checkout at the same path keeps `--resume` working for both CLIs. A concurrent run on the same repo and reviewer steps aside onto a unique path and says so; crashed leftovers — on the stable path or a stepped-aside one — are healed on the next run.
- **Fallbacks** — any creation failure falls back to a live-tree review with a loud warning and `snapshot: live-fallback` in the header. A snapshot commit holds one version per path, so when staged content diverges from the on-disk tree the script refuses the snapshot and reviews the live tree — its prepared diff carries the staged and unstaged versions separately, keeping both in scope.
- **Limits** — submodules are not populated inside the snapshot (the script warns at launch when the reviewed checkout defines them), and ignored build artifacts are absent — which is why `--exec` audits run in the live tree instead.

## Codex (`codex exec`)

Canonical reviewer invocation shape:

```bash
cd <review-root> && codex exec --sandbox read-only --json \
  -c 'web_search="live"' -c 'approval_policy="never"' \
  -c 'model_reasoning_effort="xhigh"' -o <report-file> "<prompt>" </dev/null
```

- `-o` writes exactly the final message; the file is absent on failure, so trust exit codes plus a non-empty-file check (a rare upstream bug can exit 0 without doing work).
- `--json` streams events to stdout; the session id arrives first as `{"type":"thread.started","thread_id":"<uuid>"}`. The human stderr banner disappears under `--json` — parse the event, not the banner. `agent.sh` uses this stream only for liveness and metadata while the report body stays in `-o`. The event stream carries no model field, so the report header's `model:` line (the detected model) appears only for Claude reviews; a `--model` pin shows on the `reviewer:` line and in the ledger for both CLIs.
- `</dev/null` is mandatory: with an open non-TTY stdin the CLI blocks forever on "Reading additional input from stdin...".
- A user config can default to `danger-full-access`; **always** pass the sandbox explicitly. Read-only is the default posture; `--exec` gates run `workspace-write` so the reviewer can execute tests and local infrastructure (the script flags mutations to tracked or untracked content in the report header; ignored files are not watched). Workspace-write keeps `/tmp` and `$TMPDIR` writable; `.git/` protection is sandbox-backend behavior, so verify on the target host when the test path depends on it. Run from the review root so cwd is the writable root. Outbound network — loopback included — stays blocked unless `-c 'sandbox_workspace_write.network_access=true'`, which the script passes for `--exec`. Treat network-dependent exec gates as sandbox-backend behavior: verify with `codex sandbox` on the target host when network access matters, and run the network portion outside the gate if the installed backend blocks it.
- `web_search = "live"` (top-level config key) enables the search tool; it runs outside the shell sandbox, so even read-only reviewers can check upstream documentation. The script pins it rather than trusting config drift.
- `approval_policy = "never"` keeps the non-interactive reviewer from stalling on inherited approval prompts; the sandbox and tool restrictions are the safety boundary for this read-only gate.
- The operator's curated MCP servers stay enabled and load from the user config untouched. They run outside the shell sandbox, so curation — not the gate — is the boundary for what they can reach; review the configured servers when adding new ones.
- A consult may target a non-git directory, so its codex invocation adds `--skip-git-repo-check`, and `agent.sh` defaults the repo to the current directory when it is not a git work tree.
- stderr carries harmless macOS keychain noise (`security: SecKeychainSearch...`) and hook lines; never merge it into the report.
- `codex login status` is read-only and exits 0 when logged in.

### Follow-Up Passes

```bash
cd <review-root> && codex exec resume <session-id> --json -c 'sandbox_mode="read-only"' \
  -c 'web_search="live"' -c 'approval_policy="never"' \
  -o <report-file> "<verify prompt>" </dev/null
```

- `resume` has **no `--sandbox` flag** — only the `-c sandbox_mode` override pins read-only.
- `--skip-git-repo-check` (consult only) is an `exec` option that precedes the `resume` subcommand; the `-c` overrides (sandbox, web_search) follow the session id, since `resume` accepts `-c`.
- `resume` also has no `-C` flag — the script runs it from the same stable snapshot path (or the repo, for consults) so the working directory matches the original session.
- The resumed session genuinely remembers prior passes, but its file reads may be stale — verify prompts must direct it to re-read changed files, and the script's resumed prompt does.
- Never use `--ephemeral` in gate mode: resuming an ephemeral session silently starts a blank thread instead of erroring.

## Claude (`claude -p`)

Reviewer when the host session is Codex, and the explicit Claude reviewer for diagnostics or panels:

```bash
cd <review-root> && claude -p "<prompt>" --permission-mode dontAsk --effort xhigh \
  --tools 'Read,Glob,Grep,WebSearch,WebFetch' \
  --allowedTools 'Read,Glob,Grep,WebSearch,WebFetch,mcp__<server>__*,...' \
  --add-dir <tmpdir> --output-format stream-json --verbose
```

- **Bind the prompt to `-p` directly — never as a trailing positional.** `--tools`, `--allowedTools` and `--add-dir` are variadic and silently swallow a trailing prompt; claude then starts with no prompt and an EOF stdin and parks forever in its event loop with zero output, zero network, and near-zero CPU.
- Run from the review root; cwd scopes file access.
- **Session resume is cwd-scoped**: `--resume <session-id>` finds the session only when invoked from the directory that created it ("No conversation found" otherwise) — this is why the snapshot checkout path is stable per repo and reviewer.
- **`--tools` removes built-in tools from the model; `--allowedTools` merely pre-approves.** With Bash visible but denied, the reviewer burns its turns in a denial loop that looks exactly like a hang (near-zero CPU for an hour). Restrict the toolset, pre-approve what remains, and keep `dontAsk` as the backstop — and tell the reviewer in the prompt what it has.
- **MCP tools bypass `--tools`** — configured servers expose their tools regardless, and under `dontAsk` a visible-but-unapproved MCP tool is silently denied. The operator's curated servers therefore stay enabled and get pre-approved per server: `claude mcp list` names them, and each becomes an `mcp__<name>__*` allow rule with the name sanitized the way the installed CLI does it: every character outside `[A-Za-z0-9_-]` becomes `_`, and names starting `claude.ai ` additionally collapse underscore runs and trim edge underscores, mirroring the binary's sanitizer verbatim (there is no all-server wildcard). The listing health-checks servers over the network, so the script bounds it (`KINDLY_MCP_LIST_SECS`, default 120s); if enumeration exits nonzero or times out the script discards the whole listing, warns, and proceeds with no rules — a partial listing is never treated as complete — and the stall guard backstops a reviewer that wedges on denied tools. MCP is capability here, not a boundary. Plugin-provided tools (for example LSP) also bypass `--tools`; they are part of the curated surface.
- No Bash by default: prefix allowlists cannot seal write-capable flags (`git diff --output=<path>` writes a file). The script pre-computes the scope diff into a temp file the reviewer Reads, with `--add-dir` granting access: for snapshot audits a single unified diff of the snapshot commit (untracked content included as additions); for live-tree fallbacks staged and unstaged sections separately, since `git diff HEAD` alone hides staged blobs behind a reverted worktree, plus an untracked-file list since `?? dir/` collapses whole directories. WebSearch and WebFetch stay in the toolset for upstream documentation.
- `--exec` gates add Bash inside the OS-level sandbox (`--settings '{"sandbox":{...}}'`): writes are scoped to the repo cwd plus temp dirs; outbound shell network remains constrained by the installed sandbox and proxy. `allowLocalBinding` permits binding loopback ports, but it does not by itself guarantee that a local test client can connect back to that server from inside the sandbox, so client-to-local-server integration tests should run outside the gate unless the installed Claude sandbox documents and verifies that path. Build tools that insist on writing user-level daemon or cache paths outside the repo and temp dirs can also fail under this contract; pass host validation output into the brief or configure those caches into an allowed location rather than weakening the review sandbox. `failIfUnavailable` makes a missing sandbox dependency a hard error rather than Claude's default warn-and-run-unsandboxed fallback, and `allowUnsandboxedCommands: false` blocks the per-command escape hatch — together they keep the contract fail-closed. Outside-cwd writes fail with `Operation not permitted`. WebSearch/WebFetch are native tools outside the Bash sandbox and keep working. The content-sensitive tree hash still backstops the report-don't-fix contract for in-repo writes.
- `--output-format stream-json --verbose` emits line-delimited events: the `system/init` event carries `session_id` and `model` (which the script records in the report header and ledger), the `result` event carries the report text — and the stream doubles as the liveness signal for both watchdogs. Follow-ups: `--resume <session-id>` with the same flags, from the same directory.
- If Claude exits nonzero after emitting a successful final `result`, `agent.sh` still preserves the report and records `reviewer-exit` in the report header; a nonzero wrapper exit means no usable report was written.
- The script also sheds the host session's nested markers (`CLAUDECODE`, session id, tasks/teams flags) via `env -u` — a reviewer is independent by definition.
- No turn cap: with the toolset sealed there is no denial loop left to guard against, and a cap kills a healthy deep review late with nothing to show.
- Never pass a spend cap: `--max-budget-usd` aborts the run with `Exceeded USD budget` the moment cumulative spend crosses it — before any findings return — and a real review costs more than a cautious cap allows.

## Failure Modes

- Exit 1 with a typed failure event on stdout (`type:"error"` or `type:"turn.failed"` for Codex; a failed `result` subtype, API-error fields, or Claude Code failure fields). Claude Code's documented StopFailure error names are `rate_limit`, `authentication_failed`, `oauth_org_not_allowed`, `billing_error`, `invalid_request`, `model_not_found`, `server_error`, `max_output_tokens` and `unknown`; provider/API surfaces may also expose overload, timeout, request-size, permission, usage-limit, context-window, sandbox, policy, or bad-request names. `agent.sh` copies bounded error details into the saved failed log, scans stderr on no-usable-report failures, and may add a recovery hint.
- Exit 2 — CLI usage error (bad flag combination).
- A stalled run (first event arrived, then nothing for the stall window) is killed and reported as a failure with diagnostics — inspect the failed log, then rerun or resume.
- Rate-limit exhaustion mid-gate: retry after the provider reset or switch to the other explicit reviewer and continue the gate; note the reviewer change in the report slug.
- Every outcome — report, failure, interruption — appends one line to the report directory's `index.jsonl` with timestamp, slug, mode, scope, pass, ref, reviewer, model, exec/snapshot flags, session ids, path, and exit code.

## Effort and Models

Audits and consults always run at `xhigh` reasoning on each CLI's default top-tier model — the script pins the effort on both paths (Codex via `-c 'model_reasoning_effort="xhigh"'`, Claude via `--effort xhigh`) and rejects `--effort` overrides, so depth is a constant of the gate rather than a knob. The reviewer model comes from each CLI's user config; `--model` exists to pin a specific top-tier snapshot when provenance demands it. A full run is intentionally slow and token-heavy — launch it as a background task, work the wait, and let it finish.
