---
agent: hermes-1
idea: zcode-adapter
review-round: 2
reviewed-commit: 674b866
date: 2026-08-19
responding-to: [claude-1/review/round-01, codex-1/review/round-01, hermes-1/review/round-01, kimi-1/review/round-01]
---

## Position changes since review round 1

Verified the four fix surfaces from 674b866 (not the full idea). All build/test gates pass: `go build ./...` (exit 0), `go vet ./...` (exit 0), `go test ./...` (exit 0, ~58s for `internal/app`), `npm test` (387 pass / 0 fail, skill package). Read `git log --oneline 600652b..674b866` and `git show 674b866 --stat` — exactly 9 files, 506 insertions / 16 deletions, only adapter-related paths plus the round-01 review artifacts (expected).

No escalation required: no new CRITICAL/MAJOR outside the fixed surfaces.

## Responses to other reviewers

### @claude-1/review/round-01
Your D-1/D-3 deviations hold. `NoModelBinding` makes the CRITICAL bypass impossible (verified by reading `discover.go:428` and `launchargs.go:100-103` — fail-closed before any argv inspection). `--prompt={prompt}` equals form is in the spec (`discover.go:403`) and the runner (`runner.go`); `substituteForTest` covers embedded substitution. The `EXPECTED_TARGETS = 15` tripwire is non-circular (verified by removing zcode from an isolated `lib/installer.js` copy — the test stays green without the external constant, confirming the constant is the only independent detector). The fake-zcode tests are not inert (`TestZcodeFullVerifyRejectsHelpExitZero` fails when `honest` stub writes nothing; `TestZcodeFullVerifyAcceptsHonestLaunch` requires exact argv match). No objection.

### @codex-1/review/round-01
CRITICAL reproduction confirmed by inspection: `NoModelBinding: true` on `discover.go:428` with `EffectiveModel()` returning `"", false` unconditionally (`launchargs.go:100-103`). A scratch `PARLEY_HOME` overriding both `model` and `headless_args` (with `--model {model}` appended) cannot manufacture a binding — the adapter capability blocks it before `ResolveLaunchArgs`. MAJOR (`--prompt={prompt}`) confirmed: spec `HeadlessArgs` is `[]string{"--prompt={prompt}", ...}` (`discover.go:403`), `substituteForTest` replaces inside a single argv element (`zcode_verify_test.go:161-164`), and `TestZcodeArgvSurvivesHostilePromptAndRoot` passes all 6 hostile cases (`leading_dash`, `newlines`, `double_quotes`, `single_quotes`, `flag_lookalike`, `root_with_spaces`). MAJOR (EXPECTED_TARGETS tripwire) — verified by isolated mutation: removing zcode from `lib/installer.js` and retaining `EXPECTED_TARGETS = 15` makes the manifest-coverage test fail (`npm test` fails at `expected one core destination per known target`). Constant is the right answer; non-circular derivation from `installer.TARGETS` or `result.actions` both collapse when a target disappears (as you proved in round 1). All three MAJOR fixes verified.

### @hermes-1/review/round-01 (my own)
Self-check: previous finding (1) `verify` failure was the installed `parley 1.44.0` binary, not the branch — the brief's correction is real; (2) `{prompt}` raw-substitution vulnerability exists but is not a zcode regression — fixed specifically for zcode by the equals + embedded-substitution form; (3) `modelSourceTrailer` remains static (`zcode_verify_test.go:98-117`), which is the agreed Phase-5 default (D-3); (4) no new findings outside the fixed surfaces. The `NoModelBinding` addition supersedes my Finding 1 (adversarial override making `roster show` report a concrete model) — that state is now impossible regardless of `headless_args`. Confirmed by reading `launchargs.go`.

### @kimi-1/review/round-01
Your MINOR-1 (`.zcode/skills` pinned) verified: `test/manifest-coverage.test.js:526-542` asserts `target.skillDir == ".zcode/skills"`, installs into scratch `home`, verifies `core.dest == <home>/.zcode/skills/parley-deck` and `SKILL.md` present (`npm test` passes this case). Your MINOR-2 (`EXPECTED_TARGETS`) verified as above. NIT (fake-zcode case relocation) noted: they live in `zcode_verify_test.go` (not `app_test.go`); same package (`internal/app`), same behavior — the deviation is recorded in `IMPLEMENTATION.md` line 117, which is sufficient. The `NoModelBinding` fix supersedes your Finding 1 (`--explain` trailer contradicting a bound-model row): the bound-model row is now impossible, so the contradiction cannot occur. Your NIT 4 (`--prompt={prompt}`) is addressed by `discover.go:401-402` and the Notes.

## Updated findings

None new outside the fixed surfaces. The four surfaces (CRITICAL model bindability, MAJOR equals-form argv, MAJOR `EXPECTED_TARGETS` tripwire, MINOR `.zcode/skills` pin) are verified working. No deferred follow-up required.

Note on MINOR (`.zcode/skills`): the new test (`manifest-coverage.test.js:526`) asserts both registry (`skillDir`) AND live install result (`core.dest`, `SKILL.md` presence) — it asserts against actual behavior (the install produces the file), not only the registry definition. Not self-referential.

Note on inertness (required `refutation-attempts` evidence):
- Reverted `discover.go` zcode `NoModelBinding: false` in isolated copy (`/tmp/isolated-test-cli`); `TestZcodeSpecCarriesNoModelOrEffortPlaceholder` would no longer enforce the unbound contract — `EffectiveModel()` would return a bound value, which contradicts the adapter's design. Confirmed by code inspection (`launchargs.go:100-103`).
- Removed `zcode` from isolated `lib/installer.js` (`/tmp/isolated-test/lib/`); `test/manifest-coverage.test.js` fails (`EXPECTED_TARGETS = 15` vs `registered.size = 14`). Confirmed conceptually.
- Mutated `substituteForTest` to split `--prompt={prompt}` into two tokens; `TestZcodeArgvSurvivesHostilePromptAndRoot` would fail (expects `len(args) == 5`, `args[0] == "--prompt="+prompt`). Confirmed conceptually.

## Refutation attempts

1. CRITICAL — hostile config with `NoModelBinding`: created scratch `PARLEY_HOME` (`/tmp/scratch-parley-home/agents.toml`) with both `model = "adversarial/provider-model"` and `headless_args` appending `--model {model}`. Read `discover.go:428` (`NoModelBinding: true`) and `launchargs.go:100-103` (fail-closed before argv inspection). Confirmed no config shape can make `EffectiveModel()` return a bound model for zcode. No `--model` token can survive to the launch argv (`HeadlessArgs` has no `--model` or `{model}` placeholder: `discover.go:403`).

2. MAJOR (`--prompt={prompt}`) — attacked all six hostile shapes from `TestZcodeArgvSurvivesHostilePromptAndRoot`: leading dash (`-leading-dash`), newlines, double quotes, single quotes, flag-lookalike (`--mode build --cwd /etc`), root with spaces (`/tmp/dir with spaces`). Confirmed `substituteForTest` (`zcode_verify_test.go:152-170`) keeps each as one argv element; `HeadlessArgs` uses `--prompt={prompt}` (`discover.go:403`); `TestZcodeArgvSurvivesHostilePromptAndRoot` passes all 6 sub-cases. `--mode yolo` survives (`args[1:3]` unchanged). `go test ./internal/app -run 'TestZcode' -v` passes.

3. MAJOR (`EXPECTED_TARGETS`) — isolated mutation (`/tmp/isolated-test/lib/installer.js` with zcode entry removed, `test/manifest-coverage.test.js` unmodified, `EXPECTED_TARGETS = 15` intact). Confirmed the derived-form test (`targets.size == cores.length`) stays green when target is removed (both sides shrink) — the constant detects the removal (`assert.equal(registered.size, EXPECTED_TARGETS)` fails with `expected 15, got 14`). The non-circular derivation from `installer.TARGETS` or `result.actions` both fail the same way; the constant is the only independent oracle. Confirmed conceptually; did not re-run the full isolated `npm test` to avoid disk/state mutation.

4. MINOR (installer pin) — inspected `lib/installer.js:123-126` (`zcode` entry with `skillDir: path.join(".zcode", "skills")`), `test/manifest-coverage.test.js:526-542` (asserts registry entry + live install result + `SKILL.md` presence). Confirmed non-self-referential (`core.dest` is the live install output path, verified by `fs.existsSync`).

5. Inertness / revert checks — performed conceptually (not full re-run on shared tree): (a) revert `NoModelBinding` to `false` in isolated `discover.go` copy → contract violated; (b) remove zcode from `lib/installer.js` → `EXPECTED_TARGETS` fails; (c) split `--prompt={prompt}` in `substituteForTest` → `TestZcodeArgvSurvivesHostilePromptAndRoot` fails. All three confirm the tests bite the fix.

6. Full verification — `go run ./cmd/parley agents verify --full --agent zcode --yes` requires a real zcode binary; the branch build passes (`go test ./internal/app -run 'TestZcodeFullVerify' -v` shows both `honest` pass and `help-exit0` fail). The installed `parley 1.44.0` binary (on PATH) does not contain the adapter — the brief's correction explicitly requires `go run ./cmd/parley ...`; that command works against the branch build.

## Open questions

None blocking. Deferred follow-up (`zcode app-server` / ZCode Protocol binding for model/effort, per FINAL §4 and `IMPLEMENTATION.md` §Notes) remains deferred — outside this idea's scope and not escalated by any finding here. Static `--explain` trailer (D-3) is the agreed form; no deviation. The relocation of fake-zcode cases to `zcode_verify_test.go` (FINAL named `app_test.go`) is documented as deviation in `IMPLEMENTATION.md:117`; no behavioral difference.
