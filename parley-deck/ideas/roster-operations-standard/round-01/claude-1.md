---
agent: claude-1
idea: roster-operations-standard
round: 1
date: 2026-08-06
---

## Summary

**The roster is not under-documented; it does not exist as one object.** There are three separate
things wearing the name, and the user's symptoms follow from that directly:

1. **The adapter catalogue** — every agent the CLI knows how to launch (13 rows in
   `parley agents list`, including `qwen`, `goose`, `droid`, and a not-installed `gemini`).
2. **The roster proper** — the four IDs that participate in ideas and sign off (`roster show`).
3. **The per-deck config** — what `agents.toml` layers actually set.

Nothing defines the mapping between them. `opencode` became an adapter and stayed out of the
roster, which is why "I added opencode and sessions went inconsistent" is not a bug report about
opencode at all — it is the absence of a defined promotion path.

The stable-columns request is the visible half. The half that matters more is that **the table is
currently wrong**, and a stable format that keeps reporting a model the process does not use would
be a worse outcome than today's inconsistency, because it would look authoritative.

## Proposed approach

### 1. Name the three concepts and make one of them the answer

Define in the protocol and the skill:

- **catalogue** — launchable adapters. Answer to "what could I use".
- **roster** — the ordered, ID'd set that participates. Answer to "who is in this deck".
- **effective spec** — what a given roster member would actually launch right now.

"Show me the roster" MUST mean roster ∪ effective spec. `parley agents list` becomes the
catalogue view and says so in its header.

### 2. One table, one schema, printed by one code path

A single renderer with a frozen column list, used by `parley roster show`, by anything the skill
prints, and by `--json` (same field names). My proposed columns:

| column | source | why |
|---|---|---|
| `AGENT` | roster ID | the identity used in artifact paths and signoffs |
| `COMPANY` | derived from family | the user asked; also the quota/outage blast radius |
| `FAMILY` | adapter family | `claude`, `codex`, `hermes`, `kimi`, `opencode` |
| `MODEL` | **effective** launch value | must be what the process gets |
| `EFFORT` | effective | ditto |
| `SPEED` | config | `fast`/`balanced`/`deep`/`review` |
| `AUTO` | `AutonomousEffective()` | already fail-closed since 1.39.0 |
| `STATUS` | liveness/installed | `ok`, `not-installed`, `unverified` |
| `SOURCE` | winning config layer | so "why is it this value" is answerable without a second command |

`SOURCE` is the column I would fight for. Every roster confusion this project has had ended with
someone asking "where did that value come from", and the answer required reading four files.

### 3. Fix the declared-vs-effective split at its root, not in the renderer

`discover.go:219` embeds `--model claude-opus-4-8[1m]` inside `HeadlessArgs` while `:226` sets
`Model` separately. Config sets the field; the embedded literal wins at launch. The renderer must
not paper over this by parsing argv — the **spec must not carry a model in two places**.

Proposal: built-in `HeadlessArgs` carry a `{model}` / `{effort}` placeholder the way they already
carry `{prompt}` and `{root}`, and the resolver substitutes from the winning config layer. Then one
value exists, and `MODEL` is definitionally the launched model. For agents whose CLI takes no model
flag (`codex`, `kimi` today), the honest value is not `gpt-5.6-sol` — it is
`cli-default (~/.codex/config.toml)`, and the column should say so rather than assert a pin the
launch never passes.

This is the same fix shape as `AutonomousWrite.MissingFrom` in 1.39.0, and I would reuse the
pattern: a divergence is reported, never silently reconciled.

### 4. Three verbs, and no fourth

- `parley roster show [--scope session|deck|machine] [--json]`
- `parley roster set <agent> [--model M] [--effort E] [--speed S] --scope <deck|machine>`
- `parley roster sync --from machine --to deck [--dry-run]`

`sync` must print a diff and require `--yes`, because "sync local with global" silently overwriting
a deliberate per-project pin is exactly how a deck loses a deliberate choice. `roster init` stays as
the bootstrap/discovery path; `set` and `sync` are the steady-state operations the user is missing.

### 5. Promotion path for a new adapter

Adding an adapter must either add it to the roster or state that it did not. `parley roster show`
should list catalogue members that are installed but not in the roster under a short
`not in roster:` footer, so the opencode situation is visible the moment it happens rather than
three sessions later.

## Concerns / open questions

1. **Is `SOURCE` one column or does it need per-field granularity?** `agents list` already prints
   `sources: sandbox=… approval=… model=…` per agent. Collapsing that to one column loses
   information; keeping it per-field makes the table unreadable. My lean: one column naming the
   layer that won for `MODEL`, with per-field detail in `--json` and in a `--verbose` mode.

2. **`COMPANY` derivation.** `claude`→Anthropic, `codex`→OpenAI, `kimi`→Moonshot, `agy`→Google are
   easy. `hermes` runs `glm-5p2` (Zhipu) through a LiteLLM gateway, and `opencode` is configured as
   `litellm/xai/grok-4.5` (xAI via gateway). So company is a property of the **model**, not the
   CLI, and for gateway-routed agents it is only knowable by parsing the model id. I do not think
   we should invent a mapping table we then have to maintain; I would derive where the id encodes
   it and print `unknown` otherwise. Reviewers should push back if `unknown` is unacceptable.

3. **Does `--scope session` even exist as a real store?** `roster init` accepts
   `--scope session|machine`, but I have not verified what a session-scoped roster is persisted
   as, or whether it survives a restart. This needs checking before we build `sync` on top of it.
   I own this uncertainty and am not asserting either way.

4. **Should `agents list` be kept at all?** Two commands is how we got here. But it carries launch
   diagnostics (`headless:`, `acp:`, `note:`) that do not fit a fixed table. My lean is to keep it,
   rename its header to make "catalogue" explicit, and have `roster show` be the documented answer.

## Risks

- **A frozen schema that freezes a lie.** If we standardize the columns before fixing the
  model-substitution bug, we ship a stable table that confidently reports the wrong model. Order
  matters: fix the spec, then freeze the format.
- **Scope creep into discovery.** The `{model}` placeholder change touches every built-in adapter
  spec. That is a bigger diff than "add columns" and could pull the idea into a rewrite of
  `discover.go`. A reviewer should say if this belongs in a separate idea.
- **`sync` destroying deliberate per-project pins.** The whole value of the layering is that a deck
  can differ from the machine default. A sync verb that flattens that is a regression wearing a
  feature's clothes. Dry-run-by-default and an explicit diff are not optional.
- **My own record on this file.** I found the model divergence by reading `parley agents list`
  output and one grep of `discover.go`. I have **not** run a live probe proving that a launched
  `claude` process actually receives `claude-opus-4-8[1m]`. Under §15 that makes my root-cause
  claim `PRIMARY` for the source lines and **unproven for the runtime behaviour**. Another
  participant should verify the runtime half before we design a fix around it.
