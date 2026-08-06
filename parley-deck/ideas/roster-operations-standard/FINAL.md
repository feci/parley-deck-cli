---
idea: roster-operations-standard
drafter: claude-1
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: deliberation
rounds: 2
signoff-revisions: 3
date: 2026-08-06
status: final
---

# FINAL — standardize roster operations across the CLI and the skill

## Problem

"What is the current agent roster?" had no single answer. Three CLI surfaces produced three
different tables, two independent stores held membership, and the table reported a `MODEL` the
launcher never passed. Adding `opencode` exposed all of it at once.

Measured (`PRIMARY`): **nine distinct §2 rosters across 40 decks**; 17 with no roster at all; 17
still naming retired `antigravity-1`; 3 naming `gemini-1`; 1 naming `agy-1`; one deck with no
`hermes-1` row — which is the user's "hermes works in some sessions and not others", exactly.

And the values were wrong: `claude` declared `claude-opus-5[1m]` while launching
`claude-opus-4-8[1m]`; **six of seven adapters passed no effort flag at all**. Root cause:
`applyOverride` sets `spec.Model` (`internal/config/runtime.go:594`) and never touches
`HeadlessArgs`, while `buildAgentInvocation` substitutes only `{root}` and `{prompt}`
(`internal/runner/runner.go:1097-1108`).

## Decisions

Full text and provenance in `consensus.md`. Summary of what is binding:

**D1 — three concepts, one answer.** `parley agents list` = adapter/runtime inventory (relabelled).
`parley roster show` = **the** roster answer. Run snapshot = what a run actually used.
`roster` must appear in `parley --help` and the docs; today it is dispatched but unlisted.

**D2 — `MODEL` and `EFFORT` are effective-or-`unknown`.** Never a declaration in the effective
cell. Divergence surfaces via `STATUS`.

**D3 — frozen 11-column contract v1**, identical in text and `--json`, with `schema_version`,
ordered `columns`, golden tests, additive-only evolution:

```
AGENT  ADAPTER  STATE  INSTALLED  MODEL  MODEL-FAMILY  MODEL-COMPANY  EFFORT  SPEED  AUTO  STATUS
```

`DISPLAY-NAME` leaves the table (it contradicts `MODEL` today). `SOURCE` and `ROUTE` are excluded;
provenance lives in `--explain AGENT` and JSON. `STATUS` vocabulary: `ok`, `unmapped`,
`not-installed`, `model-drift`, `model-unbound`, `effort-unknown`, `metadata-unknown`,
`masked-by-env`, `legacy-roster`, `inactive`, `stale-snapshot`.

**D4 — `modelmeta` is CLI-owned.** Versioned, tested registry; peel gateway prefixes before
deriving company; never infer company from the adapter; `unknown` + `metadata-unknown` on no match.
No deck hand-writes family/company.

**D5 — three verbs, all in `--help`:**
```
parley roster show  [--scope deck|machine] [--all] [--json] [--explain AGENT]
parley roster set   AGENT --scope deck|machine [--adapter A] [--state active|inactive]
                    [--model M] [--effort E] [--speed S] [--dry-run] [--yes]
parley roster sync  [--dir DIR] [--keep AGENT.FIELD]... [--dry-run] [--yes]   # machine -> deck only
```
Preview by default; `--yes` applies; `--yes` alone is **refused** for membership changes.
`--scope deck` writes the **committed** `parley-deck/agents.toml`, never the gitignored
`agents.local.toml`. `roster init` becomes a deprecated alias; `--scope session` a hidden alias.

**D6 — session = immutable run snapshot**, not a third scope. `sessions inspect` reports
`stale-snapshot`.

**D7 — the model-argv fix is in scope.** `{model}`/`{effort}` placeholders in built-in
`HeadlessArgs`, substituted on the existing `runner.go:1101-1103` path, plus a **legacy normalizer**
for configs that hardcode a model literal in `headless_args`. `codex` and `kimi` both accept
`-m/--model` (verified), so injection applies to them too.

**D8 — skill/CLI boundary.** The skill invokes `parley roster show` and reproduces its output. It
never parses §2, TOML, or `agents list`, and documents no second format.

**D9 — §2 becomes a generated, non-authoritative view;** `parley-deck/agents.toml` is the deck
authority. Full normative field table in `consensus.md`. Only the agent ID and the `inactive`
marker are runtime-semantic today; `workspace_dir`, `role`, `host_handle` are render-only and are
carried across verbatim. Ordering: active before inactive, then agent ID byte-ascending.

## Binding release gates

These are gates, not sequencing advice.

**G1 — rebase must not ship before the snapshot.** The change that exposes rebase MUST also persist
**and consume** the immutable effective row, with an acceptance test that creates a run, changes
machine/deck config, continues the run, and proves adapter/model/effort/auto-args unchanged. Today
`continueAuto` re-discovers config (`internal/app/app.go:1148-1160`), so this is a live hazard, not
a hypothetical. Filed independently by claude-1, hermes-1 (R1) and codex-1.

**G2 — `STATE` wiring is a prerequisite for the migration.** `resolveRoster` discards the inactive
map into `_` (`internal/app/roster.go:110`); the parser puts every row into `active` including
inactive ones. Without the wiring, marking 17 decks' retired rows inactive is a no-op that reports
success.

**G3 — the authority cutover is atomic**: committed-TOML schema, §2 generator, removal of runtime
§2 parsing, live protocol, embedded copy (`internal/protocol/defaults/`), skill snapshot
(`skills/parley-deck/references/`), skill behaviour, CLI help and docs land together or stay
feature-gated together. Three COOPERATION.md copies, per the standing drift guard.

**G4 — generated §2 is idempotent.** Two runs of the generator produce byte-identical output.

**G5 — a `meta/protocol-changelog.md` entry** in §7 format names this idea and the user-authorized
one-off.

## Implementation stages

codex-1's staging, endorsed by hermes-1 and kimi-1, adopted.

**Stage 1 — foundations (may land behind the existing surface).** `{model}`/`{effort}` resolver +
legacy normalizer; `modelmeta` registry + golden tests; resolved-row types; the 11-column contract
and JSON schema + golden tests; `STATE` consumption. The public effective-value contract is not
exposed until the resolver and `STATE` are wired.

**Stage 2 — authority cutover + ordinary operations (atomic, G3).** Committed-TOML schema; migrate
every §2 field; `roster show` / `set`; idempotent §2 generator; remove runtime §2 parsing; live +
embedded + bundled protocol; skill behaviour; CLI help; docs.

**Stage 3 — snapshot + rebase (atomic, G1).** Persist and consume the immutable effective row and
`roster_revision`; continuation test; only then expose rebase in `roster sync`.

**Stages 1-3 ship as ONE coordinated CLI + protocol + skill release.**

**Stage 4 — migration tooling, then the attended fleet operation.** Inventory, dry-run report,
compare-and-swap, file-level backups with verified restore, per-deck rollback, resumability,
compatibility/skip gates, final machine-readable report. The 40-deck mutation is a **separate
attended operation** after that release: the frozen dry-run goes to the user, only approved decks
or small batches apply, and it is never folded into the code merge.

## Recorded dissent and process notes

**codex-1 blocked twice.** Revision 1: mis-cited §7 exception, mis-classified track, ungated rebase,
incomplete §2 spec. Revision 2: the §2 field contract was stated as a requirement and not supplied.
Both blocks were upheld in full. **Neither revision survived review, and both failed the same way** —
the drafter wrote what must be true instead of making it true.

**Fifteen drafter position changes** are recorded in `consensus.md`; four came from the rounds,
eleven were forced at signoff.

**User directions** (2026-08-06) closed three questions evidence could not: rebase over additive-pin,
the §2 protocol change inside this idea rather than a separate meta idea, and authorization of the
fleet migration. The §7 deviation is an explicit **one-off** and sets no precedent.

## Interim workaround now in force

`~/.parley/agents.toml` carries a full `headless_args` override pinning
`--model claude-opus-5[1m]`, because the config `model` field could not reach the process. Verified:
launch argv now carries opus-5, `AUTO` stays `yes`, `{root}` expands per deck. **This override MUST
be removed when D7 lands** — a wholesale `headless_args` override is exactly how `hermes` silently
lost `--yolo`. The removal is part of Stage 1's legacy normalizer work.
