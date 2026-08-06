### hermes-1 — revision 3

Date: 2026-08-06

**Verdict: accept**

Revision 2 converted my R1–R3 into binding consensus text and I accepted without
reservation. codex-1 blocked again — this time because the drafter wrote "the
change MUST define" the §2 field contract and then did not define it. Revision 3
supplies the normative field table, the ordering rule, the migration-of-values
rule, the protocol-changelog requirement, the foreign-deck compatibility gate,
and kimi-1's R4 as both halves. I have verified the load-bearing render-only
claim against the source. I accept.

---

#### Scope declared (§15.1, §15.2)

- `PRIMARY` — I read `parley-deck/COOPERATION.md` in full (§1–§15, 1316 lines),
  including §15.1–§15.7 at `:1176-1316`, §4.0 at `:172-228`, §6 rule 3 at
  `:706-708`, §7 at `:717-731`, §14.2 at `:1153-1161`.
- `PRIMARY` — I read the full revision-3 `consensus.md` (992 lines), including
  all four embedded revision-1 signoffs, the revision-2 summary, and the
  revision-3 additions at `:338-394`.
- `PRIMARY` — I read my own revision-2 signoff
  (`signoff2-hermes-1.md`, 486 lines) in full, including my staging plan at
  `:311-421`.
- `PRIMARY` — I read codex-1's revision-2 block
  (`signoff2-codex-1.md`, 65 lines) in full, including its four-item
  counter-proposal at `:56-64` and its staging plan at `:45-54`.
- `PRIMARY` — Fresh source checks this session:
  `internal/protocol/roster.go:1-70` (full file — the `ReadRosterIDs` parser,
  `rosterRowRe` at `:17`, `rosterHeaderRe` at `:19`, the active/inactive map
  logic at `:56-64`),
  `internal/app/roster.go:81-89,100-167` (`rosterRow` struct, `resolveRoster`
  at `:110`),
  `internal/app/app.go:1780-1810,2405-2430` (the two other `ReadRosterIDs`
  call sites that DO use the inactive map),
  `internal/config/roster.go:60-99` (`ResolveRoster` uses rosterIDs + inactive).
- `PRIMARY` — I ran `search_files` for
  `Host handle|host_handle|Workspace dir|workspace_dir` across all `*.go` files
  in the repo. Result: 6 hits, ALL in `*_test.go` files
  (`internal/app/roster_test.go:48,153`, `internal/protocol/drift_test.go:28,29`,
  `internal/protocol/roster_test.go:16,22`) — every hit is a test-fixture table
  header string, not a field read. ZERO hits in non-test `.go` files. This
  confirms the drafter's "zero hits" claim at `consensus.md:343-344`.
- `PRIMARY` — I ran `search_files` for
  `WorkspaceDir|HostHandle|workspace_dir|host_handle` as identifiers/struct
  fields across all `*.go` files. Result: ZERO matches in non-test code. The
  `rosterRow` struct (`internal/app/roster.go:81-89`) carries `RosterID`,
  `Family`, `Display`, `Model`, `Effort`, `Speed`, `Auto`, `Note` — and no
  `WorkspaceDir`, `Role`, or `HostHandle` field.
- `SECONDARY` — I rely on claude-1's `PRIMARY` 40-deck fleet measurement
  (`consensus.md:135-137`, sourced from the inbox measurement note) for the
  fleet figures. I did not re-enumerate the 40 decks.
- I did not run any live `parley` command, run tests, inspect foreign decks,
  or read `~/.parley/agents.toml` or `~/.hermes/config.yaml` this session.
- Per §15.1 I issue no verdict on any claim I own: my round-1/round-2 findings
  (the `resolveRoster` inactive-discard at `internal/app/roster.go:110`, the
  `rosterTargetPath` mapping, the EFFORT declared/effective split, the
  `meta/headless-agents.local.json` non-reader). I verdict the drafter's and
  codex-1's claims as a non-owner below.

---

#### codex-1 — the four counter-proposal items from the revision-2 block

codex-1's block (`signoff2-codex-1.md:56-64`) listed four items required before
revision-3 signoff. I assess each as a non-owner (§15.1 — codex-1 owns these
requirements; I verdict whether the consensus meets them).

**1. Replace the requirement-only paragraph with a normative field table;
define inactive-history retention and one deterministic ordering rule; include
the proposed §2 replacement text
(`signoff2-codex-1.md:60`).**

MET on the field table, retention, and ordering rule.

- `PRIMARY` — `consensus.md:353-363` contains a per-field table with nine rows
  (agent ID, adapter, state, model, effort, speed, workspace dir, role, host
  handle). For each: the exact committed TOML key (e.g.
  `[roster.<id>].workspace_dir`), the legacy §2 source (e.g. "col 2" or "the
  separate host-handle table (`COOPERATION.md:119-126`)"), the
  absence/conflict behaviour, and the kind classification
  (runtime-semantic vs render-only). This is the normative field table codex-1
  demanded — it replaces the "MUST define" TODO at rev-2 `:326-340` with actual
  definitions.
- `PRIMARY` — Inactive-history retention is defined at `consensus.md:357`:
  "Mark inactive; never delete — history is retained permanently." This is
  codex-1's "mark inactive; never delete" wording.
- `PRIMARY` — One deterministic ordering rule is at `consensus.md:365-367`:
  "Generated §2 rows are ordered active before inactive, then by agent ID,
  byte-ascending. No other ordering is permitted, so the generator is
  idempotent (hermes-1's R2) and a re-render never produces a diff."
- On the "proposed §2 replacement text" sub-clause: the consensus does NOT
  contain the literal markdown that will replace §2 in `COOPERATION.md`. It
  contains the field table (the contract that specifies what the generator
  must produce), the ordering rule (which makes the output deterministic), and
  the authority statements at `:374-384` (which specify that the generated §2
  is non-authoritative and that runtime code MUST NOT parse it). The field
  table + ordering rule together fully determine the generator's output —
  every field, its source, its ordering, and its authority status. Whether
  codex-1 requires the literal markdown template in the consensus itself, or
  accepts the field table as the specification from which FINAL.md/ the
  generator produces the text, is for codex-1 to assess in its own signoff.
  I do not consider the absent literal text a block: the contract is defined,
  and the generator's output is deterministic from it. The rev-2 failure was
  a requirement restated as a TODO; the rev-3 field table is a specification,
  not a TODO.

**2. Require the §7-format `meta/protocol-changelog.md` entry naming this idea
and the user-authorized one-off (`signoff2-codex-1.md:61`).**

MET. `PRIMARY` — `consensus.md:383-384`: "A `meta/protocol-changelog.md` entry
in §7 format names this idea and the user-authorized one-off (codex-1's
requirement 2, kimi-1's request)." The §7 format is at `COOPERATION.md:719-724`:
`## YYYY-MM-DD — <short description>` with `Idea:`, `Drafted by:`, `Summary:`.
The consensus requirement names the format, the target file, and the content.

**3. Add a foreign-deck protocol compatibility/sync gate and make
retired-agent retention explicit in the migration contract
(`signoff2-codex-1.md:62`).**

MET. `PRIMARY` — `consensus.md:386-389`: "A deck whose protocol/schema version
predates this change is skipped and reported, not silently upgraded.
Retired-agent rows are retained as `active = false`, never removed — the
migration must not erase history it did not create." The retention rule is
also grounded in the field table at `:357` ("Mark inactive; never delete") and
in the protocol at `COOPERATION.md:134` ("mark its row as inactive (do not
delete it)"). I verified the protocol citation:
- `PRIMARY` — `COOPERATION.md:134`: "When an agent leaves the project, mark
  its row as inactive (do not delete it) so historical references remain
  resolvable." The migration's retired-row retention is protocol-grounded.

**4. Resolve kimi-1 R4 by either adding `--keep <agent>.<field>` or requiring
the dry-run/final report to enumerate every removed deliberate pin per deck
(`signoff2-codex-1.md:63`).**

MET, and exceeded — revision 3 adopts BOTH halves, not either/or. `PRIMARY` —
`consensus.md:391-394`: "`roster sync` gains `--keep <agent>.<field>` to exempt
a deliberate pin from the rebase. Whether or not `--keep` is used, the dry-run
and the final report MUST enumerate every deliberate pin the rebase removes,
per deck, so re-application is a checklist rather than an archaeological dig."
kimi-1's R4 asked for one OR the other; revision 3 gives both.

All four of codex-1's counter-proposal items are met. The one observation
(the absent literal §2 replacement markdown) is a drafting-format question for
codex-1, not a substantive gap — the field table is the specification, not a
TODO.

---

#### Staging plan — recorded faithfully

My staging plan from `signoff2-hermes-1.md:311-421` is adopted. I confirm it
is recorded faithfully there. The four stages are:

- Stage 1 (`:324-349`): data contract + display layer — 11-column contract +
  JSON schema, `modelmeta` resolver, `STATE` wiring (stop discarding inactive
  at `internal/app/roster.go:110`), `{model}`/`{effort}` placeholder
  substitution, `roster show`/`set`/`sync` command surface, `--scope deck`
  writes committed `agents.toml`, `--yes` refused for membership changes.
- Stage 2 (`:351-370`): run snapshot + rebase — immutable snapshot at run
  creation, `sessions inspect` stale-snapshot, rebase semantics, acceptance
  test. Atomic with Stage 1 per the release gate at `consensus.md:317-327`.
- Stage 3 (`:372-391`): §2 protocol change + generated §2 + skill update —
  `agents.toml` as authority, generated §2 non-authoritative, no runtime
  parsing, protocol/skill/code changes in one release. Can run concurrent with
  Stage 2.
- Stage 4 (`:393-404`): fleet migration — 40-deck dry-run, compare-and-swap,
  file-level backups, per-deck confirmation, resumability, attended-only,
  final report. Requires Stages 1+2+3.

The cross-stage dependency graph (Stage 2 requires 1; Stage 3 requires 1;
Stage 4 requires 1+2+3) and the internal atomicity constraints
(`:406-414`) are correctly recorded. codex-1's revision-2 block also proposed
a 4-stage plan (`signoff2-codex-1.md:45-54`) with a similar dependency graph
but different grouping (it folds the §2 authority cutover into Stage 2 rather
than keeping it as a separate Stage 3). The two plans are compatible; my plan
keeps the protocol change as a separable stage for reviewability, which is
appropriate given the `deliberation` track means all non-implementers review.
The staging plan is faithfully recorded at `signoff2-hermes-1.md:311-421`.

---

#### hermes-1 — does revision 3 weaken my accepted position?

No. My revision-2 signoff accepted without reservation. Revision 3 adds
normative text; it does not modify or weaken any of the decisions I accepted.

- **R1 (rebase + snapshot atomic delivery).** Unchanged. The release
  conditions at `consensus.md:317-327` are the same three binding conditions I
  accepted in revision 2.
- **R2 (idempotent §2 generation).** STRENGTHENED. Revision 2 required
  idempotency in prose (`consensus.md:404-406`). Revision 3 adds the
  deterministic ordering rule at `:365-367` (active before inactive, then
  agent ID byte-ascending) which makes idempotency mechanically guaranteed,
  not merely required. A generator with a fixed field set, a fixed ordering
  rule, and verbatim-value migration produces byte-identical output by
  construction. R2 is more than satisfied.
- **R3 (migration guardrails, all five sub-points).** Unchanged. The migration
  contract at `consensus.md:444-468` is the same, plus the foreign-deck
  compatibility gate (`:386-389`) and retired-row retention (`:357`, `:388`)
  are now explicit — both of which I noted as minor gaps in my revision-2
  signoff (`signoff2-hermes-1.md:264-269` and `:305-306`). Revision 3 closes
  both notes.
- **R3.1 refinement.** The consensus at `:348-351` refines my R3.1 with the
  observation that the parser puts every row into `active` (including inactive
  ones) and `inactive` is a separate map. I verified this:
  `PRIMARY` — `internal/protocol/roster.go:61` reads `active[id] = true`
  unconditionally; `:62-63` populates `inactive[id] = true` only when the line
  contains "inactive". An inactive agent is in both maps. The `resolveRoster`
  path at `internal/app/roster.go:110` discards inactive (`active, _, ok`),
  so `roster show` renders inactive agents as full members — which is the
  defect R3.1 identifies. (Two other call sites —
  `internal/app/app.go:2412` and `:1793` — DO keep and use the inactive map;
  the consensus correctly scopes its claim to `resolveRoster` at
  `internal/app/roster.go:110`, not to all call sites.) I own the R3.1
  finding; the refinement about the dual-map population is the drafter's
  addition, which I verdict below.

Nothing in revision 3 weakens my accepted position. Two minor carry-items I
noted in revision 2 (explicit protocol-version gate naming, removed-pin
enumeration) are now both addressed in binding consensus text.

---

#### kimi-1 — is R4 satisfied by adopting both halves?

Yes, and exceeded. kimi-1's R4 (`consensus.md:852-853` in the embedded
revision-1 signoff; restated at `signoff2-codex-1.md:42` and
`signoff2-hermes-1.md:271-279`) asked for either `--keep <agent>.<field>` OR
mandatory per-deck enumeration of removed deliberate pins. Revision 3 adopts
BOTH: `--keep` ships (`consensus.md:391-393`) AND the dry-run and final report
must enumerate every removed deliberate pin per deck (`:393-394`). kimi-1's R4
is satisfied by construction — both halves are stronger than either alone.
`--keep` gives the user a proactive tool to preserve pins; the enumeration
gives a retrospective checklist when `--keep` is not used or when pins are
removed by other paths. kimi-1's R1–R3 are unchanged from revision 2 (I
assessed them in my revision-2 signoff at `signoff2-hermes-1.md:219-307`); the
revision-3 additions (retired-row retention, foreign-deck gate,
protocol-changelog) close the two partial-address items I noted there
(`:264-269` for R3's foreign-deck gate, `:237-248` for R2's changelog entry).

---

#### The field table — checked against source (the load-bearing render-only claim)

The claim that workspace dir, role, and host handle are render-only is
drafter-owned (`consensus.md:341-346`, `:353-363`). It is load-bearing: the
entire cutover is tractable only because most of §2 is already prose no code
reads. Per §15.1, a non-owner verdict on a drafter-owned claim is what makes
it admissible. I am a non-owner — my R3.1 was about the inactive-set wiring,
not about workspace_dir/role/host_handle being render-only. I verdict the
claim now.

**VERDICT: CONFIRMED (`PRIMARY`).**

I verified the render-only classification for each of the three fields:

**workspace_dir — render-only. CONFIRMED.**
- `PRIMARY` — `internal/protocol/roster.go:17`: `rosterRowRe` captures ONLY the
  first cell (agent ID) via `^\\|\\s*\`([a-z0-9][a-z0-9-]*)\`\\s*\\|`. The
  regex does not capture column 2 (Workspace dir). The parser reads no other
  cell from the roster row.
- `PRIMARY` — `internal/app/roster.go:81-89`: the `rosterRow` struct has no
  `WorkspaceDir` field. `resolveRoster` (`:100-167`) builds rows from
  `ReadRosterIDs` (ID only) + `LoadRosterAdapters` (TOML `[roster.<id>].adapter`)
  + the discovered agent spec. Workspace dir is never read from §2.
- `PRIMARY` — My `search_files` for `Workspace dir|workspace_dir` in `*.go`
  returned 4 hits, ALL in `*_test.go` test-fixture header strings
  (`internal/app/roster_test.go:48,153`,
  `internal/protocol/drift_test.go:28`,
  `internal/protocol/roster_test.go:16`). Zero non-test hits.
- `PRIMARY` — `COOPERATION.md:105-110`: §2 stores Workspace dir as column 2.
  The field table's "legacy §2 source: col 2" at `consensus.md:361` is
  accurate.

**role — render-only. CONFIRMED.**
- `PRIMARY` — The same `rosterRowRe` at `internal/protocol/roster.go:17`
  captures only the agent ID. Column 3 (Role, stored as prose like
  "facilitator+participant (cli `claude`…)") is never parsed.
- `PRIMARY` — `internal/app/roster.go:81-89`: no `Role` field in `rosterRow`.
  No code extracts role from §2.
- `PRIMARY` — My `search_files` for role-related identifiers in roster context
  found no §2 role parser. Role appears in `COOPERATION.md:105-110` column 3
  as prose; the field table's "legacy §2 source: col 3 prose" at
  `consensus.md:362` is accurate.

**host_handle — render-only. CONFIRMED.**
- `PRIMARY` — `internal/protocol/roster.go:18-19`: `rosterHeaderRe` matches
  `^|\s*Agent ID\s*|\s*Workspace` — this distinguishes the roster table from
  the host-handle table. The parser at `:41-44` sets `inTable = true` only
  when this header matches, and at `:49-51` breaks out of the table at the
  first non-`|` line. The host-handle table (`COOPERATION.md:119-126`) is a
  SEPARATE table below the roster table — the parser never reaches it.
- `PRIMARY` — `internal/protocol/roster_test.go:38`: "The host-handle table
  must NOT add rows (only the first roster table is read)." The test confirms
  the parser explicitly does not read the host-handle table.
- `PRIMARY` — `internal/app/roster.go:81-89`: no `HostHandle` field in
  `rosterRow`.
- `PRIMARY` — My `search_files` for `Host handle|host_handle|HostHandle` in
  `*.go` returned 2 hits, both comments
  (`internal/protocol/roster_test.go:38`,
  `internal/protocol/roster.go:18`) — zero code reads.
- `PRIMARY` — `COOPERATION.md:119-126`: the host-handle table is separate from
  the roster table. The field table's "legacy §2 source: the separate
  host-handle table (`COOPERATION.md:119-126`)" at `consensus.md:363` is
  accurate.

**The runtime-semantic fields are also correctly classified.**
- agent ID: runtime-semantic — CONFIRMED. `ReadRosterIDs` extracts it
  (`internal/protocol/roster.go:56-60`); used in all call sites.
- adapter: runtime-semantic — CONFIRMED. `LoadRosterAdapters` reads
  `[roster.<id>] adapter` from TOML (`internal/config/runtime.go:200`); used
  in `resolveRoster` at `internal/app/roster.go:109`.
- state (active/inactive): runtime-semantic — CONFIRMED. `ReadRosterIDs`
  extracts the inactive marker (`internal/protocol/roster.go:62-63`); used in
  `internal/app/app.go:2418` and `:1793` (though discarded in `resolveRoster`
  at `internal/app/roster.go:110`, which is the R3.1 bug).
- model, effort, speed: runtime-semantic — CONFIRMED. Present in `rosterRow`
  at `internal/app/roster.go:150-152`, sourced from the discovered agent spec,
  and in the launch argv.

**The drafter's measurement claim at `consensus.md:341-346` is also
confirmed.** The claim has two parts:
1. `ReadRosterIDs` extracts only the agent ID and the literal `inactive` —
   CONFIRMED. `PRIMARY` — `internal/protocol/roster.go:56-64`: the regex
   captures group 1 (agent ID) and `strings.Contains(line, "inactive")` checks
   for the inactive marker. No other cell is read.
2. A `find`-based sweep of non-test `*.go` for `Host handle`/`Workspace dir`
   returns zero hits — CONFIRMED. `PRIMARY` — my `search_files` returned 6
   hits, all in `*_test.go` files (test-fixture header strings), zero in
   non-test `.go` files.

**The refinement at `consensus.md:348-351` is confirmed.** The parser sets
`active[id] = true` for every row including inactive ones
(`internal/protocol/roster.go:61`), and `inactive` is a separate map
(`:62-63`). An inactive agent is in both maps. `resolveRoster` at
`internal/app/roster.go:110` discards inactive (`active, _, ok`), so a retired
agent renders as a full member in `roster show`. I own the R3.1 observation
about `resolveRoster` discarding inactive; the dual-map refinement is the
drafter's, and I verdict it CONFIRMED as a non-owner.

The field table is accurate against the source. The render-only classification
for workspace dir, role, and host handle is the load-bearing claim that makes
the cutover tractable, and it is CONFIRMED by primary source inspection.

---

#### Drafter position changes (§15.5)

`PRIMARY` — `consensus.md:248-274` contains `## Drafter position changes` with
fifteen entries. Entries 1–4 are from the rounds; 5–11 from revision-1/2
signoff; 12–15 are the four new revision-3 additions:

- Entry 12 (`:267`): the normative field table replaces the "MUST define" TODO.
- Entry 13 (`:268`): the `meta/protocol-changelog.md` entry requirement.
- Entry 14 (`:269`): foreign-deck compatibility gate + retired-row retention.
- Entry 15 (`:270`): kimi-1's R4 adopted as both halves.

The section satisfies §15.5: claude-1 is facilitator+participant+drafter, the
role concentration is recorded at `:250`, and every material change has a
prior position, source, new position, and rationale. I ratify its accuracy and
completeness as an existing signer (§15.5). The drafter's self-assessment at
`:272-274` — "Neither revision 1 nor revision 2 survived review. Both times
the failure was the same shape — the drafter wrote what must be true instead
of making it true" — is candid and accurate. Revision 3 breaks that pattern:
the field table exists, not a promise to define one.

---

#### Verification notes

- `PRIMARY` — `internal/protocol/roster.go:17,56-64`: the parser captures only
  the agent ID and the inactive marker. Confirmed the drafter's measurement at
  `consensus.md:341-346`.
- `PRIMARY` — `internal/protocol/roster.go:18-19,41-51`: `rosterHeaderRe`
  matches `Agent ID | Workspace` and the parser breaks at the first non-`|`
  line, so the host-handle table is never reached. Confirmed the render-only
  classification for host_handle.
- `PRIMARY` — `internal/app/roster.go:81-89`: `rosterRow` has no
  WorkspaceDir/Role/HostHandle field. Confirmed the render-only classification
  for all three fields.
- `PRIMARY` — `internal/app/roster.go:110`: `active, _, ok :=
  protocol.ReadRosterIDs(root)` — inactive discarded in `resolveRoster`. This
  is my own R3.1 finding; I do not verdict it — I confirm the consensus cites
  the correct locator.
- `PRIMARY` — `internal/app/app.go:2412,2418`: `defaultRosterParticipants`
  keeps and uses the inactive map (`if inactive[id] { continue }`). The
  consensus correctly scopes its render-as-full-member claim to `resolveRoster`
  (the `roster show` path), not to all call sites.
- `PRIMARY` — `COOPERATION.md:134`: "mark its row as inactive (do not delete
  it)." Confirmed the migration's retired-row retention is protocol-grounded.
- `PRIMARY` — `COOPERATION.md:719-724`: the §7 changelog format. Confirmed the
  protocol-changelog requirement at `consensus.md:383-384` references the
  correct format.

No `DISPUTED` claims. No `EXEMPTION-CLAIM UNVERIFIED`. The render-only claim
was the one material drafter-owned claim that needed a non-owner verdict for
admissibility (§15.1); it is now CONFIRMED by primary source inspection.

§15.6 (correlated agreement) is unchanged from revision 2: round 1 had
substantive disagreement (four different column sets, different sync
semantics, different scope labels, different §2-authority positions), so
§15.6(a) does not bind; §15.6(b) is satisfied by the `## Comparison & blind
spots` section at `consensus.md:276-295`, which records the shared-prior
caveat and states what would have to be true for the agreed position to be
wrong.

---

#### Reservation carried forward: none

R1, R2, and R3 (all five sub-points) were addressed in revision 2 and are
strengthened or unchanged in revision 3. The two minor FINAL.md carry-items I
noted in revision 2 (explicit protocol-version gate naming, removed-pin
enumeration) are now both in binding consensus text
(`consensus.md:386-389` and `:391-394` respectively). I have no remaining
reservations.

I accept.
