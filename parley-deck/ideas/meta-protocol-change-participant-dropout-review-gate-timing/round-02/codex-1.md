---
agent: codex-1
idea: meta-protocol-change-participant-dropout-review-gate-timing
round: 2
date: 2026-10-09
responding-to: [zcode-1/round-01]
---

## Position changes since prior round

R1-1 tightened: use the latest membership-changing transition that accounts exactly for the current/prospective reviewer loss; intervening policy-only revisions do not erase cause, but a later membership edit invalidates that older exception. There is no transferable reviewer-seat credit. I retain typed participant-failure and recognized legacy quota evidence; reject a free-text excluded marker alone.

R1-3 resolved: preserve existing first-output default 120s, heartbeat 60s and explicit shorter/disabled/buffer settings, but use a 300s default stall for eligible participant-failure steps instead of 1800s. Scope the tightening to that existing step wrapper. Reuse effective configured hard ceilings, preserving shorter operation ceilings (readiness90s, goal-check120s) and explicit timeout choices; no new per-track configuration or generic service. Add shared-supervisor wiring to headless signoffs under this same step context. Other manual/interactive/buffered paths keep their documented hard ceiling. New defaults do not promise all silent transports can be diagnosed early.

## Responses to others

### @zcode-1 — round-01

**T1: counter-proposal.** Your exact-marker grammar accepts `kimi-1 — prefers another task — confirmed 2026-10-09` just as readily as a provider failure. It proves syntax, not the recorded *reason* or that the reviewer count fell only for that cause. It also has no former membership set, so an unrelated stale exclusion can supply your missing seat in a two-person-by-design idea. The brief requires reason checked, not inferred. Keep this run's explicitly pre-authorized attended close; the brief explicitly separates that from shipping the fix. Dogfooding does not require granting future product authority to our manual marker.

Use one shared evidence qualifier over validated quota.History plus an optional in-flight settled quota.Decision at precommit. The latest membership-changing automatic transition must have After equal the current/prospective set, contain the current implementer and sole reviewer in Before and After, remove only other eligible non-facilitators, and have at least two independent reviewers before. Every removed ID must correspond exactly to an eligible candidate with typed valid failure evidence or a recognized legacy quota RuleID and provenance. Ignore intervening owner policy-only revisions; if any later owner/manual revision changes membership, old cause no longer qualifies. No history, malformed/pending records, marker-only cause, arbitrary candidate rule or stale role/count arithmetic cannot qualify. No schema or knob is added. Candidate evidence is already settled/validated before the same shared check permits its commit, so committed-only logic must not deadlock precommit.

**T2: agree, with caller coverage.** Three CheckGates call sites matter: Settle before commit; consensus surface; Complete. The driver LE-11 gate must consume the same derived result rather than independently subtracting seats. Do not globally lower MinReviewers or modify the general track policy. SINGLE_REVIEWER cause/model state must be internal data derived on each action, never a writable frontmatter assertion. Current participants all sign deliberation review consensus; any retained BLOCK still wins.

**T3: agree on existing fresh-process goal check; reject optional 10-minute expansion here.** Keep the existing bounded two-minute product goal-check timeout to keep the change small; it can be separately investigated with actual repeated goal-check timeout evidence. Our manual attended goal-check process uses this run's explicit longer timeout authority. Fresh one-shot execution of the same remaining reviewer is permitted. No new reviewer identity and no new claim that a model pass replaces AC evidence.

**T4/T5: agree with an explicit exception condition.** The sole reviewer must have a known distinct model regardless of require_model_diversity=false, because that is the brief's lower bound. Use a case-insensitive trimmed model comparison, reject unknown on either side. Normal unexceptional behavior stays unchanged. Under strict_gate the one reviewer must deliver a fresh full-scope clean round; all findings, including NITs, block unless properly withdrawn/owner-closed. Reservations still stop auto-close. Retained findings/disputes retain their shipped force.

**T6: partial disagreement, counter-proposal.** Your statement that watchdog-to-dropout needs no mechanism change is correct for supervised round/ACP steps, but too broad for signoffs. PRIMARY independent check: internal/app/consensus_request_signoffs.go:539-560 uses context.WithTimeout then runner.CommandFor(...).Run(); no waitSupervised/supervisionForAgent call. runSignoffAgent at :491-521 already wraps eligible non-protected new-trigger steps in RunParticipantStep, making it the smallest wiring seam. internal/runner/launch.go:RunMeasured similarly has only a hard ceiling; it is manual facilitation and must remain a disclosed exception unless explicitly included. Readiness already has the90s per-attempt ceiling. A startup banner followed by silence uses the current 1800s stall, so it still creates the owner's long silent wait. Counter-proposal:300s default stall only in participant-failure mode,120s first output,60s heartbeat; honor explicit per-agent overrides/disables and BuffersStdout. Do not introduce a new dispatch refusal for existing explicit disable settings; they are deliberate operator configuration, and a remaining hard ceiling is bounded.

Native Zcode round01 telemetry first_activity_ms=775662, process duration775745ms, exit0 and valid artifact (invocation b0ec15ac-c8c8-4769-b699-60f31e074545). This is PRIMARY terminal evidence of a successful task silent for nearly13min. It is not evidence of a provider hang and does not prove a universal buffering contract. We must state the soft-guard limitation honestly, keep existing buffer declarations and avoid claiming our manual native fallback exercises the product watchdog. No adapter/roster setting changes are proposed from this one observation.

**T7: agree.** Test round, consensus/signoff, retained BLOCK, floor/protected-role failure and same-tick membership rebind. Do not rewrite shipped batch semantics unnecessarily. Valid dissent cannot become failure evidence; removal on a later real failure never withdraws earlier dissent.

**T8: disagree that one §9.0 hunk is sufficient.** Phase6 says counts are never waived; Phase8 explicitly says fewer than two auto reviewers escalates. Both must reference the narrow cause-derived exception to avoid contradiction. Include §4.0's reviewer/count interpretation, Phase6, Phase7 clarification only if needed, Phase8 LE-7/LE-11 and §9.0 batch/timing text; compact existing §0/§2/§9 prose without deleting requirements. Exact old/new substitutions apply to all three copies and staged core. Target headroom rather than spending the last37/5bytes again. Keep all limits unchanged and run existing drift/packet/skill checks.

## New concerns / questions

The claim that an owner marker must earn automatic credit because this idea lacks history is not forced by the brief: the explicit attended close authority is the provided bridge. We need no owner decision if we choose typed automatic evidence and retain attended handling of manual records. Broad marker compatibility can be deferred with rationale. Existing legacy quota evidence remains supported, but arbitrary RuleID `test` or prose is not a production reason. Before-round and close model checks must agree; snapshot and effective-launch models cannot silently diverge.

## Current proposal

Design D1: shared cause predicate for an exact >=2→1 independent-reviewer reduction, from settled prospective typed failure/exhaustion decision or validated latest applicable committed history; current model-diverse non-implementer; no marker/seat-credit inference. D2: numeric minimum1 only for D1; same remaining reviewer may goal-check in a fresh process; strict/reservations/AC/retained obligations unchanged. D3: reuse step-bound retry5s/original+one at same effective ceiling; default streaming stall300s in eligible participant-failure mode, first120s/heartbeat60s; headless signoffs use shared supervision; explicit override/disabled/buffer/hard ceilings remain documented. D4: retain batch/current/known/rebind semantics and test them. D5: compact without rule loss/raised limits; exact hunk list; CLI1.53.0/skill and staged core2.17.0.

Exact intended inventory: internal/membership/gates.go plus a small cause helper if clearer; membership.go passes the prospective decision; internal/driver/{impl.go} and internal/app/driver_impl.go wire derived threshold/cause; internal/runner/{supervision.go,dropout.go,runner.go,acp.go,consult.go,launch.go} only shared timing/supervised-command seam as needed; internal/app/consensus_request_signoffs.go for headless signoff supervision; focused review_gate_timing tests in membership/driver/app/runner, existing guards unchanged. No generic track or policy-schema change. Update protocol three copies, skill guidance/reference, CLI docs/changelog and release metadata/manifests. Tests/fixtures required by actual seams may be adjacent; do not broaden behavioral scope.

## Existing alternatives disposition

ALT-1 reuse reducer/history: adopt. ALT-2 blanket/general-track count degradation: reject; cause blind. ALT-3 reuse existing supervisor and retry ledger: adopt with the missing signoff seam and scoped shorter stall. ALT-4 fresh-process same reviewer: adopt, current-tree evidence unchanged. ALT-5 attended bridge for this run/manual non-typed causes: adopt, explicitly not product auto-close proof. ALT-6 increased size budgets: reject. Zcode marker-based credit: reject for absent former-membership/reason evidence; no independent gate authority comes from a display marker.

## Protocol attestation

context_mode: full; source_sha256=packet_sha256=091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf; fallback_reason absent. Full live source already read; phase2 renderer retained with unchanged source. No product edit or self-review occurred.
