---
idea: zcode-adapter
status: implemented
implementer: claude-1
started: 2026-08-19
completed: 2026-08-19
branch: parley-deck-cli#zcode-adapter-impl
head-commit: 9a6f8be
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

`zcode` is now a built-in adapter. `parley agents list` prints its `headless:` argv line,
`roster show` reports AUTO=yes, and the acceptance command from FINAL passes against the real
binary:

```
$ parley agents verify --full --agent zcode --yes
zcode: installed version=zcode-app-cli 3.7.7-13
probe dir: parley-deck/meta/runtime-probes/20260818T232832.862162000Z
zcode: headless probe passed
```

MODEL and EFFORT stay `unknown` / `model-unbound`, as designed.

## Implementation plan / checklist

- [x] `internal/agents/discover.go` — `Spec{ID: "zcode"}`, argv
      `--prompt {prompt} --mode yolo --cwd {root}`, `AutonomousWrite{Mode:"yolo",
      Args:["--mode","yolo"], Scope:""}`, `Model/Reasoning: CLIDefault`.
- [x] `internal/agents/modelmeta.go` — `"zai": "Zhipu AI"` producer (zcode emits `zai/<model>`,
      no hyphen; `z-ai`/`zhipu` did not match).
- [x] `internal/app/roster_view.go` — static `model source:` trailer on `--explain`.
- [x] `internal/agents/launchargs_test.go` — lock: no `{model}`/`{effort}` placeholder, exact
      argv, AutonomousWrite args present in argv, Scope not "workspace".
- [x] `internal/agents/modelmeta_test.go` — `zai/glm-5.3` → GLM / Zhipu AI.
- [x] `internal/app/zcode_verify_test.go` — @codex-1's two fake-zcode full-verify cases plus a
      trailer test.
- [x] `internal/agents/acp_specs_test.go`, `internal/app/roster_test.go` — deliberate presence locks.
- [x] Skill: `lib/installer.js` zcode target, `skills/parley-deck/SKILL.md` autonomous-write row,
      manifests rebuilt.
- [x] `~/.parley/agents.toml` — `[agents.zcode]` removed, roster `model`/`effort` removed,
      stale exit-0 note corrected.
- [x] Checks run: `go build ./...`, `go test ./...` (all green), `npm test` (386 pass / 0 fail),
      real `agents verify --full --agent zcode --yes`.

## Deviations from FINAL.md

**D-1 — `.zcode/skills` was verified, not assumed.** FINAL required a native installer target but
not a path. `~/.zcode/` has no user skills directory on this machine (all 14 skills come from the
plugin cache), so the path could not be read off the filesystem. It was taken from the literal
`.zcode/skills` present in the `zcode-app-cli` runtime bundle, corroborated by
`~/.zcode/cli/config.json` exposing `skills.roots`. Recorded because an unverified install path
would fail silently.

**D-2 — one extra file changed, in the same defect class the idea is about.**
`test/manifest-coverage.test.js` asserted `cores.length === 14` — a hardcoded target count whose
own comment claimed it was "derived from the installer rather than a hand-picked pair". Adding a
15th target failed it. It is now derived from the install result's own target set, so a future
target extends the check instead of breaking an arithmetic assertion nothing keeps in step. Bumping
14→15 would have preserved the defect.

**D-3 — `--explain` shipped in @kimi-1's static form**, as FINAL's Phase-5 default.
@codex-1's labelled live read is not implemented; its reservation stands on the record.

## Notes for reviewers

- The fake-zcode stub asserts **exact token order and `argc == 6`**, so a spec that appends or
  reorders an option zcode does not accept fails in test rather than in production.
- The help-exit-0 case is the important one: it proves full verification cannot be satisfied by a
  process that exits 0 while writing nothing. Exit status alone is never the acceptance signal.
- `agents list` shows a model for zcode when a config layer sets one; `roster show` still reports
  `model-unbound`. That is the contract working — the roster reports what the *launch* carries.
- Not done, deferred by FINAL: the `zcode app-server` (ZCode Protocol) binding route, and generic
  exit-0-with-no-artifact diagnosis.

## Fix-up cycle 1
status: complete
completed: 2026-08-19

### Fixes applied

- **@codex-1 CRITICAL — a config override could make zcode model-bound.** A wholesale
  `headless_args` override appending `--model {model}` made `roster show` display a model, drop
  `model-unbound`, and print the trailer's "never passed by parley" beside it. Fixed by making
  bindability an **adapter capability**: new `Spec.NoModelBinding`, set for zcode, and
  `EffectiveModel`/`EffectiveEffort` fail closed on it regardless of argv. Reproduced @codex-1's
  scratch attack and @kimi-1's variant against the fix: both now report `unknown` /
  `model-unbound`, and the trailer is true in that state. This also closes @kimi-1's MINOR-1.
- **@codex-1 MAJOR — the separate-token prompt form is rejected by zcode.**
  `zcode --prompt "-leading dash"` exits 1 with "Option '--prompt' argument is ambiguous";
  `--prompt=<value>` is accepted (measured on the real binary). The spec now ships
  `--prompt={prompt}`, and `buildAgentInvocation` substitutes placeholders **inside** an argv
  element without a shell. New `TestZcodeArgvSurvivesHostilePromptAndRoot` covers a leading dash,
  newlines, double and single quotes, a flag-lookalike prompt, and a root containing spaces.
- **@codex-1 MAJOR — D-2's derived count was circular, and my first fix was circular too.**
  Deriving from `result.actions` let a dropped target shrink both sides. Deriving from
  `installer.TARGETS` failed the same way one level up — proved by removing the zcode target in an
  isolated copy and watching the test stay green. **A target that disappears cannot be detected
  from inside the registry that lost it**, so the external constant is restored as
  `EXPECTED_TARGETS = 15` with the reasoning recorded. The reversion check now fails as it must.

### Fixes applied — cycle 2

- **@kimi-1 MINOR — nothing pinned zcode's installer destination.** New test asserts the registry
  entry, `skillDir === .zcode/skills`, and that a real install lands the core in
  `<home>/.zcode/skills/parley-deck` with its `SKILL.md`.
- **@kimi-1 NIT — the spec Notes now state why the equals form is required**, with the measured
  error text.

### Deviations from agreed fixes

- **@kimi-1 NIT (accepted, not actioned):** FINAL named `internal/app/app_test.go` as the home for
  the two fake-zcode cases; they live in a new `internal/app/zcode_verify_test.go`. Moving them
  would be churn with no behavioural difference; recorded rather than silently ignored.
- **@hermes-1's MAJOR findings are refuted, with measurements.** (1) "verify fails with unknown
  agent zcode" — @hermes-1 ran the *installed* `parley 1.44.0`, which predates this adapter; the
  branch build passes. That is also a facilitator error: the review brief did not say to build the
  branch. (2) "argv injection via `{prompt}`" — `args = append(args, prompt)` appends one element
  and Go `exec` does no shell parsing, so a prompt containing flags or newlines cannot split argv;
  measured with an argv spy (6 elements, prompt intact). Its `.zcode/skills` MINOR and both NITs
  stand.
- **@kimi-1 independently upgraded D-1's evidence**: `.zcode/skills` is documented as the user
  skill root by the vendor's own bundled guide inside `zcode-app-cli`
  (`vendor/packages/zcode-guide-plugin/skills/zcode-configuration-guide/SKILL.md`), not merely a
  string literal.

### Process error, recorded

**The tree moved while review round 1 was open.** Fix-up cycle 1 was applied while @kimi-1 was
still reviewing; @kimi-1 detected it, pinned both trees by mtime, re-ran the suites and filed an
addendum. That is the facilitator's error against the standing rule that the tree does not move
during an open review round.

## Fix-up cycle 3
status: complete
completed: 2026-08-19

**Budget exception:** the `standard` cap is 2 cycles. A third was applied under the owner's standing
autonomous authorisation and is recorded in
`inbox/claude-1-to-user_zcode-adapter_fixup-budget-exception.md`.

### Fixes applied

- **@codex-1 MAJOR — `NoModelBinding` was roster-only.** Cycle 1 fixed `roster show` but left
  `agents list` and the resolved argv rendering a config-supplied model, so the two surfaces
  contradicted each other and the inventory advertised a flag zcode rejects. Stripping moved into
  `ResolveLaunchArgs`, which every surface goes through: all three now agree and the bad flag never
  reaches the launch. Verified against @codex-1's exact attack.
- **@codex-1 MAJOR — no regression lock.** Added
  `TestNoModelBindingStripsConfigSuppliedFlagsEverywhere`: asserts the model value and both flags
  are absent from the resolved argv, both status bits are set, `EffectiveModel`/`EffectiveEffort`
  report unbound, the autonomous-write flag survives the stripping, and a bindable adapter is
  unaffected.

### Deferred follow-up (@codex-1, agreed)

`IMPLEMENTATION.md` frontmatter still records the first implementation commit rather than the
current head; audit-document cleanup, outside the fix-only surface.

### Incomplete participation

`kimi-1`'s review-round-2 process was killed before writing its artifact. Recorded as incomplete,
not consent. `standard` requires 2 reviewers; @codex-1 and @hermes-1 completed.
