---
agent: hermes-1
idea: named-roster-presets
review-round: 1
date: 2026-07-04
reviewed-commit: 03a7252
---

## Summary

The implementation faithfully realizes the consensus design: layered `[rosters.<slug>]`
config with per-name replace merge, pure `ResolveRoster` with fail-closed validation
against the §2 roster, `--preset` (not `--roster`) on `parley run`, `parley preset list`
with stale-member flags, and provenance as an HTML comment. `go test ./internal/config
./internal/protocol` is green; `go build ./...` and `go vet` are clean. All six block
cases from consensus item 6 are wired and tested. The non-blocking warn cases (item 7)
are not wired but are acknowledged in FINAL.md §5, which the review brief accepts.

One real defect and a few minor issues follow. The headline problem: the provenance
HTML comment is emitted *inside* the YAML frontmatter fence, where `ReadFrontmatter`
parses it as a junk key — directly contradicting consensus item 5's "not a parseable
frontmatter key." Everything else is minor or nit.

## Findings

### CRITICAL — provenance comment is inside the frontmatter fence, parsed as a junk key

`internal/protocol/workspace.go:149-155` emits the provenance line between
`participants:` and `status:`, i.e. *inside* the `---` fence:

```
participants: [claude-1,codex-1]
<!-- roster-preset: council (source: deck) -->
status: round-01
---
```

`ReadFrontmatter` (workspace.go:334) does `strings.Cut(line, ":")` on every line in
the fence. The comment contains `:`, so it is ingested as
`meta["<!-- roster-preset"] = "council (source: deck) -->"`. I confirmed this with a
throwaway test: the map contains the junk key alongside the real fields.

This contradicts consensus item 5: "Provenance = an HTML comment in 00-prompt.md ...
not a parseable frontmatter key that a tool could mistake for canonical or re-expand
from." The comment is currently a (junk) frontmatter key. The canonical fields
(`idea`, `participants`, `status`) still parse correctly because they are on their own
lines, so this is not a runtime break — but any consumer that iterates the frontmatter
map (e.g. a future re-expansion guard, or an external tool scanning for unknown keys)
will see the junk entry, and the design intent is violated.

Fix: emit the provenance line *after* the closing `---` (e.g. as the first line of the
body, above `## Problem / idea`), or at minimum below the frontmatter fence. The
`track:` line can stay in the frontmatter (it is a real canonical key); only the
comment must move outside. No test currently covers provenance placement
(`workspace_test.go` has no provenance assertion), so the regression was not caught.

### MAJOR — `parley preset list` silently drops stale-member warnings when §2 is unparseable

`internal/app/preset.go:68` gates the entire stale-member check behind `if ok`.
When `ReadRosterIDs` returns `ok=false` (missing/unparseable COOPERATION.md), `warn`
stays empty for *every* preset, even ones with members that are obviously stale. The
list prints clean presets with no indication that validation was skipped.

This is a UX trap: a user running `parley preset list` in a deck with a broken or
missing §2 table sees no ⚠ flags and may assume all presets are healthy, when in fact
no membership check ran. `parley run` in the same state *does* surface a behavior
difference (it skips membership but still blocks empty/dup — app.go:1789-1791), but
`preset list` gives no signal at all.

Fix: when `!ok`, print a one-line notice that §2 validation was skipped (e.g.
`(§2 roster unavailable — member checks skipped)`), so the absence of warnings is
not mistaken for a clean bill of health.

### MINOR — membership fail-open on unparseable §2 is not the consensus "fail closed"

`internal/app/app.go:1789-1791`: when `ReadRosterIDs` returns `ok=false`, the code
sets `rosterIDs=nil` and continues, so `ResolveRoster` skips the not-in-§2 and
inactive checks (the `len(rosterIDs) > 0` guard at roster.go:110) but still blocks
empty/duplicate. IMPLEMENTATION.md (lines 69-72) documents this as a deliberate,
narrow fail-open *only* for the membership dimension and asks reviewers to confirm.

Consensus item 6 says "if the §2 table cannot be parsed confidently, expansion fails
(no silent fallback)." The current behavior is a silent fallback for membership: a
preset with a typo'd or fabricated agent ID expands successfully when §2 is
unparseable. The implementer flags this as an open question; I lean toward
hard-stopping on `!ok` when a preset is selected (the membership check is the main
safety point of §2 validation), but this is a judgment call the author explicitly
deferred. Noting as MINOR because it is documented and scoped, but it does diverge
from the literal consensus text.

### MINOR — `--track` is undocumented in the `run` command help and printUsage

`internal/app/app.go:1730` defines `--track` and `--preset`, but the `printUsage`
Commands/Parameters section (app.go:142-241) lists neither `--preset` nor `--track`,
and the `run` usage string (app.go:1742) still reads
`parley run [--no-tui] [--no-auto] [--participants AGENTS] [--yes] TASK` with no
mention of the new flags. `parley preset list` is also absent from the Commands
listing (only the dispatch at app.go:95 wires it). A user running `parley help` or
hitting the usage error will not discover the feature. The `--preset` flag's own
help text does point to `parley preset list`, which is good, but only if the user
already knows the flag exists.

### MINOR — track-default expansion is silent to the user when no `--preset` is given

`app.go:1787`: the expansion branch fires when `--preset != ""` *or*
(`--participants == ""` and `--track != ""`). In the second case (track-default
applied implicitly), app.go:1800 prints `Roster preset %q (source: %s) → %s` — so
the preset *name* and source are shown, but there is no hint that this came from a
track default rather than an explicit `--preset`. Consensus item 8 asked for an
override hint: `track=standard → preset 'trio' (…); override with
--preset/--participants`. The current one-line print does not include the override
hint, so a user who passed only `--track standard` may not realize a preset was
auto-applied or how to override it. The provenance comment in 00-prompt.md does
record the source layer, which partially covers debuggability, but the live CLI
hint is missing.

### NIT — `ReadRosterIDs` inactive detection is a substring match on the whole row

`internal/protocol/roster.go:62` marks a row inactive if the lowercased line contains
"inactive" anywhere. This is fragile: a workspace dir or role note containing the
word "inactive" (e.g. `../inactive-archive/`) would false-positive. The test fixture
uses `(inactive — retired)` in the role cell, which works, but the matcher is
unscoped. A tighter check (e.g. matching "inactive" inside the role/status cell, or
requiring a parenthesized marker) would be safer. Low impact given §2 table
conventions, hence NIT.

### NIT — `knownPresetNames` helps the unknown-preset error, but no closest-match hint

Consensus item 6 mentioned "unknown preset (name closest match + layers searched)."
`ResolveRoster` (roster.go:99) lists known names via `knownPresetNames` but does not
suggest a closest match, and the error does not say which layers were searched. The
listed known names are a reasonable fallback, but the closest-match hint from
consensus is not implemented. Minor UX gap; the message is still actionable.

## Open questions

1. **Provenance placement (CRITICAL)** — confirmed the comment is inside the
   frontmatter and parsed as a junk key. Should it move below the `---` fence (into
   the body), or stay in the frontmatter with a parser guard that skips lines not
   matching `key: value`? My recommendation: move it to the first body line, which
   is the simplest fix and matches the consensus "HTML comment in 00-prompt.md, not
   a frontmatter key" wording literally.

2. **§2-unparseable behavior (MINOR)** — the implementer explicitly asks reviewers
   to confirm the narrow fail-open for membership. Should `parley run` hard-stop when
   `ReadRosterIDs` returns `ok=false` and a preset is selected, rather than silently
   skipping membership? This would match the consensus "no silent fallback" more
   closely but would also block preset use in decks with non-standard §2 tables.

3. **Warn cases not wired** — the three non-blocking warnings (track-reviewer
   degradation, model-diversity risk, deck override of central preset) are in
   FINAL.md §5 but not in the code. The review brief accepts "acknowledged even if
   not all wired." Should any of these be wired in round-02, or are they explicitly
   deferred? The deck-override warning in particular is cheap to add (the source
   layer is already known at `LoadRosterPresets` time) and would close the loop on
   consensus item 7.
