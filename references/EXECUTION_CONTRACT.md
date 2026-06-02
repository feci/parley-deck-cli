# Pipeline execute contract (§12.10)

How a side-effecting `action` block is executed. The boundary is strict
(COOPERATION.md §12.4): **agents write markdown; the CLI plans and ledgers; an
external driver/harness performs the actual provider call via MCP**. The Go CLI
never mutates a provider itself.

## Roles

- **Agents** deliberate the action plan (Phase 1-4) and finalize it.
- **CLI** (`parley pipeline execute`) validates preconditions, enforces gates,
  computes the idempotency key, writes the per-effect ledger entry, and emits
  the concrete `ProviderCall`. It performs no side effect.
- **Harness/driver** reads the emitted plan, performs the MCP call when the gate
  permits, then calls `parley pipeline record-effect` with the outcome.

## Flow

1. `parley pipeline execute --json SLUG BLOCK CAPABILITY TARGET` (add `--dry-run`
   to plan without intent to mutate). Output (schema_version 1):

   ```json
   {
     "schema_version": 1,
     "status": "dry_run | pending_gate | ready_for_harness",
     "provider_call": { "provider": "vercel", "tool": "mcp__...", "args": {}, "dry_run": false },
     "effect_digest": "<digest>",
     "idempotency_key": "<full key>",
     "gate": { "required": true, "state": "approved | open | not_required" }
   }
   ```

   - `pending_gate` → a production gate is open; approve with
     `parley pipeline gate approve SLUG BLOCK->execute`, then re-run. The CLI
     wrote no executing ledger entry.
   - `dry_run` → the effect is recorded `dry_run_ok`; no mutation should happen.
   - `ready_for_harness` → the effect is recorded `executing`; the harness must
     now perform `provider_call` via the named MCP tool.

2. The harness performs `provider_call.tool` with `provider_call.args` via MCP
   **only** when `gate.state == "approved"` (or the capability is non-production
   / dry-run). Production capabilities (`*.production`, `*.rollback`) are
   non-bypassable and must never run without an approved gate.

3. The harness records the outcome:

   ```
   parley pipeline record-effect --status succeeded --external-ref <id> SLUG <effect_digest>
   parley pipeline record-effect --status failed   --external-ref <id> SLUG <effect_digest>
   ```

## Idempotency & reconcile (§12.7)

- The same `idempotency_key` always maps to the same `effects/<digest>.json`, so
  re-planning the same action never creates a duplicate.
- An effect left `executing` (harness crashed before `record-effect`) or
  `failed` is ambiguous: a fresh `execute` is blocked until the harness
  reconciles external state and records a terminal status. Retrying without the
  idempotency key is prohibited.

## Vendor neutrality

`provider_call.provider`/`tool` come from a `Provider` implementation behind a
generic capability vocabulary (`deploy.preview`, `deploy.production`,
`runtime.rollback`, `monitor.alert`, `issue.create`, `notify.send`). Vercel is
the first implementation; the protocol depends on no specific vendor.
