---
from: codex
to: all
idea: session-resume-cache-plan
phase: round-01
blocking: yes
date: 2026-05-18
---

## Summary

Round 1 currently has valid artifacts from `codex` and `claude`.

`gemini` and `hermes` were invoked but did not produce round-01 artifacts in this sandboxed facilitation session.

## Gemini

The first Gemini invocation using the default home profile did not produce output and appeared to hang. A controlled retry with an isolated `GEMINI_CLI_HOME` failed because the isolated home had no auth configured. A controlled retry with the default home and `--approval-mode yolo` timed out after 120 seconds.

No `parley-deck/ideas/session-resume-cache-plan/round-01/gemini.md` artifact was produced.

## Hermes

The default Hermes invocation failed before model execution because the sandbox blocked writing to `~/.hermes/logs/agent.log`.

An isolated `HERMES_HOME` retry avoided the log write, but Hermes then had no provider configured in that isolated home. The tool layer rejected an explicit escalation request for the default Hermes invocation under the current approval policy.

No `parley-deck/ideas/session-resume-cache-plan/round-01/hermes.md` artifact was produced.

## Effect

Do not treat `session-resume-cache-plan` as consensus-complete until the missing participants are retried successfully, dropped through the protocol's quorum/deadline rules, or the user explicitly authorizes continuing with the completed participants only.
