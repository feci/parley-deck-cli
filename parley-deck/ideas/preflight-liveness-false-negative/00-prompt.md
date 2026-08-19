---
idea: preflight-liveness-false-negative
status: open
track: standard
initiator: claude-1
date: 2026-08-20
participants: [claude-1, codex-1, kimi-1, zcode-1]
rounds: 1
---

# §9.0 readiness reported a live agent as unavailable, and offered to drop it from quorum

## The observation

[PRIMARY, 2026-08-19, parley 1.45.0]

```
$ parley preflight
Roster:
  kimi   /Users/tomasfecko/.kimi-code/bin/kimi 0.36.1   no   unavailable:no-pong
Pending gates (require user confirmation):
  [exclude-agent] kimi unavailable (unavailable:no-pong) — confirm excluding it from this idea
```

The same binary, immediately afterwards:

```
$ ~/.kimi-code/bin/kimi -p "Reply with exactly: KIMI_OK" --output-format text
• KIMI_OK        (exit 0)
```

kimi then went on to file round-1, round-2 and signoff artifacts in a live idea. It was never
unavailable.

## Why this matters more than one wrong row

§9.0 exists so that a dead participant does not silently become a missing signoff. A ping that
reports a **live** agent as dead inverts the protection: it invites the facilitator to drop a
working participant and record the drop as readiness. The only thing that stopped it here is that
`[exclude-agent]` requires user confirmation — the gate held, the measurement behind it did not.

A false negative in a liveness check is worse than no check, because it is trusted.

## What this idea must establish

1. **What the PONG ping actually sends**, and whether kimi answers it. kimi's CLI has changed
   before: it now lives at `~/.kimi-code/bin/kimi` and is **not on PATH**. `-p` cannot be combined
   with `--auto` or `--yolo` — *"Cannot combine --prompt with --auto"* (PRIMARY) — so a ping that
   appends an autonomy flag to a prompt invocation would fail on kimi specifically.
2. **Whether the failure is timeout, exit code, or output-shape.** `--ping-timeout` exists;
   establish which arm fired rather than assuming.
3. **Whether other agents are affected.** `agy`, `opencode` and `zcode` were reported available in
   the same run; that does not mean the ping is correct for them, only that it did not fire.
4. **Whether a ping failure should be able to propose exclusion at all**, given that a wrong ping
   costs a participant while a slow one costs only time. Consider reporting `unknown` and letting
   the round decide, rather than proposing a drop.

## Constraint on the fix

Do **not** fix this by loosening the check until it never fires — an exclusion gate that cannot
fire is the mirror-image defect and this deck has already recorded one trigger that could not fire.
The fix must distinguish *"did not answer"* from *"answered in a shape we did not parse"*.

## Constraints

- Track `standard`: tooling defect, no protocol change.
- Repository READ-ONLY except each participant's own round file; verification in a COPY.
- Invoking agents to test the ping calls hosted backends. Keep probes short.
- English only. No secrets.

## Why this is filed separately

Found while running `roster-membership-overlay`, unrelated to its subject. The owner asked for it
to be its own idea.
