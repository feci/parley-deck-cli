---
idea: build-companion-skills
status: implemented
implementer: claude-1
started: 2026-06-22
completed: 2026-06-22
branch: parley-deck-skill#build-companion-skills
head-commit: pending
design-pr: n/a
implementation-pr: pending
---

## Summary of work

Built both companion skills as opt-in addons inside the `parley-deck-skill` repo,
per `addon-skills-research/FINAL.md`. Authored via internal Phase-5 helpers over
disjoint files (claude-1 owns the integrated result); verified independently.

Files (in `parley-deck-skill/`):
- `addons/parley-worktrees/SKILL.md` — FINAL §A: sibling `../worktrees/` layout
  (never `.git/`), lifecycle commands, decision rules, `IMPLEMENTATION.md` lock
  manifest, file-set disjointness, runtime detect-and-adopt, pitfalls, thin core seam.
- `addons/parley-tracker/SKILL.md` — FINAL §B: tracker=mirror, one-file three-audience
  `[B]/[T]/[A]` + "At a glance" + YAML schema, hybrid Gherkin/measurable-`Verify:` AC,
  tool-enforced gap-scan, DoD=AC-ids, INVEST + vertical-slice, FINAL→tracker mapping,
  opt-in connector.
- `addons/parley-tracker/templates/{epic,story,subtask}.md` — canonical skeletons.
- `addons/parley-tracker/bin/validate.js` (+ `validate.test.js`) — dependency-free
  Node readiness/lint enforcing the §B.4 gap-scan; CLI + requireable module.
- `lib/installer.js` — `discoverAddons()`, install-all-by-default, `--no-addons`,
  `--only <name>[,<name>]`, transparent per-skill output, addons wired through
  doctor/status/paths/uninstall. `test/installer.test.js` +10 tests. `package.json`
  ships `addons/` (files + pkg.assets).

## Implementation plan / checklist
- [x] parley-worktrees SKILL.md (FINAL §A)
- [x] parley-tracker SKILL.md + 3 templates (FINAL §B)
- [x] validate tool + tests
- [x] installer --no-addons/--only + install-all default + transparent output
- [x] Checks run: `npm test` → 48 pass / 0 fail (22 pre-existing + 10 new installer + 16 validate)
- [x] Genericity: repo grep (tomasfecko, /Users/, model ids, agent ids) → no leaks
- [x] CLI smoke: `--help` shows flags; dry-run shows install-all; `--no-addons` core-only;
      `--only parley-tracker` core+tracker; `validate story.md` → PASS
- [x] Phase 6 review (codex-1 installer/validate+genericity, hermes-1 templates, antigravity-1 standards)
- [x] Phase 8 fix-up cycle 1 (all agreed fixes F1–F16); re-verified
- [ ] Release skill (version bump + npm/brew/winget) after review consensus

## Fix-up cycle 1
status: complete
completed: 2026-06-22
head-commit: pending

### Fixes applied (from review/consensus.md)
- parley-tracker: epic template now passes its own validate (F1); validate rejects
  `<...>` placeholders in required slots (F2); shipped `bin/claim.js` runs the gap-scan
  and gates `status: in-progress` (F3); validate now enforces non-empty At-a-glance /
  [B] / [T] (F4), a happy-path AC + non-empty `Verify:` (F5), and the full canonical
  schema (F6); tests now exercise all 3 shipped templates + placeholder-fail (F7);
  DoD/Verification checklist added to templates (F8); [NFR]/behavioural split (F9);
  parent-resolution documented for `--strict --dir` (F10); story/subtask terminology
  fixed (F11); story `[B]` AC + 2-line At-a-glance modeled (F12).
- installer: addon selection persisted in the core marker; doctor/status/paths/uninstall
  derive the expected addon set from the marker so `--no-addons`/legacy core-only installs
  are no longer false-negatives (F13).
- parley-worktrees: integration/feature worktree provisioning is existing-or-create
  (no `-b` abort on a pre-existing branch) (F14); `git branch -d`→`-D` guidance (F15);
  file-set intersection helper snippet (F16).

### Re-verification (facilitator, final on-disk state)
- `npm test` → **77 pass / 0 fail** (the fix-up reviewers' transient "4 fail" was a
  parallel-race stale read during concurrent edits; final state is green).
- `validate` epic/story/subtask → all exit 0; required-slot `<...>` placeholder → fails
  with a field-level message; `claim.js` → writes status/assignee on a passing ticket,
  refuses a failing one.
- `doctor` after a real `--target generic --no-addons` install → `ok: true` (F13 fixed).
- Genericity sweep over `addons/` + `lib/` → no leaks.

### Deviations
- Shipped templates are FILLED, self-passing exemplars using UPPERCASE tokens
  (SLUG/AGENT-ID/REVISION) instead of `<...>` (since validate now flags `<...>`).
- claim.js is minimal (validate-gate only) per F3 scope.
- `package.json` `feci` namespace left as-is (pre-existing distribution identity;
  FORKING.md follow-up, per the genericity audit) — dismissed as a must-fix.

## Deviations from FINAL.md
- Skill names: FINAL used design names `worktree-sessions`/`ticket-tracker`; per user
  decision the shipped names are **`parley-worktrees`/`parley-tracker`** (family brand).
- validate.js is dual-purpose (CLI + requireable module) — additive, not in FINAL.
- `--pull` conflict policy documented as "surface as a conflict" (FINAL left it open).

## Notes for reviewers
Doc + installer-code change; no protocol/COOPERATION.md edits (separate idea). Verify:
(a) genericity (no vendor/model/personal leaks); (b) installer correctness + tests +
that addons are inert/opt-out works; (c) tracker templates pass `validate` and read
well for all three audiences; (d) standards conformance (INVEST/Gherkin/DoD).
Skill files are under `parley-deck-skill/addons/`.
