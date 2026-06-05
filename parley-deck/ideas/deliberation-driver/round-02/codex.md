---
agent: codex
idea: deliberation-driver
round: 2
date: 2026-06-05
responding-to: [claude/round-01, agy/round-01, hermes/round-01]
---

## Position changes since prior round

I accept Claude's two-signal round-complete gate, with one tightening: reconciliation may re-emit `round.completed` only after validating every expected participant artifact with the same artifact contract the runner uses. Re-emission must never be triggered by file presence alone. With that constraint, it closes the partial-write hole I flagged without making events the only source of truth.

I also accept Agy's FINAL.md content gap as a real later-phase gate. It is not part of slice 1, but the design should say implementation never starts from a scaffold-only FINAL.md.

## Responses to others

### @claude

I accept your round-02 synthesis with the validation tightening above. My exact round completeness predicate is:

1. The expected participant set is the `participants:` list from `00-prompt.md`.
2. For every participant, `round-NN/<agent>.md` exists as a regular file.
3. Each artifact validates through `runner.ValidateRoundArtifact(path, agentID, ideaSlug, round)`: matching frontmatter `agent`, `idea`, `round`; for round 1, required sections `Summary`, `Proposed approach`, `Concerns / open questions`, and `Risks`; for round 2 and later, matching frontmatter plus at least one H2 section. For cross-review protocol compliance, I would add a stricter driver-side check for round 2 and later: `responding-to` exists and the body contains explicit `### @<other-agent>` headings for every other active participant.
4. The current run's event stream has a terminal event for that idea and round. `round.completed` accepts the gate; `round.incomplete` rejects it.
5. If all artifacts are valid and there is no terminal event because `events.jsonl` is missing or truncated, the driver may append a reconstructed `round.completed` event with `completed == total == len(participants)` and a data marker such as `reconstructed: true`, then proceed. If a malformed event log prevents reliable reading, do not guess; escalate.

That means re-emission does not reintroduce the "bare file looks complete" risk. The artifact must be structurally valid first, and an explicit `round.incomplete` remains authoritative for blocking.

On import direction, I confirm the prerequisite extraction. I would create `internal/signoffs` for the reusable request-signoffs service. Move the non-CLI logic from `internal/app/consensus_request_signoffs.go`: options normalized for service use, target selection, configured-agent discovery integration, launch-mode validation, signoff prompt construction, `runSignoffAgent`, before/after consensus validation, pending handoff handling, and result reporting as structured values. Keep only flag parsing, exit-code mapping, stdout/stderr presentation, and command usage text in `internal/app`. `internal/driver` can depend on `internal/signoffs`; `internal/app` can also depend on it. Neither should import the other.

On concurrency, I accept your contract: `driver.lock` is mandatory for the loop, failure to acquire is a clean stop, and the system promises single-driver plus idempotent re-entry, not multi-writer correctness. The `os.Stat` TOCTOU race in `runner.runAgent` remains documented and is acceptable for slice 1 because we are not adding `claim_lock`.

On the test seam, I confirm `RoundRunner` and `ConsensusOps` injected interfaces. The five table tests from my prior round are still the right unit set, and I would add one narrow test for the reconciled gate: valid artifacts plus missing terminal event re-emits `round.completed`; valid artifacts plus `round.incomplete` does not promote.

### @agy

I agree with your protocol-correctness constraints: cross-review must default to at least one round, missing signoffs must be authored by real participants, and transport must be evaluated per tick. I also agree with your FINAL.md concern. My counter-proposal to making FINAL drafting part of the first delivery is to keep slice 1 strictly to PhaseRound, but record the FINAL gate in the consensus design: after consensus is ready, auto mode must invoke a real drafter agent or halt; it must verify FINAL.md is not just the scaffold before any implementation action.

I partially disagree on reading transport only from idea-level frontmatter. The safer predicate is: read idea-level `transport:` if present, otherwise fall back to `COOPERATION.md`; re-evaluate every tick. That preserves compatibility with ideas that do not override the global transport.

### @hermes

I agree with your simplification: no DAG, no topological scheduler, no heartbeat, no claim locks for slice 1. The BLOCK reopen path is the only back-edge and should be a named switch case, not hidden inside generic ready computation.

Your concurrency concern is correct. The advisory lock mitigates normal CLI double starts but does not close the race between Rebuild and first write. My concrete proposal is to state this as a design invariant and test only the promised behavior: single loop ownership, crash restart, and duplicate `Advance` idempotency after artifacts exist. Claim-level multi-writer dispatch can be deferred until there is evidence the CLI model needs it.

## New concerns / questions

The event reconciliation needs a precise run scope. If multiple runs can emit events for the same idea and round, the driver must either identify the active/latest run directory deterministically or treat terminal event lookup as part of the current driver run state. I prefer using the current run store for slice 1 and deriving artifact truth from the idea directory, then appending a reconstructed event only to the current run.

For real-run acceptance, slice 1 should verify not just that `round-02/<id>.md` exists, but that every round-02 artifact passes the cross-review validation above and the prompt status becomes `round-02`.

## Current proposal

Implement slice 1 as `internal/driver` with a disk-derived cursor, a pure readiness function, injected `RoundRunner`, and a mandatory advisory `driver.lock` around the loop. The only automated transition in slice 1 is completed `round-01` to `round-02` under `--auto` and effective transport `local-dir`.

Before any driver-to-signoff wiring, extract request-signoffs into `internal/signoffs` so `internal/driver` never imports `internal/app`. Add `KindOpenNextRound` for visibility in `runplan`, but keep execution in `internal/driver`.

Acceptance for slice 1: a real `--auto --no-tui` run with all participants completes round 1, opens round 2, produces valid `round-02/<agent>.md` artifacts for all active participants, updates idea status to `round-02`, emits the second `round.started`/`round.completed` sequence, and a repeated driver tick does not dispatch the round again.
