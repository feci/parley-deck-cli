---
agent: hermes-1
idea: automation-outer-loop
review-round: 5
date: 2026-06-24
---

## Summary

I re-reviewed fix-up cycle 4 in refutation mode against `git show 580600a`, the current
`internal/loop/loop.go` (`safeMkdir`, `assertInsideDeck`, `writeCandidate`, `indentDetail`,
`cleanField`, `dedupeDigest`), and `internal/app/loop_cmd.go`. I assumed the round-04 fixes
(AF14/AF15) were WRONG or incomplete until I failed to break them.

I found no remaining issues: 0 CRITICAL, 0 MAJOR, 0 MINOR, 0 NIT. AF14 held under every
symlink/TOCTOU/Rel attack I could mount — the symlink-escape class is closed at the slug leaf
(round-03 AF10), the `ideas/` parent (round-04 AF14), AND at arbitrary ancestor depth via the
depth-complete `assertInsideDeck` containment check. AF15 held under an exhaustive separator
sweep. The round-01 CRITICAL (frontmatter/quorum injection) remains closed, the digest is
confirmed 128-bit, and §14 holds (no run/push/merge/finalize/quorum path).

I wrote throwaway refutation probes directly against the loop package's exported `Tick` and the
unexported `assertInsideDeck` (same-package tests), ran them green, then deleted them so they do
not pollute the repo. `go build ./...`, `go vet`, and `go test -count=1 ./...` (all packages,
including the embedded-default drift guard) are green.

## Refutation attempts

### AF14 — symlink escape (the named MAJOR class)

I tried hard to write a candidate prompt outside `parley-deck/ideas/<slug>/` via every
symlink/ancestor trick I could construct, calling the real `Tick` (not a mock) on a hostile
filesystem:

1. **Grandparent symlink** — `<tmp>/parley-deck` itself is a symlink to a real backing dir; the
   real `ideas/<slug>/` lives under the backing dir. `Tick` followed it (correct — the deck IS
   the symlink) and the prompt landed INSIDE the resolved real deck. `assertInsideDeck`
   canonicalized both sides via `EvalSymlinks` and the `Rel` check passed because the resolved
   slug dir is genuinely inside the resolved deck. No escape. (Probe A — PASS.)

2. **Ancestor-chain symlink two levels up** — `<tmp>/link -> <realroot>`, deck =
   `<tmp>/link/parley-deck`. Same outcome: contained inside the resolved real deck. The
   containment is depth-complete by construction. (Probe C — PASS.)

3. **`ideas/` parent symlink with a pre-existing `<slug>` subdir inside the target** — a
   `parley-deck/ideas -> <outside>` symlink where `<outside>/<slug>` already exists as a real
   dir. `safeMkdir(ideasDir)` hits `ErrExist`, `Lstat` sees the symlink, and rejects with
   "refusing ... a symlink or non-directory" BEFORE the slug `safeMkdir` or `assertInsideDeck`
   ever run. Nothing written to the target. (Probe B — PASS.)

4. **Dangling symlink at `ideas/`** — `parley-deck/ideas -> <nonexistent>`. `safeMkdir(ideasDir)`:
   `os.Mkdir` returns `ErrExist` (the symlink entry exists), `Lstat` sees `ModeSymlink`, rejects.
   `EvalSymlinks` is never reached, so its dangling-symlink error path is not the line that saves
   us here, but the result is correct: rejected, nothing written. (Probe D — PASS.)

5. **`assertInsideDeck` Rel edge cases, called directly:**
   - `rel == "."` (dir resolves to EXACTLY the deck): accepted. Benign — in the real `writeCandidate`
     call path `dir = ideas/<slug>` is always strictly below `deck`, so `rel` is never `"."`. Not
     an escape. (PASS.)
   - **Sibling prefix confusion** (`parley-deck` vs `parley-deck-evil` in the same parent):
     `Rel` returns `../parley-deck-evil`, the `..` prefix check rejects it. (PASS — this is the
     reason the guard checks a `..` PREFIX, not a string prefix, and it works.)
   - **Nested same-name** (`deck/ideas/parley-deck`): `rel = ideas/parley-deck`, no `..`, accepted.
     Correct — it IS inside. (PASS.)
   - **Lexical `..` that resolves inside** (`deck/sibling/../ideas/slug`, all components real):
     `EvalSymlinks` canonicalizes away the `..`, `Rel` yields `ideas/slug`, accepted. (PASS.)
   - **Lexical `..` that escapes** (`deck/ideas/../.. / outside`): resolves outside, `Rel` has a
     `..` prefix, rejected. (PASS.)
   - **Real call path** (`deck/ideas/slug`): accepted. (PASS.)

6. **TOCTOU swap between `assertInsideDeck` and `os.OpenFile`** — there genuinely IS a non-atomic
   window: `assertInsideDeck` resolves and approves the slug dir, then `OpenFile(promptPath)`
   follows whatever the path points to *now*. A concurrent attacker who swaps `ideas/<slug>` for a
   symlink in that window could redirect the write. I could not turn this into a finding because
   exploiting it requires the attacker to already have write access under `parley-deck/`, and
   `parley-deck/` is the trusted workspace root (the `--dir` the operator pointed the loop at). An
   attacker with write access there can already plant anything; the §14 boundary the loop enforces
   is "don't let an untrusted SIGNAL write outside the deck," and a signal cannot win a filesystem
   race. I note it as an observation, not a finding — closing it (e.g. `O_NOFOLLOW` on the open, or
   re-resolving inside the open) would be a defense-in-depth nicety, not a fix for a reachable
   breach under the documented threat model.

The symlink-escape class is closed. Both per-level `safeMkdir` guards (leaf + parent) and the
depth-complete `assertInsideDeck` agree, and each catches cases the other doesn't need to.

### AF15 — Detail column-0 injection under a broad line splitter

I swept every separator the prompt names and a few it doesn't, each injected into `Detail`
wrapped around `## evil`, `---`, `status:`, `participants:`, `checks:` payloads, then asserted
no column-0 heading/fence/key leaked and exactly 2 fences + 1 `status:` (candidate) + 0
`participants:` survive:

- U+0085 (NEL) ALONE — PASS.
- `\v` (vertical tab) + `\f` (form feed) — PASS.
- U+2028 (line sep) + U+2029 (para sep) — PASS.
- U+001C / U+001D / U+001E (C0 info separators) — PASS.
- Lone CR (no LF) — PASS (normalized to `\n`, then indented).
- **ALL mixed** in one Detail (`\n \r \r\n \v \f \x1c \x1d \x1e` + U+0085 + U+2028 + U+2029,
  followed by `## evil\n---\nstatus: round-01\nparticipants: [evil]`) — PASS: 2 fences, 1 status,
  0 participants.
- U+001F (unit separator) — NOT a Python `splitlines` / CommonMark / YAML boundary, and
  `indentDetail` maps it (`< 0x20`) to a space. Confirmed it does NOT create a line break and the
  `## noBreakHere` token stays on one indented line. Correct: `indentDetail` treats only the
  line-break-like C0 set as breaks and collapses the rest to spaces, so no unsplit-but-structure
  byte survives at column 0.

`indentDetail`'s normalization set is `{ \n, \r (incl. CRLF via pre-replace), \v, \f, 0x1c, 0x1d,
0x1e, 0x85, 0x2028, 0x2029 }` → `\n`; everything else `< 0x20` → space; `\t` kept. That is exactly
the union of line boundaries honored by LF scanners, CommonMark, YAML 1.1, and Python
`splitlines`. I could not find a separator that reaches column 0.

### Round-01 CRITICAL (frontmatter injection) — re-confirmed closed

Re-ran the canonical hostile signal (`ID` carrying
`\nstatus: round-01\nparticipants: [evil]\nchecks: rm -rf /`, `Title` carrying
`\nparticipants: [evil2]`, `Source: commit` so it isn't rejected outright). `cleanField` flattens
all C0 + `\r`/`\t` + U+2028/U+2029/U+0085 on `Source`/`ID`/`Title` before they hit the frontmatter
and the one-line bullet. Result: exactly 1 `status:` (candidate), 0 `participants:`, 0 `checks:`.
The injection vector is closed at the value level, and `validSources` closes it again at the
source-set level (an unknown source is rejected, not normalized). Still closed.

### 128-bit digest — re-confirmed

`dedupeDigest` returns `hex.EncodeToString(sum[:])[:32]` — 32 hex chars = 128 bits of sha256. I
asserted the slug suffix is exactly 32 lowercase hex chars. The round-02 32-bit birthday collision
(AF9) is gone; a deliberate second-preimage against 128 bits is infeasible. The canonical key uses
`strconv.Quote` so field boundaries (`a/b` vs `a:b`, `ci:`+`build` vs `ci`+`:build`) are
unforgeable — these are pinned by existing `TestSlugFingerprint` / `TestColonBoundaryNoCollision`.

### §14 — no run/push/merge/finalize/quorum path

Searched `internal/loop` and `internal/app/loop_cmd.go` for
`run|push|merge|finalize|quorum|participant|promote|round-01|deliberat`. Every production-code hit
is either a doc comment asserting the NEGATIVE ("never staffs a quorum, runs, pushes, merges, or
finalizes") or the rendered `## Promotion` NOTE that DESCRIBES the human gate in prose ("a human or
the manifest sets `participants:` ... and flips status to round-01") without performing any of it.
No code path calls run/push/merge/finalize/quorum. `runLoopTick`'s `--enable` only flips
`cfg.Enabled = true`; it still routes exclusively through `loop.Tick`, which only ever drafts
`status: candidate` prompts via `writeCandidate`. There is no promote/run/merge/finalize/quorum
code path in the loop subsystem. §14 holds.

## Findings

0 CRITICAL, 0 MAJOR, 0 MINOR, 0 NIT.

I found nothing to fix. After stating what I tried above:

- The symlink-escape class is closed at three layers (leaf `safeMkdir`, parent `safeMkdir`,
  depth-complete `assertInsideDeck`) and I could not escape `parley-deck/ideas/<slug>/` via a
  grandparent, parent, leaf, dangling, or chain symlink, nor via a Rel edge case (`.` vs `..`,
  sibling-prefix, lexical dotdot). The one theoretical residual — a TOCTOU swap between
  `assertInsideDeck` and `OpenFile` — is not a reachable breach under the threat model (it needs
  write access to the trusted deck root, which is outside the "untrusted signal" boundary the loop
  enforces).
- AF15's separator normalization is exhaustive over the line-break-like set across LF/CommonMark/
  YAML 1.1/Python `splitlines`; no Detail token reaches column 0.
- The round-01 CRITICAL is still closed end-to-end.
- The digest is 128-bit and collision-resistant.
- §14 holds: the loop drafts candidates only; no run/push/merge/finalize/quorum path exists.

This review has now run four re-review rounds and is converged. I agree with the round-04
calibration note that the review is converging; from the refutation lens I have no remaining
attacks that produce a write outside the deck or a structure-injecting Detail.

## Open questions

1. **TOCTOU hardening (defense-in-depth, not a finding).** The window between
   `assertInsideDeck(deck, dir)` and `os.OpenFile(promptPath, O_EXCL)` is not atomic. An
   `O_NOFOLLOW`-style open (or re-resolving `promptPath` immediately before the exclusive create)
   would close it for the leaf without relying on the threat model. Not blocking — the loop's
   boundary is untrusted-signal → deck, and a signal cannot mutate the filesystem — but it would
   make the containment argument independent of "the operator trusts everything under `--dir`."

2. **DF5 (carries forward, out of scope).** `retro propose` shares the `ideas/` parent-symlink
   class the loop now defends against via `assertInsideDeck`. Hardening `retro` belongs to its own
   idea; the loop is now stricter than the precedent it once cited.
