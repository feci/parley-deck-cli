---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
round: 1
date: 2026-09-05
lens: phase-scoped packet and instruction-layer correctness
---

## Summary

The packet slice of `meta-protocol-change-phase-packet-and-fixup-budget` is `status: not started`
(`IMPLEMENTATION.md:82-84`); no packet code exists in the tree, and the prompt builders embed rule
bullets inline without ever naming the protocol. The budget slice is enforced only on the driver
path; manual facilitation and resume have no ledger. I propose one disk-derived budget ledger shared
by driver and manual paths, a heading-anchored, hash-bound packet renderer with a human-authored
applicability manifest under §7, prompt-builder wiring that never reads the protocol itself, and
shadow-mode packets during the 12-task pilot so the packet never confounds it. Verification of
my slice must come from a non-implementer per §15.1.

## What I executed vs. inspected

**Executed (PRIMARY):** `shasum -a 256 parley-deck/COOPERATION.md` → `e2319aff8948130c…`, equal to
`protocolSha256` in `parley-deck/meta/version.json`; `wc -l` → 1373 lines; `grep -c MUST` → 34
lines (line count, not occurrences). Grep runs for `packet`, `MaxFixupCycles|CapCrossReviewRounds|
CrossReviewRounds`, `COOPERATION.md` and `func .*Prompt` over `internal/`.

**Inspected only:** the FINAL and IMPLEMENTATION of the phase-packet idea; `internal/track/track.go:130-199`;
`internal/driver/driver.go:60-84`; `internal/driver/impl.go:290-307`; `internal/runner/runner.go:821-874`;
`internal/runner/phase58.go:186-245`; `internal/runner/handoff.go`; `internal/app/consensus_request_signoffs.go:414-537`;
`internal/app/protocol.go:18-70,306-335`; `COOPERATION.md` §4.0, Phase 1, Phase 8, §9, §14, §15.
I did not run `parley protocol check`, any test, or any agent.

## Existing alternatives

Mechanisms I would otherwise build by hand, and what already ships:

- **Source hashing / live resolution** — `protocolcore.Hash` (`internal/protocolcore/core.go:189`)
  and `resolveDeck` + `protocolcore.Render` used by `protocolCheck` (`internal/app/protocol.go:307-321`).
  Inherited, with one trap: on this `protocolRole: source` deck the core render is *behind* the
  on-disk file by design (`internal/app/preflight.go:430`), so `check` reports drift. The packet
  must bind to the on-disk deck file for `source` decks, never the core render. Constraint-forced.
- **Per-track caps** — `track.PolicyFor` (`track.go:193-199`: deliberation `CapCrossReviewRounds: 3`,
  `MaxFixupCycles: 5`), `HardCrossReviewCap` bounding the consensus-BLOCK back-edge
  (`driver.go:64-73`), inclusive fix-up boundary counted from published `## Fix-up cycle N`
  (`impl.go:292-307`). Inherited; reuse, do not re-implement.
- **Disk-over-cursor precedent** — `internal/driver/cursor.go:44` ("never trusted over disk").
  Inherited for the resume rule.
- **Handoff for manual facilitation** — `runner.WriteHandoffPacket` (`handoff.go:32`) writes a
  prompt file plus instructions for `runManualSignoffAgent` (`consensus_request_signoffs.go:508`).
  This is the natural carrier for a protocol packet on the manual path. Inherited.
- **Escalation note shape** — Phase 8 LE-5 durable blocking inbox note (`COOPERATION.md:673-680`).
  Inherited.
- **Structural test technique** — the §4.0 audit's "fails when a cell has no enforcing path"
  (FINAL `:186-189`). Inherited for the "builders never read the protocol" test.
- **Frontmatter validation** — `internal/protocol/roundartifact.go` validates round files. Whether
  it tolerates an extra `protocol-context:` key is unverified, so attestation goes to the run
  manifest, not the artifact.
- **Null result:** no block classifier, omission index, or applicability manifest exists anywhere
  in `internal/` (grep `packet` hits only `HandoffPacket` and comments). Hand-built is correct.

## Proposed approach

### 1. Packet renderer: `parley protocol packet`

```
parley protocol packet --dir D --idea SLUG --phase PHASE [--agent ID] [--json|--md] [--out ABS_PATH]
parley protocol packet check --dir D            # exit 1 on any unclassified heading or manifest error
```

JSON shape (nulls explicit):

```json
{"schema":1,"sourcePath":"<abs>/parley-deck/COOPERATION.md","sourceSha256":"…","sourceRole":"source",
 "coreVersion":null,"phase":"round-01","track":"deliberation",
 "flags":{"strict_gate":false,"auto_implement":false,"protocol_change":true,"transport":"github-pr"},
 "contextMode":"packet","fallbackReason":null,
 "blocks":[{"id":"h:4.0","heading":"### 4.0 — Track selection (conditional rigor)","sha256":"…","text":"<verbatim>"}],
 "omitted":[{"id":"h:11.C","class":"transport","trigger":"transport == gitlab-mr"}],
 "unclassified":[],"generatedAt":"…","packetSha256":"…"}
```

Rules, all from the ratified FINAL: verbatim block text hashed per block; block identity is the
heading, never a line number (lines drift on every edit); `blocks ∪ omitted == every heading`,
disjoint; any parser error, unknown phase/track/flag, sha mismatch against `meta/version.json`,
or unclassified heading → `contextMode: full-fallback` with reason and the entire source as the
single block. Never committed: written under the run's manifest directory (verify with
`git check-ignore` before relying on it).

### 2. Applicability manifest, under §7

`parley-deck/meta/packet-applicability.yaml`: one entry per heading, `always | phase: [...] |
track: [...] | flag: <name> | transport: <name>`, plus a `never-cut:` list transcribing FINAL
"Never cut" as heading ids (§4.0, the current phase block, §6, §7 on protocol-change ideas, §14,
§15 full for Phases 1-3/6-7 and §15.1-15.4+15.7 for 5/8, escalation, active transport's current
phase). A missing entry is always-include and fails `packet check`. Editing the manifest is a §7
change; this idea's FINAL ratifies the first version.

### 3. Instruction layer: three paths, one renderer

- **Prompt builders** (`BuildRoundOnePrompt`, `BuildRoundPrompt`, `BuildImplementationPrompt`,
  `BuildReviewPrompt`, `BuildReviewConsensusPrompt`, `BuildFixupPrompt`, consensus/final drafters)
  take a `ProtocolContext{Mode, SourceSha256, PacketPath, OmissionIndexPath, FullPath}` and emit a
  fixed block: "Protocol context: mode=packet sha=<short>; packet at <abs>; omitted index at <abs>;
  on any doubt about an omitted rule read the full protocol at <abs>." Builders never open
  `COOPERATION.md` (FINAL `:45-46`).
- **Manual path**: `WriteHandoffPacket` gains the same fields, so a human-facilitated headless call
  carries an identical packet and attestation.
- **§9 item 1** becomes: read the packet issued for your phase when one exists, otherwise the full
  file. Default stays full-read (fail open) for interactive humans. The skill's standing line lives
  in the npm package, not this worktree; changing it needs a release, which is out of scope. Record
  that as a remaining gate, not a silent omission.
- **Attestation** is recorded per attempt in the run manifest (codex-1's telemetry field), with
  `contextMode` and `packetSha256`, so a full-fallback is visible, never inferred.

### 4. Budgets on every path: one disk-derived ledger

`budget.Ledger` (new, in `internal/driver` or `internal/budget`): `FixupCyclesPublished(ideaDir)`,
`CrossReviewRoundsOpened(ideaDir)`, `Grants(ideaDir)`, `MayOpen(kind, policy) (ok, reason)`. The
driver keeps its existing checks; the manual signoff/fix-up handoff and any command that opens a
`round-NN/` or writes a fix-up prompt call `MayOpen` first. On exhaustion: write the LE-5 blocking
inbox note, exit non-zero, never close. A grant is a file
`ideas/<slug>/budget-grant-<kind>-<N>.md` with user attribution and a single extra count; it never
resets. Resume recomputes from disk, never from the cursor. The "two consecutive confirmed
material regressions → trajectory review" rule lives in the pilot's frozen policy plus a scan of
review-consensus fields; it does not edit a §4.0 cell (audit rule, FINAL `:181-185`).

### 5. Experiment controls that touch my slice

- The 12-task pilot runs `contextMode: full` in **all three arms**; packets are generated in shadow
  and logged only. A packet in the pilot arms would confound the pilot the way the audit would
  have confounded the packet experiment (FINAL §4).
- The packet's own pre-registered A/B (6+6 paired, AB/BA, three canaries plus control, zero
  obligation misses, R≤0.50 ship, `(0.67,0.80]` returned to the user unresolved) is a separate run,
  gated on the shadow-mode omission audit showing zero never-cut violations. If time runs out it is
  a listed remaining gate in the HTML, not a claimed result.
- Freeze `sourceSha256`, manifest sha, prompt-template sha and `packetSha256` per phase before
  any measurement; a mid-run protocol edit invalidates the run.

### 6. Acceptance tests (Go, table-driven; fixtures in a disposable copy)

1. Verbatim: each block text equals the source slice; per-block sha recomputes.
2. Never-cut: for every phase×track×flag combination the never-cut ids are present.
3. Unknown phase, unknown flag, sha mismatch, malformed manifest → `full-fallback` with reason;
   body equals the full source.
4. Unclassified heading in a mutated copy → included in every packet and `packet check` exits 1.
5. Omission index complete and disjoint.
6. Structural: no `ReadFile` of the protocol path inside `internal/runner`; builders compile only
   with a `ProtocolContext`.
7. Phase 5/8 packets contain §15.1-15.4 and 15.7; Phase 1 contains §15.6.
8. Manual handoff at 5 published fix-up cycles on deliberation → refused, note written; a grant
   for cycle 6 allows exactly one more; cycle 7 refused.
9. Cursor says 2, disk says 5 → refused on resume.
10. Consensus-BLOCK back-edge on the manual round-open path respects `CapCrossReviewRounds`.
11. Secret scrub: a packet built from a copy containing a fake token string never emits it.

### 7. Work boundaries

- **claude-1:** renderer, manifest, `packet check`, builder and handoff wiring, §9 text, ledger,
  tests 1-11. I may not verdict my own claims; kimi-1 or hermes-1 reviews this slice.
- **codex-1:** manifest attestation fields, usage events, experiment harness and R recomputation
  inputs.
- **hermes-1:** packet generation failure must classify as tooling error, never as agent
  no-first-output or hang.
- **kimi-1:** adversarial manifest edits (mis-scoped `phase:` entries), independent never-cut
  audit of shadow packets, independent R recomputation.

## Concerns / open questions

- Is the run manifest directory ignored by git? If tracked, packets need another
  never-committed location.
- Does `roundartifact.go` accept unknown frontmatter keys? Until verified, no attestation in
  artifacts.
- Do the pilot launch paths use the skill's standing line at all? For headless runs the builders are
  the whole instruction layer, which makes the out-of-scope skill change irrelevant to the pilot
  but still relevant to interactive humans.
- Whether §2 roster belongs in packets: participants come from `00-prompt.md`, so I would omit §2
  with trigger "roster membership question". Open to challenge.
- Word growth 9.5k→16.6k→24.9k is the prior-rounds term, untouched by this slice (FINAL §5). The
  HTML must present R as per-call only.

## Risks

- **Silent omission** is the known unfixable class: a rule absent from packet and index gives no
  in-loop signal and Phase 2 silence reads as agreement. Mitigations are always-include on unknown,
  shadow mode first, and manifest edits under §7. Not eliminated.
- **Wrong source on a source deck**: binding to the core render would deliver an older protocol
  with a valid hash. Test 3 must include this case.
- **Cap bypass via a new path**: every prior cycle of the sibling idea found one more path that
  ignored the budget. The ledger is only as good as the list of callers; a structural test that
  enumerates round-opening functions is required.
- **Confounding the pilot** if packets are enabled by default. Default off; shadow only.
- **Self-verification**: I own these claims and issue no verdicts on them.
