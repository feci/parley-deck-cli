---
idea: meta-protocol-change-designated-implementer
artifact: release-preparation audit
agent: claude-1
role: existing non-implementer independent reviewer (participant; not organizer, not implementer)
date: 2026-09-25
cli-candidate: 01fc49526ca66b771b549a01cfbbf4addeaa886b
skill-candidate: a5664d803f1fb6156ef95ff5921dcf13a3dd2031
reviewed-source: 717f3debe3aabb0f7a702d689e8d10de1e04fa8c
record-commit: ee8849c9ba6f75748fb471fe9c910bf362300e8d
integration-merge: eb5f58646f6ab7418d503c37e2fde7a0cc9f0220 (origin/main c49b464ab6a82b863cf13d5501315eaf42f91917)
staged-core: /Users/tomasfecko/.parley/staging/COOPERATION-2.14.0.md sha256 51476d69ed4291b77db63f6855554cd251461c0b8743a303c2a3b0eacaabf67a
handoff-read: release-preparation-kimi-1.md (committed b775ae7d6f98642029216bc9fb266b902a953eac)
verdict: PASS (release preparation) — 3 non-blocking observations, 0 blocking
---

Independent audit of the organizer-authorized RELEASE PREPARATION for CLI **1.50.0** / skill+core
**2.14.0**. Nothing is published, merged, tagged, pushed, installed or configured by this audit;
no source, metadata, peer artifact or signature was edited; the attended core publish command was
read and never invoked. This audit is preparation-only: **publication remains the organizer's
later step.**

## Candidate immutability

Pinned at audit start `2026-09-25T08:10:29Z` and re-pinned at end `T08:20:55Z` — **unchanged**:
CLI `01fc4952…886b`, skill `a5664d80…2031`. During the audit the branch HEAD advanced
`01fc495 → b775ae7` (Kimi's handoff). I verified the candidate is still an ancestor and that
**no source or metadata moved**: `git diff --name-only 01fc495 b775ae7 -- internal/ cmd/ VERSION
CHANGELOG.md go.mod go.sum` is **empty**; the whole delta is the single file
`release-preparation-kimi-1.md`. No mismatch to raise.

Handoff availability: absent at my first read, present at `10:16` and read in full before this
verdict. Kimi's process exit is the organizer's gate, not mine — I claim nothing about it.

## What I verified myself (PRIMARY), against what was cited

**Reviewed behaviour survives integration untouched.** `internal/app/driver_impl.go`,
`internal/app/driver_designation_test.go`, `internal/protocol/implementer.go`,
`internal/protocol/implementer_test.go` are **byte-unchanged** between reviewed `717f3de` and the
candidate. Integration `717f3de→eb5f586` touched exactly three files in `internal/`:
`app/version.go` and the two inherited Linux `procctl` repairs — released-main content, not this
idea's.

**No unrelated change.** `git diff --stat c49b464 01fc495` over source+metadata = **14 files**:
this idea's 11 designation files plus `VERSION`, `internal/app/version.go`, `CHANGELOG.md`. The
candidate metadata commit `eb5f586→01fc495` is those three metadata files only — **no code change
beyond the version const**, as claimed.

**Version metadata.** `VERSION` 1.49.1→1.50.0; const 1.49.1→1.50.0;
`TestVersionFileMatchesBinaryVersion` re-run by me at the candidate in a fresh `git archive`
checkout → `ok 0.699s`. CHANGELOG 1.50.0 — 2026-09-25 entry describes the mechanism accurately
and states "Ships UNSET".

**Six CLI artifacts** in `release-delivery/2026-09-25-designated-implementer/cli/assets/`. I
recomputed all six SHA256 — **all six match `sha256.json` and `build-provenance.json`**; exactly
six assets on disk, no extras:

| asset | sha256 |
|---|---|
| parley-v1.50.0-darwin-arm64 | `14f4c2c6…6118` |
| parley-v1.50.0-darwin-x64 | `820590b2…f24f` |
| parley-v1.50.0-linux-arm64 | `47b4c918…f7e4` |
| parley-v1.50.0-linux-x64 | `d9aed7e4…8746` |
| parley-v1.50.0-windows-arm64.exe | `1cec3b37…3bac` |
| parley-v1.50.0-windows-x64.exe | `2514cb3d…cffe` |

`go version -m` on **all six** (the handoff spot-checked two): every one carries
`vcs.revision=01fc49526ca66b771b549a01cfbbf4addeaa886b`, **`vcs.modified=false`**,
`-trimpath=true`, and the correct `GOOS`/`GOARCH` pair. Host executable check: the darwin-arm64
asset run by me prints **`parley 1.50.0`** (installed host CLI is still 1.49.1 — nothing was
installed). `build.py` asserts a clean porcelain before and after the six builds.

**Staged core — the load-bearing check, re-derived independently, not read from the provenance
note.** Base published TEMPLATE `~/.parley/protocol/core/2.13.0/COOPERATION.md`
sha256 `fc907e59…62c9f`, 109,772 B. Staged `COOPERATION-2.14.0.md` sha256 `51476d69…f67a`,
115,010 B / 1,393 lines / mode `0444`.
- `diff base staged` = **exactly 7 change blocks / 24 changed lines** at `58c58`, `292a293,298`,
  `402c408`, `426c432`, `436c442,444`, `830c838`, `874a883,888`.
- **Content equality, not just shape:** I extracted the added/removed line sets from the staged
  delta and from this idea's own reviewed protocol delta (`git diff c49b464 01fc495 --
  internal/protocol/defaults/COOPERATION.md`) and compared them — **+19 / −5 on both sides, with
  identical content hashes** (`add c5528d13…`, `del 8136c081…`). The staged core is the published
  2.13.0 template plus exactly this idea's reviewed hunks and nothing else.
- **Template zones retained, no project header leaked:** staged keeps `<workspace-name>`,
  `<transport-choice>`, `<YYYY-MM-DD>`; it carries none of the live deck header
  (`parley-deck`, `github-pr`, `Protocol synced:`, roster). `diff staged embedded@01fc495` is
  **exactly `5,6c5,6`** — the template-zone lines alone.
- **AC-1, all copies, computed by me:** deck@01fc495 `tail -n +160`, embedded@01fc495
  `tail -n +153`, skill@a5664d8 `tail -n +153`, staged core `tail -n +153` **all hash
  `3621b7a637c5183bcf9120ff73e231e2a4503859fdec6cd6b84bf9713beef30e`.**
- **Attended publish command checked, never invoked.** `parley protocol publish --version 2.14.0
  --from /Users/tomasfecko/.parley/staging/COOPERATION-2.14.0.md` is well-formed for
  `internal/app/protocol.go`'s gate, which refuses without a controlling terminal
  (`hasTTYSupported` + `platformHasTTY`). I allocated no TTY. The gate's own comment honestly
  records that it does not stop a pty-allocating agent (DF-1, unshipped) — unchanged, not this
  release's issue.

**Skill candidate and exact tarball.** `parley-deck-skill-2.14.0.tgz`, 389,330 B,
sha256 `681b9076…135b`, sha512 `897b1d35…0be7` — **both recomputed by me and matching the
handoff**. I extracted it and compared **every one of its 210 files** against
`git show a5664d8:<path>`: **210/210 byte-identical, 0 files absent from the candidate,
0 content mismatches.** `package.json` 2.14.0; `compatibility.json` skillVersion 2.14.0 /
recommendedCli `>=1.50.0`; exactly **one** `COOPERATION.md` copy, whose protocol tail hashes to
the AC-1 value above; the packaged reference documents the default as "deliberately shipped
COMMENTED OUT" and §10 as "(ships unset)". Skill trees `a5664d8` and the superseded `4835a9f`
are both `2a7fd7d8…4cb8` — the amendment was identity-only, as stated.

**Skill gates — run by me, not accepted from a log.** At my first read the delivery `logs/`
directory held **no** `npm test`, `manifest --check` or `npm pack --dry-run` evidence, so rather
than waive them I executed them at the exact candidate in an isolated `/tmp` `git archive` tree:
- `npm test` → **exit 0: node 399/399 pass, 0 fail; python 3.14 54 tests OK across 7 files;
  manifest `--check` ok for all six add-ons** (no stale manifest — the precise failure mode that
  bit the 2.13.0 release).
- `npm pack --dry-run` → exit 0, `parley-deck-skill-2.14.0.tgz`, **entryCount 210**, and its file
  list is **identical** to the delivered tarball's; the `prepack` manifest gate ran and passed.
- Independently reproduced the handoff's disclosed first-attempt failure: without
  `node_modules`, `test/design-addons.test.js` fails with `Cannot find module 'commonmark'` — a
  declared devDependency, i.e. **environment, not product**. After a local `npm ci` confined to
  `/tmp` (no global, no channel, no skill install), green. My diagnosis and Kimi's agree
  independently.
- Portable build: `dist/parley-deck-skill-v2.14.0-macos-arm64` hashes
  `3af49e8c…55a7`, **matching** `logs/portable-current-sha256.txt`. Install dry-run log ends
  `exit 0` across 12 targets; doctor JSON reports `ok: true`, 12 targets, none unhealthy. Both
  are Kimi's runs, cited as such — I did not re-execute an installer even in dry-run.

**Hosted CI, at the exact candidate SHAs.**
- CLI run **36110596716**, head `01fc4952…886b` (matches the candidate): ubuntu **success**,
  macOS **success** (it was `in_progress` for most of this audit and I did not treat it as a pass
  until it completed), windows **failure**.
- Skill run **36110644777**, head `a5664d80…2031` (matches the candidate): both jobs
  **success**.

**Windows failure — inspected, not assumed inherited.** I pulled both job logs and diffed the
failure sets. Candidate `01fc495`: 98 failing top-level tests across 14 packages plus the compile
error `internal\evidence\tree_report_test.go:97:20: undefined: syscall.Mkfifo`. Predecessor
`c49b464` (run 36062251339): 100 failing tests, the **identical 14-package set**, the **identical
compile error**. `comm` of the two sorted sets: **zero tests fail at the candidate that did not
already fail on main**; two flaked green. The two version-named failures
(`TestVersionAllJSONIncludesSkillStatus`, `TestVersionAllUsesDirFlagForProjectStatus`) are
**present at the predecessor too** and are caused by `exec: "parley-deck-skill": executable file
not found in %PATH%`, not by the version bump — indeed that test's own payload prints
`parley 1.50.0`, confirming the new const. **Conclusion: the Windows red is entirely inherited
and this delta contributes nothing to it.**

## Findings

**Three observations, none blocking. No defect found in any candidate, artifact, hash, staged
core or gate.**

**OBS-1 — `cli-release-notes.md` cites only the predecessor CI run; candidate evidence now
exists and is stronger.** The notes state the Windows leg "of the integrated predecessor state
was failing at release time (predecessor run 36062251339)" and that "Linux and macOS CI legs of
the integrated state passed (same run)". Since then the candidate's own run **36110596716 at
`01fc495`** has completed: **Linux and macOS green at the exact commit being released**, Windows
red. The notes are not false — they are stale and understate the evidence. Fix, notes-only:
cite run 36110596716 for all three legs. Organizer/implementer action before publication.

**OBS-2 — the skill gate results are narrated in the handoff but not persisted as logs.** The
delivery `logs/` directory holds the Go suite log, doctor JSON, install dry-run log and portable
hash, but no `npm test` / `manifest --check` / `npm pack --dry-run` output, unlike the 2.13.0
delivery which kept `npm-test-2.12.1.log`. The gates themselves are satisfied — I re-ran them
green above — so this is durability of channel evidence, not a gate failure. Suggest persisting
the three outputs before the final completion file.

**OBS-3 — the Windows label is correct, and its basis is stronger than the notes say.** On
Windows the `internal/app` test binary **panics** at `app_test.go:155` and aborts, so the tests
after that point in the package — including this idea's entire designation suite in
`driver_designation_test.go` — **do not execute at all**. Windows therefore yields *no* evidence
about this feature, rather than "some tests are red". That makes the
experimental/unvalidated labelling and the held CLI WinGet substantively right, and it is worth
saying plainly in the notes so a later reader does not over-read a Windows asset. Inherited
condition; `windows-portability` owns it (DF-4).

## Evidence scoping — preserved exactly, nothing re-labelled

- **Independent full suite: `d238238`** (mine, non-implementer) — not re-claimed at any later
  commit.
- **Targeted behaviour evidence: `1bad263`** (both reviewers' focused runs, M1=9 / M2=5, probe
  sets).
- **Comment-only delta: `717f3de`**; record `ee8849c`; goal-check live re-execution `804522c`.
- **New at `01fc495`:** Kimi's broad `go test ./... -count=1` on the merged tree — every package
  ok, `SUITE_EXIT=0` (`logs/go-merged-tree-01fc495.log`). This is the **producer's** evidence,
  cited as such; I did **not** duplicate it, and the integration of main's Linux repairs with
  this delta is a genuine new concern that justified it. Its named limitation stands: it ran on
  darwin/arm64, so the Linux `procctl` paths are covered by hosted Linux CI (green at the
  candidate), not by that local run. `gofmt` leftovers are the five known pre-existing files.
- **Mine at `01fc495` / `a5664d8`:** the six-asset hash + buildinfo verification, the host
  executable version, `TestVersionFileMatchesBinaryVersion`, the full staged-core derivation, the
  210-file tarball comparison, `npm test` / `npm pack --dry-run` / `manifest --check`, and the
  two-run CI failure-set diff.

## Reviewed-behaviour posture, unchanged

Source, default and fail-closed behaviour remain exactly as reviewed and signed: the four
behaviour files are byte-identical to `717f3de`; `internal/config/runtime.go:654` still emits
`# default_implementer = "agent-id"` **commented out** with the "SHIPPED UNSET ON PURPOSE"
rationale; no active `default_implementer` exists in any shipped TOML or template at the
candidate. **The product ships UNSET.** The owner's later machine default
(`default_implementer = "codex-1"` in `~/.parley/agents.toml`) is post-release owner
configuration only and has not happened. The accepted residual NIT at `IMPLEMENTATION.md:325-326`
is byte-untouched, as the closing consensus recorded.

## Verdict

**PASS — release preparation.** The committed candidates, the six CLI artifacts, the exact npm
tarball, the staged core and every skill gate I could execute are correct, internally consistent
and faithful to the reviewed candidates. The three observations above are notes/evidence items,
not defects in anything shipping. No gate was waived: the skill gates that were missing from the
delivery I executed myself rather than accept on narration.

**This verdict covers preparation only.** It authorizes nothing: publication, merge, tagging,
npm, Homebrew, WinGet, skill install and the attended core publish all remain the
organizer's/owner's later steps.

## Remaining channel audit obligations (not performed, not waived)

1. **Post-publication participant verification of every actual channel** — GitHub releases and
   asset digests, npm `dist.integrity` after publish, both Homebrew formulae, skill WinGet
   (`Feci.ParleyDeckSkill`, one application per PR from the isolated `feci-skill-2.14.0` branch;
   **CLI WinGet stays HELD**), and the runtime SKILL.md hash audit after `install --target all`.
2. **WinGet hashes must come from the final uploaded release assets** after the portable workflow
   completes — never from the local `dist/` build.
3. **The CLI release-create command is not documented in-repo** (inventory gap; the handoff marks
   its version derived and asks the organizer to confirm). Worth closing in the runbook.
4. **Attended core publish** verification after the owner runs it, then the owner-only machine
   default, verified independently afterwards.
5. **Final completion file** must state deferred/owner-only items honestly — Windows
   (DF-4/`windows-portability`), DF-1…DF-3, and OBS-1…OBS-3 above — before releasing
   `windows-portability`'s sequencing hold.
