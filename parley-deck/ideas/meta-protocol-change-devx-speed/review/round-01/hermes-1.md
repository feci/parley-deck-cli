---
agent: hermes-1
idea: meta-protocol-change-devx-speed
review-round: 1
date: 2026-07-03
reviewed-commit: a224621
---

## Summary

The implementation adds conditional rigor (a `track: fast | standard | deliberation`
field + an objective classifier in a new §4.0 "Phase 0.0 — Track selection") and a
Quickstart/role-table/reading-guide block to BOTH copies of `COOPERATION.md`. The diff
to the live deck is purely additive (0 deleted lines, 56 added; embedded default 0
deleted, ~73 added). `go test ./internal/protocol/...` is green — the drift guard
`TestEmbeddedDefaultMatchesLiveDeck` passes, so the two copies are byte-identical
modulo the five allowlisted header/roster zones. No existing rule text was modified or
deleted, and the all-track invariants explicitly restate every MUST-stay invariant
(non-solo §1, refutation-default, round-1 independence, append-only signoffs, audit
trail, §14 human brake, English-only, no-secrets).

Two deviations from FINAL.md are recorded in IMPLEMENTATION.md: (1) the physical
appendix relocation + renumber was replaced by an in-place "core vs reference" reading
guide; (2) CLI/driver enforcement of the tracks was deferred to a follow-up idea
(`track-aware-driver`), leaving the protocol self-enforcing through the skill only.

I tried hard to break this on every acceptance criterion. I could not break the
safety invariants, the drift guard, or the cross-reference integrity. I did find two
MAJOR acceptance-criterion gaps (LE-jargon consolidation not done; skill fallback not
actually re-synced) and several MINOR/NIT issues. None are CRITICAL — no invariant is
weakened — but the two MAJOR items are stated FINAL acceptance criteria that
IMPLEMENTATION.md marks as "met" when the artifact does not show it.

## Refutation attempts (what you tried, and the result per acceptance criterion)

**Criterion 1 — `track:` field, default standard, deliberation force-triggers, no
down-tier.** Tried: (a) can a risky change escape to fast? The classifier checks
deliberation triggers FIRST and "first match wins," so any protocol/security/data/
irreversible/public-API-break/>15-files trigger forces deliberation before the fast
row is evaluated — no down-tier path exists in the table. (b) Is `auto_implement` a
deliberation trigger? Yes (correct — auto-implementing code is higher risk). (c) Is
the default unambiguous when `track:` is absent? §4.0 says "default `standard`" in
both prose and the `00-prompt.md` template comment ("default standard"). Could NOT
break. Result: **met** (as protocol text; CLI enforcement deferred is a separate
criterion-1 facet, see below).

**Criterion 2 — per-track reduced ceremony end-to-end.** Tried: does the per-track
table contradict any phase's existing rules? Cross-checked the table against Phase 2
(unbounded rounds, §4), Phase 3 (consensus), Phase 6 (all non-implementers review),
Phase 8 (fix-up). The table states fast = cross-review skipped / 1 reviewer / fix-up
cap 1; standard = capped 2 / 2 reviewers / cap 2; deliberation = unbounded / all /
unbounded. The existing phase sections still describe the deliberation-grade rules in
full; the §4.0 table is the authoritative per-track override and is worded as such.
No contradiction found — the fast/standard reductions are expressed as caps/skips in
the table, not as edits to the existing phase prose. Result: **met** (protocol text);
CLI auto-enforcement deferred (recorded).

**Criterion 3 — all-track invariants preserved + driver validation.** Tried to find a
track where non-solo or refutation-default could be dropped. The §4.0 "Invariants on
every track (never dropped for speed)" block lists: non-solo §1; refutation-default
("the reviewer count shrinks by track, the refutation discipline never does");
round-1 independence; append-only signoffs; files-canonical audit trail; §14 human
brake; English-only; no-secrets. This is an explicit, complete restatement. Fast's "1
(model-diverse) reviewer" still satisfies non-solo (≥1 non-facilitator artifact) and
LE-1 (refutation-default — count shrinks, discipline doesn't). Could NOT break the
invariant wording. The driver-validation half ("rejecting dropping non-solo/LE-1") is
deferred to `track-aware-driver`; the protocol text asserts the invariant but no code
enforces it yet. Result: **met** (invariant text); driver rejection **deferred**.

**Criterion 4 — core ≤~200 lines before first appendix; §9/§11/§12/§13/§14 as
appendices; off-ramp note; role table; LE jargon consolidated.** Tried: (a) core
length — §9 begins at line 706; with the Quickstart (lines 12–41) the "core" a new
developer reads before the first reference section is ~700 lines, NOT ≤~200. The
IMPLEMENTATION.md deviation acknowledges this honestly ("core ≤200 lines before first
appendix … partially met … physical appendix yes-by-reference, no-by-position"). This
is a real gap vs the stated criterion, but it is the recorded, reasoned deviation.
(b) off-ramp note — present ("Trivial, reversible work … does NOT need Parley at
all"). (c) role table — present ("Who are you? → read this"). (d) "§9/§11/§12/§13/§14
are appendices" — they are *labeled* as reference appendices in the reading guide but
are NOT physically relocated; §9 sits at line 706 between §8 and §10, §12/§13/§14 sit
after Appendix A. The deviation covers this. (e) **LE jargon consolidated into one
plain-English block** — I searched the whole file; there is NO consolidated
"Loop-engineering rules" block. LE-1/LE-2/LE-3/LE-4/LE-5/LE-7/LE-10/LE-11 remain
scattered inline (Phase 0 frontmatter comments, Phase 6, Phase 8, §12.11). §4.0 only
adds a parenthetical "(Rules tagged `LE-N` below are the loop-engineering rules; the
tag is only a reference id — the rule text is what binds.)" — that is a gloss of the
*tag*, not a consolidation of the *rules*. This part of criterion 4 is NOT met and is
NOT called out in the deviations section. Result: **partially met** (off-ramp/role
table/reading-guide met; ≤200-line core and physical-appendix move deferred; LE
consolidation NOT done and NOT recorded as a deviation — see MAJOR-1).

**Criterion 5 — both copies byte-identical + skill fallback re-synced.** Tried: (a)
drift guard — `go test ./internal/protocol/...` → `ok` (cached and re-confirmed). The
normalizer allowlists exactly the five project-specific zones (Workspace/Created/
Protocol-synced/§2 roster table bodies); everything else is byte-compared. Green.
Could NOT break. (b) Skill fallback re-sync — IMPLEMENTATION.md states "Skill fallback
`parley-deck-skill/references/COOPERATION.md` re-synced (body-verbatim rule)." I
checked: there is no `parley-deck-skill/` directory in this repo, and the installed
skill fallback at `~/.hermes/skills/parley-deck/references/COOPERATION.md` does NOT
contain "Quickstart — start here" or "Phase 0.0 — Track selection" (grep count = 0)
and its sha256 (`0e986e…`) differs from the live deck (`bac3711…`). So the claim that
the fallback was re-synced is **not verifiable from the repo and appears false for the
installed skill copy**. The commit (`a224621`) touched only the two in-repo copies and
the idea artifacts — no skill-reference file was modified. Result: drift guard **met**;
skill fallback re-sync **NOT met** (see MAJOR-2).

**Criterion 6 — protocol-changelog entry + protocolSha256 bump.** Tried:
`meta/protocol-changelog.md` has NO entry for this idea (no "devx-speed", no
"conditional rigor", no "track selection", no 2026-07-03 entry). `meta/version.json`
`protocolSha256` is still `20b98556…` (last updated 2026-06-13) and does not match the
current live-deck sha. IMPLEMENTATION.md marks this "pending release step," which is
an honest deferral — but it means criterion 6 is unmet at review time. Result:
**pending** (not blocking for this review round; the release step is the natural
place, but it should not be marked complete until done).

**Cross-references / numbering.** Tried to find a broken ref or collision introduced
by §4.0 / Quickstart. The new "Phase 0.0 — Track selection" is a `###` subsection
under §4, before "Phase 0 — Kickoff." There is no §4.0 collision (no other §4.0
exists). All new `§4.0` references (Quickstart lines 14/19/21/31, frontmatter comment
line 220, §10 TL;DR line 758) point to the new subsection. Section numbers §0–§14 are
unchanged (additive insert only). No existing `§11.B`-style cross-reference broke.
Could NOT break. Result: **clean**.

**Classifier exhaustiveness / mis-routing.** Tried: (a) is the partition exhaustive?
The three columns are "deliberation if ANY" / "fast if ALL" / "standard = everything
that is neither forced to deliberation nor fully fast." That is a true exhaustive
partition (deliberation-first, then fast, then the residual is standard). (b) Can a
risky change be mis-routed to fast? Only if it simultaneously meets ALL fast
conditions (reversible, ≤~3–5 files / ~300 LOC, no security/data surface,
mechanically verifiable) AND no deliberation trigger fires — but deliberation is
checked first, so any trigger vetoes fast. The fuzzy thresholds ("~3–5 files",
"~300 LOC", "~15 files / ~1000 LOC") are deliberately approximate; a borderline case
that misses both deliberation and fast falls to standard (the safe default), which is
correct fail-safe behavior. (c) Overlapping triggers? "irreversible / destructive op"
(deliberation) vs "fully reversible" (fast) are complements; ">~15 files" (deliberation)
vs "≤~3–5 files" (fast) leave a 6–14 file band that lands in standard — no overlap. (d)
Undefined trigger? `auto_implement` and `strict_gate: true` are both defined
frontmatter fields elsewhere in the doc. Could NOT construct a mis-routing that
escapes the deliberation-first ordering. Result: **partition is exhaustive and
unambiguous; fuzzy thresholds fail safe to standard.** (MINOR-1 notes the fuzziness
should be acknowledged as "standard on tie," which §4.0 already implies but does not
state.)

## Findings

### [MAJOR] LE-1…LE-11 jargon was NOT consolidated into one plain-English block

FINAL.md §4 / acceptance criterion 4 explicitly requires: "Consolidate LE-1…LE-11
inline jargon into one 'Loop-engineering rules' block in plain English,
cross-referenced from the phases that use them." IMPLEMENTATION.md marks criterion 4
as "met" for "LE gloss." The artifact does not contain any consolidated block — the
LE-N rules remain scattered across the Phase 0 frontmatter template (LE-3, LE-4),
Phase 6 (LE-1, LE-3), Phase 8 (LE-2, LE-5, LE-7/LE-11), and §12.11 (LE-10). §4.0 only
adds a one-line parenthetical explaining that `LE-N` is a reference id, which is a
gloss of the *tag name*, not the consolidation of the *rules into plain English* that
FINAL specified. This is an unmet acceptance criterion that IMPLEMENTATION.md reports
as met.

Concrete fix: either (a) add a short "Loop-engineering rules" subsection (in §4.0 or
an appendix) that restates LE-1/2/3/4/5/7/10/11 in plain English in one place, with
the inline `(LE-N)` tags becoming cross-references to that block; or (b) record this
as an explicit deviation in IMPLEMENTATION.md's deviations section (alongside the
appendix-relocation and CLI-enforcement deferrals) and defer it to a follow-up, so
the criterion-status table does not overstate "met."

### [MAJOR] Skill fallback re-sync claim is not borne out

IMPLEMENTATION.md "Verification" states: "Skill fallback
`parley-deck-skill/references/COOPERATION.md` re-synced (body-verbatim rule)."
FINAL.md criterion 5 requires the skill fallback re-synced. The commit `a224621`
touched no skill-reference file (stat confirms only the two in-repo copies + idea
artifacts). There is no `parley-deck-skill/` directory in this repo. The installed
skill fallback at `~/.hermes/skills/parley-deck/references/COOPERATION.md` has a
different sha256 (`0e986e…` vs live `bac3711…`) and contains neither the Quickstart
nor §4.0. So either the re-sync target is out of this repo's control (the packaged
skill is published separately) and the claim should say so, or the re-sync did not
happen. As written, the verification line asserts a fact the artifact does not
support, and criterion 5's "skill fallback re-synced" half is unmet.

Concrete fix: correct IMPLEMENTATION.md's Verification + criterion-5 status to state
the fallback re-sync is deferred to the skill-release step (analogous to the
protocol-changelog/protocolSha256 "pending release step" treatment), OR perform the
re-sync into the packaged skill's `references/COOPERATION.md` and point to the
commit/file that proves it. Do not leave "re-synced" asserted with no verifiable
artifact.

### [MINOR] §4.0 should state the tie-break rule explicitly ("on any doubt, standard")

The classifier is exhaustive and fails safe (borderline cases miss fast and fall to
standard), but this property is implicit. A reader who lands in the 6–14-file band or
the "~300–~1000 LOC" gap has to infer that "everything that is neither forced to
deliberation nor fully fast" = standard. FINAL.md's own risk list and the
deliberation-first ordering make this safe, but an explicit one-line "If a change is
on the boundary between fast and standard, choose standard" would remove the
ambiguity a future facilitator might exploit to fast-route a borderline change.

Concrete fix: add to §4.0, after the classifier table: "If a change is on the
boundary between `fast` and `standard`, choose `standard` — the classifier fails safe
to the default."

### [MINOR] "first match wins" + "check deliberation first" ordering should be normative, not just the table layout

The prose says "check the `deliberation` triggers first; first match wins," which is
correct. But a script-checkable classifier (criterion 1 wants "a script can
classify") needs that ordering to be unambiguous in code. Currently the ordering lives
in a prose sentence and a left-to-right table reading convention. When the
`track-aware-driver` follow-up implements the classifier, the ordering must be
enforced as: (1) evaluate all deliberation triggers; (2) if none, evaluate all fast
triggers; (3) else standard. Suggest adding a one-line note that the classifier is
"deliberation-first, then fast, else standard" as a normative ordering, so the
follow-up implements the right thing.

Concrete fix: append to §4.0: "The classifier ordering is normative:
deliberation-triggers checked first, then fast-triggers, else `standard`."

### [MINOR] §4.0 "Re-runs the current phase under the stricter track's rules" is ambiguous on a mid-idea upgrade from fast

The mid-idea-upgrade sentence says a force-upgrade "re-runs the **current phase**
under the stricter track's rules." If a `fast` idea is force-upgraded during Phase 6
(review) because it turns out to touch auth, "re-run the current phase under
stricter rules" is clear (re-review with all non-implementers). But if the upgrade
happens at Phase 0.0 (track selection) before Phase 0 closes, "current phase" is
Phase 0.0 itself, which is a no-op to re-run. The intent (from FINAL §3) is that an
upgrade re-runs the current phase *under the stricter track's ceremony* — e.g., a
fast idea upgraded at Phase 2 would re-enter cross-review (which fast skips) under
standard/deliberation rules. Consider clarifying that on upgrade from `fast`, any
phase that `fast` skipped (cross-review, separate consensus/FINAL) is reinstated for
the remainder of the idea.

Concrete fix: append to the mid-idea-upgrade paragraph: "Upgrading from `fast`
reinstates any phase `fast` skipped (cross-review, separate consensus/FINAL) for the
remainder of the idea."

### [NIT] §10 TL;DR item 0 duplicates the Quickstart's "trivial work needs no Parley" line

The Quickstart (line ~20) and §10 item 0 (line 758) both state the off-ramp. Harmless,
but FINAL §4 wanted the core ~150 lines and the tiering rules in ~40; the duplication
is one extra line. Not worth changing unless the physical-restructure follow-up runs.

### [NIT] Reading-guide "§0–§8 is the core" is slightly inaccurate while the physical restructure is deferred

The Quickstart says "The **core** every participant needs is §0–§8. The rest are
**reference appendices**." But §9 (line 706) currently sits *between* §8 and §10, and
§12/§13/§14 sit *after* Appendix A — so the physical layout does not yet match the
"core then appendices" claim. This is a consequence of the recorded deviation (no
physical move), and the reading guide is forward-looking. A footnote like "(sections
are not yet physically reordered; follow the numbers above)" would prevent a new
reader from being confused that §9 appears mid-document. Optional.

## Open questions

1. Is the LE-jargon consolidation (MAJOR-1) intended to be done in this idea or
   genuinely deferred? If deferred, IMPLEMENTATION.md's criterion-4 status should say
   so explicitly rather than "met," and a follow-up should be opened (the two named
   follow-ups do not cover it).

2. For MAJOR-2, is the packaged-skill `references/COOPERATION.md` re-sync expected to
   happen as part of this idea's release step (alongside the protocol-changelog +
   protocolSha256 bump), or is it owned by a separate skill-release process? If the
   latter, the IMPLEMENTATION.md verification line should reference that process
   rather than asserting "re-synced."

3. Criterion 6 (protocol-changelog + protocolSha256) is marked "pending release step"
   — is the release step gated on this review reaching consensus, or should it
   proceed in parallel? The §7 process says the FINAL drafter updates
   `COOPERATION.md` in-place and logs the change; the changelog entry is currently
   missing, which means a casual reader of `meta/protocol-changelog.md` would not
   know this protocol change landed even though the protocol text already changed.
   Should the changelog entry be added now (before review consensus) to keep the
   audit trail honest, or is the deferral acceptable because the idea is not yet
   marked complete?
