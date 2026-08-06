---
idea: kimi-opencode-full-adapters
author: user
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: standard
status: implementation
created: 2026-08-06
---

## Decision (made by the user, not deliberated)

Promote `kimi` and `opencode` from **ACP-only catalog stubs** to **full built-in adapters with a
declared autonomous-write mode**, so `parley` can launch them headlessly instead of requiring a
hand-written Bash call.

Two questions were put to the user and answered:

1. **Does the opencode adapter pass `--auto` explicitly, or rely on the default?**
   → **Explicitly.** Rationale accepted: it keeps `AutonomousWrite.Args` a subset of
   `HeadlessArgs` as with every other adapter, and an implicit default is precisely what a vendor
   may change between versions. Cost acknowledged: `--auto` is documented as broader
   ("dangerous!") than the deck strictly needs.
2. **Direct release, or Parley Deck review first?** → **Review first**, because the change alters
   how agents are launched unattended with auto-approve.

## Why the change is needed

`AUTO=no` on both agents was never a statement about their capability. `internal/agents/discover.go`
computes the column as `AutonomousWrite.Declared()`, i.e. "is `Mode` non-empty". kimi and opencode
were defined **only** in the ACP catalog (`acp_specs.go`), and `specFromACPBackend()` never sets
`AutonomousWrite` — so the field was empty and the column read `no`, while both CLIs in fact write
unattended.

## Constraints

- `Scope` stays **empty** for both. The type's own rule: `workspace` may be claimed only where the
  CLI enforces a real sandbox. Only `codex --sandbox workspace-write` earns it.
- ACP must remain available as an alternative launch mode for both.
- No change to any other adapter.
- English only.
