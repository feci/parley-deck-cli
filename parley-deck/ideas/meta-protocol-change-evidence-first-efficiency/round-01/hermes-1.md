---
agent: hermes-1
idea: meta-protocol-change-evidence-first-efficiency
round: 1
date: 2026-09-05
---

# hermes-1 — round 1 (liveness, process supervision, false-negative prevention)

## Summary

The five priorities are implementable and mostly touch prior ratified decisions, not new
territory. My lens is liveness and supervision: priority 4 (distinguish slow/buffered success,
malformed reply, auth/provider/process error, and a real hang) and the pilot's honesty gates
(priority 5). The core defect I own is that `hostedPONG` collapses a live-but-misparsed agent
into the same `unavailable:no-pong` verdict as a dead one, and that verdict feeds an
exclusion gate that can drop quorum. I propose one shared `LivenessClass` model for both the
preflight ping and the run watchdog, a sentinel-substring "malformed-but-live" tier, and a rule
that only a confirmed real-hang/stalled class may propose exclusion. Telemetry (priority 1) must
carry this class; I define the interface, codex-1 wires the producer.

## Existing alternatives (with locators)

- **Preflight PONG** — `internal/app/preflight.go`: `hostedPONG` (line 823) returns
  `(bool, string)` with reasons `missing / command-build-error / start-error / timeout /
  exit-error / no-pong`; `isExactPONG` (line 872) requires stdout whose sole non-whitespace is
  the token `PONG`. A bullet-prefixed or fenced correct answer (the kimi `• KIMI_OK` case in
  `ideas/preflight-liveness-false-negative/00-prompt.md`) is `unavailable:no-pong` — the exact
  false negative. The `gateExcludeAgent` gate (`preflight.go:363-368`) forwards this reason string
  verbatim into "confirm excluding it from this idea", so a misparse has an exclusion surface.
- **Run watchdog** — `internal/runner/supervision.go`: three layers (first-output grace 120s,
  stall 1800s, 60s heartbeats) with sentinels `errNoFirstOutput`/`errStalled` (lines 19-20);
  activity is tracked by counting writers / ACP events, never fs probing (lines 11-15). Already
  distinguishes "no output" from "stalled" from "process exited", but `no_first_output` fires on a
  still-running, buffered, eventually-successful agent (the retrospective's 35 `no_first_output`
  terminal events are the candidate false-kill set).
- **Failure classes** — `internal/runner/failclass.go`: ordered regex table `failureRules`
  (lines 24-43: rate-limit, auth, billing, overloaded, model-not-found, context-window, sandbox,
  budget, invalid-request) over the last 4 KiB of logs; `watchdogHints` (lines 48-53). This is the
  taxonomy to reuse, not reinvent.
- **Cost telemetry gap** — `internal/driver/loop.go:183` reads `agent.usage` events
  (`loopCostUSD`, test at `loop_budget_test.go:71`), but no production code appends `agent.usage`
  (grep confirms the only writers are tests). Loop cost is a reader with no producer.
- **Path-root false negative** —
  `inbox/claude-1-to-all_agents-verify-hermes-probe_root-cause-found.md`: probe/artifact writes
  landed under `$HOME` when the verifier looked under the repo; "did not create <path>" was true
  but useless. Same class as "unavailable:no-pong" — structurally accurate, materially misleading.
- **Ratified precedent** — `ideas/meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md`
  (packet contract, pre-registered experiment, §4.0 cell audit); `ideas/completion-contracts-evidence-ledger/FINAL.md`
  (`checks:` list form, driver-populated evidence).

## Proposed approach

### 1. One shared liveness model (priority 4 — my primary ownership)

Replace `hostedPONG`'s `(bool, string)` with a typed report, and reuse the same enum in the
supervisor so preflight- and run-level liveness cannot disagree:

```go
type LivenessClass string
const (
    ClassLive            LivenessClass = "live"            // valid PONG / output growth / ACP event
    ClassQuietSuccess    LivenessClass = "quiet-success"   // exit 0, no parseable output (process succeeded)
    ClassMalformedReply  LivenessClass = "malformed-reply" // exit 0, sentinel present but not exact shape
    ClassProviderError   LivenessClass = "provider-error"  // failclass auth/billing/overloaded/rate-limit/model
    ClassProcessError    LivenessClass = "process-error"   // non-zero exit / start-/exit-error, not provider
    ClassRealHang        LivenessClass = "real-hang"       // still running, zero output/events past grace
    ClassStalled         LivenessClass = "stalled"         // had output, then silent past stall window
)

type LivenessProbe struct {
    Class       LivenessClass
    ExitCode    int
    SawSentinel bool          // "PONG" substring present but isExactPONG false
    StdoutTail  string        // <=256B, sanitized, never secrets
    StderrTail  string        // <=256B, sanitized
    Duration    time.Duration
}
```

Key addition: `SawSentinel`. `isExactPONG` stays exact for answer-shape purposes; liveness gains a
`containsSentinel` check. A reply that contains `PONG` inside a bullet/code-fence proves the agent
ran — classify `ClassLive` (or `ClassMalformedReply`), never `unavailable`. This maps the kimi
probe, the `• KIMI_OK` shape, and any wrapper prefix to "alive but parse-shape unexpected".

`hostedPONG` then assigns: exit 0 + exact PONG → Live; exit 0 + sentinel substring → MalformedReply
(SawSentinel); exit 0 empty → QuietSuccess; non-zero exit → classify exit/stdout/stderr through
`classifyFailure` (failclass.go) and map auth/billing/overloaded/rate-limit/model-not-found →
ProviderError, else ProcessError; timeout with zero activity → RealHang; timeout with prior activity
→ Stalled.

### 2. The exclusion gate must not fire on misparse (false-negative prevention)

Change the `gateExcludeAgent` construction (`preflight.go:363`) so **only** `ClassRealHang` and
`ClassStalled` (and `ClassProcessError` with a non-provider non-zero exit) gate an exclusion.
`ClassMalformedReply`, `ClassQuietSuccess`, and `ClassProviderError` render as `unknown`/`live`
warning rows — not a gate, not a proposal to drop quorum. This satisfies both sides of the
`preflight-liveness-false-negative/00-prompt.md` constraint: it distinguishes "did not answer"
(RealHang) from "answered in a shape we did not parse" (MalformedReply), and it does **not** loosen
the check into never firing — a genuinely dead agent still gates.

### 3. Supervisor: stop killing eventual success (priority 4)

`waitSupervised` (`supervision.go:141`) already returns the process's own nil error when the child
exits before the grace window — a quiet-success that *exits* is not killed. The remaining
false-kill is a still-running, buffered agent that will emit and exit after 120s. Two bounded
changes: (a) before a `no_first_output` kill, snapshot the full state (stdout/stderr byte counts,
process-alive, any ACP `MarkEvent` time) via the existing `onWatchdog` hook — already ordered
before `kill()` (D1); (b) when the state shows recent non-output activity (ACP events within a
short rollup) do not kill — a stall, not a silence. No change to defaults without the §4.0
enforcement-audit hold; this only enriches the signal already recorded.

### 4. Path-root honesty (cross-cutting, feeds priorities 1 and 4)

The verifier's "did not create <path>" must print (i) where it looked, (ii) the explicit cwd and
any `$HOME`-relative candidate of the same name it found elsewhere. Probe/artifact writes already
set `cmd.Dir = root` (`preflight.go:834`); add a guard test that a probe invoked without an
absolute target writes under the repo, not `$HOME` — the inbox root-cause regression test.

### 5. Telemetry surface I consume / define (priority 1 boundary)

Telemetry records a stable `run_id`+`attempt_id` and, for every attempt, the `LivenessClass` plus
exit code and sanitized tails. Secret/raw-env never enters the ledger. `agent.usage` has no
producer (`loop.go:183` reads only); codex-1 owns emitting it from every launch path, including
manually-facilitated headless calls, with explicit nulls. I consume `LivenessClass` into the HTML
limitations table: counts of quiet-success vs malformed vs real-hang vs provider-error, so "48
failures" becomes honest rather than an aggregate.

## Experiment controls & acceptance tests (my lens)

**Fixtures (canary set for the pilot, priority 5).** Three fake probe fixtures plus a control,
wired through `pingProbe` like the existing `withFakeProbe` (`preflight_test.go:380`):
(a) slow-but-successful — sleeps past the first-event window, then emits exact PONG; must classify
Live, no gate. (b) real-hang — never exits within timeout; must classify RealHang and gate.
(c) malformed-reply — emits `• PONG` or a fenced `PONG`; must classify Live/MalformedReply, no gate.
(d) provider-error — exits 1 with an auth string; must classify ProviderError, no gate.

**Acceptance tests (unit, in-tree):**
1. `hostedPONG` fixture table asserting the class map above (extends `preflight_test.go`).
2. `isExactPONG` unchanged byte-for-byte; new `containsSentinel` is what liveness uses — `PONG`
   inside `PONGXYZ` still rejected exact, accepted for liveness.
3. `gateExcludeAgent` fires only for RealHang/Stalled/non-provider ProcessError; malformed/quiet/
   provider produce a warning row, exit 0, no `[exclude-agent]`.
4. Supervisor: a process exiting 0 at T>grace is not `no_first_output` (returns its nil wait
   error); a process with ACP events but no stdout is not killed at the first-output window; a
   process silent with zero ACP events is killed and classified RealHang.
5. Path-root guard: probe with a relative target resolves under repo cwd, and the verifier's
   not-found message names both the searched path and the `$HOME` candidate when one exists.
6. Read-only retro reclassification: a CLI over the structured events emits counts of the 35
   `no_first_output` terminal events by new class — comparison set only, no mutation, no deletion
   of historical artifacts.

**Pilot controls I insist on:** the pilot's per-arm task budgets are equal and the roster is the
six active machine-roster IDs, not this idea's four-person quorum (per kickoff); liveness fixtures
run through the *same* binary invocation a real run uses (`CommandFor`), never a shortcut ping.
Genuine-hang and quiet-success fixtures must be distinguishable to the evaluator in the HTML —
otherwise "honest limits" is unverifiable.

## Work boundaries

- **I own:** priorities 4 in full (shared model, probe split, gate gating, supervisor signal, path
  honesty), and the pilot's liveness canary + honesty check.
- **codex-1 owns:** priority 1 telemetry producer + run/attempt IDs + priority 5 experiment
  validity/statistics; I only define the `LivenessClass` field telemetry must carry.
- **claude-1 owns:** priority 2/3 (packet + instruction-layer + iteration budgets); no §4.0 cell
  edits by anyone absent the `track-gate-enforcement-audit` (kickoff + phase-packet FINAL).
- **kimi-1 owns:** priority 2 independence/scope differentiation + adversarial fixtures. I will
  not author completion-scope or independence attestation in their stead.

## Concerns

1. `SawSentinel` substring could re-admit a hostile prompt-echo — but for *liveness* an echo of
   `PONG` is precisely evidence the agent is running; correctness-of-answer is not this layer's job.
   Keep the two connotations in separate fields.
2. "exit 0" is not "live": a wrapper that crashes cleanly with empty stdout is QuietSuccess, which
   is not excludable but must be surfaced, not silently passed. QuietSuccess stays a warning row.
3. The `no_first_output` grace extension can itself reintroduce a hang if the ACP-event rollup
   never decays; bound it (single narrow extension), and keep the hard per-agent timeout as the
   ultimate kill.
4. Provider errors (auth) are environment failures, not agent liveness — reporting them as
   "available=unknown" rather than "available=false" changes what a facilitator sees; the freshness
   distinction must survive into the HTML so a provider outage is not read as dead roster.

## Risks

1. **Loosening-to-never-fire is the mirror defect.** If the gate gating is over-conservative and no
   real-hang ever gates, we've rebuilt `trigger-that-cannot-fire`. The fixture canary (b) exists
   precisely to prove the gate still fires.
2. **Exclusion surface is shared code.** `runTaskPreflight` and `parley preflight` both call
   `preflight`; changing the gate must not diverge them (the deck's recurring "rule binds on only
   one of two entry points" defect, `preflight.go:195-207`).
3. **Historical honesty.** The 35 `no_first_output` events may include true kills and false kills;
   relabeling them without the reclassification tooling (test 6) would be the same fabrication the
   kickoff forbids. Do not edit the retrospective.
4. **Self-verification is not independence.** My fixture tests are unit tests, not the pilot's
   blind evaluation; kimi-1's adversarial pass and a non-implementer run of the pilot gate the
   liveness diagnosis independently.
5. **HTML honesty.** A "live" agent that was actually auth-blocked must appear as provider-error in
   the limitations table, or priority 1's provenance claim collapses and priority 5's "honest
   limits" is unmeetable.
