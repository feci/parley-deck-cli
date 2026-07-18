---
idea: composite-agent-naming-and-roster-reinit
author: user
created: 2026-07-18
track: deliberation
participants: [claude-1, codex-1, hermes-1, antigravity-1, kimi-1]
roles:
  claude-1: facilitation + naming-scheme coherence + skill/protocol integration
  codex-1: Go internals (discover.go, roster resolution, ID stability across artifacts)
  hermes-1: reinit UX (discovery -> per-agent model+effort selection), session vs machine scope
  antigravity-1: sanitization edge cases, collision handling, backward compat with claude-1 IDs
  kimi-1: fresh-eyes end-to-end review + migration path for existing decks and ideas
status: final
---

## Problem / idea

Two related capabilities the user asked for:

**A. Composite, self-describing agent names.** Today a roster member is an opaque stable
ID like `claude-1`, `codex-1`. The model and effort/reasoning level are hidden in config
(`~/.parley/agents.toml`, deck `headless-agents.local.json`). The user wants the agent's
NAME to carry three facts so a run is self-documenting at a glance:

    <agent> - <model> - <effort>

e.g. `claude-opus4.8-max`, `codex-gpt5.5-xhigh`, `hermes-glm5.2-high`,
`agy-gemini3.5flash-high`, `kimi-k3-max`.

Hard constraint from the user: names must contain ONLY `[a-zA-Z0-9_.-]` (letters, digits,
`_`, `.`, `-`; the dot IS allowed so version numbers stay natural — no spaces, parens,
brackets, or slashes) and must stay READABLE. Path-safety: forbid `..` and any
leading/trailing dot so a name can never act as a path-traversal segment.

**B. A roster (re)initialization command.** A command that (re)builds the roster for a
scope by (1) discovering every agent CLI available on the machine, (2) for each, listing
its available models and effort/reasoning levels, (3) asking the user to pick a model +
effort per agent, and (4) writing the resulting roster with the composite names. Two
scopes:
  - **session/project** — this deck (`parley-deck/agents.toml` + `meta/headless-agents.local.json` + §2 roster table).
  - **machine** — the central `~/.parley/agents.toml` inherited by every project.

This overlaps the existing `parley init` deck-bootstrap gate and the mandatory
roster+model+effort confirmation. The design must RECONCILE them, not duplicate.

**C. All agents run in an autonomous ("yolo"/bypass-permission) mode.** The skill (and the
generated agents config) must make explicit that every headless participant is invoked in
its non-interactive auto-approve mode, so it can WRITE its own canonical artifact
(`round-NN/<id>.md`, signoffs, review files) without a blocking permission prompt. Each CLI
has its own equivalent — there is no single "yolo" flag across vendors:
  - claude: `--permission-mode bypassPermissions`
  - codex: `-s workspace-write` (workspace-scoped writes; or full bypass where needed)
  - hermes: `--yolo`
  - agy: `--dangerously-skip-permissions`
  - kimi (Kimi Code): headless `-p` print mode already auto-approves in-workspace writes
    (verified); NOTE `--yolo`/`--auto` are MUTUALLY EXCLUSIVE with `-p`, so the "yolo
    equivalent" for kimi headless IS plain `-p`.
The design must: (1) define a per-agent "autonomous write mode" field as a first-class part
of the agent spec / discovery, (2) have the skill state the requirement generically (every
participant runs in its auto-approve mode) plus the per-CLI mapping, and (3) keep it safe —
autonomous mode is scoped to the deck/workspace, never a blanket machine-wide bypass, and
obvious-secret redaction still applies.

**D. `fast` is the standard speed, as a SEPARATE axis from effort (user decision).** The
default startup speed for every agent is `fast`, but this is a distinct knob from the
reasoning `effort` — it must NOT downgrade the confirmed effort (claude=max, codex=xhigh,
hermes=high, agy=cli-default, kimi=max stay put). Semantics = Claude Code `/fast` (same
model, same effort, faster output), NOT the legacy "fast profile = sonnet/low" that
conflated speed with a model+thinking downgrade. Implications the design must resolve:
  - Redefine `speed` (`fast|deep|...`) so it does not silently override `model`/`effort`.
    Today `~/.parley/agents.toml [defaults] speed="deep"` and per-agent `profiles.fast =
    {model: sonnet, thinking: low}` DO downgrade — that conflation must be broken or the
    `fast` default reinterpreted as "fast output at the SAME model+effort".
  - `fast` is NOT part of the composite name (the name encodes only agent-model-effort);
    speed is a launch property, switchable to `deep` for heavy ideas.
  - Per-agent expression of "fast without dropping effort": claude `/fast` (Opus, faster
    output); for headless agents map to their fast-output/streaming mode where one exists,
    else a documented no-op (never a quality downgrade). The reinit command should default
    new rosters to `speed = fast`.

## Starting proposal (claude-1) — to be critiqued, not assumed

### Naming scheme
- Format: `agent-model-effort`, all lowercase, three sections joined by single `-`.
- Each section is sanitized to `[a-z0-9.]+`: lowercase, then delete every character not in
  `[a-z0-9.]` so spaces/parens/brackets/slashes vanish but VERSION DOTS are preserved.
  Collapse any run of dots to one and strip leading/trailing dots (no `..`, path-safe).
  - model derives from a human model label, not the raw id:
    `Opus 4.8`->`opus4.8`, `GPT-5.5`->`gpt5.5`, `GLM 5.2`->`glm5.2`,
    `Gemini 3.5 Flash (High)`->`gemini3.5flash`, `K3`->`k3`.
  - effort from a fixed vocabulary: `low|medium|high|xhigh|max|ultracode|clidefault`
    (`cli-default`->`clidefault`).
- The model section may contain dots but NEVER a `-`, so the whole name still splits
  cleanly on `-` into exactly 3 tokens (or 4 with a trailing numeric instance index).
- Collisions (same agent+model+effort twice) get a numeric suffix: `claude-opus4.8-max-2`.
- Parse rule for tooling: split on `-`; token[0]=agent; if token[-1] is all digits it is
  the instance index and effort=token[-2]; else effort=token[-1]; the remaining single
  token is the model (dots allowed).
- agy special case: it has no separate effort flag; its reasoning tier is baked into the
  model label ("(High)"), so we surface that tier as the effort token -> `agy-gemini3.5flash-high`.

### reinit command (surface — to be decided)
- Candidate: `parley roster init [--scope session|machine] [--reinit]` (or fold into
  `parley init --reinit`). Interactive: discover agents -> per agent prompt model + effort
  -> write roster with composite names. Non-interactive/CI: `--from <file>` or flags.

## Design tensions to resolve (the real work)

1. **Stable ID vs composite name.** Artifact paths and cross-references use the ID today
   (`round-01/claude-1.md`, `### Signoff: claude-1`). If the composite IS the ID, changing
   a model or effort CHANGES the ID and breaks in-flight artifact continuity / history.
   Options: (a) composite replaces the ID (accept churn, add migration); (b) keep a stable
   base ID (`claude-1`) and treat the composite as a DISPLAY NAME shown in TUI/digests/§2;
   (c) composite is the ID but with a stable *prefix* (`claude-1`) plus a descriptor.
   Pick one and justify. This is the crux.
2. **Where the effort token comes from** for agents whose effort is not a per-invocation
   flag (agy = model tier; hermes = config.yaml; kimi = config.toml). Single source of truth?
3. **reinit vs `parley init` bootstrap vs §9.0 readiness ping** — one command with a
   `--reinit` path, or a distinct verb? How does re-running change existing IDs used by
   OPEN ideas (must not silently rewrite quorum of an in-flight idea).
4. **Session vs machine scope mechanics** — session writes deck files; machine writes
   `~/.parley/agents.toml`. Precedence already exists (central < deck). Does reinit-session
   write an override layer or the full deck roster?
5. **Discovery** — how to enumerate each CLI's models + effort levels generically
   (some expose `models`/`model list`, some don't). Fallback when undiscoverable.
6. **Backward compatibility & migration** — existing decks/ideas use `claude-1` etc. A
   rename must not orphan history. Provide an alias/migration story.

## Constraints
- Protocol-track feature: any change to protocol wording goes in BOTH
  `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`
  (byte-identical modulo allowlisted zones; drift guard `TestEmbeddedDefaultMatchesLiveDeck`
  must stay green) and the skill fallback `references/COOPERATION.md` re-synced.
- Spans BOTH repos: `parley-deck-cli` (Go: discovery, roster resolution, the command, ID
  handling) and `parley-deck-skill` (SKILL.md bootstrap wording, central agents.toml seed).
- Names must be filesystem-safe and stable enough for artifact paths.
- English-only for all files under `parley-deck/`.
- Keep `go build ./...`, `go vet`, `gofmt -l` clean; add tests for sanitization + parsing
  + collision + the discovery/selection flow.

## Deliverable of this idea
A ratified FINAL.md design covering all four components: (A) the exact naming scheme
(sanitization + parse + collision + ID-vs-display decision), (B) the reinit command surface
+ behavior for both scopes + migration story, (C) the per-agent autonomous-write-mode field
+ the skill wording that mandates every participant run in its yolo/bypass equivalent (with
the per-CLI mapping), and (D) `fast` as the standard speed on a separate axis from effort
(no effort downgrade; not in the name; reinit defaults new rosters to speed=fast). Plus a
staged implementation plan across CLI + skill. Implementation and a full multi-channel
release follow after FINAL.
