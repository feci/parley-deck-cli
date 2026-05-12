---
agent: codex
idea: consensus-workflow-cli
round: 1
date: 2026-05-12
---

## Summary

The first consensus workflow slice should add deterministic protocol helpers, not a full autonomous negotiator. I propose a small command surface that validates round readiness, drafts consensus skeletons from known protocol metadata, and sequentially invokes selected participants to append their own signoff blocks. This moves the product past round-01 without violating the Parley Deck rule that each agent owns its own artifacts.

## Proposed approach

Add a `parley consensus` command group focused on design consensus first:

- `parley consensus status [--dir DIR] IDEA`
  - Reads `00-prompt.md`, available `round-NN/<agent>.md`, and `consensus.md` if present.
  - Reports current phase, expected participants, missing round files, existing signoffs, blockers, reservations, and whether finalization is allowed.
  - Exits non-zero only for malformed protocol state or explicit `--strict` failures; normal incomplete status should be printable in CI.

- `parley consensus draft [--dir DIR] [--round N] IDEA`
  - Creates `consensus.md` only when all active participants have submitted the selected round.
  - Writes a conservative template with frontmatter, sections for agreed decisions/trade-offs/deferred items, and an empty signoff section.
  - Does not invent agreement text. It can include a checklist of round files to inspect, but the drafter must fill substantive content.
  - Refuses to overwrite an existing consensus unless a future explicit repair mode is designed.

- `parley consensus signoff [--dir DIR] [--agent ID] [--status accept|reserve|block] [--notes TEXT] IDEA`
  - Supports local/manual signoff append for a single participant.
  - Validates that the participant is in quorum, that `consensus.md` exists, and that the agent has not already signed.
  - For `block`, requires a counter-proposal.
  - This is useful for humans and for scripted agent invocations.

- `parley consensus request-signoffs [--dir DIR] [--participants IDS] IDEA`
  - Sequentially invokes configured agent CLIs and asks each participant to append only its own signoff block.
  - Uses existing runtime config, generous timeouts, isolated homes, and the same "one exact file" prompt discipline used for rounds.
  - Stops on the first failed append, duplicate signoff, malformed signoff, or `BLOCK`.
  - Records runtime events under a new run ID so `status` and `resume` can show progress.

For review consensus, reuse the same internal primitives but expose it conservatively:

- Either `parley review-consensus status|draft|signoff|request-signoffs IDEA`, or `parley consensus ... --review`.
- First implementation may support `status` and `draft` for review consensus, while deferring automated reviewer signoffs if that complicates prompts and validation.

Implementation should live mostly in `internal/protocol` plus thin `internal/app` command wiring:

- Add protocol structs for `IdeaFrontmatter`, `ConsensusState`, `Signoff`, and `ConsensusKind` (`design` or `review`).
- Add parsers for signoff blocks using deterministic Markdown scanning, not loose semantic interpretation.
- Add writers that append signoff blocks atomically and preserve existing content.
- Add validation helpers:
  - expected participants from `00-prompt.md`;
  - latest complete round detection;
  - missing participant files;
  - duplicate signoffs;
  - accepted/reserved/blocked aggregate status;
  - whether `FINAL.md` or fix-up may proceed.
- Add focused tests before TUI work. This is protocol-critical code.

The TUI can show consensus state later. The first slice can be CLI-only, with `parley status` optionally including a short consensus summary after the core commands exist.

## Concerns / open questions

- Command naming matters. `parley consensus --review` is compact but may be awkward; `parley review consensus` is clearer but deeper than the current dispatcher style.
- Automated `request-signoffs` can spend hosted model tokens and writes to a shared file. It must be sequential and explicit, probably gated by `--yes` unless a single local agent is selected.
- We need to decide whether `draft` should generate a nearly empty template or include a mechanical summary of round files. I prefer the empty template plus checklist because inaccurate summaries would be worse than manual work.
- Native GitHub PR review mirrors should probably remain a follow-up. Canonical files should come first.

## Risks

- If the CLI appends signoffs too freely, it can create invalid consensus and undermine the audit trail. Validation and no-overwrite behavior are mandatory.
- If `request-signoffs` is too broad, users may accidentally launch every configured hosted agent. It should require explicit participants or a confirmation gate.
- Markdown signoff parsing can become brittle. Keep the accepted signoff syntax narrow and well tested rather than accepting every possible Markdown variant.
- Review consensus and design consensus are similar but not identical. Sharing primitives is good; hiding their differences behind one leaky command is not.
