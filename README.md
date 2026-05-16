# parley-deck-cli

## Local install

Install the current checkout into `~/.parley-deck`:

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

## Project Notes

- [Agent runtime configuration](docs/agent-runtime-configuration.md)
