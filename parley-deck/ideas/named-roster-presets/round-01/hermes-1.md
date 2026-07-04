---
agent: hermes-1
idea: named-roster-presets
round: 1
date: 2026-07-04
---

## Summary

Named roster presets are worth adding as expansion sugar, but the UX
hinges on three things that the brief under-specifies: (1) the naming
gap between `agents.toml` keys (bare family names like `claude`, `agy`)
and §2 roster IDs (suffixed, like `claude-1`, `antigravity-1`); (2) the
interaction between a preset's participant count and the track's
reviewer/quorum minimums; (3) provenance — recording *which* preset
and *which config layer* produced the expanded list, so a stale preset
is debuggable without re-deriving it by hand. I support the v1 scope
(expand-at-creation, no per-preset model overrides, no dynamic rosters)
and think the track-linked default is the highest-leverage part, but it
is also the part most likely to confuse if applied silently.

## Proposed approach

**Config shape.** Reuse the existing `[agents.<id>]` namespace pattern;
add a sibling `[rosters.<slug>]` block in both `~/.parley/agents.toml`
and `parley-deck/agents.toml`. Each preset is just `participants =
[<id>, ...]`. v1 carries no model/role/effective-speed fields — those
stay on the agent and the per-idea `roles:` map. Precedence is the
existing layering: built-ins → `~/.parley/agents.toml` → deck
`agents.toml` (deck wins), merged per-key so a deck can override one
preset while inheriting the rest.

**ID resolution — the concrete decision.** The §2 roster IDs are
`claude-1`, `codex-1`, `hermes-1`, `antigravity-1`; the existing
`agents.toml` keys are `claude`, `codex`, `agy`, `hermes`. A preset
`participants` list MUST use §2 roster IDs (the canonical quorum
identity), not `agents.toml` keys. Trade-off: this forces preset
authors to know the suffixed IDs, but it keeps quorum/signoff
accountability unambiguous (§1, §5) and avoids a silent alias layer
that could map `agy` → `antigravity-1` on one deck and to nothing on
another. If we instead allowed bare family names in presets, every
expansion would need a resolution step that can fail in
multi-instance decks (`claude-1` vs a future `claude-2`); not worth
the ambiguity for v1.

**CLI surface.** `parley idea new <slug> --preset <slug>` (prefer
`--preset` over `--roster`: `--roster` reads as "the §2 roster" and
invites confusion with roster management). Flag behavior:

- `--preset <slug>`: expand, write `participants:` into `00-prompt.md`,
  refuse to write the file if expansion fails (see failure modes
  below). Print the resolved participant list and the source layer
  before writing.
- `--preset list` (or `parley preset list`): print available presets
  with participant counts, source layer (central/deck/override), and a
  one-line "fits track X" hint. Cheap discoverability win, no
  interactive prompt.
- No `--preset` and no explicit `participants:`: if a track-linked
  default exists for the chosen track, USE IT but print
  `track=standard → preset 'trio' (claude-1, codex-1, hermes-1); override with --preset or --participants` so the silent default is
  visible. If no default, behave as today (author must supply
  participants).

**Track-linked defaults.** Optional `[track_defaults]` block:
`fast = "pair"`, `standard = "trio"`, `deliberation = "council"`.
Applied ONLY when the author supplies neither `--preset` nor an
explicit `participants:` list. The default is a convenience, never a
fallback when an explicit choice fails — a typo'd `--preset councel`
MUST error, not fall back to the track default.

**Provenance.** After expansion, write two advisory fields into the
`00-prompt.md` frontmatter (mirroring the existing pattern where
`roles:` is advisory and does not change quorum):

```yaml
participants: [claude-1, codex-1, hermes-1, antigravity-1]
roster_preset: council          # advisory; participants: above is canonical
roster_preset_source: deck      # central | deck | built-in
```

`participants:` stays the single source of truth for quorum (§2, §5)
and signoff (§3). `roster_preset*` is debug-only: a reader who sees a
weird roster can trace it back to the preset and layer that produced
it, then fix the preset rather than hand-editing participants on every
future idea. Trade-off: two extra frontmatter keys that someone could
mistake for canonical. Mitigated by naming them `roster_preset*`
(advisory-sounding) and documenting them as non-quorum in the skill,
the same way `roles:` is documented. The alternative — no provenance,
expansion is opaque — is cleaner today but makes a stale preset
(`gemini` retired, `agy` renamed) a recurring support burden.

**Failure modes — what MUST fail loudly, at expansion time, before
`00-prompt.md` is written:**

1. Preset slug not found in any layer (typo). Error names the closest
   match and the layers searched.
2. Preset references an agent ID not in the §2 roster table. Error
   points at the preset definition file:line, not just "agent not
   found".
3. Preset references an agent whose §2 row is marked inactive (§2:
   "mark inactive, do not delete"). Distinct error from (2) — inactive
   is a state, not a missing row.
4. Preset would collapse to facilitator-only (violates §1 non-solo
   and Phase 0's "MUST NOT silently collapse to only the
   facilitator"). Error cites §1.
5. Track-linked default references a preset that does not exist.
   Error at `parley idea new` time, not deferred to Phase 0.

**Failure modes — what MUST warn (not block) at expansion time:**

1. Preset participant count trips a track minimum *that track policy
   applies after expansion* (per the brief). Example: `pair` preset on
   `standard` → standard's "2 reviewers" degrades to 1 when there are
   only 2 participants (§4.0 binding/challenge). This is legal but
   fragile; print a warning naming §4.0 and LE-3 (model diversity).
   Do NOT block — the brief is explicit that track policy applies
   after expansion, and blocking here would make `pair` unusable on
   `standard`, which may be a legitimate choice for a small idea.
2. Preset lists participants who all share the implementer's model and
   `require_model_diversity: true` would trip (LE-3). Warn, since the
   implementer isn't chosen yet at `parley idea new` time.
3. Deck override of a central preset with the same slug. Print
   "using deck override of rosters.<slug>" so the facilitator's
   mental model isn't silently swapped.

The split between block and warn is the core UX judgment: block on
things that make the idea invalid at Phase 0 (non-solo, unknown
agent, inactive agent); warn on things that make the idea fragile but
legal (track/reviewer degradation, model diversity). This matches the
protocol's existing stance — §1 non-solo is a hard invariant, §4.0
reviewer counts degrade by track rather than blocking.

## Concerns / open questions

1. **`--preset` vs `--roster` flag name.** I propose `--preset`. The
   brief says `--roster` "(or equivalent existing entry point)". `--roster` collides conceptually with roster management (`parley roster`
   subcommands to add/remove agents per §2). Worth a quick decision
   before implementation; I'd push back on `--roster`.

2. **Should presets carry an optional `description`/`fits_tracks`
   field?** A one-line `description = "fast pair for small fixes"` and
   `fits_tracks = ["fast"]` would make `--preset list` genuinely useful
   and let the CLI warn when `--preset pair --track deliberation` is
   chosen. This is a small additive field, not the excluded
   "per-preset model overrides" non-goal. But it is speculative
   config surface; I'd defer unless `--preset list` turns out to need
   it.

3. **Roster IDs vs `agents.toml` keys — confirm the resolution rule.**
   I assert preset `participants` use §2 roster IDs (`claude-1`), not
   `agents.toml` keys (`claude`). The existing `agents.toml` uses bare
   family names; a preset author copying from `agents.toml` will write
   `claude` and hit failure mode (2). This is correct behavior but
   needs a clear error message and skill docs; otherwise it reads as
   "the config format is inconsistent." Codex-1's config-internals
   lens should confirm whether the loader already has an alias map or
   whether we're adding a new resolution surface.

4. **Track default applied silently vs explicitly.** I argue the
   track default should be applied but printed, with an override hint.
   Alternative: refuse to create the idea without `--preset` or
   `--participants`, forcing an explicit choice. The latter is safer
   (no silent default) but reintroduces the friction presets are
   meant to remove. I lean toward "apply + print" because the
   provenance field makes the choice auditable and reversible by
   editing `00-prompt.md`.

5. **Does expansion re-run on existing ideas?** If an author edits
   `rosters.council` after creating an idea with `--preset council`,
   the idea's `participants:` is already frozen (canonical). Should
   `parley` warn on a stale `roster_preset:` field when re-opening an
   idea? My instinct: no — once expanded, the idea is detached from
   the preset, and the provenance field is debug-only. But this
   should be explicit in the skill docs so nobody expects
   retroactive re-expansion.

6. **Preset name collisions across layers.** Deck overrides central
   per-key. But what if central has `rosters.pair` and deck has
   `rosters.pair` with identical participants — is that an "override"
   or a "no-op"? The provenance `roster_preset_source: deck` would
   record it as an override even though nothing changed. Minor, but
   worth a one-line rule: deck wins regardless of content equality,
   and provenance records the winning layer.

## Risks

- **Stale preset references a retired/renamed agent.** The
  `agents.toml` comment already records that `gemini` was retired in
  favor of `agy`. A central `rosters.council` defined before that
  rename will list `gemini` and fail every new idea that uses it until
  someone fixes the preset. Failure mode (2)/(3) makes this loud, and
  the provenance field points at the preset to fix — but in a
  multi-deck setup the preset may live in `~/.parley/agents.toml` and
  break every deck at once. Mitigation: the error message should
  suggest editing the named preset file, and `parley preset list`
  should flag presets referencing inactive/missing agents so they're
  fixed before idea creation, not during.

- **Track-reviewer degradation hidden behind a friendly name.** `--preset pair --track standard` is a footgun: standard wants 2
  reviewers, pair gives 2 participants, §4.0 degrades to 1 reviewer
  (the non-implementer). A facilitator reaching for "pair" because
  it's fast may not realize standard's review gate is now as thin as
  fast's. The warning at expansion time mitigates this, but only if
  the facilitator reads it. Provenance records the preset; it does
  not record whether the warning was seen. This is the strongest
  argument for the alternative in concern (4) — refuse without
  explicit `--participants` — but that sacrifices too much of the
  feature's value.

- **Silent default masks a deliberate roster choice.** If a
  track-linked default is set and an author forgets to pass
  `--preset`, they get the default's participants without having
  chosen them. The print-at-creation mitigation assumes the author
  reads CLI output; in an automated/driver-run (§14) context there is
  no human reading it. For auto-driven idea creation, the track
  default should still apply (drivers need a deterministic roster),
  but the provenance field is the audit trail that lets a later
  reviewer see "this was a default, not a deliberate choice." LE-2
  (driver auto-advance) is consistent with this.

- **Preset becomes a second roster, fragmenting the source of
  truth.** §2 says the roster table is project-specific and
  canonical; presets add a layer that *refers to* the roster but
  feels like a parallel structure. If presets grow fields (model
  overrides, role templates) in later versions, the pressure to
  treat `rosters.*` as the real roster grows. v1's non-goals
  (no per-preset model overrides, no dynamic rosters) hold this
  line; the risk is v2 scope creep. The provenance field helps by
  making "this came from a preset" visible on every idea, so the
  roster stays obviously canonical.

- **Provenance fields mistaken for canonical.** A future agent or
  tool could read `roster_preset:` and recompute `participants:`
  from it, overriding the frozen list. This would break the
  "participants is canonical" invariant. Low likelihood (the names
  are clearly advisory) but worth a one-line note in the skill that
  these fields are read-only debug context, never a re-expansion
  source.
