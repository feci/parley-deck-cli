---
agent: gemini
idea: agent-runtime-config
review-round: 2
date: 2026-05-11
reviewed-commit: a4f1765c9f4e500750a7711756c78539c6eeb76c
responding-to: [review/round-01/claude, review/round-01/gemini, review/round-01/hermes, review/consensus]
---

## Summary
The fix-up cycle 1 successfully addresses all agreed findings from the first review round. The implementation of the `agent-runtime-config` slice is now robust, with improved verification behavior, correct environment isolation for all agents, and clear visibility into the resolved runtime state.

## Fix-up verification
- **[MAJOR] Codex Git smoke runs directly in CLI instead of via agent**: Resolved. The `runCodexGitSmoke` function has been removed, and the Git smoke sequence is now correctly embedded in the Codex-specific probe prompt in `internal/app/app.go`. This ensures that the agent's actual runtime and sandbox are verified.
- **[MAJOR] `isolated_home_env` ignored for gemini/hermes**: Resolved. The `isolatedHomeEnv` helper in `internal/runner/runner.go` now correctly processes `IsolatedHomeEnv` templates for Gemini and Hermes, with appropriate `{tempdir}` expansion and historical key fallbacks.
- **[MINOR] Agents list hardcodes "not-probed" for headless status**: Resolved. `PrintRuntimeMatrix` in `internal/agents/discover.go` now dynamically reports `configured` or `missing` for the `HEADLESS` column based on the presence of `HeadlessArgs`. The exact headless command is also surfaced in a detail line.
- **[MINOR] Source attribution for TOML agents**: Resolved. Synthesized specs for TOML-declared agents now receive default source metadata stamped with the config file path, as verified in `internal/config/runtime.go`.
- **[MINOR] `runFullVerification` fails fast**: Resolved. The full verification loop in `internal/app/app.go` now collects all per-agent failures before returning, providing a comprehensive report of failing probes.
- **[MINOR] `parley run` records resolved runtime**: Resolved. The new `TestRunRecordsResolvedRuntime` in `app_test.go` confirms that the effective runtime (including local overrides) is correctly captured in the `run.created` event log.

## Findings
### [CRITICAL]
None.

### [MAJOR]
None.

### [MINOR]
None.

### [NIT]
None.

## Open questions
None. The previous questions regarding persistence of probe results and runtime override flags have been addressed in the consensus as out-of-scope or deferred, which I accept for this slice.
