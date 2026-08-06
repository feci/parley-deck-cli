---
idea: roster-operations-standard
phase: 8 — re-review
agent: codex-1
round: 6
date: 2026-08-06
reviewed-commit: d312c53
verdict: CLEAN
---
# Re-review round 6 — codex-1

## Verdict

PRIMARY — CLEAN. No new code finding emerged from the cycle-6 delta, the two-scope CLI fixture,
the state-write diagnostics, the revert-sensitivity experiment, or inspection of the changed
source and test.

PRIMARY — R5-M1 is a claim codex-1 first filed and therefore owns. Under §15.1, codex-1 issues no
`CONFIRMED`/`WRONG` self-verdict on it. For Phase 8 progression, codex-1 withdraws R5-M1 as an
active finding at `d312c53`; that withdrawal is a review position based on the raw evidence below,
not a self-verification verdict.

PRIMARY — I ran the mandated command exactly as `go build ./... && go test ./...`; it exited 1.
`go build ./...` was silent and the shell advanced to `go test ./...`. The exact combined output
was:

```text
?   	parley-deck-cli/cmd/parley	[no test files]
ok  	parley-deck-cli/internal/acp	(cached)
ok  	parley-deck-cli/internal/agents	(cached)
ok  	parley-deck-cli/internal/app	33.617s
ok  	parley-deck-cli/internal/config	(cached)
ok  	parley-deck-cli/internal/consensus	(cached)
ok  	parley-deck-cli/internal/driver	1.075s
ok  	parley-deck-cli/internal/fsutil	(cached)
ok  	parley-deck-cli/internal/hitl	(cached)
ok  	parley-deck-cli/internal/loop	(cached)
ok  	parley-deck-cli/internal/pipeline	(cached)
ok  	parley-deck-cli/internal/procctl	0.274s
ok  	parley-deck-cli/internal/protocol	(cached)
ok  	parley-deck-cli/internal/repomap	(cached)
ok  	parley-deck-cli/internal/retro	(cached)
ok  	parley-deck-cli/internal/runaction	(cached)
ok  	parley-deck-cli/internal/runcontrol	(cached)
ok  	parley-deck-cli/internal/runmanifest	(cached)
--- FAIL: TestDurableKillEndToEndRealProcess (0.02s)
    durablekill_test.go:116: a live attributed process should be killed, got {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001 Message:process verification failed (no recorded boot id); not killed}
FAIL
FAIL	parley-deck-cli/internal/runner	7.577s
ok  	parley-deck-cli/internal/runplan	(cached)
ok  	parley-deck-cli/internal/runstate	(cached)
ok  	parley-deck-cli/internal/sessionstore	(cached)
ok  	parley-deck-cli/internal/steer	(cached)
ok  	parley-deck-cli/internal/store	(cached)
ok  	parley-deck-cli/internal/track	(cached)
ok  	parley-deck-cli/internal/tui	0.424s
FAIL
```

PRIMARY — `git diff --name-only 8c8a8f1 d312c53 -- internal/runner` printed no path. The only
failure above is in a package untouched by cycle 6; `internal/app`, which contains the changed
call site and new test, passed in 33.617s.

## Round-05 finding: fixed or not

PRIMARY — Review disposition: withdrawn after fix-up; R5-M1 is no longer carried. The source
delta replaces the unscoped `config.RosterStateSource(root)` override in `rosterExplain` with the
already-resolved `scope.Source`. The post-write `active` branch now calls
`config.RosterStateSourceForTarget(root, target)`, which recognizes the central agents file as the
machine roster authority.

PRIMARY — I built `d312c53` and created a scratch fixture through the production CLI. Its committed
deck roster declared `claude-1 active = false` with `model = "deck-model"`; its isolated central
machine roster declared the same agent `active = true` with `model = "machine-model"`.

PRIMARY — In deck scope, `roster show` reported `inactive`; the matching explain output's
membership header and `active` provenance named the committed deck file in the same output:

```text
claude-1     claude     inactive yes       deck-model             unknown        unknown       max      balanced yes  inactive,metadata-unknown

claude-1 — membership from parley-deck/agents.toml
active         inactive                 parley-deck/agents.toml
```

PRIMARY — In machine scope, `roster show` reported the opposite state, `active`; the matching
explain output's membership header and `active` provenance both named the isolated central file:

```text
claude-1     claude     active   yes       machine-model          unknown        unknown       max      balanced yes  metadata-unknown

claude-1 — membership from /private/tmp/parley-r6-codex.rzgFjb/fixture/home/agents.toml
active         active                   /private/tmp/parley-r6-codex.rzgFjb/fixture/home/agents.toml
```

PRIMARY — I then applied real state flips in both scopes. The machine write changed `active =
true` to `false`; the deck write changed `active = false` to `true`. Both exited 0, wrote their
respective authority, and emitted no `MASKED` warning. Post-write `roster show` and `--explain`
reported `inactive` from the central file in machine scope and `active` from
`parley-deck/agents.toml` in deck scope.

PRIMARY — I repeated both state flips while `$PARLEY_HEADLESS_AGENT_CONFIG` asserted an opposing
`active` value. Neither state write emitted a masking warning, because `active` follows the roster
authority for the explicitly requested scope. As a positive control, a deck-scope `--model
deck-model-2` write under the same env layer did emit the real masking warning:

```text
warning: model = "deck-model-2" is MASKED — PARLEY_HEADLESS_AGENT_CONFIG:/private/tmp/parley-r6-codex.rzgFjb/envgen/parley-deck/agents.toml sets it at a higher layer, so the effective value did not change.
  (status `masked-by-env`; see `parley roster show --explain claude-1`)
```

PRIMARY — The new test passes natively at `d312c53`:

```text
=== RUN   TestActiveProvenanceIsScopeAware
--- PASS: TestActiveProvenanceIsScopeAware (1.52s)
PASS
ok  	parley-deck-cli/internal/app	1.787s
```

PRIMARY — I extracted `d312c53` to an isolated temporary tree, replaced only
`internal/app/roster_view.go` with its `8c8a8f1` version, kept the new test, and reran that exact
test with an isolated Go build cache. It failed with the round-05 machine-scope contradiction:

```text
=== RUN   TestActiveProvenanceIsScopeAware
    roster_cycle2_test.go:369: machine scope: provenance "active         active                   parley-deck/agents.toml" contradicts its own membership header "/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/TestActiveProvenanceIsScopeAware2627708075/002/agents.toml"
--- FAIL: TestActiveProvenanceIsScopeAware (0.84s)
FAIL
FAIL	parley-deck-cli/internal/app	1.163s
FAIL
```

PRIMARY — Thus the committed test is sensitive to reversion of the production explain fix rather
than merely restating a helper result. The production CLI probes separately exercised the
target-aware masking change, including both state scopes and a true-positive layered-field
control.

## New findings (by severity, or "none")

PRIMARY — none.

SECONDARY — Per the operator brief, the skill release is handled at Phase 8 close and is not a
code finding in this re-review.
