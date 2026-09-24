---
agent: claude-1
idea: meta-protocol-change-lean-organizer
review-round: 4
date: 2026-09-24
reviewed-commit: 4df0855
responding-to: [claude-1/review/round-03, kimi-1/review/round-03]
---

## Summary

All seven signed fix-up-cycle-3 items (H1–H7) landed as authorized, and I verified each
one independently in my own trees: H1's content signal is real code with a real
regression test, and my round-03 falsifier — the condition I said would withdraw
R3-MIN-1 — **now fires** on a fresh `git archive` extraction at this commit. Both
amendments were carried as committed. Both suites are green under my own hands at the
reviewed commits (CLI 31/31 ok, exit 0; skill 399 pass / 0 fail, exit 0), FINAL is still
frozen, and nothing has been released, tagged, published, merged, or installed. All seven
of my round-03 findings are **withdrawn as fixed**; no CRITICAL and no MAJOR this round,
and no unmet frozen acceptance criterion.

I file three new items — one MINOR and two NIT — and **none of them asks for a fix-up
cycle 4.** The MINOR is a semantic gap in `wait --for review` that H1 did not cause but
did make visible (the digest and the terminal line in the same envelope now disagree); it
belongs in a deferred follow-up. My position on closure has therefore moved from NOT YET
to **READY**, conditional on the three items being dispositioned in writing and on the
LE-7 goal-done verdict landing at the closing HEAD — which is not mine to produce.

## Protocol context attestation (this review)

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase6-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"}
```

Cross-check (PRIMARY), my own hands this session:

    $ shasum -a 256 .parley-runtime/protocol-packets/full-phase6-deliberation-8ce83c….md
    8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7   (109,928 B)
    $ shasum -a 256 parley-deck/COOPERATION.md
    8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7

The attested packet IS the live deck authority, byte for byte — `full` mode, no fallback.
Read in full before any check: frozen `FINAL.md`, `IMPLEMENTATION.md` at `4df0855`
(87,818 B, all of fix-up cycle 3), both round-03 review files, the entire signed
`review/consensus.md` cycle 3 including all three signoff blocks with both amendments and
kimi-1's §15.1 SELF-CORRECTION, `00-prompt.md`, and the organizer's goal-done commission
note.

## Reviewed trees, isolation, and tooling provenance

Paired, disposable, mine alone; the live worktrees were read but never built in or tested
in. All preserved as instructed.

| tree | what | state |
|---|---|---|
| `/tmp/c1-r4/cli` | `git clone --no-hardlinks` of the live CLI worktree, checked out at `4df0855` | `git status --porcelain` empty |
| `/tmp/c1-r4/skill` | `git clone --no-hardlinks` of the skill worktree at `b06a65a` | clean; `npm ci` for devDependencies |
| `/tmp/c1-r4/fresh` | `git archive 4df0855 \| tar -x` — the fresh-checkout shape, equal mtimes | not a git tree, by design |
| `/tmp/c1-r4/fresh-r4dir` | copy of the above + an empty `review/round-04/` | control for the F10 directory path |
| `/tmp/c1-r4/falsify` | scratch copy of the clone for revert-and-rerun experiments; restored to HEAD between probes | probes removed; `git status` clean |

Binaries built by me, `go1.27.1 darwin/arm64`:

    $ go build -o /tmp/c1-r4/parley ./cmd/parley        # at 4df0855
    85b0f34c4be8df6db29d567b358a40f5d336463d620f15f343f292ec14ef0cc2
    $ go build -o /tmp/c1-r4/parley-preH1 ./cmd/parley  # 4df0855 with H1's OR-term removed

Both binaries are task-local; neither was installed. Nothing in the live deck was written
by me except this file and its `review/round-04/` directory.

## Refutation attempts

Per-item, against the signed cycle-3 plan and then against frozen FINAL A–D. I assume each
fix is wrong until a check refuses to show it.

### H1 — the content signal (claude-1 R3-MIN-1, VC-8) — ATTACKED FOUR WAYS; HOLDS

**(a) The end-to-end falsifier I offered in round 03 and at signoff. It fires.** My stated
withdrawal condition was: *"a fresh clone/checkout/archive in that state reading `await
review artifact` withdraws R3-MIN-1"*. At this commit (PRIMARY):

    $ git archive 4df0855 | tar -x -C /tmp/c1-r4/fresh
    $ stat -f '%Sm  %N' …/IMPLEMENTATION.md …/review/round-03/claude-1.md …/review/round-03/kimi-1.md
    Sep 24 10:15:38 2026  IMPLEMENTATION.md
    Sep 24 10:15:38 2026  review/round-03/claude-1.md      # ONE mtime
    Sep 24 10:15:38 2026  review/round-03/kimi-1.md        # ONE mtime
    $ ls …/review/            # round-01 round-02 round-03 — NO round-04 directory
    $ /tmp/c1-r4/parley organizer brief --idea meta-protocol-change-lean-organizer
    ## Next action (fixed enumeration)
    await review artifact                    # CORRECT — was `await implementation` pre-H1

**(b) Before/after on the identical tree, to prove the fix and not the fixture.** Same
extraction, my pre-H1 binary:

    $ /tmp/c1-r4/parley-preH1 organizer brief --idea meta-protocol-change-lean-organizer
    ## Next action (fixed enumeration)
    await implementation                     # the wrong-direction failure, reproduced

**(c) The shipped regression test genuinely catches it.** I reverted only H1's OR-term
(`fixUpAwaitingReviewRound(d.Implementation, d.Review)`) in `/tmp/c1-r4/falsify` and ran
the two tests the record names:

    $ go test ./internal/app -run 'TestPhaseDigestNextActionFreshCheckoutEqualMtimeAwaitsReview|TestPhaseDigestNextActionFixUpPublishedAwaitsReview' -count=1 -v
    wait_test.go:555: same-round newer review artifact must not mask the published cycle: want "await review artifact", got "await implementation"
    --- FAIL: TestPhaseDigestNextActionFixUpPublishedAwaitsReview
    wait_test.go:625: fresh-checkout equal-mtime fix-up-published state must read "await review artifact", got "await implementation"
    --- FAIL: TestPhaseDigestNextActionFreshCheckoutEqualMtimeAwaitsReview

Both new assertions fail pre-H1 with R3-MIN-1's exact output and pass at HEAD. The
implementer's pre/post claim is confirmed independently, by my own revert.

**(d) I tried to break the `M ≤ N` predicate itself and could not.** My worry was that the
content signal assumes the round reviewing fix-up cycle N is `round-0(N+1)`, so a deck with
more review rounds than fix-up cycles would fall through wrongly. That worry is unfounded,
and the reason is structural, not conventional: a fix-up cycle N can only be authorized by
review consensus cycle N, which can only follow review round N. So after cycle N is
published and before its review lands, the latest complete round is exactly `round-0N`
(`M = N`); once that cycle is reviewed, `M ≥ N+1`. `M ≤ N` is therefore equivalent to "no
complete round has reviewed cycle N", which is what the doc comment claims. **This
refutation failed; I record it as a strengthening of H1, not a finding.** The added term is
also a pure OR — it can only turn `await implementation` into `await review artifact`,
never the reverse, so it carries no regression surface for the states that already worked.

**Amendment 1's locators, re-verified at the reviewed commit (PRIMARY):**
`implementationNewerThanLatestReview` spans `internal/driver/phasedigest.go:182-200`; the
strict comparison `return implInfo.ModTime().After(newest)` is at `:199`; the new helpers
are at `:206` (`fixUpCycleNumberFromStatus`) and `:227` (`fixUpAwaitingReviewRound`); the
wired arm is at `:279`. Exactly as the cycle-3 record states, and the mtime branch is
unchanged — the arrival signal was kept alongside, not replaced.

**H1's mandatory record correction** (`IMPLEMENTATION.md:828-840`) is there and is honest:
it owns that the cycle-2 "retired the residual" wording was unconditional for
mtime-conditional behavior, names the operations that produce equal mtimes, credits all
three reproductions, says the content signal restores the claim unconditionally, and
explicitly refuses to claim the retirement predates cycle 3. This is the wording I
conditioned my signature on. Discharged.

### H2 — the cycle-2 deviations subsection (R3-MIN-2) — HOLDS

`### Deviations from agreed fixes (cycle 2 — subsection added in fix-up cycle 3, H2)` is at
`IMPLEMENTATION.md:450`, marked as retroactive, with one bullet for the G1 divergence
(52,295 B delivered vs "≈ 42 KB with §2", +10.2 KB / +24 %), the derivation reason, the
FINAL-frozen note, and `None` for G2–G10. I checked both of its quotations rather than
trusting them (PRIMARY): the signed plan text is at `review/consensus-cycle-02.md:121`
("the **retained whole-section total** (≈ 42 KB with §2, recomputed at the fix-up-2…") and
FINAL's C.3/R-2 row at `FINAL.md:436` ("recorded against **both** the measured floor
(≈ 42 KB with §2) and the ≤ 70,000 B guardrail in the same test run"). Both verbatim.
This matches my open-question-1 answer: the agreed-fixes bullet is the right home, and
`## Deviations from FINAL.md` stays for mechanism changes.

### H3 — close conditions (R3-MIN-3) — (a) HOLDS; (b) CORRECTLY NOT PERFORMED HERE

(a) The pointer is at `IMPLEMENTATION.md:269-283`, and it points at the signed consensus's
complete rule set rather than restating it. I re-read the code it summarizes in my own tree
(PRIMARY, `internal/driver/impl.go`): reservations triage `:283-286`, reviewer floor
`:287-289`, LE-7 goal-done gate `:294-298`. `00-prompt.md:8` is `auto_implement: true`,
`strict_gate` absent. The record's closing sentence — *"Until the goal-done verdict lands,
nothing in this idea may claim the close conditions satisfied — and this record does
not"* — is accurate: I found no claim to the contrary anywhere in the cycle-3 section.

(b) The implementer did not run its own goal check, and says so. The organizer has
commissioned it separately: `inbox/codex-1-to-kimi-1_meta-protocol-change-lean-organizer_goal-done.md`
names kimi-1 (a non-implementer, non-facilitator), pins CLI `4df0855` / skill `b06a65a` /
FINAL `120a9bf`, requires the checker's own evidence and a decisive PASS / FAIL /
INCONCLUSIVE, and forbids a manufactured pass. That artifact is kimi-1's alone; **I have
not written it and will not.** One condition the closing consensus must check rather than
assume: the commission pins `4df0855`, so it is a check at the closing HEAD **only if no
fix-up cycle 4 lands**. If any item below is dispositioned as a fix instead of a deferral,
that verdict is staled and must be re-run at the new HEAD.

I also re-verified the `review/goal-done/` inertness both reviewers asserted at signoff
(PRIMARY, code read at `4df0855`): `latestRoundSection` filters candidates on
`strings.HasPrefix(e.Name(), "round-")` (`phasedigest.go:305`), so a `goal-done/`
directory cannot be enumerated as a review round or perturb the digest, the brief, or
`wait --for review`.

### H4 — the 23,230 ↔ 23,223 reconciliation (R3-NIT-1) — HOLDS, RE-MEASURED

I did not take my own round-03 figure on trust; I re-measured it a third time, through the
shipped parser, in `/tmp/c1-r4/falsify` (PRIMARY):

    ## 1. Scope and purpose                                      3130 B   4 blocks
    ## 3. Directory layout                                       2337 B   1 blocks
    ## 8. Inbox (lightweight channel)                            1531 B   2 blocks
    ## 10. TL;DR                                                 1494 B   1 blocks
    ## 12. Pipeline blocks & action stages                       7544 B  13 blocks
    ## 13. Retrospective optimization                            5173 B   6 blocks
    ## Appendix A — Adopting this protocol in a new project      2021 B   1 blocks
    SEVEN TOP-LEVEL TOTAL = 23230 B over 28 blocks
    NINE-ENTRY TOTAL      = 30262 B

23,230 B exactly, over 28 blocks exactly — so the corrected clause at `:350-353` is right
on both halves: the delta is 7 (one per section), and a per-block newline would have added
28. The nine-entry 30,262 B in the cycle-2 record reproduces too. My §15.1 SELF-CORRECTION
from round 03 (23,223 → 23,230) stands and is now thrice-measured.

### H5 — the crossed `floor` labels (R3-NIT-2) — HOLDS, RE-MEASURED

Under the old code's own `len(b.Text)` heading-block convention (PRIMARY):

    floor (9 heading blocks, len(Text)) = 15759 B     (record: 15,759 B)
    section 2 heading block             =  3977 B     (record: 3,977 B)
    floor + section2                    = 19736 B     (record: 19,736 B)

The corrected clause at `:354-357` attributes both numbers correctly. I also confirmed the
measurement basis is legitimate: `git diff --stat 118b245..4df0855 -- parley-deck/COOPERATION.md`
is **empty**, so the protocol text is byte-identical between the commit the figures describe
and the one I measured on.

### H6 — the `runs/` residual (R3-NIT-3, Amendment 2) — HOLDS, AND IS DURABLE

The clause at `:597-606` carries the shape-based replacement both reviewers endorsed at
signoff, verbatim. Its facts at this HEAD (PRIMARY):

    20260923T202501.377412000Z  run.json:2278   events:1543
    20260923T211817.507433000Z  run.json:NONE   events:1616
    20260924T003301.850305000Z  run.json:NONE   events:1636
    20260924T025208.421874000Z  run.json:NONE   events:1636
    20260924T044128.336773000Z  run.json:NONE   events:1636
    $ created_at 2026-09-23T20:25:01.377412Z / updated_at 2026-09-23T20:25:01.377412Z  → equal

Four of five recent records write `events.jsonl` only; the single manifest is zero-width;
post-F7 records exist. The wording is falsified only by a real driver transition, which is
what I asked for — it did not rot between the signoff and this review, and it cannot rot at
the next signoff launch.

### H7 — the asserted relaxation (R3-NIT-4) — HOLDS, VIA THE STRONGER ALTERNATIVE

`TestPhaseDigestNextActionFixUpPublishedAwaitsReview` no longer ends in `t.Logf` +
enumeration validity. `wait_test.go:554-557` now asserts the same-round-newer state stays
`await review artifact` (the content signal, M=1 ≤ N=1), and `:575-578` asserts the genuine
relaxation — round-02 complete with newer artifacts, M=2 > N=1 — reads `await
implementation`. The original G4 record wording at `IMPLEMENTATION.md:399` is now backed by
real assertions, so keeping it is correct rather than convenient. The implementer took the
stronger option kimi-1 mildly preferred and I offered; I confirm it is stronger, because
both halves now fail pre-H1 (shown above) instead of one.

### Both suites, at the reviewed commits, in my own trees

**CLI** — clean clone at `4df0855`, the G3-named command:

    $ go test ./... -count=1 -timeout 2400s
    ok lines: 31   FAIL lines: 0   EXIT=0
    internal/trajectory 639.558s    internal/app 562.081s
    $ go vet ./...   → exit 0, no output

`internal/trajectory` again exceeds Go's 600 s per-package default on my machine — 639.6 s
here, 644.7 s in my round 03, 609.4 s for kimi-1, against the implementer's 577.8 s. Four
measurements straddling the cliff on three machines is stronger evidence for G3 than any
single run: the explicit `-timeout` is load-bearing and machine-dependent, not decorative.

**Skill** — clean clone at `b06a65a`:

    $ npm ci && npm test
    ℹ tests 399   ℹ pass 399   ℹ fail 0
    python 3.14: 54 tests OK across 7 files
    EXIT=0

Counts identical to the record's (399 pass / 0 fail; 54 python tests across 7 files).
**Methodological disclosure, because it could be misread as a regression:** my first run in
this clone reported 385 pass / **1 fail** / EXIT=1. The failure was
`test/design-addons.test.js` → `Error: Cannot find module 'commonmark'` — `commonmark` is a
declared devDependency (`package.json:70`) and a fresh `git clone` carries no
`node_modules`. After `npm ci` the suite is green. The record's claim was never in doubt;
my tree was. I report it because an unexplained 1-fail in a reviewer's log is exactly the
kind of thing that should not sit unexplained in an audit trail.

### Frozen FINAL A–D at this HEAD — regression sweep

- **Frozen:** `git diff 120a9bf..4df0855 -- …/FINAL.md` **empty**.
- **A:** the goal-check contract refuses a non-independent checker in code —
  `driver_impl.go:414-416` returns *"goal-check has no independent checker"* when
  `checker == o.implementer`, and `o.roleErr` short-circuits ahead of it; the
  role-ineligibility and run-plan tests passed inside my full-suite run.
- **B:** exit map live — `--for review` at the fresh tree exit 0; the same invocation with
  an empty `review/round-04/` present exits **3** naming `claude-1 (review artifact),
  kimi-1 (review artifact)`. `--json` stdout carries only `{notes, digest}` and decodes
  through `json.load`; the terminal line is on **stderr** (G2's shape). Two consecutive
  `--json` runs over an unchanged tree are `cmp`-identical (determinism). `notes` correctly
  surfaces the pre-existing unanswered to-user escalation (G5).
- **C:** `parley protocol packet check` → `ok`, exit 0, 69 blocks. `organizer brief` =
  **2,187 B** ≤ 8,192 B, byte-identical ×2, and wrote **nothing** into `parley-deck/`
  (no file under the deck postdates the extraction instant). Skill core `SKILL.md` =
  **17,802 B** ≤ 20,000 B.
- **D:** untouched by this cycle — `git diff --stat 913f8ba..4df0855 -- internal/telemetry
  internal/usage` is empty; the cycle's only source commit touches two files. The live
  attribution-window residual stands, now with H6's corrected facts.
- **Cross-cutting:** deck `COOPERATION.md` = the attested packet (no protocol-text edit →
  no restage needed, none performed); staged core `~/.parley/staging/COOPERATION-2.13.0.md`
  sha256 `fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f`, 109,772 B —
  unchanged, and the skill's bundled reference copy is byte-identical to it;
  `~/.parley/protocol/core/` still holds **only `2.10.0`** (no publish); `VERSION` =
  **1.48.0**; **no tag points at HEAD**. `git show --stat` confirms scope discipline:
  `c3baf09` touches `internal/driver/phasedigest.go` + `internal/app/wait_test.go` only;
  `4df0855` touches `IMPLEMENTATION.md` only. No reviewer file, signoff, FINAL, organizer
  artifact, or inbox note is in either commit.

### Break attempts that failed

Find a state where a published fix-up cycle still reads `await implementation` with the
next round unreviewed (found one — see R4-NIT-1 — but only under a frontmatter shape no
deck in this repo uses); break `M ≤ N` with a plausible round/cycle numbering (structurally
impossible, above); find an H-item whose recorded number does not reproduce (all five
numeric claims reproduced exactly); find FINAL edited after `120a9bf`; find a release,
tag, publish, install, or merge artifact; find a cycle-3 claim of close-condition
satisfaction; find `review/goal-done/` perturbing the digest; find a stream leak onto
`--json` stdout. None succeeded except where noted.

## Findings

### [MINOR] `wait --for review` is satisfied by the stale previous round, so it cannot gate the round the digest says is owed

`boundaryReached` (`internal/app/wait.go:481-513`) defines the review boundary as
`reviewDone := d.Review != nil && d.Review.Completed == d.Review.Total && d.Review.Total > 0`
— the *latest existing* round directory being complete. It has no notion of whether that
round reviewed the currently published fix-up cycle. Post-H1 the two halves of one `--json`
envelope therefore contradict each other in a state this deck reaches every cycle
(PRIMARY, my fresh extraction at `4df0855`):

    $ /tmp/c1-r4/parley wait --idea … --for review --timeout 2s --json
    exit=0
    stdout: {"notes":[…],"digest":{… "next":"await review artifact" …}}
    stderr: wait: boundary reached (review round complete)

`next` says a review artifact is awaited; the terminal line and exit code say the boundary
is reached. Both statements come from the same process, over the same tree, in the same
second.

Why it matters operationally rather than cosmetically: the signed cycle-3 consensus's own
dispatch plan (H3(b), step 4) instructs the organizer to *"gate arrival with the leaf,
read-only `parley wait` … e.g. `--for review` for the round-04 files"*. Run verbatim at
this commit, that returns exit 0 instantly on round-03 — a false boundary, in the one
mechanic (B: "block, do not poll") that this idea exists to provide. The gap closes as soon
as the next round's directory exists; my control (PRIMARY) shows the correct behavior
returns the moment it does:

    # same tree + an empty review/round-04/
    exit=3   wait: timeout after 2s; outstanding: claude-1 (review artifact), kimi-1 (review artifact)

Scope and severity, honestly bounded. This is **not** a regression from H1 — pre-H1 the
digest simply agreed with `wait` by being wrong too. It violates **no frozen acceptance
criterion**: FINAL B.3 specifies `--for review` and `exit 0 boundary reached` without
defining the boundary relative to the published cycle, so the command matches its
specification and misses the operator's need. And it cannot mis-advance a phase:
`boundaryReached` is referenced only by `internal/app/wait.go` (PRIMARY, repo-wide grep) —
the driver's phase machine does not use it, so `wait` remains a pure observer per FINAL
B.5. That is why this is MINOR and not MAJOR.

**Suggested disposition: a deferred follow-up, not a cycle-4 fix.** Either teach
`boundaryReached`'s review scope the same content signal H1 added (a published cycle N with
latest complete `round-0M`, `M ≤ N`, is not a reached review boundary), or — cheaper and
sufficient — document in the skill's `wait` section that the organizer creates
`review/round-0(N+1)/` when dispatching, which is what the driver already does and what
makes the F10 directory path fire. I hold no preference between them and do not condition
my closing signoff on either.

### [NIT] A quoted `status:` defeats H1's content signal, because `implSection` does not strip the quotes its own validity gate does

`implSection` (`phasedigest.go:316`) sets `Status: strings.TrimSpace(meta["status"])`, while
`protocol.ValidImplementationStatus` (`reviewartifact.go:106-109`) trims `"` and `'` before
matching. So a quoted status is **ready-for-review but unparseable**: `ReadyForReview` is
true, and `fixUpCycleNumberFromStatus` returns 0 because `strings.TrimPrefix(status,
"fix-up-cycle-")` cannot match past a leading quote. Reproduced (PRIMARY) with a probe that
is byte-identical to the shipped H1 regression test except for the quotes:

    status="\"fix-up-cycle-2\"" ready_for_review=true ValidImplementationStatus=true
    review.Label="round-02" completed=2/2
    next="await implementation"        # R3-MIN-1's failure, surviving in a narrower shape

The same normalization mismatch predates H1 on the sibling branch: a quoted
`status: "complete"` misses the exact-match `d.Implementation.Status == "complete"` arm and
reads `await review artifact` instead of `FINAL published` (probe, PRIMARY: unquoted →
`FINAL published`; quoted → `await review artifact`). `ReadFrontmatter`
(`workspace.go:390-394`) strips whitespace only, so nothing upstream normalizes for them.

Why NIT and not MINOR: **zero** IMPLEMENTATION.md files in this repo's 80-plus idea corpus
use a quoted `status:` (`grep -rl --include='*.md' -E '^status: *["'"'"']' parley-deck/`
→ 0 files), and the quote tolerance in `ValidImplementationStatus` predates this idea
(commit `815c93a`, an earlier audit fix-up), so this idea consumes a pre-existing
accommodation rather than creating one. It is real, objective and code-grounded, but it has
no observed instance and no acceptance criterion rests on it. **Suggested disposition:**
fold into the same deferred follow-up as R4-MIN-1 — one line, `strings.Trim(strings.TrimSpace(…),
"\"'")` in `implSection`, fixes both arms at the root. Not worth a fix-up cycle on its own.

### [NIT] Typo in the H2 subsection the same cycle added

`IMPLEMENTATION.md:453`: *"cycle 2 landed without one — cycle 1 carries **its** at the
equivalent position"* — "its" should be "one" or "its own". The underlying fact is correct
(cycle 1's `### Deviations from agreed fixes` is at `:815`). Purely textual; I record it
only because this idea's own discipline has corrected record defects in every cycle, and
because a closing consensus should decide knowingly whether to freeze it rather than
discover it later. **Suggested disposition: dismiss as immaterial and close over it** —
see my closure position below; I do not ask for a cycle to fix a typo.

## Open questions

1. **Is the goal-done commission pinned to the right tree?** The inbox note names `4df0855`.
   That is the closing HEAD **iff** round 4 produces no fix-up cycle. If the consensus
   disposes any item above as a fix, the verdict must be re-commissioned at the new HEAD —
   LE-7 is explicit that a stale code tree cannot close an implementation. I flag the
   dependency; the organizer and the closing consensus own the call.
2. **Does `head_commit` carrying prose matter to anyone?** The frontmatter reads
   `head-commit: c3baf09 (fix-up cycle 3 source; recorded in the cycle-3 section)`, and the
   digest surfaces that string verbatim as `head_commit`, where the Phase-5/8 prompt
   template says `<sha-or-short-sha>` (`internal/runner/phase58.go:252`). I checked and no
   code parses the field (`HeadCommit` is set at `phasedigest.go:320` and only displayed),
   so I record this as an **observation, not a finding** — the parenthetical is honest and
   more informative than a bare sha, and the record commit cannot name itself. If the
   closing consensus prefers the template shape, that is a judgment call, not a defect.

## Position changes since prior review round

- **R3-MIN-1 → withdrawn as fixed, by my own falsifier.** In round 03 and again at signoff I
  wrote that a fresh clone/checkout/archive reading `await review artifact` withdraws the
  finding, and that I had run it twice without obtaining it. At `4df0855` I obtained it.
  The finding is discharged on the evidence I myself set, not on the implementer's report.
- **VC-7: NOT YET → READY.** This is a change of position and I state it as one. My NOT-YET
  rested on seven live findings with no post-close correction channel. All seven are fixed;
  what remains is one MINOR that belongs in a follow-up idea and two NITs I am content to
  see dispositioned. My cycle-3 forward commitment said a comparable crop of record findings
  on unchanged ground at round 04 would read as churn and I would support closing over
  dispositions rather than opening another cycle. Three items, of which the only
  non-trivial one needs a *new slug* rather than an edit to this idea, is not a comparable
  crop — and I am honoring the commitment in the direction it binds me: **I do not ask for a
  fix-up cycle 4.** Route 2 (close over dispositions) was unavailable to me in cycle 3
  because it required a record edit; it is available now, because nothing I found needs one.
- **VC-8: resolved, and I record it as closed.** Both reviewers converged at MINOR-fix; the
  fix landed in its endorsed shape; the conflict has no residue.
- **My round-03 open questions 1–3 are all answered and closed** (H2's placement, H3(b)'s
  dispatch shape, H1's content-signal shape). No round-03 open question carries forward.

## Responses to other reviewers

### @kimi-1

**Your round-03 zero-findings review held up under my re-check, and your signoff
self-correction was the right call.** I re-ran the substance of your regression sweep in my
own trees and reproduced your figures where they overlap mine: both suites green, brief
byte-identical and under bound (you measured 2,516 B on your tree, I measure 2,187 B on
mine — the brief embeds run-local state, so different trees legitimately differ; both are
well inside 8,192 B), skill core 17,802 B, `packet check` ok, staged core and `VERSION`
untouched. Your `internal/trajectory` 609.4 s now has a third companion over the 600 s line
(mine: 639.6 s).

**On the one place we differed, you have my agreement — and a caveat.** In round 03 you
recorded that the live `--for review` exit-0 observation "predates the creation of the empty
`review/round-03/` directory; at the current live state the same invocation exits 3 naming
both reviewers (correct behavior, explained by deck state, not a record error)". On the
question you were answering — whether the *record* was wrong — you were right, and I agree:
it was deck state, not a record error. My R4-MIN-1 asks a different question you were not
asked: whether the *boundary semantics* are right for the organizer's dispatch. I reproduced
both halves of your state-dependence as a controlled pair (no round-04 dir → exit 0 on the
stale round; empty round-04 dir → exit 3 naming us both), and my concern is narrowly that
H1 has now made the digest disagree with the terminal line in the same envelope. I would
genuinely like your read: if you think "create the round directory when you dispatch" is
the intended contract and the skill should simply say so, I will happily sign a consensus
that disposes it that way — I have marked it deferred-follow-up, not fix-now, precisely
because I think your reading may be the better one.

**On VC-7, we are now on the same side from opposite directions.** You moved READY → NOT
YET at signoff on the strength of claude-1's evidence; I am moving NOT YET → READY on the
strength of the fixes landing. Your three carried reservations (GitHub-hosted runner
wall-clock, live attribution windows, Windows portability) remain exactly as you recorded
them — I verified none of them changed this cycle, and I do not treat any of them as
closure-blocking under the frozen FINAL, which is also your stated position.

**On your separate duty:** the organizer commissioned the LE-7 goal-done check from you by
inbox note at `4df0855`. That artifact is yours; I have not written it and will not, and
nothing in this review should be read as its verdict or as a substitute for it. My only
input is the pinning question in my open question 1.

### @zcode-1 (implementer — recorded, no response owed to me)

All seven items landed as signed, both amendments were carried as committed rather than
argued, and you took the stronger H7 alternative that was left to your choice. The pre/post
demonstration you recorded reproduced exactly under my own revert. Two of the three items I
file are in code you did not touch this cycle; the third is a typo. Per §15.1 I note that
no verdict in this review is yours to issue on your own implementation, and none of the
above relies on your assertions — every figure I quote I re-measured.

## Updated findings

| id | severity | disposition asked of the consensus |
|---|---|---|
| R3-MIN-1 (equal-mtime fall-through) | MINOR | **Withdrawn — fixed by H1**, verified by my own falsifier at `4df0855` |
| R3-MIN-2 (missing cycle-2 deviations) | MINOR | **Withdrawn — fixed by H2** |
| R3-MIN-3 (close conditions unaccounted) | MINOR | **Withdrawn — fixed by H3(a); H3(b) commissioned, verdict pending** |
| R3-NIT-1 (23,230 ↔ 23,223) | NIT | **Withdrawn — fixed by H4**, re-measured |
| R3-NIT-2 (`floor` mislabel) | NIT | **Withdrawn — fixed by H5**, re-measured |
| R3-NIT-3 (stale `runs/` residual) | NIT | **Withdrawn — fixed by H6**, durable wording verified |
| R3-NIT-4 (overstated test claim) | NIT | **Withdrawn — fixed by H7** via the stronger alternative |
| **R4-MIN-1** (`wait --for review` stale boundary) | MINOR | **Deferred follow-up** — proposed slug `wait-boundary-vs-published-fixup` (`TBD`, candidate per LE-10). Not a fix-up-cycle-4 item from me. |
| **R4-NIT-1** (quoted `status:` defeats H1 + the `complete` arm) | NIT | **Deferred follow-up**, same slug — one-line root fix in `implSection` |
| **R4-NIT-2** (typo at `:453`) | NIT | **Dismiss as immaterial**; I consent to closing over it |

Counts: round 1 → 25 (1 CRITICAL); round 2 → 15 (3 MAJOR); round 3 → 7 (0 CRITICAL /
0 MAJOR); **round 4 → 3 from me (0 CRITICAL, 0 MAJOR, 1 MINOR, 2 NIT)**, none on the code
fix-up cycle 3 changed except R4-NIT-1's H1 half, and none asking for another cycle. This
is the converging shape §4's stopping judgment describes, and I say plainly that it is time
to close.

**Closure readiness — my binding position.** **READY** for a zero-fix closing consensus at
`4df0855`, subject to four things I do not own and will check on the artifact:

1. every finding above carries a written disposition (two deferred, one dismissed) —
   a zero-Agreed-fixes consensus is available on that basis and I will sign it;
2. the LE-7 goal-done verdict from kimi-1 lands **at the closing HEAD** and passes; an
   inconclusive, failed, or unavailable checker escalates and blocks the close (LE-7 can
   only withhold a close, never establish one), and my acceptance does not substitute for it;
3. LE-11 is satisfied: all-✅ signoffs, or a recorded operator ruling accepting any 🟡
   (nothing on my side forces a 🟡 today), with the two-independent-reviewer floor met by
   kimi-1 and me;
4. nothing is released, tagged, published, installed, or merged before the organizer's
   authorized release step — verified untouched as of this review.

**Residuals I record honestly as unresolved and non-blocking under the frozen FINAL:**
GitHub-hosted runner wall-clock remains unmeasured; live attribution windows remain
test-proven only (no non-zero-width window yet produced); Windows `wait`/`usage`
portability remains macOS-verified only; DF-1…DF-4 remain inactive candidates with no
quorum staffed; and R4-MIN-1/R4-NIT-1 join that list if dispositioned as I propose. None of
these is an unmet acceptance criterion; all of them are things a reader of this idea should
be able to find stated plainly, which is why they are here.

## Review record

One file written: this one, plus its `review/round-04/` directory. No source, test, FINAL,
IMPLEMENTATION, peer-review, consensus, organizer, inbox, or run artifact was modified; no
commit, merge, tag, install, release, or publish was performed; `parley protocol publish`
was not run; the goal-done artifact was not written. My review trees are preserved at
`/tmp/c1-r4/`. Every measurement above is PRIMARY — run by me, in my own trees, with the
commands quoted — except where explicitly marked as another participant's figure.
