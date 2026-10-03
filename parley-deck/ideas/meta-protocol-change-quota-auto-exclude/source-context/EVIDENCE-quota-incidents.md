# Quota-exhaustion incidents — evidence for meta-protocol-change-quota-auto-exclude

Collected by the owner's Claude Code session (the relay, not a participant) on 2026-10-02/03.
Provenance: PRIMARY where a command output is quoted, RECALL where marked. Response headers,
session IDs and correlation IDs are omitted on purpose; no credential appears anywhere below.

## 1. Today's preflight (2026-10-02 23:04 local, CLI 1.50.0, `parley preflight --dir . --json`, main checkout a8634cc)

PRIMARY. Duration 1:33. Roster rows, verbatim fields:

| rosterId | available | class | reason |
|---|---|---|---|
| claude-1 | true | ready | |
| codex-1 | true | ready | |
| kimi-1 | false | deadline-after-output | deadline-after-output |
| zcode-1 | false | provider-failure | provider-failure:rate-limit |

Gates printed by preflight (verbatim):

- `kimi-1 readiness is unresolved (deadline-after-output) — investigate and re-check; the agent is not excluded`
- `zcode-1 reports a provider failure (provider-failure:rate-limit) — resolve the provider/auth state; the agent is not excluded`

## 2. Direct zcode-1 probe (2026-10-02 23:06 local)

PRIMARY. Command: `zcode --prompt "Reply with exactly one word: PONG" --mode yolo --cwd <tmp> < /dev/null`, exit 1.
The CLI retried the same exhausted request at least four times (`isRetryable: true`) before failing.
One of the identical response bodies, verbatim:

```
statusCode: 429
'retry-after': '176930'
responseBody: '{"error":{"message":"[glm/glm-5.3] [429]: Weekly/Monthly Limit Exhausted. Your limit will reset at 2026-10-05 06:14:57 (reset after 49h 8m 50s)","retry_after":176930,"reset_at":"2026-10-04T22:14:57.003Z"}}'
Error: Turn execution failed
```

## 3. Direct kimi-1 probe (2026-10-02 23:06 local)

PRIMARY. Command: `/Users/tomasfecko/.kimi-code/bin/kimi -p "Reply with exactly one word: PONG"`.
After about 7 minutes it had written only `kimi version 0.42.0` to stderr and nothing to stdout; it was
stopped by the relay. No quota message was visible to the CLI caller.

Context, SECONDARY (another project's canonical artifact,
`librade-algoTrader/parley-deck/ideas/runtime-recovery-service-controls-lat55-90-92/00-prompt.md:98-102`):
"kimi-1's weekly Kimi quota is exhausted (until about 2026-10-03). The owner chose "claude-1 substitutes
for kimi-1"". That project ran a two-process claude-1 arrangement (a separate claude-1 organizer and a
separate headless claude-1 participant), with the role concentration recorded under §15.5.

## 4. windows-portability, 2026-09-28 — gpt-6-astra (codex-1) mid-implementation

PRIMARY, from the organizer log `runs-handoff/windows-portability-2026-09-24/organizer-resume-2.log`:

```
ERROR: unexpected status 503 Service Unavailable: [codex/gpt-6-astra] Unavailable (reset after 5h 51m 11s), url: https://omniroute.marao.sk/v1/responses
```

The codex-1 organizer died. Hypothesis to verify, not a conclusion: `internal/runner/failclass.go` classifies
this text as `overloaded` (the 5xx rule), not as `rate-limit` or `billing`, although it is a quota window.

## 5. windows-portability, 2026-09-28 — kimi-1 mid-implementation

SECONDARY, from `worktrees/windows-portability/parley-deck/ideas/windows-portability/organizer-notes.md:142`:
"Kimi launch first failed 400 ambiguous `k3`; read-only Ego inspection of
https://omniroute.marao.sk/dashboard/providers/kimi-coding now confirms its sole OAuth connection is
**disabled**, with "You've reached your 5-hour usage limit. Your quota will reset when the current 5-hour
window ends."" The CLI-visible error was an HTTP 400, not a quota message. The run then stopped all new
launches and has not resumed (its implementer zcode-1 is now also out of quota, see 2).

## 6. evidence-first-efficiency, 2026-09-05 — kimi-1 mid-implementation

PRIMARY, from `parley-deck/inbox/codex-1-to-user_meta-protocol-change-evidence-first-efficiency_kimi-quota.md`
(present in the windows-portability worktree deck): "Implementation attempt … returned exit 1 after 370.238
seconds. The provider reported HTTP 403 and a weekly usage limit." The organizer asked the owner for a
per-idea waiver; "Until an answer is recorded, the four-participant quorum is unchanged."

## Pattern

Four of the five incidents happened mid-idea, after the quorum had locked at Phase 0. Three different
providers expressed exhaustion three different ways: 429 with `reset_at`, 503 with "reset after", and 403
"weekly usage limit". One provider showed no quota text at all to the CLI caller (a 400, or a silent hang).
Each time, a run stalled until the owner answered or the quota came back.
