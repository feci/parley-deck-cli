---
agent: codex-1
idea: zcode-adapter
review-round: 2
date: 2026-08-19
reviewed-commit: 674b866
responding-to: [claude-1/review/round-01, codex-1/review/round-01, hermes-1/review/round-01, kimi-1/review/round-01]
---

## Position changes since review round 1

- My prompt-form MAJOR is resolved. The real parser accepts the one-element `--prompt=<value>`
  form for dash-leading, newline/quote-bearing, and flag-lookalike prompts, and both halves of the
  implementation are regression-sensitive.
- My manifest-coverage MAJOR is resolved. Removing the zcode registry entry makes the external
  `EXPECTED_TARGETS = 15` assertion fail `14 !== 15`. There is no non-circular derivation from the
  production registry or its plan; a separately maintained expected-name list would be another
  external oracle, not a derivation. The constant is the right minimal removal tripwire.
- The installer-path MINOR is resolved. The new test compares the registry value and the actual
  install destination with the independently documented `.zcode/skills` path; it does not merely
  restate the implementation.
- My model-bindability CRITICAL is fixed for `roster show`, including its `--explain` trailer, but
  the fail-closed capability does not reach the adapter inventory and has no regression test. I
  narrow that remaining issue to the MAJOR below.
- The two repositories are now clean, immutable revisions (`674b866` for the CLI and `75db419` for
  the skill). `IMPLEMENTATION.md` still does not record those current revisions; I leave that
  audit-document cleanup as a deferred follow-up because it is outside this fix-only surface.

## Responses to other reviewers

### @claude-1

The prompt, count-tripwire, destination, and Notes fixes hold under attack. I do not concur that
`Spec.NoModelBinding` is fully fail-closed yet: `roster show` is correct, but `agents list` ignores
the capability and the capability can be deleted without making the current Go tests fail. Because
the standard-track two-cycle cap is exhausted, the MAJOR below requires escalation, not an
unapproved third fix-up cycle.

### @codex-1

I retain the round-1 model attack as the decisive regression case. Against `600652b` it produced a
bound model and contradictory trailer; against `674b866` the same roster command stays unknown.
The prompt and circular-oracle findings are withdrawn as fixed. The earlier repository-anchoring
finding is materially reduced by the clean CLI and skill commits and is deferred as documentation
cleanup rather than refiled in this fix-only round.

### @hermes-1

I disagree with both round-1 MAJORs as stated. The `unknown agent zcode` result came from released
Parley 1.44.0, not this branch build. The branch build discovers zcode and starts the full probe;
this sandbox stops the real CLI later at Unix-socket creation with `EPERM`. A prompt containing
`--mode build` remains one argv element and cannot replace the later `--mode yolo`; the real
equals-form parser and the exact-argv fake both confirm that. I also reverse the view that the
derived manifest count was stronger: the isolated target-removal mutation proves it was circular.

### @kimi-1

I concur with the addendum that `NoModelBinding` fixes the roster row and that equals-form prompt
substitution fixes the real parser edge. The destination test requested by the MINOR is meaningful
and mutation-sensitive. The prior open question about `agents list` is now answered adversely:
under the hostile override it displays a concrete zcode model and unsupported `--model` argv even
while `roster show` correctly refuses both.

## Updated findings

### [MAJOR] `NoModelBinding` is roster-only and the CRITICAL fix has no regression lock

With a scratch `PARLEY_HOME` containing both a model and hostile `headless_args`, `roster show
--scope machine --explain zcode-1` correctly reports `model unknown`, `model-unbound`, and a
non-contradictory trailer for separate-token, equals-token, and literal `--model` shapes. However,
`agents list` on the same resolved configuration prints:

```text
zcode ... MODEL adversarial/provider-model
headless: ... --model adversarial/provider-model
```

The inventory path at `internal/agents/discover.go` renders `result.Model` and
`ResolveLaunchArgs()` directly, neither of which consults `NoModelBinding`. It therefore presents
a concrete model and an argv that the real zcode parser rejects, despite the new adapter capability
saying that binding is impossible.

The test gap is independently reproducible: in an isolated copy I removed the field, both
`EffectiveModel`/`EffectiveEffort` guards, and the zcode initializer bit. The original CRITICAL
immediately returned under the hostile roster command, but `go test ./internal/agents
./internal/app` still passed. No test owns the capability fix.

Why it matters: the fail-closed claim is operator-visible on two first-class inventory surfaces,
and a future deletion of the critical guard can silently restore the contradictory roster state.
Concrete fix: reject model/effort-bearing zcode overrides during config resolution, or make
launch-argument resolution and inventory rendering honor `NoModelBinding` with `unknown` plus an
unsupported-override diagnostic. Add a config-backed regression test that exercises both `roster
show --explain` and `agents list`; deleting `NoModelBinding` must fail it.

This is a fresh MAJOR on the fixed surface. The standard-track fix-up cap is reached, so it forces
operator escalation rather than a third automatic fix-up cycle.

### Resolved fixed surfaces

- `--prompt={prompt}` and embedded substitution: verified; implementation-only reversions fail
  `TestZcodeSpecCarriesNoModelOrEffortPlaceholder` and/or
  `TestZcodeFullVerifyAcceptsHonestLaunch`.
- `EXPECTED_TARGETS = 15`: verified; removing zcode fails the tripwire before the explicit zcode
  test also reports the missing target.
- Pinned `.zcode/skills` destination: verified; mutating the implementation to `.zcode/skill`
  fails the independent literal assertion.

## Refutation attempts

1. **Exact reviewed state.** [PRIMARY] Reviewed clean archives of CLI `674b866` and skill
   `75db419`; the shared CLI tree acquired another reviewer's untracked root file during the run,
   so all branch-certifying Go commands and mutations used the immutable archive.
2. **Hostile model configs.** [PRIMARY] Tried three scratch machine configs: separate
   `--model {model}`/`--effort {effort}`, equals-form `--model={model}`/`--reasoning={effort}`, and
   literal `-m literal/model`/`--thinking literal-effort`, each with conflicting roster values.
   Every `roster show --scope machine --explain zcode-1` result was `unknown` /
   `model-unbound,effort-unknown`; every trailer agreed with the row.
3. **Other model display.** [PRIMARY] `agents list` with the first hostile config displayed
   `adversarial/provider-model` and resolved `--model adversarial/provider-model`, producing the
   MAJOR above.
4. **Model-fix mutation.** [PRIMARY] Removed only the `NoModelBinding` implementation in an
   isolated copy. The hostile roster row again became `adversarial/roster-model` with no
   `model-unbound` and contradicted the trailer, while `go test ./internal/agents ./internal/app`
   exited 0. The critical fix is not test-owned.
5. **Real zcode parser.** [PRIMARY] Ran equals-form leading-dash, multiline/quoted, and
   `--mode build --cwd /etc` lookalike prompts with a root containing spaces and an explicit later
   `--mode yolo`. All passed option parsing and reached ZCode's local broker; this sandbox then
   denied its Unix-socket `listen` with `EPERM`. None produced the former ambiguous-prompt error.
6. **Prompt-fix mutations.** [PRIMARY] Reverted the equals-form spec and embedded substitution
   together: the exact-argv, hostile-prompt, and honest-full-verify tests failed. Reverted only the
   runner's embedded substitution: the honest full-verify test failed. The new tests are live.
7. **Manifest removal.** [PRIMARY] Removed zcode from `lib/installer.js` in an isolated skill copy.
   `node --test test/manifest-coverage.test.js` failed the external tripwire with `14 !== 15` and
   separately failed `installer registry has no zcode target`.
8. **Destination mutation.** [PRIMARY] Changed zcode's `skillDir` to `.zcode/skill` in an isolated
   copy. The focused manifest suite failed with expected `.zcode/skills`; the assertion is not
   implementation-derived.
9. **Required gates.** [PRIMARY] On the clean `674b866` archive, `go build ./...` and `go vet ./...`
   passed. Exact `go test ./...` ran all packages but failed only the pre-existing
   `TestDurableKillEndToEndRealProcess`: this sandbox denies both `sysctl kern.boottime` and `ps`,
   leaving no boot ID. `go test ./... -skip '^TestDurableKillEndToEndRealProcess$'` passed every
   package, and all focused zcode tests passed. On clean skill commit `75db419`, `npm test` passed
   387 Node tests, 54 Python tests, and every manifest check.
10. **Required real full probe.** [PRIMARY] The branch-built `parley agents verify --full --agent
    zcode --yes` discovered `zcode-app-cli 3.7.7-13` and launched it, then failed only when this
    sandbox denied ZCode's local Unix-socket `listen` with `EPERM`. This is not the released-binary
    `unknown agent` failure, but a normal-host rerun is still required for an exit-0 acceptance.

## Open questions

- Operator decision required: authorize a bounded exception for the MAJOR above, accept/defer it
  explicitly, or stop the release. The protocol does not permit silently starting standard fix-up
  cycle 3.
- A non-sandboxed host still needs to produce green exit codes for exact `go test ./...` and the
  real zcode full probe; this environment cannot supply the macOS process identity or ZCode socket.
- Deferred follow-up: update `IMPLEMENTATION.md` frontmatter and validation evidence to name CLI
  `674b866` and skill `75db419` so the now-clean two-repository candidate is self-describing.
