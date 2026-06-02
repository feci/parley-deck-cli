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
