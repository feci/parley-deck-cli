# ASD-STE100 output format release — 2026-10-13

Project: parley-deck. Idea: meta-protocol-change-ste100-output-format.
The owner assigned codex-1 as organizer and implementer for this idea only. The usual Claude organizer default is unchanged.

CLI 1.55.0 and skill 2.19.0 are released.
CLI merge 88dc6c5da139d2deb039bda4b605acfbe42663b4:
https://github.com/feci/parley-deck-cli/pull/84
https://github.com/feci/parley-deck-cli/releases/tag/v1.55.0
Skill merge f557c5e679b0c1a03f57a885f41ae54f68348b2a:
https://github.com/feci/parley-deck-skill/pull/12
https://github.com/feci/parley-deck-skill/releases/tag/v2.19.0
Both Homebrew formulae were published in tap commit 8ef3b97bcf97b9d710fa5a26305f010bb9c33e2c. Local upgrades, formula tests, strict audit, and all 202 skill payload byte comparisons passed.
Skill-only WinGet PR 451758 was open at the final check:
https://github.com/microsoft/winget-pkgs/pull/451758
CLI WinGet remains held.

Policy: Appendix B is the binding working summary for new English protocol prose. It cites ASD-STE100 Issue 9, January 2025. Exact quotations, code, paths, frontmatter, required syntax, and mechanical output remain exempt. Agent prose is not exempt merely because an agent generated it. Older live authority wins until authorized adoption. Per-idea ste_output: false needs an English reason; later changes need owner direction or participant consensus. Style feedback is advisory. No dictionary or prose linter ships. This summary does not certify full ASD conformance. The numeric 20/25/6 limits are local adopted rules; full-standard attribution remains unverified.

All three protocol copies changed by exactly H1/H2. Staged core 2.19.0 equals staged 2.18.0 plus those hunks only.
Core SHA256: 16323e04ffb076db3b2019169010033c6668199ebab0ca6896e79924e74bfe3b.
Parley SKILL.md is 19,482 bytes; original packet fixture is 68,874 under 70,000. The meta-idea variant is 71,934 and exceeded the limit before this change.
Thirteen real prompt builders share the instruction.

Review caught a Linux argument-limit regression from the larger full protocol. Fix-up cycle 1 adds exact private-file prompt delivery at both headless boundaries when a prompt argument reaches 120 KiB. Unique mode-0600 files retain full prepared bytes, hash, and attestation. Stdin and short arguments stay exact. File publication does not prove model reading. Real child tests and Linux/macOS evidence passed.

Zcode supplied two source reviews, own ACCEPT signoffs, a fresh source goal PASS, supplemental terminal-Windows evidence, and a separate fresh channel PASS. Codex did not review its own code. Kimi had two measured HTTP403 readiness failures and was excluded under the brief. The mechanical single-reviewer exception returned false; the brief's four-part close authority applied.
Windows remains experimental. Both product CI jobs failed with pre-existing timeout, evidence-test compile failure, and known budget sharing errors. Every current failure name occurs in the union of two pre-change baseline runs. This does not prove equal causes or complete coverage. New runner policy and transport tests were not reached after an older hanging test. The skill Windows assets are an Ubuntu cross-build, not native Windows validation.

All 114 SKILL.md hashes match the six packaged skills across 15 managed targets (including undetected targets) plus four generic locations:
~/.hermes/profiles/{ldx,librade,testprofile}/skills/parley-deck
~/.config/opencode/skills/parley-deck
The real installer was used. Zcode independently re-hashed all 114 files and ran doctor. Doctor used task-local Python 3.14; system Python 3.9.6 is too old for bidding checks.

Pending owner actions, in required order: stop writers and inspect/apply the one-time D6 declaration; then publish npm and global core for each absent version 2.15.0 through 2.19.0 in order. No agent applied D6 or published npm/core. npm latest was 2.13.0; installed cores were 2.10.0 and 2.13.0. The registry also lacks 2.14.0, outside this brief's command range.
D6 preview covers 25 roots, decisions empty. Preview hash d78dfd90744fc3eb43e86726ba4df99e0bf5b5f08a0109bc0d84b9a29da868f1. Re-inspect before owner apply; never reuse a stale digest blindly.
Skill 2.19.0 npm tarball SHA256: f07102c862908b3f6da6fb65f6a01f948fd3b83ba6fafa3d014c51145dce21e4.

Canonical final handoff:
parley-deck/inbox/codex-1-to-user_meta-protocol-change-ste100-output-format_released.md
Audit PR: https://github.com/feci/parley-deck-cli/pull/85
Local delivery: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-13-ste100
Independent channel artifact: parley-deck/ideas/meta-protocol-change-ste100-output-format/source-context/channel-verification-zcode-1.md
It is process 54e2db04-8c86-4cc8-ad05-ff2f901b8ea5, exit 0, 744.227 seconds, artifact SHA256 c3e1ba900b1c2f9f0b55614a2007ab2871f3b7a62d1593c6f30a37bcf95b78cd.
Closed source artifacts stay frozen. Delivery evidence is separate. GitHub COMMENT mirrors use shared account feci and do not imply separate native approvals.
Native token/cost telemetry is unavailable. Codex client counters are cumulative session counters, not a marginal task bill.
