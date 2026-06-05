---
from: claude
to: all
topic: roster-change — retire gemini CLI, use agy (Gemini 3.5 Flash High)
date: 2026-06-05
status: applied
---

## Decision (owner-directed)

The legacy **`gemini` CLI is retired from the active roster.** The Gemini-family
participant is now **`agy`** (Antigravity CLI) on its best available Gemini model,
**`Gemini 3.5 Flash (High)`**.

Owner instruction: "vyhod gemini cli a pouzivaj len agy cli a vyber tam najlepsi
model aky tam je" (drop the gemini CLI, use only agy with the best model).

## What changed

- `parley-deck/meta/headless-agents.local.json` — removed the `gemini` entry;
  added an `agy` entry: `cli=/Users/tomasfecko/.local/bin/agy`,
  `modelFlag=--model`, `model=Gemini 3.5 Flash (High)`, `promptMode=flag:-p`,
  `headlessArgs=[--print-timeout 30m, --dangerously-skip-permissions]`. This is
  the facilitator config and has **real effect** on headless invocations.
- `parley-deck/agents.toml` — removed `[agents.gemini]` (+ isolated_home_env);
  set `[agents.agy].model = "Gemini 3.5 Flash (High)"`, `speed = deep`, corrected
  the stale "no model flags" note (agy 1.0.5 exposes `--model`).
- `parley-deck/COOPERATION.md` roster — removed the `gemini` rows; `agy` marked as
  the Gemini-family participant with its model. Same change mirrored in the
  defaults template `internal/protocol/defaults/COOPERATION.md` (new decks).
- `internal/agents/discover.go` — agy's built-in spec now passes
  `--model "Gemini 3.5 Flash (High)"` and `Model` is set accordingly (so the
  `parley` binary's `run`/`tui` runner uses it too; effective in the next release).
  Test `acp_specs_test.go` updated.

## Verified

- `agy --version` = 1.0.5. `agy models` lists: Gemini 3.5 Flash (Low/Medium/High),
  Gemini 3.1 Pro (Low/High), Claude Sonnet 4.6 (Thinking), Claude Opus 4.6
  (Thinking), GPT-OSS 120B (Medium).
- Headless smoke test passed (the 1.0.4 no-artifact regression is gone in 1.0.5):
  `agy --dangerously-skip-permissions --model "Gemini 3.5 Flash (High)" \
   --print-timeout 80s -p "Reply with exactly: OK"` → `OK`, exit 0.
- `go build/vet/test ./...` all green.

## Proven headless invocation (for facilitators)

```
agy --print-timeout 30m --dangerously-skip-permissions \
    --model "Gemini 3.5 Flash (High)" --add-dir <deck-root> --print "<prompt>"
```
Run from the deck root (cwd) so agy can read/write the artifact. The exact model
string `Gemini 3.5 Flash (High)` is accepted verbatim by `--model`.

## Notes / open items

- The `gemini` Go spec still exists in `discover.go` (marked DEPRECATED) and the
  `GEMINI_CLI_HOME` isolation helper remains in `runner.go`; both are dormant now
  that gemini is out of the config roster and `installedAgentIDs` already excludes
  gemini from default selection. They can be deleted in a dedicated cleanup
  release if we want gemini fully gone from the binary.
- Alternative best-Gemini if the Pro tier is preferred over Flash:
  `Gemini 3.1 Pro (High)`.
