---
idea: integrate-parley-bidding-addon
review-round: 8
agent: kimi-1
date: 2026-07-30
reviewed-commit: ebe269e
---

## Verdict

ACCEPT

Full-scope round at `ebe269e`, run because `00-prompt.md` sets `track: deliberation` and
`strict_gate: true` and cycles 8–9 changed installer code and tests after round 7's unanimous
accept. I confirmed the gate text in `00-prompt.md:5-6`; codex-1's second refusal was
procedurally correct, and this round is the right remedy. Everything load-bearing was
re-measured at the new commit, including what I accepted in round 7.

## Outstanding findings — closed or not

1. **F20 — preflight walked the source only for manifested add-ons (codex-1 round-1 MAJOR) —
   CLOSED.** `preflightSkillUnit` now runs `firstCopyObstacle` over `copySourcesFor(unit)` for
   every unit, manifest or not (lib/installer.js:985-995), and `installTarget` still preflights
   all units before the write loop, reporting non-blocked units as `skipped` with zero writes
   (lib/installer.js:1079-1107). I checked parity against `copyRecursive` (lib/installer.js:1364):
   the walk covers every refusal the copier has — symlink (via `lstatSync`, so dangling links
   are caught), non-regular file, unreadable or unstat-able entry — and is a strict superset of
   what gets copied (`listVisibleEntries` filters only root-level `.DS_Store`; the walk covers
   it too, fail-closed on a git-ignored file — considered, dismissed as inconsequential).
   Measured independently of the committed regressions: I staged a package with a **dangling**
   symlink nested two levels deep in a manifest-free add-on sorting last (`zz-probe`), ran
   `installCommand` against an empty HOME — `ok:false`, message `Refusing to copy symlink in
   skill payload: references/nested/dangling.md`, all six other units `skipped`, and the HOME
   left completely empty. Zero writes, not a partial fleet.
2. **F23 — cycle 8 excluded the core while claiming "every source unit" (codex-1 round-8) —
   CLOSED.** `copySourcesFor` enumerates the core's `PAYLOAD_ENTRIES` plus existing
   `OPTIONAL_PAYLOAD_ENTRIES` (lib/installer.js:1003-1016) — exactly the sources
   `copyPayloadAtomically` copies from (lib/installer.js:1271-1285; the antigravity extra copy
   reuses `CORE_SKILL_DIR/SKILL.md`, already in `PAYLOAD_ENTRIES`). `firstCopyObstacle` accepts
   single-file roots, which `SKILL.md`/`plugin.json`/`gemini-extension.json` require. The
   second regression (symlink in the core's `references/`, zero writes asserted) passes.
3. **F21/F24 — D-2's "byte-level check / untouched" claim (codex-1 round-1 MINOR, plus the two
   surviving table rows) — CLOSED.** D-2 in IMPLEMENTATION.md is narrowed to what a file count
   plus cache scan establishes and states the SHA inventory was never captured and cannot be
   reconstructed. I grepped both documents for residual `untouched`/`byte-level` claims: the
   verification-table rows now state what the evidence shows (review/consensus.md:159 says
   explicitly "**Not** proof it is byte-for-byte untouched"); remaining uses of the word are the
   historical-correction narrative or unrelated contexts.
4. **F22 — the record claimed a stricter grammar than the guard enforces (codex-1 round-1
   MINOR) — CLOSED.** The guard comment now states the continuation exception plainly
   (test/design-addons.test.js:1109-1118), and the grammar test asserts both directions: the raw
   `PUBLISHED_PYTHON` refuses a trailing `\`, the sentinel-stripped form accepts the exact
   shipped shape, an interior backslash still refuses. I confirmed both shipped multi-line
   commands exist (skills/parley-bidding/SKILL.md:111-117) and D-3 in IMPLEMENTATION.md tells
   the same story as the code.
5. **F25 — "several hermes rounds produced nothing" — CLOSED.** IMPLEMENTATION.md and the
   consensus now say round 5 only. I had the artifacts re-enumerated: hermes-1 files exist in
   rounds 1, 2, 3, 4, 6, 7 and are missing in round 5 only.
6. **Consensus verdict table — accurate.** Re-derived independently by counting severity
   headings in all 19 review files: every cell matches, including codex-1 round 1 = 3 MAJOR +
   2 MINOR, my round 1 = 1 CRITICAL + 2 MAJOR + 2 MINOR, my round 4 = ACCEPT with 3 NITs
   (retrospectively finalized, genuine), my round 5 = PENDING draft.
7. **Everything I accepted in round 7 — holds.** The diff touches only `preflightSkillUnit`,
   the two new helpers, and tests; the health/ownership/selection logic that rounds 4–7
   hardened is untouched, the payload is frozen (`git diff 714712f..ebe269e --
   skills/parley-bidding` is empty; 48 tracked files; zero caches), and the full suite that
   carries those regressions is green (below).

## New findings

None.

## Release judgement

Yes — `ebe269e` is releasable as 2.1.0. The one deferred item (only `parley-bidding` ships a
manifest, so a universal third-party install reports one `valid-unmanaged` and five
`malformed`) is stated plainly in CHANGELOG.md:68-74, referenced from README.md:81, recorded
as consensus follow-up #1, and `FINAL.md` B3.11 holds the other add-ons unaffected. It is a
documented limitation with a fail-closed direction, not a release blocker.

## What I verified

- Reviewed an isolated `git archive` of `ebe269e` in `/tmp` with `node_modules` symlinked in;
  the source worktree was never mutated (`git status --porcelain` clean before and after; all
  temp dirs removed). HEAD at review time: `ebe269e15a79d0a992d51fe568e8a9cdb895fb55`.
- The full diff `49fc3ec..ebe269e` read line-by-line (3 files, +172): the preflight source
  walk, the two symlink regressions, the guard comment and grammar assertions.
- Full `npm test` on python3 3.14.6: **316 node tests, 0 fail**; Python leg **54 OK across 7
  files (4+20+2+3+15+3+7)**; manifest check `ok (47 files,
  sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d)` — matches the
  recorded aggregate. Count arithmetic consistent: 314 at round 7 + 2 new regressions = 316.
- Node leg with `/usr/bin/python3` 3.9.6 first on PATH: **316/0**. Python leg under 3.9.6:
  refuses with `python3 is 3.9, but the add-on declares >=3.10`, exit 1 — the designed F2
  refusal, zero tests executed through it.
- Not reproduced locally: the Python leg on 3.10/3.11 — those interpreters are not installed on
  this machine. That pair rests on the implementer's measurement and the CI matrix; everything
  runnable here was run.
- End-to-end install probe from the archive into an isolated HOME: six-skill install,
  `doctor --json` `ok:true`, all six units `valid` / `managed:true` / `selected:true`; bidding
  marker aggregate matches the manifest; runtime probe `{"ok":true,"requirement":">=3.10"}`.
- Independent B5 probe (my own staging, dangling nested symlink in a manifest-free add-on):
  preflight failure, six units skipped, zero writes — the partial-fleet class is closed through
  this path.
- Record consistency: verdict table and hermes outage re-derived from the review files (both
  accurate); `track: deliberation` / `strict_gate: true` confirmed in `00-prompt.md`; deferred
  residual confirmed in CHANGELOG.md and README.md.
