---
agent: kimi-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 4
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review, round 4 — release candidate 1.44.0 + 2.8.0

## Summary

**[PRIMARY — CONFIRMED]** The budget, as the tree stands now (fix-up cycle 4 applied), survives
every attack I could construct. I probed it myself, in an isolated copy, against both the cycle-3
and the cycle-4 snapshots: single-record deletion in either direction, both-record deletion,
forged markers, odd round-dir names, corrupt/absent cursor, errored `Fixup`, failed post-fix-up
checks, crash windows on both sides of the reservation, the AF2 boundary at/above the cap, and the
dirty-tree backstop in a real git repo. Every vector is either closed or disclosed verbatim in both
CHANGELOGs under "What this is not". The AF2 off-by-one and the late reservation — the two genuine
defects I found in the cycle-3 code this session, independently of @codex-1's round-04 file (see
the validity event) — are fixed in the current tree, and I re-ran the probes against the fix.

The verdict is NOT CLEAN, narrowly, and not on behavior: **fix-up cycle 4 changed three
safety-critical branches and added no test for any of them.** I reverted each of the three in the
isolated copy and the full suite stayed green every time — the suite cannot distinguish cycle 4's
fixes from their absence, which makes cycle 4's own one-line verification claim ("go build, go vet,
full Go suite green") vacuous as evidence for those fixes. The tests exist (my probes below) and the
fix is mechanical. Plus three text NITs and one MINOR residual, none behavioral.

This is the second consecutive round in which the candidate mutated under review; everything below
is hash-pinned.

## Review-validity event (read before any finding)

**[PRIMARY — CONFIRMED]** Timeline (file mtimes, my shell history; my session opened ~15:46 CEST):

```text
15:46   review/round-04/ exists and is EMPTY; round-03/ holds codex-1 + hermes-1 only
~15:47–15:55  I read IMPLEMENTATION.md (through cycle 3), round-03/codex-1.md,
        round-03/hermes-1.md, round-02/kimi-1.md, and the full diff
~15:55  isolated copy made (rsync, no .git); it captured the CYCLE-3 code
~15:56  my adversarial probe file written — BEFORE any round-04 file existed
16:01   review/round-04/codex-1.md written   (NOT by me; not read by me until later)
16:02   review/round-04/hermes-1.md written  (NOT by me; not read by me until later)
16:02:49  impl.go + driver.go rewritten (fix-up cycle 4) — WHILE my probes were running
16:03:36  both CHANGELOGs gain the trust-boundary disclosure
16:03:58  IMPLEMENTATION.md cycle-4 section
~16:05–16:10  my probes + reversions A–D complete against the CYCLE-3 snapshot
~16:10  hash mismatch vs the tree reveals the mutation; only NOW do I read the round-04 files
16:20   review/round-03/kimi-1.md lands (21,878 bytes) — a parallel kimi-1 session; I read
        it only after all my probe/reversion evidence was captured
16:22–16:25  IMPLEMENTATION.md, both CHANGELOGs, and impl.go comments updated again
        (its F3 comment fix); current impl.go sha256 de170e52…
```

So: the AF2 equality-stranding and the reservation-ordering findings are my own derivations from
the cycle-3 code, convergent with @codex-1's round-04 review, which I had not read. The round-3
files (both) I did read at session start — they directed my attention at AF2 and the marker
provenance question; the boundary analysis is mine. The parallel round-03/kimi-1 file reached the
same findings (its F1–F4) independently of me; I did not read it until my evidence was complete.

Snapshot pins (sha256, current tree, 16:30 CEST):
`impl.go de170e52…`, `driver.go a88f7ddb…`, `cursor.go 0df19271…`, `consensus.go 1330db98…`,
`impl_test.go c35c465c…`, `track.go c5ce3122…`, CLI `CHANGELOG.md c4f65724…`, skill
`CHANGELOG.md 105d942c…`. The cycle-3 snapshot my first probe run certified: `impl.go ee40f92d…`,
`driver.go 90e5d7f1…` (contents recoverable from my session's full `git diff` capture).

**[PRIMARY — CONFIRMED]** My footprint on both working trees is zero: all probes and reversions ran
in `/tmp/kimi1-r04-iso` (module copy with its own GOCACHE/GOMODCACHE under /tmp); both
`git status --short` file-sets are identical before and after my session; no git write commands were
issued by me in either repo.

## Findings

### F1 [MAJOR] Cycle 4's three safety behaviors ship with no regression test — and the suite is blind to their reversion

**[PRIMARY — CONFIRMED]** Cycle 4 (16:02:49) changed three safety branches —
reserve-before-`Fixup` (impl.go:317-325), the AF2 gate's `>` boundary (impl.go:147-153), and the
corrupt-cursor escalation (driver.go:235-248) — and touched no test file (`impl_test.go` mtime
15:41, predating cycle 4; hash-identical to my pre-mutation copy). In the isolated copy I reverted
each one and ran the suite:

```text
revert AF2 `>` back to `>=`            → full driver suite GREEN (only my probe went red)
move the reservation back after Fixup  → full driver suite GREEN (only my probes went red)
drop the corrupt-cursor escalation     → full driver suite GREEN (only my probe went red);
                                          TestCorruptCursorIgnoredRebuildRecovers stays green —
                                          it never calls Advance
```

Three independent one-line reversions of this release's core safety semantics, none detectable by
the shipped suite. The cycle-3 fixes, by contrast, are each pinned: my reversions of the
carry-forward, the whole AF2 gate, and the exact-dir predicate each turned exactly the claimed tree
test red. Cycle 4 also broke this idea's own verification discipline: cycles 1–3 each recorded
isolated-copy reversion checks in IMPLEMENTATION.md; cycle 4 records only "go build, go vet, full
Go suite green" — a suite that, per the three reversions above, passes with or without the fixes.

This is the class the release itself exists to close: a safety boundary with no tripwire, in the
one module that has already regressed twice mid-review (the AF2 off-by-one was *introduced* by
cycle 3 while fixing another defect; nothing in the suite noticed — my probe did). The fix is
mechanical: my probe file already implements all four missing tests against the package's existing
helpers — AF2 completes at equality, AF2 escalates over the cap, an errored `Fixup` still spends
the cycle (and a failed reservation aborts before `Fixup` runs), a corrupt cursor escalates. Drop
them in, and re-run the three reversions against them as the cycle-4 verification record.

### F2 [NIT] The fix-up escalation message reports the would-be cycle, not the spent count

**[PRIMARY — CONFIRMED]** impl.go:308-309: `cycle := published + 1; if cycle > cap { … "after %d
cycle(s)", cycle … }`. At cap=5 with 5 spent it prints "after 6 cycle(s)" — 6 is the *refused*
cycle. The parallel round-03/kimi-1 file noted the same. Pre-existing wording; one-line fix (print
`published`, or "cycle %d would exceed").

### F3 [NIT] The BLOCK back-edge diagnostic still counts one unrun round

**[PRIMARY — CONFIRMED]** consensus.go:96-98 still escalates with `next-1`: for `standard`, after
cross-review rounds 2 and 3 actually ran, the refusal of round 4 prints "blocked after 3
cross-review round(s) after round 1". @codex-1's round-3 NIT; unfixed in cycles 3–4. The guard
itself is correct (re-verified: `TestBlockedConsensusRespectsTheHardCrossReviewCap` green; my
round-2 guard-removal reversion still red on today's file — consensus.go is untouched since,
sha256 `1330db98…`).

### F4 [NIT] The Cursor struct doc is now false for one field

**[PRIMARY — CONFIRMED]** cursor.go:42-44 still says "Cursor is a rebuildable cache… Rebuild
derives every field from on-disk artifacts". `FixupCyclesPublished` is precisely the field Rebuild
must NOT derive — the carry-forward in Advance exists because it is not rebuildable. The field's
own doc comment is correct; the struct doc contradicts the design. One sentence.

### F5 [MINOR] A per-marker stat error is still swallowed as "absent"

**[PRIMARY — CONFIRMED]** `markedFixupCycles` counts via `fileExists`, which maps any stat error
(permission-damaged round dir) to "no marker". The review-dir root failing IS surfaced (escalates);
a single damaged child round is not. With the cursor intact this is harmless (the max covers it);
it matters only if the cursor is also lost — two coincident independent failures, inside the
disclosed floor. Hardening is one line (`statRegular` + propagate) and can ride
`fixup-budget-trust-anchor`. Not a blocker.

## Answers to the requested questions

### Q1 — is there ANY remaining path that lowers the count or opens a fix-up past the cap?

**Within the ratified scope: no, and I probed every named vector myself. Beyond it: the disclosed
floor, now honestly documented in both CHANGELOGs.** All probes ran in the isolated copy against
the current tree (8/8 pass); the tamper directions:

- **Delete the markers** → cursor holds; still escalates at cap (tree tamper subtest 1; my
  carry-forward reversion turns exactly that subtest red). **Delete the cursor / whole run
  directory** → markers hold (my probe 2a). **Delete BOTH** → the count restarts and a cycle runs
  (my probe 2c: `action=fixup` where the cap demanded escalation). That is the theoretical floor of
  any on-disk counter without an external authority; it is now stated in plain text in both
  CHANGELOGs and deferred to `fixup-budget-trust-anchor`.
- **Forge a marker** (any round, including the current one) → the count only rises → escalates
  sooner; AF2 consults the budget first (tree tamper subtest 2 + my over-cap probe: escalates with
  "fix-up budget exceeded: 3 cycle(s) recorded against MaxFixupCycles=2", no round opened).
- **Corrupt cursor** → Advance escalates: "refusing to act on an unknown fix-up budget" (my probe
  5). **Absent cursor** → fresh run, markers only — correct, and the two cases are now
  distinguished (cycle-4 fix).
- **Rolled-back git checkout** → survives: `driver.json` is untracked (verified:
  `git ls-files parley-deck/runs | grep -c driver.json` → 0), so a checkout/reset leaves the
  cursor, and the cursor alone holds the count. Reaching the count through git needs
  `git clean -fdx`-class deletion of BOTH the untracked cursor and the untracked markers — the
  disclosed both-deleted floor, not a new vector.
- **The window between `Fixup` returning and the cursor being persisted** → closed by cycle 4. On
  the cycle-3 code I confirmed the leak live (simulate the failed persist: no cursor, no marker,
  resume re-runs the same cycle unspent — `fixupCalls=2`). On the current code the reservation is
  taken BEFORE the code-writing call: an errored `Fixup` leaves the cycle spent (my probe 3:
  cursor=1, resume escalates instead of retrying); a failed reservation aborts before `Fixup` runs
  (probe 3b: `fixupCalls=0`); a crash between reservation and confirmation conservatively burns one
  cycle — the right error direction. In a real git repo the dirty-tree gate additionally backstops
  any residual re-run (probe 6: escalate "git working tree is dirty", `fixupCalls=1`).
- **Two runs sharing an idea** → sequential second run cannot re-spend (markers cross runs; my
  probe 2b setup is exactly that shape). **Concurrent** runs can double-spend into one count
  (@codex-1's round-04 probe; consistent with my reading of the per-run lock at loop.go) — inside
  the disclosed "two concurrent runs" clause, not a new finding.
- **Ordering** → the AF2 gate runs before archive/open and the main gate before reservation; no
  path reaches `OpenReviewRound` or `Fixup` past the cap on any state I could construct.

The scope line itself is the one open dispute: @codex-1's round-04 CRITICAL requires idea-scoped,
serialized, out-of-repo anchoring; the implementer declined on scope grounds and disclosed instead.
On the merits I accept the deferral: every *careless* path is now closed, the remaining ones
require deliberate two-target destruction by a participant who could equally delete the idea
itself, and an out-of-repo trust anchor was not ratified by this idea. §15.3 applies: my
acceptance is not a sign-off on the scope call — the owner should rule on it explicitly.

### Q2 — "spent when it runs, not when it passes": is @codex-1's way right?

I am not @hermes-1, so the direct question is not mine to answer — for the record, his round-04
file says plainly: "I am persuaded. I change my position from round 3." My own position, on the
merits: **@codex-1 is right, and cycle 4 is the first version that actually implements it.** The
cap exists to interrupt churn; a fix-up that errors, crashes, or breaks the build is the churn.
Cycle 3 only moved the persist ahead of the check gate — my probes showed an errored `Fixup` and
the return-to-persist crash window both still handed the cycle back. Cycle 4's reserve-before-run
closes both (Q1). One caveat now disclosed in the CLI CHANGELOG itself: between reservation and
marker the cursor is the only record, so losing the cursor in that window loses the count — honest,
and inside the follow-up's scope.

### Q3 — is the AF2 gate correct, or does it strand a legitimate crash recovery?

**On the cycle-3 code it stranded exactly the recovery it should complete — I proved it before
cycle 4 existed. On the current code it is correct in both directions — and unpinned by any tree
test (F1).** My paired probes, same fixture, two snapshots:

- cycle-3 code (`spent >= cap`), marker at the current round with the budget spent exactly to the
  cap: `action=escalated`, "budget exhausted", `calls=[]` — the last allowed cycle's closing review
  round never opens. The normal (no-crash) path with the identical spent budget opens the round
  unconditionally (probe 1b: `action=fixup`, `open-review` called). A crash converted a lawful
  transition into an escalation — one too strict, because AF2 finishes cycle N, it does not start
  N+1.
- current code (`spent > cap`): equality completes the transition (`action=fixup`,
  `open-review`, no `Fixup` re-run — the closing review happens); strictly-over-cap still
  escalates (probe 1b). Both directions verified by swapping only the operator in the isolated
  copy.

Note @hermes-1's round-04 Q3 argues the `>=` form was correct; the probe evidence says otherwise —
at equality the transition spends nothing, and the next cycle is still refused by the ordinary
branch.

### Q4 — tests and reversion checks, all in an isolated copy

**[PRIMARY — CONFIRMED]** Everything executed by me this session; all mutations confined to
`/tmp/kimi1-r04-iso` (plus its own /tmp Go caches). Neither working tree was edited; no git write
commands anywhere.

```text
CYCLE-3 snapshot (impl.go ee40f92d…):
  go build ./... ; go vet ./...                      exit 0
  go test ./... -count=1                             26 ok + 1 no-test-files, zero FAIL
  5 first-edition probes                             AF2 stranding reproduced; crash-window leak
                                                     reproduced; both-deleted reset reproduced
  reversions (carry-forward / AF2 gate / exact-dir / persist-removal)
                                                     red exactly on the claimed tree tests —
                                                     except persist-removal: NO tree test red

CURRENT tree (cycle 4 + comment fixes; impl.go de170e52…, driver.go a88f7ddb…):
  go build ./... ; go vet ./...                      exit 0
  go test ./... -count=1                             26 ok + 1 no-test-files, zero FAIL
  go test ./internal/driver ./internal/track ./internal/protocol   ok
  8 second-edition probes, all PASS:
    AF2 at spent==cap → completes (open-review, no re-fixup)   [cycle-4 fix verified]
    AF2 over cap → escalates, no round opened                  [forgery still blocked]
    cursor deleted / markers deleted → budget survives both single deletions
    BOTH deleted → count restarts (action=fixup)               [disclosed floor]
    errored Fixup → cycle spent; resume escalates              [reserve-before-run]
    reservation-save failure → aborts BEFORE Fixup (0 calls)
    failed checks → cursor=1 persisted                         [spent-when-it-runs]
    corrupt cursor → escalates "unknown fix-up budget"
    dirty-tree backstop in a real git repo → escalate, no re-run
  reversions on the current code:
    carry-forward removed      → tree test "deleting the markers…" red      [pinned]
    AF2 gate removed           → tree test "a forged marker…" + my probe red [pinned]
    isRoundDirName → HasPrefix → TestOnlyExactRoundDirs… red                 [pinned]
    AF2 `>` → `>=`             → NO tree test red; my boundary probe red     [F1 gap]
    reservation after Fixup    → NO tree test red; my two probes red         [F1 gap]
    corrupt-cursor escalation dropped → NO tree test red; my probe red       [F1 gap]
  after each reversion: restored, sha256 re-matched to the tree, package green

SKILL repo:  node scripts/build-addon-manifest.js --check → all six skills ok
             npm test → 386 pass / 0 fail; python 3.14: 54 OK across 7 files
             (payload files predate my runs; only CHANGELOG.md changed after, re-read)

READ-ONLY tree checks: git diff --check clean in both repos; gofmt clean on all changed
Go files (the one gofmt hit, internal/app/protocol_test.go, is a pre-existing double blank
line from v1.43.0, untouched by this release); `trajectory` absent from all Go code;
`packet` matches are only the idea slug and the pre-existing signoff HandoffPacket.
```

### Q5 — are both CHANGELOGs and IMPLEMENTATION.md accurate about what ships?

**[PRIMARY — CONFIRMED]** Yes in the current texts; every load-bearing sentence was probed above.
The CLI CHANGELOG's mechanism paragraph ("maximum of two driver-authored records… reserved before
the code-writing call… once both records exist, losing either one does not lower the count… one
window with a single record…") matches the code line for line. The "What this is not" paragraph
(workspace-write participant, deleted run directory, rollback, concurrent runs) matches my probes
2a–2c and the disclosed floor. The skill CHANGELOG's "charges a cycle when it is RESERVED" and its
non-security-boundary clause match. The "Why 5" distribution I reproduced digit-for-digit in round
2; "Not in this release" (packet, trajectory payload, fast route) all re-verified true today.
IMPLEMENTATION.md's cycle-3 and cycle-4 sections match the diffs, including the honest admission
that the AF2 off-by-one was self-inflicted in cycle 3.

Residual text defects: F2/F3/F4 above. The stale call-site comment that survived three rounds
(impl.go:297-302, naming IMPLEMENTATION.md as the source) was fixed at 16:25:11 — verified by diff;
it was comment-only.

### Q6 — anything that should stop this release?

**One thing, and it is cheap.** F1: cycle 4's three safety behaviors must be pinned by tree tests
before release — the probes are written, they pass against the current code, and they turn the
three "green under reversion" gaps red. With those landed, the cycle-4 verification section
re-recorded honestly (reversions against the new tests, as cycles 1–3 did), and F2–F4 swept (three
one-liners), I expect the confirming pass to be short — provided the tree holds still for it: this
candidate has now mutated at least five times inside two consecutive review windows, and each of my
certifications is pinned to the hashes above, not to whatever the tree says next.

Not blocking, for the record: the trust-anchor scope dispute (§15.3 — owner rules, my acceptance on
the merits is logged in Q1), the F5 hardening, and the follow-up idea `fixup-budget-trust-anchor`
already named in the release notes.

## Responses to the other reviewers

**@codex-1.** Your round-04 CRITICAL/MAJORs and my independently-derived findings converge on the
two code defects (late reservation, AF2 equality) — both fixed in the current tree and re-verified
by my probes after the fix. On the CRITICAL's scope half we part ways: I accept the disclosed
deferral to `fixup-budget-trust-anchor`, for the Q1 reasons. Where I fully reinforce you is
permanent tests: my F1 is your Q6 demand, demonstrated from the other side — the suite is
*provably* blind to all three cycle-4 behaviors.

**@hermes-1.** Your round-04 CLEAN examined the cycle-3 candidate (16:02) and your Q3 defends the
`>=` gate that candidate shipped; my paired probes show that gate stranding the last allowed
cycle's closing review, and cycle 4 already adopted the correction. Your Q2 reversal is on the
record and, on the merits, right. Your trust-model argument for the floor I now partially adopt —
with the boundary stated rather than implied.

**The parallel round-03/kimi-1 file.** Landed 16:20, after my evidence was captured; I did not
write it and did not read it during my probe work. Its F1–F4 converge with my cycle-3 findings; its
"frozen pass" demand is the same as my Q6. Two independent kimi-1 sessions, same conclusion —
treat that as corroboration, not coordination.

## Independent validation record

All commands run by me this session, read-only in the trees or mutating only the `/tmp` copy.
Hashes and timeline as pinned above. The isolated copy survives at `/tmp/kimi1-r04-iso` until the
next reboot, with the full second-edition probe file at
`internal/driver/zz_kimi1_r04_probe_test.go` — the four gap-closing tests of F1 are in it in
drop-in form, written against the package's existing helpers.
