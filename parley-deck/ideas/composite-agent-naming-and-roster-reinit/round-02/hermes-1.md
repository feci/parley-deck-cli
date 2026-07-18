---
agent: hermes-1
idea: composite-agent-naming-and-roster-reinit
round: 2
date: 2026-07-18
role: reinit UX + the display-name rendering path (§2 table, TUI, digests); session vs machine scope defaults and the open-idea guard
---

## Where I land (one paragraph)

Composite = DISPLAY name; stable roster ID (`claude-1`) = identity. That is
option (b), and it is the position three of us (codex-1, hermes-1, kimi-1)
reached independently. claude-1's "frozen slug" is option (b) wearing a
different label — once the composite is frozen per-idea and the base ID is the
config key, the composite is a display name and the base ID is the identity, so
the three proposals converge on the same shape; the only real disagreement is
**whether to migrate the runner from the family ID (`claude`) to the roster ID
(`claude-1`) as the artifact identity**. I now think the right answer is: **fix
the schism, but only for NEW ideas, and only by making the roster ID the
participant key the runner matches against — never by rewriting history.**
Component B is a new `parley roster init` verb (not `parley init --reinit`).
Component C lands the per-CLI autonomous mapping AND corrects the built-in
defaults so they are actually autonomous. Component D (`fast` on a separate
axis from effort) is mostly already true in the code today — the brief's worry
about `profiles.fast = {model: sonnet, thinking: low}` downgrading does not
match this repo — but the central default still says `speed = "deep"`, which
reinit must flip to `fast`, and the `Profile` field must be redefined so it can
never silently carry a model/effort downgrade.

## Per-participant reckoning

### claude-1 — agree on the destination, disagree on the label and the migration cost

Agree: the sanitization rule (delete every char not in the allowed set,
preserve version dots, collapse `..`, strip leading/trailing dots), the
fixed effort vocabulary, the `agy` tier-as-effort special case, the numeric
collision suffix, and the parse rule (split on `-`, last token = effort or
instance index, middle = model). Agree that `agents.toml` stays keyed by the
stable base and the composite must never be a filename or frontmatter key.
Agree that autonomous mode is a first-class per-agent field with the per-CLI
mapping in config, not prose.

Disagree on two things.

(1) **The label "frozen slug" is misleading and I would drop it.** claude-1
frames the composite as "frozen into `00-prompt.md` and used verbatim for every
artifact path/signoff for that idea's whole life." But artifact paths and
signoffs use the **base ID** in every other participant's proposal and in the
actual code (`runner.go:349` writes `agent.ID+".md"`; `consensus.go:90` parses
`### Signoff: <id>`; `consensus.go:145` membership-checks against
`participants:`). If we freeze the composite into `00-prompt.md`'s
`participants:` list, then the runner's `selectedAgents`
(`runner.go:327-342`), which matches `participants:` against `agent.ID` by
exact string equality, **matches nothing** for this very idea — exactly the
split-brain kimi-1 traced in §0. So the composite must NOT be what
`participants:` lists. Counter-proposal: `participants:` lists the **roster
ID** (`claude-1`); the composite is a derived display string stored alongside
in config and rendered in the §2 table / TUI / digest, recomputed from
`model`+`reasoning` at render time (kimi-1's "derive-on-load, store nothing"
position, which I adopt). The "frozen per idea" property claude-1 wants is
satisfied differently and better: the **per-round artifact frontmatter** records
the `model:`/`effort:`/`speed:` that actually ran (kimi-1 §1.5's audit
proposal), so two years later you can see exactly what wrote round-01 without
the identity layer ever churning. That is strictly more audit than a frozen
display name, and it costs zero rename risk.

(2) **Migration: I would NOT take the one-time rename of existing artifacts.**
claude-1 says "existing decks/ideas that use `claude-1` keep working — treat
`claude-1` as a legacy alias of the `claude` base; do not mass-rename." That is
the right instinct but it contradicts the frozen-slug framing. The honest
position is: **zero renames, ever, forward-only** (kimi-1 §5). New ideas get
the new wiring (roster ID as the participant key the runner matches); old ideas
keep their old-ID artifacts untouched. The precedent is already in the deck:
`meta/roster-update_2026-06-19.md` did exactly this (`agy` → `antigravity-1`,
forward-only, old ideas kept old artifacts). I agree with claude-1's
`--dry-run` preview but not with `parley roster migrate` as a rename tool —
make it a *report* tool (`parley roster show --diff`) that shows drift, not a
rewriter.

### codex-1 — strongest Go trace; adopt the Spec split, tighten two things

Agree almost entirely. codex-1's `Spec` split (`ID` stable protocol ID,
`AdapterID` launcher family, `DisplayName` composite, `Model`/`Reasoning`
authoritative) is the right type shape and I adopt it. Agree that the composite
is "presentation and an immutable launch-profile snapshot, not a lookup key."
Agree on the `runmanifest.Manifest` participant-profile snapshot for resume
(hermes-1 round-01 did not say this; codex-1 is right to add it). Agree on the
collision rule (preserve surviving allocations, deterministic stable-ID
ordering, never auto-compact). Agree on the open-idea guard being a hard refuse
on ID rename/removal touching a live quorum.

Two counter-proposals.

(1) **On the schism: codex-1 says "keep `Spec.ID` so the many artifact and
event call sites remain keyed by stable identity" — but it leaves `Spec.ID`
pointing at the *family* (`claude`), not the *roster seat* (`claude-1`).** That
perpetuates the two-namespace split kimi-1 found: the runner writes
`round-01/claude.md`, the §2 table and `00-prompt.md` say `claude-1`, and
`consensus.go:145` rejects the runner's signoff as "unknown participant" for
any idea that lists `claude-1`. codex-1's own risk list flags exactly this
("conflating ID, adapter, and display will reintroduce path instability"), but
its proposed split does not close the gap — it just renames the family ID to
`AdapterID` and keeps `ID = "claude"`. Counter-proposal: **`Spec.ID` becomes
the roster seat ID (`claude-1`) for rostered agents; `AdapterID` is the family
(`claude`); `selectedAgents` matches `participants:` against `Spec.ID`.** That
makes the runner and the protocol agree on one key, and it is the actual bug
this idea exists to fix (hermes-1 round-01 §"the runner→§2-roster gap is the
real bug"). The cost is a one-time wiring change for NEW ideas only; existing
decks keep family-ID artifacts because the runner never rewrites history.
codex-1's concern about `runner.cleanParticipantEnv`/`isolatedAgentHome`
switching on vendor strings is handled cleanly: those switch on `AdapterID`
(the family), which is exactly what codex-1 proposed, so no behavior change
there.

(2) **codex-1's charset is narrower than the brief now requires.** codex-1
deletes dots from model tokens (`opus48`, not `opus4.8`), citing "deliberately
narrower than the user's hard character set." The brief has since been updated
to say the dot IS allowed so version numbers stay natural: the target names are
`claude-opus4.8-max`, `codex-gpt5.5-xhigh`, `hermes-glm5.2-high`,
`agy-gemini3.5flash-high`, `kimi-k3-max`. codex-1's `codex-gpt55-xhigh` and
`hermes-glm52-high` do not match. Counter-proposal: **keep the dot inside the
model token**, exactly as claude-1 round-01 proposed and as the updated brief
mandates. The model section is `[a-z0-9.]+` (dots allowed, no `-`), so the
whole name still splits cleanly on `-` into exactly 3 tokens (or 4 with a
trailing instance index). This is the single point where claude-1's round-01
was more correct than codex-1's, and the late brief update settles it. The
canonical display names for the five rostered agents are therefore:

| Roster ID       | AdapterID | Display name              |
| --------------- | --------- | ------------------------- |
| `claude-1`      | `claude`  | `claude-opus4.8-max`      |
| `codex-1`       | `codex`   | `codex-gpt5.5-xhigh`      |
| `hermes-1`      | `hermes`  | `hermes-glm5.2-high`      |
| `antigravity-1` | `agy`     | `agy-gemini3.5flash-high` |
| `kimi-1`        | `kimi`    | `kimi-k3-max`             |

All match `^[a-z0-9][a-z0-9.-]*$` ⊂ the brief's `[a-zA-Z0-9_.-]`, contain no
`..`, no leading/trailing dot, and split on `-` into exactly 3 tokens. Note the
`antigravity-1` seat's display name uses the **family** `agy` (not the roster
ID prefix `antigravity`), because model/effort are properties of the CLI
runtime, not the seat — this is kimi-1 §2's point and it is correct: a naive
prefix heuristic would produce `antigravity-gemini3.5flash-high`, which lies
about which CLI runs.

### hermes-1 (my own round-01) — what I'm changing

I keep: the distinct `parley roster` verb, the session-default-when-deck-exists
rule, the open-idea guard mechanics, the derive-on-load display, the per-CLI
autonomous mapping with `workspace-write` as the default and
`bypassPermissions` as a confirmed escalation, the TUI width fix.

I am changing two positions based on reading the others and the late brief
updates.

(1) **I over-rotated on "wire the runner to the base ID = §2 roster ID as a
mandatory migration."** Round-01 said the runner→§2 gap "is the actual bug this
idea should fix" and leaned toward taking the migration. After reading codex-1
and kimi-1, the migration is real cost and the forward-only path gets 90% of
the benefit at 10% of the risk. Revised position: **the runner matches
`participants:` against the roster seat ID (`claude-1`) for ideas created after
ratification; existing artifacts are never renamed; the `AdapterID` is what
vendor-specific code switches on.** This is a wiring change in
`selectedAgents`/`runAgent`, not a filesystem migration. Existing decks keep
working because their old ideas keep their old `participants:` lists (which the
runner will still match if we preserve family-ID matching as a legacy fallback
ONLY for ideas whose `participants:` contain a bare family ID — fail-closed on
ambiguity). kimi-1's §5.4 "contain, don't widen" is the right posture for the
schism, and I now agree the full driver-side resolver is a separate idea — but
this idea should at least make the *new* path correct.

(2) **Component D was underspecified in my round-01 (I did not address it).**
Folding in the late update below.

### kimi-1 — agree on the crux and the schism posture; adopt the derive-on-load and the built-in fix; one correction on charset

Agree on nearly everything and I adopt several of kimi-1's contributions
wholesale: the fresh-eyes two-namespace trace (§0), derive-on-load display
(§1.2), the `model_label` field distinct from `model` (§2
counter-position), the honest "codex is configured `cli-default` today so its
display is currently `codex-clidefault-clidefault`" observation, the
right-to-left fail-closed parse rule with effort validated by vocabulary
membership (so all-digit model tokens parse), the surgical (not bulk) update of
`headless-agents.local.json`, the forward-only migration with zero renames
(§5), the built-in autonomous defaults being broken today and needing the fix
(§4), and the exact skill wording block (§4). The `model_label` field is the
cleanest answer to "where does the human label come from" and I adopt it: reinit
writes `model_label = "Opus 4.8"` and `model = "claude-opus-4-8[1m]"` (or
whatever the CLI flag wants); the display is derived from `model_label`, never
from the raw `model` id, so `claude-opus-4-81m` never appears.

One disagreement, and it is the same one as with codex-1: **kimi-1 chose the
lowercase-hyphen subset `[a-z0-9-]` and deletes dots (`opus48`), but the brief
now mandates dots (`opus4.8`).** kimi-1's reasoning was (i) the §2 roster
parser accepts nothing else (`protocol/roster.go:17` = `[a-z0-9][a-z0-9-]*`),
(ii) case-insensitive filesystems make uppercase a collision vector, (iii) `_`
makes tokenization ambiguous. Points (ii) and (iii) stand and I keep them
(canonical display is lowercase, no `_`). Point (i) is the crux of the
disagreement and it needs care:

- The §2 roster parser regex (`roster.go:17`) governs the **roster ID**
  (`claude-1`), NOT the display name. The display name lives in a new column
  and is never parsed by `ReadRosterIDs`. So the roster ID stays
  `[a-z0-9][a-z0-9-]*` (no dots, no underscores) and the display name can use
  the wider `[a-z0-9.-]` charset the brief allows. **Two charsets, two
  namespaces, one table** — that is exactly the schism made explicit and
  harmless.
- The signoff parser (`consensus.go:90` = `[A-Za-z0-9._-]+`) already tolerates
  dots, so a display name with a dot would parse as a signoff header if anyone
  ever wrote one — but no one will, because signoffs use the roster ID. Still,
  it means the display charset is a subset of what the signoff regex accepts,
  so there is no parser that rejects it.

Counter-proposal: **roster ID charset = `[a-z0-9][a-z0-9-]*` (unchanged,
roster.go:17 stays); display name charset = `[a-z0-9][a-z0-9.-]*` with no `..`
and no leading/trailing dot (new, wider, matches the updated brief).** The
display name is generated, stored in the §2 table's new "Display name" column
and in config, and rendered in TUI/digest; it is never the argument to
`ReadRosterIDs`, `selectedAgents`, `AppendSignoff`, or `validation.go:42-72`.

On kimi-1's two-namespace finding (spec ID vs §2 roster ID): **this feature
should FIX the schism for new ideas, not just guarantee not to deepen it.**
kimi-1 §5.4 says "fixing the schism itself is a separate idea; this one must
not depend on it." I respect that boundary but I think we can do slightly
better at near-zero cost: the fix is to make `Spec.ID` the roster seat ID
(`claude-1`) and `AdapterID` the family (`claude`), and have `selectedAgents`
match `participants:` against `Spec.ID`. That is a ~10-line wiring change in
`runner.go:327-342` plus the `Spec` split codex-1 proposed; it does not require
a driver-side resolver or a prefix heuristic. Legacy ideas whose
`participants:` list a bare family ID (`claude`) keep working via a legacy
fallback in `selectedAgents` (if no roster-seat match, try family match;
fail-closed if both or neither). The concrete Go touch points:

- `internal/agents/discover.go`: `Spec` gains `AdapterID string` (family);
  `ID` becomes the roster seat for rostered agents (the central/deck config
  writes `[agents.claude-1]` with `adapter = "claude"`; the built-in
  `defaultBuiltinSpecs` keeps `ID = "claude"` as the unrostered family default,
  so discovery of a bare `claude` CLI still works before any roster exists).
- `internal/runner/runner.go:327-342` (`selectedAgents`): match
  `participants:` against `Spec.ID` first, then `AdapterID` as legacy
  fallback; `runAgent` (`:349`) writes `agent.ID+".md"` using the matched ID,
  so new ideas get `round-01/claude-1.md` and old ideas keep
  `round-01/claude.md`.
- `internal/runner/runner.go` `cleanParticipantEnv`/`isolatedAgentHome`:
  switch on `AdapterID` (family), not `ID` — no behavior change.
- `internal/consensus/consensus.go:145`: unchanged; it already
  membership-checks against `participants:`, which now lists roster IDs, and
  the runner writes signoffs under the roster ID, so they agree.
- `internal/protocol/roster.go:17`: unchanged (roster ID charset stays
  narrow).
- `internal/config/runtime.go`: `agentOverride` gains `adapter` and
  `display_name`/`model_label`; `applyFile` keys `[agents.<id>]` blocks by the
  roster seat ID; `[agents.claude]` (bare family) is treated as a legacy
  adapter default, not a roster seat.

This is "fix the schism for new ideas, contain it for old ideas" — stronger
than kimi-1's "contain, don't widen," and it does not require the separate
driver-side resolver idea kimi-1 (correctly) says is out of scope.

## Folding in the two late brief updates

### (1) Dot-in-names charset

Settled above. Two charsets:

- **Roster ID** (the identity, the key, the filename, the signoff, the
  `participants:` entry): `[a-z0-9][a-z0-9-]*` — unchanged, matches
  `protocol/roster.go:17`, no dots, no underscores, lowercase only.
- **Display name** (the composite, the descriptor, never a key): generated
  from `agent` (the family/AdapterID) + `model` (sanitized from `model_label`)
  + `effort` (vocabulary token), each section `[a-z0-9.]+` for agent/model and
  a vocab token for effort, joined by single `-`. Whole name matches
  `^[a-z0-9][a-z0-9.-]*$`, contains no `..`, no leading/trailing dot. Dots live
  only inside the model token so splitting on `-` is still unambiguous.

The five display names are in the table above. The parse rule (for tooling that
needs to split a display name back out — e.g. the TUI rendering a chip): split
on `-`; `token[0]` = family; if `token[-1]` is all digits it is an instance
index and effort = `token[-2]`, else effort = `token[-1]`; the single remaining
middle token is the model slug (lossy — the real model id is in config, never
reversed from the slug). Effort is validated by vocabulary membership
(kimi-1's tightening, so an all-digit model token like a hypothetical `530`
does not misparse). Anything that does not round-trip is a parse error;
tooling must never guess.

Collision rule: within one roster scope, composites are allocated in catalog
order (built-in spec order, then config-append order); first occurrence keeps
the bare name, later duplicates of the same family+model+effort get `-2`,
`-3` (the instance suffix, same mechanism as the parse rule). Persisted
allocations are preserved across reinit (codex-1's rule) so a harmless reinit
never renames surviving displays. A composite collision is cosmetic (the two
agents are still distinguished by roster ID in all paths/signoffs), but reinit
warns and suggests distinct roster IDs.

### (2) Component D — `fast` as a separate axis from effort

This is the part the round-01 files did not address (it landed after launch),
and there is one important code fact to correct first: **the brief's worry that
"`profiles.fast = {model: sonnet, thinking: low}` downgrades" does not match
this repo today.** I searched every `.toml` and `.go` file: there are no
`[profiles.fast]` / `[profiles.deep]` subtables anywhere. `Speed` and `Profile`
are flat string fields on `Spec` (`discover.go:32-33`) and `agentOverride`
(`runtime.go:110-111`), and `DefaultSpeed = "balanced"` (`discover.go:68`).
Searching for consumers of `.Speed`/`.Profile`, the only write path is
`runtime.go:341-347` where the global `[defaults].speed` is copied onto every
spec's `Speed` field — and `Speed`/`Profile` are **never read into
`HeadlessArgs` or any launch assembly**. They are inert metadata today
(`runcontrol.go:163-167` echoes them into run events for observability, but
they do not change the argv). So the "fast profile = sonnet/low downgrade" is a
legacy/hypothetical conflation from a different codebase, not a live bug here.

That said, component D still requires real work, because (a) the central
default template hardcodes `speed = "deep"` (`runtime.go:286`), (b) `DefaultSpeed`
is `"balanced"` not `"fast"`, and (c) the `Profile` field is a loaded name that
COULD become a downgrade vector if someone later wires it to a profile table.
My position:

1. **Redefine the `speed` vocabulary to `{fast, deep}` and make `fast` the
   default.** `DefaultSpeed` changes from `"balanced"` to `"fast"`; the central
   default template (`runtime.go:286`) changes from `speed = "deep"` to
   `speed = "fast"`. `balanced` is accepted as a legacy alias mapped to `fast`
   (so existing `agents.toml` files with `speed = "balanced"` do not break).
   The deck's current `speed = "balanced"` / `speed = "deep"` values
   (`parley-deck/agents.toml:10,18,29,42`) are reinterpreted: `balanced` →
   `fast`, `deep` stays `deep`.

2. **`speed` is a launch-property axis, NEVER part of the composite name and
   NEVER a model/effort override.** The design must state explicitly: `speed`
   controls output pacing only (Claude Code `/fast` semantics — same model,
   same effort, faster output); it does NOT change `Spec.Model`,
   `Spec.Reasoning`, or `HeadlessArgs`' model/effort flags. The per-agent
   expression of "fast without dropping effort":
   - claude: append `/fast` (or the headless equivalent flag) to the
     invocation — same `--model`, same `--effort`, faster output.
   - codex: map to its fast-output/streaming mode where one exists, else a
     documented no-op (never a quality downgrade).
   - hermes: `--yolo` already implies non-interactive; a `fast` speed maps to
     the streaming/fast-output flag if the CLI exposes one, else no-op.
   - agy: `--print-timeout` already bounds it; `fast` is a no-op unless agy
     exposes a streaming flag.
   - kimi: `-p` print mode is already the fast path; `fast` is the default and
     `deep` would be the opt-out.
   When no fast-output flag exists, `fast` is an honest no-op (the model and
   effort are unchanged; the agent simply is not slower). The skill must say:
   "`fast` means 'do not add extra thinking overhead,' never 'use a weaker
   model.'"

3. **Kill the `Profile` field's ability to downgrade.** `Profile` is currently
   a free-form string (`CLIDefault` by default) with no consumer. Either (a)
   remove it (it is dead — no code reads it into the launch path), or (b)
   redefine it as a strictly-named enum that is explicitly NOT a model/effort
   selector (e.g. `{cli-default, fast, deep}` mirroring speed, or remove it
   entirely and let `Speed` carry the concept). I prefer (a) remove `Profile`
   and keep `Speed` as the single pacing axis, because two fields with
   overlapping names (`Profile` and `Speed`) is exactly the kind of ambiguity
   that lets a future "fast profile" reintroduce the downgrade. If the group
   wants to keep `Profile` for backward compat, it must be documented as
   "obsolete alias of `speed`; ignored at launch" and never wired to a profile
   table.

4. **reinit defaults new rosters to `speed = fast`.** Both scopes: the central
   template writes `speed = "fast"`, and a session reinit that does not
   explicitly pick writes `speed = "fast"` into the deck override. The user can
   switch a roster (or a single idea) to `deep` for heavy ideas. `fast` is NOT
   in the composite name (the name encodes only agent-model-effort).

5. **Per-idea override.** Since `speed` is a launch property, an idea's
   `00-prompt.md` MAY carry a `speed:` frontmatter field (defaulting to the
   roster's `fast`) so a heavy idea can opt into `deep` without a roster
   reinit. The runner reads it and passes the corresponding flag. This is the
   "switchable to deep for heavy ideas" the brief asks for.

## Settling component B — command surface and both scopes

**Command surface: a new `parley roster init` verb, NOT `parley init --reinit`.**

This is the point hermes-1 round-01, codex-1, and kimi-1 all agreed on and
claude-1 round-01 leaned toward ("fold into the existing `parley init`
bootstrap"). I keep the distinct-verb position and here is the final reasoning:

- `parley init` (`app.go:362`) creates the deck scaffold + central seed +
  transport. That is a *creation* job.
- `parley roster init` rebuilds the roster on a live deck. That is a
  *maintenance* job that can destabilize open ideas and must run the open-idea
  guard.
- Overloading `init --reinit` conflates the two and makes the help text lie
  about blast radius. A `roster` verb group also leaves room for `roster show`
  / `roster diff` / future `roster add` / `roster drop` / `roster retire`
  without re-shaping the CLI.
- claude-1's concern was "one code path, two scopes" — that is satisfied by
  having `parley roster init` be the single selection engine that BOTH
  `--scope session` and `--scope machine` call, and by having `parley init`
  *chain* into it (the skill's bootstrap says "run `parley init` then `parley
  roster init --scope session`"). One selection flow, one writer, two verbs
  with different jobs.

Final surface:

```
parley roster init [--scope session|machine] [--from <file>] [--dry-run] [--yes] [--json]
parley roster show [--scope session|machine] [--json]
parley roster diff [--scope session|machine] [--json]   # what would change vs current
```

`--reinit` is NOT a separate flag: the command IS the re-init. Without
`--yes`/`--from` it runs interactively with existing values as the defaults
shown in the picker. It refuses to silently overwrite a scope that already has
roster entries — interactive confirmation or `--yes` required.

**Both scopes — exactly what each writes:**

- **`--scope machine`** writes ONLY `~/.parley/agents.toml`. It is a
  managed-block edit, not a full rewrite: preserve `[defaults]` (except
  flipping `speed` to `fast`), `[rosters.*]`, unknown keys, and comments; on
  first conversion create a backup before establishing managed-block markers
  (codex-1 concern 3: `go-toml/v2` is not comment-preserving, so use a
  managed-block/lossless editor, not a full marshal). It writes no deck file
  and never edits `COOPERATION.md`. It is the seed for future projects. Deck
  values are never copied up.
- **`--scope session`** (default when a deck exists in cwd) writes:
  1. `parley-deck/agents.toml` — the tracked deck override layer. It writes
     only the `[agents.<roster-id>]` blocks that differ from the central
     default plus the `[defaults]` knobs the user touched (override layer, not
     a full copy — keeps the deck file small and the diff readable). Each
     block: `adapter`, `model`, `model_label`, `reasoning`, `speed`,
     `display_name` (derived, optionally stored), `headless_args`,
     `sandbox_mode`, `approval_policy`, `autonomous_write`,
     `autonomous_write_args`.
  2. `parley-deck/COOPERATION.md` §2 roster table — a targeted, allowlisted
     edit that adds/refreshes a "Display name" column and a "Speed" column.
     The §2 table currently has 3 columns (`Agent ID | Workspace dir | Role`)
     with the `cli`/`model` info buried in the Role prose
     (`COOPERATION.md:105-110`); the new table is
     `Agent ID | Display name | Workspace dir | Role | Speed`. The edit shares
     one allowlist definition with the drift guard
     (`TestEmbeddedDefaultMatchesLiveDeck`) so the reinit writer and the drift
     guard agree on what the "roster zone" is (hermes-1 round-01 risk). The
     protocol default copy (`internal/protocol/defaults/COOPERATION.md`) and
     the skill fallback (`references/COOPERATION.md`) are re-synced in the
     same change so the drift guard stays green.
  3. `parley-deck/meta/headless-agents.local.json` — a SURGICAL update
     (kimi-1 §3), not a bulk rewrite: only the model/thinking/writeMode/speed
     fields for rostered IDs, and only if the file already exists. This file is
     gitignored machine-local state; bulk-rewriting it from a shared command
     would stomp per-machine CLI paths. If it does not exist, reinit creates it
     with rostered IDs only.

**Default scope when the flag is omitted:** `session` if a deck exists in the
cwd (the more common, lower-blast-radius operation, matches where the user
already is); `machine` if no deck exists (you cannot write a session roster
for a deck that isn't there — fall back to seeding the central default, with a
hint to run `parley init` then `parley roster init --scope session`).

**Open-idea guard (my lens — the UX I own):** before any write, scan
`parley-deck/ideas/*/00-prompt.md` for `status:` not in
`{final, abandoned, complete, archived}` and collect their `participants:`.
- Any roster-ID **rename or removal** touching an open idea's quorum: **hard
  refuse** by default, print the blocking ideas, exit 3. `--yes` proceeds but
  WRITES the break into each affected idea's `00-prompt.md` as an
  `excluded:`/`reincluded:` line (reusing the format the preflight-readiness
  idea ratified) — CI does not get to silently break quorum either; if
  `--yes` is set it writes the lines and exits 0, if not it exits 3 with a
  JSON report.
- Adding a new agent: not a break (it is not in any open idea's quorum yet);
  the new agent just becomes available for future ideas.
- Changing an existing agent's **model, effort, or speed** (composite or speed
  changes, roster ID stable): **never** a quorum break — that is the whole
  point of keeping the roster ID stable. Model/effort/speed re-picks are
  always safe and never trip the guard. This is also why component D's
  `speed` is safe to default to `fast` on reinit: it changes a launch
  property, not the identity, so no open idea is affected.

**Reconciliation with `parley init` and §9.0:**

- `parley init` stays bootstrap-only (scaffold + central seed + transport) and
  prints `parley roster init --scope session` as the required next step. The
  skill's §0 "MUST confirm roster + model + effort" gate is satisfied BY
  RUNNING `parley roster init`. One selection flow, one writer; `parley init`
  never grows a second divergent picker. (Optionally `parley init
  --with-roster` chains it for convenience, as kimi-1 suggested.)
- §9.0 preflight (`preflight.go`) stays read-only over roster selection: it
  pings liveness through the real configured invocation and gates on
  availability; it never re-selects models or effort. The only addition: a
  non-gating WARN when a rostered model/effort is no longer in the family
  catalog ("roster drift — run `parley roster init`"). Write path = roster
  init; verify path = preflight; both share one probe implementation
  (`pingProbe`/`hostedPONG`). reinit optionally runs a liveness probe against
  the newly-selected roster as a post-write verification, not a gate.

## Settling component C — autonomous write mode, for real

**Representation:** a first-class field on `agents.Spec` plus the per-CLI
mapping in built-in specs and the skill. I adopt codex-1's structured shape
over kimi-1's bare boolean, because the boolean alone cannot carry the
workspace-confinement verification codex-1 rightly demands:

```go
type AutonomousWriteMode struct {
    Mode             string   // bypass-permissions | workspace-write | accept-edits | print-auto
    Args             []string // launch fragments; empty when implicit
    Scope            string   // must be "workspace"
    ScopeEnforcement string   // cli-sandbox | outer-sandbox | unverified
}
```

`AutonomousWrite` on `Spec`; `agentOverride` gains the TOML fields;
`withBuiltinSources` field list (`discover.go:224-259`) and
`centralDefaultTemplate` (`runtime.go:277-313`) carry them; an `AUTO` column
in `PrintRuntimeMatrix` (`discover.go:314-384`) shows the effective mode.
roster init sets the field only after checking the configured args match the
family's known autonomous mapping (static check per family); a live
sentinel-write probe is opt-in (`parley agents verify --full`), not default
(probes are slow and cost tokens — kimi-1 §4). Preflight WARNS (not gates)
when a rostered agent lacks a verified autonomous profile.

**Fix the built-ins — today they are NOT all autonomous (kimi-1's fresh-eyes
finding, confirmed in the code):**

- claude built-in ships `--permission-mode acceptEdits`
  (`discover.go:134`) — still prompts outside edits. **Change to
  `bypassPermissions`** with `--add-dir {root}` (scope = the add-dir + cwd).
- codex built-in ships `approval_policy="on-failure"` (`discover.go:113`) —
  can still ask. **Change to `approval_policy="never"`** with
  `-s workspace-write` (workspace-scoped; never escalate to full bypass for
  protocol artifacts — codex-1 §C).
- hermes built-in lacks `--yolo` (`discover.go:203` has `--accept-hooks` but
  no `--yolo`; `--yolo` exists only in this deck's local JSON). **Add
  `--yolo`** to the built-in `HeadlessArgs`, keep `--accept-hooks`, cwd =
  root.
- agy built-in already has `--dangerously-skip-permissions`
  (`discover.go:156`) — correct; keep, with `--add-dir {root}`.
- kimi has NO built-in spec (only the ACP catalog, `acp_specs.go:37`). **Add
  a built-in with `-p` print mode** (which auto-approves in-workspace writes).
  NEVER add `--yolo`/`--auto` — mutually exclusive with `-p` (brief + kimi-1).

The per-CLI mapping table (encoded in built-ins AND in the skill):

| CLI    | Autonomous invocation (workspace-scoped)                                             |
| ------ | ------------------------------------------------------------------------------------ |
| claude | `--permission-mode bypassPermissions --add-dir {root}` (scope = add-dir + cwd)      |
| codex  | `exec --sandbox workspace-write -c approval_policy="never"` (`--cd {root}`)          |
| hermes | `--yolo --accept-hooks` (cwd = root)                                                 |
| agy    | `--dangerously-skip-permissions --add-dir {root}`                                    |
| kimi   | plain `-p` print mode (auto-approves in-workspace writes); NEVER `--yolo`/`--auto`   |

**Safety invariants (encode in code + skill):** the bypass flags live in
per-project `headless_args`, never in a vendor CLI's machine-global config;
every mapping carries its workspace constraint (`--add-dir`/`--cd`/cwd); codex
stays on `workspace-write` — this change must not escalate any agent to
full-filesystem danger; secret redaction is orthogonal to permission mode and
stays on (driver `checks:` output is already secret-scrubbed per LE-4;
`runner.SanitizeForContext` keeps reasoning fences out of reused context).
codex-1's point that "a vendor flag named 'dangerous' is not evidence of
workspace confinement" is correct: the resolved profile must DECLARE and
VERIFY confinement, and if no portable outer confinement exists the matrix
should honestly report `ScopeEnforcement: unverified` and fail the mandatory
autonomous profile rather than label it safe. claude's `bypassPermissions` is
the broadest mode; its scoping is only as good as the `--add-dir` + cwd
discipline, so reinit always writes scope+mode as a pair and the CHANGELOG
calls out the built-in default change (kimi-1 concern 6).

**Exact skill wording (normative block to add):**

> Every headless participant MUST be invoked in its non-interactive autonomous
> write mode, so it can WRITE its own canonical artifact
> (`round-NN/<id>.md`, review files, signoff appends) without a blocking
> permission prompt. There is no cross-vendor flag; use the per-CLI mapping
> recorded in the agent's spec (`autonomous_write` with the args above). The
> mode MUST be scoped to the project workspace (`--add-dir`/`--cd`/working
> directory) — never enable a machine-wide bypass, never edit the vendor
> CLI's global config. If a CLI has no autonomous write mode, it cannot be a
> headless participant. Autonomous mode changes permission prompting only:
> secret redaction of outputs, logs, and artifacts still applies in full.

## The display-name rendering path (my lens — concrete)

This is the part that touches the files I own: the §2 table, the TUI, and the
digests. The design must say exactly where the composite appears and where it
does not.

**§2 roster table (`COOPERATION.md:101-110`):** gains a "Display name" column
and a "Speed" column. New shape:

```
| Agent ID       | Display name              | Workspace dir       | Role          | Speed |
| -------------- | ------------------------- | ------------------- | ------------- | ----- |
| `claude-1`     | `claude-opus4.8-max`      | `../claude/`        | facilitator+participant (cli `claude`) | `fast` |
| `codex-1`      | `codex-gpt5.5-xhigh`      | `../codex/`         | participant (cli `codex`)               | `fast` |
| `hermes-1`     | `hermes-glm5.2-high`      | `../hermes/`        | participant (cli `hermes`)              | `fast` |
| `antigravity-1`| `agy-gemini3.5flash-high` | `../antigravity/`   | participant (cli `agy`)                 | `fast` |
| `kimi-1`       | `kimi-k3-max`             | `../kimi/`          | participant (cli `kimi`)                | `fast` |
```

The roster ID stays in the first column (parsed by `ReadRosterIDs`,
`roster.go:17`, charset unchanged). The display name is a new column — never
parsed by `ReadRosterIDs`. The `cli`/`model` info that used to be buried in
the Role prose is now also visible as the display name; the Role prose keeps
the `cli` annotation for readability but drops the redundant `model` prose
(since the display name carries it). The reinit writer touches ONLY this
table block (allowlisted zone shared with the drift guard); the protocol
default and skill fallback copies are re-synced in the same change.

**TUI runtime matrix (`PrintRuntimeMatrix`, `discover.go:314-384`):** add a
`DISPLAY` column showing the composite next to the roster ID. The current
matrix is already wide; the composite is ≤22 chars (`agy-gemini3.5flash-high`),
which fits if we widen the `AGENT` column from 8 to 13 (to fit
`antigravity-1`) and add the display column.

**Round digest (`digest.go:14-34`):** `AgentLine` gains `Display string`
(kimi-1 §6). `BuildRoundDigest` (`digest.go:48-53`) still keys on the
participant string (the roster ID) and reads `<id>.md`; the display is looked
up from config and attached. Add `Model`/`Effort`/`Speed` to `AgentLine` too
(kimi-1's audit proposal) so the digest is self-documenting. `_index.md`
(`round_index.go:96-151`) and review snapshots (`reviewsnapshot.go:67`)
inherit stability for free because they key on the roster ID.

**TUI round summary (`roundsummary.go:84-90`):** the `@%-13s` field is too
narrow for composites (`claude-opus4.8-max` = 18 chars,
`agy-gemini3.5flash-high` = 22). Required follow-up: render the **roster ID**
in the `@`-prefixed slot (stable, ≤14 chars) and the **display name** on the
same line after it (e.g. `@claude-1 claude-opus4.8-max <position>`), or widen
the field. NEVER truncate the roster ID (it is the identity); if anything is
truncated for width, truncate the display name with `…`. This is the first
place the change becomes visible and must be sequenced in the implementation
plan (kimi-1 risk).

**TUI agent tabs (`protosnap.go:253-262`):** key on `runstate.AgentState.ID`
(roster ID) — unaffected. The tab LABEL can show the display name for
readability, but the key stays the roster ID. Kill/steer flows
(`internal/tui/live.go`) address agents by the roster ID — unchanged.

**Consensus (`consensus.go:90,145`):** the signoff regex already tolerates the
composite charset (`[A-Za-z0-9._-]+`), and participant-membership checks stay
exact-match on the roster ID — no change needed. A mid-idea display drift
cannot invalidate a signoff because the signoff uses the roster ID, not the
display name.

**Per-round artifact frontmatter (audit, kimi-1 §1.5):** the round-file writer
prompts gain `model:`, `effort:`, `speed:` frontmatter fields (the values that
actually ran this round), and `validation.go:42-72` tolerates them. This is
the audit path that makes "see which model wrote round-01 two years later"
work — strictly better than encoding the model in the identity, and it does
not churn when the model bumps.

## One coherent design (lift-ready)

1. **Identity = roster ID** (`claude-1`), charset `[a-z0-9][a-z0-9-]*`, stable
   forever, the filename/signoff/`participants:`/frontmatter key. Never
   renamed; zero history rewrites; forward-only (precedent:
   `meta/roster-update_2026-06-19.md`).
2. **Display = composite** (`claude-opus4.8-max`), charset
   `[a-z0-9][a-z0-9.-]*` (dots allowed inside the model token, no `..`, no
   leading/trailing dot, lowercase, no `_`), derived from
   `family`+`model_label`+`effort` at render time, stored optionally in config
   and the §2 table's new "Display name" column, shown in TUI/digest/matrix.
   Never a key, never a filename, never a signoff.
3. **Spec split (codex-1):** `ID` = roster seat, `AdapterID` = family,
   `DisplayName` = composite, `Model`/`Reasoning`/`Speed` authoritative,
   `AutonomousWrite` structured. `selectedAgents` matches `participants:`
   against `ID` (roster seat) with `AdapterID` legacy fallback — this FIXES
   the schism for new ideas without a driver-side resolver and without
   renaming old artifacts.
4. **Component B:** `parley roster init [--scope session|machine]` (distinct
   verb), session default when a deck exists, machine default otherwise.
   Session writes deck `agents.toml` (override layer) + §2 table (allowlisted
   zone, shared drift-guard allowlist) + surgical `headless-agents.local.json`.
   Machine writes only `~/.parley/agents.toml` (managed-block, comment
   preserving, backup on first conversion). Open-idea guard: hard-refuse
   rename/removal touching a live quorum; `--yes` writes the
   `excluded:`/`reincluded:` break; model/effort/speed changes never trip it.
   `parley init` chains into `parley roster init --scope session`; §9.0 stays
   read-only with a drift WARN.
5. **Component C:** first-class `AutonomousWrite` on `Spec` (structured, with
   `Scope`/`ScopeEnforcement`); fix the built-in defaults so they are actually
   autonomous (claude `bypassPermissions`, codex `never`, hermes `--yolo`,
   agy unchanged, kimi new built-in with `-p`); per-CLI mapping in built-ins
   AND skill; workspace-scoped, never machine-wide; secret redaction
   orthogonal and preserved.
6. **Component D:** `speed` vocabulary `{fast, deep}`, `fast` default;
   `DefaultSpeed` → `"fast"`, central template `speed = "fast"`,
   `balanced` → `fast` legacy alias; `speed` is a launch-property axis, never
   part of the composite, never a model/effort override; remove or neuter
   `Profile` to prevent a future downgrade vector; per-idea `speed:` override
   for heavy ideas; `fast` = "no extra thinking overhead," never "weaker
   model."
7. **Audit:** per-round `model:`/`effort:`/`speed:` frontmatter (the values
   that actually ran) + `AgentLine.Model`/`Effort`/`Speed` in the digest.
8. **Tests:** sanitization + parse + collision (table-driven, the 5 names);
   reinit session-vs-machine write targets; open-idea guard
   (drop/rename/add/model-change/speed-change cases); composite-on-load
   derivation; autonomous-mode field round-trip + the built-in default fix;
   `selectedAgents` matching roster-ID-with-family-fallback; drift guard
   (`TestEmbeddedDefaultMatchesLiveDeck`) stays green after the §2 table
   writer lands; `go build ./...`, `go vet`, `gofmt -l` clean.

The one open question I would put to the group before FINAL: **do we remove
`Profile` entirely, or keep it as a documented-obsolete alias of `speed`?**
I lean remove (two fields with overlapping names is the ambiguity that lets a
future "fast profile" reintroduce the downgrade), but it is a backward-compat
call and the group should decide.
