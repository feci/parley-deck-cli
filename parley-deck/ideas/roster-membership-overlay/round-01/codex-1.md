---
agent: codex-1
idea: roster-membership-overlay
round: 1
date: 2026-08-19
---

## Summary

**PRIMARY — recommendation:** Build an explicit, versioned, opt-in membership overlay with both `add` and `remove` operations, while preserving every existing unmodified deck as a full replacement and preserving unmarked legacy §2 tables as legacy full replacements. The current all-or-nothing resolver cannot represent “inherit future machine membership, except for these named differences”; a full copied roster can represent today’s set but not that inheritance policy.

**PRIMARY — scope guard:** Do not reinterpret existing `[roster.*]` blocks as deltas and do not mass-convert the fleet. A current scan found 38 decks: 37 resolve from deck blocks and one resolves from the machine roster; 36 of those 37 full replacements omit only `zcode-1` relative to the current six-agent machine roster. An automatic exact-set conversion would turn an ambiguous historical omission into 36 durable removal tombstones, while an automatic base adoption would add an agent to 36 decks. Neither inference is safe.

**PRIMARY — do-nothing assessment:** NO CHANGE is viable for this deck’s zero-local-declaration case, which already resolves all six machine members as `inherited-roster`. NO CHANGE is not sufficient for a deck that needs either a deck-only specialist or a deliberate machine-member exclusion while continuing to inherit later machine additions.

## Verification record

### Commands and direct observations

**PRIMARY — tool and live roster:** I ran `command -v parley`, `parley version`, `parley roster show`, `parley roster show --scope machine`, and `parley roster show --explain zcode-1`. The binary was `/opt/homebrew/bin/parley`, version `1.45.0`. Both live tables contained six active agents. Every deck-scope row carried `inherited-roster`; the machine-scope rows did not. The relevant explanation was:

```text
zcode-1 — membership from ~/.parley/agents.toml (INHERITED — this deck declares no roster of its own)

FIELD          EFFECTIVE                SET BY
adapter        zcode                    ~/.parley/agents.toml
model          zai/glm-5.3              ~/.zcode config (agent's own, read at launch)
effort         max                      ~/.zcode config (agent's own, read at launch)
speed          deep                     ~/.parley/agents.toml
active         active                   ~/.parley/agents.toml
```

**PRIMARY — live deck configuration:** I ran `cat parley-deck/agents.toml`. It contains launch configuration but no `[roster.*]` block, and its roster comment says membership is intentionally not declared locally. I also read `internal/config/runtime.go` with `sed -n '120,210p'` and `sed -n '379,410p'`. The first range selects membership in this order: non-empty deck blocks, otherwise a readable legacy §2 table, otherwise machine members; the second range applies values from machine, deck, deck-local, and environment layers in low-to-high order.

**PRIMARY — winner-takes-all behavior:** I ran `parley roster show --dir '/Volumes/My Shared Files/AI_WORKSPACE/BYTE'`, inspected BYTE’s `[roster.*]` blocks with `rg`, and ran `parley roster show --dir '/Volumes/My Shared Files/AI_WORKSPACE/BYTE' --explain zcode-1`. BYTE declared five members, the table showed five, and the explanation command returned `zcode-1 is not in this deck's roster` even though the machine roster contains `zcode-1`.

**PRIMARY — isolated three-branch reproduction:** In `/tmp/parley-membership-behavior.bbCBjF`, with an isolated `PARLEY_HOME`, I created a two-member machine roster, a one-member deck roster, and a valid one-row §2 table. I ran these states through the installed CLI:

```text
PARLEY_HOME=... parley roster show --dir /tmp/parley-membership-behavior.bbCBjF
  -> local-a only

mv .../parley-deck/agents.toml .../parley-deck/agents.toml.saved
PARLEY_HOME=... parley roster show --dir /tmp/parley-membership-behavior.bbCBjF
  -> local-a only, STATUS=legacy-roster,unmapped

# I removed only the temporary §2 body row with apply_patch.
PARLEY_HOME=... parley roster show --dir /tmp/parley-membership-behavior.bbCBjF
  -> machine-a and machine-b, both STATUS=inherited-roster
```

**PRIMARY — executable tests:** I ran the following uncached tests in the shared tree without editing it:

```text
go test -count=1 ./internal/protocol/...
ok   parley-deck-cli/internal/protocol  0.260s

go test -count=1 -v ./internal/app -run 'TestDeckMembershipIsTheDeckFileNotTheLayeredUnion|TestRosterlessDeckMarksInheritedRows|TestLegacySection2BeatsTheMachineRoster'
--- PASS: TestLegacySection2BeatsTheMachineRoster
--- PASS: TestDeckMembershipIsTheDeckFileNotTheLayeredUnion
--- PASS: TestRosterlessDeckMarksInheritedRows

go test -count=1 ./internal/config/...
ok   parley-deck-cli/internal/config  0.454s
```

**PRIMARY — drift guard against my own temporary edits:** I copied `go.mod`, `go.sum`, `internal/`, and the live deck into `/tmp/parley-roster-overlay.kFxYN1`, edited only that copy with `apply_patch`, and ran `go test -count=1 ./internal/protocol/...` after each edit. Renaming the exact roster header failed with `appears 0 times, want exactly 1`; duplicating it failed with `appears 2 times, want exactly 1`; adding a line before the header failed at the first normalized divergence. This directly exercises `internal/protocol/drift_test.go:27-135`. The unedited shared tree passed the same test.

**PRIMARY — current fleet count and sync state:** From `/Volumes/My Shared Files/AI_WORKSPACE`, I ran the prior verification’s bounded `find` scan for `*/parley-deck/COOPERATION.md` and checked each file for the `NOT authoritative` marker. The current result was `TOTAL=38 SYNCED=37 NOT_SYNCED=1`; the unsynced deck was `ecb-meeting-2026.05`.

**PRIMARY — current fleet resolution:** I invoked `parley roster show --json` once in each of those 38 deck roots and classified the returned statuses. The result was `deck=37 inherited=1 legacy=0`. Comparing active IDs with `parley roster show --scope machine --json` produced two distinct active sets: one inherited deck equal to the six-agent machine set; one full replacement also equal to it; and 36 full replacements whose only active-set difference was removal of `zcode-1`.

**PRIMARY — current fleet working state:** I ran `git status --porcelain=v1 --untracked-files=normal` in every discovered deck root, without printing filenames from the status payload. At that observation point, 18 deck roots reported local changes, 14 were not inside Git worktrees, and only 6 reported clean. These are deck-root counts, not deduplicated repository counts.

**PRIMARY — historical fleet report inspected directly:** I ran `jq` over `parley-deck/ideas/roster-operations-standard/evidence/migrate-report-2026-08-06.json`. It contains 36 deck entries with `applied=24`, `skipped=9`, `unchanged=3`, and `failed=0`.

**SECONDARY — prior rationale:** `parley-deck/ideas/roster-operations-standard/FINAL.md:16-22,79-82` records the earlier 40-deck measurement and the decision that `agents.toml` owns deck membership while §2 is generated. I did not reconstruct the pre-migration 40-deck state, so the historical “nine rosters / 17 rosterless / 17 retired-agent rows” figures remain secondary here.

**SECONDARY — sibling precedent:** `parley-deck/ideas/protocol-overlay-local-extension/FINAL.md:23-26,121-125` records an extend-only v1 and defers replacement together with its block-addressing, target-hash, tombstone, and registry machinery. I read that final artifact; I did not reopen its underlying user ruling.

## Proposed approach

### Explicit mode, never implicit reinterpretation

**PRIMARY — proposed schema:** Add a separate, versioned membership stanza to the committed `parley-deck/agents.toml`; keep value layering in the existing `[roster.<id>]` blocks.

```toml
[membership]
mode = "overlay-v1"
base = "machine"
add = ["deck-specialist-1"]
remove = ["opencode-1"]

[roster.deck-specialist-1]
adapter = "codex"
```

**PRIMARY — proposed set semantics:** In overlay mode, the visible row universe is the union of machine IDs, `add`, and `remove`; the active quorum candidates are `(machine active IDs ∪ add) − remove`. A removed ID is a tombstone, not a deleted record: it remains visible with `STATE=inactive`, so a later machine re-add does not silently revive it. An ID in both lists is a hard error.

**PRIMARY — proposed authority split:** `add` and `remove` are the only deck membership operations in overlay mode. A `[roster.<id>]` block not named by either operation remains a value override and cannot accidentally create membership. An added ID must resolve an adapter from some value layer. Deck-level `active` is rejected in overlay mode so membership state cannot have two competing authorities; existing replacement mode retains today’s `active` behavior.

**PRIMARY — backward-compatible default:** If `[membership]` is absent, behavior remains byte-for-byte semantic compatibility: any deck `[roster.*]` blocks form a full replacement; a genuinely legacy table remains a legacy full replacement; otherwise the deck inherits the machine roster. Only the explicit `mode = "overlay-v1"` stanza enables set composition.

### Operations and migration

**PRIMARY — proposed operations:** Extend `roster set` with an explicit deck-membership action such as `--membership add|remove|inherit`. `inherit` clears that ID’s local delta. Every add, remove, revival, or conversion remains preview-first and requires both `--yes` and `--confirm-breaking`; ordinary model/effort/value changes do not become membership changes.

**PRIMARY — proposed conversion:** Provide an attended conversion preview with two distinct choices: preserve the current active set by materializing reviewed deltas, or adopt the current machine membership and produce no delta for equal rows. The preview must show current, machine, proposed delta, and effective-after sets and must refuse dirty/unreadable targets in fleet mode. It must never infer that an omission is deliberate.

**PRIMARY — migration consequence from the scan:** A preserve-set conversion would propose `remove = ["zcode-1"]` for 36 current decks. That is evidence for a human decision point, not evidence that 36 intentional exclusions exist. Existing files therefore remain replacements until individually or batch-explicitly converted.

### Legacy §2 handling

**PRIMARY — proposed rule 2:** An unmarked legacy §2 table stays a full authority until attended migration; it never becomes a delta. A modern generated §2 table stops being read as authority and is only a projection.

**PRIMARY — proposed disambiguation:** Add an exact machine-readable projection marker adjacent to the existing table, for example `<!-- parley-roster-table:projection-only/v1 -->`. Legacy files lack it. A protocol render that would add the marker to a non-empty, unmarked table with no modern membership declaration must block and require `roster migrate` or an explicit replacement declaration; it must not silently relabel the table.

**PRIMARY — drift-guard requirement:** Keep the exact §2 section and table-header anchors once each. Adding the marker is a coordinated core/default/live protocol edit, not deck-local prose. The temporary failures above show that editing only the live copy, moving/removing the header, or duplicating it fails closed.

### Visibility

**PRIMARY — proposed table contract:** Keep all eleven columns unchanged and add three documented STATUS terms: `overlay-base`, `overlay-added`, and `overlay-removed`. Every overlay row gets exactly one membership-origin term; removed rows also show `STATE=inactive`. Existing `inherited-roster` and `legacy-roster` meanings remain unchanged.

**PRIMARY — proposed provenance:** `roster show --explain <agent>` reports `mode=overlay-v1`, the machine base source, and the exact deck operation source. Text and JSON golden tests must cover the three new terms, and the skill’s closed vocabulary must be updated in the same release. No computed membership should display plain `ok` without its origin.

## Answers to the six questions

### 1. Operations

**PRIMARY — answer:** Support both add and remove. Removal is a named tombstone with inactive visibility and a breaking-change confirmation, never physical row deletion. Extend-only would leave the concrete “machine minus one while inheriting future additions” state unrepresentable.

**PRIMARY — precedent analysis:** The sibling’s caution transfers, but its operation restriction does not. A protocol `replace` can weaken a sealed global rule and needed addressing and target-hash machinery. Roster removal is also safety-sensitive, so it needs explicit syntax, provenance, preview, and a hard gate; however roster shrinking already exists through full replacement. Banning sparse removal would preserve the less-auditable mechanism where every omitted ID is ambiguous.

### 2. The existing decks

**PRIMARY — answer:** Every existing full `[roster.*]` list remains a full replacement. A delta interpretation requires the new explicit mode. There is no automatic protocol-sync rewrite and no fleet-wide semantic migration.

**PRIMARY — evidence-based migration rule:** The 36 current `zcode-1` omissions cannot be classified from file shape alone. Conversion must ask whether to adopt the machine member or record a deliberate removal, and it must skip or stop on the observed dirty/non-Git states unless the operator supplies the existing attended backup/restore path.

### 3. The legacy §2 table

**PRIMARY — answer:** It remains a full authority only while unmarked and legacy; it never becomes a delta. After attended migration or a projection marker, it stops being read for membership and `roster render` owns its body.

**PRIMARY — today’s defect boundary:** The isolated reproduction showed that any readable row currently revives rule 2 after deck blocks disappear. A projection marker separates a genuine legacy declaration from stale generated output without dropping backward compatibility.

### 4. Visibility and STATUS

**PRIMARY — answer:** Preserve the frozen column meanings and add `overlay-base`, `overlay-added`, and `overlay-removed` to the closed vocabulary. Show removed IDs as inactive rows, not by omission, and expose field-versus-membership provenance through `--explain`.

### 5. The anti-goal

**PRIMARY — mechanism:** One committed sparse delta replaces copied full lists; §2 remains generated; machine and deck membership writes retain `--confirm-breaking`; overlay mode is explicit; and per-idea `participants:` continues to lock quorum once Phase 0 closes. A machine change may affect future ideas, but it cannot rewrite an already-open idea’s participant list.

**PRIMARY — measurements:** Track, per release or fleet audit: counts of `replace` / `overlay` / `inherited` / `legacy` modes; number and size of add/remove deltas; redundant or unresolved tombstones; `section2-only` rows; distinct active-set fingerprints with the delta explaining each difference; and the number of decks affected by a proposed machine membership change. A healthy migration reduces full replacements and legacy rows without growing unexplained or redundant deltas.

### 6. Do nothing

**PRIMARY — answer:** This deck needs no overlay to follow the whole machine roster; its present zero-membership-declaration state already does that. The unserviceable cases are policy-relative membership: `machine + local specialist` or `machine − named member`, while continuing to inherit unrelated future machine changes. Today those require a full copied roster whose omissions do not say whether they are intentional.

## Concerns / open questions

**PRIMARY — machine dependence:** The base is user-local and therefore can differ across collaborators. That is already visible for `inherited-roster`, but overlay expands the number of machine-dependent decks. The run/idea must record the effective active IDs and a membership revision at kickoff; repository state alone cannot reconstruct an overlay without the machine base.

**PRIMARY — central blast radius:** A machine-scope membership change can affect every inherited or overlay deck’s next idea. The existing `--confirm-breaking` gate should be paired with an affected-deck preview where fleet discovery is available, and with an explicit generic blast-radius warning where it is not.

**PRIMARY — tombstone lifecycle:** A removed ID should stay suppressed if the machine entry disappears and later returns. Tooling still needs a precise warning and cleanup rule for tombstones that have been redundant for a defined period; automatic deletion would defeat the audit value.

**PRIMARY — conversion UX:** “Preserve current set” and “adopt machine set” must be separate named choices. A default that chooses either would silently decide whether the 36 observed omissions are intent or drift.

**PRIMARY — release coupling:** Parser, resolver, `roster set/show/explain`, STATUS documentation, JSON/text golden tests, legacy marker/render behavior, embedded/live protocol copies, and the skill vocabulary are one contract change and should ship atomically or behind one feature gate.

## Risks

**PRIMARY — silent quorum expansion:** A newly added machine member enters future inherited/overlay deck defaults. Mitigation is explicit opt-in mode, machine-change preview/confirmation, visible `overlay-base`, and the canonical Phase-0 participant lock.

**PRIMARY — silent quorum shrink:** A local remove affects future ideas. Mitigation is a named tombstone, `STATE=inactive`, `overlay-removed`, committed reviewable config, and `--confirm-breaking`; no value-only or local/env layer may change membership.

**PRIMARY — false-intent migration:** Converting copied rosters by set difference can canonize stale omissions. The current 36-deck `zcode-1` pattern demonstrates the risk; migration must be attended and choice-explicit.

**PRIMARY — legacy misclassification:** Treating every valid §2 row as legacy recreates today’s fallback surprise; treating every row as projection can silently expand a genuine legacy deck. The marker plus a blocking transition handles both directions.

**PRIMARY — contract drift:** New STATUS terms or marker prose added in only one protocol copy would break consumers or the drift guard. Atomic release tests and the existing exact-anchor guard are required.

**PRIMARY — central single point of drift:** Sparse overlays reduce copied lists but make correctness depend more heavily on the machine roster. Fleet audits must therefore measure machine-base changes and delta redundancy, not merely count distinct effective rosters.
