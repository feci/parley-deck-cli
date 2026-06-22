---
idea: build-companion-skills
author: user
created: 2026-06-22
participants: [claude-1, codex-1, hermes-1, antigravity-1]
roles:
  claude-1: implementer (Phase 5) + integration
  codex-1: reviewer — installer wiring & validate-tool correctness, genericity
  hermes-1: reviewer — ticket templates & three-audience readability
  antigravity-1: reviewer — standards conformance (INVEST/Gherkin/DoD) & genericity
status: round-01
---

## Problem / idea

Build the two vendor-neutral companion skills designed in
`parley-deck/ideas/addon-skills-research/FINAL.md` (the binding spec), as **opt-in
addons inside the `parley-deck-skill` repo**, full scope:

1. **`parley-worktrees`** — `SKILL.md` from FINAL §A (run parallel Parley sessions /
   Phase-5 implementers over one repo via git worktrees; sibling `../worktrees/`;
   branch discipline; `IMPLEMENTATION.md` lock manifest; integration-branch merge;
   cleanup; runtime detect-and-adopt; pitfalls; thin core seam noted, not edited here).
2. **`parley-tracker`** — `SKILL.md` from FINAL §B + `templates/{epic,story,subtask}.md`
   (the full instances from `addon-skills-research/round-01/{codex-1,hermes-1}.md`) +
   a `validate` readiness/lint tool. Tracker = mirror; one file, three audiences
   `[B]/[T]/[A]` + "At a glance" + YAML metadata; hybrid Gherkin/measurable-`Verify:`
   AC; tool-enforced no-assumption gap-scan; DoD = AC-id checklist; generic field
   schema mapped to trackers at the edge; opt-in connector only.
3. **Installer** (`lib/installer.js`): bundle + install both addons **by default**
   (install-all), with **`--no-addons`** (core only) and **`--only <skill>`** opt-outs,
   transparent install output listing what was installed; independent per-skill
   versioning; keep `parley-deck` core advertised identity.
4. Packaging metadata (manifest / plugin / gemini-extension / README) updated to
   describe the addons; tests for the new installer flags + `validate`.

Phase 5 implementer = claude-1 (FINAL drafter). Phases 6-8: codex-1/hermes-1/
antigravity-1 review against FINAL.md + this prompt.

## Constraints

- English only. Strictly **vendor/tracker/runtime-agnostic** (no Jira-only, no
  personal models/agents/paths — we just audited the core for this; keep the bar).
- Follow `addon-skills-research/FINAL.md` exactly; record deviations in IMPLEMENTATION.md.
- Core stays thin: NO protocol/COOPERATION.md edits here (the two core seams are a
  separate meta-protocol-change idea).
- Install-all-by-default but addons must be inert until triggered; `--no-addons` works.

## Non-goals

- No protocol/COOPERATION.md edits (separate idea).
- No live tracker API/auth in core (connector is a later opt-in add-on).
- Not changing the parley-deck core SKILL.md beyond what the addon wiring requires.
