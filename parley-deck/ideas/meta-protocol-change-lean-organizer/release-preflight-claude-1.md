---
author: claude-1
idea: meta-protocol-change-lean-organizer
date: 2026-09-24
phase: release preflight (post-close, independent second-model review)
verdict: PASS
scope: post-close release metadata deltas and release readiness only — NOT a new idea review round
---

# Release preflight — CLI 1.49.0 / skill 2.13.0 (claude-1, independent)

**Verdict: PASS.** No blockers. Five non-blocking observations below.

This report satisfies `parley-deck-skill/RELEASING.md:25` ("Ask a second model to review the
final diff before publishing"). I reviewed only the post-close release deltas; the feature
verdicts, the final suites and the LE-7 goal PASS already sit with the reviewers and the
commissioned checker, and I did not reopen them.

## Reviewed hashes

| Repo | Range reviewed | Head | Accepted source it must not drift from |
|---|---|---|---|
| CLI (`parley-deck-cli`) | `4df0855..06e563e` | `06e563e8b1fe8094132149b83e14da7be9aae51e` | `c3baf0977b82c2e9fc386b174f45de9e5a1dc1d9` |
| Skill (`parley-deck-skill`) | `b06a65a..8161e5e` | `8161e5e9d3ecdece25b7ae225e2304d3159be79b` | `b06a65adaa081ebc063046f51fcff8c00cdb08d0` |

Verified in my own detached-HEAD worktrees (`lean-organizer-review-claude-1`,
`lean-organizer-review-claude-1-skill`), both clean at those heads.

## 1. No feature drift — PASS

`git diff --name-status c3baf09..06e563e` outside `parley-deck/**`: exactly `M CHANGELOG.md`,
`M VERSION`, `M internal/app/version.go`. Nothing else.

`git diff --name-status b06a65a..8161e5e`: exactly `M CHANGELOG.md`, `M package.json`,
`M package-lock.json`, `M skills/parley-deck/parley-addon.json`,
`M skills/parley-deck/references/compatibility.json`.

All eight paths are `M` — **no additions, no deletions in either repo**. No `lib/`, `bin/`,
`internal/` (beyond the version constant), `cmd/` or SKILL/reference payload changed. Both
deltas are metadata-only as claimed.

## 2. Version metadata correctness — PASS

- `VERSION` = `1.49.0`; `internal/app/version.go` `const version = "1.49.0"`; built binary prints
  `parley 1.49.0`; `go test ./internal/app -run Version -count=1` → 10/10 PASS including
  `TestVersionFileMatchesBinaryVersion` (the guard that pins the file against the constant).
- `go build ./...` exit 0, `go vet ./...` exit 0 at `06e563e`.
- Skill: `package.json` + `package-lock.json` (both the root `version` and `packages.""` entry)
  = `2.13.0`; `compatibility.json` `skillVersion` = `2.13.0`, `recommendedCli` = `>=1.49.0`;
  `node bin/parley-deck-skill.js --version` → `2.13.0`. The `recommendedCli` bump is justified —
  the core now documents `parley wait`, which ships in 1.49.0.

## 3. Manifest integrity — PASS

- Recomputed independently: `shasum -a 256 skills/parley-deck/references/compatibility.json` =
  `19ff1df406de35678dbb76a8a93c43995fc3674211843e24616928ed6d245550`, byte-identical to the
  manifest entry. The only per-file hash that changed is the one whose payload changed.
- `npm run manifest:check` exit 0; `parley-deck` aggregate
  `d8931d0081c4d563cc400436ece9f6bb94e5ba7bfe0a78f6f67e0e02c876193d` matches the committed
  `aggregate`; the other five add-ons unchanged and ok.
- `npm test` at my tree: **399 pass / 0 fail**, python 3.14 54 tests OK across 7 files, manifest
  gate ok, **exit 0**. This independently reproduces the recorded prep result; the manifest was
  regenerated, not hand-edited.

## 4. Package whitelist + release notes — PASS

**Whitelist.** `files` in `package.json` is unchanged by the delta. `npm pack --dry-run` →
`parley-deck-skill-2.13.0.tgz`, **210 files**, 387.4 kB packed / 1.3 MB unpacked, prepack gate
green. Packed top level is exactly the whitelist: `skills/` (201), `bin/`, `lib/` (2),
`README.md`, `LICENSE`, `NOTICE.md`, `gemini-extension.json`, `plugin.json`, `package.json`.
No `.github/`, `packaging/`, `scripts/`, `test/`, `dist/`, `node_modules/`, deck artifacts or
dotfiles leak in. The Python tools and JS tests that do ship live **inside** `skills/parley-bidding`
and `skills/parley-design-check` — intended add-on payload (`RELEASING.md:15-16`) and integrity-gated.
Since the delta adds and deletes nothing, the packed file set is structurally identical to 2.12.1.

**CLI release notes.** Every 1.49.0 claim checked against the source, not the commit message:
`parley wait` exists and its help states the exact claimed exit map (0 boundary / 1 usage-IO /
3 timeout with partial digest / 4 invalid artifact); `parley protocol packet --audience
participant|facilitator` exists with the attestation fields; `parley usage ingest --agent …
--source codex-rollout|claude-jsonl` exists; the two "Fixed" items are real
(`internal/app/preflight.go:395` "cannot read workspace status" fail-closed with its test at
`facilitator_test.go:146-152`; content-keyed next action with
`TestPhaseDigestNextActionFreshCheckoutEqualMtimeAwaitsReview` and
`TestPhaseDigestNextActionFixUpPublishedAwaitsReview` in `wait_test.go`).

**Limitations section is honest and accurate.** It claims **no** measured token reduction and
names the ledger as the future instrument; it records live attribution as ambiguous; the
"Windows exercised by the CI leg" claim is true — `.github/workflows/tests.yml` matrix is
`[ubuntu-latest, windows-latest, macos-latest]`; the 45-minute suite timeout is stated rather
than hidden.

**Skill release notes.** Core `skills/parley-deck/SKILL.md` = **17,802 B** — "17.8 KB, under the
20 KB cap" is accurate (decimal KB) and the four named reference files are present. The `--json`
stream contract is documented in the core at `SKILL.md:131-143` exactly as the entry describes
(stdout envelope on 0/3/4, terminal status on stderr, exit 1 no envelope). The permissive §9.0
audience/brief sentence is present in the portable `references/COOPERATION.md` snapshot (line 880).

## 5. RELEASING.md second-model checklist

| Item | Result |
|---|---|
| cross-platform path handling | No code changed by the delta; full suite (incl. the macOS-specific firmlink/`chflags` cases) green at my tree |
| safe overwrite / update / uninstall | Exercised in an isolated temp dest (`--target generic --dest /tmp/…`, never a global install): fresh install → `doctor` `ok:true`, `status: valid`, `version 2.13.0`, `missing: 0`; `uninstall --dry-run` names the right single path. Temp artifacts removed. |
| Codex / Claude / Antigravity / Hermes / legacy Gemini targets | `install --target all --dry-run` resolves each correctly: `~/.codex/skills/`, `~/.claude/skills/`, `~/.gemini/config/plugins/` (agy), `~/.gemini/extensions/` (legacy gemini), `~/.hermes/skills/` — 12 targets total |
| npm package file whitelist | See §4 — PASS |
| portable binary build and asset upload | `.github/workflows/release-portable.yml` derives ref and asset names from the release tag, so 2.13.0 flows automatically once `v2.13.0` exists. I did not run the pkg Windows build (recorded green in prep). |
| Homebrew formula correctness | Formulas live only in the separate tap (`feci/homebrew-parley`), correctly absent from this repo. Local tap checkout is **in sync** with `origin/main` at `f20b78f` (1.48.0 / 2.12.1); both formulas still need `url` + `sha256` bumped to `v1.49.0` / `v2.13.0` at release time. |
| README install clarity | README carries **no** hardcoded version strings; the documented entry point is the version-agnostic `npx -y parley-deck-skill@latest install`. No drift possible. |

## 6. Release readiness (verified myself, read-only)

- CLI `origin/main` = `b4831d6` **is** an ancestor of `06e563e` → fast-forward possible.
- Skill `origin/main` = `d1e57d5` **is** an ancestor of `8161e5e` → fast-forward possible.
- Neither `v1.49.0` nor `v2.13.0` exists yet; `v1.48.0` / `v2.12.1` are the current tags.
- Representative cross-builds from `06e563e` succeed: `windows/amd64` (10,789,376 B) and
  `linux/arm64` (9,568,416 B) with the plan's exact `-trimpath -ldflags "-s -w"` flags.
- Runbook (`release-plan-zcode-1.md`) corrections are sound: roles are right (zcode-1 metadata
  prep only; organizer performs release operations; owner ALONE attends the core publish), the
  "reported complete only when both are done" language is gone, the attended core publish is
  stated as a **pending owner-only parallel action and explicitly NOT a channel gate**, and the
  winget channel is honest — opening the two PRs is the completable action, external catalog
  merge is not claimed as a completion gate and open PRs stay reported as open. I found **no
  invented core-publish or catalog-merge gate** anywhere in the delta; the closure record
  likewise states "Nothing in this closure performs, implies, or gates on any of them."

## 7. Observations (non-blocking, organizer's call)

1. **`doctor --target all --json` returns `ok:false` on this machine** — every detected target
   reports `status: malformed` with `missing: [references/ARTIFACT_TEMPLATES.md,
   HEADLESS_LAUNCH.md, ROSTER_AND_PROTOCOL.md]`. This is **stale local state, not a package
   defect**: the machine still has 2.12.1 installed, which predates those files. A clean 2.13.0
   install reports `ok:true` / `valid` / `missing: 0`. Two consequences worth knowing: this
   `RELEASING.md` preflight step will not read green until 2.13.0 is installed, and existing
   2.12.1 users who run `doctor` from a 2.13.0 `npx` will see "malformed" (not "outdated") until
   they re-install. The wording originates in the already-accepted source, not in this delta.
2. **Runbook §0 tap row is stale in the safe direction**: it says the tap is 1 commit behind at
   1.47.0/2.11.0 needing `git pull --ff-only`. The local checkout is already at `f20b78f`
   (1.48.0/2.12.1), so no pull is needed. The plan's own freshness correction (#5) covers the
   pattern; only this row's specifics are out of date.
3. Correction #3 cites 626.9 s for the slowest package while the commit message cites 633.7 s
   for this run — two separate measurements, both far under the 45-minute timeout. No
   substantive contradiction.
4. The CLI repo has a `docs/releases/<version>/` convention from 1.48.0 (organizer assessment
   plus independent reviews). No `docs/releases/1.49.0/` counterpart exists yet. My report is
   filed where dispatch directed instead; whether to also place an organizer assessment there is
   the organizer's decision, not a gate.
5. The CLI repo has **no** release workflow (only `tests.yml`) and no cross-compile script — the
   six binaries are manual builds per the runbook, as at 1.48.0. Verified buildable above.

## 8. Checks deliberately not re-run

The full Go suite (recorded 31/31 ok, exit 0 at this tree) and the LE-7 goal PASS are already on
record; I ran targeted version/build/vet checks instead of repeating a 45-minute run. I did not
run the pkg Windows portable build.

## Boundaries observed

No source or peer-artifact edits; no changes to the closed IMPLEMENTATION.md or the frozen FINAL;
no tag, push, publish, merge, release or global install; all installer testing confined to
temporary directories, since removed. This file is the only artifact I wrote.
