### codex-1 — revision 2

**Verdict:** block

#### Scope

- `PRIMARY` — I read the live `parley-deck/COOPERATION.md:1-1316` in full, including the binding verification rules at `:1176-1316`, the track classifier at `:172-228`, the file-ownership rule at `:704-715`, and the protocol-change rule at `:717-730`.
- `PRIMARY` — I read the complete revision-2 `parley-deck/ideas/roster-operations-standard/consensus.md:1-938`, including all four embedded revision-1 signoffs at `:425-936`, and independently re-read my standalone revision-1 signoff at `parley-deck/ideas/roster-operations-standard/signoff-codex-1.md:1-60`.
- `PRIMARY` — I also read `parley-deck/ideas/roster-operations-standard/00-prompt.md:1-85`, `parley-deck/agents.toml:1-75`, and the existing roster configuration shape at `internal/config/runtime.go:22-35,196-245` solely to scope the unresolved authority contract. I ran the read-only skill status/hash checks and file searches/reads. I did not run tests, launch agents, enumerate or inspect the 40 foreign decks, execute a migration dry-run, or run any roster or Git write command. I wrote only this signoff file.
- `PRIMARY` — I do not issue `CONFIRMED` or `WRONG` on the run-manifest or continuation claims I already own, as recorded in `signoff-codex-1.md:11-23`; §15.1 prohibits that (`COOPERATION.md:1197-1205`). This signoff evaluates whether revision 2 contains the release conditions I required, not whether my owned source claims are true.

#### My four required changes from revision 1

1. **Rebase/snapshot and legacy-run gate — met.** `PRIMARY` — Revision 2 says the change exposing rebase “MUST also persist and consume the immutable effective snapshot,” requires the configuration-change/continue acceptance test, forbids rebase shipping first, skips nonterminal pre-snapshot legacy runs without rewriting their artifacts, and requires the explicit “unsafe for pre-snapshot resumable runs” warning if the gate is relaxed (`consensus.md:305-315`). That matches `signoff-codex-1.md:19-23,57`.

2. **Correct §7 wording, correct track, and complete §2 authority contract — not fully met.**

   - `PRIMARY` — The track correction is met: frontmatter now says `track: deliberation` (`consensus.md:1-9`), and `:370-374` applies the protocol-change, data-migration, and irreversible-operation triggers from `COOPERATION.md:179-190`.
   - `PRIMARY` — The §7 correction is met: revision 2 calls the venue an explicit user-authorized one-off, says it is not a protocol exception, and says it creates no precedent (`consensus.md:354-368`). That matches the distinction in `COOPERATION.md:704-730` and my requirement at `signoff-codex-1.md:25-31`.
   - `PRIMARY` — The §2 authority contract is still not met. `consensus.md:326-340` says that, “Before ratification,” the change **MUST define** the canonical source and migration for workspace dir, role, host handle, active/inactive history, and ordering. It does not then supply those definitions. The general statement that `parley-deck/agents.toml` becomes the deck authority (`:320-324`) does not specify the keys/schema for those fields, how each existing §2 value migrates or conflicts are handled, how inactive history is retained, or what deterministic row-order rule the generator uses. Nor does revision 2 contain the promised replacement §2/protocol text. This repeats my pre-ratification requirement as a TODO instead of satisfying `signoff-codex-1.md:33,58`.
   - `PRIMARY` — Deferring those answers to `FINAL.md` cannot cure this signoff gap: consensus signoff is the Phase-3 ratification gate (`COOPERATION.md:351-368`), while `FINAL.md` is drafted only afterward in Phase 4 (`:370-398`).

3. **Fleet migration contract — met for my revision-1 requirements.** `PRIMARY` — Revision 2 now requires the exact-root/source-revision/hash/version/worktree inventory (`consensus.md:393-395`), compare-and-swap plus explicit approval (`:396-398`), a precise unclean/skip definition (`:399-401`), recorded file-level backups, verified restore, atomic writes, validation, and automatic rollback (`:402-405`), and a per-deck final machine-readable result with hashes and backup reference while prohibiting implied commits, pushes, or historical-artifact edits (`:411-414`). These cover `signoff-codex-1.md:39-45,59`.

4. **VC-1 and VC-3 closure — met.** `PRIMARY` — VC-1 is closed because kimi-1 withdrew `SOURCE` through the field-specific-provenance argument, not because of the vote count; the eleven-column contract stands (`consensus.md:155-165`), and kimi-1's own embedded reasoning is preserved at `:794-806`. VC-3 records unanimous `deck|machine` and the committed `parley-deck/agents.toml` target (`:222-230`). This matches `signoff-codex-1.md:47-53,60`.

#### hermes-1 reservations

- **R1 — addressed.** `PRIMARY` — Snapshot persistence **and consumption**, rebase atomicity, and the acceptance test are binding at `consensus.md:305-315`, answering hermes-1's R1 at `:544-550`.
- **R2 — addressed.** `PRIMARY` — The generated §2 must be byte-identical on a second run and preserve the human-readable form (`consensus.md:350-352`), answering R2 at `:575-580`.
- **R3.1 — addressed.** `PRIMARY` — Inactive-set/`STATE` wiring is a hard prerequisite and must ship with migration (`consensus.md:342-348`).
- **R3.2 — addressed.** `PRIMARY` — Apply requires per-deck or small-batch confirmation and honors `confirm-breaking` (`consensus.md:406-407`).
- **R3.3 — addressed.** `PRIMARY` — Migration is human-attended only and prohibited from loop, cron, or CI execution (`consensus.md:410`).
- **R3.4 — addressed.** `PRIMARY` — Backups are file-level copies rather than Git operations and carry hashes plus verified restore/rollback requirements (`consensus.md:402-405`).
- **R3.5 — addressed.** `PRIMARY` — Migration is resumable after a partial batch and already-migrated decks are no-ops (`consensus.md:408-409`).

#### kimi-1 reservations and VC-1

- **R1 — addressed.** `PRIMARY` — The effective snapshot must be persisted and consumed in the same delivery as rebase; nonterminal pre-snapshot runs are skipped; and relaxing the gate requires the unsafe-resume warning (`consensus.md:305-315`). This meets the release-gate substance at kimi-1's embedded signoff `:837-841`.
- **R2 — only partly addressed.** `PRIMARY` — Revision 2 correctly chooses a recorded one-time deviation rather than amending §7 (`consensus.md:354-368`), and requires a non-authoritative generated §2 plus no runtime parsing (`:326-340`). It still lacks the actual replacement protocol text and does not require the §7-format `meta/protocol-changelog.md` entry that kimi-1 expressly requested at `:842-846`.
- **R3 — only partly addressed.** `PRIMARY` — Locked participants/run artifacts, machine-readable skip classes, restore/re-resolution, non-Git backups, no implied commit/push, and fleet confirmation are covered at `consensus.md:311-312,393-414`. Two explicit requirements from kimi-1's R3 at `:847-851` remain absent: retired agents must be retained as `inactive` and never deleted, and each foreign deck must pass a compatible-protocol/sync gate rather than merely report its protocol version. `consensus.md:342-348` describes migrating the known retired rows to inactive, but does not define the retention rule; `:393-401` inventories protocol/schema versions but does not define the compatibility gate.
- **R4 — not addressed.** `PRIMARY` — Revision 2 neither adopts `--keep` nor makes per-deck enumeration of every deliberate pin removed by rebase mandatory. The inventory/diff requirements at `consensus.md:393-414` do not classify and preserve that checklist, which kimi-1 required at `:852-853`.
- **VC-1 closure — recorded as kimi-1 argued it.** `PRIMARY` — `consensus.md:155-165` attributes closure to kimi-1's withdrawal and the argument that a row-wide `SOURCE` is incoherent while a model-only source column is not worth permanent API width. The embedded signoff at `:794-806` preserves the same reasoning and expressly says the result is argument-based rather than 3-to-1.

#### Delivery shape

**Position: stage the implementation; do not review or land it as one monolithic change.** `PRIMARY` — The agreed surface spans the effective-value resolver, frozen table/JSON API, model metadata, commands, run state, generated protocol text, skill/docs, and a destructive fleet operation (`consensus.md:59-151,305-419`). Reviewable stages are safer, but the atomic groups below are release gates rather than optional sequencing advice.

1. **Internal roster foundations:** implement and test the `{model}`/`{effort}` placeholder resolver and legacy normalizer, the `modelmeta` registry, resolved-row types, the versioned eleven-column/JSON schema, and active/inactive `STATE` consumption. These may land behind the existing surface, but the public effective-value contract must not be exposed until the resolver and `STATE` semantics are wired.
2. **Authority cutover and ordinary operations — one atomic group:** finalize the complete committed-TOML schema; migrate every §2 field; implement `roster show`/`set`; generate §2 idempotently; remove all runtime parsing of generated §2; and update the live protocol, embedded protocol copy, bundled skill snapshot, skill behavior, CLI help, and docs. The authority cutover, generator, runtime consumer cutover, and protocol/skill text must land together or remain feature-gated together.
3. **Snapshot plus rebase — one atomic group:** persist the immutable effective row and `roster_revision`, consume that snapshot on every continuation, add the configuration-mutation continuation test, and only then expose rebase semantics in `roster sync`. Snapshot persistence without consumption is insufficient; rebase must not land or release separately (`consensus.md:305-310`).
4. **Migration tooling and attended fleet execution:** add inventory, dry-run report, compare-and-swap, file backup/verified restore, per-deck rollback, resumability, compatibility/skip gates, and final report after stages 1-3 are complete. The actual 40-deck mutation is a separate attended operation: present the frozen dry-run to the user, apply only approved decks/small batches, and never fold the fleet mutation into the code/protocol merge.

Stages 1-3 may be separate implementation PRs if unfinished user-visible behavior stays disabled. They should form one coordinated CLI/protocol/skill release before migration is enabled. Stage 4 depends on that released version and must not start early.

#### Counter-proposal required by this block

Before requesting revision-3 signoff:

1. Replace the requirement-only paragraph at `consensus.md:326-340` with a normative field table that names, for every roster field, its exact committed TOML key, its legacy §2 source, conflict/absence behavior, and whether it is runtime-semantic or render-only. Define inactive-history retention (“mark inactive; never delete”) and one deterministic ordering rule. Include the proposed §2 replacement text now, because the consensus itself says this is required before ratification.
2. Require the §7-format `meta/protocol-changelog.md` entry naming this idea and the user-authorized one-off.
3. Add a foreign-deck protocol compatibility/sync gate and make retired-agent retention explicit in the migration contract.
4. Resolve kimi-1 R4 by either adding `--keep <agent>.<field>` or requiring the dry-run/final report to enumerate every removed deliberate pin per deck.

The rebase decision, eleven-column contract, `deck|machine` labels, hermes-1 guardrails, and migration safety machinery do not need reopening. With the four additions above, my revision-1 block would be fully discharged.
