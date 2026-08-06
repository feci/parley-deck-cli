---
agent: hermes-1
idea: roster-operations-standard
round: 1
date: 2026-08-06
---

## Summary

The roster has two unrelated surfaces that answer two different questions, and neither is the canonical "show me the roster" the user wants. `parley roster show` answers "who is in this deck's §2 roster and what model/effort will they run?" (4 rows, roster IDs); `parley agents list` answers "what agent CLIs are installed on this machine and how are they configured?" (13 rows, family IDs). They must stay two commands, but the user-facing roster question needs ONE table, and that table must show the EFFECTIVE launched model, not the declared `spec.Model` that config layers set but the launcher ignores. The opencode inconsistency is a §2-roster vs `[roster.*]`-mapping drift: opencode-1 was hand-added to `~/.parley/agents.toml [roster.opencode-1]` but never added to the §2 table in COOPERATION.md, so `roster show` (driven by §2) omits it while `agents list`/`preflight` (driven by discovery) include it — and sessions that read §2 see a different roster than sessions that read the `[roster.*]` map. Update and sync need two command surfaces (`roster set` for per-agent pinning, `roster sync` for §2↔config reconciliation), and the skill must describe the human steps the CLI cannot enforce (confirming a roster change, choosing a model) while the CLI enforces the data contract.

## Proposed approach

### 1. The canonical roster table — column by column

One table is THE roster answer. It is produced by `parley roster show` and printed verbatim by the skill (the skill must not describe a second format). Proposed columns, in order:

| Column | Source | Why it is load-bearing |
|---|---|---|
| `AGENT` | §2 roster ID (`claude-1`) from `protocol.ReadRosterIDs` | The stable identity used in artifact paths, signoffs, frontmatter. Not the family. |
| `FAMILY` | `[roster.<id>] adapter` mapping → `spec.ID` (resolved via `config.LoadRosterAdapters` + `agents.ResolveParticipant`) | The launch/discovery dispatch key. Distinct from AGENT since composite naming. |
| `MODEL` | **Effective** model: the model token actually present in `spec.HeadlessArgs`, parsed out; fallback to `spec.Model` only when no model flag is in argv and the CLI reads its own config (kimi, codex) | This is the fix for the declared-vs-effective split. Today `roster show` prints `spec.Model` (the declared value config layers set), but the launcher runs `HeadlessArgs` verbatim and synthesizes no model flag. For claude, `spec.Model`=`claude-opus-5[1m]` while `HeadlessArgs` launches `claude-opus-4-8[1m]` — the table lies. The effective model is the one in the launch argv; when the argv has no model flag (kimi `-p`, codex `exec`), the CLI's own config file is the effective source and the column should say `<from-cli-config>` or the declared value with a flag. |
| `MODEL-FAMILY` | Derived from the effective model token via a maintained mapping table in the CLI (see §2 below) | User asked for it; must be derived, not hand-maintained per deck. |
| `MODEL-COMPANY` | Derived from model-family via the same mapping | User asked for it; same derivation rule. |
| `EFFORT` | `spec.Reasoning` (already shown) | The reasoning/effort level. |
| `SPEED` | `spec.Speed` (already shown) | Separate axis from effort (fast output at same model+effort). |
| `AUTO` | `spec.AutonomousEffective()` (already shown, fail-closed) | Whether the launch argv actually enables the declared auto-write mode. |
| `INSTALLED` | `discovery.Found` | Whether the CLI binary is on PATH. A roster row whose family is not installed is a broken deck. |
| `EFFECTIVE-MATCH` | boolean: does the model in `HeadlessArgs` equal `spec.Model`? | The divergence flag the constraints demand. When false, the row is annotated `⚠ declared=<spec.Model> but launch uses=<argv model>`. This is the `AUTO` honesty rule applied to `model`. |

`DISPLAY-NAME` (the composite `family_model_effort`) is dropped from the table. It is a derived rendering of FAMILY+MODEL+EFFORT and showing it alongside those three columns is redundant; it remains available as `--format wide` or in JSON for contexts that want the single-token form. The skill's §2 roster table in COOPERATION.md keeps its current shape (AGENT ID, workspace dir, role) — that is human-edited prose and is NOT the same table; the canonical roster table is a CLI-derived runtime view, and the skill must say "run `parley roster show` for the current roster" rather than maintaining a second copy.

**Why these columns and not more:** `AGENT`, `FAMILY`, `MODEL`, `EFFORT`, `SPEED`, `AUTO` are the user's named set plus the existing working columns. `MODEL-FAMILY` and `MODEL-COMPANY` are added per the explicit request. `INSTALLED` and `EFFECTIVE-MATCH` are the two integrity signals the constraints require (effective-vs-declared flag, installed state). Anything beyond this (timeout, sandbox, approval, backend, home isolation) is `agents list`'s job — it is the diagnostic surface for "why did this agent fail to launch", not the roster answer.

### 2. `model family` and `model company` — derived how

A single mapping table in the CLI source (e.g. `internal/agents/models.go`), keyed by a prefix or token of the effective model id:

```
claude-opus-*    -> family "opus",   company "Anthropic"
claude-sonnet-*  -> family "sonnet", company "Anthropic"
gpt-*            -> family "gpt",    company "OpenAI"
o3-*|o4-*        -> family "o-series", company "OpenAI"
glm-*            -> family "glm",    company "Zhipu AI"
grok-*           -> family "grok",   company "xAI"
gemini-*         -> family "gemini", company "Google"
kimi-code/*      -> family "k3",     company "Moonshot AI"
```

For **gateway-routed models** (hermes `glm-5p2` via a custom LitLLM gateway; opencode `litellm/xai/grok-4.5` via the user's LiteLLM gateway), the derivation strips the `litellm/` prefix and resolves the inner model token. So `litellm/xai/grok-4.5` → family "grok", company "xAI". The gateway itself is not a company; it is a routing layer. The mapping table is maintained in the CLI and updated when new model families ship; it is NOT per-deck config (the constraint says derived, not hand-maintained per deck). An unrecognized model token prints `family=?`, `company=?` — fail-open for display, but the `?` is a visible signal to update the mapping.

**Trade-off:** a hardcoded mapping in the CLI must be updated when a new model family is rostered. The alternative (a config file) just moves the maintenance and adds a place to forget. Given the set of model companies is small and changes rarely, CLI-source is the right call; it ships with the release and is versioned.

### 3. `parley roster show` and `parley agents list` — two commands, two purposes

They stay two commands. The split is correct; the problem is that neither is documented as the roster answer.

- **`parley roster show`** = "who is in this deck's roster and what will they run?" Rows = §2 roster IDs. Columns = the canonical table above. This is THE answer to "what is the current agent roster?". Add it to `parley --help` (it is registered and dispatches today but is absent from the Usage/Commands listing — `PRIMARY`: confirmed by reading `internal/app/app.go:100-101` dispatch vs `printUsage` at `:111-154` which omits the `roster` line).
- **`parley agents list`** = "what agent CLIs are installed and how are they configured?" Rows = all discovered family IDs (including uninstalled ACP backends). Columns = the runtime matrix (installed, version, launch, sandbox, approval, model, timeout, home, backend). This is the diagnostic surface for debugging launch failures and for seeing the full adapter catalog. It is NOT the roster.

The skill must say this explicitly: "`parley roster show` is the roster; `parley agents list` is the installed-CLI matrix. They answer different questions and have different row sets by design."

### 4. Update — `parley roster set`

A new `roster set` subcommand for per-agent pinning. Today the only update path is editing TOML by hand or running `roster init` (which only writes the `[roster.<id>] adapter` mapping, not model/effort). Proposed:

```
parley roster set <roster-id> --model <model> [--effort <effort>] [--scope session|machine] [--dry-run] [--yes]
parley roster set <roster-id> --family <family> [--scope session|machine] [--dry-run] [--yes]
```

- `--scope session` writes to `parley-deck/agents.toml` (deck override).
- `--scope machine` writes to `~/.parley/agents.toml` (global).
- `--dry-run` prints the diff; `--yes` writes without confirmation (same safety shape as `roster init`).

**Safety properties (same as existing `roster init`):**
- Atomic temp+rename write (`fsutil.WriteFileAtomic`).
- Structural validation of the candidate before write (`config.ValidateAgentsConfigBytes`).
- Idempotency judged against the TARGET file, not the layered stack.
- Machine scope never proposes/writes a deck-only family (`config.MachineFamilyCatalog`).
- Refuses to write a model the adapter cannot accept — but since the CLI launches `HeadlessArgs` verbatim and does NOT synthesize model flags, `roster set --model` must ALSO rewrite `headless_args` to embed the new model (for adapters that bake the model into argv: claude, hermes, agy, opencode). For adapters that read their own config (kimi, codex), it only sets `spec.Model` for display and leaves `headless_args` alone. The command must know which family is which.

**The HeadlessArgs-rewrite requirement is the load-bearing part.** This is the root-cause fix for the claude declared-vs-effective split: today you can set `model = "claude-opus-5[1m]"` in config and the launch still runs `claude-opus-4-8[1m]` because `HeadlessArgs` is a separate field that config overrides wholesale but nobody rewrites the embedded `--model` token. `roster set --model` must update BOTH the `model` field AND the `headless_args` array, so the declared and effective models can never diverge by construction. This is the same class of fix as the `AutonomousEffective` fail-closed check — the data contract should make the divergence impossible, not just detect it.

### 5. Sync — `parley roster sync`

A new `roster sync` subcommand: reconcile the §2 roster table (COOPERATION.md) with the `[roster.*]` config mapping. This is the operation that is "re-invented every session" today.

```
parley roster sync [--scope session|machine] [--dry-run] [--yes]
```

What it does:
1. Read §2 roster IDs from COOPERATION.md (`protocol.ReadRosterIDs`).
2. Read `[roster.*]` mappings from the layered config (`config.LoadRosterAdapters`).
3. Compute the diff: IDs in §2 but not in `[roster.*]` (need a mapping); IDs in `[roster.*]` but not in §2 (orphaned mapping — either add to §2 or drop the mapping).
4. `--dry-run` prints the diff. Without `--yes`, asks for confirmation. With `--yes`:
   - For IDs in §2 missing a mapping: propose a family via `proposeFamily` and write `[roster.<id>] adapter = "<family>"` to the target scope (same as `roster init` does today).
   - For IDs in `[roster.*]` missing from §2: print a warning and DO NOT auto-add them to §2 (adding a roster member is a human/protocol decision per §2 and §14.3). This is the opencode case — see §6 below.

**The sync direction is config ← §2 for mappings, and §2 ← human for membership.** Sync never adds a row to §2 automatically; it only writes the `[roster.*]` adapter mappings for IDs already in §2. Adding a roster member is a §2 table edit that the human does (or a `meta/roster-update_<date>.md` idea per §2). Sync makes the mapping layer match §2; it does not make §2 match the mapping layer.

**Safety:** same atomic-write + validation as `roster init`. Sync is idempotent: a deck where §2 and `[roster.*]` agree is a no-op.

### 6. The opencode inconsistency — explanation and fix

**Explanation (`PRIMARY`, verified by reading the live config and §2):**

- The §2 roster table in `parley-deck/COOPERATION.md` (lines 107-110) lists four agents: `claude-1`, `codex-1`, `hermes-1`, `kimi-1`. It does NOT list `opencode-1`. (`PRIMARY`: read the file, `grep -n "opencode\|hermes-1\|kimi-1\|claude-1\|codex-1"` on COOPERATION.md returns only the four rows.)
- The `[roster.*]` mapping in `~/.parley/agents.toml` lists FIVE agents: `claude-1`, `codex-1`, `hermes-1`, `kimi-1`, AND `opencode-1` (hand-added 2026-08-05). (`PRIMARY`: `cat ~/.parley/agents.toml` shows `[roster.opencode-1] adapter = "opencode"`.)
- `parley roster show` reads §2 (`protocol.ReadRosterIDs`) → sees 4 IDs → shows 4 rows. opencode is absent.
- `parley agents list` reads discovery (`config.LoadAgentSpecs` → all built-in + configured specs) → sees opencode (installed + configured) → shows it. It does not consult §2 at all.
- `parley preflight` reads `selectedParticipants(discovered)` → filters `discovered` by `agent.Found && agent.ID != "gemini"` → sees opencode (family ID, installed) → includes it. It also includes `agy` (still a built-in spec at `discover.go:238`, still installed at `~/.local/bin/agy`) even though agy was "removed from the roster" — the removal was a `[agents.agy]` config-block deletion, not a built-in-spec deletion. (`PRIMARY`: `parley preflight --no-ping` output lists `agy ... yes` and `opencode ... yes` alongside the four §2 families.)

**Why sessions are inconsistent:** sessions/flows that read §2 (e.g. `parley roster show`, the §9.0 readiness check that locks quorum from `participants:` which comes from §2) see a 4-agent roster without opencode. Sessions/flows that read discovery or the `[roster.*]` map (e.g. `parley agents list`, `parley preflight` when run standalone, or a facilitator manually reading `~/.parley/agents.toml`) see opencode as available. The two namespaces — §2 (human-edited protocol state) and `[roster.*]` (machine-edited config state) — drifted because opencode was added to one and not the other, and nothing reconciles them. This is exactly the "two different meanings of roster coexist with no defined relationship" problem in the brief.

**Fix (three parts):**

1. **`roster sync` surfaces the drift.** Running `parley roster sync --dry-run` would print: "opencode-1 is in `[roster.*]` but not in §2 — add it to §2 or drop the mapping." Today there is no command that shows this diff.
2. **`roster set` is the single update path.** Adding opencode should have been `parley roster set opencode-1 --family opencode --scope machine` (writes the `[roster.*]` mapping) PLUS a §2 table edit (the human step). The command should refuse to write a `[roster.<id>]` mapping for an ID that is not in §2 without `--allow-unmapped`, so the drift cannot happen silently.
3. **The built-in `agy` spec should be gated by §2 or removed.** `agy` was "removed from the roster" by deleting the `[agents.agy]` central-config block, but the built-in spec at `discover.go:238` still exists, so `agents list` and `preflight` still show it as an installed agent. This is a separate bug: "removed from config" ≠ "removed from the built-in catalog". Either the built-in `agy` spec is deleted (a CLI release), or `agents list`/`preflight` filter by §2 membership (which would make them roster-aware, a bigger change). I lean toward filtering `agents list` by §2 membership for the roster rows, with an `--all` flag to show the full adapter catalog — but this is an open question for round 2.

### 7. Is there a real session-scoped roster store?

No. `--scope session` is nominal: it writes to `parley-deck/agents.toml` (a deck-level file), not a per-session store. (`PRIMARY`: `internal/app/roster.go:383-389` `rosterTargetPath` maps `session` → `parley-deck/agents.toml`, `machine` → `~/.parley/agents.toml`. There is no session-id-keyed store.) "Session" means "this deck/project", not "this OS process". This is fine — there is no need for a per-session roster store; the roster is project state. The naming is misleading though: `--scope session` should arguably be `--scope deck` or `--scope project`, but renaming it is a breaking change. The skill should document that "session" = "this deck's `agents.toml`".

### 8. Skill vs CLI boundary

**The CLI enforces (data contract, machine-checked):**
- The column set and its sources (effective model from HeadlessArgs, AUTO from `AutonomousEffective`, INSTALLED from discovery).
- The `roster set` / `roster sync` safety properties (atomic writes, validation, idempotency, machine-scope family filtering).
- The `EFFECTIVE-MATCH` divergence flag.
- `roster set --model` rewriting `headless_args` so declared and effective cannot diverge.

**The skill describes (human steps, not machine-enforceable):**
- "Run `parley roster show` to see the roster; `parley agents list` to see installed CLIs. They are different commands for different questions."
- "To add an agent: (1) add a row to the §2 table in COOPERATION.md, (2) run `parley roster set <id> --family <family>`, (3) run `parley roster sync` to verify. Do not edit `[roster.*]` by hand."
- "To change an agent's model: run `parley roster set <id> --model <model> --scope session|machine`. The command updates both the declared model and the launch argv."
- "To sync the session roster from the global one: run `parley roster sync --scope session`. This makes the deck's `[roster.*]` mappings match §2; it does not copy model/effort from `~/.parley/agents.toml` (use `roster set` per agent for that, or accept the central default)."
- The §2 table edit is a human step because adding a roster member is a protocol decision (§2, §14.3), not a config write.

The boundary: the CLI owns the data contract and the mechanical safety; the skill owns the workflow and the human-in-the-loop steps. The skill must not describe a second table format or a second update path — it points to the CLI commands.

## Concerns / open questions

1. **`roster set --model` rewriting `headless_args` is family-specific.** The command needs to know, per family, whether the model is baked into `HeadlessArgs` (claude, hermes, agy, opencode) or read from the CLI's own config (kimi, codex). This is a per-adapter property that is not currently encoded in the `Spec`. It could be a `ModelInHeadlessArgs bool` field or a per-family rewrite function. This is the most complex part of the proposal and needs a concrete design in round 2. (PRIMARY: the split is visible in `discover.go` — claude/hermes/agy/opencode have `--model`/`-m` in `HeadlessArgs`; kimi/codex do not.)

2. **Should `agents list` become §2-aware?** Today it shows the full adapter catalog (13 rows including uninstalled ACP backends). If it filtered to §2 members by default with `--all` for the full catalog, the opencode/agy drift would be less visible — but it would also hide the diagnostic surface for "is this CLI installed?". Open question: is `agents list` a roster command or a diagnostic command? I argue diagnostic; it should stay as-is and the roster answer should come only from `roster show`.

3. **The built-in `agy` spec is still present.** "Removed from the roster" by config deletion did not remove the built-in spec. Is the right fix to delete the built-in `agy` spec (a CLI release, but then anyone who wants agy back needs a config-only definition), or to make `agents list`/`preflight` filter by §2? This is a separate idea or a decision for this idea's scope.

4. **`model family` / `model company` for unknown models.** A hardcoded mapping will have gaps (e.g. a new model id before the mapping is updated). The `?` fallback is fine for display, but should the CLI warn, or should an unknown family/company be a soft gate for `roster set`? I lean toward display-only `?`; the mapping is a rendering concern, not a correctness gate.

5. **`roster sync` and §14.3.** The protocol says an automated loop must not modify the active roster (§14.2). `roster sync` only writes `[roster.*]` mappings for IDs already in §2 — it does not add §2 rows — so it respects §14.2. But if `roster sync` is ever run by a cron/loop, it must not auto-add mappings either without a human gate. The command is safe for manual use; the skill should say "run `roster sync` at session start, manually".

6. **Backward compatibility for `roster show` column changes.** Adding `MODEL-FAMILY`, `MODEL-COMPANY`, `INSTALLED`, `EFFECTIVE-MATCH` changes the table width. JSON output (`--json`) is the stable machine-readable surface; the text table is for humans and can change. Existing decks keep working because the data sources are unchanged — only the rendering changes. The `--json` shape should add the new fields without removing existing ones.

## Risks

1. **`roster set --model` HeadlessArgs rewrite could get the family-specific logic wrong.** If it rewrites `headless_args` for a family that reads its model from its own config (kimi, codex), it would inject a `--model` flag the CLI may reject or ignore. The per-family property must be correct and tested. Mitigation: a per-family `ModelArgStyle` enum (`embedded` | `cli-config` | `none`) in the spec, with tests for each.

2. **The model-family/company mapping will go stale.** New model families ship faster than CLI releases. The `?` fallback prevents breakage but the table will show `family=?` for unrecognized models, which looks broken. Mitigation: the mapping is CLI-source and ships per release; an unknown model is a signal to update the CLI, not a deck-config problem.

3. **Two commands remaining could still confuse users.** Even with documentation, "roster show" vs "agents list" is two commands where one might be expected. Mitigation: `roster show` is the default; `agents list` is documented as diagnostic-only and could be renamed `agents matrix` in a future release (out of scope here).

4. **`roster sync` could mask the §2-edit requirement.** If users run `roster sync` and see "opencode-1 in config but not in §2", they might expect sync to add it to §2. Sync must NOT do that (§14.2, §2 human decision), but the refusal could feel like a broken command. Mitigation: the dry-run output must be explicit — "opencode-1 is mapped in config but not listed in §2. To add it: edit the §2 table in COOPERATION.md, then re-run sync. To drop it: remove `[roster.opencode-1]` from agents.toml."

5. **The `EFFECTIVE-MATCH` column requires parsing the model token out of `HeadlessArgs`.** This is a per-adapter parse (claude `--model <x>`, opencode `-m <x>`, hermes `--model <x>`). A regex or positional parse could break if a vendor changes its flag syntax. Mitigation: the parse is best-effort; on parse failure the column shows `?` and falls back to the declared model with a `⚠` flag, same as the `AUTO` fail-closed pattern.
