---
agent: claude
idea: consensus-request-signoffs
round: 1
date: 2026-05-13
---

## Summary

Add a `parley consensus request-signoffs [--review] --participants IDS [--yes] [--dry-run] IDEA` command that, for a given idea, resolves either the design (`consensus.md`) or review (`review/consensus.md`) file, determines which participants still owe a signoff, and invokes their configured headless agents one at a time so each appends its own canonical signoff block. After each invocation the command re-parses the consensus file to confirm exactly one new, well-formed signoff appeared for that participant; otherwise it stops with a diagnostic. `--dry-run` prints the plan without invoking; `--yes` is required before any hosted/external agent is launched. The command reuses `internal/consensus` for parsing/status and the existing runtime config + command construction used by other parley commands, keeping the file-tree under `parley-deck/` as the source of truth and staying within the constraints (sequential, no locking, no GitHub review API, no auto-mode across phases).

## Proposed approach

Success criteria for the slice:
- Running with `--dry-run` prints participants in resolution order, target consensus file, per-agent launch command preview, and exits 0 without invoking anything.
- Running with `--yes` invokes each selected agent sequentially; after each, the consensus file has gained exactly one canonical signoff for that participant.
- Stops with a non-zero exit on invocation failure, missing signoff, duplicate/malformed signoff, or a `❌ BLOCK` outcome observed during re-parse.
- `--review` switches resolution to `review/consensus.md` and is the only change in behavior between design and review modes.

Plan:

1. CLI wiring → verify: `parley consensus --help` lists `request-signoffs`; flag parsing rejects missing `IDEA` and missing `--participants` when no safe default is available.
   - Add subcommand under the existing `consensus` group. Flags: `--review` (bool), `--participants` (comma-separated IDs, required unless a safe default exists), `--yes` (bool), `--dry-run` (bool). Positional: `IDEA`.
   - Mutual handling: `--dry-run` overrides `--yes` requirement (no invocation happens). If neither `--yes` nor `--dry-run` is set and any selected agent is hosted/external, fail early with a clear message.

2. Resolve target file → verify: function returns the correct path for design vs review and errors clearly if the file is absent.
   - `--review` ⇒ `parley-deck/ideas/<IDEA>/review/consensus.md`. Otherwise `parley-deck/ideas/<IDEA>/consensus.md`.

3. Parse + decide participants → verify: parser returns the structured signoff state; selection logic produces a deterministic ordered list.
   - Reuse `internal/consensus` to parse the file into the existing signoff/status model.
   - Compute the set of participants that have *not* yet signed off (the "missing" set), preserving file/declaration order.
   - If `--participants` is provided: use it as-is, but validate every ID is a known participant for the idea; reject unknown IDs.
   - If `--participants` is omitted: default to the missing set only when it is non-empty and unambiguous; otherwise require explicit `--participants`. Document this as the "safe default".

4. Hosted/external gating → verify: hosted participants without `--yes` (and without `--dry-run`) cause an early, descriptive failure listing them.
   - Reuse runtime config to classify each participant's agent (local vs hosted/external).
   - If any selected agent is hosted/external and `--yes` is not set and `--dry-run` is not set, abort before any work.

5. Build per-agent invocation plan → verify: dry-run output matches the would-be commands and order.
   - Reuse existing command construction (the same path used elsewhere to launch a configured agent for an idea/round) to produce a `(participant, command, env, workdir)` tuple per target.
   - `--dry-run` prints: target file, ordered participant list, and the full launch summary per participant, then exits 0 with no side effects.

6. Sequential execution loop → verify: each iteration ends with exactly one new canonical signoff for that participant, or the loop aborts.
   - For each participant in order:
     a. Snapshot consensus file (in-memory parse + content hash) to detect "before" state.
     b. Invoke the configured agent for that participant, streaming stderr/stdout to the user.
     c. On non-zero exit, stop and report.
     d. Re-parse the consensus file. Validate:
        - the participant's signoff is now present,
        - it is canonical (well-formed block, single instance, expected fields),
        - it does not duplicate or overwrite another participant's block,
        - status is not `❌ BLOCK` (if BLOCK appears, stop and report).
     e. Continue to next participant.

7. Stop conditions and exit codes → verify: each failure mode produces a unique, descriptive message and a non-zero exit.
   - Invocation error, missing signoff after invocation, duplicate/malformed signoff, BLOCK status, and unknown/unselectable participants are each distinct failure paths.
   - Partial progress is preserved on disk (prior successful signoffs remain committed by their agents); no rollback in this slice.

8. Out of scope (explicit in the brief, restated for clarity): full auto across phases, GitHub native review API, GitLab MR automation, cross-process file locking, generated consensus prose, release packaging.

## Concerns / open questions

- "Safe default from missing signoffs": should the default include *all* missing participants, or only those whose agent is `local`? Hosted defaults plus omitted `--yes` would force a re-run. My current proposal includes all missing and lets the hosted/external gate fail early — but explicit confirmation is welcome.
- Hosted/external classification: which field in the existing runtime/agent config marks an agent as hosted vs local? I assume one exists since the brief differentiates them; need to confirm the exact predicate to avoid mis-gating.
- Canonical signoff shape: what is the authoritative regex/spec used by `internal/consensus` to determine "well-formed" vs "malformed"? The verification step must use the same predicate the parser uses — not a second, drifting one.
- Working directory and branch: the live branch is `idea-consensus-request-signoffs` while idea paths use `consensus-request-signoffs`. Should the command pass any branch context to invoked agents, or does the existing command construction already handle this? I'd prefer to inherit from the existing launcher rather than re-implement.
- Streaming vs capturing output: do we stream the child agent's logs live, capture to a per-round log file, or both? Affects how failures are surfaced.
- Per-invocation timeout: the effective launch timeout is 1800000 ms for *this* task; should the per-agent invocation inherit the parent's remaining budget, or take its own flag? I propose deferring a flag and inheriting until a user asks otherwise.
- Re-runs and resumption: if invocation #3 of 5 fails, can the user simply re-run with the remaining participants? My design says yes (the file is the state), but it's worth confirming we don't want a `--resume` affordance in this slice.
- Ordering determinism: when `--participants` is omitted, is "file/declaration order" the right canonical order, or should it be alphabetical? Picking file order to mirror authoring intent, but happy to switch.
- What exactly counts as duplicate detection — same participant id appearing twice, or also two blocks of the same shape from different agents? I'm assuming "same participant id, more than one signoff block."

## Risks

- Agents append duplicate or malformed signoffs: mitigated by mandatory re-parse + canonical-shape check after every invocation, with hard stop on the first offense. Residual risk: the parser's "canonical" definition diverges from what an agent emits; that should be caught in the existing parser tests rather than re-asserted here.
- Hosted agents invoked unintentionally: mitigated by requiring `--yes` and printing a hosted-agent list before launch. Residual risk: a misconfigured agent classification slips through; a unit test on the gate using fixture configs would catch the common case.
- Partial state on failure: agents commit/append independently, so a mid-list failure leaves earlier signoffs in place. This matches the "files are source of truth" model and is recoverable by re-running with the remaining participants, but users may be surprised — the dry-run output and the failure message should both make this explicit.
- Concurrent operator runs racing on the same consensus file: out of scope (no file locking), but two operators running this command at once could interleave writes. Documenting the limitation in the command's help text is cheap insurance.
- `❌ BLOCK` from a later participant invalidates the implicit "agreement" implied by earlier signoffs in the same run: the slice's behavior is to stop and surface the BLOCK; deciding what to do with the earlier signoffs is a policy question for a later slice, not this one.
- Cost/time amplification: sequential hosted invocations multiply latency and spend. The dry-run preview is the main mitigation; an explicit count + "this will invoke N hosted agents" line before launch would help, and is cheap to add.
- Drift between command construction here and the one used by other parley commands: mitigated by reusing the existing builder rather than re-implementing argv construction. If reuse turns out to be awkward, that's a refactor signal — not a reason to fork the logic in this slice.
