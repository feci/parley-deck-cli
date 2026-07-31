---
idea: integrate-parley-bidding-addon
review-round: 7
agent: kimi-1
date: 2026-07-30
reviewed-commit: 49fc3ec
---

## Verdict

ACCEPT

Both round-6 findings are fixed at `49fc3ec` and I re-verified each by measurement, not just
by reading the diff. Nothing else regressed. I concur with codex-1's ACCEPT.

## Outstanding findings — closed or not

1. **Omitted-`skill` marker exempt from the identity check (codex-1 round-6 MAJOR) —
   CLOSED.** The `!== undefined` exemption is gone (lib/installer.js:1342); an absent identity
   now gets its own message. Measured at `49fc3ec` on an isolated archive: full install, then
   delete only the `skill` field from the core and bidding markers, payload and all other
   fields untouched. Both units report `malformed` + `managed:false` ("the install marker
   carries no skill identity, so this directory cannot be confirmed as ..."), `doctor.ok ===
   false`; `install --only parley-bidding` and `uninstall --only parley-bidding` both return
   `blocked` for core and bidding. Health and both mutations give one answer — the round-5/6
   contradiction class is closed through this path.
2. **Filtered read re-labels a residual add-on as `selected:true` and green (my round-6
   MINOR) — CLOSED.** `selected` is now read from the core marker via `markerAddonNames` for
   read commands; install/uninstall keep the requested set because they are writing the
   selection (lib/installer.js:894-914). Measured: full install, then `install --force
   --only parley-design` to make bidding a residual. Unflagged doctor: bidding
   `valid-unselected`, `selected:false`, `managed:true`, `ok:false`. `doctor --only
   parley-bidding` on the same tree: narrowed to `[parley-deck, parley-bidding]`, bidding
   `valid-unselected`, `selected:false`, `ok:false`. `status --only parley-bidding` agrees
   (`valid-unselected`, `selected:false`). The scoped opt-out probe and the unflagged gate now
   assert the same recorded fact.
3. **Cycle-5 narrowing — INTACT.** The filtered read above inspects only what was asked and
   does not accuse the other recorded add-ons; the other three residuals simply do not appear
   in the `--only` read.
4. **2.0.0-shape markers — still valid.** Measured: full install, all six markers rewritten
   to the exact 2.0.0 shape (`skill` present, no `markerSchema`, no `manifest`; core keeps
   `addons`). All six units `valid`, `managed:true`, doctor green. The stricter predicate
   refuses no released marker shape.
5. **Rounds 4–6 closures (probe cache key, `managed` on missing units, wrong-skill
   ownership, `--only` as recorded selection, probe cwd, 3.9.6 record, `valid-unselected`
   semantics) — hold.** Now covered by the committed regression suite, which passes in full
   (below).

## New findings

None.

## Release judgement

Yes. `49fc3ec` is releasable as 2.1.0. The shipped payload is byte-identical to its first
integration commit, the installer's health and ownership logic now gives one consistent
answer across doctor, status, install and uninstall for every marker shape I could construct,
and the suite is green on every supported leg. No further change is required for release.

## What I verified

- Reviewed an isolated `git archive` of `49fc3ec` in `/tmp` with `node_modules` symlinked in;
  the source worktree was never mutated (`git status --porcelain` clean before and after; all
  temp dirs removed).
- Full node suite: **314 tests, 0 fail** (node v26.5.0, python3 3.14.6).
- Python leg: **54/54 across 7 files** on python3 3.14.6. With `/usr/bin/python3` 3.9.6 first
  on PATH the runner refuses as designed ("python3 is 3.9, but the add-on declares >=3.10")
  and runs zero tests.
- Manifest check: `parley-bidding: ok (47 files,
  sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d)` — matches the
  recorded aggregate.
- Payload frozen: `git diff 714712f..49fc3ec -- skills/parley-bidding` is empty; 48 tracked
  files, zero `__pycache__`/`.pyc`.
- Both fix diffs read line-by-line against their commit messages; behavior then confirmed by
  the probes above (omitted-`skill`, residual filtered read, 2.0.0-shape downgrade).
- Round-7 cross-check: codex-1's ACCEPT at the same commit reports the same closures and the
  same suite counts; no divergence.
