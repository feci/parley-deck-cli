---
agent: codex-1
idea: zcode-adapter
review-round: 1
date: 2026-08-19
reviewed-commit: 600652b
---

## Summary

Not ready for review consensus: I found one CRITICAL contract bypass, three MAJOR defects, and one MINOR release-surface gap. The built-in happy path is mostly sound, the fake-zcode tests are not inert, and `.zcode/skills` is a real ZCode discovery root. However, a config layer can make `roster show` report a concrete zcode model; a leading-dash prompt is rejected by the real CLI; D-2 replaced an independent coverage oracle with a circular one; and the reviewed CLI/skill state is not anchored to reproducible commits.

`go build ./...` passed. `npm test` in the skill package passed (386 Node tests and 54 Python tests). `go test ./...` did not pass in this sandbox because `TestDurableKillEndToEndRealProcess` cannot obtain the macOS boot ID; the same isolated test fails at parent `9a6f8be`, so this is not introduced by `600652b`, but the required gate still needs a host rerun.

## Findings

### [CRITICAL] A config override can make zcode model-bound and make `--explain` contradict the roster row

**What is wrong.** A model-only roster override behaves correctly, but the zcode invariant is bypassable through the supported wholesale `headless_args` override. In a scratch `PARLEY_HOME`, I set `zcode-1.model = "adversarial/provider-model"` and added:

```toml
[agents.zcode]
model = "adversarial/provider-model"
headless_args = ["--prompt", "{prompt}", "--mode", "yolo", "--cwd", "{root}", "--model", "{model}"]
```

`roster show --scope machine --explain zcode-1` then printed `model adversarial/provider-model`, omitted `model-unbound`, and immediately printed the static trailer saying the model is "never passed by parley". [PRIMARY: scratch command/output.] `applyOverride` accepts the vector, `ResolveLaunchArgs` binds `{model}`, and `EffectiveModel` treats the resulting `--model` token as effective even though the zcode adapter is defined by the absence of any model flag. Existing deck/global overrides can therefore restore exactly the impossible state that deleting this machine's current `[agents.zcode]` block was meant to remove.

**Why it matters.** This breaks the explicit acceptance contract that zcode remains `model-unbound` even when configuration supplies a model. It also gives mutually contradictory operator output and would launch a flag the real zcode parser rejects, so the roster presents confidence where the process cannot successfully answer.

**Concrete fix.** Represent model binding as an adapter capability rather than inferring it solely from arbitrary argv. Mark zcode as non-bindable and reject (or fail closed on) config-supplied model flags/placeholders for that adapter; `roster show` must remain `unknown` / `model-unbound`. Make the trailer conditional on the resolved row state so it cannot say "never passed" beside a concrete MODEL. Add a layered-config regression test that supplies both `model` and hostile `headless_args` and asserts the invariant plus non-contradictory `--explain` output.

### [MAJOR] The locked argv is rejected when the prompt begins with a dash

**What is wrong.** The shipped vector passes the prompt as the next token after `--prompt`. Running the real CLI with a newline/quote-bearing prompt whose first character is `-` produced:

```text
Option '--prompt' argument is ambiguous.
Did you forget to specify the option argument for '--prompt'?
To specify an option argument starting with a dash use '--prompt=-XYZ'.
```

[PRIMARY: `zcode --prompt $'-leading-dash\nline with ...' --mode yolo --cwd "$PWD"`, exit 1.] The equivalent `--prompt=$'-leading-dash\n...'` form passed option parsing and reached ZCode startup. The fake "honest" zcode accepts the separate-token form unconditionally, so its exact-order assertion locks a behavior the real parser rejects.

**Why it matters.** An otherwise valid prompt can fail before a model session starts. The current tests give false assurance precisely on the edge case this review brief required attacking.

**Concrete fix.** Emit one argv token in equals form, `--prompt={prompt}`, and extend safe placeholder substitution to replace `{prompt}` inside that single argv element without invoking a shell. Update the fake-zcode lock to parse `--prompt=<value>` and add cases for a leading dash, newlines, double/single quotes, and a root containing spaces.

### [MAJOR] D-2 makes manifest coverage circular and weaker than the fixed count

**What is wrong.** The new test derives both sides from the install result: `cores` comes from `result.actions`, and `targets` is the set of those same actions. In a scratch copy I reverted only the zcode entry in `lib/installer.js` while leaving D-2 in place; `node --test test/manifest-coverage.test.js` still passed all 15 tests. [PRIMARY: scratch mutation/check.] The old fixed cardinality would have failed when a registered target disappeared.

**Why it matters.** A target can be dropped from the registry or omitted from the plan and the coverage test shrinks its expectation to match the defect. That is the failure class the test is supposed to detect.

**Concrete fix.** Keep an independent expected target-name snapshot in the test (now including `zcode`), assert the action target set equals it, and then assert exactly one core destination for every expected name. Intentional target additions should update that oracle in review; avoiding that update removes the test's independence.

### [MAJOR] The reviewed implementation is not anchored to reproducible CLI and skill revisions

**What is wrong.** `IMPLEMENTATION.md` records `head-commit: 9a6f8be`, which is the parent of the actual implementation commit and contains none of the adapter code. Commit `600652b` claims "CLI + skill", but the skill package is a separate repository on `main` at `4412a8a` with four uncommitted modified files. The CLI commit itself contains 56 files / 9,902 insertions; 36 paths are outside `internal/` and `parley-deck/ideas/zcode-adapter/`, including three unrelated ideas, an unrelated inbox note, and the deck roster edit. [PRIMARY: `git show 600652b --stat`, `git show --name-only`, and skill-package `git status`.]

**Why it matters.** A reviewer or releaser cannot reconstruct one immutable CLI+skill candidate from the implementation record. Unrelated canonical artifacts can be merged or reverted with the adapter, and the skill half can drift after this review without changing `reviewed-commit`.

**Concrete fix.** Split/rebase the CLI implementation onto a clean base so the zcode change carries only its code/tests and its own protocol artifacts; preserve unrelated agents' artifacts in their own commits. Commit the skill changes on a feature branch, record that repository and SHA in `IMPLEMENTATION.md`, and update the CLI `head-commit` to the commit that actually contains the code. Review both immutable SHAs before release.

### [MINOR] The native installer surface is stale and lacks a direct zcode regression lock

**What is wrong.** The skill package README still says there are "fourteen named runtimes" and lists no ZCode. `test/installer.test.js`'s user/project destination assertions also omit zcode, so the new `.zcode/skills/parley-deck` destination has no direct target-specific lock. [PRIMARY: README and test inspection.]

**Why it matters.** The public installation contract disagrees with the implementation, and a later path typo could survive the general fleet tests—especially given D-2's circular count.

**Concrete fix.** Update the README to fifteen targets and include ZCode in the list/examples. Add explicit user- and project-scope zcode path assertions and an explicit install/discovery test for the zcode target.

## Refutation attempts

- **Built-in presence and argv resolution.** [PRIMARY] `go run ./cmd/parley agents list` discovered `/opt/homebrew/bin/zcode`, reported AUTO=yes, and printed `--prompt {prompt} --mode yolo --cwd {root}`. Inspection of `runner.buildAgentInvocation` confirmed exact-token substitution into an `exec.Command` argv slice, so a root containing spaces is not shell-split.
- **Prompt preservation.** [PRIMARY] Newlines and both quote styles passed ZCode option parsing and reached runtime startup; the sandbox then denied ZCode's Unix-socket listener. A leading-dash prompt failed option parsing, producing the MAJOR finding above. Equals form passed parsing.
- **Model-unbound, ordinary configuration.** [PRIMARY] In a scratch `PARLEY_HOME`, `roster set zcode-1 --scope machine --model adversarial/provider-model --yes` followed by `roster show --scope machine --json` returned MODEL `unknown` and statuses `model-unbound, effort-unknown, metadata-unknown`.
- **Model-unbound, adversarial configuration.** [PRIMARY] Adding a zcode `headless_args` override containing `--model {model}` made the same row display `adversarial/provider-model` and removed `model-unbound`; `--explain` still claimed the model was never passed. This produced the CRITICAL finding.
- **Trailer scope.** [PRIMARY] `modelSourceTrailer` returns text only for adapter `zcode`, and the focused trailer test passed. The adversarial bound-model state proves the adapter-only predicate is insufficient even though unrelated adapters do not receive the trailer.
- **Fake-zcode test liveness.** [PRIMARY] I archived `600652b` to a scratch copy, replaced only `internal/agents/discover.go` with the parent version, and ran `go test ./internal/app -run 'TestZcodeFullVerify' -count=1 -v` with a fresh `GOCACHE`. Both fake-zcode tests failed at `built-in spec zcode is missing`; they are not inert. (The explicit `-count=1` avoids reusing a cached pass from the feature-present tree.)
- **Full verification.** [PRIMARY] The exact requested installed command used `parley 1.44.0` and failed `unknown agent zcode`. The source-tree command `go run ./cmd/parley agents verify --full --agent zcode --yes` discovered zcode but could not finish because this sandbox rejects ZCode's Unix-socket `listen` with `EPERM`. The implementation record contains an earlier successful real probe, but it is not tied to an immutable binary/revision.
- **Installer target D-1.** [PRIMARY] ZCode's bundled `zcode-configuration-guide` explicitly lists user `~/.zcode/skills`. I then installed the working skill package into a scratch HOME with `--target zcode --no-addons`; `zcode --cwd /tmp skills list --json` discovered `parley-deck` at that exact scratch `.zcode/skills/parley-deck/SKILL.md` path. D-1 is independently validated.
- **D-2 mutation.** [PRIMARY] Removing the zcode target from a scratch copy while retaining the derived assertion left all 15 manifest-coverage tests green, producing the MAJOR finding.
- **Focused tests.** [PRIMARY] `go test ./internal/agents -run Zcode -v` and `go test ./internal/app -run Zcode -v` passed on the feature-present tree. `TestDeriveModelMetaResolvesZaiNamespace` is included in the full agents package run.
- **Required suites.** [PRIMARY] `go build ./...` passed. Skill-package `npm test` passed 386 Node tests, 54 Python tests, and manifest checks. `go test ./...` failed only `internal/runner.TestDurableKillEndToEndRealProcess` because no boot ID was available; the same isolated test fails unchanged at parent `9a6f8be`, establishing that this environment failure is pre-existing rather than introduced by the adapter.
- **Scope/release state.** [PRIMARY] `git show 600652b` exposed the unrelated paths and absent skill commit described above. The installed `parley-deck-skill 2.8.0` rejects `--target zcode`, confirming that release/channel acceptance has not happened yet.

## Open questions

- FINAL says unqualified `roster show` should report zcode AUTO=yes and drop `not-in-roster`, but this deck does not declare `zcode-1`: default `roster show` omits it and `--all` reports zcode `not-in-roster`; only `--scope machine` shows the intended row. Was the acceptance criterion meant to be machine-scoped? Adding zcode to the deck would change future quorum and needs explicit membership approval.
- The real full probe must be rerun outside this sandbox at the fixed source/skill SHAs because ZCode's local socket cannot bind here.
- The required Go gate needs either a normal host rerun or a pre-existing test correction that skips the durable-kill end-to-end case when `procctl.CurrentBootID()` is unavailable, as the neighboring test already does.
- FINAL's minor-version, npm/Homebrew/winget/GitHub release remains pending. It should occur only after the findings are fixed and both implementation repositories are committed and reviewed.
