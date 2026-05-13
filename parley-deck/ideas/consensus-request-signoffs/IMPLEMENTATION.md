---
idea: consensus-request-signoffs
implemented-by: codex
status: ready-for-review
date: 2026-05-13
implementation-pr: https://github.com/feci/parley-deck-cli/pull/12
branch: feature-consensus-request-signoffs
base: main
design-pr: https://github.com/feci/parley-deck-cli/pull/11
---

## Summary

Implemented the `parley consensus request-signoffs` CLI command from `FINAL.md`.

## Changes

- Added `parley consensus request-signoffs`.
- Added support for:
  - `--dir`;
  - `--review`;
  - `--participants`;
  - `--yes`;
  - `--dry-run`.
- Added participant target resolution:
  - explicit ordered target set via `--participants`;
  - default to missing signoffs in `00-prompt.md` participant order.
- Added preflight validation for unknown, already-signed, missing-runner, and non-installed participants.
- Added hosted/non-local backend confirmation gating through existing discovery metadata.
- Added dry-run output with target file, current status, selected participants, backend type, launch order, command preview, and `--yes` requirement.
- Added sequential agent invocation through `runner.CommandFor`.
- Added append-only signoff prompt generation with absolute paths.
- Added post-invocation validation using `internal/consensus` to stop on:
  - missing signoff;
  - duplicate or malformed signoff;
  - changed foreign signoff block;
  - `BLOCK`;
  - child process failure, including non-zero exit after a valid append.
- Added focused fake-CLI tests for happy path, dry-run, hosted gate, explicit already-signed rejection, review path, non-zero-after-append failure, and BLOCK stop.

## Deviations from FINAL.md

- The implementation uses the already-parsed participant order from `protocol.ReadWorkspaceStatus`, which reads `00-prompt.md` frontmatter. No new participant parser was needed.
- Durable per-agent run logs remain deferred. Child stdout/stderr stream to the parent command.
- The implementation branch uses `feature-consensus-request-signoffs` instead of protocol-preferred `feature/consensus-request-signoffs` because local Git could not create nested refs in this sandbox. The exact failed command was `git checkout -b feature/consensus-request-signoffs`; escalation retry was rejected by the harness under `on-failure`.

## Verification

- `GOCACHE=/private/tmp/parley-go-build-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go test ./internal/app -run 'TestConsensusRequestSignoffs'`
- `GOCACHE=/private/tmp/parley-go-build-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go test ./...`
- `GOCACHE=/private/tmp/parley-go-build-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go run ./cmd/parley consensus request-signoffs --dry-run consensus-request-signoffs`

## Ready for review

Review should focus on:

- participant selection semantics;
- hosted/non-local `--yes` gate behavior;
- prompt shape and ownership guarantees;
- post-invocation validation strictness;
- error codes for usage/safety vs runtime failures;
- whether the command should add durable run logs in a later slice.
