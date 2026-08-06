---
idea: roster-operations-standard
phase: 8 — re-review
agent: kimi-1
round: 5
date: 2026-08-06
reviewed-commit: 8c8a8f1
verdict: FINDINGS
---
# Re-review round 5 — kimi-1

## Verdict

FINDINGS — one MINOR. The R4-M1 fix is verified behaviorally on the surface the finding
named: with all four layers in conflict, `roster show` state, `--explain` provenance, and
the `roster set --state` warning now AGREE — `active` is attributed to the membership
authority on every authority branch (deck TOML, legacy §2, inherited machine), a write to
the authority is no longer reported as masked, and a genuinely ineffective machine-file
write still warns, truthfully. The membership gate now fires on real flips (retire,
revive — both observed refused without `--confirm-breaking`) and not on no-ops (observed
writing clean). Both new tests FAIL when their fix is reverted, proven by overlay.

But the fix overshoots on one surface it was not aimed at: the `--explain` override is
unconditional, so at `--scope machine` it attributes the machine file's own `active`
value to the deck authority — the same provenance-contradicts-displayed-state defect
class as R4-M1, re-introduced by its own fix on a secondary scope. R5-m1 below.

**§15.1 ownership disclosure.** R4-m1 (skill 2.5.1) was my round-4 finding; the brief
states it is being released at the Phase 8 close, so I record no code verdict on it here
and do not re-file it. R5-m1 is a fresh claim about code that first exists at 8c8a8f1,
verified this session; I own it and issue no self-verdict on it. FIXED labels below
adjudicate the implementer's cycle-5 claims, not my own prior positions. All claims
PRIMARY (read, run, or executed by me this session) unless tagged otherwise.

**Method.** Binary built from 8c8a8f1 to `/tmp/parley-r5`; for the before/after
comparison, a848e67 extracted read-only via `git archive` to `/tmp/a848` and built to
`/tmp/parley-r4`. Scratch decks under `/tmp/kr5*` with `PARLEY_HOME` isolation. Reversion
sensitivity proven with `go test -overlay` from `/tmp` — the reviewed tree stayed
byte-identical: `git status --short` shows only the untracked `review/round-05/`
directory this file lands in. No git write commands used.

**Mandated command, run by me at 8c8a8f1.** `go build ./... && go test ./...` exited 0.
Exact output:

```text
BUILD_OK
?   	parley-deck-cli/cmd/parley	[no test files]
ok  	parley-deck-cli/internal/acp	(cached)
ok  	parley-deck-cli/internal/agents	(cached)
ok  	parley-deck-cli/internal/app	(cached)
ok  	parley-deck-cli/internal/config	(cached)
ok  	parley-deck-cli/internal/consensus	(cached)
ok  	parley-deck-cli/internal/driver	(cached)
ok  	parley-deck-cli/internal/fsutil	(cached)
ok  	parley-deck-cli/internal/hitl	(cached)
ok  	parley-deck-cli/internal/loop	(cached)
ok  	parley-deck-cli/internal/pipeline	(cached)
ok  	parley-deck-cli/internal/procctl	(cached)
ok  	parley-deck-cli/internal/protocol	(cached)
ok  	parley-deck-cli/internal/repomap	(cached)
ok  	parley-deck-cli/internal/retro	(cached)
ok  	parley-deck-cli/internal/runaction	(cached)
ok  	parley-deck-cli/internal/runcontrol	(cached)
ok  	parley-deck-cli/internal/runmanifest	(cached)
ok  	parley-deck-cli/internal/runner	(cached)
ok  	parley-deck-cli/internal/runplan	(cached)
ok  	parley-deck-cli/internal/runstate	(cached)
ok  	parley-deck-cli/internal/sessionstore	(cached)
ok  	parley-deck-cli/internal/steer	(cached)
ok  	parley-deck-cli/internal/store	(cached)
ok  	parley-deck-cli/internal/track	(cached)
ok  	parley-deck-cli/internal/tui	0.237s
TEST_EXIT=0
```

Mostly cached, so I also ran the suite uncached (`go test -count=1 ./...`, stricter than
mandated): exit 0, all 26 packages ok, 0 failures — including `internal/runner` 9.702s
(the durable-kill case that fails only in codex-1's sandbox passes here, as in my rounds
3–4). `go vet ./...` clean.

## Round-04 findings: fixed or not

### [MAJOR, codex-1 R4-M1] effective-state provenance and masking — FIXED on the deck scope; machine-scope regression filed as R5-m1

PRIMARY — Four-layer fixture with conflicting `active` values, every layer live at once:
deck commits `claude-1 active = true` and `kimi-1` (no `active` key); machine file sets
`claude-1 active = false`; `agents.local.toml` sets `kimi-1 active = false`;
`$PARLEY_HEADLESS_AGENT_CONFIG` sets both `active = false`. At 8c8a8f1:

- `roster show --json` (deck scope): `{'claude-1': 'active', 'kimi-1': 'active'}` — all
  three value-layer retirements discarded. (Values still layer, correctly: claude-1's
  model is `machine-claude-model` from the machine file.)
- `roster show --explain claude-1`:

  ```text
  claude-1 — membership from parley-deck/agents.toml
  active         active                   parley-deck/agents.toml
  ```

  Round 4 printed the env file here, asserting the opposite. `kimi-1` identical. State
  and provenance now agree, via `layers["active"] = src` from `config.RosterStateSource`
  (roster_view.go:178-180; runtime.go:1174-1179 returns `LoadRosterScoped`'s
  `scope.Source` — the same object that decided the state).
- Write surface, codex-1's exact round-4 scenario (env asserts `active = false`, write
  `--state active` to the authority): after retiring claude-1, `roster set claude-1
  --state active --yes --confirm-breaking` wrote with NO masked warning; `roster show`
  then reported `active` and `--explain` attributed it to the deck. Round 4 emitted the
  false `MASKED — PARLEY_HEADLESS_AGENT_CONFIG:...` warning on this write; the warning's
  substance (effective value did not change) was false then and is correctly absent now.
- TRUE-positive control: `roster set kimi-1 --scope machine --state inactive --yes
  --confirm-breaking` (machine is the only non-authority file `roster set` can write —
  roster_set.go:142-158) warned:

  ```text
  warning: active = "false" is MASKED — parley-deck/agents.toml sets it at a higher layer, so the effective value did not change.
  ```

  and deck-scope `roster show` kept `kimi-1 active` — the warning's substance is true,
  and "higher layer" is accurate for every reachable case (the deck outranks the machine
  file in `configLayers`, runtime.go:385-391). The masking branch compares the target
  path against the authority's path (roster_set.go:101-116), not the raw layer stack.
- Legacy authority branch: deck with only a §2 table (`claude-1` active, `kimi-1`
  `(inactive)`) plus the env file asserting both retired → `roster show` reported
  `{'claude-1': 'active', 'kimi-1': 'inactive'}` and `--explain` printed
  `active ... COOPERATION.md §2` for both. Inherited-machine branch: no deck roster,
  machine file authoritative → machine-scope `--state inactive` write to the machine
  file produced NO masked warning and the effective state flipped. All three authority
  branches attribute consistently.

### [Gate correctness, cycle 5] no-op gating — verified both directions

PRIMARY — No-op: `roster set kimi-1 --state active --yes` (kimi-1's deck block has no
`active` key; absence already means active) produced `+ active = true`, wrote, exit 0 —
no `--confirm-breaking` demand, no warning. Round 4's defect was exactly this write
demanding the second confirmation.

PRIMARY — Real flips still gated, both directions:
`roster set claude-1 --state inactive --yes` → refused, exit 2, "this retires a roster
member"; with `--confirm-breaking` → wrote, no masked warning, show/explain flipped to
`inactive` attributed to the deck. Then `roster set claude-1 --state active --yes` →
refused, exit 2, "this reactivates a retired roster member"; with `--confirm-breaking` →
wrote, state back to `active`. The gate keys on the prior state in the file being edited
(`priorActiveIn`, roster_set.go:312-326; conditions at roster_set.go:296-301), not on the
presence of a diff line.

### [MINOR, all three] skill 2.5.1 — carried; not re-filed

PRIMARY — The brief states the release happens at the Phase 8 close. Per instructions
this is not re-filed as a code finding; nothing in 8c8a8f1's diff touches it. RECALL —
its round-4 state (skill repo unchanged, IMPLEMENTATION.md claiming release) is recorded
in all three round-4 reviews; the close must still perform it.

## New findings (by severity, or "none")

### [MINOR] R5-m1 — cycle 5's provenance override misattributes `active` at `--scope machine`

PRIMARY — I own this finding and issue no self-verdict. Same fixture as above. The
machine-scope view deliberately answers about the machine file alone: header reads
`membership from /tmp/kr5/home/agents.toml`, and the displayed state comes from the
machine file (`rosterScopeFor` machine branch, roster_view.go:33-46 — no authority
rewrite). But the cycle-5 override applies unconditionally, pulling the DECK-scoped
authority (`config.RosterStateSource(root)`, roster_view.go:178-180):

```text
$ parley roster show --scope machine --explain claude-1        # 8c8a8f1
claude-1 — membership from /tmp/kr5/home/agents.toml
active         inactive                 parley-deck/agents.toml
```

The displayed `inactive` is the machine file's own value; the SET BY cell names a file
that (a) asserts the opposite (`active = true`) and (b) plays no role in this scope's
view. It also contradicts the header on the same screen. This is R4-M1's exact defect
shape — provenance contradicting the state shown, naming a layer asserting the
opposite — moved one scope over by the fix itself.

PRIMARY — Regression, not pre-existing: the a848e67 binary on the identical fixture
prints `active inactive ~/.parley/agents.toml` — self-consistent, because
`RosterFieldSourcesScoped(root, agent, machineOnly=true)` (roster_view.go:165) scans only
the machine layer. Conditions to reproduce: a deck roster exists (authority ≠ machine)
AND the machine file sets `active` for the member AND `--scope machine --explain`. The
legacy-authority variant diverges identically (`COOPERATION.md §2` named for the machine
file's value).

PRIMARY — Severity MINOR, not MAJOR: diagnostics-only, one cell on a secondary scope;
the machine-scope `roster show` STATE is correct (`{'claude-1': 'inactive', ...}`), the
masked warning on machine writes is correct (verified above), and state resolution,
quorum, and the gate are untouched. Suggested fix: skip the override when
`opts.scope == "machine"` — the machine-scoped `RosterFieldSourcesScoped` attribution is
already correct there — or source the override from the scope object actually in view.

### Anything else cycle 5 broke — nothing found

PRIMARY — Full uncached suite green, `go vet` clean. The cycle-5 diff touches three
source files; each hunk is exercised above or by the suite. `membershipChange`'s
signature change has exactly one call site (roster_set.go:67), updated; the pre-existing
`TestMembershipGateCatchesNewBlockWrittenWithAnyField` still passes natively with the new
argument. Deck-scope `--all` value layering, machine-scope `--all` scoping (the
`machine-claude-model` row above is machine-sourced in both scopes), and the
inherited/legacy branches all behaved as at a848e67.

## Test-quality assessment

PRIMARY — Both new tests run natively at 8c8a8f1 (PASS), then against overlaid reverts of
their own fixes from `/tmp`:

- `TestActiveProvenanceAndMaskingFollowTheAuthority` — REAL for the masking half.
  Overlaying `rosterFieldMaskedBy` without its authority branch (the round-4 fall-through
  to the raw layer stack) makes it FAIL, `go test` exit 1, with round 4's exact
  signature: `write to the authority reported as masked by
  "PARLEY_HEADLESS_AGENT_CONFIG:.../env.toml"`. Coverage note, below finding threshold:
  the test asserts `RosterStateSource` and the masking path directly but never calls
  `rosterExplain` — overlaying away the roster_view.go override (the user-facing explain
  wiring) leaves the suite GREEN (verified, exit 0), so the explain surface is guarded
  only behaviorally, by the runs above. That gap is exactly how R5-m1 shipped; a
  committed test driving `rosterExplain` at both scopes would have caught it.
- `TestMembershipGateIgnoresNoOpStateWrites` — REAL. Overlaying `membershipChange` to
  ignore `priorActive` (round-4 semantics, signature kept so it compiles) makes it FAIL,
  exit 1: `writing active=true to an already-active member gated as "reactivates a
  retired roster member"`. Both real flips (retire, revive) are also asserted in-test and
  were confirmed end-to-end at the CLI above.

No tautologies: each test constructs the defect fixture and asserts the
operator-visible outcome; each fails when its fix is reverted.
