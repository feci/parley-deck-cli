---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
round: 2
date: 2026-09-05
lens: phase-scoped packet and instruction-layer correctness
---

## Evidence this round (PRIMARY, executed in this worktree)

- `shasum -a 256 parley-deck/COOPERATION.md` = `e2319aff…` = `protocolSha256` in `meta/version.json`.
- `git ls-files parley-deck/runs` lists tracked files; `git check-ignore` on a run dir: not ignored.
  Phase-packet `FINAL.md:37`: "Generated on demand; never committed."
- `grep Buffers internal/runner/supervision.go internal/runner/runner.go`: no hits. `Spec.BuffersStdout`
  is at `internal/agents/discover.go:57-60` (codex cited 61-65; off by a few lines, substance correct),
  consumed only by the TUI and `run.created`.
- `grep -rn "agent.usage\|agent.acp.usage" internal`: producer `internal/runner/acp.go:389` emits
  `agent.acp.usage`; the sole `agent.usage` reader is `internal/driver/loop.go:183`.
- `grep COOPERATION.md internal/runner`: zero hits; builders never name the protocol.
- `ls /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/`: `parley-deck-skill` source repo is a sibling.
- Read: `internal/driver/cursor.go:42-65`, `driver.go:238-256`, `impl.go:290-308,510-545`,
  `internal/runner/runner.go:535-567`, `internal/app/driver_checks.go:52-97`,
  `internal/app/protocol.go:139-162,306-345`, `internal/app/preflight.go:425-431,820-889`,
  `internal/runner/handoff.go:32-53`, `internal/driver/consensus.go:96-99`, `COOPERATION.md:876`.
- Not executed: no `go test`, no build, no agent launch, no probe.

## Retractions of my round-01

1. **Wrong:** "inclusive fix-up boundary counted from published `## Fix-up cycle N` (`impl.go:292-307`)".
   Current `impl.go:297-300,536-545` takes the maximum of the run cursor's monotonic charged count and
   the `.fixup-done` markers; heading counting was rejected in the fix-up idea's round-02
   (`impl.go:529-532`). I withdraw the disk-derived `budget.Ledger` and the rule "resume recomputes from
   disk, never from the cursor" for fix-up counts: `cursor.go:46-48` makes that field deliberately
   non-rebuildable. Test 9 stays and gains its inverse: cursor 5, markers 2 → still refused.
2. **Wrong:** "the skill's standing line lives in the npm package, not this worktree; changing it needs a
   release". The skill source is a writable sibling repository. Source edits are in scope; npm
   publication and global core publication stay out.
3. **Withdrawn:** participant-authored `budget-grant-<kind>-<N>.md` files. No file a participant can
   write may raise capacity. Exhaustion writes the LE-5 note and stops; extension is a human act
   outside the loop.
4. **Answered:** the run directory is tracked, so packet bodies cannot sit beside `handoff-prompt.md`
   (`handoff.go:36-44`) without an ignore rule.
5. **Withdrawn:** test 11 "never emits the token". Redaction under the original hash is a false attestation.

## Response to codex-1: ACCEPT all corrections, no BLOCK

- **Charged budget.** Reuse `chargedFixupAttempts` and `HardCrossReviewCap`; `consensus.go:96` is the
  only enforcement comparison and it is driver-only. Codex owns budgets; I need a read-only query the
  manual handoff can call (interfaces below).
- **Unknown applicability → full fallback,** not include-one-heading. `packet check` still exits 1
  (FINAL:44-45); at render time the reason is `unclassified:<heading>`, since its dependencies are
  unresolved (FINAL:42-43). Shadow mode makes this free until the map is complete.
- **Secrets.** Renderer returns `contextMode: refused`, reason `secret-detected`, no body, no hash
  claim. This replaces test 11.
- **Source-role authority,** stated as the rule: `protocolRole: source` → the on-disk deck file is
  the resolved authority and `protocol check` reporting `hand-edited-or-stale` (`protocol.go:323-326`)
  is expected (`preflight.go:427-431`). `consumer` → the core render must equal the on-disk file,
  otherwise `full-fallback` reason `drift`. No bundled snapshot ever (FINAL:38-39).
- **Invocation IDs.** Verified `runner.go:535,545`: marker is `RunID:agent:attemptID`, no segment,
  ordinal per call. UUID plus retained ordinal.
- **Usage.** One normalized `agent.usage` per attempt; `agent.acp.usage` stays a raw stream excluded
  from sums. A test injects both names and asserts single counting.
- **Independent evidence.** Verified `driver_checks.go:78-87`: a failed evidence write only warns
  while `allPass` still returns true. Kimi's evidence path must fail closed there.
- **Three paths in source, rollout gated.** Accepted; §9 item 1 is `COOPERATION.md:876`.
- **Ownership.** I take `internal/protocolpacket`, the packet CLI file, the applicability map,
  skill-source and §9 wording. Builder and handoff call sites are codex's serialized change.

## Response to hermes-1

- **Contradiction with codex on `RealHang`: I side with codex.** `supervision.go:11-15` counts
  writers or ACP events; nothing in `runner/` consults `BuffersStdout`, so at 120 s a declared-buffered
  agent and a hung one are indistinguishable. Name it `deadline-no-output`; it opens a readiness
  gate, never an exclusion.
- **Contradiction on `SawSentinel`:** `isExactPONG` (`preflight.go:868-871`) exists to reject
  echoed prompts. A substring may record `malformed-reply`; it is never PASS and never exit 0. Keep
  the class split, drop the Live mapping and the exit-0 warning row.
- **Interface.** A packet render failure before spawn is a `tooling:packet` failed-start record,
  never `no_first_output`. Consistent with hermes's own boundary.

## Response to kimi-1

- W3 matches my slice after the retractions above. Their `parley budget status` becomes a call to
  codex's query reading the run cursor, not the idea directory.
- Dual-read of both usage names: rejected per codex (double count when both emit).
- Attestation lives in the telemetry record, not artifact frontmatter; `roundartifact.go` validates
  only the §15.6 section, so I claim nothing about extra keys.
- Freeze file, blinding, disjoint task-author/evaluator roles and non-implementer R recomputation
  are accepted as gates on the packet A/B.

## Interfaces I need from codex-1

1. Requested-record fields: `context_mode`, `source_sha256`, `packet_sha256`, `fallback_reason` (nullable).
2. `driver.ChargedBudget(runDir, ideaDir, policy)` read-only, callable from `writeSignoffHandoff`;
   exhaustion returns an error the handoff refuses on and writes the LE-5 note.
3. `.gitignore` entry `parley-deck/runs/*/agents/*/packet/` plus a `git check-ignore` test.
4. Builder and `handoff.go` call sites accept `protocolpacket.Context`.

## Bounded plan for my slice (no code yet)

1. `internal/protocolpacket`: heading parser, `Render`, `Check`, JSON with
   `contextMode ∈ {packet, full-fallback, refused}`; tests 1-10 plus revised 11, fixtures in `t.TempDir()`.
2. `internal/app/protocol_packet.go`: `parley protocol packet` and `packet check`; dispatch
   registration by codex.
3. `parley-deck/meta/packet-applicability.yaml` with the FINAL:201 never-cut list, ratified through
   this idea under §7.
4. §9 item 1 and the skill standing line as one reviewed patch; the skill edit is in a separate
   repository and is applied only when explicitly assigned.
5. Shadow packets for every pilot phase, logged only; all pilot arms stay `full`. Kimi audits
   never-cut violations independently.
6. The packet A/B (6 paired runs in each of phases 1 and 6, AB/BA, three canaries plus control,
   zero obligation misses, ship at R ≤ 0.50) runs on codex's harness against the packet commit after
   the pilot; kimi recomputes R. If unexecuted at reporting time the HTML lists an unmet gate, not a result.

A non-implementer verifies my claims per §15.1; I issue no verdict on my own slice.
