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

`skills/parley-bidding/` **has not changed since `714712f`**, its first integration commit — 47
files, aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`,
re-verified in every round. **No round found a defect in the payload this idea exists to ship.**
The same holds for the integrity mechanism, the seven Python tools, the four platform adapters,
the test runner, the CI workflow and the documentation.

Every fix-up cycle from 10 onward was in **one mechanism**: the gate that refuses an install or
uninstall plan in which two destinations would physically collide. Seventeen cycles on a
mechanism that was not the subject of this idea, but without which shipping a 47-file
security-relevant payload could not be defended.

### The seven arms the gate now refuses

Each was measured on the tree before it was fixed, and each has a regression that fails at the
commit it discriminates against:

| arm | first measured |
|---|---|
| per-target rather than fleet-wide preflight | round 8 |
| `existsSync` calling a dangling symlink absent | round 10 |
| `--force` suppressing the only destination check | round 11 |
| existence checked but not permission | round 11 |
| create/touch checked but not **dispose** | round 10 |
| stored data (marker, manifest keys) becoming a path | rounds 11–12 |
| physical collision through symlink chains, firmlinks and nesting | rounds 15–21 |

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
3. **Four claims of mine claimed more than they showed** — "byte-level check", "every destination
   path", the Python-leg sentence, and "the scaffolding is gone" when it was not. All four are
   corrected **in place in `IMPLEMENTATION.md` with a note that they were false**, rather than
   quietly rewritten.

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
   directories.
5. **Residual disposal arms** — `uappnd` directories and delete-denying ACLs pass `access(2)`
   entirely and node exposes no `st_flags`. Under the quarantine transaction these produce debris
   rather than a partial fleet.
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
  round or the five that followed. Removed from the roster on 2026-07-30 by the user's decision.
- **Rounds 1–7 rested on an incomplete reading of round 1.** The facilitator read
  `review/round-01/codex-1.md` **while codex was still writing it**: 3.9 KB and two MAJOR when
  acted on, 9.4 KB and three MAJOR plus two MINOR when finished. Three findings went unaddressed
  for six rounds. Every prompt since carries "finish, then write".
- **Round 5** — `hermes-1` outage.
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
- **Manifest check ok** — 47 files, aggregate unchanged since `714712f`.
- **All seven collision arms refused**, each with a regression that fails at the commit it names.

## Signoffs
