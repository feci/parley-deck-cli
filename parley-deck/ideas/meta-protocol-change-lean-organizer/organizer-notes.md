# Organizer execution notes

## Phase 0 — 2026-09-23

- Organizer: codex-1; participants: claude-1, kimi-1, zcode-1, verified with `parley roster show --scope machine`. codex-1 inactive is intentional.
- Skills applied: parley-deck; supplied isolated worktrees follow parley-worktrees mechanics. OpenViking MCP tools are not available; local evidence is authoritative for this run.
- CLI baseline v1.48.0; skill installer/runtime v2.12.1. Source-role protocol drift is advisory. No protocol freshness mutation was performed.
- Organizer renderer attestation: context_mode=full, source_sha256=12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18, packet_sha256=12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18, fallback_reason absent. Reading was restricted to the facilitator reading set plus relevant §7 under the owner's explicit brief; no claim of full-body reading. Participant launches receive their own renderer context.
- Driver gap: `parley run` always creates a timestamp-derived slug, has no existing-idea/explicit-slug input, and creates a generic prompt. `parley continue` requires a run. The facilitator therefore seeds only kickoff/run metadata manually for the exact owner-required slug, then uses `continue --auto --no-implement`. No participant content is authored by the facilitator.
- Initial hosted readiness: Claude and Zcode ready, Kimi malformed-reply. A direct Kimi invocation returned a decorated PONG with banner/reasoning/session footer. A local runtime override requests stream-json; no model or effort change and no quorum exclusion.
- Direct-main/no-development-PR release is the owner's explicit per-run override. Existing unrelated historical ideas are left untouched.

- Kimi readiness resolved manually: stream-json emits exactly `{ "role": "assistant", "content": "PONG" }` between meta envelopes (source-context/kimi-readiness.jsonl). The current preflight parser rejects this real envelope. Quorum remains unchanged; this is a documented parser fallback, not a bypass of an unavailable agent. Runtime Kimi stream-json override is /tmp/lean-organizer-runtime.toml.
- Run metadata manually seeded as 20260923T202501.377412000Z with explicit frozen roster and launch args; no synthetic participant events or verdicts. The driver owns subsequent launches.
