# parley-deck-cli

> **Real multi-agent deliberation with a durable, reviewable audit trail** — not
> three agents in three terminals, and not one model role-playing a committee.

`parley` runs **Parley Deck**: a transport-agnostic protocol for getting several
AI agents to genuinely cooperate on a hard change — plus the CLI that makes it
usable instead of just specified. Each agent writes its own analysis, they
cross-review, reach a recorded consensus, implement, and review the
implementation — every step a file you can read, diff, and resume.

**Why not just spawn three agents in three terminals?** Ad-hoc multi-agent gives
you no audit trail, no conflict discipline, no consensus step, and no way to
resume. Parley Deck ships those as protocol.

**Why not ask one model to play a committee?** That's solo reasoning in a costume —
there's no second voice to actually disagree, and it reintroduces the single-model
self-preference that quorum-gated review exists to defeat. Parley Deck is non-solo
by design.

### What you get

- **An 8-phase idea lifecycle** (§4) — kickoff → independent analysis → cross-review
  → consensus → `FINAL.md` → `IMPLEMENTATION.md` → code review → fix-up.
  Append-only, and resumable from the documents alone.
- **Non-solo by design** (§1) — stable agent IDs (§2), one file per agent per round
  (§6); no agent overwrites another.
- **Compare, don't merge** — the consensus "Comparison & blind spots" lens rates
  confidence by agreement and surfaces contradictions and blind spots instead of
  averaging them away.
- **Transport-agnostic** (§0, §11) — `local-dir`, `github-pr`, or `gitlab-mr`: the
  same protocol whether agents share a filesystem or review each other through a PR.
- **Vendor/model-agnostic roster** — Claude, Codex, Gemini, GLM, and more by stable
  ID (subject to the CLIs you have installed and authorized).
- **Readiness preflight** (§9.0, `parley preflight`) — protocol-freshness plus a
  live roster ping before an idea starts; any exclusion is user-confirmed.
- **Advisory retrospectives** (§13, `parley retro`) — a quorum-gated pass over the
  deck's own history that *proposes* improvements through the normal workflow,
  never applies them automatically.
- **Supervised automation** — a live TUI and an auto-drive driver advance protocol
  *phases*; agent supervision (watchdog, stall guard, validated-artifact-beats-
  nonzero-exit) catches hung agents; code implementation and side effects stay gated.

### Inspired by — adopted & adapted

Parley Deck didn't invent these ideas; it wired them into one repository-backed,
quorum-gated protocol:

- **OpenRouter Fusion** → the compare-not-merge consensus lens
  (confidence-by-agreement, blind-spots), applied to asynchronous multi-round
  markdown instead of a real-time API ensemble.
- **OpenAI ExecPlans / PLANS.md** → resume-from-the-doc state, split into a static
  `FINAL.md` and a living `IMPLEMENTATION.md` governed by review-consensus.
- **RHO (Retrospective Harness Optimization)** → §13 retro, but advisory-only and
  quorum-gated instead of single-model self-preference.
- **kindly** → strict gates, stopping judgment, no-suppression review dispositions,
  and artifact-wins supervision.
- **Preflight readiness** → §9.0 protocol-freshness and roster liveness before each idea.

*Reference to these projects is for attribution and lineage only; no endorsement,
sponsorship, or affiliation is implied.*

## Install

Install the current checkout into `~/.parley-deck` while developing:

```bash
cd /path/to/parley-deck-cli
scripts/install-local.sh
```

The binary is installed as:

```text
~/.parley-deck/bin/parley
```

Add it to your shell path if needed:

```bash
export PATH="$HOME/.parley-deck/bin:$PATH"
```

Verify:

```bash
parley version
parley version --all
parley help
```

The release version follows semantic versioning and is recorded in `VERSION`. `parley version --all` also reports `parley-deck-skill` installer, runtime skill, and project metadata status when the skill installer is available.

Re-run `scripts/install-local.sh` after pulling or building new changes to replace the installed binary with the latest local version.

Options:

```bash
scripts/install-local.sh --dry-run
scripts/install-local.sh --prefix /tmp/parley-test
scripts/install-local.sh --bin-dir "$HOME/bin"
```

Homebrew users can install or update from the tap:

```bash
brew update
brew install feci/parley/parley-deck-cli
brew upgrade feci/parley/parley-deck-cli
```

## Using Parley Deck

Initialize a repository:

```bash
parley init --dir .
```

Inspect and verify available agents:

```bash
parley agents list --dir .
parley agents verify --dir . --agent claude
parley agents verify --dir . --full --agent codex --yes
```

Start a multi-agent round:

```bash
parley run --dir . --participants claude,agy --yes "Plan the next CLI slice"
```

Run without opening the live TUI:

```bash
parley run --dir . --no-tui --participants claude,agy --yes "Review the repo-map MVP"
```

Resume or inspect work:

```bash
parley status --dir .
parley status --dir . --idea repo-map-mvp
parley continue --dir . 20260517T120000.000000000Z
parley resume --dir . 20260517T120000.000000000Z
```

`continue` prints the planner-derived next safe workflow action for an existing
run or idea. `resume` remains the read-only detail/live-view command.

Answer a human-in-the-loop question:

```bash
parley answer --dir . 20260517T120000.000000000Z q1 "Use the conservative default"
```

Generate repository context for prompts or tools:

```bash
parley context repo-map --dir . --format markdown --max-files 50
parley context repo-map --dir . --format json --max-files 10
```

Manage consensus:

```bash
parley consensus status --dir . repo-map-mvp
parley consensus draft --dir . --round 1 --by codex repo-map-mvp
parley consensus request-signoffs --dir . --participants claude,agy --yes repo-map-mvp
parley consensus signoff --dir . --agent codex --status accept --notes "Ready" repo-map-mvp
```

Open the TUI:

```bash
parley tui --dir .
```

For the full command and flag reference:

```bash
parley help
```

## Documentation

- [Agent runtime configuration](docs/agent-runtime-configuration.md)
- [CLI reference](docs/cli-reference.md)
