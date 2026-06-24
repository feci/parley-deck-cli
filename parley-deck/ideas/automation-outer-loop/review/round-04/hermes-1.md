---
agent: hermes-1
idea: automation-outer-loop
review-round: 4
date: 2026-06-24
---

## Summary

Round-4 refutation of fix-up cycle 3 (commit d32b082: AF10–AF13 + F3). I treated
every round-03 fix as wrong until I failed to break it, then went after the
threat model itself rather than just re-reading the diff. I read
`review/consensus.md` (AF10–AF13), `git show d32b082`, the CURRENT
`internal/loop/loop.go` (writeCandidate, indentDetail, cleanField, dedupeDigest),
`internal/app/loop_cmd.go`, and the `retro propose` precedent at
`internal/app/retro.go:102-122`.

AF11, AF12, AF13 hold under every separator I could throw at them. The round-01
CRITICAL (Source/ID/Title frontmatter injection) is re-confirmed closed
end-to-end through the compiled CLI. §14 holds (no run/push/merge/finalize/quorum
path reachable from `loop.Tick` / `runLoopTick`; `internal/loop` deps are only
`internal/protocol` + `internal/fsutil`, neither of which execs).

But AF10 does NOT fully close the symlink-escape class it was filed against. I
reproduced, end-to-end through the compiled CLI, that a symlink planted at
`parley-deck/ideas/` (one level ABOVE the slug dir that AF10 guards) makes the
loop write the entire candidate tree — `ideas/<slug>/00-prompt.md` — into the
symlink target, silently, with a success return. This is the same §14
write-boundary breach ("never write outside `parley-deck/ideas/<slug>/`") as the
original AF10 vector codex broke, one directory level up, and AF10's fix does
not cover it. The root cause is `os.MkdirAll(ideasDir, 0o755)` at loop.go:212,
which follows a pre-existing symlink at `ideas/` (MkdirAll is idempotent on an
existing link-to-dir) — the slug-dir `Lstat` guard at loop.go:224 only inspects
`ideas/<slug>`, never `ideas/` itself.

One finding (MAJOR). No CRITICAL, no MINOR, no NIT. The CRITICAL and §14
invariants I was asked to re-confirm are closed/holding; this is a regression in
the AF10 fix's own threat model, not in the original CRITICAL or §14 code-path
boundary.

`gofmt`, `go build ./...`, `go vet`, `go test -count=1 ./internal/loop
./internal/app` — green (re-run from clean after removing my throwaway probe
tests).

## Refutation attempts

### AF10 — symlink-safe slug claim (BROKEN)

This was the highest-value target because AF10 was filed as closing the
"never write outside `parley-deck/ideas/<slug>/`" write boundary, and the fix
only hardens the slug dir, not its parent. I reproduced codex's original
AF10 vector (symlink AT `ideas/<slug>`) — it is correctly refused. Then I
moved the symlink one level up.

1. Symlink at `ideas/<slug>` (codex's original vector). Planted
   `parley-deck/ideas/<slug> → /tmp/target`, ran a matching signal. The
   `os.Mkdir` → `ErrExist` → `Lstat` path at loop.go:220-230 catches it:
   `fi.Mode()&os.ModeSymlink != 0` → error, nothing written to target.
   `TestTickRejectsSymlinkedSlugDir` covers this. AF10 holds FOR THE SLUG DIR.

2. Symlink AT `parley-deck/ideas/` (the parent). This is the breach.
   Reproduced end-to-end through the compiled CLI:
   - `parley init --dir $ROOT`
   - `rm -rf $ROOT/parley-deck/ideas; ln -s /tmp/outside-deck-target $ROOT/parley-deck/ideas`
   - `parley loop tick --dir $ROOT --signals s.json --enable` where s.json has
     `{"source":"commit","id":"ideas-symlink-escape"}`
   - Result: `loop tick: 1 candidate(s) drafted` (exit 0, NO error), and the
     file `00-prompt.md` materialized at
     `/tmp/.../outside-deck-target/loop-commit-fbab0e0121ae81e9ba3512a2a738e322/00-prompt.md`
     — entirely outside `parley-deck/`. The symlink is left intact (no side
     effect), and the loop reports success. A unit probe
     (`TestProbeSymlinkedIdeasDir`) reproduced the same: the prompt appeared in
     the target, not under the deck.

   Root cause: loop.go:212 `os.MkdirAll(ideasDir, 0o755)`. `MkdirAll` treats an
   existing symlink-to-directory as success (it is idempotent and does not
   reject a link at the leaf it is creating toward). Then loop.go:220
   `os.Mkdir(dir, ...)` where `dir = ideasDir/<slug>` creates the slug dir
   THROUGH the already-followed `ideas/` symlink — so the new directory and the
   `O_EXCL` prompt both land in the symlink target. AF10's `Lstat` guard
   (loop.go:224) inspects only `dir` (the slug dir), which is a freshly-created
   REAL directory INSIDE the target — it is not itself a symlink, so the guard
   passes. The parent `ideas/` link is never Lstat'd.

   Why this is the same §14 class AF10 was filed against: the consensus text for
   AF10 says the fix makes the loop "never write outside
   `parley-deck/ideas/<slug>/`" and "a loop must only ever write inside
   `parley-deck/ideas/<slug>/`, never through a planted symlink to somewhere
   else." A symlink at `ideas/` routes every write through it, so the prompt is
   not written inside `parley-deck/` at all. The threat model is "any ancestor
   symlink can redirect the write"; AF10 hardens only the immediate slug leaf.

3. Symlink AT `00-prompt.md` (slug dir real). Planted
   `ideas/<slug>/00-prompt.md → target/00-prompt.md`, ran a matching signal.
   `O_CREATE|O_EXCL` on a symlinked file: on macOS/Linux `open(O_CREAT|O_EXCL)`
   with a trailing-symlink target returns EEXIST (the link exists), so the
   existing symlink is NOT followed and nothing is written to the target; the
   signal is correctly skipped (dedupe). No breach. (Confirmed by
   `TestProbeSymlinkedPromptFile`: res Skipped, target untouched.)

4. Non-directory regular file at `ideas/<slug>`. Planted a regular file, ran a
   matching signal. `os.Mkdir` → `ErrExist` → `Lstat` → `!fi.IsDir()` → error,
   nothing written, the poison file is byte-for-byte unchanged. No leak. No
   partial write. (Confirmed by `TestProbeNonDirAtSlug`.)

5. TOCTOU swap between Mkdir and OpenFile. I could not win a real race in a
   review budget, but I tested the steady-state analogue: a real dir exists at
   `Mkdir`/`Lstat` time, then is swapped to a symlink before `OpenFile`. In the
   current code `Mkdir` returns `ErrExist`, `Lstat` sees the symlink, and the
   path rejects BEFORE `OpenFile` is reached — so the swap is caught at the
   `Lstat` gate, not at `OpenFile`. The residual TOCTOU window is between the
   `Lstat` (loop.go:224) and the `OpenFile` (loop.go:234): if the slug dir is a
   real dir at `Lstat` time and is atomically replaced by a symlink in that
   window, `OpenFile` would create `00-prompt.md` through the new symlink. This
   is a genuine but narrow window (requires winning a race against a single
   tick's own sequential calls) and the ancestor-`ideas/` vector above is a
   strictly easier and deterministic exploit, so I rank the deterministic one as
   the finding and note the TOCTOU as the same root-cause class.

6. "Matches retro propose precedent" claim. I checked `internal/app/retro.go:102-122`.
   `retro propose` does `Lstat(ideaDir)` FIRST and rejects ANY pre-existing entry
   (including a real dir), then `os.Mkdir`. It does NOT let an existing real dir
   fall through (the loop does, for AF7 empty-dir healing). More relevantly, I
   reproduced the `ideas/` symlink against `retro propose` too: it ALSO writes
   through the `ideas/` symlink into the target (retro.go:113
   `MkdirAll(filepath.Dir(ideaDir))` follows it). So the precedent AF10 claims to
   match is itself vulnerable to the same parent-symlink class — the consensus
   "matching retro propose's precedent" framing slightly overstates the
   precedent's coverage. (This is out of scope to fix in `retro`; I flag it only
   to calibrate the finding, not as a retro deliverable.)

### AF11 — indentDetail normalizes U+2028/U+2029/U+0085 + CR to \n (HOLDS)

I tried to get ANY Detail content to column 0 (a heading, `---`, `status:`, or a
frontmatter key).

1. Detail that is ONLY Unicode-separated (no `\n`): `safe<U+2028>## evil<U+2028>status: round-01`,
   and the same for U+2029 and U+0085. After normalization each separator becomes
   `\n`, so every segment gets the 4-space prefix: `    safe` / `    ## evil` /
   `    status: round-01`. No column-0 content. (Probe `TestProbeDetailOnlyUnicodeSep`.)
2. Leading Unicode separator: `<U+2028>## evil heading<U+2029>---<U+0085>status: round-01`.
   `TrimSpace` (AF12) drops the leading separator, then normalization+indent
   covers the rest. Output: `    ## evil heading` / `    ---` / `    status: round-01`. No column-0 content. (Probe `TestProbeDetailLeadingUnicodeSep`.)
3. CR-separated fence: `x\r---\rstatus: round-01`. CR → `\n`, all indented:
   `    x` / `    ---` / `    status: round-01`. Zero column-0 fences in the
   block. (Probe `TestProbeDetailCRFFence`.)
4. Vertical tab `\v` (0x0B) and form feed `\f` (0x0C) in Detail. These are
   `< 0x20` so `cleanField` flattens them for frontmatter fields, but
   `cleanField` is NOT applied to Detail (AF6), and `indentDetail` does NOT
   normalize them. So `safe\v## evil-vt\fsafe-ff\nstatus: round-01` renders as
   `    safe\v## evil-vt\fsafe-ff` (one physical line with the VT/FF mid-line)
   and `    status: round-01`. Because the VT/FF do not split the line in Go's
   `strings.Split(s, "\n")`, the `## evil-vt` stays MID-LINE after `safe\v` — it
   never starts a line, so it is NOT at column 0. No heading injection. This is
   technically incomplete coverage of the AF11 contract ("every logical line
   indented") if a downstream MARKDOWN renderer treats VT/FF as line breaks, but
   no CommonMark/Goldmark renderer I know of does, and the repo's scanners split
   only on `\n`. Not a live bypass; not even a realistic future bypass. I do not
   raise it as a finding — noting it only because I tried it and it did not break.
5. NUL byte (0x00) in Detail: `safe\x00\n## evil`. `indentDetail` does not strip
   NUL, but `\n` still splits, so `    ## evil` is on its own indented line. No
   column-0 content. (Probe `TestProbeDetailNULByte`.)

AF11 holds. The AF13 test (`TestTickDetailCannotInjectHeadingOrFence`) pins the
contract correctly (exactly 2 column-0 fences, no non-structural `## ` heading,
clean frontmatter). I confirmed the full prompt bytes for a hostile Detail and
the frontmatter stayed `status: candidate` with no `participants:`.

### AF12 — TrimSpace (HOLDS)

`indentDetail` now `TrimSpace`s before the empty check (loop.go:314), dropping
leading AND trailing blank lines (AF12 supersedes the round-03 F1 "leading blank
indented lines" cosmetic). A Detail `"\n\n## h"` no longer emits leading blank
indented lines. Cosmetic; verified. Holds.

### AF9/AF2 — digest collision resistance (HOLDS)

`dedupeDigest` (loop.go:106-115) takes 32 hex / 128 bits of sha256 over
`strconv.Quote`'d canonical key. 128 bits makes deliberate second-preimage /
birthday collision infeasible (birthday bound ~2^64 ops). I did not attempt a
128-bit collision (infeasible in review budget). Same-run dedupe still works
(two identical signals → one `O_EXCL` win, one `ErrExist` skip). Slug length
44-45 chars for valid sources, well under path-component limits. Holds.

### Round-01 CRITICAL re-confirm (CLOSED)

Reproduced the original vector end-to-end through the compiled CLI:
`parley init` → hostile signals.json with
`id="abc\n---\nstatus: round-01\nparticipants: [evil]\nchecks: rm -rf /"`,
`title="t\nparticipants: [evil2]"`, `detail="log1\n---\nstatus: round-01\nparticipants: [evil3]\n## Injected heading"`
→ `parley loop tick --enable` → `parley status`.

Result: `status=candidate  participants=` (no quorum claim, no round-01 flip).
The `source_id:` line shows the payload flattened to one scalar:
`abc --- status: round-01 participants: [evil] checks: rm -rf /`. The Detail's
hostile tokens are confined to the indented body block below the closing `---`.
No `checks:` frontmatter key is produced, so the latent RCE path
(`driver_impl.go` `checks:` → `sh -c`) stays unreachable via the loop. CRITICAL
closed. `cleanField` (loop.go:145-153) flattens `\n \r \t r<0x20` plus
U+2028/U+2029/U+0085; `validSources` rejects unknown sources. Both hold.

### §14 — no run/push/merge/finalize/quorum path (HOLDS)

`go list -deps ./internal/loop` → only `internal/protocol` and `internal/fsutil`
internal deps; `internal/fsutil` is a leaf with no exec/runner/driver imports.
`internal/loop/loop.go` imports only stdlib + `internal/protocol`. `app.go:91`
wires `case "loop"` → `runLoop` → `runLoopTick` → `loop.Tick`, which only ever
writes `00-prompt.md` with `status: candidate`. No run/push/merge/finalize/quorum
path is reachable. §14 holds — the code-ACTION boundary is intact. (The AF10
finding above is a code-OUTPUT/path boundary issue, not an action-path issue.)

## Findings

### [MAJOR] AF10 does not guard `parley-deck/ideas/` — a parent symlink silently writes the candidate outside the deck

What is wrong: AF10's fix (loop.go:216-232) hardens only the slug leaf
`ideas/<slug>`: `os.Mkdir` it, and on `ErrExist` `Lstat` and reject a
symlink/non-dir. It does NOT guard the parent `parley-deck/ideas/`. The line
above it (loop.go:212) is `os.MkdirAll(ideasDir, 0o755)`, which follows a
pre-existing symlink at `ideas/` (it is idempotent on an existing
symlink-to-dir and does not reject it). With `ideas/` as a symlink, the
subsequent `os.Mkdir(ideasDir/<slug>)` and the `O_EXCL` prompt open both create
through the symlink, so the entire candidate tree lands in the symlink target,
outside `parley-deck/`.

I reproduced this end-to-end through the compiled CLI: a planted
`parley-deck/ideas → /tmp/outside-deck-target` made `parley loop tick --enable`
report `1 candidate(s) drafted` (exit 0, no error) while writing
`00-prompt.md` into `/tmp/.../outside-deck-target/loop-commit-.../`. The symlink
is left intact and the loop reports success — a silent breach. A unit probe
reproduced the same deterministically.

Why it matters: AF10 was filed against the §14 write boundary "a loop must only
ever write inside `parley-deck/ideas/<slug>/`, never through a planted symlink
to somewhere else" (consensus.md AF10). A symlink at `ideas/` routes every
write through it, so the prompt is written outside `parley-deck/` entirely.
This is the same threat class and the same boundary as codex's original AF10
vector — one directory level up. The fix's own comment ("never follow a
pre-existing symlink for the slug dir") describes a leaf-only guard, but the
boundary it is supposed to enforce is ancestor-transitive: any symlinked
ancestor redirects the write. The slug-leaf `Lstat` cannot catch a parent link
because the slug dir it inspects is a freshly-created real dir INSIDE the
target (not itself a link).

Concrete fix: do not trust `MkdirAll` for the `ideas/` parent either. Before
the `MkdirAll(ideasDir)`, `Lstat` the `ideasDir` path (relative to the deck
root) and reject if it is a symlink; if it does not yet exist, `MkdirAll` only
up to `deck` (or `deck/parley-deck`) and then `os.Mkdir` the `ideas` dir
itself so a leaf symlink cannot be created through. Equivalently and more
robustly, after resolving all writes, verify the final `promptPath` is still
inside the deck root with `filepath.Rel`/prefix check on the cleaned
(`filepath.EvalSymlinks`) path, and reject if it escapes — a single
containment check that covers ANY ancestor symlink at any depth, not just the
two levels hardcoded today. Add a regression test mirroring the existing
`TestTickRejectsSymlinkedSlugDir` but planting the symlink at `ideas/`
(`parley-deck/ideas → /tmp/target`) and asserting: (a) `Tick` returns an error
(or at minimum does NOT report `Created`), and (b) no file appears under the
target. (I wrote and ran exactly this probe during review; it failed against
the current code, confirming the breach. Promoting it to a permanent test is
part of the fix.)

Note on scope/severity: This is NOT the round-01 CRITICAL (frontmatter
injection), which I re-confirmed closed. It is NOT an action-path breach (no
run/push/merge/finalize/quorum becomes reachable — §14's code-action boundary
holds). It is a regression in the AF10 fix's own threat model: AF10 was
accepted as closing "never write outside `parley-deck/ideas/<slug>/`" and it
does not. I rate MAJOR rather than CRITICAL because exploiting it requires an
attacker who can already plant a symlink inside `parley-deck/` (the same trust
precondition as the original AF10 vector codex broke, which was also rated
MAJOR), and because the breached boundary is the write-path containment AF10
names, not the frontmatter/quorum invariant that was the round-01 CRITICAL.

## Open questions

1. (AF10 containment depth) Is a single `filepath.EvalSymlinks` +
   `filepath.Rel(deck, realPromptPath)` containment check preferred over
   per-level `Lstat` guards? It would cover ancestor symlinks at any depth in
   one place (including a future `parley-deck/` symlink) and is easier to
   reason about than hardcoding the two levels that happened to be exploited so
   far. The per-level approach matches `retro propose`'s style but, as shown
   above, `retro propose` is itself vulnerable to the `ideas/` parent symlink —
   so the precedent is not a complete reference. A containment check would also
   let the loop adopt `retro propose`'s stricter "reject ANY pre-existing
   `ideas/<slug>` entry" without losing AF7 empty-dir healing (heal only when
   the contained real dir is empty).

2. (Carry-forward, still deferred — DF1) `protocol.ReadFrontmatter`
   last-wins vs `readFrontmatterFieldErr` first-wins on duplicate keys. AF1
   keeps the loop from producing duplicates, so it is unreachable via the loop;
   remains a pre-existing protocol-parser property for human-authored prompts.
   Still out of scope for this idea.

3. (Carry-forward, still deferred — DF4) Case-insensitive `Source`. FINAL.md
   LE-9 specifies the closed set lowercase and signals are machine-generated;
   `strings.EqualFold` is friendly future polish, not a fix.
