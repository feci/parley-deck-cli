---
from: user
to: all
idea: continuous-run-tui
date: 2026-05-23
status: authorized
---

# User-Authorized Continuation Exception

The user explicitly instructed the facilitator to continue until consensus and start implementation without waiting:

> pokracuj az kym nebude konsenzus a potom zacni implementaciu, necakaj

This follows failed real-agent attempts for the configured non-facilitator participants:

- `claude`: installed but not logged in.
- `gemini`: installed but missing configured auth in isolated `GEMINI_CLI_HOME`.
- `hermes`: installed but isolated profile has no valid model; non-isolated profile could not be used inside the current sandbox because it writes to `~/.hermes`.

This note authorizes a protocol exception for this idea only. The implementation may proceed from the facilitator-authored proposal, but the exception must remain visible in consensus, FINAL, IMPLEMENTATION, and review notes.
