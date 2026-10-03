# Release preparation — quota auto-exclude

Status: planning only, 2026-10-04. Implementation, independent review, both review-consensus signoffs and
the owner's attended-close answer remain prerequisites. No version bump, merge, tag, push, channel write,
installation or core publication has occurred for this idea.

## Current inventory

Read-only GitHub and git observations at 2026-10-03 23:04Z:

- CLI latest 1.50.0, remote main a8634cc14709614bdb67862d86e80ce41fb53251.
- Skill latest 2.14.0, remote main a5664d803f1fb6156ef95ff5921dcf13a3dd2031.
- No remote v1.51* / v2.15* tags found. Recheck before choosing release numbers.
- Homebrew checkout at 9855366, clean. Both Formula/parley-deck-cli.rb and Formula/parley-deck-skill.rb
  must change. Canonical tap: feci/homebrew-parley.
- Published local core directories: 2.10.0 and 2.13.0. Core 2.14.0 is not published locally; staged
  ~/.parley/staging/COOPERATION-2.14.0.md exists and is the brief-authorized base for core 2.15.0.

## Delivery after the attended close

1. Recheck latest release numbers and main ancestry. Integrate latest main if necessary and rerun checks
   if product content changes. Direct main merge, no development PRs.
2. Prepare the next CLI/skill metadata (planned 1.51.0 / 2.15.0) and frozen release evidence. CLI assets:
   Darwin/Linux/Windows arm64 and x64, plus SHA256 manifest. Windows remains experimental; CLI winget held.
3. Publish GitHub releases with final assets. For skill, require the portable workflow's final Windows
   executables and use those final hashes in winget, never local replacements.
4. Bump both Homebrew formulae with final archive hashes; run style/audit/tests and verify installed versions.
5. Prepare a skill-only winget version directory on a fresh isolated branch; one application per PR.
   Respect the destination repository instructions and report actual Windows validation status honestly.
6. Run skill npm test, pack dry-run, portable build, installer dry-run/doctor, and pack the exact tarball.
   npm publication needs the owner's web OTP: supply the exact `! npm publish <absolute-tarball> --access
   public` command. Do not substitute an OTP workaround or claim npm published without registry evidence.
7. Install the exact skill payload into all existing local runtimes, including dormant managed copies,
   and retain the per-runtime SKILL.md hash table for independent claude-1 verification.
8. Build ~/.parley/staging/COOPERATION-2.15.0.md from staged 2.14.0 plus exactly this idea's reviewed
   normative hunks. Preserve core template zones. Verify the diff and hashes independently, then supply
   the exact attended `! parley protocol publish --version 2.15.0 --from <absolute-file>` command.
   Never allocate a TTY or otherwise bypass that gate.
9. Separately invoke claude-1 with a short channel-evidence brief and exact output path. Its independent
   report must cover every channel and distinguish completed actions, external PR state and owner-only
   commands. Resolve findings before reporting completion.
10. Write inbox/codex-1-to-user_meta-protocol-change-quota-auto-exclude_released.md with shipped versions,
    evidence, owner-only commands, deferred items (including D6), and cumulative/nonduplicated usage.

The skill repository's RELEASING.md and package.json define its checks. The CLI's hosted workflow uses
`go build ./...` and `go test ./... -count=1 -timeout 45m` on Linux, macOS and Windows; Windows failures
must be assessed as actual evidence, never silently omitted or relabeled as passes.
