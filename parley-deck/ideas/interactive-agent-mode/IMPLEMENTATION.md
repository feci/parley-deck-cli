---
idea: interactive-agent-mode
implementer: codex
status: complete
date: 2026-05-14
implementation-pr: https://github.com/feci/parley-deck-cli/pull/14
branch: feature/interactive-agent-mode
base: main
design-pr: https://github.com/feci/parley-deck-cli/pull/13
head-commit: 30d5fe6
fix-up-cycle: 1
---

## Summary

Implemented the first staged slice from `FINAL.md`: launch-mode configuration, visibility, handoff packets, stricter artifact validation, and interactive/manual consensus signoff handling.

## Changes

- Added agent launch-mode fields:
  - `launch_mode`;
  - `interactive_mode`;
  - `interactive_command`;
  - `interactive_args`;
  - `interactive_prompt_mode`;
  - `interactive_invoke`;
  - `interactive_timeout_ms`;
  - `interactive_poll_ms`;
  - `interactive_notes`.
- Kept `headless` as the default for all built-in agents.
- Updated `agents list` runtime output to show the resolved launch mode and interactive command shape when relevant.
- Added `--mode` override to `parley consensus request-signoffs`, supporting both bare `--mode interactive` and per-agent `--mode claude=manual` forms.
- Changed hosted/non-local `--yes` gating so it applies to headless launches, not manual/interactive handoff preparation.
- Added handoff packet generation under `parley-deck/runs/<run-id>/agents/<agent>/` with:
  - `handoff-prompt.md`;
  - `handoff.md`;
  - target artifact path;
  - validation contract;
  - provider-agnostic usage caveat.
- Added interactive signoff support:
  - `print-only` writes handoff files, prints instructions, polls for the signoff, and then uses existing append-only consensus validation.
  - `spawn-tty` starts the configured interactive command attached to the parent terminal when a TTY is available.
- Added manual signoff support:
  - writes handoff files;
  - does not invoke the agent;
  - records pending handoff events;
  - exits with visible pending status.
- Added `resume --no-tui` validation for pending consensus-signoff handoffs when the signoff has been appended later.
- Added stricter round artifact validation after headless round runs: frontmatter identity plus required Round 1 sections.
- Documented launch modes and recommended local interactive config in `docs/agent-runtime-configuration.md`.

## Deviations from FINAL.md

- Mixed-mode round execution is not wired yet. This was explicitly staged after signoff support in `FINAL.md`.
- Manual mode currently exits successfully after creating pending handoff instructions instead of using a distinct pending exit code. The exact exit code remains a deferred implementation detail.
- `resume` currently validates pending consensus-signoff handoffs. Full generic resume behavior for all pending run types remains deferred.
- `spawn-tty` has a basic TTY gate and attached process launch. Process-group signal handling polish remains deferred.
- Strict validation is currently applied to Round 1 artifacts produced by the runner. Later cross-review round validation should be added with the mixed-mode round execution slice.

## Verification

- `GOCACHE=/private/tmp/parley-go-build-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go test ./...`
- `GOCACHE=/private/tmp/parley-go-build-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go run ./cmd/parley agents list --dir .`
- `git diff --check`

## Ready for review

Review should focus on:

- whether the launch-mode config surface matches `FINAL.md`;
- whether hosted confirmation is still correct for headless vs interactive/manual flows;
- handoff packet contents and paths;
- append-only consensus validation after interactive polling;
- manual pending/resume behavior;
- whether the stricter Round 1 artifact validation is too broad for existing headless runs.

## Fix-up cycle 1

Applied agreed fixes from `review/consensus.md`:

- store `before_len` and `before_sha256` on signoff handoff events and make `resume` verify append-only content before marking pending signoffs complete;
- unify placeholder expansion for handoff instructions, dry-run previews, and `spawn-tty` args;
- change manual handoff instructions to point at `parley resume <run-id>`;
- return exit code `3` for pending manual handoffs and document it;
- rename the strict round artifact validator to `ValidateRoundOneArtifact` so its scope is explicit;
- document why manual mode validates shared `interactive_*` handoff fields;
- remove the unreachable `runInteractiveTTY` command fallback;
- change the `IMPLEMENTATION.md` frontmatter key to `implementer`.

Verification:

- `GOCACHE=/private/tmp/parley-go-build-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go test ./...`
- `git diff --check`

Ready for re-review of fix-up cycle 1.

## Review closeout

Review round 02 approved fix-up cycle 1 with no remaining agreed fixes. The cycle-1 consensus with agreed fixes is preserved in `review/consensus-cycle-01.md`; the current `review/consensus.md` records zero agreed fixes and all participant signoffs.

Implementation is complete.
