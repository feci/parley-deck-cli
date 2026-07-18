---
idea: composite-agent-naming-and-roster-reinit
phase: consensus
date: 2026-07-18
participants: [claude-1, codex-1, hermes-1, kimi-1]
inactive: [antigravity-1]   # agy quota outage this idea (see inbox tooling-exception)
status: ratified
---

## Outcome

All four active participants converged. The crux was decided **against** the facilitator's
round-01 frozen-slug in favour of **option (b): the stable roster ID is the identity; the
composite `agent-model-effort` is a derived display name**. claude-1's freeze requirement is
preserved — moved down a layer into a per-idea profile snapshot (codex-1's mechanism). The
design also **fixes the pre-existing two-ID-namespace schism** kimi-1 traced (the driver keys
on spec/family IDs while the whole deck keys on roster IDs), because the roster command must
write the mapping anyway and the fix is forward-only with zero artifact renames.

## The one coherent design (ratified)

1. **Identity** — the roster ID (`claude-1`, `codex-1`, `hermes-1`, `kimi-1`,
   `antigravity-1`) is the single identity used for artifact paths, `agent:` frontmatter,
   signoffs, runstate, TUI tabs, digests, review snapshots. It never encodes model or effort.
   Canonical IDs stay in the intersection charset `[a-z0-9-]` (unchanged
   `internal/protocol/roster.go:17`).
2. **Display** — the composite `family-model-effort` is **derived at render time** from
   `model` / `model_label` / `reasoning`; shown in the §2 table (new `Display name` column),
   TUI, digests, run headers. Never a key, never stored as the source of truth (deriving, not
   storing, prevents a hand-edited `model` from forking a stale stored name).
3. **Freeze** — new ideas stamp a `participant-profiles` snapshot into `00-prompt.md` and the
   exact profiles into `runmanifest.Manifest`; resume prefers the snapshot. This delivers
   claude-1's "changing a model later must not rewrite an OPEN idea" as data, not as a key.
4. **Schism fix (minimal, here)** — split the overloaded `agents.Spec.ID` into `ID` (stable
   identity) + `AdapterID` (launch/discovery family); an explicit `[roster.<id>] adapter =
   "<family>"` mapping written by the roster command; a fail-closed resolver; the participant
   string becomes the artifact identity. Forward-only migration, no renames.
5. **Command** — `parley roster init|show|diff`; session scope writes the deck layer, machine
   scope writes only `~/.parley/agents.toml`.
6. **Autonomous write** — a first-class `AutonomousWrite{Mode,Args,Scope}` field; the built-in
   defaults are corrected to actually be autonomous; the skill states the requirement + the
   per-CLI mapping.
7. **Speed** — `fast` is the new default on a **separate axis from effort**; it never
   downgrades model/effort and is never part of the name.

## A. Naming grammar (dots allowed, path-safe)

```text
display-name := family "-" model "-" effort [ "-" instance ]
family       := [a-z0-9]+
model        := [a-z0-9]+ ("." [a-z0-9]+)*      # dots only between alphanumerics -> no "..", no edge dots
effort       := low|medium|high|xhigh|max|ultracode|clidefault
instance     := [2-9][0-9]*                     # collision suffix, >= 2
```
- **Sanitization (write-time):** from the human `model_label` (strip a parenthesized tier for
  agy), lowercase ASCII, delete every char not in `[a-z0-9.]`, collapse dot runs, strip
  edge dots, reject empty (ask for an explicit token — never invent). Effort is *selected*
  from the vocabulary, never sanitized from free text; `cli-default` -> `clidefault`.
- **Parse (read-time, fail-closed, right-to-left):** split on `-`; a trailing all-digit token
  is the instance; the next-from-right must be a vocabulary effort (hard error otherwise); the
  remaining middle token is the model (dots allowed). Family resolves via the explicit
  roster→family mapping, never a prefix heuristic (`antigravity-1`↔`agy` breaks prefixing).

| Roster ID | Family | model_label → token | Effort | **Display name** |
|---|---|---|---|---|
| `claude-1` | claude | `Opus 4.8`→`opus4.8` | max | **`claude-opus4.8-max`** |
| `codex-1` | codex | `GPT-5.5`→`gpt5.5` | xhigh | **`codex-gpt5.5-xhigh`** |
| `hermes-1` | hermes | `GLM 5.2`→`glm5.2` | high | **`hermes-glm5.2-high`** |
| `antigravity-1` | agy | `Gemini 3.5 Flash (High)`→`gemini3.5flash` | high (tier) | **`agy-gemini3.5flash-high`** |
| `kimi-1` | kimi | `K3`→`k3` | max | **`kimi-k3-max`** |

- **agy:** authoritative `reasoning` stays `cli-default` (no per-invocation flag); the display
  surfaces the label tier (`high`) as the effort token. Render rule: prefer the recorded tier
  for agy, else the reasoning field.
- **Honesty:** names follow config truth. codex today is `cli-default` in the deck, so its
  truthful display is `codex-clidefault-clidefault` until `roster init` pins `gpt5.5`/`xhigh`.
- **Code consistency:** `roster.go:17` unchanged (composites never occupy the §2 ID cell —
  they live in the new annotation column); `consensus.go:90` unchanged (already tolerates
  `[A-Za-z0-9._-]`; membership stays exact-match on roster IDs).

## B. `parley roster` command

```text
parley roster init [--dir DIR] [--scope session|machine] [--agent <family>]...
                   [--from FILE] [--dry-run] [--yes] [--force-drop] [--json]
parley roster show [--scope session|machine] [--json]
parley roster diff [--scope session|machine] [--json]
```
- **No `--reinit` flag** — the command *is* the re-init; existing values are the picker's
  defaults. It shares ONE internal `RosterInit` service with `parley init` (dedup), and reuses
  the §9.0 preflight `pingProbe`.
- **Default scope:** `session` inside a deck, `machine` outside; `--dry-run` prints exactly
  which files change per scope; session refuses to run outside a deck.
- **Session writes:** deck `agents.toml` (`[agents.<family>]` model/`model_label`/reasoning/
  speed/`autonomous_write` + `[roster.<id>] adapter`); targeted §2 table update (adds the
  `Display name` column via the allowlist-shared writer); surgical local-JSON export (rostered
  IDs only, only if present, fingerprinted — Go never parses the JSON).
- **Machine writes:** `~/.parley/agents.toml` only (managed-block markers + first-time backup);
  never touches a deck, never copies deck values up.
- **Open-idea guard (non-optional):** model/effort/speed/display changes never trip it; a
  dropped ID referenced by an open idea → default **hard refuse**; `--yes` converts the drop to
  **inactive-retention** (pinned with its profile until the idea is terminal); only
  `--force-drop` writes `excluded:`/`reincluded:` lines and proceeds.
- **Config-only families:** kimi has no built-in spec (ACP catalog/user config only) and is
  absent from this deck's §2/local-JSON — reinit must handle config-only families with a
  documented static-catalog fallback for model/effort discovery.

## C. Autonomous write

`AutonomousWrite{Mode string, Args []string, Scope string}` on `Spec` (`Scope` = `"workspace"`
or the bit is unset). No `ScopeEnforcement` taxonomy; codex-1's honesty rule kept — where
workspace confinement cannot be demonstrated, `roster init` marks the profile `unverified` and
**refuses to set the bit** (fail-closed). Surface an `AUTO` column in `PrintRuntimeMatrix`;
preflight WARNs (never gates); an opt-in live sentinel probe lives behind `parley agents verify
--full` (never default — probe cost).

**Fix the built-ins** (verified against `defaultBuiltinSpecs`):

| CLI | Today | Change to (workspace-scoped) |
|---|---|---|
| claude | `--permission-mode acceptEdits` (`discover.go:134`) | `--permission-mode bypassPermissions` + keep `--add-dir {root}` |
| codex | `approval_policy="on-failure"` (`discover.go:113,119`) | `approval_policy="never"` + keep `--sandbox workspace-write` (no full-fs escalation) |
| hermes | `--accept-hooks`, no `--yolo` (`discover.go:203`) | `--yolo --accept-hooks` |
| agy | `--dangerously-skip-permissions` (`discover.go:156`) | already autonomous; keep `--add-dir {root}` |
| kimi | no built-in spec | add spec/catalog entry: plain `-p` IS the autonomous mode; **never** `--yolo`/`--auto` with `-p` (mutually exclusive) |

**Skill wording:** adopt codex-1's normative block (required, workspace-confined, fail-closed,
redaction preserved) + the mapping table; the skill points at the spec field as the source of
truth (a vendor flag change is a config edit, not a skill revision). **Safety:** bypass flags
live in per-project config, never vendor-global; every mapping carries its workspace constraint;
secret redaction is orthogonal and stays on; changing claude's built-in is a real posture change
→ scoping flags land in the same commit + CHANGELOG note.

## D. Speed (`fast`) as a separate axis

- **Semantics:** `speed` (`fast|deep`) selects output speed only. `fast` = same model, same
  effort, faster output (Claude Code `/fast`). Confirmed efforts are invariant under speed.
- **Not in the name.** Speed is a launch property; `deep` per idea when needed.
- **Config/template:** `centralDefaultTemplate` (`runtime.go:286`) flips `speed = "deep"` →
  `"fast"` for new machine files; `roster init` writes `speed = "fast"` for new rosters (both
  scopes), respecting an explicit existing per-agent `speed` on re-run.
- **Kill the downgrade blocks:** session reinit rewrites every `profiles.{fast,deep,review}`
  in the deck JSON to carry the SAME model+thinking as the entry's top level, or (preferred)
  **drops the `profiles` table** where no distinct fast mechanism exists. (Go parses no
  `profiles.*`; today `headless-agents.local.json:28-41` pins claude `fast={sonnet,low}` and
  `:120-133` pins agy `fast={Gemini…(Low)}` — the banned legacy downgrade.)
- **Per-CLI fast mechanism — honesty rule:** the skill lists the verified fast-output
  mechanism per CLI (claude interactive `/fast`); headless agents with no verified fast flag
  are `no-op (documented)` — never a weaker model. Implementation verifies each flag before
  writing it.
- **Guard test:** a `speed="fast"` roster keeps each agent's pinned `model`/`reasoning` through
  `LoadAgentSpecs` and launch assembly — the tripwire against any future speed→model/effort map.

## Schism fix — Go touch points (complete)

`internal/agents/discover.go` (Spec split `ID`/`AdapterID`, `AutonomousWrite`, built-in fixes);
new `internal/agents/resolve.go` (fail-closed resolver: (1) exact spec-ID, (2) explicit
`[roster.*]` mapping, (3) else hard error); `internal/runner/runner.go` (`selectedAgents`→resolve,
`runAgent` artifact path = participant string, vendor env branches → `AdapterID`);
`internal/runner/validation.go` (expected `agent:` value); `internal/config/runtime.go`
(`[roster.*]` mapping, `model_label`, `autonomous_write`, template `speed`);
`internal/config/roster.go` + `internal/protocol/roster.go` (§2 writer + Display-name column +
allowlist shared with the drift guard); `internal/app/app.go` (`roster` verb group);
`internal/app/preflight.go` (probe reuse + non-gating drift WARN);
`internal/runmanifest/manifest.go` (profile snapshots); `internal/driver/digest.go` +
`internal/tui/roundsummary.go` (display rendering + width — `@"%-13s"` vs 23-char
`agy-gemini3.5flash-high`, widen/truncate the DISPLAY, never the ID).

**Tests:** sanitization/parse/collision (dot cases incl. `..` rejection, edge dots, all-digit
model, legacy `family-N`); resolver (exact/mapped/unknown/ambiguous); artifact-identity
stability across a mid-run model change; open-idea guard (drop→inactive, `--force-drop`
records); autonomous mapping per CLI incl. kimi's `-p` exclusivity; speed-invariance guard;
§2-writer round-trip against `TestEmbeddedDefaultMatchesLiveDeck` zones (stays green).

## How each disagreement was resolved
- **Crux (claude-1 frozen-slug vs the other three):** resolved to option (b). Two code facts
  sank frozen-slug — a dotted slug can't be a §2 ID (`roster.go:17`), and the driver can't
  resolve stamped composites (`runner.go:327-342`). The freeze survives as codex-1's per-idea
  profile snapshot. claude-1 conceded in round-02.
- **Store vs derive display name (codex-1 stored it; kimi-1/hermes-1 derive):** derive at
  render; accept an optional `display_name` override only for genuinely unlabelable cases.
- **`ScopeEnforcement` enum (codex-1) vs minimal struct (kimi-1):** minimal struct + honesty
  rule (unverified ⇒ bit unset); no speculative taxonomy.
- **Schism: fix now vs later (kimi-1 round-01 said "separate idea"):** fix now, minimally —
  codex-1's `AdapterID` split removed the cost objection. kimi-1 flipped in round-02.
- **`--yes` on a drop (hermes-1 recorded the break; codex-1 deferred):** merged — refuse →
  inactive-retention (`--yes`) → `--force-drop` records, ordered by blast radius.

## Deferred / round-03 residuals (non-blocking — implementation notes)
- Per-CLI headless fast mechanisms are unverified (except claude `/fast`) — verify-or-noop.
- `model_label` backfill: codex is `cli-default` today — honest display until reinit pins it.
- §2 writer + drift guard must share ONE allowlist definition (add the round-trip test).
- TUI width must land in the same release as display rendering (first visible surface).

## Signoffs
_(each active participant appends its own `### Signoff: <id> — <date>` block below)_

### Signoff: claude-1 — 2026-07-18
Status: ✅ ACCEPT

I drafted this consensus and I accept it. I conceded my round-01 frozen-slug: the three
independent code-grounded analyses were right that freezing a composite into artifact paths
fights the repo (the artifact path already uses the spec/family ID, not the §2 roster ID) and
would deepen the two-namespace schism. Option (b) with codex-1's per-idea profile snapshot
gives my freeze requirement without making a mutable string the identity. I endorse fixing the
schism in this idea (the `AdapterID` split makes it cheap and forward-only), the dot-preserving
grammar, component D as a separate axis that never downgrades effort, and the corrected
autonomous built-ins. No reservations.

### Signoff: codex-1 — 2026-07-18
Status: ✅ ACCEPT

The consensus faithfully captures my round-02 position: stable roster IDs remain protocol identity, dotted composites are derived display names, per-idea profile snapshots preserve freeze and audit semantics, and the roster command fixes the two-namespace schism now. I accept the merges that derive rather than persist the display name, replace my `ScopeEnforcement` taxonomy with the minimal `AutonomousWrite` structure plus the fail-closed honesty rule, and order roster removal as refuse → inactive retention under `--yes` → explicit `--force-drop`; these preserve the safety properties I argued for. I also accept `fast` as a separate, quality-invariant axis and the corrected autonomous built-ins. No reservations.

### Signoff: hermes-1 — 2026-07-18
Status: ✅ ACCEPT

The consensus faithfully captures my round-02 position: option (b) with the roster ID as identity and the composite as a render-derived display name; the schism fixed for new ideas via the `ID`/`AdapterID` split with `selectedAgents` matching `participants:` against the roster seat (forward-only, zero renames, legacy family-ID fallback fail-closed); the distinct `parley roster init` verb with session-default-when-deck-exists; the open-idea guard mechanics I own; component C's minimal `AutonomousWrite` struct plus the fail-closed honesty rule; and component D as a separate, quality-invariant axis with `fast` as the new default and the `Profile` downgrade vector killed. My two round-02 counter-proposals were both merged: the dot-preserving grammar (`opus4.8`, not `opus48`) is the ratified charset, and my "fix the schism for new ideas, contain it for old" posture is exactly what the consensus adopts. The merges against codex-1 (display derived not stored; minimal struct over `ScopeEnforcement`; refuse → `--yes` inactive-retention → `--force-drop`) and kimi-1 (derive-on-load, `model_label`, right-to-left fail-closed parse, forward-only migration) preserve every safety and honesty property I argued for. No reservations.

### Signoff: kimi-1 — 2026-07-18
Status: ✅ ACCEPT

The consensus faithfully captures my round-02 position: option (b) with the roster ID as sole identity and the dot-preserving composite derived at render time; claude-1's freeze delivered as codex-1's per-idea profile snapshot; the two-namespace schism fixed now and minimally via the `ID`/`AdapterID` split, explicit `[roster.*]` mapping, and fail-closed resolver — the flip I made in round-02 from "separate idea" is correctly recorded. My counters to codex-1 were merged as argued: display name derived not stored (optional override only), and the minimal `AutonomousWrite{Mode,Args,Scope}` struct with the fail-closed honesty rule in place of the `ScopeEnforcement` taxonomy. My three-way merge of the open-idea guard (hard refuse → `--yes` inactive-retention → `--force-drop` records) and my component-D analysis (speed as a pure separate axis, `profiles.*` downgrade blocks rewritten or dropped, template flips to `fast`, speed-invariance guard test) are adopted verbatim, as are the agy tier-in-label render rule, the kimi config-only-family catalog fallback, and the `codex-clidefault-clidefault` bootstrap-honesty note. No reservations.
