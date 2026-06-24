---
agent: codex-1
idea: sync-skill-protocol-fallback
round: 1
date: 2026-06-24
---

## Recommendation

I counter claude-1's strict recommendation (A) only at the header boundary. The faithful sync should be a narrow hybrid:

- Treat `internal/protocol/defaults/COOPERATION.md` as the canonical source for the protocol body.
- Replace the stale skill fallback body with the canonical body, importing every new rule and section.
- Preserve a fallback-neutral header instead of copying the CLI default's active-looking `Transport: github-pr` and `Created: <date> - created by parley init` lines.

Concretely, the skill fallback should use a header like:

```markdown
**Workspace:** `<workspace-name>`
**Parley deck:** `./parley-deck/`
**Transport:** `<transport-choice>` (pick one of `local-dir`, `github-pr`, or `gitlab-mr` at deck bootstrap - see §0)
**Created:** `<YYYY-MM-DD>` (set at deck bootstrap)
```

Everything after the header should track the CLI embedded default verbatim unless there is a deliberately documented fallback-header exception. This preserves the single-source-of-truth property for real protocol rules while avoiding a fallback file that silently claims GitHub PR transport or `parley init` provenance for agents that may not be using the Parley CLI.

Version bump: patch, `1.4.0 -> 1.4.1`. This is a bundled protocol snapshot sync / safety correction, not a new installer API or user-facing command. Update `package.json`, `package-lock.json`, and `references/compatibility.json` `skillVersion` together.

## Analysis

The diff confirms the prompt's summary: the stale skill fallback is 759 lines and stops at §12, while the CLI embedded default is 1037 lines and includes §§13-14. The unified diff is 447 lines and has substantive changes beyond the two missing sections.

Required imports from the CLI default:

- Header/template zone: adopt `<workspace-name>` from the canonical template, but keep fallback-neutral `Transport: <transport-choice>` and `Created: <YYYY-MM-DD>`.
- §0: import the one-time deck bootstrap paragraph requiring roster, model, and reasoning/effort confirmation; strongest available reasoning by default; central config inheritance; §9.0 remains liveness-only.
- §2: import the local launch config note and the current empty roster-table shape unless the implementer intentionally keeps placeholder rows as a documented fallback-only cosmetic delta.
- §3: import the stronger descriptions of `FINAL.md` as static/self-contained and `IMPLEMENTATION.md` as the living execution document.
- §4 Phase 0: import `strict_gate`, `require_model_diversity`, and `checks` frontmatter fields.
- §4 Phase 3 and Phase 7: import `Comparison & blind spots` / `Coverage & blind spots`.
- §4 Phase 4: import the expanded `FINAL.md` required sections: purpose, context, observable acceptance criteria, idempotence and recovery, risks, and the rule that complex / auto-driven / pipeline ideas must be resumable from `FINAL.md` plus `IMPLEMENTATION.md`.
- §4 Phase 5: import the expanded `IMPLEMENTATION.md` living sections: progress, decision log, surprises, validation evidence, and outcomes/retrospective.
- §4 Phase 6: import `## Refutation attempts`, the refutation-default rule, model-diversity handling, and the review-briefs/dispositions rule that prevents suppressing findings.
- §4 Phase 8: import the strict review gate, driver enforcement, stopping judgment, loop-budget escalation rule, and close-decision integrity / goal-done check.
- §5: import quorum lock at §9.0 readiness and the user-confirmed exclusion rules.
- §7: import the version-sync carve-out so upstream-ratified additive syncs are maintenance, not protocol-change ideas.
- §8: import the `Consults` note.
- §9: import §9.0 pre-idea readiness check, including protocol freshness and roster liveness.
- §12.11: import the `status: candidate` watcher-remediation rule.
- §13: import Retrospective optimization.
- §14: import Automated outer loop - the human brake.

On `parley` CLI literals: they are acceptable in the body when they are concrete examples or names of Parley ecosystem tooling, not preconditions for a non-Parley agent. Do not scrub `parley-deck/`, `parley-deck-skill status`, `parley preflight`, `parley consult`, `parley retro`, `parley run`, or `parley loop tick` out of the protocol body. Over-scrubbing would fork canonical wording and risk dropping safety semantics. The header is different: `Transport: github-pr` is an active project setting, and `created by parley init` is provenance. Those are not mere examples, so they should stay generic in a bundled fallback.

If the group chooses claude-1's pure A anyway, it is still safer than leaving the fallback stale because it imports §§13-14 and the Phase 6/8 safety rules. I just do not think an empty full-file diff is worth importing misleading header state into a portability fallback.

## Verification

Before publish, verify the protocol sync mechanically:

- `wc -l references/COOPERATION.md` should be 1037 if only the two header lines differ from the CLI default.
- `diff -u ../parley-deck-cli/internal/protocol/defaults/COOPERATION.md references/COOPERATION.md` should show only the expected fallback header hunk.
- `diff -u <(tail -n +8 ../parley-deck-cli/internal/protocol/defaults/COOPERATION.md) <(tail -n +8 references/COOPERATION.md)` should be empty.
- `rg -n "## 13\.|## 14\.|Refutation attempts|strict_gate|Loop budgets|Close-decision integrity|status: candidate" references/COOPERATION.md` should find the imported safety sections.
- `rg -n "feci|claude-1|codex-1|hermes-1|antigravity-1" references/COOPERATION.md` should return no project-specific roster leakage.

Then run the skill release preflight from `RELEASING.md`:

- `npm test`
- `npm pack --dry-run`
- `npm run build:portable:current`
- `node bin/parley-deck-skill.js install --target all --dry-run`
- `node bin/parley-deck-skill.js doctor --target all --json`

Also verify the version files agree (`package.json`, root entries in `package-lock.json`, and `references/compatibility.json`) and ask the required second-model review before `npm publish --access public`. After publishing, create and push tag `v1.4.1` per `RELEASING.md`; release-asset, WinGet, and Homebrew verification follow the existing release document.

## Open questions

- Should the implementation include a small drift test now, or leave that as a follow-up? I would keep it as a follow-up unless the implementer can add a tiny body-diff test without expanding scope.
- If maintainers strongly prefer exact full-file identity with the CLI default, are they willing to accept `Transport: github-pr` as the fallback's apparent active transport? My recommendation is no; preserve the neutral header and verify body identity instead.
