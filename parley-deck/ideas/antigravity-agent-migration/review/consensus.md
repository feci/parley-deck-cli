---
idea: antigravity-agent-migration
review-cycle: 1
drafted-by: codex
date: 2026-05-27
participants: [codex, claude, agy, hermes]
status: complete
---

## Agreed fixes

- From `claude` MAJOR: keep the CLI non-TUI default participant behavior
  scoped to this migration by excluding only legacy `gemini`; do not also
  silently drop non-headless or ACP-capable installed agents in `installedAgentIDs`.
- From `claude` MAJOR: remove Hermes `--yolo` from built-in and shared defaults.
  The migration may use stronger verified model selection, but changing the
  Hermes safety posture needs a separate explicit decision.
- From `agy` MINOR and `claude` MINOR: document that `agy` intentionally does
  not enable isolated home by default because no stable Antigravity home
  environment variable was discovered; users who enable it must provide
  `isolated_home_env`.
- From `agy` MINOR: document in code that the `agy --print` prompt argument is
  intentionally last because the flag consumes the next value.
- From `agy` NIT: annotate the retained `gemini` host mapping as legacy in both
  live and embedded protocol files.

## Deferred or dismissed

- `claude` MINOR on duplicated Claude model/effort values is accepted as a
  trade-off: runtime metadata reports the effective intended defaults, while
  `HeadlessArgs` are the actual CLI launch contract.
- `claude` NIT on `acp_specs_test.go` map formatting is dismissed because
  `gofmt` applied normal alignment after the touched test file changed.
- `agy` open questions on future Antigravity model flags and Gemini removal
  window are release-note/follow-up items, not blockers for this patch.
- `agy` telemetry placeholder is deferred until Antigravity documents a stable
  telemetry/output mode.

## Fix-up cycle 1 verification

Review round 02 from `claude`, `agy`, and `hermes` verified all five agreed
fixes with no remaining `CRITICAL`, `MAJOR`, `MINOR`, or `NIT` findings. No
additional fix-up cycle is required.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex - 2026-05-27

Status: ACCEPT
Notes: The agreed fixes are narrow, testable, and directly address the
review-blocking risks without changing the migration direction.

### Signoff: claude - 2026-05-27

Status: ACCEPT
Notes: Both MAJORs I raised are captured in agreed fixes — the
`installedAgentIDs` scope restriction and the Hermes `--yolo` removal.
The HeadlessArgs duplication dismissal is a reasonable trade-off given
that runtime metadata and CLI args serve different consumers. The `agy`
`isolate_home` documentation fix and legacy annotation on the retained
gemini host mapping are adequate resolutions for the remaining items.

### Signoff: agy - 2026-05-27

Status: ACCEPT
Notes: The proposed fixes cover all raised issues, including documenting the isolated home caveat, clarifying the `agy --print` argument order, and annotating the legacy Gemini host mapping. Open questions and the telemetry placeholder are successfully deferred.

### Signoff: hermes - 2026-05-27

Status: ACCEPT
Notes: Agreed fixes resolve the scope and safety concerns cleanly. The migration stays focused and the legacy handling is appropriately documented.
