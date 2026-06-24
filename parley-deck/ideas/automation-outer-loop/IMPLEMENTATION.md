---
idea: automation-outer-loop
status: fix-up
implementer: claude-1
started: 2026-06-24
completed: 2026-06-24
branch: parley-deck-cli#loop-engineering-impl
head-commit: (this commit)
review-round: 4
fixup-cycle: 4
---

## Summary of work

Tier 4 (the outer loop), per `FINAL.md`: LE-8 (human-brake protocol §) + LE-9
(`parley loop tick`, the one-shot human-braked discovery command). `gofmt`,
`go build ./...`, `go vet`, `go test -count=1 ./...`, and the drift guard are green.

## Implementation checklist

- [x] **LE-8 — §14 human brake** — added `## 14. Automated outer loop (loop engineering)
  — the human brake` to BOTH `COOPERATION.md` copies (live + embedded default), byte
  identical (drift guard green). It binds any automated/standing/scheduled loop to
  discover-and-draft-only: §14.1 (MAY: draft `status: candidate`, no `participants:`
  claim), §14.2 (MUST NOT without a recorded human/full-quorum gate: promote, run,
  implement, land/merge/push, finalize, edit the roster, override consensus), §14.3
  (fail-safe — when uncertain, do less).
- [x] **LE-9 — `internal/loop` package** — `Candidate`, `Config` (Enabled, disabled by
  default), `TickResult`; `ReadConfig` (absent → disabled/no-error, malformed → fail
  closed), `ReadSignals` (absent → empty/no-error), `SlugFor`/`fingerprintOf` (explicit
  fingerprint kept whole; else 8-char sha256 of source+id), and `Tick` — drafts a
  `status: candidate` 00-prompt.md per not-yet-seen signal (dedupe by slug) and returns.
  It never staffs a quorum, runs, pushes, merges, or finalizes (§14 enforced in code).
- [x] **LE-9 — `parley loop tick` command** — `internal/app/loop_cmd.go`: `runLoop` +
  `runLoopTick` (`--dir`, `--signals`, `--enable`, `--json`). Resolves `<root>/parley-deck`,
  reads config (`--enable` force-enables this one-off run, still candidate-only), reads
  signals, calls `loop.Tick`, prints a human/JSON summary. Exits 0 when disabled or when
  zero candidates (cron/idempotent-safe). Wired `case "loop"` + a usage line in `app.go`.
- [x] **Tests** — `internal/loop/loop_test.go` (disabled-writes-nothing, enabled-drafts-
  candidate, dedupe, fingerprint default + explicit, ReadSignals/ReadConfig edge cases);
  `internal/app/loop_cmd_test.go` (disabled exits 0, `--enable` drafts candidate-only,
  `--json` shape).

## Deviations from FINAL.md

- None of substance. The MVP reads a signals file (as specified); live connectors
  (GitHub/CI APIs) and optional human-confirmed `parley run` from a promoted candidate are
  explicit out-of-scope follow-ups.

## Notes for reviewers

- **The §14 brake is the security boundary.** Try to break it: can any code path in
  `loop.Tick` / `runLoopTick` stand up a quorum, set `participants:`, flip a candidate to
  `round-01`, call `parley run`, push, merge, or finalize? It must not.
- **Disabled-by-default + fail-safe.** Confirm an absent config writes nothing; a malformed
  config fails closed (disabled), never silently enables.
- **Dedupe integrity.** Confirm the slug is stable for the same signal and that an existing
  candidate dir is skipped (no overwrite). Check an unexpected `os.Stat` error fails closed.
- **Prose vs frontmatter.** The candidate prompt mentions `participants:` in its Promotion
  note (prose); the invariant is that there is no frontmatter `participants:` key.

## Fix-up cycle 1 (round-01 review consensus)

All five agreed fixes from `review/consensus.md` applied. `gofmt`, `go build ./...`,
`go vet`, `go test -count=1 ./...` (incl. drift guard) green. The CRITICAL was reproduced
before and confirmed closed after, end-to-end.

- **AF1 (CRITICAL — all 3 reviewers)** — `cleanField` flattens CR/LF/tab/control chars to
  spaces on `Source`/`ID`/`Title`/`Detail` before they reach the candidate prompt, closing
  the YAML frontmatter injection (no more attacker-injected `participants:`/`status:`/`checks:`
  keys). `Source` is also validated against `{commit, ci, issue, manual}` (`validSources`);
  an unknown source is rejected (`TickResult.Rejected`). (`internal/loop/loop.go`.)
- **AF2 (MAJOR — codex + antigravity)** — identity is now an 8-char sha256 over an
  unambiguous `strconv.Quote` canonical key (`dedupeDigest`), not a lossy sanitized string.
  Fixes `a/b`==`a:b`, the colon boundary shift, and emoji→`"x"` collapse together.
- **AF3 (MINOR — codex)** — `runLoopTick` returns the disabled result before reading signals,
  so a disabled tick is inert/cron-safe even with a malformed signals file.
- **AF4 (MINOR — codex + antigravity)** — atomic slug claim (`os.Mkdir` + `O_EXCL`) replaces
  the `os.Stat`→`MkdirAll`/`WriteFile` TOCTOU; concurrent ticks can't double-create/clobber.
- **AF5 (MINOR — hermes)** — `## Constraints`/`## Non-goals` carry a `(to be filled on
  promotion)` placeholder.

New tests: `TestTickRejectsFrontmatterInjection`, `TestTickRejectsUnknownSource`,
`TestColonBoundaryNoCollision`, extended `TestSlugFingerprint` (all in
`internal/loop/loop_test.go`).

**Deferred (own follow-ups):** DF1 reconcile `ReadFrontmatter` last-wins vs
`readFrontmatterField` first-wins (pre-existing parser inconsistency; AF1 makes it
unreachable via the loop); DF2 live connectors + human-confirmed run; DF3 require an
initialized deck for `--enable`.

## Fix-up cycle 2 (round-02 review consensus)

Round-02 confirmed the round-01 CRITICAL closed and surfaced four new agreed fixes (two of
them regressions introduced by cycle-1's own fixes — the payoff of a real re-review). All
applied; `gofmt`, `go build ./...`, `go vet`, `go test -count=1 ./...` (incl. drift) green.

- **AF6 (MAJOR)** — `Detail` is no longer `cleanField`'d (that flattened legit multi-line
  logs/traces). It is rendered in a `## Signal detail` section as a markdown indented block
  with newlines preserved; `Source`/`ID`/`Title` stay `cleanField`'d (frontmatter safety).
  (`indentDetail` in `internal/loop/loop.go`.)
- **AF7 (MAJOR)** — the atomic claim moved from the directory to the PROMPT FILE
  (`O_CREATE|O_EXCL` on `00-prompt.md`); an empty dir left by a crashed prior tick is now
  healed instead of suppressing the signal forever, and a failed write removes the partial
  file. (`writeCandidate`.)
- **AF8 (MINOR)** — `cleanField` now also flattens U+2028/U+2029/U+0085 (YAML line breaks),
  keeping its "no line break injects a key" contract true under a future YAML parser.
- **AF9 (MAJOR)** — `dedupeDigest` widened from 8 hex (32 bits, birthday-collidable — codex
  found a real collision) to 32 hex (128 bits), defeating deliberate dedupe-suppression.

New tests: `TestTickHealsPoisonedEmptyDir`, `TestTickPreservesMultilineDetail`,
`TestTickFlattensUnicodeSeparators`.

**Dismissed/deferred (round-02):** old-slug migration (N/A — new command, no deployed
candidates); rejected-source label (cleanField already strips ANSI/`<0x20`); `SlugFor`
fallback (deliberate totality of the exported helper); DF4 case-insensitive `Source`
(signals are lowercase per FINAL.md; `EqualFold` is future polish).

## Fix-up cycle 3 (round-03 review consensus)

Round-03: antigravity-1 clean, hermes-1 2 MINOR + 1 NIT (no blocker), codex-1 a NEW MAJOR
(symlink escape) the others missed. All applied; suite + drift guard green; AF10 verified
end-to-end (a planted symlink is refused, nothing written to the target).

- **AF10 (MAJOR)** — symlink-safe slug claim: `os.Mkdir` the exact slug dir; on `ErrExist`,
  `Lstat` and reject a symlink / non-directory (a loop must never write outside
  `parley-deck/ideas/<slug>/`). (`writeCandidate`.)
- **AF11 (MINOR)** — `indentDetail` normalizes U+2028/U+2029/U+0085 (and CR/CRLF) to `\n`
  before splitting, so every logical Detail line is four-space indented (no column-0
  heading/fence/key even under a markdown renderer).
- **AF12 (MINOR)** — `indentDetail` uses `TrimSpace` (drops leading blank lines too).
- **AF13 (test)** — `TestTickDetailCannotInjectHeadingOrFence` pins the no-injection contract.
- **F3 (NIT)** — clarified the `cleanField` doc comment ("C0 control characters").

New tests: `TestTickRejectsSymlinkedSlugDir`, `TestTickDetailCannotInjectHeadingOrFence`.

**Resolved/deferred:** codex's `internal/runner` durable-kill failure is the codex-sandbox
limitation (green locally); DF1/DF4 carry forward.

## Fix-up cycle 4 (round-04 review consensus)

Round-04: antigravity-1 clean ("converged"); codex-1 + hermes-1 both found the same MAJOR
(the symlink class one level up) + codex a MINOR. All applied; suite + drift guard green;
AF14 verified end-to-end (a planted `ideas/` parent symlink is refused, nothing written to
the target).

- **AF14 (MAJOR)** — a symlink at `parley-deck/ideas/` (parent of the slug) was still
  followed via `os.MkdirAll(ideasDir)`. Fix: a `safeMkdir` helper guards BOTH `ideas/` and
  `ideas/<slug>` (os.Mkdir + Lstat-reject-symlink/non-dir), plus a depth-complete
  `assertInsideDeck` containment check (`EvalSymlinks` + `Rel`) that rejects a slug dir
  resolving outside the deck through a symlink at ANY ancestor.
- **AF15 (MINOR)** — `indentDetail` now normalizes the C0 line separators (`\v`, `\f`,
  U+001C/1D/1E) in addition to CR/U+0085/U+2028/U+2029, so no Detail token reaches column 0
  under any line splitter; other C0 controls become spaces.

New tests: `TestTickRejectsSymlinkedIdeasParent`, `TestTickIndentsC0SeparatorsInDetail`.

**Out of scope (follow-up):** DF5 — `retro propose` is vulnerable to the same `ideas/`
parent-symlink class (hermes); hardening `retro` is a separate idea. The loop's
`assertInsideDeck` is now stricter than that precedent.
