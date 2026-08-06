---
idea: roster-operations-standard
phase: 7 — review consensus signoff
agent: kimi-1
date: 2026-08-06
verdict: BLOCK
---

# Signoff — kimi-1

## Verdict

BLOCK — two codex-1 findings have no disposition anywhere in the draft, so §4's
"every finding from all three reviews was reproduced and is either agreed (§2) or
deferred (§3)" is false as written:

- **B1.** codex-1's MAJOR "machine-scope writes use the wrong file whenever
  `PARLEY_HOME` is set" (round-01/codex-1.md:261-277) appears in neither §2, §3,
  nor §4. It was a real defect at v1.40.0 and fix-up cycle 1 fixed it — record it
  as fixed-in-cycle-1, the same treatment the two CRITICALs got via the baseline
  line. Actionable: one entry; no code work.
- **B2.** codex-1's MINOR-1 "pin preview/`--keep` is not sufficient against typos
  or concurrent edits" (round-01/codex-1.md:279-296) is still open at v1.40.1 and
  has no disposition. DF-3 covers *my* N-min-7 (backup/cleanliness hint) — a
  different finding. Actionable: agree it into §2 (reject unmatched/typo'd
  `--keep` tokens — a few lines) or defer it explicitly in §3 with a reason.

Both are one-paragraph additions to the draft; I expect to flip to ACCEPT once
the disposition map is complete. Everything else in the draft I would sign.

## Answers to 1-6

### 1. My findings — all dispositioned, none dropped or misrepresented

`PRIMARY` (I mapped each of my round-01 items against the draft's §2/§3 just now):

- C1 residuals M1/M2 → A4, A5 (agreed, MAJOR). C2 residual M3 → A2 (agreed, MAJOR).
- M4 → A1, escalated to CRITICAL with my position adopted (see Q3).
- M5 → A6 (severity MAJOR upheld via C-2 — correctly not by count but by defect
  class: a silently wrong `--scope` answer vs. a rejected flag).
- M6 → A13. **Downgrade, named and accepted:** I filed it MAJOR; the draft splits
  it — MAJOR for `--help`, MINOR for docs. The agreed fix still covers both
  surfaces, and `--help` is the load-bearing half (D1's "and the docs" survives
  as MINOR work). Accept.
- M7 → A3; M8 → A10; M9 → A8 — all agreed at my filed severity.
- M10 → DF-1 (deferred; see Q4).
- N-min-1 → A12; N-min-2 → A15; N-min-3 → A14; N-min-4/N-min-6/N-min-8 → A15;
  N-min-5 → DF-2; N-min-7 → DF-3; NIT-1/NIT-4 → A15; NIT-2 → DF-5; NIT-3 → DF-6;
  NIT-5 → DF-4.

The draft's characterizations of my M3, M4, M7, and M10 are accurate, including
the honest DF-1 framing. My G5 "SATISFIED" verdict is resolved against me in C-1
— handled through §15.3, correctly (see Q2).

### 2. C-1 (G5 changelog format) — resolution is sound; my verdict is not upheld

`PRIMARY` — I re-read both artifacts now. §7's template (COOPERATION.md:744-748)
mandates four plain fields: `## YYYY-MM-DD — <desc>`, `Idea: ideas/.../`,
`Drafted by: <agent-id>`, `Summary: <1–2 sentences>`. The entry
(parley-deck/meta/protocol-changelog.md:119-139) has the heading but then uses
bold `**Idea:**` with a bare slug, `**Change:**`/`**Why:**` instead of `Summary:`,
and **omits `Drafted by:` entirely**. My round-1 "in §7 format" read the entry's
completeness, not the template's field names — I was wrong; codex-1's MINOR-2 is
correct. The resolution rests on a normative artifact in the repository, and the
draft says so explicitly (consensus.md:58-59: "a majority-of-one either way would
have been meaningless. The template ... decides."). §15.3-compliant in substance.
A11's fix (add `Idea:` as a path, `Drafted by:`, `Summary:`) matches §7. See
Conditions for the one formal gap (verbatim quoting).

### 3. A1 escalation — right resolution, right scope

Position match: `PRIMARY` — the agreed resolution (deck membership = deck file;
machine seeds values never membership; display-only inheritance with every
inherited row marked; `render` refuses to bake an inherited roster into §2
without an explicit flag) is exactly my M4 position (round-01/kimi-1.md:183-186)
and also answers codex-1's open question 2 (round-01/codex-1.md:357-360).

Severity: `SECONDARY` — the drafter's quoted reproduction (a deck declaring two
members shows five rows; an empty deck inherits the full machine roster with no
marker) is materially worse than either reviewer reported, and since 1.40.1
routes participant selection through the same layered view, it silently expands
the run quorum. I did not re-run those scratch decks this session; the mechanism
(layered `config.LoadRoster` consumed by both `RosterMembership` and
`renderRosterTable` while `rosterSync` reads the deck file alone) matches what I
verified PRIMARY in round 1 (roster.go:643, roster_render.go:30,
roster_sync.go:46-50). CRITICAL is the right call: it defeats the non-solo
quorum invariant and re-creates the committed-drift vector D9 exists to kill.
Scope (show + render + membership/quorum, all three consumers of the layered
view) is complete — nothing outside those consumers needs the rule.

### 4. DF-1 deferral — acceptable

My M10 said "fix it *before* it bites 40 repositories" on the stated premise that
the fleet operation "has not run" (round-01/kimi-1.md:286). That premise no longer
holds. Fleet numbers (24 applied / 9 skipped / 3 unchanged / 0 failed, per-deck
backups, validation with rollback): `SECONDARY` — the drafter's testimony in the
draft (§3 DF-1); the fleet is ~40 external decks, not verifiable from this repo.
Corroboration, `PRIMARY`: CHANGELOG.md:32-34 (1.40.1) documents migrate's
dry-run default, file-level backups, post-write validation with automatic
rollback, and skip-and-report; this deck's own `parley-deck/agents.toml` carries
4 `[roster.*]` blocks, consistent with the attended run having begun here. Given
that, the pre-flight-blocker rationale is moot, the finding is upheld in full as
a hardening task, and the cheap guard (dirty-tree skip + `--confirm-breaking`)
closes the residual risk class — a second *unattended* run. `PRIMARY` note: that
guard does not exist in the tree yet (`rg confirm-breaking|dirty` over
roster_migrate.go → no hits) — correct, it is cycle-2 work, and the draft frames
it as "ships in this cycle", not as landed. Accept, on the recorded condition
that DF-1 stays open for any future fleet run.

### 5. §4 dismissals — nothing should be dismissed; §4 is instead *incomplete*

I considered the candidates. hermes-1's "§2 copies not identical" MINOR
(round-01/hermes-1.md:276-297) is self-described as the guard "working as
designed" — folding it into A9's drift assertions is right, not dismissable.
codex-1's PARLEY_HOME MAJOR is a real defect (verified below) and must be
recorded as fixed, not dismissed. So: no finding merits dismissal — but §4's
completeness claim fails because two codex-1 findings have no disposition at all
(B1, B2 in Verdict). Evidence:

- B1 `PRIMARY`: at v1.40.0, `rosterScopeFile` composed
  `$PARLEY_HOME/.parley/agents.toml` (`git show 203f73b:internal/app/roster_set.go`,
  lines 76-92) while the loader reads `$PARLEY_HOME/agents.toml`
  (`CentralHome`/`CentralAgentsPath`, internal/config/runtime.go:414-431). de269c4
  fixed it via `config.CentralAgentsPath()` (current roster_set.go:89-103; commit
  message attributes it to codex-1; CHANGELOG.md:24-28). Fixed, but unrecorded in
  the consensus.
- B2 `PRIMARY`: at v1.40.1 (working tree = 58db960, clean), `roster sync` still
  lowercases `--keep` tokens into a map with no validation of unmatched entries
  (internal/app/roster_sync.go:52-55) and still computes drops from one read
  (:46) while deleting from a second (:126-135) with no hash/old-value binding —
  exactly codex-1's MINOR-1, still open, no §2/§3 entry.

### 6. Collectively missed / must be in this cycle

- The two blocking items above (B1 record-keeping, B2 open safety nit).
- 1.40.1 CHANGELOG.md:8 says the defects were "found by codex-1 and hermes-1".
  Both CRITICALs were also independently corroborated in my round-1 file
  (round-01/kimi-1.md:55-106) — a sentence-level attribution fix can ride along
  with cycle 2. *Per §15.1 the independence of my own review is a claim I own; I
  issue no verdict on it. Stated for the record only.*
- Nothing else. `PRIMARY` (I read all three round-01 files in full this session):
  no third defect class exists that none of us filed and that cycle 2 must absorb.

## Conditions

1. **B1** — add a "fixed in fix-up cycle 1" record for codex-1's PARLEY_HOME
   machine-scope MAJOR (one line in §0/baseline or a §2 note), so §4's
   completeness claim becomes true.
2. **B2** — disposition codex-1's MINOR-1: agree it into §2 (reject unmatched
   `--keep` tokens; optionally bind apply to the previewed file hash) or add a
   DF entry with a stated reason.
3. §15.3 formal letter (`PRIMARY`, COOPERATION.md:1274-1277): the conflicts
   section must quote each conflicting verdict verbatim with author, tag, and
   evidence. The draft's §1 paraphrases. The verdicts are canonically preserved
   in round-01/kimi-1.md:395-396 and round-01/codex-1.md:298-307, so this is a
   paste, not new analysis — non-blocking condition.
4. §15.5 formal letter (`PRIMARY`, COOPERATION.md:1299-1304): the draft records
   the role concentration in substance (§0) but lacks the mandated
   `## Drafter position changes` section. The design-phase consensus by the same
   drafter has one (consensus.md:248); this one needs it too — the A1 escalation
   against claude-1's own 1.40.0 release claims is a drafter position change and
   belongs in it (or `None`, though I do not see how `None` would be accurate).
   Non-blocking condition.
