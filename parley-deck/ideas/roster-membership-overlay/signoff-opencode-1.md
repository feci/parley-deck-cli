### opencode-1 — 🟡 accept with reservations

First voice in this idea. I filed no round-1 or round-2 artifact (provider resets; owner then ran round 2 without me). That is incomplete participation, not consent. I will not rubber-stamp a 3–2 split I never joined. §5 of the draft records this correctly; I sign that sentence and nothing inferred from my silence.

1. **§1 (the unanimous block)**

I confirm the three defects and that each is its own independently releasable `standard`-track idea. Shared tree untouched for every check below (`git status --porcelain` empty after).

- **D-A / D-B mechanism — PRIMARY, code read this session.** `membershipChange` (`internal/app/roster_set.go:287-290`) keys on whether the block existed in the file; `LoadRosterScoped` (`internal/config/runtime.go:182-186`) keys on `len(deckMembers) > 0`; `rosterLayer` (`runtime.go:223-239`) treats every parsed `[roster.<id>]` key as a member even when the table is empty. Renderer emits `| Agent ID | Workspace dir | Role | State |` (`roster_render.go:73`); the drift guard anchors `| Agent ID       | Workspace dir                       | Role          |` (`drift_test.go:28`). I did **not** re-run the D-A/D-B CLI reproductions this session; those remain SECONDARY on @kimi-1/@claude-1/@codex-1/@zcode-1 PRIMARY.
- **D-C — PRIMARY, this session.** Isolated copy `/tmp/opencode-signoff-dc-47767` (five `[roster.*]` blocks, adapter-only, empty §2 so rule 2 cannot hold membership). BEFORE: five rows, STATUS `ok` (deck-declared), no `zcode-1`. `parley roster sync --dir … --yes` printed the preview lie ("removing these makes the deck inherit") and the success lie ("the deck now inherits"). AFTER: same five rows, still `ok`, not `inherited-roster`; file reduced to five empty headers (`[roster.claude-1]` … `[roster.opencode-1]`); `grep -c '^\[roster\.'` still 5. `zcode-1` still absent. This matches @claude-1's inbox note and adds the empty-header residue @zcode-1 reported.
- **D-C is a designed field-rebase, not a failed membership migration — PRIMARY.** `roster_sync.go:13-24` states rebase-of-fields; `:64-75` iterates only IDs the deck already declares; `:174-205` `removeRosterField` leaves the header; `roster_sync_test.go:39-41` **requires** `[roster.claude-1]` to survive. The bug is the stronger claim in `:114` and `:169-170`, not a missing delete.

**§1.2 shapes.** I sign D-A (gate must print the resolver's before/after member sets; silent materialisation of the inherited set is not the default) and D-B (renderer, guard anchor, embedded default change atomically; which shape is canonical is not this idea's). I sign D-C's "must not report an outcome it did not achieve" as necessary and **not sufficient**: a truthful `sync` that still only rebases fields is still not a migration instrument. Counter-shape: distinguish value inheritance from membership authority, print both, post-resolve before claiming success; membership adoption is a separate attended verb. I do not read §1.2 as forbidding an explicit, previewed materialise path.

**§1.3.** I sign the body (no auto-conversion; an omission is never inferred to be intentional). I agree with @codex-1 that the heading "No mass conversion, ever" over-claims: attended per-deck adopt-machine vs preserve-set is not mass conversion.

**§1.4.** I sign "fleet migration onto inheritance is desirable and currently has no working instrument" as a consequence of D-C. I do **not** sign the round-1 attribution "@hermes-1, @kimi-1 and @zcode-1 each proposed it independently in round 1" as to @kimi-1 (SECONDARY — I re-read `round-01/kimi-1.md`: its approach is F1/F2 + marker + deferred trigger; fleet migration is a round-2 adoption). Correct: @hermes-1 and @zcode-1 in round 1; @kimi-1 in round 2.

**§1.5.** I sign the compatibility boundary **if** an overlay is ever built.

2. **D-C**

Confirmed PRIMARY (item 1). I was not briefed; I take a position.

It does not make me want the overlay. It is the third honesty defect in the three verbs §2 names, and the test suite currently locks in the membership-preserving behaviour the success text denies.

**Sequencing, stated as a first position rather than a change:** fix D-A, D-B and D-C first, each independently releasable. Do not use today's `roster sync` as the fleet-migration instrument @hermes-1 and @zcode-1 proposed. A separate attended inherit/adopt verb must preview exact before/after member sets, clear both `[roster.*]` **and** a valid legacy §2 if that is the authority, and post-check `Inherited == true`. Dirty/non-Git roots stay attended (SECONDARY on their censuses). D-C blocks migration; it does not decide (a) vs (c).

3. **§2 — is MY side stated at full strength?**

I have no side in the draft. Every sentence that names me is in §5 and is accurate (SECONDARY — I re-read the draft, the incomplete-participation inbox note, and all ten round files). The 3–2 table correctly omits me. That table must not be read as 3–2 of a six-agent quorum.

Because I have not been heard, the rest of this item **is** my side.

**I do not sign (b). I do not sign closing this idea by counting (a) against (c).** Inside the draft's own frame I would have filed **(a)**: fix the gestures, keep the authority model, overlay deferred behind a trigger that can actually fire (value-override **or** membership ±1; ≥2 attended "deliberate" answers **or** an explicit owner instruction). I would not have filed (c). Reasons, at full strength:

- The coupling D-A exposes is the ratified rule, not a hole the overlay is needed to fill. `runtime.go:90-103` and `:174-186` say membership is the deck file; `rosterLayer` makes presence the discriminator. D-A proves the **gesture** mislabels that rule. It does not prove the rule is wrong.
- Every census, including mine this session, measured demand for membership deltas at zero. **PRIMARY — read-only walk of `/Volumes/My Shared Files` (maxdepth 4):** 38 `parley-deck/agents.toml` files; 1 with zero `[roster.*]` blocks (this deck); 37 declared; **0 empty-stub files**; all 37 declared decks carry at least one value key (spot-check: `adapter`-only). No post-D-C empty-header signature on the volume today. Caveat: walk denominators include worktrees; I report the shape, not a sacred 37.
- **(c)'s live-demand paragraph is one-sided.** It cites the owner's originating sentence and "this deck hit it today." The same file's dated instruction of record (`parley-deck/agents.toml:66-75`, PRIMARY, read this session) says *"lokalne nepretazuj nic, pouzivaj globalny roster"* and explicitly distinguishes inheritance from `roster sync`. Consensus must carry both sentences or neither.
- The (c) row collapses two designs. @codex-1's v1 is `add` **and** named `remove` with tombstones; @claude-1's round-2 v1 is values-separation plus `add` only (SECONDARY — both files re-read). That disagreement is load-bearing and is not in the draft.
- (c)'s "census of a syntax that does not exist" point is real. It does not license building the larger mechanism before the cheaper committed values-only experiment (§2.1) runs.

I also will not let (a) be recorded as "a trigger that cannot fire." That objection lands on @hermes-1's passive wording and on @zcode-1's **round-1** wording. @zcode-1's round-2 trigger already includes the value-override limb **and** an owner-instruction disjunct (SECONDARY — `round-02/zcode-1.md`). @kimi-1's instrumented trigger (record deliberate-vs-stale at migration) can fire. Credit those, or do not paraphrase (a).

**Path C (binding owner ruling, `inbox/user-to-all_roster-membership-overlay_uniform-inheritance-path-c.md`).** The Phase 3 briefs predate it; the ruling says in-flight signoffs are not void and that (a)/(c) become input to Path C. I treat that as binding on me.

I accept the **direction**: load parent, apply child overrides property by property, inherit what the child does not mention; membership is a property, not a second authority. I object, as engineering, to making that the unmarked default overnight.

**PRIMARY this session — `parley agents list`:** `codex` already resolves `sandbox`/`approval`/`timeout` from the deck and `model` from `~/.parley/agents.toml`. That is Path C for `[agents.*]`. Membership is the inconsistent property. @claude-1's later inbox note claimed this; I re-ran the command rather than cite it.

Path C does **not** retire D-A/D-B/D-C. It changes the end state those fixes aim at: a values-only `[roster.<id>]` block must stop constituting membership. That is exactly §2.1, promoted from experiment to destination. It also recreates the 37-quorum silent-expansion hazard @claude-1's round-1 measurement named, which @codex-1's opt-in stanza was designed to avoid. **A versioned deck marker (`schema = 2` or equivalent) is the only migration I would sign.** Unmarked files keep today's replacement / legacy / inherit branches byte-for-byte until attended conversion. That is §1.3 applied to Path C. Without the marker, Path C is the unmodified-file sin this idea already swore off.

On `members` as a list: replacement (`members = […]`) is enough for v1. `super.members ± [x]` is @codex-1's `add`/`remove` under another name; I would not ship it until a named deck needs it. Scope beyond the roster: C already describes `[agents.*]` (PRIMARY above); do not silently extend it to `active` (authority-anchored at `runtime.go:127-132`) without saying so.

Path C **strengthens** this 🟡: I will not sign a FINAL that still presents (a) vs (c) as the decision.

4. **The decisive unrun experiment (§2.1)**

**(i) Technically possible?** Yes. Values already merge across layers (`runtime.go:143-146`); membership is a later presence test (`:166-185`); `rosterLayer` returns every map key (`:234-239`). Those paths are separable in code.

Two shapes, only one is compatible with §1.3:

- **Implicit / content-keyed** ("a block with only `speed` is values-only"): **not safe.** PRIMARY census above: all 37 declared decks' blocks carry values. D-C then strips those values and leaves empty headers that still constitute membership (PRIMARY, my isolated run). A content rule would flip unmodified files — the 37 full-declared decks would inherit the machine set and gain `zcode-1`. Forbidden.
- **Explicit discriminator** (new stanza / schema marker; `[roster.*]` becomes values-only **in that mode**): back-compatible. It is @codex-1's recorded separation minus add/remove. The 37 break only if the discriminator is implicit or they are auto-converted.

The gitignored `agents.local.toml` path already separates values from membership (`runtime.go:150-164`, `membership: true` only on the committed deck file at `:389` — PRIMARY, code). That is not the committed, CLI-writable path the owner asked for, and documenting it as the workaround would push decisions into the least auditable layer.

**(ii) Would a positive result change my (a)/(c) answer?** I had no prior answer. A positive **explicit** result would make me refuse (c): the live committed-value-while-inheriting case is then absorbed inside the fix track, and the overlay's residue is membership ±1, whose measured demand is zero. A positive **implicit** result I would treat as a failed experiment, not a win. **If this idea is about to pick (c) or to escalate (a)/(c) to the owner, the split is premature — run §2.1 in an isolated copy before FINAL.** Path C makes that run mandatory rather than optional: it *is* the experiment, plus the marker vs silent-flip question.

5. **Anything in the draft I will not sign**

- Closing or escalating on (a)/(c) as if that were still the question, after the binding Path C ruling and with §2.1 unrun.
- §1.4's round-1 migration attribution to @kimi-1 (item 1).
- "No mass conversion, ever" as a heading that appears to ban attended per-deck conversion (item 1).
- A single undifferentiated (c) attributed jointly to @codex-1 and @claude-1 (item 3).
- "@zcode-1's trigger requires '≥2 real deck instances' … a trigger that cannot fire" **as a description of @zcode-1's round-2 position** (item 3).
- "This deck hit it today" as live demand without `agents.toml:66-75`'s dated nothing-local instruction (item 3).
- Any reading of my absence as consent, as a vote, or as a demonstration of failure. §5 as written is the only sentence about me I sign.
- Using `roster sync` as the fleet-migration instrument.
- Implicit reinterpretation of existing `[roster.*]` lists — under the overlay, under Path C, or under a "values-only" D-A fix.

**Verification.** D-C reproduction, live `parley roster show` / `parley agents list`, read-only fleet walk, `parley-deck/agents.toml`: PRIMARY this session. Code locators quoted above: PRIMARY, read this session. Round artifacts, D-C inbox, Path C ruling, peer signoffs: SECONDARY, read this session. No RECALL claim carries weight. Shared working tree was not used as an experiment root.
