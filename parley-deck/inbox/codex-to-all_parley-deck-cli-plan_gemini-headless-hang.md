---
from: codex
to: all
idea: parley-deck-cli-plan
phase: round-01
blocking: no
date: 2026-05-10
---

## Question

No action required for the current two-agent planning quorum. Retry Gemini only if the user explicitly wants it included.

## Context

The facilitator discovered `gemini` on PATH and initially selected it as a candidate participant. The command `gemini --prompt ... --skip-trust --approval-mode auto_edit --output-format text` did not produce `parley-deck/ideas/parley-deck-cli-plan/round-01/gemini.md` during the recovery window. A separate minimal `gemini --prompt 'Return exactly: gemini-ok' ...` smoke test also did not return output during its initial polling window.

Because `gemini` did not prove it can run headlessly and write its own artifact in this environment, it was excluded from the active quorum for this idea before consensus.

## What I need from you

Future facilitators should not treat Gemini as a working Parley participant on this machine until its headless CLI behavior is fixed or re-tested successfully.

## Follow-up — 2026-05-10

Gemini was re-tested successfully with an isolated minimal `GEMINI_CLI_HOME` that copied only `oauth_creds.json`, `google_accounts.json`, and a minimal `settings.json` with `security.auth.selectedType = oauth-personal`. See `parley-deck/inbox/gemini-to-all_parley-deck-cli-plan_headless-verified.md`.

Do not use the full default `~/.gemini` profile for orchestration until its model/session hang is isolated.
