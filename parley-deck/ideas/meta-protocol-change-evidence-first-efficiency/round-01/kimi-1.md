---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
round: 1
date: 2026-09-05
lens: independent evidence, completion scope, adversarial testing
---

# Round 01 — kimi-1: evidence-first efficiency, five priorities

## Summary

All five priorities are tractable against the current code, and three of them extend machinery
that already exists rather than inventing new systems: the runner already classifies terminal
attempts (`internal/runner/runner.go:631` emits `agent.finished`/`agent.failed` with
`failure_class`, `attempt_id`, `segment_id`), the manifest already freezes launch identity
(`internal/runmanifest/manifest.go:49-79` RosterSnapshot), and the completion-contract veto
already binds list-form `checks:` to an all-pass at current HEAD. The genuinely new builds are:
an `agent.usage` producer (the consumer exists and starves), the phase-scoped packet (Slice 2 of
`meta-protocol-change-phase-packet-and-fixup-budget`, ratified but never started), and the frozen
pilot. I propose five bounded workstreams with a hard sequencing dependency: telemetry first,
because the pilot and the cost ceilings measure through it; liveness early, because the pilot's
six-roster-ID rotation is invalid if a live agent can be read as dead.

**Evidence basis.** This round I *inspected* only: `parley-deck/COOPERATION.md` (full),
the FINAL/IMPLEMENTATION of `meta-protocol-change-phase-packet-and-fixup-budget`, FINALs of
`completion-contracts-evidence-ledger` and `meta-protocol-change-verification-integrity`,
`ideas/preflight-liveness-false-negative/00-prompt.md`, the hermes-probe inbox note, and targeted
reads/greps of `internal/app/preflight.go`, `internal/driver/loop.go`, `internal/runner/{runner,acp}.go`,
`internal/runmanifest/manifest.go`, `internal/store/events.go`, `internal/protocolcore/*.go`,
`parley-deck/agents.toml`. I executed no tests, no builds, no probes, and mutated nothing.
Every code claim below carries the locator I read; none is from memory.

## Existing alternatives (with locators)

- **Cost reader without a producer.** `internal/driver/loop.go:174-186` `loopCostUSD` sums
  `agent.usage` events; its own comment says "the runners do not yet emit agent.usage, so this
  is 0 in practice". Worse, `internal/runner/acp.go:389` emits `agent.acp.usage` — a name the
  reader never matches, so even existing ACP usage is invisible to the LE-6 budget. Any telemetry
  work must fix this seam or cost dashboards will read a confident zero.
- **Attempt identity and terminal classification exist.** `internal/runner/runner.go:495-624`
  (exec attempt loop, retry-once only for `no_first_output`) and `:700-730` (`agent.finished`/
  `agent.failed` carrying `attempt_id`, `failure_class`, `recovery_hint`, exit code). The stable
  attempt ID the kickoff demands is already the runner's unit; telemetry should reuse it, not mint a new one.
- **Requested-vs-reported identity anchor exists.** RosterSnapshot
  (`internal/runmanifest/manifest.go:66-79`) freezes adapter/model/effort/launch args at run
  creation, with `RosterRevisionOf` hashing them. "Requested configuration" is already recorded;
  the missing half is provider-reported identity per attempt.
- **Preflight taxonomy exists but conflates.** `internal/app/preflight.go:823-866` `hostedPONG`
  returns exactly three reasons — `unavailable:timeout|exit-error|no-pong`; `:872` `isExactPONG`
  demands the sole sentinel. `exit-error` lumps auth, provider and process failure; `no-pong`
  lumps malformed reply with silence; `timeout` lumps a real hang with slow-but-buffered output.
  The kimi false-negative is documented in `ideas/preflight-liveness-false-negative/00-prompt.md:17-31`.
- **Completion contract exists.** `completion-contracts-evidence-ledger/FINAL.md`: named-list
  `checks:`, driver-written `## Validation evidence`, Phase-8 fail-closed veto at current HEAD.
  Extend it; do not build a parallel `done_when`.
- **Packet experiment is frozen.** `meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md`
  §3 pre-registers 6 paired runs in each of phases 1 and 6, AB/BA, three canaries plus control,
  zero obligation misses, ship at R≤0.50 both phases, disputed refute band (0.67,0.80],
  non-implementer recomputation. Its IMPLEMENTATION.md records Slice 1 (budgets) shipped with
  `HardCrossReviewCap` and `chargedFixupAttempts`; Slice 2 (packet) not started.
- **Idempotent effects ledger exists as a model.** §12.6-12.7 (COOPERATION.md:1143-1147):
  content-keyed, append-only `attempts[]`, reconcile-before-retry. The usage ledger should copy
  this shape.
- **Live protocol resolution exists.** `internal/protocolcore/render.go:77` `Render` and the
  release store in `core.go` — the "live resolved protocol, SHA-bound" source the packet contract requires.

## Proposed approach

Five workstreams, separately closable, recorded as slices in one IMPLEMENTATION.md — a slice's
failure must never be papered over by another slice's pass.

**W1 — Telemetry (blocks W5).** Emit `agent.usage` from every launch path: exec, ACP, consult,
preflight ping, and a `parley usage record` path for manually facilitated headless calls (marked
`provenance: manual-facilitated`, never averaged with driver-measured rows). Fields: `run_id`,
`segment_id`, `attempt_id`, `retry_of`, agent/adapter, `requested{model,effort}` (from
RosterSnapshot), `reported{provider_model}` (nullable), `tokens_in/out` (nullable),
`cost_usd` (nullable) plus explicit `cost_unknown: true` — never write 0 for unknown. Wall
start/finish from the runner clock. Reconciliation: every attempt reaches `terminal` or
`unresolved` with the reason; dedupe/import idempotency key = `(run_id, agent, attempt_id,
payload_digest)` per §12.7's shape. Unify the consumer: `loopCostUSD` reads both `agent.usage`
and `agent.acp.usage` during migration, or the ACP event is renamed with an alias.
*Acceptance:* ≥20 real launch attempts reconciled; a retried `no_first_output` attempt
contributes once; re-import of the same log changes no totals; unknown spend renders as
"unknown" in reports. *Adversarial tests:* duplicate-import probe, retry-double-count probe,
tampered event file → reconciliation flags mismatch rather than silently trusting
(driver-authored records only for any number that gates — the rule this deck recorded in the
fix-up idea: a number that is a safety boundary must not be authored by the party it constrains).

**W2 — Completion binding.** Extend the ratified contract, one vocabulary:
`design-accepted | implemented-unverified | verified | deployed`. `status: complete` requires
(a) list-form `checks:` all-pass at current HEAD (exists), (b) a `## Scope reconciliation` table
in IMPLEMENTATION.md mapping every `00-prompt.md` deliverable to an evidence locator or an
explicit `partial` declaration, and (c) for each material claim a non-owner `PRIMARY` verdict
per §15.1/§15.2 — the owner cannot close its own claim. Review attestations carry the commit
SHA they reviewed; a SHA ≠ current HEAD attests nothing about the current tree.
*Acceptance (negative tests, the classes the kickoff names):* a serial test labeled
"concurrency" is rejected; a skipped check cannot parse as pass; partial scope cannot close as
complete; stale-SHA review rejected; self-verdict rejected by construction.

**W3 — Packet + budgets.** Implement Slice 2 exactly as ratified: `parley protocol packet`
rendered from the live resolved protocol (`parley protocol check`'s source), `sourceSha256`
bound, verbatim selected blocks, complete omission index, always-include on unknown
applicability, full fallback recorded as `context-mode=full-fallback`. The frozen experiment is
run unchanged against the packet commit before release; the (0.67,0.80] band goes to the user
with the number. Budget extension is *coverage*, not retuning: the enforcement that exists
(HardCrossReviewCap, chargedFixupAttempts) binds driver runs; manually facilitated and resumed
rounds need a `parley budget status` read and an escalation path, and exhaustion escalates —
never closes. No §4.0 cell edits (the enforcement-audit gate stands).

**W4 — Liveness diagnosis.** Replace the three-bucket taxonomy with classified detail:
`timeout-no-output` (hang) vs `timeout-partial-output` (slow/buffered success killed at the
deadline — detectable because bytes arrived), `exit-error` split into auth/provider vs process
by stderr signature, `no-pong` split into `malformed-reply` (non-empty, unparsable) vs `silent`
(exit 0, empty). Rule change with teeth: a ping parse failure alone reports `unknown` and opens
a gate; it never proposes exclusion — and the gate must still be able to fire on real evidence
(the `preflight-liveness-false-negative` constraint: do not loosen into never-fire). Fixture
agents: a quiet-but-successful fake CLI and a genuinely hung one must classify distinctly.
*Acceptance:* scripted fake adapters for each class; the kimi-shaped reply (PONG with framing
text) is `malformed-reply`, not `unavailable`.

**W5 — The frozen pilot.** Twelve tasks × three arms (solo, duo, full six-ID machine roster),
equal per-task arm budgets enforced through the same MaxCostUSD path, role rotation across the
six roster IDs so model×task is not confounded, blind evaluation (evaluator sees anonymized arm
labels; a non-evaluator holds the key), all versions frozen pre-measurement: task set, rubric,
model/effort (RosterSnapshot), protocol `sourceSha256`, resource ceilings. No synthetic runs
presented as live; a null result is reported as null; a gate failure is a gate failure.
14/30-day durability is registered as follow-up observations, never simulated.

**Deliverable HTML:** offline, self-contained, keyboard-accessible, filterable, printable;
carries the evidence tables *and* the uncertainty (unknown spend, unresolved cells, disputed
band). Browser verification via ego-browser only.

## Experiment controls (my lens)

- **Freeze file before any measurement:** hashes of task set, rubric, prompt templates, roster
  snapshot, protocol source. Anything not in the freeze file is not evidence.
- **Blinding hygiene:** RosterSnapshot `LaunchArgs` can deanonymize an arm; strip identity-bearing
  fields from anything evaluators see.
- **Independent recomputation:** a non-implementer recomputes every reported ratio from raw event
  logs (the packet FINAL's own precedent). Raw logs ship with the HTML.
- **Pilot self-protection:** a hard total cost/wall ceiling and a kill rule, because the pilot
  itself is a spend.
- **Sequencing honesty:** the pilot cannot report cost until W1 lands, and cannot rotate safely
  until W4 lands. Report the dependency; do not silently substitute solo-arm data.

## Concerns

- **Scope mass.** Five priorities in one umbrella is the largest deliberate scope this deck has
  attempted. The mitigation is real slice boundaries with independent close criteria — and the
  explicit rule that an unexecuted experiment or a scaffold is not completion, per the kickoff.
- **The frozen packet experiment must not be retuned** under this umbrella; W3 is implementation
  plus the pre-registered run, nothing else.
- **Manual-facilitation telemetry is honor-system capture.** It must be labeled as such and
  excluded from any cross-arm average, or the pilot's equal-budget claim is void.
- **§15.6 notice:** the four of us are related models; convergence on "obvious" taxonomy splits
  (W4) is a shared prior. The fixture-based acceptance tests are the counterweight.

## Risks

1. **Event-name migration** (`agent.acp.usage` → `agent.usage`) can silently zero the LE-6 budget
   if only one side changes; keep the dual-read until the alias is proven by a test that injects
   both event names.
2. **Retry accounting**: the runner retries `no_first_output` once (`runner.go:495-503`); if
   `retry_of` lineage is dropped, the 20-attempt reconciliation overcounts. Boundary test required.
3. **Artifact-path verification must use absolute paths and explicit cwd** — the hermes probe
   wrote real artifacts under `$HOME` while the verifier looked in the repo
   (`inbox/claude-1-to-all_agents-verify-hermes-probe_root-cause-found.md`). Telemetry that
   asserts "file written" without the rooted check repeats the incident.
4. **Reversion/adversarial probes belong in isolated copies** — I record this against myself: in
   the fix-up idea I left a revert applied in the shared tree (its IMPLEMENTATION.md, cycle 2).
   All W1–W4 tamper probes run in `/tmp` module copies.
5. **Reviewer independence in the pilot:** blind evaluation fails if evaluators authored tasks;
   task authorship and evaluation must be disjoint roles in the freeze file.
6. **HTML deliverable drift:** "presentable" must not become an excuse to hide nulls; unknown
   cost and disputed thresholds render as first-class cells, not footnotes.
