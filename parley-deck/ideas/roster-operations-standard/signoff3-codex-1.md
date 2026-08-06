### codex-1 — revision 3

**Verdict:** accept

#### Scope

- `PRIMARY` — I read the live `parley-deck/COOPERATION.md:1-1316` in full, including the Phase-3 gate at `:351-368`, §7 at `:717-730`, and the binding verification rules at `:1176-1316`.
- `PRIMARY` — I read the complete revision-3 `parley-deck/ideas/roster-operations-standard/consensus.md:1-992` and my complete revision-2 signoff at `parley-deck/ideas/roster-operations-standard/signoff2-codex-1.md:1-65`.
- `PRIMARY` — For the new field contract I read `internal/protocol/roster.go:11-69`, `internal/app/roster.go:21-166`, `internal/app/app.go:1780-1818,2405-2429`, `internal/app/preset.go:14-90`, `internal/config/roster.go:13-123`, `internal/config/runtime.go:22-35,134-153,196-245,520-617`, `internal/agents/resolve.go:28-67`, `internal/agents/discover.go:20-84`, `internal/agents/launchargs.go:5-121`, and `internal/runner/runner.go:841-871,1094-1124`. I also executed the read-only non-test Go searches quoted below.
- `PRIMARY` — I did not inspect the 40 foreign decks, run a migration dry-run, launch an agent, run tests, or execute any roster or Git write command. I wrote only this signoff file.
- `PRIMARY` — I issue no verification verdict on the run-manifest, continuation, or rebase-safety claims I already own. This signoff checks revision 3's new drafter-owned field claim and whether the text satisfies my prior counter-proposal; it does not self-verify my earlier claims, as prohibited by `COOPERATION.md:1197-1205`.

#### Revision-2 counter-proposal

1. **Normative field contract, retention, ordering, and migration — met.** `PRIMARY` — The table at `consensus.md:353-363` gives all nine requested fields an exact committed TOML key, legacy §2 source, absence/conflict rule, and runtime-semantic/render-only classification. The surrounding normative text makes `parley-deck/agents.toml` authoritative and generated §2 non-authoritative (`:374-382`), retains retired rows as `active = false` (`:386-389`), orders active before inactive and then agent ID byte-ascending (`:365-367`), and defines verbatim migration plus skip-on-unclean behavior (`:369-372`). This is the replacement contract revision 2 lacked, not another promise to define it later.
2. **§7-format changelog entry — met.** `PRIMARY` — `consensus.md:383-384` requires `meta/protocol-changelog.md` in the format specified by `COOPERATION.md:719-724`, naming this idea and the user-authorized one-off.
3. **Foreign-deck compatibility and retired history — met.** `PRIMARY` — `consensus.md:386-389` requires older-protocol/schema decks to be skipped and reported, never silently upgraded, and requires retired rows to remain as `active = false`, never removed.
4. **Kimi-1 R4 — met, with both halves.** `PRIMARY` — `consensus.md:391-394` requires `roster sync --keep <agent>.<field>` and also requires both dry-run and final report to enumerate every removed deliberate pin per deck. This is stronger than the either/or minimum in `signoff2-codex-1.md:63`.

`PRIMARY` — All four items are met. No part of my revision-2 counter-proposal remains outstanding.

#### Non-owner source verdict on the field table

**Verdict on the drafter-owned load-bearing claim: CONFIRMED (`PRIMARY`).** The scoped claim is that the current §2 roster's workspace-dir, role, and host-handle values are render-only, while membership ID and the inactive marker are the §2 values with runtime meaning.

- `PRIMARY` — `internal/protocol/roster.go:16-19` defines the roster-row capture solely around the first backtick-wrapped cell. In the parser body, the relevant passage is `id := m[1]`, `active[id] = true`, followed by `if strings.Contains(strings.ToLower(line), "inactive") { inactive[id] = true }` at `:56-64`. No workspace, role, adapter, model, or host-handle cell is extracted.
- `PRIMARY` — The same passage confirms the refinement at `consensus.md:348-351`: every parsed row enters `active`, and an inactive row separately enters `inactive`. `internal/app/roster.go:109-110` then reads `active, _, ok := protocol.ReadRosterIDs(root)`, discarding that inactive map for `roster show`; other runtime paths do consume the separate inactive set to reject or omit inactive members at `internal/config/roster.go:104-115` and `internal/app/app.go:2412-2425`.
- `PRIMARY` — I executed `find . -type f -name '*.go' ! -name '*_test.go' -exec grep -nHE 'Host handle|host_handle|Workspace dir|workspace_dir' {} +`; the relevant output was **no matches**. This independently reproduces the drafter's zero-hit sweep at `consensus.md:343-346`.
- `PRIMARY` — The exact adapter mapping is already runtime data: `internal/config/runtime.go:29-35` defines `[roster.<id>].adapter`, `:196-220` loads it, and participant resolution fails closed without an exact ID or explicit mapping at `internal/agents/resolve.go:28-67`. Model, reasoning/effort, and speed are likewise effective runtime fields at `internal/config/runtime.go:594-612`; model and effort are materialized into launch arguments by `internal/agents/launchargs.go:48-80` and `internal/runner/runner.go:1094-1112`, while speed is supplied to the participant prompt at `internal/runner/runner.go:841-871`. Their **runtime-semantic** classification is therefore consistent with the source.
- `PRIMARY` — Current §2 puts workspace dir and role in its human-readable table and host handles in the adjacent table at `COOPERATION.md:101-126`, but the parser evidence and zero-hit scan above show that none of those values governs membership resolution or launch. The table's **render-only** classification for `workspace_dir`, `role`, and `host_handle` is therefore admissible and correct for this cutover.

`PRIMARY` — The future `[roster.<id>].active`, `.model`, `.effort`, `.speed`, `.workspace_dir`, `.role`, and `.host_handle` names are normative additions rather than claims that current code already implements them. Current source exposes only `.adapter` in the roster struct (`internal/config/runtime.go:22-35`), so the authority cutover correctly remains an implementation obligation guarded by `consensus.md:374-382`; that does not contradict the field contract.

#### Prior accepted positions

- **hermes-1 — not weakened.** `PRIMARY` — Revision 3 preserves the same-release snapshot/rebase gate and continuation test at `consensus.md:316-327`, byte-idempotent §2 generation at `:365-367,404-406`, `STATE` wiring before migration at `:396-402`, and all five R3 operational guardrails at `:447-468`. The additions strengthen rather than dilute hermes-1's accepted-with-reservations position.
- **kimi-1 — not weakened; R4 is satisfied.** `PRIMARY` — Kimi-1's R1-R3 protections remain in the same release, authority, compatibility, retention, and migration clauses at `consensus.md:316-327,374-389,447-468`. R4 is satisfied by adopting **both** halves at `:391-394`: `--keep <agent>.<field>` ships, and every removed deliberate pin is enumerated per deck in both dry-run and final reporting.

#### Delivery staging

`PRIMARY` — My four-stage plan at `signoff2-codex-1.md:45-54` is recorded faithfully in substance and dependency order: runtime foundations and `STATE` correctness precede exposure (`consensus.md:46-132,396-406`); the TOML-authority cutover, generated view, runtime-reader removal, protocol copies, skill, and changelog form one same-release group (`:338-384`); snapshot persistence and consumption gate rebase (`:316-327`); and migration runs only after the CLI contract exists, with a frozen attended dry-run and the inventory/CAS/restore/resume/report machinery (`:433-468`). Stages may be reviewed separately only while unfinished user-visible behavior remains disabled; the authority cutover and snapshot/rebase boundaries remain atomic release gates.

Revision 3 discharges my block. I accept the protocol text and the staged delivery contract as written.
