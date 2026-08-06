---
agent: codex-1
idea: roster-operations-standard
review-round: 1
date: 2026-08-06
reviewed-commit: 203f73bfd83cc8c55ea14aeca50a47edbeedee95
reviewed-skill-commit: b806adabdae349b80e83f3126312a66fceae1848
---

## Summary

Request changes for `parley-deck-cli` 1.40.1 and `parley-deck-skill` 2.5.1. I found 2
CRITICAL, 5 MAJOR, 2 MINOR, and no NIT issues. The two binding release boundaries that matter
most did not land: the run snapshot is persisted but never consumed, and TOML is not the roster
authority used by run selection. `roster sync --yes` should be disabled in the patch release until
the snapshot consumer exists; legacy/partial-roster mutations should likewise fail closed until
the authority migration is atomic.

## Scope and provenance

- `PRIMARY` — I read `parley-deck/COOPERATION.md` in full, including binding §15; `00-prompt.md`;
  `FINAL.md`; the normative field table, authority text, release conditions, and migration contract
  in `consensus.md`; the six named implementation commits; the release changelogs; and the affected
  CLI, protocol, tests, docs, and skill files in tag archives of `v1.40.0` and `v2.5.0`.
- `PRIMARY` — CLI release target: `v1.40.0^{}` =
  `203f73bfd83cc8c55ea14aeca50a47edbeedee95`. Skill release target: `v2.5.0^{}` =
  `b806adabdae349b80e83f3126312a66fceae1848`.
- `PRIMARY` — The CLI worktree had a pre-existing untracked
  `internal/app/roster_migrate.go`. It is absent from `v1.40.0`, so I excluded it and did not read
  it as shipped implementation. Additional post-tag roster implementation changes appeared in the
  shared worktree while this review was running; I likewise excluded them. I modified neither
  repository outside this review file.
- `PRIMARY` — No canonical
  `parley-deck/ideas/roster-operations-standard/IMPLEMENTATION.md` exists, and `00-prompt.md`
  remains `track: standard`, `status: round-01` while FINAL records `track: deliberation`. I used
  the user-supplied release commits and FINAL as the implementation boundary; the missing Phase-5
  handoff remains audit debt.
- `PRIMARY` — Negative searches below used `rg -uuu` or `find`, not an ignore-honouring negative
  grep. In particular, `rg -n -uuu 'Generate.*Roster|generate.*Roster|Render.*Roster|render.*Roster|generated.*view' internal --glob '*.go'`
  found no generator or caller; its matches were comments/tests about roster rewriting, not a
  generator implementation.
- `PRIMARY` — I built the release tag with `go build -o <tmp>/parley ./cmd/parley`. Targeted tests
  passed for `internal/agents`, `internal/app`, `internal/config`, `internal/runcontrol`,
  `internal/runmanifest`, and `internal/runstate`. `go test ./...` reached all packages but exited 1
  at `internal/runner.TestDurableKillEndToEndRealProcess` with
  `process verification failed (no recorded boot id); not killed`; I do not use that failure as a
  roster finding.
- `PRIMARY` — I exercised the release binary only against disposable tag archives and temporary
  decks. The commands and relevant output are quoted below. The temporary tree was removed after
  review.
- `SECONDARY` — No decisive finding relies on another participant's verdict.
- `RECALL` — No finding relies on memory or unsupported reasoning.

I own the findings and recommendations in this file and therefore issue no §15 verification
verdict on those claims. The D/G dispositions below are non-owner verdicts on the implementer-owned
release claims in `CHANGELOG.md:3-56` and the released artifacts; their decisive evidence is
`PRIMARY` and is quoted or located here.

## Refutation attempts

### D1-D9 against FINAL

| Decision | Release-claim verdict | Refutation result (`PRIMARY`) |
|---|---|---|
| D1 — three concepts, one answer | **WRONG as a complete claim** | `roster show` exists and is in top-level help, but `agents list` is not explicitly labelled “not the roster”, `rg -n -uuu 'roster' docs --glob '*.md'` returned no documentation, and the “run snapshot” is not what continuations use (CRITICAL-1). |
| D2 — effective MODEL/EFFORT or unknown | **WRONG as a release-wide claim** | Mapped rows derive from resolved argv (`launchargs.go:57-121`), but a legacy unmapped row is guessed through `proposeFamily` and prints that adapter's model/effort even though the roster ID cannot resolve for launch (`roster.go:270-345`; CRITICAL-2). |
| D3 — frozen identical 11-column contract | **WRONG** | Text emits `STATUS=ok`; JSON emits `"status": null`, includes out-of-contract `display_name`/`note`, and has no JSON golden. The only contract test pins the text header (`roster_test.go:212-225`) (MAJOR-2). |
| D4 — CLI-owned modelmeta | **CONFIRMED in inspected scope** | `modelmeta.go:27-108` peels `litellm`/`openrouter`, derives company from the model reference, and returns unknown on an unmatched reference; `modelmeta_test.go:8-70` covers the ratified cases and built-in defaults. |
| D5 — three complete verbs and safety | **WRONG** | `--all` and `--explain` are undefined; show ignores `--scope`; top-level help omits sync; `roster init --scope deck` is rejected; docs are absent; and `--yes` alone adds membership (MAJOR-2, CRITICAL-2). |
| D6 — immutable run snapshot | **WRONG** | The manifest writes snapshot fields, but continuation rediscovers configuration and `sessions inspect` never calculates `stale-snapshot` (CRITICAL-1). |
| D7 — model argv plus legacy normalizer | **WRONG as a complete claim** | Built-in placeholders reach the runner, but no legacy hardcoded-argv normalizer exists. An actual hardcoded Hermes override still defeats `roster set --model` (MAJOR-3). |
| D8 — skill/CLI boundary | **WRONG as a complete claim** | The new skill section correctly says to run and reproduce `parley roster show`, but bootstrap/protocol passages still instruct writing or filling §2 (MAJOR-4). |
| D9 — TOML authority and generated §2 | **WRONG** | Run/preset selection still parses §2, no generator exists, render-only fields are not rendered/migrated, and rows are sorted only by ID rather than active-first (CRITICAL-2, MAJOR-1). |

### G1-G5 binding gates

| Gate | Release-claim verdict | Refutation result (`PRIMARY`) |
|---|---|---|
| G1 — rebase gated on persisted **and consumed** snapshot | **WRONG** | Rebase shipped. `continueAuto` calls `discoverConfigured` and passes fresh `Agents`; no consumer or config-change continuation test exists (CRITICAL-1). |
| G2 — STATE wiring | **WRONG as an end-to-end gate** | `roster show` now consumes TOML/legacy inactive state, but default run and preset membership still consume §2, so `roster set --state inactive` is not authoritative for launch selection (CRITICAL-2). |
| G3 — atomic authority cutover | **WRONG** | The three main §2 paragraphs changed, but runtime readers remain, the generator is absent, CLI/docs are incomplete, and all three protocol/skill copies retain contradictory §2-authoring instructions (CRITICAL-2, MAJOR-1, MAJOR-2, MAJOR-4). |
| G4 — idempotent §2 generator | **WRONG** | There is no §2 generator and therefore no two-run byte-idempotence test (MAJOR-1). |
| G5 — §7-format changelog entry | **WRONG on format; present in substance** | The entry names the idea and one-off, but omits the mandated `Idea: ideas/.../`, `Drafted by:`, and `Summary:` fields (MINOR-2). |

## Findings

### [CRITICAL] G1 shipped rebase with a write-only snapshot

`PRIMARY` — `internal/app/app.go:1857-1875` computes a snapshot and passes it into run creation;
`internal/runmanifest/manifest.go:48-73,120-162` persists it. But
`internal/app/app.go:1147-1184` reconstructs a continuation by calling
`discoverConfigured(ctx, root)` and supplies those freshly resolved `Agents` to every later driver
operation. It never calls `runmanifest.Load`. The only app loads are session inspection/status at
`app.go:944,1004`; neither applies a snapshot. An `rg -n -uuu` scan of all internal Go files found
no other snapshot consumer.

`PRIMARY` — The stored row is insufficient even if it were loaded: it contains adapter, model,
effort, speed, `Auto bool`, and installed state, but not the resolved launch argv or the required
autonomous-write arguments (`manifest.go:65-73`). G1 requires the latter to remain unchanged.
`runcontrol.Create` also treats the manifest write as best-effort and permits a resumable run with
no durable snapshot (`runcontrol.go:76-102`).

`PRIMARY` — The required acceptance test does not exist. The `rg -n -uuu` test scan for
`RosterSnapshot|roster_snapshot|mid-run|config.*change|frozen` found only the ordinary manifest
load in `runcontrol_test.go:53`; that test does not pass or assert a snapshot. The non-test strings
`stale-snapshot` and `masked-by-env` have zero executable occurrences; `stale-snapshot` appears only
in a manifest comment.

Why it matters: after `roster sync` removes a deck pin, a continuation may change adapter, model,
effort, or permission-bearing argv mid-run. This is the exact live hazard G1 said must block the
release.

Suggested fix: fail closed immediately by withholding `roster sync --yes` in 1.40.1 until every
continuation/resume/retry path loads a participant-only frozen invocation plan. Store the exact
secret-free resolved argv and identity required to reconstruct each selected participant; make
snapshot persistence mandatory for resumable new runs (or duplicate it durably in the event log);
consume it in all later phases; implement `stale-snapshot`; and add the required acceptance test
that mutates machine and deck config between create and continue while asserting unchanged
adapter/model/effort/autonomous args.

### [CRITICAL] The TOML authority cutover can select or collapse the wrong quorum

`PRIMARY` — New-run authority did not move. Presets call `protocol.ReadRosterIDs` directly at
`internal/app/app.go:1793-1805`; default participants call it at `app.go:2418-2442`. The existing
`TestDefaultRosterParticipants` even creates `[roster.*]` TOML blocks and then changes §2 to drive
three expected results (`roster_test.go:41-84`); the relevant app test suite passed. Explicit
selection resolves adapter mappings but never rejects `active=false` (`app.go:2445-2466`).

`PRIMARY` — The cutover boundary is unsafe for legacy decks. `config.LoadRoster` merges machine,
deck, local, and env layers, and `resolveRoster` decides legacy solely from `len(entries)==0`
(`runtime.go:64-114`; `roster.go:227-258`). Thus any machine roster entry suppresses a deck's
legacy §2 fallback. Conversely, adding one deck block makes that one block the entire roster shown
by the TOML path. In a disposable legacy copy I ran:

```text
parley roster set hermes-1 --scope deck --model xai/grok-4.5 --yes
Wrote .../parley-deck/agents.toml
parley roster show --json
... exactly one roster row: hermes-1 ...
```

`PRIMARY` — Membership changes have no second confirmation. `rosterSet` writes any newly appended
`[roster.<id>]` block whenever `--yes` is present (`roster_set.go:21-66,126-139`). A disposable
new deck produced `+ adapter = "codex"` followed by `Wrote .../agents.toml` from a single
`roster set new-1 ... --yes`. That directly contradicts D5's “`--yes` alone is refused”.

`PRIMARY` — A deck with neither authority behaves inconsistently. The release binary's
`roster show` exits 1 with:

```text
no roster: declare [roster.<id>] in parley-deck/agents.toml
(or keep a legacy §2 table in COOPERATION.md)
```

but new-run defaulting treats an unreadable/empty §2 roster as “no roster” and falls through to
all installed families (`app.go:2425-2428,2445-2447`). That is not fail-closed membership and can
silently expand quorum.

Why it matters: `show`, preset validation, default launch, and explicit launch can disagree about
who is active. A one-agent edit can implicitly retire every legacy peer, while a no-roster deck can
launch every installed adapter. Both violate D9 and can defeat the protocol's non-solo invariant.

Suggested fix: create one resolved-roster API and use it for show, presets, preflight, default
selection, explicit selection, and snapshots. Distinguish a deck's committed authority from
machine defaults explicitly. On a legacy deck, refuse partial `set`/`sync` mutations until a full,
previewed migration materializes every member and retained field; require the second membership
confirmation. With neither TOML nor valid legacy §2, all launch paths must fail closed.

### [MAJOR] G4 has no §2 generator, and `roster sync` cannot migrate a legacy deck

`PRIMARY` — Neither `roster set` nor `roster sync` calls a protocol renderer; both end after
renaming `agents.toml` (`roster_set.go:55-66`; `roster_sync.go:121-141`). The full `rg -uuu` scan
quoted in Scope found no generator. After adding `new-1` in a disposable deck,
`rg -n -uuu 'new-1' <deck>/parley-deck` found only `agents.toml`; `COOPERATION.md` remained unchanged.

`PRIMARY` — `roster sync` iterates only IDs already present in the deck TOML
(`roster_sync.go:46-70`). Against the released CLI deck, which has a hand-written §2 table and no
deck `[roster.*]` blocks, the binary reported:

```text
roster sync: .../parley-deck/agents.toml already inherits from
~/.parley/agents.toml — nothing to do
```

So the documented operation does not “move it across”. There is no verbatim carry of
`workspace_dir`, `role`, or `host_handle`; no generated-view update; and no idempotence test.
`resolveRoster` also sorts all IDs together (`roster.go:259`) rather than active rows first and
inactive rows second as the normative contract requires.

Suggested fix: implement the bounded §2 renderer from the normative field table, preserving all
text outside the generated zone; order active then inactive, each byte-ascending by ID; call it in
the same atomic operation as every committed roster mutation; and add a test that renders twice
and compares every byte. Make `sync` materialize a complete legacy migration (or refuse and point
to a separate migration verb), never report a legacy no-op as migrated.

### [MAJOR] The frozen command and JSON contracts are incomplete or contradictory

`PRIMARY` — The release binary returned `flag provided but not defined: -all` and the same for
`-explain`. `runRoster` parses `--scope` but `show` calls `rosterShow(root, ...)` without it
(`roster.go:87-112`). In a disposable deck containing only deck member `new-1`,
`roster show --scope machine --json` still returned `new-1`. Top-level help lists show and set but
not sync (`app.go:119-125`), and `roster init --scope deck` returns
`invalid --scope "deck" (want session|machine)` because the deprecated path does not apply the
alias (`roster.go:138-139,408-410`). The supposedly hidden `session` spelling is still the visible
default in flag help (`roster.go:90`).

`PRIMARY` — `rg -n -uuu 'roster' docs --glob '*.md'` returned no hits, although D1 and the 1.40.0
changelog say docs landed. The JSON row adds `display_name` and `note` beyond the eleven columns
(`roster.go:165-181`); a healthy text row prints `ok`, while the observed JSON emitted
`"status": null`. No golden exercises text and JSON together; `roster_test.go:212-225` checks only
the text header.

Why it matters: 1.40.0 advertised this as a frozen API. Scripts cannot use two specified query
modes, scope is misleading, sync is undiscoverable in the main help, and JSON is already different
from the declared v1 shape.

Suggested fix: implement and test the exact D5 grammar (`show --scope/--all/--explain`, set, sync,
and a genuinely deprecated init alias); add all three verbs to top-level help and both named docs;
validate help through the binary; and publish one golden text/JSON contract with an explicit
representation for `ok`. Remove out-of-contract row fields or formally place them in the decided
explain/provenance object. If correcting the shipped JSON is considered breaking, bump and document
the schema rather than silently redefining v1.

### [MAJOR] D7's legacy normalizer is absent, so `roster set --model` can succeed without changing the launch

`PRIMARY` — `applyOverride` replaces `HeadlessArgs` wholesale at
`internal/config/runtime.go:622-625` and independently changes `spec.Model` at `:674-676`; no step
rewrites a hardcoded `--model <literal>` to `{model}`. `EffectiveModel` then reads the literal argv
(`launchargs.go:97-107`). This is the exact legacy shape D7 required a normalizer for.

`PRIMARY` — The released repository supplies a real example: its Hermes override hardcodes
`--model glm-5p2` in `parley-deck/agents.toml:61-71`. In a disposable copy,
`roster set hermes-1 --model xai/grok-4.5 --yes` reported success, but `roster show --json` still
reported `"model": "glm-5p2"` and `model-drift`. The display name used the newly configured Grok
value while the canonical model cell stayed GLM, recreating the split this idea was meant to end.

Suggested fix: implement the ratified legacy normalizer over known value-taking model flags,
preserving every unrelated permission/prompt token; bind the normalized placeholder from the
winning model field; and add fixtures for Claude, Codex, Kimi, Hermes, Agy, and OpenCode overrides,
including an autonomous flag next to the hardcoded model. `roster set` should warn/fail if the
requested value cannot reach effective argv.

### [MAJOR] G3/2.5.0 left contradictory §2-authoring instructions in every protocol surface

`PRIMARY` — The core new skill section at `skills/parley-deck/SKILL.md:254-295` correctly says to
run/reproduce `roster show` and never parse §2. But the same skill still says bootstrap choices are
recorded in local config “and the §2 roster” (`SKILL.md:171`). The live, embedded, and bundled
protocol copies likewise retain “mirrored in the §2 roster” at lines `57/56/56`, “Fill in §2
roster” at `1062/1053/1053`, and “Modify the active roster (§2)” at `1183/1174/1174`. These results
come from `rg -n -uuu` across all three copies located by `find`.

Why it matters: G3 required an atomic authority cutover and removal of every instruction that
treats §2 as a store. A facilitator following bootstrap or Appendix A can still hand-edit the
generated, explicitly non-authoritative view, while 2.5.0 claims that this guidance was removed.

Suggested fix: update all three protocol copies and `SKILL.md` so bootstrap and project adoption
write `parley-deck/agents.toml` through `roster set`/a full migration and then invoke the generator;
refer to “active roster” without calling §2 the authority; add drift assertions for these known
contradictory phrases; and do not claim the skill reproduces a generated view until CLI G4 exists.

### [MAJOR] Machine-scope writes use the wrong file whenever `PARLEY_HOME` is set

`PRIMARY` — The config loader defines `PARLEY_HOME` as the central config directory and reads
`$PARLEY_HOME/agents.toml` (`runtime.go:411-431`). `rosterScopeFile`, however, treats the same value
as a user home and writes `$PARLEY_HOME/.parley/agents.toml`
(`roster_set.go:76-90`). In a disposable environment, `parley init` reported central defaults at
`<home>/agents.toml`; `roster set machine-only --scope machine --yes` reported success at
`<home>/.parley/agents.toml`; `roster show --scope machine` did not contain `machine-only`. `find`
showed both files.

Why it matters: machine updates can report success while writing a file no resolver consumes.
Tests currently encode the wrong nested path in `roster_sync_test.go:11-24,79-89`, so they pass the
defect.

Suggested fix: make machine-scope operations call `config.CentralAgentsPath()` directly and change
the tests to seed/read that exact path. Add a round-trip test with `PARLEY_HOME` set: machine set →
machine show → deck inheritance.

### [MINOR] Pin preview/`--keep` is not sufficient against typos or concurrent edits

`PRIMARY` — The good half is present: preview is the default, every differing value is printed,
the exact keep token is shown, and the existing keep test passes (`roster_sync.go:98-123`;
`roster_sync_test.go:44-77`). That satisfies the minimum discoverability decision for a careful
two-command operator.

`PRIMARY` — The protection is fragile for a committed file. Keep tokens are lowercased into a map
but never validated or checked for unused entries (`roster_sync.go:52-55`), so a typo such as
`--keep kimi-1.modle --yes` still removes `kimi-1.model`. The command computes drops from one read
at `:46`, rereads the file at `:126`, and deletes the named fields without checking that their
values are unchanged at `:131-135`; an edit between those reads can therefore be lost despite the
atomic rename.

Suggested fix: reject every unmatched/unknown keep token; bind apply to the previewed file hash and
field old values; and consider requiring `--drop-pins` (or a preview plan ID) in addition to
`--yes` whenever `len(pins)>0`. Git makes recovery possible in a clean repository, but it is not a
substitute for validating operator intent before rewriting a committed file.

### [MINOR] G5's entry is not in the required §7 format

`PRIMARY` — `parley-deck/meta/protocol-changelog.md:117-139` contains the substantive date,
description, idea name, and user-authorized one-off. But §7 mandates four fields:
`## DATE — description`, `Idea: ideas/.../`, `Drafted by: ...`, and `Summary: ...`
(`COOPERATION.md:748-754`). The new entry instead uses bold `**Idea:**` without the path and has no
`Drafted by:` or `Summary:` line.

Suggested fix: reformat the entry exactly as §7 requires while retaining the useful one-off
explanation beneath it.

### [NIT] No findings

No NIT findings.

## Legacy fallback answer

No: it is only conditionally correct and is not safe as shipped.

- `PRIMARY` — With no `[roster.*]` in any config layer and a readable §2 table, `roster show`
  falls back and labels every row `legacy-roster`. The disposable release run did so, but also
  labelled every row `unmapped` and guessed an adapter/model for display.
- `PRIMARY` — With no deck block but any machine roster entry, `LoadRoster` is non-empty, so the
  code does **not** fall back to §2 and does not report `legacy-roster`; the machine roster wins.
- `PRIMARY` — With neither TOML roster nor parseable §2, `roster show` fails with exit 1, while a
  no-flag run falls back to all installed families. The latter is unsafe; it should fail closed.
- `PRIMARY` — `roster sync` does not create the missing deck entries and can report “already
  inherits” on a legacy deck. `roster set` can create one entry and thereby make only that partial
  entry visible on the TOML path. Neither is a safe migration.

## `roster sync` deliberate-pin answer

The default preview plus exact `--keep` output is a reasonable minimum for an ordinary, careful
two-step use, and it matches FINAL's chosen rebase semantics. It is not sufficient for unattended
`--yes`: unmatched keep tokens do not fail, there is no preview/apply hash binding, and every
differing pin is removed by default. MINOR-1 gives the concrete hardening needed; the attended
fleet migration still requires the much stronger inventory/CAS/backup/restore contract.

## Shipped extras and missing decided work

- `PRIMARY` — Shipped but not decided as part of the frozen row contract: JSON `display_name` and
  `note`. `DISPLAY-NAME` was explicitly removed; provenance was assigned to `--explain`/JSON, but
  the implemented extras have no specified schema and `--explain` itself is absent.
- `PRIMARY` — Decided but not shipped: G1 consumption/test, `stale-snapshot`, complete D5 flags and
  safety, the legacy argv normalizer, runtime authority cutover, the §2 generator/idempotence,
  normative field migration/rendering/order, complete protocol-copy cleanup, and exact G5 format.
- `PRIMARY` — FINAL Stage 4 migration tooling and the attended 40-deck operation are not in
  `v1.40.0`. FINAL deliberately sequences them after the coordinated Stage 1-3 release, so I do
  not classify their mere absence from 1.40.0 as a release defect. They remain outstanding and
  must not use the unreviewed partial semantics above. The untracked `roster_migrate.go` was not
  evaluated.
- `PRIMARY` — The missing `IMPLEMENTATION.md`, stale `00-prompt.md` phase/track, and release before
  Phase 6 leave no canonical implementation/deviation/validation record. The user's explicit
  direction explains why review is post-release; it does not supply the missing audit artifact.

## Open questions

1. Will 1.40.1 disable mutating `roster sync` until CRITICAL-1 is fixed, or will the patch include
   the complete snapshot consumer and continuation acceptance test?
2. Is the intended machine roster a source of deck **membership**, or only inherited field
   defaults for IDs declared by the deck? FINAL says both “deck authority” and “machine → deck
   inheritance”; the current layered merge silently chooses membership inheritance. The patch
   needs one explicit rule and tests before repairing CRITICAL-2.
3. Which artifact will own the Stage-4 inventory/CAS/backup/restore report, and how will its
   reviewed implementation be separated from the pre-existing untracked migration draft?
