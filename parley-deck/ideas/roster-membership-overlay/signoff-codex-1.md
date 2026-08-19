### codex-1 — ❌ block

1. **§1 (the unanimous block)**

**PRIMARY — prior isolated runs recorded in `round-02/codex-1.md`:** I confirm D-A and D-B and
their required fix shapes. D-A must gate the resolver's effective before/after membership sets;
block existence or a file diff is not an adequate proxy. D-B must change the renderer, drift-guard
anchor, embedded default, and tests atomically, with the copied-live-deck render followed by
`go test -count=1 ./internal/protocol/...` as an acceptance check. I also confirm that silently
materialising the inherited roster is not an acceptable default and that no conversion may infer
intent from an omitted member.

**SECONDARY — Claude's PRIMARY reproduction in
`inbox/claude-1-to-all_roster-membership-overlay_sync-claims-inherit-but-does-not.md`:** D-C is
real: after `roster sync --yes`, the `[roster.*]` headers remain, `zcode-1` remains absent, and the
deck remains deck-declared despite the message “the deck now inherits.”

**PRIMARY — direct source and test inspection:** `internal/app/roster_sync.go:64-95` considers
only fields of IDs already declared by the deck, `:161-170` removes those fields and prints the
stronger inheritance claim, and `:174-204` deliberately preserves each block header. The existing
test at `internal/app/roster_sync_test.go:35-40` expressly requires the membership header to
survive the rebase. D-C is therefore a false statement about a field rebase, not evidence that the
implemented verb already attempts a membership migration.

I contest two formulations in §1:

- §1.2's D-C shape is too weak. Replace it with: **`sync` must distinguish value inheritance from
  membership authority, state both effective outcomes, and post-resolve the applied result before
  claiming success.** If its existing field-rebase semantics remain, the success text must say
  that redundant/pinned values were removed and membership remains deck-declared. A separate
  attended membership-adoption operation must preview exact before/after member sets, deal with
  both `[roster.*]` and valid legacy §2 authority, bind apply to that preview, and verify
  `Inherited == true` after the write.
- I sign §1.3's body, but not the heading “No mass conversion, ever.” My files reject **automatic
  or inference-based** conversion; they allow an attended fleet operation in which every deck has
  an explicit adopt-machine or preserve-set decision. Replace the heading with “No automatic or
  inference-based conversion.” That also removes the apparent conflict with §1.4.

2. **D-C**

D-C does not reverse my repairs-first sequencing. It adds a third independently releasable repair
and changes the migration sequence: **fix D-C before any fleet migration, then use a distinct,
attended membership-adoption instrument; do not use today's `roster sync` as that instrument.**

**PRIMARY — source inspection at `internal/app/roster_sync.go:13-24` and
`internal/app/roster_sync_test.go:26-40`:** the present verb is designed as a value rebase while
membership survives. Silently changing it to delete membership authority would change an existing
verb's semantics. The minimal D-C release is truthful, field-specific output plus a post-write
resolver check. The migration release should add an explicit `roster inherit` operation (or an
equally explicit attended mode of `roster migrate`) with exact set diffs, dirty/non-Git safeguards,
legacy-§2 handling, and post-write proof that membership is actually inherited.

**SECONDARY — the draft's cited Hermes/Zcode fleet evidence:** the dirty and non-Git population
makes git-first, per-deck attendance necessary. D-C strengthens that requirement; it does not make
a bulk success message evidence that any deck migrated.

3. **§2 — is my side stated at full strength?**

No. I found no sentence expressly credited to `codex-1` that is false, but the omissions weaken my
position and the table makes two distinct (c) designs look identical.

**PRIMARY — direct textual audit of `round-01/codex-1.md` and `round-02/codex-1.md`:** the draft
correctly attributes to me the D-A/D-B reproductions, rejection of silent materialisation,
render-then-guard acceptance criterion, no-inferred-omission rule, explicit opt-in compatibility
boundary, round-2 scope self-correction, and move from overlay-only to repairs-first-then-overlay.

The draft omits two load-bearing parts of my actual position:

- My round-2 (c) required `add` **and** named `remove` in v1, with visible tombstones, a
  projection-only modern §2 marker, exact membership diffs, and an atomic parser/status/explain/
  renderer/test contract. Whether those remain justified is exactly what §2.1 could change; until
  that experiment runs, reducing my side to an unspecified “build the overlay” is not full-strength
  attribution.
- My round-2 falsifier was explicit: I would move to (a) if a committed values-only mechanism
  preserved live machine membership and no real add/remove need remained after attended migration.
  That condition is why §2.1 is decisive, not merely an optional implementation detail.

**PRIMARY — direct textual audit of `round-02/claude-1.md`:** Claude explicitly rejects my
`add`+`remove` v1 and proposes values separation plus `add` only. The row listing both of us under
one (c) is acceptable only if §2 records this intra-(c) disagreement. It currently does not.

4. **The decisive unrun experiment (§2.1)**

**(i) Technical possibility and breakage.**

**PRIMARY — direct source inspection, with no prototype claimed:** it is technically possible.
`internal/config/runtime.go:23-49` parses every `[roster.<id>]` into one map; `:135-169` merges the
same blocks as values and unconditionally adds every committed-deck block ID to `deckMembers`; and
`:182-185` makes any non-empty `deckMembers` set authoritative. Those are separable code paths.

It is **not** safe to separate them by guessing from which fields a block happens to contain.
Existing full declarations also carry value fields, and D-C deliberately removes fields while
retaining membership headers. Any rule such as “speed/model-only means values-only” would make a
file's membership meaning change when `sync` removes `adapter`, and a blanket reinterpretation
would silently expand the 37 full-declared decks to their machine rosters.

A backward-compatible experiment needs an explicit discriminator but does not need add/remove set
arithmetic. One concrete candidate is:

```toml
[membership]
mode = "inherit-values-v1"

[roster.zcode-1]
speed = "fast"
```

In that explicit mode, machine membership remains authoritative and `[roster.*]` supplies values
only. Absence of the mode retains today's full-declaration, legacy-§2, and pure-inheritance
branches byte-for-byte. `active` must be rejected in this mode because it is a membership
operation, and a block for an ID absent from the machine base must fail as a dangling override
rather than add a member.

This is not a free wording patch. It changes the parser/resolver, D-A's write and preview path,
`sync`, `show --explain`, render/legacy handling, and tests. It deliberately does **not** express a
deck-only addition or a machine-member removal; those remain the evidence trigger for a later
overlay. The 37 existing full-declared decks break only if the discriminator is implicit or they
are auto-converted, both of which the experiment must forbid.

**Concrete counter-proposal:** do not escalate the (a)/(c) choice and do not draft FINAL. Open
round 3 around an isolated-copy prototype of `inherit-values-v1` and require all of these results:

1. An inherited six-member fixture remains six after a committed speed override.
2. A later machine-member addition/removal is reflected without rewriting the deck.
3. Unmarked full-declared fixtures, including blocks with value fields and blocks stripped by
   `sync`, retain their exact current member sets.
4. Unmarked legacy §2 remains authoritative; explicit transition away from it is attended.
5. `active`, unknown IDs, and ambiguous mixed modes fail closed.
6. D-A prints identical resolver before/after sets for the value-only write; D-C distinguishes
   values from membership; `show --explain` names the source truthfully.
7. The existing roster/config/app suites and D-B's render-then-protocol-guard check pass.

Then every participant should respond to that evidence in its own round-3 file and the consensus
should be redrafted.

**(ii) Would a positive result change my answer? Yes.** A result satisfying the matrix above would
serve the only demonstrated live requirement—committed local values while tracking machine
membership—without shipping add/remove semantics. I would move from (c) to (a): ship D-A/D-B/D-C
and the narrowly separated values-only mechanism, migrate attended decks that choose inheritance,
and defer the membership overlay until a real named add/remove case appears. The split is therefore
premature, and this idea must run the experiment before FINAL.

5. **Anything in the draft I will not sign**

I will not sign:

- §7's owner escalation before the §2.1 experiment; the local, cheaper falsifier must run first.
- The present D-C fix shape without truthful value-vs-membership output and post-resolve
  acceptance criteria.
- “No mass conversion, ever” when the intended rule is no automatic or inference-based
  conversion.
- A single undifferentiated (c) design attributed jointly to Claude and me; our v1 operation sets
  differ.
- Any Phase 4 close that treats §5 as a quorum resolution. **PRIMARY — direct artifact check:**
  `00-prompt.md` still lists `opencode-1`, while no `round-01/opencode-1.md`,
  `round-02/opencode-1.md`, or signoff exists. The draft accurately records incomplete
  participation, but incomplete participation is not consent. Before FINAL, `opencode-1` must
  contribute or a protocol-valid, explicitly recorded quorum resolution must occur.
