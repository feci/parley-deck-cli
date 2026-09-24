## Parley Deck — global roster and organizer

Owner default set 2026-09-25 (supersedes the 2026-09-09 default); the same user-global default for Claude and Codex:

- For new Parley Deck runs, the organizer/facilitator is `claude-1`, using Claude CLI with `claude/claude-opus-5-5[1m]` (Opus 5.5). Declare it in each idea with `facilitator: claude-1`: since parley 1.49.0 a declared facilitator is a pure organizer that does not implement, verify code, or sign off.
- The implementer is `codex-1` (Codex CLI, `gpt-6-astra`). Until the designated-implementer mechanism ships (idea `meta-protocol-change-designated-implementer`), the protocol's default implementer is still the FINAL drafter, so the organizer has `codex-1` take implementation through the protocol's own claim path (`inbox/codex-1-to-all_<slug>_impl-claim.md`) before Phase 5. When the mechanism ships, set `codex-1` as the global default implementer.
- The quorum participants are exactly `codex-1`, `kimi-1` (Kimi CLI, `kimi-code/k3`) and `zcode-1` (Zcode CLI, `zai/glm-5.3`). Code review falls to the non-implementers, `kimi-1` and `zcode-1`. No other agent participates by default.
- Active global membership and model settings live in `~/.parley/agents.toml`; verify them with `parley roster show --scope machine`. `claude-1` is inactive there on purpose: it organizes and is not a quorum member. Inactive historical rows are not participants.
- If a run starts from Codex, hand organization to `claude-1`; Codex stays the implementer and participant rather than silently becoming the organizer.
- The organizer assignment is an orchestration instruction; the roster CLI does not select or enforce a facilitator. Do not invent a TOML organizer setting or claim an existing session has switched models.
- This global default does not rewrite existing deck membership, historical idea artifacts, signatures, or runs already in flight (the codex-1-organized `meta-protocol-change-designated-implementer` and `windows-portability` runs started 2026-09-24).

