# round-01 — claude-1

Role: facilitation + naming-scheme coherence + skill/protocol integration.
Independent analysis (written before reading any other round-01).

## Position summary

Ship all three components, but resolve the ID/name crux with a **frozen-slug** model so the
user gets composite names everywhere WITHOUT breaking artifact continuity. Make the reinit
command an extension of `parley init` (one code path, two scopes), and make the autonomous
("yolo") write mode a first-class, per-agent field surfaced in both the config and the skill.

## A. Naming scheme

**Format:** `agent-model-effort`, lowercase, three sections joined by single `-`.

**Sanitization (per section):** lowercase → delete every char not in `[a-z0-9.]` so
spaces/parens/brackets/slashes vanish but VERSION DOTS are preserved; collapse repeated dots
to one and strip leading/trailing dots (path-safe: never `..`). Derive `model` from a human
model LABEL, not the raw id:
- `Opus 4.8` → `opus4.8`, `GPT-5.5` → `gpt5.5`, `GLM 5.2` → `glm5.2`,
  `Gemini 3.5 Flash (High)` → `gemini3.5flash`, `K3` → `k3`.
Effort from a fixed vocabulary `{low,medium,high,xhigh,max,ultracode,clidefault}`
(`cli-default` → `clidefault`).

**Resulting roster:**
- `claude-opus4.8-max`
- `codex-gpt5.5-xhigh`
- `hermes-glm5.2-high`
- `agy-gemini3.5flash-high`
- `kimi-k3-max`

All match `^[a-z0-9][a-z0-9.]*(-[a-z0-9.]+)*$` ⊂ `[a-zA-Z0-9_.-]`, readable (natural version
numbers), path-safe.

**Why the model keeps dots but never a `-`:** the dot lives INSIDE the model token, so the
name still splits on `-` into exactly 3 tokens (or 4 with a trailing numeric instance index)
— parsing stays unambiguous while `opus4.8` reads naturally. (`opus4-8` would overload the
section separator and read like a range.)

**Collisions:** same agent+model+effort twice → append `-2`, `-3`: `claude-opus4.8-max-2`.

**Parse:** split on `-`. `token[0]`=agent; if `token[-1]` is all digits it's the instance
index and effort=`token[-2]`, else effort=`token[-1]`; the remaining single token = model.
Validate agent against discovered bases and effort against the vocabulary; reject otherwise.

**agy special case:** no per-invocation effort flag; its tier lives in the model label
("(High)"). Surface that tier as the effort token → `agy-gemini35flash-high` (effort=`high`),
rather than `agy-gemini35flashhigh-clidefault`.

## B. The crux: stable ID vs composite name → **frozen slug**

Artifacts key off the ID today: `round-01/claude-1.md`, `### Signoff: claude-1`, review
paths, quorum in `00-prompt.md`. If the composite IS the live ID, swapping a model/effort
renames the agent mid-history and orphans in-flight artifacts. But the user explicitly wants
the composite to BE the name. Reconcile with a **frozen slug**:

- The composite is the agent's **name** shown everywhere a human looks (§2 roster, TUI tabs,
  digests, run headers) AND is what a NEW idea stamps into its `participants:` list.
- At idea creation, the composite is **frozen into that idea's `00-prompt.md`** and used
  verbatim for every artifact path/signoff for that idea's whole life. Changing an agent's
  model/effort later does NOT rewrite an OPEN idea — it only affects ideas created afterward.
- Keep a stable **base** (`claude`, `codex`, …) as the durable anchor for discovery, config
  keys (`[agents.claude]`), and history/aliasing. The composite = base + model + effort.

Net: composite names everywhere (user intent), zero churn for in-flight ideas (safety),
`agents.toml` still keyed by the stable base (no config thrash).

Migration: existing decks/ideas that use `claude-1` keep working — treat `claude-1` as a
legacy alias of the `claude` base; do not mass-rename historical artifacts. New ideas get
composite names. Provide a `parley roster migrate`/`--dry-run` preview.

## C. Autonomous ("yolo") write mode — first-class field

Add `autonomous_write_mode` (or `yolo_args`) to each agent spec, because there is no common
flag across vendors. Discovery + the generated config carry the per-agent mapping; the skill
states the requirement generically ("every participant is invoked in its auto-approve mode
so it can write its own artifact without a blocking prompt") plus the concrete mapping:
- claude `--permission-mode bypassPermissions`; codex `-s workspace-write`; hermes `--yolo`;
  agy `--dangerously-skip-permissions`; kimi headless `-p` (auto-approves in-workspace writes;
  `--yolo`/`--auto` are mutually exclusive with `-p`, so `-p` IS kimi's yolo-equivalent).

Safety: autonomous mode is scoped to the deck/workspace, never a blanket machine-wide bypass;
obvious-secret redaction still applies; the skill must say so explicitly.

## reinit command surface

`parley roster init [--scope session|machine] [--reinit] [--from FILE] [--yes]`:
- Discover every agent CLI (command detection); for each, enumerate models + effort levels
  where the CLI exposes them (`models`/`model list`/documented aliases), else fall back to the
  configured/`cli-default` value and let the user type an exact id.
- Interactively ask model + effort per agent; compute the composite name; write the roster.
- `--scope machine` writes `~/.parley/agents.toml`; `--scope session` writes the deck
  (`parley-deck/agents.toml` override + §2 roster + `meta/headless-agents.local.json`).
- **Fold into the existing `parley init` bootstrap** rather than a parallel implementation —
  same discovery/selection engine; `--reinit` just re-runs it against an existing deck/central
  file and never silently rewrites the quorum of an OPEN idea (frozen slugs protect them).

## Risks / open questions for the others
- codex-1: exact touch points in `discover.go`/roster resolution/`LoadAgentSpecs`; is the
  base-vs-composite split clean in the Go types, or does it ripple through artifact-path code?
- hermes-1: is `parley roster init` the right verb, or overload `parley init --reinit`? Session
  vs machine default when the user omits `--scope`?
- antigravity-1: sanitization corner cases (unicode, empty effort, duplicate collapses like two
  models that both collapse to `k3`); collision policy beyond `-2`.
- kimi-1: migration for existing decks/ideas — is legacy-alias enough, or do we need a rewrite
  tool? End-to-end: does the frozen-slug rule actually hold across TUI + driver + digests?
