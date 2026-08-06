---
idea: roster-operations-standard
phase: 8 — re-review
agent: kimi-1
round: 3
date: 2026-08-06
reviewed-commit: 7220715
verdict: FINDINGS
---
# Re-review round 3 — kimi-1

## Verdict

**FINDINGS** — one new MINOR (N-1, a narrow freeze hole introduced by the A5 fix) and one
round-02 MINOR carried open (skill 2.5.1 still unshipped; the edits are still working-tree-only).
Every other round-02 finding — the CRITICAL included — is verified **fixed**, each proven
behaviorally at 7220715, not inferred from the diff. `PRIMARY` unless tagged otherwise.

**§15.1 ownership disclosure.** A1's resolution adopted my round-01 M4 position; A12's
post-write re-resolve was my adopted position; the legacy-fallback condition was codex-1's. I
issue no verdict on those positions. What I verdict below are fresh implementation facts about
code that first exists at 7220715, gathered this session. codex-1 and hermes-1 verdict
independently.

**Method.** Binary built from 7220715 to `/tmp/parley-rev3`; scratch decks under `/tmp/kr3` with
`PARLEY_HOME` isolation; nothing written to either repository (`git status --short` in the CLI
repo shows only the untracked `review/round-03/` directory this file lands in). Where the CLI
surface cannot observe an internal path (participant selection, the continuation boundary), I
injected a probe test via `go test -overlay` from `/tmp` — the reviewed tree stays
byte-identical; the probe drives the same production functions (`defaultRosterParticipants`,
`applyRosterSnapshotToParticipants` + `agents.ResolveParticipant`, `rosterSnapshotState`,
`rosterFieldMaskedBy`).

**Mandated command, run by me at 7220715.** `go build ./...` exit 0; `go test ./...` — all 26
packages ok, 0 failures. Because the first run was partly cached I re-ran the whole suite
uncached: `go test -count=1 ./...` — all 26 packages ok, 0 failures (`internal/app` 22.4s,
`internal/runner` 9.7s — `TestDurableKillEndToEndRealProcess`, which codex-1's round-02 review
saw fail in his sandbox, **passes on this machine**). `go vet ./...` clean. The cycle-3 handoff
note about that test is accurate for this environment.

## Round-02 findings: fixed or not

**[CRITICAL, all three] A1 legacy-fallback clarification — FIXED, proven for all three
authority cases on `roster show` AND participant selection.** `PRIMARY` (behavioral). Machine
home with 5 members (claude-1, codex-1, hermes-1, kimi-1, opencode-1):

- *Committed deck block* (2 members): exactly 2 rows; the undeclared model still inherits as a
  value; no machine-only member appears.
- *Legacy §2 only* (§2 = claude-1, kimi-1 active + hermes-1 `(inactive)`; no config roster):
  exactly the 3 §2 rows, each `legacy-roster`, hermes-1 correctly `STATE=inactive`; codex-1 and
  opencode-1 (machine-only) absent; machine values still seed (`claude-opus-5`, `glm-5p2`).
- *Neither*: 5 rows, every one `inherited-roster`; `roster render --yes` refuses exit 1 and
  `--adopt-inherited` permits (re-verified).

Participant selection, which the binary cannot print without launching, was driven through the
overlay probe: `RosterMembership("/tmp/kr3/deck-legacy")` → `active={claude-1,kimi-1}`,
`inactive={hermes-1}`; `defaultRosterParticipants` → `[claude-1 kimi-1]` — the run quorum is the
§2 set, not the machine's five. The authority order lives where the fix claims it:
`LoadRosterScoped` (internal/config/runtime.go:164-196) decides deck blocks → valid legacy §2
(`Legacy: true`, runtime.go:176-192) → machine (`Inherited: true`), before any value layering;
`resolveRoster` marks `legacy-roster` from `scope.Legacy` (internal/app/roster.go:277,307-309).
`roster render` on a legacy deck no longer mis-refuses as "no roster of its own" — it
regenerates §2 from the legacy membership, preserving roles and the inactive mark (verified).

**[CRITICAL, codex-1] membership pooled every non-machine layer — FIXED.** `PRIMARY`
(behavioral + overlay). Only `parley-deck/agents.toml` carries `membership: true`
(runtime.go:306-312,365); `agents.local.toml` and `$PARLEY_HEADLESS_AGENT_CONFIG` are skipped as
membership sources (runtime.go:142-156). Deck declaring 2 members + `agents.local.toml` adding
`sneaky-local` + env config adding `sneaky-env`: `roster show` still shows exactly the committed
2 (the local layer's *value* override for claude-1 does apply — values layer, membership does
not), and the overlay probe confirms the run quorum on that deck is `[claude-1 kimi-1]`.

**[MAJOR] A3 `roster init --yes` gate — FIXED.** `PRIMARY` (behavioral). On a legacy deck:
`roster init --yes` exits 2 with "this is a membership change — re-run with --confirm-breaking
as well as --yes" and writes nothing; `--yes --confirm-breaking` writes the 3 blocks, exit 0;
the `--json` path returns `outcome: "needs-breaking-confirmation"` with exit 2. Gate at
internal/app/roster.go:554-558,604-607; `runRoster` now passes the flag through (roster.go:167).

**[MAJOR] A5 per-ID pins at the real continuation call site — FIXED at the boundary; one new
edge, see N-1.** `PRIMARY` (overlay through the runner's own resolver). `continueAuto` now calls
`applyRosterSnapshotToParticipants(run.Participants, discovered, rosterMappingFor(root), …)`
(internal/app/app.go:1184), which resolves each participant first, then freezes
(internal/app/roster_snapshot_apply.go:83-99). My probe fed it the production input shape — one
adapter-keyed discovery `claude` (model "drifted"), participants `claude-1`/`claude-2`, snapshot
pinning opus-a/`--yolo` and opus-b/`--safe` — then resolved each participant with
`agents.ResolveParticipant` exactly as the runner's `selectedAgents` does
(internal/runner/runner.go:364-372). Result: claude-1 → `opus-a`,`[--yolo]`; claude-2 →
`opus-b`,`[--safe]`. The frozen clones carry `Spec.ID = <roster-id>`, so the runner's rule-1
exact-ID match finds them (Discovery embeds Spec, internal/agents/discover.go:181-188). The
cycle-2 defect shape is dead.

**[MAJOR] A6 machine scope read deck values/provenance — FIXED.** `PRIMARY` (behavioral). Deck
sets claude-1 `model = "deck-only-model"`; machine file has `machine-write-2`. `--scope machine`
shows `machine-write-2`; `--scope machine --explain claude-1` attributes adapter/model to
`~/.parley/agents.toml` and effort/speed to `built-in default` — no deck layer anywhere. Deck
scope still shows `deck-only-model`. The scoped loaders exist as claimed
(`LoadAgentSpecsScoped`, `LoadRosterAdaptersScoped`, `RosterFieldSourcesScoped`,
runtime.go:1032-1109; consumed at roster.go:249,265 and roster_view.go:165).

**[MAJOR] A7 `display_name`/`note` outside the frozen 11 columns — FIXED.** `PRIMARY`
(behavioral). Both fields are now `json:"-"` (roster.go:205-214). A healthy row's JSON carries
exactly the 11 declared keys (`agent, adapter, state, installed, model, model_family,
model_company, effort, speed, autonomous, status`); the text table is unchanged; `scope` and
`schema_version` present. `Note` still surfaces in the text table/`--explain` (the `⚠` lines),
which is where the consensus allowed operator guidance to live.

**[MAJOR] A9 drift guard 2-of-4 surfaces — FIXED, negative control proven.** `PRIMARY`. The
extended `TestNoSection2AsAStoreInstructions` (internal/protocol/drift_test.go:261-292) now also
reads the bundled skill protocol and SKILL.md via the sibling checkout, fails loudly per
surface, and bans `roster sync` moves it over/across`. Run verbosely here: PASS with **no**
"skipping" log lines — i.e. all four surfaces were actually read (the sibling checkout exists at
`../../../parley-deck-skill`). Negative control: I overlaid the test with sibling paths pointing
at poisoned `/tmp` copies — it FAILS and names both surface and phrase ("bundled skill protocol
… 'roster sync` moves it over'"). Honest caveat (code-read): on a machine without the sibling
checkout the guard degrades to 2 surfaces with only a `t.Logf`, not a failure — structural for a
cross-repo guard, and visible in test output.

**[MAJOR, hermes-1] stale "roster sync moves it over" in all three COOPERATION.md copies —
FIXED in content on all four surfaces.** `PRIMARY`. `grep` for all six banned phrases (the four
originals + both new sync claims) across live `parley-deck/COOPERATION.md`, embedded
`internal/protocol/defaults/COOPERATION.md`, the skill repo's bundled `references/COOPERATION.md`
and `SKILL.md`: zero hits. The replacement text names the real remediation (`migrate`, or `set …
--confirm-breaking` then `render`) and explicitly says sync does not migrate. The skill bundled
copy is byte-identical to the CLI embedded copy below the header (diff-verified). **The
release-hygiene half of my F-4 / codex-1's A14 is NOT done**: in the skill repo, HEAD is still
`b806ada`, both files are still ` M` (uncommitted working-tree edits), and `package.json` still
reads 2.5.0 — skill 2.5.1 does not exist. See "New findings", carried-open item.

**[MAJOR] A10 `RosterRevisionOf` omitted LaunchArgs — FIXED, proven through the audit
surface.** `PRIMARY` (overlay). The hash now mixes `strings.Join(e.LaunchArgs, "\x1f")`
(internal/runmanifest/manifest.go:88-99). Probe: a fabricated manifest whose revision matches the
deck reports `current`; the same snapshot with one entry's argv extended (`--output-format=json`
appended, `Auto` unchanged) reports `stale-snapshot` — exactly the drift round-02 proved
invisible.

**[MAJOR/MINOR] A12 `masked-by-env` false-positive on machine-scope writes — FIXED, both
directions.** `PRIMARY` (behavioral). `rosterFieldMaskedBy` now resolves the source label via
`config.RosterSourcePath` and compares absolute paths (roster_set.go:96-118;
runtime.go:1013-1021). Machine-scope write with no higher-layer override: silent, write
effective. Machine-scope write genuinely masked by `agents.local.toml`: correct warning naming
the actual masker. Deck-scope genuine mask: still warns (round-02 behavior preserved).

**[MINOR] A16 residuals — all three FIXED.** `PRIMARY`. (a) Unmapped guidance now reads `parley
roster set <id> --scope deck --adapter <family> --yes --confirm-breaking` (roster.go:324) —
executable under the A3 gate. (b) modelmeta: the broad `{"k",…}` prefix rule is gone;
`kimiCodename` (k+digit, modelmeta.go:130-136) classifies `k2` → Kimi K/Moonshot, while
`kronos-2` → unknown/unknown (no longer misclassified) and `kimi-k2-thinking` → Kimi/Moonshot
(behavioral, all three). (c) CHANGELOG.md:7 now credits "codex-1, hermes-1 and kimi-1".

**[NIT] IMPLEMENTATION.md handoff accuracy — the accuracy half FIXED; the format half
persists.** `PRIMARY`. The cycle-2 section now reads "Ten new tests" with a coverage list that
matches exactly those ten (the eleventh, the drift assertion, is named in A9's own bullet), and
the runner-failure caveat is stated. I withdraw the count objection: under the sentence's own
scoping, ten is correct (my round-02 NIT-3 counted the drift test against a list that excluded
A9). Still open, codex-1's round-02 NIT part (a): the file still opens with two blank lines and
no Phase-5 status/head-commit metadata (lines 1-3 are blank/blank/`## Fix-up cycle 2`). Record
hygiene only.

**hermes-1's round-02 MINOR (wholesale `HeadlessArgs` replace) — unchanged, and I do not
re-raise it.** `PRIMARY` (code-read): roster_snapshot_apply.go:66-68 still replaces the argv
wholesale. hermes-1's own framing was advisory ("correct behavior for the common case … the
risk is theoretical"); cycle 3 neither claimed nor needed to change it.

## New findings (by severity, or "none")

### N-1 — [MINOR] NEW in cycle 3: a roster ID spelled exactly like its adapter family silently loses its freeze on continuation

`PRIMARY` (behavioral, overlay probe). `applyRosterSnapshotToParticipants` *appends* the frozen
per-participant clones to the discovery list, leaving the original adapter-keyed discovery in
place (roster_snapshot_apply.go:88-97). The runner's resolver returns the **first** exact-ID
match (internal/agents/resolve.go:47-53). For a canonical suffixed ID (`claude-1`) the original
adapter entry (`claude`) can never match, so the frozen clone wins — verified. But for a deck
whose roster ID **is** the bare family name (`[roster.claude]`, legal TOML and accepted by the
§2 parser's `rosterRowRe`), the original unfrozen discovery sorts first under the same ID and
shadows the clone. Probe output, participant `claude` with snapshot entry `Agent:"claude",
Model:"frozen-model"`: the post-boundary resolution returns model `"drifted-machine-model"` —
the current config, not the freeze — with **no warning** (the clone matched the snapshot, so the
not-in-snapshot warning never fires). Pre-cycle-3 this ID shape was pinned correctly (the old
call path hit `frozen["claude"]` directly on the adapter discovery — `git show 57fe9d7` diff,
PRIMARY code-read), so this is a genuine, if narrow, regression introduced by the A5 fix: the
one ID shape the new resolve-then-freeze order cannot see is also the one where the failure is
silent, and the freeze is the autonomy-posture protection. Narrow trigger (non-canonical but
legal ID; `--participants claude` alone does not hit it — the *roster* must declare the ID for
the snapshot to contain it), so MINOR not MAJOR. The fix is small: replace the adapter discovery
in place (or prepend clones) instead of appending, or have rule 1 prefer a snapshot-matched
entry.

### N-2 — [MINOR, carried open from round 02] Skill 2.5.1 still does not exist; the corrected files are still uncommitted

`PRIMARY`. My round-02 F-4(a) and codex-1's A14 MINOR, re-checked: skill repo HEAD `b806ada`,
`git status --short` → ` M skills/parley-deck/SKILL.md`,
` M skills/parley-deck/references/COOPERATION.md`, `package.json` 2.5.0, no 2.5.1 changelog
entry. The *content* is verified correct on all four surfaces (above) and is now drift-guarded,
so this is purely the release step — and to be even-handed, the CLI's own 1.40.2 does not exist
yet either (CHANGELOG.md's latest entry is 1.40.1); shipping was always the post-Phase-8 step.
I record it rather than close it because the round-02 consensus's agreed text was "ship as
1.40.2 / skill 2.5.1": if Phase 8 closes here, the uncommitted skill edits must be committed and
both versions cut **as part of that close**, or the content this cycle verified never reaches an
installed skill. Not chargeable to commit 7220715's code; chargeable to the exit checklist.

### NITs

- **NIT-1** `PRIMARY` (code-read): the A10 hash-format change silently flips every pre-fix
  manifest to `stale-snapshot` on first inspect under the new binary (old revisions were
  computed without LaunchArgs). Fail-safe direction and self-healing on the next run, but worth
  one line in the release notes so operators do not chase phantom drift across the fleet's
  in-flight runs.
- **NIT-2** (persisting, my round-02 NIT-2) `PRIMARY`: the set warning still says "(status
  `masked-by-env`; …)" and docs/cli-reference.md:80 lists the code in the STATUS vocabulary, but
  no `roster show` row can emit it — it remains warning-only vocabulary.
- **NIT-3** (persisting half of codex-1's round-02 NIT) `PRIMARY`: IMPLEMENTATION.md still lacks
  the Phase-5 status/head-commit header (see the NIT verification above).

## Test-quality assessment

`PRIMARY`. I read all six new/extended tests, ran them individually (all PASS), and checked each
against its revert. **No tautologies.**

- `TestLegacySection2BeatsTheMachineRoster` — REAL. Fixture is the exact CRITICAL shape (deck
  with no config roster + 5-member machine + valid 2-row §2; `deckWith` confirmed to put the
  second argument in `PARLEY_HOME`). Asserts both the membership set (via `RosterMembership`)
  AND the row marking (no `inherited-roster`, all `legacy-roster`). Reverting the
  `LoadRosterScoped` legacy branch fails the first half; reverting `resolveRoster`'s
  `scope.Legacy` fails the second.
- `TestSnapshotPinsSurviveParticipantResolution` — REAL, and it uses the production input shape
  (adapter-keyed discovery + roster-ID participants + mapping) that the cycle-2 test got wrong.
  Reverting the helper to adapter-keyed freezing fails it (the returned list would contain only
  `claude`). One named gap: the test stops at the helper — the one-line call site at
  app.go:1184 is still unguarded, so reverting *the wiring* (not the helper) passes the suite.
  I closed that gap for this review with the overlay probe through `agents.ResolveParticipant`;
  DF-2's full create-run→mutate→continue shape remains deferred by consensus.
- `TestRosterRevisionCoversLaunchArgs` — REAL. Exactly the defect: same fields, different argv,
  revision must differ. Revert fails it.
- `TestMachineScopeWriteIsNotReportedAsMasked` — REAL for the negative direction (revert to
  label comparison → the machine write reports itself masked → FAIL). The positive direction (a
  genuine higher-layer mask must still warn) has no shipped test; I verified it behaviorally
  (CLI + overlay). Small gap, worth filling the next time the file is touched.
- `TestRosterInitRequiresConfirmBreaking` — REAL and end-to-end through `runRoster`: exit 2,
  no file written, then the confirmed path exits 0. A gate revert fails the first assertion.
- Extended drift guard — REAL, negative control run by me (poisoned sibling paths via overlay →
  FAIL naming surface and phrase). The skip-when-sibling-absent degradation is logged, not
  silent.

**Regression sweep (did cycle 3 break anything).** Full uncached suite green; vet clean; the
three authority cases, render refusal/escape, render-on-legacy (regenerates from the §2
membership, preserving roles/inactive), init JSON outcome strings, sync-on-legacy ("already
inherits … nothing to do", matching the new protocol text), and machine-scope explain all
behave. One new narrow defect found: N-1. Nothing else.
