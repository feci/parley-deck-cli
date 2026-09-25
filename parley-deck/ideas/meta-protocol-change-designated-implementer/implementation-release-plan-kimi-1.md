---
title: Advisory release plan — meta-protocol-change-designated-implementer
author: kimi-1 (implementer)
date: 2026-09-25
status: advisory — no build, package, channel, publication, or configuration action performed
---

# Implementation release plan (advisory)

Scope and authority: this plan is **advisory only**. Per FINAL R50–R56, every release, merge,
channel, publication, and global-configuration step below is the **organizer's/owner's attended
act**. Nothing here was executed by the implementer beyond local, reversible verification
(`go build` / `go vet` / `go test` / `gofmt`, three-copy hash checks, greps).

## 1. Preconditions (all owner/organizer-gated)

- **Release order (R50):** `release-1.49.1` → **this idea** → `windows-portability`. The
  predecessor gate is satisfied: `v1.49.1` is tagged, `origin/main` at `c49b464`
  ("complete independently audited release 1.49.1"). Do not let `windows-portability` jump ahead.
- **Reviewed implementation:** Phases 6–8 to zero agreed fixes (reviewers: claude-1, zcode-1),
  then a **fresh non-implementer goal-done check** (COOPERATION.md §4 Phase 8 / LE-7).
- **Independent channel verification by a participant** of every channel actually used (R51).
- **Integration (R56):** integrate latest `origin/main` **only after** the reviewed
  implementation; local canonical files + direct main integration (per-run override; the global
  transport header is not changed).

## 2. Build / verify (CLI worktree)

```bash
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer"
go build ./... && go vet ./... && go test ./...
gofmt -l internal/ cmd/        # must print nothing
scripts/install-local.sh       # optional local install: builds ./cmd/parley into ~/.parley-deck/bin
```

Sanity spot-checks of the shipped feature (unset product must be byte-identical behaviour):

```bash
grep -n 'default_implementer' internal/config/runtime.go   # template emits it COMMENTED OUT
grep -rn "impl-claim" --include="*.go" .                   # must print nothing (R14)
shasum -a 256 parley-deck/COOPERATION.md internal/protocol/defaults/COOPERATION.md
tail -n +160 parley-deck/COOPERATION.md | shasum -a 256    # must equal the next line's hash
tail -n +153 internal/protocol/defaults/COOPERATION.md | shasum -a 256
```

## 3. Skill package preparation (skill worktree, `parley-deck-skill`)

```bash
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer-skill"
npm test                     # node tests + python3 leg + add-on manifest check (refuses stale manifest)
npm run manifest:addons      # only if an add-on payload changed (this idea changes none)
npm pack --dry-run           # prepack re-verifies every shipped add-on payload
```

The only skill-tree deltas from this idea: `skills/parley-deck/references/COOPERATION.md`
(identical hunks) and `skills/parley-deck/references/ROSTER_AND_PROTOCOL.md` (one Phase-5 line).
No installer, manifest, plugin.json, gemini-extension.json, or channel surface changed (R57).

## 4. Version selection (R55 — organizer act at staging time)

Next **minor above actual releases**, re-verified then. Observed at FINAL time: CLI `v1.49.1`,
skill `2.13.0`, core `2.13.0` — so the CLI line would be `1.50.0`, skill/core `2.14.0`,
**re-checked against reality at staging**. The implementer bumped no VERSION file and cut no tag.

## 5. Core staging and publication (owner-attended only)

- **Base (R53, verified):** published live core `~/.parley/protocol/core/2.13.0/COOPERATION.md`,
  SHA-256 `fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f`, 109,772 bytes,
  mode `0444`. Stage the next core as **that published base plus exactly this idea's reviewed
  hunks — nothing else**. The hunks are the ones now byte-identical across
  `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md` (§0 defaults
  sentence, §4.0 template field, §4 Phase 4 two edits, §4 Phase 5 mechanism, §9.0 readiness
  sub-bullet, §10 TL;DR item 6) — extract with
  `git diff <base> -- internal/protocol/defaults/COOPERATION.md` at integration time.
- **Publication (R54):** `parley protocol publish --version V --from FILE` requires a controlling
  terminal and is **escalated to the owner**. Never allocate a pty to work around the gate; no
  participant invokes core publication.

## 6. Channel preparation (organizer/owner attended; nothing done here)

- **npm** `parley-deck-skill`: `npm publish --access public` from the skill repo after its
  preflight passes; tag `v<version>` and push per `RELEASING.md`. Independent channel
  verification (R51) before this fires.
- **Homebrew:** formulas live in `feci/homebrew-parley` — version bump there is a release step.
- **WinGet:** IDs `Feci.ParleyDeckCli` (CLI) and the skill's own. **R52 owner decision —
  preserve exactly:** Windows stays labelled **experimental/unvalidated**, labelled Windows assets
  are **retained**, and **no CLI winget PR is opened** — the owner holds it until
  `windows-portability` ships. Skill channels are unaffected, so the skill's winget path is normal.
- No GitHub release, no tag push, no npm publish, no Homebrew/WinGet PR was performed.

## 7. Machine default — set / inspect (post-release, owner-only)

The product ships **UNSET** (R12): `default_implementer` exists only as a commented-out line in
the generated central template; no deck generator writes it. Setting it is the owner's chosen
post-release configuration — **never a participant act**.

Inspect the live layered value (four layers, low→high: central `~/.parley/agents.toml` →
`parley-deck/agents.toml` → `parley-deck/agents.local.toml` → `$PARLEY_HEADLESS_AGENT_CONFIG`;
later **non-empty** values win; an empty value does not clear a lower layer):

```bash
grep -n 'default_implementer' ~/.parley/agents.toml \
  parley-deck/agents.toml parley-deck/agents.local.toml "$PARLEY_HEADLESS_AGENT_CONFIG" 2>/dev/null
```

Set the machine default (owner's standing choice, e.g. `codex-1`) — edit the central file and
uncomment the shipped line:

```bash
${EDITOR:-vi} ~/.parley/agents.toml
# change:   # default_implementer = "agent-id"   ...
# to:       default_implementer = "codex-1"
```

Deck-wide suppression of a set machine default, and per-idea override:

```toml
# parley-deck/agents.toml  — suppress for this deck:
[defaults]
default_implementer = "none"
```

```yaml
# parley-deck/ideas/<slug>/00-prompt.md — per-idea designation (outranks every config layer):
implementer: <agent-id>     # or `implementer: none` to decline for this idea
```

Behaviour with nothing set anywhere is byte-identical to before this change: no event, no extra
stdout line, FINAL-drafter fallback chain unchanged.
