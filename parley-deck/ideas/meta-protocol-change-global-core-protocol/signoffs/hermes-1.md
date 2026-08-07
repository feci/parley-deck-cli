---
idea: meta-protocol-change-global-core-protocol
agent: hermes-1
date: 2026-08-07
verdict: ACCEPT
---
# Signoff — hermes-1

## Verdict

ACCEPT, with three conditions. The design is sound, the sandbox resolution is
verified by my own re-run, and the ranked scope is correct. The conditions are
procedural and gate-level gaps, not design defects — all fixable before or at
FINAL.md without reopening the deliberation.

## Answers to 1-6

### 1. Are my round 1-2 positions adopted, or rejected with a reason I accept?

Mostly adopted. Two items to name:

- ADOPTED: My round-2 ranking shape (core+render first, enforcement last) is
  reflected in §3. My round-2 self-correction on `mergePreservingZones` (it is
  migration scaffolding, not the renderer; the renderer is a fresh pure
  function) is reflected in §3 rank 1, which narrows it to six identity zones
  and replaces the `## 3.` anchor with registry IDs.

- MISATTRIBUTED (minor): §3 says "Adopting hermes-1's ranking, which all four
  accepted." The actual order — pinning (rank 2) before overlay (rank 3) — is
  codex-1's ordering (`round-02/codex-1.md` §262-284), not mine. My round-2
  ranking put overlay second and pinning third (`round-02/hermes-1.md` §90-94).
  The consensus order is better (pinning is safety-critical for open ideas;
  overlay serves 1-in-36), but the attribution is wrong. Not blocking — the
  order is right, the label is not.

- DROPPED: `protocolRole: source` inversion. I raised this in round 1
  (`round-01/hermes-1.md:288-293`) and round 2 (`round-02/hermes-1.md:77`). No
  participant addressed it; the consensus does not mention it. Today the source
  repo's deck is the protocol upstream and preflight never auto-writes it
  (`internal/app/preflight.go:389-395`). With a global core, this relationship
  inverts: the source repo's `~/.parley` core IS the authority, and the deck
  copy is generated from it. This is a real implementation gap, not a design
  defect — but it should be recorded as an implementation note or deferred
  item so it is not discovered during Phase 5. See Condition 3.

### 2. Is the VC-1 sandbox resolution sound? Re-run result.

I re-ran all four tests. Results match `round-02/claude-1.md` verbatim [PRIMARY]:

```
Test 1 (resolved path /private/tmp/fakeparley):
  /bin/sh: COOPERATION.md: Operation not permitted    # denied, exit=1

Test 2 (child process):
  /bin/sh: COOPERATION.md: Operation not permitted    # denied, inherited, exit=1

Test 3 (rm -f):
  rm: COOPERATION.md: Operation not permitted         # denied, exit=1
  file content after rm attempt: original-content     # file survives

Test 4 (unresolved path /tmp/fakeparley — symlink trap):
  exit=0                                              # NO error
  file content: tampered                              # write SUCCEEDED
```

The resolution is sound. codex-1's mechanism claim (OS sandbox denies writes,
children inherit, file survives deletion) is CONFIRMED by my own execution.
The 3-vs-1 majority — which included my own round-1 position — was wrong, and
the consensus correctly resolved it on the strength of the test rather than by
counting. This is exactly what §15.3 requires.

The two scope limits are correctly stated and I verified both:

1. Facilitator excluded: correct. The sandbox covers processes parley launches.
   The facilitator (a Claude Code process parley did not launch) is outside it.
   This is not a caveat — it is the largest unprotected surface, as I argued in
   `round-02/hermes-1.md:47`.

2. Unresolved path silently disables: CONFIRMED by my Test 4. A profile built
   from `/tmp` (symlink to `/private/tmp`) permitted the write with no error.
   The consensus correctly states this means "a configured sandbox is not
   evidence that anything is confined" and requires a preflight write-probe
   against the resolved path (D9, G3).

One refinement: the consensus D9 says the sentinel probe writes "to a sentinel
inside the protected subtree." kimi-1's round-2 concern #3
(`round-02/kimi-1.md:204-208`) noted that on an unconfined run the probe
succeeds and has then written a file into the store it was testing — the probe
must write to a dedicated sentinel name it can also delete. The consensus does
not state this cleanup requirement. It is an implementation detail, but given
this project's history of "documented as landed and wrong at the call site,"
stating it explicitly would prevent the naive implementation. Not blocking —
G3's "probe MUST use the resolved path" can be read to include this — but I
flag it for the implementer.

### 3. Does D3's near-empty open surface satisfy constraint 3?

Yes. Constraint 3 requires "local override AND extension are allowed, in the
sense of object-oriented inheritance." D3 provides:

- Override: `s6.6` (working language), which the current core already explicitly
  permits as a project choice (`COOPERATION.md:743`: "Projects that explicitly
  need a different working language may override this rule"). This is a real
  override of a defined protocol part, not a data slot.
- Extension: `ext-1`, one append point with deck-namespaced IDs.

Both verbs exist and are non-vacuous. An OOP class with one overridable method
and one extension hook is still "override and extension." The 1-in-36 fleet
evidence justifies the minimal surface — I argued this myself in round 2
(`round-02/hermes-1.md:19`): "The two-verb mechanism should exist; the initial
open set should be near-empty." Opening more blocks on speculation would be the
over-engineering CLAUDE.md §2 warns against. If a second genuine override case
appears, a meta idea opens it — the §7 path the consensus itself describes.

This is not a token gesture. It is the smallest surface that honestly satisfies
the constraint, backed by evidence.

### 4. Is the ranked implementation scope (§3) right?

Yes. The order is correct: core store + render + check (rank 1) delivers the
bulk of drift reduction; pinning (rank 2) protects open ideas before any core
change can reach them; overlay (rank 3) is deliberately third because shipping
override machinery before the generator builds for a case that barely exists;
detection-layer enforcement (rank 4) is the honest floor that every later piece
assumes.

Nothing in ranks 1-4 must move. Nothing deferred (§4) cannot safely wait:
- DF-1 (sandbox): the detection layer (rank 4) ships without it; every surface
  reports `confinement-unproven` until it lands. Safe.
- DF-2 (fleet migration): tooling from ranks 1-4 must exist first — correct
  ordering, given the 2026-08-06 data loss.
- DF-3 (opencode-1 fitness): out of scope for this idea.
- DF-4 (librade-algoTrader restoration): depends on the overlay (rank 3).
  Correctly deferred.

One note: the consensus attributes the ranking to me ("Adopting hermes-1's
ranking"), but the pinning-before-overlay order is codex-1's. See Q1 above.

### 5. Are the gates (§5) sufficient to catch a half-done implementation?

Almost. G1-G6 cover the mechanism's internal correctness: idempotent render
(G1), no write path to core (G2), confinement probe (G3), overlay fail-closed
(G4), pinning stability (G5), changelog (G6).

The class "documented as landed and wrong at the call site" is not fully
covered. G5 proves that changing the core mid-idea does not change what the
open idea resolves — but it tests the pinning mechanism, not whether the
protocol bytes that an agent actually reads at session start are the pinned
bytes. A render command that works, a pin that holds, and a preflight path
that still reads the old flat copy would pass G1 and G5 while the live read
path is wrong. That is exactly the failure shape this project has shipped
before.

Missing gate — proposed G7: an end-to-end acceptance test that verifies the
protocol bytes an agent reads at session start (the actual §9 read path, not
the `protocol render` command) are byte-identical to the pinned effective hash
recorded in the run manifest. This catches the "mechanism works in isolation,
wiring to the real call site is missing" class. See Condition 2.

### 6. Anything the four of us collectively missed.

1. `protocolRole: source` inversion — see Q1. Raised twice by me, unaddressed.
   The source repo's deck is currently the protocol upstream; with a global
   core, this inverts and needs explicit handling. Not in the consensus, not
   in deferred items.

2. §15.5 `## Drafter position changes` section is missing from `consensus.md`.
   §15.5 requires this when the facilitator is also a participant and drafts
   `consensus.md` — which is the case here (claude-1 is facilitator, drafter,
   and participant, as §0 records). The section must record every material
   change in the drafter's position since its most recent round file. claude-1's
   position on prevention changed materially from round 1 ("not achievable") to
   round 2 ("prevention for parley-launched agents") — this is quoted in VC-1
   but not in a formal `## Drafter position changes` section. This is a binding
   procedural requirement on every track (`§15.7`). See Condition 1.

3. §15.6(b) record: on `deliberation`, consensus MUST record that unanimity
   among related models is a shared prior, not independent evidence, and state
   what would have to be true for the agreed position to be wrong. The
   participants are from four different model families (Anthropic, OpenAI,
   Zhipu, Moonshot), so the "related models" concern is weaker — but the
   "state what would have to be true for the agreed position to be wrong"
   requirement still binds. The consensus does not contain this. I note it
   but do not make it a condition: the sandbox resolution already provides a
   concrete falsification condition (the unresolved-path trap), and the
   deferred-items list implicitly states what could go wrong. A one-line
   addition would satisfy it cleanly.

4. Sentinel-probe cleanup (kimi-1's round-2 concern #3): the probe writes to
   the core store on an unconfined run and must clean up after itself. D9/G3
   do not state this. Minor, but this is exactly the kind of detail that
   becomes "documented as landed, wrong at the call site" if unstated.

## Conditions

1. `consensus.md` must add a `## Drafter position changes` section per §15.5,
   recording claude-1's position change on prevention (round 1: "not
   achievable" → round 2: "prevention for parley-launched agents, detection
   for the rest"). This is a binding procedural requirement; signoffs ratify
   its accuracy, so it must exist before FINAL.md.

2. Add G7 to §5: an end-to-end acceptance test verifying that the protocol
   bytes an agent reads at session start (the actual §9 read path) are
   byte-identical to the pinned effective hash in the run manifest. This
   catches the "mechanism works in isolation, wiring to the live call site is
   missing" failure class this project has shipped before.

3. Record the `protocolRole: source` inversion as an implementation note or
   deferred item (DF-5). The source repo's deck is currently the protocol
   upstream (`preflight.go:389-395`); with a global core this inverts and
   needs explicit handling in Phase 5. It should not be discovered during
   implementation.
