# Release orchestration record

Status: preparation only. No release/version/channel mutation performed for this idea. Product implementation and verification remain participant-owned; this record is operations, not a code verdict.

## Preconditions

- Predecessor completion record: `source-context/release-1.49.1-done.md`. Its ordering hold is satisfied.
- Current observed GitHub releases (2026-09-25): CLI v1.49.1; skill v2.13.0. npm latest independently queried by organizer for operational state: 2.13.0, integrity `sha512-EHqVIXPveetgOD/xfZZ2iy+iE0ksyXcWXOXCMLZJOUZWHCYetul17mJvTwTFMjM6Iubm9HkNl/y+voVC16zYGg==`. Core published 2.13.0 per participant evidence. Expected next minors 1.50.0 / 2.14.0; recheck at staging.
- CLI origin/main CI run 36062251339 at c49b464: Linux and macOS jobs success, Windows test job failure. This is predecessor state, not evidence for this release; owner holds CLI winget and requires Windows assets labelled experimental/unvalidated.
- Phase-8 fix-up and all closure evidence are still owed. No publication starts before participant reviews, signed zero-fix consensus, independent full suite and fresh non-implementer goal check.

## Authorized delivery sequence

Integrate latest origin/main, participant-prepared metadata/builds and independent validation; direct main merge with no development PR; immutable next-minor tags and GitHub releases; Windows assets retained with CLI experimental/unvalidated labels; BOTH Homebrew formulae; skill winget one application per PR, CLI winget held; npm exact packed tarball with `npm publish --access public <tarball>`; install skill to all sessions and participant hash audit of every runtime SKILL.md; set owner-selected machine default codex-1 only after shipment and independently verify.

Stage core from published 2.13.0 TEMPLATE plus exactly reviewed hunks, independently verify, then provide exact attended `parley protocol publish --version V --from FILE` command. No TTY bypass. Owner-only core publication stays separate from completed agent-controlled channels.

Every actual delivery channel needs independent participant verification before final completion file. That file releases windows-portability's sequencing hold and must state deferred/owner-only items honestly.

## Current closing state (2026-09-25)

Closing consensus cycle 4 is signed ACCEPT by all three, runner exited 0; canonical evidence preserved at 804522c. Zcode is performing a separate fresh LE-7 goal invocation. No release preparation launched pending verdict. Read-only staging-version recheck still reports GitHub CLI 1.49.1 and skill 2.13.0, npm 2.13.0; selected next minors remain 1.50.0 / 2.14.0. Homebrew checkout is clean at 52b09f3. Winget checkout is clean on historical feci-skill-2.11.0; use an isolated updated branch/worktree for this skill PR, never reuse/mix its old branch.

## Candidate CI before publication

Integrated CLI origin/main c49b464 automatically at eb5f586; skill main was already ancestor. Kimi prepared CLI candidate 01fc49526ca66b771b549a01cfbbf4addeaa886b and skill candidate a5664d803f1fb6156ef95ff5921dcf13a3dd2031. Organizer pushed only the candidate branches (no development PR): CLI CI run 36110596716, skill CI run 36110644777. Skill CI reports success; CLI jobs pending at this observation. Branch pushes are prepublication validation; neither immutable tag nor release exists yet. Skill winget isolated sparse clone created in delivery/winget on feci-skill-2.14.0 from fresh upstream; no manifest mutation or PR yet.

## Independent candidate audit dispatched

Resumed exact Claude session 52dccda8-ad11-45b2-805b-7c9ee727cf89 with frozen Opus 5[1m]/max for read-only candidate audit while Kimi finished its suite/handoff. Source candidates are already committed; audit writes only its own canonical report. Publication waits for both process exits and a complete PASS, never merely file appearance. Kimi's merged-tree full suite at 01fc495 now reports exit 0; this is new candidate-specific producer evidence, not a relabelling of prior d238238 independent evidence. Claude assesses it and hosted CI independently. Linux candidate CI succeeded; Windows failed; macOS still running at this observation. Claude must classify the actual Windows failure, not infer inheritance from the platform name.

Publication will upload the six CLI binaries plus the sha256.json manifest named by the prepared release notes; Windows binaries keep per-asset experimental/unvalidated labels. No artifact/version publication has occurred yet.

## Published source and first channels

Both processes exited 0. Claude preparation audit PASS (3 nonblocking observations, no defect); Kimi evidence-retention follow-up exited 0 and committed 4ecff5b, preserving full skill gate outputs at frozen candidate without changing artifacts. Organizer updated ONLY external CLI release notes per OBS1/3: candidate CI36110596716 Linux/macOS pass and Windows app abort means designation tests did not execute; independent failure-set comparison found no new Windows failures.

Direct atomic pushes: CLI main c49b464 -> b775ae7 plus immutable v1.50.0 -> 01fc495; skill main8161e5e -> a5664d8 plus immutable v2.14.0 -> a5664d8. No development PR. GitHub releases published at their standard tag URLs. CLI command: gh release create v1.50.0 --repo feci/parley-deck-cli --verify-tag --title 'Parley Deck CLI 1.50.0' --notes-file DELIVERY/cli-release-notes.md SIX_BINARY_PATHS DELIVERY/cli/sha256.json; Windows path arguments carried #Windows ARM64 — experimental/unvalidated and #Windows x64 — experimental/unvalidated labels. Skill gh release create v2.14.0 --repo feci/parley-deck-skill --verify-tag --title 'Parley Deck Skill 2.14.0' --notes-file DELIVERY/skill-release-notes.md. Exact argv generation and command outcome are retained in tool history; CLI output in delivery/github-cli-publication.log.

Skill portable workflow36113011791 succeeded, final Windows pair uploaded. npm publish --access public EXACT_AUDITED_TARBALL exited1 with E404 access/not-found; owner login requested in the required inbox and asynchronous question, not retried without restored auth. Core staged verified command also escalated to owner with controlling-terminal reason; no TTY invocation by organizer.

Installer operated directly from extraction of the exact audited tarball (npm registry pending): node DELIVERY/skill/install-source/package/bin/parley-deck-skill.js install --target all --json. Exit0, ok=true; 12 detected targets / 72 skill actions. Full JSON and stderr retained in delivery. Independent runtime hashes still owed. No force/undetected override and no unrelated skill mutation. Kimi preparing both Homebrew formulae and skill-only winget metadata next; organizer retains push/PR/install/config actions.

## Remaining agent-controlled channel operations

Kimi channel-prep process exited0; canonical handoff committed f512f66. Both prepared commits have the correct actual [codex-1] prefix (the handoff's [codex] quote is a transcription slip, not the real commit message). Homebrew9855366 published to tap main; installed tap cleanly fast-forwarded, then HOMEBREW_NO_AUTO_UPDATE=1 brew upgrade feci/parley/parley-deck-cli feci/parley/parley-deck-skill exited0. Operational installed commands report parley1.50.0 and skill2.14.0. Independent brew audit/test and remote/hash checks still owed.

Winget26720f091 pushed as feci:feci-skill-2.14.0. PR https://github.com/microsoft/winget-pkgs/pull/441054 opened with repository template, exactly one application/version, explicit local Windows-validation limitation. Duplicate search returned none. Repository AGENTS.md requires local winget validation WHEN PRACTICAL and a truthful summary when unavailable; contrary to the participant handoff's stronger before-PR wording, no Windows-host prerequisite is imposed by that rule. No gate bypass: Windows validation is the PR pipeline; no local winget success claimed. No CLI winget files or PR. Current committed skill workflow/RELEASING/README require exactly the2 Windows assets; predecessor3 extra portable uploads have no documented consumer, so no additional build or clobber.

After GitHub/Homebrew publication and installed CLI1.50.0, organizer applied the owner's explicit subsequent instruction: ~/.parley/agents.toml [defaults].default_implementer from absent to codex-1 at2026-09-25T08:43:51Z. Only that key changed; no roster/model/effort setting changed. Minimal nonsecret operation record in delivery/owner-default-operation.json. Product template remains UNSET per participant audit; independent actual-config verification now owed. npm/core questions remain pending; no elapsed time treated as owner action.

Final participant audit follows. Records-only preservation commits may skip redundant CI; exact product candidates already have independently assessed hosted CI, and main's initial publication run remains observable. No source or tag change is made by those records commits. Shared-memory discovery retry reached the configured gateway but returned Session terminated; no recall/write success claimed.
