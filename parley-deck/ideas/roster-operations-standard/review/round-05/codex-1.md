---
idea: roster-operations-standard
phase: 8 — re-review
agent: codex-1
round: 5
date: 2026-08-06
reviewed-commit: 8c8a8f1
verdict: FINDINGS
---
# Re-review round 5 — codex-1

## Verdict

PRIMARY — FINDINGS. The requested deck-scope fix and the state-flip gate work behaviorally, but
cycle 5 regresses the public machine-scoped provenance surface: `roster show --scope machine`
resolves state from the machine roster while `--explain` says the same state came from the deck.
The same unscoped decision also makes an applied machine-scoped state write say its effective
value did not change when the machine-scoped row shows that it did. This is R5-M1 below.

PRIMARY — codex-1 owns R5-M1 and therefore issues no §15.1 `CONFIRMED`/`WRONG` verdict on that
new finding. `CONFIRMED FIXED` below adjudicates the cycle-5 implementation claims first made by
the implementer, not codex-1's round-04 finding.

PRIMARY — The mandated command was run exactly as `go build ./... && go test ./...`; it exited 1.
`go build ./...` succeeded silently, then `go test ./...` failed in the environment-dependent
durable-kill case. The exact combined output was:

```text
?   	parley-deck-cli/cmd/parley	[no test files]
ok  	parley-deck-cli/internal/acp	(cached)
ok  	parley-deck-cli/internal/agents	(cached)
ok  	parley-deck-cli/internal/app	30.337s
ok  	parley-deck-cli/internal/config	(cached)
ok  	parley-deck-cli/internal/consensus	(cached)
ok  	parley-deck-cli/internal/driver	0.939s
ok  	parley-deck-cli/internal/fsutil	(cached)
ok  	parley-deck-cli/internal/hitl	(cached)
ok  	parley-deck-cli/internal/loop	(cached)
ok  	parley-deck-cli/internal/pipeline	(cached)
ok  	parley-deck-cli/internal/procctl	0.253s
ok  	parley-deck-cli/internal/protocol	(cached)
ok  	parley-deck-cli/internal/repomap	(cached)
ok  	parley-deck-cli/internal/retro	(cached)
ok  	parley-deck-cli/internal/runaction	(cached)
ok  	parley-deck-cli/internal/runcontrol	(cached)
ok  	parley-deck-cli/internal/runmanifest	(cached)
--- FAIL: TestDurableKillEndToEndRealProcess (0.03s)
    durablekill_test.go:116: a live attributed process should be killed, got {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001 Message:process verification failed (no recorded boot id); not killed}
FAIL
FAIL	parley-deck-cli/internal/runner	7.413s
ok  	parley-deck-cli/internal/runplan	(cached)
ok  	parley-deck-cli/internal/runstate	(cached)
ok  	parley-deck-cli/internal/sessionstore	(cached)
ok  	parley-deck-cli/internal/steer	(cached)
ok  	parley-deck-cli/internal/store	(cached)
ok  	parley-deck-cli/internal/track	(cached)
ok  	parley-deck-cli/internal/tui	0.380s
FAIL
```

PRIMARY — `git diff --name-only a848e67 8c8a8f1 -- internal/runner` printed no path; the reviewed
delta does not touch the failing package.

## Round-04 findings: fixed or not

### 1. [MAJOR, codex-1 R4-M1] authority-state provenance and false masking

PRIMARY — The implementer's claim is CONFIRMED FIXED for the requested deck scope. I built the
8c8a8f1 CLI and used a scratch deck with two members and all four layers in conflict: the deck
made `claude-1` active and `kimi-1` inactive, while the machine file,
`parley-deck/agents.local.toml`, and `$PARLEY_HEADLESS_AGENT_CONFIG` asserted the opposite.

PRIMARY — `roster show --json` kept the deck's states. The relevant output was:

```text
      "agent": "claude-1",
      "adapter": "claude",
      "state": "active",

      "agent": "kimi-1",
      "adapter": "kimi",
      "state": "inactive",
```

PRIMARY — In the same conflicting fixture, both `roster show --explain` calls named the deciding
deck authority and agreed with the rows:

```text
claude-1 — membership from parley-deck/agents.toml
active         active                   parley-deck/agents.toml

kimi-1 — membership from parley-deck/agents.toml
active         inactive                 parley-deck/agents.toml
```

PRIMARY — `roster set claude-1 --scope deck --state active --yes` on the already-active member
exited 0, added the formerly implicit `active = true`, required no `--confirm-breaking`, and
emitted no masking warning. The complete mutation output ended with:

```text
  + active = true

Wrote /tmp/parley-r5-codex-8c8a8f1/deck/parley-deck/agents.toml
```

PRIMARY — The true no-change case, `roster set kimi-1 --scope deck --state inactive --yes`, also
exited 0 with `already matches ... nothing to do` and no warning. Thus the deck-scope row,
provenance, and state-write diagnostic agree in both authority directions.

PRIMARY — The gate still fired on both real flips. Retiring active `claude-1` and reviving
retired `kimi-1` without `--confirm-breaking` each exited 2 and left the authority unchanged;
the relevant output was:

```text
roster set: this retires a roster member — a membership change, not a settings change.
Re-run with --confirm-breaking as well as --yes.

roster set: this reactivates a retired roster member — a membership change, not a settings change.
Re-run with --confirm-breaking as well as --yes.
```

PRIMARY — Repeating each flip with `--confirm-breaking` exited 0, wrote the deck file, emitted no
masking warning, and the following `--explain` reported the new state from
`parley-deck/agents.toml`.

PRIMARY — Both new tests pass at 8c8a8f1:

```text
=== RUN   TestActiveProvenanceAndMaskingFollowTheAuthority
--- PASS: TestActiveProvenanceAndMaskingFollowTheAuthority (0.00s)
=== RUN   TestMembershipGateIgnoresNoOpStateWrites
--- PASS: TestMembershipGateIgnoresNoOpStateWrites (0.00s)
PASS
ok  	parley-deck-cli/internal/app	0.160s
```

PRIMARY — Both tests are sensitive to reversion of their core helper logic. With temporary Go
overlays restoring the old general-layer masking rule and the old presence-based gate while
leaving the reviewed tree unchanged, both failed with the old signatures:

```text
roster_cycle2_test.go:321: write to the authority reported as masked by "PARLEY_HEADLESS_AGENT_CONFIG:<scratch>/env.toml"
roster_cycle2_test.go:328: writing active=true to an already-active member gated as "reactivates a retired roster member"
```

PRIMARY — The tests are not end-to-end guards for every fixed surface. Reverting only the
`roster_view.go` call that installs `RosterStateSource` left
`TestActiveProvenanceAndMaskingFollowTheAuthority` green, because the test calls the helper and
`rosterFieldMaskedBy` but never calls `rosterExplain` (`roster_cycle2_test.go:303-323`). Replacing
only `rosterSet`'s `priorActiveIn(...)` argument with a constant likewise left
`TestMembershipGateIgnoresNoOpStateWrites` green, because it calls `membershipChange` directly
(`roster_cycle2_test.go:326-335`). The behavioral CLI runs above cover those deck-scope call sites;
R5-M1 demonstrates the unguarded scope regression.

### 2. [MINOR, all three] skill 2.5.1 release

PRIMARY — The operator brief states that skill 2.5.1 is released separately at Phase 8 close and
explicitly says not to re-file it as a code finding. This review follows that direction and makes
no claim that the future release step has already occurred.

## New findings (by severity, or "none")

### [MAJOR] R5-M1 — PRIMARY — cycle 5 makes machine-scoped state provenance unscoped

PRIMARY — codex-1 owns this finding and issues no §15.1 truth verdict on it. The reproduced
fixture has a committed deck roster and a machine roster with the opposite `active` values. After
the machine authority had `claude-1 active = true` and the deck authority had
`claude-1 active = false`, the 8c8a8f1 binary's `roster show --scope machine --json` reported:

```text
      "agent": "claude-1",
      "adapter": "claude",
      "state": "active",
```

PRIMARY — The matching 8c8a8f1 command,
`roster show --scope machine --explain claude-1`, contradicted both its own membership header and
that row:

```text
claude-1 — membership from /tmp/parley-r5-codex-8c8a8f1/home/agents.toml
active         active                   parley-deck/agents.toml
```

PRIMARY — A binary built from the baseline a848e67 archive against the identical fixture named
the machine authority instead:

```text
claude-1 — membership from /tmp/parley-r5-codex-8c8a8f1/home/agents.toml
active         active                   ~/.parley/agents.toml
```

PRIMARY — This is a cycle-5 regression, not a pre-existing machine-scope result: the delta adds
the overwrite at `internal/app/roster_view.go:176-180`; it calls
`config.RosterStateSource(root)` even though the surrounding source/spec loaders are explicitly
scoped by `opts.scope == "machine"` at lines 165 and 175.

PRIMARY — The helper itself has no scope input and always calls the deck-aware
`LoadRosterScoped(root)` (`internal/config/runtime.go:1170-1179`). When a deck roster exists, that
returns the deck authority even while `rosterScopeFor(..., "machine")` resolved the displayed row
from the central machine file.

PRIMARY — The write diagnostic has the same missing-scope problem. Applying
`roster set claude-1 --scope machine --state active --yes --confirm-breaking` wrote the machine
authority and then printed:

```text
warning: active = "true" is MASKED — parley-deck/agents.toml sets it at a higher layer, so the effective value did not change.
  (status `masked-by-env`; see `parley roster show --explain claude-1`)
```

PRIMARY — Immediately afterward, `roster show --scope machine --json` reported `state: active`.
`rosterFieldMaskedBy` receives no command scope and uses the same unscoped helper at
`internal/app/roster_set.go:96-115`, so a command explicitly changing the machine authority is
described relative to the deck authority instead of the requested scope.

PRIMARY — Suggested fix: make state-source resolution accept the requested scope (or pass the
already-resolved authority through), and use that same scoped authority in both `rosterExplain`
and the post-write masking check. Add a production-surface test with opposite deck/machine states
that invokes `roster show --scope machine --explain` and `roster set --scope machine --state`;
the current helper-only test remains green if the explain call site is reverted and never
exercises machine scope.

PRIMARY — Inspection of the remaining five source/test files in `a848e67..8c8a8f1`, the
deck-scope behavioral matrix, the targeted tests, and the full mandated command produced no
additional cycle-5 finding.
