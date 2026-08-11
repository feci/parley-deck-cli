---
agent: hermes-1
idea: protocol-overlay-local-extension
round: 1
date: 2026-08-07
---

## Summary

The overlay must ship a registry in this slice, because D2 (consensus.md:88-93) forbids heading-text addressing and forbids inline markup, and no other addressing option exists in the record. The release store is empty today (`~/.parley/protocol/core/` absent; no `protocol-lock.yaml` under AI_WORKSPACE), so adding a registry file to the release layout costs nothing — and that window closes on the first published release. The v1 surface stays near-empty: one replaceable block (`s6.6`), one extension point (`ext-1`). The §2 roster table is NOT overlay content — it belongs to `agents.toml` + `roster render`, and the overlay must not become a competing roster surface (R7.4). The loss report must keep its hard invariant ("an empty report means no line disappeared, NOT that no meaning was lost", render.go:214) and gain an `Added` field so a second content source is visible. Two hazards are disqualifying as-is: H2 (repositioning false-loss) and H9 (prose-matched table addressing, probe-confirmed). Both must be fixed before the overlay ships onto this substrate.

## Q1 — D-a: Overlay file syntax

The record fixes the path (`parley-deck/protocol-overlay.md`, consensus.md:110), the two operations (replace-a-replaceable-block-by-ID, extend-at-ext-1, consensus.md:110-111), and the provenance fields (rationale + operation-specific hash, consensus.md:111-114). It does not ratify a grammar — all three candidates are dead (brief §4 D-a). I propose YAML frontmatter + Markdown payload.

Rationale, by elimination:

- TOML sidecar is the wrong tool: the payload for an `ext-1` extension is Markdown prose (D-p shows the four DF-4 decks carry HTML-comment roster blocks and header provenance lines; D-q shows three decks improvise dated change-log structures). A TOML sidecar forces escaping multi-line Markdown into a string field, which is hostile to the one operation (extend) that exists to carry prose.
- Markdown-with-fenced-directive-blocks embeds operations inside the prose, which collides with H10's probe-confirmed finding that `findLine` takes the first line with an exact prefix at column 0 and substitution fires inside fenced code blocks (render.go:160-167, brief H10). An overlay that documents protocol syntax in fenced blocks makes the slot-location bug worse. The record also forbids inline markup in the core body (consensus.md:93); while that prohibition is scoped to the core, putting operation directives in fenced Markdown blurs the line between "content" and "control" in a way the renderer's current Markdown-blind logic cannot safely handle.
- YAML frontmatter + Markdown payload separates operations (structured, in frontmatter) from content (free Markdown, in the body). The renderer already reads the overlay as bytes passed by the caller (R3.4, render.go:43-45), so parsing frontmatter is a clean addition. Go's standard library has no YAML parser, but the project already depends on one for `agents.toml`-adjacent config — and a frontmatter parser is a smaller surface than a directive-block parser that must coexist with the core's own fenced code.

Concrete grammar:

```yaml
---
core-version-range: ">=1.0.0 <2.0.0"   # R1.4, consensus.md:205-206
operations:
  - op: replace
    target: s6.6
    base-hash: <sha256 of the core's s6.6 block bytes>
    rationale: "Project works in Slovak; English-only rule replaced per user constraint 3."
    body: |
      6. **Working language.** All content written to any file under
      `parley-deck/` MUST be in English unless the project's own overlay
      declares otherwise. ...
  - op: extend
    target: ext-1
    id: <deck-slug>.<slug>              # D-d, deck-namespaced
    depends: [s1, s2, s6, s7, s15]      # D10, defaults to all sealed if absent
    dep-hashes:
      s1: <hash>
      s7: <hash>
      s15: <hash>
    rationale: "Project-specific packaged-reference drift notes."
    body: |
      ## Project-specific packaged-reference drift
      ...
---

# (body is the concatenation of operation bodies in order;
#  frontmatter is the control plane, body is the content plane)
```

Key points:

- `body` is a literal YAML block scalar, so the Markdown payload is preserved verbatim with no escaping.
- `base-hash` is required for `replace` (R1.3, consensus.md:111-114). For `extend`, `depends` + `dep-hashes` replace it — and the revision demanding a replaced-block hash from extensions was BLOCKED by codex-1 (R1.3), so this must not be reintroduced.
- `core-version-range` satisfies R1.4 and D10's need for a baseline (consensus.md:205-206).
- An empty file (no frontmatter, no operations) is rejected per R1.5/G4 (consensus.md:115-117). Absence of the file is the only "no customization" state.
- The grammar permits at most one `replace` in v1 (D-w), since `s6.6` is the only replaceable block (D3, consensus.md:96-102). Multi-replace ordering (D-w) is deferred — see ranking below.

## Q2 — D-b: Does this slice ship the block-ID registry?

Yes. This is not close.

The record creates a hard dependency: D2 (consensus.md:88-93) says addressing is by permanent, never-reused registry block ID, never by heading text, and no inline markup in the body. D3 (consensus.md:96-102) names `s6.6` as the one replaceable block and `ext-1` as the one extension point. If there is no registry, there is no `s6.6` and no `ext-1` — the IDs are the registry. H7 (brief) is confirmed by PRIMARY inspection: `CoreFileName = "COOPERATION.md"` (core.go:26), `Publish(version, body string)` writes only that file (core.go:137, 167-175), and `grep -i registry` over `internal/protocolcore/` and `internal/app/protocol.go` returns nothing. The registry does not exist.

The brief's D-b frames this as a choice, but the record leaves no third addressing option. D2 forbade heading text. D2 forbade inline markup. If this slice ships without a registry, it ships with no addressing at all — which means the overlay cannot target `s6.6` or `ext-1`, which means the overlay does nothing. That is not a deferred follow-up; it is the absence of the feature.

The cost is currently nil: `~/.parley/protocol/core/` does not exist on this machine (PRIMARY: `ls` returns "No such file or directory"), and `find` for `protocol-lock.yaml` under AI_WORKSPACE returns nothing. Zero published releases exist, so the release layout can grow a registry file for free. The moment the first release is published, that window closes — releases are write-once (core.go:137-159, D1 consensus.md:81-86). The kickoff (00-prompt.md:37-43) names this closing window as a "why now" reason. Waiting until after the first release means either living without the overlay or breaking write-once to retrofit a registry into a published release.

Sub-decision — registry file and block extent: The registry lives alongside the core body in the release directory as `registry.yaml`. Each entry maps an ID to `sealed | replaceable | extension-point` and records the byte range (start line, end line) of the block it covers. The byte-range choice is forced by H8: `s6.6` is list item 6 inside `## 6. Conflict-avoidance mechanics` (PRIMARY: `## 6.` heading at parley-deck/COOPERATION.md:732, the English-only rule at line 743). It is not heading-delimited, so heading-based segmentation cannot work even for the single block v1 needs. Byte offsets in the registry are the only option that handles both heading-delimited sections and list-item sub-blocks. The registry is hashed alongside the body (R2.2, consensus.md:93), so a registry change is a release change.

Tombstones for deleted IDs (R2.1, consensus.md:92) are entries with `status: tombstoned` and no byte range.

## Q3 — D-k: Is the §2 roster table overlay content at all?

No. The §2 roster table is NOT overlay content. It belongs to `agents.toml` (the roster authority, parley-deck/COOPERATION.md:102-105) and `roster render` (roster_render.go). The overlay must not own the rows, and it must not own bespoke columns.

Evidence:

- 11/29 decks rewrite the §2 table schema with bespoke columns (brief §1). But `agents.toml` is the authority and 23/29 tables already disagree with it (R7.4, brief). The problem is not "the overlay can't express columns" — the problem is "the fleet hand-edited a generated view and it drifted." Creating a second surface (the overlay) to carry columns that `agents.toml` cannot express would create a second competing roster authority, which R7.4 and the brief's hard constraint explicitly forbid.
- The bespoke columns (CLI/runtime, Model, Reasoning/effort, State, Display name, Backend, Invocation) are genuine config, not protocol content. They belong in `agents.toml` via `parley roster set`. The fix for the 11/29 decks is to extend `agents.toml`'s schema and `roster render`'s output, not to build an overlay escape hatch.
- D3's open surface (consensus.md:96-102) does not cover §2. The identity slots include "the §2 roster table body" (D3, consensus.md:99-100), but that is a renderer input (R2.4, consensus.md:99-100; round-02/codex-1.md: "data substitution, not OOP-style replacement"), not an overlay override.

What the overlay DOES own for §2: the annotations around the rendered table. Three decks append roster-decision logs after the §2 table (brief §1: auftra + ldx-wt-mail-fixups as HTML comments, librade as a blockquote; ldx writes numbered project prose "so they survive the next sync"). That is genuine project-local content that `agents.toml` cannot express — dated user directives, invocation gotchas, MANUAL-Bash caveats. That content is an `ext-1` payload or a replace of a future annotation block, not a roster-table override.

So the answer to D-k is option (b) with a caveat: the overlay owns nothing in the table, the table is always rendered from `agents.toml`, and the bespoke columns are handled by extending `agents.toml` and `roster render` (a separate follow-up, not this slice). The overlay reserves `ext-1` for the prose annotations around the table that `agents.toml` genuinely cannot express.

## Q4 — D-c / D-g: Placement and loss reporting

### D-c — Where ext-1 renders

ext-1 renders at a fixed named position: after §8 (Inbox) and before §10 (TL;DR). This is codex-1's position (round-02/codex-1.md:46-48, brief D-c), and I agree with it over "end of file" (kimi-1) or "a named block ID" (my own round-02 position).

Reasoning:

- "End of file" is non-deterministic in the presence of future core versions that append sections. A position defined as "the last line" moves every time the core grows. The registry solves this: `ext-1` is an `extension-point` entry in the registry with a fixed insertion point recorded as "after block s8, before block s10". That is a named position, not a tail.
- My round-2 position ("it must be a named block ID") was right in spirit but the registry (Q2) makes it concrete: `ext-1` is a registry entry whose position is "insert after the last line of block s8". The registry holds the position, not the overlay.
- kimi-1's "after the final core section" has the same non-determinism problem as "end of file".

This is coupled to H2 (brief): rendering ext-1 at a fixed position that differs from where a deck's local section sat mid-document produces a repositioning false-loss. That hazard is real and must be handled under D-g — but it does not change the placement decision, because the alternative (rendering the payload in place at the section's current position) kills the "one declared position" simplicity that D3 ratified (consensus.md:98, "rendered at a declared position"). The false-loss is a reporting problem, not a placement problem.

### D-g — Loss report semantics with a second source

The constraint that must survive (brief, render.go:207-216): an empty report means "no line disappeared", NOT "no meaning was lost". This is non-negotiable. G1 (consensus.md:255-257) requires reporting every block replaced or removed.

I choose option (iii): accept a one-time noisy report with an explicit migration note, plus a structural change to `RenderResult`.

Why not (i) — exempting overlay-carried content from the loss report: this risks re-creating the silent-erasure class that nine review cycles bought (brief D-g, IMPLEMENTATION.md:409-416). If overlay content is exempt, a bug that drops the overlay payload silently is invisible. G1's entire purpose is to make a dropped line visible. Exempting the thing the overlay exists to carry inverts that.

Why not (ii) — suppressing losses whose lines all reappear elsewhere: this weakens the order-sensitive sequence guarantee (render.go:193-202). The LCS diff was deliberately made order-sensitive because "the same lines in a different ORDER change the meaning". Suppressing repositioned losses tells the operator "nothing happened" when something did.

Option (iii) means: on the first render after an overlay is applied (or after a core bump that repositions ext-1 content), the loss report WILL show the deck's old local section as "not carried forward" because it moved from its old position to ext-1's declared position. The report is noisy but honest — every line that disappeared from its old position is listed, and the operator is told in the report body that the content was carried to ext-1. After the first render, the report is clean (the content is now at ext-1, the deck matches, idempotence holds).

Structural change: `RenderResult` (render.go:35-39) gains an `Added []string` field and per-block provenance. Today it has `Body`, `Removed`, `Preserved` — there is no `Added` (H3, brief). An overlay that injects a whole `## EXT` section yields `REMOVED: []` and no indication that content was added or by whom (H3). G1 says "replaces or removes", but a literal reading leaves the overlay's entire payload unreported, which is unacceptable for a second content source. Both call sites (protocol.go:218, protocol.go:258) need to surface the new field. The `Added` field reports blocks the overlay contributed, identified by operation ID, so the operator can see "the overlay added ext-1 payload `<deck>.<slug>`" in the same report that shows what moved.

The stamp regex (H4) must change in the same commit. `generatedStampRe` (render.go:380) is format-locked to `core <v> (<h>)`. If the stamp gains an overlay/effective hash (D-j), the regex must change atomically or the exemption silently dies and every render reports a spurious one-line loss (H4, brief). This is a coupled requirement, not a separate decision.

## Q5 — Disqualifying hazards

Two hazards are disqualifying as-is. The overlay must not ship onto the current substrate until both are fixed.

### H9 — Prose-matched table addressing (DISQUALIFYING, probe-confirmed)

PRIMARY confirmation: `isTableHeader` requires the line to start with `| Agent ID` and contain "Workspace" or "Host handle" (render.go:129-133). Removal report groups by heading text in `map[string]int` (render.go:239-240). `roster render` locates §2 by a literal heading constant `rosterSectionHeader = "## 2. Active agents (roster)"` (roster_render.go:24, used at :157-180). All three violate D2 ("Addressing is by ID, never by heading text", consensus.md:93).

The probe-confirmed consequence (brief H9): if the deck's table header text differs from the core's (core "Workspace dir" vs deck "Workdir"), the render emits header + separator with zero data rows — the deck's whole roster is wiped, reported only as "## 2 — 2 lines not carried forward". I additionally confirmed the drift-guard/renderer column conflict: the drift guard requires `| Agent ID       | Workspace dir                       | Role          |` exactly once (drift_test.go:28, 59-61), while the renderer emits `| Agent ID | Workspace dir | Role | State |` (roster_render.go:73). These are already mutually incompatible.

D-t asks whether table zones move to ID addressing in the same slice or later. My answer: in the same slice, because D2 is already ratified and H9 is probe-confirmed. Shipping the overlay onto a substrate that already violates D2 — where a core column rename silently empties every deck's roster — is shipping a feature on top of a known data-loss bug. The registry (Q2) makes this fix cheap: table zones get registry IDs, and `isTableHeader` / `tableBodyFor` / `rosterSectionHeader` are replaced by registry-ID lookups. This is part of the registry work, not a separate follow-up.

### H2 — Repositioning false-loss (DISQUALIFYING as-is, fixable in-slice)

`droppedContent` is an LCS sequence diff grouped by nearest preceding heading (render.go:217-272, LCS at :300-335), deliberately order-sensitive (render.go:193-202). The brief's probe confirms: a deck whose local section sits mid-document, re-rendered with that section carried at ext-1 (end of file), reports it as REMOVED even though every line is carried forward. Shipping the overlay produces a false data-loss report on exactly the decks it exists to rescue.

This is disqualifying as-is because it trains operators to ignore the G1 report — if every rescue produces a false alarm, the real alarms are ignored. But it is fixable in-slice without weakening the order guarantee: the `Added` field (Q4/D-g) and the explicit migration note make the report honest. The line DID disappear from its old position; the report says so; the report ALSO says it was carried to ext-1. The operator reads a noisy-but-true report, not a false one. The order-sensitive guarantee survives because we are not suppressing the loss — we are explaining it.

### Hazards that are serious but not disqualifying

- H6 (lock parser is a prefix scan that silently ignores unknown keys, protocol.go:92-98): serious — a lock declaring `overlay-sha256:` is read as if the field did not exist, producing silent data loss. But this is entangled with rank 2 (D-i, brief §5) and the lock is not yet written by any shipped code path. Fix in this slice only to the minimum extent: refuse unknown keys in the lock parse so an old binary reading a new lock fails closed rather than silently rendering without the overlay. Full YAML parse is rank 2.
- H4 (stamp regex format-locked, render.go:380): must be fixed in the same commit as D-j (stamp content). Not disqualifying on its own — it is a coupled requirement.
- H12 (four writers, one file): `protocol render` (protocol.go:233), `roster render` (roster_render.go:146), and `preflight.syncConsumerProtocol` (preflight.go:488, using `mergePreservingZones` with a `## 3.` byte-offset anchor at preflight.go:539) all write `COOPERATION.md`. D-s asks whether `protocol render` becomes the sole writer. My position: yes, `protocol render` must become the sole writer, and `roster render` / `preflight` must feed inputs to it rather than writing independently. But this is a structural refactor that can be staged — the overlay ships with `protocol render` as the authoritative writer for overlay-relevant content, and the consolidation of `roster render` and `preflight` is a named follow-up (D-s, see ranking). Not disqualifying because the overlay's own write path is clean; the risk is a later `preflight` run reintroducing a body the core did not produce (H12), which is a pre-existing bug, not one the overlay creates.
- H13 (drift guard vs overlay on this repo, drift_test.go:46): if the CLI's own deck ever renders an ext-1 payload, `TestEmbeddedDefaultMatchesLiveDeck` fails. D-r asks whether the CLI's own deck ever carries an overlay. My position: the CLI's own deck should NOT carry an overlay — it is the source repo. Add a standing rule, not a sixth allowlisted zone. Not disqualifying — it is a test-harness decision, not a data-loss path.

### Hazards that are not disqualifying

- H1 (in-sync predicate collapses four states, protocol.go:218, :258, :261): the `check --json` status enum has two values. With an overlay, "core drifted", "overlay changed", "overlay missing", and "hand-edited" collapse into one `hand-edited-or-stale`. This is a reporting deficiency, not a data-loss path. Fix in this slice by expanding the enum, but it does not block the overlay's core function.
- H3 (additions structurally invisible): fixed by the `Added` field (Q4/D-g). Not separately disqualifying.
- H5 (stamp built from core alone, render.go:95): fixed by D-j (derive from lock, R3.3). Coupled, not disqualifying.
- H7 (nothing to address): fixed by shipping the registry (Q2). This is the registry decision, not a separate hazard.
- H8 (s6.6 is a list item, not heading-delimited): handled by byte-offset registry entries (Q2). Not disqualifying — it constrains the registry format, not the overlay's viability.
- H10 (slot location fragile and Markdown-blind): `findLine` takes the first line with an exact prefix at column 0 (render.go:160-167), substitution fires inside fenced code blocks (render.go:96-98). An overlay that documents protocol syntax in fenced blocks makes this worse. Not disqualifying because v1's open surface is one replaceable list item and one extension point — the fragile slot logic is not on the overlay's critical path. But it must be fixed before the open surface widens. Named follow-up.
- H11 (asymmetric guard between sources, render.go:65-67): a core `**Protocol synced:**` line is dropped unconditionally before slot substitution; no equivalent for the overlay. Fix in this slice — the overlay's stamp line needs the same guard. Small.
- H14 (zero production decks): this is the reason the registry is cheap now (Q2). Not a hazard to the overlay — it is a hazard to delaying the overlay.
- H15 (two hard-coded promises that must retire atomically, protocol.go:211, parley-deck/COOPERATION.md:767-768): both confirmed at exact cited lines. Both must change in the same commit the overlay ships. This is a coupled requirement, not a separate hazard.
- H16 (test coverage one layer up): `render_test.go` is 45 lines covering only the LCS primitive. Adding an overlay parameter to `Render` breaks compilation in two files. Not disqualifying — it is a test-coverage gap, not a design blocker. The behavioral tests in `protocol_test.go` are the real guard.
- H17 (fleet noise floor on first render): 27/29 decks missing the default's Blast radius block from sync lag. Not an overlay hazard — it is a pre-existing sync-lag issue. The overlay does not make it worse. Named follow-up for the fleet migration (DF-2).
- H18 (placeholder pollution): 4 decks carry `<workspace-name>`, etc. An overlay extraction that reads these as intentional local content freezes a defect. Not disqualifying — it is a migration-tooling concern (DF-2), not an overlay-design concern. The overlay must not auto-extract; R1.6 (consensus.md, round-01/codex-1.md:149-151) already forbids auto-creation.

## Ranking of remaining open decisions

### Must settle in this slice

1. D-b (registry: ship it now) — SETTLED ABOVE. No registry = no overlay.
2. D-a (overlay file syntax) — SETTLED ABOVE. YAML frontmatter + Markdown payload.
3. D-k (§2 roster table is overlay content?) — SETTLED ABOVE. No; `agents.toml` owns it.
4. D-c (where ext-1 renders) — SETTLED ABOVE. Fixed named position after §8, before §10, via registry.
5. D-g (loss report semantics) — SETTLED ABOVE. Option (iii) + `Added` field + per-block provenance.
6. D-t (prose-matched table addressing replaced now?) — YES, in this slice. H9 is disqualifying.
7. D-h (missing or unreadable overlay file) — BLOCK, like a missing pinned release (protocol.go:118-128). Rendering without it silently erases the overlay's content — a literal repeat of the 2026-08-06 incident. This is a one-line decision but it must be explicit: a missing overlay file when the lock declares `overlay-hash: <non-none>` is a block, not a render-without. An absent overlay when the lock declares `overlay-hash: none` is fine (R1.5).
8. D-j (stamp content and regex, same commit) — The stamp carries core + overlay + effective hash (D7, consensus.md:141). `generatedStampRe` (render.go:380) changes in the same commit. The stamp is one field, not an append-only log (the drift guard demands exactly one, drift_test.go:72-77). This is coupled to H4/H5 and must be atomic.
9. D-x (line endings between three sources) — The deck keeps deciding output convention and the overlay is normalized in, same as core (render.go:53-54, :57). Undecided = the exact non-convergence bug already fixed once. This is a one-line decision but it must be explicit to avoid reintroducing the bug.
10. D-w (multi-replace ordering) — State explicitly that v1 permits at most one replace (only `s6.6` is replaceable). Specify deterministic ordering (by registry ID byte order) for the future case. Decide cheaply now; costs nothing.
11. D-d (deck-namespaced extension ID format) — `<deck-slug>.<slug>`, with collision validation: the deck slug is derived from the deck directory name, the slug is author-chosen, and the full ID must not collide with any registry core ID (which are `s<N>` / `ext-<N>`). kimi-1's format was never ratified (brief D-d); I propose this and defend it. Must settle in this slice because `ext-1` payloads need IDs to be referenced in the loss report's `Added` field.
12. D-u (is user constraint 3 currently satisfied?) — This run does NOT discharge kimi-1's condition automatically. kimi-1 signed off on D3's near-empty surface "on condition that rank 3 actually ships this cycle" (brief D-u). Rank 3 did not ship. This idea IS rank 3. If this idea reaches FINAL with a shipped overlay design, the condition is discharged. If it slips again, the condition must be re-raised explicitly. Must settle in this slice because it is the framing condition for the whole idea.

### Can be deferred with a named follow-up

13. D-e (where D10's reconfirmation is recorded) — Defer to follow-up `overlay-receipt-artifact`. Define reconfirmation as a lock-field update alone for v1 (the lock already carries overlay hash, consensus.md:166-167). The receipt file (codex-1's proposal) is a nice-to-have that adds a second artifact; v1's lock-field update is sufficient because the lock is committed and byte-verified (G8, consensus.md:283-291). The risk of deferring: no standalone proof of reconfirmation outside the lock. Acceptable for v1 because the lock IS the proof.
14. D-f (authoring surface: read-only vs writer) — Defer to follow-up `overlay-authoring-cli`. Ship read-only (`protocol overlay show|validate`, file hand-written). A writer verb (`protocol overlay set|remove`) is a usability improvement, not a correctness requirement. R5.1 says the overlay is not user-only, but "not user-only" means an agent may propose an overlay change through a normal idea — it does not require a CLI writer verb. The file is hand-written and committed, same as `agents.toml` was before `roster set` existed.
15. D-i (lock parsing: real YAML parse?) — Partially in-slice (refuse unknown keys, fail closed), full YAML parse deferred to rank 2 (brief §5). The unknown-key refusal is in-slice because H6 is a silent-data-loss vector; the full parse is rank 2 because the lock's full field set belongs to rank 2.
16. D-l (Transport: identity or core?) — Defer to follow-up `transport-authority`. 27/29 decks override it, §0 calls it sticky and requires a §7 idea to change (parley-deck/COOPERATION.md:54), and it is excluded from the drift allowlist. The tension is real but it is not on the overlay's critical path — Transport is already an identity slot (render.go:9-20, `transportPrefix`). The overlay does not need to own it for v1. The follow-up decides whether Transport becomes overlay-addressable or stays protocol-governed.
17. D-m ("six identity slots" doc/impl mismatch) — Fix the count in D3's language in this slice (it is a doc fix, not a design decision), but defer the question of whether `**Parley deck:**` becomes a slot to follow-up `identity-slots-complete`. The doc fix is in-slice because shipping the overlay with a known doc/impl mismatch in the same area is sloppy; the seventh slot is not on the critical path.
18. D-n (how Workspace is derived) — Defer to follow-up `workspace-derivation`. 16/29 decks do not match their project directory. Workspace is already an identity slot (render.go:15). The overlay does not need to own it for v1. The follow-up decides: overlay stores it, or derived from project dir with overlay supplying a label.
19. D-o (host-handle table under local-dir) — Defer to follow-up `host-handle-table-policy`. 13/29 empty, 9 filled with noise. This is a rendering hygiene issue, not an overlay issue. The overlay does not touch the handle table for v1.
20. D-p (DF-4's content: ext-1 payload or core?) — Defer to follow-up `df-4-content-restoration`. The four decks carry content as HTML-comment roster blocks and header provenance lines. D2 forbids inline markup in the core body — the sub-question of whether that extends to an ext-1 payload is unsettled. My position: the payload is free-form Markdown, so HTML comments are permitted in ext-1 payloads (the D2 prohibition is scoped to the core body, consensus.md:93). But this does not block the overlay's design — it blocks DF-4's specific restoration, which is a named follow-up (consensus.md:250-251).
21. D-q (provenance/justification as first-class structure?) — Settled implicitly by Q1's grammar: the `rationale` field per operation (R1.3, consensus.md:111-114) is the home for all of it. A separate dated change-log structure is not needed for v1. The three decks that improvise it (servers, auftra, librade) will migrate their prose into `rationale` fields. If a deck needs a running change log, that is an `ext-1` payload, not a separate structure. Can be deferred — if the `rationale` field proves insufficient, a follow-up adds a change-log structure.
22. D-r (does the CLI's own deck carry an overlay?) — Settled above: no, the source repo's deck never carries an overlay. Standing rule, not a sixth allowlisted zone. This is a decision, not a deferral.
23. D-s (does protocol render become the sole writer?) — Defer the full consolidation to follow-up `single-writer-consolidation`. The overlay ships with `protocol render` as the authoritative writer for overlay-relevant content. `roster render` and `preflight.syncConsumerProtocol` continue to write independently for now, with a documented warning that running `preflight` after `protocol render` reintroduces a body the core did not produce (H12). The full consolidation is a refactor that touches three subsystems and is not on the overlay's critical path.
24. D-v (never-synced decks) — Defer to follow-up `never-synced-deck-policy`. Exclude never-synced decks from overlay extraction until they are synced (the evidence says staleness, not project-local removal). This is a migration-tooling concern (DF-2), not an overlay-design concern.

## Concerns / open questions

1. The registry's byte-range approach assumes the core body is line-stable across the version range the overlay declares. If a core release reflows lines within a block (e.g. reformatting a list item) without changing the block's ID or semantic content, the byte range changes and the `base-hash` mismatches. Is the base-hash computed over the block's raw bytes (line-sensitive) or over a normalized form? I propose raw bytes, matching how the core is hashed today (core.go:189-192, `Hash` is over the raw string). A reflow is a new release; the overlay re-confirms. This is consistent with D10 (consensus.md:200-208) but should be stated explicitly.

2. The `Added` field in `RenderResult` (Q4/D-g) is a breaking change to the `Render` function signature (render.go:46). H16 notes this breaks compilation in two files. The overlay parameter itself is also a new parameter. The signature becomes `Render(rel Release, overlay Overlay, priorDeckBody string) (RenderResult, error)`. Is `Overlay` a parsed struct or raw bytes? I propose raw bytes parsed inside `Render` (preserving purity, R3.4) — but the YAML frontmatter parse happens inside `protocolcore`, which today has no YAML dependency. This is an implementation detail, but the dependency choice should be settled before FINAL.

3. D-u (kimi-1's condition) is the framing condition for this idea. If this idea does not reach FINAL with a shippable overlay design, user constraint 3 remains unmet and the condition must be re-raised. I judge this idea capable of discharging it — but only if it reaches FINAL, not merely consensus.

## Risks

1. The registry is the single most expensive part of this slice. It changes the release layout, which is write-once. If the registry format is wrong, the first published release is wrong forever. Mitigation: zero releases exist today (H14), so the format can be iterated before the first publish. But the iteration window is the time between this idea's FINAL and the first `parley protocol publish` — which could be short.

2. The `Added` field and per-block provenance (Q4/D-g) add reporting complexity to a function that was wrong in "nine distinct ways across nine review cycles" (render.go:207-208). The honest posture is that the loss report is a line-level diff and nothing more. Adding `Added` does not change that — it reports what was added, not what it means. But it is new surface area in a function with a history of subtle bugs.

3. Fixing H9 (prose-matched table addressing) in this slice means changing `isTableHeader`, `tableBodyFor`, and `rosterSectionHeader` to registry-ID lookups. This touches the roster renderer (roster_render.go) and the drift guard (drift_test.go:28-29), which are already mutually incompatible (the 3-col vs 4-col header conflict). The fix must reconcile all three — registry, roster renderer, drift guard — or it trades one bug for another.

4. The "one-time noisy report" (D-g option iii) is honest but has a usability cost: every deck's first render after overlay adoption produces a non-empty loss report. If operators learn to ignore the first report, the G1 guarantee is eroded by fatigue rather than by silence. Mitigation: the report body explicitly states "content carried to ext-1" so the operator can distinguish a repositioning from a real loss. But this is a social mitigation, not a structural one.
