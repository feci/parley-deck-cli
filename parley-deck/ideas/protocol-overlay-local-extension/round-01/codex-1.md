---
agent: codex-1
idea: protocol-overlay-local-extension
round: 1
date: 2026-08-07
---

## Summary

The overlay should not ship as a patch on the current renderer. This slice must add the release registry, a strict overlay and lock grammar, source-aware change reporting, and one composition path that is the sole writer of `COOPERATION.md`; otherwise the first overlay merely gives the existing silent-loss paths another input. The §2 tables are typed projections of `agents.toml`, never overlay content. The record did not ratify serialization, the exact `ext-1` position, roster-column preservation, or most of the remaining D-items, so the choices below are proposals rather than inherited decisions.

## Proposed approach

### Evidence discipline

The tags below identify the basis of factual evidence. They do not assign `CONFIRMED` or `WRONG` verdicts to claims I own; §15 distinguishes evidence from verdicts and forbids owner self-verdicts (`parley-deck/COOPERATION.md:1238-1265`).

- `PRIMARY` — The ratified record requires permanent IDs, no heading-text addressing, and no inline markers in the core (`parley-deck/ideas/meta-protocol-change-global-core-protocol/consensus.md:88-102`). It fixes the overlay path and the two operation kinds, but not their serialization (`parley-deck/ideas/meta-protocol-change-global-core-protocol/consensus.md:108-117`; `parley-deck/ideas/protocol-overlay-local-extension/00-prompt.md:69-77`).
- `PRIMARY` — The current release API reads and writes only `COOPERATION.md`; `CoreFileName` is singular and `Publish(version, body)` creates that one file (`internal/protocolcore/core.go:25-26`, `internal/protocolcore/core.go:74-108`, `internal/protocolcore/core.go:132-185`).
- `PRIMARY` — The current renderer discovers identity values and both tables by global prefix/prose matching (`internal/protocolcore/render.go:129-188`), while `roster render` separately finds §2 by a literal heading and writes the file itself (`internal/app/roster_render.go:24`, `internal/app/roster_render.go:104-157`).
- `PRIMARY` — The current loss result has only `Body`, `Removed`, and `Preserved`; removal is an order-sensitive LCS over exact lines, and its documented limit is: an empty report means no line disappeared, not that no meaning was lost (`internal/protocolcore/render.go:30-39`, `internal/protocolcore/render.go:191-217`, `internal/protocolcore/render.go:293-335`).
- `PRIMARY as to the supplied survey, not an independently rerun fleet scan` — The required brief reports that 0/29 decks are byte-identical, 27 modern decks diverge only in the header and §2, 11/29 use bespoke roster schemas, and one deck has a true extra section (`parley-deck/ideas/protocol-overlay-local-extension/00-scoping-brief.md:11-19`). I use those measurements as scoped input, not as a ratification.

### 1. D-a — choose strict YAML frontmatter with Markdown literal payloads

The record never decided the syntax. I choose **YAML-frontmatter-plus-Markdown-payloads**, with each payload held in a YAML literal scalar. There is no free-form body after the closing `---`; that avoids inventing a second delimiter language and permits arbitrary Markdown headings, fences, and HTML comments inside a payload without confusing operation boundaries.

Concrete v1 grammar:

```yaml
---
schema: parley.protocol-overlay/v1
core-version-range: ">=1.0.0 <2.0.0"
operations:
  - id: deck.working-language
    kind: replace
    target: s6.6
    expected-target-sha256: "<64 lowercase hex characters>"
    rationale: "Why this deck needs a different working-language rule."
    markdown: |-
      6. **Working language.** All protocol artifacts MUST be written in English.
  - id: deck.local-procedures
    kind: extend
    target: ext-1
    dependencies:
      s1: "<64 lowercase hex characters>"
      s2: "<64 lowercase hex characters>"
      # Every sealed dependency is listed; `all-sealed` is the default set.
    rationale: "Why these project-local procedures belong in this deck."
    markdown: |-
      ## Project-local procedures

      Additional procedure text.
---
```

Normative parser rules should be deliberately smaller than YAML in general: UTF-8; one document; exact known keys; no aliases, tags, merge keys, duplicate mapping keys, comments carrying semantics, or trailing body; strings remain strings; and unknown keys fail closed. `operations` contains one or two entries in v1, at most one `replace` and at most one `extend`; the resolver ignores file order and always composes replace before extend. Every operation requires a unique `id`, non-empty `rationale`, and non-empty `markdown`; therefore an empty or zero-operation overlay is invalid.

Use `deck.<slug>` operation IDs, matched by `^deck\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$` and unique within the deck. The file itself is the deck namespace; deriving a namespace from `Workspace`, a remote URL, or a directory would make identity depend on values the evidence shows are unstable. Git history plus the normal idea is the change log; D4's required per-operation `rationale` is sufficient provenance, so D-q does not need a second dated log structure.

### 2. D-b — the registry is part of this slice

Yes. Shipping an overlay without the registry would violate the ratified addressing rule rather than implement a smaller version of it. There is no legitimate third addressing mechanism in the record, and the current release layout has no registry (`PRIMARY`: `parley-deck/ideas/meta-protocol-change-global-core-protocol/consensus.md:81-93`; `internal/protocolcore/core.go:25-26`, `internal/protocolcore/core.go:132-185`). The brief reports that the core store is still empty, so this is the last cost-free opportunity to change the release layout before write-once releases exist (`parley-deck/ideas/protocol-overlay-local-extension/00-scoping-brief.md:128-130`).

Each release should contain `COOPERATION.md` and `registry.yaml`. `registry.yaml` records the body hash, permanent IDs, policy, exact UTF-8 byte spans as half-open `[start,end)` offsets, block hashes, tombstones, the zero-width `ext-1` insertion offset, and separate typed identity-slot spans. The publisher validates: body hash; offsets on UTF-8 boundaries; in-bounds spans; unique IDs; containment without crossing ranges; block hashes; exactly one `ext-1`; and that every non-slot byte is sealed by default. The release identity must cover both files, not merely the Markdown.

Byte spans are the right v1 extent mechanism. They are release-local coordinates over immutable, hashed bytes; permanent IDs provide cross-release identity. They neither depend on headings nor add markers to the core. This also handles `s6.6`, which is currently a list item rather than a heading (`PRIMARY`: `parley-deck/COOPERATION.md:732-743`). **Exemption witness for H8:** heading-delimited segmentation fails because `s6.6` has no heading; this design's necessary preconditions are a verified body hash plus a validated byte span, so the absent-heading precondition is irrelevant. A test must rename/renumber headings without changing registry IDs and still resolve the same blocks.

The registry should also identify the five caller-supplied slots (Workspace, Created, Transport, roster table body, host-handle table body) and the generated sync-stamp position. The ratified phrase “six identity slots” should be corrected to **five typed inputs plus one generated provenance field**: the implementation already declares five fields and regenerates the sixth (`PRIMARY`: `internal/protocolcore/render.go:9-20`, `internal/protocolcore/render.go:94-99`). Slot addressing must therefore stop using prose.

### 3. D-k — §2 is not overlay content

Choose D-k(b): the overlay owns neither §2 rows nor its schema. `agents.toml` remains the roster authority, and both tables are projections supplied as typed renderer inputs. The live protocol explicitly says the table is a generated, non-authoritative view and hand-edits are overwritten (`PRIMARY`: `parley-deck/COOPERATION.md:103-122`). Allowing an overlay to own any table bytes would create the competing roster surface the kickoff forbids.

The 11 bespoke schemas are real migration requirements, but they are evidence for a richer **roster projection**, not for protocol override. Open follow-up `roster-projection-schema`: add explicitly supported columns sourced from typed `agents.toml` fields (model, effort, state, host handle, workspace, role, and any justified display-only fields). Invocation caveats and dated directives that are not roster data may become an `ext-1` procedure after human review. DF-2 must not migrate an affected deck until that follow-up can represent its intentional columns or the operator explicitly retires them.

This slice must still eliminate prose-matched table addressing. `roster render` may remain as a user-facing command, but it must delegate to the same compositor as `protocol render`; it may not splice §2 independently. That is the only credible repair for H9 and H12.

### 4. D-c / D-g — registry position and source-aware loss reporting

The record never fixed the exact position: D3 says “declared position,” while D5 says to append the payload (`PRIMARY`: `parley-deck/ideas/meta-protocol-change-global-core-protocol/consensus.md:95-100`, `parley-deck/ideas/meta-protocol-change-global-core-protocol/consensus.md:119-128`). Set the v1 `ext-1` registry offset immediately after §8 and before `## 10. TL;DR` and the reference appendices. “Append” then means append the single payload at that registered extension point, not “guess EOF.” Runtime placement uses the registry offset, never the heading text.

For D-g, choose option (ii), but reject its unsafe global-line formulation. Keep the order-sensitive LCS as the baseline, then reclassify a removed contiguous run as `relocated-to-overlay` only when there is a strict witness: the run is byte-identical to one complete overlay payload and occurs uniquely in both the prior deck and the effective output. Ambiguity stays `removed`. This is not a multiset exemption and does not forgive arbitrary text merely because similar lines occur elsewhere.

Replace `RenderResult.Removed/Preserved` with a typed change report containing `removed`, `relocated`, `replaced`, `added`, and `preserved`, each with source (`core`, `identity`, or `overlay`), block/operation ID, before/after hash, and line count. The first adoption of `auftra` can then report its uniquely matched local block as `relocated`, with zero disappeared lines, while still making the movement prominent. A steady-state second render is byte-identical and empty. The invariant remains exact: **an empty report proves only that no line disappeared; it does not prove semantic equivalence**. A non-empty `relocated` event is not a loss, but it is still a review prompt because movement can change Markdown meaning.

The synced stamp must name both sources in the same slice, for example:

```text
**Protocol synced:** core 1.0.0 (<release-hash-12>); overlay <none|overlay-hash-12>; resolver overlay-v1
```

The brief's D-j either/or incorrectly couples rank 3 to rank 2's effective hash. I would not write an effective hash into this slice's stamp before rank 2 defines and stores it. The existing core-only stamp would misidentify two different effective inputs, so retaining it is also unacceptable. Remove the regex-based forgiveness and classify the generated stamp by its registry slot; do not merely expand `generatedStampRe` (`PRIMARY`: `internal/protocolcore/render.go:370-395`). Rank 2 can add the effective hash atomically with its lock semantics.

### 5. Composition and state model

Implement one pure compositor:

```text
Compose(verified release body + registry,
        strict lock,
        optional verified overlay,
        five typed identity inputs) -> effective bytes + typed change report
```

The caller reads files; the pure package does not. The current pure-function boundary is an explicit design property (`PRIMARY`: `internal/protocolcore/render.go:41-57`), but its signature must stop deriving typed inputs from arbitrary prior Markdown.

Use a strict v2 lock in this slice:

```yaml
schema: parley.protocol-lock/v2
core:
  version: 1.0.0
  release-sha256: "<hash covering body and registry>"
overlay:
  sha256: none
  compatibility-receipt-sha256: none
resolver-version: overlay-v1
identity:
  workspace: parley-deck
  created: 2026-05-09
  transport: github-pr
```

There must be no top-level `core-version:` compatibility alias: the old prefix scanner would ignore every new field but continue rendering (`PRIMARY`: `internal/app/protocol.go:81-101`). With the nested v2 shape, an old binary finds no pin and fails closed. The new parser requires every key, rejects unknown keys, and verifies hashes before composition. This is the minimum lock work rank 3 needs; per-idea capture and snapshots remain rank 2.

Overlay state is then unambiguous:

- lock says `sha256: none` and file absent: valid no-customization state;
- lock names a hash and the file is absent, unreadable, empty, or mismatched: block;
- file exists while the lock says `none`: block and require explicit validation/adoption;
- valid file plus matching lock: compose it.

D10 reconfirmation needs a committed receipt, not an unexplained lock edit. Use `parley-deck/meta/protocol-overlay-receipt.yaml`, containing core-release hash, overlay hash, resolver version, changed dependency IDs/hashes, decision (`accepted` or `rejected`), and the canonical idea/consensus path. Hash that receipt in the lock. A lock-field update alone proves only that bytes changed, not that anyone reviewed the dependency report.

`protocol render` owns the only write service for `COOPERATION.md`. `roster render` and preflight may update their authorities or request a render, but cannot maintain separate splice algorithms. Today all three write paths are visible directly (`PRIMARY`: `internal/app/protocol.go:180-238`; `internal/app/roster_render.go:104-151`; `internal/app/preflight.go:490-539`). Workspace remains an explicit label in the lock, not a directory-derived value; Transport remains a typed, overlay-inaccessible identity value governed by the existing sticky-transport rule; roster and host-handle projections come from `agents.toml`. For CRLF, normalize each source independently to LF, compose/hash in LF, then restore the prior deck's convention; a new deck defaults to LF, matching the existing per-source rule (`PRIMARY`: `internal/protocolcore/render.go:51-57`, `internal/protocolcore/render.go:99-103`).

### 6. Disqualifying hazards

I judge H1-H13, H15, and H16 disqualifying if left in the shipping path. They are not all equally severe, but every one is either an integrity hole or a release-truth/test gate:

1. **Immediate integrity blockers: H6-H12.** A stale binary ignores new lock fields; no registry or `s6.6` extent exists; table and identity addressing are prose/first-prefix based; stamp filtering is asymmetric; and independent writers can overwrite one another. The direct code is at `internal/app/protocol.go:81-101`, `internal/protocolcore/core.go:25-26`, `internal/protocolcore/render.go:65-89`, `internal/protocolcore/render.go:129-188`, `internal/app/roster_render.go:155-201`, and `internal/app/preflight.go:501-539` (`PRIMARY`). The overlay may **not** ship on this substrate.
2. **Observability blockers: H1-H5.** `check` collapses all drift into `hand-edited-or-stale`, additions have no result field, moves look removed, and the stamp is core-only with a format-locked exemption (`PRIMARY`: `internal/app/protocol.go:241-289`; `internal/protocolcore/render.go:30-39`, `internal/protocolcore/render.go:94-99`, `internal/protocolcore/render.go:191-217`, `internal/protocolcore/render.go:370-395`). Typed state and change events are part of the overlay slice, not polish.
3. **Repository and truth blockers: H13, H15, H16.** The current drift guard byte-compares outside five prose-defined zones and would reject a real overlay (`PRIMARY`: `internal/protocol/drift_test.go:14-30`, `internal/protocol/drift_test.go:43-101`). Replace it with an overlay-aware expected-render test, not a sixth normalization escape hatch. Retire both “overlay not shipped” promises atomically (`internal/app/protocol.go:211`; `parley-deck/COOPERATION.md:767-768`). Add unit tests around registry parsing/composition as well as G7b production-dispatch tests; the core package currently has only two tests, both for LCS behavior (`internal/protocolcore/render_test.go:8-45`).

H14 is not a blocker; it is the reason to change the release format now. H17 and H18 block unattended fleet extraction, not the overlay resolver: first-render noise and placeholder classification belong to DF-2's attended migration. They must not be used to auto-create overlay operations.

### 7. Ranking the remaining open decisions

The five kickoff questions settle D-a, D-b, D-c, D-g, and D-k above. I rank every remaining item as follows.

#### Must be settled and implemented in this slice, in order

1. **D-i, strict v2 lock parsing:** required to make an old binary fail closed rather than ignore the overlay.
2. **D-h, missing/unreadable state:** use the four-state lock/file rules above; missing expected bytes block.
3. **D-s, sole writer:** all commands feed one compositor/write service.
4. **D-t, table addressing:** registry slot IDs now; no prose matching remains in a shipping path.
5. **D-e, reconfirmation receipt:** committed receipt plus receipt hash in the lock.
6. **D-j, stamp:** core + overlay + resolver now, with registry-slot classification; effective hash deferred only with rank 2.
7. **D-m, slot count:** five typed inputs plus one generated stamp; `**Parley deck:**` stays the fixed `./parley-deck/` path, not a slot.
8. **D-l, Transport:** typed identity governed by §0, never an overlay operation.
9. **D-n, Workspace:** explicit label in lock identity, not directory-derived.
10. **D-d, extension IDs:** deck-local `deck.<slug>` grammar and uniqueness above.
11. **D-q, provenance structure:** required operation rationale plus normal idea/git history; no second change log.
12. **D-w, ordering:** v1 permits at most one replace and one extend; canonical order is replace then extend.
13. **D-x, line endings:** normalize all inputs to LF internally and restore the deck convention.
14. **D-r, this repo's deck:** allow an overlay and make the drift guard render-aware; a standing “never use it here” exception would leave the feature un-dogfooded.
15. **D-f, authoring surface:** ship `protocol overlay show|validate` only; the file remains human-authored through a normal idea. Defer mutation commands to follow-up `protocol-overlay-authoring-writer` after the grammar is proven.
16. **D-u, constraint status:** the owner's constraint is **not discharged today**. It becomes discharged only when this slice, including the registry and blockers above, lands with tests. Open `protocol-overlay-audit` as DF-5 for the dropped fleet audit surface; an audit command is not required to make one deck's overlay safe.

#### May be deferred, with named follow-ups

1. **D-o → `transport-aware-host-handle-rendering`.** For v1, keep the core table and render zero rows when `agents.toml` has no handles; never preserve placeholder rows. Decide whole-table suppression for `local-dir` separately because D3 currently opens only the body slot.
2. **D-p → `df4-local-content-disposition`.** Before DF-2, review each restored block as core candidate, retired workaround, or `ext-1` payload. Overlay payloads are free-form Markdown values, so HTML comments are representable; D2's no-inline-marker rule applies to core addressing, not to payload prose. Representability is not approval.
3. **D-v → `protocol-overlay-fleet-migration` (DF-2).** Exclude never-synced decks from automatic extraction. Their absences need attended classification; do not invent a staleness marker in the resolver.

The related bespoke-column work from D-k is separately named `roster-projection-schema` and gates DF-2 for affected decks. Rank 2 remains responsible for per-idea effective hashes/snapshots; DF-1 and rank 4 remain untouched, matching the scope boundary in `00-scoping-brief.md:240-250`.

## Concerns / open questions

The largest unresolved semantic limit is not a serialization question: a machine can prove that an extension only occupies `ext-1` and cannot replace sealed bytes, but it cannot generally prove that arbitrary English prose does not contradict a core rule. The effective core must state that local extensions are subordinate and cannot waive core obligations; contradictory prose is incompatible and must be rejected through the normal idea review. The implementation must not claim semantic contradiction detection unless a later design provides a real witness.

The registry publisher also needs a precise authoring contract. This analysis chooses the frozen on-disk format and byte-span semantics, but not whether release authors hand-supply offsets or use a generator that compiles a source map. The release must contain only validated compiled coordinates; a generator is preferable, but its source format is tooling detail so long as publish independently verifies all spans and hashes.

The source-aware relocation rule is intentionally narrow. If a prior block is duplicated, partially edited, or interleaved with other prose, it stays in `removed` and the first adoption remains noisy. That is safer than declaring uncertain content preserved.

## Risks

- A permissive YAML library can silently defeat the strict grammar through duplicate keys, aliases, implicit types, or unknown fields. Parser conformance tests must exercise each refusal through the real `Run` dispatch, not only a helper.
- Byte offsets are safe only when tied to exact normalized bytes and independently validated at publish/load. Applying offsets after identity substitution or CRLF restoration would target the wrong range.
- A unified writer is a larger refactor than adding an `overlay` argument to `Render`, but retaining independent writers leaves a reproducible lost-update path and should block release.
- The compatibility receipt adds one artifact and lock field, but omitting it reduces “reconfirmed” to an unauditable assertion.
- Free-form extensions remain human-reviewed policy. Mechanical non-weakening comes from sealed-byte immutability, fixed placement, and core precedence—not from pretending the resolver understands English.
