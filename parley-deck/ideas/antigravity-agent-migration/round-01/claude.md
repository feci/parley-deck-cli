---
agent: claude
idea: antigravity-agent-migration
round: 1
date: 2026-05-27
---

## Summary

Replacing Gemini CLI with Antigravity CLI (`agy`) touches the built-in agent catalog, the runner's isolated-home logic, the project and embedded-default COOPERATION.md rosters, `agents.toml`, help text/examples, documentation, and test fixtures. The migration is straightforward mechanically, but protocol safety requires keeping `gemini` as an inactive roster entry rather than deleting it, and the new `agy` spec must be authored from actual CLI discovery (`agy --help`, `agy --version`) rather than assumed flag names.

## Proposed approach

### 1. Roster update (protocol-safe)

In both `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`:

- Mark `gemini` as **inactive** in the roster table (add a column or parenthetical note) per COOPERATION.md section 2: "When an agent leaves the project, mark its row as inactive (do not delete it) so historical references remain resolvable."
- Add `agy` as a new active participant row.
- Mirror the same change in the host-handle mapping table.

Do **not** delete the `gemini` row. Historical round files (e.g., `round-01/gemini.md` across dozens of past ideas) reference this ID and must remain resolvable.

### 2. Built-in agent spec (`internal/agents/discover.go`)

- Add a new `agy` entry to `defaultBuiltinSpecs()` using flags and modes discovered from the live `agy` CLI (not invented). If `agy --help` reveals specific `--model`, `--thinking`, or `--profile` flags, populate them; otherwise use `cli-default` per the prompt's constraint.
- Keep the existing `gemini` spec but add a `Notes` field marking it as `"DEPRECATED: legacy Gemini CLI support; prefer agy"` and move it below the active entries or to a separate `legacyBuiltinSpecs()` helper.
- If `agy` requires isolated home (like gemini did for OAuth hangs), add the equivalent `IsolateHome` / `IsolatedHomeEnv` with the correct env var name (e.g., `AGY_HOME` or whatever the CLI uses).

### 3. Runner isolated-home logic (`internal/runner/runner.go`)

- Add a `case "agy":` handler alongside the existing `case "gemini":` and `case "hermes":` blocks.
- Determine whether `agy` needs the same OAuth credential copying that `isolatedGeminiHome()` performs, or if it uses a different auth mechanism. If it shares config with Gemini (the `.antigravitycli/` symlink to `.gemini/config/` suggests some relationship), document this and handle it.
- Keep `case "gemini":` functional but add a log/warning that gemini is deprecated.

### 4. Project agent config (`parley-deck/agents.toml`)

- Add `[agents.agy]` section with appropriate defaults.
- Rename `[agents.gemini]` to include a comment marking it as legacy/deprecated, or move it to a `# Legacy agents` section. Do not remove it — existing local overrides in `agents.local.toml` may reference it.

### 5. CLI help text and examples (`internal/app/app.go`)

- Replace `gemini` with `agy` in the `--participants` example (line 176) and the `run` example (line 242) and the `consensus request-signoffs` example (line 255).

### 6. Documentation (`README.md`, `docs/`)

- Update all active documentation to list `agy` as a default participant and `gemini` as deprecated.
- Do **not** edit historical idea artifacts (per the prompt constraint).

### 7. Tests

Update test fixtures in:
- `internal/agents/acp_specs_test.go` (line 60): add `"agy"` to the required IDs if it gets an ACP entry, or update the `"gemini"` reference.
- `internal/runner/runner_test.go`: add test for `agy` isolated home env; keep gemini tests for legacy coverage.
- `internal/tui/live_test.go`: update test data to use `"agy"` instead of `"gemini"` for the default agent roster.
- `internal/hitl/hitl_test.go`: update or add a test case for the new agent name.

### 8. Release metadata

- Update npm `package.json`, Homebrew formula, and any other packaging metadata to mention `agy` as a supported participant.
- WinGet manifests should only be claimed as updated if actually created/submitted per the non-goals constraint.

## Concerns / open questions

1. **CLI discovery first.** The prompt explicitly says "Do not invent model or thinking flag names." Before writing the `agy` spec in `discover.go`, someone must run `agy --help` and `agy --version` to determine the actual flags for headless mode, model selection, thinking/reasoning, sandbox, approval, and output format. The spec cannot be assumed from gemini's flags.

2. **Auth mechanism.** The `.antigravitycli/` directory in the repo root is a symlink to `.gemini/config/projects/...`. Does `agy` share Gemini's OAuth/API-key auth, or does it have its own? This affects whether `isolatedGeminiHome()` logic can be reused or needs a distinct `isolatedAgyHome()`.

3. **ACP support.** `agy` does not appear in `ACPCatalog()` today. If Antigravity CLI supports ACP (JSON-RPC over stdio), it should be added there. If not, it should remain headless-only.

4. **In-flight ideas.** The current `00-prompt.md` for this very idea still lists `participants: [codex, claude, agy, hermes]` — meaning `agy` is already being used as a participant ID for new work. But the roster in COOPERATION.md still says `gemini`. We need to reconcile: either this idea's `participants:` list is forward-looking (treating `agy` as the approved replacement per the constraint), or the roster must be updated before this idea can complete its full lifecycle.

5. **Default participant list in `parley init`.** When a user runs `parley init`, what default roster is generated? If it currently defaults to `[codex, claude, gemini, hermes]`, this must change to `[codex, claude, agy, hermes]`.

6. **Backwards compatibility scope.** The prompt says "Keep legacy Gemini support only if needed for backwards compatibility." I recommend keeping the gemini spec discoverable (so `parley agents` still shows it if the binary is on PATH) but not as a default participant. Users with gemini in existing `agents.toml` or `participants:` lists should get a deprecation warning, not a hard error.

## Risks

1. **Breaking existing workflows.** Users who have `gemini` hardcoded in `agents.toml`, `participants:` lists, or automation scripts will break if we remove the gemini spec entirely. Mitigation: keep gemini as a deprecated-but-functional spec; emit a warning at runtime.

2. **Incorrect spec for `agy`.** If we write the `agy` headless spec without verifying the actual CLI flags, headless launches will fail silently or produce malformed output. Mitigation: mandatory `agy --help` discovery before implementation.

3. **Isolated home regression.** Gemini's `isolatedGeminiHome()` was added to work around a specific OAuth hang. If `agy` inherits the same auth mechanism but we don't replicate the isolation, we reintroduce that hang. If `agy` does not share the mechanism, we waste complexity copying it. Mitigation: test both scenarios.

4. **Release coordination.** Changing the default roster affects every new `parley init`, every CI pipeline, and every user who updates. The release notes must clearly document the migration path. Failing to update Homebrew/npm/WinGet metadata in sync creates a window where the CLI references `agy` but the packaging metadata still says `gemini`.

5. **Test coverage gap.** If we update tests to use `agy` but leave no test exercising the deprecated `gemini` path, we lose regression coverage for users still on gemini. Keep at least one test per deprecated code path.
