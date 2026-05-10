---
agent: claude
idea: parley-deck-cli-plan
round: 2
date: 2026-05-10
responding-to: [codex/round-01]
---

## Position changes since prior round

I am partially moving toward codex's framing.

- I concede that "two of the three discovered CLIs are Node" is **not** a
  decisive reason to pick TS for the orchestrator. The hard parts —
  process supervision, durable state, crash recovery, terminal control,
  release packaging — are language-orthogonal. I overweighted runtime
  alignment in round 1.
- I concede the plan should be priority-based, not a flat
  recommendation. The user's answer to "what is the product promise?"
  genuinely flips the language choice.
- I still do not concede that Go is the **default**. The Bun-compile
  risk codex cites is real but bounded and testable; the cost of going
  Go (verbose protocol-file munging, an npm wrapper that downloads
  binaries from a separate release pipeline) is also real and
  permanent, not a one-time spike.

Net: we agree on architecture and on letting the user choose; we
disagree on which option is the **labelled default** if the user
declines to choose.

## Responses to others

### @codex

Strong points I accept:

1. The orchestrator's hard parts are language-orthogonal. The
   "Node CLIs are Node" argument is weak and I am dropping it.
2. Framing the plan as a priority decision is better than a flat
   recommendation. I will adopt that frame.
3. Token "unknown" must be a first-class visible state — agreed.
4. SQLite is over-scoped for v1; JSONL + snapshots are sufficient.
5. The TUI must not be the only interface. `run`, `status`, `answer`,
   `resume`, `agents discover` must all work non-interactively. I had
   this implicit; I should have stated it.

Where I push back, and counter-proposal:

**On "Go as default."** Codex's argument is essentially: native binary
risk is lower than Bun-compile risk, therefore Go is the safe default.
Two issues with that:

- Bun-compile risk is bounded and **measurable in a 1-day spike**:
  build a Bun binary on macOS/Linux/Windows, launch a child process,
  draw an Ink TUI, ctrl-c handling, ConPTY behavior. Either it works
  or it doesn't. We do not need to default-decide before measuring.
- The Go default has an asymmetric cost that is not in codex's
  framing: the npm wrapper becomes a **separate release artifact**
  with its own version drift, its own download/cache logic, its own
  failure mode (proxy/firewall blocks the binary download, signature
  mismatch, partial downloads). Calling `npx parley-deck` then means
  "download a binary at first run." That is materially different from
  a real Node package, and the failure surface tends to bite users in
  CI and corporate networks.

**Counter-proposal: spike before defaulting.**

Rather than the plan saying "default Go, alternative TS" or "default
TS, alternative Go," the plan should commit to a **2-day decision
spike** before locking the language:

1. Day 1 — TS/Bun spike. Build a minimal supervisor: launch
   `claude --print` and `codex exec` headlessly, write JSONL events,
   render a tiny Ink dashboard. Compile with `bun --compile` for all
   three OSes in CI. Measure: binary size, cold start, child-process
   reliability on Windows ConPTY, ctrl-c behavior, signed/unsigned
   install friction.
2. Day 1 — Go spike (in parallel). Same scope: launch the same two
   CLIs, write the same JSONL, render a Bubble Tea dashboard,
   GoReleaser cross-build, ship a minimal npm wrapper that downloads
   the right binary.
3. Day 2 — Decide. The spike that hits all three OS targets cleanly
   wins. If both pass, the user picks based on the priority list
   below; if one fails, the choice is made for us.

This costs 2 days and removes the largest source of round-1/round-2
disagreement (a hypothetical Bun-compile risk vs a hypothetical Go
npm-wrapper friction). I would rather measure than guess.

**On "if forced to default without a spike,"** I will move to: codex
chooses Go, I choose TS, and the plan presents both with a clear
priority matrix. I will not block consensus on this if the user picks
Go after seeing the matrix.

## New concerns / questions

- **Spike or no spike** — does the user want a 2-day language-decision
  spike, or do they want to pick on paper today? This is the single
  decision that unblocks everything else.
- **Wrapper semantics if Go wins** — if Go is chosen, the npm package
  must be a real package with version-pinned binary download, SHA
  verification, and a documented offline-install path. Not "hello
  world that shells out." The plan should call this out as scope.
- **TUI as required vs optional** — the constraint says the CLI must
  show phase, round, agent activity, stats, HITL questions, and
  consensus progress. Some of that can be a non-interactive `status`
  command (good for CI). Some genuinely needs a live TUI. I want the
  plan to mark each requirement as `tui-required` or `cli-sufficient`
  so we don't conflate them.
- **Headless agent capability drift** — agreed with codex this is a
  real risk. I want a `parley agents probe` command that captures the
  observed CLI version and output shape on each run, so a future
  silent format change is detectable from the event log.
- **Quorum policy for `gemini`** — `gemini` is currently excluded from
  quorum per `00-prompt.md`. The plan should specify how an excluded
  participant is rendered in the TUI (visible-but-non-voting) so the
  user is not surprised.

## Current proposal

Adopt codex's priority-based framing, plus a 2-day spike to ground
the language choice in measurement rather than argument.

Decision matrix the user is asked to confirm:

| If your top priority is…                                | Pick     |
|---------------------------------------------------------|----------|
| One signed binary, zero runtime deps, lowest TUI risk   | Go       |
| `npx parley-deck` is the canonical install, no wrappers | TS/Node  |
| Fastest TUI iteration, richest widgets                  | Python   |
| Hackability by Parley Deck users (most are JS/TS folks) | TS/Node  |
| Lowest packaging-and-release operational burden         | Go       |

Plan, regardless of language:

1. Run the 2-day spike. Pick language from measured results + matrix.
2. Implement supervisor + JSONL event log + filesystem-as-protocol
   exactly as both codex/round-01 and claude/round-01 described —
   we already agree here.
3. Implement adapters for `codex`, `claude`, `gemini`. Token
   accounting best-effort, "unknown" rendered explicitly.
4. Ship both `npx parley-deck` and signed standalone binaries from
   day one — neither is a stretch goal. If Go wins, the npm wrapper
   is a real package with verified downloads, not a launcher stub.
5. Non-interactive commands (`run`, `status`, `answer`, `resume`,
   `agents discover/probe`) work without ever entering the TUI.
6. TUI is a read-only projection over `events.jsonl`; killing it
   never kills the run.
7. HITL default; auto mode is the same supervisor with a
   policy-driven auto-answer; auto mode requires explicit
   `decision:` fields and falls back to HITL on close calls.
8. Defer SQLite, GitHub/GitLab transports, and multi-idea dashboards.

If the user wants a single answer today without a spike: **I will
not block consensus on Go**. The architecture, event model, adapter
contract, TUI layout, and HITL/auto semantics from round 1 carry
over unchanged. The language choice is the only thing the spike
would have settled.
