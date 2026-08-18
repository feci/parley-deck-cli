---
agent: kimi-1
idea: addon-bidding-scope
round: 1
date: 2026-08-17
---

## Summary

**Recommendation: FLIP THE DEFAULT — but only with grandfathering, not as a bare default change.**
`parley-bidding` stays in the package, shipped, tested, and `doctor`-verified exactly as today; what
changes is that a *fresh* flag-less install stops placing it, while every *existing* install that
recorded it keeps it across upgrades. The naive flip (just change the default list) is the one
version of this idea that measurably breaks existing users: it turns every routine upgrade into a
red `doctor` fleet-wide, per the installer's own health semantics. Details and migration mechanics
below.

Why not KEEP: the package's own README already records that default-on is a cost imposed on people
who never asked (`README.md:31-37`), and the workaround it offers is not sticky — every flag-less
`install --force` re-selects all discovered add-ons (`lib/installer.js:892-901`, `:968-989`) and
rewrites the recorded selection (`:1958-1963`), so the opt-out must be re-declared on every update,
forever. A permanent apology paragraph for one's own default is a defaulted decision asking to be
made explicit. Why not CUT or SPLIT: my widened searches (below) still find no usage — but absence
of evidence cuts against removal too, the artifact is real (built from live tender work at BYTE,
not speculation), and the orphan mechanics of SPLIT/CUT are strictly worse than FLIP's.

## Proposed approach

Concretely, three small changes, all in `parley-deck-skill` (no protocol change, no new
dependency):

1. **Fresh installs default to core + four process add-ons.** `selectedAddons`
   (`lib/installer.js:892-901`) currently returns all discovered add-ons when no flag is given.
   The default set excludes `parley-bidding`; `--only parley-bidding` (or an additive `--with`)
   opts in. The skill remains packaged, discovered, manifest-validated, and covered by `npm test`
   (54 Python tests, `CHANGELOG.md` 2.1.0 entry) — nothing about its quality gate changes.
2. **Grandfather recorded selections.** For `install` with no flags, consult the core marker's
   recorded `addons` list before falling back to the package default — the read path already
   exists (`markerAddonNames`, `:930-960`; `expectedAddonNames` already does this for non-install
   commands, `:978-987`, and deliberately excludes install today). Effect: the fleet that got
   bidding by default keeps it, `doctor` stays green, nobody's runtime changes underneath them;
   only new installs see the new default. The damaged-marker path already fails closed
   (`:943-956`), so this adds no new failure mode.
3. **Make single add-on removal expressible.** Today there is no command that removes one add-on
   and keeps the rest: `uninstall --no-addons` removes only the core (test,
   `test/installer.test.js:593-604`), and `uninstall --only parley-bidding` plans the core unit
   too — `targetSkillUnits` puts the core at `units[0]` unconditionally (`:995-1016`) and
   `removeFleetAtomically` removes every planned unit (`:1712-1745`). (Code reading; no test
   covers `--only` uninstall — that gap is worth a test either way.) A user who wants bidding
   gone today must `rm -rf` the directory and then fix the marker via `--only` with the remaining
   four names. The flip should ship with either per-addon uninstall or a documented equivalent.

**Who breaks, honestly:** no workflow breaks — nothing invokes the skill automatically; it is
model-selected per task. What breaks under a *naive* flip is health signaling: an existing default
install that upgrades flag-lessly keeps the on-disk `parley-bidding/` (install never removes
unselected add-ons — `:1086-1094` comment) while the rewritten marker excludes it, so `doctor`
reports `valid-unselected` and exits non-zero on every runtime (`:2124-2131`, `:2702-2706`).
Re-including it via `--only parley-bidding` then drops the other four from the marker
(`writeMarker` records exactly the selected set, `:1958-1963`) and they go `valid-unselected`
instead — the only clean re-include is `--only` with all five names. Grandfathering (item 2)
eliminates this entirely: the marker still records bidding on upgrades, so nothing flips red.
Under SPLIT/CUT it is worse: after any flag-less upgrade the new package no longer discovers the
add-on, so it vanishes from `doctor`'s on-disk sweep (`:1103-1118` iterates `discovered` only),
`--only parley-bidding` is rejected as unknown (`:877-887`), and uninstall can reach it only
while the old marker still records it (the `authorize` second clause exists precisely for this,
`:1026-1043`) — after the marker rewrite, no shipped command can remove it; the fix is manual
`rm -rf` of each `<skills-dir>/parley-bidding`. The per-skill marker left behind
(`.parley-deck-skill-install.json`, e.g. on this machine: `markerSchema: 2`, `version: 2.8.0`,
manifest aggregate + sha256 — `~/.kimi-code/skills/parley-bidding/`) then describes an orphan no
tool manages.

**Cost of being wrong.** If default-on was silently load-bearing — a population of bidding users
who never show up in this workspace — grandfathering means they keep working; only *new* users
must pass one flag. The asymmetry is the argument: flip-with-grandfather strands nobody, while
KEEP's cost (unrequested surface on every runtime, red `doctor` on python3-less hosts) compounds
with every install. The real cost of my choice is discoverability: a user who would have found the
skill by having it, now won't. Given the only demonstrated user is the owner's own German tender
work (below), that cost is currently theoretical.

## Concerns / open questions

**1. "Availability, not permission" — attacked, and it survives only half.** For a human tool the
defence is sound: a binary in `PATH` does not act. An agent runtime is different in one specific
way the README does not address: skills are surfaced to the model by name + description, and the
model *self-selects* by matching tasks against descriptions. The gates (`SKILL.md:10-22`,
`:57-72`: read-only default, E0–E8 approvals, no credentials) bind *actions inside* the skill;
nothing gates *selection*. So the defence fully answers "can it hurt me?" — no: mis-invocation
costs some tokens and a workflow that stops at its first approval gate — and entirely dodges
"should it be competing for selection on every task of every user?" The bidding description is
long and dense with generic trigger vocabulary ("analyze requirements, delivery, suppliers,
security, contracts, or pricing" — `SKILL.md:3`), so it will match non-procurement tasks on
runtimes whose owners never wanted it. This is not hypothetical surface area: on this machine all
three user-scope runtimes carry it via default install (`~/.kimi-code/skills/`,
`~/.claude/skills/`, `~/.codex/skills/`, markers dated 2026-08-13, v2.8.0), and it is in *this*
session's own skill listing — I am one of the models doing the matching. Verdict: the defence is
true about harm and irrelevant to the actual question, which is selection pressure and consent.

**2. A fourth instance of *printed rule binds only where enforcement lives*? No — but the
facilitator's framing is half-right, and the other half is more interesting.** The three recorded
instances are printed *norms* with no mechanism at the point of action: the `standard` fix-up cap
of 2 under which `skills-cli-install-path` ran 15 cycles
(`ideas/meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md:94`; independently confirmed,
`ideas/cognee-mechanism-mining/round-02/hermes-1.md:116`); round-1 independence, explicitly "a
social one … no enforcement beyond agent discipline"
(`internal/protocol/defaults/COOPERATION.md:886`); and the later-review-rounds obligation to
address every other reviewer explicitly, complied with in 23/348 = 7% of reviewer files
(`ideas/mas-research-mining/consensus.md:204-221`; the obligation sits at `COOPERATION.md:531` in
this deck's current copy — the audit's `:527` is line drift, same sentence). A README warning is
not a norm and binds no one, so filing it as instance four is category slippage that would weaken
the class. The accurate statement: at integration, the "default-install availability expansion"
was recorded under **"Documentation duties (must appear, are not blockers)"**
(`ideas/integrate-parley-bidding-addon/FINAL.md:90-99`) — a decision-shaped concern discharged as
documentation. The three instances are print pretending to be enforcement; this is print
pretending to be a decision. Same genus — text doing the work of mechanism — different species.
Worth one more observation: the installer is the most enforcement-dense code this ecosystem has
(thirteen-plus review rounds of comment-cited hardening); the entire default-on question reduces
to one line — `if (!context.options.only) { return discovered; }` (`installer.js:897-898`). The
gap between that line's blast radius and the scrutiny the rest of the file received is itself the
best evidence that the default was inherited, not chosen.

**3. Usage evidence, widened — commands I ran.** The facilitator's search was depth-limited; mine
are not:

- `find /Volumes/.../AI_WORKSPACE` (no depth limit, pruning only `node_modules`, `.git`,
  `.gomodcache`, `.gocache`, `dist`) for the six workspace-artifact names the skill creates
  (`bid-state.json`, `procedure-profile.json`, `bid-book.md`, `requirements-register.csv`,
  `qualification-brief.md`, `pricing-worksheet.csv` — shape from
  `scripts/init_bid_workspace.py:94-135`). **Result: 32 hits, every one an `assets/templates/`
  file inside a skill copy** — the packaged tree, the BYTE ancestor, and installer-test fixtures
  under `parley-deck/.r16-kimi-scratch/` (fake `home/.codex/skills/` trees from review-round
  testing). No `bid-state.json` and no `procedure-profile.json` exists anywhere — i.e., the
  workspace initializer has never produced a surviving workspace on this machine, not even in
  test fixtures.
- `rg -uuu -l "lifecycle_state|workspace-initialized|unknown-possibly-submitted"` over the same
  root (unrestricted: no `.gitignore` honouring, hidden and binary included; excluding
  `node_modules`/archives/graphify output). **Result: 105 hits; after removing skill copies, the
  two building decks (`BYTE/.../software-bidding-multiplatform-skill/`,
  `ideas/integrate-parley-bidding-addon/`), the HITL design deck
  (`BYTE/.../dtvp-bidding-hitl-skill/`), and the `.r16-kimi-scratch` fixtures, the remainder is
  noise** — one Java heap dump (binary false positive), one coincidental field name
  (`lifecycle_state` meaning `ACTIVE`/`TOMBSTONED` in an unrelated trading-pattern project,
  `altfins/.../marao-pattern-engine-implementation-plan/round-01/hermes-1.md:454`), and one
  `round-01/` file of this very idea matching on the quoted vocabulary — path seen in the result
  list, content not opened, per round-1 independence. **No live bid workspace anywhere.**
- Timeline cross-check: IHK_PFALZ's tender artifacts are dated May 2026 (`ls -l`), predating the
  integration commit (2026-07-30, `parley-deck-skill` `714712f`) — they cannot be skill output.
  BYTE's tender work spans the integration (ideas Jul 20–Aug 6; `submission/` touched Aug 10),
  but its artifacts are hand-shaped `.docx`/negotiation files, not the skill's state machinery.
  So the lineage is: built *from* live German tender work (BYTE ↔ the shipped DTVP adapter is not
  a coincidence), integrated, and — on the evidence visible from this machine — not run in its
  packaged shape since. Limit, stated: my searches cover this machine's workspace; external users
  are unobservable from here. This remains absence of evidence — but it is now unbounded-filename
  plus unrestricted-content absence, not a depth-5 sample.

**4. A runtime without `python3` — established, not speculative.** The documented workflow invokes
`python3 scripts/*.py` at step 1 (`init_bid_workspace.py`, `SKILL.md:26`) and throughout freeze/
validate (`SKILL.md:98-118`); there is no no-interpreter fallback anywhere in `SKILL.md` or
`references/` (grep for `python3|interpreter|fallback|without python`; the one "fallback" hit is
about AI-model fallback in bid content, `references/software-bid-model.md:57`). The scripts are
stdlib-only (import scan of all seven), floor `>=3.10` declared in `parley-addon.json:3-5`. So on
such a host the skill's published commands fail at the shell; the skill text gives the agent no
instruction for that case — no fallback, no degraded mode, no "install python3" — and the honest
agent stops and reports. The safety gates are markdown and still bind, so nothing *unsafe*
happens; the skill is simply dead weight. The installer side is explicit: `doctor` probes
`python3` against the declared floor and reports the add-on `valid` **and** unavailable, exiting
non-zero (`installer.js:2185-2262`, `:2692-2708`; CHANGELOG 2.1.0: "On a Windows host where only
`python` exists, or where `python3` is the Store app-execution alias, the add-on is reported
unavailable. That is the fail-safe direction"). Consequence that belongs in the scope decision:
under default-on, every python3-less runtime — realistically a share of the Windows fleet the
winget channel targets (RECALL: python3 is not part of a default Windows install) — gets a
permanently red `doctor` for a skill its owner never asked for, and the only fix is a workaround
that does not survive the next upgrade. That is the concrete, non-hypothetical price of the
current default, and the README's defence never mentions it.

## Risks

- **Flip done naively is the worst option on the table** — worse than KEEP. It converts a silent
  relevance cost into a loud health failure on every existing install (mechanics above). If the
  deck cannot agree on grandfathering + removal ergonomics, KEEP beats bare FLIP.
- **Grandfathering freezes a consent accident.** Users who never wanted bidding keep it forever
  unless they act. Acceptable because the skill is safe-inert (read-only default, gated actions),
  and item 3 gives them a clean exit; but it should be recorded as a trade-off, not glossed.
- **`uninstall --only` claiming the core is asserted from code reading**
  (`installer.js:995-1016`, `:1712-1745`); no test covers it (`grep` of `test/installer.test.js`
  finds only the `--no-addons` cases at `:593`, `:727`). If I misread the planning path, item 3
  gets cheaper, not more expensive.
- **I did not run `doctor` on a python3-less host**; the red-doctor claim is from
  `installer.js:2185-2262`, `:2697-2708` and the 2.1.0 changelog, not from a live reproduction on
  such a machine.
- **SPLIT may be the right end state if the owner keeps building verticals** — a `parley-verticals`
  family with its own release cadence would stop this exact debate recurring per domain. Out of
  scope to decide now, but the flip mechanism (recorded-selection-aware defaults) generalizes to
  `parley-design-check`'s 615 KB, which gets its own idea; design the flag vocabulary once, not
  twice.
- **Sunk-cost arguments struck, per the prompt.** The 25+ fix-up cycles (`IMPLEMENTATION.md` —
  its own tail records cycles 26–27 and "29 fix-up cycles", `:1881-1930`) are evidence about this
  deck's review process under a new payload type, not about whether the skill stays. I use them
  for nothing. Conversely, "it forced the generic manifest mechanism" is an argument for having
  integrated it, not for keeping it default-on.
