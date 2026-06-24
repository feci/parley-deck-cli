---
agent: hermes-1
idea: automation-outer-loop
review-round: 2
date: 2026-06-24
---

## Summary

Round-02 re-review of the Tier 4 outer-loop fix-up cycle 1 (`git show 7ff7985`,
AF1–AF5). I worked in refutation mode: I assumed every fix was wrong or
incomplete and tried to break it with vectors the round-01 reviewers did NOT
try, plus regression probes. I inspected the CURRENT `internal/loop/loop.go`
(`cleanField`, `validSources`, `dedupeDigest`, `SlugFor`, `Tick`,
`writeCandidate`), `internal/app/loop_cmd.go` (`runLoopTick`), and — critically
— the two downstream frontmatter parsers that decide whether an injection
actually lands (`protocol.ReadFrontmatter` in `internal/protocol/workspace.go`,
`readFrontmatterFieldErr` in `internal/driver/cursor.go`) and the RCE path
(`RunChecks` in `internal/app/driver_impl.go:194-218`, which reads `meta["checks"]`
and runs `sh -c <checks>`).

The round-01 CRITICAL (YAML frontmatter injection via `\n` in a signal field) is
CLOSED against the actual parsers in the repo today. `cleanField` flattens
`\n`/`\r`/`\t`/NUL/VT/FF on Source/ID/Title/Detail, and `validSources` rejects
unknown sources. I confirmed the named CRITICAL vector end-to-end: an `id`
carrying `\nstatus: round-01\nparticipants: [evil]\nchecks: rm -rf /` now lands
as a single-line `source_id:` scalar and `protocol.ReadFrontmatter` parses
exactly one `status: candidate` with no `participants:`/`checks:`.

However `cleanField` has a real, if currently-non-exploiting, gap: it does NOT
flatten the Unicode line/para separators U+2028/U+2029, which ARE YAML line
breaks. Today's parsers (`bufio.ScanLines`, `strings.Split("\n")`) do not split
on them, so the vector does not exploit *now* — but the fix's stated contract
("no newline/control char can inject a YAML key") under-delivers, and a future
swap to a real YAML library reopens the round-01 CRITICAL. I also found a
genuine liveness defect in the AF4 atomic claim: a slug dir left behind by a
crashed/failed write (mkdir ok, prompt write failed) is skipped forever and
never healed, permanently suppressing that signal. AF2, AF3, AF5 hold; no
behavioral regressions.

Verdict: the CRITICAL is closed; AF4 has a residual MAJOR (poisoned-dir
liveness); AF1 has a MINOR latent hardening gap (U+2028/U+2029). Ship-blocker
count: zero CRITICAL.

## Refutation attempts

1. **AF1 — cleanField novel injection vectors (the ones round-01 did NOT try).**
   Built a Go harness exercising `cleanField` against: CR-only (`\r`), CRLF,
   NUL (`\x00`), U+2028 (line separator), U+2029 (paragraph separator), VT
   (`\x0b`), FF (`\x0c`), TAB. Results:
   - CR/CRLF/NUL/VT/FF/TAB → all flattened to spaces. PASS.
   - U+2028 and U+2029 → **NOT flattened**, pass through verbatim**.

   That pass-through is only a finding if a downstream parser splits on those
   codepoints, so I tested the REAL parsers:
   - `protocol.ReadFrontmatter` uses `bufio.NewScanner` with `bufio.ScanLines`.
     `bufio.ScanLines` splits on `\n` (and strips a trailing `\r`); it does
     NOT split on `\r`-only, NUL, or U+2028/U+2029.
   - `readFrontmatterFieldErr` uses `strings.Split(data, "\n")`; same — splits
     on `\n` only, not on U+2028/U+2029.
   So an `id` of `abc\u2028status: round-01\u2028participants: [evil]` lands as
   ONE line (`source_id: abc\u2028status: round-01\u2028participants: [evil]`),
   `strings.Cut(line, ":")` takes the first `:` → key `source_id`, value the
   rest as a literal scalar. I ran the full `Tick` → `protocol.ReadFrontmatter`
   round-trip with this exact vector: parsed `status: candidate`, no
   `participants:`, no `checks:`. **The U+2028 vector does NOT exploit against
   the current parsers.** It is a latent gap, not a live CRITICAL (see Findings
   F1).

2. **AF1 — newline-free YAML injection (no newline needed).** Round-01 was all
   about newlines. I tried a value that injects a key with NO newline: an `id`
   of `abc checks: rm -rf /`. This produces `source_id: abc checks: rm -rf /`
   on ONE line. `strings.Cut(line, ":")` takes the FIRST `:`, so key=`source_id`,
   value=`abc checks: rm -rf /`. The `checks:` token is buried inside the value,
   not a top-level key. `meta["checks"]` is therefore `""` and `RunChecks` never
   runs it. Flow-style `{...}` / anchors/aliases are the same: without a
   newline they are all literal characters inside the first key's value. I could
   NOT force a second top-level key without a line break. This holds because
   both parsers are line-oriented and split on the first `:` per line.

3. **AF1 — validSources bypass.** Probed the validation
   `validSources[strings.TrimSpace(sig.Source)]` with: surrounding whitespace
   (`" commit "`), tab (`"commit\t"`), trailing newline (`"commit\n"`), case
   (`"Commit"`, `"COMMIT"`), trailing junk (`"commitx"`), valid-plus-junk
   (`"commit evil"`), valid-plus-newline-key (`"commit\nstatus: round-01"`),
   empty, space-only. Results:
   - Whitespace/tab/newline-wrapped `"commit"` → ACCEPTED (TrimSpace normalizes).
     This is benign: the accepted value still writes through `cleanField`, and
     the closed set is honored. No bypass to an unknown source.
   - Case variants → REJECTED (`"Commit"` ≠ `"commit"`). No bypass.
   - Trailing junk / valid-plus-junk / valid-plus-newline-key → REJECTED. No
     bypass — a source carrying an injected `status:` key is rejected outright,
     which is exactly the defense-in-depth we want.
   I could NOT get an unknown source accepted. AF1's `validSources` holds.
   (Minor observation: case is rejected, which is stricter than necessary — a
   signals-file author writing `"Commit"` gets a silent reject. Not a security
   issue; see OQ1.)

4. **AF2 — digest collision via the canonical key.** Probed
   `dedupeDigest`'s `strconv.Quote` canonical key for escaping ambiguity:
   colon-shift (`ci:`+`build` vs `ci`+`:build`), trailing-space-in-id, tab-inside.
   All produced DISTINCT keys (`strconv.Quote` quotes each field with double
   quotes and escapes interior quotes, so `"ci:"`+`"build"` ≠ `"ci"`+`":build"`).
   The `TestColonBoundaryNoCollision` and extended `TestSlugFingerprint` tests
   confirm `a/b` ≠ `a:b`. I could NOT force two distinct signals to one digest.
   AF2 holds. One subtle note: the digest is only 8 hex chars of sha256
   (32 bits). Birthday collision is ~65k distinct signals — far beyond any
   realistic signals file, but worth a sentence in OQ2.

5. **AF3 — disabled tick truly inert with malformed signals.** Ran the REAL
   `runLoopTick` (package `app`) with no config (disabled) and a malformed
   `signals.json` (`{not json`). Result: exit 0, stdout `loop tick: disabled ...
   Wrote nothing.`, no `ideas/` dir created. The short-circuit at loop_cmd.go:63
   fires BEFORE `ReadSignals`. AF3 holds — a disabled tick is fully
   cron-safe even with a broken signals file.

6. **AF4 — concurrent double-create / clobber.** Ran 50 goroutines all calling
   `Tick` on the same signal with `-race`: exactly 1 create, 49 skips, exactly 1
   `00-prompt.md` on disk, race detector clean. `os.Mkdir` (ErrExist → skip) +
   `O_CREATE|O_EXCL` on the prompt fully serialize. AF4 holds for the
   happy path. BUT see F2 — the atomic claim has a poisoned-dir liveness hole.

7. **AF4 — poisoned empty dir.** Pre-created `ideas/<slug>/` with NO
   `00-prompt.md`, then ticked the matching signal. Result: `os.Mkdir` returns
   `fs.ErrExist` → `writeCandidate` returns `(false, nil)` → signal is
   `Skipped`, the prompt is NEVER written, and the dir stays empty forever.
   Every future tick of that signal skips it permanently. See F2.

8. **Regressions — did any fix break previously-working behavior?**
   - **cleanField over-strip of multi-line detail:** ticked a signal with a
     legit 3-line `Detail`. cleanField flattens `\n`→space. Content is
     preserved as a single run-on line (`"Line one of a real multi-line issue
     description with three lines."`). Readable, not useless. No functional
     regression, though prose detail loses its structure (OQ3).
   - **source validation rejects a previously-valid signal:** round-01 had no
     source validation, so ANY source was drafted. Now `"Commit"` (case) is
     rejected. This is intended (closed set), not a regression of documented
     behavior, but a signals-file using title-case sources will now silently
     lose candidates (OQ1).
   - **digest change breaks dedupe of an already-drafted candidate:** ticked the
     same signal twice. First tick `Created`, second tick `Skipped` with the
     SAME slug. `dedupeDigest` is deterministic, so dedupe is stable across
     ticks. No regression.

## Findings

### MINOR — F1: cleanField does not flatten U+2028/U+2029 (latent YAML line-break gap)

What: `cleanField` (loop.go:136-143) flattens `\n`, `\r`, `\t`, and `r < 0x20`
to spaces. The Unicode line separator U+2028 and paragraph separator U+2029
are NOT `< 0x20` (they are U+2028/U+2029), so they pass through verbatim. I
confirmed this empirically: `cleanField("abc\u2028status: round-01")` returns
`"abc\u2028status: round-01"` unchanged.

Why it matters (and why MINOR not CRITICAL): The round-01 CRITICAL was "no
newline/control char can inject a YAML key." U+2028/U+2029 ARE YAML 1.1 line
breaks and ARE line breaks in many YAML libraries. cleanField's own doc comment
promises to neutralize "control characters — newlines above all." It does not
cover the Unicode separators. HOWEVER, I verified end-to-end against the ACTUAL
parsers in this repo (`protocol.ReadFrontmatter` via `bufio.ScanLines`, and
`readFrontmatterFieldErr` via `strings.Split("\n")`): neither splits on
U+2028/U+2029, so the vector does NOT exploit today — the injected text stays
inside the `source_id:` scalar value and never becomes a top-level key. This is
a latent hardening gap, not a live bypass: if the frontmatter parser is ever
replaced with a real YAML library (go-yaml, goccy/go-yaml, etc.), U+2028/U+2029
become line breaks and the round-01 CRITICAL reopens. Given the fix's stated
defense-in-depth goal, the gap should be closed now while it is cheap.

Concrete fix: extend the `cleanField` predicate to also flatten U+2028 and
U+2029 (and optionally U+0085 NEL, another YAML 1.1 line break), e.g.:

```go
func cleanField(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' || r < 0x20 ||
			r == '\u2028' || r == '\u2029' || r == '\u0085' {
			return ' '
		}
		return r
	}, strings.TrimSpace(s))
}
```

Add a regression test feeding `\u2028status: round-01` and asserting the parsed
frontmatter has exactly one `status:` (`candidate`) and no `participants:`.

### MAJOR — F2: a failed/crashed write leaves a poisoned empty slug dir that suppresses the signal forever (AF4 liveness hole)

What: `writeCandidate` (loop.go:191-262) claims the slug with `os.Mkdir`; on
`fs.ErrExist` it returns `(false, nil)` → `Tick` records the signal as `Skipped`
and moves on. The dedupe decision is "the directory exists," NOT "the prompt
file exists." If `os.Mkdir` succeeds but the subsequent `O_EXCL` prompt write
fails (disk full, ENOSPC, EIO, a crash between the two syscalls, or an
interrupted process), the empty `ideas/<slug>/` directory is left behind with
no `00-prompt.md`. I confirmed this: pre-creating the empty dir and ticking the
matching signal yields `Skipped` with the prompt never written, and every
future tick of that signal skips it permanently. The signal is silently
swallowed for the lifetime of the deck.

Why it matters: This is the inverse of the TOCTOU AF4 fixed. AF4 closed the
double-create/clobber race but introduced (or kept) a liveness defect: the
atomic claim is on the DIRECTORY, but the meaningful artifact is the PROMPT
FILE. A failed write poisons the slug. Because the loop is the outer discovery
surface (cron/CI/MCP), a transient I/O error on one tick permanently loses a
candidate with no error surfaced — the operator sees `skipped (already
present)` and never knows the candidate was never actually drafted. This is a
silent data-loss path on the very surface whose job is to surface candidates.

Concrete fix: Make the dedupe check and the claim target the same artifact.
Either (a) claim on the PROMPT FILE, not the dir: `O_CREATE|O_EXCL` on
`00-prompt.md` directly (create the dir with `MkdirAll` first; treat an
existing file's `ErrExist` as skip), so a crash leaves no file and the next
tick retries cleanly; or (b) on a prompt-write failure, clean up the empty
slug dir (`os.Remove(dir)` before returning the error) so a future tick is not
poisoned. Option (a) is cleaner and matches "the prompt is the claim." Add a
test: create empty `ideas/<slug>/`, tick the signal, assert it is drafted (or
re-attempted), not skipped.

### NIT — F3: rejected-source summary line is built from untrusted fields and printed to stdout

What: `Tick` builds the `Rejected` entry as
`cleanField(sig.Source)+":"+cleanField(sig.ID)` (loop.go:168), and
`runLoopTick` prints it verbatim (`! <entry> (rejected ...)`). `cleanField`
flattens newlines so it cannot inject into a log line, but the value is still
untrusted-derived and could contain terminal escape sequences (ANSI) or
control-ish bytes >0x7F that a naive terminal might interpret.

Why it matters: Negligible — it is a human-readable summary on stdout, not a
parsed artifact, and `cleanField` already removed `< 0x20`. ANSI escapes
(0x1b) ARE `< 0x20` so they are flattened too. The only residual is arbitrary
printable text in a log line. Flagging only for completeness; no fix required
for the MVP.

## Open questions

1. **Case-sensitivity of `validSources`.** `"Commit"` and `"COMMIT"` are
   rejected. FINAL.md LE-9 lists the set lowercase; is a title-case source in a
   signals file a real scenario (e.g. a CI hook writing `"Issue"`)? If so, a
   case-insensitive compare (`strings.EqualFold`) would be friendlier without
   weakening the closed set. If signals are always machine-generated lowercase,
   current behavior is fine. Not a security issue either way.

2. **8-char digest birthday bound.** `dedupeDigest` returns 8 hex chars
   (32 bits). A collision needs ~65,536 distinct signals (birthday paradox).
   That is far beyond any realistic single-deck signals volume, so not a
   practical concern — but if the loop ever aggregates across many decks or
   long retention, bumping to 12 hex chars (48 bits) is trivial insurance.
   Worth a one-line note in the code comment.

3. **Multi-line `Detail` readability.** cleanField flattens a multi-line
   detail to a single space-joined line. The content survives, but a 10-line
   issue description becomes a run-on paragraph under `## Problem / idea`.
   Acceptable for a candidate prompt (it is a stub to be elaborated on
   promotion), but if detail provenance matters for the human triaging the
   candidate, consider preserving newlines in the BODY section only (frontmatter
   fields stay flat). Out of scope for the security fix; flagging for awareness.
