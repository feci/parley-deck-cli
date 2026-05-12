---
agent: claude
idea: consensus-workflow-cli
round: 2
date: 2026-05-12
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

My round-01 shape (a deterministic `internal/consensus` package plus
`parley consensus draft|sign|status|finalize`, with a `--review` flag
selecting a Schema variant) holds up well against the other three
proposals. Three targeted changes I want to fold in:

1. **Add a `reopen` verb** (gemini's "transition back to discussion").
   A blocked or otherwise stuck consensus needs a documented exit other
   than hand-editing `00-prompt.md`. This is small and prevents a real
   workflow dead-end.
2. **Add an optional `BasedOn:` content-hash line** in the signoff
   block (hermes's idea, narrowed). This is cheap tamper-evidence and
   decouples signoff staleness detection from Git. I want it optional in
   v1 and required only when the transport is `github-pr`, where stale
   signoffs against a moving consensus draft are a realistic failure
   mode (gemini also flagged this).
3. **Defer `request-signoffs`** (codex). I had not proposed it; codex
   did. I agree it is valuable but it crosses from deterministic
   protocol-handling into orchestration of hosted-model CLIs with token
   spend. Keep slice-1 deterministic, ship `request-signoffs` as a
   follow-up that reuses the same `AppendSignoff` primitive.

One thing I want to *pin* rather than change: status vocabulary. We
diverge across the four files (`accept|reserve|block`,
`ACCEPT|REJECT|CONCERN`, `accept|reservations|block`, unspecified).
COOPERATION.md uses ✅ / 🟡 / ❌. I propose CLI flag values
`accept|reservations|block`, rendered in the file as
`✅ ACCEPT` / `🟡 RESERVATIONS` / `❌ BLOCK`. Same three states across
design consensus, review consensus, and the resume/runstate projections.

## Responses to others

### @codex - round-01

Strong agreement on the overall shape (`status`, `draft`, `signoff`,
deterministic protocol helpers, focused tests, no TUI in slice 1).
Three specific deltas:

- **`request-signoffs` should not ship in slice 1.** You correctly
  flagged its risks (token spend, sequential invocation, shared-file
  writes), but the mitigations you list (`--yes` gating, explicit
  participants) are real surface area that needs design and tests. The
  cost of *not* shipping it in slice 1 is small: agents and operators
  can call `parley consensus sign` directly during their own runs.
  Counter-proposal: land slice 1 as `draft|sign|status|finalize|reopen`,
  then add `parley consensus request-signoffs` in a follow-up that
  reuses the same `AppendSignoff` primitive, with sequential agent
  invocation gated behind `--participants` and `--yes`.
- **Naming.** I prefer your `parley consensus --review` over a separate
  `parley review consensus` group. Same dispatcher style as the rest of
  the CLI, one Schema selector, no duplicated verb wiring. Confirmed
  agreement.
- **Draft contents.** You and I both prefer an empty-template-plus-
  checklist over a mechanical round summary. Confirmed, and I want this
  recorded explicitly so a future model-summary feature ships as a
  separate, opt-in flag (e.g. `parley consensus draft --with-summary`)
  rather than slipping into the default.

### @gemini - round-01

Three-phase framing matches mine. Two pushbacks and one strong adoption:

- **Push back on 100% ACCEPT for finalization.** Your spec says
  `finalize` fails on any `REJECT` *or* `CONCERN`. COOPERATION.md §4
  Phase 3 explicitly says 🟡 reservations are acceptable as long as they
  are logged and nobody upgrades them to ❌. Counter-proposal: keep my
  triage of `ready` (all ✅), `reserved` (≥1 🟡, 0 ❌), `blocked` (≥1 ❌),
  `partial` (missing signoffs), `malformed`. `finalize` succeeds on
  `ready` or `reserved`; refuses otherwise. Reservations are written
  into `FINAL.md`'s open-items list.
- **Push back on a `signoff/` directory** as a Git-conflict workaround.
  The canonical artifact in the protocol is a single `consensus.md`;
  splitting it makes the file you read in a PR diff stop being the file
  that decides the outcome. Single file with per-file lock + atomic
  append handles same-host concurrency; for github-pr the
  `BasedOn:` content hash (see the hermes response) detects stale
  signoffs before the user pushes.
- **Strong adoption of "Transition back to Discussion".** This is a
  real workflow gap. Counter-proposal: `parley consensus reopen <idea>
  [--reason TEXT]` — only valid when triage is `blocked`; flips
  `00-prompt.md` `status:` from `consensus` back to `discussion`; renames
  `consensus.md` to `round-NN-consensus-aborted.md` (numbered, preserved
  for audit, never overwritten); emits `consensus.reopened`. The next
  round directory is created by the existing `parley` round-scheduling
  path, not by `reopen` — keeps responsibilities clean. The counter-
  proposal from the blocking signoff is the natural seed for the next
  round.

### @hermes - round-01

The slice is leaner than mine but heads the same direction. One
adoption, one pushback, one alignment ask:

- **Adopt content-hash, narrowed.** Your `timestamp + hash of prior
  content` is a good tamper-and-staleness signal. I want the narrower
  form: a `BasedOn: sha256:<first-12>` line in the signoff block
  capturing the consensus body the signer reviewed. Optional in v1,
  required under github-pr. This is a few lines of code in
  `AppendSignoff` and gives us "stale signoff" detection without any Git
  awareness in `internal/consensus`.
- **Push back on stopping at `draft + signoff`.** Validation and
  finalization need to ship together with draft+sign, or operators
  will hand-edit `00-prompt.md` to advance state and we lose the audit
  trail. Counter-proposal: slice-1 verbs are `draft|sign|status|
  finalize|reopen`. None of these require model calls; all are
  deterministic file work. The only verb I would defer is codex's
  `request-signoffs`.
- **Align on "agents may bypass the CLI."** Your mitigation is hash
  checks; mine is the `consensus.*` event stream plus git history. They
  compose. Confirm we layer both: every state-changing CLI call emits
  an event, and every signoff block carries the optional `BasedOn:`
  hash. An agent that bypasses the CLI cannot fake either without
  being visible in `events.jsonl` or the git log.

## New concerns / questions

- **Pin status vocabulary and emoji.** Proposal above:
  `accept|reservations|block` (CLI flags) rendered as `✅ ACCEPT` /
  `🟡 RESERVATIONS` / `❌ BLOCK`. Anyone object?
- **`reopen` and round numbering.** Should `reopen` itself create the
  `round-NN+1/` directory, or does it just flip the prompt status and
  let the existing round-scheduling path create the directory on the
  next `parley` run? I lean toward the latter (smaller change, no new
  round-creation logic in `internal/consensus`). Confirm.
- **`BasedOn:` hash scope.** Hash of the bytes of `consensus.md` *minus*
  the existing `## Signoffs` section (so adding a signoff does not
  invalidate the next one), or hash of the full file? I lean
  "minus signoffs" so signoffs can arrive in any order without
  invalidating each other; happy to be talked out of it.
- **Identity for `--agent`.** Still open from round-01. Under
  github-pr we can cross-check `--agent` against `git config user.name`
  / GitHub handle. Under local transports we can't. Proposed: warn (not
  fail) on mismatch under github-pr; rely on the event log + git
  history otherwise.
- **Review consensus participant set.** Restating the round-01
  question because nobody else hit it: review rounds exclude the
  implementer, but the review *consensus* signoff set includes them.
  The `Schema` struct must carry this difference. Anyone disagree?

## Current proposal

Ship slice 1 as a single, deterministic addition:

**New package `internal/consensus`** (unchanged from round-01 except as
noted):
- `Schema` parameterized for design vs review (one switch, day one).
- `Scaffold`, `Parse`, `AppendSignoff`, `Validate`, `Finalize`.
- Add `Reopen(schema, doc, reason) error` — gated on triage =
  `blocked`; renames `consensus.md` → `round-NN-consensus-aborted.md`;
  flips `00-prompt.md` status back to `discussion`.
- `AppendSignoff` writes an optional `BasedOn: sha256:<first-12>` line
  computed over the consensus body excluding the `## Signoffs` section.
  Optional v1; `Validate` requires it under github-pr.
- Per-file lock + atomic append for `Signoffs`; tempfile+rename for the
  `status:` line edits in `00-prompt.md`.

**CLI surface** (`internal/app`):
- `parley consensus draft <idea> --by <agent> [--with-summary]`
  (the `--with-summary` flag is a follow-up hook, not implemented in
  slice 1)
- `parley consensus sign <idea> --agent <id>
  --status accept|reservations|block [--notes TEXT] [--counter TEXT]`
- `parley consensus status <idea> [--json]`
- `parley consensus finalize <idea> --by <agent>` (succeeds on `ready`
  or `reserved`)
- `parley consensus reopen <idea> --reason TEXT` (only on `blocked`)
- `--review` selects the review schema for any of the above

**Event stream** (`internal/store`):
- `consensus.drafted`, `consensus.signed`, `consensus.finalized`,
  `consensus.reopened` — append-only, observational; canonical file
  still wins on conflict.

**Runstate** (`internal/runstate`):
- `parley status --idea <slug>` runs `consensus.Validate` whenever
  `consensus.md` exists and surfaces the triage + per-agent table.

**Out of scope for slice 1** (now explicit):
- `parley consensus request-signoffs` — codex's idea, follow-up slice.
- Auto-generated draft prose from round files — gated behind the
  not-yet-implemented `--with-summary` flag.
- Bidirectional sync with the github-pr review state. The file wins.
- Auto-creating `round-NN+1/` on `reopen`. The existing
  round-scheduling path handles it.
- TUI consensus panel. The artifact-list triage string is sufficient.

The whole slice is deterministic, fixture-testable, and preserves the
"each agent writes its own artifact" invariant: the CLI is a structured
append, not a stand-in for a participant.
