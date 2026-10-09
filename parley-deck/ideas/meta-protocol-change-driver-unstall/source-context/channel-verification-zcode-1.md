---
agent: zcode-1
idea: meta-protocol-change-driver-unstall
date: 2026-10-09
role: independent post-publication channel verifier (fresh process, separate from all reviews, signoffs and the goal check)
product-merge: bed5ffd1049dc0f207eb39beaa3dd3ba879181c3
skill-merge: 054954458dbfe2d764bf6234a6d6bbeee12e2494
reviewed-product-commit: 134ac40cd178ebe0318a838d9357c4e6f561f925
track: deliberation
transport: github-pr
---

# Channel verification — zcode-1, fresh process, 2026-10-09 (attempt 02, ~15:35Z start)

## Role, scope and method

I am zcode-1, launched in a NEW process for the mandatory AC9 post-publication channel
verification. This file is my only canonical write. I edited no FINAL/IMPLEMENTATION
bytes, no prior artifact, signature, production source, installed skill, credential or
release; I invoked no Claude/Kimi, pruned/declared no worktree, allocated no terminal,
applied no real legacy declaration and published no npm/core. Read-only `gh`/`curl`/git
APIs only; no browser. Temporary downloads live under
`R/channel-probes/attempt02/` only. This is channel verification, not a new source
review; I re-ran no saturated source suite. Usage/token telemetry was unavailable to me;
I invent no figures.

Retry context honored: attempt 01 (`13bedb27-…`) timed out at 1800 s with **no canonical
artifact or verdict**; I treat its result as nothing and did not wait, sleep-poll or
rerun anything on its behalf. I completed within this attempt's 2400 s ceiling with no
long sleeps. My independent assessment of the attempt-01 cleanup observation is in
**Channel 8 / process note** below.

Protocol attestation: launch supplied the full live protocol verbatim (`context_mode:
full`, `source_sha256 = packet_sha256 =
0357d504982f92b713f2f86604129f46e1ef3276664e4da91d39ecf420253f0c`); the shadow packet
audit in the launch header is an unapplied diagnostic (packet_bytes 0). I re-verified
(PRIMARY, `shasum -a 256`) that `parley-deck/COOPERATION.md` still hashes to that exact
value and equals `meta/version.json.protocolSha256` (`protocolRole: source`) — no drift.

Provenance convention (§15): every check below marked PRIMARY is a command or API call I
executed myself in this process, with its result recorded here or byte-comparable to a
retained artifact whose bytes I re-hashed myself. Organizer records were used only as
claims to test against, never as evidence by themselves. Retained evidence I re-hashed
was confirmed byte-identical before use (raw log hashes and archive hashes below).

## Channel 1 — Merges, tags, source byte identity — VERIFIED

- PRIMARY `gh pr view 81 --repo feci/parley-deck-cli`: state MERGED, mergedAt
  2026-10-09T12:26:04Z, mergeCommit `bed5ffd1049dc0f207eb39beaa3dd3ba879181c3`;
  `git cat-file -p bed5ffd` shows **two parents** `a9e383d2…` + `6fc6da25…` (real merge
  commit, not squash). PR11 `--repo feci/parley-deck-skill`: MERGED 12:26:11Z, merge
  commit `054954458dbfe2d764bf6234a6d6bbeee12e2494`, parents `8ce4dec0…` + `e46e5518…`.
- PRIMARY tags: `v1.54.0^{commit}` → `bed5ffd…` and `v2.18.0^{commit}` → `0549544…`
  (resolved in the fetched worktrees); both `origin/main` heads equal the merge commits
  and `git log <merge>..origin/main` is **empty** — nothing was pushed after either
  merge.
- PRIMARY CLI byte identity: `git diff --name-only 134ac40..6fc6da2` lists 22 paths,
  all inside this idea (`IMPLEMENTATION.md`, review/consensus + cycle archive + zcode
  round files, `source-context/*` evidence, one inbox note). The same diff restricted to
  `':(exclude)parley-deck'` is **empty** — product, docs, protocol and all non-deck
  bytes equal the reviewed 134ac40. (PR81 merged head 6fc6da2, not 134ac40 itself; the
  delta is exactly the audit records named above.)
- PRIMARY skill: merge-commit tree `7bb32ae081f2f93eb2d9787a7daa7f6277d08de5` **equals**
  reviewed `e46e551871005ecac4225f7cb55a652754fbbb9b`'s tree — the merge introduced zero
  byte changes beyond the reviewed tree.
- PRIMARY frozen artifacts: `FINAL.md` at bed5ffd hashes `e29953c134dbe0f2…a2f0d` and
  `IMPLEMENTATION.md` hashes `28a6d27213002d91…be26`, both byte-equal to the current
  worktree files; with no commits after the merges and no local modification, there is
  **no closed FINAL/IMPLEMENTATION drift after the source merge**.

## Channel 2 — GitHub releases and hosted assets — VERIFIED

- PRIMARY `gh release view v1.54.0 / v2.18.0`: both exist, `isDraft: false`,
  `targetCommitish` = the merge commits; CLI release carries the six platform binaries
  plus `sha256.json`; skill release carries the npm-format `.tgz`, macos-arm64 portable
  and both Windows portables.
- PRIMARY digest triangle: all 11 live GitHub asset digests equal (a) the organizer's
  `hosted-assets-verification.json` values and (b) my `shasum -a 256` of the local
  `R/hosted-assets/*` copies — e.g. binaries `0496c050…` (darwin-arm64), `b9b581a8…`
  (win-x64), skill tgz `2d1e7901…`; the hosted `sha256.json` asset (hash `c0be61b2…`)
  internally lists exactly the six live binary digests I checked.
- PRIMARY independent re-download (fresh `curl` into `attempt02/`): CLI
  `parley-v1.54.0-darwin-arm64` → `0496c050…`, CLI `windows-x64.exe` → `b9b581a8…`,
  skill `parley-deck-skill-2.18.0.tgz` → `2d1e7901…` — byte-identical to the live
  digests.
- PRIMARY binary execution (darwin-arm64 only — I claim no execution on foreign
  architectures): `./parley-v1.54.0-darwin-arm64 --version` → `parley 1.54.0`;
  `go version -m` → `vcs.revision=bed5ffd1049dc0f207eb39beaa3dd3ba879181c3`,
  `vcs.modified=false` (clean build of the exact merge commit), `-trimpath`, CGO off.
- PRIMARY skill payload: extracted my fresh tgz → `package/package.json` version
  `2.18.0`, 210 files, and all six `skills/*/SKILL.md` hashes equal the
  `skill/payload-verification.json` values (`aad7d7d8…` bidding, `ece72723…` deck,
  `9ccf8e9d…` design, `337aafd6…` design-check, `25b7a6b2…` tracker, `91f94c40…`
  worktrees).
- PRIMARY Windows installer architecture: parsed PE headers of all four local Windows
  assets (byte-identical to live): machine `0x8664` (x64) for both `…windows-x64.exe`,
  `0xAA64` (arm64) for both `…windows-arm64.exe`. Static header check only; no foreign
  execution claimed.

## Channel 3 — Skill portable workflow run 37930155652 — VERIFIED (as cross-build)

- PRIMARY `gh run view 37930155652 --repo feci/parley-deck-skill`: `completed/success`,
  event `release`, headSha `0549544…` (the merge/tag). Jobs: `macos-tests` success
  (12:31:15→12:32:12; npm ci, npm test), `windows` success (12:32:14→12:32:53; npm ci,
  npm test, `npm run build:portable:windows`, Upload release assets) — the `windows` job
  started **after** macos-tests completed (`needs: macos-tests`).
- PRIMARY raw log `R/logs/skill-portable.log`: the `windows` job's Set-up lines record
  `Image: ubuntu-24.04`; the workflow source at the merge commit declares
  `windows: needs: macos-tests / runs-on: ubuntu-latest`. So this job **cross-builds
  Windows installers on Ubuntu after macOS tests — it is NOT a native Windows test
  run**, and I record it exactly that way.

## Channel 4 — Homebrew tap and installation — VERIFIED

- PRIMARY tap: `repos/feci/homebrew-parley/commits/main` HEAD is exactly
  `bd2a3640a9d217d4c3f48d7121c6669690ed61f0`; the compare API diff of that commit
  touches exactly two files — `Formula/parley-deck-cli.rb` and
  `Formula/parley-deck-skill.rb` — each changing **only** the `url` and `sha256` lines
  (v1.53.0→v1.54.0, v2.17.0→v2.18.0). No other file or line changed; no unrelated tap
  trust touched.
- PRIMARY formula hashes: CLI source hash `fd07f50c…` and skill source hash
  `b089855c…` each equal (a) `R/parley-deck-*-source.tar.gz`, (b) the attempt-01
  retained archives, and (c) my **fresh download** of the GitHub v1.54.0 tag archive
  (`fd07f50c…` re-derived just now).
- PRIMARY live brew state: `brew list --versions` shows `parley-deck-cli` 1.54.0 and
  `parley-deck-skill` 2.18.0 installed/linked (installed 2026-10-09), with older kegs
  (1.50–1.53 / 2.14–2.17) **retained** as backups. `/opt/homebrew/bin/parley` →
  Cellar 1.54.0 and `parley --version` → `parley 1.54.0`.
- PRIMARY installed payload, not version strings: the 2.18.0 keg's six
  `libexec/skills/*/SKILL.md` files hash **exactly** to the six release payload hashes
  (shebang-restored bytes match the tarball payload).
- Organizer steps re-read with retained logs: `brew upgrade` / `brew test` /
  `brew style` for both fully-qualified `feci/parley/…` names each exit 0
  (`homebrew-verification.json` + `logs/brew-*.log`). I did not re-run brew test myself
  (no need; installed bytes already independently verified).

## Channel 5 — Runtime skill installations — VERIFIED

- PRIMARY full re-hash: I re-hashed **every one** of the 114 recorded installed
  `SKILL.md` files myself (`sha256`): 0 missing, 0 hash mismatches against the six
  release payload hashes.
- PRIMARY destinations: the 114 files span 19 distinct roots = 15 managed runtimes
  (aionrs, claude, codebuddy, codex, cursor, factory, gemini×2, goose, hermes,
  kimi-code, opencode, qwen, vibe, zcode) **plus exactly the four required generic
  targets** `~/.hermes/profiles/{ldx,librade,testprofile}/skills` and
  `~/.config/opencode/skills` — matching `runtime-verification.json` (15 managed + 4
  generic, 114 files). No manual copying or target edits needed or performed.
- PRIMARY backup: `runtime-skills-before-2.18.0.tar.gz` exists and hashes to the
  recorded `2ea79367…`.
- PRIMARY read-only installer check: `parley-deck-skill --version` → `2.18.0`;
  `parley-deck-skill doctor --target all --scope user --json` → `ok: true`, all
  detected targets report status `valid`.

## Channel 6 — Skill-only WinGet PR — VERIFIED as open/pending, correctly disclosed

- PRIMARY `gh pr view 449541 --repo microsoft/winget-pkgs`: **OPEN** (not draft),
  author feci, title "Update: Feci.ParleyDeckSkill to 2.18.0", created 12:36:10Z,
  exactly **three added files** (`…2.18.0/Feci.ParleyDeckSkill.installer.yaml` 15
  lines, `….locale.en-US.yaml` 27, `….yaml` 6), `mergeable: MERGEABLE`.
- PRIMARY upstream checks at my check time: 01 Pull Request Validation, 02 Manifest
  Validation, 03 URLs, 04 URL Domain, 05 Manifest Policy, 06 Catalog Content,
  07 Installers Scan, license/cla — all **SUCCESS**; **08 Installation Validation
  IN_PROGRESS**, 09 Metadata / 10 Validation Completed QUEUED. Current upstream state is
  *validating, not merged* — recorded honestly; the brief requires the PR to be open and
  the state disclosed, which is satisfied. No merge/availability claim.
- PRIMARY manifest bytes: fetched the raw `installer.yaml` from the PR head — local
  `R/winget/manifests/…` copy is byte-identical (`bd3606fc…`). Its x64 `InstallerSha256
  04F82EB64EB9A78C…` and arm64 `FE2659C4E7D41432…` match the live hosted Windows asset
  digests I verified in Channel 2; URLs point at the v2.18.0 release assets;
  InstallerType portable.
- Local schema pre-validation passed per `winget-validation.json`; **native `winget
  validate`/`install` is unavailable on this macOS host and is disclosed as not
  performed** — no local Windows installation is claimed.
- PRIMARY no CLI WinGet change: `gh search prs --author feci --repo
  microsoft/winget-pkgs` → the only OPEN feci PR is 449541 (skill); the newest CLI
  package update merged is 1.48.0 — the CLI WinGet channel remains held, as required.

## Channel 7 — Owner-only pending actions, D6 state, honest status — VERIFIED

- PRIMARY npm registry (live `https://registry.npmjs.org/parley-deck-skill`):
  dist-tags latest `2.13.0`; version list ends at 2.13.0 — **2.15.0, 2.16.0, 2.17.0 and
  2.18.0 are all absent (unpublished)**. (Observation, outside this idea's recorded
  2.15–2.18 scope: 2.14.0 is also absent from npm; it predates this handoff's brief and
  its staging file exists — flagged for the owner's awareness, no scope change made.)
- PRIMARY installed core: `~/.parley/protocol/core/` contains only `2.10.0` and
  `2.13.0`; staged `COOPERATION-2.18.0.md` hashes to the recorded `9ead4475…`
  (unpublished). `owner-publication-status.json` / `owner-actions-current.md` list
  npm+core commands for 2.15→2.18 in **ascending order** with matching archive/core
  hashes; every listed archive path exists with the stated hash (2.18.0 re-hashed by
  me; earlier ones recorded and spot-consistent with the retained delivery dirs).
- PRIMARY D6 unapplied: the git-common budget area contains only the pre-existing
  `cycles-23451469…` entry; `find … -name "legacy-*"` over the common dir returns
  **zero** entries — no declaration exists. Re-checked again **after** my inspect below
  (still zero).
- PRIMARY May bytes: `parley-deck/runs/20260510T194003Z/events.jsonl` hashes
  `ee4f52b717159e54aa8f144161d14467160cfc6846aed97abbdd6768683c6906`, byte-equal to
  `git show 3ec10ac:…` — the original May history is untouched.
- PRIMARY read-only inspect reproduction (released CLI 1.54.0): `parley budget legacy
  inspect --dir <this worktree> --run parley-deck/runs/20260510T194003Z` returned
  `history: "unknown-history"`, `manifest_sha256 7283348e…`, `preview_sha256
  d7d6d890…`, `decisions: []`, and the **identical 24-root set** recorded in
  `legacy-activation-preview-current.json` / `owner-legacy-action.json`
  (`applied: false`, `matching_visible_roots: 24`). Every field matches the organizer
  record exactly; inspect created no authority. The one-time apply command in
  `owner-legacy-action.json/.md` remains explicitly **pending until the owner stops
  writers and attends it** — no stopped-writer assertion or attendance is fabricated
  anywhere in the records.

## Channel 8 — Main CI, Windows tripwire, cancelled duplicates, process note — VERIFIED

- PRIMARY `gh run view 37929997523 --repo feci/parley-deck-cli`: `completed/failure`
  (push, main, head `bed5ffd…`, the merge). Jobs: ubuntu **success** (113818091066),
  macos **success** (113818091433), windows **failure** (113818091383, ended 13:13:26).
  This is the newest main run; I did not wait for or launch any further CI.
- PRIMARY raw-bytes tripwire re-derivation: `R/logs/windows-main-merged.log` re-hashes
  to the recorded `6bb90162ead4716463215240cc32a456338f33792713734d7d5e268af23d8dbb`
  (baseline log re-hashes to the recorded `da7c0e4f…`). Applying the recorded
  extraction pattern `(?m)^[0-9T:.Z-]+[ \t]+--- FAIL: ([^\s]+)` myself: **355** direct
  distinct failing names — set-equal to `main-windows-comparison.json`'s claimed set;
  baseline (main 128e30b) has **356**; **added_failure_names = ∅** (my independent
  set difference, both directions); one baseline name absent
  (`TestConcurrentReservationsCannotOverspend` — fewer failures, not a regression);
  `panic: test timed out after 45m0s` present; log anchored on bed5ffd.
- Assessment: **no new in-scope Windows regression is demonstrated** by the main-merged
  run. Honest limits retained exactly as recorded: name-set equality proves neither
  cause equality nor per-test coverage (non-verbose summaries; opt-in/conditional skips
  exist); native CLI Windows remains **unresolved/experimental**; broad repair stays
  deferred (FINAL ALT-10); CLI WinGet stays held. Nothing was suppressed to reach this.
- PRIMARY cancelled duplicates: all four runs live-verified `completed/cancelled`
  (37929983549 PR + 37929976410 push at 6fc6da2; 37929305000 PR + 37929299196 push at
  7ff3d8e). I independently confirmed both heads are **audit-only**: `git diff
  134ac40..7ff3d8e -- ':(exclude)parley-deck'` and the 6fc6da2 equivalent are empty.
  These cancellations freed duplicate CI of byte-identical product trees; they are
  recorded as **not passes** (`not_a_pass_claim: true`), and the product runs at
  134ac40 plus the main bed5ffd run were preserved. No source/code/test suppression
  occurred anywhere I checked.
- Process note on attempt-01 cleanup (independent assessment, as instructed): the
  surviving read-only CI-polling shell/sleep after the attempt-01 timeout was
  organizer-channel-check tooling, terminated by the organizer; `channel-retry-policy.json`
  itself labels it manual cleanup, not product behavior. It is **not** evidence about
  `RunParticipantStep` supervision/cleanup (D2), no product claim in FINAL/consensus
  rests on it, and nothing in this verification depends on it. No broader claim is
  warranted from one observation; the deferred-gaps list stands.

## Honest limits of this verification

- Windows: everything above is name-set/timeout/raw-byte analysis of retained logs plus
  live statuses; I did not (and cannot, on this Mac) run native Windows tests or winget,
  and no such claim is made. Cause equality and per-test coverage remain unproven.
- The skill portable `windows` job is an Ubuntu cross-build; Windows portables were
  verified by hash and PE header only, never executed.
- npm/core publication, the one-time D6 apply, WinGet upstream completion and the
  released inbox handoff are correctly **owner/upstream-pending**, not delivered by this
  verification; the handoff did not exist when I started and does not exist at this
  writing (PRIMARY `ls parley-deck/inbox/` — no `…driver-unstall_released.md`).
- Organizer evidence was used only after byte-identity confirmation (raw logs, archives,
  manifests, payloads); organizer JSON summaries were treated as claims, and every
  load-bearing one was re-derived or re-hashed by me above.

## Verdict

Every authorized delivery requirement of this release was verified satisfied with
current PRIMARY evidence: two-parent merges and tags at the reviewed bytes; live
releases with digest-exact assets and clean `vcs.modified=false` 1.54.0 binary; the
2.18.0 skill payload identical everywhere it is installed (114/114 files); successful
portable workflow correctly described as an Ubuntu cross-build; the Homebrew tap bump
limited to exactly the two formula URL/hash pairs with installed payload bytes verified;
the skill-only WinGet PR open with exact three-file diff, correct installer hashes and
upstream state disclosed as validating-not-merged; npm/core 2.15–2.18 and the D6
activation correctly recorded as owner-pending with the inspect preview reproduced
byte-for-byte and the May history untouched; main CI honestly terminal-failed with a
clean re-derived added-name tripwire and audit-only cancellations correctly framed. The
owner-only pending actions, the WinGet state and the Windows limits are correctly stated
in the retained records. Nothing I checked broke a delivery requirement, and no
demonstrated new in-scope regression exists to assess.

CHANNEL-VERIFICATION: PASS
