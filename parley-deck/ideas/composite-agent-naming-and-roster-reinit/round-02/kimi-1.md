---
agent: kimi-1
idea: composite-agent-naming-and-roster-reinit
round: 2
date: 2026-07-18
role: end-to-end ID integrity across runner/validation/consensus/digest/TUI/runstate/snapshots/config + two-namespace resolution and migration
---

## Summary

Cross-review of all four round-01s, folding in both late updates (dot charset; component D).
Final positions, each argued from the code:

1. **CRUX: option (b) stands — stable roster ID is identity, composite is a derived display
   name.** claude-1's frozen-slug is (a)-with-a-freeze and fails on two code facts it did not
   trace: `internal/protocol/roster.go:17` cannot parse the dotted slugs it would mint, and
   `internal/runner/runner.go:327-342` cannot resolve them. What claude-1 gets *right* (the
   idea-time freeze) I adopt — moved one layer down, into codex-1's per-idea profile snapshot.
2. **Dot charset: adopt claude-1's dot-preserving sanitization** (now user-mandated), with the
   whole-name grammar below. Canonical roster IDs stay `[a-z0-9-]` — composites never enter the
   §2 ID cell, so `roster.go` needs **no change**; `consensus.go:90` already tolerates the full
   display charset. The design is consistent with both files as they are.
3. **Two-namespace schism: I change my round-01 position — this feature should FIX it, not
   merely avoid deepening it.** hermes-1's direction + codex-1's type split make the fix small,
   and the deck's own artifacts already live in the roster-ID namespace, so the fix aligns the
   driver with reality instead of migrating history. Concrete Go touch points in §4.
4. **B: distinct `parley roster` verb group** (hermes-1/codex-1/me), sharing one internal
   service with `parley init` (claude-1's valid dedup point). Session = deck layer, machine =
   `~/.parley/agents.toml` only.
5. **C: first-class `autonomous_write` struct + fix the built-ins** (my round-01 finding,
   verified: claude ships `acceptEdits` at `discover.go:134`, codex ships `on-failure` at
   `discover.go:113`, hermes lacks `--yolo` at `discover.go:203`). Mapping per the brief.
6. **D: `fast` is a pure label axis in Go today — the downgrade lives in the deck JSON's
   `profiles.fast`** (verified at `headless-agents.local.json:28-41`, `:120-133`). Fix =
   redefined semantics + template change + reinit rewriting/removing `profiles.*`, never
   touching model/effort, never in the name.

## 1. The CRUX, settled

Three of four round-01s (codex-1, hermes-1, kimi-1) chose (b): composite = display, stable ID
= identity. claude-1 chose frozen-slug: the composite IS the name, frozen per idea at creation.

### Addressing claude-1 — frozen-slug fails on two code facts; its freeze survives as a snapshot

**Where I agree:**
- The composite must be what a human sees everywhere: §2 table, TUI, digests, run headers.
  Under (b) it is — rendered, not keyed.
- The idea-time freeze is a real requirement. An idea opened under `opus4.8/max` should still
  *show* that profile after the roster moves on. I adopt this — as data, not as a key (below).
- Dot-preserving sanitization (`opus4.8` not `opus48`): now user-mandated, adopted verbatim (§2).
- One selection engine shared by bootstrap and reinit: adopted (§5).

**Where I disagree — and the counter-proposal:**

1. **A dotted frozen slug cannot be a §2 roster ID, full stop.** `roster.go:17` matches
   `` `[a-z0-9][a-z0-9-]*` `` — no dots. `claude-opus4.8-max` in the §2 first cell is *silently
   dropped* by `ReadRosterIDs`, which breaks roster-preset validation
   (`internal/config/roster.go:110-115`) by fail-closed design. The brief explicitly says the
   design must be consistent with this line. Frozen-slug therefore requires widening the §2 ID
   charset — a protocol-level change with new ambiguity surface (`agent-id` vs display strings
   become indistinguishable), bought for zero integrity gain.
   **Counter:** roster IDs stay `[a-z0-9-]`; the composite rides in a new §2 *Display name*
   column (an annotation cell, never parsed by `rosterRowRe`, which anchors on the first cell).
2. **The driver cannot resolve frozen slugs.** `selectedAgents` (`runner.go:327-342`) matches
   `participants:` against spec IDs by exact string equality. Stamping composites into new
   ideas' `participants:` means the driver selects *zero* agents for every new idea — the same
   schism this idea already suffers, now baked in deeper. The brief's minimum bar is "do not
   deepen the schism"; frozen-slug fails it.
   **Counter:** `participants:` keeps stable roster IDs; the driver learns to resolve them (§4).
3. **The freeze belongs in the snapshot layer.** claude-1 wants "changing model/effort later
   does not rewrite an OPEN idea" — codex-1's `participant-profiles:` frontmatter in
   `00-prompt.md` plus exact profiles in `runmanifest.Manifest` delivers exactly that, and
   *additionally* freezes the profile for audit (which model wrote round-01) without making a
   mutable string the lookup key. Frozen-slug freezes the *name*; the snapshot freezes the
   *facts*. The facts are what audit needs.
   **Counter (the merge): (b) + codex-1's immutable per-idea profile snapshot = claude-1's
   freeze, achieved at the display layer.** New ideas stamp `participant-profiles:
   [claude-1=claude-opus4.8-max, …]` at Phase 0; TUI/digest/§2 render the live composite;
   the idea's own pages render its frozen snapshot. No rename path exists, so no rename bug
   exists.
4. claude-1's migration (legacy alias `claude-1` → base `claude`) goes the wrong way: `claude-1`
   is not the legacy thing, it is the *roster* thing; the spec family `claude` is the legacy
   runtime key. §4's resolver covers this without an alias table.

### Addressing codex-1 — agree with the architecture; three counters

**Agree:** the `Spec` split (`ID` / `AdapterID` / `DisplayName` / `Model` / `Reasoning`) is the
correct Go shape and is exactly what §4's schism fix needs; vendor branches
(`cleanParticipantEnv`, `isolatedAgentHome`) must key off `AdapterID`; rerunnable `roster init`
with no `--reinit` flag; managed-block writes + backup for machine scope; per-run manifest
snapshot preferred on resume; the `agents.local.toml` vs `headless-agents.local.json`
single-authority concern (settled in §5); fail-closed verification of autonomous profiles.

**Disagree + counter:**
1. **Dot deletion is superseded.** codex-1's grammar (`agent/model := [a-z0-9]+`, dots deleted)
   was a defensible reading of the pre-update brief; the user's late decision allows and wants
   dots (`claude-opus4.8-max`). **Counter:** §2's dot-preserving grammar — still narrower than
   the user's full charset, still fail-closed, now with natural version numbers. codex-1's
   "generate from operator-confirmed labels, never reverse the slug" and "lossy by design"
   principles are unchanged and I keep them.
2. **Do not store `display_name` in config; derive it.** codex-1's TOML sketch stores
   `display_name = "codex-gpt55-xhigh"`. A stored display can lie (hand-edit `model` and the
   stored name forks from reality — my round-01 §1.2). **Counter:** store only the inputs
   (`model`, `model_label`, `reasoning`), derive the display at render time from the same
   fields the runner passes; accept an *optional* `display_name` override for genuinely
   unlabelable cases, defaulting to derived. hermes-1 (concern 2) leans the same way.
3. **`ScopeEnforcement: cli-sandbox|outer-sandbox` is one field too many.** The enum implies a
   verification capability we will not build in this idea. **Counter:** `AutonomousWrite{Mode,
   Args, Scope}` where `Scope` is always `"workspace"` for rosterable agents, plus codex-1's
   *honesty requirement* kept as a verification note: where workspace confinement cannot be
   demonstrated, `roster init` marks the profile `unverified` and refuses to set the bit
   (fail-closed) — same outcome, no speculative taxonomy.

### Addressing hermes-1 — agree with the UX and the diagnosis; two counters on scope

**Agree:** distinct `parley roster` verb group (`init|show|diff`); session-default scope inside
a deck, machine outside; the open-idea guard is non-optional; derive-on-load display; the §2
table writer must share one allowlist definition with the drift guard; exit codes 0/1/2/3;
post-write liveness probe reusing preflight's `pingProbe`; the `agy` tier-in-label handling.
hermes-1's trace of the runner↔§2 divorce (`runner.go:349` writes the *spec* ID; §2's
`claude-1` is consulted only by preset validation) is correct — I verified it independently.

**Disagree + counter:**
1. **The runner rewire is right, the framing is wrong.** hermes-1 presents "wire the runner to
   the §2 base ID" as a costly migration (renames `round-01/claude.md` → `round-01/claude-1.md`
   for every existing deck). But that cost assumes existing artifacts are spec-ID-keyed — in
   *this* deck they are not: the skill flow already writes `round-01/claude-1.md` by hand
   (this idea's own round-01 files are the proof), and the driver has *never* successfully run
   over this deck's `participants: [claude-1, …]` (exact-match finds nothing). Making the
   participant string the artifact identity **aligns the driver with what decks already do** —
   for deck-native rosters it is a zero-rename fix, and for legacy driver-run decks
   (`participants: [claude]`) resolution rule 1 (exact spec-ID match) keeps them working.
   **Counter:** adopt the rewire, forward-only, with the fail-closed resolver in §4 — no
   retroactive renames, no `--yes`-driven history edits. This is my answer to the brief's
   schism question: **fix it, minimally, in this idea** (reversing my round-01 §5.4 "separate
   idea" position — see self-review below).
2. **`--yes` must not break quorum by default.** hermes-1 lets `--yes` proceed on a drop and
   record `excluded:` lines. Recording a break is better than hiding it, but codex-1's deferral
   is better than either: a dropped ID referenced by an open idea becomes *inactive for new
   ideas*, pinned with its profile until the idea goes terminal — no break ever needs
   recording. **Counter (the merge):** default = hard refuse; `--yes` = convert drop →
   inactive-retention (codex-1); only an explicit `--force-drop` writes hermes-1's
   `excluded:`/`reincluded:` lines and proceeds. Three positions, one mechanism each, ordered
   by blast radius.

### Addressing kimi-1 (self-review) — what I uphold, what I change

- **Uphold:** option (b); derive-don't-store; the namespace trace; `model_label` as a distinct
  field (raw IDs like `claude-opus-4-8[1m]` sanitize to garbage under *any* rule); surgical-only
  writes to the local JSON; the built-ins-aren't-autonomous finding; forward-only migration per
  `meta/roster-update_2026-06-19.md`; the TUI width follow-up (`roundsummary.go:84-90`,
  `@"%-13s"` — now 23 chars for `agy-gemini3.5flash-high`; widen or truncate the *display*,
  never the ID).
- **Change 1 — the dot.** I narrowed the canonical form to `[a-z0-9-]`, citing `roster.go:17`.
  The user has since decided dots are allowed. My `roster.go` argument was about the *ID cell*;
  it does not constrain the display column. I accept the dot (§2) and keep only the lowercase
  canonical form (APFS case-folding makes uppercase a collision vector; the user charset allows
  case but nothing requires generating it).
- **Change 2 — the schism.** I said "contain, don't widen; fixing is a separate idea." The
  brief now forces the question, and codex-1's `AdapterID` split removes my cost objection:
  the resolver + identity fix is a few touch points, all inside files this idea already edits
  (`runner.go`, `validation.go`, `runtime.go`). Two namespaces with a documented resolver is
  strictly worse than one identity + adapter templates. **I flip: fix it here, minimally (§4).**

## 2. Late update 1 — the dot charset: final naming grammar

**Charset:** user allows `[a-zA-Z0-9_.-]`. Canonical generated form is lowercase, dots only
*inside* the model section, hyphens only *between* sections, no underscores generated (accepted
on input). Path-safety per the brief: no `..`, no leading/trailing dot — enforced per section,
which implies it for the whole name.

```text
display-name := family "-" model "-" effort [ "-" instance ]
family       := [a-z0-9]+                 # spec family (agy, claude, …), hyphen-free by construction
model        := [a-z0-9]+ ("." [a-z0-9]+)* # dots only between alphanumerics → no "..", no edge dots
effort       := low|medium|high|xhigh|max|ultracode|clidefault   # closed vocabulary
instance     := [2-9][0-9]*               # >= 2, collision suffix
```

**Sanitization (per section, write-time):** take the human `model_label` (after stripping a
parenthesized tier — agy); lowercase ASCII; delete every char not in `[a-z0-9.]`; collapse dot
runs to one; strip leading/trailing dots; reject an empty result (ask for an explicit
`name_token`; never invent). Effort is *selected* from the vocabulary, never sanitized from free
text. `cli-default` → `clidefault`.

**Parse (read-time, fail-closed, right-to-left) — updated from my round-01 for dots:**
split on `-`. 1 token → bare family (legacy). 2 tokens, token[1] all-digits → legacy instance
ID `family-N`. 3–4 tokens → composite: if the last token is all-digits it is the instance
(effort vocab can never be all-digits, so no ambiguity); the next token from the right must be
a vocabulary member (hard error otherwise — an all-digit *model* like a hypothetical `530`
still parses correctly); the single remaining middle token is the model, dots allowed. Anything
else → parse error; tooling never guesses. Family matching uses the explicit roster→family
mapping (§4), not prefix heuristics — `antigravity-1`↔`agy` already breaks prefixing.

**The five names under the final grammar** (display names; canonical IDs unchanged):

| Roster ID (canonical) | Family | `model_label` (source of truth) | Effort token | Display name |
|---|---|---|---|---|
| `claude-1` | `claude` | `Opus 4.8` → `opus4.8` | `max` | `claude-opus4.8-max` |
| `codex-1` | `codex` | `GPT-5.5` → `gpt5.5` | `xhigh` | `codex-gpt5.5-xhigh` |
| `hermes-1` | `hermes` | `GLM 5.2` → `glm5.2` | `high` | `hermes-glm5.2-high` |
| `antigravity-1` | `agy` | `Gemini 3.5 Flash (High)` → strip tier → `gemini3.5flash` | `high` (from tier) | `agy-gemini3.5flash-high` |
| `kimi-1` | `kimi` | `K3` → `k3` | `max` | `kimi-k3-max` |

agy note, reconciling D's "agy=cli-default": the authoritative `reasoning` stays `cli-default`
(no per-invocation flag exists); the *display* surfaces the label tier `high` as the effort
token (all four round-01s converged on this). Render rule: prefer the recorded tier for agy,
else the reasoning field.

**Code consistency (the brief's explicit demand):**
- `internal/protocol/roster.go:17` (`[a-z0-9][a-z0-9-]*`): **unchanged.** Composites never
  occupy the §2 ID cell; they live in a new `Display name` annotation column. Roster IDs stay
  in the strict charset forever.
- `internal/consensus/consensus.go:90` (`[A-Za-z0-9._-]+`): **unchanged.** It already tolerates
  the full display charset; membership checks (`:145-146`) stay exact-match on roster IDs.
- Documented rule: **canonical IDs live in the intersection charset (`[a-z0-9-]`); display
  names live in the user charset minus path-unsafe sequences.** No parser needs widening.

## 3. Late update 2 — component D: `fast` as a separate axis that never downgrades effort

**Code reality (traced, so the fix lands in the right layer):** in Go, `Spec.Speed` is a pure
pass-through label. `runtime.go:343-349` applies `[defaults].speed` to every spec;
`runner.go:836` renders it in the TUI header; `runcontrol.go:166` logs it in events. **No Go
code derives model or effort from speed** — `headless_args` hard-code model/effort and
`buildAgentInvocation` (`runner.go:1066`) substitutes only `{root}`/`{prompt}`. The downgrade
conflation lives entirely in config/skill conventions, and it is real: this deck's
`meta/headless-agents.local.json:28-41` pins claude-1 `profiles.fast = {model: "sonnet",
thinking: "low"}`, and `:120-133` pins agy `profiles.fast = {model: "Gemini 3.5 Flash (Low)"}`
— exactly the legacy "fast = weaker model + lower thinking" the user has now banned. Go does
not parse any `profiles.*` table (verified: no such key in `internal/config`), so those blocks
govern only the skill flow.

**Settled design for D:**

1. **Semantics:** `speed` (`fast|deep`) selects output speed only. `fast` = *same model, same
   effort, faster output* (Claude Code `/fast` semantics). It must never alter `model`,
   `reasoning`, or any materialized launch flag for them. The confirmed efforts
   (claude=max, codex=xhigh, hermes=high, agy=cli-default, kimi=max) are invariant under speed.
2. **Not in the name.** The composite encodes agent-model-effort only. Speed is a launch
   property, switchable to `deep` for heavy ideas.
3. **Config/template changes:** `centralDefaultTemplate` (`runtime.go:286`) changes
   `speed = "deep"` → `speed = "fast"` for new machine files; `roster init` writes
   `speed = "fast"` as the new-roster default (both scopes), respecting an explicit existing
   per-agent `speed` on re-run.
4. **Kill the downgrade blocks:** session-scope reinit surgically rewrites the deck JSON so
   every `profiles.{fast,deep,review}` carries the *same* model+thinking as the entry's top
   level — or, preferred, **drops the `profiles` table entirely** where no distinct fast
   mechanism exists (single source of truth: top-level `model`/`thinking`). The same rule
   applies to any `profiles.*` convention in central TOML comments/skill text.
5. **Per-agent fast mechanism — honesty rule:** the skill's table lists, per CLI, the
   fast-output mechanism *where verified* (claude interactive `/fast` per the brief); for
   headless agents with no verified fast flag, the entry is `no-op (documented)` — never a
   substitute weaker model. Implementation must verify each flag before writing it into a
   spec; unverified = no-op. Speed's only Go-side effects remain display/logging until a real
   per-CLI mechanism lands.
6. **Guard test:** table-driven test asserting that a `speed="fast"` roster keeps each agent's
   pinned `model`/`reasoning` through `LoadAgentSpecs` and launch assembly — the regression
   tripwire for any future code path that tries to map speed → model/effort.

## 4. Two-namespace schism — final position: FIX it here, minimally

The schism (my round-01 §0, hermes-1's summary, both verified): the driver keys everything on
spec IDs (`claude`, `agy`) while the deck — §2, `00-prompt.md`, `headless-agents.local.json`,
and every hand-written artifact including this idea's — uses roster IDs (`claude-1`,
`antigravity-1`). `selectedAgents` (`runner.go:331-339`) exact-matches, so a driver run over
this very idea selects zero agents. Two namespaces, no resolver.

**Why fix now:** (i) `roster init` must write the roster-ID↔family mapping *anyway* (all four
round-01s agree); the resolver consumes that same mapping — marginal cost is one function plus
call-site rewiring. (ii) The fix is zero-rename for deck-native rosters (artifacts are already
roster-ID-keyed) and backward-compatible for legacy spec-ID participants (resolution rule 1).
(iii) Leaving it means the composite feature ships on a driver that cannot run the deck's real
roster — the feature would be display-only for the skill flow and dead for the driver.

**The fix — one identity (roster/participant ID) + adapter templates:**

1. **`Spec` split (codex-1):** `ID` (stable identity), `AdapterID` (launch/discovery family).
   Vendor branches move to `AdapterID`: `cleanParticipantEnv`, `isolatedAgentHome`, Hermes env
   in `runner.go`; `defaultBuiltinSpecs`/`mergeACPCatalog` stay family-keyed.
2. **Explicit mapping, written by `roster init`** (the single authority answer to codex-1's
   concern 1): deck `agents.toml` gains `[roster.<roster-id>] adapter = "<family>"`; the §2
   table's `cli` annotation and the local JSON's per-entry `cli` field are *renders/exports* of
   it. Go's loader authority stays TOML (central < deck < local TOML < env,
   `runtime.go:127-144`); the JSON is generated local state with `schemaVersion` + source
   fingerprint, never parsed by Go, written surgically (rostered IDs only, only if it exists).
3. **Fail-closed resolver** (new `internal/agents/resolve.go`): participant string → spec by
   (1) exact spec-ID match (legacy decks), (2) explicit `[roster.*]` mapping, (3) else hard
   error naming the unresolvable participant. No prefix heuristic (`antigravity-1`↔`agy`
   breaks it). Ambiguity = error, never a guess.
4. **Artifact identity = the participant string.** `runAgent` (`runner.go:349`) writes
   `<participant>.md`; `validation.go:42-72` expects `agent: <participant>`; runstate, review
   snapshots (`reviewsnapshot.go:67`), `_index.md` (`round_index.go`), TUI tabs
   (`protosnap.go:253-262`) key on the participant ID. Digest (`digest.go:48-53`) already keys
   on the raw participant string — zero change. Consensus membership (`consensus.go:145`) —
   zero change. `preflight.rosterEntry.RosterID` is already the roster ID (codex-1).
5. **Migration: forward-only** per `meta/roster-update_2026-06-19.md`. No historical renames,
   no alias subsystem, no `git mv`. Open ideas keep their locked participants; the resolver's
   rule 1 keeps spec-ID-era ideas runnable. If a future deck insists on composite-as-ID, §2's
   join/leave mechanics (my round-01 §5.3) remain the only path — manual, confirmed, guarded.

**Go touch points (complete list):** `internal/agents/discover.go` (Spec split,
`AutonomousWrite`, built-in fixes), new `internal/agents/resolve.go`,
`internal/runner/runner.go` (`selectedAgents` → resolve; `runAgent` artifact path; vendor env
branches → `AdapterID`), `internal/runner/validation.go` (expected `agent:` value),
`internal/config/runtime.go` (`[roster.*]` mapping, `model_label`, `autonomous_write`,
template `speed`), `internal/config/roster.go` + `internal/protocol/roster.go` (§2 writer,
Display-name column, shared allowlist with the drift guard), `internal/app/app.go` (`roster`
verb group), `internal/app/preflight.go` (probe reuse + non-gating drift WARN),
`internal/runmanifest/manifest.go` (profile snapshots), `internal/driver/digest.go` +
`internal/tui/roundsummary.go` (display rendering + width). Tests: sanitization/parse/collision
(dot cases incl. `..` rejection, edge dots, all-digit model, legacy `family-N`), resolver
(exact/mapped/unknown/ambiguous), artifact-identity stability across a mid-run model change,
open-idea guard (drop→inactive, `--force-drop` records), autonomous mapping per CLI incl.
kimi's `-p` exclusivity, speed-invariance guard. `TestEmbeddedDefaultMatchesLiveDeck` stays
green via the shared allowlist.

## 5. Component B, settled

**Surface** (distinct verb group; `parley agents list|verify` stay runtime-discovery):

```text
parley roster init [--dir DIR] [--scope session|machine] [--agent <family>]...
                   [--from FILE] [--dry-run] [--yes] [--force-drop] [--json]
parley roster show [--scope session|machine] [--json]
parley roster diff [--scope session|machine] [--json]
```

- **No `--reinit` flag** (codex-1, kimi-1): the command is the re-init; existing values are the
  picker's defaults. claude-1's "--reinit re-runs it" = run `roster init` again.
- **Default scope:** `session` inside a deck, `machine` outside (hermes-1), with `--dry-run`
  printing exactly which files change per scope.
- **Session scope writes:** deck `agents.toml` (`[agents.<family>]` model/`model_label`/
  reasoning/speed/`autonomous_write` + `[roster.<id>]` mapping); targeted §2 table update
  (adds the `Display name` column, allowlist-shared writer); surgical local-JSON update
  (only rostered IDs, only if present, fingerprinted).
- **Machine scope writes:** `~/.parley/agents.toml` only — managed-block markers + first-time
  backup (codex-1); never touches a deck, never copies deck values up.
- **Reconciliation:** `parley init` stays scaffold+seed and prints `parley roster init --scope
  session` as the required next step (skill chains it; optional `--with-roster` convenience
  flag runs the *same* internal `RosterInit` service — claude-1's "one engine" point without
  overloading the verb). §9.0 preflight stays liveness-only, gains a non-gating drift WARN
  ("roster drift — run `parley roster init`"). Open-idea guard per §1 (refuse →
  inactive-retention → `--force-drop` records `excluded:`/`reincluded:`); model/effort/speed/
  display changes never trip it; open runs use their manifest snapshot.
- **Bootstrap honesty (kept from my round-01):** codex is `cli-default` in both deck layers
  today, so the truthful current display is `codex-clidefault-clidefault`; `codex-gpt5.5-xhigh`
  becomes true only when reinit pins it. Names follow config truth, never wishes. Also noted:
  kimi has no built-in spec (ACP catalog only) and is absent from this deck's §2 table and
  local JSON — reinit must handle config-only families (`applyFile`, `runtime.go:361-381`)
  with a documented catalog fallback.

## 6. Component C, settled

**Field (merge of all three proposals):** `AutonomousWrite{Mode string, Args []string,
Scope string}` on `Spec` — a struct (hermes-1/codex-1) so the argv stops being buried, but
minimal (kimi-1): `Scope` is `"workspace"` or the bit is unset; codex-1's `ScopeEnforcement`
taxonomy is dropped, its honesty rule kept (unverified confinement ⇒ bit unset, WARN).
Bookkeeping: `agentOverride` (`runtime.go:88`), `withBuiltinSources` (`discover.go:224`),
`centralDefaultTemplate`, an `AUTO` column in `PrintRuntimeMatrix`, static mapping check at
`roster init`, preflight WARN (not gate), opt-in live sentinel via `parley agents verify
--full` (never default — probe cost).

**Fix the built-ins (verified against `defaultBuiltinSpecs`):**

| CLI | Today (file:line) | Change to (workspace-scoped) |
|---|---|---|
| claude | `--permission-mode acceptEdits` (`discover.go:134`) | `--permission-mode bypassPermissions` + keep `--add-dir {root}` |
| codex | `approval_policy="on-failure"` (`discover.go:113`,`:119`) | `-c approval_policy="never"` + keep `--sandbox workspace-write` — no full-filesystem escalation |
| hermes | `--accept-hooks`, no `--yolo` (`discover.go:203`) | `--yolo --accept-hooks` |
| agy | `--dangerously-skip-permissions` (`discover.go:156`) | already autonomous; keep `--add-dir {root}` |
| kimi | no built-in spec | add spec/catalog entry: plain `-p` IS the autonomous mode; **never** `--yolo`/`--auto` with `-p` (mutually exclusive) |

**Skill wording:** adopt codex-1's normative block (round-01 §C quote) nearly verbatim — it is
the tightest statement of "required, workspace-confined, fail-closed, redaction preserved" —
with the mapping table above and one clause from hermes-1: the skill points at the spec field
as source of truth so a vendor flag change is a config edit, not a skill revision.
**Safety invariants (all four agree, restated once):** bypass flags live in per-project config,
never vendor-global; every mapping carries its workspace constraint (`--add-dir`/`--cd`/cwd);
secret redaction is orthogonal to permission mode and stays on; changing claude's built-in is a
real posture change → the scoping flags land in the same commit + CHANGELOG note (my round-01
concern 6).

## 7. The one coherent design (for consensus.md → FINAL.md)

1. **Identity:** roster ID (`claude-1`) is the single identity for artifacts, signoffs,
   frontmatter, runstate, TUI, digests. Spec family becomes `AdapterID`, a launch template.
2. **Display:** composite `family-model-effort` (dots allowed inside model; grammar in §2),
   *derived at render* from `model`/`model_label`/`reasoning`; shown in §2 (new column), TUI,
   digests, run headers; never a key, never stored as truth.
3. **Freeze:** per-idea `participant-profiles` snapshot in `00-prompt.md` + exact profiles in
   `runmanifest.Manifest`; resume prefers the snapshot (claude-1's freeze, codex-1's mechanism).
4. **Schism:** fixed minimally per §4 — explicit `[roster.*]` mapping + fail-closed resolver +
   participant-string artifact identity; forward-only migration, zero renames.
5. **Command:** `parley roster init|show|diff`; session=deck layer (TOML authority + surgical
   JSON export + §2 writer), machine=`~/.parley/agents.toml` (managed blocks); default scope
   session-in-deck; open-idea guard with inactive-retention and `--force-drop` recording.
6. **Autonomous:** `AutonomousWrite{Mode,Args,Scope}` first-class; built-ins corrected per §6
   table; static verification + WARN; skill normative block + mapping.
7. **Speed:** `fast` default for new rosters, separate axis, never downgrades model/effort,
   never in the name; `profiles.*` downgrade blocks rewritten or dropped; template flips to
   `speed = "fast"`; per-CLI fast mechanisms verified-or-noop; speed-invariance guard test.

## 8. Residual risks / open questions for round-03

- Per-CLI headless fast mechanisms are unverified (except claude interactive `/fast` per the
  brief); implementation must verify before writing any flag — else documented no-op.
- `model_label` backfill: codex today is `cli-default` — honest display until reinit pins it;
  the picker must not invent a label.
- The §2 writer and the drift guard must share one allowlist definition (hermes-1's risk); add
  a test that round-trips the writer against `TestEmbeddedDefaultMatchesLiveDeck`'s zones.
- TUI width (`roundsummary.go:84-90`) must land in the same release as display rendering —
  first visible surface of the change.
- kimi enters via ACP catalog/user config only; its `roster init` model/effort discovery needs
  the documented static-catalog fallback (hermes-1 concern 5, my concern 5).
