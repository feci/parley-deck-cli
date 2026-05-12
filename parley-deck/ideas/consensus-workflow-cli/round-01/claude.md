---
agent: claude
idea: consensus-workflow-cli
round: 1
date: 2026-05-12
---

## Summary

Ship a small, deterministic `internal/consensus` package plus a `parley consensus`
command group that scaffolds, signs, validates, and finalizes the canonical
`consensus.md` artifact — without generating any of its substantive prose. The
same primitives should be parameterized by section schema so review consensus
(`ideas/<slug>/review/consensus.md`) reuses parser/validator/append/finalize on
day one without widening the slice. GitHub PR API automation stays out of this
slice; the file remains canonical and the CLI only emits hints about the
matching native review action.

## Proposed approach

### New package: `internal/consensus`

A pure, file-driven package with no network or model calls. Mirrors the shape
of `internal/protocol` and is tested with table-driven fixtures.

Public surface (rough):

- `type Schema struct { Path, Sections []SectionSpec, Participants []string }`
  — chosen by caller as either `DesignSchema(ideaDir, prompt)` or
  `ReviewSchema(ideaDir, impl, prompt)`. The two schemas differ only in the
  required section headings and the participant set (review includes the
  implementer; design uses `participants:` from `00-prompt.md`).
- `Scaffold(schema, drafter, date) error` — writes a new `consensus.md` if and
  only if it does not exist. Refuses to overwrite (matches the round-01
  protection in this very prompt). Body contains canonical frontmatter, the
  required section headings as empty stubs, and an empty `## Signoffs` section
  with the protocol's HTML reminder comment.
- `AppendSignoff(schema, block SignoffBlock) error` — atomic append-only
  edit of the `## Signoffs` section. `O_APPEND` write of a rendered block
  guarded by a per-file lock plus a re-parse pre-check that the same agent
  does not already have a block. `block.Status ∈ {accept, reservations, block}`;
  `reservations` requires `Notes`; `block` requires `CounterProposal`.
- `Parse(path) (Doc, error)` — strict-but-forgiving parser: frontmatter via the
  existing `protocol.ReadFrontmatter` helper, sections matched by exact
  `## <Heading>` and signoff blocks matched by a fixed regex on the
  `### Signoff: <agent> — <date>` header plus the `Status:` line. Returns line
  numbers so validation can pinpoint format errors.
- `Validate(schema, doc) Report` — deterministic checks only:
  - frontmatter has required keys;
  - every required section exists and is non-empty;
  - every participant has exactly one signoff block;
  - no unknown agents have signoff blocks;
  - statuses parse and required side fields are present;
  - returns a triage: `ready` (all ✅), `partial` (missing signoffs), `reserved`
    (all signed, ≥1 🟡 but 0 ❌), `blocked` (≥1 ❌), or `malformed`.
- `Finalize(schema, doc, drafter) (finalPath, error)` — succeeds only when
  `Validate(...).Triage == ready` or `reserved`. Writes a `FINAL.md` skeleton
  with canonical frontmatter and stub `## Final plan / specification` and
  `## References` sections; the CLI does not invent prose. Updates
  `00-prompt.md`'s `status:` to `final` and the consensus draft writes
  `status: consensus`. These two frontmatter mutations are surgical
  (replace the single `status:` line), not full rewrites.

### CLI surface: `parley consensus …`

Five thin verbs in `internal/app`, all built on the package above. They follow
the existing pattern of `parley status` / `parley resume`: plain terminal
first, `--json` flag for tests/CI, no TUI required in this slice.

- `parley consensus draft <idea> --by <agent>` — calls `Scaffold`, bumps the
  prompt status to `consensus`, emits `consensus.drafted` to the run's
  `events.jsonl` if a run is active, and prints the next-step hint (the
  drafter still has to fill the Agreed sections by hand or via their own
  agent runtime — the CLI does not).
- `parley consensus sign <idea> --agent <id> --status accept|reservations|block [--notes …] [--counter …]`
  — appends a single block. Refuses if a block for `<id>` already exists.
  Emits `consensus.signed`. Prints a one-line reminder under github-pr
  ("also submit your matching PR review: Approve / Request changes").
- `parley consensus status <idea> [--json]` — runs `Validate` and prints the
  triage plus a per-participant table (signed / missing / status). Exits
  non-zero only on `malformed`. Reuses the same formatting helpers as
  `parley status --idea`.
- `parley consensus finalize <idea> --by <agent>` — gated on `ready` or
  `reserved`. Writes the `FINAL.md` skeleton, flips `00-prompt.md` `status:` to
  `final`, emits `consensus.finalized`. Does not merge the design PR — the
  github-pr transport's merge step stays a human/operator action in this slice.
- `parley consensus sign --review <idea> ...` and `parley consensus status
  --review <idea>` — same verbs, review schema. One flag, same code path,
  validates that `IMPLEMENTATION.md` exists before allowing a review consensus
  to be drafted.

### Integration with existing layers

- `internal/protocol`: reuse `ReadFrontmatter` and add a small `WriteStatus`
  helper that replaces the `status:` line atomically (tempfile + rename). No
  schema rewrite.
- `internal/store`: three new event types (`consensus.drafted`,
  `consensus.signed`, `consensus.finalized`) with `idea`, `agent`, `triage`,
  and `path` fields. Append-only, observational; the file is still source of
  truth.
- `internal/runstate`: extend the per-idea detail projection so
  `parley status --idea <slug>` calls `consensus.Validate(...)` whenever
  `consensus.md` exists and surfaces the triage and per-agent signoff state
  alongside the existing artifact list. This is the only "engine" wiring this
  slice needs.
- `internal/tui`: no new panel in this slice. The existing live TUI already
  renders artifact presence; surfacing the triage in the artifact line is a
  one-string change. A dedicated consensus panel can come later if needed.
- `internal/runner` and `internal/agents`: untouched. Agents continue to write
  their own round files; the CLI does not invoke an agent to author a signoff
  block. This preserves the non-solo invariant: `parley consensus sign` is a
  structured append the agent (or operator on its behalf) drives — the CLI is
  not standing in for participants.

### Why this shape

- Deterministic before generative. Every behavior is parser/validator/append;
  the CLI can prove what it does. Prose still comes from agents.
- Sequential signoffs are not strictly ordered by the protocol — the rule is
  one-block-per-agent and append-only. We enforce exactly that and accept
  parallel arrival. "Sequential" in the prompt is satisfied by serializing
  appends, not by forcing a phase queue.
- Review consensus shares 95% of the code via a `Schema` switch. Adding it
  later would rewrite the parser; adding it now is a few extra section names
  and one participant-set rule.
- GitHub PR automation is excluded by constraint. We emit a textual reminder
  on `sign` and let humans/agents mirror to PR reviews. A follow-up idea
  `consensus-pr-mirror` can add `gh` API calls.

### Tests (deterministic, fixture-based)

- Scaffold: refuses to overwrite an existing `consensus.md`.
- AppendSignoff: rejects duplicate agent, rejects unknown agent for the
  schema's participant set, rejects malformed status, validates the
  status-specific required fields.
- AppendSignoff concurrency: two parallel appends both land cleanly without
  truncation (use the file lock; assert via a stress test).
- Parse: missing section, missing frontmatter key, mis-cased status emoji,
  duplicate signoff blocks — each returns a precise line-numbered error.
- Validate triage: `ready`, `partial`, `reserved`, `blocked`, `malformed`.
- Finalize: refuses when triage is `partial`/`blocked`/`malformed`; succeeds
  on `ready` and `reserved`; writes a skeleton `FINAL.md`; flips `status:`
  to `final`.
- CLI: `parley consensus status --json` shape stable enough for tests;
  exit codes match `runstate`'s conventions.
- Review schema: same matrix on `review/consensus.md`.

### Out of scope for this slice

- GitHub API calls (PR reviews, labels, merge). Follow-up idea.
- Auto-generated Agreed-decisions prose from round files. Round 2 question.
- TUI consensus panel beyond a one-string triage in the artifact view.
- Auto-merging the design PR on finalize.
- Bidirectional sync between PR review state and the file. The file wins by
  protocol; mirror enforcement is a separate idea.

## Concerns / open questions

- **Identity for `sign --agent`.** In github-pr the canonical reviewer
  identity is the GitHub handle; here the CLI takes `--agent <id>` on the
  caller's word. The defense is that the consensus.md change is committed
  under the agent's git identity per transport rules, so misuse is visible in
  history. Worth round-2 challenge: should the CLI cross-check `--agent`
  against `git config user.name`/host handle when the transport is github-pr?
- **Should `parley consensus draft` pre-fill the Agreed sections from a model
  summary of the round files?** Per the prompt constraint, no — deterministic
  templates first. Flagging it explicitly so round 2 can confirm.
- **`status: consensus` bump timing.** Should `00-prompt.md` flip to
  `consensus` on `draft`, or only after the drafter has filled the sections?
  Leaning toward `draft`-time because the round files are already closed by
  the time anyone calls `draft`; want to confirm.
- **Reservations triage.** `Validate` treats `reserved` (all ✅+🟡, no ❌) as
  finalize-eligible per COOPERATION.md §4 Phase 3 ("🟡 is acceptable if the
  reservation is logged as open items and no one upgrades it to ❌"). Want
  explicit confirmation that the CLI may proceed without re-prompting.
- **Append-during-edit.** What if a participant is hand-editing
  `consensus.md` in their editor when another participant runs `sign`?
  The lockfile prevents two CLI processes from racing, but not a CLI vs. a
  human text editor. Mitigation is the post-write re-parse: if the file is
  malformed afterward, the event log says who wrote last. Worth round-2
  discussion or accept-and-move.
- **Review schema participant set.** Phase 6/7 reviewers are "all
  participants except the implementer" for review rounds but "implementer
  included" for review consensus signoffs. The schema needs to express
  this difference; my plan handles it but it is the kind of fiddly rule worth
  cross-checking.

## Risks

- **Brittle Markdown parsing.** Tolerant enough to read agent-authored files,
  strict enough to detect a missing signoff. Mitigation: lock the signoff
  block shape to a regex (header line + Status line) and surface
  line-numbered errors via `Validate`. Sections matched by exact `##` headers
  per protocol.
- **Forged signoffs.** `--agent` is on the caller's honor. Mitigation: every
  state-changing call writes a `consensus.*` event to `events.jsonl` with the
  caller's process identity; the git commit carries the host handle under
  github-pr. Document the trust boundary in `parley consensus --help` and the
  github-pr conventions section.
- **Scope creep into FINAL.md authoring.** The CLI must not invent the final
  plan. `finalize` emits a skeleton only, identical in spirit to how
  `parley init` writes COOPERATION.md without authoring the project.
- **PR review drift.** Without a mirror enforcer in this slice,
  `consensus.md` and the GitHub PR's native review state can disagree. The
  file wins by protocol; users seeing two truths is a minor confusion.
  Mitigation: `sign` prints the matching native action; a follow-up idea
  closes the loop.
- **Review consensus shape divergence later.** If we wait and ship the
  design-only version first, a follow-up will rewrite the parser to accept
  different section names. Mitigation: parameterize `Schema` on day one even
  though it is one extra struct field.
- **Concurrent append correctness.** A poorly tested lock could corrupt
  `consensus.md`. Mitigation: include the stress test in the test matrix and
  use the well-trodden tempfile-rename approach for the status-line edits;
  pure-append on a locked file for signoffs.
