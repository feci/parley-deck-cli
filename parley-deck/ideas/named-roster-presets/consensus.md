---
idea: named-roster-presets
drafted-by: claude-1
date: 2026-07-04
track: standard
participants: [claude-1, codex-1, hermes-1]
---

## Agreed decisions

Strong convergence across all three round-01 files. Presets are pure expansion sugar;
`participants:` in `00-prompt.md` stays the single canonical quorum. No protocol-text
change is required (a machine-local convenience, like `headless-agents.local.json`).

1. **Config shape** — a sibling `[rosters.<slug>]` block with `participants = [...]`,
   parsed in the same layered path as agent specs (`internal/config/runtime.go`):
   built-ins → `~/.parley/agents.toml` → deck `parley-deck/agents.toml`, **merged
   per-preset-name** (a higher layer redefining a preset REPLACES its participant list,
   does not append). Track defaults: `[defaults.track_rosters]` with `fast`/`standard`/
   `deliberation` → preset name, merged per-field.

2. **IDs are §2 canonical roster IDs** (`claude-1`, `codex-1`, …), NOT `agents.toml`
   family keys (`claude`). Rationale (codex + hermes agree): unambiguous quorum/signoff
   identity, no silent alias layer that breaks in multi-instance decks. A `protocol.
   ReadRosterIDs(root)` helper validates preset members against the §2 table; if the
   roster table cannot be parsed confidently, expansion **fails** rather than falling
   back to installed agents.

3. **Flag name `--preset <slug>`** (hermes) — not `--roster`, which collides with
   roster management. `--preset` + `--participants` together = hard error (ambiguous).

4. **Expansion happens once, at idea creation, before preflight and
   `runcontrol.Create`** (codex). The expanded concrete list is written into
   `00-prompt.md participants:`. An existing `00-prompt.md` is never re-expanded, even
   if the preset later changes (frozen quorum).

5. **Provenance = an HTML comment** in `00-prompt.md`, e.g.
   `<!-- roster-preset: council (source: deck) -->` (claude's middle path, resolving
   codex "no frontmatter keys" vs hermes "want provenance"): debuggable, but not a
   parseable frontmatter key that a tool could mistake for canonical or re-expand from.

6. **Fail loudly (block, before the idea dir is created):** unknown preset (name closest
   match + layers searched); empty or duplicate members; member not in §2; member whose
   §2 row is inactive (distinct error); expansion that collapses to facilitator-only
   (§1 non-solo); a selected track-default that names an unknown preset.

7. **Warn (not block):** preset size trips a track reviewer minimum (§4.0 degrades by
   design — legal but fragile); all members share the implementer's model (LE-3 risk,
   implementer not yet known); a deck override of a central preset (print which layer
   won).

8. **Track-linked default** applies ONLY when neither `--preset` nor explicit
   `participants:` is given AND the creation entry point knows the track; it is applied
   but printed with an override hint (`track=standard → preset 'trio' (…); override with
   --preset/--participants`). If the track is unknown at creation, no track-default
   expansion (codex). Track policy (§4.0) still applies AFTER expansion.

## Deferred follow-ups

- Optional `description` / `fits_tracks` per-preset fields for a richer `--preset list`
  (hermes concern 2) — speculative, deferred.
- A structured protocol-level roster alias map — only if bare family names in presets
  are ever wanted (explicitly rejected for v1).

## Dismissed / non-goals

- No per-preset model/reasoning/transport overrides (all three agree).
- No dynamic/conditional rosters.
- No change to the deck-bootstrap roster confirmation flow.

## Signoffs

<!-- each participant appends its own block -->

### Signoff: claude-1 — 2026-07-04
Status: ✅ ACCEPT
Drafter. Captures the §2-canonical-ID decision, --preset flag, expand-at-creation,
HTML-comment provenance, and the block-vs-warn split. No protocol change needed.

### Signoff: codex-1 — 2026-07-04
Status: ✅ ACCEPT
Accepting with the key constraint that `[defaults.track_rosters]` only expands when the idea-creation path already knows the track; otherwise no track-default preset is applied. The consensus captures this without changing `participants:` as the canonical quorum.

### Signoff: hermes-1 — 2026-07-04
Status: ✅ ACCEPT
`--preset` naming (item 3) and the block-vs-warn split (items 6–7) match my round-01 exactly. Provenance shifted from my frontmatter keys to claude's HTML comment (item 5) — still debuggable and non-re-expandable, so my "want provenance" concern is satisfied; §2 canonical IDs (item 2) confirmed.
