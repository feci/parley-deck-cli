---
idea: interactive-agent-mode
review-cycle: 1
drafted-by: codex
date: 2026-05-14
reviewed-commit: 348b693
---

## Agreed fixes

- From `claude` review round 01 [CRITICAL] `Manual + resume bypass append-only signoff validation`: store enough pre-handoff state for consensus signoff handoffs and make `resume` validate append-only content plus canonical signoff status before completing a pending handoff.
- From `claude` review round 01 [MAJOR] `Handoff "Suggested command" leaks unsubstituted placeholders`: use the same placeholder expansion for printed handoff commands, dry-run previews, and `spawn-tty` args.
- From `claude` review round 01 [MAJOR] `Manual handoff resume hint points at the wrong command`: handoff instructions and terminal output must tell the user to run `parley resume <run-id>` after appending a manual signoff. `consensus status` may remain an inspection hint only.
- From `claude` review round 01 [MINOR] and `gemini` review round 01 [MINOR] on `ValidateArtifact` Round 2+ compatibility: make the validator contract honest for this slice by renaming or narrowing it to Round 1, or by explicitly erroring for unsupported rounds.
- From `claude` review round 01 [MINOR] and `gemini` review round 01 [MINOR] on pending manual exit semantics: return a distinct pending exit code for pending manual handoffs and document it.
- From `claude` review round 01 [MINOR] `validateLaunchModes validates interactive enums for launch_mode = "manual" too`: either skip the unused invoke validation for manual mode or document why manual reuses the interactive handoff fields.
- From `claude` review round 01 [NIT] `Dead fallback in runInteractiveTTY`: remove the unreachable fallback path.
- From `claude` review round 01 [NIT] `IMPLEMENTATION.md frontmatter key drifts from protocol`: change `implemented-by` to `implementer`.

## Deferred follow-ups

- From `gemini` review round 01 [MAJOR] `Signal handling for spawn-tty`: defer full process-group signal handling to the already-deferred `spawn-tty` polish slice. The current implementation has a TTY gate and attached process launch, but complete signal semantics are not required before this first signoff-focused slice can be reviewed.
- From `gemini` review round 01 [NIT] `ReadFrontmatter parser robustness`: defer until the protocol accepts frontmatter beyond the current simple key/value shape.
- From `gemini` review round 01 open question on re-running failed headless agents during resume: defer to the mixed-mode round execution slice.
- From `claude` review round 01 open question on `--allow-interactive-noninteractive`: defer; print-only interactive currently works without TTY, while `spawn-tty` requires one.

## Dismissed findings

- From `gemini` review round 01 [NIT] `Hardcoded poll interval for interactive mode`: dismissed as not a bug. The default is 2000ms but is configurable with `interactive_poll_ms`.
- From `hermes` review round 01 [MAJOR] `Hosted confirmation flow incomplete`: dismissed. There is no `handoff` confirmation struct requirement in `FINAL.md` or `IMPLEMENTATION.md`; hosted confirmation now applies only to headless launches by design.
- From `hermes` review round 01 [MINOR] `Stricter artifact validation may break resume`: partially overlaps the agreed append-only validation fix, but the specific backward-compatibility claim is not supported by the implementation. Existing old run display is not changed by Round 1 validation.
- From `hermes` review round 01 [NIT] `Mode/config semantics docstring`: dismissed. There is no `InteractiveAgentConfig` docstring in the current code.

## Signoffs

<!-- Each active participant APPENDS their signoff block. -->

### Signoff: codex — 2026-05-14
Status: ✅ ACCEPT
Notes: Accept. The agreed fixes capture the real correctness issues before merge, especially append-only validation for manual/resume handoffs.

### Signoff: claude — 2026-05-14
Status: ✅ ACCEPT
Notes: Accept. All round-01 findings I raised are routed correctly: the CRITICAL append-only bypass for manual/resume is in agreed fixes with the right shape (persist pre-handoff hash, re-validate in resume), both MAJORs and all MINORs/NITs are agreed, and my open questions are resolved — `before_sha256` is implicitly handled by the CRITICAL fix, mixed-mode validation and `--allow-interactive-noninteractive` are appropriately deferred.

### Signoff: gemini — 2026-05-14
Status: ✅ ACCEPT
Notes: Accept. My findings on Round 2+ validation and manual exit semantics are addressed in agreed fixes. Deferring signal handling for `spawn-tty` and mixed-mode resume polish to subsequent slices is appropriate for this foundational implementation.

### Signoff: hermes — 2026-05-14
Status: ✅ ACCEPT
Notes: Accept. Round-01 findings on hosted confirmation and validation are addressed by the agreed append-only fixes and dismissals; all critical path issues for manual/resume are now covered.
