---
idea: addon-bidding-scope
status: final
drafted-by: claude-1
date: 2026-08-17
track: standard
participants: [claude-1, codex-1, hermes-1, kimi-1]
rounds: 2
signoffs: [claude-1 ✅, codex-1 🟡, hermes-1 ✅, kimi-1 ✅]
---

# FINAL — `parley-bidding` becomes opt-in, and the reason is not the one we started with

## Decision

**`parley-bidding` stays in the package and stops installing by default.** Unanimous across four
model families, held through two rounds. Nobody proposed CUT; nobody proposed SPLIT as a first
choice.

The idea opened asking whether a jurisdiction-bound, Python-carrying business vertical belongs in a
vendor-neutral multi-agent-protocol package. It closes on something narrower and better evidenced:
**two defects that exist today, both reproduced independently three times.**

## The two results the decision rests on

Re-run by @codex-1, @hermes-1 and @kimi-1 on `parley-deck-skill 2.8.0` in fresh scratch targets,
all three reporting the same outcomes.

**1. The documented opt-out does not persist.** `install --no-addons` writes `addons: false` to the
core marker; the next flag-less `install --force` rewrites it to all five add-on names.
`selectedAddons()` (`lib/installer.js:892-901`) never consults the recorded selection, and
`marker.addons` is written *from* the current run (`:1958-1963`). **The opt-out is per-invocation.**
A user who declines today gets the add-on back on the next routine upgrade — so the README's advice
is a treadmill.

**2. The naive flip is worse than doing nothing.** After a default install, an `install --force`
under the new default rewrites the marker from five names to four; **the bidding directory stays on
disk**; `doctor` then reports `selected: false, status: valid-unselected`, problem *"installed but
not part of the recorded selection"*, overall `ok: false`, **exit 1**. A bare default change turns
every routine upgrade on every existing runtime red.

**A third, from @codex-1's round-1 experiment:** installing with `PATH=/definitely-no-python` exits 0
with `installOk: true` while `doctor` exits 1 — a default-installed add-on makes fleet health depend
on an interpreter the package otherwise does not require, in a package where two other add-ons
advertise zero runtime dependencies as a feature.

## What ships

1. **Selection persistence — owed unconditionally.** Make a flag-less run *read* the recorded
   selection instead of re-selecting everything. This is a defect under **every** option including
   KEEP: someone who wants default-on still wants `--no-addons` to mean something durable. Per
   @kimi-1, it ships **bundled** rather than separately — one code path is both the persistence fix
   and the grandfather.
2. **The default flip, with grandfathering.** Existing installs that recorded `parley-bidding` keep
   it across upgrades (@kimi-1's marker-consult, reached independently by @hermes-1 in round 2), plus
   @codex-1's one-time operator notice. @kimi-1 simplified @codex-1's design by dropping the
   "migrated legacy selection" provenance field.
3. **`--with` / `--without`, in the same changeset.** @codex-1 found that `--only` means core *plus
   only* the named add-ons, so a bidding user preserving the other four must enumerate all of them,
   and uninstall planning always begins with the core unit (`:991-1003`, `:649-669`) — there is no
   safe per-add-on removal. @kimi-1 verified this live and tightened its position to must-ship: the
   documented equivalent is `rm -rf` plus a five-name `--only`, *"not a thing a README should ask of
   anyone"*. Without `--with`, the flip would make the opt-in path worse than the opt-out path it
   removes.

## What does not ship

- **SPLIT.** Declined on @codex-1's criteria: no independent maintainership, no second jurisdiction
  pack, no demand for bidding without the protocol. Revisit when one appears. *(@claude-1's stated
  reason — that a second package doubles the all-channels release tax permanently — was wrong;
  @kimi-1 corrected it: independent cadence means the vertical pays that tax on its own releases.)*
- **CUT.** It works, it is gated, it handles no credentials, and the domain material is live in this
  workspace.
- **Any inference about `parley-design-check`.** @codex-1's formulation: **"The principle transfers;
  the conclusion does not."** This decision sets a question for every default add-on — *does its
  broad expected value justify model-routing surface, runtime constraints and upgrade behaviour for
  users who did not select it?* — and design-check must answer it separately. It shares exactly one
  property with bidding: size. Size was nobody's argument.

## Two corrections against the drafter, recorded because they change how to read this

**The framing that opened the idea was wrong.** `00-prompt.md` called the README `[!IMPORTANT]` a
fourth instance of *a printed rule binds only where enforcement lives*. @codex-1 struck it:
disclosure is not a rule purporting to enforce non-installation, the real default lives in
`selectedAddons()`, and calling a notice a failed gate is pattern-matching. @hermes-1 and @kimi-1
withdrew their uses of it; @claude-1 withdrew the framing. What survives is narrower — prose about
*permission* answers a different question from whether an unrequested skill enters an agent's
*routing surface*.

**The routing argument has no verified leg.** @hermes-1 and @claude-1 both argued that availability
is surfacing and surfacing is model attention. Neither measured it. @codex-1 asserted the weaker
premise — that at least one model-facing runtime surfaces installed skills — but ran no activation
experiment, and corrected the consensus draft at signoff for implying otherwise. **The premise is
UNVERIFIED and the behavioural effect is un-measured.** This is exactly why the decision rests on the
two reproduced defects instead, and FINAL must not let those defects retroactively validate the
asymmetry argument: they are different claims with different evidence.

## Method note

Unanimity here is **partly a shared prior** — the facilitator selected the recon framing and all four
then agreed with it — recorded per `COOPERATION.md:1356-1360`. The counter-evidence is specific: all
three non-facilitators independently attacked that framing, and @kimi-1 produced the fact that
reframed the idea and that the prompt did not contain. **A prompt that generated its own refutation
has not captured the room.**

@codex-1 answered the counterfactual most usefully: it would have said **KEEP after fixing
persistence** had installed skills not been surfaced until separately enabled, and had a missing
Python runtime not failed health.

## Implementation note — not authorised by this idea

The outcome is an installer change in `parley-deck-skill`. `COOPERATION.md` is untouched. Under this
deck's standing release discipline, shipping it means a version bump and a release across every
channel with independent per-channel verification. **That is the owner's decision, not this
consensus's.**
