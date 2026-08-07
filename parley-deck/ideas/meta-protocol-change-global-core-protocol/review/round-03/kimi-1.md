---
agent: kimi-1
idea: meta-protocol-change-global-core-protocol
review-round: 3
date: 2026-08-07
reviewed-commit: 4a5c447
verdict: FINDINGS
---

## Summary

Cycle 2 is largely real and I verified it behaviorally against the built binary: all three builds,
vet, and the full suite are green; the §7 "Not yet in force" paragraph is present and byte-identical
in all three copies (skill copy: committed at 455aafe plus this cycle's edit, uncommitted in the
working tree, as the brief states); codex-1's demonstrated droppedContent fixture now reports
correctly; the store's own three symlink components are refused on publish; padded versions are
rejected; CRLF cores render clean in both deck pairings; the stamp false-report is gone; a failed
publish no longer bricks the version (verified with `ulimit -f 0` — republish succeeds). The two
tests hermes-1 asked for exist and genuinely fail on revert; my vacuous traversal test is fixed and
now fails on revert.

Two things keep this from CLEAN:

1. **The claimed "per-section" droppedContent is not per-section.** `renderCounts` is a
   document-global multiset (render.go:205-210); a deck line that exists ONLY under
   `## Requires explicit user approval` is still silently erased when the core carries the identical
   line under `## Allowed automatically` — preview AND apply report nothing, and the deck loses the
   binding restriction (demonstrated end-to-end below). That is codex-1's CRITICAL mechanism, still
   firing, and three cycle-2 texts now assert it cannot: render.go:199-200's comment, and
   IMPLEMENTATION.md:131-132's "so semantic erasure cannot hide behind a coincidental match". The
   fix is also pinned by zero tests — reverting it leaves all 18 protocol tests green.
2. **The honesty fixes again landed beside some of the text they were meant to correct.** The
   changelog still names `DETECTED-UNATTRIBUTED` part of "the shipped guarantee" (0 Go references);
   the new §7 paragraph itself claims "no agent-accessible code path writes a release" two sentences
   after admitting a pty-allocating agent can publish; and `check`'s report preamble still
   contradicts its entries (only render's was rewritten).

Command output (this machine, 4a5c447, clean tree):

```
$ go build ./...                  EXIT=0
$ GOOS=windows go build ./...     EXIT=0
$ GOOS=linux go build ./...       EXIT=0
$ go vet ./...                    EXIT=0
$ go test ./... -count=1          all 26 packages ok (incl. internal/runner 10.202s), EXIT=0
$ go test ./internal/app -run 'TestProtocol|TestPublish|TestCore|TestLoad' -count=1 -v
  18 tests, 18 PASS
```

The flaky `internal/runner` failure codex-1 saw at 8888e00 did not appear here either.

Counts: 1 MAJOR, 4 MINOR, 4 NIT.

## Round-02 findings: fixed or not

1. **[MAJOR hermes-1 + kimi-1] §7 shipped unimplemented guarantees — FIXED in all three §7 copies;
   the changelog instance is NOT fixed.** The deck copy (parley-deck/COOPERATION.md:768-773), the
   embedded default (:759-764), and the skill working tree carry the identical "Not yet in force —
   do not rely on it" paragraph; pinning and `DETECTED-UNATTRIBUTED` are now correctly labelled
   "ratified but not implemented" (diff-verified deck==embedded==skill; the drift test pins the two
   repo copies). BUT `parley-deck/meta/protocol-changelog.md:23-25` is untouched and still says
   "the shipped guarantee is: write-once releases, an attended TTY-gated publisher, no
   agent-accessible write path, and detection with `DETECTED-UNATTRIBUTED` for anything else" —
   hermes-1's round-02 MAJOR named exactly this line, `DETECTED-UNATTRIBUTED` still has zero Go
   references, and the changelog now flatly contradicts the §7 text shipped in the same commit.
   And the new §7 paragraph contains a fresh false absolute of its own — new finding 2.

2. **[CRITICAL codex-1] droppedContent section context + multiplicity — PARTIALLY FIXED; the
   erasure class still fires (new finding 1).** Verified against the real binary, release
   `## Allowed automatically / - Deploy production.` / `## Requires explicit user approval /
   - Delete the database.`:
   - Deck carrying `- Deploy production.` under BOTH sections (codex-1's exact round-02 fixture):
     now reported — `## Requires explicit user approval — 1 line not carried forward`. FIXED.
   - Deck carrying three identical local lines: reported as "4 lines not carried forward"
     (multiplicity counts). FIXED.
   - Deck carrying `- Deploy production.` ONLY under `## Requires explicit user approval`:
     **preview prints no removal section at all; `--yes` writes the deck with the line present only
     under `## Allowed automatically`.** The render's single global copy is consumed by the deck's
     single copy. NOT FIXED — the multiplicity half landed, the section half did not.

3. **[MAJOR codex-1 + kimi-1] symlink hardening — FIXED for the store's own components on Publish;
   the Load read-path is NOT fixed (new finding 3).** Through a real pty publish: `<home>`,
   `<home>/protocol`, and `<home>/protocol/core` as symlinks are all refused
   ("protocolcore: … is a symlink; the core store must not be reached through one", exit 1, outside
   dirs stayed empty). An ancestor symlink ABOVE the home still redirects the store (published into
   the symlink target) — that is the documented, deliberate scope (macOS `/var`), accurately
   described in core.go:192-201 and IMPLEMENTATION.md, so not a finding. But a symlink at
   `<store>/<version>` on the READ path — codex-1's round-02 PRIMARY #1 — is unchanged:
   `Load` (core.go:86-99) never calls the new check, `render --yes` with a lock pinning `safe-1`
   where `<store>/safe-1 → /tmp/.../outside` wrote `SYMLINK-ESCAPED-CONTENT` into the deck, exit 0.
   No test pins the new check at all (revert leaves the suite green — see test-quality).

4. **[MINOR codex-1] ValidVersion trimmed-copy validation — FIXED in code, unpinned in tests.**
   Real pty publish of `--version ' padded '` now fails `protocolcore: unsafe version " padded "`
   and creates nothing. But reverting to the 8888e00 trimming shape leaves all 18 tests green, and
   the new test comment (protocol_test.go:354-356) says untrimmed input is covered "at the API
   level … see TestPublishRejectsUnsafeVersions" — that test contains no padded case.

5. **[MINOR kimi-1] vacuous traversal test — FIXED, now real.** In a scratch tree with Load's
   `ValidVersion` reverted: `TestProtocolRejectsPathTraversalInTheLock` FAILS on its new reason
   assertion ("rejected for the wrong reason: … release not installed"), and the new
   `TestLoadRefusesToEscapeTheStore` FAILS with `Load escaped the store and read "PLANTED"`.

6. **[MAJOR hermes-1] TTY-refusal and status read-error e2e tests — FIXED, both real.**
   `TestPublishRefusesWithoutATerminal` drives `runProtocol`, asserts no release was written, and
   asserts the refusal text does not overclaim; deleting the `platformHasTTY()` gate in scratch
   makes it FAIL ("published without a controlling terminal"). `TestProtocolStatusReportsReadErrors`
   (lock path replaced by a directory) FAILS when the `pinErr` branch is reverted — status exits 0
   showing `deck pins : —`.

7. **[MINOR kimi-1] CRLF cores + stamp false-report — FIXED behaviorally; one site of the preamble
   fix missed; both fixes untested.** CRLF core + LF deck: uniform LF output (0 CRLF, 0 lone-CR),
   second render "already matches", `check` exit 0. CRLF core + CRLF deck: uniform CRLF (12 CRLF,
   0 lone LF, **0 CRCRLF**), idempotent, `check` exit 0. Version bump 1.0.0→1.0.1 with a
   byte-identical body: no more `(document header) — 1 line not carried forward`. Render's preamble
   is fixed, but **`check`'s was not** (protocol.go:268): it still prints "sections present here but
   not in the core:" directly above "## 3. Phases — 1 line not carried forward" for a section that
   IS in the core — the exact contradiction my round-02 finding named at both sites. Reverting the
   CRLF-core normalization or the stamp skip in scratch leaves all 18 tests green.

**Still open from earlier rounds, not claimed by this fix-up** (restating so they are not lost):
`render` still swallows deck-file read errors (protocol.go:180 — `prior, _ := os.ReadFile(path)`);
`docs/cli-reference.md` and the root `CHANGELOG.md` still never mention the `protocol` command;
preflight.go:584-596 still owns the legacy `**Protocol synced:**` stamp format (latent until
adoption); the mode-preserving write branch and `check --json`'s non-zero exit remain untested; the
O_NOFOLLOW half of `TestPublishRefusesExistingReleaseDirAndSymlinks` remains unreachable dead code
(round-02 MINOR); IMPLEMENTATION.md's frontmatter still says `status: implemented` /
`head-commit: (see release commit)` and neither fix-up section carries a SHA (codex-1's round-02
MINOR); "GOOS=windows and GOOS=linux builds are now part of the check" (IMPLEMENTATION.md:69) —
there is still no check: `grep -rn GOOS scripts/ .github/` → 0 matches; the source comment
codex-1 flagged (protocol.go:277-279) is unchanged.

## New findings (by severity)

### [MAJOR] Cross-section silent erasure still fires, and three cycle-2 texts now claim it can't

`droppedContent`'s "per-section" claim is not implemented: `renderCounts` is built from the WHOLE
rendered body (render.go:205-210), so a deck line counts as carried whenever the render holds an
unconsumed copy ANYWHERE. Demonstrated against the real binary at 4a5c447 (fixture in finding 2
above): a deck whose ONLY copy of `- Deploy production.` sits under `## Requires explicit user
approval`, rendered against a core carrying the identical line under `## Allowed automatically` —
preview reports nothing, apply removes the line from the approval section, exit 0, and a following
`check` certifies the result in sync. The deck's binding restriction is reclassified as allowed with
zero report: the semantic-inversion case from codex-1's CRITICAL, still reachable whenever the
render's global count of a line ≥ the deck's count in the dropped section.

What makes this MAJOR rather than a residual note is that the cycle now asserts the property it
didn't implement, in the G7b-honesty fix itself:

- render.go:199-200: "So the comparison is per-section and multiplicity-aware: **within the section
  a line belongs to**, a deck line counts as carried only if the render still has an unconsumed copy
  of it." — the code does not do this.
- IMPLEMENTATION.md:131-132: "Now per-section and multiplicity-aware (consume-a-copy), **so semantic
  erasure cannot hide behind a coincidental match**." — it can; shown above.
- The §7 G1 posture ("report every block replaced or removed", render.go:31-33) still rests on this
  report.

The fix is straightforward in shape: key `renderCounts` by (heading-path, line) — walk the render
with the same `heading()` tracker the deck walk already uses — so a copy only satisfies the section
it sits in. Whichever way it's fixed, the texts above must match the code.

### [MINOR] The new §7 "in force" sentence overclaims in all three copies

parley-deck/COOPERATION.md:772-773 (and the byte-identical embedded + skill copies): "What IS in
force today: the core store is write-once, `parley protocol publish` is attended-only, and **no
agent-accessible code path writes a release**." Two sentences earlier the same paragraph admits the
refusal "does not stop an agent that allocates a pty". I published through a pty wrapper again this
cycle (`Published core 9.9.9 …`, exit 0): for a pty-allocating agent there IS an accessible code
path that writes a release — `parley protocol publish` itself. The sentence whose entire job is
enumerating what is in force lists one thing that is not, as stated. One clause fixes it ("no code
path an ordinary agent run can reach writes a release" — which is exactly what the refusal text
already says honestly).

### [MINOR] Load still follows a symlinked release directory

codex-1's round-02 PRIMARY #1, reproduced unchanged at 4a5c447: lock `core-version: safe-1`,
`<store>/safe-1` a symlink to an outside directory holding `COOPERATION.md` → `render --yes` exits 0
and writes the outside bytes into the deck. `assertNoSymlinkComponents` is called only from
`Publish` (core.go:138) and covers only the three store-owned components — the version directory is
outside its scope by construction, and `Load` does no Lstat of its own. Impact is bounded: an
attacker who can plant that symlink already has the local write access to drop a whole fake release
directory (the lock pins no hash, so D8's missing-release block is the only barrier either way) —
hence MINOR, same as the equivalent residual I rated in round 02. But the round-02 finding asked for
the same rule on Load, and neither the fix nor the fix-up's text mentions the gap.

### [MINOR] The changelog still ships `DETECTED-UNATTRIBUTED`; the commit's own texts now disagree

protocol-changelog.md:23-25 still presents "detection with `DETECTED-UNATTRIBUTED`" as part of "the
shipped guarantee". Zero Go files reference it; §7 (this same commit) now says it is "ratified but
not implemented". A reader of the changelog — the document G6 exists to keep accurate — gets the
pre-fix claim. This is the un-fixed half of round-02 finding 1: the fix-up corrected the three §7
copies named in its bullet and left the fourth carrier of the same claim.

### [NIT] "Nineteen protocol tests" does not reproduce

IMPLEMENTATION.md:159. `protocol_test.go` holds 18 tests (15 at 8888e00 + 3 new); cycle 1's
"Fifteen tests" counted that file. 19 only appears if `TestEmbeddedDefaultMatchesLiveDeck` is now
counted — a convention cycle 1 didn't use. Trivial, but this is the repo that keeps getting bitten
by unverified numbers.

### [NIT] The stamp skip silently drops genuine deck prose that starts with the stamp prefix

droppedContent now skips any deck line with the `**Protocol synced:**` prefix (render.go:226-230).
A deck line of real prose beginning with that prefix (e.g. a runbook documenting the marker) is
dropped with no report — verified: report covers only the section heading, "1 line not carried
forward", while two lines vanish. Renderer-owned prefix, so the collision is unlikely; a shape check
(`**Protocol synced:** core `) would close it.

### [NIT] The traversal test's new comment cites coverage that doesn't exist

protocol_test.go:354-356 says untrimmed input "is rejected at the API level instead; see
TestPublishRejectsUnsafeVersions" — that test's list is `""`, `"."`, `".."`, `"../x"`, `"a/b"`,
`".hidden"`. No padded case; the trim-reject is unpinned (revert-verified green).

### [NIT] The refusal text's "hash detection" overstates what exists

protocol.go:295: "The durable guarantees today are write-once releases and hash detection." What
exists is deck-VIEW drift detection (`check`); release-tamper detection is `DETECTED-UNATTRIBUTED`,
which §7 now correctly labels not implemented. "Hash detection" sits between the two. Borderline,
one word from honest.

## Test-quality assessment

18 protocol tests (+1 drift test), all green. Revert method unchanged: scratch tree from
`git archive 4a5c447`, each fix surgically reverted, revert verified landed by assertion/compile
before trusting any result.

New/changed tests this cycle — each fails when its fix is reverted:

- `TestProtocolRejectsPathTraversalInTheLock` (modified) — **real now.** Fails on revert with
  "rejected for the wrong reason".
- `TestLoadRefusesToEscapeTheStore` (new) — **real.** Fails on revert: reads "PLANTED".
- `TestPublishRefusesWithoutATerminal` (new) — **real.** Fails on gate removal; also pins the
  honesty of the refusal text.
- `TestProtocolStatusReportsReadErrors` (new) — **real.** Fails on `pinErr`-branch removal.

**No tautologies among the new tests.** The gap this cycle is the inverse of round 2's: the tests
are honest, but five of the seven claimed fixes have NO test that notices their absence — each
revert below left all 18 protocol tests green:

- the `droppedContent` multiplicity/section rewrite (CRITICAL fix; remove the decrement → green);
- `assertNoSymlinkComponents` (MAJOR fix; remove the call → green);
- CRLF-core normalization (remove `rel.Body` normalization → green; the CRLF test only exercises a
  CRLF DECK);
- the synced-stamp skip (remove → green; no test bumps a version);
- the ValidVersion trim-reject (restore trimming → green).

The failed-publish cleanup is at least honestly disclosed as untested in IMPLEMENTATION.md — that
is the right G7b posture, and the other four should have had either tests or the same disclosure.
The pre-existing suite's strengths and the dead O_NOFOLLOW half are as assessed in round 02.

**G7b sweep (any text claiming a guarantee with no end-to-end test):** the changelog's
`DETECTED-UNATTRIBUTED` (no implementation at all); §7's "no agent-accessible code path writes a
release" (false for pty-agents, all three copies); IMPLEMENTATION.md's "semantic erasure cannot hide
behind a coincidental match" (false, demonstrated); render.go:199-200's comment claim (false);
`check`'s preamble (contradicts its entries); refusal-text "hash detection" (overstated); the
cross-build "check" that doesn't exist; "nineteen" tests. The TTY refusal and status read-errors
moved from claimed-untested to claimed-AND-tested this cycle — that part of G7b genuinely improved,
as did the store-component symlink hardening, CRLF handling, and the traversal guard tests.

## What I verified that is GOOD (so the next cycle keeps it)

- All three builds + vet + full suite green on fresh `-count=1`; no runner flake.
- codex-1's round-02 droppedContent fixture now reports correctly; multiplicity counts correctly;
  no false positives found (fully-carried deck → "already matches"; roster rows carried by slot not
  reported; second render idempotent; `check` exit 0).
- All three store-component symlink attacks refused on publish, outside dirs untouched; padded
  version rejected; happy-path pty publish works.
- CRLF core × {LF, CRLF} deck: uniform endings both ways, zero CRCRLF, idempotent, check green.
- Failed publish (`ulimit -f 0`) cleans up; the version label republishes successfully.
- Version bump no longer cries wolf on the stamp.
- §7 "Not yet in force" paragraph byte-identical across deck/embedded/skill-wt; drift test green.
- All four new/changed tests fail on revert; the honesty assertion inside the publish test is a
  genuinely good pattern.
- IMPLEMENTATION.md's disclosure of the untested cleanup, and its honest account of the
  `/var`-symlink scoping mistake, are exactly the candor this idea keeps needing.
