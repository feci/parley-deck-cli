---
idea: roster-operations-standard
phase: 6 — code review
track: deliberation
reviewer: kimi-1
round: 1
date: 2026-08-06
baseline: parley-deck-cli v1.40.0 (203f73b) under review; v1.40.1 (58db960) tagged mid-review; parley-deck-skill v2.5.0 (b806ada)
---

# Review — roster-operations-standard (kimi-1, round 1)

## 0. Scope and §15 declaration

**What I reviewed (PRIMARY unless tagged otherwise):** `FINAL.md`, `consensus.md` (incl. the
revision-3 normative field table and the revision-2 migration contract), `COOPERATION.md` §15;
CLI diff `v1.39.0..v1.40.0` and, after the baseline shifted, `v1.40.0..v1.40.1`; skill diff
`v2.4.0..v2.5.0`. I read every file the brief names, ran `go build ./...`, ran
`go test ./internal/...` (green at v1.40.0 *and* at v1.40.1), built the binary to `/tmp` and
exercised `roster show|set|sync|render`, `--help`, `agents list` against scratch decks under
`/tmp` with `PARLEY_HOME` isolation. Negative-evidence searches used `rg -uuu` (stated per
finding). **I did not write to either repository.** All scratch state was under `/tmp` and has
been deleted; `git status` is clean in both repos. The real `~/.parley/agents.toml` was only
read, never written.

**Baseline shift, disclosed.** While I was reviewing, the working tree advanced: uncommitted
`roster migrate` WIP appeared (~19:41 local), then `de269c4` (fix-up cycle 1) and the
`v1.40.1` tag landed (19:56). Every finding below was gathered independently from v1.40.0
source and binary behavior **before** I saw the 1.40.1 CHANGELOG entry; I did not read
hermes-1's review file (present in this directory since 19:45) before writing this one.
Because 1.40.1 exists, each finding carries a **status at v1.40.1**: FIXED or OPEN. OPEN
findings are the 1.40.2 / 2.5.1 candidates.

**Claims I own and do not verdict (§15.1).** From this idea's rounds/signoffs I own: the
1.39-era measurements credited to kimi-1 (roster dispatched but unlisted in `printUsage` and
absent from docs; `DISPLAY-NAME` contradicts `MODEL` via `RenderDisplayName`/`ModelLabel`),
and the *positions* later adopted as R4 (`--keep` + pin enumeration) and my round-2 adoptions
(rebase, snapshot, post-write `masked-by-env` re-resolve). I issue no verdict on those
historical claims. Where a *current-state* fact is adjacent (e.g. "sync is absent from
`--help` **today**"), it is a new claim about new code, established here by fresh PRIMARY
evidence; I flag the adjacency and note codex-1/hermes-1 have independent verdicts per the
1.40.1 fix attributions. My R4 was a design position; the implementation facts (what shipped
sync actually does) first exist now and I verdict them PRIMARY.

**Machine-state note.** The stage-1 commit message claims the four `headless_args` overrides
in `~/.parley/agents.toml` "were removed". `PRIMARY` — I read the file: no active
`headless_args` key remains (4 comment mentions only); the interim workaround is gone, as
FINAL required. Stale comments describing the removed overrides survive (operational nit,
outside repo scope).

---

## 1. CRITICAL

Both CRITICALs are **FIXED in v1.40.1** — they are recorded here because the review target is
the released v1.40.0, where both binding gates were violated, and because their residuals
(§2, findings M1/M2/M3) are still open.

### C1 — G1 violated in v1.40.0: the snapshot was persisted and never consumed; the acceptance test did not exist — FIXED in 1.40.1, two residuals (→ M1, M2)

`PRIMARY`. At v1.40.0:

- `rg -uuu -n '\.RosterSnapshot\b' --type go` matched exactly two reads — both struct-literal
  *writes* (`internal/runcontrol/runcontrol.go:88`, `internal/runmanifest/manifest.go:161`).
  Nothing read the stored snapshot back.
- `continueAuto` still re-discovered: `internal/app/app.go:1152` (v1.40.0) called
  `discoverConfigured(ctx, root)` and passed `Agents: discovered` (app.go:1163), the exact
  hazard G1 names ("a live hazard, not a hypothetical").
- `rg -uuu -ln 'RosterSnapshot|roster_snapshot|RosterRevision' --glob '*_test.go'` matched
  **zero files**: the gate's mandated acceptance test (create run → change config → continue
  → prove adapter/model/effort/auto-args unchanged) did not exist.
- The stage-3 commit (7e03a74) message says "Only then does sync ship, per the ratified
  gate" — the implementer believed persistence satisfied the gate. It did not: G1 demands
  persist **and consume** plus the test. Rebase shipped without the half that made it safe.

**Status at v1.40.1: FIXED in substance.** `applyRosterSnapshot`
(`internal/app/roster_snapshot_apply.go:22`) pins model/effort/speed from the frozen row;
`continueAuto` consumes it (`app.go:1157-1160`, `runmanifest.Load` + apply, stderr warning);
`TestContinuationUsesTheFrozenSnapshotNotTheChangedConfig`
(`internal/app/roster_snapshot_apply_test.go:16`) proves config change → frozen values win;
suite green. Residuals filed as M1 (auto-args) and M2 (per-adapter keying) — both OPEN.

### C2 — G3/G4 violated in v1.40.0: no §2 generator existed, and run membership still parsed §2 as its sole authority — FIXED in 1.40.1, one residual (→ M3)

`PRIMARY`. At v1.40.0:

- `rg -uuu` for any generator (`RenderRoster|GenerateRoster|renderSection2|renderRoster`)
  found no function and no caller. The only §2-related writer was the *uncommitted*
  `roster_migrate.go`, which writes TOML **from** §2 (migration direction) — nothing wrote
  §2 **from** TOML (the generator direction). Yet commit 3b94f85 had already changed all
  three protocol copies to call §2 "a generated, non-authoritative view". The protocol
  described a generator that did not exist — G4 ("generated §2 is idempotent") was
  unsatisfiable, and any `roster set`/`sync` instantly made the (still runtime-parsed) §2
  stale with no regeneration path.
- Runtime authority cutover was partial: `resolveRoster` read TOML-first (roster.go:232-258,
  v1.40.0), but participant defaulting (`defaultRosterParticipants`, app.go:2425) and preset
  validation (app.go:1797) both read **only** `protocol.ReadRosterIDs` — the §2 table. On a
  TOML-authority deck with no/stale §2, runs fell back to installed families or failed
  closed while `roster show` answered from TOML: the table and the run could disagree about
  who was in the roster — the two-sources-of-truth defect this idea exists to end.

**Status at v1.40.1: FIXED.** `parley roster render` (`internal/app/roster_render.go:83`)
generates §2; `RosterMembership` (roster.go:643) is the single membership authority and both
selection paths now use it (app.go diff in de269c4). Idempotency verified behaviorally
(§3, G4). Residual: the normative field table's §2-only-ID conflict rule is still
unimplemented and render silently erases such rows — M3, OPEN.

---

## 2. MAJOR

### M1 — G1 residual: snapshot consumption does not pin AUTO / autonomous-write args — OPEN

`PRIMARY` (read of `internal/app/roster_snapshot_apply.go:40-48`): the frozen entry's `Auto`
field is **never read**; `headless_args`, `approval_policy` and the rest of the launch shape
come from fresh discovery. A machine-config change that drops `--yolo`/`--dangerously-skip-
permissions` (or a `roster sync` that removes a deck pin carrying one) changes the
continuation's autonomy posture — the launch is then *not* what the run froze, and the row's
own AUTO answer stops being true. G1's acceptance demand names "adapter/model/effort/
**auto-args** unchanged"; the mechanism pins three of four and the test
(`roster_snapshot_apply_test.go:16-35`) proves only those three. The gate is not fully
discharged. (Adapter pinning is by construction — matching is per adapter family — subject
to M2.)

### M2 — G1 residual: frozen rows are keyed by adapter, collapsing per-roster-ID pins — OPEN

`PRIMARY` (`roster_snapshot_apply.go:26-29`): `frozen[e.Adapter] = e` — the map key is the
adapter family, so two roster IDs that share an adapter with *different* frozen models
overwrite each other; both continuations then launch the last entry's model. Per-ID pinning
is a supported configuration by the implementation's own contract (roster.go comment:
"Per-roster-ID settings beat the adapter-family default, so two roster IDs can share an
adapter and still run different models"), and `RosterSnapshot` stores entries per agent ID
— the consumer throws the distinction away. Key the map by `Agent` and match against the
roster-ID resolution, or the freeze is wrong exactly where it should be strongest. Edge-case
in today's single-ID-per-adapter decks, structural the day a second ID shares one.

### M3 — Normative field-table conflict rule unimplemented: §2-only IDs are silently dropped, and `roster render` erases them without report — OPEN

`PRIMARY` (behavioral, v1.40.1 binary, scratch deck, since deleted): deck with TOML roster
(claude-1, kimi-1, antigravity-1) **and** a §2 table containing a fourth ID `ghost-1`:

```
$ parley roster show --dir deckC        # → no ghost-1 row, no `unmapped` status
ghost-1: ABSENT from show (silently dropped)
$ parley roster render --dir deckC      # preview lists only the 4 config IDs
$ parley roster render --dir deckC --yes && parley roster render --dir deckC --yes
§2 already matches the roster — nothing to do   (second run: BYTE-IDENTICAL — G4 itself holds)
```

`ghost-1`'s row — including its render-only workspace-dir/role cells — vanished from §2 with
no preview mention, no report, no skip. The field table (`consensus.md:355`) decides
otherwise: "TOML wins; **a §2-only ID is reported `unmapped`, never auto-added**". The
verbatim-carry rule (:369-372) exists precisely so render-only project data is not lost.
Note the inconsistent safety postures inside the same release: `roster migrate` *skips and
reports* an unresolvable ID as `unclean`; `roster render` *erases* it silently. Any deck in
the measured fleet state (17 naming retired `antigravity-1`, 3 `gemini-1`, 1 `agy-1`, one
deck missing `hermes-1`) hits this the first time render runs on it. Severity MAJOR: a
ratified absence/conflict rule is unimplemented and the new generator makes the loss
permanent and silent.

### M4 — Undecided membership layering leaks the machine roster into deck §2 generation and the default run quorum — OPEN

`PRIMARY`. `RosterMembership` (roster.go:643-656) and `renderRosterTable`
(`roster_render.go:30`) both use `config.LoadRoster` — the *layered* machine+deck view —
while `rosterSync` deliberately operates on the deck *file* only
(`roster_sync.go:46-50`). Consequences, verified on the real deck (dry-run only, nothing
written):

```
$ parley roster render --dir .        # this repo's own deck, whose [roster.*] has 4 members
… table includes `opencode-1` | – | participant | active   ← from ~/.parley/agents.toml
```

A deck with no deck-file roster **silently inherits the whole machine roster** (verified at
v1.40.0: deck with neither §2 nor `[roster.*]` but a rostered machine printed the machine
roster, exit 0, no `legacy-roster`, no marker of any kind); render then *commits* that
inheritance into `COOPERATION.md`, where it goes stale the next machine change — the drift
vector this idea exists to kill, re-created in committed files; and since 1.40.1 routes
participant defaulting through the same layered view, the machine roster also becomes the
deck's default **quorum** (this repo's deck would default-run `opencode-1`, a fifth
participant nobody put on the deck). FINAL/D9 decided "agents.toml is the deck authority"
but never decided whether deck membership means *deck-file* or *layered*; the release
implements both in different verbs. This needs an explicit decision; my position: deck
membership is the deck file, machine seeds *values* (per sync's rebase model), and render/
membership/quorum should use the deck view with inheritance made explicit in output, not
baked into committed tables.

### M5 — D5 surface gaps: `show` has no `--all`/`--explain`, and `--scope` is advertised but silently ignored — OPEN

`PRIMARY` (v1.40.1 binary):

```
$ parley roster show --dir deckC --all
flag provided but not defined: -all
$ parley roster show --dir deckC --explain claude-1
flag provided but not defined: -explain
$ parley roster show --dir deckC --scope machine     # prints the DECK roster, exit 0, no warning
```

D5 freezes `roster show [--scope deck|machine] [--all] [--json] [--explain AGENT]`. None of
the three works; `--all` was the decided answer to the `opencode` invisibility problem
(consensus.md:112-113), and `--explain` is where D3 parked per-field provenance (the
`SOURCE` substitute). Worst is `--scope`: `printUsage` and `rosterUsage` both advertise it
for `show`, it parses, and it changes nothing — a silently wrong answer is worse than a
rejected flag. (Same state at v1.40.0.)

### M6 — Discoverability decided in D1/D5 is still unmet: `sync` (and now `render`, `migrate`) absent from `--help`; docs contain zero roster mentions — OPEN

`PRIMARY` (v1.40.1): `--help` lists only `roster show` and `roster set`
(`internal/app/app.go:123-124`); `rg -uuu -n 'roster sync' internal/app/app.go` → no hit;
`rg -uuu -c 'roster' docs/cli-reference.md docs/agent-runtime-configuration.md` → **zero
matches in both**. D1: "`roster` must appear in `parley --help` **and the docs**." D5:
"three verbs, all in `--help`". The verb a legacy-deck operator most needs (`sync`) — and
now the generator (`render`) and fleet tool (`migrate`) — are undiscoverable; the
documentation debt the idea measured is reproduced for its own fix. *Ownership adjacency
noted per §0: the 1.39-era version of this measurement was mine; this finding is fresh
evidence about 1.40.x and codex-1/hermes-1 verdict independently.*

### M7 — Membership-gate bypass: a new member added without `--adapter` writes with `--yes` alone — OPEN (introduced by the 1.40.1 fix)

`PRIMARY` (v1.40.1 binary, scratch deck):

```
$ parley roster set another-9 --scope deck --dir deckC --adapter kimi --yes
roster set: this adds a new roster member … Re-run with --confirm-breaking as well as --yes.   exit=2   ✓ gate works
$ parley roster set sneaky-9  --scope deck --dir deckC --model k3 --yes
Wrote …/agents.toml                                                                            exit=0   ✗ bypassed
```

`membershipChange` (`roster_set.go:236-246`) proxies "new member" as "+ adapter = " in the
change list; a block created with only `--model/--effort/--speed` (or `--state`) slips the
second confirmation D5 mandates for **membership changes**. The result is a real member
(`LoadRoster` returns it) that renders `unmapped` and cannot launch — visible, so impact is
limited, but the ratified gate's intent is defeated by shape, not by operator intent. Gate
on "block did not exist before", not on which key was written. (The v1.40.0 state was
worse — `--yes` alone added members, D5 violated outright — and is covered by the fixed
half of this finding.)

### M8 — D6 half-shipped: `sessions inspect` never reports `stale-snapshot` — OPEN

`PRIMARY`: `rg -uuu -n 'stale-snapshot|StaleSnapshot' --type go` → only an aspirational
comment in `runmanifest/manifest.go:56` and unrelated review-snapshot sweep code. Nothing
compares `RosterRevision` against the live deck; `inspectSession`
(`internal/app/app.go:865` area) embeds the manifest but computes nothing. D6: "`sessions
inspect` reports `stale-snapshot`" — the audit-trail half of the snapshot story — is absent
in both 1.40.0 and 1.40.1, and the frozen STATUS vocabulary (D3) therefore ships a code no
surface can emit (see also M9/N1).

### M9 — D7's decided legacy normalizer never shipped — OPEN

`PRIMARY`: D7's adopted mechanism is "{model}/{effort} placeholders … **plus a legacy
normalizer** for configs that hardcode a model literal in `headless_args`"
(`consensus.md:125-127`). The placeholder half shipped and is verified (below, "What
lands"). The normalizer half does not exist: `applyOverride`
(`internal/config/runtime.go:622-625`) still replaces `HeadlessArgs` wholesale with no
literal-model detection or rewrite (`rg -uuu` for any such logic → none), so a deck/machine
override carrying `--model <literal>` silently outranks the configured `model` field —
surfaced as `model-drift` (good) but never normalized (decided). The one known instance was
removed by hand (§0), which handles the user, not the mechanism.

### M10 — `roster migrate` (shipped in 1.40.1 as Stage-4 tooling) deviates from the ratified migration contract — OPEN, and now is exactly when to fix it

`PRIMARY` (read of `internal/app/roster_migrate.go` against `consensus.md:444-468`):

- **No compare-and-swap** between dry-run and apply — the contract requires skipping a deck
  whose source roster or target files changed since the dry-run; the tool re-reads live at
  apply time with no frozen comparison.
- **Bulk `--yes` across all decks** — the contract requires per-deck or small-batch
  confirmation honoring `roster_change_policy = "confirm-breaking"`; `rosterMigrate` applies
  every non-skipped deck on one flag and never reads the policy.
- **Foreign-deck compatibility gate absent** — decks whose protocol/schema version predates
  the change must be "skipped and reported" (:386-388); no version is read at all.
- **Inventory is thin** — no frozen source roster revision, no protocol/schema version, no
  worktree state (a dirty-tree deck is migrated in place); the JSON report has
  before/after hashes and backup paths only.
- **Verbatim-carry incomplete** — only `adapter`/`active` are written; `workspace_dir`,
  `role`, `host_handle` are never carried into the committed TOML keys the field table
  defines (:361-363). Render's legacy-cell carry prevents *visible* loss, but the authority
  file never gains the data.
- The tool also changes *meaning*, not only storage — rows whose adapter is not in the
  machine roster become `active = false` (`roster_migrate.go:200-217`, self-acknowledged in
  its comment as user-authorized). Plausibly inside the user's "zosúladím s aktuálnym
  rosterom"; flagging so the authorization is consciously scoped before the fleet run.

Stage 4's fleet operation is a *separate attended operation* that has not run — so unlike
the two CRITICALs, this one can still be fixed **before** it bites 40 repositories.

---

## 3. MINOR

- **N-min-1 — `masked-by-env` is in the frozen STATUS vocabulary (D3) but nothing emits it.**
  `PRIMARY`: `rg -uuu -n 'masked-by-env' --type go` → zero hits. `roster set` performs no
  post-write re-resolution, so a deck pin masked by a higher layer
  (`$PARLEY_HEADLESS_AGENT_CONFIG`, `agents.local.toml`) reports success and the user never
  learns the effective row did not change. (Post-write re-resolve was my own round-2 adopted
  *position* — noted, not verdicted as such; the vocabulary membership is D3 text.) — OPEN
- **N-min-2 — `agents list` was never relabelled** "adapter/runtime inventory — not the
  roster" (D1). `PRIMARY`: header still `AGENT INSTALLED VERSION LAUNCH …` (binary, HEAD);
  `rg -uuu -n 'inventory|not the roster'` in the agents-list path → zero. — OPEN
- **N-min-3 — legacy-deck remediation is self-contradictory across surfaces.** `PRIMARY`
  (v1.40.0 binary): `roster sync --dir <legacy-deck>` → "already inherits … nothing to do"
  (no `[roster.*]` → zero redundant overrides) — yet consensus decision 9 names
  `roster sync --dry-run` "the documented remediation", and skill 2.5.0's SKILL.md tells
  agents a legacy deck keeps working "until `parley roster sync` moves it across". Sync
  moves nothing across; the real paths are `init` (a *deprecated* alias per D5) and, as of
  1.40.1, `migrate`/`render` — neither named in the skill. Skill text is factually wrong
  about released behavior → 2.5.1. — OPEN
- **N-min-4 — continuation silently proceeds unfrozen when the manifest is unreadable.**
  `PRIMARY` (`app.go:1158-1160`, v1.40.1): `if m, mErr := runmanifest.Load(...); mErr == nil`
  — on load error the run continues with zero protection and no warning; pre-snapshot runs
  are the *intended* no-op case, but a corrupt `run.json` deserves a stderr line. — OPEN
- **N-min-5 — the G1 acceptance test is a unit test of the pinning function, not the
  gate's literal end-to-end shape** (create run → mutate config → continue → prove
  unchanged). The wired call site gives the property (code-read), but the gate mandated the
  run-level proof; the current test would not catch a future unwiring of
  `continueAuto`. — OPEN
- **N-min-6 — file-mode regression on roster writes.** `PRIMARY` (observed):
  `writeRosterFileAtomic` (roster_set.go:200-218) renames a `0600` temp file over a `0644`
  target — the machine file came back `-rw-------` after one `roster set`. Harmless for
  single-user, wrong on shared/team setups; `fsutil.WriteFileAtomic` with explicit mode
  exists and is used elsewhere. — OPEN
- **N-min-7 — `roster sync` edits a committed file with no backup or cleanliness check.**
  Per the adopted R4 design the preview + enumeration + `--keep` protection is otherwise
  sufficient (see §5), and the apply run *does* enumerate every removed pin before writing
  (verified); git history is the only rollback. A one-line "commit or stash first" hint when
  the deck tree is dirty would close it. — OPEN, MINOR by the decision's own lights
- **N-min-8 — stale guidance strings.** `unmapped` rows still advise "run `parley roster
  init`" (roster.go) though init is the deprecated alias and a TOML-authority deck needs
  `roster set`; `roster init` itself still demands the pre-1.40 `--scope session|machine`
  spellings while set/sync speak `deck|machine`, and emits no deprecation warning. — OPEN

## 4. NIT

- **NIT-1 — `modelmeta` prefix rule `{"k", …}` precedes `{"kimi", …}`** (`modelmeta.go:65-66`):
  any id starting with "k" is "Kimi K"; the `kimi` rule is unreachable for such ids, and an
  unrelated model starting with k would misclassify. Ordering/breadth, one-line fix. — OPEN
- **NIT-2 — agy's built-in `Model` is the display label** `"Gemini 3.6 Flash (High)"`
  (`discover.go:257`, pre-existing): the frozen contract's MODEL cell — and the launch argv
  itself — now carry a human label with spaces and parens. Made *visible* by D2's honesty,
  not caused by it. — OPEN
- **NIT-3 — legacy `(inactive)` role-suffixes are carried verbatim into the generated
  Role column** alongside the new State column (observed: `participant (inactive)` |
  `inactive`), double-marking retired rows forever. Compliant with verbatim-carry, cosmetically
  odd. — OPEN
- **NIT-4 — `membershipChange` message can name the wrong direction**: a reactivation
  (`- active = false` → `+ active = true`) matches the `"active = false"` substring first
  and reports "retires a roster member". Still gated; wording only. — OPEN
- **NIT-5 — `RosterSnapshot` drops `Display`/provenance**: the snapshot could carry the
  derived display name for artifact rendering; it re-derives instead. Not decided either
  way; noting for completeness. — OPEN

---

## 5. The brief's specific questions

**Does each of D1-D9 land as decided?**

| D | Verdict | Basis |
|---|---|---|
| D1 — three concepts, one answer | **Partial.** `roster show` is the answer and is in `--help`; the `agents list` relabel and **all doc updates** are missing (N-min-2, M6). | PRIMARY (binary, rg) |
| D2 — effective-or-`unknown` | **Lands.** Verified: configured-but-not-launched model → `model-drift`; unbound → `unknown` + `model-unbound`/`effort-unknown`; cells never carry a bare declaration. | PRIMARY (scratch runs) |
| D3 — frozen 11-column contract | **Lands**, minus two unemittable vocabulary codes (`masked-by-env`, `stale-snapshot` — M8, N-min-1). Header pinned by golden test (roster_test.go:223); `--json` carries `schema_version: 1` + ordered `columns` (verified). | PRIMARY |
| D4 — CLI-owned `modelmeta` | **Lands.** Gateway peel (`litellm/xai/grok-4.5`→xAI/LiteLLM) and `glm-5p2`→Zhipu AI verified live (hermes-1 scratch row); golden tests exist; no adapter-inference. NIT-1 noted. | PRIMARY |
| D5 — three verbs + safety | **Partial.** Verbs exist and preview-default/`--yes`/atomic writes hold (tested); but `--all`/`--explain`/working `--scope` absent (M5), sync absent from `--help` (M6), and the membership second-confirmation was absent in 1.40.0, fixed with a bypass in 1.40.1 (M7). `--scope deck` writes the committed file, never `agents.local.toml` ✓; `--scope session` hidden alias ✓. | PRIMARY |
| D6 — immutable run snapshot | **Persisted ✓, consumed in 1.40.1 with residuals (M1/M2); `stale-snapshot` reporting absent (M8).** | PRIMARY |
| D7 — model-argv fix + normalizer | **Placeholders land** (built-ins carry `{model}`/`{effort}`, runner substitutes at runner.go:1101; `codex`/`kimi` `-m` verified in `agents list` output; unbound flag-drop tested). **Normalizer absent** (M9). | PRIMARY |
| D8 — skill/CLI boundary | **Lands** in skill 2.5.0: SKILL.md names `roster show` as THE answer, forbids parsing §2/TOML/`agents list`, documents the same 11-column contract (not a second format). One factual error about sync (N-min-3). | PRIMARY (skill diff b806ada) |
| D9 — §2 generated view, TOML authority | **Substantially lands as of 1.40.1** (render + `RosterMembership` + three protocol copies + drift guard green). Residuals: §2-only-ID rule (M3), layering leak (M4). | PRIMARY |

**G1 (rebase gated on snapshot).** At v1.40.0: **violated** — persist without consume, no
acceptance test, `continueAuto` re-discovered (C1). The implementer's belief ("Only then
does sync ship") covered only the writing half. At v1.40.1: consumption and a test exist;
the AUTO-args clause and per-ID keying remain open (M1, M2) and the test is not the
gate-shaped end-to-end proof (N-min-5).

**G2 (STATE wiring).** **Satisfied as of 1.40.1.** The `_` discard is gone; `resolveRoster`
renders STATE/inactive correctly (verified: `antigravity-1` row STATE=inactive, status
`inactive`); `RosterMembership` feeds inactive-aware exclusion to both participant-selection
paths (app.go:2431-2433 logic preserved). In 1.40.0 the show-half worked but selection was
§2-bound (part of C2).

**G3 (atomic authority cutover).** v1.40.0: **violated** — protocol copies changed without
a generator and with selection still on §2 (C2). The three-copy drift guard is green at both
tags and the embedded/skill §2 sections are byte-identical (`shasum` match); live differs
only in the allowlisted roster zone. v1.40.1 completes the code side; M3/M4 are the
remaining contract gaps.

**G4 (idempotent §2 generator).** v1.40.0: **no generator existed at all** (C2). v1.40.1:
`roster render` exists and I verified two consecutive applies are **byte-identical** with a
clean no-op second run, deterministic ordering (active-first, then byte-ascending — the
normative rule), prose and host-handle table untouched, legacy cells carried. Gate met,
modulo M3's silent-erasure caveat.

**G5 (protocol-changelog entry).** **Satisfied.** `parley-deck/meta/protocol-changelog.md:119-139`
names the idea, the track, the §7 one-off, and the no-precedent wording, in §7 format.

**The legacy fallback.** Correct and safe in both pure modes (PRIMARY, behavioral): a deck
with no `[roster.*]` *anywhere in the stack* renders §2 rows flagged `legacy-roster`
(including §2-only IDs, marked `unmapped`, inactive detected); a deck with NEITHER (and a
rosterless machine) exits 1 with a clear "no roster: declare `[roster.<id>]` …" — no silent
empty answer. The unsafe middle is not the fallback but the *mixed* state: once any machine
or deck layer has `[roster.*]`, §2 is ignored entirely and §2-only IDs vanish without the
decided `unmapped` report (M3), and a rosterless deck on a rostered machine silently adopts
the machine roster with no marker (M4).

**`roster sync` removing deliberate pins.** The preview/`--keep` protection is **sufficient
as decided** (PRIMARY, behavioral): preview is the default; pins are enumerated with the
exact `--keep AGENT.FIELD` spelling; the apply run prints the same enumeration before
writing; writes are atomic; membership survives (test suite + my runs); the machine file is
never written. This matches the adopted R4 design — noting §0 that R4 was my position; the
*implementation facts* are what I verdict. Residuals: no backup/cleanliness check on a
committed-file edit (N-min-7), and "nothing to do" on legacy decks conflicts with the
documented-remediation story (N-min-3).

**Shipped that FINAL did not decide / decided and not shipped.** Nothing shipped outside
the decisions' envelope — `render` and `migrate` are the decided G4 generator and Stage-4
tooling (migrate shipped with contract deviations, M10; the deck's own self-migration
committed in de269c4 was the attended operation beginning, outside code scope). Decided but
not shipped (all OPEN at v1.40.1): doc updates and the `agents list` relabel (D1);
`--all`/`--explain`/working `--scope` and sync-in-help (D5); auto-args pinning (G1);
`stale-snapshot` reporting (D6); the legacy normalizer (D7); `masked-by-env` emission (D3);
the §2-only-ID conflict rule and full verbatim-carry (field table); compare-and-swap,
per-deck confirmation, foreign-deck gate, full inventory (migration contract).

---

## 6. What lands correctly (verified, no findings)

- `{model}`/`{effort}` end-to-end: built-ins carry placeholders (discover.go:202,225,250,298,
  329,361), the runner substitutes (runner.go:1101), unbound placeholders drop the
  introducing flag (launchargs.go:86-91, tested), `AUTO` computed from resolved argv
  (fail-closed, discover.go:138-144). The interim machine overrides are gone (§0).
- The frozen contract renders identically in text and JSON with `schema_version` and ordered
  `columns`; golden tests pin it; `DISPLAY-NAME` is out of the table and derived per-row.
- `modelmeta`: gateway peel + producer namespace + prefix rules behave as decided on the
  worked cases; `unknown` + `metadata-unknown` on no match (kimi-1 unmapped scratch row).
- `roster set`: preview-default, idempotent no-op, comment-preserving line-based patch,
  atomic write, deck scope = committed file, inactive marks (never deletes), and as of
  1.40.1 `--confirm-breaking` gates the adapter-shaped membership paths.
- `roster sync`: rebase semantics, one direction, pin enumeration + `--keep`, machine file
  untouched — each covered by its own test.
- `roster render`: byte-idempotent, deterministic order, prose-preserving, legacy-cell
  carry (G4 itself holds).
- Three protocol copies + changelog + drift guard: consistent and green (G3 text side, G5).
- Full suite `go test ./internal/...` green at v1.40.0 and at v1.40.1 (ran both).

## 7. Summary for the facilitator

Two CRITICALs against v1.40.0 — both binding-gate violations (G1 consume-half; G3/G4
generator/selection authority) — independently corroborate the findings 1.40.1 already
fixes; both fixes verified with two open residuals each (M1/M2, M3). Ten MAJORs remain open
at v1.40.1, of which M3 (silent §2-only-ID erasure), M4 (undecided membership layering
leaking into committed §2 and the default quorum), M7 (membership-gate bypass), and M10
(migration-contract deviations **before** the attended fleet run) are the ones I would not
let reach 1.40.2's ship line unaddressed. MINORs: 8. NITs: 5. No severity is empty.
