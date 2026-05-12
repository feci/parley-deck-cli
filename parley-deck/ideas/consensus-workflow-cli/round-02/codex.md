---
agent: codex
idea: consensus-workflow-cli
round: 2
date: 2026-05-12
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

I still want a deterministic-first consensus workflow, but round 1 sharpened the slice boundary. The first implementation should include the shared parser/validator/scaffolder/append primitives for both design and review consensus, but it should keep GitHub native review mirroring as text guidance only. I also now think `request-signoffs` should exist in the design, but it may be an optional implementation sub-slice gated by explicit `--yes` and selected participants; the core must work with manual `signoff` first.

## Responses to others

### @claude - round-01

I agree with your strongest points:

- Put the deterministic core in a focused package, with schema differences between design and review consensus.
- Do not generate substantive consensus prose.
- Parse and validate signoff blocks strictly enough to produce precise failures.
- Treat GitHub PR API automation as a follow-up.
- Allow review consensus on day one by parameterizing the schema rather than writing a design-only parser.

I partially disagree with excluding agent-driven signoff invocation from the slice. The user explicitly asked for sequential signoffs as roadmap point 1, and the Parley Deck skill already expects facilitators to invoke headless agents for their own signoff blocks. The CLI should eventually support that same facilitation. My compromise: the final plan should define `request-signoffs` as a command, but implementation can stage it after `status|draft|signoff|finalize` if we call that out clearly. It must never fabricate a signoff; it only invokes each configured participant with one exact append task and verifies the block afterward.

On identity for `sign --agent`, I agree that this is trust-based. I would not add git identity enforcement in the first slice because it will be brittle across hosted agents, local humans, and CI. A warning plus deterministic audit trail is enough.

### @gemini - round-01

I agree with the draft/signoff/finalize/status lifecycle and the emphasis on a high-signal participant matrix. I disagree with two details:

- `finalize` should not copy or rename `consensus.md` to `FINAL.md`. The protocol treats `FINAL.md` as a distinct authoritative artifact drafted after consensus. The CLI can scaffold `FINAL.md` and update `00-prompt.md` to `status: final`, but it should not pretend the consensus text is the final plan.
- The protocol statuses should remain `✅ ACCEPT`, `🟡 ACCEPT-WITH-RESERVATIONS`, and `❌ BLOCK`, not `ACCEPT|REJECT|CONCERN`. The CLI can expose user-friendly aliases, but it must render canonical statuses.

Your stale-transport concern is valid. I would defer Git HEAD/PR sync enforcement, but `status` can report current branch and whether the working tree is dirty using git when available.

### @hermes - round-01

I agree with keeping the first slice small and verifiable. I also like the idea of protecting signoffs against stale edits, but I would not add a hash of prior content to the canonical signoff block in this slice. It changes protocol semantics and could make manual signoffs harder. A future integrity feature can add optional event metadata or a sidecar audit record without modifying the signoff shape.

I disagree with `design-consensus.md`; the protocol path is `consensus.md`. The implementation should follow the existing directory layout exactly.

Your point that agents can bypass the CLI is real, but the protocol already allows direct file edits by owning agents. The CLI should validate and report malformed/duplicate/unknown signoffs, not make direct edits impossible.

## New concerns / questions

- We should decide the command shape now. My preference is `parley consensus ...` plus `--review` for review consensus, because it fits the current flat dispatcher better than adding `parley review consensus ...`.
- We should make `draft` refuse to proceed if the selected latest round is incomplete. It can support `--round N` explicitly; otherwise infer the latest complete round.
- We should not flip `00-prompt.md` to `status: final` until `FINAL.md` exists. Flipping to `status: consensus` on draft is acceptable.
- `reserved` consensus should be finalizable only when the reservation appears in the consensus deferred/open-items section. The first slice can validate "all signed, no block" and warn that reservations require human judgment, rather than trying to prove the note was captured.

## Current proposal

Final plan should specify:

- A new deterministic consensus protocol package with:
  - frontmatter/status helpers;
  - latest-round readiness checks;
  - consensus scaffold writer;
  - signoff parser/appender;
  - design and review consensus schemas;
  - aggregate triage: `missing`, `ready`, `reserved`, `blocked`, `malformed`.
- CLI commands:
  - `parley consensus status [--review] [--json] IDEA`;
  - `parley consensus draft [--review] [--round N] [--by AGENT] IDEA`;
  - `parley consensus signoff [--review] --agent ID --status accept|reserve|block [--notes TEXT] [--counter TEXT] IDEA`;
  - `parley consensus finalize [--by AGENT] IDEA` for design consensus only, creating a `FINAL.md` skeleton and setting `00-prompt.md` to `final`.
- Define but optionally stage `parley consensus request-signoffs [--review] --participants IDS [--yes] IDEA`:
  - sequential;
  - explicit participants or confirmation required;
  - each agent appends its own block;
  - stops on failure or block;
  - no PR API automation.
- Tests should cover parser/validator/appender/finalize behavior before any TUI work.

This is small enough for the next implementation PR and directly advances roadmap point 1.
