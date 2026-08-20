---
idea: protocol-and-skill-audit
implementer: claude-1
status: ready-for-review
date: 2026-08-20
---

## Summary of work

Fixed the confirmed findings from the audit. Every fix has a test that fails without it; where a
reversion check was run it was run in an **isolated copy**, never the shared tree.

`go test ./...` green (26 packages) after every commit. `npm test` 388 pass / 0 fail.

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

**Still open:** codex-1/F3, F6, F7, F8, F10, F11, F13; zcode-1/F1, F14.

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

## Verification

- `go test ./...` — 26 packages, green.
- `npm test` — 388 pass, 0 fail.
- `parley agents verify --full --agent hermes --yes` — passes for the first time.
- `parley init` writes a real workspace name and date; verified end to end.
- Reversion checks in isolated copies for the probe path, F20, F22, F5 and F24 — each fails
  without its fix.
