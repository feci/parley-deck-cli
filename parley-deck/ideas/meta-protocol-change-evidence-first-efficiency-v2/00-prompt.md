---
idea: meta-protocol-change-evidence-first-efficiency-v2
author: codex-1
created: 2026-09-15
track: deliberation
participants: [codex-1, claude-1, kimi-1, zcode-1]
status: round-02
amends: meta-protocol-change-evidence-first-efficiency
---

# Prospective Audit Quorum and Four-Participant Pilot Amendment

## User direction

The user approved: "Yes, both proposed steps." Translated from Slovak.
The question proposed replacing Hermes with Zcode in the remaining audit quorum,
keeping Codex as organizer, and preparing a full-six to full-four pilot amendment,
while preserving history, twelve tasks, equal ceilings, blind evaluation and the
separate context-packet experiment. The recorded answer authorizes both steps.
The user separately answered the explicit budget question: "15 minutes. USD 15."
Translated from Slovak shorthand. This means USD 15 TOTAL for prospective live
experiments and 15 minutes for one task in one arm, not USD 15 per invocation.
Implementation and review costs already incurred are not this experiment budget.

## Problem and narrow scope

The original signed FINAL D8 and AC-X1 require a full-six arm. Preserve that FINAL,
its kickoff, rounds and signatures unchanged. Produce a new explicit amendment
which supersedes only those prospective roster references. The remaining audit
quorum is now codex-1, claude-1, kimi-1 and zcode-1 under the user's express waiver.
Historical Hermes contributions retain their attribution; Zcode catches up on the
signed design and authors its own current review/signatures. No old observer note
is promoted to a quorum signature. This amendment is not implementation acceptance.

Use solo, duo and full-four arms across the original twelve task definitions.
Retain equal enforced ceilings, full protocol in every arm, counterbalanced
identities/pairs/drafters/critics, owned phase artifacts, blind nonauthor grading,
retained failures/timeouts and paired analysis. Preserve the separate packet
experiment exactly: six matched AB/BA pairs in EACH of phases 1 and 6, three
packet canaries plus a full control, existing ship threshold and decision bands.
The 15-minute task-arm limit is shared by all participants and phases in an arm,
including required orchestration; it is not reset at each call. All experiment
spending, including retries and grading, must fit USD 15 total. Unknown spend
must stop or carry a defensible enforced conservative reservation, never zero.

## Decisions to settle before freeze or treatment

1. Exact 12-task arm/identity/pair/role/order allocation and shared elapsed clock.
2. A feasible dollar allocation across 36 pilot cells, 28 packet/canary/control
   calls and grading. Show actual enforcement and worst-case bounding; a guidance
   word limit, post-hoc token total or CLI cost estimate is not a hard spend cap.
3. Two grading-only configurations disjoint from task/candidate authors, preserved
   blindness and explicit handling of grader disagreement within the same budget.
4. What to record if the unchanged experiment cannot fit these ceilings. Retain
   planned denominators and not-run cells; no silent scope shrink or extra spend.
5. Reuse existing preparation scripts and tests where applicable; freeze inputs,
   rubrics, models/efforts, source/test hashes and limits before measurement.

## Sources and existing alternatives

Read ../meta-protocol-change-evidence-first-efficiency/FINAL.md, especially D2,
D5, D8, AC-P2, AC-X1, AC-X2 and recovery. Read the separately binding packet FINAL
at ../meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md. All original
acceptance criteria except prospective full-six membership remain obligations.

Existing preparation lives at the sibling evaluation repository's
`delivery/2026-09-05/pilot/` and `scripts/pilot_{acceptance,grading,analysis}.py`.
Its DRAFT.md says it is preparation, not preregistration; no treatment calls.
The relevant DRAFT is copied to `source-context/pilot-preparation.md` before
participant launches, with its hash and origin. Hidden tests and grading keys
are not participant inputs and must never appear in candidate task workspaces.

## Startup and execution

Transport github-pr, dedicated idea branch based on original implementation
checkpoint f0e7d7b. Live source protocol is authoritative, full mode, SHA256
4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a.
Prior source-role status and sync dry-run made no metadata changes.
Three actual measured CLI implementation calls were still running when this
amendment opened; finish and verify them before reusing their participants.
Readiness comes from those actual completed calls or a new bounded check; do not
infer health from silence. Every participant owns round-01 independently, then
cross-review and its own signoff. No proxy authorship or internal subagents.

Configured choices: organizer codex-1 / Codex CLI gpt-6-astra; claude-1 / Claude
CLI claude-opus-5[1m] max; kimi-1 / Kimi CLI kimi-code/k3 with max from config;
zcode-1 / Zcode CLI zai/glm-5.3 with max from config. These are requested choices,
not a claim about the current organizer session's resolved model. The machine
codex-1 row currently says inactive; the explicit organizer authorization is
separate and does not mutate global membership. Each future launch must record
its effective configuration and reported identity honestly.

## Non-goals

No historical rewrite, global roster mutation, protocol-core publication, package
release, deployment, reduction of task count or packet pairs, retrospective
reattribution, default packet rollout, treatment before freeze or fake elapsed
follow-up. Original implementation PR #73 stays incomplete until its full gates
are satisfied. This v2 design changes the specified future experiment only.
