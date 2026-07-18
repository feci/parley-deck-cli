# round-02 (cross-review) — claude-1

I read codex-1, hermes-1, and kimi-1's round-01 in full. **I concede my frozen-slug: the
composite must be a DISPLAY name, not the identity.** Three independent code-grounded
analyses converged against me and they are right.

## Why I moved (the code, not the vote)

My frozen-slug froze the composite into each idea's artifact paths. But the repo does not
work that way today, so freezing would actively FIGHT the code and *deepen* the very schism
we should be closing:
- hermes-1 & kimi-1: `internal/runner/runner.go:349` already writes `agent.ID+".md"` where
  `agent.ID` is the **spec/family** ID (`claude`, `codex`), not the §2 roster ID (`claude-1`).
- codex-1: `agents.Spec.ID` is overloaded as BOTH a vendor discriminator
  (`cleanParticipantEnv`/`isolatedAgentHome` switch on `claude`/`hermes`) AND the artifact
  identity. Baking model+effort into that one string breaks both jobs.
- kimi-1's clincher: model/effort are **mutable launch properties**; an identity that
  changes when you swap a model is not an identity. A display name recomputed at render from
  the single source of truth (layered `agents.toml`) can never drift, lie, or rewrite an
  OPEN idea's quorum.

So: **stable ID = identity (unchanged); composite `agent-model-effort` = derived display
name**, rendered in §2, TUI tabs, digests, run headers; never in an artifact path or the
`agent:` frontmatter key.

## Per-participant

- **codex-1** — Agree on splitting the overloaded `Spec.ID`. Concretely I support three
  distinct concepts: `identity` (stable roster key, artifact/signoff/event key), `family`
  (vendor discriminator for env/home/invocation), `display` (the composite, derived). Agree
  autonomous profiles need immutable per-run snapshots.
- **hermes-1** — Agree the pre-existing ID divorce is the real target and the composite is
  the display vehicle. Your point that the artifact filename is ALREADY not the §2 ID is
  what flipped me.
- **kimi-1** — Strong agree on the two-namespace trace. My one push: don't leave it at "just
  don't deepen it." The feature should ship a **resolver** (family-prefix match: `claude-1`
  → family `claude`; fail-closed on ambiguity) so a driver run over `participants:[claude-1]`
  stops matching nothing. Migration stays non-destructive (no mass rename of history).

## Folding in the two late updates

- **Dot in names:** because the composite is display-only, §2 keeps `claude-1` (satisfies the
  `internal/protocol/roster.go:17` `[a-z0-9][a-z0-9-]*` regex kimi-1 found), and the dotted
  composite (`claude-opus4.8-max`) lives in the config/display layer. Validate the display
  string against `^[a-z0-9]([a-z0-9.]*[a-z0-9])?(-[a-z0-9]([a-z0-9.]*[a-z0-9])?)*$` (dots
  allowed inside a section, never leading/trailing, never `..`). `consensus.go:90` already
  tolerates `[A-Za-z0-9._-]`, so if a signoff ever carries a display name it won't break —
  but signoffs should still use the stable ID.
- **Component D (fast axis):** agree it is orthogonal to effort and must not downgrade it.
  Redefine `speed` so `fast` means "fastest output at the SAME model+effort" (Claude Code
  `/fast` semantics), and BREAK the legacy conflation where `profiles.fast={model:sonnet,
  thinking:low}` silently downgrades. `speed` stays OUT of the composite name; reinit
  defaults new rosters to `speed=fast`, `--deep` per idea when needed.

## Component B / C positions

- **B:** I withdraw `parley init --reinit`; a distinct **`parley roster init [--scope
  session|machine] [--reinit] [--dry-run]`** verb (codex-1, kimi-1) is clearer — but it must
  SHARE the §0 bootstrap engine and the §9.0 probe, not fork them. `--dry-run` prints exactly
  which files change per scope; session refuses to run outside a deck; hard-refuse ID changes
  referenced by any open idea.
- **C:** agree — change the built-in defaults so they are actually autonomous (kimi-1: claude
  ships `acceptEdits`, codex `on-failure`), add a first-class `autonomous_write` field on the
  spec, and have the skill state the invariant + per-CLI mapping, scope-paired with the mode.

## Convergence check
We now agree on the crux (display-name), the dot handling, component D orthogonality, the
`roster init` verb, and autonomous defaults. Remaining for consensus: the exact Go type
split (codex-1 to specify), whether the resolver ships now or is scoped as a follow-up, and
the staged implementation order across CLI + skill.
