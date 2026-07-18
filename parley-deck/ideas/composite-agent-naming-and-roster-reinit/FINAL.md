---
idea: composite-agent-naming-and-roster-reinit
phase: final
date: 2026-07-18
drafter: claude-1
participants: [claude-1, codex-1, hermes-1, kimi-1]
inactive: [antigravity-1]   # agy quota outage this idea
signoffs: all ACCEPT (unanimous, no reservations)
status: final
supersedes: none
---

## Design of record

The full ratified design is `consensus.md` (unanimously signed off). This FINAL.md restates
the decisions as an implementation contract and adds the staged plan, scope, and non-goals.

### The seven decisions (authoritative)

1. **Identity = roster ID.** `claude-1`, `codex-1`, `hermes-1`, `kimi-1`, `antigravity-1` are
   the single identity for artifact paths, `agent:` frontmatter, signoffs, runstate, TUI,
   digests, review snapshots. Canonical IDs stay `[a-z0-9-]`; `internal/protocol/roster.go:17`
   unchanged.
2. **Display = derived composite** `family-model-effort`, computed at render from `model` /
   `model_label` / `reasoning`; shown in a new §2 `Display name` column, TUI, digests, run
   headers. Never a key; never stored as truth (an optional `display_name` override exists only
   for genuinely unlabelable cases).
3. **Freeze = per-idea snapshot.** New ideas stamp `participant-profiles` into `00-prompt.md`
   and exact profiles into `runmanifest.Manifest`; resume prefers the snapshot.
4. **Schism fixed here, minimally, forward-only.** Split `agents.Spec.ID` → `ID` (identity) +
   `AdapterID` (launch/discovery family); explicit `[roster.<id>] adapter = "<family>"` mapping;
   new fail-closed resolver `internal/agents/resolve.go` (exact spec-ID → explicit mapping →
   hard error; never a prefix heuristic); the participant string becomes the artifact identity.
   No renames; legacy spec-ID participants keep working via resolver rule 1.
5. **Command:** `parley roster init|show|diff` (no `--reinit` flag — the command is the reinit).
   Session scope writes the deck layer; machine scope writes only `~/.parley/agents.toml`.
   Shares one `RosterInit` service with `parley init`; reuses the §9.0 `pingProbe`. Open-idea
   guard: refuse → `--yes` inactive-retention → `--force-drop` records.
6. **Autonomous write:** first-class `AutonomousWrite{Mode,Args,Scope}`; built-in defaults
   corrected to be actually autonomous (claude `bypassPermissions`, codex `approval_policy=never`
   + workspace sandbox, hermes `--yolo`, agy already ok, kimi plain `-p`); fail-closed honesty
   (unverified confinement ⇒ bit unset). Skill states the requirement + per-CLI mapping.
7. **Speed:** `fast` = new default on a separate axis; same model + same effort, faster output;
   never in the name; the legacy `profiles.*` downgrade vector is rewritten or dropped; a
   speed-invariance guard test locks it.

Naming grammar, the five concrete display names, the full Go touch-point list, and the
per-disagreement resolutions are in `consensus.md` §A–§7 and are part of this contract.

## Staged implementation plan

The stages are independently shippable and ordered so each keeps `go build ./... && go vet &&
gofmt -l && go test ./...` green. Default implementer: claude-1 (FINAL drafter); codex-1 may
claim the Go-core stages.

- **S1 — Naming core (pure, no wiring).** `internal/agents/naming.go`: `Sanitize`,
  `DisplayName(spec)`, `ParseDisplayName` (right-to-left, fail-closed), collision suffix.
  Full unit tests (dot cases, `..`/edge-dot rejection, all-digit model, legacy `family-N`,
  agy tier). No behaviour change elsewhere. Ships alone.
- **S2 — Spec split + resolver.** `discover.go`: `Spec.ID` vs `AdapterID`; vendor branches
  (`cleanParticipantEnv`, `isolatedAgentHome`, hermes env) key off `AdapterID`. New
  `resolve.go`. `runner.go` `selectedAgents` → resolver; `runAgent` artifact path +
  `validation.go` expected `agent:` = participant string. `[roster.*]` mapping in
  `runtime.go`. Tests: resolver (exact/mapped/unknown/ambiguous), artifact identity stable
  across a mid-run model change. Forward-only; legacy decks still run.
- **S3 — Config fields + speed semantics.** `runtime.go`: `model_label`, `autonomous_write`,
  template `speed = "fast"`; kill the `profiles.*` downgrade convention; `LoadAgentSpecs`
  keeps model/reasoning invariant under speed. Speed-invariance guard test.
- **S4 — Autonomous built-ins.** Correct `defaultBuiltinSpecs` per the §C table; `AUTO` column
  in `PrintRuntimeMatrix`; preflight WARN; static mapping check. CHANGELOG note for the claude
  posture change (`acceptEdits` → `bypassPermissions`, workspace-scoped).
- **S5 — `parley roster` command.** `internal/app/roster.go` (`init|show|diff`), the shared
  `RosterInit` service, both scopes, open-idea guard, `--dry-run`/`--yes`/`--force-drop`,
  config-only-family catalog fallback (kimi). `parley init` chains `roster init --scope session`.
- **S6 — §2 writer + display rendering.** §2 `Display name` column via the allowlist shared
  with the drift guard; `digest.go` + `roundsummary.go` render the display (widen/truncate the
  DISPLAY, never the ID; `@"%-13s"` vs 23-char `agy-gemini3.5flash-high`). §2-writer round-trip
  test against `TestEmbeddedDefaultMatchesLiveDeck` zones (drift guard stays green).
- **S7 — Skill side (`parley-deck-skill`).** SKILL.md: the autonomous-write requirement + the
  per-CLI mapping (source of truth = the spec field); the `parley roster init` bootstrap step;
  the display-name convention; `fast`-default note. If any protocol wording changes, edit BOTH
  `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md` (drift guard) +
  re-sync `references/COOPERATION.md`; bump `meta/version.json` protocolSha256 + changelog.

## Scope / non-goals
- **In scope:** components A–D, the minimal schism fix, the `parley roster` command, the skill
  wording, tests, and the drift-guard-safe §2 change.
- **Non-goals (explicit):** composite-as-ID / any historical artifact rename; a live per-CLI
  autonomous-write sentinel by default (opt-in `parley agents verify --full` only); building
  real headless "fast output" mechanisms beyond claude `/fast` (verified-or-noop); widening the
  §2 ID charset; an alias subsystem for `claude-1`↔`claude`.

## Verification (definition of done)
`go build ./... && go vet ./... && gofmt -l` clean; `go test ./...` green including the new
naming/resolver/identity/guard/speed-invariance/§2-writer tests; `TestEmbeddedDefaultMatchesLiveDeck`
green; a real `parley roster init --dry-run` in this deck prints the correct per-scope file
changes and the five display names; the skill installs clean (`node bin/parley-deck-skill.js
doctor --target all`). Then: full multi-channel release of `parley` (minor bump) and, if the
skill changed, `parley-deck-skill` — per the standing release discipline.
