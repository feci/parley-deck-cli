---
agent: hermes-1
idea: roster-membership-overlay
round: 1
date: 2026-08-19
---

## Summary

**NO CHANGE to the membership authority model.** The owner's 99% case is served by inheritance
(PRIMARY: `parley roster show` reports six agents all `inherited-roster`; `agents.toml` has no
`[roster.*]` blocks). The overlay's only new value — a deck wanting the global roster plus/minus
one member — has zero demonstrated demand (no delta syntax appears in any of 41 deck directories
found on this volume; only full declarations or pure inheritance). The live disease is stale full
copies (measured below), cured by migration onto inheritance, not by a new authority layer. I
recommend NO CHANGE with a recorded trigger that reopens a narrow, explicit extend-only overlay
only when a post-migration deck demonstrates a real ±delta need.

## Proposed approach (what I would build — nothing in the authority model)

1. **Migration, not authority change.** Migrate the 37 full-declared decks (measured, see Q2) back
to inheritance. Do this per-deck, attended, git-first (18 of 26 tracked repos are dirty, 15 of 41
have no git work tree at all — bulk edit would create unauditable state).
2. **A safe transition command** (separate small idea, fast-track). Today removing `[roster.*]`
without also clearing the §2 generated table falls to legacy-rule 2 (reproduced empirically,
see Q3). A `parley roster inherit`-style verb that empties both atomically, refusing if the
result would be an empty quorum, removes the footgun.
3. **Record the re-open trigger** (this file, Q1 and Q6). After migration, when a real deck
expresses a demonstrable need for "machine roster ± one agent" (two+ instances, or explicit
user instruction), open a new deliberation idea scoped narrowly to: a new explicit syntax (e.g.
`[roster-overlay]` with `add:` / `remove:`), never a reinterpretation of existing `[roster.*]`
blocks; add-only in v1; full-override preserved as back-compat; fail-closed when both forms
are present; a new additive STATUS term (`overlay-roster`) rather than changing what
`inherited-roster` means; and the render/rule-2 interaction guarded so a removed overlay
never silently resolves to a stale §2 authority.

## Answers to the six questions

### 1. Operations (extend-only vs +/−)

**Not built now; if ever, extend-only in v1.** The sibling precedent (`protocol-overlay-local-extension`,
ratified 2026-08-08: extend-only per user ruling, `replace` deferred) transfers its risk logic:
removal silently shrinks quorum, which is the operation the deck's review/signoff system is
designed to detect on committed files, not through hidden delta expressions. The mechanical
argument does NOT fully transfer — protocol `replace` needs a block-ID registry that does not
exist; membership `-` names discrete agent IDs and needs none — but the risk still outweighs the
benefit when the benefit is measured at zero instances. A full `[roster.*]` declaration already
expresses "global minus this agent" explicitly and reviewably; an overlay's `-` would only save
a file line at the cost of hidden state. If the user wants removal later, include it explicitly
in the trigger-scoped design, not in v1.

### 2. What happens to a deck that declares members today (the 36 / 41 decks)

[PRIMARY, verified] Commands run from repo root on 2026-08-19 (shared volume, maxdepth 4):
`find ... -name parley-deck` → 41 directories. Classification by parsing committed TOML blocks
(`grep` on `[roster.`) and the rendered §2 body (reading `COOPERATION.md` §2 in each) yields:
37 full-declared (5–10 `[roster.*]` blocks each), 0 legacy-table-only, 4 inherit/empty.
[PRIMARY] Spot-check of first 12 declaring decks: all declare the same core five IDs
(`claude-1`, `codex-1`, `hermes-1`, `kimi-1`, `opencode-1`); 4 also include `antigravity-1`
(still listed despite being retired per `agents.toml` comment); 2 include unprefixed duplicates
(`claude-code`, `codex` etc.); **none include `zcode-1`** (the newest machine-roster member).
That is the stale-copy shape: frozen copies of an earlier global state.

Therefore a full `[roster.*]` list **must stay a full override**. Reinterpreting it as a delta
would silently change the quorum for 37 of 41 decks in one release (machine ∪ local instead of
local), which is Q2's dangerous case with a measured blast radius. The only safe path is a
new explicit syntax; anything implicit is out.

### 3. Rule 2, the legacy §2 table (verified by running)

[PRIMARY] `sed -n '160,210p' internal/config/runtime.go` confirms the authority order: (1)
committed TOML, (2) valid legacy §2, (3) machine. [PRIMARY] Reproduced on `/tmp` copy: deleting
`[roster.*]` but keeping the §2 rows makes `roster show` fall to `legacy-roster` and show only
those rows — all machine members vanish. Removing the rows too reaches inheritance. [PRIMARY]
`go test ./internal/protocol/...` fails on three mutations: duplicated header, prose inserted
before the header (outside the normalized zone), missing header — all fail closed. [PRIMARY]
The facilitator's `49afe45` commit message confirms this is the exact trap that required the
manual two-file cleanup.

Position: **leave rule 2 untouched in this idea.** It guards a near-empty population here
(0 of 41 rely on it; 37 TOML-declared, 4 inherited), but it remains load-bearing for pre-cutover
decks elsewhere in the fleet. What should change is the **ergonomics** — the safe one-command
inherit transition (Q6/proposed approach) — not the authority semantics.

### 4. Visibility / STATUS (verified by running)

[PRIMARY] `parley roster show` → 6 agents, STATUS `inherited-roster` on all six; `parley roster
show --scope machine` → `ok`. `parley roster show --explain hermes-1` reports the exact
layered per-field mapping the prompt quotes (`adapter` from `~/.parley/agents.toml`, `model`
`fireworks/inkling`, `speed` `deep`, `status` `inherited-roster`). [SECONDARY, from `COOPERATION.md`
§4.0 and the skill's `SKILL.md`]: STATUS vocabulary is closed; adding a term is allowed, changing
an existing column's meaning is not. [SECONDARY, from `protocol-overlay-local-extension/FINAL.md`
B6 / `roster-projection-schema` deferred item]: a new `overlay-roster` term (additive, not
replacing `inherited-roster`) is the shape a future overlay would use.

Under NO CHANGE: no contract change needed. The current output already answers the user's
question in one table, per-agent, with provenance. If an overlay ever ships: the roster must
be a computed set, so each row needs provenance (`machine-base` vs `local-add` vs `local-remove`);
a new STATUS term (`overlay-roster`) must be added additively; changing what `inherited-roster`
means is out of bounds per the frozen-column contract.

### 5. The anti-goal (nine rosters across 40 decks — verified status)

[SECONDARY, named: `COOPERATION.md` §2 — the historical measurement is not reproducible here;
no `parley fleet` command exists, no fleet registry, and no archive of the 40-deck state from
the 1.41.0 cutover.] [PRIMARY, this session's census]: 41 directories found; 37 full-declared
with identical core IDs; 4 inherited; 0 relying on the legacy table; the active damage is stale
full-copies that predate the current machine-roster shape (no `zcode-1`). [PRIMARY, code]
`runtime.go` explains why the union approach reproduced the disease: pre-1.41.0 `LoadRoster`
unioned layers, so `roster render` wrote resolved inherited members into §2, which then froze
as authority and went stale on the next machine change — exactly the mechanism the anti-goal
describes.

How NO CHANGE avoids reopening it: inheritance keeps one canonical source (machine file);
no second grammar; no computed set whose render output could become an authority; no hidden
delta. What to measure post-migration (quarterly, one command — my census command, saved as the
auditing pattern): (a) count of committed membership sets different from machine base without
recorded rationale; should drop from 37 toward 0. (b) Per declared override: dated rationale in
the file (this deck's `opencode-1` model is the pattern — `agents.toml` carries a comment
explaining why it was added/removed). (c) After any future overlay: that removing the overlay
file never resolves to a stale §2 authority (the #7 trap, guarded by the fail-closed design).

### 6. Do nothing — and the concrete situation it does not serve

The honest case where NO CHANGE leaves a deck worse off: a deck that wants "machine roster plus
(or minus) exactly one agent" without freezing a full 6-member copy. Today that deck must declare
a full list and maintain it through every machine change — a real ergonomic cost. Across 41
decks on this volume, that shape appears **zero times**; only full copies (37) or pure inheritance
(4). [SECONDARY, from `COOPERATION.md` §2 / prompt §5: 36 decks synced 2026-08-06; no fleet
registry verifies the exact subset, but the measured pattern — identical five-ID declarations —
suggests most of those 36 are the same stale-copy population, not deliberate overrides.]

Recommendation: **NO CHANGE**, with approach §3's trigger: once migration completes and a real
post-migration deck demonstrates a ±1 need (two+ instances or explicit user instruction), open a
narrow deliberation scoped to an explicit add/remove syntax, never implicit reinterpretation.
That keeps the design ready without building a contract change against a zero-demand case.

## Concerns / open questions

- **Census scope limitation.** I counted one SMB volume (`/Volumes/My Shared Files`) via `find`,
not a fleet registry. The 36-deck `fleet-protocol-sync` subset (2026-08-06, SECONDARY, named:
prompt §5 / no `parley fleet` command / no registry file) is unverified. Worktree clones inflate
the 41 count; independent non-repo directories (`aditoLeads`, `lustrator`, `rev-kimi-scratch`,
`design-mail-fe`) corroborate the stale-copy shape but reduce the true denominator. Any action
beyond NO CHANGE should be based on the verified machine file and this deck's state, not on the
uncertain fleet count.
- **Copy vs deliberate declaration.** The 12-deck spot check shows identical five-ID blocks
without `zcode-1` or `antigravity-1` removal; that strongly suggests `roster init` output or a
past `roster sync`, not a deliberate per-deck choice. The migration (proposed approach #1) must
confirm per deck before removing a declaration — a deck that genuinely wants a different quorum
keeps its full override, which remains a supported (if rare) state.
- **`active` does not layer (verified from code).** `internal/config/runtime.go` (read past the
quoted authority block) shows `applyAuthorityState` anchors `active` to the authority layer,
not to a value layer. A gitignored file cannot deactivate a member. If the owner's sentence
"apply local if something is changed locally" includes activation, that is a separate open
question — this idea should answer it explicitly rather than inherit it silently.
- **Unprefixed duplicate IDs.** Two declaring decks carry `claude-code` / `codex` / `hermes` /
`kimi` without the numbered suffixes, alongside the numbered forms. That looks like adapter-
name residue, is orthogonal to the overlay question, and deserves a separate small audit during
any migration.
- **The prompt's framing is the facilitator's own, and I wrote the response, not a peer.** Per
§15.6 (correlated agreement / shared prior): if this round closes unanimous NO CHANGE, the
consensus must record that the agreement stems from a shared reading of the measured census,
not from the facilitator's declared framing — and must record the steelman for the overlay
(this file's approach §3 with the trigger condition: post-migration decks demonstrating ±1
demand, which is exactly what would change my recommendation).

## Risks

- **Of NO CHANGE:** the ±1 residue stays unserviceable. Mitigated by the explicit trigger (post-
migration real demand → open a new deliberation idea with the design sketch in Q6). Cost if it
fires: one stale deck, visible in the next census.
- **Of the alternative (build overlay now, which I reject):** silent semantic change to 35–37
existing full declarations (the dangerous case of Q2, measured); two grammars coexisting in one
file; the render/rule-2 interaction (the #7 trap, verified by test); an additive STATUS contract
change; and — per `protocol-overlay-local-extension/FINAL.md` (B4 / round-1 reversal): a registry
and a `replace` mechanism were unanimously wanted and then unanimously dropped once scoped,
with the explicit finding that "round 1's unanimity was worth very little." The measured cost
of building against zero demonstrated demand is disproportionate.
- **Of the proposed migration (if adopted separately):** editing 37 committed files across repos,
18 of which are dirty and 15 of which have no git work tree, risks unauditable state. Mitigation:
per-deck, attended, git-first ordering, the one-command transition verb, and no bulk script.
- **Verification integrity (§15) checks applied to this file:** all measured claims (`parley
roster show`, `sed` on `internal/config/runtime.go`, `go test`, census commands, `/tmp` reproduction
of rule 2, drift-guard mutations) tagged [PRIMARY] with the command or file locator; historical
claims (1.41.0 disease, 36-deck sync, 9-rosters measurement) tagged [SECONDARY] with named
source (`COOPERATION.md` §2 / prompt §5 / `protocol-overlay-local-extension/FINAL.md` / `runtime.go`
comment); no untagged [RECALL] claims used to support the recommendation. The owner sentence
is quoted from the prompt and treated as a desired-behavior statement, not a design — consistent
with the prompt's own instruction ("do not treat it as the answer").
