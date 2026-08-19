---
agent: kimi-1
idea: zcode-adapter
review-round: 1
date: 2026-08-19
reviewed-commit: 600652b
---

## Summary

The implementation matches FINAL on every ratified point I could measure, and none of the
attacks below produced a CRITICAL or MAJOR. The argv is exactly
`--prompt {prompt} --mode yolo --cwd {root}`; substitution is exact-token with no shell, so
the prompt travels as one argv element; the roster contract holds — no config layer that sets
only the `model`/`effort` fields can make `roster show` display a model for zcode. The two
fake-zcode tests are not inert: removing the spec, corrupting the argv, or appending an
unsupported option all fail them in a scratch copy. The acceptance command passes against the
real binary (`zcode: installed version=zcode-app-cli 3.7.7-13`, probe passed, rc=0), and
`--agent zcode-1` fails rc=2 `unknown agent` as FINAL predicted. `.zcode/skills` (D-1) is
corroborated by the vendor's own documentation bundled inside `zcode-app-cli`
(`vendor/packages/zcode-guide-plugin/skills/zcode-configuration-guide/SKILL.md` documents
`~/.zcode/skills/` as the user skill root), not merely by a string literal.

Checks run on the reviewed state: `go build ./...` + `go vet ./...` clean, `go test ./...`
green (exit 0), `npm test` in the skill package 386 pass / 0 fail.

**Reviewed-state caveat:** BOTH trees moved while I was reviewing. (1) Skill package
(uncommitted working-tree state): the `manifest-coverage.test.js` I first diffed derived the
target count from `result.actions` (circular — see Findings); during the review the file was
revised (mtime 1787096914) to an `EXPECTED_TARGETS = 15` tripwire plus a registry/plan
comparison, with comments attributing the circularity catch to codex-1's round-1 MAJOR. I
re-ran `npm test` against the revised state (again 386/0) and assessed the revised form. The
skill side of this review covers the tree as of mtimes 1787096793 (`lib/installer.js`) /
1787096914 (`test/manifest-coverage.test.js`). (2) CLI repo working tree: uncommitted edits
to `discover.go`, `launchargs.go`, `runner.go`, `launchargs_test.go`, `zcode_verify_test.go`
landed at mtimes 1787096701-1787097235 — see the addendum at the bottom. All findings below
are stated against 600652b as reviewed; the addendum records what the in-flight edits change.
Pin both states before merge.

## Findings

### MINOR — the `--explain` trailer lies in the operator-override state

`modelSourceTrailer` (internal/app/roster_view.go:220) keys only on the adapter name, so it
prints unconditionally for zcode. Reproduced live in a scratch deck:

```toml
[agents.zcode]
model = "zai/glm-5.3"
headless_args = ["--prompt", "{prompt}", "--mode", "yolo", "--cwd", "{root}", "--model", "{model}"]
```

`roster show --explain zcode-1` then prints `model zai/glm-5.3` (bound, effective) and, two
lines later, the trailer: "never passed by parley (zcode has no --model flag), so the
effective model cannot be bound or observed from here". Both claims are false in that
configuration — parley IS passing the model and the row above observes it. Why it matters:
`--explain` is the provenance surface this idea created to be scrupulously honest; an
operator who overrides the argv (the documented escape hatch) gets internally contradictory
output exactly when they are most likely to be reading it. Concrete fix: gate the trailer on
the unbound condition, e.g. print it only when `row.Model == agents.Unknown` (or when
`EffectiveModel()` reports no model), so the trailer describes the state it claims to
describe. Rewording to "with the built-in argv, ..." is the weaker alternative.

### MINOR — nothing pins the zcode installer's resolved destination

D-1's own risk statement is "an unverified install path would fail silently", yet the skill
suite's destination-pinning test (`test/installer.test.js`, "resolves user and project target
paths") was not extended for zcode — it names codex, claude, agy, gemini, hermes, kimi, qwen,
droid, aionrs and stops. I verified the resolved dest by executing it
(`~/.zcode/skills/parley-deck`, `detected: true`), and the path is vendor-documented, but a
future typo in `skillDir` (`zcode/skills`, `.zcode/skill`) passes the entire suite: the
manifest-coverage test counts cores, it does not assert WHERE they land. Concrete fix: add
zcode to that existing test (user and project scope), one assertion each. Note the test file
was revised mid-review; check whether the current tree already added this before acting.

### NIT — fake-zcode cases live in `zcode_verify_test.go`, not the FINAL-named `app_test.go`

FINAL's change set names `internal/app/app_test.go` as the home of @codex-1's two fake-zcode
cases; they shipped in the new `internal/app/zcode_verify_test.go`. Same package, same
effect, arguably cleaner — but IMPLEMENTATION.md records three deviations (D-1..D-3) and this
relocation is not among them, so a reader diffing FINAL against the deviation list can miss
it. Fix: none required for code; record the relocation (or don't, and accept this NIT as the
record).

### NIT — dash-leading prompt values: zcode rejects them, and the spec Notes could say so

Measured against the real binary: `zcode --prompt "-n ..." --mode yolo --cwd .` exits 1 with
"Option '--prompt' argument is ambiguous ... use '--prompt=-XYZ'". This is fail-closed for
parley (nonzero exit, no artifact, full verify rejects it), and no parley-generated prompt
starts with a dash today. But the spec `Notes` record "rejected flags exit 1" without noting
that a rejected prompt VALUE has the same shape; one clause ("prompts must not begin with `-`;
the `--prompt={prompt}` form would lift that") would make the Notes complete. No behavioural
change required.

## Refutation attempts

1. **Argv lock — substitution and intactness.** Fake `zcode` on PATH recording argc and every
   argv element, driven through the real `parley agents verify --full --agent zcode --yes`
   (scratch deck, scratch PARLEY_HOME). Captured: `argc=6`,
   `--prompt|<multi-line prompt as ONE element>|--mode|yolo|--cwd|<root>`. The multi-line
   probe prompt arrived byte-intact. `{root}` substituted to `.` when run from the deck root
   (consistent with `cmd.Dir = root`) and to the absolute `/tmp/zcode-review/deck` when run
   from a foreign cwd via `-dir`. Exact-token substitution confirmed in
   internal/runner/runner.go:1101-1112 (`exec.CommandContext` with an args slice — no shell,
   so quotes/newlines cannot split).
2. **Argv lock — hostile prompt against the REAL zcode.** `--prompt "-n reply ..."` → exit 1
   with an explicit parser diagnostic (see NIT above). My first measurement reported rc=0
   because it piped through `tail` — the exact measurement trap FINAL's corrections section
   describes; re-measured without a pipe. Rejected `--model` also exits 1 ("Unknown option
   '--model'"), confirming the spec Notes. `zcode --help` confirms no model flag and documents
   `--mode` "default: yolo for --prompt", matching FINAL's rationale for passing it
   explicitly.
3. **Acceptance command, both ends.** Real repo, real binary:
   `agents verify --full --agent zcode --yes` → rc=0,
   `zcode: installed version=zcode-app-cli 3.7.7-13`, `headless probe passed`.
   `--agent zcode-1` → rc=2 `unknown agent zcode-1` (selectDiscoveries matches discovery IDs,
   as FINAL says). `agents list` prints the `headless:` argv line for zcode.
4. **Contract claim — five attempts to make `roster show` display a model.** (a) roster-level
   `[roster.zcode-1] model/effort`, (b) deck-level `[agents.zcode] model/effort`, (c)
   machine-level `[agents.zcode] model` (scratch PARLEY_HOME): all three still render
   `unknown / model-unbound,effort-unknown`. (d) full `headless_args` override carrying a
   `--model <literal>`: the literal is normalized to `{model}` and then bound from the
   machine layer's model field — the displayed model IS what the launch passes. (e) override
   plus deck-level `model`: displays the deck value, again matching the launch. (d) and (e)
   are the contract working (MODEL = what argv carries), not violations; the only dishonest
   surface in that state is the static trailer — Finding 1.
5. **Test inertness — spec revert (scratch clone of 600652b at /tmp, shared tree untouched).**
   Removed the zcode spec block from `discover.go`: `TestDefaultSpecsMergesACPCatalog`,
   `TestZcodeSpecCarriesNoModelOrEffortPlaceholder`, `TestMachineFamilyCatalogHasBuiltins`,
   `TestZcodeFullVerifyRejectsHelpExitZero`, `TestZcodeFullVerifyAcceptsHonestLaunch` ALL
   fail. Restored.
6. **Test inertness — two mutations.** (a) Neutered the sentinel-prefix check in
   `runHeadlessProbe`: `TestZcodeFullVerifyRejectsHelpExitZero` still passes — because the
   artifact-EXISTENCE check (`os.ReadFile(outputPath)`) fires first; the help-exit0 fake
   writes nothing. The test is guarded by the existence check, with the sentinel as
   defense-in-depth; not inert, but worth knowing which check actually bites. (b) Appended
   `--max-turns 5` to the spec argv: both the fake-zcode honest case (argc assertion, exit
   64) and the exact-argv lock fail. The argv lock is live in both places.
7. **`.zcode/skills` (D-1) independently.** `~/.zcode/skills` indeed does not exist on disk
   (D-1's premise confirmed). But `zcode skills list` runs and discovers 14 skills, including
   one from `~/.agents/skills` — the sibling user root the vendor's bundled
   zcode-configuration-guide documents alongside `~/.zcode/skills/`
   ("User `~/.zcode/skills`", "Workspace `.zcode/skills`"). The installer's resolved dest
   executes to `/Users/tomasfecko/.zcode/skills/parley-deck` with `detected: true`.
   `~/.zcode/cli/config.json` exposes `skills.roots` (currently `[]`). The target choice is
   vendor-documented, not inferred — D-1 holds. (I did not install into the real `~/.zcode`;
   that would mutate state outside the working directory.)
8. **D-2 — derived vs hardcoded count.** The form I first diffed derived the expected count
   from `result.actions` — circular: a target dropped from the install plan shrinks both
   sides and the test stays green, and no other test pins the full target list (the "13
   targets" in bidding-addon.test.js is a comment, not an assertion). Mid-review the tree was
   revised to `EXPECTED_TARGETS = 15` (registry-size tripwire) plus `planned` deep-equal
   `registered` plus `cores === EXPECTED_TARGETS` — which closes exactly the circularity, and
   its comment records a revert-proof. The current form is as strong as the hardcoded 14 was,
   on the addition AND removal directions. `npm test` green against both states (386/0).
9. **`--explain` trailer scope.** Present for zcode rows (deck and machine scope), absent for
   claude-1 (checked live), absent by construction for all other adapters (switch default).
   Truthful in every default-configuration state I could construct; false only under an
   operator argv override — Finding 1.
10. **Manifest integrity.** `parley-addon.json`'s SKILL.md sha256 matches the file
    (`a8e389fd...`), and `npm test`'s addon checks pass (`parley-deck: ok`, aggregate
    `b530f589...`). Machine-config obligations from FINAL verified at ~/.parley/agents.toml:
    `[agents.zcode]` block removed (comment at line 101), `[roster.zcode-1]` carries only
    `adapter`/`speed` (lines 201-205), stale exit-0 note corrected (lines 106-111).

## Open questions

- The skill package was edited during this review (the D-2 test rewrite and the `TARGETS`
  export appeared with codex-1 round-1 attributions while I was working). Which tree state is
  the candidate for merge, and will the review round re-run cleanly once it stops moving?
- `agents list` still shows a *configured* model for zcode when a config layer sets one
  (with a `sources:` attribution) while `roster show` correctly refuses it. IMPLEMENTATION.md
  frames this as the contract working; agreed — but is the divergence between the two
  surfaces documented anywhere an operator will find, other than this idea's files?
- The full end-to-end proof of the installer target (install into a scratch HOME, then
  `zcode skills list` discovering parley-deck from `~/.zcode/skills`) was not run — the
  pieces are individually verified, the composition is not. Worth one line in the release
  checklist since D-1's whole premise was "would fail silently".

## Addendum — tree drift observed after the findings above were written

Uncommitted edits landed in the CLI working tree during this review (mtimes
1787096701-1787097235, i.e. after my `go test ./...` pass began; attributions in the new
comments cite codex-1's round-1 review). I did not make or request these edits; I only
measured them:

- `Spec.NoModelBinding` (discover.go, launchargs.go): when set, `EffectiveModel`/
  `EffectiveEffort` report unbound regardless of argv, and the zcode spec sets it. This is a
  STRONGER fix for my Finding 1 than the one I proposed: I re-ran the exact attack that
  produced the finding (operator `headless_args` override carrying `--model {model}` plus a
  deck-level `model`) against a fresh build of the current tree, and the row now renders
  `model unknown` / `model-unbound` with the trailer truthful again. Finding 1 stands against
  600652b and is addressed in-flight.
- `--prompt={prompt}` equals form plus embedded-placeholder substitution in
  `buildAgentInvocation` (runner.go): addresses my NIT 4 (dash-leading prompt values),
  measured against the real parser by the implementer and by me. Note this changes the zcode
  argc from 6 to 5, so the fake-zcode assertions necessarily moved with it
  (`zcode_verify_test.go` +77 lines in the same drift) — the in-flight test changes were not
  re-derived by me.
- Re-ran against the current tree: `go build ./...` clean; `go test ./internal/agents/
  ./internal/app/ ./internal/runner/` all ok. I did not re-run the full suite, the
  scratch-copy revert, or the acceptance command against the drifted tree — those results in
  this review certify 600652b as written, and the drift itself needs its own review pass
  (the `NoModelBinding` semantics in particular change `EffectiveModel` for every adapter,
  currently gated by a flag only zcode sets).
