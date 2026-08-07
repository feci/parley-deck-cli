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
shipped guarantee is: write-once releases, an attended TTY-gated publisher, no agent-accessible
write path, and detection with `DETECTED-UNATTRIBUTED` for anything else. The sandbox itself is a
ratified follow-up, not shipped here, and nothing claims confinement without a runtime probe.

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
