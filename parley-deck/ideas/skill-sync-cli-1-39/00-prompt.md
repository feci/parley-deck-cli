---
idea: skill-sync-cli-1-39
author: user
created: 2026-08-06
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: standard
status: round-01
---

## Problem / idea

`parley-deck-cli` 1.39.0 shipped. **Decide what `parley-deck-skill` must carry as a result**, then
release it. The skill is a separate package with its own channels and it has now been under-shipped
twice in one day, so the owner made two standing rules: every cli/skill update goes through a real
Parley Deck run, and a CLI release always ships the skill in the same turn.

### What changed in the CLI

1. **`kimi` and `opencode` promoted from ACP-only catalog stubs to full built-in adapters** with a
   declared autonomous-write mode. `kimi` launches as `kimi -p <prompt>`; `opencode` as
   `opencode run --auto <prompt>`. Both keep ACP as an alternative launch mode.
   - `kimi --auto -p …` **exits 1** — `Cannot combine --prompt with --auto` — so `-p` is the only
     autonomous headless shape kimi has.
   - `opencode run` writes unattended even without `--auto`; `--auto` is passed **explicitly** by
     user decision, because an implicit default is what a vendor may change between versions.
2. **The `AUTO` signal now fails closed.** A config layer overriding `headless_args` replaces them
   wholesale and never touches `autonomous_write`, so parley could declare an agent autonomous
   while the launched command never passed the enabling flag. **`hermes` was live in that state** —
   an override had dropped `--yolo` while it still reported `AUTO=yes`. `AUTO` now reports `no` and
   names the missing args, and `agents list` prints the **effective** argv instead of a built-in
   label.

### What is stale or missing in the skill (facts, verified)

- `skills/parley-deck/SKILL.md` "Autonomous Execution" table has rows for claude, codex, hermes,
  agy and kimi — **no `opencode` row**.
- The same section's fail-closed sentence covers only *workspace confinement* ("if workspace
  confinement cannot be demonstrated … treat its autonomous bit as unset"). It does **not** cover
  the new case: a declared mode whose enabling flag the effective launch does not pass.
- `skills/parley-deck/references/compatibility.json` says `skillVersion: 1.4.3` while
  `package.json` is `2.3.0` — stale by four releases. Nothing checks it.
- The bundled protocol snapshot `skills/parley-deck/references/COOPERATION.md` is **current** —
  1.39.0 changed no protocol text (`git diff v1.38.0..v1.39.0 -- …COOPERATION.md` is empty, and the
  normalized diff against the live deck is clean). Do not change it without a reason.

## What round 1 must produce

For each item, a position: **carry it into the skill / leave it out (why) / already covered (cite
the text)**. Plus:

1. **Scope discipline.** The skill is instructions for *agents*, not a changelog. Which of these
   facts actually change what a facilitator or participant DOES? Anything that does not should stay
   out, however true it is.
2. **The fail-closed wording.** Propose the exact sentence(s). It must tell a facilitator what to
   check before trusting that a participant can write its artifact.
3. **`compatibility.json`.** Is bumping `skillVersion` enough, or does the silent staleness need a
   guard? Note the constraint: no new tooling unless justified.
4. **Anything else 1.39.0 invalidates** that this brief has not listed. The brief was written by
   the facilitator and is not exhaustive.

## Constraints

- The skill is vendor-neutral instruction text. No new transports, roster entries, or skills.
- `references/COOPERATION.md` is the portability fallback and must stay byte-identical to the
  regenerated CLI bootstrap template — do not hand-edit it.
- English only.

## Non-goals

- Re-deciding the CLI change. It is shipped and reviewed.
- Documenting CLI internals (`AutonomousWrite.MissingFrom`, spec structs) in the skill.
