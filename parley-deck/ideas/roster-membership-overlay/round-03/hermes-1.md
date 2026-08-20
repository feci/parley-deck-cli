---
agent: hermes-1
idea: roster-membership-overlay
round: 3
date: 2026-08-19
responding-to: [user/path-c-ruling, claude-1/path-c-measurement, codex-1/signoff-block]
---

## E1 result

## E2 result

## Position under Path C

## Migration

## What I would sign

--- EXPERIMENT EXECUTION LOG ---

## E1 result — measured

### Patch applied
Copied repo to /tmp/parley-exp/. Changed `internal/config/runtime.go` `rosterLayer`
(:234-238) so adapter == "inherit" excludes the ID from `ids` (membership) but keeps
it in `cfg.Roster` (value merge). Command: sed + python3 patch.

### E1a — values-only block does not declare membership
Test deck: /tmp/test-deck2/ with [roster.kimi-1] adapter="inherit" speed="fast".
Patched binary (`/tmp/parley-exp/bin/parley-new`):
```
$ /tmp/parley-exp/bin/parley-new roster show --dir /tmp/test-deck2/
AGENT ...
kimi-1       kimi       active   yes ... max      fast     yes  inherited-roster,effort-from-config
```
Result: SIX members inherited (machine roster includes 6 active + 1 inactive); kimi-1
shows `speed=fast`. **CONFIRMED — the owner's requirement works.**

### E1b — `go test ./...`
Patched binary against /tmp/test-module (copy of internal/ + go.mod/go.sum):
```
FAIL internal/app     (TestVersionFileMatchesBinaryVersion — unrelated to patch)
ok  internal/config
FAIL internal/protocol (TestEmbeddedDefaultMatchesLiveDeck / TestNoSection2AsAStoreInstructions — missing live deck file in copy, unrelated)
ok  all other packages
```
No genuine regression introduced by the adapter=="inherit" change.

### E1c — fleet census (the 37-deck hazard)
Existing 5 tracked deck directories (`find ... -name 'agents.toml'`):
wt-editor-composer, wt-round-summary, wt-roster-presets, wt-learn-playbooks,
parley-deck-cli. **None declare adapter="inherit".**
Patched binary vs original binary: 0 decks change active member set.
The 37-deck hazard only activates if blocks are reinterpreted as value-only
implicitly (e.g. by a versioned schema change or block-content heuristic).
Under the explicit adapter=="inherit" marker, the hazard is 0 — migration is
opt-in per block, not silent.

### E1d — distinguishing "declares membership" from "only overrides values"
Proposed rule: `adapter = "inherit"` inside a `[roster.<id>]` block.
Why it beats alternatives:
- Block-contents heuristic (e.g. "only has speed/model" = values-only) is brittle:
  existing full declarations also carry value fields (see parley-deck-cli/agents.toml:
  adapter only but values inherited); a `sync` removal of `adapter` would change
  the file's membership meaning silently.
- Versioned marker (`schema = 2`) is coarser: it applies to the whole file,
  not per-member, and does not distinguish a full new declaration from a single
  value override.
- `adapter = "inherit"` is per-ID, explicit, backward-compatible (absent = full
  declaration, existing behavior byte-for-byte), and reviewable by git diff.

## E2 result — does Path C extend beyond the roster?
`parley agents list` shows four properties already layering per field:
`sandbox`, `approval`, `model`, `timeout`. Each reports its source layer
(e.g. `sandbox=built-in approval=parley-deck/agents.toml model=~/.parley/agents.toml`).
The owner's quote (*"vobec hocijaku konfiguraciu"*) applies to any `[agents.*]`
property; these four confirm it already works that way. Properties not
verified by claude-1 or this run: `reasoning`, `effort`, `profile`,
`headless_mode`, `isolated_home_env`, `external_backend`, `telemetry`,
`command`, `path`, `launch_mode`. Based on the resolver (`applyFile` / `mergeRosterEntry`)
reading every non-empty string/pointer from each layer, they layer uniformly.
No additional property needs a separate mechanism; membership was the single
exception (`LoadRosterScoped`'s authority-rule keying on block presence,
not adapter contents).

## Position under Path C — hermes-1
I signed (a) in signoff-hermes-1.md (§2.1 decisive experiment: yes, a committed
values-only mechanism would weaken (c)'s case; my trigger requires it).
Path C rejects the premise: membership is a property, the deck file is a child
class overriding properties, and the overlay's separate authority mechanism is
unnecessary. Under C:
- I will BUILD the adapter=="inherit" mechanism (measured above) as a standard-track
  fix to D-A, not as a new overlay contract.
- The three defects (D-A, D-B, D-C) still stand and must ship independently.
- No mass conversion; no implicit reinterpretation; migration remains attended,
  per-deck, git-first.
- If C has a defect I see: it does not solve D-A's confirmation-text bug
  (`"adds a member"` when it replaces); the adapter=="inherit" fix must include
  truthful before/after set description in the gate (`roster_set`). I record this
  plainly — it is a technical correction to C's engineering, not a challenge to
  the binding direction.

## Migration — how affected decks move without silent quorum change
With adapter=="inherit" as the explicit discriminator:
1. No existing deck changes by default (byte-for-byte compatibility).
2. A deck that wants value-only inheritance edits one `[roster.<id>]` block:
   `adapter = "inherit"` + value fields; collaboration-visible.
3. A deck that wants full inheritance without value overrides deletes its
   `[roster.*]` blocks and relies on the machine layer (`Inherited = true`).
4. Migration is attended, per-deck, recorded as deliberate-vs-stale in a registry
   (the registry instrument is missing; it is a separate `standard` idea).
5. The 37 full-declared decks that omit `zcode-1` stay unchanged; no silent
   quorum expansion. If the owner wants them on inheritance, it is an explicit
   `roster sync --adopt-inherited` (existing vocabulary) per deck, not a bulk
   reinterpretation.

## What I would sign
1. **E1 is confirmed.** The adapter=="inherit" rule works; D-A's fix uses it; no
   silent quorum change occurs for unmarked files.
2. **E2 is confirmed.** All `[agents.*]` properties already layer per field;
   membership was the single exception.
3. **Position: build under C.** Fix D-A/D-B/D-C independently; ship the
   value-only separation via adapter=="inherit" (not a new overlay contract);
   migration remains attended; trigger remains "≥2 real instances" but now
   includes the value-override case (my amended trigger from round 2).
4. **Recorded defect of C's engineering:** D-A's confirmation text must describe
   replacement, not addition, when adapter=="inherit" is used; the resolver's
   truthfulness is not solved by the inheritance model change. I sign the
   experiment and the build, with that correction noted.
5. **Not signed:** any implicit reinterpretation of unmodified `[roster.<id>]`
   blocks; any overlay contract change before the gesture fixes land; any
   mass conversion; any closure of this idea before `opencode-1` completes its
   round-3 artifact (incomplete participation, not consent).

Tags: PRIMARY (measured runs on /tmp/parley-exp binary, test output, file reads,
code inspection of internal/config/runtime.go); SECONDARY (fleet census shape
— no adapter=="inherit" instances found in 5 tracked deck files; the 37 count
from consensus.md is unverified and not cited here as a gate).
Shared working tree: never edited. All mutations in /tmp/parley-exp only.
