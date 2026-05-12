---
idea: consensus-workflow-cli
drafted-by: codex
date: 2026-05-12
---

## Agreed decisions

- Build a deterministic `internal/consensus` package for consensus file handling. It owns schema selection, scaffold creation, signoff parsing/appending, validation, triage, finalization helpers, and reopen helpers.
- Support design consensus and review consensus from day one with one shared schema model:
  - design consensus path: `parley-deck/ideas/<slug>/consensus.md`;
  - review consensus path: `parley-deck/ideas/<slug>/review/consensus.md`;
  - review consensus uses the same signoff rules, but reports review-specific sections and agreed fixes.
- Add CLI commands under `parley consensus`:
  - `parley consensus status [--dir DIR] [--review] [--json] IDEA`;
  - `parley consensus draft [--dir DIR] [--review] [--round N] [--by AGENT] IDEA`;
  - `parley consensus signoff [--dir DIR] [--review] --agent ID --status accept|reserve|reservations|block [--notes TEXT] [--counter TEXT] IDEA`;
  - `parley consensus finalize [--dir DIR] [--by AGENT] IDEA` for design consensus;
  - `parley consensus reopen [--dir DIR] IDEA --reason TEXT`.
- `draft` writes an empty, structured template and checklist. It must not generate substantive consensus prose or summarize round files by default.
- `draft` updates `00-prompt.md` to `status: consensus` for design consensus. It should refuse or warn clearly when the selected round is incomplete unless the user passes an explicit override already present in the CLI style.
- `status` reports deterministic triage:
  - `ready`: every active participant has a valid ✅ signoff;
  - `reserved`: all participants signed, at least one 🟡 signoff, and no ❌ signoff;
  - `blocked`: at least one ❌ signoff;
  - `partial`: missing signoffs;
  - `malformed`: parse or validation error.
- Status flags are user-friendly aliases, but file output uses the canonical protocol values exactly:
  - `accept` -> `✅ ACCEPT`;
  - `reserve` or `reservations` -> `🟡 ACCEPT-WITH-RESERVATIONS`;
  - `block` -> `❌ BLOCK`.
- `signoff` appends one structured block for one agent. It must validate known participant IDs, reject duplicate signoffs unless an explicit replacement path is added later, and require notes for 🟡 or ❌. A ❌ signoff also requires a counter-proposal.
- `finalize` validates the consensus state and creates `FINAL.md` as a separate scaffold. It never renames or copies `consensus.md` as the final plan. `consensus.md` remains the audit record.
- `finalize` succeeds for `ready`. It may also succeed for `reserved` only when reservations are visible in the consensus open-items section; exact validation can be conservative, but the command must not silently ignore reservations.
- `finalize` updates `00-prompt.md` to `status: final` only after `FINAL.md` has been written successfully.
- `reopen` is deterministic escape handling for blocked consensus. It is valid only for `blocked` consensus, records the reason, preserves the old `consensus.md` under a numbered aborted filename, and moves `00-prompt.md` back to a discussion state. The existing round scheduler remains responsible for creating the next round directory.
- Expose consensus triage in `parley status --idea <slug>` when a consensus file exists. A TUI panel is not required in this slice.
- Emit append-only events for state-changing commands where the existing store/event pattern supports it: `consensus.drafted`, `consensus.signed`, `consensus.finalized`, and `consensus.reopened`.
- Keep GitHub PR integration as a mirror, not the source of truth. The CLI may print the native review action expected by the active `github-pr` transport, but full GitHub API review submission is out of scope for this slice.

## Agreed trade-offs

- Deterministic file workflow is the first priority. Model-generated summaries, natural-language synthesis, and autonomous signoff orchestration are follow-up features.
- A single `consensus.md` remains the canonical signoff file. Separate signoff directories would reduce append conflicts but would weaken the human-readable protocol artifact.
- `parley consensus --review` is preferred over a separate `parley review consensus` command group because it reuses one parser/validator and keeps the dispatcher shallow.
- The first implementation collects sequential signoffs through explicit `signoff` invocations. Automated `request-signoffs` is valuable, but it crosses into agent orchestration, token spend, timeout policy, and write-permission handling; it should reuse the deterministic primitives after they are tested.
- Content hashes such as a `BasedOn:` line are useful for stale-signoff detection but are not required in slice 1. Git history, protocol files, and events are enough for the first implementation.
- Reservations are valid protocol outcomes when logged. The CLI should make them visible and gated rather than treating them as hard failures.
- Concurrency protection should be pragmatic: avoid corrupting append operations and use atomic file replacement for frontmatter/status edits. Exact lock mechanics can be chosen during implementation based on the existing dependency and platform constraints.

## Open items deferred to implementation

- Exact Go API names and data structures inside `internal/consensus`.
- Exact console wording for status tables and validation errors.
- Exact `FINAL.md` scaffold content, as long as it clearly points back to `consensus.md`, round files, and any reservations/open items.
- Exact event payload fields for `consensus.*` events.
- Whether `sign` is added as an alias for `signoff`. The primary documented verb is `signoff`.
- The precise aborted consensus filename created by `reopen`, as long as it is numbered, preserved, and never overwrites an existing file.

## Deferred follow-up ideas

- `consensus-request-signoffs`: `parley consensus request-signoffs [--review] --participants IDS [--yes] [--dry-run] IDEA`, invoking selected agents sequentially and verifying that each appends its own canonical block.
- `consensus-stale-signoff-detection`: optional `BasedOn: sha256:<prefix>` support or equivalent stale-signoff metadata.
- `consensus-github-review-mirror`: native GitHub review submission/approval/request-changes mirroring from consensus signoffs.
- `consensus-generated-summary`: opt-in draft summary generation such as `parley consensus draft --with-summary`.
- `consensus-tui-panel`: live TUI/status panel for consensus triage and missing signoffs.
- `consensus-identity-enforcement`: stronger checks between `--agent`, git identity, and remote review identity when transport tools can prove the mapping.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-12
Status: ✅ ACCEPT
Notes: Accept. The consensus keeps slice 1 deterministic, covers design and review consensus with one shared package, and defers model orchestration until the file primitives are tested.

### Signoff: claude — 2026-05-12
Status: ✅ ACCEPT
Notes: Accept. The draft adopts the slice-1 shape from my round-02 (`draft|signoff|status|finalize|reopen`, `--review` flag, deterministic file work, events for state changes, GitHub as mirror) and renders status values verbatim from COOPERATION.md (`✅ ACCEPT` / `🟡 ACCEPT-WITH-RESERVATIONS` / `❌ BLOCK`), which is strictly better than my round-02 shorthand. Deferring the optional `BasedOn:` content hash to the `consensus-stale-signoff-detection` follow-up is acceptable because the `consensus.*` event stream plus git history give us a sufficient audit layer for slice 1; I will champion the follow-up when the primitives land. The `signoff` vs `sign` verb question is correctly listed as a deferred implementation detail.

### Signoff: gemini — 2026-05-12
Status: ✅ ACCEPT
Notes: Accept. The consensus effectively captures the requirement for a deterministic, file-based workflow that supports both design and review cycles. It preserves the protocol's non-solo nature by ensuring each agent appends its own signoff block via the CLI, while providing clear triage and finalization paths. Deferring the automated 'request-signoffs' to a follow-up ensures the core primitives are robustly tested first.
