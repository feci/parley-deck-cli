---
idea: integrate-parley-bidding-addon
review-round: 6
agent: kimi-1
date: 2026-07-30
reviewed-commit: c673111
---

## Verdict

BLOCK

One blocker, and it is not new: the omitted-`skill` ownership hole codex-1 filed this round.
I spotted the asymmetry while reading the predicate before opening his file, then reproduced
it by measurement (probe P2 below) — at `c673111` a readable marker with our package name but
no `skill` property is `valid` + `managed:true` to `doctor` while install and uninstall both
refuse the same directory. That is the round-5 contradiction one field-deletion away, so
`c673111` must not ship as-is.

Process note: the implementer committed `083f4f6` (fix-up cycle 6) while this review was
running. I verified that fix by measurement (suite 312/312, plus my P2 scenario now reports
`malformed` + `managed:false` for core and add-on alike, doctor red). This BLOCK is therefore
already satisfied in the working tree; the next round should review `083f4f6` and can converge
quickly. My one MINOR below reproduces at `083f4f6` too and is the only item still open there.

## Outstanding findings — closed or not

My own round-4/5 suspects (my round-5 review never left draft, so these were all open):

1. **Probe cache key collision (SOH separator can occur inside values) — CLOSED.** The key is
   now `JSON.stringify([workingDir, sortedPairs])` over the whole effective environment plus
   the resolved cwd (lib/installer.js:1471), injective by construction. Measured
   env-sensitivity: same PATH, a stub `python3` echoing `$STUB_VER`, verdicts 3.12-ok then
   3.9-fail in one process — no cache bleed.
2. **`managed` absent on `missing` units — CLOSED.** Measured: a deleted add-on directory
   reports `status:"missing"`, `managed:false`, doctor red.
3. **Ownership agreement for the wrong-skill marker (codex-1 round-5 MAJOR) — CLOSED.**
   Measured: flipping the bidding marker's `skill` to `parley-design` gives
   `malformed` + `managed:false` + the identity problem in `doctor`, `blocked` from install,
   `blocked` from uninstall. One answer.
4. **Read-command `--only`/`--no-addons` treated as the recorded selection (codex-1 round-5
   MAJOR) — CLOSED.** Measured on a healthy full install: unflagged doctor green with 6 units;
   `doctor --only parley-bidding` narrows to `[parley-deck, parley-bidding]`, zero problems,
   green; `doctor --no-addons` narrows to `[parley-deck]`, green; `status --only` narrows
   identically. Residual detection still runs for unflagged reads. One residual interaction
   remains — see the MINOR below.
5. **Probe cwd (codex-1 round-5 MINOR) — CLOSED.** Measured: relative `PATH=bin`, cwd with
   3.12 stub vs cwd with 3.9 stub, same env, both call orders — correct pass/fail/first-verdict
   again. `spawnSync` receives the resolved cwd and the key distinguishes it.
6. **Validation record overstated 3.9.6 (codex-1 round-5 MINOR) — CLOSED.** Re-measured: with
   `/usr/bin/python3` 3.9.6 first the Python runner refuses (`python3 is 3.9, but the add-on
   declares >=3.10`) and runs zero tests; node leg 309/309. On python3 3.14.6: node 309/309,
   Python 54/54, manifest check 47 files at
   `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`. The record now
   says exactly this.
7. **`valid-unselected` (cycle-4 follow-up) — VERIFIED.** Measured: full install then
   `install --force --only parley-design` leaves four residuals reporting `valid-unselected`,
   `managed:true`, `selected:false`, `missing:[]`, `runtime:null`, doctor red; stderr names
   "installed but outside the recorded selection"; `status` carries them and still computes
   compatibility. Re-including restores `valid`. A residual with a deleted payload file is
   `malformed` (payload verdict wins), not `valid-unselected`.
8. **Omitted-`skill` marker (codex-1 round-6 MAJOR) — CONFIRMED, NOT CLOSED at `c673111`.**
   My independent probe: full install, delete only `skill` from the core and bidding markers,
   payload and all other fields untouched. `doctor.ok === true`, both units `valid`,
   `managed:true`, no problems; `install --only parley-bidding` and
   `uninstall --only parley-bidding` both return `blocked` for core and bidding. Health claims
   ownership that both mutations deny — the exact defect class round 5 blocked on, through the
   `!== undefined` exemption at lib/installer.js:1342. No released marker needs the exemption
   (all 29 release tags, v1.0.0 through v2.0.0, wrote `skill` — confirmed from each tag's
   `writeMarker`; and my legacy-shape probe with `skill` present stays green and upgrades
   without `--force`). Fixed at `083f4f6`;
   verified there: both units `malformed`, `managed:false`, "carries no skill identity",
   doctor red.

## New findings

### [MINOR] A filtered read re-labels a residual add-on as `selected:true` and green

**Where:** lib/installer.js:894-903 (`selected: true` by construction for flag-derived units),
lib/installer.js:920-937 (residual detection skipped under an explicit filter)

**What:** After `install --force --only parley-design` on a full install, the unflagged read
correctly reports the residual `parley-bidding` as `selected:false`, `valid-unselected`,
health red. But `doctor --only parley-bidding` on the *same tree* reports that same unit
`selected:true`, `status:"valid"`, zero problems, `doctor.ok === true`. The `selected` field
is defined by its own problem text as "part of the recorded selection" — under a filter it
asserts that recorded fact and gets it wrong, because it is derived from the flag, not from
the core marker.

**Why it matters:** The cycle-5 fix made the flag a filter over what to inspect, which is
right — but it also silenced the one fact the filter did not ask about. A user or script
verifying the bidding opt-out with a scoped probe (`doctor --only parley-bidding --json`)
gets a green `valid` + `selected:true` for a unit the recorded selection excludes and the
unflagged gate fails on. Two read commands give opposite answers about a recorded fact for the
same directory. Narrow trigger: the filter must name exactly the residual unit, and the
documented unflagged health gate is correct — hence MINOR, not MAJOR.

**Evidence:** Measured at `c673111` and reproduced identically at `083f4f6`:
unflagged → `parley-bidding: valid-unselected, selected:false, doctor.ok:false`;
`--only parley-bidding` → `valid, selected:true, managed:true, doctor.ok:true`.

**Fix:** Derive `selected` for flag-requested units from the recorded selection in the core
marker, not from the flag. If a filtered read names a unit that is on disk but absent from the
recorded selection, report it `selected:false` / `valid-unselected` (and let health reflect
it) instead of skipping residual awareness for it. The narrowing itself — not accusing the
*other* recorded add-ons — must stay as cycle 5 fixed it.

## What I verified and found correct

- Reviewed an isolated `git archive` of `c673111` in `/tmp`; the source worktree was never
  mutated. The tree moved to `083f4f6` mid-review; that successor was archived separately for
  the fix verification only.
- Full suite at `c673111`: node **309/309** (python3 3.14.6); Python **54/54** on 3.14.6;
  3.9.6 refusal as designed; manifest check 47 files, sha256 as above. At `083f4f6`: node
  **312/312** (includes the two previously-uncommitted tests and the new no-skill-identity
  regression — that regression fails at `c673111` with `valid` vs `malformed`, as expected).
- Legacy markers *with* `skill` (v2.0.0-shaped core recording all five add-ons, v1.4.0-shaped
  add-ons): all six units `valid`, doctor green, full upgrade without `--force` proceeds. The
  stricter predicate refuses no released marker shape.
- `valid-unmanaged`: universal marker-less copy of all six skills → bidding
  `valid-unmanaged`/`managed:false`, the other five `malformed`, doctor red; install over the
  copies blocked without `--force`. The deferred B3.11 state is unchanged.
- `unmanagedButVerified` remains anchored to the packaged source (round-4 laundering probes
  stay `malformed` — covered by the committed suite, re-run here).
- Text rendering: doctor/status print `valid-unselected` units with their integrity problem;
  no unhandled branch; doctor's stderr names the selection reason separately from
  "missing or malformed".
- `uninstall --only <residual>` does remove the residual add-on — and also the core skill.
  The core-always behavior predates this idea (v2.0.0's `targetSkillUnits` already built the
  unit list core-first), so it is not a new defect; noted as a question below.
- Probes were run as standalone scripts against the archived copies with stub interpreters;
  no state outside temp directories was touched.

## Open questions for the implementer

1. `uninstall --only parley-bidding` on a residual bidding removes the core skill too
   (pre-existing semantics). The `valid-unselected` remediation advice ("remove the directory,
   or re-run install including it") pointedly does not mention `uninstall --only` — is that
   deliberate because it takes the core with it? If so, fine; a doc line might save a user a
   surprise. Not a finding.
2. With `083f4f6` already committed and verified here, should round 7's reviewed commit be
   `083f4f6`? The only open item I have there is the MINOR above.
