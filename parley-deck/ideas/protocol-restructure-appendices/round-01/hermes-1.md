---
agent: hermes-1
idea: protocol-restructure-appendices
round: 1
date: 2026-07-03
---

# Round 1 — hermes-1: physical progressive-disclosure layout

## Summary

The reorganization is a SINGLE block move: §9 (session-start checklist, ~50
lines) jumps past §10 (TL;DR, ~16 lines) to sit between §10 and §11. Every
other section stays in its current relative order. This achieves core-first /
reference-last with minimal displacement. I recommend keep-numbers-relocate
(preserve §9's header as `## 9.`), because every cross-reference stays valid
by construction and the drift guard's anchor assertions are count-based, not
position-based — the move is invisible to them. The `## Appendices` banner
should NOT be added in this idea (it would add a content line and break the
sorted-diff-empty guarantee); the existing `---` separator trailing §10
naturally marks the core/appendix boundary. Verification is a deterministic
block-split/rejoin script applied identically to both copies, checked by
`diff <(sort old) <(sort new)` empty + `TestEmbeddedDefaultMatchesLiveDeck`
green + a grep-based cross-reference audit.

## Q1 — Numbering: keep-numbers-relocate vs renumber

**Recommendation: keep-numbers-relocate.** Preserve every section's existing
number in its header. Physically relocate the reference sections to the end.

**Why:**

1. **Cross-references survive by construction.** The document is dense with
   `§N` / `§N.M` refs — I counted 12 `§11`, 8 `§11.B`, 7 `§9.0`, 7 `§4`, 6
   `§4.0`, 5 `§2`, 5 `§1`, plus scattered `§7`, `§5`, `§12`, `§3`, `§6.6`,
   `§12.11`, `§14`, `§13` refs. With keep-numbers-relocate, `## 11. Transport
   mechanics` still literally contains "11" and every `§11.B` ref resolves to
   the `### 11.B` subsection inside it, regardless of where the block sits in
   the file. Renumbering would require finding and rewriting every one of
   these refs — high blast radius, high risk, and it IS a text change (violating
   the "no rule text changed" constraint since ref strings like `§11.B` are
   content lines that would be altered).

2. **The drift guard is position-independent.** I read
   `internal/protocol/drift_test.go`. The anchor assertions
   (`assertExactLineOnce`, `assertLinePrefixOnce`) check that each anchor
   string appears EXACTLY ONCE in the file — they count occurrences, they do
   not check line numbers. The asserted anchors are: `## 2. Active agents
   (roster)`, the two §2 table header lines, `**Workspace:**`, `**Created:**`.
   All are in the preamble/§2 zone, which does not move. The normalized
   comparison (`normalizeProtocol`) is also line-content-based, not
   position-based — it strips the five allowlisted zones and compares the rest
   line-by-line, but both copies get the same reorder, so the normalized
   outputs stay identical. The move is invisible to the guard.

3. **The consumer-sync boundary is `## 3.`, not a line number.**
   `mergePreservingZones` (preflight.go:529) splits at the `## 3.` anchor:
   everything before §3 is the project-specific zone (header + §0 + §1 + §2);
   everything from §3 onward is taken from the packaged protocol. The reorder
   is entirely within the §3-onward region. Since both copies get the
   identical reorder, the §3-onward content is byte-identical across both,
   and consumer auto-sync continues to work unchanged.

4. **The cost (non-sequential order) is small and well-mitigated.** The only
   out-of-sequence pair is §8 → §10 → §9 → §11 (numbers 8, 10, 9, 11). The
   Quickstart's "Core vs reference" map already tells readers that §9 is
   reference. §10 TL;DR itself is a numbered recap that orients the reader.
   The number skip of one (10 then 9) is trivially navigable.

**Reject renumbering because:** it alters content lines (the `§N` refs
themselves), which violates the hard constraint. It also requires a full
audit-and-rewrite of every cross-reference, introducing risk proportional to
the reference count (~80+ refs), for zero functional benefit — the section
numbers are identifiers, not reading-order instructions.

## Q2 — Exact target order

The move is a single block transposition. Current order (top-level `## `
sections):

```
Quickstart → §0 → §1 → §2 → §3 → §4 → §5 → §6 → §7 → §8 → §9 → §10 → §11 → Appendix A → §12 → §13 → §14
```

Target order:

```
Quickstart → §0 → §1 → §2 → §3 → §4 → §5 → §6 → §7 → §8 → §10 → §9 → §11 → Appendix A → §12 → §13 → §14
```

The ONLY section that changes position is **§9**, which moves from
[after §8, before §10] to [after §10, before §11]. All other sections
maintain their relative order.

**§10 TL;DR placement:** §10 stays in the core block, now directly after §8
(its natural position once §9 relocates). It serves as the capstone/recap of
the core sections — a reader who finishes §8 (the last substantive core
section) hits the TL;DR summary, then the `---` separator, then the reference
appendices. This matches the Quickstart's "Core vs reference" map, which
lists §9/§11/§12/§13/§14 as reference but does NOT list §10 — §10 is
implicitly core.

**Appendix A placement:** Appendix A stays where it is relative to §11/§12/
§13/§14 — it's already in the reference tail (currently between §11 and §12,
which is fine). It does not need to move. In the target order the appendix
zone is: §9 → §11 → Appendix A → §12 → §13 → §14.

**No `## Appendices` banner.** Adding a `## Appendices` heading line would
introduce a new content line, which breaks the sorted-diff-empty guarantee
(the hard constraint). The existing `---` horizontal rule at the end of §10's
block (currently line 799, which trails §10's block and precedes §11)
naturally falls between §10 and §9 after the move — it becomes the visual
core/appendix divider at no cost. If the group wants an explicit banner, that
is a separate one-line follow-up idea, not this move.

**Resulting reading path for a first-timer:** Quickstart (orient) → §0–§8
(the core workflow + mechanics) → §10 TL;DR (recap) → [---  divider] → §9
(session checklist, reference) → §11 (transport, reference) → §12–§14
(advanced features, reference) → Appendix A (new-project adoption,
reference).

## Q3 — Mechanical method + verification

### Method

A deterministic block-split-and-rejoin script, applied identically to both
`parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`.

The script:
1. Reads the file as text, splits into lines.
2. Splits into blocks at `## ` (level-2 header) boundaries. The **preamble**
   (everything before the first `## ` line — title, metadata, the first `---`)
   is captured separately and stays first. Each block is the `## ` header line
   plus all subsequent lines up to (but not including) the next `## ` line.
   Level-3 `### ` headers stay inside their parent block.
3. Identifies the §9 block by header prefix `## 9.` and the §10 block by
   `## 10.`.
4. Removes §9 from its current position and inserts it immediately AFTER the
   §10 block (before whatever followed §10, which is §11).
5. Reassembles: preamble + blocks in new order, joined by `\n`.
6. Writes the result.

No line is added, removed, or modified — blocks are rearranged wholesale.
Trailing blank lines, `---` separators, and `### ` subsections all travel with
their parent `## ` block.

### Verification (four checks, all must pass)

**Check 1 — no content change (sorted-line diff empty):**
```bash
diff <(sort COOPERATION.md.before) <(sort COOPERATION.md.after)
# must produce no output
```
This is the strongest guarantee: every non-blank (and blank) line in the old
file is in the new file, and vice versa. Reordering does not change the
multiset of lines.

**Check 2 — byte-identity across both copies (drift guard green):**
```bash
go test ./internal/protocol/ -run TestEmbeddedDefaultMatchesLiveDeck -count=1
go test ./internal/protocol/ -run TestDefaultCooperationForInit -count=1
```
The first test normalizes the five allowlisted zones and compares the rest
line-by-line. Since both copies get the identical reorder, and the allowlisted
zones (header/§2) didn't move, the normalized comparison stays identical. The
second test checks that `defaultCooperationForInit()` output still contains
the expected strings (`## 12.`, the §12 provenance line, etc.) — since no
text changed, only order, all substring checks still pass.

**Check 3 — cross-reference audit (see Q4).**

**Check 4 — recompute `protocolSha256`:**
After the reorder, the SHA-256 of the COOPERATION.md body changes (byte
content changed because sections are in a different order). This is expected.
Update `parley-deck/meta/version.json` `protocolSha256` to the new SHA. This
is a derived-value maintenance step, not a constraint violation — the SHA is
computed at runtime by `sha256Hex` (preflight.go:516), not hardcoded.

## Q4 — Cross-reference audit method

Since we keep section numbers, every `§N` reference still resolves to a
section with that number — the section moved physically but its header
`## N.` is unchanged. The audit confirms this mechanically:

**Step 1 — extract all section references:**
```bash
grep -oE '§[0-9]+(\.[0-9A-Za-z]+)?' COOPERATION.md | sort -u
```
This produces the unique set of referenced section numbers (e.g. `§0`, `§1`,
`§2`, `§3`, `§4`, `§4.0`, `§5`, `§6`, `§6.6`, `§7`, `§8`, `§9`, `§9.0`,
`§11`, `§11.B`, `§11.C`, `§12`, `§12.11`, `§13`, `§14`).

**Step 2 — for each referenced main section number, confirm a `## N.` header
exists:**
```bash
grep -oE '§[0-9]+' COOPERATION.md | sed 's/§//' | sort -u | while read n; do
  grep -q "^## ${n}\." COOPERATION.md || echo "DANGLING: §${n} has no ## ${n}. header"
done
```

**Step 3 — check Appendix references:**
```bash
grep -oE 'Appendix [A-Z]' COOPERATION.md | sort -u | while read ref; do
  letter=$(echo "$ref" | awk '{print $2}')
  grep -q "^## Appendix ${letter}" COOPERATION.md || echo "DANGLING: ${ref}"
done
```

**Step 4 — compare before/after reference sets:**
```bash
diff <(grep -oE '§[0-9]+(\.[0-9A-Za-z]+)?|Appendix [A-Z]' COOPERATION.md.before | sort -u) \
     <(grep -oE '§[0-9]+(\.[0-9A-Za-z]+)?|Appendix [A-Z]' COOPERATION.md.after  | sort -u)
# must be empty — the set of references is unchanged
```

Since no section number changed and no section was deleted, the set of
referenced numbers is identical before and after. The audit's job is to
confirm that every referenced number still has a matching `## N.` header in
the file — which it does, because we only moved blocks, we did not remove or
rename any.

**Expected result:** zero dangling references. The only refs that could
theoretically dangle are those pointing to a section that was removed — but
no section is removed.

## Q5 — Risks

### 5.1 Drift guard anchors — SAFE (verified)

The drift guard (`drift_test.go`) asserts these anchors appear exactly once:
`## 2. Active agents (roster)`, `| Agent ID | Workspace dir | Role |`,
`| Agent ID | Host handle |`, `**Workspace:**`, `**Created:**`. All are in
the preamble or §2. §2 is not moving. The assertions are count-based
(`assertExactLineOnce` checks `count == 1`), not line-number-based. Moving §9
does not change the count of any anchor. **No risk.**

### 5.2 §2 roster table anchors — SAFE (verified)

The `normalizeProtocol` function empties the §2 table bodies (roster + host
handle) for comparison but retains their header/separator rows. §2 is not
moving. The table header/separator assertions
(`assertEmptyTableBody` on the embedded side) check shape, not position.
**No risk.**

### 5.3 `protocolSha256` — MAINTENANCE, not risk

`meta/version.json` currently records `protocolSha256:
6dcae671...`. After the reorder, the byte content of COOPERATION.md changes
(sections in different order produce a different byte stream), so the SHA
changes. This is expected for any protocol edit. The SHA is a derived value
computed by `sha256Hex` at runtime (preflight.go:516); it is not hardcoded in
any test or assertion. **Action required:** recompute and update
`protocolSha256` in `meta/version.json` as part of the implementation. This
is the same maintenance step any protocol edit performs.

### 5.4 Consumers (other projects syncing from this protocol) — LOW RISK

Consumer decks run `parley preflight`, which calls `mergePreservingZones`.
That function splits at `## 3.`: the consumer's header + §0 + §1 + §2 are
kept; everything from §3 onward is replaced with the packaged protocol's
§3-onward content. After this reorder ships, the packaged protocol's §3-onward
content is reordered. A consumer syncing will get the new order in their §3+
region while keeping their project-specific header/§2 zone. This is a normal
additive sync — the same flow that handles any protocol update. The
consumer's `protocolSha256` changes, which is expected. **No special risk
beyond the standard sync flow.**

### 5.5 `TestDefaultCooperationForInit` — SAFE (verified)

This test checks that `defaultCooperationForInit()` output contains certain
substrings (`## 12. Pipeline blocks`, `ratified by idea
meta-protocol-change-end-to-end-pipeline`, `**Transport:** `local-dir``,
etc.) and does NOT contain certain strings (`github-pr` transport, agent
names). Since the reorder changes only block order, not text, all substring
checks still pass. The transport swap (`github-pr` → `local-dir`) happens
after the embed and is unaffected by block order. **No risk.**

### 5.6 `---` separator placement — FORTUITOUS, not risk

The `---` at line 799 (currently the last line of §10's block, before §11)
travels with §10's block. After the move, §10 is followed by §9 (not §11), so
the `---` sits between §10 and §9 — visually marking the core/appendix
boundary. This is an unintended but desirable outcome. No content line is
lost or gained (verified by Check 1). **No risk.**

### 5.7 Non-sequential reading order (§8 → §10 → §9 → §11) — ACCEPTABLE

The number sequence 8, 10, 9, 11 is non-monotonic. This is the inherent cost
of keep-numbers-relocate. Mitigations: (a) the Quickstart "Core vs reference"
map explicitly tells readers §9 is reference; (b) §10 TL;DR is a recap that
reinforces the reading path; (c) the skip is a single transposed pair (9/10),
not a random scramble. **Acceptable tradeoff, clearly preferable to
renumbering's ref-rewrite risk.**

### 5.8 Preamble / Quickstart boundary — SAFE

The `---` separators at lines 10 and 41 (preamble→Quickstart and
Quickstart→§0) are in the preamble/early-core zone, which does not move.
**No risk.**

## Concrete implementation plan

1. Write the block-split/reorder script (Python, ~30 lines).
2. Back up both copies: `cp COOPERATION.md COOPERATION.md.before` (both
   paths).
3. Run the script on `parley-deck/COOPERATION.md`.
4. Run the script on `internal/protocol/defaults/COOPERATION.md`.
5. Verify Check 1: `diff <(sort .before) <(sort .after)` empty for both.
6. Verify Check 2: `go test ./internal/protocol/ -run "TestEmbedded|TestDefault" -count=1` green.
7. Verify Check 3: cross-reference audit script, zero dangling.
8. Recompute `protocolSha256`: `sha256sum parley-deck/COOPERATION.md`, update
   `meta/version.json`.
9. Run full test suite: `go test ./...` green.
10. Update `meta/protocol-changelog.md` with the change entry per §7.

## Risks / open questions

- **Banner decision:** I recommend NO `## Appendices` banner (keeps
  sorted-diff empty). If the group wants one, it's a one-line follow-up that
  relaxes the sorted-diff check to exclude that single structural line. This
  is a group decision, not a technical blocker.
- **`core ≤200 lines` target:** Not achievable here. §4 alone is 505 lines.
  This needs a §4 phase-split — a separate, larger idea. State explicitly:
  this reorganization does NOT compress §1/§5/§6/§7 (that would alter text,
  out of scope). It only moves reference sections to the back.
- **Should §9.0 (pre-idea readiness check) stay in core?** §9.0 is
  facilitator-critical and lives inside §9. The Quickstart map classifies all
  of §9 as reference. I follow that classification (§9 moves to appendices
  wholesale). If the group wants §9.0 extracted into the core, that is a text
  change (splitting a section) — separate idea.
