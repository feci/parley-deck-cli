---
idea: parley-deck-cli-plan
status: final
author: codex
consensus-date: 2026-05-10
participants: [codex, claude]
---

## Final plan / specification

### Assumptions

- The first product promise is a local, self-contained CLI/TUI that helps a human supervise Parley Deck workflows.
- The source of truth remains the Parley Deck filesystem protocol under `parley-deck/`; runtime state is supplemental.
- The first release supports local CLI agents only. GitHub PR and GitLab MR transports are later work.
- Token and cost accounting are best effort because agent CLIs differ in whether they expose usage metadata.
- Human-in-the-loop is the default. Automatic mode must be available, explicit, auditable, and conservative.

### Recommendation

Use a priority-based stack decision.

If implementation must start before a spike, use:

- **Language/runtime:** Go.
- **TUI:** Bubble Tea + Bubbles + Lip Gloss.
- **Distribution:** native binaries for macOS, Linux, and Windows; npm package as a first-class launcher/downloader for `npx parley-deck`; Homebrew tap/cask; direct release archives; later Scoop/WinGet/Chocolatey.

Why: this best matches "completely self contained" and "installable as executable" while still supporting the Node path via npm.

If the user wants stronger evidence before choosing, run a short decision spike before locking the language:

1. Build a minimal Go + Bubble Tea supervisor that launches two local agents, writes JSONL events, renders a tiny dashboard, and cross-builds native binaries.
2. Build the same minimal TypeScript + Ink supervisor and test `npx`, Node single executable or Bun-compiled binaries, child process handling, Ctrl+C behavior, and Windows terminal behavior.
3. Pick the stack from measured results and the priority answers below.

### Stack comparison

| Stack | Best fit | TUI options | Distribution fit | Main risk | Verdict |
| --- | --- | --- | --- | --- | --- |
| Go | self-contained binary, subprocess supervision, predictable release packaging | Bubble Tea, Bubbles, Lip Gloss; tview as simpler alternative | excellent native binaries; Homebrew straightforward; npm requires wrapper/downloader | npm path is a second artifact and must be reliable | recommended provisional default |
| TypeScript / Node.js | npm-native CLI, React-style contributors, fastest `npx` path | Ink; Blessed/neo-blessed for lower-level terminal control | excellent npm path; native executable story is improving but needs validation | standalone executable and terminal edge cases are less certain than Go | best if `npx` is the canonical product |
| Python | fastest prototyping and rich widgets | Textual, Rich, prompt_toolkit | PyInstaller can make executables, but packages are heavier; `npx` is unnatural | distribution and startup size | do not choose unless TUI iteration matters more than packaging |
| Rust | strongest native binary and performance | Ratatui | excellent binaries; npm wrapper needed | slower implementation unless team wants Rust | optional later alternative, not v1 default |

### Initial feature set

1. `parley init`
   - Creates `parley-deck/COOPERATION.md`, `ideas/`, `inbox/`, `meta/`, and optional local config.
   - Verifies transport is `local-dir` for v1.

2. `parley agents discover` / `parley agents probe`
   - Finds configured CLIs.
   - Runs non-destructive `--help` and `--version` probes.
   - Records executable path, version, headless mode, write mode, model options if discoverable, telemetry support, and warnings.

3. `parley run "<task>"`
   - Creates a new idea and kickoff prompt.
   - Selects participants.
   - Starts round 1 in parallel.
   - Opens the TUI by default unless `--no-tui` is passed.

4. `parley status`
   - Prints current idea, phase, round, participant states, pending questions, and artifact completion.
   - Works in CI and plain terminals.

5. `parley resume <idea-or-run>`
   - Rebuilds state from durable events and protocol files.
   - Reattaches to live processes if possible, otherwise shows recoverable state.

6. `parley answer <question-id>`
   - Lets the user answer a HITL question without opening the TUI.
   - Writes the answer into durable run state and the next protocol artifact as required by the cooperation protocol.

7. `parley auto`
   - Runs the same protocol engine with automatic consensus policy enabled.
   - Requires explicit policy fields and stops on unresolved blockers, destructive actions, credentials, or private-data decisions.

### Architecture

Use five layers:

1. `protocol`
   - Parses `COOPERATION.md` and idea frontmatter.
   - Knows phase and round rules.
   - Computes expected artifacts.
   - Enforces "one agent, one file" and signoff sequencing.

2. `runner`
   - Launches agent CLIs with configured headless args.
   - Streams stdout/stderr.
   - Watches for output artifact creation.
   - Detects timeout, auth failure, empty output, non-zero exit, and blocked questions.

3. `adapters`
   - `codex`, `claude`, `gemini` first.
   - Each adapter describes command, version probe, headless mode, write mode, model/thinking flags when discoverable, and telemetry parser.
   - Unknown telemetry is explicit, not guessed.

4. `store`
   - `parley-deck/runs/<run-id>/events.jsonl`: append-only runtime events.
   - `state.json`: derived snapshot for fast resume.
   - `agents/<agent-id>/stdout.log` and `stderr.log`: raw process streams.
   - `questions/<question-id>.md`: pending HITL questions and answers.

5. `tui`
   - Reads derived state and event stream.
   - Sends commands to the supervisor: answer, pause, resume, retry, cancel, switch mode.
   - Does not replace protocol files as source of truth.

### Core state model

Agent states:

- `queued`
- `probing`
- `running`
- `writing-artifact`
- `blocked-question`
- `blocked-error`
- `timed-out`
- `done`
- `signed-off`
- `skipped-non-quorum`

Protocol states:

- `kickoff`
- `round-NN`
- `cross-review-NN`
- `consensus`
- `final`
- `implementation`
- `review-round-NN`
- `review-consensus`
- `fix-up-NN`
- `complete`

Metric fields:

- always available: elapsed time, start/end time, process status, retries, stdout/stderr byte counts, artifact path, artifact validity.
- best effort: input tokens, output tokens, total tokens, cost, model, reasoning/profile.
- never infer exact model tokens from raw text unless the backend reports them.

### TUI design

Initial layout:

```text
┌ Parley Deck: parley-deck-cli-plan ─ phase: round-01 ─ mode: HITL ─ elapsed 08:14 ┐
│ Timeline: kickoff > round-01 running > consensus pending                         │
├ Agents ────────────────────────┬ Live Activity ──────────────────────────────────┤
│ codex   running  2m14s  9.1k   │ codex: writing proposed architecture...          │
│ claude  blocked  question 4.2k │ claude asks: should npm be primary install path? │
│ gemini  skipped  non-quorum    │ gemini headless probe did not complete           │
├ Protocol ──────────────────────┼ Stats ──────────────────────────────────────────┤
│ 00-prompt.md                   │ input 18.2k  output 4.1k  total 22.3k            │
│ round-01/codex.md done         │ cost unknown wall 8m14s failures 0               │
│ round-01/claude.md done        │ per-agent table                                  │
│ consensus.md pending           │                                                   │
├ HITL Inbox ──────────────────────────────────────────────────────────────────────┤
│ [1] claude asks: Pick priority: easiest install, richest TUI, or fastest MVP?     │
└ Enter=respond  Tab=panel  a=auto  p=pause  r=retry  q=quit-but-keep-running ─────┘
```

First release uses one dashboard. Later releases may add tabs:

- `Overview`: phase, timeline, agent table, global stats.
- `Agents`: per-agent logs, token/time data, retries, current prompt.
- `Artifacts`: protocol tree and validation state.
- `Questions`: HITL queue and answer history.
- `Settings`: mode, policy, participants, timeout.

### HITL and automatic mode

HITL default:

- Pause when an agent asks a user-facing question or needs permission.
- Show the question in the TUI and `parley answer`.
- Record the user's answer in durable state.
- Quote the answer into the next round/review artifact when protocol requires it.

Automatic mode:

- Must be explicitly selected with `--auto` or from the TUI.
- Uses the same protocol engine.
- Advances only when all required artifacts exist and the configured consensus policy is satisfied.
- Falls back to HITL when there is a blocker, missing artifact, timeout, destructive action, credentials, or sensitive-data decision.
- Should default to unanimous consensus for design choices. Majority/weighted policy can be added later, but should not be v1 default.

### Packaging plan

For Go default:

- Build with Go.
- Release cross-platform archives and checksums.
- Use Homebrew tap/cask for macOS.
- Publish `@parley-deck/cli` to npm with a `bin` entry.
- The npm package should install or dispatch the exact platform binary with checksum verification and clear offline instructions.
- Keep npm wrapper tests in CI because `npx parley-deck` is a required path, not a demo.

For TypeScript alternative:

- Publish npm package directly.
- Evaluate Node single executable applications and/or Bun compile during the spike.
- Test Windows ConPTY, Ctrl+C, process cleanup, log streaming, and asset bundling before selecting this route.

For Python alternative:

- Use Textual for UI only if distribution constraints are relaxed.
- Use PyInstaller for executable packaging if chosen.

### Milestones

1. **M0: repo bootstrap**
   - `go mod init` or TS package setup after stack decision.
   - Basic CLI parser.
   - Tests for config parsing and protocol paths.

2. **M1: protocol and event store**
   - Read idea frontmatter.
   - Compute expected round artifacts.
   - Append JSONL events.
   - Rebuild state from events.

3. **M2: agent probes and one runner**
   - Implement `agents discover`.
   - Implement Codex adapter first.
   - Launch one agent, capture logs, detect artifact written.

4. **M3: parallel multi-agent run**
   - Add Claude and Gemini adapters.
   - Launch agents in parallel.
   - Track running/done/error/timed-out.
   - Validate exactly one artifact per agent.

5. **M4: dashboard TUI**
   - Render overview layout.
   - Show phase, round, timeline, agent table, live logs, stats, artifact state.
   - Quit/reopen without losing state.

6. **M5: HITL**
   - Detect questions.
   - Pause relevant agent/phase.
   - Answer via TUI and `parley answer`.
   - Persist answer and protocol quote.

7. **M6: consensus and auto mode**
   - Draft/validate consensus state.
   - Sequential signoffs.
   - Conservative auto policy.
   - Fallback to HITL on ambiguity.

8. **M7: release**
   - Native binaries.
   - `npx parley-deck`.
   - Homebrew.
   - Smoke tests on macOS, Linux, and Windows.

### Questions for the user

Please rank these before implementation:

1. Is the top priority **one self-contained executable** or **`npx parley-deck` as the canonical path**?
2. Do you want to run a short Go-vs-TypeScript spike first, or choose the stack now?
3. Should v1 optimize for the **smallest reliable dashboard** or a **more polished TUI**?
4. Should auto mode be allowed to continue through low-risk design questions, or should any agent question pause for you by default?
5. Is best-effort token tracking acceptable, with `unknown` shown when a CLI does not expose usage?
6. Which agent CLIs must be supported first: Codex, Claude, Gemini, something else?
7. Which OS is the release gate for v1: macOS only, macOS + Linux, or all three including Windows?
8. Should v1 include only local-directory transport, or must GitHub/GitLab PR/MR integration be in the first usable release?

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/, ./round-02/
- Node single executable applications: https://nodejs.org/api/single-executable-applications.html
- npm exec / npx behavior: https://docs.npmjs.com/cli/v11/commands/npm-exec/
- Bubble Tea: https://github.com/charmbracelet/bubbletea
- Ink: https://github.com/vadimdemedes/ink
- Textual: https://textual.textualize.io/
- PyInstaller operating mode: https://pyinstaller.org/en/stable/operating-mode.html
- GoReleaser npm publishing notes: https://www.goreleaser.com/customization/publish/npm/
- Homebrew taps: https://docs.brew.sh/Taps
