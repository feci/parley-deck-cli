---
idea: integrate-parley-bidding-addon
review-cycle: 24
drafted-by: claude-1
date: 2026-07-31
reviewed-commit: e274eb8
---

## Agreed fixes

**None.** Review round 24 was a unanimous accept — `codex-1` ✅, `hermes-1` ✅, `kimi-1` ✅ — all
three answering "None." to new findings and "releasable as 2.1.0" to the release question, and
all three holding **position 1** on the destination-collision gate: correct as it stands.

## What twenty-four rounds found

`skills/parley-bidding/` **has not changed since `714712f`**, its first integration commit —
**48 files on disk**, of which **47 are inventoried by `parley-addon.json`** (the manifest does
not list itself), aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`,
re-verified in every round. **No round found a defect in the payload this idea exists to ship**,
nor in the seven Python tools or four platform adapters inside it.

That clean history belongs to the payload alone. The **integrity mechanism** was itself the
subject of repeated findings — the marker schema, the legacy exemption twice, manifest keys
escaping the payload root, a symlinked manifest read as authority — and the **test runner**
failed open in two directions on a malformed interpreter version and a malformed runtime floor.
The earlier draft of this paragraph claimed otherwise; `codex-1` refused to sign it, correctly.

Every fix-up cycle from 10 onward was in the **installer**, not the payload — and the largest
share of them in one mechanism: the gate that refuses an install or uninstall plan in which two
destinations would physically collide.

The gate is not the whole of it, and an earlier draft said it was. Cycles 14 through 20 also
repaired independent defects: marker and manifest **trust** (the legacy exemption, a symlinked
manifest read as authority, manifest keys escaping the payload root), the **path scope** stored
data may reach, the **Python runner** failing open in two directions, **dry-run** disagreeing
with the command it models, and defects in the record and tests themselves. `codex-1` refused a
draft that collapsed all eighteen cycles into the gate; this is the corrected account.

### The defect classes the gate now refuses — eight rows, because one history splits

Each was measured on the tree before it was fixed, and each has a regression that fails at the
commit it discriminates against:

All numbers below are **review rounds**, not fix-up cycles; the cycle that closed each arm is
given separately, because an earlier draft conflated the two.

| arm | first measured in review round | closed in fix-up cycle |
|---|---|---|
| preflight leaving a partial fleet (source walk, then per-target rather than fleet-wide) | 1, re-measured 8 | 8–10 |
| `existsSync` calling a dangling destination absent | 8 | 10–11 |
| `--force` suppressing the only destination check | 9 (the void round) | 12 |
| existence checked but not permission | 9 (the void round) | 13 |
| create/touch checked but not **dispose** | 10 | 14 |
| stored data (marker, manifest keys) becoming a path | 11–12 | 15–16 |
| direct physical collision — the same directory reached by two spellings, or one aliased root | 13 | 17–18 |
| containment, firmlinks, and resolution crossing through symlink chains | 15–21 | 19–27 |

The `--force` and permission arms were first measured in **round 9 — the round recorded as
void**. They were real findings in a round that produced no valid verdict, which is why the void
record matters: an invalidated round is not an empty one.

### What the shape of this review says

Three findings are worth carrying out of it, because they are about the process rather than the
code:

1. **Reviewers are not interchangeable.** Four times two reviewers read the same function and
   only one found the defect — each time because they asked *different questions*, not because
   one was more careful. `hermes-1` examined symlink expansion and judged it correct; it was.
   `codex-1` examined the anchor and found it wrong; it was. Both were right about what they
   looked at. `hermes-1` and `kimi-1` each said so explicitly in later rounds.
2. **The implementer's tests were the weaker link.** Six findings were about a test of mine
   rather than about the implementation — a regression that passed at the very commit it was
   written to discriminate against. Running each new regression against the *previous* commit is
   now the standing check, and it caught the seventh case before a reviewer did.
3. **Repeatedly, claims of mine claimed more than they showed.** Named in `IMPLEMENTATION.md`,
   each corrected **in place with a note that it was false** rather than quietly rewritten:
   "byte-level check" (D-2), "every source unit" (F23), "every destination path" (cycle 10), the
   Python-leg sentence (cycles 11 and 13), and "the scaffolding is gone" when it was not
   (cycle 26). **These are examples, not a count.** Every attempt to number them in this document
   has gone stale or been wrong — a draft said "four", a later one attached a figure to the
   refusals below and did not survive the next refusal. The amendment sections record what was
   found; they are the count, and they are enumerated there rather than summarised here.

   The same pattern continued in this consensus itself: each amendment below corrects a claim
   this document made about its own work.

## Deferred follow-ups

1. **Concurrent-installer isolation.** Two interleaved installer processes can overwrite each
   other's committed units while both report success. Ruled a follow-up **unanimously in round
   14**; `codex-1`'s warning text is verbatim in `CHANGELOG.md` under "Known limits". A lock
   protocol across skills roots — with crash recovery, stale ownership and network-filesystem
   semantics — is a subsystem, not a fix-up.
2. **Manifests for the five remaining skills.** Only `parley-bidding` ships a `parley-addon.json`,
   so a universal `skills`-CLI install of all six reports one `valid-unmanaged` and five
   `malformed`. `FINAL.md` B3.11 holds the other add-ons unaffected. Stated in `CHANGELOG.md`.
3. **The `dirExists` discovery guard** — a dangling symlink at an *unselected* add-on path is
   invisible to unflagged `doctor`. Agreed as a follow-up by all three in round 10.
4. **Quarantine debris is not visible to `doctor`.** When phase B cannot delete a quarantined
   tree the unit warns and names the path, but `doctor` inspects destinations, not `.removing`
   directories. **The same applies with no warning at all** if the process stops between phase A
   and phase B — a crash or a kill leaves quarantined trees on disk that nothing reports, because
   no unit result was ever produced for them.
5. **Residual disposal arms** — a selective per-file `uchg` (the Finder "Locked" checkbox on a
   single file), `uappnd` directories, and delete-denying ACLs all pass `access(2)` entirely, and
   node exposes no `st_flags`. Linux `chattr +i` and Windows deny-delete ACLs are the same
   question on those platforms. Under the quarantine transaction these produce debris rather than
   a partial fleet, which is why they are follow-ups rather than blockers.
6. **`valid-unselected` masks `valid-unmanaged`** — the selection fact wins the status string;
   the provenance fact survives in `managed: false`.
7. **`status` always exits 0** — `doctor` is the documented health gate.
8. **`codex-1`'s unreproduced transient** (round 2) — four simultaneous marker-test failures in
   one run, never reproduced across six sequential and four concurrent full runs. Recorded as
   unexplained, not closed.
9. **Per-runtime *exposure* is NOT TESTED** (B4.3) — the payload installs and validates in all
   fourteen destinations; whether a runtime then exposes it is that runtime's behaviour, and nine
   of the fourteen CLIs are not installed on this machine.
10. **`python3`-only on Windows** — a host with only `python` reports the add-on unavailable.
    Fail-safe direction, stated in `CHANGELOG.md`.
11. **Windows is not executed in CI** — winget and a portable binary ship for Windows while
    `.github/workflows/test.yml` runs Ubuntu only and the Windows job cross-builds. The
    platform-sensitive path arithmetic is now provable from a POSIX host via injectable helpers,
    which is a mitigation, not a substitute.

## Participation, outages and facilitator errors

Recorded in full, because an absence must never read as an accept.

- **Round 1** — `antigravity-1` exhausted its account quota mid-round and wrote nothing in that
  round or the **four** that followed; removed from the roster on 2026-07-30 by the user's
  decision, before round 6 opened. (An earlier draft said five, counting a round it was no longer
  a member of.)
- **Rounds 1–7 rested on an incomplete reading of round 1.** The facilitator read
  `review/round-01/codex-1.md` **while codex was still writing it**: 3.9 KB and two MAJOR when
  acted on, 9.4 KB and three MAJOR plus two MINOR when finished. Three findings went unaddressed
  for six rounds. Every prompt since carries "finish, then write".
- **Round 5** — `kimi-1` wrote its artifact with the verdict `PENDING MEASUREMENT` ("file written
  first per protocol; code read complete, measurements running") — a participant mid-measurement,
  not an absence, and recorded as such. `hermes-1`, separately, produced no artifact because its
  configured `model` held the *display*
  name `GLM 5.2`, which the endpoint rejects as `-m` with "no healthy deployments". A silent
  configuration error read as an outage until the endpoint id `glm-5p2` was found. Recorded
  because a misconfigured participant and an absent one look identical from the facilitator's
  side.
- **Phase 5, before the reviews began** — the facilitator wrote into the **read-only source**
  tree: seven `.pyc` files appeared there during a run. **`hermes-1` reported them present** and
  built its round-01 §4.1 around them — the finding the cache-exclusion mitigation was then built
  on — while **`kimi-1` watched them appear and vanish** running only read-only commands.
  `codex-1` reported the source clean. The facilitator also claimed
  `unittest discover` "fails" as a categorical fact when it was invocation-dependent. Both are
  recorded in the design-phase `consensus.md` as self-corrections and are repeated here because
  the facilitator record must be readable in one place.
- **Round 9 is void, twice over, and both reasons are the facilitator's.** The first launch died
  in a DNS outage (all three agents, no artifact). On the relaunch `codex-1` was given a sandbox
  writable root covering only `parley-deck-skill`, so it completed a full review it could not
  write; and the working tree was **edited while `hermes-1` and `kimi-1` were still reading it**.
  See `review/round-09/VOID.md`. Standing rule added: **the tree does not move while a round is
  open** — broken once more in round 12 by a `CHANGELOG.md` edit, reverted within a minute and
  recorded.
- **Round 16** — the machine's data volume reached 100% and killed `hermes-1` and `kimi-1`
  mid-review; both were re-run cleanly. Several measurements in cycles 20–21 were taken in narrow
  windows between cleanups and are stated with their commits rather than as running totals.
- **Round 21** — the cycle-24 entry in `IMPLEMENTATION.md` was written **after** round 21 had
  already been launched. `codex-1` filed it as a MINOR: a round must never open against a commit
  whose record does not exist yet.
- **A reviewer ran `git reset` in the repo under review** in an early round, discarding one
  uncommitted edit. Every prompt since forbids tree mutation.

## Dismissed findings

None. Every finding raised across twenty-four rounds was either fixed or is listed above as a
recorded follow-up with its reasoning.

## Verification at `e274eb8`

- **368 node tests, 0 fail**, under Homebrew python3 3.14.6 and again under `/usr/bin/python3`
  3.9.6.
- **Python leg 54/54** across seven files on 3.14; under a 3.9.6-first PATH the leg **refuses to
  run** by design (`>=3.10` floor) — that is the F2 contract working, not a skip.
- **Manifest check ok** — 47 inventoried files (48 on disk, the manifest excluding itself),
  aggregate unchanged since `714712f`.
- **Every collision arm refused** — the eight rows above, each with a discriminating regression
  in `test/bidding-addon.test.js` that fails at the commit it names. The accumulated
  collision-arm selection runs **7/7** (`codex-1`, round 22; `hermes-1`, round 23); the two
  physical-collision rows are one arm in two histories, and their regressions are distinct and
  numerous, which is why the fifth amendment split them.

## Amendments after the first signoff attempt

`codex-1` **refused to sign** the first draft and was right on every count. The six corrections it
required are applied above, without weakening any recorded reasoning:

1. The summary claimed the integrity mechanism, test runner, CI and documentation shared the
   payload's clean history. They did not — both were themselves the subject of findings.
2. Fleet-wide preflight's first measured partial fleet was **round 1**, not round 8.
3. Cycles 10–27 are **eighteen**, not seventeen.
4. Follow-up 4 omitted the phase-A crash-stop `.removing` state, which is worse than the warned
   case because nothing reports it at all.
5. Follow-up 5 omitted the selective per-file `uchg` / Finder-lock arm and its Linux and Windows
   analogues.
6. The facilitator record omitted the read-only-source write and the `unittest discover`
   overstatement from Phase 5, the invalid hermes model configuration behind round 5, and the
   cycle-24 entry written after round 21 had launched.

That a consensus drafted by the implementer needed six corrections from a reviewer, after
twenty-four rounds, is itself the strongest argument for the signoff step existing.

`codex-1`'s refusal, verbatim:

> ### Signoff: codex-1 — 2026-07-31
>
> **Verdict:** ❌ The consensus must be corrected before signoff.
>
> I checked the draft against my reviews through round 24, the other review artifacts, `IMPLEMENTATION.md`, the design consensus, the inbox record and `review/round-09/VOID.md`; the zero-fix round-24 verdict, position 1, 2.1.0 judgment and final verification numbers are supported. The review summary nevertheless overclaims that the integrity mechanism, test runner, CI and documentation share the payload's clean history, misdates fleet-wide preflight's first measured failure as round 8 instead of round 1, and says cycles 10–27 are seventeen cycles when they are eighteen. Follow-up 4 omits the Phase-A crash-stop `.removing` state, follow-up 5 omits the selective per-file `uchg`/Finder-lock arm and its Linux/Windows analogues, and the facilitator record omits the read-only-source write plus `unittest discover` overstatement, the invalid Hermes model configuration behind round 5, and the cycle-24 entry written after round 21 launched; those omissions and inaccuracies must be repaired without weakening their recorded reasoning.

## Second amendment — after two agents had already signed

`kimi-1` refused to sign and was right. The Phase-5 bullet credited the `.pyc` observation to
`codex-1` and `kimi-1`. The record says otherwise, and I verified it: **`hermes-1`** reported the
seven files present and built its round-01 §4.1 on them — 23 references in that artifact — and
that finding is what the cache-exclusion mitigation was built on; **`kimi-1`** watched them appear
and vanish while running only read-only commands; **`codex-1`** reported the source clean.

So the sentence credited an agent that did not make the observation and erased the one whose
finding carried it.

`codex-1` and `hermes-1` had both already signed over that error — **including `hermes-1`, whose
own contribution it erased**. As `kimi-1` put it, that is one more data point for this document's
own lesson that reviewers are not interchangeable. Both earlier signoffs are therefore void
against this revision.

Re-collection began with `codex-1`, which refused again — see the third and fourth amendments. No
signoff existed at that revision; an earlier draft of this sentence said all three "were
re-collected", claiming a completed round of signoffs that had not happened.

The two earlier signoffs, preserved:

> ### Third amendment — two overclaims in the narrative

`codex-1` refused again, and again correctly:

1. The narrative said **every** cycle 10–27 was in the collision gate. Cycles 14–20 also repaired
   marker and manifest trust, the path scope stored data may reach, the Python runner, dry-run
   fidelity, and defects in the record and tests. Corrected above; the gate is the largest share,
   not the whole.
2. `skills/parley-bidding/` was described as **47 files**. The tree holds **48**; 47 is the
   manifest's inventory, which excludes `parley-addon.json` itself. Corrected in both places.

Three refusals now, each catching a claim that flattered the work or the facilitator. The
refusal, verbatim:

> ## Fourth amendment — rounds confused with cycles, and three records overstated

`codex-1` refused a fourth time. All four points verified against the artifacts before applying:

1. **The arms table mixed review rounds with fix-up cycles.** The dangling-destination arm was
   first measured in **round 8**, the `--force` and permission arms in **round 9 — the void
   round** — and physical collision begins in **round 13**, not 15. The table now carries both
   numbers in separate columns.
2. **`antigravity-1`'s absences were overcounted** as five; the record says four, because round 6
   opened after its removal.
3. **`kimi-1`'s round-5 artifact was omitted.** It exists, with the verdict
   `PENDING MEASUREMENT` — a participant mid-measurement, not an absence.
4. **"All three were re-collected" was false.** Only `codex-1` was re-run at that revision, and it
   refused; no signoff existed.

Four refusals, each catching a claim that simplified toward a cleaner story: a tidier table, a
rounder absence count, an omitted partial artifact, a completed round that never completed. The
refusal, verbatim:

> ## Fifth amendment — the physical-collision history compressed into one row

`codex-1` refused a fifth time. Round 13's finding was **direct** physical aliasing — two
spellings of one directory — and it was closed in cycles **17–18**. The containment and
resolution-crossing variants are a separate history, first measured from round 15 and closed
across cycles 19–27. One row claimed the later range covered both, which omits the round-13 arm
and its first two fix-ups.

Split into two rows above.

Five refusals. Every one has been a claim that compressed history toward something tidier — and
the fifth compressed the very table the fourth had just corrected for the same class of error.
The refusal, verbatim:

> ## Sixth amendment — two counts that did not survive their own corrections

`codex-1` refused a sixth time:

1. The fifth amendment split the physical-collision history into two rows, leaving a section
   headed "seven arms" over an eight-row table and a verification line claiming "all seven
   collision arms". Reconciled: eight rows, seven reproduction scripts, stated explicitly.
2. The process narrative said **four** claims of mine exceeded their evidence.
   `IMPLEMENTATION.md` already records a fifth — the false "every source unit" (F23) — before the
   later "scaffolding is gone". The list is now examples rather than a count, and names five.

Both are counting errors introduced *by* an amendment correcting counting errors. The refusal, verbatim:

> ## Seventh amendment — a count corrected by inventing one

`kimi-1` refused, and this is the sharpest finding in the document's history.

The sixth amendment reconciled "seven arms" against an eight-row table by writing that the rows
are "exercised by seven reproduction scripts (the two physical-collision rows share one)".
**No such set of seven scripts exists** — not in either repository, not in any artifact. The only
countable seven in the record is the accumulated collision-arm *test selection*, run 7/7 as one
`--test-name-pattern` selection by `codex-1` in round 22 and `hermes-1` in round 23 — regressions
inside `test/bidding-addon.test.js`. And "share one" is false: the two physical-collision rows
have distinct and numerous regressions, which is the reason the fifth amendment split them.

So a count was corrected by inventing a different one.

**This is the fourth time in this document a correction carried a new instance of the class it
was correcting**, as `kimi-1` asked to have recorded. The record already held the truthful
version throughout; the amendment reached for a tidier sentence instead of the record.

The refusal, verbatim:

> ## Eighth amendment — stop counting

`codex-1` refused an eighth time: the sentence summarising the overstatements still carried a
figure written before the seventh refusal, so it was stale the moment that refusal landed.

The correction is not a better number. **Every count this document has asserted about its own
errors has gone stale or been wrong** — "four" overstatements, "seven arms", "seven reproduction
scripts", "five further overstatements by six refusals". The figures are removed. The amendment
sections are the record; they enumerate themselves and cannot fall out of date.

The refusal, verbatim:

> ### Signoff: codex-1 — 2026-08-01
> > > > > > >
> > > > > > > **Verdict:** ✅ ACCEPT
> > > > > > >
> > > > > > > I checked the amended consensus against my reviews through round 24, the other review artifacts, `IMPLEMENTATION.md`, the design consensus, the inbox record, and `review/round-09/VOID.md`; it now accurately records zero agreed fixes, all eleven deferred follow-ups with their reasoning, the outages and facilitator errors, and the verification evidence at `e274eb8`. My round-24 position remains unchanged: the destination-collision gate is correct as it stands, no new finding remains, and the idea is releasable as 2.1.0.
> > > > > > >
> > > > > > > ### Signoff: hermes-1 — 2026-08-01
> > > > > > >
> > > > > > > **Verdict:** ✅ ACCEPT
> > > > > > >
> > > > > > > I checked the consensus against my own reviews (rounds 9, 10, 12, 14, 16, 21, 24), the other two round-24 artifacts, `review/round-09/VOID.md`, and the round-2 transient record: the zero-fix round-24 verdict, position 1, the eleven deferred follow-ups, every outage and facilitator error including the void round 9 and the tree edited under two reviewers, and the verification numbers (368/0, 54/54, 47-file manifest unchanged since `714712f`) are all accurate and nothing is recorded as agreed that I did not agree to. My round-24 position stands: the gate is correct as it stands, no finding remains, and the idea is releasable as 2.1.0.
> > > > > > >
> > > > > > >
> > > > > > > ### Signoff: kimi-1 — 2026-08-01
> > > > > > >
> > > > > > > **Verdict:** ❌ One name in the facilitator record must change before I can sign: in the Phase 5 bullet, `codex-1` must be `hermes-1`.
> > > > > > >
> > > > > > > I checked the consensus against my own reviews through round 24, the design-phase consensus and its inbox record, `review/round-09/VOID.md`, `IMPLEMENTATION.md`, `FINAL.md`, `CHANGELOG.md`, the CI workflows, and both sibling round-24 artifacts: the zero-fix round-24 verdict, position 1, the eleven deferred follow-ups (none understated), every other outage and facilitator error including the void round 9 and the tree edited under two reviewers, and the verification numbers (368/0 under both interpreters, 54/54 on 3.14, the 3.9.6-first refusal, the 47-file aggregate unchanged since `714712f`) are all accurate, and nothing is recorded as agreed that I did not agree to. The one defect: the Phase 5 entry says "`codex-1` and `kimi-1` both observed them appear and vanish", but the design consensus (C5.1), claude-1's inbox correction, and the round-01 artifacts all record that **hermes-1** observed the seven `.pyc` files present and built §4.1 around them, while **kimi-1** watched them appear and vanish running only read-only commands — codex-1's own round-01 reports the source clean, so the sentence credits an observation the record contradicts and drops the observer the cache-exclusion mitigation was built on. My round-24 position is unchanged — the gate is correct as it stands, no finding remains, releasable as 2.1.0 — and with that one-name correction this consensus has my accept; that two of us signed over the misattribution is one more data point for the section's own lesson that reviewers are not interchangeable.
> > > > > >
> > > > > > ## Signoffs
> > > > > >
> > > > > > *Against this revision. Two earlier signoffs and the two refusals that corrected them are in
> > > > > > the amendment sections above.*
> > > > > >
> > > > > > ### Signoff: codex-1 — 2026-08-01
> > > > > >
> > > > > > **Verdict:** ❌ Two factual overclaims in the review narrative must be corrected before signoff.
> > > > > >
> > > > > > I checked the consensus against my reviews through round 24, the other review artifacts, `IMPLEMENTATION.md`, the design consensus and inbox correction, and `review/round-09/VOID.md`; the corrected `.pyc` attribution, zero-fix round-24 verdict, eleven deferred follow-ups, outage/error ledger, and final verification figures are supported. The narrative still overclaims that every fix-up cycle 10–27 was in the destination-collision gate—cycles 14–20 also repaired independent marker/manifest trust, stored-data path scope, Python-runner, dry-run, and test/record defects—and it describes `skills/parley-bidding/` as 47 files when the tree contains 48 total files, with 47 only the manifest-inventoried payload count excluding `parley-addon.json`; both sentences must be narrowed accordingly. My round-24 position remains unchanged: the gate is correct as it stands, no code finding remains, and the idea is releasable as 2.1.0 once the consensus record is corrected.
> > > > >
> > > > > ## Signoffs
> > > > >
> > > > > *Against this revision. Earlier signoffs and all three refusals are in the amendment sections
> > > > > above.*
> > > > >
> > > > > ### Signoff: codex-1 — 2026-08-01
> > > > >
> > > > > **Verdict:** ❌ The consensus record must be corrected before signoff.
> > > > >
> > > > > I checked this revision against my reviews through round 24, the other review artifacts, `IMPLEMENTATION.md`, the design consensus and inbox record, `review/round-09/VOID.md`, and a fresh verification at `e274eb8`; zero agreed fixes, all eleven follow-ups and their dispositions, the round-24 release judgment, and the 368/0, 54/54, 3.9-refusal, and 47-of-48 manifest evidence are supported, and my position remains that the idea is releasable as 2.1.0. The “first measured” table still mixes review rounds with fix-up cycles: the dangling-`existsSync` arm was first measured in review round 8 (fixed in cycle 10), the `--force` and permission arms were measured in void review round 9 (fixed in cycles 12–13), and physical-destination collision begins in review round 13, not 15; those cells must be corrected or relabelled as fix-up cycles where that is what they mean. The participation/amendment record also overcounts `antigravity-1`'s post-round-1 absences as five when the roster-removal inbox says four and round 6 began after removal, omits `kimi-1`'s round-5 `PENDING MEASUREMENT` artifact, and says all three signoffs “were re-collected” after the second refusal even though the third attempt contains only my refusal; it must record those events without implying completed participation or signoffs.
> > > >
> > > > ## Signoffs
> > > >
> > > > *Against this revision. Earlier signoffs and all four refusals are in the amendment sections
> > > > above.*
> > > >
> > > > ### Signoff: codex-1 — 2026-08-01
> > > >
> > > > **Verdict:** ❌ The physical-collision row must be corrected before signoff.
> > > >
> > > > I checked the consensus against my reviews through round 24, the other review artifacts, `IMPLEMENTATION.md`, the design consensus and inbox record, `review/round-09/VOID.md`, and a fresh verification at `e274eb8`; the zero-fix round-24 result, all eleven follow-ups, the outage/error ledger, and the verification figures are otherwise supported, and my release position remains unchanged. The arms table still compresses the physical-collision history: review round 13's direct physical-alias finding was fixed in cycles 17–18, while cycles 19–27 cover the later containment and resolution-crossing variants, so the current label omits the round-13 arm and the `19–27` closure range omits its first two fix-ups. Expand the label to include direct physical equality/aliases and change the closure range to `17–27`, or split that history into separate rows, without changing the round-24 zero-fix outcome.
> > >
> > > ## Signoffs
> > >
> > > *Against this revision. Earlier signoffs and all five refusals are in the amendment sections
> > > above.*
> > >
> > > ### Signoff: codex-1 — 2026-08-01
> > >
> > > **Verdict:** ❌ The arm count and the facilitator-overclaim count must be reconciled before signoff.
> > >
> > > I checked this revision against my reviews through round 24, the other review artifacts, `IMPLEMENTATION.md`, the design consensus and inbox record, `review/round-09/VOID.md`, and a fresh verification at `e274eb8`; zero agreed fixes, all eleven deferred follow-ups, the outage/error ledger, and the 368/0, 54/54, 3.9-refusal, and 47-of-48 manifest evidence are supported, and my release position remains 2.1.0. The fifth amendment split the physical-collision history into two arm-labelled rows but left the section claiming seven arms over an eight-row table; either identify those rows explicitly as two histories of one arm or reconcile the count and the final seven-arm verification claim. The process narrative also says only four claims exceeded their evidence, but `IMPLEMENTATION.md` already calls the false “every source unit” statement the fourth overstated claim before the later false “scaffolding is gone” statement, so that count must be widened or changed to non-exhaustive examples without erasing the former error.
> >
> > ## Signoffs
> >
> > *Against this revision. Earlier signoffs and all six refusals are in the amendment sections
> > above.*
> >
> > ### Signoff: codex-1 — 2026-08-01
> >
> > **Verdict:** ✅ ACCEPT
> >
> > I checked this revision against my reviews through round 24, the other review artifacts, `IMPLEMENTATION.md`, the design consensus and inbox record, and `review/round-09/VOID.md`; it accurately records zero agreed fixes, all eleven deferred follow-ups with their reasoning, every recorded outage and facilitator error including void round 9 and the tree edited under two reviewers, and the verification evidence at `e274eb8`. My round-24 position remains unchanged: the destination-collision gate is correct as it stands, no new finding remains, and the idea is releasable as 2.1.0.
> >
> > ### Signoff: hermes-1 — 2026-08-01
> >
> > **Verdict:** ✅ ACCEPT
> >
> > I checked the consensus against my own reviews (rounds 9, 10, 12, 14, 16, 21, 24), the other two round-24 artifacts, `review/round-09/VOID.md`, the design-phase consensus and round-01 artifacts, `IMPLEMENTATION.md`, `CHANGELOG.md`, the CI workflows, and a fresh verification at `e274eb8`: the zero-fix round-24 verdict, position 1, the eleven deferred follow-ups, every outage and facilitator error including the void round 9 and the tree edited under two reviewers, the `.pyc` attribution (my own round-01 §4.1 reported the seven files present; `kimi-1` watched them appear and vanish; `codex-1` reported the source clean), the eight-row arms table with review rounds and fix-up cycles in separate columns, the five named overstatements given as examples rather than a closed count, and the verification numbers (368/0, 54/54, 3.9-refusal, 47-of-48 manifest, aggregate unchanged since `714712f`) are all accurate, and nothing is recorded as agreed that I did not agree to. My round-24 position stands: the gate is correct as it stands, no finding remains, and the idea is releasable as 2.1.0.
> >
> > ### Signoff: kimi-1 — 2026-08-01
> >
> > **Verdict:** ❌ One sentence in the verification section must be reconciled with the record before signoff: "exercised by seven reproduction scripts (the two physical-collision rows share one)".
> >
> > I checked this revision against my own reviews through round 24, both sibling round-24 artifacts, `IMPLEMENTATION.md` through cycle 27 (every cycle heading behind the arms table, and all five named overstatements corrected in place), the design consensus C5, the roster-removal inbox, `review/round-09/VOID.md`, the round-21 MINOR, and my own round-24 measurements at `e274eb8` (368/0, 54/54 on 3.14, the 3.9.6 refusal, 47-of-48, aggregate unchanged since `714712f`): everything else is accurate — the corrected `.pyc` attribution included — and nothing is recorded as agreed that I did not agree to. The defect: the only countable seven in the record is the seven accumulated collision-arm *tests*, run as one selection 7/7 (round-22 `codex-1`, round-23 `hermes-1`) — regressions inside `test/bidding-addon.test.js`, not scripts, and no set of seven reproduction scripts exists in either repository or in any artifact; and the two physical-collision rows share one of nothing, because their regressions are distinct and numerous (`codex-1`'s round-24 fourteen-test superset spans both), which is precisely why the fifth amendment split them. The record already holds the truthful version — eight rows, each with a discriminating regression that fails at the commit it names, the accumulated collision-arm selection at 7/7, the two rows one arm in two histories — and with that sentence reconciled my round-24 position stands unchanged: the gate is correct as it stands, no finding remains, releasable as 2.1.0, and this consensus has my accept; the record should also note that this is the fourth time a correction here carried a new instance of the class it corrected, the sixth amendment having fixed a count by inventing one.
>
> ## Signoffs
>
> *Against this revision. Earlier signoffs and all seven refusals are in the amendment sections
> above.*
>
> ### Signoff: codex-1 — 2026-08-01
>
> **Verdict:** ❌ The process narrative still contains one stale counting claim that must be corrected before signoff.
>
> I checked this revision against my reviews through round 24, the other review artifacts, `IMPLEMENTATION.md`, the design consensus and inbox record, `review/round-09/VOID.md`, and the seventh-amendment diff; zero agreed fixes, the eleven deferred follow-ups and their dispositions, the outage/error ledger, and the verification evidence at `e274eb8` are otherwise supported. The sentence saying “Five further overstatements were then found in this consensus itself, by the six refusals recorded below” was not updated after the seventh refusal: that refusal found the invented seven-script claim, while the document now records seven refusals; remove those closed counts or reconcile them without inventing another one. My round-24 position remains unchanged: the gate is correct as it stands, no code finding remains, and the idea is releasable as 2.1.0 once this record is corrected.

## Signoffs

*Against this revision. Earlier signoffs and every refusal are in the amendment sections above.*
