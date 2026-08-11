# PROTOCOL OVERLAY — Scoping Brief

**Status:** rank 3 of 4 in the ratified staging (consensus.md:230-232). Design not started. Hard prerequisite for DF-2 fleet migration (IMPLEMENTATION.md:471-473).

---

## 1. What the overlay must do

**The ratified job.** Give a deck a committed, machine-checked way to carry project-local protocol content across a core render, without letting any deck weaken the core. One file per deck at `parley-deck/protocol-overlay.md`, two operations — replace a replaceable block by ID, extend at `ext-1` — with operation-specific provenance (consensus.md:108-114). Absence of the file is the only canonical "no customization" state; an empty file is forbidden (consensus.md:115-117). Composition is core-owned in five fixed steps, and a contradictory extension is *incompatible*, not higher-precedence (consensus.md:118-128).

**The measured job.** The empirical taxonomy over 29 real decks says what the fleet actually needs:

- **Zero of 29 decks are byte-identical to the packaged default.** Every deck carries project-local content.
- For the 27 modern decks, **all divergence lives in exactly two zones**: the header block (default L3–L7) and §2 Active agents (default L130–L146). Median 4–5 hunks. This is the strongest single piece of evidence the near-empty v1 surface (D3: one replaceable block, one extension point) is right.
- The largest genuine override class is **not** ratified as an override at all: **11/29 decks rewrite the §2 roster table schema** with bespoke columns (CLI/runtime ×5, Model ×7, Reasoning/effort, State ×4, Display name, Backend, Invocation-CLI). Only 3/29 carry the header `parley roster render` actually emits.
- Only **1/29 decks** adds a whole new `##` section (librade-algoTrader, "Project-specific packaged-reference drift"). That is the single true `ext-1` case in the fleet — consistent with the 1-in-36 evidence D3 was built on (consensus.md:103-105).
- **3/29 decks** append roster-decision logs after the §2 table (auftra + ldx-wt-mail-fixups as HTML comments, librade as a blockquote); **1/29** (ldx) writes numbered project prose inside §2 explicitly "so they survive the next sync" (ldx/parley-deck/COOPERATION.md:147). Durability across sync is the user requirement in the users' own words.
- The failure mode is proven, not hypothetical: the servers deck's pre-sync backup shows a whole-file sync destroying an in-body `Transport: local-dir` declaration (bak L67-68) and silently flipping the header transport `github-pr` → `local-dir` (bak L5-5 → cur L5-15).
- Today's tooling **names** the loss but does not prevent it: rendering `auftra` reported 13 lines lost in §2 and applying it dropped genuine content (IMPLEMENTATION.md:465-473).

**The deliverable must therefore be:** a registry-addressed, fail-closed second content source that (a) survives a core version bump, (b) subsumes or coexists with the five zones the drift guard already enforces, and (c) makes the auftra 13-line loss report go to zero without weakening G1's loss detection.

---

## 2. Requirement list

Each requirement is traceable; nothing here is invented.

### R1 — File and shape
- **R1.1** Exactly one overlay file per deck, committed, at the fixed deck-relative path `parley-deck/protocol-overlay.md`. Not configurable, not a lock field for its location. (consensus.md:109-110)
  - *Record discrepancy to settle explicitly:* kimi-1's round-01 said `protocol.overlay.md` (round-01/kimi-1.md:20); D4 says `protocol-overlay.md` (consensus.md:110). D4 governs; settle it in writing rather than by whoever types first.
- **R1.2** Exactly two operations: replace-a-replaceable-block-by-ID and extend-at-`ext-1`. No third operation was ratified. (consensus.md:110-111)
- **R1.3** Operation-specific provenance: a replace records the expected hash of the target block; an extend records D10's dependency set and their hashes. Each records a rationale. (consensus.md:111-114) — A revision demanding a replaced-block hash from extensions was BLOCKED by codex-1; do not reintroduce it.
- **R1.4** The overlay declares the **core version range** it was written against. (consensus.md:205-206, DPC-4)
- **R1.5** Absence = no customization. An empty overlay file MUST be rejected. `protocol status` must distinguish absent (fine) from empty (error). (consensus.md:115-117, G4)
- **R1.6** No overlay is created at bootstrap or migration unless reviewed local content actually exists. (round-01/codex-1.md:149-151)

### R2 — Addressing
- **R2.1** Addressing is by permanent, never-reused registry block ID; never by heading text or section position. Deleted IDs get tombstones. (consensus.md:88-93)
- **R2.2** No inline markup in the core body. The registry lives outside the prose, hashed alongside the body. (consensus.md:93)
- **R2.3** v1 open surface: replaceable = `s6.6` (working language) **only**; extension-point = exactly one, `ext-1`, deck-namespaced IDs; six identity slots as typed renderer inputs; everything else sealed, explicitly including §7 and §15. (consensus.md:96-102)
- **R2.4** Identity slots are **not** overlay overrides — they are typed renderer inputs. (consensus.md:99-100; round-02/codex-1.md:120-122: "They are data substitution, not OOP-style replacement of a declared protocol part.")

### R3 — Composition
- **R3.1** Five-step core-owned order: load+verify the pinned release → fill the six identity slots → apply replace → append the single `ext-1` payload → validate sealed blocks present, hash the effective bytes. Identity fill precedes overlay application. (consensus.md:120-126)
- **R3.2** Never last-writer-wins. A contradictory extension is incompatible, not a higher-precedence exception. (consensus.md:118, 127-128)
- **R3.3** The renderer is a **NEW pure function** taking (core release, overlay, identity slots) → bytes, with the synced-stamp derived from the **deck lock** (not from `rel.SHA256` as today, render.go:95). `mergePreservingZones` is zone-extraction scaffolding for migration only; its `## 3.` anchor does not survive. (consensus.md:224-229)
- **R3.4** Purity is preserved: overlay bytes are read by the caller in `internal/app/protocol.go` and passed in, never opened inside `protocolcore`. (render.go:43-45 — "no filesystem, no clock … which is what makes G1's idempotence testable at all")
- **R3.5** The overlay needs its own explicit CRLF normalization line, per the per-source convention already established for deck (render.go:53-54) and core (render.go:57, comment: "A CRLF CORE must be normalized as well, or the render emits mixed endings, never converges across two runs").

### R4 — Fail-closed
- **R4.1** Unknown ID, duplicate provider, sealed target, or changed base hash → fail closed. (consensus.md:114-115, G4 at consensus.md:263-264)
- **R4.2** Compatibility check differs per operation. Replace: target exists + still replaceable + base-block hash matches. Extend: declared dependency set, **defaulting to all sealed blocks**; any change there produces a reviewable change report requiring reconfirmation. (consensus.md:200-206)
- **R4.3** Outcomes: mismatch → re-confirm; missing / tombstoned / now-sealed target → **block**. Auto-pass only on a zero-change report. **Never auto-migrate prose.** An incompatible deck stays pinned on its old core — that pinning *is* the quarantine. (consensus.md:206-208)

### R5 — Authority and gating
- **R5.1** A deck-overlay change goes through a **normal idea in that deck**. It must NOT be forced through the meta-idea + user-ratification path reserved for the core. (consensus.md:212-214; parley-deck/COOPERATION.md:767-768)
- **R5.2** G2's TTY-attended gate is scoped to `~/.parley/protocol/core/` only. `parley-deck/` is outside it, and `protocol render --yes` already writes the deck with no TTY check (protocol.go:222-238). Do not silently widen G2.
- **R5.3** The overlay must be incapable of editing core rules — §7 and §15 sealed (consensus.md:101-102). Empirically 28/29 decks never touched a line outside the header and §2, so the seal costs the fleet nothing measurable.

### R6 — Guarantees and tests
- **R6.1** G1 stays binding: `protocol render` MUST be idempotent (byte-identical on a second run) and MUST report every block it replaces or removes, in preview and on apply. (consensus.md:255-257)
- **R6.2** G7b: every overlay guarantee needs an end-to-end test through the **real** command entry point — "A guarantee without such a test MUST NOT be documented as landed." (consensus.md:277-282). The pattern to instance is `TestProtocolIsReachableThroughProductionDispatch` (protocol_test.go:439-457), which drives `Run`, not `runProtocol`.
- **R6.3** Extend the existing fixtures: `protocolFixture` / `protocolFixtureWith` (protocol_test.go:57-78, 669-688) — an overlay is one extra file write per fixture.
- **R6.4** Mirror the existing security-test patterns for any overlay path/ID input: traversal rejection asserted **by reason** (`TestProtocolRejectsPathTraversalInTheLock`, protocol_test.go:352-371, asserts on the string "unsafe version"), and a planted real target so the refusal cannot pass for the wrong reason (`TestLoadRefusesToEscapeTheStore`, protocol_test.go:375-393).
- **R6.5** Every fix proven by a **compiling, actually-applied** revert going red — the discipline recorded after one "failure" came from a non-compiling revert and one "pass" from a revert that never applied (IMPLEMENTATION.md:382-388).

### R7 — Must not break what exists
- **R7.1** The five zones the drift guard already enforces (Workspace, Created, Protocol synced, §2 roster body, §2 host-handle body — drift_test.go:16-30) must be subsumed or preserved, or `TestEmbeddedDefaultMatchesLiveDeck` (drift_test.go:46) fails.
- **R7.2** The embedded default must stay a bootstrap template with empty table bodies and generic placeholders (`assertEmptyTableBody`, drift_test.go:37-41, 149-179). The overlay cannot be implemented by seeding example content into the default.
- **R7.3** Prose around the §2 table must survive regeneration (roster_render.go:156-157). Extra table *columns* currently do not — 11/29 decks depend on those columns.
- **R7.4** `agents.toml` stays the roster authority; the overlay must not become a second competing roster surface. 23/29 tables already disagree with agents.toml (dominant cause: `opencode-1` present in toml, absent from §2 in 18 decks).
- **R7.5** Legacy hand-written decks keep working without migration as a precondition (defaults/COOPERATION.md:123-126).

---

## 3. Design hazards

These are the specific places today's code becomes ambiguous or wrong the moment a second content source exists. All from the readers' probes and cited lines.

### H1 — The in-sync predicate collapses four distinct states into one
Both Render call sites use `res.Body == string(prior)` as the sole no-op / in-sync test (`internal/app/protocol.go:218` and `:258`, status string at `:261`). With an overlay, this silently conflates *core drifted*, *overlay changed*, *overlay missing on this machine*, and *hand-edited* into one `hand-edited-or-stale`. The `check --json` status enum has exactly two values (protocol.go:258-274).

### H2 — The repositioning false-loss (the most important single hazard)
`droppedContent` is an LCS **sequence** diff grouped by nearest preceding heading (render.go:217-272, LCS at :300-335), deliberately order-sensitive (render.go:193-202: "the same lines in a different ORDER change the meaning"). Probe-confirmed: a deck whose local section sits mid-document, re-rendered with that section carried at `ext-1` (end of file), reports `REMOVED: ["## 99. Local rule — 1 line not carried forward", "## 4 — 2 lines not carried forward"]` **even though every line is carried forward**. Shipping the overlay produces a false data-loss report on exactly the decks it exists to rescue — including all four DF-4 decks.

### H3 — Additions are structurally invisible
`RenderResult` has three fields: `Body`, `Removed`, `Preserved` (render.go:35-39). There is no `Added`, no per-line provenance, no effective hash, no overlay identity. Probe: injecting a whole `## EXT` section yields `REMOVED: []` / `PRESERVED: ["Created"]` — the operator sees no indication that content was added or by whom. G1 says "replaces or removes", so a literal reading leaves the overlay's entire payload unreported.

### H4 — The stamp regex is format-locked and fails silently
The exemption that forgives the regenerated stamp line requires `generatedStampRe` = `^\*\*Protocol synced:\*\* core \S+ \([0-9a-f]+\)\s*$` (render.go:380). Probe: `"… core 1.0.0 (abc123def456)"` matches; `"… (abc123def456) + overlay (deadbeef)"` does **not**; `"**Protocol synced:** effective abc123 (core 1.0.0)"` does **not**. If the stamp gains an overlay/effective hash (D7, consensus.md:141) without the regex changing in the same commit, `headerStamp` returns `""` forever, `forgive` is permanently false (render.go:233-235), and **every** render reports a spurious one-line loss under `(document header)` (render.go:274) — training operators to ignore the G1 report. `headerStamp` also requires the stamp be the line immediately after `**Created:**` (render.go:384-395), so anything inserted between them disables the exemption too.

### H5 — The stamp itself would lie
The stamp is built from the core alone: `fmt.Sprintf("%s core %s (%s)", syncedPrefix, rel.Version, ShortHash(rel.SHA256))` (render.go:95). Two decks with different overlays under the same core produce an identical stamp. R3.3 already requires deriving it from the lock instead.

### H6 — The lock parser is a prefix scan that silently ignores unknown keys
`pinnedVersion` loops lines and `strings.CutPrefix(strings.TrimSpace(l), "core-version:")`, returning on first hit (protocol.go:92-98); the path is fixed at protocol.go:101. A lock declaring `overlay-sha256:` / `effective-sha256:` / `resolver-version:` is read by today's binary **as if those fields did not exist** — it renders without the overlay and reports the overlay's content as not carried forward. That is a forward-compatibility silent-data-loss path, and it is exactly the documented-not-wired failure G8 exists to prevent (consensus.md:283-291). G8 is not implemented today (core.go:84-85; FINAL.md:73).

### H7 — There is nothing to address
A release directory contains exactly one file: `CoreFileName = "COOPERATION.md"` (core.go:26, path at core.go:100, Publish writes only that at core.go:167-175). `Publish(version, body string)` (core.go:137) has nowhere to put a registry. `grep -i registry` over `internal/protocolcore/` and `internal/app/protocol.go` returns no matches. D2's entire addressing substrate is absent.

### H8 — The one v1 replaceable block is not heading-delimited
`s6.6` (working language) is an ordered **list item** inside `## 6. Conflict-avoidance mechanics`: `parley-deck/COOPERATION.md:743` — `6. **English only.** All content written to any file under parley-deck/ MUST be in English…`, under the `## 6.` heading at line 732. Block segmentation cannot be heading-based even for the single block v1 needs.

### H9 — The renderer already violates D2 in three places
- Table zones located by **prose match**: `isTableHeader` requires the line to start with `| Agent ID` and contain "Workspace" or "Host handle" (render.go:129-133) — directly against "Addressing is by ID, never by heading text" (consensus.md:93).
- Removal report grouped by heading **text** in `map[string]int` (render.go:239-240) — two sections with the same heading merge into one group.
- `roster render` locates §2 by a literal heading constant: `rosterSectionHeader = "## 2. Active agents (roster)"` (roster_render.go:24, used at :157-180).

Probe-confirmed consequence: if the deck's table header text differs from the core's (core "Workspace dir" vs deck "Workdir"), the render emits header + separator with **zero data rows** — the deck's whole roster is wiped, reported only as `REMOVED: ["## 2 — 2 lines not carried forward"]`.

### H10 — Slot location logic is fragile and Markdown-blind
- `findLine` takes the **first** line in the whole document with an exact prefix at column 0, untrimmed (render.go:160-167) — not scoped to the header block, so any later line (including overlay content) with that prefix can become the deck's identity.
- Inconsistent trimming: `findLine` does not trim (render.go:162), `isTableHeader`/`tableRows` do (render.go:174). Probe: "an indented identity line is not found (findLine requires column 0)".
- Substitution fires inside fenced code blocks. Probe: `**Workspace:** \`real-deck\`` appears both in the header and inside a ``` fence; `PRESERVED: [Workspace Created Workspace]`. Same for stamp insertion — `strings.Replace(body, created, created+"\n"+stamp, 1)` (render.go:96-98), probe: "stamp landed INSIDE the code fence (stamp@50 fence@18)". An overlay that documents protocol syntax makes this far more likely.
- `tableBodyFor` tests "Host handle" first (render.go:136-146) while `ExtractIdentity` calls `tableRows(body, "Workspace")` (render.go:155-156) — one table can silently feed two slots.
- Separator detection requires a literal `| -` prefix (render.go:76, :178, :182). Probe: a `|---|` separator is consumed as a data row.

### H11 — Asymmetric guard between sources
Any **core** line starting with `**Protocol synced:**` is dropped unconditionally before slot substitution (render.go:65-67). No equivalent guard would exist for the overlay, so an overlay line with that prefix survives into the output.

### H12 — Four writers, one file
`protocol render` (protocol.go:233), `roster render` (roster_render.go:146), and `preflight.syncConsumerProtocol` (preflight.go:509) all write `parley-deck/COOPERATION.md` today, each with its own location logic; the overlay is a fourth content source into a file with no single owner. Worse, `syncConsumerProtocol` still rewrites from the **packaged skill body** using `mergePreservingZones` with a `## 3.` byte-offset anchor (preflight.go:522-560) — a coarser preservation model that consensus.md:225-228 explicitly demotes to scaffolding. Running preflight after `protocol render` reintroduces a body the core did not produce.

### H13 — Drift guard vs the overlay, on this very repo
`TestEmbeddedDefaultMatchesLiveDeck` (drift_test.go:46) byte-compares `../../parley-deck/COOPERATION.md` against the embedded default outside five normalized zones (drift_test.go:102). If parley-deck-cli's own deck ever renders an `ext-1` payload, that test fails. Separately, the guard and `parley roster render` are already mutually incompatible: the guard requires `rosterHeaderLine = "| Agent ID       | Workspace dir                       | Role          |"` exactly once (drift_test.go:28, :59-61), while the renderer emits `| Agent ID | Workspace dir | Role | State |` (roster_render.go:73).

### H14 — The design lands on a system with zero production decks
`~/.parley/protocol/core/` does not exist on this machine; `find … -name protocol-lock.yaml -maxdepth 6` returns nothing anywhere under AI_WORKSPACE, including this repo's own deck. The core store is empty (IMPLEMENTATION.md:475-479). Mitigating: zero published releases means changing the release layout to add a registry costs nothing **today** and will not stay that way.

### H15 — Two hard-coded promises that must retire atomically
- `internal/app/protocol.go:211` — `"  (project-local content; the overlay that will carry it is ratified but not shipped)"`.
- `parley-deck/COOPERATION.md:767-768` — "the deck's own overlay, **once that ships**".
Both must change in the same commit the overlay ships, or the CLI and the protocol text lie in opposite directions.

### H16 — Test coverage is one layer up
`internal/protocolcore/render_test.go` is 45 lines covering only the LCS primitive (`TestLargeDocumentDoesNotExhaustMemory`, `TestLCSSeesOrder`). All ~35 behavioural tests live in `internal/app/protocol_test.go`. Adding an overlay parameter to `Render` breaks compilation in exactly two files and has almost no unit-level guard rail inside `protocolcore` itself.

### H17 — Fleet noise floor on first render
27/29 decks are missing the default's `Blast radius` block (default L745-769) purely from sync lag (default mtime Aug 7, deck mtimes Aug 6) — not an override. The live deck's stamp is in the legacy format (`parley-deck/COOPERATION.md:7`), which render.go:377-379 deliberately reports as not carried forward. So **every** deck's first real render produces a non-empty removal report before the overlay adds anything. And one deck (`ecb-meeting-2026.05`, 684L vs 1360L, 38 missing headings including all of §15, 43 hunks) never received the sync at all.

### H18 — Placeholder pollution would be frozen as "project data"
4 decks still carry `<workspace-name>`, 5 carry `<date> — created by parley init`, 2 carry `<agent-id-1>` roster rows, 7 carry `@<host-user>` handle rows. An overlay extraction that reads these as intentional local content freezes a defect class.

---

## 4. Genuinely open decisions

All either/or. Nothing here was ratified — where a reader says a scheme was proposed and killed, that is noted so nobody re-proposes a corpse.

### D-a — Overlay file syntax (never decided)
D4 fixes the path, the two operations and the required provenance fields; **no artifact ratifies how an operation is written on disk**. All three candidate syntaxes are dead: hermes-1's `## override §N` withdrawn (round-02/hermes-1.md:13), codex-1's `provide <target-id>` rejected and the rejection accepted (signoffs/codex-1.md:37-39), kimi-1's `Override §<anchor>` superseded by IDs. `00-prompt.md:77` made format choice an explicit non-goal and codex-1 kept serialization deliberately open (round-01/codex-1.md:266-268).
**Choose:** Markdown-with-fenced-directive-blocks / YAML-frontmatter-plus-Markdown-payloads / TOML sidecar. Pick one and write the grammar.

### D-b — Registry: ship it in this slice, or not?
**Either** the overlay slice also ships D2's block-ID registry (release layout changes; releases are write-once, but zero exist today so migration cost is currently nil — H7, H14) **or** it ships without one — in which case there is *no ratified third addressing option in the record*, because D2 forbids heading text and forbids inline markup.
**Sub-decision if yes:** what file, and how is a block's extent delimited in the Markdown — given `s6.6` is a list item, not a heading (H8)? Sub-heading/list-item granularity, byte offsets in the registry, or restructure the core so every addressable block IS a heading before the overlay ships.

### D-c — Where `ext-1` renders (three unratified proposals in conflict)
D3 says only "rendered at a declared position"; D5 step 4 says "append". Round 2 never converged:
- codex-1: "after §8 and before the TL;DR/reference appendices" (round-02/codex-1.md:46-48)
- kimi-1: "after the final core section" (round-02/kimi-1.md:57)
- hermes-1: "end of file" is non-deterministic; it must be a named block ID (round-02/hermes-1.md:73)

**And this decision is coupled to H2:** rendering `ext-1` at a fixed tail position produces the repositioning false-loss on every DF-4 deck. **Either** render the payload in place at the section's current position (kills the "one declared position" simplicity) **or** accept a fixed position and handle the false loss under D-g.

### D-d — Deck-namespaced extension ID format
D3 says "deck-namespaced IDs". kimi-1's `<deck-slug>.<slug>` (round-02/kimi-1.md:59-60) was **never ratified**. Decide the scheme and its collision/validation rules.

### D-e — Where D10's reconfirmation is recorded
codex-1 proposed a committed compatibility **receipt** keyed by core hash + overlay hash + resolver version, with the receipt hash in the lock (round-01/codex-1.md:113-119, round-02/codex-1.md:209). D8's ratified lock fields are core version + hash, overlay hash or `none`, resolver version, effective hash — **no receipt** (consensus.md:166-167). D10 requires "reconfirmation" without naming the artifact that proves it happened.
**Either** add the receipt file + lock field **or** define reconfirmation as a lock-field update alone.

### D-f — Authoring surface
No CLI verb for creating or editing an overlay was ratified; the shipped surface is `status | render | check | publish` (protocol.go:57-69, IMPLEMENTATION.md:37). **Either** read-only (`protocol overlay show|validate`, file hand-written) **or** a writer (`protocol overlay set|remove --dry-run --yes`, matching the `roster set` posture at roster.go:97-101). If a writer: is `--confirm-breaking` (roster.go:101) the right analogue for a replace that changes a rule? Note D11/R5.1 says the overlay is *not* user-only, so an attended gate is not implied.

### D-g — Loss report semantics with a second source
G1 requires reporting every block replaced or removed, and today's report already flags genuine project content (IMPLEMENTATION.md:467-470).
**Choose exactly one:** (i) exempt overlay-carried content from the loss report (risks re-creating the silent-erasure class nine review cycles bought — IMPLEMENTATION.md:409-416); (ii) suppress losses whose lines all reappear elsewhere in the render (weakens the order-sensitive sequence guarantee, render.go:193-202); (iii) accept a one-time noisy report with an explicit migration note. Constraint that must survive: "an empty report means 'no line disappeared', NOT 'no meaning was lost'" (render.go:207-216).
**Related:** does `RenderResult` grow an `Added` field and/or per-block provenance (core / identity / overlay)? Both call sites need new output if so.

### D-h — Missing or unreadable overlay file
**Either** block (like a missing pinned release, protocol.go:118-128) **or** render without it. Rendering without it silently erases the overlay's content from the committed deck and reports it as "not carried forward" — a literal repeat of the 2026-08-06 incident that motivated this whole design.

### D-i — Lock parsing: does it become a real YAML parse?
**Either** required-field validation + unknown-key refusal, so an old binary reading a new lock fails closed **or** leave the prefix scan and accept H6's silent-data-loss vector for anyone on a stale binary. Note this is entangled with rank 2 — see §5.

### D-j — Stamp content and its regex, in the same commit
**Either** the stamp stays core-only (render.go:95) **or** it carries core + overlay + effective hash (D7, consensus.md:141). Whichever: `generatedStampRe` (render.go:380) changes in the same commit or the exemption silently dies (H4). Separately: is `**Protocol synced:**` one field or an append-only log? The drift guard demands exactly one (drift_test.go:72-77); librade-algoTrader already stacks three.

### D-k — The §2 roster table: is it overlay content at all?
This is the highest-volume real case (11/29 bespoke schemas) and it is **not** covered by D3's open surface.
**Choose:** (a) overlay owns the rows; (b) overlay owns nothing and the table is always rendered from `agents.toml`; (c) overlay owns only the annotations *around* a rendered table. If (b)/(c), the bespoke columns (CLI/runtime, Model, Reasoning/effort, State, Display name, Backend, Invocation) need a home — likely `parley roster render` growing columns sourced from `agents.toml`, with the overlay reserved for what `agents.toml` genuinely cannot express (invocation gotchas, dated user directives, MANUAL-Bash caveats).
Hard constraint on all three: `agents.toml` is the authority (defaults/COOPERATION.md:102-105) and 23/29 tables already disagree with it — do not create a second competing roster surface.

### D-l — Is Transport identity, or core?
27/29 decks override it, yet §0 calls it sticky and requires a §7 idea to change (defaults/COOPERATION.md:54), and it is **deliberately excluded** from the drift allowlist. **Either** an overlay/identity field with a recorded change-reason **or** a protocol-governed value the overlay may not touch. Evidence for the tension: the servers deck improvised an 11-line HTML comment recording that the sync flipped transport without recording it, "which §0 forbids" (servers/parley-deck/COOPERATION.md:6).
**Note the reader conflict here.** The renderer reader lists Transport as one of the five preserved `IdentitySlots` fields (render.go:9-20); the real-decks reader states "Transport is deliberately NOT in it [the drift allowlist]" (drift_test.go:24-29 lists `workspacePrefix, createdPrefix, protocolSyncPrefix, rosterHeaderLine, handleHeaderLine`). Both are correct about different mechanisms — the renderer preserves it, the drift guard does not exempt it. That asymmetry is itself a decision to make explicit.

### D-m — "Six identity slots" is already a doc/impl mismatch
`IdentitySlots` declares **five** fields (Workspace, Created, Transport, RosterTable, HandleTable — render.go:9-20) while its own doc comment and D3 both say six; the sixth (`**Protocol synced:**`) is regenerated, not preserved. **Either** correct the count in D3's language **or** make the stamp a real slot. Also unresolved: `**Parley deck:** ./parley-deck/` (parley-deck/COOPERATION.md:4) is a genuine per-deck header value that is *not* a slot and would be silently replaced for any deck whose deck dir differs.

### D-n — How is Workspace derived?
16/29 decks do not match their project directory (11 literally say `parley-deck`, 4 are placeholders). And it is used for free text: design-mail L3 is `**Workspace:** \`design-mail\` (3 repos: design-mail-fe, design-mail-fe-dialogportal, design-mail-fatclient-backend-schema)`. **Either** the overlay stores it **or** it is derived from the project directory with the overlay supplying an optional label.

### D-o — Host-handle table under `local-dir`
13/29 empty, 9 filled with `n/a` / `local-dir` / `(local-dir)` noise, 7 still holding raw `@<host-user>` template rows, exactly one deck with real handles (the CLI dogfood, `['feci']`). **Either** suppress the table entirely when transport is `local-dir` **or** keep it and accept the noise.

### D-p — DF-4's content: `ext-1` payload, or core?
Three positions on record, unreconciled. consensus.md:250-251 ratifies restoring librade-algoTrader's destroyed section as an `ext-1` payload. codex-1 wanted it reviewed as a candidate **global** rule or a retired workaround rather than auto-converted (round-02/codex-1.md:285-287). kimi-1 wanted its substance ratified into the **core** (round-01/kimi-1.md:165-166). In practice DF-4 was closed by restoring raw deck text from `git show HEAD:` in **four** decks (IMPLEMENTATION.md:434-444), so the question is live again and larger than the ratified text.
**Coupled sub-question:** the four decks carry that content as HTML-comment roster blocks and header provenance lines (IMPLEMENTATION.md:440-444). D2 forbids inline markup in the **core** body — does that prohibition extend to an `ext-1` **payload**, or is the payload free-form Markdown?

### D-q — Provenance/justification as first-class overlay structure?
Three decks already improvise it with three different conventions: servers' HTML-comment TRANSPORT RECONCILIATION, auftra's `<!-- ROSTER UPDATE 2026-07-30 (user): "…" -->` with verbatim user quotes, librade's `> **2026-07-18 roster swap:**` blockquote. D4 requires a *rationale* per operation — **either** that field is the home for all of it **or** a separate dated change-log structure is needed.

### D-r — Does the CLI's own deck ever carry an overlay?
**Choose:** a sixth allowlisted zone in `drift_test.go`, an overlay-aware drift guard, or a standing rule that the source repo's deck never carries an overlay (H13).

### D-s — Does `protocol render` become the sole writer of COOPERATION.md?
**Either** `roster render` (roster_render.go:146) and `preflight.syncConsumerProtocol` (preflight.go:509) keep writing independently **or** they feed inputs to a single renderer. Four location logics over one file is the shape that produced the drift this idea exists to end (H12).

### D-t — Is prose-matched table addressing replaced now, or later?
**Either** table zones move to ID addressing in the same slice **or** the overlay ships onto a substrate that already violates D2 — where a core column rename silently empties every deck's roster (H9, probe-confirmed).

### D-u — Is user constraint 3 currently satisfied?
kimi-1's signoff judged D3's near-empty surface adequate "on condition that rank 3 actually ships this cycle (… if it slipped to deferred, constraint 3 would be unmet)" (signoffs/kimi-1.md:95-97). Rank 3 did not ship; neither the rev-2 nor rev-3 signoff revisits the condition (signoffs/rev3/kimi-1.md:15-53 covers only fixes 1-4). **Either** this run discharges that condition **or** it must be re-raised.
**Related unfinished business:** kimi-1's `parley protocol audit` surface (overlay count, targeted IDs, core-version spread) was dropped with no deferred-follow-up number; kimi-1's own signoff suggested it become DF-5 (signoffs/kimi-1.md:58-60) and **no DF-5 exists**.

### D-v — Never-synced decks
Does the overlay apply to `ecb-meeting-2026.05` (684L, 38 missing headings)? Extracting an overlay from it requires deciding whether its absent sections are project-local removals or pure staleness. The evidence says staleness; nothing in the file records that. **Either** exclude never-synced decks from overlay extraction until they are synced **or** define a staleness marker.

### D-w — Multi-replace ordering (future-proofing, decide cheaply)
D5 step 3 says "apply the replace operation" (singular) and D4 forbids a duplicate provider per ID, but multi-replace ordering is unspecified — relevant the moment a future core opens a second replaceable block. **Either** specify deterministic ordering now **or** state explicitly that v1 permits at most one replace.

### D-x — Line endings between three sources
Today the deck decides output convention (render.go:53, 100-102) and each source is normalized on its own line (render.go:54, :57). **Either** the deck keeps deciding and the overlay is normalized in **or** a global LF-normalization rule. Undecided = the exact non-convergence bug already fixed once for the core.

---

## 5. What is NOT in scope

**Not rank 2 (per-idea protocol pinning).** D7's manifest fields and D8's full five-field lock belong to rank 2 (consensus.md:141, 166-167; IMPLEMENTATION.md:51-53: "the remaining fields belong to rank 2"). The overlay slice touches the lock **only** to the minimum extent D-i decides, and must state which lock fields it writes. Do not build per-idea pin capture, run-manifest protocol recording, or effective-hash-per-idea here.

**Not DF-1 (the sandbox).** Out of scope entirely; it is a separately ratified follow-up, unbuilt (FINAL.md:52-54).

**Not DF-2 (fleet migration).** The overlay is a *prerequisite* for it, not part of it: "**Therefore DF-2 must not run until the overlay ships.**" (IMPLEMENTATION.md:471-473). No fleet-wide pass in this slice. Migration *tooling design* may be sketched only insofar as D-g requires knowing how reported-lost lines become overlay payloads.

**Not rank 4 (detection-layer enforcement).** Ratified and scheduled, not built here (FINAL.md:52-54).

**Not the core content itself.** No new sections, no rewriting §7 or §15 — both are sealed (consensus.md:101-102). The only permitted core-text change is retiring the two promissory notes in H15.

**Not a widened open surface.** Every wider proposal was retracted on the 1-in-36 evidence — claude-1's §12/§11-parameters/timeouts/§4.0-thresholds and kimi-1's §0/§6/§8/§11/§13/Appendix A (consensus.md:103-105, round-02/kimi-1.md:13-19, round-02/hermes-1.md:19). The 29-deck taxonomy independently confirms it: 28/29 decks never touched a line outside the header and §2. Do not reopen.

**Not a general patch language.** No JSON Patch, no unified diff, no line ranges, no heading-text patching, no conditionals, includes or parameterization; whole-block replacement only (round-01/codex-1.md:71-74, round-01/kimi-1.md:55-60).

**Not extension of sealed sections.** Extension happens only at `ext-1`, so no checker ever has to judge whether appended prose weakens a sealed rule (round-02/hermes-1.md:67, consensus.md:101-102 + D5 step 4). kimi-1's round-1 "purely additive local procedures under sealed sections" was struck.

**Not multiple overlay files, and not auto-created overlays.** One file only (round-01/kimi-1.md:57-58); no overlay at bootstrap or migration unless reviewed local content exists, and no auto-conversion of unknown prose (round-01/codex-1.md:149-151).

**Not `mergePreservingZones` as the renderer.** It survives only as migration zone-extraction scaffolding; its `## 3.` anchor does not survive (consensus.md:224-229).

**Not a second roster store.** Model and reasoning already have a home in `agents.toml` via `parley roster set` (defaults/COOPERATION.md:56); the overlay must not become a competing one (R7.4).

**Not widening G2.** The attended TTY gate stays scoped to `~/.parley/protocol/core/` (consensus.md:257-260). A deck-overlay change goes through a normal idea (D11).