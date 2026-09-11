## 2026-09-10 — §9 item 1: a `refused` launch is a stop, not a fallback (UNRELEASED)
Idea: ideas/meta-protocol-change-evidence-first-efficiency/
Drafted by: claude-1
Summary: §9 item 1 previously said "Without an attestation, or on `refused`, read all of
`parley-deck/COOPERATION.md` and record `context_mode=full-fallback`", which read as permission to
convert a refusal into an ordinary fallback. The item now separates the two outcomes. `full-fallback`
stays exactly as before: a valid, visible result that reads the live authority in full, records its
reason, and proceeds. `refused` (unprovable authority, a detected secret) is a **stop** — never emit
the refused content, never substitute another authority for it (a bundled snapshot, a cached or stale
copy, a hand-assembled excerpt), never continue that launch on unattested text; resolve at the
renderer and re-render, or report the blocker. A protocol task launch carrying no attestation is
unresolved the same way and obtains one from the renderer before the task starts; the read-the-live-
source fallback applies only where no renderer is reachable. This states FINAL D4's existing rules
("Missing/unprovable authority blocks rather than substituting a bundled snapshot"; "Detected secrets
refuse external context emission") at the point an agent acts on them.
Unchanged: no authority is broadened, no new context mode exists, `full` remains the default, an
optimized `packet` remains the ratified trial's explicit input, and no
`meta/packet-applicability.yaml` classification changed (the edit is prose inside the existing
`### 9.0 Pre-idea readiness check …` block; no heading was touched).
Mirrored identically into `internal/protocol/defaults/COOPERATION.md` and the parley-deck-skill
reference copy `skills/parley-deck/references/COOPERATION.md`; `skills/parley-deck/SKILL.md`
standing instructions carry the matching stop/no-substitution and missing-attestation wording.
The same commit mirrors the LE-7 close-integrity change below into the skill reference copy,
superseding that entry's "not yet mirrored into the parley-deck-skill reference copy" status note.

**Status: UNRELEASED.** Deck source, embedded default and skill source only — not published to a
global core and not in a package release. Checks actually performed: PRIMARY source reads with
file:line locators. No build, no `go test`, no drift-guard run and no skill add-on hash-manifest
regeneration — no shell was available in this session; the facilitator regenerates the manifest and
runs the tests. No signoff, no Phase-6 verdict and no acceptance is claimed.

## 2026-09-10 — §4 LE-7: a goal-done check can withhold a close, never establish one (UNRELEASED)
Idea: ideas/meta-protocol-change-evidence-first-efficiency/
Drafted by: claude-1
Summary: The §4.0.1 LE-7/LE-11 line and the Phase-8 "Close-decision integrity" paragraph no longer
say the goal-done check is "fail-open on its own error (a broken or inconclusive checker never
blocks a review-clean idea)". A checker that is missing, is the implementer, or cannot be resolved
and launched; an execution that fails or exits non-zero; and a verdict that is inconclusive or only
a pass-with-reservations each now leave completion **unverified** and escalate for a human decision.
The paragraph also states what the check is NOT: a textual verdict is defense in depth on top of the
review consensus and never substitutes for the current-tree independent criterion evidence a close
already requires (this idea's FINAL D3 — a self verdict, a stale code tree, a skipped/no-execution
report, a missing criterion or a partial original scope cannot close an implementation).
Unchanged: the check still fires only under `auto_implement` or `strict_gate`; the
`ACCEPT-WITH-RESERVATIONS` and fewer-than-two-independent-reviewer refusals stand as written; a
design-only idea keeps the lighter close. No §4.0 track cell, quorum rule, signoff rule or
`meta/packet-applicability.yaml` classification changed. Mirrored byte-identically into
`internal/protocol/defaults/COOPERATION.md` (the drift guard requires both copies to match).

**Why:** the in-tree driver already fails closed on exactly these conditions
(`internal/app/driver_impl.go` `GoalCheck`, whose `false` returns `internal/driver/impl.go` turns
into `ActionEscalated`), so the protocol text in force asserted the opposite of the code. This entry
carries the §7 text change that code needed; it does not certify the code. Known remaining gaps,
owned elsewhere: a stale in-tree comment at `internal/driver/impl.go` still says "fail-open inside
GoalCheck", and the verdict-aggregation correction is in progress on the codex-1 slice.

**Ratification:** the design is accepted in this idea's `consensus.md` (codex-1, claude-1, hermes-1,
kimi-1) and D3 states the rule; this is the Phase 5 protocol-text change on the claude-1 slice,
pending independent review. No signoff, no overall acceptance and no Phase-6 verdict is claimed.

**Status: UNRELEASED.** Deck source plus embedded default only — not published to a global core, not
in a package release, and not yet mirrored into the parley-deck-skill reference copy (separate
worktree, outside this slice). Checks actually performed: PRIMARY source reads with file:line
locators and a line-by-line comparison confirming both protocol copies now carry identical text. No
build, no `go test`, and no drift-guard run — no shell was available in this session.

## 2026-09-05 — §9 item 1: launch context comes from the shared packet renderer with attestation
Idea: ideas/meta-protocol-change-evidence-first-efficiency/
Drafted by: claude-1
Summary: §9 item 1 now requires an official launch to receive its protocol context from the shared
renderer (`parley protocol packet`) with an attestation (`context_mode`, `source_sha256`,
`packet_sha256`, `fallback_reason`) rendered from the live resolved authority, never a bundled
snapshot; without an attestation an agent reads the full file and records `full-fallback` with the
reason. Full context stays the default and the ratified packet experiment (phase-packet FINAL §3:
phases 1 and 6, six matched AB/BA pairs each, three canaries plus a full control, ship at R ≤ 0.50
in both phases) is unchanged; an optimized packet is that trial's explicit experimental input, not
an enabled release. `meta/packet-applicability.yaml` is the ratified applicability map and is
protocol: a classification change is a §7 change. Mirrored into the embedded default and the
skill reference copy. Runner/handoff prompt wiring is integration-owned and is NOT claimed here.

**Ratification:** design accepted in the idea's consensus.md (codex-1, claude-1, hermes-1, kimi-1);
this is the Phase 5 source change on the claude-1 slice, pending independent review.

## 2026-08-07 — §7 blast radius: a core change is not a deck change
Idea: ideas/meta-protocol-change-global-core-protocol/
Drafted by: claude-1
Summary: The protocol moves to a single global core under `~/.parley/protocol/core/<version>/`, of
which each deck's COOPERATION.md is a generated view; §7 now distinguishes a CORE change (meta idea
plus explicit user ratification, user-only) from a DECK overlay change (a normal idea in that deck).

**Ratification:** track `deliberation`, 2 rounds, 3 consensus revisions, accepted by claude-1,
codex-1, hermes-1, kimi-1. opencode-1 was invoked four times, produced no artifact, and is recorded
absent rather than agreeing
(`inbox/claude-1-to-all_meta-protocol-change-global-core-protocol_opencode-timeout.md`).

**Why:** measured across 36 decks before a one-off sync — eight different `deckVersion` values, §15
present in 5 of 36, the §2 roster-authority change in 1 of 36. The per-deck copy-as-store model had
already failed, and the hand-written sync that repaired it was not a mechanism. Only ONE genuine
local protocol section existed in the whole fleet, and it was governance about how the protocol is
synced — content that belongs in the core, not in a deck.

**Enforcement, stated honestly.** Prevention of an agent writing the core is real for
parley-launched participants under an OS sandbox (verified: a macOS seatbelt profile denies the
write, the denial is inherited by children, and `rm` is denied too) — but a profile built from an
UNRESOLVED path silently denies nothing, and the facilitator is not launched by parley. So the
shipped guarantee is: write-once releases (per release directory, `O_EXCL|O_NOFOLLOW`, symlinked
store components refused on read and write), an attended publisher that refuses without a
controlling terminal, and no agent-accessible write path. **`DETECTED-UNATTRIBUTED` and per-idea
pinning are ratified but NOT implemented** — they are ranks 2 and 4 — and the attended refusal
stops an ordinary agent run but not one that allocates a pty. The sandbox is the ratified DF-1
follow-up.

## 2026-07-04 - Add §13.5 Playbooks (parley learn distillation)

Idea: ideas/parley-learn-playbooks/ (parley-learn-playbooks)
Drafted by: claude-1
Summary: Additive. Extends §13 with §13.5 Playbooks - an advisory, non-canonical retro
output. `parley learn <closed-idea-slug>` distills a COMPLETED idea into a reusable
parley-deck/playbooks/<topic>.md (proven shape: track, roster, checklist, gotchas,
verification), idea-specific specifics stripped. Playbooks sit beside consults in the
advisory tier - never quorum, never overriding protocol; referencing one in Phase 0 is
optional context. `parley learn` is a read-only tooling command that writes exactly one
new fail-closed playbook file (Lstat symlink guard), NOT a Parley round. Unanimous
deliberation-track consensus (claude-1, codex-1, hermes-1, antigravity-1 - accept x4).
Mirrored into the embedded default + skill fallback.

## 2026-07-04 — Completion contract: list-form `checks:` + driver-written evidence

Idea: ideas/completion-contracts-evidence-ledger/ (completion-contracts-evidence-ledger)
Drafted by: claude-1
Summary: Additive, backward-compatible. `checks:` in 00-prompt.md now accepts either a
scalar command (unchanged) or an optional named list of {name, command} criteria; the
list form activates a completion contract. The driver runs each criterion, writes a
per-criterion result table (exit, duration, secret-scrubbed truncated output) into the
existing `## Validation evidence` section of IMPLEMENTATION.md each cycle (overwrite;
git history preserves prior cycles), and vetoes `status: complete` while any criterion
fails at the current HEAD — the same fail-closed shape as strict_gate, independent of it.
No new `done_when:` key and no separate evidence artifact (rejected in round-02). Scalar
or absent `checks:` is byte-for-byte today's behavior. Unanimous deliberation-track
consensus (claude-1, codex-1, hermes-1, antigravity-1 — ✅ ×4). Protocol text: LE-4 +
Phase-5 template + Phase-8, scoped to the list shape; mirrored into the embedded default
and the skill fallback snapshot.

## 2026-07-03 — Progressive-disclosure layout (relocate §9 after §10)

Idea: ideas/protocol-restructure-appendices/ (protocol-restructure-appendices)
Drafted by: claude-1
Summary: Pure content-preserving reorder — §9 (session-start checklist) relocated to sit after
§10 (TL;DR) so the document reads core-first (§0–§8, §10) then reference-last (§9, §11, §12, §13,
§14, Appendix A). Keep-numbers-relocate: every section keeps its number so all `§N` cross-refs
resolve; no rule text added/removed/changed (sorted-line diff empty); no `## Appendices` banner.
Unanimous multi-agent review (✅ ×3). Not a rule change — a layout move (still logged here for the
audit trail). `core ≤200 lines` compression + §4 phase-split are deferred to separate ideas.

## 2026-07-03 — Add §4.0 conditional-rigor tracks + developer Quickstart

Idea: ideas/meta-protocol-change-devx-speed/ (meta-protocol-change-devx-speed)
Drafted by: claude-1
Summary: Additive DevX + speed change. Added a `track: fast | standard | deliberation` field
(default `standard`) with an objective, fail-safe classifier (§4.0) that scales ceremony to
risk: `fast` = one model-diverse reviewer, collapsed consensus/FINAL, cross-review skipped,
≤1 fix-up, ~5-min timeouts; `standard` = 2 reviewers, cross-review capped at 2, ~15-min;
`deliberation` = the unchanged full lifecycle, forced by protocol/security/data-migration/
irreversible/strict_gate/auto_implement/pipeline/API-break triggers. §4.0's per-track table is
the single authoritative gate overriding the full-lifecycle defaults in §4/§5/§9.0/§11; all
MUST-stay invariants (non-solo, refutation-default, §14 human brake, audit trail, English-only,
no-secrets, round-1 independence) hold on every track. Also added a top-of-doc Quickstart, a
role table, a core-vs-appendix reading guide, an "off-ramp" (trivial reversible work needs no
Parley), and a consolidated LE-N glossary (§4.0.1). Unanimous design signoff (claude-1, codex-1,
hermes-1, antigravity-1 — ✅ ×4); implemented on the deliberation track with a full Phase-6
review round. Deferred to named follow-ups: physical appendix relocation/renumber
(`protocol-restructure-appendices`) and CLI/driver enforcement of tracks (`track-aware-driver`).

## 2026-06-02 — Add §12 Pipeline blocks & action stages

Idea: ideas/2026-06-02T12-07-14-meta-protocol-ch/ (meta-protocol-change-end-to-end-pipeline)
Drafted by: claude
Summary: Added an additive, opt-in §12 that composes the unchanged Phase 0–8 engine into typed pipeline blocks (deliberation/implementation/action/watcher) driven by a `pipelines/<slug>/pipeline.yaml` manifest, turning a single idea into an automatic idea→business-spec→technical-spec→impl-design→implementation→deployment→operations→monitoring pipeline. Agents stay markdown-only; a driver performs side effects behind a provider-agnostic interface, supervised-first gates (production mutations non-bypassable), durable cursor + per-effect idempotent ledger with reconcile-on-resume, capability-dispatch-halts-not-degrades, linear v1. Unanimous signoff (codex, claude, hermes); agy excluded (headless print-mode produced no artifact). Implementation is staged; nothing in Sections 0–11 changed.

## 2026-05-27 — Replace Gemini defaults with Antigravity CLI

Idea: ideas/antigravity-agent-migration/
Drafted by: codex
Summary: Added `agy` as the active Antigravity CLI participant, moved `gemini` to inactive legacy status for historical compatibility, and updated shared runtime defaults so new workflows prefer Antigravity while retaining explicit Gemini overrides.

## 2026-05-10 — Switch transport to GitHub PR

Idea: ideas/meta-protocol-change-github-pr-transport/
Drafted by: codex
Summary: The user created `https://github.com/feci/parley-deck-cli` and requested GitHub usage. Future Parley Deck coordination should use `github-pr` transport while keeping `parley-deck/` files canonical.

## 2026-05-14 — Adopt lightweight team coordination guidance

Idea: ideas/meta-protocol-change-agent-teams-patterns/
Drafted by: codex
Summary: Added advisory per-idea roles/lenses, internal-helper accountability, participant sizing guidance, Phase 5 plan-gate guidance, and inbox mirroring rules inspired by agent-team workflows while preserving Parley Deck's vendor-neutral canonical artifact model. See `ideas/meta-protocol-change-agent-teams-patterns/`.

## 2026-05-15 — Clarify helper identity boundaries

Idea: user-follow-up to `ideas/meta-protocol-change-agent-teams-patterns/`
Drafted by: codex
Summary: Clarified that participant-spawned helpers may contribute only through the owning participant and must not create canonical round, review, consensus, or signoff files under a separate helper identity unless that identity is explicitly listed in `participants:`.

## 2026-05-25 — Concrete roster and local headless config note

Idea: ideas/meta-protocol-change-roster-headless-config/
Drafted by: codex
Summary: Replaced placeholder roster rows with `codex`, `claude`, `gemini`, and `hermes`; marked host handles as not mapped; and documented that `parley-deck/meta/headless-agents.local.json` is optional, gitignored, machine-local launch configuration rather than canonical project state.

## 2026-06-12 — review-gate honesty (idea: meta-protocol-change-review-gate-honesty)

- Phase 6 gains **"Review briefs and dispositions"**: briefs must never suppress
  findings; known-finding dispositions travel openly with a standard shape and
  the reviewer states concurrence per disposition. Disputed findings close only
  via reviewer withdrawal, normal review consensus, or a verbatim-quoted
  operator ruling.
- Phase 8 gains the opt-in **"Strict review gate"** (`strict_gate: true` in
  00-prompt.md frontmatter): closing requires a fresh full-scope review pass
  with zero findings of any severity; fix-verification passes converge but
  never close; findings must be objective and code-grounded; mutability only
  via consensus or recorded operator direction.
- Phase 8 gains **"Stopping judgment"**: trajectory over pass counters
  (converging / churning / blocked); MaxFixupCycles is an escalation threshold,
  never a close criterion.
- §8 gains **"Consults"**: parley-deck/consults/ artifacts are advisory and
  non-canonical — never quorum evidence.
- Mirrored into the embedded default protocol
  (internal/protocol/defaults/COOPERATION.md). The external parley-deck-skill
  bundled snapshot still needs a sync — flagged via inbox.

## 2026-08-06 — §2 roster authority moves to `parley-deck/agents.toml`
Idea: ideas/roster-operations-standard/
Drafted by: claude-1
Summary: §2's roster table stops being the hand-edited membership store and becomes a generated,
non-authoritative view; `[roster.<id>]` blocks in `parley-deck/agents.toml` own membership and
per-agent settings, with retired members marked rather than deleted.

**Ratification:** track `deliberation`, 2 rounds, 3 signoff revisions, accepted by claude-1,
codex-1, hermes-1, kimi-1.

**Change:** §2's roster table stops being the hand-edited membership store and becomes a
generated, non-authoritative view. `[roster.<id>]` blocks in `parley-deck/agents.toml` carry
membership plus adapter/model/effort/speed (runtime-semantic) and workspace_dir/role/host_handle
(render-only). Retired agents are marked `active = false`, never deleted. Decks with only a legacy
table keep working and report `legacy-roster`.

**Why:** measured across 40 decks — nine distinct rosters, 17 with no roster at all, 17 still
naming an agent retired months earlier. The store the protocol told humans to maintain by hand was
the store that drifted.

**Venue deviation, recorded.** §7 requires a protocol change to run as its own
`meta-protocol-change-*` idea. It ran inside `roster-operations-standard` on the **user's explicit
one-off authorization**. This is NOT a general exception to §7 — §6 rule 3's direct-user-instruction
exception is scoped to editing another agent's file — and it sets **no precedent**: the next
protocol change needs its own meta idea unless the user again directs otherwise. Full participant
ratification was still required and obtained.
