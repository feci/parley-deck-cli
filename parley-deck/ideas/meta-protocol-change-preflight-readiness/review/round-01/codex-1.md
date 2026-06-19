---
agent: codex-1
idea: meta-protocol-change-preflight-readiness
review-round: 1
date: 2026-06-19
reviewed-commit: c9e872d
---

## Summary
The run-path hosted-PONG default is acceptable and matches the operator ruling, and the source-role freshness path does avoid writing `COOPERATION.md`. The implementation still has blocking faithfulness gaps: the run gate validates a different participant set than the idea will use, the advertised confirmation path is not implemented, and several exit/freshness/availability paths fail open.

## Findings
### [CRITICAL] `parley run --participants ...` can bypass the §1 two-participant hard stop
`runTask` computes the actual idea participants from `--participants` at `internal/app/app.go:1713`, but it calls `runTaskPreflight` with only `discovered` at `internal/app/app.go:1728`. `runTaskPreflight` then re-expands to every found runtime via `selectedParticipants(discovered)` at `internal/app/preflight.go:185-187`, while `runcontrol.Create` later creates the idea with the original `participants` at `internal/app/app.go:1736-1740`. On a machine with two or more installed agents, `parley run --participants codex ...` can pass preflight against the full installed set and then create a one-participant `00-prompt.md` with no recorded solo exception. Fix by passing the exact selected roster/participant IDs into preflight and enforcing `>= 2` on the set that will be written to `00-prompt.md`, unless an explicit solo exception is recorded before `runcontrol.Create`.

### [MAJOR] The advertised `--yes` confirmation path is dead and records no exclusions
`parley preflight` parses `--yes` into `preflightOptions.Yes` at `internal/app/preflight.go:112-130`, but no later code reads `opts.Yes`. All gates print `parley preflight --dir ... --yes` as the confirm command (`internal/app/preflight.go:244-248`, `299-303`, `317-321`), yet rerunning it recomputes the same gate and writes no `excluded:` entry, no role backfill, and no confirmed freshness decision. This fails FINAL.md acceptance criterion 2 and §9.0's requirement that exclusions/re-inclusions be explicit and recorded in the idea. Fix by implementing a real gate-confirm flow and plumbing confirmed exclusions/freshness decisions into `runcontrol.Create`/`CreateIdea` before the idea is opened, or stop advertising `--yes` as a confirm command until it does that work.

### [MAJOR] `parley preflight --json` masks pending-gate and hard-failure exit codes
When `--json` is set, `runPreflight` returns `printJSON(...)` directly at `internal/app/preflight.go:143-145`; `printJSON` returns `0` on successful encoding. That means a pending gate (`3`) or hard failure (`1`) becomes process exit `0` for JSON callers, which breaks the specified automation contract. Fix by printing the JSON payload, then returning the preflight `code` unless JSON encoding itself fails.

### [MAJOR] Missing version metadata and non-workspaces fail open instead of producing the required gate/hard failure
`classifyAndSyncFreshness` treats a missing `parley-deck/meta/version.json` as advisory and returns no gate at `internal/app/preflight.go:275-280`; `protocol.InitWorkspace` also does not create a `version.json` with `protocolRole: "consumer"` (`internal/protocol/workspace.go:37-58`). As a result, fresh consumer workspaces never enter the required role/backfill path, and standalone `parley preflight --dir <empty-dir> --no-ping` can report ready instead of exit `1` for "no workspace" when installed runtimes are present. Fix by making `parley init` write consumer metadata, validating that the workspace exists before reporting readiness, and treating absent metadata in an existing deck as the one-time role/backfill gate described in FINAL.md.

### [MAJOR] Hosted-PONG accepts any stdout containing `PONG`, not the exact sentinel
`hostedPONG` marks an agent available with `strings.Contains(out.String(), pongSentinel)` at `internal/app/preflight.go:616-618`. The spec requires the exact sentinel; a successful process that echoes the prompt, prints explanatory text such as "cannot return PONG", or includes `PONG` in any other output will be classified available. Fix by accepting only trimmed stdout equal to `PONG` (or a single exact `PONG` line with no other non-whitespace output) and add tests for echoed-prompt/commentary false positives.
