---
agent: kimi-1
idea: composite-agent-naming-and-roster-reinit
round: 1
date: 2026-07-18
---

## Summary

Position: adopt the composite `agent-model-effort` name as a **derived display name**, never
as the agent's canonical ID (CRUX: option (b)). The ID stays the stable roster key
(`claude-1`); model+effort live in one source of truth (the layered `agents.toml`) and the
composite is recomputed at render time — so it can never drift, lie, or rewrite an OPEN
idea's quorum. Component B should be a new `parley roster init` verb (session scope writes
the deck layer, machine scope writes only `~/.parley/agents.toml`) that *is* the §0
bootstrap confirmation flow, reuses the §9.0 probe, and hard-refuses ID changes referenced
by any open idea. Component C should be a first-class `autonomous_write` assertion on
`agents.Spec` plus a per-CLI mapping table in the skill — with built-in defaults changed so
they actually are autonomous (today claude ships `acceptEdits`, codex ships `on-failure`).

Fresh-eyes finding the brief misses: the deck already runs **two ID namespaces**. The Go
driver keys everything on the spec ID (`claude`, `agy`), while §2, `00-prompt.md`, and
`headless-agents.local.json` use `claude-1`, `antigravity-1`. Any naming scheme that does
not say which namespace it joins will bake the schism in deeper.

## Proposed approach

### 0. What the code actually does today (fresh-eyes trace)

I traced every place an agent identity crosses a boundary. The ID model must hold in all
of them:

- **Artifact path** — `internal/runner/runner.go:349`:
  `outputPath := filepath.Join(opts.Idea.Path, opts.RoundLabel, agent.ID+".md")`, where
  `agent.ID` is the **spec ID** (`claude`, `codex`, `agy`, `hermes` from
  `internal/agents/discover.go:105-222`), not the §2 roster ID. `selectedAgents`
  (`runner.go:327-342`) matches `participants:` against spec IDs by exact string equality —
  so a driver run over this very idea's `participants: [claude-1, …]` would match *nothing*.
  The skill flow writes roster-ID files by hand; the driver writes spec-ID files. Two
  namespaces, no resolver.
- **Frontmatter validation** — `internal/runner/validation.go:42-72`: `agent:` must equal
  the invoking ID exactly, or the artifact is rejected.
- **Signoffs** — `internal/consensus/consensus.go:90` parses
  `### Signoff: <id>` with `[A-Za-z0-9._-]+` (note: already tolerates `.` and `_`), but
  `consensus.go:145-146` and `:383-388` require the signer to be an exact member of
  `00-prompt.md`'s `participants:` — otherwise "unknown participant".
- **§2 roster parse** — `internal/protocol/roster.go:17`: roster rows must match
  `` `[a-z0-9][a-z0-9-]*` `` — **lowercase + hyphens only**. An ID containing uppercase or
  `_` is silently dropped from `ReadRosterIDs`, which breaks roster-preset validation
  (`internal/config/roster.go:110-115`).
- **Round digest** — `internal/driver/digest.go:48-53`: `AgentLine.Agent` is the raw
  participant string; the digest reads `<participant>.md` and has no model/effort field.
- **TUI** — `internal/tui/roundsummary.go:84-90` renders digest rows with `@"%-13s"`;
  `internal/tui/protosnap.go:253-262` keys agent tabs on `runstate.AgentState.ID`
  (`internal/runstate/runstate.go:50-52`); kill/steer flows (`internal/tui/live.go`) address
  agents by the same ID.
- **Review snapshots & logs** — `internal/runner/reviewsnapshot.go:67` and
  `runs/<run>/agents/<id>/` are keyed by agent ID; `_index.md`
  (`internal/runner/round_index.go:96-151`) rows are keyed by `AgentID`.
- **Config layers** — `internal/config/runtime.go:127-144`: central `~/.parley/agents.toml`
  < deck `agents.toml` < deck `agents.local.toml` < env. `[agents.<family>]` is keyed by
  spec family; unknown families are appended as new specs (`runtime.go:361-381`) — which is
  how `kimi` enters at all, since it has **no built-in spec** today.
- **`parley init`** — `internal/app/app.go:362-397`: writes the central seed
  (`EnsureCentralDefault`, `runtime.go:251-272`) + workspace skeleton. Non-interactive; the
  interactive model/effort confirmation currently lives only in the skill.
- **§9.0 ping** — `internal/app/preflight.go`: `checkRoster` → `hostedPONG` →
  `runner.CommandFor` runs each agent's *real configured invocation*; preflight never
  selects models (per §0).
- **Precedent for renames** — `parley-deck/meta/roster-update_2026-06-19.md`: the deck was
  already renamed once (`agy` → `antigravity-1`, etc.) under direct user instruction,
  **forward-only** ("from the next idea onward"); old ideas kept their old-ID artifacts.
  That file even carries a "Write mode" column — component C's mapping table already exists
  in embryo there.

### 1. The CRUX: the composite must NOT be the ID

I pick **option (b): a stable base ID (`claude-1`) as the canonical key; the composite is a
derived display name** shown in the TUI, digests, §2 roster annotations, and run output.

Justification, grounded in the trace above:

1. **The risk of a mid-idea rename is identity split-brain plus a smuggled quorum change.**
   Every cross-reference in the protocol is an exact string: artifact paths, `agent:`
   frontmatter, `responding-to:` lists, `### Signoff:` blocks, inbox filenames
   (`<from>-to-<to>_<topic>.md`), commit-message prefixes, review-snapshot dirs, TUI tabs.
   If the composite IS the ID, then flipping `claude` from Opus 4.8/max to Sonnet/high
   mid-idea renames `claude-opus48-max` → `claude-sonnet-high`, and in one move: (i) the
   OPEN idea's `participants:` (locked at Phase 0, §5) no longer contains any matching ID —
   the next signoff append fails `consensus.go:145` "unknown participant", and
   `consensus status` reports the old signer's ✅ as dangling; (ii) the driver's
   "artifact already exists" skip (`runner.go:361-374`) misfires — the old file sits under
   the old name, the new name looks missing, the agent re-runs and the idea now has two
   identities' artifacts for one seat; (iii) every `@claude-opus48-max` reference in prior
   rounds dangles. Renaming an agent mid-idea is therefore indistinguishable from
   **excluding one participant and adding another without the §9.0 user-confirmed gate** —
   a silent rewrite of an OPEN idea's quorum, which the brief itself lists as a hard
   requirement not to do. On case-insensitive filesystems (default APFS on macOS), a
   case-only rename additionally aliases two IDs onto one file.
2. **A stored composite name can lie; a derived one cannot.** If the name is the ID, it is
   written once (at reinit) and hand-edited config later (`agents.toml` says `opus` but the
   ID still says `opus48-max`… which opus? which effort?) forks the display from reality.
   Under (b) the display is recomputed from the same `Spec.Model`/`Spec.Reasoning` the
   runner actually passes on the command line — the single source of truth stays the
   layered config, and the display self-heals on every render.
3. **What the user actually wants — "a run is self-documenting at a glance" — is a display
   property**, fully satisfied by rendering `claude-1 · claude-opus48-max` in the §2 table,
   TUI, and digests, without paying any rename cost.
4. **Option (c) (stable prefix + descriptor) is the worst of both.** It still renames on
   every model/effort change (paths churn), and the moment tooling parses the prefix back
   out to find the real key, the prefix *is* the ID and the descriptor *is* a display name —
   you have reinvented (b) with a fragile parse step in the middle.
5. **The audit argument for (a) is better served elsewhere.** "See which model wrote
   round-01 two years later" belongs in the artifact, not the key: add `model:` /
   `effort:` to the round-file frontmatter (cheap: extend the writer prompts and tolerate
   the fields in `validation.go`) and add `Model`/`Effort` to `digest.AgentLine`
   (`internal/driver/digest.go:15-20`). Frontmatter records what *actually ran* per round —
   strictly better audit than a name, which records what was *rostered*.

Concrete display integration points: §2 roster table gains a `Display name` column (the §2
table is a project zone under preflight's `mergePreservingZones`,
`internal/app/preflight.go:524-560`, so editing it is safe); `digest.AgentLine` gains
`Display`; `roundsummary.go` renders the composite (see §6 for the width fix it forces);
`parley agents list` (`PrintRuntimeMatrix`, `discover.go:314-384`) prints it per family.

### 2. Component A — naming scheme (sanitize, parse, collide)

**Charset decision (narrowing the hard constraint, deliberately).** The user allows
`[a-zA-Z0-9_-]`; I choose the canonical form to be the **lowercase-hyphen subset
`[a-z0-9-]`**, because (i) the §2 roster parser accepts nothing else
(`protocol/roster.go:17`), (ii) case-insensitive filesystems make uppercase a collision
vector, and (iii) mixing `_` into a hyphen-joined composite makes the tokenization below
ambiguous for no readability win. Display names are always lowercase; the human model label
keeps its case only in prose contexts (§2 "Model" column, TUI detail panes).

**Sanitization (write-time, per section).** Each of the three sections is produced
independently:

1. Take the *human model label* (see the new `model_label` field below) — **after removing
   any parenthesized tier annotation** (agy: `Gemini 3.5 Flash (High)` → `Gemini 3.5
   Flash`; the stripped tier becomes the effort token).
2. Lowercase ASCII; delete every character not in `[a-z0-9]` (dots, spaces, slashes,
   brackets vanish; adjacent alphanumerics join). No hyphen is ever introduced *inside* a
   section — that is what keeps the whole name split-clean.
3. Effort is not sanitized from free text; it is **selected from a fixed vocabulary**:
   `low | medium | high | xhigh | max | ultracode | clidefault`.
4. Agent/family section: the **spec ID** (already `[a-z0-9-]` by construction; reinit must
   reject or sanitize any configured family outside it). Note the roster ID and the family
   differ today (`antigravity-1` vs spec `agy`) — the composite uses the **family**
   (`agy-…`), because model/effort are properties of the CLI runtime, not of the seat.

**Parse rule (read-time, fail-closed, right-to-left):**

- Split on `-`. 1 token → bare family (legacy spec ID). Exactly 2 tokens with token[1]
  all-digits → legacy instance ID `family-N` (no model/effort). 3 or 4 tokens → composite.
- In a composite: token[0] = family. If the last token is all digits it is the instance
  index and is stripped; the next token from the right must be a member of the effort
  vocabulary (hard error otherwise — do *not* infer effort as "whatever isn't digits",
  which misparses an all-digit model token like a hypothetical `530`); the single remaining
  middle token is the model.
- Anything else → parse error. Tooling must never guess.

This tightens claude-1's sketch in two places: effort is validated by vocabulary
membership (so all-digit model tokens parse correctly), and the legacy 2-token `family-N`
form is special-cased (claude-1's rule reads `claude-1` as agent=`claude`,
instance=`1`, effort=`claude` — broken).

**Collision rules (deterministic, write-time):**

- Within one roster scope, composites are allocated in catalog order (built-in spec order,
  then config-append order from `runtime.go:361-381`): first occurrence keeps the bare
  name; later duplicates of the same family+model+effort get `-2`, `-3`, … (which is
  exactly the parse rule's instance suffix, so no second mechanism).
- Effort vocabulary can never be all-digits, so the instance suffix is unambiguous.
- A composite that collides with an existing *legacy* ID (`family-N`) or another scope's
  name is not an error — namespaces are per-deck; but reinit must warn if a generated
  composite equals an existing *different* roster entry.

**The five names under these rules** (display names; canonical IDs unchanged):

| Canonical ID    | Family   | Model label (source of truth)        | Effort token | Display name              |
| --------------- | -------- | ------------------------------------ | ------------ | ------------------------- |
| `claude-1`      | `claude` | `Opus 4.8` → `opus48`                | `max`        | `claude-opus48-max`       |
| `codex-1`       | `codex`  | `GPT-5.5` → `gpt55`                  | `xhigh`      | `codex-gpt55-xhigh`       |
| `hermes-1`      | `hermes` | `GLM 5.2` → `glm52`                  | `high`       | `hermes-glm52-high`       |
| `antigravity-1` | `agy`    | `Gemini 3.5 Flash (High)` → strip `(High)` → `gemini35flash` | `high` (from the stripped tier) | `agy-gemini35flash-high` |
| `kimi-1`        | `kimi`   | `K3` → `k3`                          | `max`        | `kimi-k3-max`             |

**Counter-position on the model token's source:** the brief says "model derives from a
human model label, not the raw id" — I agree, but today no such label field exists, and
the config contradicts itself: deck `agents.toml` pins `model = "opus"` for claude and
`model = "glm-5p2"` for hermes, while `headless-agents.local.json` says `GLM 5.2` and the
built-in spec says `claude-opus-4-8[1m]` (which sanitizes to the unreadable
`claude-opus-4-81m`). So: add an explicit `model_label` (human label) to `agents.Spec` /
`agentOverride`, written by reinit, distinct from `model` (the value passed to the CLI
flag). And one honesty check the framing glosses: **codex is configured `cli-default` in
both deck layers today** — the honest current display is `codex-clidefault-clidefault`.
`codex-gpt55-xhigh` becomes truthful only after reinit *pins* `model`/`reasoning` for codex
(§0's bootstrap rule already demands this confirmation). The name must follow config truth,
never wishes; under (b) it mechanically does.

### 3. Component B — `parley roster init`

**Command surface** (new verb group; `parley agents list|verify` stays runtime-discovery):

```
parley roster init [--dir DIR] [--scope session|machine]        # default: session
                   [--agent <family>]...                         # subset; default all discovered
                   [--from <file>]                               # non-interactive: family → {model, effort}
                   [--dry-run] [--yes] [--json]
parley roster show [--scope session|machine] [--json]           # effective roster + display names
```

`--reinit` is not a separate flag: the command *is* the re-init; without `--yes`/`--from`
it runs interactively and existing values are the defaults shown in the picker. It refuses
to overwrite a scope that already has roster entries unless interactive-confirmed or
`--yes` (no silent overwrite).

**Flow:** `agents.Discover(config.LoadAgentSpecs(root))` (this picks up built-ins *and*
config-only families like kimi) → per-family model/effort menu from a new
`ModelCatalog(family)` adapter in `internal/agents` (real enumeration where a CLI offers it
— e.g. `agy models`; a curated static catalog where it doesn't; `cli-default` always
offered; undiscoverable → free-text + `clidefault`, exactly §0's fallback) → user picks →
write → run the §9.0 probe (`checkRoster`/`hostedPONG`, reused from
`internal/app/preflight.go`) so a dead CLI cannot be rostered.

**Scope writes differ precisely:**

- `--scope session` (default): writes the **tracked deck layer** `parley-deck/agents.toml`
  (`[agents.<family>]`: `model`, `model_label`, `reasoning`, `headless_args`,
  `sandbox_mode`, `approval_policy`, `autonomous_write`); rewrites the §2 roster table's
  annotation columns (model/display name — a project zone, preserved across protocol
  syncs); and **surgically** updates `meta/headless-agents.local.json` (model/thinking/
  writeMode fields for rostered IDs only, and only if the file exists). Counter-position:
  the brief bundles the local JSON into session scope wholesale — no; that file is
  gitignored operator-local state (§2), bulk-rewriting it from a shared command would stomp
  per-machine CLI paths. Surgical-only.
- `--scope machine`: writes **only** `~/.parley/agents.toml` (same TOML shape; reuse/extend
  `centralDefaultTemplate`, `runtime.go:277-313`). It never touches the deck, and deck
  values are never copied up. Precedence needs no new logic — `configLayers`
  (`runtime.go:127-144`) already resolves central < deck < local < env; machine scope
  defines inherited defaults, session scope defines this deck's truth.

**Reconciliation (no duplication, no silent quorum rewrite):**

- *vs `parley init`* — init stays bootstrap-only (central seed + skeleton) and prints
  `parley roster init` as the required next step; the skill's §0 "MUST confirm roster +
  model + effort" gate is satisfied **by running this command**. One selection flow, one
  writer; init never grows a second, divergent picker. (Optionally `parley init
  --with-roster` chains it for convenience.)
- *vs §9.0 preflight* — preflight stays **read-only** over roster selection: it pings
  liveness through the real configured invocation and gates on availability; it never
  re-selects models (protocol already says so). The only addition: a non-gating WARN when a
  rostered model/effort is no longer in the family catalog ("roster drift — run `parley
  roster init`"). Write path = roster init; verify path = preflight. Both share one probe
  implementation.
- *OPEN-idea guard (the hard rule)* — before any write, scan `ideas/*/00-prompt.md`
  (same frontmatter reader as `internal/protocol/workspace.go`) for `status:` not in
  `{final, abandoned, complete}` and collect their `participants:`. **Any ID change
  (rename/removal) touching an open-idea participant: hard refuse**, printing the blocking
  ideas. Model/effort value changes with IDs untouched: allowed, with a printed note that
  open ideas keep their locked quorum (§5) and the change applies from the next idea. This
  is the mechanical twin of §9.0's user-confirmed exclusion gate — roster mutation is never
  silent where quorum is live.

### 4. Component C — autonomous ("yolo") write mode

**Representation (first-class, per agent):** add `AutonomousWrite bool` (TOML
`autonomous_write`) to `agents.Spec` — the *assertion* that this agent's headless
invocation writes workspace-scoped artifacts with zero blocking prompts. The *mechanism*
stays in the existing fields (`HeadlessArgs`, `SandboxMode`, `ApprovalPolicy`,
`internal/agents/discover.go:28-29`); do not invent a parallel flag vector. Bookkeeping
touched: `agentOverride` (`internal/config/runtime.go:88-122`), the `withBuiltinSources`
field list (`discover.go:224-259`), `centralDefaultTemplate` (`runtime.go:277-313`), and an
`AUTO` column in `PrintRuntimeMatrix` (`discover.go:314-384`). `roster init` sets the bit
only after checking the configured args match the family's known autonomous mapping
(static check per family; a live sentinel-write probe is available via `parley agents
verify --full` as an opt-in, not a default — probes are slow and cost tokens). Preflight
warns (not gates) when a rostered agent lacks `autonomous_write = true`.

**Fix the built-ins — today they are not all autonomous** (fresh-eyes finding):
claude's built-in ships `--permission-mode acceptEdits` (`discover.go:134`), which still
prompts outside edits; codex ships `approval_policy="on-failure"` (`discover.go:113`),
which can still ask; hermes' built-in lacks `--yolo` (it exists only in this deck's local
JSON); kimi has no built-in spec at all. The mapping, to be encoded both in built-ins and
in the skill:

| CLI    | Autonomous invocation (workspace-scoped)                                             |
| ------ | ------------------------------------------------------------------------------------ |
| claude | `--permission-mode bypassPermissions --add-dir {root}` (scope = the add-dir + cwd)   |
| codex  | `exec --sandbox workspace-write -c approval_policy="never"` (`--cd {root}`)          |
| hermes | `--yolo --accept-hooks` (cwd = root)                                                 |
| agy    | `--dangerously-skip-permissions --add-dir {root}`                                    |
| kimi   | plain `-p` print mode auto-approves in-workspace writes; **never** add `--yolo`/`--auto` — mutually exclusive with `-p` |

Safety invariants (encode in code + skill): the bypass flags live in **per-project**
`headless_args`, never in a vendor CLI's machine-global config; every mapping carries its
workspace constraint (`--add-dir` / `--cd` / cwd); codex stays on `workspace-write` — this
change must not escalate any agent to full-filesystem danger; secret redaction is
orthogonal to permission mode and stays on (driver `checks:` output is already
secret-scrubbed per LE-4; `runner.SanitizeForContext` keeps reasoning fences out of reused
context).

**Exact skill wording (normative block to add):**

> Every headless participant MUST be invoked in its non-interactive autonomous write mode,
> so it can WRITE its own canonical artifact (`round-NN/<id>.md`, review files, signoff
> appends) without a blocking permission prompt. There is no cross-vendor flag; use the
> per-CLI mapping recorded in the agent's spec (`autonomous_write = true` with the args
> above). The mode MUST be scoped to the project workspace (`--add-dir`/`--cd`/working
> directory) — never enable a machine-wide bypass, never edit the vendor CLI's global
> config. If a CLI has no autonomous write mode, it cannot be a headless participant.
> Autonomous mode changes permission prompting only: secret redaction of outputs, logs,
> and artifacts still applies in full.

### 5. Migration path for existing decks and ideas (my lens)

1. **Zero renames.** Existing roster IDs (`claude-1`, …) stay canonical forever. Display
   names are derived from config at render time, so adoption touches no historical file.
2. **Forward-only, per the established precedent** (`meta/roster-update_2026-06-19.md`):
   from ratification, new §2 annotations/digests/TUI show composites; old ideas keep their
   old artifacts and signoff strings untouched. No `git mv`, no history rewrite, no
   alias table.
3. **If a deck later insists on composite-as-canonical** (I advise against): only via §2's
   existing join/leave mechanics — add the new ID row, mark the old row *inactive* (do not
   delete; "historical references remain resolvable"), record user confirmation in
   `meta/roster-update_<date>.md`, and only when the open-idea guard (§3) finds zero
   references. Do not build a new aliasing subsystem for this.
4. **Namespace schism (driver spec IDs vs deck roster IDs) — contain, don't widen.**
   `roster init` must write the roster-ID ↔ family mapping **explicitly** (§2 `cli` column
   + the local JSON's per-entry `cli` field), never rely on a prefix heuristic —
   `antigravity-1` ↔ `agy` already breaks naive prefixing. Fixing the schism itself
   (e.g. a driver-side resolver `family-N` → spec, fail-closed on ambiguity) is a separate
   idea; this one must not depend on it.
5. **Config backfill:** reinit adds `model_label` + `autonomous_write` where absent;
   missing `model_label` falls back to sanitizing `model` (ugly but honest), and the
   fallback is visible in `--dry-run` output so the user can correct the label.

### 6. Verification: does the (b) ID model hold across TUI, driver, digests?

Yes, with two required display-layer changes and one pre-existing asymmetry noted:

- **Driver/digests:** `BuildRoundDigest` keys on the participant string and reads
  `<id>.md` (`digest.go:48-53`) — stable IDs keep position maps coherent across model
  changes; add `Display`/`Model`/`Effort` to `AgentLine` rather than re-keying. `_index.md`
  (`round_index.go`) and review snapshots (`reviewsnapshot.go:67`) inherit stability for
  free.
- **TUI:** agent tabs/`protosnap` key on `AgentState.ID` — unaffected. The digest panel's
  `@"%-13s"` (`roundsummary.go:84-90`) is too narrow for 17–22-char composites
  (`claude-opus48-max` = 17, `agy-gemini35flash-high` = 22): required follow-up — widen the
  field or truncate the display with `…`; never truncate the canonical ID.
- **Consensus:** the signoff regex already tolerates the composite charset
  (`consensus.go:90`), and participant-membership checks stay exact-match on the stable ID —
  no change needed, and mid-idea display drift cannot invalidate a signoff.
- **Asymmetry to record (not introduced by this design):** the roster parser is stricter
  (`[a-z0-9-]`, `protocol/roster.go:17`) than the signoff parser (`[A-Za-z0-9._-]+`,
  `consensus.go:90`). My lowercase-hyphen canonical form sits in their intersection, so it
  parses everywhere; document the intersection as the ID charset rule.

## Concerns / open questions

1. The two-namespace schism (§0 trace) is the largest latent risk in this area and is out
   of scope here; at minimum this idea must ratify "composite derives from the **family
   (spec ID)** layer" and "roster IDs stay the deck's keys", plus the explicit mapping
   write in reinit. A driver-side resolver deserves its own idea.
2. Model catalogs go stale (vendor releases); the catalog is a cache refreshed at reinit,
   never a runtime dependency. Display derivation must read only parley config — never
   parse vendor files (`config.yaml`, `config.toml`) at render time; for hermes/kimi the
   recorded effort is a claim made at reinit, and the §9.0 drift WARN is the staleness
   tripwire.
3. agy's tier-in-label parse (`(High)` → `high`) is brittle against vendor label changes;
   store the tier at reinit and treat label-parsing as a one-time suggestion, not a render
   rule.
4. Effort vocabulary extensibility: new vendor tiers (beyond the seven listed) require a
   vocabulary bump; parse fails closed on unknown tokens, so a stale parser errors loudly
   rather than misreading.
5. How interactive should `roster init` be inside an agent-driven session (the facilitator
   is itself an agent)? `--from`/`--yes` cover non-interactive runs; the interactive TUI
   picker is for the human operator. The skill should route accordingly.
6. `claude` currently runs `bypassPermissions` per this deck's local JSON while the
   built-in says `acceptEdits` — ratifying C means changing a built-in default with a
   genuine security posture difference; the workspace-scoping flags must land in the same
   commit, and the CHANGELOG should call it out.

## Risks

- **Rename temptation.** Once display names exist, users will ask to "make the paths
  pretty". That path exists (§5.3) but must stay manual, confirmed, and guarded by the
  open-idea check — otherwise we reintroduce split-brain through the back door.
- **claude's `bypassPermissions` is the broadest vendor mode**; its scoping is only as good
  as the `--add-dir` + cwd discipline. Mitigation: reinit always writes scope+mode as a
  pair; the skill text states the invariant; verification probes run in a scratch subdir.
- **Display drift by hand-edit.** Users editing `agents.toml` directly change the display —
  acceptable and self-consistent under (b) (the display simply follows the new truth), but
  the §2 table's rendered copy can go stale until the next reinit; `roster show --dry-run`
  should diff it.
- **TUI layout regression** until `roundsummary.go` is widened — cosmetic, but it is the
  first place the change becomes visible; sequence it in the implementation plan.
- **Probe cost.** Making the sentinel-write probe mandatory would slow every reinit and
  burn tokens on hosted backends; keep it opt-in (`verify --full`) and rely on the static
  mapping check by default.
- **Scope confusion.** A user running `--scope machine` expecting per-project effect (or
  vice versa) gets surprised; mitigate with `--dry-run` printing exactly which files change
  per scope, and the command refusing to run outside a deck for session scope.
