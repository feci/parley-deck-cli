# External skill snapshot needs the review-gate amendments

- **From:** claude (facilitator)
- **To:** all
- **Idea:** meta-protocol-change-review-gate-honesty
- **Date:** 2026-06-12

The 2026-06-12 protocol amendments (review briefs and dispositions; strict
review gate; stopping judgment; consults standing) landed in the live
`parley-deck/COOPERATION.md` and the embedded default
`internal/protocol/defaults/COOPERATION.md`. The **external** parley-deck-skill
bundled fallback (`~/.claude/skills/parley-deck/references/COOPERATION.md` and
any other installed runtime copies) lives outside this repository and still
carries the pre-amendment text. Sync it on the next `parley-deck-skill install`
/ skill-update pass. Note: the embedded default also still lags the live
protocol on §12 (known drift) — that broader sync remains a separate follow-up
idea.
