---
agent: hermes-1
idea: automation-outer-loop
review-round: 3
date: 2026-06-24
---

## Summary

Re-ran refutation against `git show 14f8295` (fix-up cycle 2: AF6–AF9) with the §14
human brake as the named boundary. I treated every round-02 fix as wrong until I
failed to break it, traced the hostile Detail through the repo's REAL frontmatter
parsers (`protocol.ReadFrontmatter` at `internal/protocol/workspace.go:296` and the
driver's `readFrontmatterFieldErr` at `internal/driver/cursor.go:290`, the latter is the
promotion-gate `status` reader), and reproduced the original CRITICAL end-to-end through
`parley init` → `parley loop tick --enable` → `parley status` with a hostile
signals file.

All four round-02 fixes hold. The round-01 CRITICAL (frontmatter injection via
Source/ID newlines) is re-confirmed closed end-to-end. §14 holds (no
run/push/merge/finalize/quorum path reachable from `loop.Tick` / `runLoopTick`).
No CRITICAL, no MAJOR. Two MINOR observations (neither defeats §14; one is a
defense-in-depth suggestion, one is a test-coverage gap), one NIT.

`gofmt`, `go build ./...`, `go vet`, `go test -count=1 ./...` — green (re-run from
clean after removing my throwaway probe tests).

## Refutation attempts

### AF6 — Detail as a 4-space indented block, not cleanField'd

This was the highest-value target: AF6 deliberately REMOVED sanitization from a
free-form body field, so I tried hardest here.

1. Break out into the candidate's OWN frontmatter. The frontmatter closes at the
   second `---` the scanner sees. Because `writeCandidate` emits the real closing
   `---` immediately after the frontmatter scalars (loop.go:245) and ONLY THEN
   renders the Detail block below `## Signal detail`, every Detail line — including
   a line that is exactly `---` — appears AFTER the closing fence. Both parsers stop
   at the second `---`, so no Detail content can become a frontmatter key. Verified
   with a Detail of `"log1\n---\nstatus: round-01\nparticipants: [evil]\n## Injected\n---\ntail"`:
   `ReadFrontmatter` returned `status=candidate`, no `participants`; the driver-style
   reader returned `status=candidate`, no `participants`.
2. Detail that STARTS with `---` (first rendered line `    ---`). Tried
   `"---\nstatus: round-01\nparticipants: [evil]"` and a Detail that is exactly `"---"`.
   `TrimSpace("    ---") == "---"` IS recognized as a fence by both parsers, BUT it is
   the THIRD `---` in the file (opener line 1, real closer line 8, then this one) —
   the parser already returned at the real closer, so it is never seen. Status stayed
   `candidate`, no `participants`. No breakout.
3. Leading blank lines + heading. `indentDetail` only `TrimRight`s newlines
   (loop.go:292), so a Detail `"\n\n## I am a heading\nstatus: round-01"` renders as
   blank indented lines then `    ## I am a heading`. Every content line carries the
   4-space prefix, so no column-0 `## ` heading is produced; the only `## ` lines in
   the file are the five structural sections. Frontmatter parsed clean.
4. The RCE path (`checks:` → `sh -c`, `driver_impl.go:205-209`). The loop's prompt
   carries NO `checks:` frontmatter key, and `ReadFrontmatter` reads only up to the
   closing `---`, so a `checks:` token placed in the Detail body is never read as
   frontmatter. The original latent RCE stays unreachable via the loop.

AF6 holds. The indented-block approach is sound because the frontmatter is closed
above the block, not because indentation is parsed as literal by the repo's scanners
(they don't care about indentation — they only look for a top-level `---` fence, and
the block can't produce one that the scanner reaches).

### AF7 — atomic claim on the prompt file, MkdirAll first, defer-cleanup

1. Two concurrent ticks double-create/clobber. The claim is `O_CREATE|O_EXCL` on
   `00-prompt.md` (loop.go:216). `O_EXCL` is the kernel-atomic create-or-fail; exactly
   one tick wins the fd, the other gets `fs.ErrExist` → `created=false` (skip). No
   clobber. Correct.
2. TOCTOU between `MkdirAll` and `OpenFile`. `MkdirAll` is idempotent and treats an
   already-existing dir as success; a concurrent creator cannot make it fail in a way
   that matters. The real atomicity lives in `O_EXCL`, not in the dir create. No TOCTOU.
3. Defer-cleanup removes a VALID prompt? `wrote` is set true ONLY after a fully
   successful `f.WriteString` (loop.go:281); the defer (loop.go:224-229) removes only
   when `wrote==false`. `*os.File.WriteString` loops internally and returns `err!=nil`
   on a short write, so `err==nil` implies the full prompt was written. A partial write
   errors → `wrote` stays false → `os.Remove` cleans the partial file → next tick
   retries. The defer `Close`s before `Remove`, so the unlink is clean. I could not
   construct a path where a valid, complete prompt is removed.
4. Empty-dir healing. A pre-existing empty `ideas/<slug>/` (simulated crash) is
   healed: `MkdirAll` is a no-op, `OpenFile` creates the missing file.
   `TestTickHealsPoisonedEmptyDir` covers this; I re-ran it green.

AF7 holds.

### AF8 — U+2028/U+2029/U+0085 flattened

YAML 1.1's line-break set is exactly `{LF x0A, CR x0D, NEL x85, LS x2028, PS x2029}`.
`cleanField` (loop.go:145-153) flattens `\n`, `\r`, `\t`, `r < 0x20` (covers 0x00-0x1F,
including 0x0A/0x0D/0x0B/0x0C), plus explicitly `U+2028`, `U+2029`, `U+0085`. `U+0085`
(0x85 = 133) is above 0x20 so it MUST be listed explicitly — and it is. The set is
complete for YAML 1.1. `TestTickFlattensUnicodeSeparators` exercises real U+2028/U+2029;
re-ran green. No remaining YAML 1.1 line break is missed.

### AF9 — 32 hex / 128-bit digest

1. Collision. 128 bits makes a deliberate second-preimage / birthday collision
   infeasible (birthday bound ~2^64 ops). codex's 32-bit collision is gone. I did not
   attempt to find a 128-bit collision (would be infeasible in a review budget); the
   width is the fix.
2. Same-run dedupe still works. Two identical signals in one `Tick` (same source+id)
   produce one slug → one `OpenFile` succeeds, the second gets `ErrExist` →
   `Created=1, Skipped=1`. The digest change did not break in-run dedupe. Verified
   (`TestProbeSameRunDedupe`, run green then removed).
3. Slug length. `loop-` + sanitize(source) + `-` + 32 hex. For `source=commit` that is
   44 chars; for the longest valid source (`manual`) 45. Well under any filesystem
   (255) or path-component limit. `parley status` and `readIdeas` handle it fine
   (verified end-to-end). Acceptable.

AF9 holds.

### Round-01 CRITICAL re-confirm + §14

Reproduced the original vector end-to-end: `parley init` → hostile signals.json with
`id="abc\n---\nstatus: round-01\nparticipants: [evil]\nchecks: rm -rf /"`,
`title="t\nparticipants: [evil2]"`, `detail="log1\n---\nstatus: round-01\nparticipants: [evil3]\n## Injected heading"`
→ `parley loop tick --enable` → `parley status`.

Result: `status=candidate  participants=` (no quorum claim, no round-01 flip). The
`source_id:` line shows the payload flattened to one scalar line
(`abc --- status: round-01 participants: [evil] checks: rm -rf /`). The Detail's
hostile tokens are confined to the indented body block. CRITICAL closed.

§14: `internal/loop/loop.go` imports only stdlib + `internal/protocol` — no
`runner`, `driver`, `runcontrol`, `exec`, `CreateIdea`, push/merge/finalize. `app.go:91`
wires `case "loop"` → `runLoop` → only `tick` → `runLoopTick` → `loop.Tick`, which
only ever writes `00-prompt.md` with `status: candidate`. No run/push/merge/finalize/quorum
path is reachable. §14 holds.

## Findings

### F1 — MINOR — `indentDetail` preserves leading newlines as blank indented lines (cosmetic, not a §14 breach)

`indentDetail` (loop.go:289-301) does `strings.TrimRight(s, "\n")` but no left-trim. A
Detail that begins with `\n\n...` renders as blank `    ` lines before the first real
content line. This is harmless for security (every line is still 4-space indented, no
column-0 heading possible) and harmless for the frontmatter parsers (the block is below
the closing fence), but it produces a slightly ugly `## Signal detail` section with
leading blank indented lines. Not worth a fix cycle on its own; if touched, a
`strings.TrimSpace` before the empty-check (keeping the `(no detail provided)` sentinel)
would tidy it. I do NOT recommend blocking on this.

Concrete fix (optional): in `indentDetail`, replace
`s = strings.TrimRight(s, "\n")` with `s = strings.TrimSpace(s)` — the subsequent
`TrimSpace==""` empty-check is unchanged and the block loses leading/trailing blank
lines. Pure cosmetic.

### F2 — MINOR — no regression test asserts the Detail block cannot inject a REAL markdown heading into the candidate's body

The cycle-2 suite added `TestTickPreservesMultilineDetail` (newlines survive) but no
test asserts the negative: that a Detail carrying `## Foo` / `---` / `status:` cannot
produce a column-0 `## ` heading or a top-level `---` that a future MARKDOWN renderer
(not the current frontmatter scanner) might misread. The security argument today rests
on the frontmatter scanners stopping at the closing `---`, which is correct, but the
AF6 contract as written in consensus.md ("indentation keeps it literal, so it cannot
inject a heading") is a markdown-literal claim that is not directly tested. A small
negative test (assert no line in the body other than the five structural sections
starts with `## ` at column 0, and no body line trims to `---` EXCEPT inside the
indented block) would pin the contract and catch a future regression if someone
"improves" `indentDetail` to strip the prefix from code-fence-aware lines.

Concrete fix: add a test in `internal/loop/loop_test.go` that drafts a candidate with
`Detail: "x\n## evil\n---\nstatus: round-01"` and asserts (a) every line beginning with
`## ` in the file is one of the five structural headings, (b) `ReadFrontmatter` still
returns `status=candidate` with no `participants`. (I wrote and ran exactly this probe
during review; it passed. Promoting it to a permanent test is the suggestion.)

### F3 — NIT — `cleanField` comment says "control characters" but the code flattens `r < 0x20` which is the C0 control block; correct, just phrasing

The comment at loop.go:137-144 says "control characters — newlines above all — are
flattened". The predicate `r < 0x20` covers C0 controls (0x00-0x1F); it does NOT cover
the C1 controls (0x80-0x9F), of which U+0085 (NEL) is one and is handled explicitly.
This is correct and complete for the YAML line-break contract, but a reader might
infer "all control characters" are flattened when in fact the C1 block (minus the three
explicit runes) is not. None of the other C1 runes are YAML 1.1 line breaks, so this is
purely a documentation precision issue. No code change needed; optionally clarify the
comment to "C0 controls plus the YAML 1.1 line breaks U+2028/U+2029/U+0085".

## Open questions

1. (Carry-forward, still deferred — DF1) `protocol.ReadFrontmatter` is last-wins on
   duplicate keys; `readFrontmatterFieldErr` is first-wins. AF1 makes this unreachable
   via the loop (the loop can no longer PRODUCE duplicate keys), and AF6 keeps Detail
   out of frontmatter, so neither parser sees duplicates from a loop candidate. The
   inconsistency remains a pre-existing protocol-parser property for human-authored
   prompts. Still out of scope for this idea; flagged for the protocol/parser follow-up.

2. (AF6 markdown-renderer surface) No repo tool currently RENDERS `00-prompt.md` as
   markdown in a way that gates a §14 action — the gates are frontmatter-scanner based
   (`status`, `participants`, `checks`), all confirmed safe. If a future tool renders
   the prompt body as markdown and acts on headings, the indented-block approach stays
   safe (4-space prefix = code block in CommonMark/Goldmark, not a heading), but that
   contract is currently untested (F2). Worth pinning before such a renderer lands.

3. (AF9 operational) A 128-bit digest makes accidental collision negligible and
   deliberate collision infeasible, but the dedupe is still hash-only — there is no
   human-readable disambiguation in the slug beyond the source hint. Two distinct
   signals from the same source that happen to share a 128-bit prefix will silently
   suppress one another; this is astronomically unlikely but, unlike the
   fingerprint-explicit case, is not human-inspectable in the slug. Acceptable for an
   MVP; noting it as a future "show full digest in the prompt body" polish, not a fix.
