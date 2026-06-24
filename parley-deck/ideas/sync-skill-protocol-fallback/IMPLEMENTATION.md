---
idea: sync-skill-protocol-fallback
status: complete
implementer: claude-1
started: 2026-06-24
completed: 2026-06-24
target-repo: parley-deck-skill
review-round: 1
---

> **Review complete (Phase 6/7):** round-01 unanimously clean — codex-1, antigravity-1,
> hermes-1 all 0 findings (the RELEASING.md "second-model review" gate). See
> `review/consensus.md`. Proceeding to publish (npm + tag v1.4.1).

## Summary of work

Synced `parley-deck-skill/references/COOPERATION.md` to the canonical protocol per `FINAL.md`:
a verbatim copy of `parley-deck-cli/internal/protocol/defaults/COOPERATION.md` with only the
two header lines (`Transport`, `Created`) kept in the skill's vendor-neutral form. Bumped the
skill to 1.4.1. Implementation lives in the **parley-deck-skill** repo; this idea's artifacts
live in the CLI deck.

## What changed (parley-deck-skill)

- `references/COOPERATION.md` — replaced (759 → 1037 lines). Now byte-identical to the CLI
  embedded default from line 7 onward; lines 5–6 are the neutral `<transport-choice>` /
  `<YYYY-MM-DD>` header. Imports §0 bootstrap, the full §4 Phase-6/8 + loop-budget +
  close-integrity rework, §9.0 readiness, §12.11 candidate-remediation, §13 retrospective
  optimization, and §14 the automated-outer-loop human brake.
- `package.json` + `package-lock.json` (root + packages[""]) → `1.4.1` (via `npm version`).
- `references/compatibility.json` `skillVersion` → `1.4.1`.

## Verification (all gates green)

- `diff defaults/COOPERATION.md references/COOPERATION.md` → ONLY the 2-line Transport/Created
  header hunk.
- `diff <(tail +7 default) <(tail +7 ref)` → empty (body identity).
- Anchors: `^## 13.`=1, `^## 14.`=1; `Refutation attempts`=2, `strict_gate`=7, `status:
  candidate`=4 present.
- Leakage: `grep -i feci|claude-1|codex-1|hermes-1|antigravity-1` → 0.
- `npm test` → 77 pass / 0 fail. `npm pack --dry-run` → version 1.4.1, ships the 79.1kB
  references/COOPERATION.md, 22 files. `install --target all --dry-run` ok; `doctor --json`
  → ok: true.

## Deviations from FINAL.md

- None. Header neutralized to exactly the two specified lines; body verbatim (body `parley`
  literals kept, per consensus).

## Notes for reviewers

- **Faithfulness:** confirm the body is byte-identical to the CLI embedded default
  (`diff <(tail -n +7 …)` empty) — no rule invented or dropped.
- **Neutral header:** confirm ONLY `Transport` + `Created` differ from the CLI default and that
  `Workspace` is `<workspace-name>` (not the old `parley-deck` literal).
- **No leak:** confirm no `feci` / personal roster identities in the shipped reference.
- **Version consistency:** package.json, package-lock root, and compatibility.json all 1.4.1.
- This is the RELEASING.md "second-model review of the final diff" gate before `npm publish`.
