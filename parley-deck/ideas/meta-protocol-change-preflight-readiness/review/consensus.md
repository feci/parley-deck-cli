---
idea: meta-protocol-change-preflight-readiness
review-cycle: 1
drafted-by: claude-1
date: 2026-06-19
reviewed-commit: c9e872d
---

Phase 6 reviewers: codex-1 + hermes-1 (claude-1 = implementer; agy waived — same
signoff-append hang as in consensus). Both reviews confirmed the core safety machinery
is correct (probe lifecycle/process-group kill, source fail-closed no-write, gates never
auto-answered, unattended hard-stop never reads stdin) but found a CRITICAL plus several
MAJOR faithfulness/fail-open gaps. Fixing all agreed items before release.

## Agreed fixes
- **[CRITICAL] §1 bypass via `--participants`** (codex-1). `runTask` passes the full
  `discovered` set to `runTaskPreflight`, which re-expands via
  `selectedParticipants(discovered)` and checks ≥2 against that — but the idea is
  created with the `--participants` subset. So `parley run --participants codex` can
  pass preflight and then open a 1-participant idea with no solo exception. **Fix:** pass
  the exact selected participant IDs into preflight; enforce the §1 ≥2 hard-stop on the
  set that will be written to `00-prompt.md`.
- **[MAJOR] `--yes` confirm path is dead; exclusions/role-backfill not recorded**
  (codex-1 + hermes-1). `opts.Yes` is set but never read; the advertised
  `parley preflight … --yes` re-emits the same gate; no `excluded:` / role backfill is
  written. **Fix:** make `--yes` actually confirm — suppress the exclude/re-include and
  unknown-role gates (treat as confirmed), proceed with exit 0 when ≥2 participants
  remain, backfill `protocolRole`, and plumb confirmed exclusions into `00-prompt.md` at
  idea creation (`runcontrol.Create`). (If a full plumb is too large this cycle, at
  minimum: `--yes` clears the standalone gates + records the role backfill, and the run
  path writes the confirmed `excluded:` list into the created `00-prompt.md`.)
- **[MAJOR] `--json` masks exit codes** (codex-1). `runPreflight` returns `printJSON(...)`
  (always 0), so a pending gate (3) / hard failure (1) becomes exit 0 for JSON callers.
  **Fix:** print the JSON payload, then return the preflight `code` (only an encode error
  is a separate failure).
- **[MAJOR] Missing `version.json` / non-workspace fails open** (codex-1). Absent
  metadata is treated as advisory (no gate); `parley init` doesn't write
  `protocolRole:"consumer"`; an empty dir can report ready. **Fix:** `parley init` writes
  consumer metadata; validate the workspace exists before reporting ready (else exit 1);
  absent metadata in an existing deck → the one-time role/backfill gate.
- **[MAJOR] Hosted-PONG substring match** (codex-1). `strings.Contains(out, "PONG")`
  passes on echoed prompts / "cannot return PONG". **Fix:** accept only trimmed stdout
  exactly `PONG` (or a single exact `PONG` line); add false-positive tests
  (echoed-prompt, commentary).
- **[MAJOR] `classifyBump` fails open** (hermes-1). Returns additive on unparseable
  versions → could auto-sync a breaking consumer bump. **Fix:** return `bumpMajor` (gate)
  on parse failure; wire/remove the dead `bumpUnknown` constant.

## Deferred follow-ups
- **Roster-ID ↔ runtime-ID `-1` reconciliation** (both, [MINOR]). Reports use runtime
  IDs (`codex`) while §2 now uses `-1` IDs; the `-1` names aren't in any source the Go
  runtime reads. Separate follow-up idea (it spans the roster-naming layer, not just
  this tool). `TBD`.
- **MINOR perf** (hermes-1): skip the `parley-deck-skill status` shell-out for
  `role==source` + `--no-ping`; differentiate per-probe vs global deadline. `TBD`
  (optional polish; not blocking).

## Dismissed findings
- None outright; the two MINORs above are deferred (real but non-blocking), not dismissed.

## Signoffs
<!-- Each active participant (implementer included) APPENDS their signoff block. -->

### Signoff: claude-1 — 2026-06-19
Status: ✅ ACCEPT
Notes: Implementer. All 6 agreed fixes applied + tested in fix-up cycle 1 (verified by
me): [CRITICAL] §1 hard-stop now evaluated on the selected `--participants` set that
`runcontrol.Create` writes (`TestRunParticipantsSubsetHardStopsSolo` passes);
`--yes` confirms gates + records `excluded:` into `00-prompt.md` via
`CreateOptions.Excluded`; `--json` returns the real exit code; `parley init` writes
`protocolRole:consumer`, non-workspace→exit 1, absent-metadata→gate; hosted-PONG uses
`isExactPONG` (echoed-prompt/commentary false-positives tested); `classifyBump`
fail-closed to `bumpMajor`. `go build` + full `go test ./...` + drift guard green; real
`parley preflight` re-verified. The two MINORs (roster-ID `-1` reconciliation; perf) are
deferred follow-ups per this consensus. Zero remaining agreed fixes.

### Signoff: codex-1 — 2026-06-19
Status: ✅ ACCEPT
Notes: Verified the accepted codex-1 CRITICAL and MAJOR findings are resolved: preflight now evaluates the exact selected participant set, `--yes` records confirmed exclusions, JSON preserves the real exit code, missing metadata/non-workspaces fail closed, and hosted-PONG uses exact matching. The hermes-1 `classifyBump` parse-failure issue also fails closed to `bumpMajor`, with targeted tests covering the fixes.

### Signoff: hermes-1 — 2026-06-19
Status: ✅ ACCEPT
Notes: Verified both my MAJORs are resolved in code + tests. `opts.Yes` is now read: unavailable agents → `report.Excluded` as `<roster-id> — reason — confirmed <date>`, unknown/absent role → `protocolRole=consumer` backfill, and `CreateIdeaWithExclusions` plumbs the confirmed `excluded:` list into `00-prompt.md` via runcontrol (the §1 <2-participant hard-stop stays non-waivable). `classifyBump` fails closed to `bumpMajor` on unparseable versions and the dead `bumpUnknown` constant is gone (breaking bump now gates, never auto-writes); table-driven test asserts parse-failure→`bumpMajor`. `go build` + full `go test ./...` + drift guard green. The two MINORs (roster-ID `-1` reconciliation, perf) remain deferred follow-ups per this consensus — non-blocking. Zero remaining agreed fixes.
