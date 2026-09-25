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
