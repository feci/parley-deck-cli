---
idea: automation-outer-loop
phase: review-consensus
drafter: claude-1
review-round: 1
date: 2026-06-24
---

## Review consensus — round 1

Three non-implementer reviewers ran Phase 6 in refutation mode against
`git diff 1a09459..HEAD` (LE-8 §14 human brake + LE-9 `parley loop tick`), with the §14
brake as the named security boundary:

- **codex-1** — 1 CRITICAL, 1 MAJOR, 2 MINOR (`review/round-01/codex-1.md`)
- **hermes-1** — 1 CRITICAL, 1 MAJOR, 2 MINOR (`review/round-01/hermes-1.md`)
- **antigravity-1** — 1 CRITICAL, 1 MINOR, 1 NIT (`review/round-01/antigravity-1.md`)

The dogfood worked decisively: **all three independently found the same CRITICAL** — a YAML
frontmatter injection that lets an untrusted signal write `participants:` / `status: round-01`
(and a latent `checks:` → `sh -c` RCE) into a loop-drafted candidate, defeating the very §14
invariant the loop exists to uphold. codex demonstrated it end-to-end (`parley status` then
reported `status=round-01 participants=codex-1`); hermes traced the downstream RCE path. The
code *action* boundary held (no run/push/merge/finalize/roster path) — the breach was the
*output* artifact.

## Agreed fixes (all applied in fix-up cycle 1)

| ID | Severity | Raised by | Fix |
|----|----------|-----------|-----|
| **AF1** | CRITICAL | codex + hermes + antigravity | `cleanField` flattens control chars / CR / LF / tab to spaces on every signal-derived value (`Source`, `ID`, `Title`, `Detail`) before it is written into the candidate prompt, so newlines can no longer inject extra YAML frontmatter keys. **Plus** `Source` is validated against the closed set `{commit, ci, issue, manual}` (`validSources`) — an unknown source is rejected (`TickResult.Rejected`), never normalized. Verified end-to-end: an `id` carrying `\nstatus: round-01\nparticipants: [evil]\nchecks: rm -rf /` now lands as a single-line `source_id:` scalar; `parley status` reports `status=candidate participants=`. |
| **AF2** | MAJOR | codex + antigravity | Identity is now a collision-resistant 8-char sha256 over an **unambiguous** canonical key (`strconv.Quote`), not a lossy sanitized string. Fixes the explicit-fingerprint collision (`a/b` vs `a:b` both → `loop-manual-a-b`), the default-hash colon boundary shift (`"ci:"+"build"` vs `"ci"+":build"`), and the emoji-fingerprint → `"x"` collapse (antigravity NIT) in one change. The slug keeps a readable source hint + the digest. |
| **AF3** | MINOR | codex | `runLoopTick` short-circuits the disabled case **before** reading the signals file, so a disabled tick is fully inert and cron-safe even when `signals.json` is malformed (no more spurious exit 1). |
| **AF4** | MINOR | codex + antigravity (OQ) | The slug claim is now atomic: `os.Mkdir` (treat `fs.ErrExist` as a clean skip) + an `O_CREATE|O_EXCL` prompt write, replacing the `os.Stat`→`MkdirAll`/`WriteFile` TOCTOU. Two concurrent ticks can no longer both create or clobber the same candidate. |
| **AF5** | MINOR | hermes | The candidate prompt's `## Constraints` / `## Non-goals` sections carry a `(to be filled on promotion)` placeholder instead of being empty. |

New regression tests (the security suite now probes the boundary with hostile input):
`internal/loop/loop_test.go` — `TestTickRejectsFrontmatterInjection` (AF1, the named CRITICAL
vector), `TestTickRejectsUnknownSource` (AF1), `TestColonBoundaryNoCollision` (AF2), and the
extended `TestSlugFingerprint` (AF2 `a/b` vs `a:b`). This directly answers hermes F2 (MAJOR:
"a security-boundary test suite that only tests benign inputs is theater").

## Deferred follow-ups (out of scope — documented, not blocking)

- **DF1** — `protocol.ReadFrontmatter` (last-wins on duplicate keys) vs `readFrontmatterField`
  (first-wins) disagree on duplicate keys (hermes OQ#2 / codex). This is a **pre-existing**
  protocol-parser inconsistency, not introduced here; AF1 removes the loop's ability to
  produce duplicate keys, so it is no longer reachable via the loop. General reconciliation
  belongs to its own protocol/parser idea.
- **DF2** — live connectors (GitHub/CI APIs) and an optional human-confirmed `parley run` from
  a promoted candidate (FINAL.md "Out of scope"). The MVP reads a signals file.
- **DF3** — `--enable` creating `parley-deck/` in an uninitialized directory (codex OQ). It
  stays candidate-only and harmless; requiring an initialized deck is a small future polish.

## Verification after fix-up cycle 1

`gofmt`, `go build ./...`, `go vet`, `go test -count=1 ./...` (all packages incl. the
embedded-default drift guard) — green. The CRITICAL was reproduced before and confirmed
closed after, end-to-end.

## Round-02 re-review — new agreed fixes (fix-up cycle 2)

All three reviewers re-ran refutation against `git show 7ff7985`. **The round-01 CRITICAL
is confirmed closed** end-to-end against the repo's parsers. They surfaced four new fixes (the
fixes themselves introduced two of them — the value of a real re-review):

| ID | Severity | Raised by | Fix |
|----|----------|-----------|-----|
| **AF6** | MAJOR (antigravity; codex/hermes OQ) | A regression from AF1: `cleanField` was flattening `Detail` (a free-form BODY field — stack traces, logs), destroying multi-line readability for no security benefit (the frontmatter is already closed above it). Fix: `Detail` is no longer `cleanField`'d; it is rendered in a `## Signal detail` section as a markdown **indented block** with newlines preserved (indentation keeps it literal, so it cannot inject a heading). `Source`/`ID`/`Title` (frontmatter + one-line bullet) stay `cleanField`'d. |
| **AF7** | MAJOR (all three) | A liveness hole in AF4: the atomic claim was on the **directory**, but the durable artifact is `00-prompt.md`. A crash/failed write after `os.Mkdir` left an empty `ideas/<slug>/` that future ticks treated as a dedupe hit, **silently suppressing the signal forever**. Fix: the claim is now the PROMPT FILE (`O_CREATE\|O_EXCL` on `00-prompt.md`, `MkdirAll` the dir first); an empty dir from a prior crash is healed (the absent file is re-created), and a failed write removes the partial file so the next tick retries. |
| **AF8** | MINOR (codex + hermes) | `cleanField` did not flatten the Unicode line/paragraph/next-line separators U+2028/U+2029/U+0085. Not a live bypass (the repo's scanners split only on `\n`, verified end-to-end), but the fix's contract ("no line break can inject a key") under-delivered and a future real-YAML-parser swap would reopen the CRITICAL. Fix: `cleanField` now flattens them too. |
| **AF9** | MAJOR (codex) | AF2's digest was only **8 hex / 32 bits**. codex found a CONCRETE collision by birthday search (`probe-55599` and `probe-100565` both → `f3b52266`) and confirmed one candidate suppressed the other — the same dedupe data-loss class, moved to truncated hashing. Fix: the digest is now **32 hex / 128 bits**, making a deliberate second-preimage collision infeasible. |

New regression tests: `TestTickHealsPoisonedEmptyDir` (AF7), `TestTickPreservesMultilineDetail`
(AF6), `TestTickFlattensUnicodeSeparators` (AF8, real U+2028/U+2029 via `string(rune(...))`).

### Dismissed / deferred (round-02)

- **Old-slug migration** (codex F4 / antigravity F3) — N/A: `parley loop` is a brand-new
  command; no loop-drafted candidates exist in any deck, so there is nothing to migrate.
- **Rejected-source summary uses untrusted fields** (hermes NIT F3) — already neutralized:
  `cleanField` strips all `< 0x20` bytes (incl. the ANSI `ESC` 0x1b) from the printed label.
- **`SlugFor` `""→"signal"` fallback "dead code"** (antigravity NIT F4) — kept deliberately:
  `SlugFor` is an exported pure helper; the fallback keeps it total for an unvalidated input
  rather than depending on `Tick`'s upstream `validSources` check.
- **Case-insensitive `Source`** (codex + hermes OQ) — deferred: FINAL.md LE-9 specifies the
  closed set lowercase and signals are machine-generated; `strings.EqualFold` is a friendly
  future polish, not a fix. (DF4.)

## Verification after fix-up cycle 2

`gofmt`, `go build ./...`, `go vet`, `go test -count=1 ./...` (incl. drift guard) — green.
Round-03 re-review requested from all three reviewers.

## Signoffs

(Phase 7 signoffs appended below after the re-review converges to zero agreed fixes.)
