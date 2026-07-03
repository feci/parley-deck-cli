---
agent: codex-1
idea: protocol-restructure-appendices
round: 1
date: 2026-07-03
---

## Summary

Recommend keep-numbers-relocate. Renumbering is incompatible with the stated "physical MOVE only" constraint because it would require rewriting every `§N`, `§11.B`, `§9.0`, and similar reference. The achievable improvement is progressive disclosure by order: core first, reference material last. The `core <=200 lines` target is explicitly not achievable in this idea because current §4 is about 505 lines by itself; that requires a separate §4 phase-split idea.

One important constraint issue: I would not add a new `## Appendices` banner in this round's implementation. A new heading is a new non-blank content line, so the required sorted non-blank-line diff would not be empty. If the group wants the banner, make that an explicitly allowed text addition or a follow-up idea.

## 1. Numbering

Pick keep-numbers-relocate.

Reasons:

- It preserves every existing reference literally: `§9.0`, `§11.B`, `§11.C`, `§12.11`, `§6.6`, and `Appendix A` do not need edits.
- It satisfies the "no rule text changed" constraint. Renumbering is not just layout; it changes visible protocol text and forces broad reference rewrites.
- It keeps consumer and skill documentation stable. A non-sequential reading order is acceptable because Quickstart already states that §0-§8 are core and §9/§11/§12/§13/§14 are reference appendices.
- It minimizes review surface. The intended diff becomes a block move, not a semantic protocol edit.

The trade-off is that readers see §10 before §9. That is less elegant than sequential numbering, but it is much safer than touching every reference in a protocol document whose whole purpose is durable cross-agent citation.

## 2. Final Section Order

Target order:

1. Title and metadata header
2. Quickstart
3. §0 Choose the transport
4. §1 Scope and purpose
5. §2 Active agents (roster)
6. §3 Directory layout
7. §4 Protocol - phases of an idea
8. §5 Quorum and async participation
9. §6 Conflict-avoidance mechanics
10. §7 Changing this protocol
11. §8 Inbox (lightweight channel)
12. §10 TL;DR
13. §9 Session-start checklist for every agent
14. §11 Transport mechanics
15. Appendix A - Adopting this protocol in a new project
16. §12 Pipeline blocks & action stages
17. §13 Retrospective optimization
18. §14 Automated outer loop (loop engineering) - the human brake

Placement rationale:

- §10 TL;DR belongs immediately after the core sections as a closing summary before the reference block.
- §9 becomes the first reference section after §10 because it is operational session-start material, not core conceptual flow.
- Appendix A stays in its current relative position between §11 and §12. It is bootstrap/reference material, and moving it elsewhere adds churn without improving the core path.
- No new `## Appendices` heading in this no-text-change version. The appendix/reference boundary is the transition after §10, already explained by the Quickstart map.

## 3. Mechanical Method And Verification

Implementation method:

1. Work from exact heading anchors, not line numbers. Split the file into byte slices at top-level `## ` headings, with the title/header plus Quickstart as the preamble block.
2. Assert every expected top-level block appears exactly once before writing: Quickstart, §0 through §14, and `## Appendix A`.
3. Reassemble only by permuting whole byte slices in the target order above. Do not normalize whitespace inside blocks. Do not edit references. Do not edit prose.
4. Apply the same permutation to both in-repo copies: `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`.
5. Preserve project-specific zones as they already differ: live deck has real workspace/header/roster data and `Protocol synced`; embedded default has bootstrap placeholders and empty §2 tables. The invariant is byte identity outside the allowed normalizations enforced by `TestEmbeddedDefaultMatchesLiveDeck`, not raw whole-file identity.

Verification gates:

- Per-block SHA check: for each named block in each file, compute SHA-256 before and after. The same block's bytes must be identical; only its position may change.
- Multiset non-blank-line check: strip blank lines, sort with `LC_ALL=C`, and diff old vs new. It must be empty for each file.
- Normalized copy drift check: run `go test ./internal/protocol -run TestEmbeddedDefaultMatchesLiveDeck` from the repo root. This guards live vs embedded default after the allowlisted zones are normalized.
- Init template check: run `go test ./internal/protocol -run TestDefaultCooperationForInit` so `parley init` still gets the expected embedded default behavior.
- Whole test slice if time allows: run `go test ./internal/protocol ./internal/app` because preflight owns consumer sync and protocol freshness behavior.
- Metadata restamp: after final content is accepted, update `meta/version.json` `protocolSha256` and, if the packaged/default body changed, `packagedProtocolSha256`. For this source deck, stale hashes are advisory, but consumers use them for freshness decisions.
- External fallback sync: re-sync any packaged/installed skill fallback `references/COOPERATION.md` from the accepted body and verify body identity according to that packaging's expected header policy.

## 4. Cross-Reference Audit

Use an automated audit before and after the move.

1. Extract references with a regex such as `§[0-9]+(?:\\.[0-9A-Z]+)*|Appendix [A-Z]`.
2. Extract targets from headings:
   - `## N. ...` defines `§N`.
   - `### N.X ...` defines targets such as `§4.0`, `§9.0`, `§11.B`, and `§12.11`.
   - `## Appendix A ...` defines `Appendix A`.
3. Also support ordered-list clause targets inside a base section, because the live document has `§6.6`, which resolves to item 6 under §6 rather than a Markdown heading.
4. Compare the before and after reference sets. They should be identical, because no prose changes are allowed.
5. For every reference, confirm the target exists after the move. For this document the known set includes `Appendix A`, `§0`, `§1`, `§2`, `§3`, `§4`, `§4.0`, `§5`, `§6`, `§6.6`, `§7`, `§8`, `§9`, `§9.0`, `§11`, `§11.B`, `§11.C`, `§12`, `§12.11`, `§13`, and `§14`.

This audit should fail on dangling references, but it should not require section numbers to appear in ascending order. Non-sequential order is the deliberate design.

## 5. Risks

- Drift guard anchors: `internal/protocol/drift_test.go` asserts exact single occurrences of `## 2. Active agents (roster)`, the roster table header, and the host-handle table header. It also expects table separators immediately after those headers. The move must not duplicate, rename, or reformat those lines.
- Consumer sync anchor: `preflight` preserves project-specific zones by cutting at the first line starting with `## 3.`. The reorganization keeps §3 in the core before appendices, so this should remain valid. Accidentally changing the §3 heading shape would be a real break.
- New banner risk: adding `## Appendices` breaks the sorted-line invariant and may also introduce a new top-level heading that simple section splitters did not expect. Avoid it here.
- `protocolSha256`: any reorder changes the file hash even with no rule text changes. `meta/version.json` must be restamped, or source decks will show freshness noise and consumer decks may make the wrong sync classification.
- Consumers: consumer decks auto-merge packaged protocol text from §3 onward while preserving header, §0, §1, and §2. A reorder is safe for project zones, but it creates a noisy downstream diff and may overwrite consumer-local edits outside the preserved zones, as designed.
- Raw byte identity wording: live and embedded defaults cannot be raw-byte identical because the embedded default intentionally has generic header lines, no `Protocol synced`, and empty §2 tables. The enforceable target is byte identity after the existing allowlist normalization plus body identity for any fallback snapshot that is supposed to be a direct body copy.
- Human line-number references: any out-of-band references to old line numbers become stale. Protocol-internal `§N` references survive; line-number links in PR comments, docs, or external notes do not.
- Sequential-order assumptions: any consumer that assumes `§N` appears before `§N+1` will be exposed by this move. The known implementation anchors are heading strings and freshness zones, but a focused grep for section-heading parsers should still be part of implementation review.
