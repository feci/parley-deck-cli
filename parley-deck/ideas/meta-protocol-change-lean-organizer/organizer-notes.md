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

- Additional driver startup gap: `continue --auto` on a freshly seeded run waits for initial round artifacts instead of launching them (run has only run.created after startup). First round is therefore hand-launched through the three configured CLIs; the running driver will own subsequent rounds/drafting/signoffs. Manual prompts only for this unsupported initial step.

- Correction to the initial continuation diagnosis: the driver was waiting because two round-01 artifacts used suffixed `## Existing alternatives (...)` headings. The existing validator requires the exact heading. `parley status` showed a complete on-disk round but concealed this validity reason. A disposable runtime helper called the existing Go artifact validator and printed its verdicts; no code tests/review were run by the organizer. Kimi and Zcode were each asked to repair their own heading. The idle driver (no child processes) was stopped with a diagnostic SIGQUIT; it will be resumed after validation, preserving run state and all artifacts.
- Usage boundaries can be reconstructed from driver event timestamps and the latest preceding client token_count event, so a blocking driver does not require organizer polling for accounting. The ledger labels observations and retains client accounting timestamp; this is cumulative client evidence, not model self-report or a pricing claim.

## Phase 2 fallback — 2026-09-23

- All three round-01 artifacts now pass the existing Go validator. The driver reconstructed round.completed and a deterministic round.digest.
- The driver then refused to launch round-02 before any participant launch: cross-review accounting reports unavailable historical worktree `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/f2repo`. Its durable escalation is `inbox/claude-to-user_meta-protocol-change-lean-organizer_driver-error.md` (the driver labels the author `claude`, despite codex-1 organizing this run).
- Owner authorization: “When the driver cannot do a step, fall back manually and record why in the idea (that gap is evidence for this idea).” We take that manual canonical workflow fallback. No worktree history is pruned, no unknown-history declaration is forged, no budget ledger is reset/relaxed, and no product code is changed to bypass the guard. Manual cross-review count starts at 1 (round-02), within deliberation's 3 cross-review limit; all three participants remain required. Driver status shows terminal/outcome completed for round completion while the idea is still pre-consensus; this is not treated as idea completion.
- These are operational observations for A–D, not new mandatory implementation scope. Review participants may discuss whether existing validators/status adequately express these situations under scope B.
