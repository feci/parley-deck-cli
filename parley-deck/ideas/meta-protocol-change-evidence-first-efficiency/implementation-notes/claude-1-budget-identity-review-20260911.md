---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
type: partial-source-review
role: source reviewer
reviewed-commit: 3ea8693 (brief-reported; I ran no git command)
responding-to: codex-1-budget-identity-dispositions-20260911.md
supersedes: nothing — claude-1-budget-corrections-review-20260911.md and every earlier
  claude-1 artifact stand unchanged
scope: internal/budget/{lock.go,lock_unix.go,lock_windows.go,ledger.go,review_test.go,ledger_test.go},
  internal/app/{budget.go,budget_test.go,budget_attended_other.go,budget_attended_windows.go,
  consensus_request_signoffs.go,app_test.go,app.go,termios_*.go}, internal/store/events.go,
  internal/fsutil/{sync_darwin.go,replace_unix.go,replace_windows.go}, docs/agent-runtime-configuration.md
signoff: none
---

# Partial source review — budget lock identity, read-only inspection, attended control, signoff evidence

## 0. What this is and is not

Bounded partial source review of one slice. Not a Phase-6 round, not a review-consensus
signoff, not an AC closure, not a signature. It closes no finding by itself.

I had **native file tools only and no shell. I executed nothing** — no test, no build, no
`git`, no runtime observation. I made no source edit, no git mutation, no global-config
change, and invoked no other agent. I touched no other owner's artifact.

Codex's dispositions are proposals. I evaluated each against source, not against the
disposition text, and say concur or disagree with reasons.

## 1. Provenance (§15.2) and test limitations

- **PRIMARY** — every claim about *what the code says* comes from reading the named files at
  the stated locators in this worktree. The negative claim "`internal/budget` still has exactly
  one non-test importer" is from a repo-wide `*.go` search for `internal/budget` → only
  `internal/app/budget.go` and `internal/app/budget_test.go`. Same method for
  `Snapshot.Exposure()` having **no remaining definition or reference** anywhere.
- **Brief-reported, not verified by me** — that this tree is `3ea8693`. I re-anchored every
  locator by reading the current files, so the locators hold for this tree whatever its SHA.
- **SECONDARY (facilitator-executed, explicitly NOT mine)** — budget suite PASS (1.957 s),
  PASS on the shared volume (1.407 s), `-race` PASS (5.089 s), focused app
  Budget/ConsensusRequestSignoffs/Help PASS (10.151 s), `vet` PASS on app/budget/runner, and
  the Windows `app` cross-**build**. The brief itself states these focused results do not
  certify every AC. Where a verdict of mine depends on one, I say so.
- **RECALL (UNVERIFIED)**, tagged where used: (a) POSIX `flock` attaches to the open file
  description, so two `open()`s in one process conflict — this is what makes the two-handle
  probe meaningful; (b) **Windows `LockFileEx` exclusive byte-range locks block read access to
  the locked range through every *other* handle, including other handles in the locking
  process.** (b) is the basis of D-MAJOR-1 and is uncorroborated — nobody has run Windows, and
  a cross-build is not runtime evidence; (c) Go's `syscall.Errno.Is` maps
  `ERROR_ALREADY_EXISTS`/`ERROR_FILE_EXISTS` to `os.IsExist`, which is what makes the
  `MoveFileEx` loser path (`lock_windows.go:32` → `lock.go:248`, `:315`) fall through to a
  re-check instead of erroring. I did not read Go's stdlib or the Win32 docs.

I assert no NOVELTY and no EXEMPTION claim. Nothing below is a general result about filesystem
locking — only what these files do.

## 2. Dispositions re-evaluated

### 2.1 Concur — closed

| Disposition | Verdict | PRIMARY evidence |
| --- | --- | --- |
| **C-MAJOR-2** Inspect not read-only | **Closed. Concur, fully.** | `ledger.go:95-111` takes no lock, creates no directory and never calls `pinLockOrigin`; it only `read`s the atomically published file and then checks scope. Comment corrected at `:92-94`; `docs:268-271` corrected. The regression is stronger than I asked for: `review_test.go:110-145` deletes `lock-origin` from a ledger that already carries a charge, inspects at the right **and** wrong scope, and asserts after each that the directory still holds **exactly one** file named `ledger.json` — so a re-pin, a recreated lock or a stray temp file all fail the test. It then asserts the subsequent write refuses with "has no lock origin" and that the ledger bytes are unchanged. |
| **C-MINOR-1** zero ceiling | **Closed at the level I proposed. Concur.** | `budget.go:72-74` prints a specific notice naming what a zero ceiling means for exposure; `budget_test.go:61-84` asserts both halves of the text and that the observation is not rewritten. It is a notice, not a second confirmation — which is exactly the minimum I offered. |
| **C-MINOR-4** anonymous contention | **Closed. Concur.** | `ErrLockContention` (`ledger.go:29`) is wrapped with the resolved lock path at `lock.go:93`, and only via `interrupted()` when `waited` is true (`:90-96`, `:148`) — so a caller whose context was already dead at entry still gets its own error, undiluted. `review_test.go:162-175` requires **both** `errors.Is(ErrLockContention)` and `errors.Is(context.DeadlineExceeded)`. The pre-existing `ledger_test.go:189` assertion still holds through the wrap. |
| **C-MINOR-5** config doc | **Closed. Concur.** | `docs:118` now names `parley consensus request-signoffs` as the only surface that spawns automatically, before the operator chooses `spawn-tty` rather than after. |
| NIT lossy `Exposure()` | **Closed. Concur.** | The method is gone from `ledger.go`; no definition or caller remains in the repo. Only `ExposureError` survives (`ledger.go:242-261`). |
| NIT missing synopsis entry | **Closed. Concur.** | `app.go:127` (Usage) and `:229-231` (Commands). |
| NIT validation order | **Closed. Concur.** | `budget.go:64-67` (attendance) precedes `:68-71` (reconcile-specific flags). Shared flag parsing still precedes both, which the disposition states. |
| NIT disclosures | **Concur.** | Local cache-path visibility `docs:285-286`; per-entry decision-ID scope `docs:261-263`; byte/ledger bounds `docs:262-263`, `:306-309`; crash orphans and the absence of a sweeper `docs:294-295`. Each is stated as a limitation, not as a fix. |

### 2.2 Concur — closed, with residuals I name below

**C-MAJOR-1 (origin identity can repeat or vanish).** **Concur that both mechanisms I named are
closed, and that a third I had not named was found and closed.**

1. *Recreated cache inode.* `lock.go:82` now opens with `O_RDWR` and **no `O_CREATE`**, and
   `lockIdentity(path, newOrigin)` (`:73`, `:289-295`) refuses to create an identity whenever the
   shared origin already exists. A cache wipe under a live ledger therefore fails closed instead
   of minting a second lock. `review_test.go:33-62` deletes the pinned lock *while it is held*,
   requires the refusal to contain "refusing recreation", requires the file to still be absent,
   and requires the ledger bytes to be byte-identical.
2. *Repeated `hostname + cache path`.* The v2 origin is
   `version\nhost\nlockPath\ntoken\n` (`lock.go:187`) where `token` is a random 256-bit nonce
   published once per permanent local lock (`:296-318`). Two containers sharing a hostname and
   `HOME` now hold different tokens, and the second refuses (`:222-224`). This is the durable-token
   design I suggested, implemented more cleanly than I proposed: the identity file **is** the lock
   file, so there is no second artifact to drift.
3. *The read/open/lock race Codex found during implementation.* `verifyLockIdentity`
   (`lock.go:162-177`) binds the actual held descriptor to both the pinned token and the current
   inode — `Lstat` + `f.Stat()` + `os.SameFile` + `Size()==65` + a positional `ReadAt` content
   compare — and it runs **four** times: after open (`:86`), after acquisition (`:111`), on the
   probe (`:124`), and again after the probe (`:140`). `review_test.go:86-108` injects a `take`
   that rewrites the token **in place, preserving the inode**, proving an inode-only check would
   not have caught it.

I checked the remaining ways to split exclusion inside this design and could not find one that
does not require an operation the docs forbid: an external actor recreating the lock with a
*different* valid token is caught by the origin comparison; deleting the shared origin under a
live ledger is caught by `lock.go:67-71`; deleting the origin *and* the ledger loses no charge
because there are no charges left. A **copied** identity file is indistinguishable by
construction — disclosed at `lock.go:254-257` and `docs:284-285`. Concur as disclosed.

**C-MAJOR-3 (attended recovery platform-gated, message misdirects).** **Concur, closed better than
the minimum I asked for.** `budget_attended_windows.go:7-14` gives Windows a real attendance probe
(`GetConsoleMode` on the input handle, failing closed when `GetStdHandle` errors), and it is a
*separate* function from the protocol-publication policy — `termios_other.go` and
`protocol.go:370-381` are untouched, so the ratified publish gate did not move. `budget.go:16-19`
now separates "unavailable on this platform" from `:64-66` "refusing unattended cost adjustment",
which was the two-branch shape I asked for. Residuals → D-MINOR-4.

**C-MINOR-2 (mismatch cannot say what differed).** **Partially concur.** The parsed reason exists
and is correct code (`lock.go:213-226`), and the absence of a migration/re-pin path is disclosed
without euphemism (`lock.go:68`, `docs:288-295`). But in the cases an operator is most likely to
hit, that diagnostic is never reached → **D-MINOR-1**.

**C-MINOR-3 (wedge preventable).** **Concur** — `Limits.RequireKnownCost` (`ledger.go:49`,
enforced at `:153`, fixtured at `review_test.go:147-160`) is the second option I offered, and the
ordering is right: `ErrReserved` at `:135-137` still dominates, which the test pins at `:157-159`.
The disposition's own caveat ("does not invent an observed provider cost") is accurate. One
consequence of that caveat is not stated anywhere and is not what I would call disclosed →
**D-MINOR-2**.

**C-MINOR-6 (no positive event after a valid append).** **Concur, and the implementation is
stronger than my suggestion.** I proposed calling `consensusContainsSignoff` on the error path;
Codex instead gated the event on the **full shared validator**
(`consensus_request_signoffs.go:186-205`). That matters: `validateRequestedSignoff` (`:622-662`)
requires a clean parse, a strict append-only prefix (`:664-667`), that the agent was *not* already
signed (`:632-634`), that **exactly one** new signoff appeared and it is this agent's (`:635-637`),
that no other agent's block changed (`:642-653`), and a canonical status. Only then is
`agent.signoff.artifact-present-after-failure` appended with the artifact path, its SHA-256 and the
launch mode — and the function still returns the failure (`:204`), with an event-append error
joined rather than swallowed (`:202`). `agent.handoff.completed` remains reachable only through the
interactive poll loop (`:508-517`), so no consumer can read the failure as a completion. The event
name claims *presence*, not causation, which is the honest claim for a before/after file diff.
`app_test.go:1244-1293` pins exit 1, the error text, and exactly one event with the right agent,
status and non-empty hash. Residuals → D-MINOR-3, NITs.

### 2.3 Stated boundary — re-checked, still accurate

No automatic launch/action callers (one non-test importer, and `docs:236-238` says the command's
existence is not evidence that launches are budgeted); no distributed/cloned-writer guarantee
(`lock.go:20-27`, `:254-257`, `docs:284-285`); no Windows runtime validation (`docs:271-274`,
`:293-294`); origin-mismatch/legacy migration explicitly unimplemented (`lock.go:68`,
`docs:288-292`); no human authentication by terminal (`budget.go:23-25`, `docs:266-268`). All five
are stated in both code and docs. The findings below are defects **inside** that boundary; none is
a complaint that the boundary is narrow.

## 3. Refutation attempts (inspection only — I executed nothing)

1. *"A wiped or replaced cache lock cannot silently re-authorize work."* Tried deletion under a
   live hold, replacement with a fresh valid token, and an in-place token swap between read and
   acquisition. **Could not break it** — all three refuse before the ledger is touched.
2. *"Inspection has no side effects."* Traced `Inspect` to `read` only. **Could not break it.**
3. *"The probe proves real exclusion on every supported platform."* **Broken on Windows if
   RECALL (b) holds** — the probe's own identity read overlaps the byte the holder locked. See
   **D-MAJOR-1**.
4. *"A refusal names what actually differs."* **Broken** — a second machine, a moved ledger and a
   v1 deck all report a missing *local cache* lock instead. See **D-MINOR-1**.
5. *"A known conservative reserve keeps bounding exposure."* **Broken by a settle** — a nil actual
   discards it, and an operator ceiling may then land below it. See **D-MINOR-2**.
6. *"Every validated append after a failed process leaves positive evidence."* **Broken for
   BLOCK.** See **D-MINOR-3**.
7. *"A refusal never advises something impossible."* **Broken narrowly** on unsupported-attendance
   targets. See **D-MINOR-4**.
8. *"An operator decision cannot erase a charge or repeat an action."* Re-ran my prior attempts
   against the new code (reconcile→re-reserve, reconcile a known observation, conflicting replay,
   write failure). **Could not break it.**

## 4. New findings

No CRITICAL confirmed. One MAJOR (conditional on an unverified platform claim), five MINOR, seven
NIT.

### [MAJOR] D-MAJOR-1 — the probe's identity re-read overlaps the locked byte, which (RECALL) makes Windows acquisition fail

**Locators:** `internal/budget/lock.go:118-129` (open probe → `verifyLockIdentity(probe,…)` →
`take(probe)`), `lock.go:169` (`f.ReadAt(data[:], 0)`), `internal/budget/lock_windows.go:12`.

`tryLock` on Windows locks exactly byte 0: `LockFileEx(h, EXCLUSIVE|FAIL_IMMEDIATELY, 0, 1, 0,
&Overlapped{})` → `bytesLow=1`, `Overlapped.Offset=0`. `verifyLockIdentity` reads bytes `[0,65)`.
So after `f` holds the lock, the probe's verification at `:124` reads a range `f` has locked
exclusively, **through a different handle**.

Per RECALL (b), an exclusive `LockFileEx` range denies read access through every other handle,
including another handle in the same process. If that recollection is right, `ReadAt` returns
`ERROR_LOCK_VIOLATION`, which is neither nil nor `io.EOF`, so `:170-172` returns it and `:124-129`
aborts the acquisition — meaning **no budget lock is ever acquirable on Windows**, and the
`GetConsoleMode` work that just made `parley budget reconcile` reachable there would be nullified
one layer down. The same root cause makes a *waiting* Windows process fail in
`lockIdentity`'s `read()` (`lock.go:276`) and in `pinLockOrigin`'s `check()` (`:209`) with a raw
lock-violation error instead of waiting and reporting `ErrLockContention`.

This is **unverified**: no one has run Windows, the cross-build proves compilation only, and I read
neither the Win32 documentation nor Go's stdlib. On unix nothing here is wrong — `flock` is
whole-file advisory and does not gate reads, which is why the suite passes on macOS.

**Suggested correction (harmless if I am wrong):** lock a byte range that cannot overlap the
identity, e.g. `Overlapped{Offset: 1 << 20}` with length 1 in both `tryLock` and `unlock`
(`lock_windows.go:12`, `:19`). Unix is unaffected (whole-file flock), the two-handle probe keeps
its exact meaning, and every identity read stays legal on both platforms. A cheaper partial
mitigation — dropping the *content* compare from the probe check and keeping `Lstat`+`SameFile` —
fixes acquisition but not the contender's read, so I would not prefer it. Either way this is the
single highest-value candidate for the first real Windows run.

### [MINOR] D-MINOR-1 — the new mismatch diagnostics are unreachable in the cases operators actually hit

**Locators:** `internal/budget/lock.go:61-79` (ordering), `:293-295` (the message that fires
instead), `:213-226` (the diagnostics that do not fire). **Doc at stake:** `docs:289`.

`lockIdentity` runs **before** `pinLockOrigin`, and when the shared origin exists it refuses
outright if the *local* lock is missing or malformed. Walking the realistic cases:

- **Second machine / different account.** The local lock cannot exist there → the operator is told
  `established local budget lock is missing at <their own cache path>; refusing recreation`. Nothing
  was ever "established" on that machine; the real cause — this ledger belongs to another
  environment — is never printed.
- **Relocated ledger.** The lock path is derived from the canonical ledger dir (`:60`), so a moved
  ledger looks for a lock file that does not exist → the same local-cache message. The precise
  reason `"cache path or ledger location changed"` (`:221`) is never reached.
- **A pre-v2 deck.** A v1 lock file is empty, so `read()` fails the `Size() != 65` check (`:264`)
  with `invalid local budget lock identity`. The reason `"unsupported or malformed origin version"`
  (`:214`) is never reached.

Of the four parsed reasons, only `"hostname changed"` (a renamed host, same cache, same ledger) and
`"local lock identity changed or was copied"` (the `review_test.go:64-84` fixture) are reachable in
practice. Every path still fails closed and preserves charges, so this is diagnostic quality, not
correctness — but it is the exact thing C-MINOR-2 asked for, and `docs:289` ("Old v1 origins and
relocated ledgers refuse rather than silently repinning") is true about the refusal while the
message points at the wrong artifact. With no migration command, a misdirected message is the
operator's only signal.

**Suggested correction:** when `!newOrigin` and the local identity is missing or invalid, read the
origin first and report a host or `lockPath` difference from it — no token is needed for that
comparison — and fall back to the local-identity message only when host and path both match.

### [MINOR] D-MINOR-2 — a settle with unknown actual discards a known reserve, re-wedging the scope and letting an operator ceiling land under it

**Locator:** `internal/budget/ledger.go:245-251`; interacts with `:194-200` and `:153`.

`ExposureError` reads `n := entry.ReserveMicros; if entry.Settled { n = entry.ActualMicros }`. An
entry reserved at a known, cap-validated conservative maximum and later settled with a nil actual
therefore contributes **unknown**, not its reserve. Two consequences:

1. `RequireKnownCost` does not keep exposure known. It blocks nil *reservations* only, so the
   documented and expected case — "a provider may still report unknown terminal cost"
   (`ledger.go:47-48`, `docs:299-301`) — still wedges every later cost-limited reservation until an
   attended reconcile, once per poisoned entry. The docs are honest that the price stays unknown;
   neither the docs nor `Limits` say the *previously known bound stops counting*.
2. Because `observed` is `ActualMicros` for a settled entry (`:195-197`), such an entry is
   reconcilable, and `:249-251` then prefers the operator ceiling. A 7 USD reserve can become 0
   exposure with `--ceiling-micros 0`, and the zero-ceiling notice (`budget.go:72-74`) does not
   mention that a known prior reserve is being superseded. `ReconcileUnknown`'s contract — "cannot
   replace a known monetary observation" — holds for the *observation* and not for the bound.

**Suggested correction (one change covers both):** for a settled entry with a nil actual, fall back
to `ReserveMicros` before consulting reconciliations, or use `max(ReserveMicros, latest ceiling)`.
Exposure then stays known, the scope does not wedge, and no operator decision can silently drop
below a bound the cap already validated. If the current semantics are deliberate, state them in
`Limits.RequireKnownCost`'s comment and at `docs:299-301`, and name the prior reserve in the
`budget.go:72-74` notice.

### [MINOR] D-MINOR-3 — a valid BLOCK appended by a failing process still leaves no positive event

**Locator:** `internal/app/consensus_request_signoffs.go:654-660` vs `:186-205`.

`validateRequestedSignoff` treats a BLOCK as an **error** (`:658-660`) and returns a zero
`Signoff`. So when a process appends a well-formed BLOCK and then exits non-zero, control takes the
`validateErr != nil` branch at `:187-193` and the positive event is never written — the same
evidence gap C-MINOR-6 closed for ACCEPT, now confined to the verdict that most needs a durable
record. The error string names it, but a string in a returned error is not the event trail; a
consumer reconciling from events sees only `agent.handoff.failed` (interactive) or nothing at all
(headless).

**Suggested correction:** split the append-validity check from the verdict classification, so a
structurally valid BLOCK emits its own event (e.g. `agent.signoff.block-recorded`, or the same
event with `signoff_status` set) before the command fails as it already does.

### [MINOR] D-MINOR-4 — the unsupported-platform refusal advises something that cannot be done, and has no test seam

**Locator:** `internal/app/budget.go:14-20`; `internal/app/budget_attended_other.go:5-7`.

On a target matched by `termios_other.go` that is not Windows (illumos/Solaris, AIX, plan9, wasm),
`supported` is the constant `false`, so `reconcile` refuses with "preserve the ledger and use the
supported original environment". But the ledger's origin is pinned to *that* environment
(`lock.go:182-252`), and relocating it refuses with no migration available (`lock.go:68`,
`docs:288-292`). For a ledger created there, the advice is unsatisfiable and the honest answer is
"there is no recovery on this target yet".

Separately, `supported` is not threaded through `runBudgetControl` the way `attended` is (`:26`),
so that branch cannot be exercised from a test on any host — `budget_test.go` covers the attended
and unattended paths only.

**Suggested correction:** say plainly in that message (and at `docs:271-274`) that a ledger
originating on an unsupported target has no attended recovery until the migration control lands;
and pass `supported` into `runBudgetControl` so the message has a regression.

### [MINOR] D-MINOR-5 — lock-free `Inspect` now reports a benign concurrent replacement as an error

**Locator:** `internal/budget/ledger.go:103` → `:327-342`.

Removing the lock from `Inspect` was right, and `ReplaceSyncedFile` guarantees a reader sees a
complete old or new file. But `read` fails hard if the file is replaced between its `Lstat` and its
`Open`/`SameFile` check (`:339-341`, "budget ledger changed during open"), and `Inspect` does not
retry. Before this change the lock serialized inspection, so this race did not exist; now a
diagnostic read concurrent with any reservation can fail with a message that reads like tampering.
It fails closed, hence MINOR. `docs:270-271` discloses the *staleness* of the snapshot but not this
failure mode.

**Suggested correction:** retry the read a small bounded number of times inside `Inspect` before
surfacing the error, or map it to a distinct "ledger is being replaced, retry" error.

### [NIT]

- `signoff_status` in the new event (`consensus_request_signoffs.go:200`) records the **raw**
  parsed status, not `consensus.CanonicalStatus`, which is computed one function away at `:654`.
  `app_test.go:1277` asserts `"accept"`. A machine consumer will have to canonicalize spellings
  that real agents write as `✅ ACCEPT`. Recording the canonical value (or both) costs nothing now
  and cannot be changed later without breaking the trail's readers.
- No regression covers the *negative* half of the C-MINOR-6 contract — a process that exits
  non-zero **and** appends invalid or no content must emit no event. It is correct by structure
  (`:187-193` returns first), and every existing forged/rewrite fixture exits 0, so the combination
  is untested. A fake CLI that forges a second signoff and exits 7, asserting zero events, closes
  it.
- For an **all-headless** selection, `run.created` is written only under
  `hasNonHeadlessLaunch` (`:139-153`), so a headless failure-with-valid-append creates a run
  directory whose only line is the new event. `runstate` derives idea/mode from `run.created`
  (`runstate.go:130-138`), so the run lists with empty fields and `resumePendingConsensusSignoffs`
  skips it (`app.go:1351-1353`) — inert, but a run record with no header.
- `publishOrigin` (`lock_unix.go:22`, `lock_windows.go:22`) now publishes the lock **identity**
  too (`lock.go:315`); the Windows comment "a competing complete origin wins and must match" reads
  oddly for that second use. A neutral name (`publishExclusive`) would match both callers.
- `verifyLockIdentity` compares `err != io.EOF` directly (`lock.go:170`) rather than
  `errors.Is(err, io.EOF)`. Correct for `os.File.ReadAt` today; brittle if the read is ever
  wrapped.
- Windows attendance probes **stdin only** (`budget_attended_windows.go:8-13`) while unix probes
  stdin **or** stdout (`termios_unix.go:16-20`). Windows is the stricter of the two, so it fails
  closed, but the asymmetry is undocumented at `docs:271-273`.
- `review_test.go:241-274` now runs the 12-process fixture **twice** (24 spawned processes) against
  a lock that polls every 20 ms with no backoff or queueing (`lock.go:149`), with a 10 s child
  deadline (`:285`). Facilitator-reported passing, including under `-race`; on a slow shared mount
  this remains the first candidate to flake, and a starved child surfaces as
  `t.Errorf("child failure")` rather than as a cap violation, which would misdescribe the failure.

## 5. Summary

- **Concur closed:** C-MAJOR-1, C-MAJOR-2, C-MINOR-1, C-MINOR-4, C-MINOR-5, C-MINOR-6, and all
  four NIT dispositions. C-MAJOR-1's third element (the identity read/open/lock race) was found and
  closed by the implementer, not by me, and the fixture at `review_test.go:86-108` proves it.
- **Concur partially:** C-MINOR-2 (correct code, unreachable in the likely cases → D-MINOR-1);
  C-MINOR-3 (opt-in works; exposure still goes unknown through `Settle` → D-MINOR-2); C-MAJOR-3
  (Windows genuinely fixed in source; exotic targets and the test seam remain → D-MINOR-4).
- **Disagree:** nothing. No disposition overstates what the code does, and each one I checked is
  narrower than or equal to what I can verify in source.
- **New this pass:** D-MAJOR-1 (conditional on RECALL (b)); D-MINOR-1..D-MINOR-5; seven NITs.
- **Unresolved / not verifiable by me:** every Windows runtime behavior, including D-MAJOR-1, the
  `MoveFileEx`/`os.IsExist` fallthrough, and `os.SameFile` semantics on Windows; every test result
  (SECONDARY on the facilitator); the tree's SHA; behavior on a case-insensitive or network
  filesystem other than the shared volume the facilitator used.

## 6. Standing obligations this note does not touch

Unchanged and still open: per-launch/action budget wiring across launch/manual/driver/resume/BLOCK;
the origin-mismatch/legacy migration control (which D-MINOR-1 and D-MINOR-4 both make more
load-bearing, and which I would treat as a prerequisite for wiring real callers rather than as
generic follow-up); complete launch coverage (AC-T1); the packet experiment (AC-P2); the
≥20-attempt reconciliation (AC-T2); the twelve-task pilot (AC-X1/X2); Kimi typed-evidence
integration; readiness/liveness validation; any quorum amendment; and the full Phase-6/7/8 cycle.
**Nothing here waives any of them, nothing here closes a finding, and nothing here is a signature.**
