---
author: codex-1
idea: meta-protocol-change-lean-organizer
phase: release-only resume
status: complete
---

# Release 1.49.1 orchestration

Owner authorization: inbox/user-to-codex-1_meta-protocol-change-lean-organizer_windows-scope.md. Scope: release CLI 1.49.1 through GitHub and Homebrew; retain Windows assets labelled experimental/unvalidated; hold CLI winget. Separate windows-portability owns the future repair. Skill 2.13.0 unchanged, npm publication remains with the owner-facing Claude session. No A–D deliberation or product code change is authorized here.

Organizer is codex-1; participant authorship and independent verification remain with the existing roster. zcode-1 edits release wording, kimi-1 reviews and audits. claude-1 remains a participant; no additional deliberation is commissioned.

Manual launch fallback: installed `parley consult --timeout 25m zcode-1` rejected the roster ID as not installed/configured, despite active roster and installed adapter. Direct one-shot configured CLI launch is used for the auxiliary release task; no closed workflow is restarted. Transport is local artifact files, with direct main integration and no development PR per owner.

Effective zcode launch: /opt/homebrew/bin/zcode --mode yolo --cwd <lean-organizer workspace> --prompt=<task>; model.main observed zai/glm-5.3 in local config, effort from config, 25-minute timeout. Prompt scopes edits to CHANGELOG 1.49.1 and its own candidate/report; no product changes or publication.

Effective Kimi launch planned from discovered runtime: /Users/tomasfecko/.kimi-code/bin/kimi --output-format stream-json -m kimi-code/k3 -p <task>; print mode provides autonomous workspace writes, effort from config, 25-minute timeout, same workspace. Only its own review/audit reports may be written.

Raw task logs are outside the repository at /tmp/parley-release-1491. Build and channel evidence will be preserved under release-delivery/2026-09-24-lean-organizer/release-1.49.1/.

## Release operations completed

- zcode-1 wording commit: `54e07981740d495f08cd90842168d9228c58c90f`, exactly its declared three metadata/report paths. kimi-1 independently PASSed that exact commit and proposed release notes in `release-wording-1.49.1-review-kimi-1.md`.
- Atomic direct push fast-forwarded remote main from `9134c7a` to `54e0798` and created `v1.49.1` there. No development PR and no published tag move.
- Six fresh binaries built in an isolated clean checkout of `54e0798`, CGO disabled, `go build -trimpath -ldflags "-s -w"`, darwin/linux/windows × arm64/amd64. Build source stayed clean before and after; source, script, provenance and SHA256 map retained outside the repository in the evidence directory. Prior candidate binaries were not reused.
- GitHub release https://github.com/feci/parley-deck-cli/releases/tag/v1.49.1 published with exactly six assets; Windows ARM64/x64 assets each carry the explicit experimental/unvalidated label. Notes are the participant-reviewed CHANGELOG section plus the companion installer link.
- Homebrew CLI-only formula commit `52b09f39c8f4a98612dbf052b1b417bd2f043eb2` published; real tag archive SHA256 `6c0760611a4c4242418015c70e7de1508eaacbaf763a80c27a5a549a1b7a31f4`. Skill formula untouched at2.13.0. Local tap fast-forwarded, then `HOMEBREW_NO_AUTO_UPDATE=1 brew upgrade feci/parley/parley-deck-cli` completed successfully.
- No CLI winget PR prepared/opened. Skill winget PR440360 observed MERGED. No npm login/publication attempted; npm latest was2.12.1 at initial read. Final live registry status is recorded in the completion handoff.
- kimi-1 commissioned for a separate independent channel audit: tag/review alignment, all six real downloaded hashes and Go VCS metadata, notes/asset labels, real Homebrew tag archive hashes, Cellar resolution/version, CLI winget absence, current skill winget/npm/skill GitHub state. The organizer performs release operations and records participant verdicts, not product-code verification.

## Independent closure

kimi-1 returned **PASS**, no findings, in `release-1.49.1-audit-kimi-1.md`. All six downloaded asset SHA256 values match staged values and GitHub digests, all binaries embed the exact reviewed commit with `vcs.modified=false`, tag and source archives agree with the released metadata, Cellar CLI reports1.49.1, Windows labels and winget hold are confirmed. Final organizer command `parley --version` printed `parley 1.49.1`. Final live npm read still reports2.12.1 as latest; owner-facing Claude owns the remaining skill publication. No login attempted. Own session accounting appended to organizer-usage.md. Records-only closure is published after the immutable release commit; the completion file is materialized last so the next release cannot start before the record push finishes.
