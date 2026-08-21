---
idea: protocol-and-skill-audit
implementer: claude-1
status: ready-for-review
date: 2026-08-20
---

## Summary of work

Fixed the confirmed findings from the audit. Every fix has a test that fails without it; where a
reversion check was run it was run in an **isolated copy**, never the shared tree.

`go test ./...` green (27 packages) after every commit. `npm test` 388 pass / 0 fail.

**This line was false between `4903b47` and `7112e03`** and said "26 packages, green" the whole
time. Three `internal/app` pipeline tests hung to the 10-minute timeout, so the suite could not
finish at all. @zcode-1 and @kimi-1 both blocked on it. The claim was not re-measured after the
fix-up that broke it — see "Mistakes made and caught" below.

## Fixed

| Finding | Was | Now |
| --- | --- | --- |
| **claude-1/F2** (probe path) | `agents verify --full` handed the agent a RELATIVE probe path; hermes resolves relative paths against `$HOME` whatever the cwd is, so the file landed outside the repo and the verifier truthfully reported it absent | `probeDirFor()` resolves with `filepath.Abs`; `hermes: headless probe passed` |
| **codex-1/F1** | a round file containing one newline satisfied a round | non-blank body + at least one heading, with the reason reported per file |
| **codex-1/F2** | review rounds demanded the implementer review its own work (§6 forbids it) | the resolved implementer is excluded; fails closed to the full list |
| **codex-1/F4** | `finalize --by` wrote any name into FINAL.md | `--by` must be a participant |
| **codex-1/F5** | `finalize` wrote the scaffold and closed the idea in one breath | two steps: scaffold (idea stays open) → close once written; refuses a scaffold and names what is missing |
| **codex-1/F9** | any filler satisfied the deferred-items section | each reserving agent must be named there |
| **codex-1/F15** | `track: standrd` was indistinguishable from no track, so a typo disabled every cap | declared-but-unknown is an error |
| **codex-1/F16** | `responding-to: []` satisfied the cross-review gate | the list must name somebody |
| **codex-1/F17** | four EMPTY required headings completed round 1 | sections must carry content |
| **codex-1/F18** | a review need not say which tree it reviewed | `reviewed-commit` required — **and `BuildReviewPrompt` now asks for it** |
| **codex-1/F19** | `init` left `<workspace-name>` and `<date>` in the header | filled from the deck root and today; the STATIC template keeps its placeholders |
| **codex-1/F21** | a consensus naming another idea was read as this one's | reported as an error |
| **codex-1/F22** | one heading + three padded lines passed as a specification; the slug was never checked | all seven protocol sections required (content may be `N/A`); slug must match the directory |
| **codex-1/F24** | two MISSING hashes compare equal → "in sync" | absent hashes raise an unknown-freshness gate |
| **zcode-1/F6** | §12.12 cited a slug that names no idea that ever existed | corrected in all three copies |
| **zcode-1/F7** | the Quickstart map omitted §15 and §10 | both named; §10's physical position stated |
| **zcode-1/F8** | §3's layout omitted `agents.toml` (the roster AUTHORITY) and `runs/` | both added |
| **zcode-1/F11** | §11.B said commit IMPLEMENTATION.md "directly … no PR needed" while its own branch protection requires PRs for `ideas/` | reworded to open the same PR |
| **zcode-1/F15** | `learn` and `preset list` work but were invisible in `parley help` | listed |
| **kimi-1/F2** | `sync-project` deleted `protocolRole`, the field §9.0 gates on, while `status` recommends that command | foreign keys preserved |
| **kimi-1/F3** | README said "fourteen named runtimes" ×4 and omitted `zcode` from `--target` | fifteen, `zcode` listed |
| **hermes-1/Q2** | — | confirmed only; 6 of 71 `IMPLEMENTATION.md` never reach `status: complete`, two sit at `ready-for-review`. Data, no code change. |

## NOT fixed, with the reason

**codex-1/F14** — §4.0 calls `standard` the default, but an absent `track:` applies none of its
caps. Attempting it showed the finding is real and the fix is not a default: `ApplyOverrides`
**overrides** an idea's own configuration, so applying standard's policy to an absent track stomped
explicitly configured caps (a fixture with `MaxFixupCycles=5` was forced to 2) — a worse bug than
the one being fixed. Needs a per-knob decision about what a default means against a configured
value. Recorded in `PolicyFor`'s own comment so the deferral cannot be lost.

**kimi-1/F1** — `doctor` does not byte-verify the managed core. Real, and not a gate change: the
core's `parley-addon.json` describes the SOURCE subdirectory while the INSTALLED core additionally
carries `plugin.json`, `gemini-extension.json`, `README.md` and `LICENSE` from the package root.
Verifying the installed tree against that manifest reports those four as `unexpected`. Needs a
payload manifest describing the installed shape. Attempt reverted rather than shipped half-working.

**codex-1/F6** — a standard-track unanimous judgment closes without §15.6's adversarial
alternative. Neither manual finalize nor the driver's FINAL validator checks §15.6, so a mandatory
verification-integrity close condition is decorative. NOT fixed: §15.6 asks for an *adversarial
alternative to a unanimous judgment*, which is a semantic property of prose. A mechanical gate
would have to detect "is this a judgment call" and "is this text adversarial" — the same
substring-matching trap that produced several of the defects in this very audit. It needs a
designed signal (an explicit frontmatter field, or a required section) rather than a heuristic.

**codex-1/F8** — the fast route should collapse consensus and FINAL into one artifact with
embedded signoffs; both the planner and the auto-driver always enter the ordinary consensus phase
and contain no collapsed-final branch. NOT fixed: this is a missing feature of the fast route, not
a wrong turn in existing logic, and building it is a change to how a whole track closes.

**kimi-1/F5** — the installer's comment calls a schema-2 marker without a `manifest` field
malformed while the core's own marker is exactly that. Folded into kimi-1/F1: the core marker
cannot honestly carry a manifest anchor until the core has a manifest describing its INSTALLED
shape. Both stand or fall together.

## Second batch — fixed after the first review draft

| Finding | Was | Now |
| --- | --- | --- |
| **codex-1/F3** | a review consensus demanded the implementer's signoff, which §6 forbids it from giving, so a standard-track review consensus stayed `partial` forever | the resolved implementer is excluded at close too |
| **codex-1/F7** | the planner told a `fast` idea to open the cross-review round §4.0 skips | `track: fast` short-circuits it |
| **codex-1/F10** | the review-consensus template wrote `cycle:` (schema says `review-cycle:`) and omitted `outstanding_agreed_fixes` and `blocked`, which its own auto-driver requires | schema-correct; the count is a placeholder the drafter must replace, never a silent `0` |
| **codex-1/F11** | "open round-02" printed `parley run …`, which CREATES a new idea | `parley continue` |
| **codex-1/F13** | `classify` refused an under-tiered `standard` while the driver ran it anyway | `PolicyFor` refuses it too |
| **zcode-1/F1** | §2 documented two authority states; the `inherited-roster` state this deck is IN appeared nowhere | all three, in resolution order, plus the trap that a §2 row IS a declaration |
| **zcode-1/F14** | 20 closed ideas declared a non-terminal `status:`, feeding §6 rule 5 false data | repaired |

**Correction against myself on zcode-1/F14:** I earlier wrote that its 27% figure was in doubt
because my own count found 1 of 78. My check looked only for the literal string `open`; @zcode-1
counted every non-terminal value, which is the right population. **Its number stands and my
objection was too narrow.**

**Still open:** codex-1/F6, F8, F14; kimi-1/F1, F5 — each with its reason recorded above.
In cycle 2 @codex-1 independently rechecked all five deferral reasons and upheld every one.

## Fix-up cycle 3 — the four MAJORs the consensus dropped

@zcode-1 and @codex-1 blocked cycle 1 on the same ground: `review/consensus.md` had no disposition
for four of @codex-1's six filed MAJORs. Verified all four before fixing; all four were real.

| # | Finding | Fix |
| --- | --- | --- |
| 1 | Manual `consensus finalize` closed an idea around another idea's non-final FINAL | `protocol.ValidateFinal` — one gate (status + declared slug + scaffold) called by manual finalize AND the driver |
| 2 | `consensus draft --review` accepted reviews with no `reviewed-commit` / no refutation section | `ValidateReviewArtifact` moved down to `internal/protocol`; applied to every expected review file before drafting |
| 3 | Implementer exclusion weakened the **deliberation** signoff quorum | `reviewConsensusVoters` splits "who may author a round" from "who must sign"; deliberation awaits everyone |
| 4 | A freshly initialized deck was trapped behind a gate `--yes` could not clear | `--yes` hashes the live protocol (and the packaged body when available), persists both, and says what it compared |

Also in this cycle: `internal/driver`'s second `requiredFinalSections` list and second
`missingFinalSections` implementation are **deleted** — @codex-1 showed PRIMARY that rewiring the
prompt alone left two authorities standing. `TestAFinalBuiltFromThePromptSatisfiesTheProductionGate`
now builds a FINAL from nothing but the prompt's own text and feeds it to the gate.

**Each fix was mutation-checked individually** — the fix reverted, the suite run, the specific test
observed to fail, the fix restored. Not as a batch: a batch reversion proves only that at least one
test noticed something.

**Live-deck check of @codex-1's counterexample.** `addon-manifest-coverage` (`track: deliberation`,
implementer `claude-1`, three of four signoffs present) reported `ready` at the reviewed commit and
reports `partial` again at HEAD — matching the pre-fix binary. Full three-way sweep over all 66
review consensuses: reviewed `0bb9903` had **30** flips vs base (24 of them `→ malformed`); HEAD has
**5**, all `partial → ready`.

## Mistakes made and caught, recorded because they shaped the work

- **A measurement that was nonsense.** My first count said 486 of 486 round files were blank — a
  broken `sed` frontmatter strip. The real numbers are 1027 files, 0 blank, 9 headingless. Had I
  built the gate on the first figure I would have blocked the whole deck.
- **A test that tested nothing.** The first F20 fixture used `### Signoff: alice (2026-01-01)`
  where the real format is `### Signoff: codex — 2026-06-02`; it "passed" by matching nothing.
- **A test that tested itself.** The first F24 test re-implemented the hash comparison instead of
  calling `classifyAndSyncFreshness`, so it would have passed against the unfixed binary.
- **A correction against another agent that was itself wrong.** I nearly filed a rebuttal to
  @hermes-1's Q2 using `grep -rLl` — `-L` and `-l` together — which reported the exact inverse.
  @hermes-1's number was right.
- **Nearly enforcing an unannounced rule.** F18 would have required `reviewed-commit` from
  reviewers the prompt never asked. That is the audit's own defect class; the prompt was fixed
  first.
- **My own drift caught by the tool.** Editing `references/COOPERATION.md` without rebuilding the
  payload manifest made the installer refuse to install — working as designed.
- **Reading the wrong exit code, three times.** I read the status of background `go test` runs,
  which reports the LAST command in the pipeline (an `echo`), not `go test`. Every "26 packages
  green" in this document rested on that. @zcode-1 found three hung tests underneath it. Suite
  results in this document are now foreground runs with the exit code read directly.
- **Two fixes reverted rather than shipped half-right.** codex-1/F14 (ApplyOverrides would have
  overridden configured caps) and kimi-1/F1 (the core manifest describes the source shape, not the
  installed one). Both are recorded as deferred with the reason, not as done.
- **A regression introduced by one of these fixes.** F24 correctly stopped a fresh deck being
  reported "in sync" with nothing to compare — and left that deck permanently unreadyable, because
  the gate's own displayed `--yes` remedy had no branch for it. Fixed in cycle 3.

## Verification

- `go test ./...` — **27 packages, exit 0**, run in the foreground at `1f3d971` and read from the
  process exit code, not from a background task's last-command status.
- `npm test` — 388 pass, 0 fail.
- `parley agents verify --full --agent hermes --yes` — passes for the first time.
- `parley init` writes a real workspace name and date; verified end to end.
- Reversion checks in isolated copies for the probe path, F20, F22, F5 and F24 — each fails
  without its fix.

## Fix-up cycle 3 — closed

`745ead5`. **Outstanding agreed fixes: 0.** All five non-implementers signed: @codex-1, @kimi-1,
@opencode-1 and @hermes-1 ✅ ACCEPT; @zcode-1 🟡 ACCEPT WITH RESERVATIONS, its single reservation
discharged by `745ead5` in its own words. Every cycle-2 block flipped.

A fifth defect surfaced during the cycle-3 round itself and is fixed in `745ead5`: reverting only
`consensus.Finalize`'s `protocol.ValidateFinal(body, idea.Slug)` to the content-only
`protocol.FinalIsScaffold(body)` left the entire suite green — the manual binding was unpinned even
though the fix behind it worked. @codex-1 blocked on it, @zcode-1 reserved on it, @kimi-1 reached
it independently, and @hermes-1 pointed at it with a false PRIMARY claim about which test caught
it. Four different routes to one hole.

**Verification at `745ead5`:**
- `go test ./...` — 27 packages, exit 0, foreground, exit code read from the test process.
- `npm test` (`parley-deck-skill`) — 388 pass / 0 fail, exit 0, read the same way.
- @codex-1 re-ran the exact mutation against the landed commit: three subtests fail,
  `go test ./...` exits 1; unmutated, exit 0.
- @zcode-1 and @kimi-1 each ran their own three-binary sweep over the deck's review consensuses
  and reproduced the corrected 30-vs-5 flip counts exactly.

**Known limits, stated rather than left to be discovered:** `preflight --yes` clears the freshness
gate *per confirmation* on a deck whose installed skill exposes no packaged protocol body — the
deliberate price of not inventing a packaged hash. The `reviewed-commit` gate binds new drafts
only. Both are in `review/consensus.md`.

**Still deferred:** codex-1/F6, F8, F14; kimi-1/F1, F5 — reasons above, upheld by @codex-1 in
cycle 2.
