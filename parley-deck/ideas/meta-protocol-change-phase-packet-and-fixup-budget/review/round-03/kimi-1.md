---
agent: kimi-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 3
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review, round 3 — release candidate 1.44.0 + 2.8.0

## Summary

**[PRIMARY — CONFIRMED]** The mandated round-3 object — the fix-up-cycle-2 candidate (budget
counted from the driver's `.fixup-done` markers alone) — is **NOT CLEAN on two independent
defects I confirmed with my own hands before the tree moved**: (1) the marker-only count is still
fail-open, because deleting markers lowers it — the editable state moved, it did not close
(MAJOR); (2) the CLI CHANGELOG's "Wrong unit" paragraph still described the prose-heading source
that round 2 killed, i.e. the release notes re-asserted the CRITICAL design as the shipped fix
(MAJOR). Both were found independently here; the implementer's later IMPLEMENTATION.md sections
attribute the same points to @codex-1 and @hermes-1 round-3/4 reviews, which I have **not read**
(provenance note below).

**[PRIMARY — CONFIRMED]** For the second consecutive round, **the working tree mutated while this
review was running** — three times: 15:40–15:41 (fix-up cycle 3: MAX of cursor + markers, AF2
budget gate, exact round-dir match), 15:44:47 and 16:03:36 (CHANGELOG rewrites), 16:02:49
(fix-up cycle 4: reserve-before-Fixup, AF2 `>` boundary fix, cursor-error escalation), 16:03:58
(IMPLEMENTATION.md cycle-4 section). My green/red evidence below is pinned to snapshots by hash
and mtime; nothing about the tree is certified beyond what is pinned.

**[PRIMARY — CONFIRMED]** The current tree (cycle-4 state, sha256 below) is the strongest version
yet, and I verified it end to end: full suite green (26 ok + 1 no-test), all four reversion
checks red exactly as claimed, eight adversarial probes of my own in an isolated copy mapping
every tamper direction, skill package green (386/386, manifest `--check` ok). Both CHANGELOGs are
now accurate, including an honest "What this is not" trust-boundary paragraph. One NIT remains
even in the newest state (a stale call-site comment), and two cycle-4 safety behaviors are pinned
by no tree test. Verdict is NOT CLEAN for the round-3 candidate as mandated; the in-flight
successor needs one **frozen** pass to certify, and it should be short.

Provenance: I have not read `review/round-03/codex-1.md`, `review/round-03/hermes-1.md`, or
anything under `round-04/`/`round-05/`. My findings under Q1/Q5 were derived from the pre-mutation
snapshot and my own probes before I saw any of the successor code. What I know of the other
reviewers' verdicts comes only from the implementer's own summaries in IMPLEMENTATION.md — tagged
[SECONDARY] where relied on at all.

## Review-validity event (read before any finding)

**[PRIMARY — CONFIRMED]** Timeline (file mtimes via `ls -lT`, my shell history; my session opened
15:25:51 CEST):

```text
~15:26–15:31  I read the diff and code of the cycle-2 candidate; go build/vet green
~15:32–15:35  go test ./... green (27/27); focused driver/track tests green  [snapshot S2]
15:36:38      review/round-03/hermes-1.md written   (not by me, not read by me)
15:40:04      review/round-03/codex-1.md written    (not by me, not read by me)
15:40:36–15:41:59  cursor.go, impl.go, impl_test.go rewritten   (fix-up cycle 3)
15:44:47      CHANGELOG.md rewritten ("Wrong unit, and then a wrong source")
~15:45        my isolated copy made — it already contained the cycle-3 code
15:46:12      IMPLEMENTATION.md: "## Fix-up cycle 3" section
16:02:49      impl.go + driver.go rewritten again   (fix-up cycle 4)
16:03:36      both CHANGELOGs gain the "not a security boundary" disclosure
16:03:58      IMPLEMENTATION.md: cycle-4 section (cites a round-4 codex-1 NOT CLEAN, hermes-1 CLEAN)
```

Snapshot pins (sha256, current tree at 16:12 CEST — **S4**, the cycle-4 state):
`impl.go e6aa4c7d…`, `driver.go a88f7ddb…`, `cursor.go 0df19271…`, `consensus.go 1330db98…`,
`impl_test.go 6418965c…`, `CHANGELOG.md f9247104…`. The cycle-2 snapshot (S2) content I reviewed is
recoverable from my round-3 tool record: full `git diff` output plus Reads of `impl.go` lines
100–340 and 470–525 taken ~15:28–15:31.

**[PRIMARY — CONFIRMED]** My footprint on both working trees is zero: every reversion and probe ran
in `mktemp -d` isolated copies of the module (the discipline my own round-2 process failure
helped establish); both `git status` outputs at the end of this review are byte-identical in file
set to those at the start. The copies were deleted after the checks.

## Findings

### F1 [MAJOR] The cycle-2 marker-only count was still fail-open: deletion lowers it

**[PRIMARY — CONFIRMED]** Against S2, `publishedFixupCycles` (S2 `impl.go:506-525`) counted
`.fixup-done` markers and nothing else. The four round-2 prose failure modes were closed (missing
file → 0 is unreachable past the `review/consensus.md` gate; unreadable dir → escalates; no prose
to misparse; the implementer's document no longer feeds the count). But the source files are
untracked dotfiles in participant-writable round directories — not `.gitignore`d, I checked — so
**deleting the markers resets the budget to zero**, the same unsafe direction as round 2, reachable
without touching any document (wipe/restore of round dirs, `git clean -fdx`). A careless heading
edit was replaced by a careless directory operation: less plausible, same class. A number that is a
safety boundary was still authored (and un-authored) by the party it constrains.

Status in S4: **FIXED, verified by me** — the count is `MAX(run-cursor counter, markers)`
(`impl.go` `publishedFixupCycles`, S4), the cursor lives in the run directory and is carried
forward monotonically across `Rebuild` (`driver.go:235-247`), and `TestFixupBudgetIsTamperFailSafe`
pins both directions. My probes confirm the residual boundary precisely (Q1): with **both** records
gone the count still resets — disclosed verbatim in both CHANGELOGs ("not a security boundary…")
and punted to the named follow-up `fixup-budget-trust-anchor`. That disclosure is honest and the
deferral is defensible; it is also exactly the scope line @codex-1's round-4 CRITICAL disputes
[SECONDARY — the implementer's summary], so the round in progress should rule on it, not I.

### F2 [MAJOR] The 1.44.0 CHANGELOG described the round-2-CRITICAL design as the shipped fix

**[PRIMARY — CONFIRMED]** Against S2 (CHANGELOG read ~15:29), the "Wrong unit" paragraph said the
count "is now derived from the `## Fix-up cycle N` records actually published in
`IMPLEMENTATION.md`" — the derivation round 2 proved fail-open in four ways. My round-2 review
warned verbatim that this paragraph "becomes stale and must be rewritten in the same patch" if the
marker counter landed. It landed; the paragraph was not rewritten until 15:44:47. Release notes
misdescribing the safety mechanism — and re-asserting the killed design — are this deck's core
defect class.

Status in S4: **FIXED, verified by me** — the entry now reads "Wrong unit, and then a wrong
source", describes the MAX-of-two-records design, the reserve-before-write rule, and carries a
"What this is not" paragraph whose every clause I probed (Q5). The skill CHANGELOG matches.

### F3 [NIT] The call-site comment still describes the cycle-1 prose derivation — even in S4

**[PRIMARY — CONFIRMED]** `impl.go:296-302` (S4, sha256 `e6aa4c7d…`) still says: "The cycle number
is derived from the fix-up cycles actually PUBLISHED in IMPLEMENTATION.md, not from the
review-round ordinal." Nothing below that comment has read `IMPLEMENTATION.md` since cycle 2; the
`publishedFixupCycles` doc comment four screens down says the opposite, correctly. This is the
round-1-NIT class (comments describing the old behavior) sitting in the one function this whole
idea is about. One sentence to fix; should not survive the release.

### F4 [observation] AF2's first budget gate was one off in the harmful direction — fixed at 16:02

**[PRIMARY — CONFIRMED]** The 15:41 AF2 gate used `spent >= cap`, which strands the legitimate
crash recovery of the **last allowed** cycle: at the inclusive boundary the marker at the current
round means cycle N's budget is already spent, and AF2 is finishing it, not starting N+1. I built
the boundary probe (cap 2, markers in rounds 1–2, crash window at round 2): on the 15:41 code it
escalates ("budget exhausted at 2 cycle(s)", no round opened — a stranded final cycle); on the
16:02 code (`spent > cap`) it completes the transition (`action=fixup`, `open-review` called, no
Fixup re-run). Verified in both directions by swapping the two `impl.go` versions in the isolated
copy against the same probe. No action needed; recorded because it was introduced and fixed inside
this review's window, and my probe is its only published evidence.

## Answers to the requested questions

### Q1 — is the fail-open counter actually closed?

**Against S2 (the mandated candidate): no — F1.** Deleting the markers extends the budget;
`fileExists` also swallows per-marker stat errors, so a permission-damaged round dir silently
under-counts. Forgery inflates (safe direction); oddly named `round-*` dirs inflate (safe); the
missing-dir → 0 branch is unreachable past the consensus gate and benign; an unreadable review dir
escalates. The one unsafe direction — deletion — was live.

**Against S4: closed to the standard this deck applies, with two mapped residuals, both disclosed
or minor.** Each attack from the brief, run by me on S4 in the isolated copy:

- **Missing review dir** → `IsNotExist` → 0. Unreachable past the `review/consensus.md` existence
  gate in `advanceReview`; benign where reachable. Closed.
- **Unreadable review dir** (`chmod 000 review/`) → escalates, no fix-up runs (it fails even
  earlier, at the impl-status gate, because `Rebuild` can no longer see any rounds). Closed.
- **Unreadable/corrupt cursor** → `Advance` escalates: "refusing to act on an unknown fix-up
  budget" (probe 7; cycle-4 fix). A genuinely **absent** cursor is a fresh run — markers only,
  correct. Closed.
- **Delete the markers** → cursor holds the count; still escalates at the cap
  (`TestFixupBudgetIsTamperFailSafe` subtest 1; my reversion of the carry-forward turns exactly
  this subtest red). **Delete the cursor** → markers hold. **Delete BOTH** → the count restarts
  (probe 1: `action=fixup` where the cap demanded escalation). That is the disclosed residual:
  CHANGELOG "What this is not", IMPLEMENTATION.md cycle-3/cycle-4, follow-up
  `fixup-budget-trust-anchor`. A participant with write access to both trees can also edit the
  driver binary; the trust-model boundary is stated, not papered over.
- **Forge a marker** (any round, including the current one) → the count only rises → escalates
  sooner; the AF2 branch now consults the budget before opening anything (subtest 2; my removal of
  that gate turns exactly this subtest red). Safe direction only.
- **Oddly named round dir** (`round-backup`, `round-x`, `round-`, `rounds-03`) → not counted:
  `isRoundDirName` accepts only `round-<digits>`; `TestOnlyExactRoundDirsCountAsPublishedCycles`
  pins 6 names → 2; my prefix-matching revert turns it red (counted 5). Closed.
- **AF2 interaction** → three probes: the legitimate crash window still completes (marker at the
  current round, budget not exhausted → archives, opens next round, no Fixup re-run); the inclusive
  boundary is correct (F4, both directions); a forged current-round marker past the cap escalates.
  The round-2 collision I flagged (fixture planting a marker at the live round) is gone from the
  fixtures and from production reach.
- **Per-marker stat swallow** → `fileExists` reads any stat error as "absent". With the cursor
  also gone, a permission-damaged round dir silently drops the count by one and the budget extends
  (probe 3: `action=fixup` at cap 2 with 2 truly spent). Requires two coincident independent
  failures; **MINOR residual** — one-line hardening available later (`statRegular` + propagate),
  not a blocker.

So: nothing reachable by careless operation of the protocol extends the budget anymore. The two
remaining extensions need either two independent failures or deliberate workspace tampering, and
both are disclosed in the release notes with a named follow-up.

### Q2 — does counting only gate-passing fix-ups hand back budget after a failed fix-up?

**It did in S2; S4 spends the cycle when it runs, and that is the right call — verified.**

My independent pre-mutation analysis of the cycle-2 semantics leaned the other way: a fix-up whose
post-fix-up checks fail publishes nothing to reviewers, and AF1 escalates immediately, putting a
human in the loop per failure — so the marker definition looked defensible, arguably closer to
"published" than the implementer's heading (which is written before checks run). The counter-
reading [SECONDARY — attributed to @codex-1 by the implementer's cycle-3/4 notes] is the one the
implementer adopted, and on the merits it is correct and I now prefer it: the cap exists to
interrupt exactly the churn of a failing fix-up, "never resets the count" reads most safely as
"never un-spends a run", and any design where a failing fix-up is free depends on every escalation
actually halting the loop. S4 reserves the cycle **before** the code-writing call
(`impl.go:317-325`) and persists it before the check gate, so neither an errored `Fixup` nor a
crash in the window hands the cycle back; a pre-confirmation crash conservatively burns one cycle
— the right error direction. My probe 8: an erroring `Fixup` at 1-spent/cap-2 leaves the next
`Advance` escalating ("after 3 cycle(s) (MaxFixupCycles=2)") instead of retrying cycle 2 for free.
The failed-fix-up loop against an undepleting ceiling is closed.

One cosmetic note from that probe output: the message says "after 3 cycle(s)" when 2 are spent and
the 3rd is the one refused — pre-existing wording, uses the would-be cycle number. Not worth a
finding; mention only because a human reading the escalation may misread the count.

### Q3 — is `standard`'s newly bounded BLOCK back-edge correct and disclosed in both CHANGELOGs?

**[PRIMARY — CONFIRMED]** Yes on both, unchanged since S2 and re-verified on S4.

- Code: `driver.go:143-148` (S2 hunk, still present) sets `HardCrossReviewCap` from
  `CapCrossReviewRounds` for every track whose policy prints one — `standard` 2, `deliberation` 3;
  `consensus.go:92-96` escalates at `next > 1+cap` before `RunRound`. `standard`: rounds 2–3 open,
  the third post-round-1 round (round 4) escalates — exactly its printed "capped at 2, then
  escalate/upgrade". `deliberation`: round 5 escalates;
  `TestBlockedConsensusRespectsTheHardCrossReviewCap` asserts escalation, names `§4.0 cap=3`, and
  asserts **zero runner calls**. Absent-track ideas keep `HardCrossReviewCap=0` and the legacy
  `MaxRounds` bound — backward compatible. My round-2 guard-removal reversion (red) still stands;
  the file is untouched by cycles 3–4 (sha256 `1330db98…` throughout my session).
- Disclosure: CLI CHANGELOG — "This also bounds `standard`'s back-edge at its own printed cap of
  2, which it previously ignored". Skill CHANGELOG — "binds the consensus-BLOCK back-edge as well,
  on `standard` (cap 2) as well as `deliberation` — previously that path ignored both." Both
  accurate (the back-edge was governed by `MaxRounds` alone). The round-2 NIT is also fixed: the
  escalation message now ends "escalating for human review" — no "trajectory" promise.

### Q4 — tests run by me, including the marker-count reversion

**[PRIMARY — CONFIRMED]** Everything below executed by me this session; all destructive checks in
`mktemp -d` isolated copies of the module, never the working tree (copies deleted afterward).

```text
TREE (S2, ~15:32):  go test ./... -count=1            26 ok + 1 no-test, zero FAIL
TREE (S2, ~15:35):  focused driver/track (7 tests)    all PASS (incl. all 6 boundary subtests)
TREE (S4, ~16:10):  go build/vet                      exit 0
TREE (S4, ~16:10):  go test ./... -count=1            26 ok + 1 no-test, zero FAIL
SKILL (unchanged since 15:23):  node --test           386 pass / 0 fail
                    python suite                      54 OK;  build-addon-manifest.js --check  ok (all six skills)
PROTOCOL TEXT:      the two changed §4.0 cells        byte-identical across all three COOPERATION.md copies

REVERSIONS (isolated copy; each reverted, confirmed red, restored, re-verified green):
  prose-parsing revert of publishedFixupCycles   → red: all 3 one-past-cap subtests,
       TestPhaseReviewMaxFixupCyclesEscalates, BOTH tamper subtests   [on 15:41 state AND on S4]
  carry-forward removed (driver.go)              → red: exactly "deleting the markers…" subtest
  AF2 budget gate removed                        → red: exactly "a forged marker…" subtest  [both states]
  isRoundDirName → HasPrefix("round-")           → red: TestOnlyExactRoundDirs… (counted 5, want 2)

MY PROBES (isolated copy, S4, all pass; file zz_kimi_probe_test.go, 8 tests):
  1 both records deleted        → count restarts (action=fixup)         [disclosed residual boundary]
  2 review/ chmod 000           → escalates, no fix-up                  [fail-closed]
  3 round dir chmod 000, no cursor → swallowed stat under-counts by 1   [MINOR residual]
  4 AF2 crash window            → completes: open-review, no fixup re-run
  5 escalation text             → names count + cap; no "trajectory"
  6 AF2 at spent==cap           → completes on S4; STRANDED-escalation on 15:41 code (F4)
  7 corrupt cursor              → escalates "unknown fix-up budget"     [cycle-4 claim; no tree test]
  8 errored Fixup               → reservation spent; next Advance escalates  [cycle-4 claim; no tree test]
```

Gap noted: the two cycle-4 behaviors (cursor-error escalation, reserve-before-Fixup) are pinned
by **no tree test** — my probes 7–8 are their only behavioral evidence. Both are safety branches;
a one-test-each pin belongs in the release or a named fast follow.

### Q5 — is anything in either CHANGELOG or IMPLEMENTATION.md still overstated?

**[PRIMARY — CONFIRMED]** No overstated claim found in the current texts; every load-bearing
sentence I could probe, I probed:

- "maximum of two driver-authored records… Losing one record does not lower the count, forging one
  can only raise it" — tamper subtests + my reversions A/B. "an unreadable cursor escalates
  instead of counting as zero" — probe 7. "reserved before the code-writing call" — probe 8.
  "crash-recovery path consults the cap… while still being allowed to *finish* the last allowed
  cycle" — probes 4 and 6 (both directions).
- "What this is not… not a security boundary: a participant with workspace write, a deleted run
  directory, a repository rollback, or two concurrent runs of the same idea can still reduce or
  duplicate the count" — consistent with my probes 1 and 3; the honest boundary statement F1
  required.
- "Why 5" distribution — reproduced digit-for-digit in round 2 (`0×17, 1×34, 2×7, 3×2, 4×3, 5×2,
  then 9, 14, 15, 25`); the calibration counted headings while enforcement now counts spent
  cycles — a slight unit drift at the edges, immaterial to the 6–8-band argument.
- "Not in this release" — packet not started; trajectory payload absent (`trajectory` appears only
  in changelogs/idea text, no Go payload); `fast` Phase-8 route manual (`track.go` still rejects
  `fast + auto_implement`). All still true on S4.
- IMPLEMENTATION.md cycle-3 verification claims (three isolated-copy reversions) — independently
  reproduced by me, red exactly as claimed. Cycle-4 claims — my probes 6/7/8; note the cycle-4
  section itself claims only build/vet/suite and offers no reversions (gap above).
- Remaining text defects: F3 (stale call-site comment, S4) and the cosmetic "after 3 cycle(s)"
  wording (Q2). Nothing else.

### Q6 — should anything stop this release?

**Yes, narrowly.** The mandated round-3 candidate carried F1+F2 (both MAJOR) — that is the verdict
of record for this round. Against the **current** tree, what remains is smaller but real:

1. **F3 (NIT)** — the stale call-site comment at `impl.go:296-302` contradicts the function it
   annotates, in the exact function this release exists about. Trivial fix; must not ship as-is.
2. **A frozen pass.** The candidate mutated three times during this review (second consecutive
   round with a mid-review mutation; round-04 and round-05 directories already exist while round 3
   was still collecting). Whatever is certified needs to hold still for one full review pass; my
   S4 evidence is pinned to sha256 `e6aa4c7d…` and should make that pass short if the tree freezes
   there.
3. **Non-blocking, for the record:** pin tree tests for the two cycle-4 safety branches (probe 7/8
   behaviors); the `statRegular` hardening for the per-marker stat swallow (probe 3); and the
   scope dispute over codex-1's round-4 CRITICAL (idea-scoped, serialized, out-of-repo anchor)
   belongs to the round now in progress — the deferral to `fixup-budget-trust-anchor` is honestly
   disclosed, and I do not independently object to it, but §15.3 means my not-objecting is not a
   sign-off on that scope call.

## Independent validation record

All commands run by me, this session, in the trees (read-only) or in isolated copies (mutations).
Tree state never modified by me; both `git status` file-sets identical before/after. Skill repo
untouched since 15:23 and fully green. Hashes and timeline as pinned above. The isolated copies
(`mktemp -d`, module + go.mod/go.sum + embedded protocol defaults) were deleted after the checks.
