# CLI Reference

`parley` orchestrates Parley Deck workflows from the terminal. Commands operate
on the current directory by default; pass `--dir DIR` to target another
workspace root.

## Common Workflow

```bash
parley init --dir .
parley agents list --dir .
parley agents verify --dir . --agent claude
parley run --dir . --participants claude,agy --yes "Plan the next CLI slice"
parley status --dir .
parley continue --dir . <run-id-or-idea>
parley resume --dir . <run-id-or-idea>
```

## Commands

### `parley init [--dir DIR]`

Create the repository-local `parley-deck/` workspace with protocol, ideas,
inbox, metadata, and run directories.

- `--dir DIR`: workspace root. Defaults to `.`.

### `parley agents list [--dir DIR]`

Print the effective runtime matrix for configured agents. The output includes
installed state, version probe, launch mode, sandbox and approval policy,
model/profile, timeout, isolated-home use, backend type, and config sources.

- `--dir DIR`: workspace root. Defaults to `.`.

### `parley agents verify [--dir DIR] [--agent ID] [--full] [--yes]`

Verify configured agent CLIs.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--agent ID`: verify only one stable agent ID.
- `--full`: run behavioral headless probes, not just version probes.
- `--yes`: confirm probes that may call hosted backends.

Examples:

```bash
parley agents verify --dir . --agent codex
parley agents verify --dir . --full --agent hermes --yes
```


### `parley roster show [--scope deck|machine] [--dir DIR] [--all] [--json] [--explain AGENT]`

**This is the answer to "what is my roster?"** — the single canonical table. `agents list` is the
adapter/runtime inventory (what this machine *can* launch); `roster show` is the team (who this
deck *does* run). Do not parse `parley-deck/agents.toml` or the §2 table to answer the question:
reproduce this command's output.

The eleven columns are a versioned contract (`schema_version: 1`), identical in text and `--json`:

```
AGENT  ADAPTER  STATE  INSTALLED  MODEL  MODEL-FAMILY  MODEL-COMPANY  EFFORT  SPEED  AUTO  STATUS
```

`MODEL` and `EFFORT` hold the value the launch **actually passes**, or `unknown` — never a
declaration wearing the effective cell. Divergence and absence surface through `STATUS`:

| Status | Meaning |
| --- | --- |
| `ok` | nothing to report |
| `model-drift` | the configured model is not the one the launch passes |
| `model-unbound` / `effort-unknown` | the placeholder could not be bound; the flag is dropped |
| `metadata-unknown` | no `modelmeta` rule matched the model reference |
| `unmapped` | no adapter resolves this roster ID |
| `section2-only` | declared only in the legacy §2 table, which is no longer authoritative |
| `legacy-roster` | this deck has no `[roster.*]` block; §2 is the fallback |
| `inherited-roster` | this deck declares no roster; rows come from `~/.parley/agents.toml` |
| `inactive` | retired member, kept for history |
| `masked-by-env` | a higher config layer overrides this value |

Flags:

- `--scope deck` (default) — the roster this project declares. `--scope machine` — the user-global
  roster in `~/.parley/agents.toml`.
- `--all` — additionally list configured adapters that no roster declares. Use this when you
  installed an agent and it does not appear: the roster legitimately excludes it.
- `--explain AGENT` — per-field provenance: which config layer set each value.

### `parley roster set AGENT --scope deck|machine [--adapter A] [--model M] [--effort E] [--speed S] [--state active|inactive] [--dry-run] [--yes] [--confirm-breaking]`

Changes one member's settings. Preview is the default; nothing is written without `--yes`.

A **membership change** — adding a member, retiring one, or reviving one — additionally requires
`--confirm-breaking`, because it changes who deliberates and therefore a future idea's quorum.
Retiring sets `active = false`; members are marked, never deleted, so past ideas stay readable.

### `parley roster sync [--dir DIR] [--keep AGENT.FIELD]... [--dry-run] [--yes]`

Rebases the deck onto the machine defaults, one direction only (machine → deck); the machine file
is never written. It removes deck overrides so the deck inherits: redundant ones silently, and
**deliberate pins only after listing each one** with the exact `--keep` token that protects it. A
`--keep` token that matches no override is an error, not a no-op — a typo must not silently fail to
protect the field it names. Apply is bound to the preview: if the file changed in between, the
command refuses rather than discarding the edit.

### `parley roster render [--dir DIR] [--dry-run] [--yes] [--adopt-inherited]`

Regenerates the §2 roster table in `COOPERATION.md` from the authority. §2 is a **generated,
non-authoritative view**; do not hand-edit it. Generation is idempotent — rendering twice produces
byte-identical output. Rows present in §2 but absent from the roster are **reported** before they
are removed. A deck that declares no roster of its own will not have the inherited machine roster
written into its committed `COOPERATION.md` without `--adopt-inherited`.

### `parley roster migrate [--dir ROOT] --backup-dir DIR [--dry-run] [--yes --confirm-breaking] [--json]`

One-shot fleet migration: converts legacy §2-only decks under `ROOT` to `[roster.*]` blocks. Every
write is backed up first and validated after, with automatic rollback on failure; a deck whose
adapters cannot be resolved is skipped and reported rather than guessed. A deck with uncommitted
changes is skipped, so git history remains a usable rollback. `--yes` alone is refused: this
rewrites every deck under the root and is attended-only.

### `parley run [--dir DIR] [--no-tui] [--auto] [--participants IDS] [--yes] TASK`

Create a new idea from `TASK` and launch round-01 with selected participants.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--no-tui`: run headlessly and print results instead of opening the TUI.
- `--auto`: enable automatic low-risk HITL handling during the run.
- `--participants IDS`: comma-separated agent IDs, for example `claude,agy`.
- `--yes`: confirm selected hosted/non-local backend launches.
- `TASK`: free-form work request.

Example:

```bash
parley run --dir . --no-tui --participants claude,agy --yes "Review the repo-map MVP"
```

### `parley resume [--dir DIR] [--no-tui] RUN_OR_IDEA`

Resume a previous run by run ID or idea slug. This also validates pending
manual or interactive consensus signoff handoffs when possible.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--no-tui`: print run detail instead of opening the TUI.
- `RUN_OR_IDEA`: run ID under `parley-deck/runs/` or idea slug.

### `parley continue [--dir DIR] [--json] RUN_OR_IDEA`

Inspect a previous run by run ID or idea slug and print planner-derived next
actions for continuing the workflow. This first slice is read-only: it tells
the user which safe command to run next, such as answering a HITL question,
drafting consensus, requesting signoffs, finalizing, or retrying one missing
agent step in a later implementation slice.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--json`: print machine-readable JSON containing the run summary and actions.
- `RUN_OR_IDEA`: run ID under `parley-deck/runs/` or idea slug.

Example:

```bash
parley continue --dir . 20260517T120000.000000000Z
```

### `parley status [--dir DIR] [--run RUN_ID] [--idea SLUG] [--json]`

Show workspace, idea, consensus, run, and HITL question state.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--run RUN_ID`: inspect one run.
- `--idea SLUG`: inspect one idea or its latest run.
- `--json`: print machine-readable JSON.

### `parley answer [--dir DIR] RUN_ID QUESTION_ID ANSWER...`

Answer a pending human-in-the-loop question for a run.

- `--dir DIR`: workspace root. Defaults to `.`.
- `RUN_ID`: run directory ID.
- `QUESTION_ID`: question file ID.
- `ANSWER...`: answer text.

Example:

```bash
parley answer --dir . 20260517T120000.000000000Z q1 "Use the conservative default"
```

### `parley context repo-map [--dir DIR] [--format markdown|json] [--max-files N]`

Emit a deterministic map of the local repository. Markdown output is useful in
prompts; JSON output is useful for tools.

- `--dir DIR`: repository root to map. Defaults to `.`.
- `--format markdown|json`: output format. Defaults to `markdown`.
- `--max-files N`: maximum number of files to include. Defaults to `1000`.

Examples:

```bash
parley context repo-map --dir . --format markdown --max-files 50
parley context repo-map --dir . --format json --max-files 10
```

### `parley consensus status [--dir DIR] [--review] [--json] IDEA`

Inspect consensus state for an idea.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--review`: inspect `review/consensus.md` instead of `consensus.md`.
- `--json`: print machine-readable JSON.
- `IDEA`: idea slug under `parley-deck/ideas/`.

### `parley consensus draft [--dir DIR] [--review] [--round N] [--by AGENT] IDEA`

Draft a consensus file from submitted round or review files.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--review`: draft review consensus.
- `--round N`: source round number.
- `--by AGENT`: drafter agent ID recorded in frontmatter.
- `IDEA`: idea slug.

### `parley consensus signoff [--dir DIR] [--review] --agent ID --status STATUS [--notes TEXT] [--counter TEXT] IDEA`

Append one signoff block to the target consensus file.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--review`: sign review consensus.
- `--agent ID`: signing participant ID.
- `--status STATUS`: `accept`, `reserve`, `reservations`, or `block`.
- `--notes TEXT`: signoff notes.
- `--counter TEXT`: counter-proposal for `block` signoffs.
- `IDEA`: idea slug.

### `parley consensus request-signoffs [--dir DIR] [--review] [--participants IDS] [--yes] [--dry-run] IDEA`

Ask missing participants to append their own signoff blocks.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--review`: request review consensus signoffs.
- `--participants IDS`: comma-separated participant IDs to request.
- `--yes`: confirm hosted/non-local backend launches.
- `--dry-run`: print the planned requests without invoking agents.
- `IDEA`: idea slug.

### `parley consensus finalize [--dir DIR] [--by AGENT] IDEA`

Create or update `FINAL.md` from accepted consensus.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--by AGENT`: final plan author.
- `IDEA`: idea slug.

### `parley consensus reopen [--dir DIR] [--review] --reason TEXT IDEA`

Reopen a consensus file and preserve the previous blocked draft.

- `--dir DIR`: workspace root. Defaults to `.`.
- `--review`: reopen review consensus.
- `--reason TEXT`: reason recorded in the reopened file.
- `IDEA`: idea slug.

### `parley tui [--dir DIR]`

Open the project TUI.

- `--dir DIR`: workspace root. Defaults to `.`.

### `parley version [--dir DIR] [--all] [--json]`

Print the CLI version.

- `--dir DIR`: workspace root for project metadata in `--all` output.
  Defaults to `.`.
- `--all`: include `parley-deck-skill` installer, runtime skill, and project
  metadata status when available.
- `--json`: print machine-readable JSON.

## Exit Codes

- `0`: success.
- `1`: runtime failure, failed probe, malformed workspace, or agent failure.
- `2`: usage error or missing required argument.
- `3`: pending manual/interactive handoff for `consensus request-signoffs`.
