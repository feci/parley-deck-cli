# Consult result — zcode-1 advisory release-readiness inventory (meta-protocol-change-designated-implementer)

Advisory only; no verification verdict, no authorization, no mutation performed. Sources: owner runbook `release-orchestration-codex-1.md`, Kimi's advisory plan `implementation-release-plan-kimi-1.md`, predecessor `source-context/release-1.49.1-done.md`, both worktrees, and predecessor delivery evidence. Versions are tentative **CLI 1.50.0 / skill & core 2.14.0**, recheck against GitHub/npm/core-store latest at staging (R55).

## Version-bearing files
- CLI (repo `feci/parley-deck-cli`): `VERSION` (currently `1.49.0` in this worktree — lags published v1.49.1, so integrate origin/main first), `internal/app/version.go:3` `const version = "1.49.0"` (baked into binaries), `CHANGELOG.md` heading.
- Skill (`feci/parley-deck-skill`): `package.json` `"version": "2.13.0"`, both `"version"` entries in `package-lock.json`, `CHANGELOG.md` heading. npm name `parley-deck-skill`.
- Core: version chosen at publish; store `~/.parley/protocol/core/` currently holds 2.10.0 and 2.13.0.

## CLI build, assets, workflows
- Predecessor pattern (evidence `release-delivery/2026-09-24-lean-organizer/release-1.49.1/build.py`): clean `git clone --shared --no-checkout` + `checkout --detach <commit>`, then per platform `CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go build -trimpath -ldflags "-s -w" -o parley-v1.50.0-<os>-<label>[.exe] ./cmd/parley` — six targets: darwin-arm64, darwin-x64, linux-arm64, linux-x64, windows-arm64.exe, windows-x64.exe; SHA256 json + provenance recorded; binaries embed `vcs.revision`, `vcs.modified=false`.
- Local validation: `go build ./... && go vet ./... && go test ./...`; `gofmt -l internal/ cmd/` empty; plan's spot-checks (default_implementer commented out in `internal/config/runtime.go`, no `impl-claim`, tail-hash equality of the two COOPERATION.md copies).
- Workflows: CLI has only `.github/workflows/tests.yml` (push any branch/PR/dispatch; matrix ubuntu/windows/macos; `go test ./... -count=1 -timeout 45m`). No release-automation workflow and **no RELEASING.md in the CLI repo** — the exact `gh release create` invocation is not captured anywhere I could find (evidence `github-publish.txt` holds only the resulting release URL). Flag: release must re-derive the release-creation command; do not invent one here.
- Windows leg was failing at predecessor main (run 36062251339 @ c49b464): a direct main merge will likely show a red Windows CI leg; both Windows assets stay individually labelled **experimental/unvalidated** in notes and asset labels (owner decision, unchanged).

## Skill build, assets, workflows (`RELEASING.md`)
- Preflight: `npm test` (Node + python3≥3.10 leg — fails without python3 — + add-on manifest check), `npm run manifest:addons` only if an add-on payload changed (this idea changes none), `npm pack --dry-run`, `npm run build:portable:current`, `bin/parley-deck-skill.js install --target all --dry-run`, `… doctor --target all --json`.
- `.github/workflows/release-portable.yml`: triggers `release: published` **and** `workflow_dispatch` with `tag` input; macOS tests → ubuntu build leg → `gh release upload "$TAG" dist/* --clobber`. **Clobber risk**: assets are overwritten on re-run/re-publish, so winget hashes must come from final release assets (`gh release view "v${VERSION}" --json assets --jq '.assets[] | [.name, .digest] | @tsv'`), never a local build. Expected assets: `dist/parley-deck-skill-v2.14.0-windows-x64.exe` and `-windows-arm64.exe`. `test.yml` fires on any push/PR — direct main merge triggers it.
- This idea's skill-tree deltas are only `skills/parley-deck/references/COOPERATION.md` (identical hunks) and one Phase-5 line in `references/ROSTER_AND_PROTOCOL.md`; six skills ship (parley-bidding, parley-deck, parley-design, parley-design-check, parley-tracker, parley-worktrees).

## npm exact-tarball publication
Owner-authorized form: `npm pack`, then `npm publish --access public parley-deck-skill-2.14.0.tgz` (exact tarball — overrides RELEASING.md's plain `npm publish --access public`). Verify independently with `npm view parley-deck-skill version dist.integrity dist.tarball --json` reporting 2.14.0.

## Homebrew — BOTH formulae, one tap repo
- `feci/homebrew-parley/Formula/parley-deck-cli.rb` — `brew (install|upgrade) feci/parley/parley-deck-cli`; url `https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.50.0.tar.gz` + its SHA256.
- `feci/homebrew-parley/Formula/parley-deck-skill.rb` — url/sha256 from the `v2.14.0` tag tarball (`curl -fsSL … | shasum -a 256`); minimal install/symlink block per RELEASING.md.
- The skill repo must keep **no** formula copy (drift rule). Checks: `brew style`, `brew audit --strict --online feci/parley/<formula>`, `brew upgrade`, `brew test`, `parley --version` / `parley-deck-skill --version`.

## WinGet
- CLI (`Feci.ParleyDeckCli`): **held** — no PR until `windows-portability` ships (R52; predecessor confirmed none open).
- Skill (`Feci.ParleyDeckSkill`): normal path — draft manifest at `packaging/winget/manifests/f/Feci/ParleyDeckSkill/2.14.0/` (in-repo drafts stop at 1.3.1; 2.13.0 went straight to fork PR — expect to author a fresh draft), hashes from final GitHub release assets only, copy into a fork of `microsoft/winget-pkgs`, validate on Windows (`winget validate`, `winget install --manifest`), **one application per PR**. Precedent: PR440360 (2.13.0) merged 2026-09-24.

## Runtime skill install targets + hash inventory
`install --target all` covers the `TARGETS` registry in `lib/installer.js:18`: codex (`.codex/skills`), claude (`.claude/skills`), agy (`.gemini/config/plugins`), gemini (`.gemini/extensions`), hermes (`.hermes/skills`), qwen, codebuddy, goose, kimi (`.kimi-code/skills`), droid (`.factory/skills`), vibe, cursor, opencode, aionrs, zcode (`.zcode/skills`) — 15 named targets plus `auto`/`all`/`generic`. Owner sequence: install to all sessions, then **participant hash audit of every runtime SKILL.md**; `doctor --target all --json` and the schema-2 add-on manifest marker (`npm run manifest:check`) supply the hash surface.

## Core staging + attended publish
- Base (R53): published `~/.parley/protocol/core/2.13.0/COOPERATION.md`, SHA-256 `fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f`, 109,772 bytes, mode 0444 — verified present in the store.
- Stage = that base plus **exactly** the reviewed hunks (§0 defaults sentence, §4.0 template field, §4 Phase-4 two edits, §4 Phase-5 mechanism, §9.0 readiness sub-bullet, §10 TL;DR item 6), extracted at integration via `git diff <base> -- internal/protocol/defaults/COOPERATION.md`; hunks are byte-identical across `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`.
- Publish shape: `parley protocol publish --version 2.14.0 --from FILE` — attended, requires a controlling TTY (`internal/app/protocol.go:362-375`, gate G2); never a pty workaround; owner-only, separate from agent-controlled channels.

## Post-release machine default
Product ships **UNSET** (`default_implementer` only a commented line in the generated central template). Owner-only, after shipment and independent verification: uncomment `default_implementer = "codex-1"` in `~/.parley/agents.toml`. Layers low→high: `~/.parley/agents.toml` → `parley-deck/agents.toml` → `agents.local.toml` → `$PARLEY_HEADLESS_AGENT_CONFIG` (non-empty wins; empty never clears); `default_implementer = "none"` suppresses deck-wide; per-idea `implementer:` in `00-prompt.md` outranks all.

## Gaps / flags
1. CLI release-creation command (`gh release create` with notes/labels) is undocumented in-repo — predecessor evidence holds only the URL; derive and record it at staging.
2. Worktree VERSION/version.go/CHANGELOG lag v1.49.1; integration of origin/main precedes any bump.
3. Skill in-repo winget drafts end at 1.3.1 — no committed 2.13.0 draft to copy from.
4. Red Windows CI leg expected on merge until `windows-portability`; labeling policy unchanged.

Reminder: implementation review, zero-fix signatures, fresh goal check, and integration gates all precede any release preparation; nothing above starts that work.
