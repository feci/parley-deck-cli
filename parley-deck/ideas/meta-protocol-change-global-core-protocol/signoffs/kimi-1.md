---
idea: meta-protocol-change-global-core-protocol
agent: kimi-1
date: 2026-08-07
verdict: BLOCK
---

# Signoff — kimi-1

## Verdict

BLOCK — two actionable items, both cheap to fix:

1. **D7's conditional snapshot materialization re-opens the hole three of four participants
   converged to close, and can violate user constraint 1.** D7 writes the rendered body only
   "when the render is not reproducible from those inputs." In the common case (all inputs
   present at pin time) the body is never written. Later loss of the pinned release — pruning,
   machine migration, fresh clone — then makes the open idea's governing bytes unrecoverable:
   the body was never stored, the release is gone, D8 blocks substitution, and claude-1's
   "materialize it then" is impossible because the body can no longer be rendered. Constraint 1
   ("an idea that is already open runs to completion under the protocol version it started
   with", `00-prompt.md:38-39`) fails in a mundane new-laptop scenario. The rationale given
   (claude-1 round-2: "copying ~80 KB into every idea directory buys nothing") argues against
   codex-1's *round-1* per-idea copy — which codex-1 itself replaced in round-2 with
   dedup-by-effective-hash ("I no longer propose one physical Markdown copy per idea",
   `round-02/codex-1.md:53-55`). The consensus adopted a counter to a superseded proposal over
   a stated three-way convergence (codex-1, hermes-1, kimi-1 round-2). Fix, either:
   (a) always write the effective body, deduped by effective hash at a deck-level
   content-addressed path (~60 KB per *distinct* effective protocol — the storage objection
   does not survive dedup), and state that snapshots are append-only (my round-2 concern 5,
   currently unanswered); or (b) keep conditionality only with a specified release-recovery
   channel and a never-prune guarantee for releases. (a) is strictly simpler.
2. **`consensus.md` lacks the mandatory `## Drafter position changes` section.** §15.5
   (`parley-deck/COOPERATION.md:1303-1309`): when the facilitator is also a participant and
   drafts `consensus.md`, the artifact MUST contain that section, or `None` — and "existing
   signoffs ratify its accuracy and completeness." The section is absent (PRIMARY: full read of
   `consensus.md`, 221 lines, sections 0–5 + Signoffs only), so there is nothing to ratify,
   and it is not obviously `None`: claude-1's round-2 requirement that "the overlay declares
   the core version range it targets" (`round-02/claude-1.md:139`) has no disposition in
   D4/D10. Fix: append the section; one line either adopting the version-range declaration or
   recording its drop with reason.

## Answers to 1-6

**1. My positions.** Adopted, accurately credited: write-once versioned store (D1);
registry + permanent IDs + tombstones + no inline markup (D2, my merged proposal);
near-empty surface with my five-section retraction recorded correctly (D3); two-verb single
overlay with expected-base-hash (D4, my round-1 stale-override rule); core-owned resolution
(D5); committed generated view reporting every replaced/removed block (D6); pinning +
`stale-protocol` (D7/D8); per-launch-path enforcement with diagnosis never blocked (D9 — my
rule, claude-1's narrowing accepted: blocking `protocol check` on a tampered core bricks
diagnosis exactly when needed); never auto-migrate, quarantine = stay pinned (D10); §7
blast-radius (D11); attended Stage-4 migration + librade restore (DF-2/DF-4); quorum hygiene
re: opencode-1 (§0). Rejected with reason I accept: my ratification-record attribution layer —
codex-1's round-2 critique (an unsandboxed writer rewrites hash+receipt together; no local
trust anchor, SECONDARY: `round-02/codex-1.md:136-141`) mirrors my own round-1 hash
self-reference concern, and D8's committed lock supplies the independent anchor. Downgraded
without adequate reason: the always-materialized snapshot (block item 1). Dropped, minor:
`parley protocol audit` fleet surface (my round-1 mitigation for N-version spread — suggest
DF-5); sentinel-probe cleanup rule (round-2 concern 3 — half-adopted via the word "sentinel").
Misrepresentation, minor: §3 claims "hermes-1's ranking, which all four accepted" — hermes-1's
round-2 order is overlay 2nd / pinning 3rd, mine was hash-floor 2nd / overlay 4th, codex-1's
was pin 2nd / migration 3rd / overlay 4th (PRIMARY: `round-02/hermes-1.md:85-99`,
`round-02/kimi-1.md:251-263`, `round-02/codex-1.md:271-284`). The consensus order is closest
to codex-1's; correct the attribution.

**2. VC-1.** Ownership note (§15.1): the resolved claim is codex-1's; my round-1 claim was the
rejected one, so I may verdict. I re-ran the probe on this machine (macOS), same profile shape,
`/tmp/fakeparley` target: resolved path `/private/tmp/fakeparley` — direct write EPERM, child
write EPERM (inherited), `rm -f` EPERM with file surviving; additionally rename and `chmod`
also EPERM (two of codex-1's conformance items). Unresolved path `/tmp/fakeparley` (`/tmp` →
`/private/tmp` symlink): write succeeded, exit 0, **no error** — the silent-disable trap
reproduces exactly. (PRIMARY: commands and outputs executed by me this session.) Verdict on
the resolution: CONFIRMED — an executed test over three argued positions is §15.3 working as
specified. Both scope limits are correctly stated; on this machine `~/.parley` resolves to
`/Users/tomasfecko/.parley` with no symlink traversal (PRIMARY: `os.path.realpath`), but the
trap class is proven. One limit is *missing*, though raised in round 2 by codex-1 and me:
delegation — a sandboxed child can ask an unconfined broker (`osascript`/Finder, `ssh
localhost`, a pre-existing MCP server) to write; the sentinel probe proves direct-write denial
only, so a run can report confinement-proven while a delegation path exists. D9's "proves it"
must be scoped to direct writes, or DF-1's profile/conformance suite must cover IPC and
inherited writable descriptors (condition below).

**3. D3 vs constraint 3.** Ownership note: I co-own D3's shape, so I ground this in the
constraint text, not preference. Constraint 3 (`00-prompt.md:41-43`, PRIMARY) requires that
override AND extension be *allowed*, OOP-style — i.e. declared-overridable parts, not a minimum
target count. A base class with one `virtual` method and one documented extension hook honors
that. The mechanism is real, not decorative: both verbs exist with a live, evidenced target
each (s6.6's override allowance is in the current core at `parley-deck/COOPERATION.md:743`,
PRIMARY; ext-1 takes the one payload the fleet ever produced — DF-4 plans exactly that
restore, `00-prompt.md:20-26`); validation fails closed (G4); opening further blocks is a
core-version decision, so the expansion path is structural. The honest limits, which the
protocol text should state rather than gloss: override has exactly one target in v1, and
extension appends at one point only, so a local rule elaborating §15 lives at the end with a
cross-reference. That is a usability compromise chosen on evidence, not a token gesture — on
condition that rank 3 actually ships this cycle (it is in §3 scope; if it slipped to deferred,
constraint 3 would be unmet).

**4. Ranked scope.** The order is right. Rank 1 before overlay is correct (35/36 decks need no
overlay); pinning at 2 is what makes adoption safe under constraint 1; detection-layer at 4 is
acceptable *only because* D8's lock implies byte-hash verification ships inside ranks 1–2 —
which is why G8 (below) must make that executable rather than implied; my round-2 rank-2
"honest floor" concern is covered iff G8 exists. Rank-1 wording risk: "`mergePreservingZones`
… becomes the empty-overlay special case" — three reviewers argued it is migration scaffolding,
not the renderer (PRIMARY defects: heading-literal anchor `preflight.go:539`; over-preservation
of everything before `## 3.` `preflight.go:522-548`; `time.Now()` stamp `preflight.go:585`).
G1 catches the nondeterminism, but the consensus should say the renderer is a new pure function
with the synced-stamp as a lock input. Deferred items can all wait — DF-2 *must* wait for
tooling (the 2026-08-06 data loss argues it), DF-1 ships honestly as `confinement-unproven`
meanwhile — except the D7 fix, which is in-cycle rank 2 (block item 1).

**5. Gates.** G1–G6 are good but none targets the failure class this project has shipped —
"documented as landed, wrong at the call site." Add:

- **G7 (call-site truth):** every guarantee named in protocol text or CLI output
  (`stale-protocol`, `confinement-unproven`, `drifted-core`, "blocks rather than substituting",
  "reports every block it replaces or removes") MUST be asserted by an end-to-end test driving
  the real command entry point (`parley preflight` / `protocol render` / `protocol check` /
  `sessions inspect`) against a fixture deck, asserting exit code, reported status, and
  resulting file state — not only unit tests of internals. A guarantee without such a test
  MUST NOT be documented as landed.
- **G8 (lock byte-verification):** adoption and continuation MUST byte-compare the on-disk
  release against the deck lock's core hash and refuse on mismatch; a test MUST install a
  same-version-label / different-bytes release and prove refusal. Without this, D8's lock
  records a hash nothing checks — the exact documented-not-wired failure.

**6. Collectively missed.** (a) Release distribution/recovery across machines — part of block
item 1. (b) Snapshot/release retention rule — part of block item 1. (c) `protocolRole: source`
is undispatched: today it means "advisory only, never writes COOPERATION.md"
(`internal/app/preflight.go:389-393`, PRIMARY); under a global core the per-deck upstream
concept dissolves, and rank-1 implementation will trip on that switch (hermes-1 round-2
concern 3, never answered). (d) The synced-stamp must become a lock input, not `time.Now()`
(`preflight.go:585`, PRIMARY). (e) The delegation/IPC sandbox limit (see Q2). (f) Sentinel
cleanup: on probe success the sentinel must be deletable and the result recorded as
`confinement-unproven`, never as damage.

## Conditions (if any)

Both block items above, plus, before I flip to ACCEPT:

1. G7 and G8 added to §5 as written in answer 5.
2. D9/DF-1 scoped: the preflight probe proves *direct-write* denial; either say so, or extend
   the profile/conformance suite to delegation paths and inherited writable descriptors.
3. §3 rank 1 states the renderer is a new pure function (synced-stamp from the lock);
   `mergePreservingZones` is zone-extraction scaffolding only.
4. §3 attribution corrected (the ranking is the synthesis closest to codex-1's, not
   "hermes-1's, accepted by all four").
5. `## Drafter position changes` includes a disposition for the overlay version-range
   declaration (part of block item 2).
6. `protocolRole: source` disposition recorded (retire, map, or keep — one line; hermes-1's
   unanswered round-2 concern 3).
