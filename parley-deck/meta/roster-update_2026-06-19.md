# Roster update — 2026-06-19

- **Date:** 2026-06-19
- **Performed by:** claude-1
- **Authority:** direct user maintenance instruction
- **Type:** roster change

## Summary

Set the Parley Deck roster to a fixed **4-participant quorum** with stable agent IDs
`claude-1`, `codex-1`, `hermes-1`, `antigravity-1` (renamed from the prior
`claude` / `codex` / `hermes` / `agy`). Updated both `COOPERATION.md` §2 tables (roster
+ host-handle) and the backend map `meta/headless-agents.local.json`.

| Agent ID        | CLI      | Role                      | Model                      | Headless | Write mode                         |
| --------------- | -------- | ------------------------- | -------------------------- | -------- | ---------------------------------- |
| `claude-1`      | `claude` | `facilitator+participant` | `opus`                     | `-p`     | `--permission-mode bypassPermissions` |
| `codex-1`       | `codex`  | `participant`             | `cli-default`              | `exec -C "$(pwd)" -s workspace-write - < <promptfile>` | `-s workspace-write` |
| `hermes-1`      | `hermes` | `participant`             | `GLM 5.2`                  | `-z`     | `--yolo --accept-hooks`            |
| `antigravity-1` | `agy`    | `participant`             | `Gemini 3.5 Flash (High)`  | `-p`     | `--dangerously-skip-permissions` (always `--add-dir "$(pwd)"`, `--print-timeout 10m`) |

Host handles (transport is `github-pr`): all four map to **`feci`** (operator-chosen).

## Notes

- **Quorum is now 4** (`claude-1`, `codex-1`, `hermes-1`, `antigravity-1`) from the
  next idea onward (§5).
- **Model strength order** (strongest → weakest): `claude-1` > `codex-1` > `hermes-1` >
  `antigravity-1`. When an idea uses an advisory `roles:`/lens map, assign lenses by
  criticality in that order (strongest takes the most demanding/adversarial lens);
  `claude-1` leads final synthesis regardless. Lenses are per-idea and do **not** change
  quorum, signoff weight, or artifact ownership.
- `hermes-1` is on **GLM 5.2** (not grok). `antigravity-1` uses the `agy` CLI with
  Gemini 3.x and MUST always pass `--add-dir "$(pwd)"` or it hangs; never the `gemini`
  CLI / gemini-2.5.
- No phases, quorum rules, artifact shapes, or transport were changed by this update.
