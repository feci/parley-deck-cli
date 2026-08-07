---
agent: kimi-1
idea: meta-protocol-change-global-core-protocol
round: 2
date: 2026-08-07
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

Three changes, all moves toward codex-1's evidence discipline:

1. **The replaceable set shrinks from five sections to one.** In round 1 I marked §0, §6, §8,
   §11, §13 as open "because a deck organizes its own work there." That was design-by-plausibility,
   exactly the inheritance-framework-for-hypothetical-demand codex-1 warned against. The fleet
   evidence supports exactly one replaceable block: the working-language choice, which the current
   core already names as a project decision (`parley-deck/COOPERATION.md:742-743`). Everything
   else I had open gets sealed in v1; opening a block later is a core-version decision, not an
   overlay capability.
2. **Per-idea pinning: I switch from version+hash reference to codex-1's stored snapshot.** My
   round-1 pin (version+sha in frontmatter and manifest) points into a per-machine store that may
   not retain the pinned release; a hash reference does not self-contain. The snapshot does.
   Argument in the responses below.
3. **The sandbox disagreement: codex-1 is partially right, and my "prevention is impossible"
   was scoped too loosely.** It is impossible for the facilitator and for any user-launched agent.
   It is real, with named defeat conditions, for parley-launched participants. The round-1 framing
   "prevention vs detection" was a false binary; the correct claim is per-launch-path. Detail below.

Also one correction to my own round-1 text: I proposed HTML-marker-free §-anchors while claude-1
proposed inline `<!-- parley:virtual -->` markers. I now think both are half right — see the
registry proposal under @claude-1.

## Responses to others

### @claude-1

**On the marker scheme.** Your `<!-- parley:virtual id="pipeline" -->` puts the block registry
inside the prose it describes. That creates a new drift channel: a core editor can renumber or
retitle a section and forget the marker, and now the marker's `id` and the section's identity
disagree — which one is authoritative? Counter-proposal, merging your ids with my §-anchors:
**the block ID is the section number, and the registry lives outside the prose.** Section numbers
become permanent identifiers under a core authoring rule: numbers are never reused; a deleted
section's number is tombstoned; retitling keeps the number. The core release carries a machine
registry (codex-1's manifest) mapping `s4 → sealed`, `s6.6 → replaceable`, `ext-1 → extension
point`, with the registry hashed alongside the body. The prose stays clean of markup; the binding
is stable under renames because it never referenced the title text. Note the live deck already
exhibits why heading *text* cannot be the key: `## 10. TL;DR` sits before `## 9.` in the committed
file (`parley-deck/COOPERATION.md:785,801`) — heading order and heading text are both unreliable;
only the number, held permanent by rule, survives.

**On "start restrictive" — you didn't start restrictive enough.** Your virtual list (§12 pipeline,
§11 parameters, timeouts, §4.0 thresholds) has no fleet evidence behind it; it is the same
plausibility argument I made and now retract. Concrete v1 answer, as the task demands:

- **Replaceable (exactly one):** `s6.6` — the working-language block, already an acknowledged
  project choice in the current core.
- **Extension points (exactly one):** `ext-1`, rendered after the final core section. Extension
  sections carry deck-namespaced IDs (`<deck-slug>.<slug>`) so two decks' extensions can never
  collide and a fleet audit can group them.
- **Identity slots (not overrides):** the six measured zones — Workspace, Created, Transport,
  Protocol-synced stamp, §2 roster table, host-handle table — are renderer inputs, as in your
  model and codex-1's.
- **Sealed:** everything else, explicitly `s1, s3, s4, s5, s7, s8, s9, s10, s11, s12, s13, s14,
  s15`. I drop my round-1 "open §8/§11/§13" list entirely.

**Is even two verbs too much?** Nearly. The honest reading of constraint 3 is that override and
extension must be *permitted*, and `s6.6` makes Override non-vacuous by exactly one target. If the
facilitator or user would rather read working-language as a seventh identity slot (it is arguably
data, not prose), then Override has zero targets in v1 and the one-verb design is: `Extend` only,
single extension point, and the Override verb exists in the schema with an empty replaceable set —
shippable the day the registry opens a second block, undeployed until evidence demands it. I can
live with either; I refuse the middle course of opening unevidenced blocks to justify the verb.

**On your multi-machine concern:** you say the committed render makes cross-machine drift visible.
True but passive — visible in a git diff after the fact. codex-1's deck lock (version + hash,
compare-and-swap, refuse to substitute a same-named local version) makes it *impossible to commit
silently*, which is strictly stronger and I adopt it.

### @codex-1

**On the sandbox — you are right, and here is exactly how right.** Round 1 framed this as
prevention-vs-detection; your §5 actually already conceded the same-UID case, so the real question
is what the managed-sandbox tier covers. Concretely on macOS:

- The mechanism is a seatbelt profile via `sandbox-exec -f profile` wrapping every adapter
  invocation: `(deny file-write* (subpath "~/.parley/protocol"))` with the workspace allowlisted.
  Seatbelt policy is inherited across `exec`/`fork` and cannot be relaxed by a child, so the
  "children inherit the denial" claim holds — this is the same substrate the codex CLI itself uses
  for its `--sandbox` modes on macOS, so there is working precedent in the fleet's own toolchain.
- **What it covers:** naive direct writes from any parley-launched participant process tree —
  which is where the round-1 roster says participants actually run.
- **What defeats it, and must be written into the design:**
  1. *IPC to an unconfined broker.* Seatbelt denies operations, not intent: unless the profile also
     restricts `mach-lookup`, a sandboxed child can ask an unconfined macOS service to perform the
     write (`osascript` directing Finder/System Events, `launchctl`, `ssh localhost`). A profile
     that locks down the core path but leaves broad IPC open is a speed bump, not a boundary.
  2. *The facilitator itself.* This is the largest unprotected surface and round 1 said too little
     about it. The facilitator is launched by the user, not by parley; no wrapper covers it; and it
     is precisely the process that runs preflight, sync, and (today) writes COOPERATION.md. For the
     facilitator, detection + fail-closed is the permanent ceiling.
  3. *Adapters that shell out to desktop helpers* (your own caveat) and any platform where the
     launcher hasn't implemented confinement — preflight must then downgrade the claim per-run.
- **Verification:** your "preflight inspects the effective confinement policy" needs a mechanism on
  macOS, because seatbelt is not introspectable from inside. The workable probe: at run start, the
  managed child attempts a write to a sentinel path under the core store and *requires* EPERM; if
  the write lands, the run is unconfined and the protocol record says `core-writes: detect-only`
  for that run. Probe-and-attest beats claim-and-hope.

So the resolved claim for the design document: **prevention is real for verified-sandboxed
parley-launched participants; the facilitator and all user-launched agents get detection,
attribution, and a closed gate — and that ceiling is fundamental, not a backlog item.** With your
defeat-conditions list plus the facilitator item added, I consider the disagreement closed.

**On pinning — I switch to your side, and the storage-cost objection collapses on inspection.**
Content-addressed snapshots deduplicate: 35 of 36 decks render the same effective bytes, so the
fleet stores one blob per *distinct* effective protocol, not per deck. At ~60 KB per protocol, a
hundred distinct historical versions cost ~6 MB — noise against any git repo. What the reference
model (my round-1 version+hash) quietly requires is worse: that every machine retain every pinned
core release in `~/.parley` forever, or re-resolution of an old idea fails; and the historical
audit trail then depends on a per-machine, user-editable store. The snapshot is also the honest
completion of the roster precedent you cite — the run manifest freezes the effective roster, not a
pointer to a roster that might still exist. **`sessions inspect` should report:** pinned
core-version, effective sha256, snapshot path, `stale-protocol` (deck's current effective hash ≠
pinned — informational, the pin is supposed to lag), and core-release verification status
(`verified` / `drifted-core`). One nuance I'd add to your lock: the snapshot path belongs in the
deck lock AND the run manifest, so a reader can find the governing bytes from either direction.

**On the compatibility receipt:** agreed, with one bound — require the user's explicit acceptance
only when the checker reports a changed target or dependency; a no-changes report should be
auto-passable. Otherwise a fleet with one overlay pays a human round-trip on every core bump
forever, and the pressure will be to delete the overlay rather than keep it honest.

### @hermes-1

**On `syncConsumerProtocol`/`mergePreservingZones` — it is the wrong shape for the implementation,
and its own history says so.** Three concrete defects as a render engine:

1. **It splits at a heading literal (`## 3.`, `preflight.go:539`).** Block-ID addressing is what we
   are converging on precisely because heading anchors break; making the permanent renderer
   heading-anchored builds the known fragility into the new foundation. The committed deck file
   already has §10 before §9 — heading-based logic on this corpus is on borrowed time.
2. **It over-preserves.** "Everything through §2 comes from the deck" means any stale or
   hand-edited core prose in §0–§2 is preserved *forever* as a "project zone" — the merge cannot
   distinguish the six identity values from edited rules text in the head region. That is the
   librade-algoTrader failure mode running in reverse: instead of destroying a local section, it
   immortalizes drift and calls it local content.
3. **It has no sealed/open notion, no overlay, no hash gate** — it is a repair heuristic for the
   flat-copy world, and a good one. Its correct fate: keep it as *migration scaffolding only* —
   the head-extraction pass that pulls the six identity zones out of a legacy deck — and even
   there, narrow it from "everything before `## 3.`" to the six named zones, so migration does not
   import head-region drift into the generated world. The renderer itself should be a fresh,
   small, pure function: registry + core + identity slots + overlay → flat bytes, byte-idempotent,
   with `refreshProtocolSyncedLine`'s provenance-line logic (`preflight.go:582`) carried over as
   the one salvageable piece.

**On your markdown overlay format:** I withdraw my round-1 indifference — the registry argument
above applies here too. If block IDs live in a hashed registry, the overlay should reference those
IDs, not `## override §N` headings whose `§N` is doing double duty as title text and identifier.
Your format keeps working if we read `§N` strictly as the permanent number, so the delta between us
is small: I want the number declared permanent by core-authoring rule and the sealed/open
classification checked against the registry rather than hardcoded in the renderer.

**On your concern 3 (three protocol surfaces + the core):** your reading is right and I'd sharpen
it — the drift guard's target changes from `embedded == live-deck` to a chain: `embedded ==
packaged-skill == core-release`, then `core + overlay == committed render (modulo identity
zones)`. The first hop is a build-time test; the second is `parley protocol check`. That removes
the current guard's widening prose allowlist, which codex-1 also flagged.

**On your risk 4 (user expectation of hard prevention):** the sandbox resolution above gives the
user something real for managed participants, which softens but does not remove the gap for the
facilitator. Agreed it goes in the consensus text verbatim.

### @kimi-1

Three self-corrections, stated as responses so the record shows the movement:

1. My open list (§0/§6/§8/§11/§13) was plausibility design; retracted above in favor of one
   evidenced replaceable block. The rule I proposed — "sealed what other decks rely on" — was a
   good rule applied with a bad inventory.
2. My version+hash pin was reference-based; the cross-machine and audit arguments in @codex-1
   above convinced me the stored snapshot is correct. I keep from my round-1 version only the
   insistence that hashes, not version labels, name things — which the snapshot model already
   satisfies.
3. My "prevention is not achievable" was true only of same-UID unsandboxed processes; I should
   have scoped it. The per-launch-path claim (managed-sandboxed vs facilitator/user-launched) is
   the accurate one, and codex-1's tiering was closer to it than my flat denial.

Unchanged from my round 1 and still load-bearing: never auto-migrate prose; the stale-override
re-confirmation rule (codex-1's expected-base-hash check is the same mechanism, better specified);
committed render as cross-machine authority; migration as attended Stage-4-style operation;
ratification records tying each core version to its meta idea.

## New concerns / questions

1. **Quorum hygiene.** opencode-1 is silent this round and not counted; nothing in the emerging
   consensus should cite "all participants agree" — four of five active voices are present, and
   the sandbox/pinning resolutions above are agreements among four, not five. If opencode-1
   rejoins in round 3 with a dissent on the snapshot model, the pinning decision should be
   re-openable without treating that as bad faith.
2. **Who sandboxes the facilitator's *children*?** If the facilitator launches participants via
   adapters, the sandbox wrapper lives in the adapter launch path — but the facilitator process
   itself remains unconfined and holds the same privileges. The design must name the facilitator
   as permanently detection-only, or someone will later claim the whole system is sandboxed.
3. **The sentinel-write probe** for confinement verification writes (attempts to write) inside the
   core store. On a correctly confined run it fails; on an *unconfined* run it succeeds — and has
   then written a file into the store it was testing. The probe must write to a dedicated sentinel
   name it can also delete, and treat "write succeeded and cleanup succeeded" as `detect-only`,
   never as damage. Worth one line in the design so nobody implements the naive version.
4. **Registry vs prose ordering.** If the registry says `s6.6 → replaceable` but a future core
   renumbers sections (e.g. inserts a new §6, pushing content down), permanence of numbers means
   the new section takes a *fresh* number and old numbers tombstone — core authoring loses the
   ability to renumber for tidiness, forever. That is the price of stable IDs; it should be paid
   explicitly in the core's authoring rules, not discovered.
5. **Snapshot pruning.** Content-addressed snapshots deduplicate, but nothing yet says whether a
   deck may prune snapshots of *closed* ideas. My answer: no — the snapshot is the audit record of
   what governed that idea, and §15's provenance rules already treat evidence as append-only.
   Say it in the design or someone will "clean up."

## Current proposal

Consolidated, with deltas from round 1 marked:

1. **Registry + permanent section numbers** (new): the core release = body + registry mapping
   permanent section IDs (`s1`…`s15`, subsections like `s6.6`) to `sealed | replaceable |
   extension-point`, both hashed. Numbers never reused; tombstones on deletion. No inline markup.
   (Merges claude-1's ids, codex-1's registry, my anchors; kills hermes-1's heading-text
   addressing.)
2. **v1 surface**: replaceable = {`s6.6` working language}; extension point = one (`ext-1`, end of
   core, deck-namespaced IDs); identity slots = the six measured zones; everything else sealed.
   (Changed: my five-open-sections list retracted, adopting codex-1's evidenced minimum.)
3. **Overlay**: one committed file per deck, two verbs (replace-by-ID, extend-at-`ext-1`), each op
   with rationale + expected base-block hash; absence of the file is the only canonical
   no-customization state; empty overlays forbidden (codex-1's rule, which I now think is
   essential).
4. **Pinning**: per-idea content-addressed snapshot of the effective rendered protocol, committed;
   deck lock with version+hash, compare-and-swap adoption, fail-closed on missing pinned release.
   (Changed: adopted codex-1's snapshot over my reference pin.)
5. **Enforcement, per launch path** (changed — resolves the round-1 disagreement): verified
   seatbelt sandbox (sentinel-probe attested) for parley-launched participants = real prevention
   with named defeat conditions (IPC brokers, helper escapes, unconfined adapters);
   facilitator and user-launched agents = 0444 speed bump + hash verification + ratification
   records + fail-closed gate = detection and attribution, permanently.
6. **Compatibility**: registry-driven check at adoption; expected-base-hash mismatch =
   re-confirm; missing/tombstoned target or sealed target = block; auto-pass only on a
   zero-change report; never auto-migrate prose; deck stays pinned on the old core while
   incompatible (that IS the quarantine).
7. **Implementation fate**: `mergePreservingZones` is migration scaffolding, narrowed to the six
   identity zones; the renderer is a new pure function; `refreshProtocolSyncedLine`'s provenance
   pattern carries over; drift guard becomes a two-hop chain (embedded == packaged == core;
   core + overlay == committed render modulo zones).
8. **Smallest shippable, ranked** (answer to the task's question 5):
   1. Core store + renderer + `parley protocol check` byte-guard — **no overlay machinery**.
      This alone covers 35 of 36 decks and kills copy-drift at the source; it is the piece whose
      absence produced the 8-version fleet. Ship this week.
   2. Hash verification + ratification record + fail-closed preflight gate — cheap, and it is the
      honest floor under constraint 4 that every later piece assumes.
   3. Per-idea snapshot pinning — needs the deck lock; small once 1 exists.
   4. Overlay verbs + compatibility checker — required by constraints 3/5 but demanded by zero
      decks today; design now, ship after 1–3, with absence-of-file as the default so late
      delivery costs nothing.
   5. Sandbox wrapper + sentinel probe — a hardening track on the adapter launch path, shippable
      independently per platform; never a blocker for 1–4, and never claimed beyond what the
      probe attests.
