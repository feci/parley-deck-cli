# Release assessment: CLI 1.48.0 and skill 2.12.0

Date: 2026-09-18. Organizer: codex-1.

The user directed the prepared implementation to be merged and released directly,
without new development PRs. The two prepared branch heads are included in this
release history: integration `6f885d793131cb910d3a3e6f756a86d75e6e1837` and v2
`02b220182523f461d93ecf2e47cb44ce627c2ce3`.

This is a software release, not closure of the empirical audit. Its amendment,
12-task pilot, packet experiment and longitudinal follow-ups remain unfinished.
Missing historical worktrees remain explicitly unknown. No historical FINAL or
participant signature is changed by this release.

## Independent review

Claude and Zcode wrote the accompanying reports themselves. Their reports retain
their original scope and observations, including findings subsequently corrected.
Kimi was invoked for a focused admission/recovery review but reached the 900-second
limit without producing a report; that attempt supplies no review or signoff.

- Corrected the skill Core Rule to use renderer-first context, including a clear stop
  for reachable renderer errors without attestation (Zcode MAJOR-1 and MINOR-1).
- Restored portable header placeholders, added an advisory CLI recommendation,
  clarified omitted fallback reasons, removed a stale dispatch comment and linked
  previous release notes (Claude and Zcode documentation findings).
- Changed portable builds to invoke pkg through Node directly, avoiding Windows
  `.cmd` execution without a shell (Claude N4).
- The applicability map already states that conditions combine with any-match OR;
  its conservative inclusion semantics are unchanged (Zcode map observation).
- R3 remains a disclosed timing-window NIT: a concurrent scope edit can consume a
  verification ticket before later refusal. Charges remain retained. Claude's
  concurrence is a release opinion, not historical amendment ratification.

## Validation

- Full current-source `go test -count=1 -timeout=30m -json ./...` passed: 2,627
  test cases including subtests; 31 tested packages and one package with no tests.
- `go vet ./...` passed.
- App verifier recovery/relaunch tests passed under the race detector.
- Full budget/trajectory/runner/telemetry/protocolpacket race check passed.
- Installer: 391 Node tests, 54 Python tests, and all payload manifest checks passed.
- npm dry pack and actual tarball inspected; no dependency tree, Git or runtime
  state shipped. Actual npm-tarball and portable installs plus doctor passed for
  all supported targets in isolated projects.
- All six CLI and five existing portable installer targets built; host binaries
  reported the intended versions. Final CLI binaries are rebuilt after committing.
- `parley protocol packet check` passed; full context remains the default.

The first full Go run exposed a macOS test-fixture path comparison (`/var` versus
`/private/var`). The fixture now uses the same canonical-root identity as production;
no production check was weakened. The fixed lifecycle test and subsequent full run
passed. The large trajectory package exceeded the initial default 10-minute package
timeout while progressing; the successful full run used an explicit 30-minute cap.

Homebrew's byte-preserving shebang restoration is retained through its current
`post_install_steps` API. A simulation rewrote ten payload shebangs, ran the actual
serialized Homebrew step and confirmed all 198 payload file hashes were restored.
Style checks pass. Final archive hashes and channel checks are recorded in the
published release evidence after tagging.

## Distribution

GitHub releases contain the supported portable binaries. The existing npm package
is `parley-deck-skill`; this release does not invent a new CLI npm package. Homebrew
formulas live in `feci/homebrew-parley`. WinGet IDs are `Feci.ParleyDeckCli` and
`Feci.ParleyDeckSkill`; catalog availability requires Microsoft's submission process.
Preparation, submission and availability are reported separately. npm authentication
and the user's no-PR restriction must be resolved for those affected channels;
neither a local tarball nor prepared WinGet YAML is called a publication.
