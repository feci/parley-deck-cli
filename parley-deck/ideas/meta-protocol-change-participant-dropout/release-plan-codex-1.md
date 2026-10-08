# Release preparation — participant dropout

Prepared by codex-1, 2026-10-09 Europe/Berlin. This is a preparation record, not a release or acceptance claim. Independent review and current-tree validation are still in progress. FINAL and the controlling brief authorize delivery only after both final review-consensus ACCEPT blocks, zero open CRITICAL/MAJOR, and independent criterion evidence.

## Candidate and channels

- CLI 1.52.0 candidate: e4681cf5b9144ed786b269d8ab14b8d65b23d581; implementation PR https://github.com/feci/parley-deck-cli/pull/76.
- Skill 2.16.0 candidate: efe296c7acf13a147ab820ce6cbf8e6705b68691; PR https://github.com/feci/parley-deck-skill/pull/9.
- Latest releases observed 2026-10-08 22:40Z: CLI 1.51.0 and skill 2.15.0. Recheck before tagging. npm remains 2.13.0; 2.15.0 is unpublished.
- Homebrew canonical checkout is `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley`, clean at 320a1d70ee59b9eb47a1c37d2f38360eba313809. Update both formulae after obtaining final source archive hashes.
- Old shared WinGet checkout is clean but on an unrelated 2.11.0 branch. Use a fresh sparse release clone, not that branch. One skill package per PR; no CLI WinGet while Windows is experimental. Final Windows skill asset hashes must come from the successful hosted portable-release workflow.
- Managed installer is 2.15.0. Install the exact final npm payload to all managed targets using `--target all --include-undetected --force`; then use the four explicit `install --target generic --dest ... --force` commands from FINAL. Hash every installed core and companion SKILL.md and run doctor. Do not change roster/models.

## Final delivery sequence

1. After signed close, merge both implementation PRs using merge commits. Preserve all canonical review files and fix consensus histories. Do not squash or prune the owner worktrees.
2. Build CLI six-platform binaries in a clean packaging clone at the final release commit; require VCS commit identity and vcs.modified=false. Publish the GitHub release and hashes. Mark Windows experimental in notes.
3. Pack exact skill npm tarball, build/smoke-test macOS portable, publish GitHub release, wait for hosted Windows portable assets and verify live hashes. npm remains an owner command.
4. Update both Homebrew formulae; style, strict online audit, real upgrade and tests. Keep source archives and hashes.
5. Submit skill-only WinGet 2.16.0 PR with final release-asset hashes and schema validation. Native Windows validation/installation is reported honestly if unavailable on this Mac.
6. Back up existing managed and four named generic installations privately, install exact release payload by the required CLI routes, run doctor and compare every SKILL.md hash.
7. Reproduce staged core 2.16.0 from `/Users/tomasfecko/.parley/staging/COOPERATION-2.15.0.md` (SHA256 0d81fd807114e4ff67f5fa98bba73622b58c86d66a2bf6b57d1def68cd099c09) plus only final reviewed P1–P6. Preserve generic zones. Stage but never publish through an invented terminal.
8. Launch zcode-1 separately with short channel verification brief, exact final refs and all channel evidence. Zcode authors its own report. Resolve material delivery issues before finalizing.
9. Save verified reusable OpenViking outcome and check persistence. Write `inbox/codex-1-to-user_meta-protocol-change-participant-dropout_released.md` with channel evidence, limitations, usage and exact owner-only npm/core commands (2.15.0 before 2.16.0 if still unpublished).

Delivery root: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-participant-dropout`. No release channel mutation or global installation has occurred at this preparation point.
