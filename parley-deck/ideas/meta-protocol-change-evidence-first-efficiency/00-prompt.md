---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
created: 2026-09-05
track: deliberation
participants: [codex-1, claude-1, hermes-1, kimi-1]
roles:
  codex-1: runtime telemetry, integration, experiment validity
  claude-1: phase-scoped packet and instruction-layer correctness
  hermes-1: liveness, process supervision and false-negative prevention
  kimi-1: independent evidence, completion scope and adversarial testing
status: round-01
---

# Implement All Five Self-Assessment Priorities

## User Direction and Scope

The user directed: "Implement all the priorities you proposed and give me the
output in a presentable format, for example a single HTML." Translated from Slovak.
This authorizes implementation and verification of the five priorities below,
not an unattended production deployment, a new package release, a global-core
publication or changing existing operator model defaults. Those are not needed
to complete this task. Preserve the historical evaluation and unrelated work.

1. Capture model, time, tokens, estimated/actual cost provenance and retries for
   every supported launch path, including manually facilitated headless calls.
2. Bind "done" to independently checked acceptance criteria and delivered scope,
   not document existence or a self-issued pass. Distinguish accepted design,
   partial implementation, verified implementation and deployment.
3. Complete the already-ratified phase-scoped context packet and enforce iteration
   budgets, including manual facilitation, resume and consensus back-edges.
4. Fix liveness diagnosis before excluding agents: slow/buffered success,
   malformed reply, auth/provider/process errors and a real hang are distinct.
5. Actually run the frozen 12-task / three-arm pilot (solo, duo, full roster),
   equal per-task arm budgets, blind evaluation, role rotation and honest limits.
   Full experiment roster is the six active machine-roster IDs, not necessarily
   this implementation idea's four-person quorum. No synthetic/mock model runs
   may be presented as the requested live comparison.

Deliver a self-contained, offline single HTML with the historical baseline,
implementation evidence, live experiment results, limitations and remaining gates.
An unexecuted experiment, a scaffold or passing narrow tests is not full completion.

## Existing Evidence and Prior Decisions

Base code commit: `257ef8c75ac5e478edf42b59543e89cf94730de6`, CLI 1.47.0.
Installed skill 2.11.0 and all runtime installs are valid. The deck is protocolRole
source; its difference from the packaged reference is expected. Do not overwrite it.
Transport remains github-pr. No open native PRs at startup. Work happens in an
isolated sibling worktree; the ordinary main checkout stays undisturbed.

The frozen retrospective found 673 logical ideas, 463 duplicate instances, 5,178
participant artifacts, 29 legacy flat review artifacts. Review coverage is an
artifact proxy (400/453), not product correctness. Median canonical words grew
9,535 -> 16,620 -> 24,944 for June/July/August initiation cohorts. Age/task mix
prevents a causal quality/efficiency claim. Deduplicated structured events contain
86 terminal events, 48 failures, 35 classified no_first_output. This is a small
selected CLI subset, not all usage. Never rank models by filenames or word count.

The full frozen assessment is at:
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck-evaluation/assessments/2026-09-05/`.
Its REPORT.md, METHOD.md, KPIS.md and CASES.md are advisory evidence, not protocol.
The relevant propositions are copied above so the task does not require trusting
an inaccessible outside document. Further source-derived claims require locators.

Read the existing authoritative decisions before proposing replacements:
- `ideas/meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md`, especially
  the exact live-source packet contract, 6 paired runs in each of phases 1 and 6,
  AB/BA order, three canaries plus control, zero obligation misses, R<=0.50 ship
  gates and unresolved (0.67,0.80] speed region. Do not retune the experiment.
- Its IMPLEMENTATION.md: only the budget slice shipped, the packet did not.
- `ideas/preflight-liveness-false-negative/00-prompt.md` and its existing rounds.
- `inbox/claude-1-to-all_agents-verify-hermes-probe_root-cause-found.md`: historical
  relative-path writes under HOME. Use absolute target paths and explicit cwd.
- `ideas/completion-contracts-evidence-ledger/FINAL.md` and existing implementation.
- `ideas/meta-protocol-change-verification-integrity/FINAL.md` and protocol section 15.

Relevant current code: internal/runner (exec, supervision, ACP, consult, prompts),
internal/app/preflight.go (hostedPONG and isExactPONG), app/driver_checks.go,
internal/driver (loop cost reader exists but no agent.usage producer),
internal/protocolcore (live resolution), internal/store and internal/runmanifest.

## Acceptance Boundaries

- Telemetry uses stable run/attempt IDs and explicit nulls for unavailable fields;
  provider-reported/resolved identity is distinct from requested configuration.
  Secrets, raw environment and credentials never enter public artifacts or HTML.
- At least 20 actual launch attempts have a reconciled terminal/unresolved state;
  retries and repeated imports do not double-count usage. Unknown spend is visible.
- A serial test cannot claim concurrency, skipped checks cannot become pass,
  incomplete scope cannot close as complete, stale/different-commit review cannot
  attest the current implementation, self-verification cannot supply independence.
- Packet source is the live resolved protocol, SHA-bound, verbatim selected blocks,
  complete omission index, always-include on unknown applicability, full fallback
  on uncertainty. Preserve applicable section 6/14/15 and open objections.
  Preserve historical artifacts; this is not lossy history compaction.
- Budget exhaustion escalates, never closes; two consecutive confirmed material
  patch-induced regressions trigger a trajectory review in the pilot policy.
  Do not alter section 4.0 cells without its required enforcement audit.
- Quiet but successful and genuinely hung fixture agents are distinguishable;
  ping parse failure alone cannot silently alter membership/quorum.
- Live experiments freeze task/test/rubric/model/effort/protocol versions and
  resource ceilings before measurement. Neither gate failure nor a null result
  may be relabeled success. Real 14/30-day delivery durability is future follow-up,
  never simulated as elapsed; register follow-up observations if not yet due.
- Final HTML is offline, legible on narrow/low viewports, keyboard accessible,
  filterable, printable and carries evidence and uncertainty. Browser testing uses
  ego-browser only, never Chrome. No hidden remote assets or telemetry.

## Process

Design quorum is four active roster identities. Each writes its own independent
round-01 before reading the others. Give an integrated proposal, with extra depth
on your lens. Cross-review, canonical signoffs, FINAL and non-implementer review
are required. External agents must author their own artifacts. Implementation
file-set claims will be recorded in IMPLEMENTATION.md after design closes.

Existing unrelated open ideas remain open. This umbrella coordinates the requested
work; it does not retroactively close their missing artifacts or rewrite closed FINALs.
Global core publication remains a separate explicit human act. Normal native host
mirrors may be created through the configured GitHub account; it cannot approve
its own PR, so canonical agent signoffs must never be impersonated as distinct
GitHub user approvals.

## Non-goals

No package release, Homebrew/WinGet deployment, production mutation, credential
change, alteration of historical assessment inputs or unrequested global roster
change. No claim of statistically definitive model superiority from a pilot.
