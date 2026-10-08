# Release delivery evidence — 2026-10-08

Organizer-owned evidence for CLI 1.51.0 and skill/core 2.15.0. The separate
claude-1 delivery report supplies the independent verdict. No implementation or
existing reviewer artifact is changed by these delivery records.

The complete local delivery directory is
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-08-quota-auto-exclude`.
It contains build assets, tag archives, validation logs, exact npm payload,
workflow state, WinGet checkout and the private pre-install backup. Raw participant
transcripts and pre-install payloads are not committed here.

The command-line verifier initially selected eight executable-detected runtimes;
`--target all --include-undetected` then updated and checked all 15 existing managed
installations. Runtime hashes and 2.15.0 markers passed; full core/addon doctor passed.

Homebrew audit output is empty on success. Formula tests include actual skill
installation and doctor, alongside CLI version/help. Name-based strict online audit
passed after the first path-based invocation was rejected by modern Homebrew.

WinGet YAML uses schema 1.12.0 and final workflow-produced Windows EXE hashes.
Native winget validation/install was unavailable locally; actual upstream state is
captured at finalization. CLI WinGet remains held. npm and global core are staged
for owner-only publication; their hashes are evidence, not a published claim.

The committed Homebrew test log normalizes terminal CR line endings and trailing
spaces; the original command output remains in the local delivery logs.
