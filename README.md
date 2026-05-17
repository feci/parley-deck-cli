# parley-deck-cli

`parley` is the local CLI for Parley Deck multi-agent cooperation. It creates
and maintains the `parley-deck/` workspace, discovers local agent CLIs, starts
or resumes rounds, tracks human-in-the-loop questions, manages consensus
signoffs, and emits repository context for agents.

## Install

Install the current checkout into `~/.parley-deck` while developing:

```bash
cd /Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-cli
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
parley run --dir . --participants claude,gemini --yes "Plan the next CLI slice"
```

Run without opening the live TUI:

```bash
parley run --dir . --no-tui --participants claude,gemini --yes "Review the repo-map MVP"
```

Resume or inspect work:

```bash
parley status --dir .
parley status --dir . --idea repo-map-mvp
parley resume --dir . 20260517T120000.000000000Z
```

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
parley consensus request-signoffs --dir . --participants claude,gemini --yes repo-map-mvp
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
