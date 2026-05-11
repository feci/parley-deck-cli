---
idea: agent-runtime-config
review-cycle: 1
drafted-by: codex
date: 2026-05-11
reviewed-commit: d38e7f738471caa61042ca087c445b4df5bd2606
---

## Agreed fixes

- From `claude` review round 01 [MAJOR] `Configured isolated_home_env is ignored for gemini and hermes`: thread configured `IsolatedHomeEnv` through the Gemini and Hermes isolated-home helpers, falling back to the historical env var names only when no template is configured. Add focused coverage.
- From `claude` review round 01 [MAJOR] and `gemini` review round 01 [MINOR] on the constant `HEADLESS` column: stop printing an unconditional `not-probed` value. The matrix should show useful configured headless state, and full verification should report probe results in command output.
- From `gemini` review round 01 [MAJOR] `Codex Git smoke runs directly in CLI instead of via agent`: move the Codex Git smoke requirement into the Codex full headless probe prompt so the Codex runtime is the actor asked to perform the Git operations. Keep a probe sentinel, and adjust tests around the prompt/probe contract rather than direct CLI Git execution.
- From `claude` review round 01 [MINOR] `New TOML-declared agents lose source attribution for unset fields`: populate default source metadata for synthesized agents instead of falling back to `discovered`.
- From `claude` review round 01 [MINOR] `{tempdir} placeholder has two different meanings depending on the field`: document the exact placeholder semantics and add/adjust tests so the split is intentional rather than accidental.
- From `claude` review round 01 [MINOR] `runFullVerification fails fast on the first error`: collect per-agent full verification errors, print all probe outcomes, and return non-zero at the end if any probe failed.
- From `claude` review round 01 [MINOR] `parley agents verify` strict missing-agent behavior: document the strict default behavior and print a hint to use `--agent ID` for partial verification.
- From `claude` review round 01 [MINOR] `No app-level test that parley run actually injects resolved runtime`: add an app-level run test that verifies resolved runtime appears in `run.created`.
- From `claude` review round 01 [NIT] `Compatibility aliases agents discover|probe are not exercised by tests`: add focused alias coverage.

## Deferred follow-ups

- From `claude` review round 01 [NIT] `valueOr` helper duplication`: defer small helper consolidation until another shared formatting pass.
- From `claude` review round 01 [NIT] `text/tabwriter`: defer matrix formatting polish unless the fixed headless column still renders poorly.
- From `gemini` open question: persisting full-probe state for future `agents list` output is useful but out of scope for this fix-up.
- From `gemini` open question: runtime override flags for model/profile/timeout are already a FINAL.md non-goal.
- From `hermes` open question: printing exact resolved command lines during full verification is useful but can expose sensitive local paths; defer until a redaction policy exists.

## Dismissed findings

- From `claude` review round 01 [MINOR] `Shared agents.toml re-declares built-in defaults`: no functional fix required. The checked-in file intentionally records project defaults for auditability, even when they currently match built-ins. The fix-up should avoid relying on this file for correctness, but it may remain.
- From `hermes` review round 01 [MINOR] `agents.toml missing some fields present in runtime model`: dismissed. `agents.toml` is shared defaults, not an exhaustive schema example; the user-facing docs carry the field guidance.
- From `hermes` review round 01 [NIT] `Hardcoded timeout in verify smoke test`: dismissed as not applicable to the current implementation; no fixed 30s Codex Git smoke timeout exists in the reviewed code.
- From `gemini` review round 01 [NIT] `Default values for non-proven fields`: dismissed. `speed = balanced` is an explicit Parley default and `external_backend = unknown` is a domain value, not a model/reasoning/profile invention.

## Signoffs

<!-- Each active participant APPENDS their signoff block. -->

### Signoff: codex — 2026-05-11
Status: ✅ ACCEPT
Notes: I accept this review consensus and will implement the agreed fixes in fix-up cycle 1 before requesting re-review.

### Signoff: claude — 2026-05-11
Status: ✅ ACCEPT
Notes: Both my MAJORs (isolated_home_env ignored; HEADLESS placeholder column) are accepted, six of seven MINORs and one NIT are agreed, and all three open questions are resolved. The dismissed `agents.toml` re-declaration MINOR has a defensible rationale (intentional auditability with no correctness dependency) — non-blocking.

### Signoff: gemini — 2026-05-11
Status: ✅ ACCEPT
Notes: My MAJOR regarding Codex Git smoke isolation and my MINOR regarding the HEADLESS column are both addressed. The dismissal of my NIT regarding default values is acceptable as those values are indeed intentional Parley domain defaults.
