---
idea: meta-protocol-change-lean-organizer
author: kimi-1
role: independent reviewer
artifact: release-wording-1.49.1-review-kimi-1
date: 2026-09-24
subject: Wording-only commit 54e07981740d495f08cd90842168d9228c58c90f (parent 868825f20bd0988abde1a6c38b9efe405c5123d1) — owner-authorized release wording for CLI 1.49.1
basis: owner decision parley-deck/inbox/user-to-codex-1_meta-protocol-change-lean-organizer_windows-scope.md (status authorized); prior review release-candidate-1.49.1-review-kimi-1.md (PASS at 481fb65); prior source acceptance taken as given — no product build/tests run
---

# Independent review — CLI 1.49.1 release wording — kimi-1

## VERDICT: **PASS**

Reviewed commit: `54e07981740d495f08cd90842168d9228c58c90f` against parent
`868825f20bd0988abde1a6c38b9efe405c5123d1`. Release-only wording resume; no
A–D or product-code review repeated.

## Checks performed (all pass)

1. **Changed paths exactly three, all declared.** `git diff --name-only
   868825f..54e0798` = `CHANGELOG.md`,
   `release-candidate-1.49.1-zcode-1.md`,
   `release-wording-1.49.1-zcode-1.md` (new). Nothing else.
2. **CHANGELOG: single hunk, 1.49.1 only.** One `@@` hunk in the 1.49.1
   "Release scope and limitations" bullet 1: pre-decision "choice is
   between… publication held pending" replaced by the recorded
   authorization. `git diff v1.49.0..54e0798 -- CHANGELOG.md` shows **0
   content-deletion lines** — the frozen 1.49.0 entry is byte-identical.
3. **Wording matches the owner decision verbatim in substance.** CHANGELOG,
   candidate report (intro, § 1, § 3, § 5, § 6, § 9) and the zcode wording
   report all record: release CLI 1.49.1 **now for macOS and Linux through
   GitHub and Homebrew**; **Windows assets kept and explicitly labelled
   experimental/unvalidated**; **CLI winget submission HELD** until a
   separate reviewed windows-portability track fixes the defects, turns
   hosted windows-latest CI green, and removes the label (until then every
   CLI release keeps the label and opens no CLI winget submission); **skill
   unaffected**; skill 2.13.0 npm stays with the **owner-facing Claude
   session and is reclassified as not a CLI gate** (§ 6 gate 3). All
   consistent with the authorized inbox decision.
4. **Proposed release notes checked.**
   `release-delivery/2026-09-24-lean-organizer/release-1.49.1/release-notes.md`
   is byte-identical to the CHANGELOG 1.49.1 entry body (minus the
   `# Changelog` header, plus the existing companion-skill-installer line
   "Skill npm publication is handled separately"). The notes carry the
   experimental/unvalidated Windows label and the winget hold in text;
   per the dispatch the GitHub assets will also be labelled
   experimental/unvalidated at publish time (organizer/owner action, not
   part of this wording commit).
5. **Historical statements preserved.** § 8 of the candidate report
   unchanged (dated record corrections, incl. the then-true "publication
   held pending"); § 9 quotes the pre-decision gate list as reviewed and
   records that the decision lifts the hold. Front matter keeps the
   original `amended:` line and adds a dated `amended-2:` line.
6. **Self-reference accurate.** § 6 gate 1 records "kimi-1 PASS at 481fb65,
   recorded by the organizer in 868825f" — matches my prior review and the
   parent commit.
7. **Tags unmoved; version metadata consistent.** `v1.49.0` = `06e563e`
   locally and at `git ls-remote origin`; **no `v1.49.1` tag** exists.
   `VERSION` = `1.49.1`, `internal/app/version.go` = `"1.49.1"`, CHANGELOG
   top heading `## 1.49.1 — 2026-09-24` — unchanged from the accepted
   candidate.
8. **No unrelated drift.** Closed artifacts (FINAL/IMPLEMENTATION/consensus)
   untouched by this commit (0 paths); peers, inbox, ledger, `runs/`
   untouched; working tree shows only the pre-existing untracked
   organizer/peer records — all preserved.

## Findings

None. No corrections requested.

## Out of scope / not re-done

Prior source acceptance (hosted Linux/macOS green at 9134c7a, token-identity
of 519951e) taken as given; no product builds, tests, or A–D review
repeated. Release operations (tag, builds, GitHub release, Homebrew, asset
labelling, winget hold enforcement) remain organizer/owner-only; this review
performs none of them and is not release completion.
